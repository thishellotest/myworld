package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_AssistantUsecase_GetJsonFromAiResultForAssistant(t *testing.T) {
	str := "After reviewing the information provided, I need to update the VA statement because new treatment facility and medication information has been provided.\n\n```json\n{\n  \"update_required\": true,\n  \"updated_statement\": \"I am currently taking Bupropion 300 MG and Escitalopram 30 MG for my depression, though I have been experiencing deficiencies in most areas of my life including work, school, family relations, and other important aspects of daily functioning. The lack of treatment has allowed my condition to progress unchecked, worsening the symptoms and their impact on my daily life.\"\n}\n```"
	a := biz.GetJsonFromAiResultForAssistant(str)
	lib.DPrintln(a)
}

func Test_AssistantUsecase_HandleSaveAllStatements(t *testing.T) {
	lib.DPrintln("ss")
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	tClient, _, _ := UT.DataComboUsecase.ClientWithCase(*tCase)
	UT.AssistantUsecase.HandleSaveAllStatements(*tCase, *tClient)
}

func Test_AssistantUsecase_HandleAssistant(t *testing.T) {
	task, _ := UT.AiTaskUsecase.GetByCond(builder.Eq{"id": 1932})
	err := UT.AssistantUsecase.HandleAssistant(context.TODO(), task)
	lib.DPrintln(err)
}
