package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type RemindUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	TUsecase          *TUsecase
	TaskCreateUsecase *TaskCreateUsecase
	ClientCaseUsecase *ClientCaseUsecase
	UserUsecase       *UserUsecase
	LogUsecase        *LogUsecase
}

func NewRemindUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	UserUsecase *UserUsecase,
	LogUsecase *LogUsecase) *RemindUsecase {
	uc := &RemindUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		TUsecase:          TUsecase,
		TaskCreateUsecase: TaskCreateUsecase,
		ClientCaseUsecase: ClientCaseUsecase,
		UserUsecase:       UserUsecase,
		LogUsecase:        LogUsecase,
	}

	return uc
}

type FollowingUpSignMedicalTeamFormsEmailVo struct {
	Email   string
	Subject string
	Body    string
}

func (c *RemindUsecase) FollowingUpSignMedicalTeamFormsEmailBody(tCase *TData, tUser *TData, ContractSentOn *time.Time) (vo *FollowingUpSignMedicalTeamFormsEmailVo, err error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	subject := "[VBC] You Have A New Reminder (Subject: Following Up with Clients to Sign Medical Team Forms)"

	fullName := tUser.CustomFields.TextValueByNameBasic("full_name")
	email := tUser.CustomFields.TextValueByNameBasic("email")
	if email == "" {
		return nil, errors.New("user - email is nil")
	}

	gid := tCase.CustomFields.TextValueByNameBasic("gid")

	ContractSentOnStr := ""
	if ContractSentOn != nil {
		ContractSentOnStr = ContractSentOn.In(configs.VBCDefaultLocation).Format("Mon, 02 Jan 2006 15:04 PM")
	}

	items := lib.TypeList{{
		"Label": "Subject",
		"Value": "Following Up with Clients to Sign Medical Team Forms",
	}, {
		"Label": "Client Case Name",
		"Value": tCase.CustomFields.TextValueByNameBasic("deal_name"),
	}, {
		"Label": "Contract Sent On",
		"Value": ContractSentOnStr,
	}, {
		"Label": "Client Case URL",
		"Value": "<a target=\"_blank\" href=\"https://base.vetbenefitscenter.com/tab/cases/" + gid + "\">https://base.vetbenefitscenter.com/tab/cases/" + gid + "</a>",
	}}

	body, err := FollowingUpSignMedicalTeamFormsEmailBody(subject, fullName, items)
	if err != nil {
		return nil, err
	}
	vo = &FollowingUpSignMedicalTeamFormsEmailVo{
		Email:   email,
		Subject: subject,
		Body:    body,
	}
	return vo, nil
}

func (c *RemindUsecase) FollowingUpUploadedDocumentEmailBody(tCase *TData,
	email string,
	updateFiles []*ReminderClientUpdateFilesEventVoItem,
) (vo *MailMessageVo, err error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	//if tUser == nil {
	//	return nil, errors.New("tUser is nil")
	//}

	subject := "[VBC] You Have A New Reminder (Subject: Following up with the client who has updated a new document)"

	//fullName := tUser.CustomFields.TextValueByNameBasic("full_name")
	fullName := "All"
	//email := tUser.CustomFields.TextValueByNameBasic("email")
	//if email == "" {
	//	return nil, errors.New("user - email is nil")
	//}

	//gid := tCase.CustomFields.TextValueByNameBasic("gid")

	//LatestUpdatesFiles := ""
	//for _, v := range updateFiles {
	//	LatestUpdatesFiles += fmt.Sprintf("<tr><td style=\"font-size:14px;line-height:170%%\"><a href=\"https://veteranbenefitscenter.app.box.com/file/%s\" target=\"_blank\">%s</a></td><td style=\"font-size:14px;line-height:170%%\"><a href=\"https://veteranbenefitscenter.app.box.com/file/%s\" target=\"_blank\">%s</a></td></tr>", v.SourceBoxResId, v.SourceBoxPath, v.BoxResId, v.BoxPath)
	//}

	//LatestUpdates := fmt.Sprintf("<table><tr><td style=\"font-size:14px;line-height:170%%\">Original Source File</td><td style=\"font-size:14px;line-height:170%%\">Copied File Destination</td></tr>%s</table>", LatestUpdatesFiles)

	//items := lib.TypeList{{
	//	"Label": "Client Case Name",
	//	"Value": tCase.CustomFields.TextValueByNameBasic("deal_name"),
	//}, {
	//	"Label": "Latest Updates",
	//	"Value": LatestUpdates,
	//}}

	body, err := FollowingUpUploadedDocumentEmailBody(subject, fullName, tCase.CustomFields.TextValueByNameBasic("deal_name"), updateFiles)
	if err != nil {
		return nil, err
	}
	vo = &MailMessageVo{
		Email:   email,
		Subject: subject,
		Body:    body,
	}
	return vo, nil
}

// CreateUnfinishedFeeContract Creating Client does not finish fee contract remind via VS
func (c *RemindUsecase) CreateUnfinishedFeeContract(clientId int32) error {

	if configs.IsProd() {
		return nil
	}
	// 提醒流程暂不支持
	return nil

	client, err := c.TUsecase.DataById(Kind_client_cases, clientId)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("client is nil.")
	}
	clientData := client.CustomFields

	user, err := c.TUsecase.Data(Kind_users, Eq{"user_gid": clientData.TextValueByNameBasic("user_gid")})
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user is nil")
	}
	userData := user.CustomFields
	//subject := fmt.Sprintf("Please remind client(%s, %s) to sign the contract",
	//	clientData.TextValueByNameBasic("first_name"),
	//	clientData.TextValueByNameBasic("last_name"),
	//)
	//	body := MailAutomationBodyHeader(subject) + `<div>Hi ` + userData.TextValueByNameBasic("first_name") + `:</div>
	//<div style="line-height:10px;">&nbsp;</div>
	//<div>Please remind client here: <a href="https://www.veteranbenefitscenter.com/" target="_blank">notification</div>
	//` + MailAutomationBodyBottom()

	email := userData.TextValueByNameBasic("email")
	if email == "" {
		return errors.New("user email is empty.")
	}
	nextAt := time.Now().Unix() + 48*3600
	if !configs.IsProd() {
		nextAt = time.Now().Unix() + 10
	}
	return c.TaskCreateUsecase.CreateTaskMail(clientId, MailGenre_NotifyVsRemindClient, 0, nil, nextAt, "", "")
	//return c.TaskCreateUsecase.CreateCustomTaskMail(clientId, &MailMessage{
	//	To:      email,
	//	Subject: subject,
	//	Body:    body,
	//}, nextAt)
}

// CreateUnfinishedIntakeForm Creating Client does not finish intake form remind via VS
func (c *RemindUsecase) CreateUnfinishedIntakeForm(clientId int32) error {
	if configs.IsProd() {
		return nil
	}
	return nil
	client, err := c.TUsecase.DataById(Kind_client_cases, clientId)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("client is nil.")
	}
	clientData := client.CustomFields

	user, err := c.TUsecase.Data(Kind_users, Eq{"asana_user_gid": clientData.TextValueByNameBasic("assignee_gid")})
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user is nil")
	}
	userData := user.CustomFields
	subject := fmt.Sprintf("Please remind client(%s, %s) to complete intake form",
		clientData.TextValueByNameBasic("first_name"),
		clientData.TextValueByNameBasic("last_name"),
	)
	body := MailAutomationBodyHeader(subject) + `<div>Hi ` + userData.TextValueByNameBasic("first_name") + `:</div>
<div style="line-height:10px;">&nbsp;</div>
<div>Please remind client here: <a href="https://www.veteranbenefitscenter.com/complete-intake-form" target="_blank">notification</div>
` + MailAutomationBodyBottom()
	email := userData.TextValueByNameBasic("email")
	if email == "" {
		return errors.New("user email is empty.")
	}
	nextAt := time.Now().Unix() + 14*24*3600
	return c.TaskCreateUsecase.CreateCustomTaskMail(clientId, &MailMessage{
		To:      email,
		Subject: subject,
		Body:    body,
	}, nextAt)
}

func (c *RemindUsecase) SubmissionToGoogleDriveFailed(dealName string, caseGid string, toEmail string, cc []string) MailMessage {

	subject := "VBC: Medical Team – Issue with Private Exams Submission to Google Drive"
	content := `The Client Case Information:<br />
<br />
Client Case Name: <a href="` + configs.Domain + "/tab/cases/" + caseGid + `" target="_blank">` + dealName + `</a><br />
Date & Time: ` + TimeFormatToString(time.Now()) + `<br />
Submission Status: Incomplete<br />
Possible Cause: <br />
1. The client's folder is not a BOX Folder.<br />
2. The client's Medical Team - Forms information has not been successfully received yet.<br />
<br />
Please be informed.`

	body := MailAutomationBodyHeader(subject) + `
<div style="line-height:10px;">&nbsp;</div>
<div style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">` + content + `</div>
` + MailAutomationBodyBottom()
	return MailMessage{
		To:      toEmail,
		Subject: subject,
		Body:    body,
		Cc:      cc,
	}
}

func (c *RemindUsecase) CreateTaskForSubmissionToGoogleDriveFailed(tCase TData) error {

	email := configs.SubmissionToGoogleDriveFailedNotifyEmail
	if configs.IsDev() {
		email = "glliao@vetbenefitscenter.com"
	}
	cc := []string{"info@vetbenefitscenter.com", "engineering@vetbenefitscenter.com"}
	//var cc []string
	nextAt := time.Now().Unix()
	mailMessage := c.SubmissionToGoogleDriveFailed(tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name), tCase.Gid(), email, cc)

	return c.TaskCreateUsecase.CreateCustomTaskMail(0, &mailMessage, nextAt)
}

type ITFExpirationsEmailItem struct {
	Days  int
	TCase TData
}
type ITFExpirationsEmailVo struct {
	Overdues  []ITFExpirationsEmailItem //  (Days < 0)
	DueSoons  []ITFExpirationsEmailItem // (0–30 Days)
	MidTerms  []ITFExpirationsEmailItem //(31–60 Days)
	LongTerms []ITFExpirationsEmailItem // (61–90 Days)
}

func (c *RemindUsecase) ITFExpirationsEmail(userGid string) (iTFExpirationsEmailVo ITFExpirationsEmailVo, err error) {

	var cases TDataList
	if userGid == "" {
		cases, err = c.ClientCaseUsecase.ItfCases()
		if err != nil {
			return iTFExpirationsEmailVo, err
		}
	} else {
		cases, err = c.ClientCaseUsecase.ItfCasesByUserGid(userGid)
	}
	now := time.Now()
	now = now.In(configs.VBCDefaultLocation)
	now, _ = time.ParseInLocation(time.DateOnly, now.Format(time.DateOnly), configs.VBCDefaultLocation)

	for _, v := range cases {
		itf := v.CustomFields.TextValueByNameBasic(FieldName_itf_expiration)
		tTime, _ := time.ParseInLocation(time.DateOnly, itf, configs.VBCDefaultLocation)
		diff := tTime.Sub(now)
		days := int(math.Ceil(diff.Hours() / 24))
		iTFExpirationsEmailItem := ITFExpirationsEmailItem{
			Days:  days,
			TCase: v,
		}
		if days < 0 {
			iTFExpirationsEmailVo.Overdues = append(iTFExpirationsEmailVo.Overdues, iTFExpirationsEmailItem)
		} else if days <= 30 {
			iTFExpirationsEmailVo.DueSoons = append(iTFExpirationsEmailVo.DueSoons, iTFExpirationsEmailItem)
		} else if days <= 60 {
			iTFExpirationsEmailVo.MidTerms = append(iTFExpirationsEmailVo.MidTerms, iTFExpirationsEmailItem)
		} else {
			iTFExpirationsEmailVo.LongTerms = append(iTFExpirationsEmailVo.LongTerms, iTFExpirationsEmailItem)
		}
	}
	return
}

func (c *RemindUsecase) DoITFExpirationsEmail(userGid string) (string, error) {

	vo, err := c.ITFExpirationsEmail(userGid)
	if err != nil {
		return "", err
	}
	text := "<h2 style=\"margin-bottom: 10px;\">Daily Summary of Cases by ITF Expiration</h2>"
	if len(vo.Overdues) > 0 {
		text += `<div style="background-color: #ffebee; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
      <strong style="display: block; margin-bottom: 5px;">🟥 Overdue (Days &lt; 0)</strong>
      <table width="100%" cellpadding="4" cellspacing="0" border="0">`

		for _, v := range vo.Overdues {
			text += `<tr style="padding:0;margin:0">
          <td style="font-size:13px;padding:0;margin:0"><a target="_blank" href="https://base.vetbenefitscenter.com/tab/cases/` + v.TCase.Gid() + `" style="color:#333;font-size:13px;text-decoration: none;"><span style="width: 250px; display: inline-block;">` + v.TCase.CustomFields.TextValueByNameBasic(FieldName_deal_name) + `</span></a><span style="width: 300px; display: inline-block;">` + v.TCase.CustomFields.DisplayValueByNameBasic(FieldName_stages) + `</span> (ITF Expired ` + InterfaceToString(v.Days) + ` days ago)</td>
        </tr>`
		}

		text += `</table>
    </div>`
	}
	if len(vo.DueSoons) > 0 {

		text += `<div style="background-color: #fff8e1; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
      <strong style="display: block; margin-bottom: 5px;">🟧 Due Soon (0–30 Days)</strong>
      <table width="100%" cellpadding="4" cellspacing="0" border="0">`

		for _, v := range vo.DueSoons {
			text += `<tr style="padding:0;margin:0">
          <td style="font-size:13px;padding:0;margin:0"><a target="_blank" href="https://base.vetbenefitscenter.com/tab/cases/` + v.TCase.Gid() + `" style="color:#333;font-size:13px;text-decoration: none;"><span style="width: 250px; display: inline-block;">` + v.TCase.CustomFields.TextValueByNameBasic(FieldName_deal_name) + `</span></a> <span style="width: 300px; display: inline-block;">` + v.TCase.CustomFields.DisplayValueByNameBasic(FieldName_stages) + `</span> (ITF Expired ` + InterfaceToString(v.Days) + ` days ago)</td>
        </tr>`
		}

		text += `</table>
    </div>`
	}
	if len(vo.MidTerms) > 0 {
		text += `<div style="background-color: #e3f2fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
      <strong style="display: block; margin-bottom: 5px;">🟦 Mid-Term (31–60 Days)</strong>
      <table width="100%" cellpadding="4" cellspacing="0" border="0">`

		for _, v := range vo.MidTerms {
			text += `<tr style="padding:0;margin:0">
          <td style="font-size:13px;padding:0;margin:0"><a target="_blank" href="https://base.vetbenefitscenter.com/tab/cases/` + v.TCase.Gid() + `" style="color:#333;font-size:13px;text-decoration: none;"><span style="width: 250px; display: inline-block;">` + v.TCase.CustomFields.TextValueByNameBasic(FieldName_deal_name) + `</span></a> <span style="width: 300px; display: inline-block;">` + v.TCase.CustomFields.DisplayValueByNameBasic(FieldName_stages) + `</span> (ITF Expired ` + InterfaceToString(v.Days) + ` days ago)</td>
        </tr>`
		}

		text += `</table>
    </div>`
	}
	if len(vo.LongTerms) > 0 {
		text += `<div style="background-color: #e8f5e9; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
      <strong style="display: block; margin-bottom: 5px;">🟩 Long-Term (61–90 Days)</strong>
      <table width="100%" cellpadding="4" cellspacing="0" border="0">`

		for _, v := range vo.LongTerms {
			text += `<tr style="padding:0;margin:0">
          <td style="font-size:13px;padding:0;margin:0"><a target="_blank" href="https://base.vetbenefitscenter.com/tab/cases/` + v.TCase.Gid() + `" style="color:#333;font-size:13px;text-decoration: none;"><span style="width: 250px; display: inline-block;">` + v.TCase.CustomFields.TextValueByNameBasic(FieldName_deal_name) + `</span></a> <span style="width: 300px; display: inline-block;">` + v.TCase.CustomFields.DisplayValueByNameBasic(FieldName_stages) + `</span> (ITF Expired ` + InterfaceToString(v.Days) + ` days ago)</td>
        </tr>`
		}

		text += `</table>
    </div>`
	}

	return text, nil
}
func (c *RemindUsecase) HandleCreateTaskForITFExpirations() error {

	key := "HandleCreateTaskForITFExpirations"
	str, _ := c.CommonUsecase.RedisClient().Get(context.TODO(), key).Result()
	if str == "" {
		email := configs.YnEmail + ";" + configs.EdEmail
		err := c.CreateTaskForITFExpirations("", email, nil)
		if err != nil {
			return err
		}
		er := c.CommonUsecase.RedisClient().Set(context.TODO(), key, "1", time.Hour*24).Err()
		if er != nil {
			c.log.Error(er)
		}
	}
	er := c.HandleCreateTaskForITFExpirationsForVSTeam()
	if er != nil {
		c.log.Error(er)
	}
	return nil
}

func (c *RemindUsecase) HandleCreateTaskForITFExpirationsForVSTeam() error {

	vsTeamUsers, err := c.UserUsecase.VSTeamUsers()
	if err != nil {
		c.log.Error(err)
	}
	for _, v := range vsTeamUsers {
		key := "HandleCreateTaskForITFExpirations:" + InterfaceToString(v.Id())
		str, _ := c.CommonUsecase.RedisClient().Get(context.TODO(), key).Result()
		if str == "" {
			email := v.CustomFields.TextValueByNameBasic(UserFieldName_email)
			err = c.CreateTaskForITFExpirations(v.Gid(), email, nil)
			if err != nil {
				return err
			}
			er := c.CommonUsecase.RedisClient().Set(context.TODO(), key, "1", time.Hour*24).Err()
			if er != nil {
				c.log.Error(er)
			}
		}
	}
	return nil
}

// CreateTaskForITFExpirations email := lib.YnEmail + ";" + lib.EdEmail
func (c *RemindUsecase) CreateTaskForITFExpirations(userGid string, email string, cc []string) error {

	c.LogUsecase.SaveLog(0, "CreateTaskForITFExpirations", map[string]interface{}{
		"userGid": userGid,
		"email":   email,
		"cc":      cc,
	})

	subject := "VBC: Upcoming ITF Expirations Overview"
	//if lib.IsDev() {
	//	subject += ":" + email + ":" + userGid
	//	email = "liaogling@gmail.com"
	//}
	if configs.IsProd() {
		if userGid != "" {
			cc = append(cc, "engineering@vetbenefitscenter.com")
		}
	}
	nextAt := time.Now().Unix()
	content, err := c.DoITFExpirationsEmail(userGid)
	if err != nil {
		return err
	}
	body := MailAutomationBodyHeader(subject) + `
<div style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;max-width: 800px;">` + content + `</div>
` + MailAutomationBodyBottom()
	mailMessage := MailMessage{
		To:      email,
		Subject: subject,
		Body:    body,
		Cc:      cc,
	}

	return c.TaskCreateUsecase.CreateCustomTaskMail(0, &mailMessage, nextAt)
}

func MailAutomationBodyHeader(subject string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <title>` + subject + `</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
<body>`
}

func MailAutomationBodyBottom() string {
	return `<br/>
--<br/>
<div style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">
Automation Workflow<br/>
Veteran Benefits Center LLC<br/>
</div>
</body>
</html>`
}
