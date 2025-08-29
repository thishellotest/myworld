package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/uuid"
)

func Test_aaa22(t *testing.T) {
	aa := uuid.UuidWithoutStrike()
	lib.DPrintln(aa)
}

func Test_MetadataHttpUsecase_BizFields(t *testing.T) {
	userFacade := biz.UserFacade{}
	a, err := UT.MetadataHttpUsecase.BizFields("", userFacade)
	lib.DPrintln(a)
	lib.DPrintln(err)
}

func Test_MetadataHttpUsecase_BizOptions(t *testing.T) {
	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	aa, err := UT.MetadataHttpUsecase.BizOptions(*userFacade, biz.Kind_users, "role_gid", nil, nil)
	lib.DPrintln(aa, err)
}
