package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_OptionUsecase_JotformIdsOptions(t *testing.T) {
	//UT.
	tCase, _ := UT.TUsecase.DataByGid(biz.Kind_client_cases, "d1fbcc1328424c3699057dd71f14e970")
	a, err := UT.OptionUsecase.JotformIdsOptions(tCase, "1", "")
	lib.DPrintln(a, err)
}

func Test_OptionUsecase_NewJotformIdsOptions(t *testing.T) {
	//UT.
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	a, err := UT.OptionUsecase.NewVersionJotformIdsOptions(tCase, "")
	for _, v := range a {
		lib.DPrintln(v.OptionLabel)
	}
	lib.DPrintln(err)
	//lib.DPrintln(a, err)
}
