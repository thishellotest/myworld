package tests

import (
	"os"
	"testing"
	"vbc/lib"
)

func Test_PdfUsecase_PdfInfo(t *testing.T) {

	f, _ := os.Open("./tmp/0.pdf")
	aa, err := UT.PdfUsecase.PdfInfo(f, "0.pdf")
	lib.DPrintln(err)
	lib.DPrintln(aa.Dimensions)
}
