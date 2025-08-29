package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type MulticasesUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewMulticasesUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *MulticasesUsecase {
	uc := &MulticasesUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}
