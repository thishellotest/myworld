package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/oauth2"
)

type ZoomTokenUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func ZoomAccountId() string {
	return configs.EnvZoomAccountIdKey()
}

func ZoomClientId() string {
	return configs.EnvZoomClientIdKey()
}

func ZoomSecret() string {
	return configs.EnvZoomSecretKey()
}

func NewZoomTokenUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ZoomTokenUsecase {
	uc := &ZoomTokenUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	return uc
}

func (c *ZoomTokenUsecase) OauthToken() (*oauth2.Token, error) {
	info, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: info.GetString("access_token"),
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(3000 * time.Second),
	}, nil
}

func (c *ZoomTokenUsecase) GetAccessToken() (lib.TypeMap, error) {
	url := "https://zoom.us/oauth/token?grant_type=account_credentials&account_id=" + ZoomAccountId()

	headers := make(map[string]string)
	headers["Host"] = "ZoomToken.us"
	headers["Authorization"] = "Basic " + lib.BasicAuth(ZoomClientId(), ZoomSecret())

	res, _, err := lib.Request("POST", url, nil, headers)
	lib.DPrintln(res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("res is nil")
	}
	return lib.ToTypeMapByString(*res), nil

}
