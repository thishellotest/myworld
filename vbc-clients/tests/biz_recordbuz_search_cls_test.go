package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RecordbuzSearchCls_hasPermissionRow(t *testing.T) {
	cls := UT.RecordbuzSearchUsecase.NewRecordbuzSearchCls(false)
	tUser, _ := UT.TUsecase.DataByGid(biz.Kind_users, "6159272000000453669")

	e, _ := UT.KindUsecase.GetByKind(biz.Kind_client_cases)
	r, err := cls.HasPermissionRow(biz.UserFacade{TData: *tUser}, *e, "b703679a066446d5910c4c43277248bc")
	lib.DPrintln(r, err)
}
