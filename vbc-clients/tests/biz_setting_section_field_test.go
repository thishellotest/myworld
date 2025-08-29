package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_SettingSectionFieldUsecase_GetSettingSectionsByKind(t *testing.T) {
	vo, err := UT.SettingSectionFieldUsecase.GetSettingSectionsByKind("", "65f75c7f965f4387af33da18d5e12036", biz.Record_From_detail)
	lib.DPrintln(vo, err)
}

func Test_SettingSectionFieldUsecase_GetSettingSectionsByKindForFormColumn(t *testing.T) {
	//vo, err := UT.SettingSectionFieldUsecase.GetSettingSectionsByKindForFormColumn(biz.Kind_users, "441f19d51858417cb948cc286ef1b585", biz.RecordHttpFormColumnRequest_Type_Form_mgmt_user)
	//lib.DPrintln(vo, err)
}
