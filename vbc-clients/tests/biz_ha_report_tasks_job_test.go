package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_HaReportTaskJobUsecase_RunJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	ctx := context.TODO()
	err := UT.HaReportTaskJobUsecase.RunHandleCustomJob(ctx, 1, time.Second*5,
		UT.HaReportTaskJobUsecase.WaitingTasks,
		UT.HaReportTaskJobUsecase.Handle,
	)
	lib.DPrintln(err)

	UT.HaReportPageJobUsecase.RunHandleCustomJob(ctx, 1, time.Second*5, UT.HaReportPageJobUsecase.WaitingTasks,
		UT.HaReportPageJobUsecase.Handle)

	wait.Wait()
}

func Test_HaReportTaskJobUsecase_WaitingTasks(t *testing.T) {
	rows, err := UT.HaReportTaskJobUsecase.WaitingTasks(context.TODO())
	lib.DPrintln(err)
	//_, list, err := lib.SqlRowsTrans(rows)
	//lib.DPrintln(list)
	//var tasks []*biz.HaReportTaskEntity
	tasks, err := lib.SqlRowsToEntities[biz.HaReportTaskEntity](UT.CommonUsecase.DB(), rows)
	lib.DPrintln(err, "__")
	lib.DPrintln(tasks)
}

//func Test_HaReportTaskJobUsecase_WaitingTasksByCreatingPdf(t *testing.T) {
//
//	rows, err := UT.HaReportTaskJobUsecase.WaitingTasksByCreatingPdf(context.TODO())
//	tasks, err := lib.SqlRowsToEntities[biz.HaReportTaskEntity](UT.CommonUsecase.DB(), rows)
//	lib.DPrintln(err, "__")
//	lib.DPrintln(tasks)
//}
//
//func Test_HaReportTaskJobUsecase_HandleByCreatingPdf(t *testing.T) {
//
//	a, _ := UT.HaReportTaskUsecase.GetByCond(builder.Eq{"id": 12})
//	err := UT.HaReportTaskJobUsecase.HandleByCreatingPdf(context.TODO(), a)
//	lib.DPrintln(err)
//}

func Test_HaReportTaskJobUsecase_Handle(t *testing.T) {

	a, _ := UT.HaReportTaskUsecase.GetByCond(builder.Eq{"id": 12})
	err := UT.HaReportTaskJobUsecase.Handle(context.TODO(), a)
	lib.DPrintln(err)
}

func Test_HaReportTaskJobUsecase_HandleExec(t *testing.T) {

	a, _ := UT.HaReportTaskUsecase.GetByCond(builder.Eq{"id": 12})
	err := UT.HaReportTaskJobUsecase.HandleExec(context.TODO(), a)
	lib.DPrintln(err)
}

func Test_HaReportTaskJobUsecase_HandleAiReport(t *testing.T) {
	task, _ := UT.HaReportTaskUsecase.GetByCond(builder.Eq{"id": 1})
	blobSlice, _ := UT.BlobSliceUsecase.GetByCond(builder.Eq{"blob_gid": task.BlobGid, "gid": "3d978012b99340f58d2713d34c174c79"})
	err := UT.HaReportTaskJobUsecase.HandleAiReport(context.TODO(), task, blobSlice)
	lib.DPrintln(err)
}
