package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AiHttpUsecase_BizTaskHandle(t *testing.T) {
	r, err := UT.AiHttpUsecase.BizTaskHandle(5511)
	lib.DPrintln(r, err)
}

func Test_AiHttpUsecase_BizTasks(t *testing.T) {
	a, err := UT.AiHttpUsecase.BizTasks(5431)
	lib.DPrintln(a, err)
}

func Test_BizTaskRenew(t *testing.T) {
	a, err := UT.AiHttpUsecase.BizTaskRenew(6)
	lib.DPrintln(a, err)
}
