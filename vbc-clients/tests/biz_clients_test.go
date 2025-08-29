package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ClientUsecase_BizChangeStagesToGettingStartedEmail(t *testing.T) {
	er := UT.ClientUsecase.BizChangeStagesToGettingStartedEmail("6159272000005582016")
	lib.DPrintln(er)
}

func Test_ClientUsecase_HandleChangeStagesToGettingStartedEmail(t *testing.T) {
	er := UT.ClientUsecase.HandleChangeStagesToGettingStartedEmail("1206398481017098")
	lib.DPrintln(er)
}

func Test_ClientUsecase_BizChangeStagesToMiniDBQsFinalized(t *testing.T) {
	er := UT.ClientUsecase.BizChangeStagesToMiniDBQsFinalized("6159272000005582016")
	lib.DPrintln(er)
}

func Test_ClientUsecase_HandleBizChangeStagesToMiniDBQsFinalized(t *testing.T) {
	er := UT.ClientUsecase.HandleChangeStagesToMiniDBQsFinalized("abc")
	lib.DPrintln(er)
}

func Test_ClientUsecase_GetOneClientPipeline(t *testing.T) {
	str := "6159272000001008204"
	str = "abc"
	res, err := UT.ClientUsecase.GetOneClientPipeline(str)
	lib.DPrintln(err)
	lib.DPrintln(res)
}

func Test_ClientUsecase_GetClientsPipelines(t *testing.T) {
	clientGids := []string{
		"d0e39b4ae48246cbb91087fed6dd3973",
		"a",
		"6159272000001008204",
	}
	res, err := UT.ClientUsecase.GetClientsPipelines(clientGids)
	lib.DPrintln(err)
	lib.DPrintln(res)
}
