package tests

import (
	"testing"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientTaskBuzUsecase_HandleCompleteTask(t *testing.T) {
	err := UT.ClientTaskBuzUsecase.HandleCompleteTask(8)
	lib.DPrintln(err)
}

//func Test_ClientTaskBuzUsecase_HandleAutoCreateTask(t *testing.T) {
//	err := UT.ClientTaskBuzUsecase.HandleAutoCreateTask(context.TODO())
//	lib.DPrintln(err)
//}

func Test_ClientTaskBuzUsecase_NeedCreateTask(t *testing.T) {
	now := time.Now().In(configs.GetVBCDefaultLocation())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, configs.GetVBCDefaultLocation())
	ok, subDay, err := UT.ClientTaskBuzUsecase.NeedCreateTask("2024-10-10", now)
	lib.DPrintln(ok, subDay, err)
}

func Test_ClientTaskBuzUsecase_CreateClientTask(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5450)
	if err != nil {
		panic(err)
	}
	gid, r, err := UT.ClientTaskBuzUsecase.CreateClientTask(tCase, 3)
	lib.DPrintln(gid, r, err)
}

func Test_ClientTaskBuzUsecase_HandleOtherItfExpTask(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5264)
	if err != nil {
		panic(err)
	}
	err = UT.ClientTaskBuzUsecase.HandleOtherItfExpTask(tCase, "ITF Expiration within 2 days")
	lib.DPrintln(err)
}
