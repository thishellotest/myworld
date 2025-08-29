package lib

import (
	"context"
	"github.com/pkg/errors"
)

type JobManager struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	jobFunc    func(ctx context.Context)
}

func NewJobManager(jobFunc func(ctx context.Context)) *JobManager {

	return &JobManager{
		jobFunc: jobFunc,
	}
}

func (c *JobManager) Start(ctx context.Context) error {
	if c.ctx != nil {
		return errors.New("JobManager 任务没有关闭，禁止启动")
	}
	c.ctx, c.cancelFunc = context.WithCancel(ctx)
	c.jobFunc(c.ctx)
	return nil
}

func (c *JobManager) Cancel() {
	if c.ctx != nil {
		c.cancelFunc()
	}
	c.ctx = nil
}
