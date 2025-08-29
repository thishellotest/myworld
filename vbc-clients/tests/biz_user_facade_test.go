package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_UserFacade_GetTimezonesEntity(t *testing.T) {
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	userFacade := biz.UserFacade{
		TData: *tUser,
	}
	a, err := userFacade.GetTimezonesEntity(UT.TimezonesUsecase)
	lib.DPrintln(a, err)
	a, err = userFacade.GetTimezonesEntity(UT.TimezonesUsecase)
	lib.DPrintln(a, err)
}
