package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"regexp"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	VbcAI_Host  = "http://20.190.193.236:8050"
	VbcAi_Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTA4LTI3VDAxOjA4OjIxLjMxMTg2MTc4OFoiLCJyZXFfaWQiOiJhZmM3MjhiMGM1YmM0MTY4OGYzZmYwOTIwNTIwMWVjYiIsInVzZXJfaWQiOjR9.mTDBXkhEpEWwo0Fz8mGg47FmhyKykXWtMgD3eb3k2eE"
)

/*
此方法主要解决，不使用本地IP直接调用 云厂商AI服务的问题
*/
type VbcAIUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
}

func NewVbcAIUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *VbcAIUsecase {
	uc := &VbcAIUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func GetJotformSubmissionIdFromFileName(newFileName string) string {
	re := regexp.MustCompile(`-(\d+)\.pdf`)
	matches := re.FindStringSubmatch(newFileName)
	if len(matches) > 0 {
		return matches[1] // 打印匹配到的数字
	}
	return ""
}

func GetJsonFromAiResult(aiResult string) string {

	re := regexp.MustCompile(`(?s)\{.*?\}`) // (?s) 让 . 匹配换行
	match := re.FindString(aiResult)
	if match != "" {
		return match
	} else {
		return ""
	}
}

func (c *VbcAIUsecase) Claude3(prompt string, text string, fromPromptKey string) (string, int32, error) {

	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        2048,
		Messages: []Claude3Message{
			{
				Role: "user",
				Content: []Claude3Content{
					{
						Type: "text",
						Text: text,
					},
				},
			},
		},
	}
	if prompt != "" {
		payload.SystemPrompt = prompt
	}
	return c.Claude3Invoke(payload, fromPromptKey)

	//url := fmt.Sprintf("%s%s", VbcAI_Host, "/api/ai/claude3")
	//params := make(lib.TypeMap)
	//params.Set("prompt", prompt)
	//params.Set("text", text)
	//res, _, err := lib.Request("POST", url, params.ToBytes(), map[string]string{
	//	"authorization": "Bearer " + VbcAi_Token,
	//})
	//if err != nil {
	//	return "", err
	//}
	//if res == nil {
	//	return "", errors.New("res is nil")
	//}
	//data := lib.ToTypeMapByString(*res)
	//if data.GetInt("code") != 200 {
	//	return "", errors.New(data.GetString("code") + ":" + data.GetString("message"))
	//}
	//return data.GetString("ai_response"), nil
}

func (c *VbcAIUsecase) Claude3Invoke(claude3Request Claude3Request, fromPromptKey string) (aiResponse string, aiResultId int32, err error) {

	url := fmt.Sprintf("%s%s", VbcAI_Host, "/api/ai/claude3")
	lib.DPrintln("url:", url)
	params := make(lib.TypeMap)
	params.Set("ai_request", InterfaceToString(claude3Request))
	params.Set("prompt_key", fromPromptKey)
	res, _, err := lib.Request("POST", url, params.ToBytes(), map[string]string{
		"authorization": "Bearer " + VbcAi_Token,
	})
	if err != nil {
		return "", 0, err
	}
	if res == nil {
		return "", 0, errors.New("res is nil")
	}
	data := lib.ToTypeMapByString(*res)
	if data.GetInt("code") != 200 {
		return "", 0, errors.New(data.GetString("code") + ":" + data.GetString("message"))
	}
	return data.GetString("ai_response"), data.GetInt("ai_result_id"), nil
}
