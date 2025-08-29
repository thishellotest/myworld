package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

const ClientTaskSubject_ITFExpirationWithin90Days = "ITF Expiration within 90 days"

const ClientTaskSubject_ITFExpirationWithPrefix = "ITF Expiration within"

type ClientTaskBuzUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	TUsecase          *TUsecase
	ZohoUsecase       *ZohoUsecase
	ClientTaskUsecase *ClientTaskUsecase
	DataEntryUsecase  *DataEntryUsecase
	StageTransUsecase *StageTransUsecase
	LogUsecase        *LogUsecase
	FieldUsecase      *FieldUsecase
	MapUsecase        *MapUsecase
}

func NewClientTaskBuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	ZohoUsecase *ZohoUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	DataEntryUsecase *DataEntryUsecase,
	StageTransUsecase *StageTransUsecase,
	LogUsecase *LogUsecase,
	FieldUsecase *FieldUsecase,
	MapUsecase *MapUsecase) *ClientTaskBuzUsecase {
	uc := &ClientTaskBuzUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		TUsecase:          TUsecase,
		ZohoUsecase:       ZohoUsecase,
		ClientTaskUsecase: ClientTaskUsecase,
		DataEntryUsecase:  DataEntryUsecase,
		StageTransUsecase: StageTransUsecase,
		LogUsecase:        LogUsecase,
		FieldUsecase:      FieldUsecase,
		MapUsecase:        MapUsecase,
	}

	return uc
}

// HandleCompleteTask 此处关闭保留（ITF）和当前Stage的任务，其余都自动关闭
func (c *ClientTaskBuzUsecase) HandleCompleteTask(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	var dbStage string
	if configs.IsProd() && configs.VBC_CRM_RELEASE == false {
		dealMap, err := c.ZohoUsecase.GetDeal(tCase.Gid())
		if err != nil {
			return err
		}
		zohoStage := dealMap.GetString("Stage")
		if zohoStage == "" {
			return errors.New("zohoStage is empty")
		}
		dbStage, err = c.StageTransUsecase.ZohoStageToDBStage(zohoStage)
		if err != nil {
			c.log.Error(err, zohoStage)
			return err
		}
	} else {
		dbStage = tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	}
	var cond Cond
	cond = And(Eq{"biz_deleted_at": 0, "what_id_gid": tCase.Gid()},
		In("status",
			config_zoho.ClientTaskStatus_Waitingforinput,
			config_zoho.ClinetTaskStatus_NotStarted,
			config_zoho.ClientTaskStatus_Deferred,
			config_zoho.ClientTaskStatus_InProgress),
		Expr("subject not like '"+ClientTaskSubject_ITFExpirationWithPrefix+"%'"))

	//cond = cond.And(Neq{"aa": "22"})
	//cond = cond.And(Neq{"aa": "221"})

	tasks, err := c.TUsecase.ListByCond(Kind_client_tasks, cond)
	if err != nil {
		return err
	}

	var currentStageTasks []*TData
	for k, v := range tasks {
		if !config_vbc.JudgeTaskWhetherBelongsStage(v.CustomFields.TextValueByNameBasic("subject"), dbStage) {

			c.HandleCompletedTask(tCase, dbStage, v.Gid(), v.Id())

			//var data lib.TypeList
			//row := make(lib.TypeMap)
			//row["id"] = v.CustomFields.TextValueByNameBasic("gid")
			//row["Status"] = config_zoho.ClientTaskStatus_Completed
			//data = append(data, row)
			//records := make(lib.TypeMap)
			//records.Set("data", data)
			//_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
			//if err != nil {
			//	c.log.Error(err, InterfaceToString(data))
			//} else {
			//
			//	c.LogUsecase.SaveLog(v.Id(), "HandleCompleteTask", map[string]interface{}{
			//		"dbStage": dbStage,
			//		"caseId":  tCase.Gid(),
			//	})
			//
			//	dataEntry := make(TypeDataEntry)
			//	dataEntry["gid"] = v.CustomFields.TextValueByNameBasic("gid")
			//	dataEntry["status"] = config_zoho.ClientTaskStatus_Completed
			//	_, err = c.DataEntryUsecase.UpdateOne(Kind_client_tasks, dataEntry, FieldName_gid, nil)
			//	if err != nil {
			//		c.log.Error(err, InterfaceToString(dataEntry))
			//	}
			//}
		} else {
			currentStageTasks = append(currentStageTasks, tasks[k])
		}
	}
	// 当前stage有重复任务，只保留一个
	if len(currentStageTasks) >= 2 {
		for i := 0; i < len(currentStageTasks)-1; i++ {
			v := currentStageTasks[i]

			c.HandleCompletedTask(tCase, dbStage, v.Gid(), v.Id())

			//dbCompleted := false
			//if lib.IsProd() && lib.VBC_CRM_RELEASE == false {
			//	var data lib.TypeList
			//	row := make(lib.TypeMap)
			//	row["id"] = v.CustomFields.TextValueByNameBasic("gid")
			//	row["Status"] = config_zoho.ClientTaskStatus_Completed
			//	data = append(data, row)
			//	records := make(lib.TypeMap)
			//	records.Set("data", data)
			//	_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
			//	if err != nil {
			//		c.log.Error(err, InterfaceToString(data))
			//	} else {
			//		dbCompleted = true
			//	}
			//} else {
			//	dbCompleted = true
			//}
			//if dbCompleted {
			//	c.LogUsecase.SaveLog(v.Id(), "HandleCompleteTaskMore", map[string]interface{}{
			//		"dbStage": dbStage,
			//		"caseId":  tCase.Gid(),
			//	})
			//	dataEntry := make(TypeDataEntry)
			//	dataEntry["gid"] = v.CustomFields.TextValueByNameBasic("gid")
			//	dataEntry["status"] = config_zoho.ClientTaskStatus_Completed
			//	_, err = c.DataEntryUsecase.UpdateOne(Kind_client_tasks, dataEntry, FieldName_gid, nil)
			//	if err != nil {
			//		c.log.Error(err, InterfaceToString(dataEntry))
			//	}
			//}
		}
	}

	return nil
}

func (c *ClientTaskBuzUsecase) HandleCompletedTask(tCase *TData, dbStage string, taskGid string, taskId int32) {
	var err error
	dbCompleted := false
	if configs.IsProd() && configs.VBC_CRM_RELEASE == false {
		var data lib.TypeList
		row := make(lib.TypeMap)
		row["id"] = taskGid
		row["Status"] = config_zoho.ClientTaskStatus_Completed
		data = append(data, row)
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
		if err != nil {
			c.log.Error(err, InterfaceToString(data))
		} else {
			dbCompleted = true
		}
	} else {
		dbCompleted = true
	}
	if dbCompleted {
		c.LogUsecase.SaveLog(taskId, "HandleCompleteTaskMore", map[string]interface{}{
			"dbStage": dbStage,
			"caseId":  tCase.Gid(),
		})
		dataEntry := make(TypeDataEntry)
		dataEntry["gid"] = taskGid
		dataEntry["status"] = config_zoho.ClientTaskStatus_Completed
		_, err = c.DataEntryUsecase.UpdateOne(Kind_client_tasks, dataEntry, FieldName_gid, nil)
		if err != nil {
			c.log.Error(err, InterfaceToString(dataEntry))
		}
	}
}

var AutoCreateItfTasksDays = []int{90, 80, 70, 60, 50, 40, 30, 20, 10, 6, 3}

//
//func (c *ClientTaskBuzUsecase) HandleAutoCreateTask(ctx context.Context) error {
//	builder := Dialect(MYSQL).Select("*").From("client_cases").Where(Eq{"deleted_at": 0})
//	builder.And(Eq{"biz_deleted_at": 0})
//	builder.And(NotIn("stages", vbc_config.Stages_AwaitingDecision,
//		vbc_config.Stages_AwaitingPayment,
//		vbc_config.Stages_Completed,
//		vbc_config.Stages_Terminated,
//		vbc_config.Stages_Dormant,
//	))
//	now := time.Now().In(lib.GetVBCDefaultLocation())
//	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, lib.GetVBCDefaultLocation())
//	begin := now.Format(time.DateOnly)
//	end := now.AddDate(0, 0, 91).Format(time.DateOnly)
//
//	builder.And(And(Neq{"itf_expiration": ""}, NotNull{"itf_expiration"}, Gte{"itf_expiration": begin}, Lte{"itf_expiration": end}))
//
//	sql, err := builder.ToBoundSQL()
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	//return nil
//	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	defer sqlRows.Close()
//
//	caches := lib.CacheInit[*TData]()
//
//	fields, err := c.FieldUsecase.ListByKind(Kind_client_cases)
//	if err != nil {
//		return err
//	}
//
//	for sqlRows.Next() {
//		_, row, err := lib.SqlRowsToMap(sqlRows)
//		if err != nil {
//			c.log.Error(err)
//			return err
//		}
//		tData := &TData{
//			CustomFields: c.TUsecase.GenTFields(&caches, Kind_client_cases, row, fields),
//		}
//		itfExpiration := tData.CustomFields.TextValueByNameBasic("itf_expiration")
//		isOk, subDay, err := c.NeedCreateTask(itfExpiration, now)
//		if err != nil {
//			c.log.Error(err)
//			return err
//		}
//		if isOk {
//			caseId := tData.Id()
//			lib.DPrintln(tData.Id(), isOk, subDay)
//			key := MapKeyMapITFClientTask(caseId, subDay, itfExpiration)
//			mval, err := c.MapUsecase.GetForString(key)
//			if err != nil {
//				c.log.Error("caseId: ", caseId, " ", err)
//				return err
//			}
//			if mval != "1" {
//				_, _, err = c.CreateClientTask(tData, subDay)
//				if err != nil {
//					c.log.Error(err)
//					return err
//				}
//				err = c.MapUsecase.Set(key, "1")
//				if err != nil {
//					c.log.Error(err)
//					return err
//				}
//
//				// 关闭其它的ITF
//				retainSubject := GenITFExpirationSubject(subDay)
//				err = c.HandleOtherItfExpTask(tData, retainSubject)
//				if err != nil {
//					c.log.Error("HandleOtherItfExpTask: ", tData.Id(), " ", err)
//				}
//			}
//
//		}
//	}
//	return nil
//}

func (c *ClientTaskBuzUsecase) HandleOtherItfExpTask(tCase *TData, retainSubject string) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	tasks, err := c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"biz_deleted_at": 0, "what_id_gid": tCase.Gid()},
		In("status",
			config_zoho.ClientTaskStatus_Waitingforinput,
			config_zoho.ClinetTaskStatus_NotStarted,
			config_zoho.ClientTaskStatus_Deferred,
			config_zoho.ClientTaskStatus_InProgress),
		Neq{"subject": retainSubject},
		Like{"subject", "ITF Expiration within"},
	))
	if err != nil {
		c.log.Error(err)
		return err
	}
	for _, v := range tasks {
		var data lib.TypeList
		row := make(lib.TypeMap)
		row["id"] = v.CustomFields.TextValueByNameBasic("gid")
		row["Status"] = config_zoho.ClientTaskStatus_Completed
		data = append(data, row)
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
		if err != nil {
			c.log.Error(err, InterfaceToString(data))
		} else {
			c.LogUsecase.SaveLog(v.Id(), "HandleOtherItfExpTask", map[string]interface{}{
				"caseId": tCase.Gid(),
			})
			dataEntry := make(TypeDataEntry)
			dataEntry["gid"] = v.CustomFields.TextValueByNameBasic("gid")
			dataEntry["status"] = config_zoho.ClientTaskStatus_Completed
			_, err = c.DataEntryUsecase.UpdateOne(Kind_client_tasks, dataEntry, FieldName_gid, nil)
			if err != nil {
				c.log.Error(err, InterfaceToString(dataEntry))
			}
		}
	}
	return nil
}

func GenITFExpirationSubject(days int) string {
	return fmt.Sprintf("ITF Expiration within %d days", days)
}

func (c *ClientTaskBuzUsecase) CreateClientTask(tCase *TData, days int) (gid string, row lib.TypeMap, err error) {

	if tCase == nil {
		return "", nil, errors.New("tCase is nil")
	}
	if days <= 0 {
		return "", nil, errors.New("days is wrong")
	}

	caseGid := tCase.CustomFields.TextValueByNameBasic("gid")
	clientGid := tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid)
	userGid := tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid)
	if clientGid == "" {
		return "", nil, errors.New("clientGid is empty")
	}
	if caseGid == "" {
		return "", nil, errors.New("caseGid is empty")
	}
	if userGid == "" {
		return "", nil, errors.New("userGid is empty")
	}

	record := make(lib.TypeMap)
	record.Set("Status", config_zoho.ClientTaskStatus_Completed)
	record.Set("Owner", userGid)
	record.Set("Created_By", config_vbc.User_Dev_gid)
	record.Set("Modified_By", config_vbc.User_Dev_gid)
	record.Set("Who_Id", clientGid)
	record.Set("What_Id", caseGid)
	record.Set("Priority", "High")
	record.Set("Status", "Not Started")
	record.Set("Subject", GenITFExpirationSubject(days))
	record.Set("$se_module", "Deals")
	record.Set("Due_Date", time.Now().Format(time.DateOnly))

	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")

	gid, r, err := c.ZohoUsecase.PostRecordV1(config_zoho.Tasks, record)
	if err != nil {
		c.log.Error("caseId: ", tCase.Gid(), " ", err)
	}

	c.LogUsecase.SaveLog(tCase.Id(), "CreateClientTask", map[string]interface{}{
		"caseId":        tCase.Gid(),
		"clientTaskGid": gid,
	})

	return gid, r, err
}

func (c *ClientTaskBuzUsecase) NeedCreateTask(caseItfExpiration string, currentTime time.Time) (isOk bool, subDay int, err error) {
	if caseItfExpiration == "" {
		return false, 0, errors.New("caseItfExpiration is wrong")
	}
	itfTime, err := time.ParseInLocation(time.DateOnly, caseItfExpiration, configs.GetVBCDefaultLocation())
	if err != nil {
		return false, 0, err
	}
	subDayDuration := itfTime.Unix() - currentTime.Unix()
	//lib.DPrintln("subDayDuration:", subDayDuration)
	subDay = int(subDayDuration / (24 * 3600))

	for _, v := range AutoCreateItfTasksDays {
		if subDay == v {
			return true, subDay, nil
		}
	}
	return false, subDay, nil
}
