package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ClientReviewUsecase(t *testing.T) {
	err := UT.ClientReviewUsecase.ImportExcel()
	lib.DPrintln(err)
}
