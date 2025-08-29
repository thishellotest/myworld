package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/api/oauth2/v2"
	"vbc/internal/conf"
)

type LoginBindingUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	UserUsecase   *UserUsecase
}

func NewLoginBindingUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	UserUsecase *UserUsecase) *LoginBindingUsecase {
	uc := &LoginBindingUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		UserUsecase:   UserUsecase,
	}

	return uc
}

// WithGoogle {"email":"glliao@vetbenefitscenter.com","hd":"vetbenefitscenter.com","id":"109700294120500814385","picture":"https://lh3.googleusercontent.com/a-/ALV-UjUAoGtECEKS6gdl9R_76OvyQuNx8XmvEL1IoLx1Eqg0DB6alQ=s96-c","verified_email":true}
func (c *LoginBindingUsecase) WithGoogle(userInfo *oauth2.Userinfo) (tUser *TData, err error) {
	if userInfo == nil {
		return nil, errors.New("userInfo is nil")
	}
	email := userInfo.Email
	if email == "" {
		return nil, errors.New("email is empty")
	}
	return c.UserUsecase.GetByEmail(email)
}

func (c *LoginBindingUsecase) WithMicrosoft(email string) (tUser *TData, err error) {
	if email == "" {
		return nil, errors.New("email is empty")
	}
	return c.UserUsecase.GetByMicrosoftEmail(email)
}
