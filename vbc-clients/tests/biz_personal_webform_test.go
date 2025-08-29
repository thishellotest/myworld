package tests

import (
	"testing"
	"vbc/lib"
)

func Test_PersonalWebformUsecase_NeedUseNewPersonalWebForm(t *testing.T) {
	flag, err := UT.PersonalWebformUsecase.NeedUseNewPersonalWebForm(5222)
	lib.DPrintln(flag, err)
}

func Test_PersonalWebformUsecase_ManualHistoryData(t *testing.T) {
	UT.PersonalWebformUsecase.ManualHistoryData()
}
