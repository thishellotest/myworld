package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib/oauth2"
	"vbc/lib/oauth2/endpoints"
	"vbc/lib/oauth2/google"
	"vbc/lib/oauth2/microsoft"
)

const (
	Oauth2_AppId_docusign  = "docusign"
	Oauth2_AppId_asana     = "asana"
	Oauth2_AppId_box       = "box"
	Oauth2_AppId_adobesign = "adobesign"
	Oauth2_AppId_google    = "google"
	Oauth2_AppId_xero      = "xero"
	Oauth2_AppId_zohocrm   = "zohocrm"
	Oauth2_AppId_vbcapp    = "vbcapp" // 使用google帐号登录
	Oauth2_AppId_zoom      = "zoom"

	// https://learn.microsoft.com/zh-cn/entra/identity-platform/scenario-web-app-sign-user-app-registration?tabs=aspnetcore
	// https://developer.microsoft.com/en-us/graph/graph-explorer
	Oauth2_AppId_microsoft_vbcapp    = "msvbcapp"
	Oauth2_microsoft_vbcapp_ClientId = "ae4a7ecc-c768-437e-9408-89f8b99789aa"
)

type Oauth2ClientEntity struct {
	ID           int32 `gorm:"primaryKey"`
	ClientId     string
	ClientSecret string
	AppId        string
	DeletedAt    int64
}

func (Oauth2ClientEntity) TableName() string {
	return "oauth2_clients"
}

func (c *Oauth2ClientEntity) Oauth2Config(conf *conf.Data) *oauth2.Config {

	domain := conf.Domain
	if configs.AppEnvType() == configs.ENV_TYPE_DEV_TEST {
		//domain = "http://123.57.61.3:8050"
	}

	config := &oauth2.Config{
		RedirectURL:  domain + "/oauth2/callback?app_id=" + c.AppId,
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
	}
	if c.AppId == Oauth2_AppId_google {
		config.Scopes = []string{"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/drive.readonly",
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/spreadsheets.readonly"}
		config.Endpoint = google.Endpoint
		config.Endpoint.AuthURL += "?prompt=consent&access_type=offline" // 这样有可能每次都返回 refresh_token
		// 文章：https://medium.com/starthinker/google-oauth-2-0-access-token-and-refresh-token-explained-cccf2fc0a6d9
		return config
	}
	if c.AppId == Oauth2_AppId_vbcapp {
		config.Scopes = []string{"https://www.googleapis.com/auth/userinfo.email"}
		config.Endpoint = google.Endpoint
		config.Endpoint.AuthURL += "?prompt=consent&access_type=offline" // 这样有可能每次都返回 refresh_token
		// 文章：https://medium.com/starthinker/google-oauth-2-0-access-token-and-refresh-token-explained-cccf2fc0a6d9
		return config
	}
	if c.AppId == Oauth2_AppId_microsoft_vbcapp {
		// "https://graph.microsoft.com/profile", https://graph.microsoft.com/email 这两个都应该没用
		config.Scopes = []string{"https://graph.microsoft.com/User.Read"}
		config.Endpoint = microsoft.AzureADEndpoint(Oauth2_microsoft_vbcapp_ClientId)
		return config
	}

	if c.AppId == Oauth2_AppId_xero {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:   "https://login.xero.com/identity/connect/authorize?access_type=offline_access",
			TokenURL:  "https://identity.xero.com/connect/token",
			AuthStyle: oauth2.AuthStyleInHeader, // : 需要设置
		}
		// https://developer.xero.com/documentation/guides/oauth2/scopes/
		config.Scopes = []string{"offline_access", "openid", "profile", "email", "accounting.transactions",
			"accounting.transactions.read", "accounting.reports.read",
			"accounting.settings", "accounting.settings.read",
			"accounting.contacts", "accounting.contacts.read", "accounting.attachments",
			"accounting.attachments.read"}
		return config
	}
	if c.AppId == Oauth2_AppId_zohocrm {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:   "https://accounts.zoho.com/oauth/v2/auth?access_type=offline",
			TokenURL:  "https://accounts.zoho.com/oauth/v2/token",
			AuthStyle: oauth2.AuthStyleInHeader, // : 需要设置
		}
		// "AaaServer.profile.Read",
		config.Scopes = []string{"ZohoCRM.modules.ALL",
			"ZohoCRM.settings.READ", "ZohoCRM.modules.leads.ALL",
			"ZohoCRM.modules.deals.ALL",
			"ZohoCRM.Users.ALL",
			"ZohoCRM.settings.ALL",
			"ZohoCRM.notifications.ALL",
			"ZohoCRM.share.leads.ALL",
			"ZohoCRM.share.contacts.ALL",
			"ZohoCRM.share.deals.ALL",
			"ZohoCRM.share.accounts.ALL",
		}
		return config
	}

	if c.AppId == Oauth2_AppId_docusign {
		config.Scopes = []string{
			"signature",
			//"openid",
			//"extended",
			//"impersonation",
		}
		if configs.IsDev() {
			config.Endpoint = oauth2.Endpoint{
				AuthURL:   "https://account-d.docusign.com/oauth/auth",
				TokenURL:  "https://account-d.docusign.com/oauth/token",
				AuthStyle: oauth2.AuthStyleInHeader, // : 需要设置
			}
		} else {
			// todo: use prod
			config.Endpoint = oauth2.Endpoint{
				AuthURL:   "https://account-d.docusign.com/oauth/auth",
				TokenURL:  "https://account-d.docusign.com/oauth/token",
				AuthStyle: oauth2.AuthStyleInHeader, // : 需要设置
			}
		}
	} else if c.AppId == Oauth2_AppId_asana {
		if configs.IsDev() {
			config.Endpoint = oauth2.Endpoint{
				AuthURL:   "https://app.asana.com/-/oauth_authorize",
				TokenURL:  "https://app.asana.com/-/oauth_token",
				AuthStyle: oauth2.AuthStyleInHeader,
			}
		} else {
			config.Endpoint = oauth2.Endpoint{
				AuthURL:   "https://account.docusign.com/oauth/auth",
				TokenURL:  "https://account.docusign.com/oauth/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			}
		}
	} else if c.AppId == Oauth2_AppId_box {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:   "https://account.box.com/api/oauth2/authorize",
			TokenURL:  "https://api.box.com/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		}
	} else if c.AppId == Oauth2_AppId_adobesign {

		config.Scopes = []string{
			"agreement_read:account",
			"agreement_write:account",
			"agreement_send:account",
			"widget_read:account",
			"widget_write:account",
			"library_read:account",
			"library_write:account",
			"webhook_read:account",
			"webhook_write:account",
			"webhook_retention:account",
		}
		config.Endpoint = oauth2.Endpoint{
			AuthURL:    "https://secure.na3.adobesign.com/public/oauth/v2",
			TokenURL:   "https://api.na3.adobesign.com/oauth/v2/token",
			RefreshURL: "https://api.na3.adobesign.com/oauth/v2/refresh",
			AuthStyle:  oauth2.AuthStyleInParams,
		}
	} else if c.AppId == Oauth2_AppId_zoom {
		config.Endpoint = oauth2.Endpoint{
			AuthURL:   endpoints.Zoom.AuthURL,
			TokenURL:  endpoints.Zoom.TokenURL,
			AuthStyle: oauth2.AuthStyleInHeader,
		}
	}
	return config
}

type Oauth2ClientUsecase struct {
	CommonUsecase *CommonUsecase
	log           *log.Helper
	DBUsecase[Oauth2ClientEntity]
	conf *conf.Data
}

func NewOauth2ClientUsecase(CommonUsecase *CommonUsecase, logger log.Logger, conf *conf.Data) *Oauth2ClientUsecase {
	uc := &Oauth2ClientUsecase{
		CommonUsecase: CommonUsecase,
		log:           log.NewHelper(logger),
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *Oauth2ClientUsecase) GetByAppId(appId string) (*Oauth2ClientEntity, error) {

	var entity Oauth2ClientEntity
	err := c.CommonUsecase.DB().Where("app_id=? and deleted_at=0", appId).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, nil
}

func (c *Oauth2ClientUsecase) GetByClientId(clientId string) (*Oauth2ClientEntity, error) {

	var entity Oauth2ClientEntity
	err := c.CommonUsecase.DB().Where("client_id=? and deleted_at=0", clientId).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, nil
}

// AuthUrl auth
func (c *Oauth2ClientUsecase) AuthUrl(appId string) (authUrl string, err error) {
	client, _ := c.GetByAppId(appId)
	if client == nil {
		return "", errors.New("client is nil.")
	}
	conf := client.Oauth2Config(c.conf)
	authUrl = conf.AuthCodeURL("")
	return
}

func (c *Oauth2ClientUsecase) Exchange(appId string, code string) (*oauth2.Token, error) {
	client, _ := c.GetByAppId(appId)
	if client == nil {
		return nil, errors.New("client is nil.")
	}
	return client.Oauth2Config(c.conf).Exchange(context.TODO(), code)
}
