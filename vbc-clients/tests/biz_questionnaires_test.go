package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_QuestionnairesUsecase_BizHttpList(t *testing.T) {
	a, err := UT.QuestionnairesUsecase.BizHttpList(5280)
	lib.DPrintln(a, err)
}

func Test_QuestionnairesUsecase_t(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5280)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))

	url := UT.QuestionnairesUsecase.LinkInitialIntake(*tClient, *tCase)
	lib.DPrintln(url)
}
