package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZohoDealScan2JobUsecase_BatchHandle(t *testing.T) {
	var lastModifiedTime string
	UT.ZohoDealScan2JobUsecase.BatchHandle(&lastModifiedTime, "", 1)
}

func Test_ZohoDealScan2JobUsecase_BizRunJob(t *testing.T) {
	err := UT.ZohoDealScan2JobUsecase.BizRunJob()
	lib.DPrintln(err)
}
