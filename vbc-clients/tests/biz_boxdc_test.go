package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_BoxdcUsecase_RecordReviewFirstSubFolderByName(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	folderId, err := UT.BoxdcUsecase.RecordReviewFirstSubFolderByName("VA Medical Records", tCase)
	lib.DPrintln(folderId, err)
	folderId, err = UT.BoxdcUsecase.RecordReviewFirstSubFolderByName("Service Treatment Records", tCase)
	lib.DPrintln(folderId, err)
	folderId, err = UT.BoxdcUsecase.RecordReviewFirstSubFolderByName("Private Medical Records", tCase)
	lib.DPrintln(folderId, err)
}

func Test_BoxdcUsecase_RecordReviewFirstSubFolderByName_notexist(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	folderId, err := UT.BoxdcUsecase.RecordReviewFirstSubFolderByName("VA Medical Records1", tCase)
	lib.DPrintln(folderId, err)
}
