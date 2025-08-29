package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type AutomaticTaskCreationUsecase struct {
	log               *log.Helper
	conf              *conf.Data
	TUsecase          *TUsecase
	ClientTaskUsecase *ClientTaskUsecase
	EventBus          *EventBus
}

func NewAutomaticTaskCreationUsecase(logger log.Logger,
	conf *conf.Data,
	TUsecase *TUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	EventBus *EventBus,
) *AutomaticTaskCreationUsecase {
	uc := &AutomaticTaskCreationUsecase{
		log:               log.NewHelper(logger),
		conf:              conf,
		TUsecase:          TUsecase,
		ClientTaskUsecase: ClientTaskUsecase,
		EventBus:          EventBus,
	}
	// 有顺序问题，移入到queue处理
	//uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleAutomaticTaskCreation)
	return uc
}

func (c *AutomaticTaskCreationUsecase) HandleAutomaticTaskCreation(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList) {
	if configs.IsProd() {
		return
	}

	if kindEntity.Kind == Kind_client_cases {
		for gid, v := range dataEntryOperResult {
			if v.IsUpdated {
				for fieldName, _ := range v.DataEntryModifyDataMap {
					if fieldName == FieldName_stages {
						//newValue := v1.GetNewVal(FieldType_dropdown)
						err := c.CreateByGid(gid)
						if err != nil {
							c.log.Error(err, "CreateByGid: gid", gid)
						}
					}
				}
			}
		}
	}
}

func (c *AutomaticTaskCreationUsecase) CreateByGid(gid string) error {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.CreateByCase(tCase)

}

func (c *AutomaticTaskCreationUsecase) Create(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	return c.CreateByCase(tCase)
}

func (c *AutomaticTaskCreationUsecase) CreateByCase(tCase *TData) error {
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if values, ok := config_vbc.AutomaticCreationTaskSubjectRelationStages[stages]; ok {
		for _, v := range values {
			userGid := tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid)
			if v.AssignUserGid != "" {
				userGid = v.AssignUserGid
			}
			caseGid := tCase.Gid()
			clientGid := tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid)
			c.ClientTaskUsecase.Create(v.Subject, userGid, v.PlusDays, clientGid, caseGid)
		}
	}

	return nil
}
