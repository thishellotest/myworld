package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type ContractReminderUsecase struct {
	log                   *log.Helper
	CommonUsecase         *CommonUsecase
	conf                  *conf.Data
	TUsecase              *TUsecase
	TaskCreateUsecase     *TaskCreateUsecase
	LogUsecase            *LogUsecase
	ClientEnvelopeUsecase *ClientEnvelopeUsecase
}

func NewContractReminderUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	TaskCreateUsecase *TaskCreateUsecase,
	LogUsecase *LogUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase) *ContractReminderUsecase {
	uc := &ContractReminderUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		TUsecase:              TUsecase,
		TaskCreateUsecase:     TaskCreateUsecase,
		LogUsecase:            LogUsecase,
		ClientEnvelopeUsecase: ClientEnvelopeUsecase,
	}

	return uc
}

func (c *ContractReminderUsecase) GetCaseForReminder(tCaseId int32) (*TData, error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, tCaseId)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	deleted, err := IsDeletedCase(tCase)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	if deleted {
		c.log.Debug("the case was deleted")
		return nil, nil
	}
	return tCase, nil
}

func (c *ContractReminderUsecase) AmIntakeFormReminderFirstReminder(tCaseId int32) error {

	tCase, err := c.GetCaseForReminder(tCaseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tCase == nil {
		return nil
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages != config_vbc.Stages_AmInformationIntake {
		c.LogUsecase.SaveLog(tCaseId, "AmIntakeFormReminderFirstReminder", map[string]interface{}{
			"stages": stages,
		})
		return nil
	}
	reminderConfig, err := config_vbc.GetReminderConfigs().GetVo(config_vbc.AmIntakeFormReminderFirst)
	if err != nil {
		return err
	}
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAtTime, err := reminderConfig.ReminderTime(time.Now(), timeLocation)
	if err != nil {
		return err
	}

	err = c.TaskCreateUsecase.CreateGroupReminderAmIntakeFormTask(tCaseId, config_vbc.AmIntakeFormReminderFirst, nextAtTime.Unix())

	if err != nil {
		return err
	}
	return nil
}

func (c *ContractReminderUsecase) FirstReminder(tCaseId int32) error {

	tCase, err := c.GetCaseForReminder(tCaseId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tCase == nil {
		return nil
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages != config_vbc.Stages_FeeScheduleandContract && stages != config_vbc.Stages_AmContractPending {
		c.LogUsecase.SaveLog(tCaseId, "FirstReminder", map[string]interface{}{
			"stages": stages,
		})
		return nil
	}
	contractReminderConfig, err := config_vbc.GetContractReminderConfigs().GetVo(config_vbc.ContractReminderFirst)
	if err != nil {
		return err
	}
	timeLocation := GetCaseTimeLocation(tCase, c.log)
	nextAtTime, err := contractReminderConfig.ContractReminderTime(time.Now(), timeLocation)
	if err != nil {
		return err
	}

	err = c.TaskCreateUsecase.CreateTask(tCaseId, map[string]interface{}{
		"ContractReminderType": config_vbc.ContractReminderFirst,
	}, Task_Dag_HandleContractReminder, nextAtTime.Unix(), "", "")

	//err = c.TaskCreateUsecase.CreateTaskMail(tCaseId,
	//	MailGenre_ContractReminder,
	//	contractReminderConfig.EmailTplSubId,
	//	nil, nextAtTime.Unix())
	if err != nil {
		return err
	}
	return nil
}
