package tests

import (
	"encoding/json"
	"testing"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
	"vbc/lib/to"
)

func Test_RemindUsecase_CreateUnfinishedFeeContract(t *testing.T) {
	err := UT.RemindUsecase.CreateUnfinishedFeeContract(12)
	lib.DPrintln(err)
}

func Test_RemindUsecase_CreateUnfinishedIntakeForm(t *testing.T) {
	err := UT.RemindUsecase.CreateUnfinishedIntakeForm(22)
	lib.DPrintln(err)
}

func Test_RemindUsecase_FollowingUpSignMedicalTeamFormsEmailBody(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	c, err := UT.RemindUsecase.FollowingUpSignMedicalTeamFormsEmailBody(tCase, tUser, to.Ptr(time.Now()))
	lib.DPrintln(c, err)
}

func Test_RemindUsecase_FollowingUpUploadedDocumentEmailBody(t *testing.T) {

	items, _ := UT.ReminderEventUsecase.AllByCond(builder.In("id", 9, 10))

	var updateFiles []*biz.ReminderClientUpdateFilesEventVoItem
	for _, v := range items {
		r := v.GetReminderClientUpdateFilesEventVo()
		updateFiles = append(updateFiles, r.Items...)
	}

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	//tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	lib.DPrintln("updateFiles:", updateFiles)
	c, err := UT.RemindUsecase.FollowingUpUploadedDocumentEmailBody(tCase, "lialing@foxmail.com", updateFiles)
	lib.DPrintln(c, err)

	serviceConfig := &biz.MailServiceConfig{
		Name:        "Dev",
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    "glliao@vetbenefitscenter.com",
		Password:    configs.EnvMailGlliaoPWD(),
		FromAddress: "glliao@vetbenefitscenter.com",
	}
	MailMessage := &biz.MailMessage{
		To: "lialing@foxmail.com",
		//To:      "liaogling@gmail.com",
		Subject: c.Subject,
		Body:    c.Body,
	}

	err = UT.MailUsecase.SendEmail(serviceConfig, MailMessage, "", nil)
	lib.DPrintln("SendEmail:", err)
}

func Test_RemindUsecase_CaseWithoutTasksEmailBody(t *testing.T) {

	subject := "VBC: List of Client Cases without a Task"

	var result []*biz.CaseWithoutTaskVo
	str := `[{"StagesName":"Getting Started Email","Items":[{"ClientCaseName":"Shi Li@Case1","CreatedTime":"2024-03-29T23:06:04+08:00","Gid":"6159272000000820046","IsNew":false}]},{"StagesName":"Fee Schedule and Contract","Items":[{"ClientCaseName":"TestForZohoLn TestForZohoOwner","CreatedTime":"2024-04-01T22:54:51+08:00","Gid":"6159272000000881012","IsNew":true}]},{"StagesName":"1. Fee Schedule and Contract","Items":[{"ClientCaseName":"Test Ln-80","CreatedTime":"2024-04-06T08:51:30+08:00","Gid":"6159272000001066012","IsNew":true}]},{"StagesName":"3. Awaiting Client Records","Items":[{"ClientCaseName":"TestGary TestLiao-0","CreatedTime":"2024-05-05T17:34:11+08:00","Gid":"6159272000002199006","IsNew":true}]},{"StagesName":"4. Record Review","Items":[{"ClientCaseName":"TestFN TestLN-0#5076","CreatedTime":"2024-05-18T15:59:47+08:00","Gid":"6159272000002701086","IsNew":true}]},{"StagesName":"12. MedTeam Forms","Items":[{"ClientCaseName":"TestFN TestLN-0#5076","CreatedTime":"2024-05-18T15:59:47+08:00","Gid":"6159272000002701086","IsNew":true}]},{"StagesName":"13. Mini-DBQs Finalized","Items":[{"ClientCaseName":"TestFN TestLN-80#5093","CreatedTime":"2024-06-05T09:38:00+08:00","Gid":"6159272000003416077","IsNew":true}]},{"StagesName":"23. Completed","Items":[{"ClientCaseName":"TestFn LnN-0","CreatedTime":"2024-04-09T22:05:21+08:00","Gid":"6159272000001184003","IsNew":true}]}]`
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		panic(err)
	}

	body, err := biz.CaseWithoutTasksEmailBody(subject, result)

	lib.DPrintln(body)
	return
	if err != nil {
		panic(err)
	}

	serviceConfig := &biz.MailServiceConfig{
		Name:        "Dev",
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    "glliao@vetbenefitscenter.com",
		Password:    configs.EnvMailGlliaoPWD(),
		FromAddress: "glliao@vetbenefitscenter.com",
	}
	MailMessage := &biz.MailMessage{
		//To: "lialing@foxmail.com",
		To: "liaogling@gmail.com",
		//To:      "ywang@vetbenefitscenter.com",
		Subject: subject,
		Body:    body,
	}

	err = UT.MailUsecase.SendEmail(serviceConfig, MailMessage, "", nil)
	lib.DPrintln("SendEmail:", err)
}

func Test_RemindUsecase_CreateTaskForSubmissionToGoogleDriveFailed(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	err := UT.RemindUsecase.CreateTaskForSubmissionToGoogleDriveFailed(*tCase)
	lib.DPrintln(err)
}

func Test_RemindUsecase_CreateTaskForITFExpirations(t *testing.T) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	//err := UT.RemindUsecase.CreateTaskForITFExpirations()
	//lib.DPrintln(err)
}

func Test_RemindUsecase_HandleCreateTaskForITFExpirations(t *testing.T) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	err := UT.RemindUsecase.HandleCreateTaskForITFExpirations()
	lib.DPrintln(err)
}
