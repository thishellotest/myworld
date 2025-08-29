package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"net/url"
	"strings"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HttpOauth2Usecase struct {
	log                 *log.Helper
	Oauth2TokenUsecase  *Oauth2TokenUsecase
	Oauth2ClientUsecase *Oauth2ClientUsecase
	CommonUsecase       *CommonUsecase
	UserUsecase         *UserUsecase
	LoginBindingUsecase *LoginBindingUsecase
	AppTokenUsecase     *AppTokenUsecase
	LogUsecase          *LogUsecase
	MicrosoftUsecase    *MicrosoftUsecase
}

func NewHttpOauth2Usecase(logger log.Logger, Oauth2TokenUsecase *Oauth2TokenUsecase, Oauth2ClientUsecase *Oauth2ClientUsecase,
	CommonUsecase *CommonUsecase,
	UserUsecase *UserUsecase,
	LoginBindingUsecase *LoginBindingUsecase,
	AppTokenUsecase *AppTokenUsecase,
	LogUsecase *LogUsecase,
	MicrosoftUsecase *MicrosoftUsecase) *HttpOauth2Usecase {
	return &HttpOauth2Usecase{
		log:                 log.NewHelper(logger),
		Oauth2TokenUsecase:  Oauth2TokenUsecase,
		Oauth2ClientUsecase: Oauth2ClientUsecase,
		CommonUsecase:       CommonUsecase,
		UserUsecase:         UserUsecase,
		LoginBindingUsecase: LoginBindingUsecase,
		AppTokenUsecase:     AppTokenUsecase,
		LogUsecase:          LogUsecase,
		MicrosoftUsecase:    MicrosoftUsecase,
	}
}

func (c *HttpOauth2Usecase) Callback(ctx *gin.Context) {

	reply := CreateReply()
	appId := ctx.Query("app_id")
	code := ctx.Query("code")

	if appId == Oauth2_AppId_vbcapp || appId == Oauth2_AppId_microsoft_vbcapp {
		var appTokenEntity *AppTokenEntity
		var err error
		if appId == Oauth2_AppId_vbcapp {
			appTokenEntity, err = c.BizVbcapp(ctx, code)
		} else if appId == Oauth2_AppId_microsoft_vbcapp {
			appTokenEntity, err = c.BizMSVbcapp(ctx, code)
		}
		if err != nil {
			reply.CommonError(err)
			goto end
		}
		if appTokenEntity == nil {
			reply.CommonStrError("appTokenEntity is nil")
			goto end
		}
		redirectUrl, err := ctx.Cookie(Cookie_Redirect_Url)
		if err != nil {
			reply.CommonError(err)
			goto end
		}
		if redirectUrl != "" {

			redirectUrl, _ = url.QueryUnescape(redirectUrl)
			if strings.Index(redirectUrl, "?") >= 0 {
				redirectUrl += "&jwt=" + appTokenEntity.AccessToken
			} else {
				redirectUrl += "?jwt=" + appTokenEntity.AccessToken
			}
			ctx.Redirect(302, redirectUrl)
			return
		} else {
			params := make(lib.TypeMap)
			params.Set("jwt", appTokenEntity.AccessToken)
			reply.Merge(params)
		}
	} else {
		a, err := c.Oauth2ClientUsecase.Exchange(appId, code)
		if err != nil {
			reply.CommonError(err)
			goto end
		}

		if a != nil {
			er := c.Oauth2TokenUsecase.UpdateByToken(appId, a)
			if er != nil {
				c.log.Error(er)
			}
		}
	}

	reply.Success()
end:
	ctx.JSON(200, reply)
}

func (c *HttpOauth2Usecase) BizMSVbcapp(ctx context.Context, code string) (appTokenEntity *AppTokenEntity, err error) {

	appId := Oauth2_AppId_microsoft_vbcapp
	a, err := c.Oauth2ClientUsecase.Exchange(appId, code)
	if err != nil {
		return nil, err
	}

	msEmail, err := c.MicrosoftUsecase.UserEmail(a.AccessToken)
	c.log.Info("msEmail:", msEmail, "err:", err)
	// todo:lgl
	if err != nil {
		return nil, err
	}

	err = c.LogUsecase.SaveLog(0, "Oauth2:"+Oauth2_AppId_microsoft_vbcapp, map[string]interface{}{
		"msEmail": InterfaceToString(msEmail),
	})
	if err != nil {
		c.log.Error(err)
	}

	tUser, err := c.LoginBindingUsecase.WithMicrosoft(msEmail)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("\"" + msEmail + "\" is also bound to the system, please contact us.")
	}
	return c.AppTokenUsecase.CreateToken(tUser.CustomFields.NumberValueByNameBasic("id"))
}

func (c *HttpOauth2Usecase) BizVbcapp(ctx context.Context, code string) (appTokenEntity *AppTokenEntity, err error) {

	appId := Oauth2_AppId_vbcapp
	a, err := c.Oauth2ClientUsecase.Exchange(appId, code)
	if err != nil {
		return nil, err
	}
	srv, err := oauth2.NewService(ctx, option.WithTokenSource(a.StaticTokenSource()))
	if err != nil {
		return nil, err
	}
	userSrv := oauth2.NewUserinfoV2Service(srv)
	userInfo, err := userSrv.Me.Get().Do()
	// {"email":"glliao@vetbenefitscenter.com","hd":"vetbenefitscenter.com","id":"109700294120500814385","picture":"https://lh3.googleusercontent.com/a-/ALV-UjUAoGtECEKS6gdl9R_76OvyQuNx8XmvEL1IoLx1Eqg0DB6alQ=s96-c","verified_email":true}
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	if userInfo != nil {
		err = c.LogUsecase.SaveLog(0, "Oauth2:"+Oauth2_AppId_vbcapp, map[string]interface{}{
			"userInfo": InterfaceToString(userInfo),
		})
		if err != nil {
			c.log.Error(err)
		}
	}

	tUser, err := c.LoginBindingUsecase.WithGoogle(userInfo)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	if tUser == nil {
		return nil, errors.New("The email is also bound to the system, please contact us.")
	}
	return c.AppTokenUsecase.CreateToken(tUser.CustomFields.NumberValueByNameBasic("id"))
}

func (c *HttpOauth2Usecase) AuthList(ctx *gin.Context) {

	reply := CreateReply()
	var tList lib.TypeList
	records, err := c.Oauth2ClientUsecase.AllByCond(And(Eq{"deleted_at": 0}))
	if err != nil {
		reply.InternalError(err)
		goto end
	}
	for _, v := range records {
		tMap := make(lib.TypeMap)
		tMap.Set("app_id", v.AppId)
		a, _ := c.Oauth2ClientUsecase.AuthUrl(v.AppId)
		tMap.Set("auth_url", a)
		tList = append(tList, tMap)
	}
	reply["list"] = tList
	reply.Success()
end:
	ctx.JSON(200, reply)
}
