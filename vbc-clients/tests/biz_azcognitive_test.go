package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_AzcognitiveUsecase_PrebuiltRead(t *testing.T) {

	blob, err := UT.AzstorageUsecase.DownloadStream(context.TODO(), "blobs/8093553f034943af8f1296a1f2212c8c/3.pdf")
	if err != nil {
		panic(err)
	}
	a, err := UT.AzcognitiveUsecase.PrebuiltRead(blob.Body)
	lib.DPrintln(a, err)
}

func Test_AzcognitiveUsecase_GetPrebuiltReadResultWithBlock(t *testing.T) {
	res, err := UT.AzcognitiveUsecase.GetPrebuiltReadResultWithBlock(context.TODO(), "https://documentintelligenceeu2s0.cognitiveservices.azure.com/formrecognizer/documentModels/prebuilt-read/analyzeResults/4f9fe37a-a1df-4703-a8c8-113e3bf33930?api-version=2023-07-31")
	lib.DPrintln(err)
	lib.DPrintln(res)
	lib.DPrintln(lib.StringToBytesFormat(lib.InterfaceToString(res)))
}
