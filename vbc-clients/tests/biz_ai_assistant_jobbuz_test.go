package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AiAssistantJobBuzUsecase_BizHttpApplyJob(t *testing.T) {
	jobUuid := "statementSection:5511:7:Medication"
	jobUuid = "statementCondition:5511:10"
	a, err := UT.AiAssistantJobBuzUsecase.BizHttpApplyJob(jobUuid)
	lib.DPrintln(a, err)
}
