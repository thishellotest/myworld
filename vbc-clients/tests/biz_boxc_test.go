package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_BoxcUsecase_GetNewEvidenceFolderId(t *testing.T) {

	primaryCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5076)
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5093)

	a, err := UT.BoxcUsecase.GetNewEvidenceFolderId(primaryCase, tCase)
	lib.DPrintln(a, err)
}
