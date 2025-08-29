package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

const Default_Timezones_CodeValue = "America/Los_Angeles"

type TimezonesEntity struct {
	ID        int32 `gorm:"primaryKey"`
	CodeValue string
	Title     string
}

func (TimezonesEntity) TableName() string {
	return "timezones"
}

type TimezonesUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TimezonesEntity]
}

func NewTimezonesUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *TimezonesUsecase {
	uc := &TimezonesUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}
