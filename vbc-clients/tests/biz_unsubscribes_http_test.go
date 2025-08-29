package tests

import (
	"testing"
	"vbc/lib"
)

func Test_UnsubscribesHttpUsecase_BizList(t *testing.T) {
	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	data, err := UT.UnsubscribesHttpUsecase.BizList(*userFacade, 1, 20)
	lib.DPrintln(data)
	lib.DPrintln(err)
}
