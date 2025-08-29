package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Timelines_HandleOne(t *testing.T) {
	UT.DataEntryUsecase.HandleOne(biz.Kind_timelines, map[string]interface{}{
		"gid":              "g124",
		"notes":            "{}",
		"kind":             "3331aaaa",
		"kind_gid":         "aaacc",
		"action":           "added",
		"related_kind":     "kin2",
		"related_kind_gid": "g1",
		"created_by":       "u1",
		"modified_by":      "u3",
		//"id":               3,
	}, "gid", nil)
}

func Test_TimelineUsecase_Create(t *testing.T) {
	gid, err := UT.TimelineUsecase.Create(biz.Kind_client_cases, "g1", "added", "", "", "ss",
		nil)
	lib.DPrintln(gid, err)
}
