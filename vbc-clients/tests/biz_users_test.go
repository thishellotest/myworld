package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_UserUsecase_GetByFullName(t *testing.T) {
	user, err := UT.UserUsecase.GetByFullName("Engineering Team")
	lib.DPrintln(user, err)
}

func Test_UserUsecase_GetByGid(t *testing.T) {
	user, err := UT.UserUsecase.GetByGid("6159272000000453669")
	lib.DPrintln(user, err)
}

func Test_UserUsecase_InitPassword(t *testing.T) {
	UT.UserUsecase.InitPassword("glliao@vetbenefitscenter.com")
}

func Test_UserUsecase_GetUserWithCache(t *testing.T) {
	//userCaches := lib.CacheInit[*biz.TData]()
	//
	//a, err := UT.UserUsecase.GetUserWithCache(userCaches, 2)
	//lib.DPrintln(a, err)
	//a, err = UT.UserUsecase.GetUserWithCache(userCaches, 2)
	//lib.DPrintln(a, err)
	//a, err = UT.UserUsecase.GetUserWithCache(userCaches, 10)
	//lib.DPrintln(a, err)
	//a, err = UT.UserUsecase.GetUserWithCache(userCaches, 10)
	//lib.DPrintln(a, err)
}

func Test_UserUsecase_ListByCond(t *testing.T) {
	a, err := UT.TUsecase.ListByCond(biz.Kind_users, nil)
	lib.DPrintln(a, err)
}

func Test_UserUsecase_Get(t *testing.T) {

	// 6159272000000453669 gary user Gid
	// 6159272000001027094 lili user Gid
	a, err := UT.TUsecase.DataByGid(biz.Kind_users, "6159272000001027094")
	//UT.FieldUsecase.ListByKind(biz.kin)
	lib.DPrintln(a, err)
	profile, err := a.RelaData(UT.BUsecase, biz.User_FieldName_profile_gid)
	lib.DPrintln(profile, err)
	c := profile.CustomFields.NumberValueByNameBasic(biz.Profile_FieldName_is_admin)

	lib.DPrintln(c)

	c = profile.CustomFields.NumberValueByNameBasic(biz.Profile_FieldName_is_admin)
	lib.DPrintln(profile, err)
	lib.DPrintln(c)
}

func Test_UserUsecase_VSTeamUsers(t *testing.T) {
	res, err := UT.UserUsecase.VSTeamUsers()
	lib.DPrintln(err)
	for _, v := range res {
		lib.DPrintln(v.Gid(), v.CustomFields.TextValueByNameBasic(biz.FieldName_email))
	}
}
