package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	//. "vbc/lib/builder"
)

func Test_ActionOnceUsecase_StageGettingStartedEmailToAwaitingClientFiles(t *testing.T) {

	err := UT.ActionOnceUsecase.StageGettingStartedEmailToAwaitingClientFiles(5005)
	lib.DPrintln(err)

}

func Test_ActionOnceUsecase_HandleDataCollectionFolder(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleDataCollectionFolder(5160)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleCopyRecordReviewFiles(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleCopyRecordReviewFiles(56)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_NoPrimaryCaseInit(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5058)
	lib.DPrintln(err)
	err = UT.ActionOnceUsecase.NoPrimaryCaseInit(tCase)
	lib.DPrintln(err)
}

/*
func Test_ActionOnceUsecase_HandleReleaseOfInformation(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleReleaseOfInformation(5005)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleCopyReleaseOfInformationFile(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleCopyReleaseOfInformationFile(5005)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandlePatientPaymentForm(t *testing.T) {
	err := UT.ActionOnceUsecase.HandlePatientPaymentForm(5005)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleCopyPatientPaymentFormFile(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleCopyPatientPaymentFormFile(5005)
	lib.DPrintln(err)
}*/

func Test_ActionOnceUsecase_BizCopyMedicalTeamForms(t *testing.T) {

	//tClientCase, err := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 5004})
	//if err != nil {
	//	panic(err)
	//}
	//
	//tClient, _, err := UT.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//a1, a2, err := UT.ActionOnceUsecase.BizCopyMedicalTeamForms("1550028092926", tClient, "264658754393")
	//lib.DPrintln(a1, a2, err)
}

func Test_ActionOnceUsecase_HandlePrivateExamsSubmittedFirstStep(t *testing.T) {

	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	lib.DPrintln(err)
	err = UT.ActionOnceUsecase.HandlePrivateExamsSubmittedFirstStep(tCase)
	lib.DPrintln("HandlePrivateExamsSubmittedFirstStep:", err)
}

func Test_ActionOnceUsecase_HandlePrivateExamsSubmitted(t *testing.T) {
	err := UT.ActionOnceUsecase.HandlePrivateExamsSubmitted(context.TODO(), 75)
	lib.DPrintln("HandlePrivateExamsSubmitted:", err)
}

func Test_ActionOnceUsecase_HandleMedicalTeamForms(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleMedicalTeamForms(5301)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleMedicalTeamFormsTest(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleMedicalTeamFormsTest(5301)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandlePersonalStatementsFile(t *testing.T) {
	err := UT.ActionOnceUsecase.HandlePersonalStatementsFile(5302)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleDoDocEmailFile(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleDoDocEmailFile(5301)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleDoCopyDocEmailFile(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleDoCopyDocEmailFile(5586)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleMedicalTeamFormsReminderEmail(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleMedicalTeamFormsReminderEmail(5105)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_MultiCasesBaseInfoSync(t *testing.T) {
	err := UT.ActionOnceUsecase.MultiCasesBaseInfoSync(5514)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleUpcomingContactInformation(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleUpcomingContactInformation(5112)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_DoUpcomingContactInformation(t *testing.T) {
	err := UT.ActionOnceUsecase.DoUpcomingContactInformation(5112)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_DoCopyPersonalStatementsDoc(t *testing.T) {
	err := UT.ActionOnceUsecase.DoCopyPersonalStatementsDoc(5293)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_DoPersonalStatementsReadyforYourReview(t *testing.T) {
	UT.ActionOnceUsecase.DoPersonalStatementsReadyforYourReview(5511)
}

func Test_ActionOnceUsecase_DoPleaseReviewYourPersonalStatementsinSharedFolder(t *testing.T) {
	UT.ActionOnceUsecase.DoPleaseReviewYourPersonalStatementsinSharedFolder(5511)
}

func Test_ActionOnceUsecase_DoPleaseReviewYourPersonalStatementsinSharedFolder1(t *testing.T) {
	err := UT.ActionOnceUsecase.CancelAutomationCrontabEmailTasks(5511)
	lib.DPrintln(err)
}

func Test_ActionOnceUsecase_HandleHelpUsImproveSurvey(t *testing.T) {
	err := UT.ActionOnceUsecase.HandleHelpUsImproveSurvey(5511)
	if err != nil {
		panic(err)
	}
}

func Test_ActionOnceUsecase_DoHandleHelpUsImproveSurvey(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)

	err = UT.ActionOnceUsecase.DoHandleHelpUsImproveSurvey(*tCase)
	if err != nil {
		panic(err)
	}
}

func Test_ActionOnceUsecase_HandleVAForm2122aSubmission(t *testing.T) {
	UT.ActionOnceUsecase.HandleVAForm2122aSubmission(5511)
}
