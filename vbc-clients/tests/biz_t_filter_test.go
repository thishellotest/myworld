package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_TFilter_a(t *testing.T) {

	str := `{
	"comparator": "eq",
	"field": {
		"field_name": "active_duty"
	},
	"value": [{
		"value": "Yes",
		"label": "Yes"
	}, {
		"value": "No",
		"label": "No"
	}]
}`
	a := lib.StringToT[biz.TListCondition](str)
	ccc := a.Value()
	ddd := ccc.ValueForOptions()

	lib.DPrintln(ddd)
}
