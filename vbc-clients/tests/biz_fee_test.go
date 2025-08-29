package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_FeeUsecase_ClientCaseAmount(t *testing.T) {
	tClientCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 122)
	amount, err, noVo := UT.FeeUsecase.ClientCaseAmount(tClientCase)
	lib.DPrintln(amount, err, noVo)
}

func Test_FeeUsecase_InvoiceAmount(t *testing.T) {
	amount, err := UT.FeeUsecase.InvoiceAmount(5004, false, 0, 80)
	lib.DPrintln(amount, err)
}

func Test_FeeUsecase_NotPrimaryCaseAmount(t *testing.T) {
	tClientCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5051)
	amount, err := UT.FeeUsecase.NotPrimaryCaseAmount(tClientCase, 100)
	lib.DPrintln(amount, err)
}

func Test_FeeUsecase_VBCFees(t *testing.T) {
	tClientCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	a, err := UT.FeeUsecase.VBCFees(tClientCase)
	lib.DPrintln(a, err)
}

func Test_FeeUsecase_GetIncreaseAmount(t *testing.T) {
	a, err := UT.FeeUsecase.GetIncreaseAmount(100, 100)
	lib.DPrintln(a, err)
}
