package tests

import (
	"testing"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_LeadVSChangeUsecase_HandleLeadVSChangeForClaimAnalysisToScheduleCall(t *testing.T) {

	a, _ := UT.ChangeHisUsecase.GetByCond(builder.Eq{"id": 209422})
	err := UT.LeadVSChangeUsecase.HandleLeadVSChangeForClaimAnalysisToScheduleCall(*a)
	lib.DPrintln(err)
}

func Test_LeadVSChangeUsecase_HandleLeadVSSyncClient(t *testing.T) {

	err := UT.LeadVSChangeUsecase.HandleLeadVSSyncClient(5112)
	lib.DPrintln(err)
}
