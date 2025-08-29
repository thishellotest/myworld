package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
)

type ZohobuzTaskUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	UsageStatsUsecase *UsageStatsUsecase
	ZohoUsecase       *ZohoUsecase
	QueueUsecase      *QueueUsecase
}

func NewZohobuzTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	UsageStatsUsecase *UsageStatsUsecase,
	ZohoUsecase *ZohoUsecase,
	QueueUsecase *QueueUsecase) *ZohobuzTaskUsecase {
	uc := &ZohobuzTaskUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		UsageStatsUsecase: UsageStatsUsecase,
		ZohoUsecase:       ZohoUsecase,
		QueueUsecase:      QueueUsecase,
	}
	return uc
}

type ClientTaskEntity struct {
	ID           int32 `gorm:"primaryKey"`
	Gid          string
	BizDeletedAt int64
	UpdatedAt    int64
}

func (c *ZohobuzTaskUsecase) RunJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ZohobuzTaskUsecase:RunJob:Done")
				return
			default:
				c.log.Info("ZohobuzTaskUsecase_SyncTasksDeletes: ", time.Now().Format(time.RFC3339))
				err := c.SyncTasksDeletes()
				if err != nil {
					c.log.Error(err)
				}
				time.Sleep(5 * 60 * time.Second)
			}
		}
	}()
	return nil
}

// SyncTasksDeletes 同步状态
func (c *ZohobuzTaskUsecase) SyncTasksDeletes() error {

	sqlRows, err := c.CommonUsecase.DB().Table("client_tasks").Select("id,gid,biz_deleted_at,updated_at").
		Where("biz_deleted_at=0 and deleted_at=0 and status!='Completed'").Rows()
	if err != nil {
		return err
	} else {
		if sqlRows != nil {
			maxLimit := 200
			var records []*ClientTaskEntity
			for sqlRows.Next() {
				var task ClientTaskEntity
				err = c.CommonUsecase.DB().ScanRows(sqlRows, &task)
				if err != nil {
					c.log.Error(err)
					continue
				}
				records = append(records, &task)

				if len(records) >= maxLimit {
					err = c.HandleZohoTasksDelete(records)
					if err != nil {
						c.log.Error(err)
					}
					records = records[:0]
				}
			}
			if len(records) > 0 {
				err = c.HandleZohoTasksDelete(records)
				if err != nil {
					c.log.Error(err)
				}
			}
			err = sqlRows.Close()
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ZohobuzTaskUsecase) HandleZohoTasksDelete(clientTasks []*ClientTaskEntity) error {

	c.log.Debug(len(clientTasks))
	if len(clientTasks) == 0 {
		return nil
	}

	fields := config_zoho.DealLayout().DealApiNames()
	params := make(url.Values)
	for _, v := range clientTasks {
		params.Add("ids", v.Gid)
	}

	//c.log.Debug("HandleZohoTasksDelete params:", params)
	c.UsageStatsUsecase.Stat("ZohobuzTaskUsecase_HandleZohoTasksDelete", time.Now(), 1)
	records, err := c.ZohoUsecase.GetRecords(config_zoho.Tasks, fields, params)
	if err != nil {
		return err
	}
	listMaps := records.GetTypeList("data")
	if len(listMaps) == 0 {
		return nil
	}
	maps := make(map[string]bool)
	for _, v := range listMaps {
		maps[v.GetString("id")] = true
	}
	var needsDeleteGid []string
	for _, v := range clientTasks {
		if _, ok := maps[v.Gid]; !ok {
			needsDeleteGid = append(needsDeleteGid, v.Gid)
		}
	}
	if len(needsDeleteGid) > 0 {

		err = c.CommonUsecase.DB().Table("client_tasks").
			Where("gid in ?", needsDeleteGid).
			Updates(map[string]interface{}{
				"biz_deleted_at": time.Now().Unix(),
				"updated_at":     time.Now().Unix(),
			}).Error

		er := c.QueueUsecase.TaskDueDateChange(needsDeleteGid)
		if er != nil {
			c.log.Error(er, " needsDeleteGids:", needsDeleteGid)
		}

		return err
	}

	return nil
}
