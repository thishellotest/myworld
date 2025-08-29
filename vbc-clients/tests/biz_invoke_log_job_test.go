package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_InvokeLogJobUsecase_FormatZohoNoteContent(t *testing.T) {
	str := `crm[user#6159272000000453640#847424422]crm The Claims folder is ready for submission, link here:https://veteranbenefitscenter.app.box.com/folder/288850667064

Please note that we could not find any evidence for his urinary frequency condition.
@[Engineering Team](6159272000000453669)`
	r := UT.InvokeLogJobUsecase.FormatZohoNoteContent(str)
	lib.DPrintln(r)
}

func Test_InvokeLogJobUsecase_RunHandleCustomJob(t *testing.T) {
	var waitGroup sync.WaitGroup
	UT.InvokeLogJobUsecase.RunHandleCustomJob(context.TODO(), 1, 3*time.Second,
		UT.InvokeLogJobUsecase.WaitingTasks,
		UT.InvokeLogJobUsecase.Handle,
	)
	waitGroup.Add(1)
	waitGroup.Wait()
}

func Test_InvokeLogJobUsecase_HandleExec(t *testing.T) {
	entity, err := UT.InvokeLogUsecase.GetByCond(Eq{"id": 5320})
	lib.DPrintln(err)
	err = UT.InvokeLogJobUsecase.Handle(context.TODO(), entity)
	lib.DPrintln(err)
}
