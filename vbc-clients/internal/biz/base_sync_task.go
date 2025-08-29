package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type SyncTask struct {
	RedisQueue      string
	RedisProcessing string
	RedisClient     *redis.Client
	Log             *log.Helper
	Handle          func(ctx context.Context, str string) error
	MaxBatchLimit   int
}

func (c *SyncTask) LPushSyncTaskQueue(ctx context.Context, tasks ...string) error {
	return c.RedisClient.LPush(ctx, c.RedisQueue, tasks).Err()
}

func (c *SyncTask) LPushSyncTaskProcessing(ctx context.Context, tasks ...string) error {
	return c.RedisClient.LPush(ctx, c.RedisProcessing, tasks).Err()
}

func (c *SyncTask) FinishSyncTask(ctx context.Context) error {
	return c.RedisClient.Del(ctx, c.RedisProcessing).Err()
}

// InitSyncTask 把未处理完毕的任务重新推回队列
func (c *SyncTask) InitSyncTask(ctx context.Context) error {
	processing := c.RedisClient.LRange(ctx, c.RedisProcessing, 0, -1)
	if processing.Err() != nil && processing.Err() != redis.Nil {
		return processing.Err()
	}
	processingSlices, _ := processing.Result()
	if len(processingSlices) > 0 {
		err := c.LPushSyncTaskQueue(ctx, processingSlices...)
		if err != nil {
			return err
		}
		c.RedisClient.Del(ctx, c.RedisProcessing)
	}
	return nil
}

func (c *SyncTask) GetTasks(ctx context.Context) (tasks []string) {
	timer := time.NewTimer(10 * time.Second)
	if c.MaxBatchLimit <= 0 {
		c.MaxBatchLimit = 100
	}
	tasksMap := make(map[string]bool)
	for {
		select {
		case <-timer.C:
			return
		default:
			for {
				r := c.RedisClient.RPop(ctx, c.RedisQueue)
				if r.Err() != nil {
					time.Sleep(1 * time.Second)
					if r.Err() == redis.Nil {
						break
					}
					if r.Err() != redis.Nil {
						c.Log.Error(r.Err())
					}
				} else {
					if _, ok := tasksMap[r.Val()]; !ok {
						tasksMap[r.Val()] = true

						lpushErr := c.LPushSyncTaskProcessing(ctx, r.Val())
						if lpushErr != nil {
							c.Log.Error(lpushErr)
						}
						tasks = append(tasks, r.Val())
						if len(tasks) >= c.MaxBatchLimit {
							return
						}
					}
				}
			}
		}
	}
}

func (c *SyncTask) RunSyncTaskJob(ctx context.Context) error {
	err := c.InitSyncTask(ctx)
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
						err := c.Handle(ctx, v)
						if err != nil {
							c.Log.Error(err)
						}
					}
					err := c.FinishSyncTask(ctx)
					if err != nil {
						c.Log.Error(err)
					}
				}
			}
		}
	}()
	return nil
}
