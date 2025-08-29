package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RoleUsecase_CacheAllRoles(t *testing.T) {
	aa, err := UT.RoleUsecase.CacheAllRoles()
	lib.DPrintln(aa, err)
	aa, err = UT.RoleUsecase.CacheAllRoles()
	lib.DPrintln(aa, err)
}

func Test_RoleUsecase_HandleChildrenRoles(t *testing.T) {
	aa, _ := UT.RoleUsecase.CacheAllRoles()
	tRole, _ := UT.TUsecase.DataByGid(biz.Kind_roles, "926b0b34c5134d6c98ae3ba484b8afbb")
	var dest []*biz.TData
	UT.RoleUsecase.HandleChildrenRoles(aa, tRole, &dest)
	for _, v := range dest {
		lib.DPrintln(v.Gid(), " ", v.CustomFields.TextValueByNameBasic("role_name"))
	}
}

func Test_RoleUsecase_ChildrenRoles(t *testing.T) {
	tRole, _ := UT.TUsecase.DataByGid(biz.Kind_roles, "926b0b34c5134d6c98ae3ba484b8afbb")
	dest, _ := UT.RoleUsecase.ChildrenRoles(tRole)
	for _, v := range dest {
		lib.DPrintln(v.Gid(), " ", v.CustomFields.TextValueByNameBasic("role_name"))
	}
}
func Test_RoleUsecase_ChildrenRolesGids(t *testing.T) {
	tRole, _ := UT.TUsecase.DataByGid(biz.Kind_roles, "926b0b34c5134d6c98ae3ba484b8afbb")
	dest, _ := UT.RoleUsecase.ChildrenRolesGids(tRole)
	lib.DPrintln(dest)
}
