package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientCasebuzUsecase_CreateACaseByClientGid(t *testing.T) {

	operUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)

	err := UT.ClientCasebuzUsecase.CreateACaseByClientGid("6024b6aa842942489e3939350bc782eb", operUser)
	lib.DPrintln(err)
}
