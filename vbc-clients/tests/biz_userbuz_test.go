package tests

import (
	"testing"
	"vbc/lib"
)

func Test_UserbuzUsecase_HandleAllPassword(t *testing.T) {
	err := UT.UserbuzUsecase.HandleAllPassword()
	lib.DPrintln("err:", err)
}
