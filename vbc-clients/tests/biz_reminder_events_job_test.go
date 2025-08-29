package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ReminderEventsJobUsecase_Handle(t *testing.T) {
	err := UT.ReminderEventsJobUsecase.Handle()
	lib.DPrintln(err)
}

func Test_ReminderEventsJobUsecase_FinishReminder(t *testing.T) {
	lib.DPrintln("ccc  aaa", "ccc")
}
