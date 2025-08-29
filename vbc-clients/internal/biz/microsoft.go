package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type MicrosoftUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
}

func NewMicrosoftUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *MicrosoftUsecase {
	uc := &MicrosoftUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func (c *MicrosoftUsecase) BaseUrl() string {
	return "https://graph.microsoft.com"
}

type MicrosoftUser struct {
	Mail string `json:"mail"`
}

func (c *MicrosoftUsecase) UserEmail(accessToken string) (email string, err error) {

	// todo:lgl
	res, _, err := lib.Request("GET", c.BaseUrl()+"/v1.0/me", nil, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("res is nil")
	}

	//var user MicrosoftUser
	user, err := lib.StringToTE[MicrosoftUser](*res, MicrosoftUser{})
	if err != nil {
		return "", err
	}
	return user.Mail, nil
}
