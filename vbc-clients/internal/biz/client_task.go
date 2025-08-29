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
	"vbc/lib/uuid"
)

const (
	TaskFieldName_subject     = "subject"
	TaskFieldName_user_gid    = "user_gid"
	TaskFieldName_priority    = "priority"
	TaskFieldName_due_date    = "due_date"
	TaskFieldName_status      = "status"
	TaskFieldName_what_id_gid = "what_id_gid" // caseGid
	TaskFieldName_who_id_gid  = "who_id_gid"  // clientGid
	TaskFieldName_closed_at   = "closed_at"

	TaskFieldName_re_kind   = "re_kind"
	TaskFieldName_se_module = "se_module"
)

type ClientTaskUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	TUsecase         *TUsecase
	ZohoUsecase      *ZohoUsecase
	LogUsecase       *LogUsecase
	KindUsecase      *KindUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewClientTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	ZohoUsecase *ZohoUsecase,
	LogUsecase *LogUsecase,
	KindUsecase *KindUsecase,
	DataEntryUsecase *DataEntryUsecase) *ClientTaskUsecase {
	uc := &ClientTaskUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		ZohoUsecase:      ZohoUsecase,
		LogUsecase:       LogUsecase,
		KindUsecase:      KindUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *ClientTaskUsecase) Create(subject string, userGid string, plusDays int, clientGid string, caseGid string) error {
	data := make(lib.TypeMap)
	data.Set(DataEntry_gid, uuid.UuidWithoutStrike())
	data.Set(TaskFieldName_subject, subject)
	data.Set(TaskFieldName_priority, "High")
	data.Set(TaskFieldName_status, "Not Started")
	data.Set(TaskFieldName_what_id_gid, caseGid)
	data.Set(TaskFieldName_who_id_gid, clientGid)
	data.Set(TaskFieldName_se_module, "Deals")
	data.Set(TaskFieldName_user_gid, userGid)
	t := time.Now()
	t = t.AddDate(0, 0, plusDays)
	t = t.In(configs.GetVBCDefaultLocation())
	data.Set(TaskFieldName_due_date, t.Format(time.DateOnly))
	_, err := c.DataEntryUsecase.HandleOne(Kind_client_tasks, TypeDataEntry(data), DataEntry_gid, nil)
	return err
}

func (c *ClientTaskUsecase) TasksByCaseGid(caseGid string) (tList []*TData, err error) {
	return c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"what_id_gid": caseGid,
		"biz_deleted_at": 0,
		"se_module":      "Deals"}, In("status",
		config_zoho.ClientTaskStatus_Waitingforinput,
		config_zoho.ClinetTaskStatus_NotStarted,
		config_zoho.ClientTaskStatus_Deferred,
		config_zoho.ClientTaskStatus_InProgress)))
}

func (c *ClientTaskUsecase) HandleAutomationCompleteTask(clientCaseId int32) error {
	clientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": clientCaseId})
	if err != nil {
		return err
	}
	if clientCase == nil {
		return errors.New("clientCase is nil.")
	}
	stages := clientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	clientCaseGid := clientCase.CustomFields.TextValueByNameBasic(FieldName_gid)

	//if clientCaseGid != "6159272000001184003" {
	//	return nil
	//}

	clientTasks, err := c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"what_id_gid": clientCaseGid,
		"biz_deleted_at": 0,
		"se_module":      "Deals"}, In("status",
		config_zoho.ClientTaskStatus_Waitingforinput,
		config_zoho.ClinetTaskStatus_NotStarted,
		config_zoho.ClientTaskStatus_Deferred,
		config_zoho.ClientTaskStatus_InProgress)))
	if err != nil {
		return err
	}
	var data lib.TypeList
	for _, v := range clientTasks {
		if config_vbc.JudgeTaskNeedCompleteBySubject(v.CustomFields.TextValueByNameBasic("subject"), stages) {
			gid := v.CustomFields.TextValueByNameBasic("gid")
			row := make(lib.TypeMap)
			row["id"] = gid
			row["Status"] = config_zoho.ClientTaskStatus_Completed
			data = append(data, row)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
		if err != nil {
			return errors.New(InterfaceToString(data) + " : " + err.Error())
		}
		return nil
	}
	return nil
}

func (c *ClientTaskUsecase) HandleAutomationCompleteTaskSpecify(clientCaseId int32, oldStage string) error {
	clientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": clientCaseId})
	if err != nil {
		return err
	}
	if clientCase == nil {
		return errors.New("clientCase is nil.")
	}
	//stages := clientCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	clientCaseGid := clientCase.CustomFields.TextValueByNameBasic(FieldName_gid)

	//if clientCaseGid != "6159272000001184003" {
	//	return nil
	//}

	clientTasks, err := c.TUsecase.ListByCond(Kind_client_tasks, And(Eq{"what_id_gid": clientCaseGid,
		"biz_deleted_at": 0,
		"se_module":      "Deals"}, In("status",
		config_zoho.ClientTaskStatus_Waitingforinput,
		config_zoho.ClinetTaskStatus_NotStarted,
		config_zoho.ClientTaskStatus_Deferred,
		config_zoho.ClientTaskStatus_InProgress)))
	if err != nil {
		return err
	}
	var data lib.TypeList
	for _, v := range clientTasks {
		if config_vbc.JudgeTaskNeedCompleteBySubjectSpecify(v.CustomFields.TextValueByNameBasic("subject"), oldStage) {
			gid := v.CustomFields.TextValueByNameBasic("gid")
			row := make(lib.TypeMap)
			row["id"] = gid
			row["Status"] = config_zoho.ClientTaskStatus_Completed

			er := c.LogUsecase.SaveLog(v.CustomFields.NumberValueByNameBasic("id"), Log_FormType_ClientTasks, map[string]interface{}{
				"taskGid":      gid,
				"clientCaseId": clientCaseId,
			})
			if er != nil {
				c.log.Error(er)
			}

			data = append(data, row)
		}
	}
	if len(data) > 0 {
		records := make(lib.TypeMap)
		records.Set("data", data)
		_, err = c.ZohoUsecase.PutRecordsV1(config_zoho.Tasks, records)
		if err != nil {
			return errors.New(InterfaceToString(data) + " : " + err.Error())
		}
		return nil
	}
	return nil
}

type DueDatesResult map[string]DueDatesItem

type DueDatesItem struct {
	DueDate string
}

func (c *ClientTaskUsecase) TasksForOpen(fromKindEntity KindEntity, gid string) (TDataList, error) {
	whichFieldName := ""
	if fromKindEntity.Kind == Kind_client_cases {
		whichFieldName = TaskFieldName_what_id_gid
	} else if fromKindEntity.Kind == Kind_clients {
		whichFieldName = TaskFieldName_who_id_gid
	} else {
		return nil, errors.New("Nonsupport")
	}
	kindEntity, _ := c.KindUsecase.GetByKind(Kind_client_tasks)

	query := MySQL().Select("*").
		From(kindEntity.KindTableName()).
		Where(And(Eq{DataEntry_biz_deleted_at: 0,
			DataEntry_deleted_at: 0,
		}, In(whichFieldName, gid), Neq{TaskFieldName_status: config_zoho.ClientTaskStatus_Completed})).
		OrderBy(DataEntry_updated_at + " desc")
	sql, err := query.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	return c.TUsecase.ListByRawSql(Kind_client_tasks, sql)
}

func (c *ClientTaskUsecase) TasksForClose(fromKindEntity KindEntity, gid string) (TDataList, error) {
	whichFieldName := ""
	if fromKindEntity.Kind == Kind_client_cases {
		whichFieldName = TaskFieldName_what_id_gid
	} else if fromKindEntity.Kind == Kind_clients {
		whichFieldName = TaskFieldName_who_id_gid
	} else {
		return nil, errors.New("Nonsupport")
	}
	kindEntity, _ := c.KindUsecase.GetByKind(Kind_client_tasks)

	query := MySQL().Select("*").
		From(kindEntity.KindTableName()).
		Where(And(Eq{DataEntry_biz_deleted_at: 0,
			DataEntry_deleted_at: 0,
		}, In(whichFieldName, gid), Eq{TaskFieldName_status: config_zoho.ClientTaskStatus_Completed})).
		OrderBy(DataEntry_updated_at + " desc")
	sql, err := query.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	return c.TUsecase.ListByRawSql(Kind_client_tasks, sql)
}

func (c *ClientTaskUsecase) TasksForIndex(fromKindEntity KindEntity, gid string) (TDataList, error) {
	whichFieldName := ""
	if fromKindEntity.Kind == Kind_client_cases {
		whichFieldName = TaskFieldName_what_id_gid
	} else if fromKindEntity.Kind == Kind_clients {
		whichFieldName = TaskFieldName_who_id_gid
	} else {
		return nil, errors.New("Nonsupport")
	}
	kindEntity, _ := c.KindUsecase.GetByKind(Kind_client_tasks)

	query := MySQL().Select("*").
		From(kindEntity.KindTableName()).
		Where(And(Eq{DataEntry_biz_deleted_at: 0,
			DataEntry_deleted_at: 0,
		}, In(whichFieldName, gid), Neq{TaskFieldName_due_date: "",
			TaskFieldName_status: config_zoho.ClientTaskStatus_Completed})).OrderBy(TaskFieldName_due_date + " asc")
	sql, err := query.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	return c.TUsecase.ListByRawSql(Kind_client_tasks, sql)
}

func (c *ClientTaskUsecase) DueDatesByWhoGids(values []string) (DueDatesResult, error) {
	return c.DueDates(TaskFieldName_who_id_gid, values)
}

func (c *ClientTaskUsecase) DueDatesByWhatGids(values []string) (DueDatesResult, error) {
	return c.DueDates(TaskFieldName_what_id_gid, values)
}

func (c *ClientTaskUsecase) DueDates(whichFieldName string, values []string) (DueDatesResult, error) {
	if len(values) == 0 {
		return nil, nil
	}

	kindEntity, err := c.KindUsecase.GetByKind(Kind_client_tasks)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}

	query := MySQL().Select(fmt.Sprintf("%s,min(%s) %s", whichFieldName, TaskFieldName_due_date, TaskFieldName_due_date)).
		From(kindEntity.KindTableName()).
		Where(And(Eq{DataEntry_biz_deleted_at: 0,
			DataEntry_deleted_at: 0,
		}, In(whichFieldName, values), Neq{TaskFieldName_due_date: "",
			TaskFieldName_status: config_zoho.ClientTaskStatus_Completed})).
		GroupBy(whichFieldName)
	sql, err := query.ToBoundSQL()
	c.log.Debug(sql)
	if err != nil {
		return nil, err
	}
	res, err := c.TUsecase.ListByRawSql(Kind_client_tasks, sql)
	if err != nil {
		return nil, err
	}
	dueDatesResult := make(DueDatesResult)
	for _, v := range res {
		dueDatesResult[v.CustomFields.TextValueByNameBasic(whichFieldName)] = DueDatesItem{
			DueDate: v.CustomFields.TextValueByNameBasic(TaskFieldName_due_date),
		}
	}
	for _, v := range values {
		if _, ok := dueDatesResult[v]; !ok {
			dueDatesResult[v] = DueDatesItem{
				DueDate: "",
			}
		}
	}
	return dueDatesResult, nil
}

func (c *ClientTaskUsecase) HandleCloseTime(taskGids []string) error {

	if len(taskGids) == 0 {
		return nil
	}

	res, err := c.TUsecase.ListByCond(Kind_client_tasks, In(DataEntry_gid, taskGids))
	if err != nil {
		return err
	}
	for _, v := range res {
		if v.CustomFields.TextValueByNameBasic(TaskFieldName_status) == config_zoho.ClientTaskStatus_Completed {
			data := make(TypeDataEntry)
			data[DataEntry_gid] = v.Gid()
			data[TaskFieldName_closed_at] = time.Now().Unix()
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_tasks, data, DataEntry_gid, nil)
			if err != nil {
				c.log.Error(err, InterfaceToString(data))
			}
		}
	}

	return nil
}
