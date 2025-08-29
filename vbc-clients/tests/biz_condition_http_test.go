package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ConditionHttpUsecase_BizList(t *testing.T) {
	r, err := UT.ConditionHttpUsecase.BizList(1, 20, nil)
	lib.DPrintln(r, err)
}

func Test_ConditionHttpUsecase_BizSources(t *testing.T) {
	r, err := UT.ConditionHttpUsecase.BizSources(1012)
	lib.DPrintln(r, err)
}
