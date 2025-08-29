package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

func TimeFormatToString(time2 time.Time) string {

	ut := USA_TIMEZONE_PT
	loc, _ := time.LoadLocation("America/Los_Angeles")
	a := time2.In(loc).Format("January 2, 2006, 03:04 PM")
	a = fmt.Sprintf("%s (%s)", a, ut)
	return a
}

type DialpadbuzInternalUsecase struct {
	log               *log.Helper
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	DialpadUsecase    *DialpadUsecase
	TaskCreateUsecase *TaskCreateUsecase
	LogUsecase        *LogUsecase
}

func NewDialpadbuzInternalUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	DialpadUsecase *DialpadUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	LogUsecase *LogUsecase,
) *DialpadbuzInternalUsecase {
	uc := &DialpadbuzInternalUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		DialpadUsecase:    DialpadUsecase,
		TaskCreateUsecase: TaskCreateUsecase,
		LogUsecase:        LogUsecase,
	}

	return uc
}

func (c *DialpadbuzInternalUsecase) ClientSMSReplyAlertForAdmin(clientName string, phone string, time time.Time, actionType string) error {

	text := c.ClientSMSReplyAlertText(clientName, phone, time, actionType)
	//dialpadUserId := "5243348849786880" // +13109719619

	lib.DPrintln("text:", text)
	//lib.DPrintln("dialpadUserId:", dialpadUserId)

	YNPhone := "+18056604465"
	EDPhone := "+16198005543"

	c.LogUsecase.SaveLog(0, "DialpadbuzInternal:ClientSMSReplyAlertForAdmin", map[string]interface{}{
		"ReceivePhone":        YNPhone,
		"ReceiveText":         text,
		"SenderDialpadUserid": SysDialpadUserid,
		"SenderUserFullName":  SysDialpadName,
	})

	c.LogUsecase.SaveLog(0, "DialpadbuzInternal:ClientSMSReplyAlertForAdmin", map[string]interface{}{
		"ReceivePhone":        EDPhone,
		"ReceiveText":         text,
		"SenderDialpadUserid": SysDialpadUserid,
		"SenderUserFullName":  SysDialpadName,
	})

	if configs.Enable_SMS_New_Version_Debug {
		return nil
	}

	err := c.DialpadUsecase.SendSmsNoFilter(YNPhone, text, SysDialpadUserid, 0, "ClientSMSReplyAlertForAdmin")
	if err != nil {
		c.log.Error(err)
	}
	err = c.DialpadUsecase.SendSmsNoFilter(EDPhone, text, SysDialpadName, 0, "ClientSMSReplyAlertForAdmin")
	if err != nil {
		c.log.Error(err)
	}
	return nil
}

func (c *DialpadbuzInternalUsecase) ClientSMSReplyAlertText(clientName string, phone string, time time.Time, actionType string) string {

	if actionType == SMS_Text_action_stop {
		actionType = strings.ToUpper(actionType) + " (Reject SMS)"
	} else if actionType == SMS_Text_action_unstop {
		actionType = strings.ToUpper(actionType) + " (Accept SMS)"
	} else {
		actionType = strings.ToUpper(actionType)
	}
	str := "The client SMS reply action:\n\nClient Case Name: " + clientName + "\nPhone: " + phone + "\nDate & Time: " + TimeFormatToString(time) + "\nAction Type: " + actionType + "\n\nPlease be informed."
	return str
}

// HandleClientSMSReplyAlertTextForEmailAdmin 通知管理员
func (c *DialpadbuzInternalUsecase) HandleClientSMSReplyAlertTextForEmailAdmin(clientName string, phone string, time3 time.Time, actionType string) error {

	//	str := c.ClientSMSReplyAlertText(clientName, phone, time3, actionType)
	//	subject := "VBC: Client SMS Reply Action Notification"
	//	content := strings.ReplaceAll(str, "\n", "<br />")
	//
	//	body := MailAutomationBodyHeader(subject) + `
	//<div style="line-height:10px;">&nbsp;</div>
	//<div>` + content + `</div>
	//` + MailAutomationBodyBottom()

	email := "ywang@vetbenefitscenter.com;ebunting@vetbenefitscenter.com"

	return c.CreateTaskClientSMSReplyAlertForEmail(clientName, phone, time3, actionType, email)
	//nextAt := time.Now().Unix()
	//return c.TaskCreateUsecase.CreateCustomTaskMail(0, &MailMessage{
	//	To:      email,
	//	Subject: subject,
	//	Body:    body,
	//}, nextAt)
}

func (c *DialpadbuzInternalUsecase) CreateTaskClientSMSReplyAlertForEmail(clientName string, phone string, time3 time.Time, actionType string, email string) error {

	str := c.ClientSMSReplyAlertText(clientName, phone, time3, actionType)
	subject := "VBC: Client SMS Reply Action Notification: " + strings.ToUpper(actionType)
	content := strings.ReplaceAll(str, "\n", "<br />")

	body := MailAutomationBodyHeader(subject) + `
<div style="line-height:10px;">&nbsp;</div>
<div style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">` + content + `</div>
` + MailAutomationBodyBottom()

	nextAt := time.Now().Unix()

	//if lib.Enable_SMS_New_Version_Debug {
	//	email = "liaogling@gmail.com;lialing@foxmail.com"
	//}

	return c.TaskCreateUsecase.CreateCustomTaskMail(0, &MailMessage{
		To:      email,
		Subject: subject,
		Body:    body,
	}, nextAt)
}
