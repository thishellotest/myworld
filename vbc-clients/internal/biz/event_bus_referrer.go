package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	//. "vbc/lib/builder"
)

type EventBusReferrerUsecase struct {
	log                *log.Helper
	conf               *conf.Data
	TUsecase           *TUsecase
	ClientTaskUsecase  *ClientTaskUsecase
	EventBus           *EventBus
	DataEntryUsecase   *DataEntryUsecase
	RecordLogUsecase   *RecordLogUsecase
	ChangeHisUsecase   *ChangeHisUsecase
	FieldOptionUsecase *FieldOptionUsecase
	ReferrerLogUsecase *ReferrerLogUsecase
	ClientCaseUsecase  *ClientCaseUsecase
}

func NewEventBusReferrerUsecase(logger log.Logger,
	conf *conf.Data,
	TUsecase *TUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	EventBus *EventBus,
	DataEntryUsecase *DataEntryUsecase,
	RecordLogUsecase *RecordLogUsecase,
	ChangeHisUsecase *ChangeHisUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	ReferrerLogUsecase *ReferrerLogUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
) *EventBusReferrerUsecase {
	uc := &EventBusReferrerUsecase{
		log:                log.NewHelper(logger),
		conf:               conf,
		TUsecase:           TUsecase,
		ClientTaskUsecase:  ClientTaskUsecase,
		EventBus:           EventBus,
		DataEntryUsecase:   DataEntryUsecase,
		RecordLogUsecase:   RecordLogUsecase,
		ChangeHisUsecase:   ChangeHisUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
		ReferrerLogUsecase: ReferrerLogUsecase,
		ClientCaseUsecase:  ClientCaseUsecase,
	}
	// 有顺序问题，移入到queue处理
	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleEventBusReferrer)
	uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleEventBusReferrerByInsertData)
	return uc
}

func (c *EventBusReferrerUsecase) HandleEventBusReferrerByInsertData(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {

	if kindEntity.Kind == Kind_clients && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsNewRecord {
				lib.DPrintln(gid)
			}
		}
	}

}

func (c *EventBusReferrerUsecase) HandleEventBusReferrer(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {

	if kindEntity.Kind == Kind_clients && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsUpdated {
				for fieldName, v1 := range v.DataEntryModifyDataMap {
					if fieldName == FieldName_referrer_gid {

						//c.ReferrerLogUsecase.Upsert()
						//lib.DPrintln("EventBusReferrerUsecase", gid, v1)
						newValue := v1.GetNewVal(FieldType_dropdown)
						if newValue != "" {
							err := c.HandleReferrerStat(gid, newValue)
							if err != nil {
								c.log.Error(gid, newValue, err)
							}
						}
						//err := c.UpdateByGid(gid, modifiedBy)
						//if err != nil {
						//	c.log.Error(err, "CreateByGid: gid", gid)
						//}
					}
				}
			}
		}
	}
}

// HandleReferrerStat 处理介绍人信息 clientGid：初介绍人Gid， ReferringGid：介绍人Gid
func (c *EventBusReferrerUsecase) HandleReferrerStat(clientGid string, ReferringGid string) error {

	lib.DPrintln("HandleReferrerStat:", clientGid, ReferringGid)

	introducedPerson, err := c.TUsecase.DataByGid(Kind_clients, clientGid)
	if err != nil {
		return err
	}
	if introducedPerson == nil {
		return errors.New("Introduced person is nil")
	}

	introducer, err := c.TUsecase.DataByGid(Kind_clients, ReferringGid)
	if err != nil {
		return err
	}
	if introducer == nil {
		return errors.New("Introducer nil")
	}
	introducerCase, err := c.ClientCaseUsecase.GetOldestCaseByClientGid(ReferringGid)
	if err != nil {
		return err
	}
	ReferrerStage := ""
	if introducerCase != nil {
		ReferrerStage = introducerCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	}
	return c.ReferrerLogUsecase.ReferringClient(clientGid, ReferrerStage, ReferringGid)

}
