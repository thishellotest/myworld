package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RecordLogbuzUsecase_ManualHandleDueDate(t *testing.T) {
	err := UT.RecordLogbuzUsecase.ManualHandleDueDate()
	if err != nil {
		panic(err)
	}
}

func Test_RecordLogbuzUsecase_ManualHandleDueDateRow(t *testing.T) {

	a, err := UT.TUsecase.DataByGid(biz.Kind_client_cases, "6159272000001008007")
	if err != nil {
		panic(err)
	}

	err = UT.RecordLogbuzUsecase.ManualHandleDueDateRow(a)
	lib.DPrintln(err)
}
