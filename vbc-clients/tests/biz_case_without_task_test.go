package tests

import (
	"testing"
	"vbc/lib"
)

func Test_CaseWithoutTaskUsecase_GetCasesForNotify(t *testing.T) {
	UT.CaseWithoutTaskUsecase.GetCasesForNotify()
}

func Test_CaseWithoutTaskUsecase_GetCases(t *testing.T) {
	a, err := UT.CaseWithoutTaskUsecase.GetCases()
	lib.DPrintln(err)
	lib.DPrintln(a)
}

func Test_CaseWithoutTaskUsecase_CaseWithoutTaskUsecase(t *testing.T) {
	a, err := UT.CaseWithoutTaskUsecase.NotifyEmailBody()
	lib.DPrintln(a, err)
}

func Test_CaseWithoutTaskUsecase_ReminderManager(t *testing.T) {
	err := UT.CaseWithoutTaskUsecase.ReminderManager()
	lib.DPrintln(err)
}
