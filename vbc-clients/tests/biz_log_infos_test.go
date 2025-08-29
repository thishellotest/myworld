package tests

import (
	"testing"
	"vbc/lib"
)

func Test_LogInfoUsecase_SaveLogInfo(t *testing.T) {
	err := UT.LogInfoUsecase.SaveLogInfo(0, "aaa", map[string]interface{}{
		"a": 1,
	})
	lib.DPrintln(err)
}
