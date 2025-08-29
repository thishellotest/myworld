package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ChangeHistoryEntity struct {
	ID        int32 `gorm:"primaryKey"`
	Kind      string
	IncrId    int32
	FieldName string
	OldValue  string
	NewValue  string
	CreatedAt int64
}

func ChangeHistoryValueFormat(fieldType string, val interface{}) string {
	v := lib.InterfaceToString(val)
	if fieldType == FieldType_timestamp {
		if v == "0" {
			v = ""
		}
	}
	return v
}

func (ChangeHistoryEntity) TableName() string {
	return "change_history"
}

type ChangeHisUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[ChangeHistoryEntity]
}

func NewChangeHisUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *ChangeHisUsecase {
	uc := &ChangeHisUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type ChangeHistoryUseacse struct {
	log                             *log.Helper
	MapUsecase                      *MapUsecase
	CommonUsecase                   *CommonUsecase
	TaskUsecase                     *TaskUsecase
	TaskCreateUsecase               *TaskCreateUsecase
	TUsecase                        *TUsecase
	DocuSignUsecase                 *DocuSignUsecase
	conf                            *conf.Data
	XeroInvoiceUsecase              *XeroInvoiceUsecase
	TaskFailureLogUsecase           *TaskFailureLogUsecase
	MaCongratsEmailUsecase          *MaCongratsEmailUsecase
	DataComboUsecase                *DataComboUsecase
	FeeUsecase                      *FeeUsecase
	ZohobuzUsecase                  *ZohobuzUsecase
	LogUsecase                      *LogUsecase
	ZohoinfoSyncUsecase             *ZohoinfoSyncUsecase
	ActionOnceUsecase               *ActionOnceUsecase
	ClientTaskUsecase               *ClientTaskUsecase
	DbqsUsecase                     *DbqsUsecase
	MiscUsecase                     *MiscUsecase
	DataEntryUsecase                *DataEntryUsecase
	ClientCaseUsecase               *ClientCaseUsecase
	CronTriggerCreateUsecase        *CronTriggerCreateUsecase
	ClientTaskBuzUsecase            *ClientTaskBuzUsecase
	ItfexpirationUsecase            *ItfexpirationUsecase
	AutomaticTaskCreationUsecase    *AutomaticTaskCreationUsecase
	LeadVSChangeUsecase             *LeadVSChangeUsecase
	AiTaskUsecase                   *AiTaskUsecase
	AiTaskbuzUsecase                *AiTaskbuzUsecase
	LeadConversionSummaryBuzUsecase *LeadConversionSummaryBuzUsecase
	AmUsecase                       *AmUsecase
}

func NewChangeHistoryUseacse(logger log.Logger, MapUsecase *MapUsecase,
	CommonUsecase *CommonUsecase,
	TaskUsecase *TaskUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	TUsecase *TUsecase,
	DocuSignUsecase *DocuSignUsecase,
	conf *conf.Data,
	XeroInvoiceUsecase *XeroInvoiceUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	MaCongratsEmailUsecase *MaCongratsEmailUsecase,
	DataComboUsecase *DataComboUsecase,
	FeeUsecase *FeeUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
	LogUsecase *LogUsecase,
	ZohoinfoSyncUsecase *ZohoinfoSyncUsecase,
	ActionOnceUsecase *ActionOnceUsecase,
	ClientTaskUsecase *ClientTaskUsecase,
	DbqsUsecase *DbqsUsecase,
	MiscUsecase *MiscUsecase,
	DataEntryUsecase *DataEntryUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	CronTriggerCreateUsecase *CronTriggerCreateUsecase,
	ClientTaskBuzUsecase *ClientTaskBuzUsecase,
	ItfexpirationUsecase *ItfexpirationUsecase,
	AutomaticTaskCreationUsecase *AutomaticTaskCreationUsecase,
	LeadVSChangeUsecase *LeadVSChangeUsecase,
	AiTaskUsecase *AiTaskUsecase,
	AiTaskbuzUsecase *AiTaskbuzUsecase,
	LeadConversionSummaryBuzUsecase *LeadConversionSummaryBuzUsecase,
	AmUsecase *AmUsecase,
) *ChangeHistoryUseacse {
	return &ChangeHistoryUseacse{
		log:                             log.NewHelper(logger),
		MapUsecase:                      MapUsecase,
		CommonUsecase:                   CommonUsecase,
		TaskUsecase:                     TaskUsecase,
		TaskCreateUsecase:               TaskCreateUsecase,
		TUsecase:                        TUsecase,
		DocuSignUsecase:                 DocuSignUsecase,
		conf:                            conf,
		XeroInvoiceUsecase:              XeroInvoiceUsecase,
		TaskFailureLogUsecase:           TaskFailureLogUsecase,
		MaCongratsEmailUsecase:          MaCongratsEmailUsecase,
		DataComboUsecase:                DataComboUsecase,
		FeeUsecase:                      FeeUsecase,
		ZohobuzUsecase:                  ZohobuzUsecase,
		LogUsecase:                      LogUsecase,
		ZohoinfoSyncUsecase:             ZohoinfoSyncUsecase,
		ActionOnceUsecase:               ActionOnceUsecase,
		ClientTaskUsecase:               ClientTaskUsecase,
		DbqsUsecase:                     DbqsUsecase,
		MiscUsecase:                     MiscUsecase,
		DataEntryUsecase:                DataEntryUsecase,
		ClientCaseUsecase:               ClientCaseUsecase,
		CronTriggerCreateUsecase:        CronTriggerCreateUsecase,
		ClientTaskBuzUsecase:            ClientTaskBuzUsecase,
		ItfexpirationUsecase:            ItfexpirationUsecase,
		AutomaticTaskCreationUsecase:    AutomaticTaskCreationUsecase,
		LeadVSChangeUsecase:             LeadVSChangeUsecase,
		AiTaskUsecase:                   AiTaskUsecase,
		AiTaskbuzUsecase:                AiTaskbuzUsecase,
		LeadConversionSummaryBuzUsecase: LeadConversionSummaryBuzUsecase,
		AmUsecase:                       AmUsecase,
	}
}

// RunChangeHistoryJob Handle property event
func (c *ChangeHistoryUseacse) RunChangeHistoryJob(ctx context.Context) error {
	go func() {
		fmt.Println("ChangeHistoryUseacse:RunChangeHistoryJob:running")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ChangeHistoryUseacse:RunChangeHistoryJob:Done")
				return
			default:
				divideId, err := c.MapUsecase.GetForInt(Map_Change_histories_divide)
				if err != nil {
					c.log.Error(err)
				} else {
					// 原来延后18秒，改为延后10秒
					sqlRows, err := c.CommonUsecase.DB().Table(ChangeHistoryEntity{}.TableName()).
						Where("id>? and created_at<=?",
							divideId, time.Now().Unix()-10).Rows()
					if err != nil {
						c.log.Error(err)
					} else {
						if sqlRows != nil {
							newDivideId := int32(0)
							for sqlRows.Next() {
								var entity ChangeHistoryEntity
								err = c.CommonUsecase.DB().ScanRows(sqlRows, &entity)
								if err != nil {
									c.log.Error(err)
								} else {
									newDivideId = entity.ID

									// GenTask和GenTaskCRM任务有复，不能同时开启

									if configs.IsProd() {
										err = c.GenTask(&entity)
										if err != nil {
											c.log.Error(err)
										}
									} else {
										err = c.GenTaskCRM(&entity)
										if err != nil {
											c.log.Error(err)
										}
									}

									err = c.GenNormalTask(&entity)
									if err != nil {
										c.log.Error(err)
									}

									c.MapUsecase.Set(Map_Change_histories_divide, lib.InterfaceToString(newDivideId))
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

const AM_PricingVersion_V1 = "33.33%"

func (c *ChangeHistoryUseacse) HandleFlowForAM(changeHistoryEntity *ChangeHistoryEntity) error {

	c.log.Info("HandleFlowForAM ")
	if changeHistoryEntity == nil {
		return errors.New("changeHistoryEntity is nil")
	}

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("tClientCase  is nil:" + strconv.FormatInt(int64(changeHistoryEntity.IncrId), 10))
	}
	email := tClientCase.CustomFields.TextValueByNameBasic(FieldName_email)
	if email == "" {
		return nil
	}

	fieldData := tClientCase.CustomFields
	stages := fieldData.TextValueByNameBasic("stages")
	// 验证是否满足条件
	if fieldData.TextValueByNameBasic("user_gid") == "" {
		return nil
	}
	if stages != config_vbc.Stages_AmInformationIntake && stages != config_vbc.Stages_AmContractPending {
		return nil
	}
	contract := fieldData.TextValueByNameBasic(FieldName_ContractSource)
	if contract != "" && contract != ContractSource_AM {
		return nil
	}

	isPrimaryCaseCalc, _, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
	if err != nil {
		return err
	}
	if !isPrimaryCaseCalc {
		return nil
	}

	if contract == "" {
		// 使用AM
		err = c.ClientCaseUsecase.SaveContractSource(tClientCase.Gid(), ContractSource_AM)
		if err != nil {
			c.log.Error("SaveContractSource error:", err, changeHistoryEntity.ID)
		}
	}

	if stages == config_vbc.Stages_AmInformationIntake {
		err = c.AmUsecase.HandleAmInformationIntake(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandleAmInformationIntake error:", err, changeHistoryEntity.ID)
		}
	} else if stages == config_vbc.Stages_AmContractPending {
		err = c.AmUsecase.HandleAmContractPending(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandleAmContractPending error:", err, changeHistoryEntity.ID)
		} else {
			err = c.ClientCaseUsecase.SavePricingVersion(tClientCase, AM_PricingVersion_V1)
			if err != nil {
				c.log.Error("SaveContractSource error:", err, changeHistoryEntity.ID)
			}
		}
	}

	c.log.Info("ContractSource: ", ContractSource_AM, " ", changeHistoryEntity.IncrId, " ", changeHistoryEntity.ID)

	return nil
}

func (c *ChangeHistoryUseacse) HandleEnvelopeAndFeeEmail(changeHistoryEntity *ChangeHistoryEntity) error {

	c.log.Info("HandleEnvelopeAndFeeEmail ")

	if changeHistoryEntity == nil {
		return errors.New("changeHistoryEntity is nil")
	}

	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("tClientCase  is nil:" + strconv.FormatInt(int64(changeHistoryEntity.IncrId), 10))
	}
	email := tClientCase.CustomFields.TextValueByNameBasic(FieldName_email)
	if email == "" {
		return nil
	}

	fieldData := tClientCase.CustomFields
	// 验证是否满足条件
	if fieldData.NumberValueByName("effective_current_rating") == nil ||
		fieldData.TextValueByNameBasic("user_gid") == "" ||
		fieldData.TextValueByNameBasic("stages") != config_vbc.Stages_FeeScheduleandContract {
		return nil
	}

	contract := fieldData.TextValueByNameBasic(FieldName_ContractSource)
	if contract != "" && contract != ContractSource_VBC {
		return nil
	}

	c.log.Info("ContractSource: ", ContractSource_VBC, " ", changeHistoryEntity.IncrId, " ", changeHistoryEntity.ID)

	// Important: Here the business is sequential
	// HandleEnvelope: Handle Pricing Version, then HandleFeeScheduleCommunicationMail
	err = c.HandleEnvelope(changeHistoryEntity, email)
	if err != nil {
		c.TaskFailureLogUsecase.Add(TaskType_HandleEnvelope, 0, map[string]interface{}{
			"changeHistoryEntityID": changeHistoryEntity.ID,
			"err":                   err.Error(),
		})
		c.log.Error(err)
	}

	err = c.HandleFeeScheduleCommunicationMail(changeHistoryEntity, email)
	if err != nil {
		c.TaskFailureLogUsecase.Add(TaskType_HandleFeeScheduleCommunicationMail, 0, map[string]interface{}{
			"changeHistoryEntityID": changeHistoryEntity.ID,
			"IncrId":                changeHistoryEntity.IncrId,
			"err":                   err.Error(),
		})
		c.log.Error(err)
	}
	return nil
}

// GenTaskCRM Crm的生成任务
func (c *ChangeHistoryUseacse) GenTaskCRM(changeHistoryEntity *ChangeHistoryEntity) error {

	if changeHistoryEntity == nil {
		return errors.New("GenTaskCRM changeHistoryEntity is nil.")
	}
	// 改为直接直执行，不延时
	//if changeHistoryEntity.Kind == Kind_client_cases &&
	//	changeHistoryEntity.FieldName == FieldName_stages {
	//	err := c.AutomaticTaskCreationUsecase.Create(changeHistoryEntity.IncrId)
	//	if err != nil {
	//		c.log.Error(err, " ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId)
	//	}
	//}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {
		if configs.Enable_Client_Task_ForCRM {
			err := c.ClientTaskBuzUsecase.HandleCompleteTask(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error(err, " ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId)
				er := c.TaskFailureLogUsecase.Add(TaskType_HandleAutomationCompleteTask, 0, map[string]interface{}{
					"ChangeHistoryId": changeHistoryEntity.ID,
					"IncrId":          changeHistoryEntity.IncrId,
					"err":             err.Error(),
				})
				if er != nil {
					c.log.Error(er)
				}
				er = c.LogUsecase.SaveLog(changeHistoryEntity.ID, Log_FromType_HandleClientTask, map[string]interface{}{
					"ChangeHistoryId": changeHistoryEntity.ID,
					"IncrId":          changeHistoryEntity.IncrId,
					"err":             err.Error(),
				})
				if er != nil {
					c.log.Error(er)
				}
			}
		}
	}

	return nil
}

// GenNormalTask 每个环境都可以执行，不支持分布式
func (c *ChangeHistoryUseacse) GenNormalTask(changeHistoryEntity *ChangeHistoryEntity) error {
	//var tCase *TData
	var tClient *TData
	var err error
	if changeHistoryEntity == nil {
		return errors.New("changeHistoryEntity is nil.")
	}
	if changeHistoryEntity.Kind == Kind_client_cases {
		//tCase, err = c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
		//if err != nil {
		//	return err
		//}
	} else if changeHistoryEntity.Kind == Kind_clients {
		tClient, err = c.TUsecase.DataById(Kind_clients, changeHistoryEntity.IncrId)
		if err != nil {
			return err
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {
		if changeHistoryEntity.NewValue == config_vbc.Stages_ScheduleCall ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmScheduleCall {
			err = c.ActionOnceUsecase.HandleUpcomingContactInformation(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleUpcomingContactInformation: ", err, " ", changeHistoryEntity.ID)
			}
		}
	}

	if changeHistoryEntity.Kind == Kind_clients && tClient != nil {
		if changeHistoryEntity.FieldName == FieldName_state || changeHistoryEntity.FieldName == FieldName_city {
			if tClient.CustomFields.TextValueByNameBasic(FieldName_state) != "" &&
				tClient.CustomFields.TextValueByNameBasic(FieldName_city) != "" {
				err = c.AiTaskbuzUsecase.HandleReturnTimezone(tClient)
				if err != nil {
					c.log.Error("HandleReturnTimezone: ", err, " ", changeHistoryEntity.ID)
				}
			}
		}
	}

	return nil
}

// GenTask (prod)生成任务，不支持分布式
func (c *ChangeHistoryUseacse) GenTask(changeHistoryEntity *ChangeHistoryEntity) error {
	if changeHistoryEntity == nil {
		return errors.New("changeHistoryEntity is nil.")
	}

	// 此任务必须在最前面，因为价格版本依懒此处
	if changeHistoryEntity.Kind == Kind_client_cases &&
		(changeHistoryEntity.FieldName == "stages" ||
			changeHistoryEntity.FieldName == "client_gid" ||
			changeHistoryEntity.FieldName == "effective_current_rating" ||
			changeHistoryEntity.FieldName == "user_gid" ||
			changeHistoryEntity.FieldName == "email") {

		c.log.Info("HandleAmount: -HandleEnvelopeAndFeeEmail ", changeHistoryEntity.IncrId, " FieldName: ", changeHistoryEntity.FieldName, " NewValue: ", changeHistoryEntity.NewValue)
		err := c.HandleEnvelopeAndFeeEmail(changeHistoryEntity)

		if err != nil {
			c.log.Error("HandleEnvelopeAndFeeEmail:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
		}

		err = c.HandleFlowForAM(changeHistoryEntity)
		if err != nil {
			c.log.Error("HandleFlowForAM:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
		}

	}

	// 事件触发
	if changeHistoryEntity.Kind == Kind_client_cases &&
		(changeHistoryEntity.FieldName == "stages" ||
			changeHistoryEntity.FieldName == "client_gid") {

		// 签合同后，再生成文件，且使用box copy folder api
		err := c.HandleCreateFolderInBoxAndMail(changeHistoryEntity)
		if err != nil {
			c.TaskFailureLogUsecase.Add(TaskType_HandleCreateFolderInBoxAndMail, 0, map[string]interface{}{
				"changeHistoryEntityID": changeHistoryEntity.ID,
				"err":                   err.Error(),
			})
			c.log.Error(err)
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases && ((changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingPayment ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingPayment)) ||
		changeHistoryEntity.FieldName == FieldName_new_rating ||
		changeHistoryEntity.FieldName == FieldName_client_gid) {
		err := c.MaCongratsEmailUsecase.HandleInputTask(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err)
			err = c.TaskFailureLogUsecase.Add(TaskType_MaCongratsEmail_HandleInputTask, 0,
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

	/*
		FieldName_email    = "email"
			FieldName_phone    = "phone"
			FieldName_ssn      = "ssn"
			FieldName_dob      = "dob"
			FieldName_state    = "state"
			FieldName_city     = "city"
			FieldName_address  = "address"
			FieldName_zip_code = "zip_code"
	*/
	// client与client case互相同步
	if (changeHistoryEntity.Kind == Kind_client_cases || changeHistoryEntity.Kind == Kind_clients) &&
		(changeHistoryEntity.FieldName == FieldName_email ||
			changeHistoryEntity.FieldName == FieldName_phone ||
			changeHistoryEntity.FieldName == FieldName_ssn ||
			changeHistoryEntity.FieldName == FieldName_dob ||
			changeHistoryEntity.FieldName == FieldName_state ||
			changeHistoryEntity.FieldName == FieldName_city ||
			changeHistoryEntity.FieldName == FieldName_address ||
			changeHistoryEntity.FieldName == FieldName_zip_code ||
			changeHistoryEntity.FieldName == FieldName_place_of_birth_city ||
			changeHistoryEntity.FieldName == FieldName_place_of_birth_country ||
			changeHistoryEntity.FieldName == FieldName_place_of_birth_state_province ||
			changeHistoryEntity.FieldName == FieldName_current_occupation) {

		// 任务先关闭，把数据全量同步之后开启
		//return nil
		// 无问题方案：必须改那个字段同步哪个字段
		if !configs.StoppedZoho {
			err := c.ZohoinfoSyncUsecase.Sync(changeHistoryEntity.Kind, changeHistoryEntity.IncrId, changeHistoryEntity.FieldName)
			if err != nil {
				c.log.Error(err)

				err = c.TaskFailureLogUsecase.Add(TaskType_ZohoinfoSync, 0,
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
	}

	if changeHistoryEntity.Kind == Kind_client_cases &&
		changeHistoryEntity.FieldName == FieldName_deal_name {
		err := c.ActionOnceUsecase.InitClientCase(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("InitClientCase:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
			err = c.TaskFailureLogUsecase.Add(TaskType_InitClientCase, 0,
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

	if changeHistoryEntity.Kind == Kind_client_cases && ((changeHistoryEntity.FieldName == FieldName_stages &&
		changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingPayment) ||
		changeHistoryEntity.FieldName == FieldName_new_rating ||
		changeHistoryEntity.FieldName == FieldName_effective_current_rating) {

		if configs.IsProd() {
			c.log.Info("HandleInvoice: ", changeHistoryEntity)
			err := c.XeroInvoiceUsecase.HandleInvoice(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error(err)
				err = c.TaskFailureLogUsecase.Add(TaskType_XeroCreateInvoice, 0,
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
	}

	if changeHistoryEntity.Kind == Kind_client_cases && ((changeHistoryEntity.FieldName == FieldName_stages &&
		changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingPayment) ||
		changeHistoryEntity.FieldName == FieldName_am_invoice_amount) {

		if configs.IsProd() {
			c.log.Info("HandleAmInvoice: ", changeHistoryEntity)
			err := c.XeroInvoiceUsecase.HandleAmInvoice(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error(err, " ", "HandleAmInvoice", " ", changeHistoryEntity.IncrId, " ", changeHistoryEntity.ID)
				err = c.TaskFailureLogUsecase.Add(TaskType_XeroAmCreateInvoice, 0,
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
	}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {

		err := c.ClientTaskBuzUsecase.HandleCompleteTask(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId)
			er := c.TaskFailureLogUsecase.Add(TaskType_HandleAutomationCompleteTask, 0, map[string]interface{}{
				"ChangeHistoryId": changeHistoryEntity.ID,
				"IncrId":          changeHistoryEntity.IncrId,
				"err":             err.Error(),
			})
			if er != nil {
				c.log.Error(er)
			}
			er = c.LogUsecase.SaveLog(changeHistoryEntity.ID, Log_FromType_HandleClientTask, map[string]interface{}{
				"ChangeHistoryId": changeHistoryEntity.ID,
				"IncrId":          changeHistoryEntity.IncrId,
				"err":             err.Error(),
			})
			if er != nil {
				c.log.Error(er)
			}
		}
	}

	// 处理同一个client多个client cases时基本信息同步
	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {
		err := c.ActionOnceUsecase.MultiCasesBaseInfoSync(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err)
			err = c.LogUsecase.SaveLog(changeHistoryEntity.ID, Log_FromType_MultiCasesBaseInfoSync, map[string]interface{}{
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

		if changeHistoryEntity.NewValue != "" && (changeHistoryEntity.NewValue != config_vbc.Stages_IncomingRequest &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmIncomingRequest &&
			changeHistoryEntity.NewValue != config_vbc.Stages_FeeScheduleandContract &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmInformationIntake &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmContractPending &&
			changeHistoryEntity.NewValue != config_vbc.Stages_GettingStartedEmail &&
			changeHistoryEntity.NewValue != config_vbc.Stages_Dormant &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmDormant &&
			changeHistoryEntity.NewValue != config_vbc.Stages_Terminated &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmTerminated) { // 当新状态生成时，使用前面判断方法有留下逻辑问题，
			err := c.ActionOnceUsecase.HandleDataCollectionFolder(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleDataCollectionFolder error:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
			}
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_RecordReview ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmRecordReview) {
		err := c.ActionOnceUsecase.HandleCopyRecordReviewFiles(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandleCopyRecordReviewFiles error:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
		}
	}
	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQ_Forms ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQ_Forms) {

		err := c.ActionOnceUsecase.HandleMedicalTeamForms(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandleMedicalTeamForms error:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
			err = c.TaskFailureLogUsecase.Add(TestType_HandleMedicalTeamForms, 0,
				map[string]interface{}{
					"ChangeHistoryId": changeHistoryEntity.ID,
					"IncrId":          changeHistoryEntity.IncrId,
					"err":             err.Error(),
				})
			if err != nil {
				c.log.Error(err)
			}

		} else {
			err = c.ActionOnceUsecase.HandleMedicalTeamFormsReminderEmail(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleMedicalTeamFormsReminderEmail error:", err, " IncrId: ", changeHistoryEntity.IncrId, " ID: ", changeHistoryEntity.ID)
			}
		}

	}

	if changeHistoryEntity.Kind == Kind_client_cases &&
		changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_MedicalTeam || changeHistoryEntity.NewValue == config_vbc.Stages_AmMedicalTeam) {
		err := c.ActionOnceUsecase.HandlePrivateExamsSubmitted(context.TODO(), changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error("HandlePrivateExamsSubmitted ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId, " ", err)
			err = c.TaskFailureLogUsecase.Add(TaskType_HandlePrivateExamsSubmitted, 0,
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

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_MedicalTeamExamsScheduled ||
			changeHistoryEntity.NewValue == config_vbc.Stages_VerifyEvidenceReceived ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingDecision ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmMedicalTeamExamsScheduled ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmVerifyEvidenceReceived ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingDecision) {
		err := c.MiscUsecase.HandleMiscThingsToKnowCPExam(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId)
			err = c.TaskFailureLogUsecase.Add(TaskType_HandleMiscThingsToKnowCPExam, 0,
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

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages &&
		(changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingPayment ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingPayment) {

		err := c.MiscUsecase.HandleRemoveMiscThingsToKnowCPExam(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err, " ID: ", changeHistoryEntity.ID, " IncrId: ", changeHistoryEntity.IncrId)
			err = c.TaskFailureLogUsecase.Add(TaskType_HandleRemoveMiscThingsToKnowCPExam, 0,
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
		tCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err)
		}
		if tCase == nil {
			c.log.Error("tCase is nil: ", changeHistoryEntity.IncrId)
		} else {
			err := c.ZohobuzUsecase.HandleSyncZohoPricingVersion(tCase)
			if err != nil {
				c.log.Error(err, " : ", changeHistoryEntity.IncrId, " : ", changeHistoryEntity.ID)
			}
		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases &&
		changeHistoryEntity.FieldName == FieldName_stages {

		enabledTwoBySMS, err := c.ClientCaseUsecase.EnabledTwoBySMS(changeHistoryEntity.IncrId)
		if err != nil {
			c.log.Error(err)
		}
		if enabledTwoBySMS && err == nil {

			// 任何阶段的变更，都需要取消任务，todo:lgl 如果stages发生改变后，会影响此任务。需要尽快采用mapping方式
			err := c.CronTriggerCreateUsecase.CancelDialpadSMSTasks(changeHistoryEntity.IncrId)
			er := c.ActionOnceUsecase.CancelAutomationCrontabEmailTasks(changeHistoryEntity.IncrId)
			if er != nil {
				c.log.Error("CancelAutomationCrontabEmailTasks err: ", er, " caseId: ", changeHistoryEntity.IncrId)
			}
			if err != nil {
				c.log.Error("CancelDialpadSMSTasks err: ", err, " caseId: ", changeHistoryEntity.IncrId)
			} else {
				if changeHistoryEntity.NewValue == config_vbc.Stages_GettingStartedEmail ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingClientRecords ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingClientRecords ||
					changeHistoryEntity.NewValue == config_vbc.Stages_STRRequestPending ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmSTRRequestPending ||
					changeHistoryEntity.NewValue == config_vbc.Stages_RecordReview ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmRecordReview ||
					changeHistoryEntity.NewValue == config_vbc.Stages_StatementsFinalized ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmStatementsFinalized ||
					changeHistoryEntity.NewValue == config_vbc.Stages_CurrentTreatment ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmCurrentTreatment ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingDecision ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingDecision ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingPayment ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingPayment ||
					changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQ_Forms ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQ_Forms ||
					changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQs_Draft ||
					changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQs_Draft {

					if changeHistoryEntity.NewValue == config_vbc.Stages_GettingStartedEmail {

						err = c.CronTriggerCreateUsecase.CreateGettingStartedEmailByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateGettingStartedEmailByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}

						err = c.CronTriggerCreateUsecase.CreateGettingStartedEmailTaskLongerThan30DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateGettingStartedEmailTaskLongerThan30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}

					} else if changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingClientRecords || changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingClientRecords {
						err = c.CronTriggerCreateUsecase.CreateAwaitingClientRecordsLongerThan30DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateAwaitingClientRecordsLongerThan30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_STRRequestPending || changeHistoryEntity.NewValue == config_vbc.Stages_AmSTRRequestPending {
						err = c.CronTriggerCreateUsecase.CreateSTRRequestPending30DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateSTRRequestPending30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_RecordReview || changeHistoryEntity.NewValue == config_vbc.Stages_AmRecordReview {
						err = c.CronTriggerCreateUsecase.CreateScheduleCallByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateScheduleCallByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}

					} else if changeHistoryEntity.NewValue == config_vbc.Stages_StatementsFinalized ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmStatementsFinalized {
						err = c.CronTriggerCreateUsecase.CreateStatementsFinalizedByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateStatementsFinalizedByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
						err = c.CronTriggerCreateUsecase.CreateStatementsFinalizedEvery14DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateStatementsFinalizedEvery14DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
						err = c.ActionOnceUsecase.HandlePersonalStatementsReadyforYourReview(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("HandlePersonalStatementsReadyforYourReview err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
						err = c.ActionOnceUsecase.HandlePleaseReviewYourPersonalStatementsinSharedFolder(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("HandlePleaseReviewYourPersonalStatementsinSharedFolder err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_CurrentTreatment ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmCurrentTreatment {
						err = c.CronTriggerCreateUsecase.CreateCurrentTreatmentByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateCurrentTreatmentByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
						//err = c.CronTriggerCreateUsecase.CreateCurrentTreatment30DaysByCaseId(changeHistoryEntity.IncrId)
						//if err != nil {
						//	c.log.Error("CreateCurrentTreatment30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						//}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingDecision ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingDecision {
						err = c.CronTriggerCreateUsecase.CreateAwaitingDecision30DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateAwaitingDecision30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_AwaitingPayment ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmAwaitingPayment {

						err = c.CronTriggerCreateUsecase.CreateAwaitingPaymentByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateAwaitingPaymentByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}

						err = c.CronTriggerCreateUsecase.CreateAwaitingPaymentAfter14DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateAwaitingPaymentAfter14DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
						err = c.CronTriggerCreateUsecase.CreateAwaitingPaymentTaskOpen30DaysByCaseId(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateAwaitingPaymentTaskOpen30DaysByCaseId err: ", err, " caseId: ", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQ_Forms ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQ_Forms {
						err = c.CronTriggerCreateUsecase.CreateSendSMSTextMedTeamForms(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateSendSMSTextMedTeamForms:", err, ":", changeHistoryEntity.ID, ":", changeHistoryEntity.IncrId)
						}
					} else if changeHistoryEntity.NewValue == config_vbc.Stages_MiniDBQs_Draft ||
						changeHistoryEntity.NewValue == config_vbc.Stages_AmMiniDBQs_Draft {
						err = c.CronTriggerCreateUsecase.CreateSendSMSTextMiniDBQsDrafts(changeHistoryEntity.IncrId)
						if err != nil {
							c.log.Error("CreateSendSMSTextMiniDBQsDrafts:", err, ":", changeHistoryEntity.ID, ":", changeHistoryEntity.IncrId)
						}

					}
				}
			}

		}
	}

	if changeHistoryEntity.Kind == Kind_client_cases && changeHistoryEntity.FieldName == FieldName_stages {

		if changeHistoryEntity.NewValue == config_vbc.Stages_ClaimAnalysis ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmClaimAnalysis {
			err := c.ActionOnceUsecase.HandleClaimsAnalysisFile(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleClaimsAnalysisFile:", err, changeHistoryEntity.ID)
			}
		}

		if changeHistoryEntity.NewValue == config_vbc.Stages_StatementDrafts ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmStatementDrafts {
			//err := c.ActionOnceUsecase.HandlePersonalStatementsFile(changeHistoryEntity.IncrId)
			//if err != nil {
			//	c.log.Error("HandlePersonalStatementsFile:", err, changeHistoryEntity.ID)
			//}
			err := c.ActionOnceUsecase.HandleDoDocEmailFile(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleDoDocEmailFile:", err, changeHistoryEntity.ID)
			}
		}

		if changeHistoryEntity.NewValue == config_vbc.Stages_StatementsFinalized ||
			changeHistoryEntity.NewValue == config_vbc.Stages_AmStatementsFinalized {
			c.log.Info("HandleDoCopyDocEmailFile:", changeHistoryEntity.ID, changeHistoryEntity.IncrId)
			err := c.ActionOnceUsecase.HandleDoCopyDocEmailFile(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleDoCopyDocEmailFile error:", err, changeHistoryEntity.ID, " caseId: ", changeHistoryEntity.IncrId)
			}
			err = c.ActionOnceUsecase.HandleDoCopyReadPriorToYourDoctorVisitFile(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleDoCopyReadPriorToYourDoctorVisitFile error:", err, changeHistoryEntity.ID)
			}
			err = c.ActionOnceUsecase.HandleCopyPersonalStatementsDoc(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandleCopyPersonalStatementsDoc:", err, changeHistoryEntity.ID, " caseId: ", changeHistoryEntity.IncrId)
			}
		}

		if changeHistoryEntity.NewValue != "" && (changeHistoryEntity.NewValue != config_vbc.Stages_IncomingRequest &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmIncomingRequest &&
			changeHistoryEntity.NewValue != config_vbc.Stages_FeeScheduleandContract &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmInformationIntake &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmContractPending &&
			changeHistoryEntity.NewValue != config_vbc.Stages_GettingStartedEmail &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AwaitingClientRecords &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmAwaitingClientRecords &&
			changeHistoryEntity.NewValue != config_vbc.Stages_STRRequestPending &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmSTRRequestPending &&
			changeHistoryEntity.NewValue != config_vbc.Stages_Terminated &&
			changeHistoryEntity.NewValue != config_vbc.Stages_Dormant &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmTerminated &&
			changeHistoryEntity.NewValue != config_vbc.Stages_AmDormant) {

			err := c.HandlePrimaryCase(changeHistoryEntity.IncrId)
			if err != nil {
				c.log.Error("HandlePrimaryCase:", err, changeHistoryEntity.ID)
			}
		}
	}

	return nil
}

func (c *ChangeHistoryUseacse) HandlePrimaryCase(clientCaseId int32) error {

	key := fmt.Sprintf("%s%s:%d", Map_ActionOnce, "HandlePrimaryCase", clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if val != "1" {
		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("tClientCase is nil")
		}
		clientCaseFields := tClientCase.CustomFields
		usePrimaryCaseCalc := true
		if HasEnabledPrimaryCase(clientCaseFields.TextValueByNameBasic("client_gid")) {
			usePrimaryCaseCalc, _, err = c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				return err
			}
			if usePrimaryCaseCalc { // 说明是primary case需要入库
				entity := make(lib.TypeMap)
				entity.Set(FieldName_is_primary_case, Is_primary_case_YES)
				entity.Set(FieldName_gid, clientCaseFields.TextValueByNameBasic("gid"))
				_, err := c.DataEntryUsecase.HandleOne(Kind_client_cases, TypeDataEntry(entity), FieldName_gid, nil)
				if err != nil {
					return err
				}
			}
		}
		c.MapUsecase.Set(key, "1")
	}
	return nil
}

func (c *ChangeHistoryUseacse) HandleEnvelope(changeHistoryEntity *ChangeHistoryEntity, email string) error {
	if changeHistoryEntity == nil {
		return errors.New("changeHistoryEntity is nil.")
	}
	checkExistKey := fmt.Sprintf("%s%d:%s", Map_CreateEnvelope, changeHistoryEntity.IncrId, email)
	a, _ := c.MapUsecase.GetForInt(checkExistKey)
	if a != 1 {
		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, changeHistoryEntity.IncrId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("tClientCase  is nil:" + strconv.FormatInt(int64(changeHistoryEntity.IncrId), 10))
		}
		fieldData := tClientCase.CustomFields

		_, tContactFields, err := c.DataComboUsecase.Client(fieldData.TextValueByNameBasic("client_gid"))
		if err != nil {
			return err
		}
		if tContactFields == nil {
			return errors.New("tContactFields is nil.")
		}
		//email := tContactFields.TextValueByNameBasic("email")
		//if email == "" {
		//	return errors.New("email is empty.")
		//}

		// 验证是否满足条件
		if fieldData.NumberValueByName("effective_current_rating") == nil ||
			fieldData.TextValueByNameBasic("user_gid") == "" ||
			fieldData.TextValueByNameBasic("stages") != config_vbc.Stages_FeeScheduleandContract {
			return nil
		}
		contract := fieldData.TextValueByNameBasic(FieldName_ContractSource)
		if contract != "" && contract != ContractSource_VBC {
			return nil
		}

		if HasEnabledPrimaryCase(tClientCase.CustomFields.TextValueByNameBasic("client_gid")) {
			isPrimaryCaseCalc, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				return err
			}
			if !isPrimaryCaseCalc { // todo:lgl 需要验证下
				// 设置非primary cases的值
				if primaryCase != nil {
					err = c.ClientCaseUsecase.SavePricingVersion(tClientCase, primaryCase.CustomFields.TextValueByNameBasic(FieldName_s_pricing_version))
					if err != nil {
						c.log.Error(err.Error(), " caseId:", tClientCase.Id())
					}

				} else {
					c.log.Error("primaryCase is nil", " caseId:", tClientCase.Id())
				}
				return nil
			}
		}

		// 只有Primary Case才设置
		if contract == "" {
			// 使用VBC
			err = c.ClientCaseUsecase.SaveContractSource(tClientCase.Gid(), ContractSource_VBC)
			if err != nil {
				c.log.Error("SaveContractSource error:", err, changeHistoryEntity.ID)
			}
		}

		tUser, err := c.TUsecase.Data(Kind_users, And(Eq{"gid": fieldData.TextValueByNameBasic("user_gid")}, Eq{"deleted_at": 0}))
		if err != nil {
			return err
		}
		if tUser == nil {
			return errors.New("tUser is nil")
		}
		// 验证是否满足条件
		if tUser.CustomFields.TextValueByNameBasic("email") == "" {
			return errors.New("User email is empty.")
		}

		// 客户使用价格版本开始
		contractIndex, templateId, pricingVersion, err := c.DocuSignUsecase.ContractTemplateId(tClientCase)
		if err != nil {
			return err
		}
		// 此处较关键，需要先行设置使用版本，使用都使用此版本
		err = c.ClientCaseUsecase.SavePricingVersion(tClientCase, pricingVersion)
		if err != nil {
			return err
		}
		//err = c.MapUsecase.Set(MapKeyClientCaseContractPricingVersion(fieldData.NumberValueByNameBasic("id")), pricingVersion)
		//if err != nil {
		//	c.log.Error(err)
		//	return err
		//}
		// 客户使用价格版本结束

		typeMap := make(lib.TypeMap)
		if c.conf.UseBoxSign {

			typeMap.Set("contractIndex", contractIndex)
			typeMap.Set("templateId", templateId)
			typeMap.Set("signType", Sign_type_box)
			typeMap.Set("clientFirstName", tContactFields.TextValueByNameBasic("first_name"))
			typeMap.Set("clientLastName", tContactFields.TextValueByNameBasic("last_name"))
			typeMap.Set("clientEmail", email)

			if email != "liaogling@gmail.com" && email != "lialing@foxmail.com" && email != "18891706@qq.com" && configs.IsProd() {
				typeMap.Set("agentFirstName", "Edward")
				typeMap.Set("agentLastName", "Bunting Jr.")
				typeMap.Set("agentEmail", "ebunting@vetbenefitscenter.com")
			} else {
				typeMap.Set("agentFirstName", tUser.CustomFields.TextValueByNameBasic("first_name"))
				typeMap.Set("agentLastName", tUser.CustomFields.TextValueByNameBasic("last_name"))
				typeMap.Set("agentEmail", tUser.CustomFields.TextValueByNameBasic("email"))
			}
		} else if c.conf.UseAdobeSign {
			// todo: handle templateId logic CBJCHBCAABAAxqRFmyaz9biyPYcmsoRYRwOxIAW8nifY
			typeMap.Set("templateId", templateId)
			typeMap.Set("signType", Sign_type_adobe)
			typeMap.Set("clientFirstName", fieldData.TextValueByNameBasic("first_name"))
			typeMap.Set("clientLastName", fieldData.TextValueByNameBasic("last_name"))
			typeMap.Set("clientEmail", fieldData.TextValueByNameBasic("email"))
			typeMap.Set("agentFirstName", tUser.CustomFields.TextValueByNameBasic("first_name"))
			typeMap.Set("agentLastName", tUser.CustomFields.TextValueByNameBasic("last_name"))
			typeMap.Set("agentEmail", tUser.CustomFields.TextValueByNameBasic("email"))
		} else {
			typeMap.Set("docusignTemplateId", templateId)
			typeMap.Set("clientName", fieldData.TextValueByNameBasic("first_name")+" "+fieldData.TextValueByNameBasic("last_name"))
			typeMap.Set("clientEmail", fieldData.TextValueByNameBasic("email"))
			typeMap.Set("agentName", tUser.CustomFields.TextValueByNameBasic("name"))
			typeMap.Set("agentEmail", tUser.CustomFields.TextValueByNameBasic("email"))
		}

		err = c.TaskUsecase.CommonUsecase.DB().Save(&TaskEntity{
			IncrId:    changeHistoryEntity.IncrId,
			TaskInput: typeMap.ToString(),
			Event:     Task_Dag_CreateEnvelopeAndSent,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}).Error
		if err != nil {
			return err
		} else {
			err = c.MapUsecase.SetInt(checkExistKey, 1)
			if err != nil {
				return err
			}
			tClientCaseFields := tClientCase.CustomFields
			clientCaseContractBasicDataVo, err := CreateClientCaseContractBasicDataVoByCase(tClientCaseFields)
			if err != nil {
				c.log.Error(err)
				return err
			}
			clientCaseContractBasicDataKey := fmt.Sprintf("%s%d", Map_ClientCaseContractBasicData, changeHistoryEntity.IncrId)
			err = c.MapUsecase.Set(clientCaseContractBasicDataKey, InterfaceToString(clientCaseContractBasicDataVo))
			c.log.Info("HandleAmount: -clientCaseContractBasicDataKey ", clientCaseContractBasicDataKey, " ", InterfaceToString(clientCaseContractBasicDataVo))
			if err != nil {
				c.log.Error(err)
				return err
			}

			er := c.ZohobuzUsecase.HandleAmount(changeHistoryEntity.IncrId)
			if er != nil {
				c.log.Error(er, "HandleAmount IncrId:", changeHistoryEntity.IncrId, " Id: ", changeHistoryEntity.ID)
			}
		}
	}
	return nil
}

// HandleFeeScheduleCommunicationMail todo:lgl 多价格版本
func (c *ChangeHistoryUseacse) HandleFeeScheduleCommunicationMail(changeHistoryEntity *ChangeHistoryEntity, email string) error {

	checkExistKey := fmt.Sprintf("%s%d:%s", Map_mail_FeeScheduleCommunication, changeHistoryEntity.IncrId, email)
	a, _ := c.MapUsecase.GetForInt(checkExistKey)
	if a != 1 {
		tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": changeHistoryEntity.IncrId})
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("tClientCase  is nil:" + strconv.FormatInt(int64(changeHistoryEntity.IncrId), 10))
		}

		if HasEnabledPrimaryCase(tClientCase.CustomFields.TextValueByNameBasic("client_gid")) {
			isPrimaryCaseCalc, _, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				return err
			}
			if !isPrimaryCaseCalc {
				return nil
			}
		}

		fieldData := tClientCase.CustomFields

		clientGid := fieldData.TextValueByNameBasic("client_gid")
		if clientGid == "" {
			return nil
		}
		tContact, err := c.TUsecase.Data(Kind_clients, Eq{"gid": clientGid})
		if err != nil {
			return err
		}
		if tContact == nil {
			return errors.New("tContact is nil.")
		}
		//tContactFieldData := tContact.CustomFields
		//email := tContactFieldData.TextValueByNameBasic("email")
		//if email == "" {
		//	return errors.New("email is empty.")
		//}

		// 验证是否满足条件
		if fieldData.NumberValueByName("current_rating") == nil ||
			fieldData.NumberValueByName("effective_current_rating") == nil ||
			fieldData.TextValueByNameBasic("user_gid") == "" ||
			fieldData.TextValueByNameBasic("stages") != config_vbc.Stages_FeeScheduleandContract {
			return nil
		}

		subId, err := c.FeeUsecase.FeeScheduleCommunicationSubId(tClientCase)
		if err != nil {
			return err
		}

		typeMap := make(lib.TypeMap)
		typeMap.Set("Genre", MailGenre_FeeScheduleCommunication)
		typeMap.Set("SubId", subId)
		typeMap.Set("Email", email)

		err = c.TaskCreateUsecase.CreateTask(changeHistoryEntity.IncrId, typeMap, Task_Dag_BuzEmail, 0, "", "")
		if err != nil {
			c.log.Error(err)
		} else {
			err = c.MapUsecase.SetInt(checkExistKey, 1)
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ChangeHistoryUseacse) HandleCreateFolderInBoxAndMail(changeHistoryEntity *ChangeHistoryEntity) error {

	checkExistKey := fmt.Sprintf("%s%d", Map_HandleCreateFolderInBoxAndMail, changeHistoryEntity.IncrId)
	a, err := c.MapUsecase.GetForString(checkExistKey)
	if err != nil {
		return err
	}
	if a == "" {
		tClientCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"id": changeHistoryEntity.IncrId})
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("tClientCase is nil.")
		}
		customFields := tClientCase.CustomFields
		stages := tClientCase.CustomFields.TextValueByNameBasic("stages")
		//contractSource := tClientCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource)

		if stages != config_vbc.Stages_GettingStartedEmail && stages != config_vbc.Stages_AmAwaitingClientRecords {
			return nil
		}

		if HasEnabledPrimaryCase(tClientCase.CustomFields.TextValueByNameBasic("client_gid")) {
			usePrimaryCaseCalc, _, err := c.FeeUsecase.UsePrimaryCaseCalc(tClientCase)
			if err != nil {
				return err
			}
			if !usePrimaryCaseCalc {
				return nil
			}
		}

		clientGid := customFields.TextValueByNameBasic("client_gid")
		if clientGid == "" {
			return nil
		}
		tContact, err := c.TUsecase.Data(Kind_clients, Eq{"gid": clientGid})
		if err != nil {
			return err
		}
		if tContact == nil {
			return errors.New("tContact is nil.")
		}
		tContactFieldData := tContact.CustomFields
		email := tContactFieldData.TextValueByNameBasic("email")

		if email == "" {
			return errors.New("email is empty.")
		}

		typeMap := make(lib.TypeMap)
		typeMap["ClientId"] = changeHistoryEntity.IncrId
		err = c.TaskCreateUsecase.CreateTask(changeHistoryEntity.ID, typeMap, Task_Dag_BoxCreateFolderForNewClient, 0, "", "")
		if err != nil {
			return err
		}

		if stages == config_vbc.Stages_GettingStartedEmail {
			err = c.TriggerGettingStartedEmail(*tClientCase)
			if err != nil {
				return err
			}
			err = c.MapUsecase.Set(checkExistKey, "1")
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	return nil
}

func (c *ChangeHistoryUseacse) TriggerGettingStartedEmail(tCase TData) error {

	typeMap1 := make(lib.TypeMap)
	if IsAmContract(tCase) {
		return nil
		typeMap1.Set("Genre", MailGenre_AmGettingStartedEmail)
	} else {
		typeMap1.Set("Genre", MailGenre_GettingStartedEmail)
	}
	err := c.TaskCreateUsecase.CreateTask(tCase.Id(), typeMap1, Task_Dag_BuzEmail, 0, "", "")
	if err != nil {
		return err
	}
	return nil
}

func IsAmContract(tCase TData) bool {
	if tCase.CustomFields.TextValueByNameBasic(FieldName_ContractSource) == ContractSource_AM {
		return true
	}
	return false
}
