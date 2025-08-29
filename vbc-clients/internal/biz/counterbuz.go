package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type CounterbuzUsecase struct {
	log            *log.Helper
	CommonUsecase  *CommonUsecase
	conf           *conf.Data
	CounterUsecase *CounterUsecase
}

func NewCounterbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	CounterUsecase *CounterUsecase) *CounterbuzUsecase {
	uc := &CounterbuzUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		CounterUsecase: CounterUsecase,
	}

	return uc
}

func (c *CounterbuzUsecase) ClientUploadFilesLimit(caseId int32, maxCount int) (hasLimit bool, err error) {
	key := CounterKeyClientUploadFiles(caseId)
	return c.CounterUsecase.HasLimit(key, maxCount)
}

func (c *CounterbuzUsecase) ClientUploadFilesStat(caseId int32) error {
	key := CounterKeyClientUploadFiles(caseId)
	return c.CounterUsecase.Stat(key, 1)
}
