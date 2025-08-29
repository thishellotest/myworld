package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AutomaticUpdateDueDateUsecase_UpdateByGid(t *testing.T) {
	err := UT.AutomaticUpdateDueDateUsecase.UpdateByGid("444edb70887a4d65bfa3f79cd2189603", "")
	lib.DPrintln(err)
}
