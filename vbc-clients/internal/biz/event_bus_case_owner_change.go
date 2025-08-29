package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type CaseOwnerChangeUsecase struct {
	log                        *log.Helper
	conf                       *conf.Data
	CommonUsecase              *CommonUsecase
	FieldOptionUsecase         *FieldOptionUsecase
	FieldUsecase               *FieldUsecase
	UserUsecase                *UserUsecase
	TUsecase                   *TUsecase
	DataEntryUsecase           *DataEntryUsecase
	EventBus                   *EventBus
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase
	CollaboratorbuzUsecase     *CollaboratorbuzUsecase
	MapUsecase                 *MapUsecase
}

func NewCaseOwnerChangeUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	FieldUsecase *FieldUsecase,
	UserUsecase *UserUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	EventBus *EventBus,
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase,
	CollaboratorbuzUsecase *CollaboratorbuzUsecase,
	MapUsecase *MapUsecase,

) *CaseOwnerChangeUsecase {
	uc := &CaseOwnerChangeUsecase{
		log:                        log.NewHelper(logger),
		CommonUsecase:              CommonUsecase,
		conf:                       conf,
		FieldOptionUsecase:         FieldOptionUsecase,
		FieldUsecase:               FieldUsecase,
		UserUsecase:                UserUsecase,
		TUsecase:                   TUsecase,
		DataEntryUsecase:           DataEntryUsecase,
		EventBus:                   EventBus,
		BoxCollaborationBuzUsecase: BoxCollaborationBuzUsecase,
		CollaboratorbuzUsecase:     CollaboratorbuzUsecase,
		MapUsecase:                 MapUsecase,
	}
	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.EventHandle)

	// case插入第一次只会更新gid，后
	//uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleContractSource)

	return uc
}

func (c *CaseOwnerChangeUsecase) HandleContractSource(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {

	if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsNewRecord {
				tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
				if err != nil {
					c.log.Error(err)
				}
				if tCase != nil {

				}
			}
		}
	}
}

func (c *CaseOwnerChangeUsecase) HandleSyncInitCase(tCase TData) error {

	key := MapKeySyncInitCase(tCase.Id())
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val == "" {
		err := c.DoHandleSyncInitCase(tCase)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *CaseOwnerChangeUsecase) DoHandleSyncInitCase(tCase TData) error {

	dataEntry := make(TypeDataEntry)
	dataEntry[DataEntry_gid] = tCase.Gid()
	isOk := false

	if tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs) == "" {
		edUser, err := c.UserUsecase.GetByGid(config_vbc.GetUserEdwardGid())
		if err != nil {
			c.log.Error(err, " GetByGid: ", config_vbc.GetUserEdwardGid())
		}
		if edUser == nil {
			c.log.Error("edUser is nil")
		} else {
			isOk = true
			dataEntry[FieldName_primary_vs] = edUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
		}
	}

	if tCase.CustomFields.TextValueByNameBasic(FieldName_lead_co) == "" {
		elissaUser, err := c.UserUsecase.GetByGid(config_vbc.GetUserElissaGid())
		if err != nil {
			c.log.Error(err, " GetByGid: ", config_vbc.GetUserElissaGid())
		}
		if elissaUser == nil {
			c.log.Error("elissaUser is nil")
		} else {
			isOk = true
			dataEntry[FieldName_lead_co] = elissaUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
		}
	}

	if tCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp) == "" {
		victoriaUser, err := c.UserUsecase.GetByGid(config_vbc.GetUserVictoriaGid())
		if err != nil {
			c.log.Error(err, " GetByGid: ", config_vbc.GetUserVictoriaGid())
		}
		if victoriaUser == nil {
			c.log.Error("victoriaUser is nil")
		} else {
			isOk = true
			dataEntry[FieldName_primary_cp] = victoriaUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname)
		}
	}

	if isOk {
		_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}

// AfterInsertDoContractSource 不能修改了
//func (c *CaseOwnerChangeUsecase) AfterInsertDoContractSource(tCase TData) error {
//	stage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
//
//	contractSource := tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
//	if contractSource != "" {
//		return nil
//	}
//	fieldOption, err := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, FieldName_stages, stage)
//	if err != nil {
//		c.log.Error(err)
//	} else {
//		if fieldOption == nil {
//			c.log.Error("fieldOption is nil")
//		} else {
//			newContractSource := PipelinesToContractSource(fieldOption.Pipelines)
//			if newContractSource == "" {
//				newContractSource = NewCaseDefaultContractSource
//			}
//			if newContractSource != "" {
//				data := make(TypeDataEntry)
//				data[DataEntry_gid] = tCase.Gid()
//				data[FieldName_ContractSource] = newContractSource
//				_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
//				if err != nil {
//					c.log.Error(err)
//				}
//
//			} else {
//				c.log.Error("newContractSource is empty")
//			}
//		}
//	}
//	return nil
//}

func (c *CaseOwnerChangeUsecase) EventHandle(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	modifiedBy string) {

	if kindEntity.Kind == Kind_client_cases && recognizeFieldName == DataEntry_gid {
		for gid, v := range dataEntryOperResult {
			if v.IsUpdated {
				for fieldName, val := range v.DataEntryModifyDataMap {
					if fieldName == FieldName_stages {
						tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
						if err != nil {
							c.log.Error(err)
						}
						if tCase != nil {
							err = c.DoCaseOwnerChange(tCase)
							if err != nil {
								c.log.Error("DoCaseOwnerChange:", err)
							}
							err = c.DoContractSource(*tCase)
							if err != nil {
								c.log.Error("DoContractSource:", err)
							}

							stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
							if !IsForLeadCOStages(stages) {
								leadCo := tCase.CustomFields.TextValueByNameBasic(FieldName_lead_co)
								if leadCo != "" {
									dataEntry := make(TypeDataEntry)
									dataEntry[DataEntry_gid] = tCase.Gid()
									dataEntry[FieldName_lead_co] = ""
									_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, dataEntry, DataEntry_gid, nil)
									if err != nil {
										c.log.Error(err, " Gid: ", tCase.Gid())
									}
								}
							}
						}
					} else if fieldName == FieldName_user_gid {
						tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
						if err != nil {
							c.log.Error(err, " gid: ", gid)
						}
						if tCase != nil {
							c.CollaboratorbuzUsecase.OperationCollaborator(*tCase, val.GetOldVal(fieldName))
						}
					} else if fieldName == FieldName_primary_vs ||
						fieldName == FieldName_primary_cp ||
						fieldName == FieldName_support_cp ||
						fieldName == FieldName_lead_co {
						tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
						if err != nil {
							c.log.Error(err, " gid: ", gid)
						}
						if tCase != nil {
							c.CollaboratorbuzUsecase.OperationCollaboratorByFullName(*tCase, val.GetOldVal(fieldName))
						}
					} else if fieldName == FieldName_client_gid {

						tCase, err := c.TUsecase.DataByGid(Kind_client_cases, gid)
						if err != nil {
							c.log.Error(err, " gid: ", gid)
						}
						if tCase != nil {
							err = c.HandleSyncInitCase(*tCase)
							if err != nil {
								c.log.Error(err, " gid: ", gid)
							}
						}
					}
				}
			}
		}
	}

}

func (c *CaseOwnerChangeUsecase) DoContractSource(tCase TData) error {

	stage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	//if vbc_config.ShouldSkipContractUpdate(stage) {
	//	return nil
	//}

	contractSource := tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)
	if contractSource != "" {
		return nil
	}

	fieldOption, err := c.FieldOptionUsecase.GetByFieldName(Kind_client_cases, FieldName_stages, stage)
	if err != nil {
		c.log.Error(err)
	} else {
		if fieldOption == nil {
			c.log.Error("fieldOption is nil")
		} else {
			newContractSource := PipelinesToContractSource(fieldOption.Pipelines)
			if newContractSource == "" {
				newContractSource = NewCaseDefaultContractSource
			}

			if newContractSource != "" {

				data := make(TypeDataEntry)
				data[DataEntry_gid] = tCase.Gid()
				data[FieldName_ContractSource] = newContractSource

				_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
				if err != nil {
					c.log.Error(err)
				}

			} else {
				c.log.Error("newContractSource is empty")
			}
		}
	}
	return nil
}

func (c *CaseOwnerChangeUsecase) HandleCaseOwnerChange(changeHistoryEntity ChangeHistoryEntity) error {

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
		if err != nil {
			return err
		}
		return c.DoCaseOwnerChange(tCase)
	}
	return nil
}

func (c *CaseOwnerChangeUsecase) UpdateCaseOwnerTo(tCase *TData, userGid string) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	data := make(TypeDataEntry)
	data[DataEntry_gid] = tCase.Gid()
	data[FieldName_user_gid] = userGid

	_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *CaseOwnerChangeUsecase) DoCaseOwnerChange(tCase *TData) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}

	stageValue := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)

	if stageValue == config_vbc.Stages_ClaimAnalysisReview ||
		stageValue == config_vbc.Stages_AmClaimAnalysisReview {
		return c.UpdateCaseOwnerTo(tCase, config_vbc.GetUserSharikaGid())
	} else if stageValue == config_vbc.Stages_StatementUpdatesDraft {
		return c.UpdateCaseOwnerTo(tCase, config_vbc.User_Lili_gid)
	} else if stageValue == config_vbc.Stages_27_AwaitingBankReconciliation ||
		stageValue == config_vbc.Stages_Am27_AwaitingBankReconciliation {
		return c.UpdateCaseOwnerTo(tCase, config_vbc.User_Victoria_gid)
	} else if stageValue == config_vbc.Stages_AwaitingPayment ||
		stageValue == config_vbc.Stages_PreparingDocumentsTinnitusLetter ||
		stageValue == config_vbc.Stages_AwaitingNexusLetter ||
		stageValue == config_vbc.Stages_ClaimAnalysis ||
		stageValue == config_vbc.Stages_FileHLRDraft ||
		stageValue == config_vbc.Stages_AmAwaitingPayment ||
		stageValue == config_vbc.Stages_AmPreparingDocumentsTinnitusLetter ||
		stageValue == config_vbc.Stages_AmAwaitingNexusLetter ||
		stageValue == config_vbc.Stages_AmClaimAnalysis ||
		stageValue == config_vbc.Stages_AmFileHLRDraft {
		return c.UpdateCaseOwnerTo(tCase, config_vbc.User_Edward_gid)
	}

	fieldStruct, err := c.FieldUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return err
	}
	stageField := fieldStruct.GetByFieldName(FieldName_stages)
	if stageField == nil {
		return errors.New("stageField is nil")
	}
	fieldOptionStruct, err := c.FieldOptionUsecase.CacheStructByKind(Kind_client_cases)
	if err != nil {
		return err
	}
	fieldOption := fieldOptionStruct.AllByFieldName(*stageField).GetByValue(stageValue)
	if fieldOption == nil {
		return errors.New("fieldOption is nil")
	}
	var tUser *TData
	if fieldOption.OptionColor == Option_VS_Color {
		tUser, err = c.UserUsecase.GetByFullName(tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs))
		if err != nil {
			return err
		}
	} else if fieldOption.OptionColor == Option_CP_Color {
		tUser, err = c.UserUsecase.GetByFullName(tCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp))
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	if tUser == nil {
		if stageValue == config_vbc.Stages_IncomingRequest || stageValue == config_vbc.Stages_AmIncomingRequest {
			return nil
		}
		return errors.New("tUser is nil: " + stageValue + " : " + InterfaceToString(tCase.Id()))
	}
	data := make(TypeDataEntry)
	data[DataEntry_gid] = tCase.Gid()
	data[FieldName_user_gid] = tUser.Gid()

	_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
	if err != nil {
		return err
	}

	return nil
}
