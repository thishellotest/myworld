package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_TimelinesbuzUsecase_List(t *testing.T) {
	r, total, err := UT.TimelinesbuzUsecase.List(biz.Kind_clients, "78fbe690068a46d3a2b8d76270fdea0b", nil, 1, 3)
	lib.DPrintln(r)
	lib.DPrintln(total)
	lib.DPrintln(err)
}
