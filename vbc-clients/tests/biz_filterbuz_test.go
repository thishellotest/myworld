package tests

import (
	"testing"
	"vbc/lib"
)

func Test_FilterbuzUsecase_FilterDelete(t *testing.T) {
	err := UT.FilterbuzUsecase.FilterDelete("abc", []int32{1})
	lib.DPrintln(err)
}

func Test_FilterbuzUsecase_BizFilterList(t *testing.T) {
	r, err := UT.FilterbuzUsecase.BizFilterList("abc", "", "")
	lib.DPrintln(r)
	lib.DPrintln(err)
}

func Test_BizFilterSave_BizFilterSave(t *testing.T) {
	r, err := UT.FilterbuzUsecase.BizFilterSave("abc", "", "ssss", "ccc")
	lib.DPrintln(r)
	lib.DPrintln(err)
}
