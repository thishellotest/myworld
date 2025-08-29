package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZohoContactScanJobUsecase_BizRunJob(t *testing.T) {
	err := UT.ZohoContactScanJobUsecase.BizRunJob()
	lib.DPrintln(err)
}
