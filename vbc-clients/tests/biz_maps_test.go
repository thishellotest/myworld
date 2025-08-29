package tests

import (
	"fmt"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_MapUsecase_Set(t *testing.T) {
	ccc := time.Now().Unix()
	lib.DPrintln(ccc)
	UT.MapUsecase.Set("aa", biz.InterfaceToString(ccc))
	a, err := UT.MapUsecase.GetForInt("aa")
	lib.DPrintln(a, err)
}

func Test_MapUsecase_Get(t *testing.T) {
	a, e := UT.MapUsecase.GetForString("aa1")
	fmt.Println(a, e)
}

func Test_MapUsecase_GetByCond(t *testing.T) {
	entity, err := UT.MapUsecase.GetByCond(And(Eq{"mval": "256143829415"}, Like{"mkey", "ClientBoxFolderId:%"}))
	lib.DPrintln(entity, err)
}
