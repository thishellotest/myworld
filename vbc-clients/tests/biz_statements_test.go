package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_SplitCaseStatements(t *testing.T) {
	str := `70 - Insomnia Disorder (increase)

30 - Bilateral pes planus aggravated by service (str)

20 - Back pain -secondary to bilateral pes planus (opinion, BVA)

10 - Hypertension (str)

10 - Tinnitus (new)

0 - Erectile  dysfunction secondary to (new) insomnia disorder (opinion, BVA)
0 - Erectile  dysfunction secondary to (new222) insomnia disorder
0 - Erectile  dysfunction secondary to (first) (second)
0 - Erectile  dysfunction secondary to insomnia disorder (final)
0 - Erectile  dysfunction secondary to (new222) insomnia disorder
0 - Erectile  dysfunction secondary to (new222) insomnia disorder()
50 - Headaches aggravated by Bipolar disorder (opinion)

30 - IBS chronic diarrhea (str)

20 - Lumbosacral Strain (increase) 

20 - Radiculopathy in left lower extremity secondary to Lumbosacral Strain (opinion)

20 - Radiculopathy in right lower extremity secondary to Lumbosacral Strain (opinion)

20 - Left hip pain (str)
-------Supplemental------------
20 - Right knee strain (str)

20 - Left knee pain (str)

10 - TMJ Disorder (str)
-- NO PRIVATE EXAMS --
10 - Shortness of breath
NO PRIVATE EXAMS
0 – Erectile dysfunction (ED) secondary to Bipolar disorder



`
	a, err := biz.SplitCaseStatements(str)
	lib.DPrintln(a, err)
}

func Test_SplitCaseStatements2(t *testing.T) {

	str := `
NO PRIVATE EXAMS
0 - Erectile dysfunction (ED) secondary to Bipolar disorder`

	str = `70 - Insomnia Disorder1233333 (increase)
0 - Erectile  dysfunction secondary to insomnia disorder (opinion, BVA)
-------Supplemental-------
30 - Bilateral pes planus aggravated by service (str)
10 - Tinnitus (new)
-------NO PRIVATE EXAMS-------
20 - Back pain secondary to bilateral pes planus (opinion, BVA)
10 - Hypertension (sssss)
`
	a, err := biz.SplitCaseStatements(str)
	lib.DPrintln(a, err)
}

func Test_StatementUsecase_GetStatementBaseInfo(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	r, err := UT.StatementUsecase.GetStatementBaseInfo(tCase, tClient)
	lib.DPrintln(r, err)
}

func Test_StatementUsecase_HandleStatementToBox(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5662)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	veteranSummary, err := UT.AiTaskbuzUsecase.HandleVeteranSummary(context.TODO(), tCase)
	if err != nil {
		panic(err)
	}
	err = UT.StatementUsecase.HandleStatementToBox(tCase, tClient, veteranSummary)
	lib.DPrintln(err)
}

func Test_StatemtUsecase_AllLatestStatements(t *testing.T) {
	res, err := UT.StatemtUsecase.AllLatestStatements(5511)
	lib.DPrintln(res, err)
}

func Test_StatemtUsecase_ObtainAvailableVersionID(t *testing.T) {
	res, err := UT.StatemtUsecase.ObtainAvailableVersionID(0)
	lib.DPrintln(res, err)
}

func Test_StatementUsecase_BizStatementDetail(t *testing.T) {
	a, _ := UT.UserUsecase.GetUserFacadeById(4)
	res, err := UT.StatementUsecase.BizStatementDetail(false, &a.TData, "6159272000013713061")
	lib.DPrintln(res, err)
}

func Test_StatementUsecase_GetListStatementDetail(t *testing.T) {
	//a, _ := UT.UserUsecase.GetUserFacadeById(4)
	//a, _ := UT.TUsecase.DataByGid(biz.Kind_client_cases, "d1fbcc1328424c3699057dd71f14e970")
	//res, err := UT.StatementUsecase.GetListStatementDetail(a)
	//lib.DPrintln(res, err)
}

func Test_StatementUsecase_PersonalStatementPassword(t *testing.T) {
	a, err := UT.StatementUsecase.PersonalStatementPassword(5511)
	lib.DPrintln(a, err)
}

func Test_ParseAiStatementCondition(t *testing.T) {
	vo := biz.ParseAiStatementCondition(StatementStringVar)
	lib.DPrintln(vo)
}

func Test_StatementUsecase_GenerateNewStatementVersion(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))

	//veteranSummary, err := UT.AiTaskbuzUsecase.HandleVeteranSummary(context.TODO(), tCase)
	//if err != nil {
	//	panic(err)
	//}
	err := UT.StatementUsecase.GenerateNewStatementVersion(*tCase, *tClient)
	lib.DPrintln(err)
}

func Test_StatementUsecase_BizStatementVerifyPassword(t *testing.T) {

	isOk, err := UT.StatementUsecase.BizStatementVerifyPassword(5511, "ndd8ZHyn")
	lib.DPrintln(isOk, err)
}

func Test_StatementUsecase_GetUpdatePSTextForAiParam(t *testing.T) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5572)
	//tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	//
	//a, _ := UT.StatementConditionUsecase.GetByCond(builder.Eq{"id": 27})
	//referenceContent, text, err := UT.StatementUsecase.GetUpdatePSTextForAiParam(*tClient, *tCase, *a)
	//lib.DPrintln(referenceContent)
	//lib.DPrintln(text)
	//lib.DPrintln(err)
}

func Test_ParseAiVeteranSummary(t *testing.T) {
	//VeteranSummaryVar
	a := `# Veteran Summary

• **Full Name:** Katelan McClough
• **Unique ID:** 5572
• **Branch of Service:** Marine Corps
• **Years of Service:** 2018-2024
• **Retirement Status:** Did not retire
• **Deployments:** No deployments to presumptive list areas
• **Marital Status:** Not married, no divorces
• **Children:** None
• **Occupation in Service:** 3521 Automotive maintenance technician

Thank you for submitting your initial information. To create personalized VA disability personal statements, I'll need additional information about specific conditions for which you're seeking disability benefits. Please provide details about each condition, including when it started, current symptoms, treatments, and how it impacts your daily life.`
	baseInfo := biz.ParseAiVeteranSummary(a)
	lib.DPrintln(baseInfo)
}

func Test_StatementUsecase_GetNewStatement(t *testing.T) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5662)
	//tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	//
	//veteranSummary, err := UT.AiTaskJobUsecase.HandleVeteranSummary(context.TODO(), tCase)
	//if err != nil {
	//	panic(err)
	//}
	//
	//abc, err := UT.StatementUsecase.GetNewStatement(*tCase, *tClient, veteranSummary)
	//lib.DPrintln(abc)
	//lib.DPrintln(err)
}

func Test_StatementUsecase_NeedUseNewPersonalWebForm(t *testing.T) {
	//a, err := UT.StatementUsecase.NeedUseNewPersonalWebForm(5511)
	//lib.DPrintln(a, err)
}

func Test_StatementUsecase_SaveCaseStatement(t *testing.T) {

	//caseGid := "cc6cbd62537843fc93169e9091c9a9ae"
	//str := `[{"id":"12","rating":"3","condition":"222288888","association":"increase","category":"General"},{"id":"new_1752662048782","rating":"2","condition":"22","association":"increase","category":"Supplemental"},{"id":"new_1752662043207","rating":"1","condition":"11","association":"new","category":"NO PRIVATE EXAMS"}]`
	//err := UT.StatementUsecase.SaveCaseStatement(caseGid, str)
	//lib.DPrintln(err)
}

func Test_StatementUsecase_ManualSyncStatement(t *testing.T) {
	gid := "6159272000003883233"
	err := UT.StatementUsecase.ManualSyncStatement(gid)
	lib.DPrintln(err)
}
