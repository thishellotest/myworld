package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type CustomTaskParams struct {
	UniqueKey string // 任务唯一值， 当队列有相同UniqueKey时， 只做一次任务
	Params    string // 任务业务参数
}

func (c CustomTaskParams) MarshalBinary() ([]byte, error) { // 实现redis  implement encoding.BinaryMarshaler
	return json.Marshal(c)
}

func (c *CustomTaskParams) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

type CustomTask struct {
	RedisQueue      string
	RedisProcessing string
	RedisClient     *redis.Client
	Log             *log.Helper
	Handle          func(ctx context.Context, str CustomTaskParams) error
	MaxBatchLimit   int
	WindowSeconds   time.Duration
}

func (c *CustomTask) LPushCustomTaskQueue(ctx context.Context, tasks ...CustomTaskParams) error {

	var dest []string
	for _, v := range tasks {
		bytes, _ := v.MarshalBinary()
		dest = append(dest, string(bytes))
	}
	return c.RedisClient.LPush(ctx, c.RedisQueue, dest).Err()
}

func (c *CustomTask) LPushCustomTaskProcessing(ctx context.Context, tasks ...CustomTaskParams) error {

	var dest []string
	for _, v := range tasks {
		bytes, _ := v.MarshalBinary()
		dest = append(dest, string(bytes))
	}
	return c.RedisClient.LPush(ctx, c.RedisProcessing, dest).Err()
}

func (c *CustomTask) FinishCustomTask(ctx context.Context) error {
	return c.RedisClient.Del(ctx, c.RedisProcessing).Err()
}

// InitCustomTask 把未处理完毕的任务重新推回队列
func (c *CustomTask) InitCustomTask(ctx context.Context) error {
	processing := c.RedisClient.LRange(ctx, c.RedisProcessing, 0, -1)
	if processing.Err() != nil && processing.Err() != redis.Nil {
		return processing.Err()
	}
	var processingSlices []CustomTaskParams
	processing.ScanSlice(&processingSlices)
	if len(processingSlices) > 0 {
		err := c.LPushCustomTaskQueue(ctx, processingSlices...)
		if err != nil {
			return err
		}
		c.RedisClient.Del(ctx, c.RedisProcessing)
	}
	return nil
}

func (c *CustomTask) GetTasks(ctx context.Context) (tasks []CustomTaskParams) {
	if c.WindowSeconds <= 0 {
		c.WindowSeconds = 10
	}

	timer := time.NewTimer(c.WindowSeconds * time.Second)
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
						c.Log.Info(r.Err())
					}
				} else {

					var customTaskParams CustomTaskParams
					err := r.Scan(&customTaskParams)
					if err != nil {
						c.Log.Error(err)
					}
					if customTaskParams.UniqueKey == "" {
						continue
					}

					if _, ok := tasksMap[customTaskParams.UniqueKey]; !ok {
						tasksMap[customTaskParams.UniqueKey] = true

						lpushErr := c.LPushCustomTaskProcessing(ctx, customTaskParams)
						if lpushErr != nil {
							c.Log.Error(lpushErr)
						}
						tasks = append(tasks, customTaskParams)
						if len(tasks) >= c.MaxBatchLimit {
							return
						}
					}
				}
			}
		}
	}
}

func (c *CustomTask) RunCustomTaskJob(ctx context.Context) error {
	err := c.InitCustomTask(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("CustomTask:RunCustomTaskJob:Done")
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
					err := c.FinishCustomTask(ctx)
					if err != nil {
						c.Log.Error(err)
					}
				}
			}
		}
	}()
	return nil
}
