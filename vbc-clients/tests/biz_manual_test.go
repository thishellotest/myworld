package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ManualUsecase_SyncCreateInvoiceBehavior(t *testing.T) {
	err := UT.ManualUsecase.SyncCreateInvoiceBehavior()
	lib.DPrintln(err)
}
