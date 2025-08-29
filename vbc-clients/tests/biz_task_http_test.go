package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_TaskHttpUsecase_BizList(t *testing.T) {
	tUser, _ := UT.TUsecase.DataByGid(biz.Kind_users, "6159272000000453669")
	var taskHttpListRequest biz.TaskHttpListRequest
	r, err := UT.TaskHttpUsecase.BizList(biz.UserFacade{TData: *tUser}, biz.Kind_client_cases,
		"b703679a066446d5910c4c43277248bc",
		taskHttpListRequest)
	lib.DPrintln(r, err)
}
