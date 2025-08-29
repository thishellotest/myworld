package tests

import (
	"context"
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/esign/v2.1/envelopes"
)

func GetCred() *biz.DocuSignCredential {
	token, err := UT.Oauth2TokenUsecase.GetByAppId(biz.Oauth2_AppId_docusign)
	if err != nil {
		fmt.Println(err)
	}

	cred := biz.NewDocuSignCredential(token, UT.Conf.Docusign.AppAccountId)
	return cred
}

func Test_DocumentsGet(t *testing.T) {
	cred := GetCred()
	//envelopes.ListStatusChangesOp
	srv := envelopes.New(cred)
	// d4872d28-3c39-44da-9002-669c4b10549e
	// 75f10664-e171-4bcd-bf9d-36bd7edb9590
	// documentId: 1
	// documentId: certificate
	// documentID: certificate
	documentDownload, err := srv.DocumentsGet("certificate",
		"641eafbd-1d03-409b-b092-37219af0ae41").
		Do(context.Background())
	lib.DPrintln(err)
	lib.DPrintln("documentDownload:", documentDownload)
	if err != nil {
		return
	}

	_, err = UT.BoxUsecase.UploadFile("241927737195", documentDownload, "2.pdf")
	lib.DPrintln(err)

	/*
		file, err := os.Create("1.pdf")
		for {
			var bytes = make([]byte, 128)
			_, err := documentDownload.Read(bytes)
			lib.DPrintln(err)
			if err == io.EOF {
				break
			}
			_, err = file.Write(bytes)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		file.Close()*/
}

func Test_DocumentsList(t *testing.T) {
	cred := GetCred()
	//envelopes.ListStatusChangesOp
	srv := envelopes.New(cred)
	// 641eafbd-1d03-409b-b092-37219af0ae41
	//
	a, e := srv.DocumentsList("641eafbd-1d03-409b-b092-37219af0ae41").Do(context.Background())

	// /envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/attachments
	// /envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/documents
	//a, e := srv.DocumentsGet("/envelopes/2adc1cbb-5616-4810-bfa4-b7e658cb28ce/attachments",
	//	"2adc1cbb-5616-4810-bfa4-b7e658cb28ce").Do(context.Background())

	lib.DPrintln(a, e)
	//l, err := srv.Get("2adc1cbb-5616-4810-bfa4-b7e658cb28ce").
	//	Do(context.Background())
	//lib.DPrintln(l)
	//lib.DPrintln(err)
}
