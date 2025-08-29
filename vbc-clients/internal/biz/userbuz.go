package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type UserbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
	LogUsecase       *LogUsecase
}

func NewUserbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	LogUsecase *LogUsecase,
) *UserbuzUsecase {
	uc := &UserbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
		LogUsecase:       LogUsecase,
	}

	return uc
}

func (c *UserbuzUsecase) HandleAllPassword() error {

	users, err := c.TUsecase.ListByCond(Kind_users, nil)
	if err != nil {
		return err
	}
	for _, v := range users {
		mailPassword := v.CustomFields.TextValueByNameBasic(UserFieldName_MailPassword)
		if mailPassword != "" {
			data := make(TypeDataEntry)
			data[DataEntry_gid] = v.Gid()
			a, err := EncryptSensitive(mailPassword)
			if err != nil {
				c.log.Error(err)
			} else {
				data[UserFieldName_MailPassword] = a
				c.DataEntryUsecase.HandleOne(Kind_users, data, DataEntry_gid, nil)
				c.LogUsecase.SaveLog(0, "EncryptSensitiveLog", map[string]interface{}{
					"id":              v.Id(),
					"mailPassword":    mailPassword,
					"newMailPassword": a,
				})
			}
		}
	}
	return nil
}
