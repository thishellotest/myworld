package biz

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"regexp"
	"vbc/internal/conf"
	"vbc/lib"
)

type UserHttpUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	DialpadUsecase   *DialpadUsecase
	JWTUsecase       *JWTUsecase
	UserUsecase      *UserUsecase
	DataEntryUsecase *DataEntryUsecase
	MailUsecase      *MailUsecase
}

func NewUserHttpUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DialpadUsecase *DialpadUsecase,
	JWTUsecase *JWTUsecase,
	UserUsecase *UserUsecase,
	DataEntryUsecase *DataEntryUsecase,
	MailUsecase *MailUsecase,
) *UserHttpUsecase {
	uc := &UserHttpUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DialpadUsecase:   DialpadUsecase,
		JWTUsecase:       JWTUsecase,
		UserUsecase:      UserUsecase,
		DataEntryUsecase: DataEntryUsecase,
		MailUsecase:      MailUsecase,
	}

	return uc
}

func (c *UserHttpUsecase) VerifyEmailOutbox(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizVerifyEmailOutbox(userFacade, body.GetString("user_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *UserHttpUsecase) BizVerifyEmailOutbox(userFacade UserFacade, userGid string) (lib.TypeMap, error) {

	tProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}
	if !IsAdminProfile(tProfile) {
		return nil, errors.New("No permission management")
	}
	tUser, err := c.TUsecase.DataByGid(Kind_users, userGid)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	mailSender := tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailSender)
	if mailSender == "" {
		return nil, errors.New("Google Mail Username is empty")
	}
	mailPassword := tUser.CustomFields.TextValueByNameBasic(UserFieldName_MailPassword)
	if mailPassword == "" {
		return nil, errors.New("Google App Password is empty")
	}
	mailPassword, err = DecryptSensitive(mailPassword)
	if err != nil {
		return nil, err
	}

	serviceConfig := &MailServiceConfig{
		Name:        tUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
		Host:        "smtp.gmail.com",
		Port:        587,
		Username:    mailSender,
		Password:    mailPassword,
		FromAddress: mailSender,
	}

	message := &MailMessage{
		To:      "engineering@vetbenefitscenter.com",
		Subject: "Email Configuration Test from The Base",
		Body: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Configuration Test</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <p>Dear User,</p>
    
    <p>This is a test email from The Base to verify that your email configuration is set up correctly.</p>

    <p>If you have received this email, your email system is functioning properly.</p>

    <p>Best regards,</p>
    <p>The Base Team</p>
</body>
</html>
`,
	}

	err = c.MailUsecase.SendEmail(serviceConfig, message, "", nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *UserHttpUsecase) SyncDailpad(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizSyncDailpad(userFacade, body.GetString("user_gid"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func FormatPhoneNumberV2(raw string) (string, error) {
	// 正则匹配 +1 后的10位数字
	re := regexp.MustCompile(`^\+1(\d{3})(\d{3})(\d{4})$`)
	matches := re.FindStringSubmatch(raw)
	if len(matches) != 4 {
		return "", fmt.Errorf("invalid phone number format")
	}
	return fmt.Sprintf("+1 %s-%s-%s", matches[1], matches[2], matches[3]), nil
}

func (c *UserHttpUsecase) BizSyncDailpad(userFacade UserFacade, userGid string) (lib.TypeMap, error) {

	tProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}
	if !IsAdminProfile(tProfile) {
		return nil, errors.New("No permission management")
	}
	tUser, err := c.TUsecase.DataByGid(Kind_users, userGid)
	if err != nil {
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("tUser is nil")
	}
	email := tUser.CustomFields.TextValueByNameBasic(UserFieldName_email)
	if email == "" {
		return nil, errors.New("Email is empty")
	}
	dialpadEmail, dialpadPhone, dialpadId := c.DialpadUsecase.UserInfoByEmail(email)
	if dialpadEmail == "" || dialpadPhone == "" || dialpadId == "" {
		return nil, errors.New("No Dailpad information was obtained through email address " + email)
	}
	if dialpadEmail != email {
		return nil, errors.New("Please contact the administrator. An issue has occurred")
	}

	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = tUser.Gid()
	dataEntry[UserFieldName_dialpad_phonenumber] = dialpadPhone
	dataEntry[UserFieldName_dialpad_userid] = dialpadId
	if tUser.CustomFields.TextValueByNameBasic(UserFieldName_mobile) == "" {
		phone, err := FormatPhoneNumberV2(dialpadPhone)
		if err != nil {
			return nil, err
		}
		dataEntry[UserFieldName_mobile] = phone
	}

	_, err = c.DataEntryUsecase.HandleOne(Kind_users, dataEntry, DataEntry_gid, &userFacade.TData)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
