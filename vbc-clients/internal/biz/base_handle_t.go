package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
	"time"
	"vbc/lib"
	//. "vbc/lib/builder"
)

func AppendHandleResultDetail(tData *TData, err error) string {
	if tData == nil {
		return ""
	}
	detail := tData.CustomFields.TextValueByNameBasic("handle_result_detail")
	if detail != "" {
		detail += "\n\n"
	}
	detail += err.Error() + " ||| " + time.Now().Format(time.RFC3339)
	return detail
}

type BaseHandleT[T any] struct {
	Log *log.Helper
}

func (c *BaseHandleT[T]) RunHandleCustomJob(ctx context.Context, gLimitNum int, windowTime time.Duration,
	source func(ctx context.Context) ([]T, error),
	handle func(ctx context.Context, t T) error) error {

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
				sources, err := source(ctx)
				if err != nil {
					c.Log.Error(err)
				} else {
					if sources != nil {
						gLimit := lib.NewGLimit(gLimitNum)
						var waitGroup sync.WaitGroup
						for k, _ := range sources {
							waitGroup.Add(1)
							gLimit.Run(func() {
								err := handle(ctx, sources[k])
								if err != nil {
									c.Log.Error(err)
								}
								waitGroup.Done()
							})
							time.Sleep(1)
						}
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
