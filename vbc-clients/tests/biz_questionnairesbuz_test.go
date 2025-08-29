package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_QuestionnairesbuzUsecase_Manual(t *testing.T) {
	err := UT.QuestionnairesbuzUsecase.Manual()
	lib.DPrintln(err)
}

func Test_QuestionnairesbuzUsecase_GetJotformSubmissionsForGenStatement(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)

	condition, _ := UT.ConditionUsecase.GetByCond(builder.Eq{
		"condition_name": "Tinnitus",
		"type":           biz.Condition_Type_Condition})
	a, b, err := UT.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatement(tCase, condition)
	lib.DPrintln(a, b, err)
}

func Test_QuestionnairesbuzUsecase_GetJotformSubmissionsForGenStatementNew(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	var StatementCondition biz.StatementCondition
	StatementCondition.ConditionValue = "Tinnitus"

	_, b, err := UT.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatementNew(tCase, StatementCondition)
	lib.DPrintln(err)

	for _, v := range b {
		lib.DPrintln(v.SubmissionId)
	}
}
