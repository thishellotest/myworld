package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type ResourceUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewResourceUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ResourceUsecase {
	uc := &ResourceUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func (c *ResourceUsecase) ResPath() (resPath string) {
	resPath = c.conf.ResourcePath
	return
}
