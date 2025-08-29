package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

type AiUsecase struct {
	log               *log.Helper
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	Awsclaude3Usecase *Awsclaude3Usecase
	AiPromptUsecase   *AiPromptUsecase
	AiResultUsecase   *AiResultUsecase
}

func NewAiUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	AiPromptUsecase *AiPromptUsecase,
	AiResultUsecase *AiResultUsecase,
) *AiUsecase {
	uc := &AiUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		Awsclaude3Usecase: Awsclaude3Usecase,
		AiPromptUsecase:   AiPromptUsecase,
		AiResultUsecase:   AiResultUsecase,
	}

	return uc
}

const (
	AnthropicVersion_2023    = "bedrock-2023-05-31"
	Prompt_associate_jotform = "prompt_associate_jotform"
)

func (c *AiUsecase) ExecuteWithClaude3New(promptKey string, dynamicParamsExample lib.TypeMap, claude3Contents []Claude3Content) (parseResult string, aiResultId int32, err error) {
	prompt, text, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, dynamicParamsExample)
	if err != nil {
		return "", aiResultId, err
	}
	if prompt == "" && text == "" {
		return "", aiResultId, errors.New("prompt is wrong")
	}
	var contents []Claude3Content
	if text != "" {
		contents = append(contents, Claude3Content{
			Type: "text",
			Text: text,
		})
	}
	contents = append(contents, claude3Contents...)
	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        2048,
		Messages: []Claude3Message{
			{
				Role:    "user",
				Content: contents,
			},
		},
	}
	if prompt != "" {
		payload.SystemPrompt = prompt
	}
	res, aiResultId, err := c.Awsclaude3Usecase.AskInvoke(context.TODO(), payload, AiFrom_ExecuteWithClaude3New, promptKey, false, nil)
	if err != nil {
		return "", 0, err
	}
	if len(res.ResponseContent) > 0 {
		parseResult = res.ResponseContent[0].Text
	}
	return parseResult, aiResultId, nil
}

/*
ExecutePromptAssociateJotform
```json

	{
	  "related_entries": [
	    "5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf",
	    "5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf"
	  ]
	}
*/
func (c *AiUsecase) ExecutePromptAssociateJotform(condition string, dataList string) (parseResult string, aiResultId int32, err error) {

	return c.ExecuteWithClaude3(Prompt_associate_jotform, lib.TypeMap{
		"condition": condition,
		"data_list": dataList,
	})
}

func (c *AiUsecase) ExecuteWithClaude3(promptKey string, dynamicParamsExample lib.TypeMap) (parseResult string, aiResultId int32, err error) {
	prompt, text, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, dynamicParamsExample)
	if err != nil {
		return "", aiResultId, err
	}
	if prompt == "" && text == "" {
		return "", aiResultId, errors.New("prompt is wrong")
	}

	res, aiResultId, err := c.Awsclaude3Usecase.Ask(context.TODO(), prompt, text, AiFrom_ExecuteWithClaude3, promptKey)
	if err != nil {
		c.log.Warn(err)
		if strings.Index(err.Error(), "Too many tokens, please wait before trying again") >= 0 {
			time.Sleep(time.Second * 5)
			c.log.Debug("ExecuteWithClaude3: sleep 5 seconds")
		}
		return "", 0, err
	}
	//c.log.Info("AiHttpUsecase err:", err)
	//c.log.Info("AiHttpUsecase:", res)
	if len(res.ResponseContent) > 0 {
		parseResult = res.ResponseContent[0].Text
	}

	//entity := &AiResultEntity{
	//	Prompt:      prompt,
	//	Text:        text,
	//	Result:      InterfaceToString(res),
	//	ParseResult: parseResult,
	//	ErrResult:   InterfaceToString(err),
	//	AiFrom:      "ExecuteWithClaude3:" + promptKey,
	//	CreatedAt:   time.Now().Unix(),
	//	UpdatedAt:   time.Now().Unix(),
	//}
	//
	//err = c.AiResultUsecase.CommonUsecase.DB().Save(&entity).Error
	//if err != nil {
	//	c.log.Warn(err)
	//}
	return parseResult, aiResultId, nil
}
