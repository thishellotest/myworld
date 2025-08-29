package tests

import (
	"testing"
	"vbc/lib"
)

func Test_NotificationHttpUsecase_BizList(t *testing.T) {
	a, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	r, err := UT.NotificationHttpUsecase.BizList(*a, 0, 0)
	lib.DPrintln(r, err)

}

//func Test_aaaaa(t *testing.T) {
//	biz.GetFacadesByGids[*biz.UserFacade](UT.TUsecase, biz.Kind_users, []string{"6159272000000453669"})
//}
