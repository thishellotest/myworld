package biz

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	"vbc/lib/builder"
	"vbc/lib/gomail"
)

type MailServiceConfig struct {
	Name        string // 发件人名称
	Host        string // 发件服务器地址或ip*
	Port        int    // 发件服务端口*
	Username    string // smtp登录的用户名
	Password    string // smtp登录的密码
	FromAddress string // 显示的发件箱
	ReplayTo    string // 回复邮箱
}

type MailMessageVo struct {
	Email   string
	Subject string
	Body    string
}

type MailMessage struct {
	To       string   // gomail 暂不支持多个收件人；多个收件人使用 ; 分隔 11@qq.com;22@qq.com
	MailType string   // 邮件类型  html 或其它
	Subject  string   // 主题
	Body     string   // 发送内容
	Cc       []string // 需要抄送的对象
}

func (c *MailMessage) GetSubject() string {
	if !configs.IsProd() {
		//return fmt.Sprintf("[%s] %s", "Development", c.Subject)
	}
	return c.Subject
}

// MailDynamicParams 获取动态参数
func MailDynamicParams(tplText string) []string {
	pattern := `{.*?}`
	re := regexp.MustCompile(pattern)
	return re.FindAllString(tplText, -1)
}

// MailReplaceDynamicParams Replace dynamic params
func MailReplaceDynamicParams(tplText string, tData lib.TypeMap) string {

	dynamicParams := MailDynamicParams(tplText)
	for _, v := range dynamicParams {
		newStr := ""
		if tData != nil {
			fieldName := strings.ReplaceAll(v, "{", "")
			fieldName = strings.ReplaceAll(fieldName, "}", "")
			newStr = tData.GetString(fieldName)
		}
		tplText = strings.ReplaceAll(tplText, v, newStr)
	}
	return tplText
}

type MailUsecase struct {
	TUsecase                 *TUsecase
	conf                     *conf.Data
	AccessControlWorkUsecase *AccessControlWorkUsecase
	DataComboUsecase         *DataComboUsecase
	FeeUsecase               *FeeUsecase
	ClientEnvelopeUsecase    *ClientEnvelopeUsecase
	MailFeeContentUsecase    *MailFeeContentUsecase
	log                      *log.Helper
	UserUsecase              *UserUsecase
	StatementUsecase         *StatementUsecase
	SendVa2122aUsecase       *SendVa2122aUsecase
	AttorneyUsecase          *AttorneyUsecase
}

func NewMailUsecase(TUsecase *TUsecase, conf *conf.Data, AccessControlWorkUsecase *AccessControlWorkUsecase, DataComboUsecase *DataComboUsecase,
	FeeUsecase *FeeUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	MailFeeContentUsecase *MailFeeContentUsecase,
	logger log.Logger,
	UserUsecase *UserUsecase,
	StatementUsecase *StatementUsecase,
	SendVa2122aUsecase *SendVa2122aUsecase,
	AttorneyUsecase *AttorneyUsecase) *MailUsecase {

	return &MailUsecase{
		TUsecase:                 TUsecase,
		conf:                     conf,
		AccessControlWorkUsecase: AccessControlWorkUsecase,
		DataComboUsecase:         DataComboUsecase,
		FeeUsecase:               FeeUsecase,
		ClientEnvelopeUsecase:    ClientEnvelopeUsecase,
		MailFeeContentUsecase:    MailFeeContentUsecase,
		log:                      log.NewHelper(logger),
		UserUsecase:              UserUsecase,
		StatementUsecase:         StatementUsecase,
		SendVa2122aUsecase:       SendVa2122aUsecase,
		AttorneyUsecase:          AttorneyUsecase,
	}
}

const (
	MailAttach_No                        = ""
	MailAttach_HowtoGuide                = "HowtoGuide"
	MailAttach_VBCClothingAllowanceGuide = "VBCClothingAllowanceGuide"
)

type MailAttachmentInputs []MailAttachmentInput

type MailAttachmentInput struct {
	Name     string
	Reader   io.Reader
	Settings []gomail.FileSetting
}

func (c *MailUsecase) SendEmail(mailServiceConfig *MailServiceConfig, message *MailMessage, MailAttach string, mailAttachmentInputs MailAttachmentInputs) error {

	var d *gomail.Dialer

	d = gomail.NewDialer(mailServiceConfig.Host, mailServiceConfig.Port, mailServiceConfig.Username, mailServiceConfig.Password)

	if configs.IsDev() {
		d = gomail.NewDialerWithProxy(mailServiceConfig.Host, mailServiceConfig.Port,
			mailServiceConfig.Username, mailServiceConfig.Password,
			gomail.Proxy{
				Address: "127.0.0.1:7890",
			})
		// 使用本地代socket5代码才能发送
	}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailServiceConfig.FromAddress, mailServiceConfig.Name)
	emails := strings.Split(message.To, ";")
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", message.GetSubject())
	if message.MailType == "text" {
		m.SetBody("text/plain", message.Body)
	} else {
		m.SetBody("text/html", message.Body)
	}
	if mailServiceConfig.ReplayTo != "" {
		m.SetHeader("Reply-To", mailServiceConfig.ReplayTo)
	}
	if len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}

	if MailAttach == MailAttach_HowtoGuide {
		m.Attach(c.conf.ResourcePath + "/How-to-Guide v7.3.pdf")
	} else if MailAttach == MailAttach_VBCClothingAllowanceGuide {
		//m.Attach(c.conf.ResourcePath + "/VBC Clothing Allowance Guide.pdf")
	}

	for _, v := range mailAttachmentInputs {
		m.AttachReader(v.Name, v.Reader, v.Settings...)
	}

	if false && configs.IsWorkflowDebug(emails[0]) { // 开启测试不真实发送
		return nil
	} else {
		return d.DialAndSend(m)
	}
}

func InitMailServiceConfig() *MailServiceConfig {

	serviceConfig := &MailServiceConfig{
		Name:        "Dev",
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    "glliao@vetbenefitscenter.com",
		Password:    configs.EnvMailGlliaoPWD(),
		FromAddress: "glliao@vetbenefitscenter.com",
	}
	if !configs.IsDev() {
		serviceConfig.Name = "Yannan Wang"
		serviceConfig.Username = "ywang@vetbenefitscenter.com"
		serviceConfig.Password = configs.EnvMailYwangPWD()
		serviceConfig.FromAddress = "ywang@vetbenefitscenter.com"
	}
	return serviceConfig

}

const (
	AmSenderEamil = "team@augustusmiles.com"
)

func InitAmMailServiceConfig() *MailServiceConfig {

	serviceConfig := &MailServiceConfig{
		Name:        "Augustus Miles Team",
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    AmSenderEamil,
		Password:    configs.EnvMailTeamAgsPWD(),
		FromAddress: AmSenderEamil,
	}
	return serviceConfig
}

// VerifyUserEmailConfig 验证用户的邮件配置是否正确
func (c *MailUsecase) VerifyUserEmailConfig(tUser *TData) bool {
	if tUser == nil {
		return false
	}
	if tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailSender) == "" ||
		tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailPassword) == "" {
		return false
	}
	return true
}

func (c *MailUsecase) SendEmailWithData(clientCase *TData, tpl *TData, mailTaskInput *MailTaskInput) (err error, subject string, body string, toEmail string, senderEmail string, senderName string) {

	var cc []string
	if tpl == nil {
		return errors.New("Email Tpl is empty."), "", "", "", "", ""
	}
	if clientCase == nil {
		return errors.New("ClientCase is empty."), "", "", "", "", ""
	}

	tClientCaseFields := clientCase.CustomFields
	clientGid := clientCase.CustomFields.TextValueByNameBasic("client_gid")

	_, tContactFields, err := c.DataComboUsecase.Client(clientGid)
	if err != nil {
		return err, "", "", "", "", ""
	}
	if tContactFields == nil {
		return errors.New("tContactFields is nil."), "", "", "", "", ""
	}

	email := tContactFields.TextValueByNameBasic("email")

	if tpl.CustomFields.TextValueByName("tpl") != nil &&
		(*tpl.CustomFields.TextValueByName("tpl") == MailGenre_FeeScheduleCommunication ||
			*tpl.CustomFields.TextValueByName("tpl") == MailGenre_StartYourVADisabilityClaimRepresentation) {
		email = mailTaskInput.Email
	}

	if email == "" {
		lib.DPrintln("SendEmailWithData:Email format is wrong", email)
		return errors.New("SendEmailWithData:Email format is wrong."), "", "", "", "", ""
	}

	if tpl.CustomFields.TextValueByName("subject") != nil {
		subject = *(tpl.CustomFields.TextValueByName("subject"))
	}
	if tpl.CustomFields.TextValueByName("body") != nil {
		body = *(tpl.CustomFields.TextValueByName("body"))
	}
	contactMap := tContactFields.ToDisplayMaps()
	caseMap := tClientCaseFields.ToDisplayMaps()
	clientMap := lib.TypeMapMerge(contactMap, caseMap)

	if mailTaskInput != nil {
		for k, v := range mailTaskInput.DynamicParams {
			clientMap.Set(k, InterfaceToString(v))
		}
	}

	// 此处因为字段相同冲突出bug了，在此修复
	clientMap.Set("email", email)
	clientMap.Set("phone", contactMap.GetString("phone"))
	clientMap.Set("ssn", contactMap.GetString("ssn"))
	clientMap.Set("dob", contactMap.GetString("dob"))
	clientMap.Set("state", contactMap.GetString("state"))
	clientMap.Set("city", contactMap.GetString("city"))
	clientMap.Set("address", contactMap.GetString("address"))
	clientMap.Set("zip_code", contactMap.GetString("zip_code"))

	//userGid := tClientCaseFields.TextValueByNameBasic("user_gid")

	var tUser *TData
	useAmSendConfig := false
	if tpl != nil {
		tplValue := tpl.CustomFields.TextValueByNameBasic("tpl")

		if tplValue == MailGenre_CongratulationsNewRating {
			tUser, _ = c.TUsecase.Data(Kind_users, builder.Eq{"gid": config_vbc.User_Edward_gid})
		} else if tplValue == MailGenre_UpcomingContactInformation {
			if mailTaskInput != nil {
				leadVSChangeLogValue := mailTaskInput.DynamicParams.GetString("LeadVSChangeLog")
				leadVSChangeLogVo := lib.StringToTDef(leadVSChangeLogValue, LeadVSChangeLogVo{})
				var newUser *TData
				newUser, _ = c.UserUsecase.GetUserByLeadVS(clientCase)
				if newUser == nil {
					return errors.New("LeadVS User is nil: " + InterfaceToString(clientCase.Id())), "", "", "", "", ""
				}

				clientMap.Set("upcoming_contact_vs:full_name", newUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname))
				clientMap.Set("upcoming_contact_vs:title", newUser.CustomFields.TextValueByNameBasic(UserFieldName_title))
				clientMap.Set("upcoming_contact_vs:mobile", newUser.CustomFields.TextValueByNameBasic(UserFieldName_mobile))
				clientMap.Set("upcoming_contact_vs:email", newUser.CustomFields.TextValueByNameBasic(UserFieldName_email))

				if leadVSChangeLogVo.PreviousVSUserGid == "" {
					tUser = newUser
				} else {
					tUser, _ = c.UserUsecase.GetByGid(leadVSChangeLogVo.PreviousVSUserGid)
				}
				verifyUser := c.VerifyUserEmailConfig(tUser)
				if !verifyUser {
					// 需要通知所有的all collaborators
					c.log.Error("LeadVS config is wrong:" + InterfaceToString(clientCase.Id()))
					return errors.New("LeadVS config is wrong:" + InterfaceToString(clientCase.Id())), "", "", "", "", ""
				}

				if tUser.Gid() != newUser.Gid() {
					cc = append(cc, newUser.CustomFields.TextValueByNameBasic(UserFieldName_email))
				}
			}
		} else if tplValue == MailGenre_AmCongratulationsNewRating ||
			tplValue == MailGenre_StartYourVADisabilityClaimRepresentation ||
			tplValue == MailGenre_AmContractReminder ||
			tplValue == MailGenre_AmIntakeFormReminder ||
			tplValue == MailGenre_VAForm2122aSubmission {
			useAmSendConfig = true
		} else {
			tUser, _ = c.UserUsecase.GetUserByLeadVS(clientCase)
			if !NeedUseSystemEmailConfig(tpl.CustomFields.TextValueByNameBasic("tpl")) {
				verifyUser := c.VerifyUserEmailConfig(tUser)
				if !verifyUser {
					// 需要通知所有的all collaborators
					c.log.Error("LeadVS config is wrong:" + InterfaceToString(clientCase.Id()))
					return errors.New("LeadVS config is wrong:" + InterfaceToString(clientCase.Id())), "", "", "", "", ""
				}
			}
		}
	}

	mailServiceConfig := InitMailServiceConfig()

	if useAmSendConfig {
		mailServiceConfig = InitAmMailServiceConfig()
	} else if tUser != nil {
		//clientMap.Set("users:email", tUser.CustomFields.TextValueByName("email"))
		clientMap.Set("users:first_name", tUser.CustomFields.TextValueByName("first_name"))
		clientMap.Set("users:last_name", tUser.CustomFields.TextValueByName("last_name"))
		clientMap.Set("users:title", tUser.CustomFields.TextValueByName("title"))
		clientMap.Set("users:mobile", tUser.CustomFields.TextValueByName("mobile"))
		if tUser.CustomFields.TextValueByNameBasic("mail_username") != "" {
			mailServiceConfig.Name = tUser.CustomFields.TextValueByNameBasic("full_name")
			mailServiceConfig.Username = tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailSender)
			mailPassword := tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailPassword)
			mailPassword, err = DecryptSensitive(mailPassword)
			if err != nil {
				return err, "", "", "", "", ""
			}
			mailServiceConfig.Password = mailPassword
		}
	}
	//return err, "", "", "", "", ""
	MailAttach := MailAttach_No

	tplvalue := tpl.CustomFields.TextValueByNameBasic("tpl")

	var mailAttachmentInputs MailAttachmentInputs

	if tpl.CustomFields.TextValueByName("tpl") != nil &&
		*tpl.CustomFields.TextValueByName("tpl") == MailGenre_FeeScheduleCommunication {
		clientMap.Set("current_rating", tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating))

		currentEvaluation, err := c.MailFeeContentUsecase.GetCurrentEvaluation(int(tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating)))
		if err != nil {
			return err, "", "", "", "", ""
		}
		clientMap.Set("current_evaluation", lib.NumberEnglishPrinter(int64(currentEvaluation)))

		fees, err := c.FeeUsecase.VBCFees(clientCase)
		if err != nil {
			return err, "", "", "", "", ""
		}
		feeBody := ""
		for _, v := range fees {
			feeBody += fmt.Sprintf("<li style=\"font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;\">%d%% evaluation: $%s (Payable only upon receiving a %d%% evaluation)</li>", v.Rating, lib.NumberEnglishPrinter(int64(v.Fee)), v.Rating)
		}
		//lib.DPrintln(feeBody)
		clientMap.Set("fee_body", feeBody)
		earningBody, err := c.MailFeeContentUsecase.GenContent(int(mailTaskInput.SubId))
		c.log.Debug("MailFeeContentUsecase:", earningBody)
		if err != nil {
			c.log.Error(err, "mailTaskInput.SubId:", mailTaskInput.SubId)
			return err, "", "", "", "", ""
		}
		clientMap.Set("earning_body", earningBody)

	} else if tpl.CustomFields.TextValueByName("tpl") != nil &&
		*tpl.CustomFields.TextValueByName("tpl") == MailGenre_GettingStartedEmail {
		firstName := tContactFields.TextValueByNameBasic("first_name")
		lastName := tContactFields.TextValueByNameBasic("last_name")
		clientMap.Set("first_name", firstName)
		clientMap.Set("last_name", lastName)
		uniqCode := tClientCaseFields.TextValueByNameBasic(FieldName_uniqcode)
		// online version
		//intakeFormUrl := fmt.Sprintf("https://docs.google.com/forms/d/e/1FAIpQLScxvX2bB-dfEY-VGDhkzuXvZfRai9bqvnkQNEHHU9ToN1lqWA/viewform?usp=pp_url&entry.1764873638=%s&entry.1434405649=%s&entry.738045322=%s", firstName, lastName, uniqCode)

		// zoho version
		//intakeFormUrl := fmt.Sprintf("https://docs.google.com/forms/d/e/1FAIpQLSe2nn2rjlIB1yC7f044cDgDKDlQIMJvYQ7xfLyYx6JtLvBxmA/viewform?usp=pp_url&entry.738045322=%s&entry.1764873638=%s&entry.1434405649=%s", uniqCode, firstName, lastName)

		//intakeFormUrl := fmt.Sprintf("https://docs.google.com/forms/d/e/1FAIpQLScxvX2bB-dfEY-VGDhkzuXvZfRai9bqvnkQNEHHU9ToN1lqWA/viewform?usp=pp_url&entry.1764873638=%s&entry.1434405649=%s&entry.738045322=%s", firstName, lastName, uniqCode)

		// Enable Jotform Intake Form
		intakeFormUrl := fmt.Sprintf("https://hipaa.jotform.com/242466899584074?name[first]=%s&name[last]=%s&vbcCase=%s", firstName, lastName, uniqCode)

		clientMap.Set("intake_form_url", intakeFormUrl)
		MailAttach = MailAttach_HowtoGuide
	} else if tplvalue == MailGenre_AmGettingStartedEmail {

		firstName := tContactFields.TextValueByNameBasic("first_name")
		lastName := tContactFields.TextValueByNameBasic("last_name")
		clientMap.Set("first_name", firstName)
		clientMap.Set("last_name", lastName)
		MailAttach = MailAttach_HowtoGuide

	} else if tplvalue == MailGenre_StartYourVADisabilityClaimRepresentation {

		clientMap.Set("attorney_full_name", tClientCaseFields.TextValueByNameBasic(FieldName_attorney_uniqid))
		clientMap.Set("current_rating", tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating))
		currentEvaluation, err := c.MailFeeContentUsecase.GetCurrentEvaluation(int(tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating)))
		if err != nil {
			return err, "", "", "", "", ""
		}
		clientMap.Set("current_evaluation", lib.NumberEnglishPrinter(int64(currentEvaluation)))
		firstName := tContactFields.TextValueByNameBasic("first_name")
		lastName := tContactFields.TextValueByNameBasic("last_name")
		clientMap.Set("first_name", firstName)
		clientMap.Set("last_name", lastName)
		uniqCode := tClientCaseFields.TextValueByNameBasic(FieldName_uniqcode)
		intakeFormUrl := fmt.Sprintf("https://hipaa.jotform.com/251865711410149?name[first]=%s&name[last]=%s&vbcCase=%s", firstName, lastName, uniqCode)
		clientMap.Set("intake_form_url", intakeFormUrl)

	} else if tplvalue == MailGenre_AmIntakeFormReminder {

		clientMap.Set("current_rating", tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating))
		firstName := tContactFields.TextValueByNameBasic("first_name")
		lastName := tContactFields.TextValueByNameBasic("last_name")
		clientMap.Set("first_name", firstName)
		clientMap.Set("last_name", lastName)
		uniqCode := tClientCaseFields.TextValueByNameBasic(FieldName_uniqcode)
		intakeFormUrl := fmt.Sprintf("https://hipaa.jotform.com/251865711410149?name[first]=%s&name[last]=%s&vbcCase=%s", firstName, lastName, uniqCode)
		clientMap.Set("intake_form_url", intakeFormUrl)

	} else if tpl.CustomFields.TextValueByName("tpl") != nil &&
		*tpl.CustomFields.TextValueByName("tpl") == MailGenre_NotifyVsRemindClient {

		if tUser == nil {
			return errors.New("tUser is nil."), "", "", "", "", ""
		}
		userEmail := tUser.CustomFields.TextValueByNameBasic("email")
		if userEmail == "" || !lib.VerifyEmail(userEmail) {
			return errors.New("userEmail is  incorrect."), "", "", "", "", ""
		}
		clientMap.Set("email", userEmail)

		a := SpawnRemindFeeContractSigningByEmail(RemindFeeContractSigningParams{})
		payload := AccessControlWorkPayload{}
		payload.Tasks = append(payload.Tasks, a)
		now := time.Now()
		expiredAt := now.Add(24 * time.Hour)
		token, err := c.AccessControlWorkUsecase.CreateAccessControlWork(WorkType_remind_fee_contract_signing,
			InterfaceToString(tClientCaseFields.NumberValueByNameBasic("id")),
			payload,
			expiredAt)
		if err != nil {
			return err, "", "", "", "", ""
		}

		remindUrl := fmt.Sprintf("%s/process/remind.html?token=%s", c.conf.Domain, token)
		clientMap.Set("remind_url", remindUrl)

	} else if tpl.CustomFields.TextValueByName("tpl") != nil &&
		*tpl.CustomFields.TextValueByName("tpl") == MailGenre_CongratulationsNewRating {

		currentRating := tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating)
		newRating := tClientCaseFields.NumberValueByNameBasic(FieldName_new_rating)
		congratsText := ""
		if currentRating > 0 {
			congratsText = fmt.Sprintf(`Congrats on your new VA disability rating increase from %d%% to %d%%!`, currentRating, newRating)
		} else {
			congratsText = fmt.Sprintf(`Congrats on your VA disability rating increase to %d%%!`, newRating)
		}
		clientMap.Set("congrats_text", congratsText)
		MailAttach = MailAttach_VBCClothingAllowanceGuide

	} else if tplvalue == MailGenre_AmCongratulationsNewRating {

		currentRating := tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating)
		newRating := tClientCaseFields.NumberValueByNameBasic(FieldName_new_rating)
		congratsText := ""
		if currentRating > 0 {
			congratsText = fmt.Sprintf(`Congrats on your new VA disability rating increase from %d%% to %d%%!`, currentRating, newRating)
		} else {
			congratsText = fmt.Sprintf(`Congrats on your VA disability rating increase to %d%%!`, newRating)
		}
		clientMap.Set("congrats_text", congratsText)
		MailAttach = MailAttach_VBCClothingAllowanceGuide

	} else if tpl.CustomFields.TextValueByName("tpl") != nil &&
		*tpl.CustomFields.TextValueByName("tpl") == MailGenre_ContractReminder {

		contractDateOn, err := c.ClientEnvelopeUsecase.ContractDateOn(tClientCaseFields.NumberValueByNameBasic("id"), false)
		if err != nil {
			return err, "", "", "", "", ""
		}
		clientMap.Set("contract_date", contractDateOn)

	} else if tplvalue == MailGenre_AmContractReminder {

		contractDateOn, err := c.ClientEnvelopeUsecase.ContractDateOn(tClientCaseFields.NumberValueByNameBasic("id"), true)
		if err != nil {
			return err, "", "", "", "", ""
		}
		clientMap.Set("contract_date", contractDateOn)

	} else if tplvalue == MailGenre_HelpUsImproveSurvey {
		firstName := tContactFields.TextValueByNameBasic("first_name")
		lastName := tContactFields.TextValueByNameBasic("last_name")
		clientMap.Set("first_name", firstName)
		clientMap.Set("last_name", lastName)
		uniqCode := tClientCaseFields.TextValueByNameBasic(FieldName_uniqcode)
		bizUrl := fmt.Sprintf("https://hipaa.jotform.com/251598744529169?fullName[first]=%s&fullName[last]=%s&uniqueId=%s", firstName, lastName, uniqCode)
		clientMap.Set("biz_url", bizUrl)

	} else if tplvalue == MailGenre_VAForm2122aSubmission {

		attorneyUniqid := clientCase.CustomFields.TextValueByNameBasic(FieldName_attorney_uniqid)
		if attorneyUniqid == "" {
			return errors.New("attorneyUniqid is empty"), "", "", "", "", ""
		}
		uniqcode := clientCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
		attorneyEntity, _ := c.AttorneyUsecase.GetByGid(attorneyUniqid)
		if attorneyEntity == nil {
			return errors.New("attorneyEntity is nil"), "", "", "", "", ""
		}
		clientMap.Set("email", attorneyEntity.RoEmail)
		vaPdfBytes, err := c.SendVa2122aUsecase.GetAmSignedVA2122aBytes(clientCase.Id())
		if err != nil {
			return err, "", "", "", "", ""
		}
		filename := "Signed VA 21-22a " + uniqcode + ".pdf"
		headers := map[string][]string{
			"Content-Type":        {"application/pdf"},
			"Content-Disposition": {`attachment; filename="` + filename + `"`},
		}
		mailAttachmentInputs = append(mailAttachmentInputs, MailAttachmentInput{
			Name:   filename,
			Reader: bytes.NewBuffer(vaPdfBytes),
			Settings: []gomail.FileSetting{
				gomail.SetHeader(headers),
			},
		})
	} else if tpl.CustomFields.TextValueByName("tpl") != nil &&
		(*tpl.CustomFields.TextValueByName("tpl") == MailGenre_PersonalStatementsReadyforYourReview ||
			*tpl.CustomFields.TextValueByName("tpl") == MailGenre_PleaseReviewYourPersonalStatementsinSharedFolder) {

		clientMap.Set("statement_url", PersonalStatementManagerUrl(clientCase.Gid()))
		password, err := c.StatementUsecase.PersonalStatementPassword(clientCase.Id())
		if err != nil {
			c.log.Error(err)
		}
		if password == "" {
			c.log.Error("caseId: ", clientCase.Id(), " password is empty")
		}

		clientMap.Set("statement_password", password)

	}
	e, s, b, to := c.SendEmailWithTypeMap(clientMap, tpl, mailServiceConfig, MailAttach, cc, mailAttachmentInputs)
	return e, s, b, to, mailServiceConfig.Username, mailServiceConfig.Name
}

func (c *MailUsecase) SendEmailWithTypeMap(typeMap lib.TypeMap, tpl *TData, mailServiceConfig *MailServiceConfig, MailAttach string, cc []string, mailAttachmentInputs MailAttachmentInputs) (err error, subject string, body string, toEmail string) {

	if tpl == nil {
		return errors.New("Email Tpl is empty."), "", "", ""
	}
	if typeMap == nil {
		return errors.New("TypeMap is empty."), "", "", ""
	}
	email := typeMap.GetString("email")

	if len(email) == 0 || !lib.VerifyEmail(email) {
		return errors.New("SendEmailWithTypeMap: Email format is wrong."), "", "", ""
	}

	if tpl.CustomFields.TextValueByName("subject") != nil {
		subject = *(tpl.CustomFields.TextValueByName("subject"))
	}
	if tpl.CustomFields.TextValueByName("body") != nil {
		body = *(tpl.CustomFields.TextValueByName("body"))
	}

	if configs.IsDev() {
		//body = email_config.GettingStartedEmail
	}
	//email = "liaogling@gmail.com"
	//email = "yannanwang@gmail.com"
	//email = "lialing@foxmail.com"
	mailMessage := &MailMessage{
		To:      email,
		Subject: MailReplaceDynamicParams(subject, typeMap),
		Body:    MailReplaceDynamicParams(body, typeMap),
		//Cc:      []string{"team@vetbenefitscenter.com"},
	}
	if !configs.IsDev() {
		//if mailServiceConfig.Username != AmSenderEamil {
		cc = append(cc, "info@vetbenefitscenter.com")
		//}
	}
	if len(cc) > 0 {
		mailMessage.Cc = cc
	}
	//mailMessage.Subject = "Testing stmp is working."
	//mailMessage.Body = "Testing body."
	return c.SendEmail(mailServiceConfig, mailMessage, MailAttach, mailAttachmentInputs), mailMessage.Subject, mailMessage.Body, email
}

// SendSystemMessage 系统邮件
func (c *MailUsecase) SendSystemMessage(subject string, body string, email string, ccs []string) error {
	mailMessage := &MailMessage{
		To:      email,
		Subject: subject,
		Body:    body,
	}
	if len(ccs) > 0 {
		mailMessage.Cc = ccs
	}
	mailServiceConfig := InitMailServiceConfig()
	return c.SendEmail(mailServiceConfig, mailMessage, "", nil)
}
