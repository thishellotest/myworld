package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type MonitorUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	TUsecase      *TUsecase
}

func NewMonitorUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
) *MonitorUsecase {
	uc := &MonitorUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
	}

	return uc
}

func (c *MonitorUsecase) DoMonitorVSUsers() error {

	c.log.Info("DoMonitorVSUsers", time.Now().Format(time.RFC3339))
	users, err := c.TUsecase.ListByCond(Kind_users, And(Eq{"biz_deleted_at": 0, "deleted_at": 0},
		In("role_gid", "540135b887ed4a6da407fd4fdee0c4af", "ba8b82363dd646e0856c62b00402f978")))
	if err != nil {
		return err
	}
	for _, v := range users {
		alterMsg := ""
		name := v.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
		if name == "" {
			alterMsg += "Full Name is empty; "
		}
		mobile := v.CustomFields.TextValueByNameBasic(UserFieldName_mobile)
		if mobile == "" {
			alterMsg += "Phone is empty; "
		}
		dialpadUserid := v.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_userid)
		if dialpadUserid == "" {
			alterMsg += "Dialpad ID is empty; "
		}
		dialpadPhoneNumber := v.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_phonenumber)

		if dialpadPhoneNumber == "" {
			alterMsg += "Dialpad Mobile is empty; "
		}

		mailSenderUsername := v.CustomFields.TextValueByNameBasic(UserFieldName_MailSender)
		if mailSenderUsername == "" {
			alterMsg += "Google Mail Username is empty; "
		}
		mailPassword := v.CustomFields.TextValueByNameBasic(UserFieldName_MailPassword)
		if mailPassword == "" {
			alterMsg += "Google App Password is empty; "
		}

		if alterMsg != "" {
			alterMsg = name + " | " + v.CustomFields.TextValueByNameBasic(UserFieldName_email) + " : " + alterMsg
			c.log.Error(alterMsg)
		}
	}
	return nil
}
