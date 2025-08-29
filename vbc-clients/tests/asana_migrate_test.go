package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Asana_Migrate(t *testing.T) {
	sqlRows, err := UT.CommonUsecase.DB().Raw(biz.AsanaMigrateSql("")).Rows()
	lib.DPrintln(err)
	_, list, err := lib.SqlRowsTrans(sqlRows)
	lib.DPrintln(list)
}
