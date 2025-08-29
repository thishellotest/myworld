package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZohobuzTaskUsecase_SyncTasksDeletes(t *testing.T) {
	err := UT.ZohobuzTaskUsecase.SyncTasksDeletes()
	lib.DPrintln(err)
}
