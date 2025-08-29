package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	UniqueCodeGenerator_Type_ClientUniqCode = "ClientUniqCode"
)

type UniqueCodeGeneratorEntity struct {
	ID        int32 `gorm:"primaryKey"`
	Type      string
	Uuid      string
	CreatedAt int64
}

func (UniqueCodeGeneratorEntity) TableName() string {
	return "unique_code_generator"
}

type UniqueCodeGeneratorUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[UniqueCodeGeneratorEntity]
}

func NewUniqueCodeGeneratorUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *UniqueCodeGeneratorUsecase {
	uc := &UniqueCodeGeneratorUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *UniqueCodeGeneratorUsecase) GenUuid(typ string, loopCount int) (string, error) {
	if typ == UniqueCodeGenerator_Type_ClientUniqCode {
		uuid := lib.UuidNumeric()
		err := c.CommonUsecase.DB().Create(&UniqueCodeGeneratorEntity{
			Type:      typ,
			Uuid:      uuid,
			CreatedAt: time.Now().Unix(),
		}).Error
		if err != nil {
			loopCount++
			if loopCount > 20 {
				return "", err
			} else {
				return c.GenUuid(typ, loopCount)
			}
		} else {
			return uuid, nil
		}
	}
	return "", nil
}
