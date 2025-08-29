package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

/*


CREATE TABLE `ai_tasks` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Id',
  `handle_status` tinyint(4) NOT NULL DEFAULT '0',
  `handle_result` int(11) NOT NULL DEFAULT '0',
  `handle_result_detail` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_type` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务来源',
  `from_code` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '来源唯一标识',
  `input` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'json格式',
  `current_result_id` int(11) NOT NULL DEFAULT '0',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Created At',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Updatd At',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Deleted At',
  PRIMARY KEY (`id`),
  KEY `idxHS` (`handle_status`),
  KEY `uniqcode` (`from_code`(191)),
  KEY `uniq` (`from_code`(191),`from_type`(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ai_tasks';

*/

const (
	AiTaskFromType_generate_doc_email = "generate_doc_email"
	AiTaskFromType_veteran_summary    = "veteran_summary"
	AiTaskFromType_statement          = "statement"
	AiTaskFromType_update_statement   = "update_statement"

	AiTaskFromType_return_timezone = "Return Time Zone"

	AiTaskFromType_Assistant = "assistant"
)

const (
	AiAssistantJobBizType_statementSection   = "statementSection"
	AiAssistantJobBizType_statementCondition = "statementCondition"
	AiAssistantJobBizType_genDocEmail        = "genDocEmail"
	AiAssistantJobBizType_allStatements      = "allStatements"
)

const (
// AiTaskInputAssistant_Biz_StatementCondition = "StatementCondition"
)

const (
	AiAssistantBizType_SetAllStatement  = "SetAllStatement"
	AiAssistantBizType_UpdateMedication = "UpdateMedication"

	AiAssistantBizType_CreateNewStatement     = "CreateNewStatement"
	AiAssistantBizType_StandardHeaderRevision = "StandardHeaderRevision"

	AiAssistantBizType_DocEmailRenew = "RenewDocEmail"

	AssistantInternalBizType_ForAllStatemts = "ForAllStatemts"
	AssistantInternalBizType_AutoApply      = "AutoApply"
)

var AiAssistantStatementSectionQuickOptions = []AiAssistantJobOption{
	//{
	//	Value: AiAssistantBizType_UpdateMedication,
	//	Label: "Update Medication",
	//},
	{
		Value: AiAssistantBizType_CreateNewStatement,
		Label: "Create New Statement",
	},
	{
		Value: AiAssistantBizType_StandardHeaderRevision,
		Label: "Standard Header Revision",
	},
}

var AiAssistantAllStatementsQuickOptions = []AiAssistantJobOption{
	{
		Value: AiAssistantBizType_SetAllStatement,
		Label: "Create All New Statements",
	},
}

var AiAssistantGenDocEmailQuickOptions = []AiAssistantJobOption{
	{
		Value: AiAssistantBizType_DocEmailRenew,
		Label: "Create New DocEmail",
	},
}

func (c *AiAssistantJobDetail) HandleQuickOptions() {
	arr := strings.Split(c.JobUuid, ":")
	if len(arr) > 0 {
		if arr[0] == AiAssistantJobBizType_statementCondition {
			//if len(arr) == 4 && arr[3] == Statemt_Section_Medication {
			c.QuickOptions = AiAssistantStatementSectionQuickOptions
			//}
		} else if arr[0] == AiAssistantJobBizType_allStatements {
			c.QuickOptions = AiAssistantAllStatementsQuickOptions
		} else if arr[0] == AiAssistantJobBizType_genDocEmail {
			c.QuickOptions = AiAssistantGenDocEmailQuickOptions
		}
	}
}

type AiTaskInputAssistant struct {
	AssistantBiz         string `json:"assistant_biz"`
	BizType              string `json:"biz_type"`
	InternalBizType      string `json:"internal_biz_type"`
	UserInputText        string `json:"user_input_text"`
	CaseId               int32  `json:"case_id"`
	StatementConditionId int32  `json:"statement_condition_id"`
	SectionType          string `json:"section_type"`
}

type AiTaskInputGenerateStatement struct {
	StatementConditionId int32 `json:"statement_condition_id"`
	//StatementCondition StatementCondition `json:"statement_condition"`
}

var AiTaskFromTypeConfigs = lib.TypeMap{
	AiTaskFromType_generate_doc_email: "Base Generated Doc Email",
	AiTaskFromType_veteran_summary:    "Veteran Summary",
	AiTaskFromType_statement:          "Base Generated Statement",
	AiTaskFromType_update_statement:   "Update Statement",
}

const (
	AiTask_HandleStatus_In_process = 0
	AiTask_HandleStatus_Complelte  = 1
	AiTask_HandleStatus_Waiting    = 2 // 当AiTaskFromType_statement时，需要补充Associate Jotform
)

type AiTaskEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	FromType           string
	FromCode           string
	TaskUniqcode       string
	CaseId             int32
	ClientId           int32
	Input              string
	CurrentResultId    int32
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
	NextRetryAt        int64
	RetryCount         int
	SerialNumber       int
	ToPsform           int
}

const (
	AiTask_ToPsform_Yes = 1
	AiTask_ToPsform_No  = 0
)

func (c *AiTaskEntity) GetAiTaskInputAssistant() (r AiTaskInputAssistant) {
	r = lib.StringToTDef(c.Input, r)
	return
}

func (c *AiTaskEntity) GetAiTaskInputUpdateStatementTask() (r AiTaskInputUpdateStatementTask) {
	r = lib.StringToTDef(c.Input, r)
	return
}

func (c *AiTaskEntity) GetIput() lib.TypeMap {
	return lib.ToTypeMapByString(c.Input)
}

func (c *AiTaskEntity) GetAiTaskInputGenerateStatement() (r AiTaskInputGenerateStatement) {
	r = lib.StringToTDef(c.Input, r)
	return
}

type AiTaskList []AiTaskItem
type AiTaskItem struct {
	TaskId             int32  `json:"task_id"`
	HandleStatus       int    `json:"handle_status"`
	HandleResult       int    `json:"handle_result"`
	Detail             string `json:"detail"`
	HandleResultDetail string `json:"handle_result_detail"`
	FromType           string `json:"from_type"`
	FromTypeLabel      string `json:"from_type_label"`
	Input              string `json:"input"`
	CreatedAt          int32  `json:"created_at"`
	UpdatedAt          int32  `json:"updated_at"`
	NextRetryAt        int32  `json:"next_retry_at"`
	RetryCount         int    `json:"retry_count"`
	SerialNumber       int    `json:"serial_number"`
}

func (c *AiTaskEntity) ToAiTaskItem() AiTaskItem {

	item := AiTaskItem{
		TaskId:             c.ID,
		HandleStatus:       c.HandleStatus,
		HandleResult:       c.HandleResult,
		HandleResultDetail: c.HandleResultDetail,
		FromType:           c.FromType,
		FromTypeLabel:      InterfaceToString(AiTaskFromTypeConfigs[c.FromType]),
		Input:              c.Input,
		CreatedAt:          int32(c.CreatedAt),
		UpdatedAt:          int32(c.UpdatedAt),
		NextRetryAt:        int32(c.NextRetryAt),
		RetryCount:         c.RetryCount,
		SerialNumber:       c.SerialNumber,
	}
	if c.FromType == AiTaskFromType_update_statement {
		personalStatementOneVo := c.GetAiTaskInputUpdateStatementTask().PersonalStatementOneVo
		vo, _ := personalStatementOneVo.ToParseUpdateStatementVo()
		item.Detail = fmt.Sprintf("%d: %s", c.SerialNumber, vo.NameOfDisability.Value)
	}
	return item
}

func (AiTaskEntity) TableName() string {
	return "ai_tasks"
}

func (c *AiTaskEntity) AppendHandleResultDetail(str string) {
	//if c.HandleResultDetail == "" {
	c.HandleResultDetail = str
	//} else {
	//	c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	//}
}

type AiTaskUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[AiTaskEntity]
}

func NewAiTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *AiTaskUsecase {
	uc := &AiTaskUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type AiTaskInputStatement struct {
	SubmissionId []string
	Condition    string
}

type AiTaskInputGenerateDocEmail struct {
	CaseStatements string `json:"case_statements"`
}

type AiAssistantJobStatusVoMap map[string]AiAssistantJobStatusVo

// GetGenerateStatementForJobStatusWithJobUuids 只能是同一个case的jobUuids
func (c *AiTaskUsecase) GetGenerateStatementForJobStatusWithJobUuids(jobUuids []string) (records AiAssistantJobStatusVoMap, err error) {

	if len(jobUuids) == 0 {
		return nil, nil
	}
	jobUuidVos := make(map[string]JobUuidForStatementCondition)
	var caseId int32
	var statementConditionIds []int32
	for _, v := range jobUuids {
		if strings.Index(v, AiAssistantJobBizType_statementCondition) == 0 {
			vo := FormatJobUuidForStatementCondition(v)
			jobUuidVos[v] = vo
			caseId = vo.CaseId
			statementConditionIds = append(statementConditionIds, vo.StatementConditionId)
		}
	}
	aiTasks, _ := c.GetGenerateStatements(caseId, statementConditionIds)

	getAiTask := func(StatementConditionId int32) *AiTaskEntity {
		for k, v := range aiTasks {
			if v.TaskUniqcode == InterfaceToString(StatementConditionId) {
				return aiTasks[k]
			}
		}
		return nil
	}

	records = make(AiAssistantJobStatusVoMap)
	for _, v := range jobUuidVos {
		jobUuid := v.ToJobUuid()
		entity := getAiTask(v.StatementConditionId)
		records[jobUuid] = AiTaskToAiAssistantJobStatusVo(entity)
	}
	return records, nil
}

func AiTaskToAiAssistantJobStatusVo(aiTask *AiTaskEntity) (vo AiAssistantJobStatusVo) {
	vo.JobStatus = AiAssistantJob_JobStatus_Normal
	if aiTask != nil && aiTask.HandleStatus == HandleStatus_waiting {
		vo.JobStatus = AiAssistantJob_JobStatus_Running
	}
	return vo
}

func (c *AiTaskUsecase) GetGenerateStatementForJobStatus(caseId int32, StatementConditionId int32) (vo AiAssistantJobStatusVo, err error) {
	entity, err := c.GetGenerateStatement(caseId, StatementConditionId)
	if err != nil {
		return vo, nil
	}
	vo = AiTaskToAiAssistantJobStatusVo(entity)
	return vo, nil
}

func (c *AiTaskUsecase) GetGenerateStatements(caseId int32, statementConditionIds []int32) ([]*AiTaskEntity, error) {
	return c.AllByCond(And(Eq{"from_type": AiTaskFromType_statement, "deleted_at": 0, "case_id": caseId}, In("task_uniqcode", statementConditionIds)))
}

func (c *AiTaskUsecase) GetGenerateStatement(caseId int32, StatementConditionId int32) (*AiTaskEntity, error) {
	return c.GetByCond(Eq{"from_type": AiTaskFromType_statement, "deleted_at": 0, "case_id": caseId, "task_uniqcode": StatementConditionId})
}

func (c *AiTaskUsecase) CreateGenerateStatement(tCase *TData, statementCondition StatementCondition, handleStatus *int) (*AiTaskEntity, error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	input := &AiTaskInputGenerateStatement{
		StatementConditionId: statementCondition.StatementConditionId,
	}
	return c.CreateTask(AiTaskFromType_statement, "", tCase.Id(), 0, InterfaceToString(input), InterfaceToString(statementCondition.StatementConditionId), handleStatus, to.Ptr(statementCondition.Sort))
}

func (c *AiTaskUsecase) CreateGenerateDocEmail(tCase *TData) (*AiTaskEntity, error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	aiTaskInputGenerateDocEmail := &AiTaskInputGenerateDocEmail{
		CaseStatements: strings.TrimSpace(tCase.CustomFields.TextValueByNameBasic(FieldName_statements)),
	}
	return c.CreateTask(AiTaskFromType_generate_doc_email, "", tCase.Id(), 0, InterfaceToString(aiTaskInputGenerateDocEmail), "", nil, nil)
}

type AiTaskInputUpdateStatementTask struct {
	PersonalStatementOneVo PersonalStatementOneVo
}

func (c *AiTaskUsecase) CreateUpdateStatementTask(tCase TData, personalStatementOneVo PersonalStatementOneVo, serialNumber int) (*AiTaskEntity, error) {
	aiTaskInputGenerateDocEmail := &AiTaskInputUpdateStatementTask{
		PersonalStatementOneVo: personalStatementOneVo,
	}
	return c.CreateTask(AiTaskFromType_update_statement,
		"",
		tCase.Id(),
		0,
		InterfaceToString(aiTaskInputGenerateDocEmail),
		"",
		nil,
		&serialNumber)
}

//func (c *AiTaskUsecase) CreateUpdateStatementTask(tCase TData, personalStatementsVo PersonalStatementsVo) (*AiTaskEntity, error) {
//	aiTaskInputGenerateDocEmail := &AiTaskInputUpdateStatementTask{
//		PersonalStatementsVo: personalStatementsVo,
//	}
//	return c.CreateTask(AiTaskFromType_update_statement, "", tCase.Id(), 0, InterfaceToString(aiTaskInputGenerateDocEmail), "", nil)
//}

func (c *AiTaskUsecase) GetVeteranSummary(tCase *TData) (*AiTaskEntity, error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	return c.GetByCond(Eq{"from_type": AiTaskFromType_veteran_summary, "case_id": tCase.Id(), "deleted_at": 0})
}

func (c *AiTaskUsecase) CreateGenerateVeteranSummary(tCase *TData) (*AiTaskEntity, error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	aiTaskInputGenerateDocEmail := &AiTaskInputGenerateDocEmail{
		CaseStatements: strings.TrimSpace(tCase.CustomFields.TextValueByNameBasic(FieldName_statements)),
	}
	status := AiTask_HandleStatus_Complelte
	return c.CreateTask(AiTaskFromType_veteran_summary, "", tCase.Id(), 0, InterfaceToString(aiTaskInputGenerateDocEmail), "", &status, nil)
}

func (c *AiTaskUsecase) GetReturnTimezone(tClient *TData) (*AiTaskEntity, error) {
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	return c.GetByCond(Eq{"from_type": AiTaskFromType_return_timezone, "client_id": tClient.Id(), "deleted_at": 0})
}

func (c *AiTaskUsecase) CreateReturnTimezone(tClient *TData) (*AiTaskEntity, error) {
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	//input := make(lib.TypeMap)
	//input.Set("value", fmt.Sprintf("%s, %s", tCase.CustomFields.TextValueByNameBasic(FieldName_state), tCase.CustomFields.TextValueByNameBasic(FieldName_city)))
	return c.CreateTask(AiTaskFromType_return_timezone, "", 0, tClient.Id(), "", "", nil, nil)
}

func (c *AiTaskUsecase) CreateTask(fromType string, fromCode string, caseId int32, clientId int32, input interface{}, taskUniqcode string, handleStatus *int, serialNumber *int) (*AiTaskEntity, error) {

	aiTaskEntity := &AiTaskEntity{
		FromType:     fromType,
		FromCode:     fromCode,
		Input:        InterfaceToString(input),
		CaseId:       caseId,
		ClientId:     clientId,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		TaskUniqcode: taskUniqcode,
	}
	if handleStatus != nil {
		aiTaskEntity.HandleStatus = *handleStatus
	}
	if serialNumber != nil {
		aiTaskEntity.SerialNumber = *serialNumber
	}
	err := c.CommonUsecase.DB().Save(&aiTaskEntity).Error
	if err != nil {
		return nil, err
	}
	return aiTaskEntity, nil
}
