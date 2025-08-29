package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"regexp"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AssistantUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	StatementUsecase          *StatementUsecase
	AiPromptUsecase           *AiPromptUsecase
	Awsclaude3Usecase         *Awsclaude3Usecase
	TUsecase                  *TUsecase
	DataComboUsecase          *DataComboUsecase
	StatementConditionUsecase *StatementConditionUsecase
	AiAssistantJobUsecase     *AiAssistantJobUsecase
	AiTaskbuzUsecase          *AiTaskbuzUsecase
	FeeUsecase                *FeeUsecase
	JotformSubmissionUsecase  *JotformSubmissionUsecase
	AiUsecase                 *AiUsecase
	AiTaskUsecase             *AiTaskUsecase
	AiResultUsecase           *AiResultUsecase
}

func NewAssistantUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	StatementUsecase *StatementUsecase,
	AiPromptUsecase *AiPromptUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	AiAssistantJobUsecase *AiAssistantJobUsecase,
	AiTaskbuzUsecase *AiTaskbuzUsecase,
	FeeUsecase *FeeUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	AiUsecase *AiUsecase,
	AiTaskUsecase *AiTaskUsecase,
	AiResultUsecase *AiResultUsecase,
) *AssistantUsecase {
	uc := &AssistantUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		StatementUsecase:          StatementUsecase,
		AiPromptUsecase:           AiPromptUsecase,
		Awsclaude3Usecase:         Awsclaude3Usecase,
		TUsecase:                  TUsecase,
		DataComboUsecase:          DataComboUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
		AiAssistantJobUsecase:     AiAssistantJobUsecase,
		AiTaskbuzUsecase:          AiTaskbuzUsecase,
		FeeUsecase:                FeeUsecase,
		JotformSubmissionUsecase:  JotformSubmissionUsecase,
		AiUsecase:                 AiUsecase,
		AiTaskUsecase:             AiTaskUsecase,
		AiResultUsecase:           AiResultUsecase,
	}

	return uc
}

func GetJsonFromAiResultForAssistant(aiResult string) string {

	re := regexp.MustCompile(`(?s)\{.*?\}`) // (?s) 让 . 匹配换行
	match := re.FindString(aiResult)
	if match != "" {
		return match
	} else {
		return ""
	}
}

func (c *AssistantUsecase) ExplainFormatJobUuidForStatementCondition(jobUuid string) (tCase *TData, tClient *TData, statementConditionEntity *StatementConditionEntity, vo JobUuidForStatementCondition) {

	vo = FormatJobUuidForStatementCondition(jobUuid)
	tCase, _ = c.TUsecase.DataById(Kind_client_cases, vo.CaseId)
	if tCase != nil {
		tClient, _, _ = c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	}
	statementConditionEntity, _ = c.StatementConditionUsecase.GetByCond(Eq{"id": vo.StatementConditionId})

	return
}

func (c *AssistantUsecase) ExplainJobUuidForStatementSection(jobUuid string) (tCase *TData, tClient *TData, statementConditionEntity *StatementConditionEntity, vo JobUuidForStatementSection) {

	vo = FormatJobUuidForStatementSection(jobUuid)
	tCase, _ = c.TUsecase.DataById(Kind_client_cases, vo.CaseId)
	if tCase != nil {
		tClient, _, _ = c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	}
	statementConditionEntity, _ = c.StatementConditionUsecase.GetByCond(Eq{"id": vo.StatementConditionId})

	return
}

func (c *AssistantUsecase) HandleAssistant(ctx context.Context, task *AiTaskEntity) error {
	if task == nil {
		return errors.New("task is nil")
	}
	return c.DoHandleAssistant(ctx, task)
}

func (c *AssistantUsecase) HandleSaveAllStatements(tCase TData, tClient TData) error {

	c.log.Info("invokeSvStatemetns AllStatements", "HandleSaveAllStatements:", tCase.Id())
	statementConditions, err := c.StatementConditionUsecase.AllConditions(tCase.Id())
	if err != nil {
		return err
	}
	var jobUuids []string
	for _, v := range statementConditions {
		jobUuid := v.ToStatementConditionJobUuid()
		jobUuids = append(jobUuids, jobUuid)
	}
	//lib.DPrintln(jobUuids)
	assistants, err := c.AiAssistantJobUsecase.AllByCond(In("job_uuid", jobUuids))
	if err != nil {
		return err
	}
	isDone := true
	for _, v := range assistants {
		c.log.Info("invokeSvStatemetns ", " JobUuid:", v.JobUuid, " v.JobStatus: ", v.JobStatus)
		if v.JobStatus == AiAssistantJob_JobStatus_Running {
			isDone = false
		}
	}
	c.log.Info("invokeSvStatemetns AllStatements", "HandleSaveAllStatements:", tCase.Id(), " isDone: ", isDone)
	if isDone { // 需要更新statements
		er := c.StatementUsecase.GenerateNewStatementVersionForWebForm(tCase, tClient)
		if er != nil {
			c.log.Error("GenerateNewStatementVersion:", er, " caseId: ", tCase.Id())
		}
		//}
		//}
		err = c.StatementUsecase.GenerateDocument(tCase, tClient)
		if err != nil {
			c.log.Error(err)
		}

		jobUuid := GenAllStatementsJobUuid(tCase.Id())

		aiAssistantJobEntity, er := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
		if er != nil {
			c.log.Error(er)
		}
		if aiAssistantJobEntity != nil {
			er = c.AiAssistantJobUsecase.UpdatesByCond(map[string]interface{}{
				"job_status": AiAssistantJob_JobStatus_Normal,
			}, Eq{"id": aiAssistantJobEntity.ID})
			if er != nil {
				c.log.Error(er)
			}
		}

		return err
	}
	return nil
}

func (c *AssistantUsecase) DoHandleAssistant(ctx context.Context, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("taskEntity is nil")
	}
	aiTaskInputAssistant := task.GetAiTaskInputAssistant()
	var tClient *TData
	var tCase *TData
	var err error
	tCase, err = c.TUsecase.DataById(Kind_client_cases, aiTaskInputAssistant.CaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err = c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}

	if aiTaskInputAssistant.AssistantBiz == AiAssistantJobBizType_statementCondition {
		if aiTaskInputAssistant.BizType == AiAssistantBizType_CreateNewStatement {
			statementConditionEntity, err := c.StatementConditionUsecase.GetByCond(Eq{"id": aiTaskInputAssistant.StatementConditionId})
			if err != nil {
				return err
			}
			if statementConditionEntity == nil {
				return errors.New("statementConditionEntity is nil")
			}
			_, err = c.HandleGenStatementFromAiTask(ctx, tCase, task, *statementConditionEntity, aiTaskInputAssistant.UserInputText)
			if err != nil {
				c.log.Error(err)
				return err
			}
		} else if aiTaskInputAssistant.BizType == AiAssistantBizType_StandardHeaderRevision {
			statementConditionEntity, err := c.StatementConditionUsecase.GetByCond(Eq{"id": aiTaskInputAssistant.StatementConditionId})
			if err != nil {
				return err
			}
			if statementConditionEntity == nil {
				return errors.New("statementConditionEntity is nil")
			}
			_, err = c.HandleGenStatementForStandardHeaderRevisionFromAiTask(ctx, *tClient, *tCase, task, *statementConditionEntity, aiTaskInputAssistant.UserInputText)
			if err != nil {
				c.log.Error(err)
				return err
			}
		} else {
			return errors.New(aiTaskInputAssistant.BizType + ":BizType is wrong")
		}
	} else if aiTaskInputAssistant.AssistantBiz == AiAssistantJobBizType_genDocEmail {
		//_, err = c.HandleGenStatemaentFromAiTask(ctx, tCase, task, *statementConditionEntity, aiTaskInputAssistant.UserInputText)
		//if err != nil {
		//	c.log.Error(err)
		//	return err
		//}
		return c.HandleGenDocEmail(ctx, *tCase, *tClient, task)
	} else if aiTaskInputAssistant.AssistantBiz == AiAssistantJobBizType_statementSection {
		statementConditionEntity, err := c.StatementConditionUsecase.GetByCond(Eq{"id": aiTaskInputAssistant.StatementConditionId})
		if err != nil {
			return err
		}
		if statementConditionEntity == nil {
			return errors.New("statementConditionEntity is nil")
		}
		_, _, aiResultId, err := c.AssistantUpdatePSFromStatementCondition(ctx, *tClient, *tCase, *statementConditionEntity, aiTaskInputAssistant.UserInputText, aiTaskInputAssistant.SectionType)
		if err != nil {
			return err
		}
		task.CurrentResultId = aiResultId
	} else {
		return errors.New("InputAssistant.AssistantBiz  is wrong")
		//if aiTaskInputAssistant.BizType == AiAssistantBizType_UpdateMedication {
		//	_, _, aiResultId, err := c.GenUpdatePSFromStatementCondition(ctx, *tClient, *tCase, *statementConditionEntity, aiTaskInputAssistant.UserInputText, aiTaskInputAssistant.SectionType)
		//	if err != nil {
		//		return err
		//	}
		//	task.CurrentResultId = aiResultId
		//
		//} else {
		//	_, _, aiResultId, err := c.AssistantUpdatePSFromStatementCondition(ctx, *tClient, *tCase, *statementConditionEntity, aiTaskInputAssistant.UserInputText, aiTaskInputAssistant.SectionType)
		//	if err != nil {
		//		return err
		//	}
		//	task.CurrentResultId = aiResultId
		//}
	}

	return nil
}

func (c *AssistantUsecase) HandleGenDocEmail(ctx context.Context, tCase TData, tClient TData, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	aiTaskInputAssistant := task.GetAiTaskInputAssistant()
	if aiTaskInputAssistant.BizType == AiAssistantBizType_DocEmailRenew {
		aiResultId, err := c.ExecGenerateDocEmail(ctx, &tCase, aiTaskInputAssistant.UserInputText)
		if err != nil {
			return err
		}
		task.CurrentResultId = aiResultId
	} else {
		return errors.New("Assistant.BizType is wrong")
	}

	return nil
}

func (c *AssistantUsecase) AssistantUpdatePSFromStatementCondition(ctx context.Context, tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, userCustomPromptText string, sectionType string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	_, statements, err := c.StatementUsecase.GetUpdatePSTextForAiParamWithAssistant(tClient, tCase, statementConditionEntity, sectionType)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if statements == "" {
		return Claude3Response{}, "", 0, errors.New("The statement does not exist.")
	}
	promptKey := CurrentGenStatementPrompt
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var claude3Messages []Claude3Message

	var contents []Claude3Content

	intakeText, referenceContent, err := c.AiTaskbuzUsecase.GetGenStatementReferenceMaterials(ctx, tCase, statementConditionEntity.ToStatementCondition(), true)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	if intakeText != "" {
		claude3Messages = append(claude3Messages, Claude3Message{
			Role: "user",
			Content: []Claude3Content{
				{
					Type: "text",
					Text: intakeText,
				},
			},
		})
	}
	veteranSummaryVo, err := c.AiTaskbuzUsecase.HandleVeteranSummary(ctx, &tCase)
	var veteranSummary string
	if err == nil {
		veteranSummary = veteranSummaryVo.ToString()
		//claude3Messages = append(claude3Messages, Claude3Message{
		//	Role: "assistant",
		//	Content: []Claude3Content{
		//		{
		//			Type: "text",
		//			Text: veteranSummary,
		//		},
		//	},
		//})
	} else {
		c.log.Error(err, "caseId: ", tCase.Id())
	}

	var userInput string
	//	if referenceContent != "" {
	//		userInput = "Reference materials:\n" + referenceContent + "\n\n"
	//	}
	//	// 如何有用户自定义Instructions， 可以放在下面
	//	userInput += `Please update the VA statement only when "Instructions" are provided. Apply the following instructions precisely:
	//Instructions:`
	//	if userCustomPromptText != "" {
	//		userInput += "\n-" + userCustomPromptText + "\n\n"
	//	}
	//	userInput += "Expected JSON:\n```json{\n  \"update_required\": true or false,\n  \"updated_statement\": \"Full updated VA statement or null\"\n}```\n"
	//	//userInput += fmt.Sprintf("\nReference: \n%s", referenceContent)
	//	userInput += fmt.Sprintf("\nVA Statement: \n%s", statements)

	userInput = fmt.Sprintf(`You are tasked with updating the **VA Statement** based on provided "Instructions" and, if available, "Reference materials."

### Key Rules:

- **Only update the VA Statement when "Instructions" are provided.**
- **Reference materials are optional.** If provided, you must incorporate relevant information into the updated statement.
- **If no Reference materials are provided, you should still update the VA Statement according to the Instructions.**
- **If no update is required, return the original VA Statement unchanged.**
- **If the VA Statement begins with "SERVICE CONNECTION:", you must retain this prefix exactly as provided.**

### Writing Style:

- Follow the writing style, tone, and clarity guidelines from the **VA Disability Personal Statement Generator** system prompt:
  - Use a respectful, first-person, natural, conversational tone.
  - Write in plain English suitable for a middle-aged veteran with a high school education.
  - Avoid using phrases like "as I mentioned" or "like I said."
  - Do not add formatting, special symbols, markdown, or headings.
  - Keep the sentence structure simple and easy to follow.

### Input Structure:

- **Reference materials (optional):** %s
- **Instructions (required):** %s
- **VA Statement (required):** %s

### Output JSON Format:`, veteranSummary+"\n"+referenceContent+"\n\n", userCustomPromptText+"\n\n", statements+"\n\n")

	userInput += "\n\n```json\n{\n  \"update_required\": true or false,\n  \"updated_statement\": \"Full updated VA statement or null\"\n}\n```"

	userInput += "\n\n- Set `\"update_required\"` to `true` if the VA Statement is updated.\n- Set `\"update_required\"` to `false` and `\"updated_statement\"` to `null` if no changes are needed."

	contents = append(contents, Claude3Content{
		Type: "text",
		Text: userInput,
	})

	claude3Messages = append(claude3Messages, Claude3Message{
		Role:    "user",
		Content: contents,
	})
	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        4048,
		Messages:         claude3Messages,
	}

	if prompt != "" {
		payload.SystemPrompt = prompt
	}

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenUpdateStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

/*
Update the VA Statement based on the provided 'Current Treatment Facility' and/or 'Current Medication'. Follow these rules:
- If 'Current Treatment Facility' is provided, replace all instances of the previous treatment facility in the VA Statement with the new one, including historical references.
- If 'Current Medication' is provided and no medication is mentioned in the original VA Statement, append a new sentence at the end of the relevant paragraph (where treatment is discussed) stating: 'I am currently prescribed [medication names] to manage related symptoms.' If medications already exist, replace them with the new ones in the same sentence structure.
- Maintain the original structure, formatting, and all other content of the VA Statement unchanged.
- Set 'update_required' to true if any updates are made; otherwise, set it to false.
- Return the result in the following JSON format:
```json

	{
	  \"update_required\": true or false,
	  \"updated_statement\": \"Full updated VA statement or null\"
	}
*/
func (c *AssistantUsecase) GenUpdatePSFromStatementCondition(ctx context.Context, tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, userCustomPromptText string, sectionType string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	referenceContent, statements, err := c.StatementUsecase.GetUpdatePSTextForAiParamWithMedication(tClient, tCase, statementConditionEntity, sectionType)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if statements == "" {
		return Claude3Response{}, "", 0, errors.New("The statement does not exist.")
	}
	promptKey := CurrentPromptPSUpdate
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var claude3Messages []Claude3Message

	var contents []Claude3Content

	//
	// 如何有用户自定义Instructions， 可以放在下面
	userInput := `You must update the VA statement if new "Current Treatment Facility" or "Current Medication" is provided.
Instructions:
- If either "Current Treatment Facility" or "Current Medication"" is provided, update "VA Statement" accordingly.
- Always replace the old values with the new ones.
- Keep all other parts of the statement unchanged.
- If the "VA Statement" needs to be updated, set "update_required" to true; otherwise, set it to false.
`
	userInput = `Update the VA Statement based on the provided 'Current Treatment Facility' and/or 'Current Medication'. Follow these rules:
- If 'Current Treatment Facility' is provided, replace all instances of the previous treatment facility in the VA Statement with the new one, including historical references.
- If a "Current Treatment Facility" is provided and no treatment facility is mentioned in the original VA Statement, append a sentence in the relevant paragraph (where treatment or diagnosis is discussed) to include the new facility, e.g., "I sought treatment at [Current Treatment Facility] where I was officially diagnosed with Obstructive Sleep Apnea."
- If 'Current Medication' is provided, update 'VA Statement' accordingly.
- Maintain the original structure, formatting, and all other content of the VA Statement unchanged.
- Set 'update_required' to true if any updates are made; otherwise, set it to false.
`

	userInput = `Update the VA Statement based on the provided 'Current Treatment Facility' and/or 'Current Medication'. Follow these rules:
- If 'Current Treatment Facility' is provided, replace all instances of the previous treatment facility in the VA Statement with the new one, including historical references.
- If a "Current Treatment Facility" is provided and no treatment facility is mentioned in the original VA Statement, append a sentence in the relevant paragraph (where treatment or diagnosis is discussed) to include the new facility.
- If 'Current Medication' is provided, update 'VA Statement' accordingly.
- Maintain the original structure, formatting, and all other content of the VA Statement unchanged.
- Set 'update_required' to true if any updates are made; otherwise, set it to false.
`
	// Set 'update_required' to true if any updates are made; otherwise, set it to false.
	if userCustomPromptText != "" {
		userInput += "\n-" + userCustomPromptText + "\n\n"
	}
	userInput += "Expected JSON:\n```json{\n  \"update_required\": true or false,\n  \"updated_statement\": \"Full updated VA statement or null\"\n}```"
	userInput += fmt.Sprintf("\nThis is the new \"Current Treatment Facility\" or \"Current Medication\": \n%s", referenceContent)
	userInput += fmt.Sprintf("\nVA Statement: \n%s", statements)

	contents = append(contents, Claude3Content{
		Type: "text",
		Text: userInput,
	})

	claude3Messages = append(claude3Messages, Claude3Message{
		Role:    "user",
		Content: contents,
	})
	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        4048,
		Messages:         claude3Messages,
	}

	if prompt != "" {
		payload.SystemPrompt = prompt
	}

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenUpdateStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

func (c *AssistantUsecase) GenUpdatePSFromStatementConditionV1(ctx context.Context, tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, userCustomPromptText string, sectionType string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	referenceContent, statements, err := c.StatementUsecase.GetUpdatePSTextForAiParamWithMedication(tClient, tCase, statementConditionEntity, sectionType)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if statements == "" {
		return Claude3Response{}, "", 0, errors.New("The statement does not exist.")
	}
	promptKey := CurrentPromptPSUpdate
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var claude3Messages []Claude3Message

	var contents []Claude3Content

	// 如何有用户自定义Instructions， 可以放在下面
	userInput := `You must update the VA statement if new "Current Treatment Facility" or "Current Medication" is provided.
Instructions:
- If either "Current Treatment Facility" or "Current Medication"" is provided, update "VA Statement" accordingly.
- Always replace the old values with the new ones.
- Keep all other parts of the statement unchanged.
- If the "VA Statement" needs to be updated, set "update_required" to true; otherwise, set it to false.
`
	if userCustomPromptText != "" {
		userInput += "\n-" + userCustomPromptText + "\n\n"
	}
	userInput += "Expected JSON:\n```json{\n  \"update_required\": true or false,\n  \"updated_statement\": \"Full updated VA statement or null\"\n}```"
	userInput += fmt.Sprintf("\nThis is the new \"Current Treatment Facility\" or \"Current Medication\": \n%s", referenceContent)
	userInput += fmt.Sprintf("\nVA Statement: \n%s", statements)

	contents = append(contents, Claude3Content{
		Type: "text",
		Text: userInput,
	})

	claude3Messages = append(claude3Messages, Claude3Message{
		Role:    "user",
		Content: contents,
	})
	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        4048,
		Messages:         claude3Messages,
	}

	if prompt != "" {
		payload.SystemPrompt = prompt
	}

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenUpdateStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

func (c *AssistantUsecase) GetAiResultEntityFromAiTaskId(aiTaskId int32) (*AiTaskEntity, *AiResultEntity, error) {
	aiTaskEntity, err := c.AiTaskUsecase.GetByCond(Eq{"id": aiTaskId, "deleted_at": 0})
	if err != nil {
		return nil, nil, err
	}
	if aiTaskEntity == nil {
		return nil, nil, nil
	}
	resultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTaskEntity.CurrentResultId})
	return aiTaskEntity, resultEntity, err
}
