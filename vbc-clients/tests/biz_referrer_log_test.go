package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ReferrerLogUsecase_ReferringClient(t *testing.T) {
	err := UT.ReferrerLogUsecase.ReferringClient("s", "cc", "aa")
	lib.DPrintln(err)
}

func Test_ReferrerLogUsecase_Upsert(t *testing.T) {
	err := UT.ReferrerLogUsecase.Upsert("a", "s", "cc", "aa")
	lib.DPrintln(err)
}
