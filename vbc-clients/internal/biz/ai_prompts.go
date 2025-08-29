package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AiPromptEntity struct {
	ID                   int32 `gorm:"primaryKey"`
	PromptKey            string
	Notes                string
	DynamicParamsExample string
	Prompt               string
	Text                 string
	Example              string
	CreatedAt            int64
	UpdatedAt            int64
}

func (AiPromptEntity) TableName() string {
	return "ai_prompts"
}

func (c *AiPromptEntity) AiInfo(dynamicParamsExample lib.TypeMap) (prompt, text string) {

	prompt = c.Prompt
	text = c.Text
	for k, v := range dynamicParamsExample {
		prompt = strings.ReplaceAll(prompt, fmt.Sprintf("{{%s}}", k), InterfaceToString(v))
		text = strings.ReplaceAll(text, fmt.Sprintf("{{%s}}", k), InterfaceToString(v))
	}
	return prompt, text
}

type AiPromptUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[AiPromptEntity]
}

func NewAiPromptUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AiPromptUsecase {
	uc := &AiPromptUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *AiPromptUsecase) GetByPromptKey(promptKey string) (AiPromptEntity *AiPromptEntity, err error) {
	return c.GetByCond(Eq{"prompt_key": promptKey})
}

func (c *AiPromptUsecase) GetAiInfoByPromptKey(promptKey string, dynamicParamsExample lib.TypeMap) (prompt string, text string, err error) {
	aiPromptEntity, err := c.GetByPromptKey(promptKey)
	if err != nil {
		return "", "", err
	}
	if aiPromptEntity == nil {
		return "", "", nil
	}
	prompt, text = aiPromptEntity.AiInfo(dynamicParamsExample)
	return
}
