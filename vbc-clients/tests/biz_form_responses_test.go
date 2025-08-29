package tests

import (
	"context"
	sync2 "sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_FormResponseUsecase_Handle(t *testing.T) {
	var entity biz.FormResponseEntity
	UT.CommonUsecase.DB().Take(&entity)
	err := UT.FormResponseUsecase.Handle(context.TODO(), &entity)
	lib.DPrintln(err)
}

func Test_FormResponseUsecase_HandleExec(t *testing.T) {
	var entity biz.FormResponseEntity
	UT.CommonUsecase.DB().Where("id=5026").Take(&entity)
	err := UT.FormResponseUsecase.HandleExec(context.TODO(), &entity)
	lib.DPrintln(err)
}

func Test_FormResponseUsecase_Run(t *testing.T) {
	var sync sync2.WaitGroup
	err := UT.FormResponseUsecase.RunHandleJob(context.TODO())
	lib.DPrintln(err)
	sync.Add(1)
	sync.Wait()
}
