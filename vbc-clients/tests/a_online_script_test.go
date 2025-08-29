package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_Sync_Pricing_Version(t *testing.T) {
	// 已经处理完毕
	cases, err := UT.TUsecase.ListByCond(biz.Kind_client_cases, And(Eq{
		biz.FieldName_pricing_version: "",
		biz.FieldName_biz_deleted_at:  0,
	}, Neq{biz.FieldName_s_pricing_version: ""}))
	if err != nil {
		panic(err)
	}
	for k, v := range cases {
		//lib.DPrintln(v.CustomFields.NumberValueByNameBasic("id"), v.CustomFields.TextValueByNameBasic("pricing_version"), v.CustomFields.TextValueByNameBasic("s_pricing_version"))

		err = UT.ZohobuzUsecase.HandleSyncZohoPricingVersion(cases[k])
		if err != nil {
			panic(err)
			lib.DPrintln(err, ":", v.CustomFields.NumberValueByNameBasic("id"))
		} else {
			lib.DPrintln("ok:", v.CustomFields.NumberValueByNameBasic("id"))
		}
	}
	lib.DPrintln(len(cases))
}
