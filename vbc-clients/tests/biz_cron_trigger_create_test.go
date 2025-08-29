package tests

import (
	"testing"
	"vbc/lib"
)

func Test_CronTriggerCreateUsecase_CreateSendSMSTextMedTeamForms(t *testing.T) {
	err := UT.CronTriggerCreateUsecase.CreateSendSMSTextMedTeamForms(5004)
	lib.DPrintln(err)
}

func Test_CronTriggerCreateUsecase_CreateGettingStartedEmail(t *testing.T) {
	err := UT.CronTriggerCreateUsecase.CreateGettingStartedEmailByCaseId(5511)
	lib.DPrintln(err)
}
