package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/uuid"
)

var AppJWTSecret = []byte("71Xbzsa%ZWQM8IyxfMRt$kBXBCKd2oE$")

type AppTokenEntity struct {
	ID           int32 `gorm:"primaryKey"`
	IncrId       int32
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    int64
}

func (AppTokenEntity) TableName() string {
	return "app_tokens"
}

type AppTokenUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[AppTokenEntity]
}

func NewAppTokenUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *AppTokenUsecase {
	uc := &AppTokenUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AppTokenUsecase) CreateToken(incrId int32) (appTokenEntity *AppTokenEntity, err error) {

	exp := time.Now().AddDate(0, 0, 30)

	requestId := uuid.UuidWithoutStrike()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": incrId,
		"req_id":  requestId,
		"exp":     exp,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(AppJWTSecret)
	if err != nil {
		return nil, err
	}

	appTokenEntity = &AppTokenEntity{
		IncrId:      incrId,
		AccessToken: tokenString,
		ExpiresAt:   exp.Unix(),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = c.CommonUsecase.DB().Create(&appTokenEntity).Error
	if err != nil {
		return nil, err
	}
	return
}

func (c *AppTokenUsecase) Parse(accessToken string) (userId int32, isExp bool, err error) {

	//lib.DPrintln("accessToken:", accessToken)
	token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return AppJWTSecret, nil
	})
	if err != nil {
		return 0, false, err
	}
	if token == nil {
		return 0, false, errors.New("token is nil")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		uId := claims["user_id"]
		exp := claims["exp"]
		if uId == nil || exp == nil {
			c.log.Error("userId or exp is nil")
			return 0, false, errors.New("userId or exp is nil")
		}
		tUserId, _ := strconv.ParseInt(InterfaceToString(uId), 10, 32)
		if tUserId <= 0 {
			return 0, false, errors.New("tUserId is 0")
		}
		expTime, err := time.Parse(time.RFC3339Nano, InterfaceToString(exp))
		if err != nil {
			return 0, false, err
		}
		if expTime.Before(time.Now()) {
			return 0, true, nil
		}
		return int32(tUserId), false, nil
	} else {
		return 0, false, errors.New("jwt.MapClaims is wrong")
	}
}

func (c *AppTokenUsecase) ParseNoDB(accessToken string) error {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return AppJWTSecret, nil
	})
	if err != nil {
		c.log.Error(err)
		//return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			c.log.Error(err)
			//return err
		}

		lib.DPrintln(claims["user_id"], claims["exp"], expirationTime)
	} else {
		return errors.New("jwt.MapClaims is wrong")
	}
	return nil
}
