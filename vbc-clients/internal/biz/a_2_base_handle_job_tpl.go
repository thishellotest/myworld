package biz

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
)

type TTemplateJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TTemplateEntity]
	BaseHandle[TTemplateEntity]
}

func NewTTemplateJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *TTemplateJobUsecase {
	uc := &TTemplateJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandle.Log = log.NewHelper(logger)
	uc.BaseHandle.TableName = TTemplateEntity{}.TableName()
	uc.BaseHandle.DB = CommonUsecase.DB()
	uc.BaseHandle.Handle = uc.Handle

	return uc
}

func (c *TTemplateJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	return c.CommonUsecase.DB().
		Table(TTemplateEntity{}.TableName()).
		Where("handle_status=? and next_at<=? and deleted_at=0",
			HandleStatus_waiting, time.Now().Unix()).Rows()

}

func (c *TTemplateJobUsecase) Handle(ctx context.Context, task *TTemplateEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.AppendHandleResultDetail(err.Error())
	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Save(task).Error
}

func (c *TTemplateJobUsecase) HandleExec(ctx context.Context, task *TTemplateEntity) error {

	return nil
}
