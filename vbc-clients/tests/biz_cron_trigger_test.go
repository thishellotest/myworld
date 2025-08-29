package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_CronTriggerUsecase_Handle(t *testing.T) {
	var cronTriggerVo biz.CronTriggerVo
	err := UT.CronTriggerUsecase.Handle(biz.HandleSendSMSTextStatementFinalizedEvery14Days, 5301, cronTriggerVo)
	lib.DPrintln(err)
}

func Test_CronTriggerUsecase_VerifyClientTasksCondition(t *testing.T) {
	entity, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5530)
	verify, err := UT.CronTriggerUsecase.VerifyClientTasksCondition(biz.HandleSendSMSTextAwaitingClientRecordsLongerThan30Days, entity)
	lib.DPrintln(verify, err)
}
