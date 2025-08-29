package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

type ManualUsecase struct {
	log             *log.Helper
	CommonUsecase   *CommonUsecase
	conf            *conf.Data
	LogUsecase      *LogUsecase
	BehaviorUsecase *BehaviorUsecase
}

func NewManualUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	LogUsecase *LogUsecase,
	BehaviorUsecase *BehaviorUsecase) *ManualUsecase {
	uc := &ManualUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		LogUsecase:      LogUsecase,
		BehaviorUsecase: BehaviorUsecase,
	}

	return uc
}

func (c *ManualUsecase) SyncCreateInvoiceBehavior() error {

	sqlRows, err := c.LogUsecase.DB.Raw("select * from  log  where from_type='Xero_CreateInvoice'").Rows()
	if err != nil {
		return err
	}
	defer sqlRows.Close()

	for sqlRows.Next() {
		var logEntity LogEntity
		err = c.LogUsecase.DB.ScanRows(sqlRows, &logEntity)
		if err != nil {
			c.log.Error(err)
			continue
		}
		//lib.DPrintln(logEntity.ID)
		err := c.BehaviorUsecase.BehaviorForCreateInvoice(logEntity.FromId, time.Unix(logEntity.CreatedAt, 0), "")
		if err != nil {
			panic(err)
		}
	}

	return nil
}
