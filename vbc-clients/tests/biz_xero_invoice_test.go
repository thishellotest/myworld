package tests

import (
	"testing"
	"vbc/lib"
)

func Test_XeroInvoiceUsecase_HandleInvoice(t *testing.T) {
	err := UT.XeroInvoiceUsecase.HandleInvoice(153)
	lib.DPrintln(err)
}
