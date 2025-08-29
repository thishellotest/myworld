package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

const (
	// 发送合同
	BehaviorType_fee_schedule_contract = "fee_schedule_contract"
	BehaviorType_sent_am_contract      = "sent_am_contract"

	// 发送Intake form
	BehaviorType_sent_intake_form              = "sent_intake_form"
	BehaviorType_sent_am_getting_started_email = "sent_am_getting_started_email"
	// 完成合同
	BehaviorType_complete_fee_schedule_contract = "complete_fee_schedule_contract"
	BehaviorType_complete_am_contract           = "complete_am_contract"

	BehaviorType_sent_medical_team_forms_contract     = "sent_medical_team_forms_contract"
	BehaviorType_complete_medical_team_forms_contract = "complete_medical_team_forms_contract"

	//BehaviorType_complete_release_of_information_contract = "complete_release_of_information_contract"
	//BehaviorType_complete_patient_payment_form_contract   = "complete_patient_payment_form_contract"

	// 完成Intake form
	BehaviorType_complete_intake_form    = "complete_intake_form"
	BehaviorType_complete_am_intake_form = "complete_am_intake_form"

	BehaviorType_sign_fee_schedule_contract_first_remind = "sign_fee_schedule_contract_first_remind"

	BehaviorType_amintakeform_reminder_first  = "amintakeform_reminder_first"
	BehaviorType_amintakeform_reminder_second = "amintakeform_reminder_second"
	BehaviorType_amintakeform_reminder_third  = "amintakeform_reminder_third"

	BehaviorType_contract_reminder_first  = "contract_reminder_first"
	BehaviorType_contract_reminder_second = "contract_reminder_second"
	BehaviorType_contract_reminder_third  = "contract_reminder_third"
	BehaviorType_contract_reminder_fourth = "contract_reminder_fourth"
	BehaviorType_contract_non_responsive  = "contract_non_responsive"

	BehaviorType_am_contract_reminder_first  = "am_contract_reminder_first"
	BehaviorType_am_contract_reminder_second = "am_contract_reminder_second"
	BehaviorType_am_contract_reminder_third  = "am_contract_reminder_third"
	BehaviorType_am_contract_reminder_fourth = "am_contract_reminder_fourth"
	BehaviorType_am_contract_non_responsive  = "am_contract_non_responsive"

	BehaviorType_contract_reminder_first_sms  = "contract_reminder_first_sms"
	BehaviorType_contract_reminder_second_sms = "contract_reminder_second_sms"
	BehaviorType_contract_reminder_third_sms  = "contract_reminder_third_sms"
	BehaviorType_contract_reminder_fourth_sms = "contract_reminder_fourth_sms"

	BehaviorType_am_contract_reminder_first_sms  = "am_contract_reminder_first_sms"
	BehaviorType_am_contract_reminder_second_sms = "am_contract_reminder_second_sms"
	BehaviorType_am_contract_reminder_third_sms  = "am_contract_reminder_third_sms"
	BehaviorType_am_contract_reminder_fourth_sms = "am_contract_reminder_fourth_sms"

	BehaviorType_prefix_sms = "sms:"

	BehaviorType_CreateInvoice   = "create_invoice"
	BehaviorType_AmCreateInvoice = "am_create_invoice"
)

type BehaviorEntity struct {
	ID           int32 `gorm:"primaryKey"`
	IncrId       int32
	BehaviorType string
	BehaviorAt   int64
	Notes        string
	CreatedAt    int64
}

func (BehaviorEntity) TableName() string {
	return "behaviors"
}

type BehaviorUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BehaviorEntity]
	ClientUsecase                               *ClientUsecase
	TUsecase                                    *TUsecase
	TaskFailureLogUsecase                       *TaskFailureLogUsecase
	DbqsUsecase                                 *DbqsUsecase
	ActionOnceHandleCopyMedicalTeamFormsUsecase *ActionOnceHandleCopyMedicalTeamFormsUsecase
	GlobalEventBus                              *GlobalEventBus
}

func NewBehaviorUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ClientUsecase *ClientUsecase,
	TUsecase *TUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	DbqsUsecase *DbqsUsecase,
	ActionOnceHandleCopyMedicalTeamFormsUsecase *ActionOnceHandleCopyMedicalTeamFormsUsecase,
	GlobalEventBus *GlobalEventBus) *BehaviorUsecase {
	uc := &BehaviorUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		ClientUsecase:         ClientUsecase,
		TUsecase:              TUsecase,
		TaskFailureLogUsecase: TaskFailureLogUsecase,
		DbqsUsecase:           DbqsUsecase,
		ActionOnceHandleCopyMedicalTeamFormsUsecase: ActionOnceHandleCopyMedicalTeamFormsUsecase,
		GlobalEventBus: GlobalEventBus,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *BehaviorUsecase) GetOne(incrId int32, behaviorType string) (*BehaviorEntity, error) {
	return c.GetByCond(Eq{"incr_id": incrId, "behavior_type": behaviorType})
}

func (c *BehaviorUsecase) MedicalTeamFormsContractSentAt(incrId int32) *time.Time {
	beh, err := c.GetOne(incrId, BehaviorType_sent_medical_team_forms_contract)
	if err != nil {
		return nil
	}
	if beh == nil {
		return nil
	}
	ti := time.Unix(beh.BehaviorAt, 0)
	return to.Ptr(ti)
}

func (c *BehaviorUsecase) Add(incrId int32, behaviorType string, behaviorAt time.Time, notes string) error {

	e := &BehaviorEntity{
		IncrId:       incrId,
		BehaviorType: behaviorType,
		BehaviorAt:   behaviorAt.Unix(),
		Notes:        notes,
	}
	err := c.CommonUsecase.DB().Save(&e).Error
	if err == nil {
		// 此处不卡流程
		er := c.HandleCompleteBoxSign(behaviorType, incrId)
		if er != nil {
			c.TaskFailureLogUsecase.Add(TaskType_CompleteBoxSignBehavior, 0,
				map[string]interface{}{
					"incrId":       incrId,
					"behaviorType": behaviorType,
					"err":          er.Error(),
				})
		}
	}
	return err
}

func (c *BehaviorUsecase) HandleCompleteBoxSign(behaviorType string, clientCaseId int32) error {

	if behaviorType == BehaviorType_complete_fee_schedule_contract {
		tClientCase, er := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if er != nil {
			return er
		} else {
			if tClientCase == nil {
				return errors.New("tClientCase is nil.")
			} else {
				gid := tClientCase.CustomFields.TextValueByNameBasic("gid")
				er = c.ClientUsecase.HandleChangeStagesToGettingStartedEmail(gid)

				c.GlobalEventBus.Bus.Publish(GlobalEventBus_AfterHandleCompleteBoxSign, gid)

				if er != nil {
					return er
				}
			}
		}
	} else if behaviorType == BehaviorType_complete_am_contract {
		tClientCase, er := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if er != nil {
			return er
		} else {
			if tClientCase == nil {
				return errors.New("tClientCase is nil.")
			} else {
				gid := tClientCase.CustomFields.TextValueByNameBasic("gid")
				er = c.ClientUsecase.HandleChangeStagesForAm(gid)

				c.GlobalEventBus.Bus.Publish(GlobalEventBus_AfterHandleCompleteAmContract, gid)

				if er != nil {
					return er
				}
			}
		}
	} else if behaviorType == BehaviorType_complete_medical_team_forms_contract {
		medicalTeamFormsBeh, err := c.GetByCond(Eq{"behavior_type": BehaviorType_complete_medical_team_forms_contract, "incr_id": clientCaseId})
		if err != nil {
			return err
		}
		if medicalTeamFormsBeh == nil {
			return errors.New("medicalTeamFormsBeh is nil")
		}
		err = c.ActionOnceHandleCopyMedicalTeamFormsUsecase.HandleCopyMedicalTeamForms(clientCaseId)
		if err != nil {
			return err
		}

		tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tClientCase == nil {
			return errors.New("HandleCompleteBoxSign: tClientCase is nil")
		}
		gid := tClientCase.CustomFields.TextValueByNameBasic("gid")
		err = c.ClientUsecase.HandleChangeStagesToMiniDBQsFinalized(gid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BehaviorUsecase) BehaviorForCreateInvoice(incrId int32, behaviorAt time.Time, notes string) error {
	return c.Upsert(incrId, BehaviorType_CreateInvoice, behaviorAt, notes)
}

func (c *BehaviorUsecase) BehaviorForAmCreateInvoice(incrId int32, behaviorAt time.Time, notes string) error {
	return c.Upsert(incrId, BehaviorType_AmCreateInvoice, behaviorAt, notes)
}

func (c *BehaviorUsecase) Upsert(incrId int32, behaviorType string, behaviorAt time.Time, notes string) error {
	e, err := c.GetByCond(Eq{"behavior_type": behaviorType, "incr_id": incrId})
	if err != nil {
		return err
	}
	if e == nil {
		return c.Add(incrId, behaviorType, behaviorAt, notes)
	}
	return nil
}
