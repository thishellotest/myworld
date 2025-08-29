package tests

import (
	"net/url"
	"testing"
	"time"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

func Test_ZohoDealScanJobUsecase_BizRunJob(t *testing.T) {
	err := UT.ZohoDealScanJobUsecase.BizRunJob()
	lib.DPrintln(err)
}

func Test_ZohoDealScanJobUsecase_a(t *testing.T) {

	perPage := 100
	fields := config_zoho.DealLayout().DealApiNames()
	time.Sleep(time.Second)
	params := make(url.Values)
	params.Add("page", lib.InterfaceToString(1))
	params.Add("per_page", lib.InterfaceToString(perPage))
	lib.DPrintln(params)
	records, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	lib.DPrintln(records)
	lib.DPrintln(err)
}
