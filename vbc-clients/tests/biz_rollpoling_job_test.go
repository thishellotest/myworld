package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RollpoingJobUsecase(t *testing.T) {
	var wait sync.WaitGroup
	err := UT.RollpoingJobUsecase.RunHandleCustomJob(context.TODO(), 3, 0,
		UT.RollpoingJobUsecase.WaitingTasks,
		UT.RollpoingJobUsecase.Handle)
	fmt.Println(err)
	wait.Add(1)
	wait.Wait()
}

func Test_RollpoingJobUsecase_Handle(*testing.T) {
	var task biz.RollpoingEntity
	UT.CommonUsecase.DB().Where("id=26").Take(&task)
	err := UT.RollpoingJobUsecase.Handle(context.TODO(), &task)
	fmt.Println(task.VendorUniqId)
	lib.DPrintln(err)
}

func Test_RollpoingJobUsecase_HandleExec(t *testing.T) {
	var task biz.RollpoingEntity
	UT.CommonUsecase.DB().Where("id=23").Take(&task)

	isDone, err := UT.RollpoingJobUsecase.HandleExec(context.TODO(), &task)
	fmt.Println(task.VendorUniqId)
	lib.DPrintln(isDone, err)
}
