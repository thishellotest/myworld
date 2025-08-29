package tests

import (
	"testing"
	"vbc/lib"
)

func Test_BoxcontractUsecase_a(t *testing.T) {
	folderId, err := UT.BoxcontractUsecase.ContractFolderId(33)
	lib.DPrintln(folderId, err)
}
