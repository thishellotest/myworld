package biz

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_GenJotformNewFileNameForAI(t *testing.T) {

	a, _ := UT.JotformSubmissionUsecase.GetLatestFormInfo("6208935729788488694")
	//lib.DPrintln(a)
	b, err := biz.GenJotformNewFileNameForAI(a)
	lib.DPrintln(b, err)
}
