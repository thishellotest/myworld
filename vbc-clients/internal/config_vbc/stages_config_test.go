package config_vbc

import (
	"testing"
	"vbc/lib"
)

func Test_GetUnderStages(t *testing.T) {
	r := GetUnderStages(Stages_MedicalTeamCallVet)
	lib.DPrintln(r)
}

func Test_JudgeTaskNeedCompleteBySubject(t *testing.T) {
	f := JudgeTaskNeedCompleteBySubject("Current Treatment Follow-up", Stages_StatementsFinalized)
	lib.DPrintln(f)
	f = JudgeTaskNeedCompleteBySubject("Welcome Email Follow-up", Stages_StatementsFinalized)
	lib.DPrintln(f)
}

//func Test_StagesToNumber(t *testing.T) {
//	c, err := StagesToNumber(Stages_AwaitingClientRecords)
//	fmt.Println(c, err)
//}
