package tests

import (
	"testing"
	"vbc/internal/biz"
)

func Test_SettingHttpUsecase_BizCustomView(t *testing.T) {
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	var settingHttpCustomViewRequest biz.SettingHttpCustomViewRequest
	UT.SettingHttpUsecase.BizCustomView("", biz.UserFacade{
		TData: *tUser,
	}, settingHttpCustomViewRequest)
}
