package tests

import (
	"testing"
	"vbc/lib"
)

func Test_MailFeeContentUsecase_GetCurrentEvaluation(t *testing.T) {
	a, err := UT.MailFeeContentUsecase.GetCurrentEvaluation(100)
	lib.DPrintln(a, err)
}

func Test_MailFeeContentUsecase_GenContent(t *testing.T) {
	a, err := UT.MailFeeContentUsecase.GenContent(90)
	lib.DPrintln(a, err)
}
