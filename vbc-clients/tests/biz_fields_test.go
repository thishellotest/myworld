package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_FieldUsecase_ListByKind(t *testing.T) {
	a, err := UT.FieldUsecase.ListByKind("")
	lib.DPrintln(a, err)
}

func Test_FieldUsecase_CacheStructByKind1(t *testing.T) {
	a, err := UT.FieldUsecase.CacheStructByKind(biz.Kind_clients)
	lib.DPrintln(a, err)
	lib.DPrintln(a.GetByFieldName("full_name").Dependence)
	c := a.GetFieldDependList()
	ccc := c.GetFieldNameIsDependent("last_name")
	lib.DPrintln("__", ccc)

	r1 := c.GetFieldNamesIsDependent([]string{"last_name", "first_name", "info"})
	lib.DPrintln("___::__:: ", r1)
}

func Test_FieldUsecase_CacheStructByKind(t *testing.T) {
	a, err := UT.FieldUsecase.CacheStructByKind("")
	lib.DPrintln(a, err)
	c := a.GetByFieldName("deal_name")
	lib.DPrintln(c)
}
