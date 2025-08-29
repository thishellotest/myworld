package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_BlobbuzUsecase_HandleBoxFile(t *testing.T) {
	a, _, _ := UT.BoxUsecase.GetFileInfoForTypeMap("1467179637164")
	_, err := UT.BlobbuzUsecase.HandleBoxFile(context.TODO(), a, "", "")
	lib.DPrintln(err)
}

func Test_BlobbuzUsecase_HandleBlobSlices(t *testing.T) {
	blob, _ := UT.BlobUsecase.GetByUniqblob("1560495055079_1714677775079")
	fileInfo, _, _ := UT.BoxUsecase.GetFileInfoForTypeMap("1560495055079")
	err := UT.BlobbuzUsecase.HandleBlobSlices(context.TODO(), blob, fileInfo.GetString("id"), fileInfo.GetString("file_version.id"))
	lib.DPrintln(err)
}

func Test_BlobbuzUsecase_BizRecordReviewTasksProgress(t *testing.T) {
	//user, _ := UT.UserUsecase.GetUserFacadeById(4)
	//UT.BlobbuzUsecase.BizRecordReviewTasksProgress(context.TODO(), *user, []string{"072c021699a042188564a6ca2edae696", "563a05bafc1b4994942f2fa015979c34"})
}
