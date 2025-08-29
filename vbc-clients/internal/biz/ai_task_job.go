package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	Ai_Prompt3_0                           = "prompt3_0"
	Ai_Prompt4_0                           = "prompt4_0_muscul"
	Ai_prompt4_0_muscul_summary_outjosn    = "prompt4_0_muscul_summary_outjosn"
	Ai_Prompt6_0                           = "prompt6_0_muscul"
	Ai_prompt7_0                           = "prompt7_0"
	Ai_prompt_ps_update                    = "prompt_ps_update"
	Ai_prompt_ps_update3_0                 = "prompt_ps_update3_0"
	Ai_prompt_ps_coco_update3_0            = "prompt_ps_update3_0_coco"
	Ai_prompt_ps_coco_update3_0_V1         = "prompt_ps_update3_0_coco_v1"
	Ai_prompt_StandardHeaderRevisionPrompt = "prompt_StandardHeaderRevisionPrompt"
	Ai_prompt_8_0                          = "prompt8_0"
	Ai_prompt_10_0                         = "prompt10_0"
	Ai_prompt_11_0                         = "prompt11_0"
	Ai_prompt_11_1                         = "prompt11_1"

	Ai_prompt_musculoskeletal_checklist = "prompt_musculoskeletal_checklist"

	CurrentPromptForDocEmail = Ai_Prompt3_0 //  Ai_Prompt3_0 Ai_prompt_11_1

	CurrentPromptMusculoskeletalChecklist = Ai_prompt_musculoskeletal_checklist

	CurrentGenStatementPrompt = Ai_prompt_11_1 // 上一次版本是Ai_prompt_8_0->Ai_prompt_11_0

	// https://base.vetbenefitscenter.com/ps-gen?case_id=5373
	CurrentPromptPSUpdate = Ai_prompt_ps_coco_update3_0_V1

	Current_StandardHeaderRevisionPrompt = Ai_prompt_StandardHeaderRevisionPrompt
)

type AiTaskJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[AiTaskEntity]
	BaseHandleCustom[AiTaskEntity]
	AiStatementUsecase        *AiStatementUsecase
	TUsecase                  *TUsecase
	AiUsecase                 *AiUsecase
	JotformSubmissionUsecase  *JotformSubmissionUsecase
	DocEmailUsecase           *DocEmailUsecase
	Awsclaude3Usecase         *Awsclaude3Usecase
	StatementUsecase          *StatementUsecase
	DataComboUsecase          *DataComboUsecase
	AiTaskUsecase             *AiTaskUsecase
	AiResultUsecase           *AiResultUsecase
	AiTaskbuzUsecase          *AiTaskbuzUsecase
	FeeUsecase                *FeeUsecase
	StatementConditionUsecase *StatementConditionUsecase
	AssistantUsecase          *AssistantUsecase
	AiAssistantJobUsecase     *AiAssistantJobUsecase
	MapUsecase                *MapUsecase
}

func NewAiTaskJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AiStatementUsecase *AiStatementUsecase,
	TUsecase *TUsecase,
	AiUsecase *AiUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	DocEmailUsecase *DocEmailUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	StatementUsecase *StatementUsecase,
	DataComboUsecase *DataComboUsecase,
	AiTaskUsecase *AiTaskUsecase,
	AiResultUsecase *AiResultUsecase,
	AiTaskbuzUsecase *AiTaskbuzUsecase,
	FeeUsecase *FeeUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	AssistantUsecase *AssistantUsecase,
	AiAssistantJobUsecase *AiAssistantJobUsecase,
	MapUsecase *MapUsecase) *AiTaskJobUsecase {
	uc := &AiTaskJobUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		AiStatementUsecase:        AiStatementUsecase,
		TUsecase:                  TUsecase,
		AiUsecase:                 AiUsecase,
		JotformSubmissionUsecase:  JotformSubmissionUsecase,
		DocEmailUsecase:           DocEmailUsecase,
		Awsclaude3Usecase:         Awsclaude3Usecase,
		StatementUsecase:          StatementUsecase,
		DataComboUsecase:          DataComboUsecase,
		AiTaskUsecase:             AiTaskUsecase,
		AiResultUsecase:           AiResultUsecase,
		AiTaskbuzUsecase:          AiTaskbuzUsecase,
		FeeUsecase:                FeeUsecase,
		StatementConditionUsecase: StatementConditionUsecase,
		AssistantUsecase:          AssistantUsecase,
		AiAssistantJobUsecase:     AiAssistantJobUsecase,
		MapUsecase:                MapUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *AiTaskJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	sql := fmt.Sprintf(`select * from ai_tasks where handle_status=%d and deleted_at=0 and next_retry_at<=%d`,
		HandleStatus_waiting,
		time.Now().Unix())
	return c.CommonUsecase.DB().Raw(sql).Rows()

}

func (c *AiTaskJobUsecase) Handle(ctx context.Context, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		errString := err.Error()
		task.AppendHandleResultDetail(errString)
		if strings.Index(errString, "Bedrock Runtime: InvokeModel, exceeded maximum number of attempts") > 0 {
			task.RetryCount += 1
			task.NextRetryAt = time.Now().Unix() + 15
			task.HandleStatus = HandleStatus_waiting
		} else {
			c.log.Error(err, " id:", task.ID)
			task.HandleResult = HandleResult_failure
		}
	} else {
		task.HandleResult = HandleResult_ok
	}
	//c.log.Info("AiTaskJobUsecase 1:", InterfaceToString(task))
	err = c.CommonUsecase.DB().Omit("deleted_at").Save(task).Error
	if err != nil {
		return err
	}
	c.AfterAiTaskJobHandle(ctx, task)
	return nil
}

func (c *AiTaskJobUsecase) AfterAiTaskJobHandle(ctx context.Context, task *AiTaskEntity) {

	if task == nil {
		return
	}
	if task.FromType == AiTaskFromType_Assistant {
		if task.HandleStatus == HandleStatus_done {
			entity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"ai_task_id": task.ID})
			if err != nil {
				c.log.Error(err, task.ID)
			} else if entity != nil {
				if task.HandleResult == HandleResult_ok {
					entity.JobStatus = AiAssistantJob_JobStatus_Done
				} else {
					entity.JobStatus = AiAssistantJob_JobStatus_Failure
				}
				entity.UpdatedAt = time.Now().Unix()
				err = c.CommonUsecase.DB().Save(&entity).Error
				if err != nil {
					c.log.Error(err, task.ID)
				}
				aiTaskInputAssistant := entity.ToAiAssistantJobInput()

				if aiTaskInputAssistant.AssistantBiz == AiAssistantJobBizType_statementCondition &&
					aiTaskInputAssistant.InternalBizType == AssistantInternalBizType_ForAllStatemts {

					tCase, _ := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
					if tCase == nil {
						c.log.Error("tCase is nil")
					} else {
						tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
						if tClient == nil {
							c.log.Error("tClient is nil")
						} else {
							err = c.AssistantUsecase.HandleSaveAllStatements(*tCase, *tClient)
							c.log.Info("invokeSvStatemetns HandleSaveAllStatements")
							if err != nil {
								c.log.Error(err)
							}

						}
					}
				} else if aiTaskInputAssistant.AssistantBiz == AiAssistantJobBizType_genDocEmail &&
					aiTaskInputAssistant.InternalBizType == AssistantInternalBizType_AutoApply {

					jobUuid := GenDocEmailJobUuid(task.CaseId)
					aiAssistantJobEntity, er := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
					if er != nil {
						c.log.Error(er)
					}
					if aiAssistantJobEntity != nil {

						_, aiResult, err := c.AssistantUsecase.GetAiResultEntityFromAiTaskId(aiAssistantJobEntity.AiTaskId)
						if err != nil {
							c.log.Error(err)
						} else {
							if aiResult != nil {
								er = c.DocEmailUsecase.SetLatestDocEmailResult(task.CaseId, aiResult.ID)
								if er != nil {
									c.log.Error(er)
								}
							}
							er = c.AiAssistantJobUsecase.UpdatesByCond(map[string]interface{}{
								"job_status": AiAssistantJob_JobStatus_Normal,
							}, Eq{"id": aiAssistantJobEntity.ID})
							if er != nil {
								c.log.Error(er)
							}
						}
					}
				}
			}
		}
	}
}

func (c *AiTaskJobUsecase) HandleExec(ctx context.Context, task *AiTaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	if task.FromType == AiTaskFromType_generate_doc_email {
		err := c.HandleExecGenerateDocEmail(ctx, task)
		if err == nil {
			//go func() {
			//	tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
			//	if err != nil {
			//		c.log.Error(err)
			//	} else {
			//		err = c.DocEmailUsecase.HandleDocEmailToBox(tCase)
			//		if err != nil {
			//			c.log.Error(err)
			//		}
			//	}
			//
			//}()
		}
		return err
	} else if task.FromType == AiTaskFromType_statement {

		veteranSummary, err := c.HandleExecGenerateStatement(ctx, task)

		if err == nil {
			go func() {
				tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
				if err != nil {
					c.log.Error(err)
				} else {
					tClient, _, _ := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
					err = c.StatementUsecase.HandleStatementToBox(tCase, tClient, veteranSummary)
					if err != nil {
						c.log.Error(err)
					}
				}
			}()
		}
		return err

	} else if task.FromType == AiTaskFromType_return_timezone {
		return c.AiTaskbuzUsecase.HandleReturnTimezoneJob(ctx, task)
	} else if task.FromType == AiTaskFromType_update_statement {
		err := c.AiTaskbuzUsecase.HandleUpdateStatement(ctx, task)
		if err == nil {

			go func() {
				tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
				if err != nil {
					c.log.Error(err)
				} else {
					if tCase == nil {
						c.log.Error("tCase is nil")
					} else {
						tClient, _, _ := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
						if tClient == nil {
							c.log.Error("tClient is nil")
						} else {
							err = c.AiTaskbuzUsecase.HandleUpdateStatementDocToBox(*tClient, *tCase)
							if err != nil {
								c.log.Error(err)
							}
						}
					}
				}
			}()
		}
		return err
	} else if task.FromType == AiTaskFromType_Assistant {
		return c.AssistantUsecase.HandleAssistant(ctx, task)
	} else if task.FromType == AiTaskFromType_veteran_summary {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
		if err != nil {
			return err
		}
		aiResultId, err := c.AiTaskbuzUsecase.DoHandleVeteranSummary(ctx, tCase)
		if err != nil {
			return err
		}
		task.CurrentResultId = aiResultId
		return nil

	} else {
		return errors.New("task.FromType is wrong")
	}
	return nil
}

func (c *AiTaskbuzUsecase) DoHandleVeteranSummary(ctx context.Context, tCase *TData) (aiResultId int32, err error) {
	//promptKey := Ai_Prompt4_0
	promptKey := Ai_prompt4_0_muscul_summary_outjosn
	_, _, aiResultId, err = c.Awsclaude3Usecase.GenVeteranSummary(ctx, tCase, promptKey)
	if err != nil {
		return 0, err
	}
	return aiResultId, nil
}

type VeteranSummaryVo struct {
	FullName            string `json:"FullName"`
	UniqueID            string `json:"UniqueID"`
	BranchOfService     string `json:"BranchOfService"`
	YearsOfService      string `json:"YearsOfService"`
	RetirementStatus    string `json:"RetirementStatus"`
	Deployments         string `json:"Deployments"`
	MaritalStatus       string `json:"MaritalStatus"`
	Children            string `json:"Children"`
	OccupationInService string `json:"OccupationInService"`
}

func (c *VeteranSummaryVo) ToString() string {
	r := "Full Name: " + c.FullName + "\n"
	r += "Branch Of Service: " + c.BranchOfService + "\n"
	r += "Years Of Service: " + c.YearsOfService + "\n"
	r += "Retirement Status: " + c.RetirementStatus + "\n"
	r += "Deployments: " + c.Deployments + "\n"
	r += "Marital Status: " + c.MaritalStatus + "\n"
	r += "Children: " + c.Children + "\n"
	r += "Occupation In Service: " + c.OccupationInService
	return r
}

func (c *AiTaskbuzUsecase) HandleVeteranSummary(ctx context.Context, tCase *TData) (veteranSummaryVo VeteranSummaryVo, err error) {

	if tCase == nil {
		return veteranSummaryVo, errors.New("tCase is nil")
	}
	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return veteranSummaryVo, errors.New("tClient is nil")
	}

	aiTaskEntity, err := c.AiTaskUsecase.GetVeteranSummary(tCase)
	if err != nil {
		return veteranSummaryVo, err
	}
	if aiTaskEntity == nil {

		aiResultId, err := c.DoHandleVeteranSummary(ctx, tCase)
		if err != nil {
			return veteranSummaryVo, err
		}

		aiTaskEntity, err = c.AiTaskUsecase.CreateGenerateVeteranSummary(tCase)
		if err != nil {
			return veteranSummaryVo, err
		}
		aiTaskEntity.CurrentResultId = aiResultId
		err = c.CommonUsecase.DB().Save(&aiTaskEntity).Error
		if err != nil {
			return veteranSummaryVo, err
		}
	}

	aiResult, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTaskEntity.CurrentResultId})
	if err != nil {
		return veteranSummaryVo, err
	}
	if aiResult == nil {
		return veteranSummaryVo, nil
	}

	str := GetJsonFromAiResultForAssistant(aiResult.ParseResult)
	veteranSummaryVo = VeteranSummaryJsonToVo(str)

	// webform的优化级更高, 获取BaseInfo
	statementDetail, err := c.StatementUsecase.GetListStatementDetail(false, *tClient, *tCase, 0)
	if err != nil {
		return veteranSummaryVo, err
	}
	if statementDetail.BaseInfo.BranchOfService != "" {
		veteranSummaryVo.BranchOfService = statementDetail.BaseInfo.BranchOfService
	}
	if statementDetail.BaseInfo.YearsOfService != "" {
		veteranSummaryVo.YearsOfService = statementDetail.BaseInfo.YearsOfService
	}
	if statementDetail.BaseInfo.RetiredFromService != "" {
		veteranSummaryVo.RetirementStatus = statementDetail.BaseInfo.RetiredFromService
	}
	if statementDetail.BaseInfo.Deployments != "" {
		veteranSummaryVo.Deployments = statementDetail.BaseInfo.Deployments
	}
	if statementDetail.BaseInfo.MaritalStatus != "" {
		veteranSummaryVo.MaritalStatus = statementDetail.BaseInfo.MaritalStatus
	}
	if statementDetail.BaseInfo.Children != "" {
		veteranSummaryVo.Children = statementDetail.BaseInfo.Children
	}
	if statementDetail.BaseInfo.OccupationInService != "" {
		veteranSummaryVo.OccupationInService = statementDetail.BaseInfo.OccupationInService
	}

	return veteranSummaryVo, nil
}

func VeteranSummaryJsonToVo(json string) (veteranSummaryVo VeteranSummaryVo) {
	veteranSummaryVo = lib.StringToTDef(json, veteranSummaryVo)
	if IsEmptyResultForStatement(veteranSummaryVo.FullName) {
		veteranSummaryVo.FullName = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.UniqueID) {
		veteranSummaryVo.UniqueID = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.BranchOfService) {
		veteranSummaryVo.BranchOfService = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.YearsOfService) {
		veteranSummaryVo.YearsOfService = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.RetirementStatus) {
		veteranSummaryVo.RetirementStatus = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.Deployments) {
		veteranSummaryVo.Deployments = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.MaritalStatus) {
		veteranSummaryVo.MaritalStatus = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.Children) {
		veteranSummaryVo.Children = ""
	}
	if IsEmptyResultForStatement(veteranSummaryVo.OccupationInService) {
		veteranSummaryVo.OccupationInService = ""
	}
	return veteranSummaryVo
}

func (c *AiTaskJobUsecase) HandleExecGenerateStatement(ctx context.Context, task *AiTaskEntity) (veteranSummary VeteranSummaryVo, err error) {

	if task == nil {
		return veteranSummary, errors.New("task is nil")
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
	if err != nil {
		return veteranSummary, err
	}
	if tCase == nil {
		return veteranSummary, errors.New("tCase is nil")
	}
	inputVo := lib.StringToTDef[*AiTaskInputGenerateStatement](task.Input, nil)
	if inputVo == nil {
		return veteranSummary, errors.New("inputVo is nil")
	}
	if inputVo.StatementConditionId <= 0 {
		return veteranSummary, errors.New("StatementCondition is empty")
	}

	statementConditionEntity, err := c.StatementConditionUsecase.GetByCond(Eq{"id": inputVo.StatementConditionId})
	if err != nil {
		return veteranSummary, err
	}
	if statementConditionEntity == nil {
		return veteranSummary, errors.New("statementConditionEntity is nil")
	}

	veteranSummary, err = c.AssistantUsecase.HandleGenStatementFromAiTask(ctx, tCase, task, *statementConditionEntity, "")
	if err != nil {
		return veteranSummary, err
	}

	/*
		veteranSummary, err = c.AiTaskbuzUsecase.HandleVeteranSummary(ctx, tCase)
		if err != nil {
			return veteranSummary, errors.New("HandleVeteranSummary: " + err.Error())
		}

		//return nil
		//return errors.New("HandleVeteranSummary test")

		promptKey := CurrentGenStatementPrompt
		_, _, aiResultId, err := c.AiTaskbuzUsecase.GenStatement(ctx, tCase, statementConditionEntity.ToStatementCondition(), veteranSummary.ToString(), promptKey)
		if err != nil {
			return veteranSummary, err
		}

		//c.log.Info("GenStatement parseResult:", parseResult)
		//c.log.Info("GenStatement aiResultId:", aiResultId)
		aiTaskResultEntity := AiTaskResultEntity{
			AiTaskId:   task.ID,
			AiResultId: aiResultId,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		}
		err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
		if err != nil {
			return veteranSummary, err
		}
		task.CurrentResultId = aiResultId*/

	return veteranSummary, nil
}

func (c *AssistantUsecase) HandleGenStatementFromAiTask(ctx context.Context, tCase *TData, task *AiTaskEntity, statementConditionEntity StatementConditionEntity, userInputPrompt string) (veteranSummary VeteranSummaryVo, err error) {

	veteranSummary, err = c.AiTaskbuzUsecase.HandleVeteranSummary(ctx, tCase)
	if err != nil {
		return veteranSummary, errors.New("HandleVeteranSummary: " + err.Error())
	}
	//return veteranSummary, errors.New("HandleVeteranSummary: test error")

	promptKey := CurrentGenStatementPrompt
	_, _, aiResultId, err := c.AiTaskbuzUsecase.GenStatement(ctx, tCase,
		statementConditionEntity.ToStatementCondition(),
		veteranSummary.ToString(), promptKey, userInputPrompt)
	if err != nil {
		return veteranSummary, err
	}

	//c.log.Info("GenStatement parseResult:", parseResult)
	//c.log.Info("GenStatement aiResultId:", aiResultId)
	aiTaskResultEntity := AiTaskResultEntity{
		AiTaskId:   task.ID,
		AiResultId: aiResultId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
	if err != nil {
		return veteranSummary, err
	}
	task.CurrentResultId = aiResultId

	return veteranSummary, nil
}

func (c *AssistantUsecase) HandleGenStatementForStandardHeaderRevisionFromAiTask(ctx context.Context, tClient TData, tCase TData, task *AiTaskEntity, statementConditionEntity StatementConditionEntity, userInputPrompt string) (veteranSummary VeteranSummaryVo, err error) {

	veteranSummary, err = c.AiTaskbuzUsecase.HandleVeteranSummary(ctx, &tCase)
	if err != nil {
		return veteranSummary, errors.New("HandleVeteranSummary: " + err.Error())
	}
	//return veteranSummary, errors.New("HandleVeteranSummary: test error")

	_, _, aiResultId, err := c.AiTaskbuzUsecase.GenStatementForStandardHeaderRevision(ctx, tClient, tCase,
		statementConditionEntity,
		veteranSummary, userInputPrompt)
	if err != nil {
		return veteranSummary, err
	}

	//c.log.Info("GenStatement parseResult:", parseResult)
	//c.log.Info("GenStatement aiResultId:", aiResultId)
	aiTaskResultEntity := AiTaskResultEntity{
		AiTaskId:   task.ID,
		AiResultId: aiResultId,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
	if err != nil {
		return veteranSummary, err
	}
	task.CurrentResultId = aiResultId

	return veteranSummary, nil
}

func (c *AiTaskJobUsecase) HandleExecGenerateDocEmail(ctx context.Context, task *AiTaskEntity) error {

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

	inputVo := lib.StringToTDef[*AiTaskInputGenerateDocEmail](task.Input, nil)
	if inputVo == nil {
		return errors.New("inputVo is nil")
	}
	if inputVo.CaseStatements == "" {
		return errors.New("Statements is empty")
	}
	//var uniqcodes []string
	//uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	//uniqcodes = append(uniqcodes, uniqcode)
	//
	//isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	//if err != nil {
	//	return err
	//}
	//if !isPrimaryCase {
	//	uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	//}
	//
	//intakeForm, err := c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	//if err != nil {
	//	return err
	//}
	//if intakeForm == nil {
	//	return errors.New("intakeForm is nil")
	//}
	//
	//promptKey := Ai_Prompt3_0
	//r, aiResultId, err := c.AiUsecase.ExecuteWithClaude3New(promptKey, lib.TypeMap{
	//	"text": "", // 必须有key,但值为空
	//}, []Claude3Content{
	//	{
	//		Type: "text",
	//		Text: FormatJotformAnswers(lib.ToTypeMapByString(intakeForm.Notes)),
	//	},
	//	{
	//		Type: "text",
	//		Text: fmt.Sprintf("pls generate doc email:\n%s", inputVo.CaseStatements),
	//	},
	//})
	//
	//if err != nil {
	//	return err
	//}
	//c.log.Info("ExecuteWithClaude3 r:", r)
	//c.log.Info("ExecuteWithClaude3 aiResultId:", aiResultId)

	aiResultId, err := c.AssistantUsecase.ExecGenerateDocEmail(ctx, tCase, "")
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

	er := c.DocEmailUsecase.SetLatestDocEmailResult(tCase.Id(), aiResultId)
	if er != nil {
		c.log.Error(er)
	}
	return nil
}

func (c *AssistantUsecase) ExecGenerateDocEmail(ctx context.Context, tCase *TData, userInputPrompt string) (aiResultId int32, err error) {

	if tCase == nil {
		return 0, errors.New("tCase is nil")
	}

	statementConditions, err := c.StatementConditionUsecase.AllConditions(tCase.Id())
	if err != nil {
		return 0, err
	}

	docEmailStatements, _ := StatementConditionsToTextNoDivide(statementConditions)
	if docEmailStatements == "" {
		return 0, errors.New("docEmailStatements is empty")
	}

	//if task == nil {
	//	return errors.New("task is nil")
	//}
	//tCase, err := c.TUsecase.DataById(Kind_client_cases, task.CaseId)
	//if err != nil {
	//	return err
	//}
	//if tCase == nil {
	//	return errors.New("tCase is nil")
	//}
	//
	//inputVo := lib.StringToTDef[*AiTaskInputGenerateDocEmail](task.Input, nil)
	//if inputVo == nil {
	//	return errors.New("inputVo is nil")
	//}
	//if inputVo.CaseStatements == "" {
	//	return errors.New("Statements is empty")
	//}
	var uniqcodes []string
	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	uniqcodes = append(uniqcodes, uniqcode)

	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return 0, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	intakeForm, err := c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return 0, err
	}
	if intakeForm == nil {
		return 0, errors.New("intakeForm is nil")
	}

	promptKey := CurrentPromptForDocEmail

	inputString := fmt.Sprintf("pls generate doc email:\n%s", docEmailStatements)

	//inputString = "Please generate exactly **" + InterfaceToString(conditionCount) + "** separate DocEmail documents based on the **" + InterfaceToString(conditionCount) + " medical conditions** listed below — one DocEmail per condition.\n\n"
	//
	//inputString += "Each condition is already clearly defined and should not be split further.\n\n"
	//
	//inputString += "Conditions:\n"
	//inputString += docEmailStatements
	//inputString += "\n\n**Pls generate Doctor's Email per condition.**"

	//inputString := "**Please generate one DocEmail per condition** from the list below:\n\n"
	//inputString += docEmailStatements

	//if userInputPrompt != "" {
	//	inputString += "\n\n" + userInputPrompt
	//}

	r, aiResultId, err := c.AiUsecase.ExecuteWithClaude3New(promptKey, lib.TypeMap{
		"text": "", // 必须有key,但值为空
	}, []Claude3Content{
		{
			Type: "text",
			Text: FormatJotformAnswers(lib.ToTypeMapByString(intakeForm.Notes)),
		},
		{
			Type: "text",
			Text: inputString,
		},
	})

	if err != nil {
		return 0, err
	}
	c.log.Info("ExecuteWithClaude3 r:", r)
	c.log.Info("ExecuteWithClaude3 aiResultId:", aiResultId)

	return aiResultId, nil
	//
	//aiTaskResultEntity := AiTaskResultEntity{
	//	AiTaskId:   task.ID,
	//	AiResultId: aiResultId,
	//	CreatedAt:  time.Now().Unix(),
	//	UpdatedAt:  time.Now().Unix(),
	//}
	//err = c.CommonUsecase.DB().Save(&aiTaskResultEntity).Error
	//if err != nil {
	//	return err
	//}
	//task.CurrentResultId = aiResultId

	//c.log.Info("AiTaskJobUsecase:", InterfaceToString(task))

	//return nil
}
