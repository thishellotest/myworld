package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

const (
	TaskType_Task                               = ""
	TaskType_ChangeStagesToGettingStartedEmail  = "ChangeStagesToGettingStartedEmail"
	TaskType_ChangeStagesToMiniDBQsFinalized    = "ChangeStagesToMiniDBQsFinalized"
	TaskType_CompleteBoxSignBehavior            = "CompleteBoxSignBehavior"
	TaskType_XeroCreateInvoice                  = "XeroCreateInvoice"
	TaskType_XeroAmCreateInvoice                = "XeroAmCreateInvoice"
	TaskType_MaCongratsEmail_HandleInputTask    = "MaCongratsEmail_HandleInputTask"
	TaskType_HandleFeeScheduleCommunicationMail = "HandleFeeScheduleCommunicationMail"
	TaskType_HandleCreateFolderInBoxAndMail     = "HandleCreateFolderInBoxAndMail"
	TaskType_HandleEnvelope                     = "HandleEnvelope"
	TaskType_HandleAmount                       = "HandleAmount"
	TaskType_HandleExecWebsite                  = "HandleExecWebsite"
	TaskType_ZohoinfoSync                       = "ZohoinfoSync"
	TaskType_InitClientCase                     = "InitClientCase"
	TaskType_BoxWebhookLog                      = "BoxWebhookLog"
	//TaskType_HandlePatientPaymentForm           = "HandlePatientPaymentForm"
	//TaskType_HandleReleaseOfInformation         = "HandleReleaseOfInformation"
	TestType_HandleMedicalTeamForms             = "HandleMedicalTeamForms"
	TaskType_HandleMiscThingsToKnowCPExam       = "HandleMiscThingsToKnowCPExam"
	TaskType_HandleRemoveMiscThingsToKnowCPExam = "HandleRemoveMiscThingsToKnowCPExam"
	TaskType_HandleAutomationCompleteTask       = "HandleAutomationCompleteTask"
	TaskType_HandlePrivateExamsSubmitted        = "HandlePrivateExamsSubmitted"
)

type TaskLogEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	TaskType           string
	TaskId             int32
	Notes              string
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	CreatedAt          int64
}

func (TaskLogEntity) TableName() string {
	return "task_failure_log"
}

type TaskFailureLogUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TaskLogEntity]
}

func NewTaskFailureLogUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *TaskFailureLogUsecase {
	uc := &TaskFailureLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *TaskFailureLogUsecase) Add(taskType string, taskId int32, notes interface{}) error {
	c.log.Error("taskType: ", taskType, " taskId: ", taskId, " notes: ", InterfaceToString(notes))
	return c.CommonUsecase.DB().Create(&TaskLogEntity{
		TaskType:  taskType,
		TaskId:    taskId,
		Notes:     InterfaceToString(notes),
		CreatedAt: time.Now().Unix(),
	}).Error
}
