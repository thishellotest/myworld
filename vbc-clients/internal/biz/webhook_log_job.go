package biz

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type WebhookLogJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	BaseHandleCustom[WebhookLogEntity]
	WebsiteUsecase         *WebsiteUsecase
	TaskFailureLogUsecase  *TaskFailureLogUsecase
	JotformbuzUsecase      *JotformbuzUsecase
	UnsubscribesbuzUsecase *UnsubscribesbuzUsecase
}

func NewWebhookLogJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	WebsiteUsecase *WebsiteUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	JotformbuzUsecase *JotformbuzUsecase,
	UnsubscribesbuzUsecase *UnsubscribesbuzUsecase) *WebhookLogJobUsecase {
	uc := &WebhookLogJobUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		WebsiteUsecase:         WebsiteUsecase,
		TaskFailureLogUsecase:  TaskFailureLogUsecase,
		JotformbuzUsecase:      JotformbuzUsecase,
		UnsubscribesbuzUsecase: UnsubscribesbuzUsecase,
	}

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}
func (c *WebhookLogJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	return c.CommonUsecase.DB().
		Table(WebhookLogEntity{}.TableName()).
		Where("handle_status=?  and deleted_at=0",
			HandleStatus_waiting).Rows()

}

func (c *WebhookLogJobUsecase) Handle(ctx context.Context, task *WebhookLogEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	isDone, err := c.HandleExec(ctx, task)

	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = err.Error()
		return c.CommonUsecase.DB().Save(task).Error
	} else if isDone {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_ok
		return c.CommonUsecase.DB().Save(task).Error
	} else {
		task.HandleStatus = HandleStatus_done
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = "不重试"
		return c.CommonUsecase.DB().Save(task).Error
	}
}

func (c *WebhookLogJobUsecase) HandleExec(ctx context.Context, task *WebhookLogEntity) (isDone bool, err error) {

	if task == nil {
		return false, errors.New("task is nil")
	}
	if task.From == WebhookLog_From_website {
		isDone, err = c.HandleExecWebsite(ctx, task)
		if err != nil { // 加入报警
			c.log.Error(err, InterfaceToString(task))
			c.TaskFailureLogUsecase.Add(TaskType_HandleExecWebsite, 0, map[string]interface{}{
				"WebhookLogId": task.ID,
				"err":          err.Error(),
			})
		}
		return
	} else if task.From == WebhookLog_From_jotform {

		isDone, err = c.HandleExecJotform(ctx, task)
		if err != nil { // 加入报警
			c.log.Error(err, InterfaceToString(task))
			c.TaskFailureLogUsecase.Add("HandleExecJotform", 0, map[string]interface{}{
				"WebhookLogId": task.ID,
				"err":          err.Error(),
			})
		}
		return
	} else if task.From == WebhookLog_From_dialpad {
		isDone, err = c.HandleExecDialpad(ctx, task)
		if err != nil { // 加入报警
			c.log.Error(err, InterfaceToString(task))
			c.TaskFailureLogUsecase.Add("HandleExecDialpad", 0, map[string]interface{}{
				"WebhookLogId": task.ID,
				"err":          err.Error(),
			})
		}
		return
	} else {
		return true, nil
	}
}

func (c *WebhookLogJobUsecase) HandleExecDialpad(ctx context.Context, task *WebhookLogEntity) (isDone bool, err error) {
	if task == nil {
		return false, errors.New("task is nil")
	}

	cliams, err := JWTParse(task.Body, config_vbc.DialpadWebhookSecret())
	if err != nil {
		return false, err
	}

	// 此处方便修改短信内容测试
	if task.NeatBody == "" {
		task.NeatBody = InterfaceToString(cliams)
	}

	task.UpdatedAt = time.Now().Unix()
	err = c.CommonUsecase.DB().Save(task).Error
	if err != nil {
		return false, err
	}

	err = c.UnsubscribesbuzUsecase.HandleFromDialpadWebhookEvent(lib.ToTypeMapByString(task.NeatBody))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *WebhookLogJobUsecase) HandleExecWebsite(ctx context.Context, task *WebhookLogEntity) (isDone bool, err error) {

	if task == nil {
		return false, errors.New("task is nil")
	}
	err = c.WebsiteUsecase.SyncToZohoOrVBCRM(task.Body)
	return true, err
}

func (c *WebhookLogJobUsecase) HandleExecJotform(ctx context.Context, task *WebhookLogEntity) (isDone bool, err error) {

	if task == nil {
		return false, errors.New("task is nil")
	}
	neatBody := lib.ToTypeMapByString(task.NeatBody)
	formID := neatBody.GetString("formID")
	submissionID := neatBody.GetString("submissionID")
	//if formID != JotformIntakeFormID {
	//	return false, errors.New("HandleExecJotform: formID is wrong")
	//}
	if submissionID == "" {
		return false, errors.New("HandleExecJotform: submissionID is wrong")
	}
	err = c.JotformbuzUsecase.HandleSubmission(submissionID, "", formID)
	return true, err
}
