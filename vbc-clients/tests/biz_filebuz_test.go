package tests

import (
	"testing"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_FilebuzUsecase_DCRecordReviewFileHandle_1(t *testing.T) {
	a, _ := UT.BoxWebhookLogUsecase.GetByCond(Eq{"id": 55462})
	typeMap := lib.ToTypeMapByString(a.Body)
	err := UT.FilebuzUsecase.DCRecordReviewFileHandle(typeMap)
	lib.DPrintln(err)
}

func Test_FilebuzUsecase_DCRecordReviewFileHandle(t *testing.T) {
	typeMap := lib.ToTypeMapByString(FilebuzTriggerFileRename)
	err := UT.FilebuzUsecase.DCRecordReviewFileHandle(typeMap)
	lib.DPrintln(err)
}
func Test_FilebuzUsecase_DCRecordReviewFileHandle1(t *testing.T) {
	typeMap := lib.ToTypeMapByString(FilebuzTriggerFIleUploaded)
	err := UT.FilebuzUsecase.DCRecordReviewFileHandle(typeMap)
	lib.DPrintln(err)
}

func Test_FilebuzUsecase_DCRecordReviewFileHandle2(t *testing.T) {
	typeMap := lib.ToTypeMapByString(FilebuzTriggerFolderMove)
	err := UT.FilebuzUsecase.DCRecordReviewFileHandle(typeMap)
	lib.DPrintln(err)
}

func Test_FilebuzUsecase_DCRecordReviewFileHandle3(t *testing.T) {
	typeMap := lib.ToTypeMapByString(FilebuzTriggerTrashed)
	err := UT.FilebuzUsecase.DCRecordReviewFileHandle(typeMap)
	lib.DPrintln(err)
}
