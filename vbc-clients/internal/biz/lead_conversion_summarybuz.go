package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

type LeadConversionSummaryBuzUsecase struct {
	log                          *log.Helper
	conf                         *conf.Data
	CommonUsecase                *CommonUsecase
	LeadConversionSummaryUsecase *LeadConversionSummaryUsecase
	TUsecase                     *TUsecase
	ChangeHisUsecase             *ChangeHisUsecase
}

func NewLeadConversionSummaryBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	LeadConversionSummaryUsecase *LeadConversionSummaryUsecase,
	TUsecase *TUsecase,
	ChangeHisUsecase *ChangeHisUsecase,
) *LeadConversionSummaryBuzUsecase {
	uc := &LeadConversionSummaryBuzUsecase{
		log:                          log.NewHelper(logger),
		CommonUsecase:                CommonUsecase,
		conf:                         conf,
		LeadConversionSummaryUsecase: LeadConversionSummaryUsecase,
		TUsecase:                     TUsecase,
		ChangeHisUsecase:             ChangeHisUsecase,
	}

	return uc
}

func (c *LeadConversionSummaryBuzUsecase) ManualAll() error {

	res, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"biz_deleted_at": 0, "deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range res {
		er := c.DoOne(v.Id())
		if er != nil {
			c.log.Info("LeadConversionSummaryBuzUsecase DoOne", er)
		}
	}

	return nil
}

func (c *LeadConversionSummaryBuzUsecase) DoOne(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	stage := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)

	createdAt := tCase.CustomFields.NumberValueByNameBasic(FieldName_created_at)
	Created := time.Unix(int64(createdAt), 0).In(configs.VBCDefaultLocation)
	CreatedTime := Created.Format(time.DateOnly)

	ContractTime, err := c.GetContractTime(caseId)
	if err != nil {
		return err
	}
	if ContractTime != "" {
		if CreatedTime > ContractTime {
			ContractTime = CreatedTime
		}
	} else {
		if stage != config_vbc.Stages_IncomingRequest &&
			stage != config_vbc.Stages_FeeScheduleandContract &&
			stage != config_vbc.Stages_Dormant &&
			stage != config_vbc.Stages_Terminated {
			ContractTime = CreatedTime
		}
	}

	InvoiceTime, err := c.GetInvoiceTime(caseId)
	if err != nil {
		return err
	}
	if InvoiceTime != "" {
		if ContractTime == "" {
			ContractTime = CreatedTime
		}
		if ContractTime >= InvoiceTime {
			InvoiceTime = ContractTime
		}
	} else {
		if stage == config_vbc.Stages_AwaitingPayment ||
			stage == config_vbc.Stages_27_AwaitingBankReconciliation ||
			stage == config_vbc.Stages_Completed ||
			stage == config_vbc.Stages_AmAwaitingPayment ||
			stage == config_vbc.Stages_Am27_AwaitingBankReconciliation ||
			stage == config_vbc.Stages_AmCompleted {
			InvoiceTime = ContractTime
		}
	}
	submitVaTime, err := c.GetSubmitVaTime(caseId)
	if err != nil {
		if err != nil {
			return err
		}
	}
	if submitVaTime != "" && InvoiceTime != "" {
		if submitVaTime > InvoiceTime {
			submitVaTime = InvoiceTime
		}
	}

	return c.LeadConversionSummaryUsecase.Upsert(caseId, CreatedTime, ContractTime, submitVaTime, InvoiceTime)
}

func (c *LeadConversionSummaryBuzUsecase) GetContractTime(caseId int32) (string, error) {

	entity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(And(Eq{"kind": Kind_client_cases, "incr_id": caseId,
		"field_name": FieldName_stages},
		NotIn("new_value", config_vbc.Stages_IncomingRequest,
			config_vbc.Stages_FeeScheduleandContract,
			config_vbc.Stages_Terminated,
			config_vbc.Stages_Dormant)), "id")
	if err != nil {
		return "", err
	}
	if entity == nil {
		return "", nil
	}
	val := time.Unix(entity.CreatedAt, 0).In(configs.GetVBCDefaultLocation()).Format(time.DateOnly)
	return val, nil
}

func (c *LeadConversionSummaryBuzUsecase) GetSubmitVaTime(caseId int32) (string, error) {

	entity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(And(Eq{"kind": Kind_client_cases, "incr_id": caseId,
		"field_name": FieldName_stages},
		In("new_value", config_vbc.Stages_VerifyEvidenceReceived)), "id")
	if err != nil {
		return "", err
	}
	if entity == nil {
		return "", nil
	}
	val := time.Unix(entity.CreatedAt, 0).In(configs.GetVBCDefaultLocation()).Format(time.DateOnly)
	return val, nil
}

func (c *LeadConversionSummaryBuzUsecase) GetInvoiceTime(caseId int32) (string, error) {

	entity, err := c.ChangeHisUsecase.GetByCondWithOrderBy(And(Eq{"kind": Kind_client_cases, "incr_id": caseId,
		"field_name": FieldName_stages},
		In("new_value", config_vbc.Stages_AwaitingPayment,
			config_vbc.Stages_27_AwaitingBankReconciliation,
			config_vbc.Stages_Completed,
			config_vbc.Stages_AmAwaitingPayment,
			config_vbc.Stages_Am27_AwaitingBankReconciliation,
			config_vbc.Stages_AmCompleted,
		)), "id")
	if err != nil {
		return "", err
	}
	if entity == nil {
		return "", nil
	}
	val := time.Unix(entity.CreatedAt, 0).In(configs.GetVBCDefaultLocation()).Format(time.DateOnly)
	return val, nil
}
