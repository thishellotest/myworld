package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_TaskCreateUsecase_CreateCustomTaskMail(t *testing.T) {
	message := &biz.MailMessage{
		To:      "18891706@qq.com",
		Subject: "test",
		Body:    "body",
	}
	UT.TaskCreateUsecase.CreateCustomTaskMail(0, message, 0)
}

func Test_TaskCreateUsecase_CreateTaskMail(t *testing.T) {
	err := UT.TaskCreateUsecase.CreateTaskMail(12, biz.MailGenre_SignFeeContractFirstRemind, 0, nil, 0, "", "")
	lib.DPrintln(err)
}

func Test_TaskCreateUsecase_CreateTask1(t *testing.T) {
	clientCaseId := int32(5004)
	UT.TaskCreateUsecase.CreateTask(clientCaseId,
		map[string]interface{}{"CaseId": clientCaseId},
		biz.Task_Dag_ReminderMedicalTeamFormsContractSent, 0, "", "")
}

func Test_TaskCreateUsecase_CreateTaskWithFrom(t *testing.T) {
	UT.TaskCreateUsecase.CreateTaskWithFrom(5004,
		biz.CronTriggerVo{
			HandleSendSMSType: biz.HandleSendSMSTextGettingStartedEmail,
		}, biz.Task_Dag_CronTrigger, 0,
		biz.Task_FromType_DialpadSMS, "5004")
}
