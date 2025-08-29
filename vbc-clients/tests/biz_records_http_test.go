package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RecordHttpUsecase_BizList(t *testing.T) {
	str := `{"filter":{"operator":"AND","group":[{"comparator":"gt","field":{"field_name":"updated_at"},"value":["2024-12-20"]}]}}`
	tUser, _ := UT.TUsecase.DataByGid(biz.Kind_users, "6159272000000453669")
	userFacade := biz.UserFacade{
		TData: *tUser,
	}
	detail, err := UT.RecordHttpUsecase.BizList("", userFacade, []byte(str), 1, 3, "")
	lib.DPrintln(detail, err)
}

func Test_RecordHttpUsecase_BizDetail(t *testing.T) {
	var tUser biz.TData
	str := `{
		"gid":"6159272000000820046"
	}`
	uf := biz.UserFacade{
		TData: tUser,
	}
	detail, err := UT.RecordHttpUsecase.BizDetail("", "6159272000000820046", uf, []byte(str), biz.Record_From_detail)
	lib.DPrintln(detail, err)
}

func Test_RecordHttpUsecase_BizSave(t *testing.T) {

	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	//data := make(lib.TypeMap)
	//data.Set("", "")
	str := `{
"collaborators":{"value":"6159272000000453001,6159272000000453669,6159272000001027094"},
"effective_current_rating":{"value":"30"},
"active_duty":{"value":"No"},
"referrer":{"value":"1"},
"client_gid":{"value":"gid03"}}`

	str = `{
"effective_current_rating":{"value":"40"},
"active_duty":{"value":"Yes"},
"referrer":{"value":"123"},
"client_gid":{"value":"6159272000005519042"}}`

	data := lib.ToTypeMapByString(str)
	gid := "6159272000011891231"
	aaa, r, err := UT.RecordHttpUsecase.BizSave(biz.Kind_client_cases, gid, *userFacade, data, true, nil, false)
	lib.DPrintln(aaa, r, err)
}

func Test_RecordHttpUsecase_BizRelatedClient(t *testing.T) {
	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	aa, err := UT.RecordHttpUsecase.BizRelatedClient(*userFacade, "6159272000000708025")
	lib.DPrintln(aa, err)
}

func Test_RecordHttpUsecase_BizTimelines(t *testing.T) {
	userFacade, _ := UT.UserUsecase.GetUserFacadeByGid("6159272000000453669")
	aa, err := UT.RecordHttpUsecase.BizTimelines(biz.Kind_clients, "78fbe690068a46d3a2b8d76270fdea0b", *userFacade, nil, 1, 100)
	lib.DPrintln(aa, err)
}
