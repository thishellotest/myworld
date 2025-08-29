package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_DbqsUsecase_BizLeadCPEmail(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
	a, err := UT.DbqsUsecase.BizLeadCPEmail(tCase)
	lib.DPrintln(a, err)
}
func Test_DbqsUsecase_LeadCPEmail(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
	a, err := UT.DbqsUsecase.LeadCPEmail(tCase)
	lib.DPrintln(a, err)
}

func Test_DbqsUsecase_LeadVSEmail(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
	a, err := UT.DbqsUsecase.LeadVSEmail(tCase)
	lib.DPrintln(a, err)
}

func Test_DbqsUsecase_ReleaseOfInformationPrefillTags(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
	client, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	res, err := UT.DbqsUsecase.ReleaseOfInformationPrefillTags(tCase, client)

	lib.DPrintln(res, err)
}

func Test_DbqsUsecase_MedicalTeamFormsPrefillTagsV2(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5301)
	client, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	res, boxSignTmpId, err := UT.DbqsUsecase.MedicalTeamFormsPrefillTagsV2(tCase, client)

	lib.DPrintln(res, err)
	lib.DPrintln(boxSignTmpId)
}
