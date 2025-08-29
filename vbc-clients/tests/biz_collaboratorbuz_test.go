package tests

import (
	"testing"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_CollaboratorbuzUsecase_DoCollaboratorByChangeHistory(t *testing.T) {
	a, _ := UT.ChangeHisUsecase.GetByCond(builder.Eq{"id": 8777})
	err := UT.CollaboratorbuzUsecase.DoCollaboratorByChangeHistory(*a)
	lib.DPrintln(err)
}
