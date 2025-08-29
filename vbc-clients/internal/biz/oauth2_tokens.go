package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	oauth22 "golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
	"vbc/internal/conf"
	"vbc/lib/oauth2"
	"xorm.io/builder"
)

//
//CREATE TABLE `oauth2_tokens` (
//`id` int NOT NULL AUTO_INCREMENT,
//`client_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`expires_at` int NOT NULL DEFAULT '0',
//`create_at` int NOT NULL DEFAULT '0',
//`update_at` int NOT NULL DEFAULT '0',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uniq_c` (`client_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

type Oauth2TokenEntity struct {
	ID           int32 `gorm:"primaryKey"`
	ClientId     string
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresAt    int64
	CreatedAt    int64
	UpdatedAt    int64
}

func (Oauth2TokenEntity) TableName() string {
	return "oauth2_tokens"
}

func (c *Oauth2TokenEntity) TokenSourceMod() (oauth22.TokenSource, error) {
	return oauth22.StaticTokenSource(&oauth22.Token{
		AccessToken: c.AccessToken,
		Expiry:      time.Unix(c.ExpiresAt, 0),
	}), nil
}

type Oauth2TokenUsecase struct {
	CommonUsecase       *CommonUsecase
	log                 *log.Helper
	Oauth2ClientUsecase *Oauth2ClientUsecase
	conf                *conf.Data
	ZoomTokenUsecase    *ZoomTokenUsecase
}

func NewOauth2TokenUsecase(CommonUsecase *CommonUsecase, logger log.Logger,
	Oauth2ClientUsecase *Oauth2ClientUsecase,
	conf *conf.Data,
	ZoomTokenUsecase *ZoomTokenUsecase) *Oauth2TokenUsecase {
	return &Oauth2TokenUsecase{
		CommonUsecase:       CommonUsecase,
		log:                 log.NewHelper(logger),
		Oauth2ClientUsecase: Oauth2ClientUsecase,
		conf:                conf,
		ZoomTokenUsecase:    ZoomTokenUsecase,
	}
}

func (c *Oauth2TokenUsecase) GetByClientId(clientId string) (*Oauth2TokenEntity, error) {
	var entity Oauth2TokenEntity
	err := c.CommonUsecase.DB().Where("client_id=?", clientId).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, nil
}

func (c *Oauth2TokenUsecase) GetByAppId(appId string) (*Oauth2TokenEntity, error) {
	client, err := c.Oauth2ClientUsecase.GetByAppId(appId)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("Oauth2 Client is nil.")
	}

	var entity Oauth2TokenEntity
	err = c.CommonUsecase.DB().Where("client_id=?", client.ClientId).Take(&entity).Error
	if err == nil {
		return &entity, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, nil
}

func (c *Oauth2TokenUsecase) GetAccessToken(appId string) (accessToken string, err error) {
	oauth2TokenEntity, err := c.GetByAppId(appId)
	if err != nil {
		return "", err
	}
	if oauth2TokenEntity == nil {
		return "", errors.New("oauth2TokenEntity is nil")
	}
	return oauth2TokenEntity.AccessToken, nil
}

func (c *Oauth2TokenUsecase) UpdateByToken(appId string, token *oauth2.Token) error {
	if token == nil {
		return errors.New("Token is nil")
	}
	client, _ := c.Oauth2ClientUsecase.GetByAppId(appId)
	if client == nil {
		return errors.New("Oauth2Client is nil")
	}
	tokenEntity, err := c.GetByClientId(client.ClientId)
	if err != nil {
		return err
	}
	if tokenEntity == nil {
		tokenEntity = &Oauth2TokenEntity{
			ClientId:  client.ClientId,
			CreatedAt: time.Now().Unix(),
		}
	}
	if len(token.RefreshToken) > 0 {
		tokenEntity.RefreshToken = token.RefreshToken
	}
	tokenEntity.AccessToken = token.AccessToken
	tokenEntity.TokenType = token.TokenType
	tokenEntity.ExpiresAt = token.Expiry.Unix()
	tokenEntity.UpdatedAt = time.Now().Unix()
	return c.CommonUsecase.DB().Save(&tokenEntity).Error
}

func (c *Oauth2TokenUsecase) RefreshAccessToken(oauth2ClientEntity *Oauth2ClientEntity, oauth2TokenEntity *Oauth2TokenEntity) error {
	if oauth2ClientEntity == nil {
		return errors.New("oauth2ClientEntity is nil.")
	}
	if oauth2TokenEntity == nil {
		return errors.New("oauth2TokenEntity is nil.")
	}
	var newToken *oauth2.Token
	var err error

	if oauth2ClientEntity.AppId == Oauth2_AppId_zoom {
		newToken, err = c.ZoomTokenUsecase.OauthToken()
		if err != nil {
			return err
		}
	} else {
		config := oauth2ClientEntity.Oauth2Config(c.conf)
		tokenSource := config.TokenSource(context.TODO(), &oauth2.Token{
			RefreshToken: oauth2TokenEntity.RefreshToken,
		})
		newToken, err = tokenSource.Token()
		if err != nil {
			return err
		}
	}
	return c.UpdateByToken(oauth2ClientEntity.AppId, newToken)
}

func (c *Oauth2TokenUsecase) WaitingRefreshToken() (records []*Oauth2TokenEntity, err error) {

	t2 := Oauth2ClientEntity{}.TableName()
	sql, err := builder.MySQL().Select("t.*").
		From(Oauth2TokenEntity{}.TableName(), "t").
		InnerJoin(Oauth2ClientEntity{}.TableName(), fmt.Sprintf("%s.client_id=t.client_id", t2)).
		Where(builder.Eq{t2 + ".deleted_at": 0}).And(builder.Lte{"t.expires_at": time.Now().Unix() + 600}).ToBoundSQL()
	// 提前10分钟结束
	if err != nil {
		return nil, err
	}
	err = c.CommonUsecase.DB().Raw(sql).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return
}

func (c *Oauth2TokenUsecase) RunRefreshTokenJob(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Oauth2TokenUsecase:RunRefreshTokenJob is done.")
				return
			default:
				c.HandleRefreshToken()
				time.Sleep(60 * time.Second)
			}
		}
	}()
	return nil
}

func (c *Oauth2TokenUsecase) HandleRefreshToken() {
	tokens, err := c.WaitingRefreshToken()
	if err != nil {
		c.log.Error(err)
	} else {
		for k, v := range tokens {
			client, err := c.Oauth2ClientUsecase.GetByClientId(v.ClientId)
			if err != nil {
				c.log.Error("GetByClientId: ", err)
			}
			if client != nil {
				err = c.RefreshAccessToken(client, tokens[k])
				if err != nil {
					c.log.Error("HandleRefreshToken: ", v.ClientId, " RefreshAccessToken: ", err, " ClientId: ", client.ClientId, " ID: ", tokens[k].ID)
				}
			}
		}
	}
}
