package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
)

type TaskFailureLogJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	BaseHandleCustom[TaskLogEntity]
	MailUsecase *MailUsecase
}

func NewTaskFailureLogJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MailUsecase *MailUsecase) *TaskFailureLogJobUsecase {
	uc := &TaskFailureLogJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		MailUsecase:   MailUsecase,
	}
	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *TaskFailureLogJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	return c.CommonUsecase.DB().
		Table(TaskLogEntity{}.TableName()).
		Where("handle_status=? and next_at<=?",
			HandleStatus_waiting, time.Now().Unix()).Rows()

}

func (c *TaskFailureLogJobUsecase) Handle(ctx context.Context, task *TaskLogEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = err.Error()
	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Save(task).Error
}

func (c *TaskFailureLogJobUsecase) HandleExec(ctx context.Context, task *TaskLogEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	config := InitMailServiceConfig()

	str := fmt.Sprintf("ID: %d taskType: %s taskID: %d Notes:%s", task.ID, task.TaskType, task.TaskId, task.Notes)

	env := "Test"
	if configs.IsProd() {
		env = "Production"
	}

	if configs.StopNotifyCaseInDebug {
		return nil
	}

	return c.MailUsecase.SendEmail(config, &MailMessage{
		To:      "engineering@vetbenefitscenter.com",
		Subject: "[" + env + "] VBC Notification",
		Body: `<!DOCTYPE html>
<html lang="en">
<head>
    <title>VBC Notification</title>
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta content="width=device-width, initial-scale=1.0" name="viewport" />
</head>
<body>
<div>` + str + `</div>
</body>
</html>`,
	}, MailAttach_No, nil)
	return nil
}
