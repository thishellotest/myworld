package tests

import (
	"testing"
	"vbc/lib"
)

func Test_UnsubscribesUsecase_CanSendSms(t *testing.T) {
	a, err := UT.UnsubscribesUsecase.CanSendSms("+14159389005")
	lib.DPrintln(a, err)

	a, err = UT.UnsubscribesUsecase.CanSendSms("(402) 215-6064")
	lib.DPrintln(a, err)
}
