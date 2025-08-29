package tests

import (
	"fmt"
	"os"
	"testing"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

func Test_MailUsecase_SendEmailMS(t *testing.T) {

	// https://learn.microsoft.com/zh-cn/exchange/clients-and-mobile-in-exchange-online/opt-in-exchange-online-endpoint-for-legacy-tls-using-smtp-auth
	// smtp-legacy.office365.com
	serviceConfig := &biz.MailServiceConfig{
		Name:        "Gary",
		Host:        "smtp-legacy.office365.com", // smtp.office365.com
		Port:        587,
		Username:    "gliao@vetbenefitscenter.onmicrosoft.com",
		Password:    "",
		FromAddress: "gliao@vetbenefitscenter.onmicrosoft.com",
	}

	message := &biz.MailMessage{
		To:      "gengling.liao@hotmail.com", // liaogling@gmail.com lialing@foxmail.com
		Subject: "[VBC] Testing GApp Pwd",
		Body:    "Testing body",
	}
	//r, err := os.ReadFile("../templates/remind_email_tpl.html")
	//if err != nil {
	//	panic(err)
	//}
	//message.Body = string(r)

	err := UT.MailUsecase.SendEmail(serviceConfig, message, biz.MailAttach_No, nil)
	lib.DPrintln(err)
}

func Test_MailUsecase_SendEmail(t *testing.T) {
	//serviceConfig := &biz.MailServiceConfig{
	//	Name:        "Mayra Olivares-Iglesias",
	//	Host:        "smtp.gmail.com",
	//	Port:        587,
	//	Username:    "molivares@vetbenefitscenter.com",
	//	Password:    "",
	//	FromAddress: "molivares@vetbenefitscenter.com",
	//}
	serviceConfig := &biz.MailServiceConfig{
		Name:        "Dev",
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    "glliao@vetbenefitscenter.com",
		Password:    configs.EnvMailGlliaoPWD(),
		FromAddress: "engineering@vetbenefitscenter.com",
	}

	serviceConfig = biz.InitAmMailServiceConfig()
	message := &biz.MailMessage{
		To:      "glliao@vetbenefitscenter.com", // liaogling@gmail.com lialing@foxmail.com
		Subject: "[VBC] Testing GApp Pwd",
		Body:    "Testing body",
	}
	//r, err := os.ReadFile("../templates/remind_email_tpl.html")
	//if err != nil {
	//	panic(err)
	//}
	//message.Body = string(r)

	err := UT.MailUsecase.SendEmail(serviceConfig, message, biz.MailAttach_No, nil)
	lib.DPrintln(err)
}

func Test_MailUsecase_MailReplaceDynamicParams(t *testing.T) {
	// 38 12
	tData, err := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 38})
	fmt.Println(*tData.CustomFields.NumberValueByName("current_rating"))
	return

	lib.DPrintln(tData, err)
	str := "aaaaa_{first_name}_scccc\n_{last_name}_\nafdafasfs_{email}_af"
	cc := biz.MailReplaceDynamicParams(str, tData.CustomFields.ToDisplayMaps())
	lib.DPrintln(cc)
}

func Test_MailUsecase_SendEmailWithData(t *testing.T) {

	// 这里可以测试发送发邮件，假如是：MailGenre_GettingStartedEmail， 需要注意修改mail
	data, err := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 5217})
	lib.DPrintln(err)
	tpl, err := UT.TUsecase.Data(biz.Kind_email_tpls, And(Eq{"tpl": biz.MailGenre_GettingStartedEmail}))
	lib.DPrintln(err)
	err, _, _, _, _, _ = UT.MailUsecase.SendEmailWithData(data, tpl, nil)
	fmt.Println(err)
	return
}

func Test_MailUsecase_SendEmailWithData_FeeScheduleCommunication(t *testing.T) {

	// 这里可以测试发送发邮件，假如是：MailGenre_GettingStartedEmail， 需要注意修改mail
	data, err := UT.TUsecase.DataByGid(biz.Kind_client_cases, "d1fbcc1328424c3699057dd71f14e970")
	lib.DPrintln(err)
	rating := 90
	subId := -1

	//rating := 90
	//subId := 100

	data.CustomFields.SetNumberValueByName(biz.FieldName_effective_current_rating, to.Ptr(int32(rating)))
	data.CustomFields.SetNumberValueByName(biz.FieldName_current_rating, to.Ptr(int32(rating)))
	email := data.CustomFields.TextValueByNameBasic(biz.FieldName_email)
	lib.DPrintln(email)
	if true {
		//return
	}
	tpl, err := UT.TUsecase.Data(biz.Kind_email_tpls, And(Eq{"tpl": biz.MailGenre_FeeScheduleCommunication, "sub_id": subId}))
	lib.DPrintln(err)
	var mailTaskInput biz.MailTaskInput
	mailTaskInput.Email = email
	mailTaskInput.Genre = biz.MailGenre_FeeScheduleCommunication
	mailTaskInput.SubId = int32(subId)
	err, _, _, _, _, _ = UT.MailUsecase.SendEmailWithData(data, tpl, &mailTaskInput)
	fmt.Println(err)
	return
}

func Test_MailUsecase_SendEmailWithData_CongratulationsNewRating(t *testing.T) {

	// 这里可以测试发送发邮件，假如是：MailGenre_GettingStartedEmail， 需要注意修改mail
	data, err := UT.TUsecase.DataByGid(biz.Kind_client_cases, "d1fbcc1328424c3699057dd71f14e970")
	lib.DPrintln(err)
	tpl, err := UT.TUsecase.Data(biz.Kind_email_tpls, And(Eq{"tpl": biz.MailGenre_CongratulationsNewRating}))
	lib.DPrintln(err)
	err, _, _, _, _, _ = UT.MailUsecase.SendEmailWithData(data, tpl, nil)
	fmt.Println(err)
	return
}

func Test_MailUsecase_SendEmailWithData_TestTemplate(t *testing.T) {

	bo, err := os.ReadFile("../resource/email_templates/Fee Schedule Communication.html")
	body := string(bo)
	lib.DPrintln(err)
	data, err := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 5217})
	lib.DPrintln(err)
	lib.DPrintln(data)
	tpl, err := UT.TUsecase.Data(biz.Kind_email_tpls, Eq{"tpl": "FeeScheduleCommunication"})
	lib.DPrintln(err)
	lib.DPrintln(tpl)

	tpl.CustomFields.SetTextValueByName("body", &body)
	err, _, _, _, _, _ = UT.MailUsecase.SendEmailWithData(data, tpl, nil)
	fmt.Println(err)
}

func Test_MailUsecase_SendSystemMessage(t *testing.T) {
	err := UT.MailUsecase.SendSystemMessage("sss", "ccc", "lialing@foxmail.com", nil)
	lib.DPrintln(err)
}
