package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/utils"
	"vbc/lib/builder"
)

type CronTriggerUsecase struct {
	log                                *log.Helper
	CommonUsecase                      *CommonUsecase
	conf                               *conf.Data
	DialpadbuzUsecase                  *DialpadbuzUsecase
	TUsecase                           *TUsecase
	TaskCreateUsecase                  *TaskCreateUsecase
	SendsmsConditionUsecase            *SendsmsConditionUsecase
	SendsmsClientTasksConditionUsecase *SendsmsClientTasksConditionUsecase
}

func NewCronTriggerUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	DialpadbuzUsecase *DialpadbuzUsecase,
	TUsecase *TUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	SendsmsConditionUsecase *SendsmsConditionUsecase,
	SendsmsClientTasksConditionUsecase *SendsmsClientTasksConditionUsecase) *CronTriggerUsecase {
	uc := &CronTriggerUsecase{
		log:                                log.NewHelper(logger),
		CommonUsecase:                      CommonUsecase,
		conf:                               conf,
		DialpadbuzUsecase:                  DialpadbuzUsecase,
		TUsecase:                           TUsecase,
		TaskCreateUsecase:                  TaskCreateUsecase,
		SendsmsConditionUsecase:            SendsmsConditionUsecase,
		SendsmsClientTasksConditionUsecase: SendsmsClientTasksConditionUsecase,
	}
	return uc
}

func (c *CronTriggerUsecase) Handle(handleSendSMSType HandleSendSMSType, caseId int32, cronTriggerVo CronTriggerVo) error {

	// 2024-09-09 关闭了
	if handleSendSMSType == HandleSendSMSTextCurrentTreatment30Day ||
		handleSendSMSType == HandleSendSMSTextCurrentTreatmentFollowingEvery30Day ||
		handleSendSMSType == HandleSendSMSTextAwaitingPaymentTaskOpen30Days {
		return nil
	}

	if c.DialpadbuzUsecase.NeedLimit(handleSendSMSType) {
		isLimit, err := c.DialpadbuzUsecase.IsLimit(handleSendSMSType, caseId)
		if err != nil {
			return err
		}
		if isLimit {
			c.log.Info(fmt.Sprintf("%s caseId:%d is limited", handleSendSMSType, caseId))
			return nil
		}
	}

	tCase, err := c.TUsecase.Data(Kind_client_cases, builder.Eq{"id": caseId, "biz_deleted_at": 0, "deleted_at": 0})
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	// 判断任务是否应该继续执行
	verify, err := c.VerifyCondition(handleSendSMSType, tCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if !verify {
		c.log.Info(InterfaceToString(caseId) + " VerifyCondition:" + InterfaceToString(verify))
		return nil
	}

	// 判断client tasks是否满足条件，满足后发送，不满足此次停止

	verifyClientTasks, err := c.VerifyClientTasksCondition(handleSendSMSType, tCase)
	if err != nil {
		c.log.Error(err)
		return err
	}
	c.log.Info(InterfaceToString(caseId) + " VerifyClientTasksCondition:" + InterfaceToString(verifyClientTasks))
	if verifyClientTasks {
		err = c.DialpadbuzUsecase.HandleSendSMS(handleSendSMSType, caseId, cronTriggerVo)
		if err != nil {
			c.log.Error(err, " caseId:", InterfaceToString(caseId), "handleSendSMSType:", handleSendSMSType)
			return err
		}
	}

	err = c.HandleNextCronTrigger(handleSendSMSType, tCase)
	if err != nil {
		c.log.Error(err, " caseId:", InterfaceToString(caseId), "handleSendSMSType:", handleSendSMSType)
		return err
	}

	return nil
}

// VerifyCondition 判断任务是否满足条件
func (c *CronTriggerUsecase) VerifyCondition(handleSendSMSType HandleSendSMSType, tCase *TData) (verify bool, err error) {
	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	// todo:lgl 需要完成所有条件的验证
	if handleSendSMSType == HandleSendSMSTextGettingStartedEmail ||
		handleSendSMSType == HandleSendSMSTextGettingStartedEmailTaskLongerThan30Days {
		return c.SendsmsConditionUsecase.VerifyTextGettingStartedEmail(tCase)
	} else if handleSendSMSType == HandleSendSMSTextAwaitingClientRecordsLongerThan30Days {
		return c.SendsmsConditionUsecase.VerifyTextAwaitingClientRecords(tCase)
	} else if handleSendSMSType == HandleSendSMSTextSTRRequestPendingLongerThan30Days ||
		handleSendSMSType == HandleSendSMSTextSTRRequestPending45Days {
		return c.SendsmsConditionUsecase.VerifyTextSTRRequestPending(tCase)
	} else if handleSendSMSType == HandleSendSMSTextStatementFinalized {
		return c.SendsmsConditionUsecase.VerifyTextStatementsFinalized(tCase)
	} else if handleSendSMSType == HandleSendSMSTextCurrentTreatment ||
		handleSendSMSType == HandleSendSMSTextCurrentTreatment30Day ||
		handleSendSMSType == HandleSendSMSTextCurrentTreatmentFollowingEvery30Day {
		return c.SendsmsConditionUsecase.VerifyTextCurrentTreatment(tCase)
	} else if handleSendSMSType == HandleSendSMSTextAwaitingDecision30Days ||
		handleSendSMSType == HandleSendSMSTextAwaitingDecisionEveryFollowing30Days {
		return c.SendsmsConditionUsecase.VerifyTextAwaitingDecision(tCase)
	} else if handleSendSMSType == HandleSendSMSTextAwaitingPayment ||
		handleSendSMSType == HandleSendSMSTextAwaitingPaymentAfter14Days ||
		handleSendSMSType == HandleSendSMSTextAwaitingPaymentTaskOpen30Days {
		return c.SendsmsConditionUsecase.VerifyTextAwaitingPayment(tCase)
	}

	return true, nil
}

// VerifyClientTasksCondition 判断任务是否满足条件(只有非实时触发的才能判断)
func (c *CronTriggerUsecase) VerifyClientTasksCondition(handleSendSMSType HandleSendSMSType, tCase *TData) (verify bool, err error) {

	if handleSendSMSType == HandleSendSMSTextGettingStartedEmailTaskLongerThan30Days ||
		handleSendSMSType == HandleSendSMSTextAwaitingClientRecordsLongerThan30Days ||
		handleSendSMSType == HandleSendSMSTextSTRRequestPendingLongerThan30Days ||
		handleSendSMSType == HandleSendSMSTextSTRRequestPending45Days {

		return true, nil

		//existTask, err := c.SendsmsClientTasksConditionUsecase.WhetherExistsTask(
		//	tCase.CustomFields.TextValueByNameBasic("gid"),
		//	tCase.CustomFields.TextValueByNameBasic(FieldName_stages))
		//if err != nil {
		//	c.log.Error(err)
		//	return false, err
		//}
		//if existTask {
		//	return true, nil
		//}
		//return false, nil
	} else if handleSendSMSType == HandleSendSMSTextStatementFinalizedEvery14Days {
		if tCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_StatementsFinalized &&
			tCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_AmStatementsFinalized {
			return false, nil
		}
	} else if handleSendSMSType == HandleSendSMSTextSTRRequestPending45Days {
		if tCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_STRRequestPending &&
			tCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_AmSTRRequestPending {
			return false, nil
		}
	}
	return true, nil
}
func (c *CronTriggerUsecase) HandleNextCronTrigger(handleSendSMSType HandleSendSMSType, tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	caseId := tCase.CustomFields.NumberValueByNameBasic("id")
	if handleSendSMSType == HandleSendSMSTextSTRRequestPendingLongerThan30Days ||
		handleSendSMSType == HandleSendSMSTextSTRRequestPending45Days {
		// 改为60天，所以需要+30天
		// HandleSendSMSTextSTRRequestPending45Days 此任务，改为每隔30天就执行一次
		nextTime, err := utils.CalIntervalDayTime(time.Now(), 30, "08:00", timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTaskWithFrom(caseId,
			CronTriggerVo{
				HandleSendSMSType: HandleSendSMSTextSTRRequestPending45Days,
			}, Task_Dag_CronTrigger, nextTime.Unix(),
			Task_FromType_DialpadSMS, InterfaceToString(caseId),
		)
		if err != nil {
			c.log.Error(err)
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextCurrentTreatment30Day ||
		handleSendSMSType == HandleSendSMSTextCurrentTreatmentFollowingEvery30Day {

		// 关闭2024-09-09
		//nextTime, err := utils.CalIntervalDayTime(time.Now(), 30, "08:05")
		//if err != nil {
		//	return err
		//}
		//err = c.TaskCreateUsecase.CreateTaskWithFrom(caseId,
		//	CronTriggerVo{
		//		HandleSendSMSType: HandleSendSMSTextCurrentTreatmentFollowingEvery30Day,
		//	}, Task_Dag_CronTrigger, nextTime.Unix(),
		//	Task_FromType_DialpadSMS, InterfaceToString(caseId),
		//)
		//if err != nil {
		//	c.log.Error(err)
		//	return err
		//}

	} else if handleSendSMSType == HandleSendSMSTextAwaitingDecision30Days ||
		handleSendSMSType == HandleSendSMSTextAwaitingDecisionEveryFollowing30Days {
		// 9月9日：30 to 45 from ED
		nextTime, err := utils.CalIntervalDayTime(time.Now(), 45, "08:00", timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTaskWithFrom(caseId,
			CronTriggerVo{
				HandleSendSMSType: HandleSendSMSTextAwaitingDecisionEveryFollowing30Days,
			}, Task_Dag_CronTrigger, nextTime.Unix(),
			Task_FromType_DialpadSMS, InterfaceToString(caseId),
		)
		if err != nil {
			c.log.Error(err)
			return err
		}
	} else if handleSendSMSType == HandleSendSMSTextStatementFinalizedEvery14Days {
		nextTime, err := utils.CalIntervalDayTime(time.Now(), 14, "08:00", timeLocation)
		if err != nil {
			return err
		}
		err = c.TaskCreateUsecase.CreateTaskWithFrom(caseId,
			CronTriggerVo{
				HandleSendSMSType: HandleSendSMSTextStatementFinalizedEvery14Days,
			}, Task_Dag_CronTrigger, nextTime.Unix(),
			Task_FromType_DialpadSMS, InterfaceToString(caseId),
		)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	return nil
}
