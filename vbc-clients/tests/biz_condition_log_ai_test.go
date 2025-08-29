package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ConditionLogAiUsecase_AddLogPromptAiCondition(t *testing.T) {
	err := UT.ConditionLogAiUsecase.AddLogPromptAiCondition(1, 1, "22")
	lib.DPrintln(err)
}
