package tests

import (
	"testing"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_ZoomMeetingUsecase_FolderName(t *testing.T) {
	e, _ := UT.ZoomMeetingUsecase.GetByCond(builder.Eq{"id": 2})
	r, err := UT.ZoomMeetingUsecase.FolderName(e)
	lib.DPrintln(r, err)
}

func Test_ZoomMeetingUsecase_RenameBoxFolderName(t *testing.T) {

	return
	for i := 4; i <= 162; i++ {
		e, _ := UT.ZoomMeetingUsecase.GetByCond(builder.Eq{"id": i})
		err := UT.ZoomMeetingUsecase.RenameBoxFolderName(e)
		lib.DPrintln(err)
	}
}
