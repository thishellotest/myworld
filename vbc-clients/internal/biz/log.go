package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

const Log_FromType_Envelope_completed = "Envelope_completed"
const Log_FromType_AdobeSign_completed_failed = "AdobeSign_completed_failed"
const Log_FromType_Box_CreateContract = "Box_CreateContract"
const Log_FromType_Box_CreateContractAm = "Box_CreateContractAm"

//const Log_FromType_Box_CreatePatientPaymentFormContract = "Box_CreatePatientPaymentFormContract"
// const Log_FromType_Box_ReleaseOfInformationSignRequests = "Box_ReleaseOfInformationSignRequests"

const Log_FromType_Box_MedicalTeamFormsSignRequests = "Box_MedicalTeamFormsSignRequests"

const Log_FromType_Xero_CreateInvoice = "Xero_CreateInvoice"
const Log_FromType_Xero_AmCreateInvoice = "Xero_AmCreateInvoice"
const Log_FromType_Asana_SyncTaskInfo = "Asana_SyncTaskInfo"
const Log_FromType_HandleClientCaseName = "HandleClientCaseName"
const Log_FromType_ZohoinfoSync = "ZohoinfoSync"
const Log_FromType_Xero_InvoicesApi = "Xero_InvoicesApi"
const Log_FromType_HandleClientTask = "HandleClientTask"
const Log_FromType_MultiCasesBaseInfoSync = "MultiCasesBaseInfoSync"
const Log_FromType_RecordReviewLPushRedisQueue = "RecordReviewLPushRedisQueue"
const Log_FromType_RecordReviewBizHandleTask = "RecordReviewBizHandleTask"
const Log_FromType_RecordReviewSyncExistsSameFileName = "RecordReviewSyncExistsSameFileName"
const Log_FormType_TransferGeneral = "TransferGeneral"
const Log_FormType_TransferPsych = "TransferPsych"
const Log_FormType_ClientTasks = "ClientTasks"

type LogEntity struct {
	ID        int32 `gorm:"primaryKey"`
	FromId    int32
	FromType  string
	Notes     string
	CreatedAt int64
}

func (LogEntity) TableName() string {
	return "log"
}

func GenLog(fromId int32, fromType string, notes string) *LogEntity {
	return &LogEntity{
		FromId:    fromId,
		FromType:  fromType,
		Notes:     notes,
		CreatedAt: time.Now().Unix(),
	}
}

type LogUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[LogEntity]
}

func NewLogUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *LogUsecase {
	uc := &LogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *LogUsecase) SaveLog(fromId int32, fromType string, notes interface{}) error {
	log := GenLog(fromId, fromType, InterfaceToString(notes))
	return c.CommonUsecase.DB().Save(log).Error
}
