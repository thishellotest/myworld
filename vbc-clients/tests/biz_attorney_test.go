package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AttorneyUsecase_GetByName(t *testing.T) {
	a, err := UT.AttorneyUsecase.GetByName("Niralkumar Patel")
	if a != nil {
		lib.DPrintln(a.ToContractAttorneyVo())
	}
	lib.DPrintln(err)
}
