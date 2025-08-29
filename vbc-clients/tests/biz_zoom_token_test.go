package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZoomTokenUsecase_GetAccessToken(t *testing.T) {
	info, err := UT.ZoomTokenUsecase.GetAccessToken()
	lib.DPrintln(info, err)
}
