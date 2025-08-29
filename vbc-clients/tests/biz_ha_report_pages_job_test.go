package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_HaReportPageJobUsecase_Job(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.HaReportPageJobUsecase.RunHandleCustomJob(context.TODO(), 1, 10*time.Second, UT.HaReportPageJobUsecase.WaitingTasks,
		UT.HaReportPageJobUsecase.Handle)
	wait.Wait()
}

func Test_HaReportPageJobUsecase_Handle(t *testing.T) {

	a, _ := UT.HaReportPageUsecase.GetByCond(builder.Eq{"id": 359})
	err := UT.HaReportPageJobUsecase.Handle(context.TODO(), a)
	lib.DPrintln(err)
}
