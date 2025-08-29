package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	. "vbc/lib/builder"
)

type DueDateUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewDueDateUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *DueDateUsecase {
	uc := &DueDateUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}
	return uc
}

func (c *DueDateUsecase) SyncDueDate() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases,
		And(Eq{"biz_deleted_at": 0, "deleted_at": 0},
			NotIn(FieldName_stages, config_vbc.Stages_Completed,
				config_vbc.Stages_Terminated, config_vbc.Stages_Dormant)))
	if err != nil {
		return err
	}

	for _, v := range cases {
		caseGid := v.Gid()
		cond := And(Eq{"biz_deleted_at": 0, "what_id_gid": caseGid},
			In("status",
				config_zoho.ClientTaskStatus_Waitingforinput,
				config_zoho.ClinetTaskStatus_NotStarted,
				config_zoho.ClientTaskStatus_Deferred,
				config_zoho.ClientTaskStatus_InProgress),
			Expr("subject not like '"+ClientTaskSubject_ITFExpirationWithPrefix+"%'"))

		clientTask, err := c.TUsecase.Data(Kind_client_tasks, cond)
		if err != nil {
			c.log.Error(err)
		} else {
			if clientTask != nil {
				dueData := clientTask.CustomFields.TextValueByNameBasic(TaskFieldName_due_date)
				if dueData != "" {
					err = c.SyncDueDateOne(caseGid, dueData)
					if err != nil {
						c.log.Error(err)
					}
				}
			}
		}
	}
	return nil
}

func (c *DueDateUsecase) SyncDueDateOne(caseGid string, taskDueDate string) error {
	if taskDueDate == "" {
		return nil
	}
	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseDueDate := tCase.CustomFields.TextValueByNameBasic(DataEntry_sys__due_date)
	isOk := false
	if caseDueDate != "" {
		if caseDueDate < taskDueDate {
			isOk = true
		}
	}
	if isOk {
		destData := make(TypeDataEntry)
		destData[DataEntry_gid] = caseGid
		destData[DataEntry_sys__due_date] = taskDueDate
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, destData, DataEntry_gid, nil)
		if err != nil {
			return err
		}
		c.log.Info("SyncDueDateOne: ", caseGid, " oldDueDate: ", caseDueDate, " newDueDate: ", taskDueDate)
	}
	return nil
}
