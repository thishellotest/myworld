package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ManualThingstoknowUsecase_DestClientCases(t *testing.T) {
	cases, err := UT.ManualThingstoknowUsecase.DestClientCases()
	lib.DPrintln(cases, err)
}

func Test_ManualThingstoknowUsecase_ThingsToKnowFileByApi(t *testing.T) {
	// 283402271401
	// 283402276201 正确
	info, err := UT.ManualThingstoknowUsecase.ThingsToKnowFileByApi("283402276201")
	lib.DPrintln(info)
	lib.DPrintln(err)
}

func Test_ManualThingstoknowUsecase_HandleUploadNewThingsToKnowFileAllCases(t *testing.T) {
	err := UT.ManualThingstoknowUsecase.HandleUploadNewThingsToKnowFileAllCases()
	lib.DPrintln(err)
}

func Test_ManualThingstoknowUsecase_UploadNewThingsToKnowFile(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5369)
	lib.DPrintln(tCase, err)
	err = UT.ManualThingstoknowUsecase.UploadNewThingsToKnowFile(tCase)
	lib.DPrintln(err)
}
