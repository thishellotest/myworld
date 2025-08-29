package tests

import (
	"testing"
	"vbc/lib"
)

func Test_MetadataUsecase_BizHttpConditions(t *testing.T) {
	r, err := UT.MetadataUsecase.BizHttpConditions()
	lib.DPrintln(r, err)
}

func Test_MetadataUsecase_BizBasicdata(t *testing.T) {
	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	r, err := UT.MetadataUsecase.BizBasicdata(*userFacade)
	lib.DPrintln(r, err)
}
