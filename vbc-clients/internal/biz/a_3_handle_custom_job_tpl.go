package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
)

/*
alter table task_failure_log

	add `handle_status` tinyint NOT NULL DEFAULT '0',
	add  `handle_result` tinyint NOT NULL DEFAULT '0',
	add  `handle_result_detail` text,
	add   `next_at` int NOT NULL DEFAULT '0' COMMENT '下次执行时间',
	add  `timeout` int NOT NULL DEFAULT '10800' COMMENT '超时时间 秒';

ALTER TABLE `task_failure_log` ADD INDEX `idx_hs` (`handle_status`);

HandleStatus       int

	HandleResult       int
	NextAt             int64
	Timeout            int32
	HandleResultDetail string
*/
type TTemplateCustomJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TTemplateEntity]
	BaseHandleCustom[TTemplateEntity]
}

func NewTTemplateCustomJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *TTemplateCustomJobUsecase {
	uc := &TTemplateCustomJobUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *TTemplateCustomJobUsecase) Handle(ctx context.Context, task *TTemplateEntity) error {

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

func (c *TTemplateCustomJobUsecase) HandleExec(ctx context.Context, task *TTemplateEntity) error {

	return nil
}
