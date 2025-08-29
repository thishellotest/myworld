package tests

import (
	"testing"
	"vbc/lib"
)

func Test_WebpUsecase_JpgToWebp(t *testing.T) {
	jpgFile := UT.ResourceUsecase.ResPath() + "/123.jpg"
	outFile := "/tmp/aa.webp"
	width, height, err := UT.WebpUsecase.JpgToWebp(jpgFile, outFile)
	lib.DPrintln(width, height, err)
}

func Test_WebpUsecase_TestWebp(t *testing.T) {
	err := UT.WebpUsecase.TestWebp()
	lib.DPrintln(err)
}
