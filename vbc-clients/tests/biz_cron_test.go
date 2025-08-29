package tests

import (
	"fmt"
	"sync"
	"testing"
	"vbc/lib"
)

func Test_CronUsecase_AddFunc(t *testing.T) {
	var wait sync.WaitGroup
	a, err := UT.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 9 * * 1", func() { // 每周一早上9点
		fmt.Println("ssss")
	})
	if err != nil {
		panic(err)
	}
	wait.Add(1)
	wait.Wait()
	lib.DPrintln(a)
}
