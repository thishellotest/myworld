package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
)

type AccessControlWorkPayloadTaskUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	BehaviorUsecase   *BehaviorUsecase
	TaskCreateUsecase *TaskCreateUsecase
}

func NewAccessControlWorkPayloadTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BehaviorUsecase *BehaviorUsecase,
	TaskCreateUsecase *TaskCreateUsecase) *AccessControlWorkPayloadTaskUsecase {
	return &AccessControlWorkPayloadTaskUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		BehaviorUsecase:   BehaviorUsecase,
		TaskCreateUsecase: TaskCreateUsecase,
	}
}

func (c *AccessControlWorkPayloadTaskUsecase) Handle(work *AccessControlWorkEntity, task *AccessControlWorkPayloadTask) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	if work == nil {
		return errors.New("work is nil.")
	}
	if work.WorkType == WorkType_remind_fee_contract_signing &&
		task.Type == AccessControlWorkPayloadTask_Type_remind_by_email {
		clientId := work.ClientId()
		if clientId <= 0 {
			return errors.New("clientId is wrong.")
		}
		// Check whether it had did.
		behavior, err := c.BehaviorUsecase.GetOne(clientId, BehaviorType_sign_fee_schedule_contract_first_remind)
		if err != nil {
			return err
		}
		if behavior != nil {
			return errors.New("The customer has been reminded, please do not repeat the reminder.")
		}
		err = c.TaskCreateUsecase.CreateTaskMail(clientId, MailGenre_SignFeeContractFirstRemind, 0, nil, 0, "", "")
		if err != nil {
			return err
		} else {
			return c.BehaviorUsecase.Add(clientId, BehaviorType_sign_fee_schedule_contract_first_remind, time.Now(), "")
		}
	} else {
		return errors.New(task.Type + " does not support.")
	}
}
