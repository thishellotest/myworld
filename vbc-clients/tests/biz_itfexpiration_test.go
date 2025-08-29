package tests

import (
	"testing"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/lib"
)

//func Test_ItfexpirationUsecase_HandleCompleteItfTasks(t *testing.T) {
//	UT.ItfexpirationUsecase.HandleCompleteItfTasks(5217)
//}
//
//func Test_ItfexpirationUsecase_ExecuteCompleteItfTasks(t *testing.T) {
//	err := UT.ItfexpirationUsecase.ExecuteCompleteItfTasks()
//	lib.DPrintln(err)
//}

func Test_ItfexpirationUsecase_HandleITFExpireReminder(t *testing.T) {
	UT.ItfexpirationUsecase.HandleITFExpireReminder()
}

func Test_ItfexpirationUsecase_CreateReminderITFExpireEmailTask(t *testing.T) {
	a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	UT.ItfexpirationUsecase.CreateReminderITFExpireEmailTask(*a)
}

func Test_ItfexpirationUsecase_CreateReminderITFExpireTextTask(t *testing.T) {
	a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	UT.ItfexpirationUsecase.CreateReminderITFExpireTextTask(*a)
}

func Test_ItfexpirationUsecase_WaitingReminderCases(t *testing.T) {
	destTime := time.Now().In(configs.GetVBCDefaultLocation())
	destTime = destTime.AddDate(0, 0, 90)
	destITFDate := destTime.Format(time.DateOnly)
	abc, err := UT.ItfexpirationUsecase.WaitingReminderCases(destITFDate)
	lib.DPrintln(abc, err)
}
