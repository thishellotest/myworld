package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ConditionSourcebuzUsecase_Handle(t *testing.T) {
	err := UT.ConditionSourcebuzUsecase.Handle()
	lib.DPrintln(err)
}
