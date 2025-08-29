package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type BUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	TUsecase           *TUsecase
	FieldUsecase       *FieldUsecase
	FieldOptionUsecase *FieldOptionUsecase
}

func NewBUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase) *BUsecase {
	uc := &BUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		TUsecase:           TUsecase,
		FieldUsecase:       FieldUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
	}

	return uc
}
