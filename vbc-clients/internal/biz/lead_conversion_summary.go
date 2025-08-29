package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type LeadConversionSummaryEntity struct {
	ID           int32 `gorm:"primaryKey"`
	CaseId       int32
	CreatedTime  string
	ContractTime string
	SubmitVaTime string
	InvoiceTime  string
	CreatedAt    int64
	UpdatedAt    int64
}

func (LeadConversionSummaryEntity) TableName() string {
	return "lead_conversion_summary"
}

type LeadConversionSummaryUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[LeadConversionSummaryEntity]
}

func NewLeadConversionSummaryUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *LeadConversionSummaryUsecase {
	uc := &LeadConversionSummaryUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *LeadConversionSummaryUsecase) Upsert(CaseId int32, CreatedTime string, ContractTime string, SubmitVaTime string, InvoiceTime string) error {

	entity, err := c.GetByCond(Eq{"case_id": CaseId})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &LeadConversionSummaryEntity{
			CaseId:       CaseId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			CreatedTime:  CreatedTime,
			ContractTime: ContractTime,
			SubmitVaTime: SubmitVaTime,
			InvoiceTime:  InvoiceTime,
		}
		return c.CommonUsecase.DB().Save(entity).Error
	} else {
		entity.CreatedTime = CreatedTime
		entity.ContractTime = ContractTime
		entity.SubmitVaTime = SubmitVaTime
		entity.InvoiceTime = InvoiceTime
		entity.UpdatedAt = time.Now().Unix()
		return c.CommonUsecase.DB().Save(entity).Error
	}
	return nil
}
