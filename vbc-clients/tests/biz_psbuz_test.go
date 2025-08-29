package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_PsbuzUsecase_HandleUpdateStatement(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	err := UT.PsbuzUsecase.HandleUpdateStatement(*tClient, *tCase)
	lib.DPrintln(err)
}
