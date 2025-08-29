package biz

import (
	"errors"
	"time"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type TaskCreateUsecase struct {
	CommonUsecase *CommonUsecase
}

func NewTaskCreateUsecase(CommonUsecase *CommonUsecase) *TaskCreateUsecase {
	return &TaskCreateUsecase{
		CommonUsecase: CommonUsecase,
	}
}

// CreateTaskWithFrom CreateTaskWithFrom
func (c *TaskCreateUsecase) CreateTaskWithFrom(incrId int32, taskInput interface{}, event string, nextAt int64, fromType string, fromId string) error {
	a := &TaskEntity{
		IncrId:    incrId,
		TaskInput: InterfaceToString(taskInput),
		Event:     event,
		NextAt:    nextAt,
		FromType:  fromType,
		FromId:    fromId,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Save(&a).Error
}

func (c *TaskCreateUsecase) CreateGroupReminderAmIntakeFormTask(incrId int32, reminderType config_vbc.ReminderType, nextAt int64) error {
	return c.CreateTask(incrId, map[string]interface{}{
		"ReminderType": reminderType,
	}, Task_Dag_HandleReminder, nextAt, config_vbc.GroupReminderAmIntakeForm, InterfaceToString(incrId))
}

// CreateTask CreateTask
func (c *TaskCreateUsecase) CreateTask(incrId int32, taskInput interface{}, event string, nextAt int64, fromType string, fromId string) error {
	a := &TaskEntity{
		IncrId:    incrId,
		TaskInput: InterfaceToString(taskInput),
		Event:     event,
		NextAt:    nextAt,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		FromType:  fromType,
		FromId:    fromId,
	}
	return c.CommonUsecase.DB().Save(&a).Error
}

type MailTaskInput struct {
	Genre         string       `json:"Genre"`
	SubId         int32        `json:"SubId"`
	Email         string       `json:"Email"`
	DynamicParams lib.TypeMap  `json:"DynamicParams"`
	MailMessage   *MailMessage `json:"MailMessage"`
}

func (c *TaskCreateUsecase) CreateTaskMail(incrId int32, genre string, subId int32, dynamicParams lib.TypeMap, nextAt int64, fromType string, fromId string) error {
	mailTaskInput := &MailTaskInput{
		Genre:         genre,
		SubId:         subId,
		DynamicParams: dynamicParams,
	}
	return c.CreateTask(incrId, mailTaskInput, Task_Dag_BuzEmail, nextAt, fromType, fromId)
}

func (c *TaskCreateUsecase) CreateCustomTaskMail(incrId int32, MailMessage *MailMessage, nextAt int64) error {
	mailTaskInput := &MailTaskInput{
		Genre:       MailGenre_Custom,
		MailMessage: MailMessage,
	}
	return c.CreateTask(incrId, mailTaskInput, Task_Dag_BuzEmail, nextAt, "", "")
}

type CreateBoxCreateFolderForNewClientTask struct {
	ClientId int32 `json:"ClientId"`
}

func (c *TaskCreateUsecase) CreateBoxCreateFolderForNewClientTask(
	task *CreateBoxCreateFolderForNewClientTask) (err error) {
	if task == nil {
		return errors.New("CreateBoxCreateFolderForNewClientTask:task is nil")
	}
	return c.CreateTask(task.ClientId,
		task, Task_Dag_BoxCreateFolderForNewClient, 0, "", "")
}
