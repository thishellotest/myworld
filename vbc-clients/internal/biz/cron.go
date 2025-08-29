package biz

import (
	"github.com/robfig/cron/v3"
)

type CronUsecase struct {
	cron *cron.Cron
}

func NewCronUsecase() *CronUsecase {
	cron := cron.New()
	cron.Start()
	return &CronUsecase{
		cron: cron,
	}
}

func (c *CronUsecase) Cron() *cron.Cron {
	return c.cron
}

func (c *CronUsecase) CleanAll() {
	entries := c.Cron().Entries()
	for _, v := range entries {
		c.Cron().Remove(v.ID)
	}
}
