package tests

import (
	"context"
	"sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_WebhookLogJobUsecase(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.WebhookLogJobUsecase.RunHandleCustomJob(context.TODO(),
		2, 0, UT.WebhookLogJobUsecase.WaitingTasks,
		UT.WebhookLogJobUsecase.Handle)
	wait.Wait()
}

func Test_WebhookLogJobUsecase_HandleExecWebsite(t *testing.T) {
	var entity biz.WebhookLogEntity
	err := UT.CommonUsecase.DB().First(&entity, 8595).Error
	lib.DPrintln(err)
	a, err := UT.WebhookLogJobUsecase.HandleExecWebsite(context.TODO(), &entity)
	lib.DPrintln(a, err)
}

func Test_WebhookLogJobUsecase_HandleExecJotform(t *testing.T) {

	var entity biz.WebhookLogEntity
	err := UT.CommonUsecase.DB().First(&entity, 37139).Error
	lib.DPrintln(err)
	lib.DPrintln(entity)
	if err != nil {
		panic(err)
	}
	isDone, err := UT.WebhookLogJobUsecase.HandleExecJotform(context.TODO(), &entity)
	lib.DPrintln(isDone, err)
}

func Test_WebhookLogJobUsecase_HandleExecDialpad(t *testing.T) {

	var entity biz.WebhookLogEntity
	err := UT.CommonUsecase.DB().First(&entity, 5415).Error
	lib.DPrintln(err)
	lib.DPrintln(entity)
	if err != nil {
		panic(err)
	}
	isDone, err := UT.WebhookLogJobUsecase.HandleExecDialpad(context.TODO(), &entity)
	lib.DPrintln(isDone, err)
}

func Test_WebhookLogJobUsecase_HandleExec(t *testing.T) {

	var entity biz.WebhookLogEntity
	err := UT.CommonUsecase.DB().First(&entity, 5414).Error
	lib.DPrintln(err)
	lib.DPrintln(entity)
	if err != nil {
		panic(err)
	}
	isDone, err := UT.WebhookLogJobUsecase.HandleExec(context.TODO(), &entity)
	lib.DPrintln(isDone, err)
}
