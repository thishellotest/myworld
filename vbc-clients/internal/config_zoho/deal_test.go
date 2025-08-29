package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_DealLayout2_DealFieldInfos(t *testing.T) {
	DealFieldInfos := DealLayout().DealFieldInfos()
	lib.DPrintln(DealFieldInfos)
}

func Test_DealLayout2_DealApiNames(t *testing.T) {
	res := DealLayout().DealApiNames()
	lib.DPrintln(res)
	lib.DPrintln(len(res))
}

func Test_DealLayout2_DealApiNames2(t *testing.T) {
	res := DealLayout().DealApiNames2()
	lib.DPrintln(res)
	lib.DPrintln(len(res))
}

func Test_DealLayout2_PrintDealApiNames(t *testing.T) {
	DealLayout().PrintDealApiNames()
}

func Test_DealLayout_FieldByApiName(t *testing.T) {
	a := DealLayout()
	c := a.FieldByApiName("Lead_Source")
	pickListValues := c.GetTypeList("pick_list_values")
	for _, v := range pickListValues {
		lib.DPrintln(v.GetString("display_value"))
		//lib.DPrintln(v.GetString("display_value"), "|", v.GetString("actual_value"), " | ", v.GetString("id"))
	}
}

func Test_DealLayout_FieldByApiName_Branch(t *testing.T) {
	a := DealLayout()
	c := a.FieldByApiName("Branch")
	pickListValues := c.GetTypeList("pick_list_values")
	for _, v := range pickListValues {
		lib.DPrintln(v.GetString("display_value"))
	}
}
