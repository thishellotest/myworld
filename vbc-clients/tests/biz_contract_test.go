package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ClientCaseContractBasicDataUsecase(t *testing.T) {
	er := UT.ClientCaseContractBasicDataUsecase.BizHttpHandleHistory("1")
	lib.DPrintln(er)
}
