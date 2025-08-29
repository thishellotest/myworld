package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_JotformSubmissionUsecase_ManualHandleFormId(t *testing.T) {
	err := UT.JotformSubmissionUsecase.ManualHandleFormId()
	if err != nil {
		panic(err)
	}
}

func Test_JotformSubmissionUsecase_AllLatestByUniqcode(t *testing.T) {
	res, err := UT.JotformSubmissionUsecase.AllLatestByUniqcode("2517537132")
	lib.DPrintln(res, err)
}

func Test_JotformSubmissionUsecase_AllLatestByUniqcodeExceptIntake(t *testing.T) {
	//tCase, _ := UT.TUsecase.Data(biz.Kind_client_cases, builder.Eq{"uniqcode": "2517537132"})
	res, err := UT.JotformSubmissionUsecase.AllLatestByUniqcodeExceptIntake([]string{"2517537132"})
	lib.DPrintln(err)
	for _, v := range res {
		aa, err := v.JotformNewFileNameForAI()
		lib.DPrintln(aa, err)
	}
}

func Test_JotformSubmissionUsecase_AllLatestUpdateQuestionnaires(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5373)
	uniqcode := tCase.CustomFields.TextValueByNameBasic(biz.FieldName_uniqcode)
	res, _ := UT.JotformSubmissionUsecase.AllLatestUpdateQuestionnaires(uniqcode)
	for _, v := range res {
		aa, err := v.JotformNewFileNameForAI()
		lib.DPrintln(aa, err)
	}
	lib.DPrintln("length:", len(res))
}
