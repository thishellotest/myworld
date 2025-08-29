package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
	"vbc/internal/config_vbc"
	"vbc/lib"
	"vbc/lib/builder"
)

type SyncAsanaTaskUsecase struct {
	CommonUsecase    *CommonUsecase
	log              *log.Helper
	AsanaUsecase     *AsanaUsecase
	DataEntryUsecase *DataEntryUsecase
	TUsecase         *TUsecase
	LogUsecase       *LogUsecase
}

func NewSyncAsanaTaskUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	AsanaUsecase *AsanaUsecase,
	DataEntryUsecase *DataEntryUsecase,
	TUsecase *TUsecase,
	LogUsecase *LogUsecase) *SyncAsanaTaskUsecase {
	return &SyncAsanaTaskUsecase{
		CommonUsecase:    CommonUsecase,
		log:              log.NewHelper(logger),
		AsanaUsecase:     AsanaUsecase,
		DataEntryUsecase: DataEntryUsecase,
		TUsecase:         TUsecase,
		LogUsecase:       LogUsecase,
	}
}

func (c *SyncAsanaTaskUsecase) LPushSyncTaskQueue(ctx context.Context, tasks ...string) error {
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_sync_asana_tasks_queue, tasks).Err()
}

func (c *SyncAsanaTaskUsecase) LPushSyncTaskProcessing(ctx context.Context, tasks ...string) error {
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_sync_asana_tasks_processing, tasks).Err()
}

func (c *SyncAsanaTaskUsecase) FinishSyncTask(ctx context.Context) error {
	return c.CommonUsecase.RedisClient().Del(ctx, Redis_sync_asana_tasks_processing).Err()
}

// InitSyncAsanaTask 把未处理完毕的任务重新推回队列
func (c *SyncAsanaTaskUsecase) InitSyncAsanaTask(ctx context.Context) error {
	processing := c.CommonUsecase.RedisClient().LRange(ctx, Redis_sync_asana_tasks_processing, 0, -1)
	if processing.Err() != nil && processing.Err() != redis.Nil {
		return processing.Err()
	}
	processingSlices, _ := processing.Result()
	if len(processingSlices) > 0 {
		err := c.LPushSyncTaskQueue(ctx, processingSlices...)
		if err != nil {
			return err
		}
		c.CommonUsecase.RedisClient().Del(ctx, Redis_sync_asana_tasks_processing)
	}
	return nil
}

// GetTasks 获取任务,且把任务加入处理队列
func (c *SyncAsanaTaskUsecase) GetTasks(ctx context.Context) (tasks []string) {
	timer := time.NewTimer(10 * time.Second)
	maxCount := 100
	tasksMap := make(map[string]bool)
	for {
		select {
		case <-timer.C:
			return
		default:
			for {
				r := c.CommonUsecase.RedisClient().RPop(ctx, Redis_sync_asana_tasks_queue)
				if r.Err() != nil {
					time.Sleep(1 * time.Second)
					if r.Err() == redis.Nil {
						break
					}
					if r.Err() != redis.Nil {
						c.log.Error(r.Err())
					}
				} else {
					if _, ok := tasksMap[r.Val()]; !ok {
						tasksMap[r.Val()] = true

						lpushErr := c.LPushSyncTaskProcessing(ctx, r.Val())
						if lpushErr != nil {
							c.log.Error(lpushErr)
						}
						tasks = append(tasks, r.Val())
						if len(tasks) >= maxCount {
							return
						}
					}
				}
			}
		}
	}
}

func (c *SyncAsanaTaskUsecase) RunSyncTaskJob(ctx context.Context) error {
	err := c.InitSyncAsanaTask(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("SyncAsanaTaskUsecase:RunSyncTaskJob:Done")
				return
			default:
				tasks := c.GetTasks(ctx)
				if len(tasks) > 0 {
					for _, v := range tasks {
						err := c.SyncTask(v)
						if err != nil {
							c.log.Error(err)
							er := c.LogUsecase.SaveLog(0, Log_FromType_Asana_SyncTaskInfo, map[string]interface{}{
								"asana_task_gid": v,
							})
							if er != nil {
								c.log.Error(er)
							}
						}
					}
					err := c.FinishSyncTask(ctx)
					if err != nil {
						c.log.Error(err)
					}
				}
			}
		}
	}()

	return nil
}

// SyncTask 同步单个task
func (c *SyncAsanaTaskUsecase) SyncTask(taskGid string) error {
	asanaGetATaskVo, isDel, err := c.AsanaUsecase.GetATask(taskGid)
	if err != nil {
		return err
	} else {
		dataEntry := make(TypeDataEntry)
		if isDel {
			dataEntry[FileName_asana_task_gid] = taskGid
			dataEntry[FieldName_biz_deleted_at] = time.Now().Unix()
			tClient, err := c.TUsecase.Data(Kind_client_cases, builder.Eq{FileName_asana_task_gid: taskGid})
			if err != nil {
				return err
			}
			if tClient == nil {
				return nil
			}
			if tClient.CustomFields.NumberValueByNameBasic(FieldName_biz_deleted_at) > 0 {
				return nil
			}
		} else {
			if asanaGetATaskVo == nil {
				return errors.New("asanaGetATaskVo不存在：" + taskGid)
			}
			dataEntry = asanaGetATaskVo.ToDataEntry()
			dataEntry[FieldName_biz_deleted_at] = 0
		}
		// Sync users
		if dataEntry["assignee_gid"] != "" {
			err := c.UserLPushSyncTaskQueue(context.TODO(), lib.InterfaceToString(dataEntry["assignee_gid"]))
			if err != nil {
				c.log.Error(err)
			}
		}

		var dataEntryList TypeDataEntryList
		dataEntryList = append(dataEntryList, dataEntry)
		_, err = c.DataEntryUsecase.Handle(Kind_client_cases, dataEntryList, FileName_asana_task_gid, nil)
		if err != nil {
			return err
		} else {
			// if Source is empty, Set to manual.
			if _, ok := dataEntry["source"]; ok {
				sourceVal := InterfaceToString(dataEntry["source"])
				if sourceVal == "" {
					if _, ok := dataEntry[FileName_asana_task_gid]; ok {
						asanaTaskGid := InterfaceToString(dataEntry[FileName_asana_task_gid])
						field := config_vbc.GetAsanaCustomFields()
						gid := field.GetByName(config_vbc.Asana_Field_Source).GetGid()
						eGid := field.GetByName(config_vbc.Asana_Field_Source).GetEnumGidByName(config_vbc.Source_Manual)
						_, err := c.AsanaUsecase.PutATask(asanaTaskGid, lib.TypeMap{
							gid: eGid,
						}, "")
						if err != nil {
							c.log.Error(err)
						}
					}
				}
			}
		}
	}
	return nil
}

// SyncUser 同步单个user
func (c *SyncAsanaTaskUsecase) SyncUser(userGid string) error {
	response, err := c.AsanaUsecase.GetAUser(userGid)
	if err != nil {
		return err
	} else if response == nil {
		return errors.New("response不存在：" + userGid)
	} else {
		sourceData := make(map[string]interface{})
		sourceData["asana_user_gid"] = response.Get("data.gid")
		sourceData["email"] = response.Get("data.email")
		sourceData["name"] = response.Get("data.name")
		var dataEntryList TypeDataEntryList
		dataEntryList = append(dataEntryList, sourceData)

		_, err = c.DataEntryUsecase.Handle(Kind_users, dataEntryList, FileName_asana_user_gid, nil)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *SyncAsanaTaskUsecase) UserLPushSyncTaskQueue(ctx context.Context, tasks ...string) error {
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_sync_asana_users_queue, tasks).Err()
}

func (c *SyncAsanaTaskUsecase) UserLPushSyncTaskProcessing(ctx context.Context, tasks ...string) error {
	return c.CommonUsecase.RedisClient().LPush(ctx, Redis_sync_asana_users_processing, tasks).Err()
}

func (c *SyncAsanaTaskUsecase) FinishSyncUser(ctx context.Context) error {
	return c.CommonUsecase.RedisClient().Del(ctx, Redis_sync_asana_users_processing).Err()
}

// UserInitSyncAsanaTask 把未处理完毕的任务重新推回队列
func (c *SyncAsanaTaskUsecase) UserInitSyncAsanaTask(ctx context.Context) error {
	processing := c.CommonUsecase.RedisClient().LRange(ctx, Redis_sync_asana_users_processing, 0, -1)
	if processing.Err() != nil && processing.Err() != redis.Nil {
		return processing.Err()
	}
	processingSlices, _ := processing.Result()
	if len(processingSlices) > 0 {
		err := c.LPushSyncTaskQueue(ctx, processingSlices...)
		if err != nil {
			return err
		}
		c.CommonUsecase.RedisClient().Del(ctx, Redis_sync_asana_users_processing)
	}
	return nil
}

func (c *SyncAsanaTaskUsecase) RunSyncUserJob(ctx context.Context) error {
	err := c.UserInitSyncAsanaTask(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("SyncAsanaTaskUsecase:RunSyncUserJob:Done")
				return
			default:
				users := c.GetUsers(ctx)
				if len(users) > 0 {
					for _, v := range users {
						err := c.SyncUser(v)
						if err != nil {
							c.log.Error(err)
						}
					}
					err := c.FinishSyncUser(ctx)
					if err != nil {
						c.log.Error(err)
					}
				}
			}
		}
	}()

	return nil
}

// GetUsers 获取任务,且把任务加入处理队列
func (c *SyncAsanaTaskUsecase) GetUsers(ctx context.Context) (users []string) {
	timer := time.NewTimer(10 * time.Second)
	maxCount := 100
	tasksMap := make(map[string]bool)
	for {
		select {
		case <-timer.C:
			return
		default:
			for {
				r := c.CommonUsecase.RedisClient().RPop(ctx, Redis_sync_asana_users_queue)
				if r.Err() != nil {
					time.Sleep(1 * time.Second)
					if r.Err() == redis.Nil {
						break
					}
					if r.Err() != redis.Nil {
						c.log.Error(r.Err())
					}
				} else {
					if _, ok := tasksMap[r.Val()]; !ok {
						tasksMap[r.Val()] = true

						lpushErr := c.UserLPushSyncTaskProcessing(ctx, r.Val())
						if lpushErr != nil {
							c.log.Error(lpushErr)
						}
						users = append(users, r.Val())
						if len(users) >= maxCount {
							return
						}
					}
				}
			}
		}
	}
}
