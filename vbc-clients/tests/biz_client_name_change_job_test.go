package tests

import (
	"context"
	"sync"
	"testing"
	"vbc/lib"
)

func Test_ClientNameChangeJobUsecase_RunJob(t *testing.T) {

	var wait sync.WaitGroup
	wait.Add(1)
	err := UT.ClientNameChangeJobUsecase.RunCustomTaskJob(context.TODO())
	if err != nil {
		panic(err)
	}
	wait.Wait()
}

func Test_ClientNameChangeJobUsecase_Do(t *testing.T) {

	err := UT.ClientNameChangeJobUsecase.Do([]string{"d1fbcc1328424c3699057dd71f14e970"})
	lib.DPrintln(err)
}
