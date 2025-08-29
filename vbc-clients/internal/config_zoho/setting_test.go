package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_ContactLayout(t *testing.T) {
	a := ContactLayout()
	fieldApiNames := a.FieldApiNamesByApiName(Contact_Sections_ApiName_Client_Information)
	lib.DPrintln(fieldApiNames)
}

func Test_DealLayout(t *testing.T) {
	a := DealLayout()
	fieldApiNames := a.FieldApiNamesByApiName(Deal_Sections_ApiName_ClientCaseInformation)
	lib.DPrintln(fieldApiNames)
}

func Test_DealLayout1(t *testing.T) {
	a := DealLayout()
	fieldApiNames := a.FieldInfoByApiName(Deal_Sections_ApiName_Description_Information)
	lib.DPrintln(fieldApiNames)
}

func Test_ContactLayout2_DealFieldInfos(t *testing.T) {
	FieldInfos := ContactLayout().ContactFieldInfos()
	lib.DPrintln(FieldInfos)
}

func Test_ContactFieldInfos(t *testing.T) {
	r := ContactLayout().ContactFieldInfos()
	lib.DPrintln(r)
}

func Test_aa(t *testing.T) {

}
