package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

type AiTaskResultEntity struct {
	ID         int32 `gorm:"primaryKey"`
	AiTaskId   int32
	AiResultId int32
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  int64
}

func (AiTaskResultEntity) TableName() string {
	return "ai_task_results"
}

type AiTaskResultUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[AiTaskResultEntity]
}

func NewAiTaskResultUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AiTaskResultUsecase {
	uc := &AiTaskResultUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AiTaskResultUsecase) Create(AiTaskId int32, AiResultId int32) error {

	entity := &AiTaskResultEntity{
		AiTaskId:   AiTaskId,
		AiResultId: AiResultId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Save(&entity).Error
}
