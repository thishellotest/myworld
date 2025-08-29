package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

//

type ConditionSourceEntity struct {
	ID         int32 `gorm:"primaryKey"`
	From       string
	Content    string
	CaseId     int32
	ContentMd5 string
	CreatedAt  int64
	UpdatedAt  int64
}

func (ConditionSourceEntity) TableName() string {
	return "condition_sources"
}

type ConditionSourceUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ConditionSourceEntity]
}

func NewConditionSourceUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ConditionSourceUsecase {
	uc := &ConditionSourceUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
