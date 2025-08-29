package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/utils"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ItfexpirationUsecase struct {
	log                         *log.Helper
	CommonUsecase               *CommonUsecase
	conf                        *conf.Data
	TUsecase                    *TUsecase
	ZohoUsecase                 *ZohoUsecase
	LogUsecase                  *LogUsecase
	DataEntryUsecase            *DataEntryUsecase
	handleITFExpireReminderLock sync.RWMutex
	TaskCreateUsecase           *TaskCreateUsecase
	CronTriggerCreateUsecase    *CronTriggerCreateUsecase
}

func NewItfexpirationUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	ZohoUsecase *ZohoUsecase,
	LogUsecase *LogUsecase,
	DataEntryUsecase *DataEntryUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	CronTriggerCreateUsecase *CronTriggerCreateUsecase) *ItfexpirationUsecase {
	uc := &ItfexpirationUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		TUsecase:                 TUsecase,
		ZohoUsecase:              ZohoUsecase,
		LogUsecase:               LogUsecase,
		DataEntryUsecase:         DataEntryUsecase,
		TaskCreateUsecase:        TaskCreateUsecase,
		CronTriggerCreateUsecase: CronTriggerCreateUsecase,
	}

	return uc
}

func HandleITFExpireReminderHasReminderKey(caseId int32) string {
	return fmt.Sprintf("ITFExpireReminderHasReminder:%d", caseId)
}

func (c *ItfexpirationUsecase) NeedReminderForCase(caseId int32) (bool, error) {
	key := HandleITFExpireReminderHasReminderKey(caseId)
	val, err := c.CommonUsecase.RedisClient().Get(context.TODO(), key).Result()
	if err == redis.Nil {
		return true, nil
	} else if err != nil {
		return false, err
	}
	if val != "" {
		return false, nil
	}

	return true, nil
}

func (c *ItfexpirationUsecase) SetReminderForCase(caseId int32) error {
	key := HandleITFExpireReminderHasReminderKey(caseId)
	return c.CommonUsecase.RedisClient().Set(context.TODO(), key, key, time.Hour*24).Err()
}

func (c *ItfexpirationUsecase) WaitingReminderCases(destITFDate string) (tList []*TData, err error) {

	return c.TUsecase.ListByCond(Kind_client_cases, And(Eq{
		FieldName_itf_expiration: destITFDate,
		"biz_deleted_at":         0,
	}, NotIn(FieldName_stages,
		config_vbc.Stages_AmIncomingRequest,
		config_vbc.Stages_IncomingRequest,
		config_vbc.Stages_FeeScheduleandContract,
		config_vbc.Stages_AmInformationIntake,
		config_vbc.Stages_AmContractPending,
		config_vbc.Stages_AwaitingDecision,
		config_vbc.Stages_AmAwaitingDecision,
		config_vbc.Stages_AwaitingPayment,
		config_vbc.Stages_AmAwaitingPayment,
		config_vbc.Stages_27_AwaitingBankReconciliation,
		config_vbc.Stages_Am27_AwaitingBankReconciliation,
		config_vbc.Stages_Completed,
		config_vbc.Stages_AmCompleted,
		config_vbc.Stages_Terminated,
		config_vbc.Stages_AmTerminated,
		config_vbc.Stages_Dormant,
		config_vbc.Stages_AmDormant,
		config_vbc.Stages_AwaitingClientRecords,
		config_vbc.Stages_AmAwaitingClientRecords,
		config_vbc.Stages_STRRequestPending,
		config_vbc.Stages_AmSTRRequestPending,
		config_vbc.Stages_RecordReview,
		config_vbc.Stages_AmRecordReview,
		config_vbc.Stages_ClaimAnalysis,
		config_vbc.Stages_AmClaimAnalysis,
		config_vbc.Stages_ClaimAnalysisReview,
		config_vbc.Stages_AmClaimAnalysisReview,
		config_vbc.Stages_ScheduleCall,
		config_vbc.Stages_AmScheduleCall,
		config_vbc.Stages_StatementNotes,
		config_vbc.Stages_AmStatementNotes,
		config_vbc.Stages_StatementDrafts,
		config_vbc.Stages_AmStatementDrafts,
		config_vbc.Stages_StatementReview,
		config_vbc.Stages_AmStatementReview,
		config_vbc.Stages_StatementsFinalized,
		config_vbc.Stages_AmStatementsFinalized,
	)))
}

func (c *ItfexpirationUsecase) HandleITFExpireReminder() error {
	c.handleITFExpireReminderLock.Lock()
	c.handleITFExpireReminderLock.Unlock()

	destTime := time.Now().In(configs.GetVBCDefaultLocation())
	destTime = destTime.AddDate(0, 0, 90)
	destITFDate := destTime.Format(time.DateOnly)

	lib.DPrintln("destITFDate:", destITFDate)
	records, err := c.WaitingReminderCases(destITFDate)
	//
	//records, err := c.TUsecase.ListByCond(Kind_client_cases, And(Eq{
	//	FieldName_itf_expiration: destITFDate,
	//	"biz_deleted_at":         0,
	//}, NotIn(FieldName_stages,
	//	config_vbc.Stages_AmIncomingRequest,
	//	config_vbc.Stages_IncomingRequest,
	//	config_vbc.Stages_FeeScheduleandContract,
	//	config_vbc.Stages_AmInformationIntake,
	//	config_vbc.Stages_AmContractPending,
	//	config_vbc.Stages_AwaitingDecision,
	//	config_vbc.Stages_AmAwaitingDecision,
	//	config_vbc.Stages_AwaitingPayment,
	//	config_vbc.Stages_AmAwaitingPayment,
	//	config_vbc.Stages_27_AwaitingBankReconciliation,
	//	config_vbc.Stages_Am27_AwaitingBankReconciliation,
	//	config_vbc.Stages_Completed,
	//	config_vbc.Stages_AmCompleted,
	//	config_vbc.Stages_Terminated,
	//	config_vbc.Stages_AmTerminated,
	//	config_vbc.Stages_Dormant,
	//	config_vbc.Stages_AmDormant,
	//)))
	if err != nil {
		c.log.Error(err)
		return err
	}
	for k, v := range records {
		caseId := v.Id()
		flag, err := c.NeedReminderForCase(caseId)
		if err != nil {
			c.log.Error("NeedReminderForCase:", err, " caseId: ", caseId)
		}
		if flag {
			er := c.CreateReminderITFExpireEmailTask(*records[k])
			if er != nil {
				c.log.Error("CreateReminderITFExpireEmailTask:", er, " caseId: ", caseId)
			}
			er = c.CreateReminderITFExpireTextTask(*records[k])
			if er != nil {
				c.log.Error("CreateReminderITFExpireEmailTask:", er, " caseId: ", caseId)
			}
			c.SetReminderForCase(caseId)
		}
	}

	return nil
}

func (c *ItfexpirationUsecase) CreateReminderITFExpireEmailTask(tCase TData) error {
	err := c.TaskCreateUsecase.CreateTaskMail(tCase.Id(), MailGenre_ITFDeadlineIn90Days, 0, nil, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (c *ItfexpirationUsecase) CreateReminderITFExpireTextTask(tCase TData) error {
	caseId := tCase.Id()
	timeLocation := GetCaseTimeLocation(&tCase, c.log)
	nextAt := utils.CalDelayDayTime(time.Now(), timeLocation)
	return c.TaskCreateUsecase.CreateTaskWithFrom(caseId, CronTriggerVo{
		HandleSendSMSType: HandleSendSMSTextITFDeadlineIn90Days,
	}, Task_Dag_CronTrigger, nextAt.Unix(),
		Task_FromType_DialpadSMS, InterfaceToString(caseId))
}

//
//func (c *ItfexpirationUsecase) ExecuteCompleteItfTasks() error {
//
//	c.LogUsecase.SaveLog(0, "ExecuteCompleteItfTasks", map[string]interface{}{
//		"time": time.Now().Format(time.RFC3339),
//	})
//	tasks, err := c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"biz_deleted_at": 0},
//		In("status",
//			config_zoho.ClientTaskStatus_Waitingforinput,
//			config_zoho.ClinetTaskStatus_NotStarted,
//			config_zoho.ClientTaskStatus_Deferred,
//			config_zoho.ClientTaskStatus_InProgress),
//		Expr("subject like '"+ClientTaskSubject_ITFExpirationWithPrefix+"%'")))
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	caseGids := make(map[string]bool)
//	for _, v := range tasks {
//		caseGids[v.CustomFields.TextValueByNameBasic("what_id_gid")] = true
//	}
//	if len(caseGids) > 0 {
//		var gids []string
//		for k, _ := range caseGids {
//			gids = append(gids, k)
//		}
//		cases, err := c.TUsecase.ListByCond(Kind_client_cases, In("gid", gids))
//		if err != nil {
//			c.log.Error(err)
//			return err
//		}
//		for _, v := range cases {
//			err := c.HandleCompleteItfTasks(v.Id())
//			if err != nil {
//				c.log.Error(err, "caseId:", v.Id())
//			}
//		}
//	}
//	return nil
//}

//func (c *ItfexpirationUsecase) HandleCompleteItfTasks(caseId int32) error {
//
//	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	if tCase == nil {
//		return errors.New("tCase is nil")
//	}
//	needCloseAll := false
//	dbStage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
//	if dbStage == config_vbc.Stages_VerifyEvidenceReceived ||
//		dbStage == config_vbc.Stages_AwaitingDecision ||
//		dbStage == config_vbc.Stages_AwaitingPayment ||
//		dbStage == config_vbc.Stages_27_AwaitingBankReconciliation ||
//		dbStage == config_vbc.Stages_Completed ||
//		dbStage == config_vbc.Stages_Terminated ||
//		dbStage == config_vbc.Stages_Dormant ||
//		dbStage == config_vbc.Stages_AmAwaitingDecision ||
//		dbStage == config_vbc.Stages_AmAwaitingPayment ||
//		dbStage == config_vbc.Stages_Am27_AwaitingBankReconciliation ||
//		dbStage == config_vbc.Stages_AmCompleted ||
//		dbStage == config_vbc.Stages_AmTerminated ||
//		dbStage == config_vbc.Stages_AmDormant {
//		needCloseAll = true
//	} else {
//		itfTime := tCase.CustomFields.TextValueByNameBasic("itf_expiration")
//		if itfTime != "" {
//			ti, err := time.ParseInLocation(time.DateOnly, itfTime, configs.GetVBCDefaultLocation())
//			if err != nil {
//				c.log.Error(err, " caseId:", tCase.Id())
//				return err
//			}
//			lib.DPrintln(err)
//			currentTime := time.Now().In(configs.GetVBCDefaultLocation())
//			currentTimeStr := currentTime.Format(time.DateOnly)
//			currentTime, _ = time.ParseInLocation(time.DateOnly, currentTimeStr, configs.GetVBCDefaultLocation())
//			if currentTime.After(ti) {
//				needCloseAll = true
//			}
//		}
//	}
//	if needCloseAll {
//		tasks, err := c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"biz_deleted_at": 0, "what_id_gid": tCase.Gid()},
//			In("status",
//				config_zoho.ClientTaskStatus_Waitingforinput,
//				config_zoho.ClinetTaskStatus_NotStarted,
//				config_zoho.ClientTaskStatus_Deferred,
//				config_zoho.ClientTaskStatus_InProgress),
//			Expr("subject like '"+ClientTaskSubject_ITFExpirationWithPrefix+"%'")))
//		if err != nil {
//			c.log.Error(err)
//		}
//		for _, v := range tasks {
//			var data lib.TypeList
//			row := make(lib.TypeMap)
//			row["id"] = v.CustomFields.TextValueByNameBasic("gid")
//			row["Status"] = config_zoho.ClientTaskStatus_Completed
//			data = append(data, row)
//			records := make(lib.TypeMap)
//			records.Set("data", data)
//			_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
//			if err != nil {
//				c.log.Error(err, InterfaceToString(data))
//			} else {
//				c.LogUsecase.SaveLog(v.Id(), "HandleCompleteItfTasks", map[string]interface{}{
//					"dbStage": dbStage,
//					"caseId":  tCase.Gid(),
//				})
//				dataEntry := make(TypeDataEntry)
//				dataEntry["gid"] = v.CustomFields.TextValueByNameBasic("gid")
//				dataEntry["status"] = config_zoho.ClientTaskStatus_Completed
//				_, err = c.DataEntryUsecase.UpdateOne(Kind_client_tasks, dataEntry, FieldName_gid, nil)
//				if err != nil {
//					c.log.Error(err, InterfaceToString(dataEntry))
//				}
//			}
//		}
//	}
//	return nil
//}
