package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
)

type ClientTaskJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TTemplateEntity]
	BaseHandleCustom[TTemplateEntity]
}

func NewClientTaskJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ClientTaskJobUsecase {
	uc := &ClientTaskJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *ClientTaskJobUsecase) Handle(ctx context.Context, task *TTemplateEntity) error {

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

func (c *ClientTaskJobUsecase) HandleExec(ctx context.Context, task *TTemplateEntity) error {

	return nil
}
