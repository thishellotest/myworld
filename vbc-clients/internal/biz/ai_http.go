package biz

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AiHttpUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	JWTUsecase                *JWTUsecase
	Awsclaude3Usecase         *Awsclaude3Usecase
	AiResultUsecase           *AiResultUsecase
	AiUsecase                 *AiUsecase
	ConditionbuzUsecase       *ConditionbuzUsecase
	TUsecase                  *TUsecase
	AiTaskUsecase             *AiTaskUsecase
	DataComboUsecase          *DataComboUsecase
	StatementUsecase          *StatementUsecase
	RelasLogUsecase           *RelasLogUsecase
	ConditionUsecase          *ConditionUsecase
	CommonUsecase             *CommonUsecase
	PsHttpUsecase             *PsHttpUsecase
	StatementConditionUsecase *StatementConditionUsecase
}

func NewAiHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	AiResultUsecase *AiResultUsecase,
	AiUsecase *AiUsecase,
	ConditionbuzUsecase *ConditionbuzUsecase,
	TUsecase *TUsecase,
	AiTaskUsecase *AiTaskUsecase,
	DataComboUsecase *DataComboUsecase,
	StatementUsecase *StatementUsecase,
	RelasLogUsecase *RelasLogUsecase,
	ConditionUsecase *ConditionUsecase,
	CommonUsecase *CommonUsecase,
	PsHttpUsecase *PsHttpUsecase,
	StatementConditionUsecase *StatementConditionUsecase) *AiHttpUsecase {
	return &AiHttpUsecase{
		log:                       log.NewHelper(logger),
		conf:                      conf,
		JWTUsecase:                JWTUsecase,
		Awsclaude3Usecase:         Awsclaude3Usecase,
		AiResultUsecase:           AiResultUsecase,
		AiUsecase:                 AiUsecase,
		ConditionbuzUsecase:       ConditionbuzUsecase,
		TUsecase:                  TUsecase,
		AiTaskUsecase:             AiTaskUsecase,
		DataComboUsecase:          DataComboUsecase,
		StatementUsecase:          StatementUsecase,
		RelasLogUsecase:           RelasLogUsecase,
		ConditionUsecase:          ConditionUsecase,
		CommonUsecase:             CommonUsecase,
		PsHttpUsecase:             PsHttpUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
	}
}

func (c *AiHttpUsecase) HandleOnceConditionSourceWithAi(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))

	go func() {
		for i := 1; i <= 30; i++ {
			err := c.ConditionbuzUsecase.HandleOnceConditionSourceWithAi()
			c.log.Info("AiHttpUsecase:HandleOnceConditionSourceWithAi:Done Times=" + InterfaceToString(i))
			if err != nil {
				c.log.Warn(err)
				break
			}
			time.Sleep(1)

		}

	}()
	//c.log.Info("parseResult:", parseResult)
	reply.Success()
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) TestAi(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))

	go func() {
		_, _, err := c.AiUsecase.ExecuteWithClaude3(body.GetString("prompt_key"), body)
		if err != nil {
			c.log.Warn(err)
		}
		c.log.Info("AiHttpUsecase:TestAi:Done")
	}()
	//c.log.Info("parseResult:", parseResult)
	reply.Success()
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) TaskLaunch(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	aiTaskId := body.GetInt("ai_task_id")
	data, err := c.BizTaskLaunch(aiTaskId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTaskLaunch(aiTaskId int32) (lib.TypeMap, error) {

	aiTaskEntity, err := c.AiTaskUsecase.GetByCond(Eq{"id": aiTaskId})
	if err != nil {
		return nil, err
	}
	if aiTaskEntity == nil {
		return nil, errors.New("aiTaskEntity is nil")
	}
	if aiTaskEntity.HandleStatus != AiTask_HandleStatus_Waiting {
		return nil, errors.New("aiTaskEntity.HandleStatus is wrong")
	}
	if aiTaskEntity.FromType != AiTaskFromType_statement {
		return nil, errors.New("aiTaskEntity.FromType is wrong")
	}
	aiTaskInputGenerateStatement := aiTaskEntity.GetAiTaskInputGenerateStatement()

	sourceId := ""

	tCase, err := c.TUsecase.DataById(Kind_client_cases, aiTaskEntity.CaseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	sourceId = InterfaceToString(aiTaskInputGenerateStatement.StatementConditionId)

	exists, err := c.RelasLogUsecase.ConditionExists(sourceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("ConditionExists: Please associate jotform first")
	}
	aiTaskEntity.HandleStatus = AiTask_HandleStatus_In_process
	aiTaskEntity.UpdatedAt = time.Now().Unix()
	err = c.CommonUsecase.DB().Save(&aiTaskEntity).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *AiHttpUsecase) TaskHandle(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	caseId := body.GetInt("case_id")
	data, err := c.BizTaskHandle(caseId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTaskHandle(caseId int32) (lib.TypeMap, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)

	if stages != config_vbc.Stages_StatementDrafts && stages != config_vbc.Stages_AmStatementDrafts {
		return nil, errors.New("The client case's stage is not \"Statement Drafts\"")
	}
	//statements := tCase.CustomFields.TextValueByNameBasic(FieldName_statements)
	//if statements == "" {
	//	return nil, errors.New("Statements is empty")
	//}

	a, err := c.AiTaskUsecase.GetByCond(Eq{"from_type": AiTaskFromType_generate_doc_email, "case_id": tCase.Id(), "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if a == nil {
		_, err = c.AiTaskUsecase.CreateGenerateDocEmail(tCase)
		if err != nil {
			return nil, err
		}
	}

	if configs.IsDev() || true { // todo:lgl 生产环境暂不开放

		statementConditionList, err := c.StatementConditionUsecase.AllConditions(caseId)
		if len(statementConditionList) == 0 {
			return nil, errors.New("Conditions has not been initialized yet")
		}
		//statementConditionList, err := SplitCaseStatements(statements)
		if err != nil {
			return nil, err
		}
		for k, v := range statementConditionList {
			entity, err := c.AiTaskUsecase.GetByCond(Eq{"from_type": AiTaskFromType_statement,
				"case_id":       tCase.Id(),
				"deleted_at":    0,
				"task_uniqcode": InterfaceToString(v.ID),
			})
			if err != nil {
				return nil, err
			}
			if entity == nil {

				handleStatus := AiTask_HandleStatus_Waiting
				if configs.NewPSGen {
					conditionUuid := InterfaceToString(v.ID)

					if configs.EnableAiAutoAssociationJotform {
						handleStatus = AiTask_HandleStatus_In_process
					} else {
						exists, err := c.RelasLogUsecase.ConditionExists(conditionUuid)
						if err != nil {
							return nil, err
						}
						if exists {
							handleStatus = AiTask_HandleStatus_In_process
						}
					}
				} else {
					conditionEntity, err := c.ConditionUsecase.ConditionGet(statementConditionList[k].ConditionValue)
					if err != nil {
						return nil, err
					}
					if conditionEntity != nil {
						exists, err := c.RelasLogUsecase.ConditionExists(InterfaceToString(conditionEntity.ID))
						if err != nil {
							return nil, err
						}
						if exists {
							handleStatus = AiTask_HandleStatus_In_process
						}
					}
				}
				_, err = c.AiTaskUsecase.CreateGenerateStatement(tCase, v.ToStatementCondition(), &handleStatus)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}

func (c *AiHttpUsecase) TaskResult(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	aiTaskId := body.GetInt("ai_task_id")
	data, err := c.BizTaskResult(aiTaskId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTaskResult(aiTaskId int32) (lib.TypeMap, error) {

	aiTask, err := c.AiTaskUsecase.GetByCond(Eq{"id": aiTaskId, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if aiTask == nil {
		return nil, errors.New("aiTask is nil")
	}
	if aiTask.CurrentResultId == 0 {
		return nil, errors.New("CurrentResultId is 0")
	}
	aiResult, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTask.CurrentResultId})
	if err != nil {
		return nil, err
	}
	if aiResult == nil {
		return nil, errors.New("aiResult is nil")
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, aiTask.CaseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	data := make(lib.TypeMap)
	if aiTask.FromType == AiTaskFromType_generate_doc_email {

		result := DocEmailExtractHealthIssues(aiResult.ParseResult)
		if len(result) > 0 {

			var wordLineList WordLineList
			wordLineList = append(wordLineList, WordDocEmailTop...)

			for _, v := range result {
				wordLineList = append(wordLineList, WordLine{Type: WordLine_Type_List,
					Value: v})
			}

			wordLineList = append(wordLineList, WordDocEmailBottom...)
			wordLineList = append(wordLineList, WordLine{Type: WordLine_Type_Normal,
				Value: tClient.CustomFields.TextValueByNameBasic(FieldName_full_name)})
			data.Set("data.parse_result", wordLineList.ToString())
		} else {
			data.Set("data.parse_result", aiResult.ParseResult)
		}
	} else if aiTask.FromType == AiTaskFromType_statement {
		text := aiResult.GetStatement()
		re := regexp.MustCompile(`(?m)^#{1,2}\s*`)
		// (?m) 开启多行模式，这样 ^ 代表每一行的开头。
		// 接着是 1 或 2 个 #（#{1,2}）
		// 后面也允许有空格（\s*）
		cleanText := re.ReplaceAllString(text, "")
		data.Set("data.parse_result", cleanText)
		data.Set("data.source_parse_result", text)
	} else {
		data.Set("data.parse_result", aiResult.GetStatement())
	}
	return data, nil
}

func StatementExtract(text string) string {
	return strings.TrimPrefix(text, "# ")
}

func (c *AiHttpUsecase) TaskDelete(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	aiTaskId := body.GetInt("ai_task_id")
	data, err := c.BizTaskDelete(aiTaskId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTaskDelete(aiTaskId int32) (lib.TypeMap, error) {

	aiTask, err := c.AiTaskUsecase.GetByCond(Eq{"id": aiTaskId, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if aiTask == nil {
		return nil, errors.New("aiTask is nil")
	}
	aiTask.DeletedAt = time.Now().Unix()
	err = c.AiTaskUsecase.CommonUsecase.DB().Save(&aiTask).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *AiHttpUsecase) TaskRenew(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	aiTaskId := body.GetInt("ai_task_id")
	data, err := c.BizTaskRenew(aiTaskId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTaskRenew(aiTaskId int32) (lib.TypeMap, error) {

	aiTask, err := c.AiTaskUsecase.GetByCond(Eq{"id": aiTaskId, "deleted_at": 0})
	if err != nil {
		return nil, err
	}
	if aiTask == nil {
		return nil, errors.New("aiTask is nil")
	}

	aiTask.DeletedAt = time.Now().Unix()
	err = c.AiTaskUsecase.CommonUsecase.DB().Save(&aiTask).Error
	if err != nil {
		return nil, err
	}
	//if aiTask.FromType == AiTaskFromType_update_statement {
	//	_, err = c.PsHttpUsecase.BizHandleUpdateStatement(aiTask.CaseId)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	newAiTask, err := c.AiTaskUsecase.CreateTask(aiTask.FromType, aiTask.FromCode, aiTask.CaseId, 0, aiTask.Input, aiTask.TaskUniqcode, nil, &aiTask.SerialNumber)
	if err != nil {
		return nil, err
	}
	if newAiTask == nil {
		return nil, errors.New("newAiTask is nil")
	}
	//}
	return nil, nil
}

func (c *AiHttpUsecase) Tasks(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	caseId := body.GetInt("case_id")
	data, err := c.BizTasks(caseId)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizTasks(caseId int32) (lib.TypeMap, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	var aiTaskList AiTaskList
	records, err := c.AiTaskUsecase.AllByCondWithOrderBy(And(Eq{"deleted_at": 0, "case_id": caseId}, In("from_type", []string{
		AiTaskFromType_generate_doc_email,
		AiTaskFromType_veteran_summary,
		AiTaskFromType_statement,
		AiTaskFromType_update_statement,
	})), "from_type desc, serial_number asc ,id desc", 1000)
	if err != nil {
		return nil, err
	}
	for _, v := range records {
		vo := v.ToAiTaskItem()
		if v.FromType == AiTaskFromType_statement {
			voInput := v.GetAiTaskInputGenerateStatement()
			statementConditionEntity, _ := c.StatementConditionUsecase.GetByCond(Eq{"id": voInput.StatementConditionId})
			if statementConditionEntity != nil {
				t := statementConditionEntity.ToStatementCondition()
				vo.Detail = t.ToOriginValue()
			}
		}
		aiTaskList = append(aiTaskList, vo)
	}
	typeMap := make(lib.TypeMap)
	typeMap.Set("tasks", aiTaskList)

	return typeMap, nil
}

func (c *AiHttpUsecase) Claude3(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))

	aiRequest := body.GetString("ai_request")
	promptKey := body.GetString("prompt_key")
	var claude3Request Claude3Request
	json.Unmarshal([]byte(aiRequest), &claude3Request)

	//c.log.Debug("Claude3:", aiRequest, InterfaceToString(claude3Request))

	// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	//lib.DPrintln(moduleName)
	go func() {
		_, err := c.BizClaude3(claude3Request, promptKey)
		if err != nil {
			c.log.Error(err)
		}
	}()

	//reply.Merge(data)
	reply.Success()

	ctx.JSON(200, reply)
}

func (c *AiHttpUsecase) BizClaude3(claude3Request Claude3Request, promptKey string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)

	//go func() {
	res, aiResultId, err := c.Awsclaude3Usecase.AskInvoke(context.TODO(), claude3Request, AiFrom_BizClaude3, promptKey, false, nil)
	if err != nil {
		c.log.Warn(err)
	}
	//c.log.Info("AiHttpUsecase err:", err)
	//c.log.Info("AiHttpUsecase:", res)

	data.Set("ai_response", InterfaceToString(res))
	data.Set("ai_result_id", aiResultId)

	return data, nil
}
