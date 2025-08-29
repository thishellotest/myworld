package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type TimezoneUsecase struct {
	log            *log.Helper
	CommonUsecase  *CommonUsecase
	conf           *conf.Data
	GoCacheUsecase *GoCacheUsecase
}

func NewTimezoneUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	GoCacheUsecase *GoCacheUsecase) *TimezoneUsecase {
	uc := &TimezoneUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		GoCacheUsecase: GoCacheUsecase,
	}

	return uc
}
