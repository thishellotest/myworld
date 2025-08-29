package tests

import (
	"testing"
	"vbc/lib"
)

func Test_BoxUserBuzUsecase_ScanBoxUser(t *testing.T) {
	err := UT.BoxUserBuzUsecase.ScanBoxUser()
	lib.DPrintln(err)
}
