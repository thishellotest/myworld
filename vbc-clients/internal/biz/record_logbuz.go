package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type RecordLogbuzUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	RecordLogUsecase *RecordLogUsecase
	TUsecase         *TUsecase
	ChangeHisUsecase *ChangeHisUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewRecordLogbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	RecordLogUsecase *RecordLogUsecase,
	TUsecase *TUsecase,
	ChangeHisUsecase *ChangeHisUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *RecordLogbuzUsecase {
	uc := &RecordLogbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		RecordLogUsecase: RecordLogUsecase,
		TUsecase:         TUsecase,
		ChangeHisUsecase: ChangeHisUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *RecordLogbuzUsecase) ManualHandleDueDate() error {
	records, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		c.log.Warn(err)
		return nil
	}
	for k, _ := range records {
		err = c.ManualHandleDueDateRow(records[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RecordLogbuzUsecase) ManualHandleDueDateRow(tCase *TData) error {

	entity, err := c.RecordLogUsecase.BizCrmStagesLatest(tCase.Gid())
	if err != nil {
		c.log.Warn(err)
		return err
	}
	if entity == nil {
		changeEntity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(
			Eq{"incr_id": tCase.Id(),
				"kind":       "",
				"field_name": FieldName_stages},
			"id desc")
		if err != nil {
			return err
		}
		if changeEntity != nil {

			stages := changeEntity.NewValue
			currentTime := time.Unix(changeEntity.CreatedAt, 0)
			dueDate := GetDueDate(stages, currentTime)
			modifiedBy := ""

			startTime := int64(0)
			endTime := int64(0)
			if dueDate != "" {
				startTime = TimeDateOnlyToTimestamp(currentTime.Format(time.DateOnly))
				endTime = TimeDateOnlyToTimestamp(dueDate)
			}
			_, err = c.RecordLogUsecase.AddBizCrmStages(tCase.Gid(), stages, startTime, endTime, modifiedBy)
			if err != nil {
				c.log.Error(err)
			}
			_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry{
				DataEntry_gid:           tCase.Gid(),
				DataEntry_sys__due_date: dueDate,
			}, DataEntry_gid, nil)
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}
