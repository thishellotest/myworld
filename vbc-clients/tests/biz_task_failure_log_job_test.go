package tests

import (
	"context"
	"sync"
	"testing"
)

func Test_TaskFailureLogJobUsecase_(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.TaskFailureLogJobUsecase.RunHandleCustomJob(context.TODO(), 2, 0,
		UT.TaskFailureLogJobUsecase.WaitingTasks, UT.TaskFailureLogJobUsecase.Handle)
	wait.Wait()
}
