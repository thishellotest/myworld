package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/utils"
)

type CronTriggerCreateUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	TaskCreateUsecase *TaskCreateUsecase
	TUsecase          *TUsecase
}

func NewCronTriggerCreateUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TaskCreateUsecase *TaskCreateUsecase,
	TUsecase *TUsecase) *CronTriggerCreateUsecase {
	uc := &CronTriggerCreateUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		TaskCreateUsecase: TaskCreateUsecase,
		TUsecase:          TUsecase,
	}

	return uc
}

// CancelDialpadSMSTasks 取消任务，阶段修改后，后续任务不需要了
func (c *CronTriggerCreateUsecase) CancelDialpadSMSTasks(caseId int32) error {

	return c.CommonUsecase.DB().Model(&TaskEntity{}).
		Where("from_id = ? and from_type =? and task_status=?",
			InterfaceToString(caseId),
			Task_FromType_DialpadSMS,
			Task_TaskStatus_processing).
		Updates(map[string]interface{}{
			"task_status": Task_TaskStatus_cancel,
			"updated_at":  time.Now().Unix()}).Error

}

func (c *CronTriggerCreateUsecase) CreateAfterSignedContract(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	// 做延时处理
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAfterSignedContract,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		"", "")
	// 此处不加 fromType和fromId 保证任务不会被取消
}

func (c *CronTriggerCreateUsecase) CreateGettingStartedEmailByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateGettingStartedEmail(tCase)
}

func (c *CronTriggerCreateUsecase) CreateGettingStartedEmail(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAt := utils.CalDelayDayTime(time.Now().Add(30*time.Second), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextGettingStartedEmail,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

}

// CreateScheduleCallByCaseId 改为：RecordReview触发了
func (c *CronTriggerCreateUsecase) CreateScheduleCallByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateScheduleCall(tCase)
}

func (c *CronTriggerCreateUsecase) CreateScheduleCall(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	// 做延时处理
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextScheduleCall,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

}

func (c *CronTriggerCreateUsecase) CreateGettingStartedEmailTaskLongerThan30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateGettingStartedEmailTaskLongerThan30Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateGettingStartedEmailTaskLongerThan30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextGettingStartedEmailTaskLongerThan30Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateAwaitingClientRecordsLongerThan30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateAwaitingClientRecordsLongerThan30Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateAwaitingClientRecordsLongerThan30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAwaitingClientRecordsLongerThan30Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateSTRRequestPending30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateSTRRequestPending30Days(tCase)
}
func (c *CronTriggerCreateUsecase) CreateSTRRequestPending30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextSTRRequestPendingLongerThan30Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

}

func (c *CronTriggerCreateUsecase) CreateStatementsFinalizedByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateStatementsFinalized(tCase)
}

func (c *CronTriggerCreateUsecase) CreateStatementsFinalized(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	// 做延时处理
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextStatementFinalized,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateStatementsFinalizedEvery14DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateStatementsFinalizedEvery14Days(tCase)
}

func GetStatementsFinalizedEvery14DaysTime(tCase *TData, log *log.Helper) (time.Time, error) {
	timeLocation := GetCaseTimeLocation(tCase, log)
	return utils.CalIntervalDayTime(time.Now(), 14, "08:00", timeLocation)
}

func (c *CronTriggerCreateUsecase) CreateStatementsFinalizedEvery14Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeAt, err := GetStatementsFinalizedEvery14DaysTime(tCase, c.log)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextStatementFinalizedEvery14Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateCurrentTreatmentByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateCurrentTreatment(tCase)
}

func (c *CronTriggerCreateUsecase) CreateCurrentTreatment(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	// 做延时处理
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextCurrentTreatment,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateCurrentTreatment30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateCurrentTreatment30Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateCurrentTreatment30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextCurrentTreatment30Day,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateAwaitingDecision30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateAwaitingDecision30Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateAwaitingDecision30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:05", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAwaitingDecision30Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPaymentByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateAwaitingPayment(tCase)
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPayment(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAwaitingPayment,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPaymentAfter14DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateAwaitingPaymentAfter14Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPaymentAfter14Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 14, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAwaitingPaymentAfter14Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPaymentTaskOpen30DaysByCaseId(caseId int32) error {
	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.CreateAwaitingPaymentTaskOpen30Days(tCase)
}

func (c *CronTriggerCreateUsecase) CreateAwaitingPaymentTaskOpen30Days(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")

	timeLocation := GetCaseTimeLocation(tCase, c.log)
	timeAt, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}

	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextAwaitingPaymentTaskOpen30Days,
	}, Task_Dag_CronTrigger, timeAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

func (c *CronTriggerCreateUsecase) CreateSendSMSTextMedTeamForms(caseId int32) error {

	tCase, _ := c.TUsecase.DataById(Kind_client_cases, caseId)
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextMedTeamForms,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

}

func (c *CronTriggerCreateUsecase) CreateSendSMSTextMiniDBQsDrafts(caseId int32) error {

	tCase, _ := c.TUsecase.DataById(Kind_client_cases, caseId)
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextMiniDBQsDrafts,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))

}
