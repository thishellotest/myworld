package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Notes_HandleOne(t *testing.T) {
	UT.DataEntryUsecase.HandleOne(biz.Kind_notes, map[string]interface{}{
		"gid":         "ea2a010935a3489fbdb6bacb9687b963",
		"content":     "test1332221",
		"created_by":  "abc",
		"modified_by": "abc1",
		"kind":        "k1",
		"kind_gid":    "aaa",
	}, "gid", nil)
}

func Test_NotesUsecase_Save(t *testing.T) {
	gid, err := UT.NotesUsecase.Save("3f25dbcc6ed348e2a13178c7fa1c4f30", biz.Kind_client_cases, "2a75814a1a5242bc862a27ccad1bbf9a", "sssf", nil)
	lib.DPrintln(gid, err)
}
