package tests

import (
	"testing"
	"vbc/lib"
)

func Test_HttpSettingsUsecase_BizHttpFields(t *testing.T) {
	aa, err := UT.HttpSettingsUsecase.BizHttpFields("")
	lib.DPrintln(aa, err)
}
