package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_XeroUsecase_Accounts(t *testing.T) {
	err := UT.XeroUsecase.Accounts()
	lib.DPrintln(err)
}

func Test_XeroUsecase_GetContacts(t *testing.T) {
	err := UT.XeroUsecase.GetContacts()
	lib.DPrintln(err)
}
func Test_XeroUsecase_GetInvoice(t *testing.T) {
	err := UT.XeroUsecase.GetInvoice("e85d1165-040f-47d8-99c4-20113ef45b25")
	lib.DPrintln(err)
}

func Test_XeroUsecase_BizCreateOrUpdateContact(t *testing.T) {
	tClient, err := UT.TUsecase.DataById(biz.Kind_client_cases, 2)
	lib.DPrintln(err, tClient.CustomFields.NumberValueByNameBasic("id"))
	//a, err := UT.XeroUsecase.BizCreateOrUpdateContact(tClient)
	//lib.DPrintln(a, err)
}

func Test_XeroUsecase_CreateOrUpdateContact(t *testing.T) {
	tClient, err := UT.TUsecase.DataById(biz.Kind_client_cases, 44)
	lib.DPrintln(err)
	res, err := UT.XeroUsecase.CreateOrUpdateContact(tClient, "48d47f54-f562-4618-a386-cf9b8a5c3963")
	lib.DPrintln(res, err)
}

//func Test_XeroUsecase_DoBizCreateInvoice(t *testing.T) {
//	tClient, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
//	lib.DPrintln(err)
//	err = UT.XeroUsecase.DoBizCreateInvoice(tClient)
//	fmt.Println(err)
//}

func Test_XeroUsecase_BizCreateInvoice(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
	lib.DPrintln(err)
	a, b, err := UT.XeroUsecase.BizCreateInvoice(tCase)
	fmt.Println(a, b, err)
}

func Test_XeroUsecase_BizAmCreateInvoice(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5750)
	lib.DPrintln(err)
	a, b, err := UT.XeroUsecase.BizAmCreateInvoice(tCase)
	fmt.Println(a, b, err)
}

func Test_XeroUsecase_CreateInvoice(t *testing.T) {
	r, err := UT.XeroUsecase.CreateInvoice(12, "7452b18f-223a-4e79-b6db-660db3176f3e", biz.Xero_VBC_BrandingThemeID)
	lib.DPrintln(r, err)
}

func Test_XeroUsecase_BrandingThemes(t *testing.T) {
	res, err := UT.XeroUsecase.BrandingThemes()
	lib.DPrintln(res, err)
}
