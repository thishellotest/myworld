package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_BlobSliceJobUsecase_RunHandleCustomJob(t *testing.T) {

	var wait sync.WaitGroup
	wait.Add(1)
	UT.BlobSliceJobUsecase.RunHandleCustomJob(context.TODO(), 2, time.Second*10,
		UT.BlobSliceJobUsecase.WaitingTasks,
		UT.BlobSliceJobUsecase.Handle)
	wait.Wait()
}

func Test_BlobSliceJobUsecase_Handle(t *testing.T) {

	task, _ := UT.BlobSliceUsecase.GetByCond(builder.Eq{"id": 718})
	err := UT.BlobSliceJobUsecase.Handle(context.TODO(), task)
	lib.DPrintln(err)
}

func Test_BlobSliceJobUsecase_HandleExec(t *testing.T) {
	task, _ := UT.BlobSliceUsecase.GetByCond(builder.Eq{"id": 1234})
	err := UT.BlobSliceJobUsecase.HandleExec(context.TODO(), task)
	lib.DPrintln(err)
}

func Test_BlobSliceJobUsecase_HandleOperationLocation(t *testing.T) {
	task, _ := UT.BlobSliceUsecase.GetByCond(builder.Eq{"id": 382})
	ol := "https://documentintelligenceeu2s0.cognitiveservices.azure.com/formrecognizer/documentModels/prebuilt-read/analyzeResults/06785338-934c-4274-9685-20d307312424?api-version=2023-07-31"
	err := UT.BlobSliceJobUsecase.HandleOperationLocation(context.TODO(), task, ol)
	lib.DPrintln(err)
}
