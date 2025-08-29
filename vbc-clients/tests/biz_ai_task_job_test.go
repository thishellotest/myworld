package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_AiTaskJobUsecase_HandleExec(t *testing.T) {
	var entity biz.AiTaskEntity
	err := UT.CommonUsecase.DB().Where("id=1448").Take(&entity).Error
	//lib.DPrintln(err)
	//lib.DPrintln(entity)
	err = UT.AiTaskJobUsecase.HandleExec(context.TODO(), &entity)
	lib.DPrintln(err)
}

func Test_AiTaskJobUsecase_Handle(t *testing.T) {
	var entity biz.AiTaskEntity
	err := UT.CommonUsecase.DB().Where("id=1346").Take(&entity).Error
	lib.DPrintln(err)
	lib.DPrintln(entity)
	err = UT.AiTaskJobUsecase.Handle(context.TODO(), &entity)
	lib.DPrintln(err)
}
func Test_AiTaskJobUsecase_AfterAiTaskJobHandle(t *testing.T) {
	var entity biz.AiTaskEntity
	err := UT.CommonUsecase.DB().Where("id=1346").Take(&entity).Error
	lib.DPrintln(err)
	lib.DPrintln(entity)
	UT.AiTaskJobUsecase.AfterAiTaskJobHandle(context.TODO(), &entity)
}

func Test_AiTaskJobUsecase_Handle1(t *testing.T) {
	var entity biz.AiTaskEntity
	err := UT.CommonUsecase.DB().Where("id=6").Take(&entity).Error
	lib.DPrintln(err)
	er := UT.CommonUsecase.DB().Omit("deleted_at").Save(&entity).Error
	lib.DPrintln(er)
}
