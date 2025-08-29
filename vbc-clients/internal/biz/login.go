package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

const Cookie_Redirect_Url = "__red_url"

type LoginUsecase struct {
	log                 *log.Helper
	CommonUsecase       *CommonUsecase
	conf                *conf.Data
	Oauth2ClientUsecase *Oauth2ClientUsecase
}

func NewLoginUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	Oauth2ClientUsecase *Oauth2ClientUsecase) *LoginUsecase {
	uc := &LoginUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		Oauth2ClientUsecase: Oauth2ClientUsecase,
	}
	return uc
}

func (c *LoginUsecase) HttpLogin(ctx *gin.Context) {

	redirectUrl := ctx.Query("redirect_url")
	from := ctx.Query("from")
	url, err := c.Login(from)
	if err != nil {
		c.log.Error(err)
		ctx.AbortWithError(500, err)
		return
	}
	ctx.SetCookie(Cookie_Redirect_Url, redirectUrl, 3600, "/", c.conf.Domain, false, true)
	ctx.Redirect(302, url)
}

func (c *LoginUsecase) Login(from string) (url string, err error) {

	if from == Oauth2_AppId_vbcapp {
		return c.Oauth2ClientUsecase.AuthUrl(from)
	} else if from == Oauth2_AppId_microsoft_vbcapp {
		return c.Oauth2ClientUsecase.AuthUrl(from)
	} else {
		return "", errors.New("from is wrong: " + from)
	}

	return
}
