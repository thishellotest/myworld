package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AiAssistantJobBuzUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	AiAssistantJobUsecase     *AiAssistantJobUsecase
	AiTaskUsecase             *AiTaskUsecase
	AiResultUsecase           *AiResultUsecase
	AssistantUsecase          *AssistantUsecase
	StatementUsecase          *StatementUsecase
	TUsecase                  *TUsecase
	StatementConditionUsecase *StatementConditionUsecase
	MapUsecase                *MapUsecase
	DocEmailUsecase           *DocEmailUsecase
}

func NewAiAssistantJobBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	AiAssistantJobUsecase *AiAssistantJobUsecase,
	AiTaskUsecase *AiTaskUsecase,
	AiResultUsecase *AiResultUsecase,
	AssistantUsecase *AssistantUsecase,
	StatementUsecase *StatementUsecase,
	TUsecase *TUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	MapUsecase *MapUsecase,
	DocEmailUsecase *DocEmailUsecase,
) *AiAssistantJobBuzUsecase {
	uc := &AiAssistantJobBuzUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		AiAssistantJobUsecase:     AiAssistantJobUsecase,
		AiTaskUsecase:             AiTaskUsecase,
		AiResultUsecase:           AiResultUsecase,
		AssistantUsecase:          AssistantUsecase,
		StatementUsecase:          StatementUsecase,
		TUsecase:                  TUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
		MapUsecase:                MapUsecase,
		DocEmailUsecase:           DocEmailUsecase,
	}

	return uc
}

func (c *AiAssistantJobBuzUsecase) HttpGetJobDetail(ctx *gin.Context) {
	reply := CreateReply()
	jobUuid := ctx.Query("job_uuid")
	data, err := c.BizHttpGetJobDetail(jobUuid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiAssistantJobBuzUsecase) HttpApplyJob(ctx *gin.Context) {
	reply := CreateReply()
	jobUuid := ctx.Query("job_uuid")
	data, err := c.BizHttpApplyJob(jobUuid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiAssistantJobBuzUsecase) BizHttpApplyJob(jobUuid string) (lib.TypeMap, error) {
	data := make(lib.TypeMap)

	entity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("The task does not exist.")
	}
	entity.JobStatus = AiAssistantJob_JobStatus_Normal
	err = c.CommonUsecase.DB().Model(&entity).Update("job_status", AiAssistantJob_JobStatus_Normal).Error
	if err != nil {
		return nil, err
	}

	var vo AiAssistantJobDetail
	var aiAssistantJobStatusVo AiAssistantJobStatusVo
	vo.JobUuid = jobUuid
	vo = entity.ToJobDetail(c.AiTaskUsecase, c.AiResultUsecase)
	aiAssistantJobStatusVo = entity.ToJobStatusInfo()
	_, aiResultEntity := entity.GetAiResultEntity(c.AiTaskUsecase, c.AiResultUsecase)
	if aiResultEntity != nil {
		aiAssistantJobInput := entity.ToAiAssistantJobInput()
		if aiAssistantJobInput.AssistantBiz == AiAssistantJobBizType_statementCondition {
			tCase, tClient, statementConditionEntity, _ := c.AssistantUsecase.ExplainFormatJobUuidForStatementCondition(jobUuid)

			var parseAiStatementConditionVo ParseAiStatementConditionVo
			if aiAssistantJobInput.BizType == AiAssistantBizType_StandardHeaderRevision {
				str := GetJsonFromAiResultForAssistant(aiResultEntity.ParseResult)
				shrAiResultVo := StringToStandardHeaderRevisionAiResult(str)
				parseAiStatementConditionVo = shrAiResultVo.ToParseAiStatementConditionVo()
			} else {
				parseAiStatementConditionVo = ParseAiStatementCondition(aiResultEntity.ParseResult)
			}
			err := c.StatementUsecase.UpdateOneConditionStatement(*tCase, *tClient, *statementConditionEntity, parseAiStatementConditionVo)
			if err != nil {
				c.log.Error(err)
			}
		} else if aiAssistantJobInput.AssistantBiz == AiAssistantJobBizType_genDocEmail {

			jobUuidCommon := FormatJobUuidForCommon(jobUuid)
			if jobUuidCommon.CaseId <= 0 {
				return nil, errors.New("jobUuidCommon.CaseId is wrong")
			}

			er := c.DocEmailUsecase.SetLatestDocEmailResult(jobUuidCommon.CaseId, aiResultEntity.ID)
			if er != nil {
				c.log.Error(er)
			}

		} else if aiAssistantJobInput.AssistantBiz == AiAssistantJobBizType_statementSection {
			parseResultString := GetJsonFromAiResultForAssistant(aiResultEntity.ParseResult)
			parseResultStringMap := lib.ToTypeMapByString(parseResultString)
			if parseResultStringMap.GetString("update_required") == "true" {
				updatedStatement := strings.TrimSpace(parseResultStringMap.GetString("updated_statement"))
				if updatedStatement != "" {
					tCase, tClient, statementConditionEntity, vo := c.AssistantUsecase.ExplainJobUuidForStatementSection(jobUuid)
					if tCase != nil && tClient != nil && statementConditionEntity != nil {
						err := c.StatementUsecase.UpdateOneSectionStatement(*tCase, *tClient, *statementConditionEntity, vo.SectionType, updatedStatement)
						if err != nil {
							c.log.Error(err)
						}
					}
				}
			}
		} else {
			return nil, errors.New("aiAssistantJobInput.AssistantBiz is wrong")
		}
	}

	data.Set("job_detail", vo)
	data.Set("job_status", aiAssistantJobStatusVo)

	return data, nil
}

func (c *AiAssistantJobBuzUsecase) HttpClearJob(ctx *gin.Context) {
	reply := CreateReply()
	jobUuid := ctx.Query("job_uuid")
	data, err := c.BizHttpClearJob(jobUuid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiAssistantJobBuzUsecase) BizHttpClearJob(jobUuid string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	entity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("The task does not exist.")
	}
	entity.JobStatus = AiAssistantJob_JobStatus_Normal
	err = c.CommonUsecase.DB().Model(&entity).Update("job_status", AiAssistantJob_JobStatus_Normal).Error
	if err != nil {
		return nil, err
	}

	var vo AiAssistantJobDetail
	var aiAssistantJobStatusVo AiAssistantJobStatusVo
	vo.JobUuid = jobUuid
	vo = entity.ToJobDetail(c.AiTaskUsecase, c.AiResultUsecase)
	aiAssistantJobStatusVo = entity.ToJobStatusInfo()

	data.Set("job_detail", vo)
	data.Set("job_status", aiAssistantJobStatusVo)

	return data, nil
}

func (c *AiAssistantJobBuzUsecase) BizHttpGetJobDetail(jobUuid string) (lib.TypeMap, error) {

	entity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
	if err != nil {
		return nil, err
	}
	var vo AiAssistantJobDetail
	var aiAssistantJobStatusVo AiAssistantJobStatusVo
	vo.JobUuid = jobUuid
	if entity != nil {
		vo = entity.ToJobDetail(c.AiTaskUsecase, c.AiResultUsecase)
		aiAssistantJobStatusVo = entity.ToJobStatusInfo()
	}
	vo.HandleQuickOptions()

	data := make(lib.TypeMap)
	data.Set("job_detail", vo)
	data.Set("job_status", aiAssistantJobStatusVo)
	return data, nil
}

func (c *AiAssistantJobBuzUsecase) HttpGetJobStatus(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	jobUuids, _ := lib.ConvertTypeListInterface[string](body.GetTypeListInterface("job_uuids"))
	data, err := c.BizHttpGetJobStatus(jobUuids)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiAssistantJobBuzUsecase) BizHttpGetJobStatus(jobUuids []string) (lib.TypeMap, error) {
	if len(jobUuids) == 0 {
		return nil, nil
	}

	jobs, err := c.AiAssistantJobUsecase.AllByCond(In("job_uuid", jobUuids))
	if err != nil {
		return nil, err
	}
	list := make(map[string]AiAssistantJobStatusVo)
	for _, v := range jobs {
		list[v.JobUuid] = v.ToJobStatusInfo()
	}

	//conditionJobs, err := c.AiTaskUsecase.GetGenerateStatementForJobStatusWithJobUuids(jobUuids)
	//if err != nil {
	//	return nil, err
	//}
	//for k, v := range conditionJobs {
	//	list[k] = v
	//}

	data := make(lib.TypeMap)
	data.Set("job_status_map", list)
	return data, nil
}

func (c *AiAssistantJobBuzUsecase) HttpCreate(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	jobUuid := ctx.Query("job_uuid")
	var aiAssistantJobInput AiAssistantJobInput
	aiAssistantJobInput = lib.BytesToTDef(rawData, aiAssistantJobInput)
	data, err := c.BizHttpCreate(jobUuid, aiAssistantJobInput)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiAssistantJobBuzUsecase) BizHttpCreateForStatementCondition(jobUuid string) (lib.TypeMap, error) {

	var jobUuidForStatementCondition JobUuidForStatementCondition
	jobUuidForStatementCondition = FormatJobUuidForStatementCondition(jobUuid)

	tCase, err := c.TUsecase.DataById(Kind_client_cases, jobUuidForStatementCondition.CaseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	statementConditionEntity, err := c.StatementConditionUsecase.GetCondition(jobUuidForStatementCondition.CaseId, jobUuidForStatementCondition.StatementConditionId)
	if err != nil {
		return nil, err
	}
	if statementConditionEntity == nil {
		return nil, errors.New("statementConditionEntity is nil")
	}

	aiTaskEntity, err := c.AiTaskUsecase.GetGenerateStatement(jobUuidForStatementCondition.CaseId, jobUuidForStatementCondition.StatementConditionId)
	if err != nil {
		return nil, err
	}
	needCreateTask := false
	if aiTaskEntity != nil {
		if aiTaskEntity.HandleStatus == HandleStatus_done {

			aiTaskEntity.DeletedAt = time.Now().Unix()
			err := c.CommonUsecase.DB().Save(&aiTaskEntity).Error
			if err != nil {
				return nil, err
			}
			needCreateTask = true
		}
	} else {
		needCreateTask = true
	}
	if needCreateTask {
		handleStatus := AiTask_HandleStatus_In_process
		_, err := c.AiTaskUsecase.CreateGenerateStatement(tCase, statementConditionEntity.ToStatementCondition(), &handleStatus)
		if err != nil {
			return nil, err
		}
	}

	aiAssistantJobStatusVo, err := c.AiTaskUsecase.GetGenerateStatementForJobStatus(jobUuidForStatementCondition.CaseId, jobUuidForStatementCondition.StatementConditionId)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("job_status", aiAssistantJobStatusVo)

	return data, nil
}

func (c *AiAssistantJobBuzUsecase) HandleAllStatements(caseId int32, userInputText string) error {

	conditions, err := c.StatementConditionUsecase.AllConditions(caseId)
	if err != nil {
		return err
	}

	for k, _ := range conditions {
		err := c.HandleStatementConditionJob(*conditions[k], userInputText)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AiAssistantJobBuzUsecase) HandleStatementConditionJob(statementConditionEntity StatementConditionEntity, userInputText string) error {

	jobUuid := statementConditionEntity.ToStatementConditionJobUuid()
	jobEntity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
	if err != nil {
		return err
	}
	if jobEntity == nil {
		jobEntity = &AiAssistantJobEntity{
			JobUuid:   jobUuid,
			CreatedAt: time.Now().Unix(),
		}
	}

	var aiTaskInputAssistant AiTaskInputAssistant
	aiTaskInputAssistant.AssistantBiz = AiAssistantJobBizType_statementCondition
	aiTaskInputAssistant.BizType = AiAssistantBizType_CreateNewStatement
	aiTaskInputAssistant.InternalBizType = AssistantInternalBizType_ForAllStatemts
	aiTaskInputAssistant.UserInputText = userInputText
	aiTaskInputAssistant.CaseId = statementConditionEntity.CaseId
	aiTaskInputAssistant.StatementConditionId = statementConditionEntity.ID

	taskEntity, err := c.AiTaskUsecase.CreateTask(AiTaskFromType_Assistant, "", statementConditionEntity.CaseId, 0, InterfaceToString(aiTaskInputAssistant), "", nil, nil)
	if err != nil {
		return err
	}
	if taskEntity == nil {
		return errors.New("Failed to create the task")
	}

	var aiAssistantJobInput AiAssistantJobInput
	aiAssistantJobInput.AssistantBiz = aiTaskInputAssistant.AssistantBiz
	aiAssistantJobInput.InternalBizType = aiTaskInputAssistant.InternalBizType
	aiAssistantJobInput.UserInputText = aiTaskInputAssistant.UserInputText
	aiAssistantJobInput.BizType = aiTaskInputAssistant.BizType

	jobEntity.AiTaskId = taskEntity.ID
	jobEntity.JobStatus = AiAssistantJob_JobStatus_Running
	jobEntity.UpdatedAt = time.Now().Unix()
	jobEntity.JobInput = InterfaceToString(aiAssistantJobInput)
	err = c.CommonUsecase.DB().Save(&jobEntity).Error
	return err
}

func (c *AiAssistantJobBuzUsecase) BizHttpCreate(jobUuid string, aiAssistantJobInput AiAssistantJobInput) (lib.TypeMap, error) {
	can, jobEntity, err := c.AiAssistantJobUsecase.CanRunNewJob(jobUuid)
	if err != nil {
		return nil, err
	}
	if !can {
		return nil, errors.New("The task already exists.")
	}
	var aiTaskInputAssistant AiTaskInputAssistant
	needCreateTask := false
	aiAssistantJobInput.UserInputText = strings.TrimSpace(aiAssistantJobInput.UserInputText)
	if strings.Index(jobUuid, AiAssistantJobBizType_statementSection) == 0 {
		needCreateTask = true
		aiAssistantJobInput.AssistantBiz = AiAssistantJobBizType_statementSection
		var obUuidForStatementSection JobUuidForStatementSection
		obUuidForStatementSection = FormatJobUuidForStatementSection(jobUuid)
		if obUuidForStatementSection.CaseId <= 0 || obUuidForStatementSection.StatementConditionId <= 0 {
			return nil, errors.New("The format of the UUID is incorrect")
		}
		if aiAssistantJobInput.BizType == "" && aiAssistantJobInput.UserInputText == "" {
			return nil, errors.New("Please enter the prompt")
		}
		aiTaskInputAssistant.AssistantBiz = aiAssistantJobInput.AssistantBiz
		aiTaskInputAssistant.BizType = aiAssistantJobInput.BizType
		aiTaskInputAssistant.UserInputText = aiAssistantJobInput.UserInputText
		aiTaskInputAssistant.CaseId = obUuidForStatementSection.CaseId
		aiTaskInputAssistant.StatementConditionId = obUuidForStatementSection.StatementConditionId
		aiTaskInputAssistant.SectionType = obUuidForStatementSection.SectionType

	} else if strings.Index(jobUuid, AiAssistantJobBizType_statementCondition) == 0 {
		needCreateTask = true
		aiAssistantJobInput.AssistantBiz = AiAssistantJobBizType_statementCondition
		if aiAssistantJobInput.BizType == "" {
			return nil, errors.New("Quick Options is required here")
		}
		var jobUuidForStatementCondition JobUuidForStatementCondition
		jobUuidForStatementCondition = FormatJobUuidForStatementCondition(jobUuid)
		aiTaskInputAssistant.AssistantBiz = aiAssistantJobInput.AssistantBiz
		aiTaskInputAssistant.BizType = aiAssistantJobInput.BizType
		aiTaskInputAssistant.UserInputText = aiAssistantJobInput.UserInputText
		aiTaskInputAssistant.CaseId = jobUuidForStatementCondition.CaseId
		aiTaskInputAssistant.StatementConditionId = jobUuidForStatementCondition.StatementConditionId

	} else if strings.Index(jobUuid, AiAssistantJobBizType_allStatements) == 0 {
		needCreateTask = false
		aiAssistantJobInput.AssistantBiz = AiAssistantJobBizType_allStatements
		if aiAssistantJobInput.BizType == "" {
			return nil, errors.New("Quick Options is required here")
		}

		var jobUuidVo JobUuidForAllStatements
		jobUuidVo = FormatJobUuidForAllStatements(jobUuid)
		aiTaskInputAssistant.AssistantBiz = aiAssistantJobInput.AssistantBiz
		aiTaskInputAssistant.BizType = aiAssistantJobInput.BizType
		aiTaskInputAssistant.UserInputText = aiAssistantJobInput.UserInputText
		aiTaskInputAssistant.CaseId = jobUuidVo.CaseId

		err = c.HandleAllStatements(aiTaskInputAssistant.CaseId, aiTaskInputAssistant.UserInputText)
		if err != nil {
			return nil, err
		}
	} else if strings.Index(jobUuid, AiAssistantJobBizType_genDocEmail) == 0 {
		needCreateTask = true
		aiAssistantJobInput.AssistantBiz = AiAssistantJobBizType_genDocEmail
		if aiAssistantJobInput.BizType == "" {
			return nil, errors.New("Quick Options is required here")
		}
		var jobUuidVo JobUuidForCommon
		jobUuidVo = FormatJobUuidForCommon(jobUuid)
		aiTaskInputAssistant.AssistantBiz = aiAssistantJobInput.AssistantBiz
		aiTaskInputAssistant.BizType = aiAssistantJobInput.BizType
		aiTaskInputAssistant.CaseId = jobUuidVo.CaseId
		aiTaskInputAssistant.UserInputText = aiAssistantJobInput.UserInputText

	} else {
		return nil, errors.New("The format of the UUID is incorrect")
	}

	if jobEntity == nil {
		jobEntity = &AiAssistantJobEntity{
			JobUuid:   jobUuid,
			CreatedAt: time.Now().Unix(),
		}
	}

	if needCreateTask {
		taskEntity, err := c.AiTaskUsecase.CreateTask(AiTaskFromType_Assistant, "", aiTaskInputAssistant.CaseId, 0, InterfaceToString(aiTaskInputAssistant), "", nil, nil)
		if err != nil {
			return nil, err
		}
		if taskEntity == nil {
			return nil, errors.New("Failed to create the task")
		}

		jobEntity.AiTaskId = taskEntity.ID
	}

	jobEntity.JobStatus = AiAssistantJob_JobStatus_Running
	jobEntity.UpdatedAt = time.Now().Unix()
	jobEntity.JobInput = InterfaceToString(aiAssistantJobInput)
	err = c.CommonUsecase.DB().Save(&jobEntity).Error
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("job_status", jobEntity.ToJobStatusInfo())
	jobDetail := jobEntity.ToJobDetail(c.AiTaskUsecase, c.AiResultUsecase)
	data.Set("job_detail", jobDetail)

	return data, nil
}
