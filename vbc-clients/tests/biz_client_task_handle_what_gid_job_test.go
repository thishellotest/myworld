package tests

import (
	"context"
	"sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientTaskHandleWhatGidJobUsecase_RunCustomTaskJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.ClientTaskHandleWhatGidJobUsecase.RunCustomTaskJob(context.TODO())
	wait.Wait()
}

func Test_ClientTaskHandleWhatGidJobUsecase_Do(t *testing.T) {
	UT.ClientTaskHandleWhatGidJobUsecase.Do([]string{"444edb70887a4d65bfa3f79cd2189603", "s1", "s2"})
}

func Test_ClientTaskHandleWhatGidJobUsecase_BizHandleTask(t *testing.T) {
	var customTaskParams []biz.CustomTaskParams
	customTaskParams = append(customTaskParams, biz.CustomTaskParams{
		UniqueKey: "6159272000009972111",
	})
	err := UT.ClientTaskHandleWhatGidJobUsecase.BizHandleTask(context.TODO(), customTaskParams)
	lib.DPrintln(err)
}
