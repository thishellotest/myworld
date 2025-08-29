package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"sync"
	"time"
	"vbc/lib"
	//. "vbc/lib/builder"
)

type BaseHandleCustom[T any] struct {
	Log *log.Helper
	DB  *gorm.DB
}

func (c *BaseHandleCustom[T]) RunHandleCustomJob(ctx context.Context, gLimitNum int, windowTime time.Duration,
	source func(ctx context.Context) (*sql.Rows, error),
	handle func(ctx context.Context, t *T) error) error {

	if gLimitNum < 1 {
		panic("RunHandleCustomJob:gLimit must be greater than 0")
	}
	if windowTime < 0 {
		panic("RunHandleCustomJob:windowTime must be greater than or equal to 0")
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("BaseHandle RunHandleJob Done")
				return
			default:
				sqlRows, err := source(ctx)
				if err != nil {
					c.Log.Error(err)
				} else {
					if sqlRows != nil {
						gLimit := lib.NewGLimit(gLimitNum)
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
								err := handle(ctx, &task)
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
				if windowTime == 0 {
					time.Sleep(2 * time.Second)
				} else {
					time.Sleep(windowTime)
				}
			}
		}
	}()
	return nil
}
