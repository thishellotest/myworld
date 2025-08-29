package tests

import (
	"context"
	"sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientTaskHandleWhoGidJobUsecase_RunCustomTaskJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.ClientTaskHandleWhoGidJobUsecase.RunCustomTaskJob(context.TODO())
	wait.Wait()
}

func Test_ClientTaskHandleWhoGidJobUsecase_Do(t *testing.T) {
	UT.ClientTaskHandleWhoGidJobUsecase.Do([]string{"6159272000005519042", "s1", "s2"})
}

func Test_ClientTaskHandleWhoGidJobUsecase_BizHandleTask(t *testing.T) {
	var customTaskParams []biz.CustomTaskParams
	customTaskParams = append(customTaskParams, biz.CustomTaskParams{
		UniqueKey: "6159272000009972111",
	})
	err := UT.ClientTaskHandleWhoGidJobUsecase.BizHandleTask(context.TODO(), customTaskParams)
	lib.DPrintln(err)
}
