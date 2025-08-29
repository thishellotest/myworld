package tests

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"vbc/lib/gomail"
)

type ServiceConfig struct {
	MailName        string `json:"mail_name"`         // 发件人名称
	MailHost        string `json:"mail_host"`         // 发件服务器地址或ip*
	MailPort        int    `json:"mail_port"`         // 发件服务端口*
	MailUsername    string `json:"mail_username"`     // smtp登录的用户名
	MailPassword    string `json:"mail_password"`     // smtp登录的密码
	MailFromAddress string `json:"mail_from_address"` // 显示的发件箱
	MailReplayTo    string `json:"mail_replay_to"`    // 回复邮箱
	MailEncryptType string `json:"mail_encrypt_type"` // 加密方式：只能允许 ssl / tls / ""
	MailVerifyType  string `json:"mail_verify_type"`  // 如：LOGIN
	MailSendType    int32  `json:"mail_send_type"`    // 0默认 1营销类
}

type MessageVar struct {
	To              string `json:"to"` // gomail 暂不支持多个收件人；多个收件人使用 ; 分隔 11@qq.com;22@qq.com
	FromName        string `json:"from_name"`
	ReplyTo         string `json:"reply_to"`     // 回复邮件地址，不设置时默认使用注册的配置
	CallbackUrl     string `json:"callback_url"` // 任务回调的地址
	CallbackPayload string `json:"callback_payload"`

	Mailtype string `json:"mailtype"` // 邮件类型  html 或其它
	Subject  string `json:"subject"`  // 主题
	Body     string `json:"body"`     // 发送内容
}

/*
Gmail SMTP server address: smtp.gmail.com
Gmail SMTP name: Your full name
Gmail SMTP username: Your full Gmail address (e.g. you@gmail.com)
Gmail SMTP password: The password that you use to log in to Gmail
Gmail SMTP port (TLS): 587
Gmail SMTP port (SSL): 465
*/

func Test_send_mail_google(t *testing.T) {
	var d *gomail.Dialer
	serviceConfig := &ServiceConfig{
		MailName:        "Liao Gary",
		MailHost:        "smtp.gmail.com",
		MailPort:        587,
		MailUsername:    "liaogling@gmail.com",
		MailPassword:    "",
		MailFromAddress: "liaogling@gmail.com",
	}

	messageVar := &MessageVar{
		To:      "18891706@qq.com",
		Subject: "Test gomail",
		Body:    "test body",
	}

	d = gomail.NewDialer(serviceConfig.MailHost, serviceConfig.MailPort, serviceConfig.MailUsername, serviceConfig.MailPassword)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()

	// 使用带名称方式
	fromDisplay := serviceConfig.MailFromAddress
	fromDisplay = serviceConfig.MailName

	m.SetAddressHeader("From", serviceConfig.MailFromAddress, fromDisplay)
	emails := strings.Split(messageVar.To, ";")
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", messageVar.Subject)
	if messageVar.Mailtype == "html" {
		m.SetBody("text/html", messageVar.Body)
	} else {
		m.SetBody("text/plain", messageVar.Body)
	}
	if messageVar.ReplyTo != "" {
		m.SetHeader("Reply-To", messageVar.ReplyTo)
	} else if serviceConfig.MailReplayTo != "" {
		m.SetHeader("Reply-To", serviceConfig.MailReplayTo)
	}

	err := d.DialAndSend(m)
	fmt.Println(err)
}
