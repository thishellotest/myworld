package tests

import (
	"testing"
	"vbc/lib"
)

func Test_LeadConversionSummaryBuzUsecase_DoOne(t *testing.T) {
	err := UT.LeadConversionSummaryBuzUsecase.DoOne(5511)
	lib.DPrintln(err)
}

func Test_LeadConversionSummaryBuzUsecase_ManualAll(t *testing.T) {
	err := UT.LeadConversionSummaryBuzUsecase.ManualAll()
	lib.DPrintln(err)
}
