package tests

import (
	"testing"
	"vbc/lib"
)

func Test_CounterbuzUsecase_ClientUploadFilesStat(t *testing.T) {
	err := UT.CounterbuzUsecase.ClientUploadFilesStat(1)
	lib.DPrintln("r:", err)
	r, err := UT.CounterbuzUsecase.ClientUploadFilesLimit(1, 3)
	lib.DPrintln(r, err)
}
