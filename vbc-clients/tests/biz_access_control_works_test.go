package tests

import (
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	//. "vbc/lib/builder"
)

func Test_AccessControlWorkUsecase_CreateAccessControlWork(t *testing.T) {

	a := biz.SpawnRemindFeeContractSigningByEmail(biz.RemindFeeContractSigningParams{})
	payload := biz.AccessControlWorkPayload{}
	payload.Tasks = append(payload.Tasks, a)
	lib.DPrintln(payload)
	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	r, er := UT.AccessControlWorkUsecase.CreateAccessControlWork(biz.WorkType_remind_fee_contract_signing,
		"12",
		payload,
		expiredAt)
	lib.DPrintln(r, er)
}

func Test_SpawnRemindFeeContractSigning(t *testing.T) {
	a := biz.SpawnRemindFeeContractSigningByEmail(biz.RemindFeeContractSigningParams{})
	payload := biz.AccessControlWorkPayload{}
	payload.Tasks = append(payload.Tasks, a)
	lib.DPrintln(payload)
}
