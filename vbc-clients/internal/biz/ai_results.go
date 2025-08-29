package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

type AiResultEntity struct {
	ID            int32 `gorm:"primaryKey"`
	ModelId       string
	FromPromptKey string
	Prompt        string
	Text          string
	AiRequest     string
	Result        string
	ParseResult   string
	ErrResult     string
	AiFrom        string
	ApiDuration   int
	CreatedAt     int64
	UpdatedAt     int64
}

func (c *AiResultEntity) GetStatement() string {
	return c.ParseResult
	//return StatementExtract(c.ParseResult)
}

func (AiResultEntity) TableName() string {
	return "ai_results"
}

type AiResultUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[AiResultEntity]
}

func NewAiResultUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AiResultUsecase {
	uc := &AiResultUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
