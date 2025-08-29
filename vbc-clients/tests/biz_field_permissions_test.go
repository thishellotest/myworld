package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_FieldPermissionUsecase(t *testing.T) {
	e, _ := UT.FieldPermissionUsecase.GetByCond(builder.Eq{"id": 3})
	lib.DPrintln(e.Permission)
}

func Test_FieldPermissionUsecase_CacheStructByKind(t *testing.T) {
	a, err := UT.FieldPermissionUsecase.CacheStructByKind(biz.Kind_client_cases, "65f75c7f965f4387af33da18d5e12036")
	lib.DPrintln(a, err)
	c := a.GetByFieldName("deal_name")
	lib.DPrintln(c)

	a, err = UT.FieldPermissionUsecase.CacheStructByKind(biz.Kind_client_cases, "65f75c7f965f4387af33da18d5e12036")
	//lib.DPrintln(a, err)
	c = a.GetByFieldName("deal_name")
	lib.DPrintln(c)

	c = a.GetByFieldName("deal_name22")
	lib.DPrintln(c)
}

func Test_FieldPermissionUsecase_CacheFieldPermissionCenter(t *testing.T) {

	a, err := UT.FieldPermissionUsecase.FieldPermissionCenter(biz.Kind_client_cases, "65f75c7f965f4387af33da18d5e12036")
	lib.DPrintln(a, err)
	lib.DPrintln(a.PermissionByFieldName("deal_name"))
	lib.DPrintln(a.PermissionByFieldName("user_gid"))
	lib.DPrintln(a.PermissionByFieldName("created_by"))

	lib.DPrintln(a.PermissionByFieldName("amount"))
	lib.DPrintln(a.PermissionByFieldName("pricing_version"))
	lib.DPrintln(a.PermissionByFieldName("s_pricing_version"))
}

func Test_FieldPermissionUsecase_CacheFieldPermissionCenter2(t *testing.T) {

	a, err := UT.FieldPermissionUsecase.FieldPermissionCenter(biz.Kind_users, "441f19d51858417cb948cc286ef1b585")
	lib.DPrintln(a, err)
	lib.DPrintln(a.PermissionByFieldName("profile_gid"))

}
