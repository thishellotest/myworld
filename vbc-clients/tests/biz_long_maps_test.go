package tests

import (
	"testing"
	"vbc/lib"
)

func Test_LongMapUsecase_GetForString(t *testing.T) {
	aa, err := UT.LongMapUsecase.GetForString("sck1")
	lib.DPrintln(aa, err)
}
