package tests

import (
	"testing"
	"vbc/lib"
)

func Test_StageTransUsecase_ZohoStageToDBStage(t *testing.T) {
	//vbc_config.Stages_GettingStartedEmail
	a, err := UT.StageTransUsecase.ZohoStageToDBStage("2. Getting Started Email")
	lib.DPrintln(a, err)
}

func Test_StageTransUsecase_DBStageToZohoStage(t *testing.T) {
	//vbc_config.Stages_GettingStartedEmail
	a, err := UT.StageTransUsecase.DBStageToZohoStage("Getting Started Email")
	lib.DPrintln(a, err)
}
