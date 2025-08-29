package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ReminderEventUsecase_AddClientUpdateFilesEvent(t *testing.T) {
	vo := biz.ReminderClientUpdateFilesEventVo{
		Items: []*biz.ReminderClientUpdateFilesEventVoItem{
			{
				BoxResName: "TestFN TestLN - Release of Information Form",
				BoxResId:   "1550118786574",
				BoxResType: "file",
			},
			{
				BoxResName: "TestFN TestLN - Release of Information Form",
				BoxResId:   "1550118786574",
				BoxResType: "file",
			},
			{
				BoxResName: "Folder1",
				BoxResId:   "269153665339",
				BoxResType: "folder",
			},
		},
	}

	err := UT.ReminderEventUsecase.AddClientUpdateFilesEvent(5004, &vo)
	lib.DPrintln(err)
}
