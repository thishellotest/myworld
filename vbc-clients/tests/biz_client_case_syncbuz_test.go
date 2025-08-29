package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientCaseSyncbuzUsecase_ClientToCases(t *testing.T) {
	err := UT.ClientCaseSyncbuzUsecase.ClientToCases("gid03", biz.ClientCaseSyncVo{
		FieldName:  "dob",
		FieldValue: "2024-12-10",
	}, nil)
	lib.DPrintln(err)
}

func Test_ClientCaseSyncbuzUsecase_CaseToClient(t *testing.T) {
	err := UT.ClientCaseSyncbuzUsecase.CaseToClient(0, biz.ClientCaseSyncVo{
		FieldName:  "dob",
		FieldValue: "2024-12-20",
	}, nil)
	lib.DPrintln(err)
}

func Test_ClientCaseSyncbuzUsecase_UpdatePersonalStatementManagerUrl(t *testing.T) {
	err := UT.ClientCaseSyncbuzUsecase.UpdatePersonalStatementManagerUrl(5511)
	lib.DPrintln(err)
}
