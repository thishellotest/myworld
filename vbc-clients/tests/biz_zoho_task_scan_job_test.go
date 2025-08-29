package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZohoTaskScanJobUsecase_BizRunJob(t *testing.T) {
	err := UT.ZohoTaskScanJobUsecase.BizRunJob()
	lib.DPrintln(err)
}
