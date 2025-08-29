package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	ReferrerLog_Type_ReferringClient = "ReferringClient"
)

type ReferrerLogEntity struct {
	ID            int32 `gorm:"primaryKey"`
	Type          string
	Uniqid        string
	ReferrerStage string
	ReferrerValue string
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     int64
}

func (ReferrerLogEntity) TableName() string {
	return "referrer_log"
}

type ReferrerLogUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ReferrerLogEntity]
}

func NewReferrerLogUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ReferrerLogUsecase {
	uc := &ReferrerLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ReferrerLogUsecase) ReferringClient(Uniqid string, ReferrerStage string, ReferrerValue string) error {
	return c.Upsert(ReferrerLog_Type_ReferringClient, Uniqid, ReferrerStage, ReferrerValue)
}

func (c *ReferrerLogUsecase) Upsert(Type string, Uniqid string, ReferrerStage string, ReferrerValue string) error {

	entity, err := c.GetByCond(Eq{
		"type":   Type,
		"uniqid": Uniqid,
	})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &ReferrerLogEntity{
			Type:      Type,
			Uniqid:    Uniqid,
			CreatedAt: time.Now().Unix(),
		}
	}
	entity.UpdatedAt = time.Now().Unix()
	//if ReferrerStage != "" {
	entity.ReferrerStage = ReferrerStage
	//}
	//if ReferrerValue != "" {
	entity.ReferrerValue = ReferrerValue
	//}
	return c.CommonUsecase.DB().Save(&entity).Error
}
