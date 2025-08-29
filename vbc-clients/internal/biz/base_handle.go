package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"sync"
	"time"
	"vbc/lib"
	//. "vbc/lib/builder"
)

const HandleStatus_waiting = 0
const HandleStatus_done = 1

const HandleResult_ok = 0
const HandleResult_failure = 1

type BaseHandle[T any] struct {
	DB        *gorm.DB
	Log       *log.Helper
	TableName string
	Handle    func(ctx context.Context, t *T) error
}

func (c *BaseHandle[T]) RunHandleJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("BaseHandle RunHandleJob Done")
				return
			default:
				sqlRows, err := c.DB.Table(c.TableName).
					Where("handle_status=? and deleted_at=0",
						HandleStatus_waiting).Rows()
				if err != nil {
					c.Log.Error(err)
				} else {
					if sqlRows != nil {
						// 需要保证顺序，所以只能一个一个执行，解决files的问题
						gLimit := lib.NewGLimit(1)
						var waitGroup sync.WaitGroup
						for sqlRows.Next() {
							var task T
							err = c.DB.ScanRows(sqlRows, &task)
							if err != nil {
								c.Log.Error(err)
								continue
							}
							waitGroup.Add(1)
							gLimit.Run(func() {
								err := c.Handle(ctx, &task)
								if err != nil {
									c.Log.Error(err)
								}
								waitGroup.Done()
							})
							time.Sleep(1)
						}
						err = sqlRows.Close()
						waitGroup.Wait()
						if err != nil {
							c.Log.Error(err)
						}
					}
				}
				time.Sleep(2 * time.Second)
			}
		}
	}()
	return nil
}
