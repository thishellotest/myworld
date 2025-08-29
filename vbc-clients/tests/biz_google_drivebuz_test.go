package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_GoogleDrivebuzUsecase_TransferPaymentForm(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5123)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	file, err := UT.GoogleDrivebuzUsecase.TransferPaymentForm(context.TODO(), tCase, tClient)
	lib.DPrintln(err)
	lib.DPrintln(file)
}

func Test_GoogleDrivebuzUsecase_TransferPsych(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5123)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	err := UT.GoogleDrivebuzUsecase.TransferPsych(context.TODO(), tCase, tClient)
	lib.DPrintln(err)
}

func Test_GoogleDrivebuzUsecase_TransferGeneral(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5123)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	err := UT.GoogleDrivebuzUsecase.TransferGeneral(context.TODO(), tCase, tClient)
	lib.DPrintln(err)
}
