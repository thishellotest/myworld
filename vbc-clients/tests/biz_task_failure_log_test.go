package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_TaskFailureLogUsecase_GetByCond(t *testing.T) {
	a, er := UT.TaskFailureLogUsecase.GetByCond(Eq{"id": 1})
	lib.DPrintln(a, er)
}

func Test_TaskFailureLogUsecase_Add(t *testing.T) {
	er := UT.TaskFailureLogUsecase.Add(biz.TaskType_ChangeStagesToGettingStartedEmail, 0, "sss")
	lib.DPrintln(er)
}
