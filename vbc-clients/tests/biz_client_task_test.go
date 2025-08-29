package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ClientTaskUsecase_TasksByCaseGid(t *testing.T) {
	records, err := UT.ClientTaskUsecase.TasksByCaseGid("aaa")
	lib.DPrintln(records, err)
}

func Test_ClientTaskUsecase_HandleAutomationCompleteTask(t *testing.T) {
	err := UT.ClientTaskUsecase.HandleAutomationCompleteTask(5102)
	lib.DPrintln(err)
}

func Test_ClientTaskUsecase_DueDates(t *testing.T) {
	aa, err := UT.ClientTaskUsecase.DueDates("what_id_gid", []string{"6159272000009972111", "6159272000009898006", "sss"})
	lib.DPrintln(aa)
	lib.DPrintln(err)
}

func Test_ClientTaskUsecase_DueDatesByWhatGids(t *testing.T) {
	aa, err := UT.ClientTaskUsecase.DueDatesByWhatGids([]string{"6159272000009972111", "6159272000009898006", "sss"})
	lib.DPrintln(aa)
	lib.DPrintln(err)
}

func Test_ClientTaskUsecase_DueDatesByWhoGids(t *testing.T) {
	aa, err := UT.ClientTaskUsecase.DueDatesByWhoGids([]string{"6159272000005519042", "6159272000009898006", "sss"})
	lib.DPrintln(aa)
	lib.DPrintln(err)
}
