package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AiTaskbuzUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	AiTaskUsecase             *AiTaskUsecase
	Awsclaude3Usecase         *Awsclaude3Usecase
	AiPromptUsecase           *AiPromptUsecase
	TUsecase                  *TUsecase
	AiResultUsecase           *AiResultUsecase
	FieldOptionUsecase        *FieldOptionUsecase
	FieldUsecase              *FieldUsecase
	DataEntryUsecase          *DataEntryUsecase
	AiUsecase                 *AiUsecase
	JotformSubmissionUsecase  *JotformSubmissionUsecase
	RelasLogUsecase           *RelasLogUsecase
	QuestionnairesbuzUsecase  *QuestionnairesbuzUsecase
	ConditionUsecase          *ConditionUsecase
	FeeUsecase                *FeeUsecase
	WordUsecase               *WordUsecase
	StatementUsecase          *StatementUsecase
	BoxUsecase                *BoxUsecase
	DataComboUsecase          *DataComboUsecase
	StatementConditionUsecase *StatementConditionUsecase
}

func NewAiTaskbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	AiTaskUsecase *AiTaskUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	AiPromptUsecase *AiPromptUsecase,
	TUsecase *TUsecase,
	AiResultUsecase *AiResultUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	FieldUsecase *FieldUsecase,
	DataEntryUsecase *DataEntryUsecase,
	AiUsecase *AiUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	RelasLogUsecase *RelasLogUsecase,
	QuestionnairesbuzUsecase *QuestionnairesbuzUsecase,
	ConditionUsecase *ConditionUsecase,
	FeeUsecase *FeeUsecase,
	WordUsecase *WordUsecase,
	StatementUsecase *StatementUsecase,
	BoxUsecase *BoxUsecase,
	DataComboUsecase *DataComboUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
) *AiTaskbuzUsecase {
	uc := &AiTaskbuzUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		AiTaskUsecase:             AiTaskUsecase,
		Awsclaude3Usecase:         Awsclaude3Usecase,
		AiPromptUsecase:           AiPromptUsecase,
		TUsecase:                  TUsecase,
		AiResultUsecase:           AiResultUsecase,
		FieldOptionUsecase:        FieldOptionUsecase,
		FieldUsecase:              FieldUsecase,
		DataEntryUsecase:          DataEntryUsecase,
		AiUsecase:                 AiUsecase,
		JotformSubmissionUsecase:  JotformSubmissionUsecase,
		RelasLogUsecase:           RelasLogUsecase,
		QuestionnairesbuzUsecase:  QuestionnairesbuzUsecase,
		ConditionUsecase:          ConditionUsecase,
		FeeUsecase:                FeeUsecase,
		WordUsecase:               WordUsecase,
		StatementUsecase:          StatementUsecase,
		BoxUsecase:                BoxUsecase,
		DataComboUsecase:          DataComboUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
	}

	return uc
}

func (c *AiTaskbuzUsecase) HandleReturnTimezone(tClient *TData) error {
	if tClient == nil {
		return errors.New("tCase is nil")
	}
	aiTask, err := c.AiTaskUsecase.GetReturnTimezone(tClient)
	if err != nil {
		return err
	}
	if aiTask == nil {
		_, err = c.AiTaskUsecase.CreateReturnTimezone(tClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AiTaskbuzUsecase) HandleUpdateStatement(ctx context.Context, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	aiTaskInputUpdateStatementTask := task.GetAiTaskInputUpdateStatementTask()

	textForAi, err := aiTaskInputUpdateStatementTask.PersonalStatementOneVo.ToTextForAi()
	if err != nil {
		c.log.Error(err)
		return err
	}
	_, parseResult, aiResultId, err := c.GenUpdatePS(ctx, *tCase, textForAi)
	if err != nil {
		return err
	}

	aiTaskResultEntity := AiTaskResultEntity{
		AiTaskId:   task.ID,
		AiResultId: aiResultId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
	if err != nil {
		return err
	}
	task.CurrentResultId = aiResultId

	_, _, err = UpdatePersonalStatementsAiResultFormat(parseResult)
	if err != nil {
		return err
	}

	//listParseUpdateStatementVo, err := c.HandleUpdateStatementDoc(aiTaskInputUpdateStatementTask.PersonalStatementsVo, aiResultId)
	//if err != nil {
	//	return err
	//}
	//
	//err = c.CreateUpdateStatementDoc(*tClient, *tCase, aiTaskInputUpdateStatementTask.PersonalStatementsVo, listParseUpdateStatementVo)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (c *AiTaskbuzUsecase) HandleUpdateStatementDocToBox(tClient TData, tCase TData) error {

	c.log.Debug("HandleUpdateStatementDocToBox caseId:", tCase.Id())
	a, err := c.AiTaskUsecase.GetByCond(And(Eq{"case_id": tCase.Id(),
		"from_type":  AiTaskFromType_update_statement,
		"deleted_at": 0,
	}, Or(Eq{"handle_status": 0}, Eq{"handle_result": 1, "handle_status": 1})))
	if err != nil {
		return err
	}
	if a != nil {
		c.log.Debug("The task is still unfinished or has failed")
		return nil
	}

	return c.DoUpdateStatementDocToBox(tClient, tCase)
}

func (c *AiTaskbuzUsecase) DoUpdateStatementDocToBox(tClient TData, tCase TData) error {

	tasks, err := c.AiTaskUsecase.AllByCondWithOrderBy(Eq{"case_id": tCase.Id(),
		"from_type":  AiTaskFromType_update_statement,
		"deleted_at": 0}, "serial_number", 1000)
	if err != nil {
		return err
	}
	var listParseUpdateStatementVo ListParseUpdateStatementVo
	var personalStatementsForCreateWordVo *PersonalStatementsForCreateWordVo
	for _, v := range tasks {

		if personalStatementsForCreateWordVo == nil {
			inputVo := v.GetAiTaskInputUpdateStatementTask()
			personalStatementsForCreateWordVo = &PersonalStatementsForCreateWordVo{
				BaseInfo: inputVo.PersonalStatementOneVo.BaseInfo,
			}
		}
		resultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": v.CurrentResultId})
		if err != nil {
			return err
		}

		if resultEntity == nil {
			return errors.New("resultEntity is nil: " + InterfaceToString(v.CurrentResultId))
		}
		noUpdated, vo, err := UpdatePersonalStatementsAiResultFormat(resultEntity.ParseResult)
		if err != nil {
			c.log.Error(err, " resultId: ", InterfaceToString(resultEntity.ID))
			return err
		}
		if noUpdated {
			inputVo := v.GetAiTaskInputUpdateStatementTask()
			vo, err = inputVo.PersonalStatementOneVo.ToParseUpdateStatementVo()
			if err != nil {
				return err
			}
			listParseUpdateStatementVo = append(listParseUpdateStatementVo, vo)
		} else {
			listParseUpdateStatementVo = append(listParseUpdateStatementVo, vo)
		}
	}
	if personalStatementsForCreateWordVo == nil {
		return errors.New("personalStatementsForCreateWordVo is nil")
	}
	return c.CreateUpdateStatementDoc(tClient, tCase, *personalStatementsForCreateWordVo, listParseUpdateStatementVo)
}

// UpdatePersonalStatementsAiResultFormat 此方法可以验证是否有问题
func UpdatePersonalStatementsAiResultFormat(text string) (noUpdated bool, parseUpdateStatementVo ParseUpdateStatementVo, err error) {

	if strings.Index(text, "no changes are needed") >= 0 { // 不需要修改
		return true, parseUpdateStatementVo, nil
	}

	parseUpdateStatementVo, err = SplitUpdatePersonalStatementsAiResult(text)
	return
}

func (c *AiTaskbuzUsecase) CreateUpdateStatementDoc(tClient TData, tCase TData, personalStatementsForCreateWordVo PersonalStatementsForCreateWordVo, listParseUpdateStatementVo ListParseUpdateStatementVo) error {
	c.log.Debug("CreateUpdateStatementDoc caseId:", tCase.Id())
	wordReader, err := c.WordUsecase.CreateUpdatePersonalStatementsWord(personalStatementsForCreateWordVo, listParseUpdateStatementVo)
	if err != nil {
		return err
	}
	dCPersonalStatementsFolderId, updateStatementFileName, boxFileId, err := c.StatementUsecase.DocUpdateStatementBoxFileId(tClient, tCase)
	if err != nil {
		return err
	}
	if boxFileId == "" {
		boxFileId, err = c.BoxUsecase.UploadFile(dCPersonalStatementsFolderId, wordReader, updateStatementFileName)
		if err != nil {
			return err
		}
	} else {
		_, err = c.BoxUsecase.UploadFileVersion(boxFileId, wordReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AiTaskbuzUsecase) HandleUpdateStatementDoc(personalStatementsVo PersonalStatementsVo, aiResultId int32) (listParseUpdateStatementVo ListParseUpdateStatementVo, err error) {

	resultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": aiResultId})
	if err != nil {
		return nil, err
	}
	if resultEntity == nil {
		return nil, errors.New("resultEntity is nil")
	}

	resultList, err := ParseUpdateStatementResult(resultEntity.ParseResult)
	if err != nil {
		return nil, err
	}

	listParseUpdateStatementVo, err = personalStatementsVo.ToStatements()

	for k, _ := range listParseUpdateStatementVo {
		newParseUpdateStatementVo := ReplaceStatement(listParseUpdateStatementVo[k], resultList)
		listParseUpdateStatementVo[k] = newParseUpdateStatementVo
	}
	//lib.DPrintln("listParseUpdateStatementVo:  ", listParseUpdateStatementVo, err)

	//docReader, err := c.WordUsecase.CreateUpdatePersonalStatementsWord(personalStatementsVo, listParseUpdateStatementVo)
	//
	//if err != nil {
	//	return nil, err
	//}
	////c.log.Info(docReader)
	//file, _ := os.Create("aaa.docx")
	//io.Copy(file, docReader)
	//defer file.Close()

	return
}

func ReplaceStatement(parseUpdateStatementVo ParseUpdateStatementVo, listParseUpdateStatementVo ListParseUpdateStatementVo) (newParseUpdateStatementVo ParseUpdateStatementVo) {

	newParseUpdateStatementVo = parseUpdateStatementVo
	for _, v := range listParseUpdateStatementVo {
		if strings.ToLower(v.NameOfDisability.Value) == strings.ToLower(parseUpdateStatementVo.NameOfDisability.Value) {
			newParseUpdateStatementVo.CurrentTreatment.Value = v.CurrentTreatment.Value
			newParseUpdateStatementVo.CurrentMedication.Value = v.CurrentMedication.Value
		}
	}

	return
}

func (c *AiTaskbuzUsecase) HandleReturnTimezoneJob(ctx context.Context, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	promptKey := "prompt_return_timezone"
	//input := task.GetIput()

	tClient, err := c.TUsecase.DataById(Kind_clients, task.ClientId)
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	if tClient.CustomFields.TextValueByNameBasic(FieldName_timezone_id) != "" {
		return nil
	}

	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return err
	}
	payload := Claude3Request{
		AnthropicVersion: AnthropicVersion_2023,
		MaxTokens:        2048,
		Messages: []Claude3Message{
			{
				Role: "user",
				Content: []Claude3Content{
					{
						Type: "text",
						Text: fmt.Sprintf("%s, %s", tClient.CustomFields.TextValueByNameBasic(FieldName_state), tClient.CustomFields.TextValueByNameBasic(FieldName_city)),
					},
				},
			},
		},
	}
	if prompt != "" {
		payload.SystemPrompt = prompt
	}

	_, aiResultId, err := c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_ReturnTimeZone, promptKey, true, func(parseResult string) bool {
		jsonMap := lib.ToTypeMapByString(GetJsonFromAiResult(parseResult))
		timezone := jsonMap.GetString("timezone")
		if timezone == "" {
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	aiTaskResultEntity := AiTaskResultEntity{
		AiTaskId:   task.ID,
		AiResultId: aiResultId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
	if err != nil {
		return err
	}
	task.CurrentResultId = aiResultId

	err = c.HandleCaseTimeZone(tClient, aiResultId)
	if err != nil {
		return err
	}
	return nil
}

func (c *AiTaskbuzUsecase) HandleCaseTimeZone(tClient *TData, aiResultId int32) error {

	if tClient == nil {
		return errors.New("tClient is nil")
	}
	if tClient.CustomFields.TextValueByNameBasic(FieldName_timezone_id) != "" {
		return nil
	}
	aiResult, err := c.AiResultUsecase.GetByCond(Eq{"id": aiResultId})
	if err != nil {
		return err
	}
	if aiResult == nil {
		return errors.New("aiResult is nil")
	}
	jsonMap := lib.ToTypeMapByString(GetJsonFromAiResult(aiResult.ParseResult))

	timezone := jsonMap.GetString("timezone")
	if timezone == "" {
		c.log.Error("aiResultId: " + InterfaceToString(aiResultId))
		return errors.New("timezone is empty")
	}

	fieldEntity, err := c.FieldUsecase.GetByFieldName(Kind_client_cases, FieldName_timezone_id)
	if err != nil {
		return err
	}
	if fieldEntity == nil {
		return errors.New("fieldEntity is nil")
	}
	option, err := c.FieldOptionUsecase.GetByEntity(*fieldEntity, timezone)
	if err != nil {
		return err
	}
	if option == nil {
		return errors.New(timezone + " does not match option")
	}

	data := make(TypeDataEntry)
	data[DataEntry_gid] = tClient.Gid()
	data[FieldName_timezone_id] = timezone
	_, err = c.DataEntryUsecase.HandleOne(Kind_clients, data, DataEntry_gid, nil)
	if err != nil {
		return err
	}
	//lib.DPrintln(jsonMap)
	return nil
}

func (c *AiTaskbuzUsecase) HandleAutoAssociateJotform(tCase TData, statementCondition StatementCondition) (submissions []*JotformSubmissionEntity, err error) {

	conditionUuid := InterfaceToString(statementCondition.StatementConditionId)
	//uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	submissions, err = c.JotformSubmissionUsecase.AllByUniqcodeAndConditionUniqid(conditionUuid)
	if err != nil {
		return nil, err
	}
	if len(submissions) == 0 {
		err := c.DoAutoAssociateJotform(tCase, statementCondition)
		if err != nil {
			return nil, err
		}
		submissions, err = c.JotformSubmissionUsecase.AllByUniqcodeAndConditionUniqid(conditionUuid)
		if err != nil {
			return nil, err
		}
	}

	return submissions, nil
}

func AssociateJotformGetDataList(jotformSubmissions []*JotformSubmissionEntity) (dataList string, err error) {
	for _, v := range jotformSubmissions {
		newFile, err := v.JotformNewFileNameForAI()
		if err != nil {
			return "", err
		}
		if dataList == "" {
			dataList = newFile
		} else {
			dataList += "\n" + newFile
		}
	}
	return dataList, nil
}

func (c *AiTaskbuzUsecase) DoAutoAssociateJotform(tCase TData, statementCondition StatementCondition) error {

	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)

	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(&tCase)
	if err != nil {
		return err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	jotformSubmissions, err := c.JotformSubmissionUsecase.AllLatestByUniqcodeExceptIntake(uniqcodes)
	if err != nil {
		return err
	}
	if len(jotformSubmissions) == 0 {
		return errors.New("jotformSubmissions is nil")
	}

	dataList, err := AssociateJotformGetDataList(jotformSubmissions)
	if err != nil {
		return err
	}
	parseResult, _, err := c.AiUsecase.ExecutePromptAssociateJotform(statementCondition.ConditionValue, dataList)
	if err != nil {
		return err
	}
	err = c.HandleAssociateJotform(parseResult, tCase, statementCondition)
	if err != nil {
		return err
	}
	return nil
}

func ParseAssociateJotformResult(parseResult string) (submissionIds []string) {

	parseResultMap := lib.ToTypeMapByString(GetJsonFromAiResult(parseResult))
	relatedEntries := parseResultMap.GetTypeListInterface("related_entries")
	for _, v := range relatedEntries {
		submissionId := GetJotformSubmissionIdFromFileName(InterfaceToString(v))
		if submissionId != "" {
			submissionIds = append(submissionIds, submissionId)
		}
	}
	return submissionIds
}

func (c *AiTaskbuzUsecase) HandleAssociateJotform(parseResult string, tCase TData, statementCondition StatementCondition) error {

	//uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	//parseResultMap := lib.ToTypeMapByString(GetJsonFromAiResult(parseResult))
	submissionIds := ParseAssociateJotformResult(parseResult)
	if len(submissionIds) == 0 {
		return errors.New("The associated data that was not obtained")
	}

	for _, v := range submissionIds {
		submissionId := v
		submissionEntity, err := c.JotformSubmissionUsecase.GetLatestFormInfoWithUniqcode(submissionId)
		if err != nil {
			return err
		}
		if submissionEntity != nil {
			sourceId := InterfaceToString(statementCondition.StatementConditionId)
			_, err = c.RelasLogUsecase.ConditionUpsert(sourceId, InterfaceToString(submissionId), nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *AiTaskbuzUsecase) GetGenStatementReferenceMaterials(ctx context.Context, tCase TData, StatementCondition StatementCondition, isAssistant bool) (intakeText string, referenceContent string, err error) {

	var intakeSubmission *JotformSubmissionEntity
	var otherSubmissions []*JotformSubmissionEntity
	if configs.EnableAiAutoAssociationJotform {
		_, err = c.HandleAutoAssociateJotform(tCase, StatementCondition)
		if err != nil {
			if !isAssistant {
				return "", "", errors.New("HandleAutoAssociateJotform error: " + err.Error())
			}
		}
	}

	intakeSubmission, otherSubmissions, err = c.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatementNew(&tCase, StatementCondition)
	if err != nil {
		return "", "", err
	}
	if len(otherSubmissions) == 0 {
		if !isAssistant {
			return "", "", errors.New("The jotform submissions need to be associated")
		}
	}

	var res []*JotformSubmissionEntity
	res = append(res, otherSubmissions...)

	notes := lib.ToTypeMapByString(intakeSubmission.Notes)
	intakeText = FormatJotformAnswersForGenStatement(notes)

	for _, v := range res {
		notes = lib.ToTypeMapByString(v.Notes)
		str := FormatJotformAnswersForGenStatement(notes)
		if referenceContent == "" {
			referenceContent = str
		} else {
			referenceContent += "\n" + str
		}
	}
	return intakeText, referenceContent, nil
}

func (c *AiTaskbuzUsecase) GenStatement(ctx context.Context, tCase *TData, StatementCondition StatementCondition, veteranSummary string, promptKey string, userInputPrompt string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var claude3Messages []Claude3Message
	var contents []Claude3Content

	//var intakeSubmission *JotformSubmissionEntity
	//var otherSubmissions []*JotformSubmissionEntity

	//if lib.EnableAiAutoAssociationJotform {
	//	_, err = c.HandleAutoAssociateJotform(*tCase, StatementCondition)
	//	if err != nil {
	//		return Claude3Response{}, "", 0, errors.New("HandleAutoAssociateJotform error: " + err.Error())
	//	}
	//}

	//intakeSubmission, otherSubmissions, err = c.QuestionnairesbuzUsecase.GetJotformSubmissionsForGenStatementNew(tCase, StatementCondition)
	//if err != nil {
	//	return Claude3Response{}, "", 0, err
	//}
	//if len(otherSubmissions) == 0 {
	//	return Claude3Response{}, "", 0, errors.New("The jotform submissions need to be associated")
	//}
	//

	//var res []*JotformSubmissionEntity
	//res = append(res, otherSubmissions...)

	//notes := lib.ToTypeMapByString(intakeSubmission.Notes)
	//intakeText := FormatJotformAnswersForGenStatement(notes)

	intakeText, referenceContent, err := c.GetGenStatementReferenceMaterials(ctx, *tCase, StatementCondition, false)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	claude3Messages = append(claude3Messages, Claude3Message{
		Role: "user",
		Content: []Claude3Content{
			{
				Type: "text",
				Text: intakeText,
			},
		},
	})
	claude3Messages = append(claude3Messages, Claude3Message{
		Role: "assistant",
		Content: []Claude3Content{
			{
				Type: "text",
				Text: veteranSummary,
			},
		},
	})

	//referenceContent := ""
	//for _, v := range res {
	//	notes = lib.ToTypeMapByString(v.Notes)
	//	str := FormatJotformAnswersForGenStatement(notes)
	//	if referenceContent == "" {
	//		referenceContent = str
	//	} else {
	//		referenceContent += "\n" + str
	//	}
	//}

	userInput := "Reference materials:\n" + referenceContent + "\n\n"
	userInput += "Use only this Condition: " + StatementCondition.ToOriginValue() + "\nOther content may be referenced, but ignore any other conditions."
	userInputPrompt = strings.TrimSpace(userInputPrompt)
	userInputPrompt = strings.Trim(userInputPrompt, ".")
	if userInputPrompt != "" {
		userInput += " " + userInputPrompt + "."
	}

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

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

type StandardHeaderRevisionAiResult struct {
	SpecialNotes                           string `json:"Special Notes"`
	Introduction                           string `json:"Introduction"`
	OnsetAndServiceConnection              string `json:"Onset and Service Connection"`
	CurrentSymptomsSeverityAndFrequency    string `json:"Current Symptoms Severity and Frequency"`
	Medication                             string `json:"Medication"`
	ImpactOnDailyLife                      string `json:"Impact on Daily Life"`
	ProfessionalImpact                     string `json:"Professional Impact"`
	NexusBetweenServiceAndCurrentCondition string `json:"Nexus Between Service and Current Condition"`
	Request                                string `json:"Request"`
}

func (c *StandardHeaderRevisionAiResult) ToParseAiStatementConditionVo() (vo ParseAiStatementConditionVo) {
	vo.SpecialNotes = c.SpecialNotes
	vo.IntroductionParagraph = c.Introduction
	vo.OnsetAndServiceConnection = c.OnsetAndServiceConnection
	vo.CurrentSymptomsSeverityFrequency = c.CurrentSymptomsSeverityAndFrequency
	vo.Medication = c.Medication
	vo.ImpactOnDailyLife = c.ImpactOnDailyLife
	vo.ProfessionalImpact = c.ProfessionalImpact
	vo.NexusBetweenSC = c.NexusBetweenServiceAndCurrentCondition
	vo.Request = c.Request
	return vo
}

func StringToStandardHeaderRevisionAiResult(str string) (vo StandardHeaderRevisionAiResult) {
	return lib.StringToTDef(str, vo)
}

func (c *AiTaskbuzUsecase) GenStatementForStandardHeaderRevision(ctx context.Context, tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, veteranSummary VeteranSummaryVo, userInputPrompt string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

	promptKey := Current_StandardHeaderRevisionPrompt
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	var claude3Messages []Claude3Message
	var contents []Claude3Content

	statementDetail, err := c.StatementUsecase.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var userInput string
	var cFaci, cMedi string

	for _, v := range statementDetail.Statements {
		if v.StatementCondition.StatementConditionId == statementConditionEntity.ID {
			userInput = v.Rows.ToStringForStandardHeaderRevision()
			cFaci, cMedi = v.Rows.GetCurrentTreatmentFacilityAndCurrentMedication()

			break
		}
	}
	baseInfo := veteranSummary.ToString()
	baseInfo += "\nCurrent Treatment Facility: " + cFaci
	baseInfo += "\nCurrent Medication: " + cMedi

	claude3Messages = append(claude3Messages, Claude3Message{
		Role: "assistant",
		Content: []Claude3Content{
			{
				Type: "text",
				Text: baseInfo,
			},
		},
	})

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

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}

func (c *AiTaskbuzUsecase) TestGenUpdatePSFromStatementCondition() {
	//tCase, _ := c.TUsecase.DataById(Kind_client_cases, 5572)
	//tClient, _, _ := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	//
	//if tCase == nil {
	//	return
	//}
	//if tClient == nil {
	//	return
	//}
	//condition, _ := c.StatementConditionUsecase.GetByCond(Eq{"id": 27})
	//if condition == nil {
	//	return
	//}
	//_, b, aiResultId, err := c.GenUpdatePSFromStatementCondition(context.TODO(), *tClient, *tCase, *condition)
	//if err != nil {
	//	lib.DPrintln(err)
	//	return
	//}
	//lib.DPrintln("TestGenUpdatePSFromStatementCondition", b, aiResultId)
}

func (c *AiTaskbuzUsecase) GenUpdatePSFromStatementCondition(ctx context.Context, tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, sectionType string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {

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

	// old
	//userInput := "Reference materials:\n" + referenceContent + "\n\n"
	//userInput += "VA statement(s):\n" + statements

	// new

	// 如何有用户自定义Instructions， 可以放在下面
	userInput := `You must update the VA statement if new Current Treatment Facility or Current Medication is provided.
Instructions:
- If either Current Treatment Facility or Current Medication is provided, set "update_required": true and update the statement accordingly.
- Always replace the old values with the new ones.
- Keep all other parts of the statement unchanged.

` + "Expected JSON:\n```json{\n  \"update_required\": true or false,\n  \"updated_statement\": \"Full updated VA statement or null\"\n}```"
	userInput += fmt.Sprintf("\nReference: \n%s", referenceContent)
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

func (c *AiTaskbuzUsecase) GenUpdatePS(ctx context.Context, tCase TData, statements string) (claude3Response Claude3Response, parseResult string, aiResultId int32, err error) {
	promptKey := CurrentPromptPSUpdate
	prompt, _, err := c.AiPromptUsecase.GetAiInfoByPromptKey(promptKey, lib.TypeMap{"text": ""})
	if err != nil {
		return Claude3Response{}, "", 0, err
	}

	var claude3Messages []Claude3Message

	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	submissions, err := c.JotformSubmissionUsecase.AllLatestUpdateQuestionnaires(uniqcode)
	if err != nil {
		return Claude3Response{}, "", 0, err
	}
	if len(submissions) == 0 {
		return Claude3Response{}, "", 0, errors.New("Jotform Submissions is empty")
	}

	var contents []Claude3Content

	referenceContent := ""
	for _, v := range submissions {
		notes := lib.ToTypeMapByString(v.Notes)
		str := FormatJotformAnswersForGenStatement(notes)
		if referenceContent == "" {
			referenceContent = str
		} else {
			referenceContent += "\n" + str
		}
	}

	//	statements := `Alexander Bagarry-20#5373
	//
	//
	//1.70 - Mixed Anxiety and Depressive Disorder secondary to Tinnitus, Hearing loss, left ear, and Asthma (new)
	//
	//2.50 - Headaches secondary to tinnitus (opinion)
	//
	//3.30 - Dizziness with staggering secondary to tinnitus (opinion)
	//
	//4.30 - Asthma (increase)
	//
	//5.20 - Low back pain (str)
	//
	//6.20* - Right knee sprain with limitation of flexion and extension (str)
	//
	//7.10 - Left foot plantar fasciitis (str)
	//
	//8.10 - Right foot plantar fasciitis secondary to left foot plantar fasciitis (opinion)
	//
	//9.80 - Gastroesophageal reflux disease (GERD)
	//
	//10.30 - Irritable Bowel Syndrome (Gulf War)
	//
	//11.20* - Left knee pain with limitation of flexion and extension secondary to right knee sprain (opinion)
	//
	//12.10 - Hemorrhoids (str, opinion)
	//
	//
	//
	//• Full Name: Alexander Bagarry
	//• Unique ID: 5373
	//• Branch of Service: Navy
	//• Years of Service: 1994-1998
	//• Retired from service: No
	//• Deployments: Persian Gulf (USS Vandergrift FFG-48 1995, USS Shiloh CG-67 1997), Operation Desert Strike
	//• Marital Status: Married (2007-Present)
	//• Children: 2 young children
	//• Occupation in service: Torpedoman's Mate (TM)
	//
	//
	//Name of Disability/Condition: Mixed Anxiety and Depressive Disorder secondary to Tinnitus, Hearing loss, left ear, and Asthma
	//Current Treatment Facility:
	//Current Medication: Lamictal 200 MG
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of Mixed Anxiety and Depressive Disorder secondary to Tinnitus, Hearing loss, left ear, and Asthma. This condition has been a significant challenge in my life since its onset during my active-duty service. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of this disorder have become increasingly debilitating, affecting not just my mental well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//My mental health challenges developed as a direct result of my service-connected tinnitus, hearing loss, and asthma. While stationed aboard the USS Vandergrift, I was exposed to extreme noise when the 76 MM gun fired unexpectedly above me. This incident led to immediate tinnitus and hearing issues that persist to this day. Additionally, after returning from deployment to the Persian Gulf, I developed breathing problems that were later diagnosed as asthma.
	//
	//The constant ringing in my ears, difficulty hearing, and breathing problems have created a perfect storm that has severely impacted my mental health. The unrelenting nature of these physical conditions has led to significant anxiety and depression that affect every aspect of my life.
	//
	//Currently, I'm taking Lamictal 200 MG for my condition. While this provides some relief, it is often short-term. My level of occupational and social impairment is characterized by deficiencies in most areas such as work, school, and family relations.
	//
	//I experience gross impairment in thought processes and communication. My thoughts often jump around randomly, and I have difficulty focusing on tasks or completing them. Written instructions often appear jumbled to me, and my voice trembles when I'm nervous. I frequently have difficulty finding the energy to start even simple tasks.
	//
	//I have developed significant memory issues, including trouble remembering people's names and details from my service. I sometimes forget what unit I served with, and planning ahead has become increasingly difficult. My judgment is impaired, leading to impulsive purchases at times.
	//
	//My behavior has changed dramatically. I always need to sit facing the door in restaurants due to heightened suspiciousness. I repeatedly check locks for safety and experience intense anxiety episodes in heavy traffic. I occasionally snap at people without warning and sometimes laugh at inappropriate moments. I've noticed shadows moving in my peripheral vision and find myself frequently planning defensive tactics.
	//
	//Sleep has become a major issue. I wake up frequently throughout the night to check the house. The combination of my physical conditions and anxiety makes it difficult to maintain regular sleep patterns.
	//
	//I struggle with basic self-care and daily activities. I occasionally wear the same clothes repeatedly and forget to eat regular meals. I regularly miss appointments and have difficulty navigating without GPS. During particularly challenging periods, I experience suicidal thoughts, with death occasionally feeling like an escape.
	//
	//My ability to maintain relationships has been severely impacted. I have difficulty maintaining friendships and find new situations extremely anxiety-provoking. I fear I cannot be the best spouse, parent, or friend due to my condition.
	//
	//Stressful situations or loud noises frequently trigger flashbacks to my time in service. I have difficulty crying, even at emotional events, and often feel disconnected from those around me.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Headaches secondary to tinnitus
	//Current Treatment Facility:
	//Current Medication: OTC Excedrin Migraine
	//
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of headaches secondary to tinnitus. This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of these headaches have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While stationed aboard the USS Vandergrift during deployment, I was exposed to extreme noise when the 76 MM gun fired unexpectedly above me. I was standing on the deck below, next to the torpedo tubes, without hearing protection. My ears immediately began ringing, and I lost my hearing for several minutes. Later, after returning from deployment to the Persian Gulf, I reported to medical that I had experienced headaches for two days along with cold symptoms. I did not have headaches before service or before developing tinnitus, and the condition has progressively worsened over time.
	//
	//Currently, I experience prostrating and prolonged attacks of headache pain 3-4 days per month, with each episode typically lasting a full day. The headaches have worsened both in frequency and severity since their onset. The pain is constant and pulsating or throbbing in nature, occurring on both sides of my head and worsening with physical activity.
	//
	//During these headache episodes, I also experience several non-headache symptoms including nausea, vomiting, sensitivity to light and sound, changes in vision, and high irritability. My symptoms worsen with exposure to bright or flickering lights and loud noises or prolonged exposure to noise.
	//
	//This condition severely affects my daily life. The headaches are significantly impacted by my tinnitus, with the constant ringing often preceding or accompanying the headache pain. I experience sensitivity to light and sound that limits my activities, and I face challenges in parenting or caring for dependents during episodes. During severe episodes, I need to retreat to a quiet, dark environment.
	//
	//My sleep is frequently disturbed by the headache pain, which often prevents me from falling asleep or causes early awakening. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I experience reduced productivity during headache episodes and have increased absenteeism due to the severity and frequency of the attacks. To manage the condition, I take over-the-counter pain medication and try to rest in a quiet, dark room or take naps when possible.
	//
	//The persistent pain and disruptions to my daily routine have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Dizziness with staggering secondary to tinnitus
	//Current Treatment Facility:
	//Current Medication:
	//
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of dizziness with staggering secondary to tinnitus. This condition has been a significant challenge in my life since its onset during my active-duty service in 1994. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of dizziness and staggering have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//During Basic Training, I developed lightheadedness that caused dizziness and staggering, requiring medical treatment. Later, while stationed aboard the USS Vandergrift in 1995-1996, I was exposed to extreme noise when the 76 MM gun fired unexpectedly above me. I was standing on the deck below, next to the torpedo tubes, without hearing protection. My ears immediately began ringing, and I lost my hearing for several minutes.
	//
	//After this incident, I began experiencing severe bouts of dizziness that caused me to stagger, often requiring me to find a place to sit or grab hold of something to avoid falling. I did not have any dizziness issues before service, and the condition has progressively worsened over time.
	//
	//Currently, my dizziness and staggering episodes occur a few times a week, lasting up to 15 minutes each time. The condition has worsened both in frequency and severity since its onset. I experience symptoms related to Meniere's syndrome and other vestibular conditions, including hearing impairment with vertigo, attacks of vertigo with clumsy staggering walking movements, and chronic tinnitus.
	//
	//This condition severely affects my daily life. I face significant balance issues that increase my risk of falls. There have been instances where my condition was so severe that I had to be taken to the Emergency Room from work due to dizzy spells where I kept falling down. These episodes are particularly challenging as they can occur without warning.
	//
	//My sleep is frequently disturbed by vertigo and dizziness. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work, primarily due to safety concerns related to falling. The unpredictable nature of these episodes makes it difficult to maintain consistent work performance and creates anxiety about potential incidents in the workplace.
	//
	//The persistent symptoms and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Asthma
	//Current Treatment Facility:
	//Current Medication: Albuterol Inhaler, Spiriva Inhaler
	//
	//
	//I am respectfully requesting an increase in Veteran Affairs benefits for my condition of asthma. This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of asthma have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While stationed aboard the USS Vandergrift during deployment to the Persian Gulf in 1995 as part of Operation Desert Strike, I was exposed to environmental conditions that impacted my respiratory health. After returning from this deployment, I began experiencing breathing problems that felt like my throat was closing up, accompanied by wheezing. I was later stationed aboard the USS Shiloh and deployed again to the Persian Gulf in 1997, which further exacerbated my condition. I did not have any breathing issues before deploying to the Persian Gulf, and my condition has progressively worsened over time.
	//
	//Currently, my asthma has worsened both in frequency and severity since its onset. I experience nearly constant episodes occurring daily, requiring daily inhalational or oral bronchodilator therapy and inhalational anti-inflammatory medication. I use both an Albuterol Inhaler and a Spiriva Inhaler to manage my symptoms.
	//
	//This condition severely affects my daily life. I experience shortness of breath that limits my physical activities, and chronic coughing or wheezing impacts my social interactions. These symptoms significantly disrupt my ability to engage in normal daily activities and maintain an active lifestyle.
	//
	//My sleep is frequently disturbed by nighttime coughing fits that cause me to wake up. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring physical exertion and experience reduced productivity due to fatigue and breathing difficulties. The constant management of my symptoms affects my ability to maintain consistent work performance.
	//
	//The persistent symptoms and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for an increase in benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Low back pain
	//Current Treatment Facility:
	//Current Medication: OTC Tylenol, OTC Motrin
	//
	//SERVICE CONNECTION: My service treatment records reflect that I suffered from LBP (Low Back Pain) while on active duty. All legal requirements for establishing service connection for LBP (Low Back Pain) have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of low back pain. This condition has been a significant challenge in my life since its onset during my active-duty service in 1994. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of low back pain have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While in Basic Training, I developed low back pain from the intense physical fitness requirements. The condition was severe enough that I went to medical when I also experienced lightheadedness. My back pain continued throughout my service and has progressively worsened over time. It's important to note that I did not have any low back pain before Basic Training.
	//
	//Currently, my back condition has worsened both in frequency and severity since its onset. I experience weakened movement, pain on movement, and have significantly less movement than normal. The pain interferes with both sitting and standing, and I'm easily fatigued. I find myself guarding to avoid flare-ups, and my walking motion is affected. I also experience incoordination, making it difficult to execute movements smoothly.
	//
	//My range of motion is severely limited. During normal conditions, I can barely bend forward, and during flare-ups, I cannot bend forward at all. I cannot bend backwards at all, even under normal conditions. When attempting to bend to either side, I have very limited movement, and during flare-ups, this becomes even more restricted. Additionally, when trying to twist my body, I have very limited movement that becomes even more restricted during flare-ups.
	//
	//I experience moderate bilateral radiculopathy, with symptoms including sharp, shooting pain that radiates along the affected nerve paths, numbness or decreased sensation in the areas supplied by the nerves, and tingling or pins and needles sensations in both legs.
	//
	//This condition severely affects my daily life. I face significant challenges in sitting or standing for extended periods. The pain and reduced mobility impact my intimate relationships, and I struggle to play with my children. I occasionally wear a back brace to manage the condition.
	//
	//My sleep is frequently disturbed, with frequent waking due to pain or discomfort. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I face challenges in maintaining posture during long periods of sitting, have difficulty with jobs requiring heavy lifting or physical labor, and experience limitations in work-related travel due to difficulty sitting for long periods.
	//
	//The persistent pain and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Right knee sprain with limitation of flexion and extension
	//Current Treatment Facility:
	//Current Medication: OTC Tylenol, OTC Motrin
	//
	//SERVICE CONNECTION: My service treatment records reflect that I suffered from Right Knee Contusion/Sprain while on active duty. All legal requirements for establishing service connection for Right Knee Contusion/Sprain have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of right knee sprain with limitation of flexion and extension. This condition has been a significant challenge in my life since its onset during my active-duty service in 1997. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of my knee pain have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While stationed on the USS Vandergrift during a port call at Seal Beach, CA, I was rollerblading with another sailor when I fell and twisted my right knee off to the side. My shipmate helped me get back to the ship where I was seen by medical. I did not have any knee issues before service, and my knee pain has progressively worsened over time.
	//
	//Currently, my knee condition has worsened both in frequency and severity since its onset. I experience incoordination, weakened movement, and pain on movement. I have significantly less movement than normal, though sometimes I experience more movement than normal, which is concerning. The condition interferes with both sitting and standing, and I'm easily fatigued. My walking motion is affected, and I find myself guarding to avoid flare-ups. I also experience knee "locking" and sudden sharp pain when putting weight on the knee.
	//
	//My range of motion is severely limited. Under normal conditions, I can barely bend my knee backward, and during flare-ups, I can hardly move it at all. When trying to straighten my leg, I have significant limitations that persist even during flare-ups. I also experience moderate instability in my right knee.
	//
	//This condition severely affects my daily life. I experience pain when walking, especially on stairs or inclines. I have limitations in playing with my children, and bending down hurts my knees. To manage the condition, I occasionally wear a knee brace.
	//
	//My sleep is frequently disturbed, requiring supportive pillows or devices to elevate my legs. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged standing or walking, face challenges in roles requiring frequent kneeling or squatting, and experience limitations in jobs requiring lifting or carrying heavy objects.
	//
	//The persistent pain and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Left foot plantar fasciitis
	//Current Treatment Facility:
	//Current Medication: OTC Tylenol, OTC Motrin
	//
	//SERVICE CONNECTION: My service treatment records reflect that I suffered from Left foot pain while on active duty. All legal requirements for establishing service connection for Left foot pain have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of left foot plantar fasciitis. This condition has been a significant challenge in my life since its onset during my active-duty service in 1994. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of plantar fasciitis have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//During Basic Training at Great Lakes, IL, we were required to run and march in our issued "boon dockers" shoes. The intense physical demands of these activities led to severe foot pain in my left foot. I went to medical where they diagnosed me as having plantar fasciitis. I did not have any foot issues before service, and the pain has progressively worsened over time.
	//
	//Currently, my condition has worsened both in frequency and severity since its onset. I experience a constant sharp pain on the bottom of my foot. I have less movement than normal and weakened movement. I experience pain on both movement and weight-bearing, as well as pain during non-weight-bearing activities. The condition interferes with standing, and I'm easily fatigued. I also experience incoordination, making it difficult to execute movements smoothly.
	//
	//This condition severely affects my daily life. I experience significant pain when walking or standing for extended periods. The condition impacts my posture and gait, potentially affecting other body parts. I have limitations in activities like dancing or running. To manage the condition, I wear over-the-counter arch supports and specially designed supportive shoes. Without these supports, the pain becomes intolerable.
	//
	//My sleep is frequently disturbed due to foot discomfort and pain. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged standing or walking, and I need frequent breaks to rest my feet.
	//
	//The persistent pain and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Right foot plantar fasciitis secondary to left foot plantar fasciitis
	//Current Treatment Facility:
	//Current Medication: OTC Tylenol, OTC Motrin
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of right foot plantar fasciitis secondary to left foot plantar fasciitis. This condition has been a significant challenge in my life since its onset during my active-duty service in 1994. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of plantar fasciitis have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//During Basic Training at Great Lakes, IL, I developed severe foot pain in my left foot from running and marching in our issued "boon dockers" shoes. While still in Basic Training, my right foot began to develop the same symptoms as my left foot. Due to the warrior mindset instilled in me early on with a large emphasis on toughness and self-reliance, I did not seek medical attention for my right foot as I was already being treated for the same issue in my left foot and didn't want to get into more trouble for missing training.
	//
	//Currently, my condition has worsened both in frequency and severity since its onset. I experience a constant sharp pain on the bottom of my foot. I have less movement than normal and weakened movement. I experience pain on both movement and weight-bearing, as well as pain during non-weight-bearing activities. The condition interferes with standing, and I'm easily fatigued. I also experience incoordination, making it difficult to execute movements smoothly.
	//
	//This condition severely affects my daily life. I experience significant pain when walking or standing for extended periods. The condition impacts my posture and gait, potentially affecting other body parts. I have limitations in activities like dancing or running. To manage the condition, I wear over-the-counter arch supports and specially designed supportive shoes. Without these supports, the pain becomes intolerable.
	//
	//My sleep is frequently disturbed due to foot discomfort and pain. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged standing or walking, and I need frequent breaks to rest my feet.
	//
	//The persistent pain and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Gastroesophageal reflux disease (GERD)
	//Current Treatment Facility:
	//Current Medication: Omeprazole 40 MG, Gaviscon
	//
	//SERVICE CONNECTION: My service treatment records reflect that I suffered from Frequent Indigestion while on active duty. All legal requirements for establishing service connection for Frequent Indigestion have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of gastroesophageal reflux disease (GERD). This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of GERD have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While deployed aboard the USS Vandergrift to the Persian Gulf, we made port call in Hong Kong in 1996. During liberty, I ate local food that I believe was contaminated, as I became severely ill for several days afterward, experiencing vomiting and extreme diarrhea. After the initial recovery, I started developing GERD symptoms including acid reflux, heartburn, regurgitation, and difficulty swallowing, with food often getting lodged in my throat. Throughout the remainder of my service, I reported frequent indigestion with these GERD symptoms and treated them with over-the-counter medications.
	//
	//Currently, my GERD has worsened both in frequency and severity since its onset. I experience daily reflux and regurgitation, which continues even during sleep. I suffer from nausea about half a dozen times a year, and I vomit two to three times a week, with episodes typically lasting less than a day. More seriously, I experience vomiting blood two or three days a year, and I have dark stools containing digested blood two or three times a month, lasting several days or more.
	//
	//The condition severely affects my daily life. I have difficulty swallowing, which impacts my enjoyment of meals, and I experience chronic pain and discomfort in my chest area. My sleep is significantly disrupted due to reflux when lying down, forcing me to elevate the head of my bed, which affects my sleep posture. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//My condition has significantly impacted my ability to work, with reduced productivity due to constant discomfort and pain. The persistent symptoms and medical interventions have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life. My GERD has worsened to the point that I have required several procedures, including a dilation endoscopy in 1994, which was unsuccessful. My doctor is now recommending surgery due to the failure of other treatments.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Irritable Bowel Syndrome (IBS)
	//Current Treatment Facility:
	//Current Medication: Linzess 72 MCG
	//
	//SERVICE CONNECTION: My treatment records reflect that I have a diagnosis of IBS. Additionally, I served in Southwest Asia during the Persian Gulf War Era, and I am entitled to the application of presumptive provisions of 38 U.S. Code § 1117. All legal requirements for establishing service connection for IBS have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of Irritable Bowel Syndrome (IBS). This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of IBS have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While stationed aboard the USS Vandergrift, we deployed to the Persian Gulf in 1995 as part of Operation Desert Strike. After returning from this deployment, I began experiencing functional gastrointestinal issues including abdominal pain and distress, constipation, gas, and bloating. I was later stationed aboard the USS Shiloh and deployed to the Persian Gulf again in 1997, where these symptoms continued.
	//
	//I reported these gastrointestinal issues to medical staff, but they focused primarily on treating my hemorrhoids. After leaving the service, my symptoms continued to worsen, and I eventually received an IBS diagnosis from my primary care physician. It's important to note that I did not have any functional gastrointestinal issues before deploying to the Persian Gulf.
	//
	//Currently, my IBS has worsened both in frequency and severity since its onset. I experience nearly constant episodes occurring daily, with symptoms including alternating diarrhea and constipation. These symptoms significantly impact my daily functioning and quality of life.
	//
	//This condition severely affects my daily life. The symptoms are unpredictable, disrupting my daily plans, and I experience regular abdominal pain and discomfort. I must take Linzess daily and rush to the bathroom when it takes effect, as I remain constipated and cannot defecate without this medication.
	//
	//My sleep is frequently disturbed by nighttime abdominal pain and discomfort. This chronic sleep disturbance compounds the challenges I face during the day, leading to fatigue and decreased overall well-being.
	//
	//The condition has significantly impacted my ability to work. I require frequent, sometimes urgent bathroom breaks which disrupt my workflow. My productivity is reduced due to the discomfort and necessary bathroom breaks, and I face challenges in maintaining focus due to my symptoms.
	//
	//The persistent symptoms and disruptions to my daily routine have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Left knee pain with limitation of flexion and extension secondary to right knee sprain
	//Current Treatment Facility:
	//Current Medication: OTC Tylenol, OTC Motrin
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of left knee pain with limitation of flexion and extension secondary to right knee sprain. This condition has been a significant challenge in my life since its onset during my active-duty service in 1997. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of my knee pain have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While stationed on the USS Vandergrift at Seal Beach, CA, I injured my right knee in a rollerblading accident. My left knee began to hurt almost immediately as I had to support and compensate for my injured right knee. Initially, I didn't seek treatment for my left knee as it started as a minor annoyance, but the condition has since escalated significantly. I did not have any knee issues before service.
	//
	//Currently, my knee condition has worsened both in frequency and severity since its onset. I experience incoordination, weakened movement, and pain on movement. I have significantly less movement than normal, though sometimes I experience more movement than normal, which is concerning. The condition interferes with both sitting and standing, and I'm easily fatigued. My walking motion is affected, and I find myself guarding to avoid flare-ups. I also experience knee "locking" and sudden sharp pain when putting weight on the knee.
	//
	//My range of motion is severely limited. Under normal conditions, I can barely bend my knee backward, and during flare-ups, I can hardly move it at all. When trying to straighten my leg, I have significant limitations that persist even during flare-ups. I also experience moderate instability in my left knee.
	//
	//This condition severely affects my daily life. I experience pain when walking, especially on stairs or inclines. I have limitations in playing with my children, and bending down hurts my knees. To manage the condition, I occasionally wear a knee brace.
	//
	//My sleep is frequently disturbed, requiring supportive pillows or devices to elevate my legs. This chronic sleep disturbance compounds the challenges I face during the day.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged standing or walking, face challenges in roles requiring frequent kneeling or squatting, and experience limitations in jobs requiring lifting or carrying heavy objects.
	//
	//The persistent pain and limitations have caused me to lose interest in activities I once enjoyed, significantly diminishing my quality of life.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.
	//
	//
	//Name of Disability/Condition: Hemorrhoids
	//Current Treatment Facility:
	//Current Medication: Preparation H, Sitz bath, Tucks
	//
	//SERVICE CONNECTION: My service treatment records reflect that I suffered from External Hemorrhoids while on active duty. All legal requirements for establishing service connection for External Hemorrhoids have been met; service connection for such disease is warranted.
	//
	//I am respectfully requesting Veteran Affairs benefits for my condition of hemorrhoids. This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of hemorrhoids have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	//
	//While deployed aboard the USS Vandergrift to the Persian Gulf and later during a port call in Hong Kong in 1996, I ate local food that I believe was contaminated. I became severely ill for several days, experiencing vomiting, extreme diarrhea, and straining to defecate. I developed pain in my rectum and went to medical, where they documented that I had hemorrhoids. I did not have hemorrhoids before this deployment, and since then, my condition has progressively worsened.
	//
	//Currently, my hemorrhoids have worsened both in frequency and severity since their onset. My symptoms are moderate in nature, characterized by large or thrombotic, irreducible hemorrhoids with excessive redundant tissue and frequent recurrences. I experience persistent bleeding along with severe pain, swelling, and itching. I have undergone radiation treatment to remove them, but the surgeon would not perform surgery due to their large size.
	//
	//This condition severely affects my daily life. I experience chronic pain and discomfort that affects my overall well-being, and I have difficulty sitting for extended periods. The persistent symptoms have caused me to lose interest in activities I once enjoyed.
	//
	//My sleep is frequently disturbed due to discomfort and pain, leading to fatigue and irritability during the day. This chronic sleep disturbance compounds the challenges I face in my daily life.
	//
	//The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged sitting and physical exertion. These limitations affect my productivity and career opportunities.
	//
	//I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.`

	userInput := "Reference materials:\n" + referenceContent + "\n\n"
	userInput += "VA statement(s):\n" + statements

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

	claude3Response, aiResultId, err = c.Awsclaude3Usecase.AskInvoke(ctx, payload, AiForm_GenStatement, promptKey, false, nil)
	if err != nil {
		return
	}
	if len(claude3Response.ResponseContent) > 0 {
		parseResult = claude3Response.ResponseContent[0].Text
	}
	return
}
