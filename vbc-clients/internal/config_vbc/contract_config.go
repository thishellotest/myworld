package config_vbc

import (
	"github.com/pkg/errors"
	"time"
	"vbc/internal/utils"
)

type ReminderType string

const GroupReminderAmIntakeForm = "GroupReminderAmIntakeForm"

const AmIntakeFormReminderFirst = ReminderType("AmIntakeFormReminderFirst")
const AmIntakeFormReminderSecond = ReminderType("AmIntakeFormReminderSecond")
const AmIntakeFormReminderThird = ReminderType("AmIntakeFormReminderThird")

func IsGroupReminderAmIntakeForm(reminderType ReminderType) bool {
	if reminderType == AmIntakeFormReminderFirst ||
		reminderType == AmIntakeFormReminderSecond ||
		reminderType == AmIntakeFormReminderThird {
		return true
	}
	return false
}

type ReminderVo struct {
	ReminderType        ReminderType
	EmailTplSubId       int32
	ReminderIntervalDay int    // 提醒与上一次执行任务的间隔天数
	DayTimeAt           string // 当天在第几天提醒 16:00 (下午4点)
}
type ReminderConfigs []ReminderVo

func GetReminderConfigs() ReminderConfigs {

	return []ReminderVo{
		{
			ReminderType:        AmIntakeFormReminderFirst,
			EmailTplSubId:       1,
			ReminderIntervalDay: 1,
			DayTimeAt:           "08:05",
		},
		{
			ReminderType:        AmIntakeFormReminderSecond,
			EmailTplSubId:       2,
			ReminderIntervalDay: 3,
			DayTimeAt:           "08:05",
		},
		{
			ReminderType:        AmIntakeFormReminderThird,
			EmailTplSubId:       3,
			ReminderIntervalDay: 6,
			DayTimeAt:           "08:05",
		},
	}
}

// ReminderTime 提醒日期
func (c *ReminderVo) ReminderTime(currentTime time.Time, timeLocation time.Location) (dest time.Time, err error) {
	return utils.CalIntervalDayTime(currentTime, c.ReminderIntervalDay, c.DayTimeAt, timeLocation)
}

func (c ReminderConfigs) GetVo(typ ReminderType) (ReminderVo, error) {
	for k, v := range c {
		if v.ReminderType == typ {
			return c[k], nil
		}
	}
	return ReminderVo{}, errors.New("ContractReminderType is error")
}

type ContractReminderType string

const ContractReminderFirst = ContractReminderType("ContractReminderFirst")
const ContractReminderSecond = ContractReminderType("ContractReminderSecond")
const ContractReminderThird = ContractReminderType("ContractReminderThird")
const ContractReminderFourth = ContractReminderType("ContractReminderFourth")

type ContractReminderVo struct {
	ContractReminder    ContractReminderType
	EmailTplSubId       int32
	ReminderIntervalDay int    // 提醒与上一次执行任务的间隔天数
	DayTimeAt           string // 当天在第几天提醒 16:00 (下午4点)
}

// ContractReminderTime 合同提醒日期
func (c *ContractReminderVo) ContractReminderTime(currentTime time.Time, timeLocation time.Location) (dest time.Time, err error) {
	return utils.CalIntervalDayTime(currentTime, c.ReminderIntervalDay, c.DayTimeAt, timeLocation)
}

type ContractReminderConfigs []ContractReminderVo

func (c ContractReminderConfigs) GetVo(typ ContractReminderType) (ContractReminderVo, error) {
	for k, v := range c {
		if v.ContractReminder == typ {
			return c[k], nil
		}
	}
	return ContractReminderVo{}, errors.New("ContractReminderType is error")
}

func GetContractReminderConfigs() ContractReminderConfigs {

	return []ContractReminderVo{
		{
			ContractReminder:    ContractReminderFirst,
			EmailTplSubId:       1,
			ReminderIntervalDay: 7,
			DayTimeAt:           "08:05",
		},
		{
			ContractReminder:    ContractReminderSecond,
			EmailTplSubId:       2,
			ReminderIntervalDay: 7,
			DayTimeAt:           "08:05",
		},
		{
			ContractReminder:    ContractReminderThird,
			EmailTplSubId:       3,
			ReminderIntervalDay: 7,
			DayTimeAt:           "08:05",
		},
		{
			ContractReminder:    ContractReminderFourth,
			EmailTplSubId:       4,
			ReminderIntervalDay: 7,
			DayTimeAt:           "08:05",
		},
	}
}
