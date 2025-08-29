package tests

import (
	"context"
	sync2 "sync"
	"testing"
	"time"
	"vbc/lib"
)

func Test_BlobJobUsecase_WaitingTasks(t *testing.T) {
	abc, err := UT.BlobJobUsecase.WaitingTasks(context.TODO())
	lib.DPrintln(abc)
	lib.DPrintln(err)
}

func Test_BlobJobUsecase_RunHandleCustomJob(t *testing.T) {
	var sync sync2.WaitGroup
	sync.Add(1)
	ctx := context.TODO()
	UT.BlobJobUsecase.RunHandleCustomJob(
		ctx,
		2,
		5*time.Second,
		UT.BlobJobUsecase.WaitingTasks,
		UT.BlobJobUsecase.Handle)

	UT.BlobSliceJobUsecase.RunHandleCustomJob(ctx, 2, 5*time.Second,
		UT.BlobSliceJobUsecase.WaitingTasks,
		UT.BlobSliceJobUsecase.Handle)

	sync.Wait()
}

func Test_BlobJobUsecase_Handle(t *testing.T) {
	a, _ := UT.BlobUsecase.GetByGid("9aeae242abf14ed0952a0f223b0a6ce6")
	err := UT.BlobJobUsecase.Handle(context.TODO(), *a)
	lib.DPrintln(err)
}
