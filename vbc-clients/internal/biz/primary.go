package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

type PrimaryUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	FeeUsecase    *FeeUsecase
}

func NewPrimaryUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	FeeUsecase *FeeUsecase,
) *PrimaryUsecase {
	uc := &PrimaryUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		FeeUsecase:    FeeUsecase,
	}
	return uc
}

// GetPrimaryCase primaryCase 一定有值，isPrimaryCase： ture-说明tCase就是Primary Case
func (c *PrimaryUsecase) GetPrimaryCase(tCase *TData) (primaryCase *TData, isPrimaryCase bool, err error) {
	if tCase == nil {
		return nil, false, errors.New("tCase is nil")
	}
	isP, PCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return nil, false, err
	}
	if isP {
		return tCase, true, nil
	} else {
		return PCase, false, nil
	}
}
