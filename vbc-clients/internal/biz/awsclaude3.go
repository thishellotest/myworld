package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"sort"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type Awsclaude3Usecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	JotformSubmissionUsecase *JotformSubmissionUsecase
	VbcAIUsecase             *VbcAIUsecase
	AiResultUsecase          *AiResultUsecase
	AiPromptUsecase          *AiPromptUsecase
	QuestionnairesbuzUsecase *QuestionnairesbuzUsecase
	ConditionUsecase         *ConditionUsecase
	CacheLogUsecase          *CacheLogUsecase
	FeeUsecase               *FeeUsecase
}

func NewAwsclaude3Usecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	VbcAIUsecase *VbcAIUsecase,
	AiResultUsecase *AiResultUsecase,
	AiPromptUsecase *AiPromptUsecase,
	QuestionnairesbuzUsecase *QuestionnairesbuzUsecase,
	ConditionUsecase *ConditionUsecase,
	CacheLogUsecase *CacheLogUsecase,
	FeeUsecase *FeeUsecase) *Awsclaude3Usecase {
	uc := &Awsclaude3Usecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		VbcAIUsecase:             VbcAIUsecase,
		AiResultUsecase:          AiResultUsecase,
		AiPromptUsecase:          AiPromptUsecase,
		QuestionnairesbuzUsecase: QuestionnairesbuzUsecase,
		ConditionUsecase:         ConditionUsecase,
		CacheLogUsecase:          CacheLogUsecase,
		FeeUsecase:               FeeUsecase,
	}

	return uc
}

func (c *Awsclaude3Usecase) BedrockRuntimeClient() (*bedrockruntime.Client, error) {

	region := "us-east-1"
	var cfg aws.Config
	var err error
	if configs.IsProd() {
		//cfg, err = config.LoadDefaultConfig(context.Background(),
		//	config.WithRegion(region),
		//	//config.WithSharedConfigProfile("config"),
		//	config.WithSharedCredentialsFiles([]string{
		//		"/app/.sck/credentials",
		//	}))
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region))
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region))
	}
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	brc := bedrockruntime.NewFromConfig(cfg)
	return brc, nil
}

type Claude3Response struct {
	ID              string                   `json:"id,omitempty"`
	Model           string                   `json:"model,omitempty"`
	Type            string                   `json:"type,omitempty"`
	Role            string                   `json:"role,omitempty"`
	ResponseContent []Claude3ResponseContent `json:"content,omitempty"`
	StopReason      string                   `json:"stop_reason,omitempty"`
	StopSequence    string                   `json:"stop_sequence,omitempty"`
	Usage           Claude3Usage             `json:"usage,omitempty"`
}

type Claude3Request struct {
	AnthropicVersion string           `json:"anthropic_version"`
	MaxTokens        int              `json:"max_tokens"`
	Messages         []Claude3Message `json:"messages"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
	TopK             int              `json:"top_k,omitempty"`
	StopSequences    []string         `json:"stop_sequences,omitempty"`
	SystemPrompt     string           `json:"system,omitempty"`
}

func (c *Claude3Request) GetCacheKey(aiFrom string) string {
	str := fmt.Sprintf("%s:%s:%s", aiFrom, InterfaceToString(c.Messages), c.SystemPrompt)
	return lib.MD5Hash(str)
}

type Claude3Message struct {
	Role    string           `json:"role,omitempty"`
	Content []Claude3Content `json:"content,omitempty"`
}

type Claude3Content struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type Claude3ResponseContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
type Claude3Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

const AiFrom_GetContentByAsk = "GetContentByAsk"
const AiFrom_AskV2 = "AskV2"
const AiFrom_ExecuteWithClaude3New = "ExecuteWithClaude3New"
const AiFrom_ExecuteWithClaude3 = "ExecuteWithClaude3"
const AiFrom_BizClaude3 = "BizClaude3"
const AiForm_GenStatement = "GenStatement"
const AiForm_GenUpdateStatement = "GenUpdateStatement"
const AiForm_GenVeteranSummary = "GenVeteranSummary"
const AiForm_ReturnTimeZone = "ReturnTimeZone"

func (c *Awsclaude3Usecase) GetContentByAsk(ctx context.Context, systemConfig string, medicalText string) (string, error) {

	if configs.IsDev() {
		return `Here is the JSON format for the conditions found in the given text:

{
  "Conditions": [
    "HEADACHE",
    "RUNNY NOSE",
    "SORE THROAT",
    "Nausea",
    "vomiting",
    "loss of Appetite",
    "viral gastroenteritis"
  ]
}`, nil
	}

	claude3Response, _, err := c.Ask(ctx, systemConfig, medicalText, AiFrom_GetContentByAsk, "")
	if err != nil {
		return "", err
	}
	for _, v := range claude3Response.ResponseContent {
		return v.Text, nil
	}
	return "", errors.New("ResponseContent is empty")
}

func (c *Awsclaude3Usecase) Ask(ctx context.Context, systemPrompt string, medicalText string, aiFrom string, promptKey string) (claude3Response Claude3Response, aiResultId int32, err error) {

	//modelID := "anthropic.claude-3-sonnet-20240229-v1:0"
	//
	//brc, err := c.BedrockRuntimeClient()
	//if err != nil {
	//	return claude3Response, 0, err
	//}

	//msg := "Hello, what's your name?"

	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        2048,
		Messages: []Claude3Message{
			{
				Role: "user",
				Content: []Claude3Content{
					{
						Type: "text",
						Text: medicalText,
					},
				},
			},
		},
	}
	if systemPrompt != "" {
		payload.SystemPrompt = systemPrompt
	}

	return c.AskInvoke(ctx, payload, aiFrom, promptKey, false, nil)
	//if !lib.IsProd() {
	//	claude3ResponseString, err := c.VbcAIUsecase.Claude3(systemPrompt, medicalText)
	//	if err != nil {
	//		return claude3Response, err
	//	}
	//	err = json.Unmarshal([]byte(claude3ResponseString), &claude3Response)
	//	if err != nil {
	//		return claude3Response, err
	//	}
	//	return claude3Response, nil
	//}
	//
	//payloadBytes, err := json.Marshal(payload)
	//if err != nil {
	//	c.log.Error("Awsclaude3Usecase Marshal error:", err)
	//	return claude3Response, err
	//}
	//
	//output, err := brc.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
	//	Body:        payloadBytes,
	//	ModelId:     aws.String(modelID),
	//	ContentType: aws.String("application/json"),
	//})
	//
	//if err != nil {
	//	c.log.Error("Awsclaude3Usecase InvokeModel error:", err)
	//	return claude3Response, err
	//}
	////lib.DPrintln(output.Body)
	//
	//var resp Claude3Response
	//
	//err = json.Unmarshal(output.Body, &resp)
	//
	//if err != nil {
	//	c.log.Info("Awsclaude3Usecase Unmarshal error:", err)
	//	return claude3Response, err
	//}
	//
	////fmt.Println("Awsclaude3Usecase response payload: ", string(output.Body))
	////fmt.Println("Awsclaude3Usecase response string: ", resp.ResponseContent[0].Text)
	//
	//return resp, nil
}

const (
	Aws_AiModel_claude3         = "anthropic.claude-3-sonnet-20240229-v1:0"
	Aws_AiModel_claude3_5       = "anthropic.claude-3-5-sonnet-20240620-v1:0"
	Aws_AiModel_claude3_5_haiku = "anthropic.claude-3-5-haiku-20241022-v1:0"  // 不能调用
	Aws_AiModel_claude3_7       = "anthropic.claude-3-7-sonnet-20250219-v1:0" // 不能调用 https://docs.aws.amazon.com/bedrock/latest/userguide/inference-profiles-support.html
	Aws_AiModel_us_claude3_7    = "us.anthropic.claude-3-7-sonnet-20250219-v1:0"
)

// AskInvoke aiFrom+payload 当enableCache启用时，相同的参数会从缓存里面找
func (c *Awsclaude3Usecase) AskInvoke(ctx context.Context, payload Claude3Request, aiFrom string, promptKey string, enableCache bool, isResultOkForEnableCache func(parseResult string) bool) (claude3Response Claude3Response, aiResultId int32, err error) {

	modelID := Aws_AiModel_us_claude3_7
	beginTime := time.Now().Unix()

	brc, err := c.BedrockRuntimeClient()
	if err != nil {
		return claude3Response, aiResultId, err
	}

	//msg := "Hello, what's your name?"

	//payload := Claude3Request{
	//	AnthropicVersion: "bedrock-2023-05-31",
	//	MaxTokens:        1024,
	//	Messages: []Claude3Message{
	//		{
	//			Role: "user",
	//			Content: []Claude3Content{
	//				{
	//					Type: "text",
	//					Text: medicalText,
	//				},
	//			},
	//		},
	//	},
	//}
	//if systemPrompt != "" {
	//	payload.SystemPrompt = systemPrompt
	//}
	cacheKey := payload.GetCacheKey(aiFrom)
	if enableCache {
		c.log.Info("cacheKey:", cacheKey)
		cacheLogEntity, _ := c.CacheLogUsecase.GetForAiResult(cacheKey)
		if cacheLogEntity != nil {
			aiResultEntity, _ := c.AiResultUsecase.GetByCond(Eq{"id": cacheLogEntity.ResultId})
			if aiResultEntity != nil {
				var resp Claude3Response
				err = json.Unmarshal([]byte(aiResultEntity.Result), &resp)
				if err != nil {
					return claude3Response, aiResultId, err
				}
				return resp, aiResultEntity.ID, nil
			}
		}
	}

	if !configs.IsProd() {
		claude3ResponseString, aiResultId, err := c.VbcAIUsecase.Claude3Invoke(payload, promptKey)
		if err != nil {
			return claude3Response, aiResultId, err
		}
		err = json.Unmarshal([]byte(claude3ResponseString), &claude3Response)
		if err != nil {
			return claude3Response, aiResultId, err
		}
		return claude3Response, aiResultId, nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.log.Error("Awsclaude3Usecase Marshal error:", err)
		return claude3Response, 0, err
	}

	output, err := brc.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
	})

	endTime := time.Now().Unix()
	apiDuration := int(endTime - beginTime)
	if err != nil {
		//c.log.Error("Awsclaude3Usecase InvokeModel error:", err)
		entity := &AiResultEntity{
			ModelId:       modelID,
			FromPromptKey: promptKey,
			//Prompt:      systemPrompt,
			//Text:        text,
			AiRequest:   InterfaceToString(payload),
			Result:      InterfaceToString(output),
			ParseResult: "",
			ErrResult:   InterfaceToString(err),
			AiFrom:      aiFrom,
			ApiDuration: apiDuration,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}
		er := c.AiResultUsecase.CommonUsecase.DB().Save(&entity).Error
		if er != nil {
			c.log.Error(er)
		}
		return claude3Response, aiResultId, err
	}
	//lib.DPrintln(output.Body)
	var resp Claude3Response
	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		c.log.Error("Awsclaude3Usecase Unmarshal error:", err)

		entity := &AiResultEntity{
			ModelId:       modelID,
			FromPromptKey: promptKey,
			//Prompt:      systemPrompt,
			//Text:        text,
			AiRequest:   InterfaceToString(payload),
			Result:      InterfaceToString(output),
			ParseResult: "",
			ErrResult:   InterfaceToString(err),
			AiFrom:      aiFrom,
			ApiDuration: apiDuration,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}

		er := c.AiResultUsecase.CommonUsecase.DB().Save(&entity).Error
		if er != nil {
			c.log.Error(er)
		}

		return claude3Response, aiResultId, err
	}

	//fmt.Println("Awsclaude3Usecase response payload: ", string(output.Body))
	//fmt.Println("Awsclaude3Usecase response string: ", resp.ResponseContent[0].Text)

	parseResult := ""
	if len(resp.ResponseContent) > 0 {
		parseResult = resp.ResponseContent[0].Text
	}

	entity := &AiResultEntity{
		ModelId:       modelID,
		FromPromptKey: promptKey,
		//Prompt:      systemPrompt,
		//Text:        text,
		AiRequest:   InterfaceToString(payload),
		Result:      InterfaceToString(resp),
		ParseResult: parseResult,
		ErrResult:   InterfaceToString(err),
		AiFrom:      aiFrom,
		ApiDuration: apiDuration,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = c.AiResultUsecase.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return
		c.log.Warn(err)
	}

	if enableCache {
		err = c.CacheLogUsecase.AddForAiResult(cacheKey, InterfaceToString(entity.ID))
		if err != nil {
			c.log.Error(err)
		}
	}

	return resp, entity.ID, nil
}

func FormatJotformAnswersForGenStatement(notes lib.TypeMap) string {
	answers := notes.GetTypeMap("content.answers")
	//lib.DPrintln(answers)

	var questions lib.TypeList
	for _, v := range answers {

		fieldInfo := lib.ToTypeMap(v)
		//row := make(lib.TypeMap)
		//row.Set("order", fieldInfo.GetString("order"))
		//row.Set("question", fieldInfo.GetString("text"))
		//row.Set("answer", fieldInfo.GetString("answer"))
		//row.Set("type", fieldInfo.GetString("type"))

		questions = append(questions, fieldInfo)
	}

	sort.SliceStable(questions, func(i, j int) bool {
		return questions[i].GetInt("order") < questions[j].GetInt("order")
	})

	//lib.DPrintln(questions)

	return InterfaceToString(questions)
}

func FormatJotformAnswers(notes lib.TypeMap) string {
	answers := notes.GetTypeMap("content.answers")
	//lib.DPrintln(answers)

	var questions lib.TypeList
	for _, v := range answers {

		fieldInfo := lib.ToTypeMap(v)
		//row := make(lib.TypeMap)
		//row.Set("order", fieldInfo.GetString("order"))
		//row.Set("question", fieldInfo.GetString("text"))
		//row.Set("answer", fieldInfo.GetString("answer"))
		//row.Set("type", fieldInfo.GetString("type"))

		questions = append(questions, fieldInfo)
	}

	sort.SliceStable(questions, func(i, j int) bool {
		return questions[i].GetInt("order") < questions[j].GetInt("order")
	})

	//lib.DPrintln(questions)

	return InterfaceToString(questions)
}

func (c *Awsclaude3Usecase) GenVeteranSummary(ctx context.Context, tCase *TData, promptKey string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}
	intakeSubmission, err := c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if intakeSubmission == nil {
		return Claude3Response{}, "", 0, errors.New("intakeSubmission is nil")
	}
	var contents []Claude3Content
	var res []*JotformSubmissionEntity
	res = append(res, intakeSubmission)
	for _, v := range res {
		notes := lib.ToTypeMapByString(v.Notes)
		str := FormatJotformAnswers(notes)
		contents = append(contents, Claude3Content{
			Type: "text",
			Text: str,
		})
	}
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

	claude3Response, aiResultId, err = c.AskInvoke(ctx, payload, AiForm_GenVeteranSummary, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

func (c *Awsclaude3Usecase) GenStatementOld(ctx context.Context, tCase *TData, StatementCondition StatementCondition, veteranSummary string, promptKey string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var contents []Claude3Content

	var intakeSubmission *JotformSubmissionEntity
	var otherSubmissions []*JotformSubmissionEntity

	if configs.NewPSGen {
		intakeSubmission, otherSubmissions, err = c.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatementNew(tCase, StatementCondition)
		if err != nil {
			return Claude3Response{}, "", 0, err
		}
	} else {
		conditionEntity, err := c.ConditionUsecase.ConditionUpsert(StatementCondition.ConditionValue)
		if err != nil {
			return Claude3Response{}, "", 0, err
		}
		if conditionEntity == nil {
			return Claude3Response{}, "", 0, errors.New("conditionEntity is nil")
		}
		intakeSubmission, otherSubmissions, err = c.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatement(tCase, conditionEntity)
		if err != nil {
			return Claude3Response{}, "", 0, err
		}
	}

	var res []*JotformSubmissionEntity
	res = append(res, intakeSubmission)
	c.log.Info(len(otherSubmissions))
	res = append(res, otherSubmissions...)
	for _, v := range res {
		notes := lib.ToTypeMapByString(v.Notes)
		str := FormatJotformAnswers(notes)
		contents = append(contents, Claude3Content{
			Type: "text",
			Text: str,
		})
	}

	if tCase.Id() == 5431 && false {
		//	contents = append(contents, Claude3Content{
		//		Type: "text",
		//		Text: `- Full Name: MarkDean Ronduen
		//- Unique ID: 5431
		//- Branch of Service: Marine Corps
		//- Years of Service: 2014-2018
		//- Retirement Status: Did not retire from service
		//- Deployments: None listed on presumptive lists
		//- Marital Status: Not married, never divorced
		//- Children: None
		//- Occupation in Service: 3051 Inventory Management Specialist`,
		//	})
		/*
			• Full Name: MarkDean Ronduen
			• Unique ID: 5431
			• Branch of Service: Marine Corps
			• Years of Service: 2014-2018
			• Retirement status: Did not retire from service
			• Deployments: No deployments to conflict areas on presumptive lists
			• Marital Status: Not married, never divorced
			• Children: No children
			• Occupation in service: 3051 Inventory Management Specialist
		*/
	}

	//return
	//contents = append(contents, Claude3Content{
	//	Type: "text",
	//	Text: StatementCondition.OriginValue,
	//})
	//contents = append(contents, Claude3Content{
	//	Type: "text",
	//	Text: "",
	//})

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

	claude3Response, aiResultId, err = c.AskInvoke(ctx, payload, AiForm_GenStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

func (c *Awsclaude3Usecase) GenStatementTest(ctx context.Context, text string, promptKey string) (claude3Response Claude3Response, aiResultId int32, err error) {

	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, 0, err
	}

	//modelID := "anthropic.claude-3-sonnet-20240229-v1:0"
	//
	//brc, err := c.BedrockRuntimeClient()
	//if err != nil {
	//	return claude3Response, 0, err
	//}

	//msg := "Hello, what's your name?"

	var contents []Claude3Content

	res, _ := c.JotformSubmissionUsecase.AllByCond(In("submission_id", "6036951604218420113",
		"6036971094217009889",
		"6036973064211091711",
		"6036977804211490647"))

	for _, v := range res {
		notes := lib.ToTypeMapByString(v.Notes)
		str := FormatJotformAnswers(notes)
		contents = append(contents, Claude3Content{
			Type: "text",
			Text: str,
		})
	}
	//return
	contents = append(contents, Claude3Content{
		Type: "text",
		Text: "Migraine headaches secondary to tinnitus",
	})

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

	return c.AskInvoke(ctx, payload, AiFrom_AskV2, promptKey, false, nil)
}
