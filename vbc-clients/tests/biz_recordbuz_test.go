package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RecordbuzUsecase_List(t *testing.T) {
	str := `{"filter":{"operator":"AND","group":[{"comparator":"is_not_empty","field":{"field_name":"sys__due_date"},"value":[]}]}}`
	tListRequest := lib.StringToTDef[*biz.TListRequest](str, nil)
	list, err := UT.RecordbuzUsecase.List(biz.Kind_client_cases, nil, tListRequest, 1, 10, nil, false)

	lib.DPrintln(list)
	lib.DPrintln(err)
}
