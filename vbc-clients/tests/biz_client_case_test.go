package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientCaseUsecase_PrimaryCase(t *testing.T) {
	a, err := UT.ClientCaseUsecase.PrimaryCase("6159272000004377287")
	lib.DPrintln(a, err)
}

func Test_ClientCaseUsecase_NotPrimaryCases(t *testing.T) {
	a, err := UT.ClientCaseUsecase.NotPrimaryCases("ssss", "bb1")
	lib.DPrintln(a, err)
}

func Test_ClientCaseUsecase_CurrentCaseInProgress(t *testing.T) {
	a, err := UT.ClientCaseUsecase.CurrentCaseInProgress("6159272000001008204")
	lib.DPrintln(a, err)
}

func Test_ClientCaseUsecase_GetPricingVersion(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	lib.DPrintln(err)
	pv, err := UT.ClientCaseUsecase.GetPricingVersion(tCase)
	lib.DPrintln(pv, err)
}

func Test_ClientCaseUsecase_SavePricingVersion(t *testing.T) {

	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	lib.DPrintln(err)
	err = UT.ClientCaseUsecase.SavePricingVersion(tCase, "v201")
	lib.DPrintln(err)
}

func Test_ClientCaseUsecase_EnabledTwoBySMS(t *testing.T) {

	flag, err := UT.ClientCaseUsecase.EnabledTwoBySMS(5004)
	lib.DPrintln(flag, err)
}

func Test_ClientCaseUsecase_BizHttpClaimsInfo(t *testing.T) {
	a, err := UT.ClientCaseUsecase.BizHttpClaimsInfo("2483184881")
	lib.DPrintln(a, err)
}

func Test_ClientCaseUsecase(t *testing.T) {
	a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5369)
	c := a.CustomFields.TextValueByNameBasic("deal_name")
	c1 := a.CustomFields.TextValueByNameBasic("collaborators")
	lib.DPrintln(c, c1)
}

func Test_ClientCaseUsecase_GetByPhone(t *testing.T) {
	a, err := UT.ClientCaseUsecase.GetByPhone("+18134594230")
	lib.DPrintln(a, err)
	if a != nil {
		lib.DPrintln(a.CustomFields.TextValueByNameBasic(biz.FieldName_primary_vs))
	}
}

func Test_ClientCaseUsecase_GetLeadVSByPhone(t *testing.T) {
	a, err := UT.ClientCaseUsecase.GetLeadVSByPhone("+18134594230")
	lib.DPrintln(err)
	if a != nil {
		lib.DPrintln(a.CustomFields.TextValueByNameBasic(biz.UserFieldName_fullname))
		lib.DPrintln(a.CustomFields.TextValueByNameBasic(biz.UserFieldName_dialpad_userid))
	}
}

func Test_ClientCaseUsecase_ItfCases(t *testing.T) {
	a, err := UT.ClientCaseUsecase.ItfCases()
	lib.DPrintln(a, err)
}

func Test_ClientCaseUsecase_ItfCasesByUserGid(t *testing.T) {
	UT.ClientCaseUsecase.ItfCasesByUserGid("c9ce3ecee21640e7978a373c08d21292")
}
