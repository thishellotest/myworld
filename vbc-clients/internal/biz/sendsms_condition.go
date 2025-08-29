package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type SendsmsConditionUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewSendsmsConditionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *SendsmsConditionUsecase {
	uc := &SendsmsConditionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func (c *SendsmsConditionUsecase) VerifyTextGettingStartedEmail(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	if tCase.CustomFields.TextValueByNameBasic(FieldName_stages) == config_vbc.Stages_GettingStartedEmail {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextAwaitingClientRecords(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_AwaitingClientRecords ||
		stages == config_vbc.Stages_AmAwaitingClientRecords {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextSTRRequestPending(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_STRRequestPending ||
		stages == config_vbc.Stages_AmSTRRequestPending {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextStatementsFinalized(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_StatementsFinalized ||
		stages == config_vbc.Stages_AmStatementsFinalized {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextCurrentTreatment(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_CurrentTreatment ||
		stages == config_vbc.Stages_AmCurrentTreatment {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextAwaitingDecision(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_AwaitingDecision ||
		stages == config_vbc.Stages_AmAwaitingDecision {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}

func (c *SendsmsConditionUsecase) VerifyTextAwaitingPayment(tCase *TData) (bool, error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	if stages == config_vbc.Stages_AwaitingPayment ||
		stages == config_vbc.Stages_AmAwaitingPayment {
		// 可能需要验证 client_tasks
		return true, nil
	}
	return false, nil
}
