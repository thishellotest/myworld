package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

type LogInfoEntity struct {
	ID        int32 `gorm:"primaryKey"`
	FromId    int32
	FromType  string
	Notes     string
	CreatedAt int64
}

func (LogInfoEntity) TableName() string {
	return "log_infos"
}

func GenLogInfo(fromId int32, fromType string, notes string) *LogInfoEntity {
	return &LogInfoEntity{
		FromId:    fromId,
		FromType:  fromType,
		Notes:     notes,
		CreatedAt: time.Now().Unix(),
	}
}

type LogInfoUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewLogInfoUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *LogInfoUsecase {
	uc := &LogInfoUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	return uc
}

func (c *LogInfoUsecase) SaveLogInfo(fromId int32, fromType string, notes interface{}) error {
	log := GenLogInfo(fromId, fromType, InterfaceToString(notes))
	return c.CommonUsecase.DB().Save(log).Error
}
