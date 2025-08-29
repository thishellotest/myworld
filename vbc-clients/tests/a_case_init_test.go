package tests

import (
	"testing"
	"vbc/lib"
)

func Test_InitCase_InitClientCase(t *testing.T) {
	caseId := int32(5251)
	var err error
	err = UT.ActionOnceUsecase.InitClientCase(caseId)
	lib.DPrintln(err)

	//	err = UT.ZohobuzUsecase.HandleAmount(caseId)
	//	lib.DPrintln(err)
	//
	//	err = UT.ZohobuzUsecase.HandleClientCaseName(caseId)
	//	lib.DPrintln(err)
}
