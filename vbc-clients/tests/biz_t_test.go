package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_TUsecase_t(t *testing.T) {
	a, er := UT.TUsecase.List("", nil, &biz.TListRequest{}, 1, 10)
	lib.DPrintln(a, er)
	b, er := UT.TUsecase.Total("")
	lib.DPrintln(b, er)
}

//func Test_TUsecase_List(t *testing.T) {
//	a, err := UT.TUsecase.ListByCond("", Neq{"stages": ""})
//	lib.DPrintln(a)
//	lib.DPrintln(err)
//}

func Test_TUsecase_Data(t *testing.T) {

	a, er := UT.TUsecase.Data(biz.Kind_users, And(Eq{"id": 4}, Eq{"id": 4}))
	lib.DPrintln(a, er)
	if a != nil {
		c := a.CustomFields.DisplayValueByName("id")
		lib.DPrintln(c)
		m := a.CustomFields.ToMaps()
		lib.DPrintln(m)
		m = a.CustomFields.ToDisplayMaps()
		lib.DPrintln(m)
		boxUserId := a.CustomFields.TextValueByNameBasic(biz.UserFieldName_box_user_id)
		lib.DPrintln(boxUserId)
	}
}

func Test_TUsecase_Data_1(t *testing.T) {

	a, er := UT.TUsecase.DataByGid(biz.Kind_client_cases, "355047448cdf4503a4348d3b68d7163b")
	lib.DPrintln(a.CustomFields.ToMaps())
	lib.DPrintln(er)
}

func Test_TUsecase_Data_EmtailTpls(t *testing.T) {

	a, er := UT.TUsecase.Data(biz.Kind_email_tpls, And(Eq{"tpl_type": "InitialContact"}))
	lib.DPrintln(a, er)
	if a != nil {
		c := a.CustomFields.DisplayValueByName("id")
		lib.DPrintln(c)
		m := a.CustomFields.ToMaps()
		lib.DPrintln(m)
		m = a.CustomFields.ToDisplayMaps()
		lib.DPrintln(m)
	}
}

func Test_TUsecase_Blob(t *testing.T) {

	a, er := UT.TUsecase.Data(biz.Kind_blobs, And(Eq{"id": 17}))
	lib.DPrintln(a, er)
	//if a != nil {
	//	c := a.CustomFields.DisplayValueByName("id")
	//	lib.DPrintln(c)
	//	m := a.CustomFields.ToMaps()
	//	lib.DPrintln(m)
	//	m = a.CustomFields.ToDisplayMaps()
	//	lib.DPrintln(m)
	//}
}

func Test_TUsecase_Cases(t *testing.T) {

	a, er := UT.TUsecase.DataById(biz.Kind_client_cases, 8)
	lib.DPrintln(a, er)

	lib.DPrintln("DisplayValueByNameBasic:", a.CustomFields.DisplayValueByNameBasic("created_at"))
	lib.DPrintln("ValueByName:", a.CustomFields.ValueByName("created_at"))
	lib.DPrintln("DisplayValueByNameBasic:", a.CustomFields.DisplayValueByNameBasic("created_at"))

	//if a != nil {
	//	c := a.CustomFields.DisplayValueByName("id")
	//	lib.DPrintln(c)
	//	m := a.CustomFields.ToMaps()
	//	lib.DPrintln(m)
	//	m = a.CustomFields.ToDisplayMaps()
	//	lib.DPrintln(m)
	//}
}

func Test_TUsecase_DataByGidWithCaches(t *testing.T) {
	caches := lib.CacheInit[*biz.TData]()

	// 5360,5399

	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_users, "6159272000001027142")
	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_users, "6159272000001027142")
	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_users, "6159272000001027142")
	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_client_cases, "123")
	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_client_cases, "123")
	UT.TUsecase.DataByGidWithCaches(&caches, biz.Kind_users, "6159272000001027142")
}

func Test_TUsecase_DataById_User(t *testing.T) {
	a, er := UT.TUsecase.DataById(biz.Kind_users, 4)
	lib.DPrintln(a, er)
}

func Test_TUsecase_ListByCond_(t *testing.T) {
	a, er := UT.TUsecase.ListByCond(biz.Kind_users, Eq{"id": 4})
	for _, v := range a {
		c := v.CustomFields.ValueByName("profile_gid")
		lib.DPrintln(c)
		c1 := v.CustomFields.ValueByName("role_gid")
		lib.DPrintln(c1)
	}
	lib.DPrintln(a, er)
}

func Test_Taa(t *testing.T) {
	a, err := UT.TUsecase.DataByGid(biz.Kind_attorneys, "5f6e50a11a2a40739028a0323c50cf95")
	lib.DPrintln(a.CustomFields.ToMaps(), err)
}
