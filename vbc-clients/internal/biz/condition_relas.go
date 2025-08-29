package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type ConditionRelaAiEntity struct {
	ID                   int32 `gorm:"primaryKey"`
	PrimaryConditionId   int32
	SecondaryConditionId int32
	Count                int
	PromptKey            string
	CreatedAt            int64
	UpdatedAt            int64
}

func (ConditionRelaAiEntity) TableName() string {
	return "condition_relas_ai"
}

type ConditionRelaAiUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ConditionRelaAiEntity]
}

func NewConditionRelaAiUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ConditionRelaAiUsecase {
	uc := &ConditionRelaAiUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ConditionRelaAiUsecase) Upsert(PrimaryConditionId int32, SecondaryConditionId int32, promptKey string) error {
	entity, err := c.GetByCond(Eq{"primary_condition_id": PrimaryConditionId,
		"secondary_condition_id": SecondaryConditionId,
		"prompt_key":             promptKey})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &ConditionRelaAiEntity{
			PrimaryConditionId:   PrimaryConditionId,
			SecondaryConditionId: SecondaryConditionId,
			Count:                1,
			CreatedAt:            time.Now().Unix(),
			PromptKey:            promptKey,
		}
	} else {
		entity.Count += 1
	}
	entity.UpdatedAt = time.Now().Unix()
	return c.CommonUsecase.DB().Save(&entity).Error
}
