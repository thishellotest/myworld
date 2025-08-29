package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ChangeHistoryNodelayJobUseacse struct {
	log                             *log.Helper
	MapUsecase                      *MapUsecase
	conf                            *conf.Data
	CommonUsecase                   *CommonUsecase
	ClientCaseSyncbuzUsecase        *ClientCaseSyncbuzUsecase
	TUsecase                        *TUsecase
	DataComboUsecase                *DataComboUsecase
	ChangeHisUsecase                *ChangeHisUsecase
	ZohobuzUsecase                  *ZohobuzUsecase
	TaskFailureLogUsecase           *TaskFailureLogUsecase
	QueueUsecase                    *QueueUsecase
	LeadVSChangeUsecase             *LeadVSChangeUsecase
	CaseOwnerChangeUsecase          *CaseOwnerChangeUsecase
	ActionOnceUsecase               *ActionOnceUsecase
	LeadConversionSummaryBuzUsecase *LeadConversionSummaryBuzUsecase
	AiHttpUsecase                   *AiHttpUsecase
	StatementConditionBuzUsecase    *StatementConditionBuzUsecase
	StatementUsecase                *StatementUsecase
	AiAssistantJobBuzUsecase        *AiAssistantJobBuzUsecase
	PersonalWebformUsecase          *PersonalWebformUsecase
	BoxCollaborationBuzUsecase      *BoxCollaborationBuzUsecase
	CollaboratorbuzUsecase          *CollaboratorbuzUsecase
	LeadcobuzUsecase                *LeadcobuzUsecase
	CollaboratorClientbuzUsecase    *CollaboratorClientbuzUsecase
}

func NewChangeHistoryNodelayJobUseacse(logger log.Logger, MapUsecase *MapUsecase,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ClientCaseSyncbuzUsecase *ClientCaseSyncbuzUsecase,
	TUsecase *TUsecase,
	DataComboUsecase *DataComboUsecase,
	ChangeHisUsecase *ChangeHisUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	QueueUsecase *QueueUsecase,
	LeadVSChangeUsecase *LeadVSChangeUsecase,
	CaseOwnerChangeUsecase *CaseOwnerChangeUsecase,
	ActionOnceUsecase *ActionOnceUsecase,
	LeadConversionSummaryBuzUsecase *LeadConversionSummaryBuzUsecase,
	AiHttpUsecase *AiHttpUsecase,
	StatementConditionBuzUsecase *StatementConditionBuzUsecase,
	StatementUsecase *StatementUsecase,
	AiAssistantJobBuzUsecase *AiAssistantJobBuzUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
	BoxCollaborationBuzUsecase *BoxCollaborationBuzUsecase,
	CollaboratorbuzUsecase *CollaboratorbuzUsecase,
	LeadcobuzUsecase *LeadcobuzUsecase,
	CollaboratorClientbuzUsecase *CollaboratorClientbuzUsecase,
) *ChangeHistoryNodelayJobUseacse {
	return &ChangeHistoryNodelayJobUseacse{
		log:                             log.NewHelper(logger),
		MapUsecase:                      MapUsecase,
		conf:                            conf,
		CommonUsecase:                   CommonUsecase,
		ClientCaseSyncbuzUsecase:        ClientCaseSyncbuzUsecase,
		TUsecase:                        TUsecase,
		DataComboUsecase:                DataComboUsecase,
		ChangeHisUsecase:                ChangeHisUsecase,
		ZohobuzUsecase:                  ZohobuzUsecase,
		TaskFailureLogUsecase:           TaskFailureLogUsecase,
		QueueUsecase:                    QueueUsecase,
		LeadVSChangeUsecase:             LeadVSChangeUsecase,
		CaseOwnerChangeUsecase:          CaseOwnerChangeUsecase,
		ActionOnceUsecase:               ActionOnceUsecase,
		LeadConversionSummaryBuzUsecase: LeadConversionSummaryBuzUsecase,
		AiHttpUsecase:                   AiHttpUsecase,
		StatementConditionBuzUsecase:    StatementConditionBuzUsecase,
		StatementUsecase:                StatementUsecase,
		AiAssistantJobBuzUsecase:        AiAssistantJobBuzUsecase,
		PersonalWebformUsecase:          PersonalWebformUsecase,
		BoxCollaborationBuzUsecase:      BoxCollaborationBuzUsecase,
		CollaboratorbuzUsecase:          CollaboratorbuzUsecase,
		LeadcobuzUsecase:                LeadcobuzUsecase,
		CollaboratorClientbuzUsecase:    CollaboratorClientbuzUsecase,
	}
}

// RunChangeHistoryNodelayJobJob Handle property event
func (c *ChangeHistoryNodelayJobUseacse) RunChangeHistoryNodelayJobJob(ctx context.Context) error {
	go func() {
		fmt.Println("ChangeHistoryNodelayJobUseacse:RunChangeHistoryNodelayJobJob:running")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ChangeHistoryNodelayJobUseacse:RunChangeHistoryNodelayJobJob:Done")
				return
			default:
				divideId, err := c.MapUsecase.GetForInt(Map_Change_histories_divide_nodelay)
				if err != nil {
					c.log.Error(err)
				} else {
					sqlRows, err := c.CommonUsecase.DB().Table(ChangeHistoryEntity{}.TableName()).
						Where("id>? ",
							divideId).Rows()
					if err != nil {
						c.log.Error(err)
					} else {
						if sqlRows != nil {
							newDivideId := int32(0)
							// 此处只能单
							for sqlRows.Next() {
								var entity ChangeHistoryEntity
								err = c.CommonUsecase.DB().ScanRows(sqlRows, &entity)
								if err != nil {
									c.log.Error(err)
								} else {
									newDivideId = entity.ID
									err = c.Do(&entity)
									if err != nil {
										c.log.Error(err)
									}

									c.MapUsecase.Set(Map_Change_histories_divide_nodelay, lib.InterfaceToString(newDivideId))
								}
							}
							err = sqlRows.Close()
							if err != nil {
								c.log.Error(err)
							}
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return nil
}

func (c *ChangeHistoryNodelayJobUseacse) GetLastValueByFieldName(caseId int32, clientId int32, fieldName string) (needUpdate bool, lastValue string, err error) {
	if caseId == 0 && clientId == 0 {
		return false, "", errors.New("GetLastValueByFieldName params is wrong")
	}
	var tClient *TData

	if caseId > 0 {
		tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": caseId, "biz_deleted_at": 0})
		if err != nil {
			return false, "", err
		}
		if tCase == nil {
			return false, "", errors.New("GetLastValueByFieldName: tCase is nil")
		}
		tClient, _, err = c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
		if err != nil {
			return false, "", err
		}

	} else {
		tClient, err = c.TUsecase.DataById(Kind_clients, clientId)
		if err != nil {
			return false, "", err
		}
	}
	if tClient == nil {
		return false, "", errors.New("tClient is nil")
	}

	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_client_gid: tClient.Gid(), DataEntry_biz_deleted_at: 0})
	if err != nil {
		return false, "", err
	}
	var conds []Cond
	destClientId := tClient.Id()
	conds = append(conds, Eq{"kind": Kind_clients, "incr_id": destClientId})
	var caseIds []int32
	for _, v := range cases {
		caseIds = append(caseIds, v.Id())
	}
	if len(caseIds) > 0 {
		conds = append(conds, And(Eq{"kind": Kind_client_cases}, In("incr_id", caseIds)))
	}

	entity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(And(Or(conds...), Eq{"field_name": fieldName}), "id desc")
	if err != nil {
		return false, "", err
	}
	if entity == nil {
		return false, "", errors.New("GetLastValueByFieldName: entity is nil")
	}
	return true, entity.NewValue, nil
}

// Do 生成任务
func (c *ChangeHistoryNodelayJobUseacse) Do(changeHistoryEntity *ChangeHistoryEntity) error {
	if !configs.IsProd() || configs.StoppedZoho { // zoho没有下线时不能开启 clients与client cases同步
		//c.log.Info(ChangeHistoryEntity.ID)
		if changeHistoryEntity.Kind == Kind_client_cases {
			fieldNames := config_vbc.SyncFieldNamesForCase()
			if lib.InArray(changeHistoryEntity.FieldName, fieldNames) {

				needUpdate, newValue, err := c.GetLastValueByFieldName(changeHistoryEntity.IncrId, 0, changeHistoryEntity.FieldName)
				if err != nil {
					return err
				}
				if needUpdate {
					err = c.ClientCaseSyncbuzUsecase.CaseToClient(changeHistoryEntity.IncrId, ClientCaseSyncVo{
						FieldName:  changeHistoryEntity.FieldName,
						FieldValue: newValue,
					}, nil)
					if err != nil {
						c.log.Error(err)
					}
				}
			}
		} else if changeHistoryEntity.Kind == Kind_clients {
			fieldNames := config_vbc.SyncFieldNamesForClient()
			if lib.InArray(changeHistoryEntity.FieldName, fieldNames) {
				tClient, _ := c.TUsecase.DataById(Kind_clients, changeHistoryEntity.IncrId)
				if tClient != nil {
					needUpdate, newValue, err := c.GetLastValueByFieldName(0, changeHistoryEntity.IncrId, changeHistoryEntity.FieldName)
					if err != nil {
						return err
					}
					if needUpdate {
						err = c.ClientCaseSyncbuzUsecase.ClientToCases(tClient.Gid(), ClientCaseSyncVo{
							FieldName:  changeHistoryEntity.FieldName,
							FieldValue: newValue,
						}, nil)
						if err != nil {
							c.log.Error(err)
						}
					}
				}
			}
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases &&
		(changeHistoryEntity.FieldName == FieldName_effective_current_rating ||
			changeHistoryEntity.FieldName == FieldName_new_rating ||
			changeHistoryEntity.FieldName == "amount" ||
			changeHistoryEntity.FieldName == FieldName_stages ||
			changeHistoryEntity.FieldName == FieldName_deal_name ||
			changeHistoryEntity.FieldName == DataEntry_sys__due_date) { // 在初始化case时不能触发， 所有加入新的触发

		// 处理延时处理，因为初始化时，合同信息还没有存储
		c.log.Info("HandleAmount: ", changeHistoryEntity.IncrId, " FieldName: ", changeHistoryEntity.FieldName, " NewValue: ", changeHistoryEntity.NewValue)
		err := c.ZohobuzUsecase.HandleAmount(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " IncrId:", changeHistoryEntity.IncrId, " Id: ", changeHistoryEntity.ID)
			err = c.TaskFailureLogUsecase.Add(TaskType_HandleAmount, 0,
				map[string]interface{}{
					"ChangeHistoryId": changeHistoryEntity.ID,
					"IncrId":          changeHistoryEntity.IncrId,
					"err":             err.Error(),
				})
			if err != nil {
				c.log.Error(err)
			}
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {
		err := c.HandleInitClientCaseChangeHistory(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandleInitClientCaseChangeHistory err: ", err, " caseId: ", changeHistoryEntity.IncrId)
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases &&
		(changeHistoryEntity.FieldName == FieldName_current_rating ||
			changeHistoryEntity.FieldName == FieldName_effective_current_rating) {
		err := c.ZohobuzUsecase.HandleClientCaseName(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " HandleClientCaseName: ", InterfaceToString(changeHistoryEntity.IncrId))
		}
	}

	if changeHistoryEntity.Kind == Kind_clients && changeHistoryEntity.FieldName == FieldName_full_name {
		err := c.HandleClientNameSyncOthers(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " ", InterfaceToString(changeHistoryEntity.IncrId), " ", InterfaceToString(changeHistoryEntity.ID))
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases {

		if changeHistoryEntity.FieldName == FieldName_primary_vs {
			tCase, _ := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
			if tCase != nil {
				err := c.LeadVSChangeUsecase.DoHandleLeadVSSyncClient(*tCase)
				if err != nil {
					c.log.Error(err, " HandleLeadVSSyncClient ", InterfaceToString(changeHistoryEntity.IncrId), " ", InterfaceToString(changeHistoryEntity.ID))
				}

				err = c.LeadVSChangeUsecase.HandleLeadVSChangeForClaimAnalysisToScheduleCall(*changeHistoryEntity)
				if err != nil {
					c.log.Error("HandleLeadVSChangeForClaimAnalysisToScheduleCall:", err, changeHistoryEntity.ID)
				}

				err = c.CollaboratorClientbuzUsecase.DoCollaboratorByChangeHistory(*changeHistoryEntity, *tCase)
				if err != nil {
					c.log.Error("DoCollaboratorByChangeHistory.DoCollaboratorByChangeHistory:", err, changeHistoryEntity.ID)
				}
			} else {
				c.log.Error("tCase is nil: ", changeHistoryEntity.ID)
			}
		}

		if changeHistoryEntity.FieldName == FieldName_lead_co {
			tCase, _ := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
			if tCase != nil {
				err := c.LeadcobuzUsecase.HandleLeadCOSyncClient(*tCase)
				if err != nil {
					c.log.Error("LeadcobuzUsecase.HandleLeadCOSyncClient:", err, changeHistoryEntity.ID)
				}
				err = c.CollaboratorClientbuzUsecase.DoCollaboratorByChangeHistory(*changeHistoryEntity, *tCase)
				if err != nil {
					c.log.Error("DoCollaboratorByChangeHistory.DoCollaboratorByChangeHistory:", err, changeHistoryEntity.ID)
				}
			} else {
				c.log.Error("tCase is nil: ", changeHistoryEntity.ID)
			}
		}

		if changeHistoryEntity.FieldName == FieldName_stages {
			tCase, _ := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
			if tCase != nil {
				err := c.LeadcobuzUsecase.HandleLeadCOSyncClient(*tCase)
				if err != nil {
					c.log.Error("LeadcobuzUsecase.HandleLeadCOSyncClient:", err, changeHistoryEntity.ID)
				}

				err = c.LeadVSChangeUsecase.DoHandleLeadVSSyncClient(*tCase)
				if err != nil {
					c.log.Error("LeadVSChangeUsecase.DoHandleLeadVSSyncClient:", err, changeHistoryEntity.ID)
				}

			} else {
				c.log.Error("tCase is nil: ", changeHistoryEntity.ID)
			}
		}

		if (changeHistoryEntity.FieldName == FieldName_stages && CanHelpUsImproveSurvey(changeHistoryEntity.NewValue)) || changeHistoryEntity.FieldName == FieldName_new_rating {
			err := c.ActionOnceUsecase.HandleHelpUsImproveSurvey(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleHelpUsImproveSurvey:", err, changeHistoryEntity.ID)
			}
		}

		if changeHistoryEntity.FieldName == FieldName_stages {
			if changeHistoryEntity.NewValue == config_vbc.Stages_RecordReview || changeHistoryEntity.NewValue == config_vbc.Stages_AmRecordReview {
				err := c.ActionOnceUsecase.HandleYourRecordsReviewProcessHasBegun(changeHistoryEntity.IncrId)
				if err != nil {
					c.log.Error("HandleYourRecordsReviewProcessHasBegun:", err, changeHistoryEntity.ID)
				}
			}
			if changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQs_Draft ||
				changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQs_Draft {

				err := c.ActionOnceUsecase.HandleEmailMiniDBQsDrafts(changeHistoryEntity.IncrId)
				if err != nil {
					c.log.Error("HandleEmailMiniDBQsDrafts:", err, changeHistoryEntity.ID)
				}

			}
			if changeHistoryEntity.NewValue == config_vbc.Stages_CurrentTreatment ||
				changeHistoryEntity.NewValue == config_vbc.Stages_AmCurrentTreatment {

				err := c.ActionOnceUsecase.HandlePleaseScheduleYourDoctorAppointments(changeHistoryEntity.IncrId)
				if err != nil {
					c.log.Error("HandlePleaseScheduleYourDoctorAppointments:", err, changeHistoryEntity.ID)
				}

			}
			err := c.LeadConversionSummaryBuzUsecase.DoOne(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("LeadConversionSummaryBuzUsecase:DoOne:", err, changeHistoryEntity.ID)
			}
			if changeHistoryEntity.NewValue == config_vbc.Stages_StatementDrafts ||
				changeHistoryEntity.NewValue == config_vbc.Stages_AmStatementDrafts {

				//err := c.PersonalWebformUsecase.HandleUseNewPersonalWebForm(changeHistoryEntity.IncrId)
				//if err != nil {
				//	c.log.Error("HandleUseNewPersonalWebForm:", err, changeHistoryEntity.ID)
				//}

				//tCase, _ := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
				//if tCase != nil {
				//	err = c.StatementConditionBuzUsecase.DoInitStatementCondition(*tCase)
				//	if err != nil {
				//		c.log.Error("DoInitStatementCondition:", err, changeHistoryEntity.ID)
				//	}
				//}

				// ps-gen begin
				//_, err = c.AiHttpUsecase.BizTaskHandle(changeHistoryEntity.IncrId)
				//if err != nil {
				//	c.log.Error("AiHttpUsecase.BizTaskHandle:", err, changeHistoryEntity.ID)
				//}
				// ps-gen end

				// pw begin
				docEmailJobUuid := GenDocEmailJobUuid(changeHistoryEntity.IncrId)
				_, err = c.AiAssistantJobBuzUsecase.BizHttpCreate(docEmailJobUuid, AiAssistantJobInput{
					BizType:         AiAssistantBizType_DocEmailRenew,
					InternalBizType: AssistantInternalBizType_AutoApply,
				})
				if err != nil {
					c.log.Error("AiAssistantJobBuzUsecase.BizHttpCreate:", err, changeHistoryEntity.ID)
				}

				jobUuid := GenAllStatementsJobUuid(changeHistoryEntity.IncrId)
				_, err = c.AiAssistantJobBuzUsecase.BizHttpCreate(jobUuid, AiAssistantJobInput{
					BizType: AiAssistantBizType_SetAllStatement,
				})
				if err != nil {
					c.log.Error("AiAssistantJobBuzUsecase.BizHttpCreate:", err, changeHistoryEntity.ID)
				}
				// pw end

				err = c.ClientCaseSyncbuzUsecase.UpdatePersonalStatementManagerUrl(changeHistoryEntity.IncrId)
				if err != nil {
					c.log.Error("UpdatePersonalStatementManagerUrl:", err, changeHistoryEntity.ID)
				}
			}
		}

		if changeHistoryEntity.FieldName == FieldName_user_gid {
			if changeHistoryEntity.OldValue != "" {
				err := c.BoxCollaborationBuzUsecase.RunHandleDeleteCollaborationByCaseId(Box_collaboration_ow, changeHistoryEntity.IncrId, changeHistoryEntity.OldValue)
				if err != nil {
					c.log.Error("RunHandleDeleteCollaborationByCaseId:", err, " ", changeHistoryEntity.ID)
				}
			}
			c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(changeHistoryEntity.IncrId)

		} else if changeHistoryEntity.FieldName == FieldName_primary_vs {
			if changeHistoryEntity.OldValue != "" {
				err := c.BoxCollaborationBuzUsecase.RunHandleDeleteCollaborationByUserFullName(Box_collaboration_vs, changeHistoryEntity.IncrId, changeHistoryEntity.OldValue)
				if err != nil {
					c.log.Error("RunHandleDeleteCollaborationByUserFullName:", err, " ", changeHistoryEntity.ID)
				}
			}
			c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(changeHistoryEntity.IncrId)
		} else if changeHistoryEntity.FieldName == FieldName_primary_cp {
			if changeHistoryEntity.OldValue != "" {
				err := c.BoxCollaborationBuzUsecase.RunHandleDeleteCollaborationByUserFullName(Box_collaboration_cp, changeHistoryEntity.IncrId, changeHistoryEntity.OldValue)
				if err != nil {
					c.log.Error("RunHandleDeleteCollaborationByUserFullName:", err, " ", changeHistoryEntity.ID)
				}
			}
			c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(changeHistoryEntity.IncrId)
		} else if changeHistoryEntity.FieldName == FieldName_support_cp {
			if changeHistoryEntity.OldValue != "" {
				err := c.BoxCollaborationBuzUsecase.RunHandleDeleteCollaborationByUserFullName(Box_collaboration_support_cp, changeHistoryEntity.IncrId, changeHistoryEntity.OldValue)
				if err != nil {
					c.log.Error("RunHandleDeleteCollaborationByUserFullName:", err, " ", changeHistoryEntity.ID)
				}
			}
			c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(changeHistoryEntity.IncrId)
		} else if changeHistoryEntity.FieldName == FieldName_lead_co {
			if changeHistoryEntity.OldValue != "" {
				err := c.BoxCollaborationBuzUsecase.RunHandleDeleteCollaborationByUserFullName(Box_collaboration_lead_co, changeHistoryEntity.IncrId, changeHistoryEntity.OldValue)
				if err != nil {
					c.log.Error("RunHandleDeleteCollaborationByUserFullName:", err, " ", changeHistoryEntity.ID)
				}
			}
			c.BoxCollaborationBuzUsecase.DoAddPermissionForBox(changeHistoryEntity.IncrId)
		}
	}

	return nil
}

func (c *ChangeHistoryNodelayJobUseacse) HandleInitClientCaseChangeHistory(clientCaseId int32) error {

	InitClientCaseChangeHistoryKey := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "InitClientCaseChangeHistory", clientCaseId)
	val, err := c.MapUsecase.GetForString(InitClientCaseChangeHistoryKey)
	if err != nil {
		return err
	}
	if val != "1" {
		err = c.ZohobuzUsecase.HandleClientCaseName(clientCaseId)
		if err != nil {
			c.log.Error("err: ", err, " clientCaseId: ", clientCaseId)
		}
		c.MapUsecase.Set(InitClientCaseChangeHistoryKey, "1")
	}
	return nil
}

func (c *ChangeHistoryNodelayJobUseacse) HandleClientNameSyncOthers(clientId int32) error {
	tClient, err := c.TUsecase.DataById(Kind_clients, clientId)
	if err != nil {
		return err
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}
	clientGid := tClient.Gid()

	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_client_gid: clientGid, DataEntry_biz_deleted_at: 0, DataEntry_deleted_at: 0})
	if err != nil {
		return err
	}
	var caseGids []string
	for _, v := range cases {
		err = c.ZohobuzUsecase.HandleClientCaseName(v.Id())
		if err != nil {
			c.log.Error(err, " ", v.Id())
		}
		caseGids = append(caseGids, v.Gid())

	}
	err = c.QueueUsecase.PushClientNameChangeJobTasks(context.TODO(), caseGids)
	if err != nil {
		c.log.Error(err)
	}
	return nil
}
