package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AttorneybuzUsecase_GetAnAttorney(t *testing.T) {
	a, err := UT.AttorneybuzUsecase.GetAnAttorney()
	lib.DPrintln(a.ToContractAttorneyVo())
	lib.DPrintln(err)

}
