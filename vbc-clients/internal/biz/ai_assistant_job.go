package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
CREATE TABLE `ai_assistant_job` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`job_uuid` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`is_running` tinyint(4) NOT NULL DEFAULT '0',
	`result` text COLLATE utf8mb4_unicode_ci COMMENT 'json格式，如何为空时，说明恢复到初始状态',
	`ai_task_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联的内部任务',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	KEY `idx` (`job_uuid`(191))

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ai_assistant_job';
*/
type AiAssistantJobEntity struct {
	ID        int32 `gorm:"primaryKey"`
	JobUuid   string
	JobInput  string
	JobStatus int
	AiTaskId  int32
	CreatedAt int64
	UpdatedAt int64
}

func (c *AiAssistantJobEntity) ToAiAssistantJobInput() (vo AiAssistantJobInput) {
	vo = lib.StringToTDef(c.JobInput, vo)
	return vo
}

const (
	AiAssistantJob_JobStatus_Normal  = 0
	AiAssistantJob_JobStatus_Running = 1
	AiAssistantJob_JobStatus_Done    = 2
	AiAssistantJob_JobStatus_Failure = 3
)

type JobUuidForCommon struct {
	CaseId int32 `json:"case_id"`
}

func FormatJobUuidForCommon(jobUuid string) (vo JobUuidForCommon) {
	arr := strings.Split(jobUuid, ":")
	if len(arr) == 2 {
		a, _ := strconv.ParseInt(arr[1], 10, 32)
		vo.CaseId = int32(a)
	}
	return
}

type JobUuidForAllStatements struct {
	CaseId int32 `json:"case_id"`
}

func FormatJobUuidForAllStatements(jobUuid string) (vo JobUuidForAllStatements) {
	arr := strings.Split(jobUuid, ":")
	if len(arr) == 2 {
		a, _ := strconv.ParseInt(arr[1], 10, 32)
		vo.CaseId = int32(a)
	}
	return
}

type JobUuidForStatementCondition struct {
	CaseId               int32 `json:"case_id"`
	StatementConditionId int32 `json:"statement_condition_id"`
}

func (c *JobUuidForStatementCondition) ToJobUuid() string {
	return fmt.Sprintf("%s:%d:%d", AiAssistantJobBizType_statementCondition, c.CaseId, c.StatementConditionId)
}

func FormatJobUuidForStatementCondition(jobUuid string) (vo JobUuidForStatementCondition) {
	arr := strings.Split(jobUuid, ":")
	if len(arr) == 3 {
		a, _ := strconv.ParseInt(arr[1], 10, 32)
		vo.CaseId = int32(a)
		a, _ = strconv.ParseInt(arr[2], 10, 32)
		vo.StatementConditionId = int32(a)
	}
	return
}

type JobUuidForStatementSection struct {
	CaseId               int32  `json:"case_id"`
	StatementConditionId int32  `json:"statement_condition_id"`
	SectionType          string `json:"section_type"`
}

func FormatJobUuidForStatementSection(jobUuid string) (vo JobUuidForStatementSection) {
	arr := strings.Split(jobUuid, ":")
	if len(arr) == 4 {
		a, _ := strconv.ParseInt(arr[1], 10, 32)
		vo.CaseId = int32(a)
		a, _ = strconv.ParseInt(arr[2], 10, 32)
		vo.StatementConditionId = int32(a)
		vo.SectionType = arr[3]
	}
	return
}

type AiAssistantJobStatusVo struct {
	JobStatus int `json:"job_status"`
}

type AiAssistantJobResult struct {
	Text string `json:"text"`
}

type AiAssistantJobOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type AiAssistantJobDetail struct {
	JobUuid        string                 `json:"job_uuid"`
	JobStatusInfo  AiAssistantJobStatusVo `json:"job_status_info"`
	JobResult      AiAssistantJobResult   `json:"job_result"`
	SelectedOption string                 `json:"selected_option"`
	Prompt         string                 `json:"prompt"`
	QuickOptions   []AiAssistantJobOption `json:"quick_options"`
}

func (c *AiAssistantJobEntity) GetAiResultEntity(AiTaskUsecase *AiTaskUsecase, AiResultUsecase *AiResultUsecase) (*AiTaskEntity, *AiResultEntity) {
	if AiTaskUsecase != nil && AiResultUsecase != nil {
		task, _ := AiTaskUsecase.GetByCond(Eq{"id": c.AiTaskId})
		if task != nil {
			if task.HandleResult == HandleResult_ok {
				aiResultEntity, _ := AiResultUsecase.GetByCond(Eq{"id": task.CurrentResultId})
				return task, aiResultEntity
			} else {
				return task, nil
			}
		}
	}
	return nil, nil
}

func (c *AiAssistantJobEntity) ToJobDetail(AiTaskUsecase *AiTaskUsecase, AiResultUsecase *AiResultUsecase) (vo AiAssistantJobDetail) {

	vo.JobUuid = c.JobUuid
	vo.JobStatusInfo = c.ToJobStatusInfo()
	if c.JobStatus == AiAssistantJob_JobStatus_Done || c.JobStatus == AiAssistantJob_JobStatus_Failure {
		if AiTaskUsecase != nil && AiResultUsecase != nil {
			aiTaskEntity, aiResultEntity := c.GetAiResultEntity(AiTaskUsecase, AiResultUsecase)
			if aiResultEntity != nil {
				vo.JobResult.Text = aiResultEntity.ParseResult
			}
			if aiTaskEntity != nil {
				if aiTaskEntity.HandleResult == HandleResult_failure {
					vo.JobResult.Text = aiTaskEntity.HandleResultDetail
				}
			}
		}
	}
	input := c.ToAiAssistantJobInput()
	vo.Prompt = input.UserInputText
	vo.SelectedOption = input.BizType
	vo.HandleQuickOptions()

	return vo
}

func (c *AiAssistantJobEntity) ToJobStatusInfo() (aiAssistantJobStatusVo AiAssistantJobStatusVo) {
	aiAssistantJobStatusVo.JobStatus = c.JobStatus
	return aiAssistantJobStatusVo
}

type AiAssistantJobInput struct {
	AssistantBiz    string `json:"assistant_biz"`
	InternalBizType string `json:"internal_biz_type"`
	BizType         string `json:"biz_type"`
	UserInputText   string `json:"user_input_text"`
}

func (AiAssistantJobEntity) TableName() string {
	return "ai_assistant_job"
}

type AiAssistantJobUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[AiAssistantJobEntity]
}

func NewAiAssistantJobUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *AiAssistantJobUsecase {
	uc := &AiAssistantJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

// CanRunNewJob 是否可以运行新任务
func (c *AiAssistantJobUsecase) CanRunNewJob(jobUuid string) (can bool, entity *AiAssistantJobEntity, err error) {
	entity, err = c.GetByCond(Eq{"job_uuid": jobUuid})
	if err != nil {
		return
	}
	if entity == nil {
		return true, nil, nil
	}
	if entity.JobStatus != 0 {
		return false, entity, nil
	}
	return true, entity, nil
}
