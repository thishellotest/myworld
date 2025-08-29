package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

func Test_AsanaUsecase_GetATask(t *testing.T) {
	// 1206237216158413
	// 1206234446219801
	// GetATask 1206343093612429 生产环境
	// 1206545140591678 生产环境
	// 1206398481017098 // test 还在
	// 1206234446219801 test
	// 1206398481017098
	// 1206237216158413
	// 1206398481017067 已删
	// 1205468380506421 生产，asana看不到
	// 1206163668062111 生产，stages为空
	a, isDel, er := UT.AsanaUsecase.GetATask("1206721348864689")
	lib.DPrintln(a, isDel, er)
	return
	c := a.ToDataEntry()
	lib.DPrintln(c)
	var dataList biz.TypeDataEntryList
	dataList = append(dataList, c)

	_, err := UT.DataEntryUsecase.Handle("", dataList, "asana_task_gid", nil)
	fmt.Println(err)
}

func Test_AsanaUsecase_CreateATask(t *testing.T) {

	field := config_vbc.GetAsanaCustomFields()
	firstNameGid := field.GetByName("First Name").GetGid()
	fmt.Println("firstNameGid:", firstNameGid)
	lastNameGid := field.GetByName("Last Name").GetGid()
	fmt.Println("lastNameGid:", lastNameGid)
	emailGid := field.GetByName("Email").GetGid()
	fmt.Println("emailGid:", emailGid)
	phoneNumberGid := field.GetByName("Phone Number").GetGid()
	fmt.Println("phoneNumberGid:", phoneNumberGid)
	stateGid := field.GetByName("Address - State").GetGid()
	fmt.Println("stateGid:", stateGid)

	customFields := make(lib.TypeMap)
	customFields.Set(firstNameGid, "First")
	customFields.Set(lastNameGid, "Last")
	customFields.Set(emailGid, "aa@qq.com")
	r, err := UT.AsanaUsecase.CreateATask(customFields, "First", "Last", "I already havea rating, but I think I could be underrated.")
	lib.DPrintln(r)
	lib.DPrintln(err)
}

func Test_AsanaUsecase_change_stages(t *testing.T) {
	//taskInfo, err := UT.AsanaUsecase.GetATask("1206398481017098")

	field := config_vbc.GetAsanaCustomFields()
	fGid := field.GetByName("Stages").GetGid()
	eGid := field.GetByName("Stages").GetEnumGidByName(config_vbc.Stages_GettingStartedEmail)
	lib.DPrintln(fGid, eGid)
	customFields := make(lib.TypeMap)
	customFields.Set(fGid, eGid)
	r, er := UT.AsanaUsecase.PutATask("1206398481017098", customFields, "")
	lib.DPrintln(r, er)

	eGid1 := field.GetByName("Source").GetEnumGidByName(config_vbc.Source_Manual)
	eGid2 := field.GetByName("Source").GetEnumGidByName(config_vbc.Source_Website)
	fmt.Println(eGid1, eGid2, "===")
}

func Test_AsanaUsecase_PutATask_only_test(t *testing.T) {

	r, er := UT.AsanaUsecase.PutATask("1206398481017098", nil, "")
	lib.DPrintln(r, er)
}

func Test_AsanaUsecase_AsanaUsecase(t *testing.T) {
	UT.AsanaUsecase.AddAProjectToATask("1206398481017098", UT.Conf.Asana.ProjectGidCp)
}

func Test_AsanaUsecase_RemoveAProjectToATask(t *testing.T) {
	UT.AsanaUsecase.RemoveAProjectToATask("1206398481017098", UT.Conf.Asana.ProjectGid)
}

func Test_AsanaUsecase_PutATask(t *testing.T) {

	customFields := make(lib.TypeMap)
	customFields.Set("1205964025409311", "New_ssn val")
	r, err := UT.AsanaUsecase.PutATask("1206237015703342", customFields, "")
	lib.DPrintln(r)
	lib.DPrintln(err)
}

func Test_TypeMapAndList(t *testing.T) {
	a, er := UT.AsanaUsecase.GetAUser("1205444097333494")
	lib.DPrintln(a, er)
	c := a.Get("data")
	fmt.Println("rrr:", c)
	ccc := lib.ToTypeList(a.Get("data.workspaces"))

	for _, v := range ccc {
		t1 := v.Get("name")
		fmt.Println(t1)
	}

	lib.DPrintln("+++", ccc)
	//
	//for _, v := range ccc {
	//	cccc := v.(map[string]interface{})
	//	lib.DPrintln(cccc)
	//}
}

func Test_AsanaUsecase_GetAUser(t *testing.T) {
	a, er := UT.AsanaUsecase.GetAUser("1205444097333494")
	lib.DPrintln(a, er)
	c := a.Get("data")
	fmt.Println("rrr:", c)
	ccc := lib.ToTypeList(a.Get("data.workspaces"))

	for _, v := range ccc {
		t1 := v.Get("name")
		fmt.Println(t1)
	}

	lib.DPrintln("+++", ccc)
	//
	//for _, v := range ccc {
	//	cccc := v.(map[string]interface{})
	//	lib.DPrintln(cccc)
	//}

	//for _, v := range ccc {
	//	aaa := lib.TypeMap(v)
	//	lib.DPrintln(aaa)
	//}
	//fmt.Println("rrr:", c)
	return
	//lib.DPrintln(a, er)
	//c := a.ToDataEntry()
	//lib.DPrintln(c)
	//var dataList biz.TypeDataEntryList
	//dataList = append(dataList, c)
	//
	//err := UT.DataEntryUsecase.Handle("", dataList, "asana_task_gid")
	//fmt.Println(err)
}

func Test_AsanaUsecase_ListWebhooks(t *testing.T) {
	a, er := UT.AsanaUsecase.ListWebhooks()
	lib.DPrintln(er)
	lib.DPrintln(a)
}

func Test_AsanaUsecase_(t *testing.T) {
	UT.AsanaUsecase.GetAllProjects(config_vbc.AsanaWorkspaceGid)

}
