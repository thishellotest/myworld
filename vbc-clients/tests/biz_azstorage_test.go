package tests

import (
	"context"
	"io"
	"os"
	"testing"
	"vbc/lib"
)

func Test_AzstorageUsecase_UploadStream(t *testing.T) {

	a, err := os.Open("./tmp/STR Full_1.pdf")
	if err != nil {
		panic(err)
	}
	defer a.Close()

	ab, err := UT.AzstorageUsecase.UploadStream(context.TODO(), "tmp/STR Full_1.pdf", a)
	if err != nil {
		panic(err)
	}
	lib.DPrintln(ab.ClientRequestID)
}

func Test_AzstorageUsecase_UploadStream_From_Box(t *testing.T) {

	a, err := UT.BoxUsecase.DownloadFile("1481495306120", "")
	if err != nil {
		panic(err)
	}
	defer a.Close()
	ab, err := UT.AzstorageUsecase.UploadStream(context.TODO(), "tmp/USMC Medical Records_fromBox.pdf", a)
	if err != nil {
		panic(err)
	}
	lib.DPrintln(ab.ClientRequestID)
}

func Test_AzstorageUsecase_DownloadStream(t *testing.T) {
	response, err := UT.AzstorageUsecase.DownloadStream(context.TODO(),
		"blobs_dev/5e1140c7fd0f48948088529b7539065b/2_699c21cad73546808b050e0a0b21ec1b.json")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	file, err := os.Create("tmp/2_699c21cad73546808b050e0a0b21ec1b.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
}

func Test_AzstorageUsecase_DeleteBlob(t *testing.T) {
	err := UT.AzstorageUsecase.DeleteBlob(context.TODO(), "blobs_dev/5e1140c7fd0f48948088529b7539065b/fb60ab26bba146db98014a8eb69ed113.json")
	lib.DPrintln(err)
}

func Test_AzstorageUsecase_SasReadUrl(t *testing.T) {
	file := "blobs_dev/0ed0b3755f4f4ea0ad972e853e3bc431/d1c8b5733c374b2a94b42bc0f3509922.pdf"
	file = "blobs/a767d206e4f0424aa69c564e5db7c14f/176_93f6b5dc9f4b4db5a048814625d29106.pdf"
	file = "blobs/a767d206e4f0424aa69c564e5db7c14f/176_58673b148d044bbe97dd12d1df80c348.json"

	file = "blobs_dev/dc87883221884396ba374be83617d756/1_80b1204c69e240bbb1e7f915fd45e9d6.json"
	//file = "blobs/a767d206e4f0424aa69c564e5db7c14f/175_31cc4354ac72407fa4eda475248261c2.json"
	url, err := UT.AzstorageUsecase.SasReadUrl(file, nil)
	lib.DPrintln(url)
	lib.DPrintln(err)
}
