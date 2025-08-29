package tests

import (
	"testing"
	"vbc/lib"
)

func Test_VbcDataVerifyUsecase_VerifyContract(t *testing.T) {
	err := UT.VbcDataVerifyUsecase.VerifyContract()
	lib.DPrintln(err)
}
