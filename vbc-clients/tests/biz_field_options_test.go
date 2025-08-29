package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_FieldOptionUsecase_CacheStructByKind(t *testing.T) {
	a, er := UT.FieldOptionUsecase.CacheStructByKind(biz.Kind_client_cases)
	lib.DPrintln(a, er)

	fieldStruct, _ := UT.FieldUsecase.StructByKind(biz.Kind_client_cases)
	//field1 := fieldStruct.GetByFieldName("stages")
	//cccc := a.AllByFieldName(*field1)
	//lib.DPrintln(cccc)

	field2 := fieldStruct.GetByFieldName("active_duty")
	cccc := a.AllByFieldName(*field2)
	lib.DPrintln(cccc)

	field2 = fieldStruct.GetByFieldName("agent_orange")
	cccc = a.AllByFieldName(*field2)
	lib.DPrintln(cccc)

	//c := a.AllByFieldName("stages")
	//lib.DPrintln(c.GetByLabel("1. Fee Schedule and Contract"))
}
