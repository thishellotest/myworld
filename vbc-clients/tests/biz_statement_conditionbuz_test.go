package tests

import (
	"testing"
	"vbc/lib"
)

func Test_StatementConditionBuzUsecase_InitStatementCondition(t *testing.T) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	//err := UT.StatementConditionBuzUsecase.InitStatementCondition(*tCase)
	//lib.DPrintln(err)
}

func Test_StatementConditionBuzUsecase_StatementConditionBuzUsecase(t *testing.T) {
	err := UT.StatementConditionBuzUsecase.UpdateCaseStatement(5004)
	lib.DPrintln(err)
}
