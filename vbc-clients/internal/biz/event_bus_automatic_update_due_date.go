package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type AutomaticUpdateDueDateUsecase struct {
	log                *log.Helper
	conf               *conf.Data
	TUsecase           *TUsecase
	ClientTaskUsecase  *ClientTaskUsecase
	EventBus           *EventBus
	DataEntryUsecase   *DataEntryUsecase
	RecordLogUsecase   *RecordLogUsecase
	ChangeHisUsecase   *ChangeHisUsecase
	FieldOptionUsecase *FieldOptionUsecase
}

func NewAutomaticUpdateDueDateUsecase(logger log.Logger,
	conf *conf.Data,
	TUsecase *TUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	EventBus *EventBus,
	DataEntryUsecase *DataEntryUsecase,
	RecordLogUsecase *RecordLogUsecase,
	ChangeHisUsecase *ChangeHisUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
) *AutomaticUpdateDueDateUsecase {
	uc := &AutomaticUpdateDueDateUsecase{
		log:                log.NewHelper(logger),
		conf:               conf,
		TUsecase:           TUsecase,
		ClientTaskUsecase:  ClientTaskUsecase,
		EventBus:           EventBus,
		DataEntryUsecase:   DataEntryUsecase,
		RecordLogUsecase:   RecordLogUsecase,
		ChangeHisUsecase:   ChangeHisUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
	}
	// 有顺序问题，移入到queue处理
	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleAutomaticUpdateDueDate)
	uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleAutomaticUpdateDueDateByInsertData)
	return uc
}

func (c *AutomaticUpdateDueDateUsecase) HandleAutomaticUpdateDueDateByInsertData(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {
	if configs.Enable_Client_Task_ForCRM {
		return
	}
	//lib.DPrintln("dataEntryOperResult:", dataEntryOperResult, modifiedBy)
	//lib.DPrintln("sourceData:", sourceData, modifiedBy)
	if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsNewRecord {
				c.UpdateByGid(gid, modifiedBy)
			}
		}
	}

}

func (c *AutomaticUpdateDueDateUsecase) HandleAutomaticUpdateDueDate(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {

	if configs.Enable_Client_Task_ForCRM {
		return
	}

	if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsUpdated {
				for fieldName, v1 := range v.DataEntryModifyDataMap {
					if fieldName == DataEntry_sys__due_date { // 此处须在FieldName_stages前面操作
						newVal := InterfaceToString(v1.NewVal)
						if newVal != "" {
							endTime := TimeDateOnlyToTimestamp(newVal)
							err := c.RecordLogUsecase.UpdateBizCrmStages(gid, modifiedBy, endTime)
							if err != nil {
								c.log.Error(err, "UpdateBizCrmStages: gid", gid)
							}
						}
					} else if fieldName == FieldName_stages {
						//newValue := v1.GetNewVal(FieldType_dropdown)
						err := c.UpdateByGid(gid, modifiedBy)
						if err != nil {
							c.log.Error(err, "CreateByGid: gid", gid)
						}
					}
				}
			}
		}
	}
}

func (c *AutomaticUpdateDueDateUsecase) UpdateByGid(gid string, modifiedBy string) error {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.UpdateByCase(tCase, modifiedBy)

}

func (c *AutomaticUpdateDueDateUsecase) Create(caseId int32, modifiedBy string) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.UpdateByCase(tCase, modifiedBy)
}

// NoDueDateStages 没有due date的stages
func NoDueDateStages() []string {
	return []string{
		config_vbc.Stages_Completed,
		config_vbc.Stages_Dormant,
		config_vbc.Stages_Terminated,
		config_vbc.Stages_AmCompleted,
		config_vbc.Stages_AmDormant,
		config_vbc.Stages_AmTerminated,
	}
}

func GetDueDate(stages string, currentTime time.Time) (dueDate string) {
	if !lib.InArray(stages, NoDueDateStages()) {
		if values, ok := config_vbc.AutomaticCreationTaskSubjectRelationStages[stages]; ok {
			for _, v := range values {
				t := currentTime
				t = t.AddDate(0, 0, v.PlusDays)
				t = t.In(configs.GetVBCDefaultLocation())
				dueDate = t.Format(time.DateOnly)

			}
		}
		// 没有的统一设置3天
		if dueDate == "" {
			t := currentTime
			t = t.AddDate(0, 0, 3)
			t = t.In(configs.GetVBCDefaultLocation())
			dueDate = t.Format(time.DateOnly)
		}
	}
	return
}

func (c *AutomaticUpdateDueDateUsecase) UpdateByCase(tCase *TData, modifiedBy string) error {
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)

	var dueDate string

	changeEntity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(Eq{"kind": Kind_client_cases,
		"incr_id":    tCase.Id(),
		"field_name": FieldName_stages},
		"id desc")
	if err != nil {
		return err
	}
	currentTime := time.Now().In(configs.GetVBCDefaultLocation())
	if changeEntity != nil {
		if changeEntity.OldValue != "" {
			oldStage, _ := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, FieldName_stages, changeEntity.OldValue)
			if oldStage != nil {
				newStage, _ := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, FieldName_stages, changeEntity.NewValue)
				if newStage != nil {
					if newStage.OptionSort < oldStage.OptionSort {

						t := currentTime
						t = t.AddDate(0, 0, 1)
						t = t.In(configs.GetVBCDefaultLocation())
						dueDate = t.Format(time.DateOnly)
					}
				}
			}
		}
	}

	if dueDate == "" {
		dueDate = GetDueDate(stages, currentTime)
	}

	//if !lib.InArray(stages, NoDueDateStages()) {
	//	if values, ok := vbc_config.AutomaticCreationTaskSubjectRelationStages[stages]; ok {
	//		for _, v := range values {
	//			t := currentTime
	//			t = t.AddDate(0, 0, v.PlusDays)
	//			t = t.In(lib.GetVBCDefaultLocation())
	//			dueDate = t.Format(time.DateOnly)
	//
	//		}
	//	}
	//	// 没有的统一设置3天
	//	if dueDate == "" {
	//		t := currentTime
	//		t = t.AddDate(0, 0, 3)
	//		t = t.In(lib.GetVBCDefaultLocation())
	//		dueDate = t.Format(time.DateOnly)
	//	}
	//}
	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = tCase.Gid()
	dataEntry[DataEntry_sys__due_date] = dueDate
	_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
	if err != nil {
		c.log.Error(err)
	}

	startTime := int64(0)
	endTime := int64(0)
	if dueDate != "" {

		startTime = TimeDateOnlyToTimestamp(currentTime.Format(time.DateOnly))
		endTime = TimeDateOnlyToTimestamp(dueDate)
	}

	err = c.RecordLogUsecase.CloseBizCrmStages(tCase.Gid(), modifiedBy)
	if err != nil {
		c.log.Error(err)
	}
	_, err = c.RecordLogUsecase.AddBizCrmStages(tCase.Gid(), stages, startTime, endTime, modifiedBy)
	if err != nil {
		c.log.Error(err)
	}

	return nil
}
