package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type ZoomUserEntity struct {
	ID                int32 `gorm:"primaryKey"`
	UserId            string
	FirstName         string
	LastName          string
	DisplayName       string
	Email             string
	Type              string
	Pmi               string
	Timezone          string
	Verified          string
	CcreatedAt        string
	LastLoginTime     string
	LastClientVersion string
	PicUrl            string
	Language          string
	PhoneNumber       string
	Status            string
	RoleId            string
	UserCreatedAt     string
	CreatedAt         int64
	UpdatedAt         int64
	DeletedAt         int64
}

func (ZoomUserEntity) TableName() string {
	return "zoom_users"
}

type ZoomUserUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ZoomUserEntity]
}

func NewZoomUserUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ZoomUserUsecase {
	uc := &ZoomUserUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
