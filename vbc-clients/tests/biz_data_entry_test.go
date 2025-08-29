package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_DataEntryUsecase_UpdateOne(t *testing.T) {
	row := biz.TypeDataEntry{
		"name":     "33333",
		"case_gid": "ss",
		"uniqblob": "3",
		"gid":      "d4be2a63c3a243e08be7547289aeda2b",
	}
	_, err := UT.DataEntryUsecase.UpdateOne(biz.Kind_blobs, row, "gid", nil)
	fmt.Println(err)
}

func Test_DataEntryUsecase_InsertOne(t *testing.T) {
	row := biz.TypeDataEntry{
		"name":     "%1'\"",
		"case_gid": "ss",
		"uniqblob": "3",
	}
	gid, err := UT.DataEntryUsecase.InsertOne(biz.Kind_blobs, row, nil)
	fmt.Println(err, gid)
}

func Test_DataEntryUsecase_InsertOneData(t *testing.T) {

	row := biz.TypeDataEntry{
		"name":     "%1'\"",
		"case_gid": "ss",
	}
	structField, err := UT.FieldUsecase.CacheStructByKind(biz.Kind_blobs)
	if err != nil {
		fmt.Println(err)
	}
	err = UT.DataEntryUsecase.InsertOneData("blobs", row, structField, nil)
	fmt.Println(err)
}

func Test_DataEntryUsecase_InsertData(t *testing.T) {

	datalist := biz.TypeDataEntryList{
		{
			"id":         11,
			"first_name": "%1'\"",
			"new_rating": "ss",
		},
		//{
		//	"first_name": "2#",
		//	"new_rating": "4",
		//},
	}
	structField, err := UT.FieldUsecase.CacheStructByKind("")
	if err != nil {
		fmt.Println(err)
	}

	_, err = UT.DataEntryUsecase.InsertData(biz.KindEntity{}, "clients", datalist, structField, "id", nil)
	fmt.Println(err)
}

func Test_DataEntryUsecase_UpdateModifyData(t *testing.T) {
	structField, err := UT.FieldUsecase.CacheStructByKind("")
	if err != nil {
		fmt.Println(err)
	}

	ke, _ := UT.KindUsecase.GetByKind(biz.Kind_clients)
	UT.DataEntryUsecase.UpdateModifyData(*ke, "clients", map[string]interface{}{
		"id":         "9",
		"first_name": "aaa1",
		"new_rating": "101",
		"stages":     "stage1",
	}, map[string]interface{}{

		"first_name": "a%#a\\a\"2'",
		"new_rating": "11",
		"stages":     "stage2",
	}, structField, "id", nil)
}

func Test_DataEntryUsecase_Handle(t *testing.T) {
	datalist := biz.TypeDataEntryList{
		{
			//"new_rating": "100",
			"id":  200,
			"gid": "gid03",
		},
		{
			"asana_task_gid": "100",
			"new_rating":     "488",
			"id":             201,
			"gid":            "gid04",
		},
	}
	//
	result, err := UT.DataEntryUsecase.Handle("", datalist, "gid", nil)
	lib.DPrintln("result: ", result)
	lib.DPrintln(err)
}

func Test_DataEntryUsecase_Handle1Insert(t *testing.T) {
	datalist := biz.TypeDataEntryList{
		{
			//"new_rating": "100",
			"new_rating": "40",
			"first_name": "F1_",
			//"full_name":  "XX ss",
			"gid": "gid15",
		},
		{
			"asana_task_gid": "100",
			"new_rating":     "48",
			"first_name":     "F2_",
			"last_name":      "L2_",
			"gid":            "gid16",
		},
	}
	//
	result, err := UT.DataEntryUsecase.Handle(biz.Kind_clients, datalist, "gid", nil)
	lib.DPrintln("result: ", result)
	lib.DPrintln(err)
}

func Test_DataEntryUsecase_Handle1(t *testing.T) {
	datalist := biz.TypeDataEntryList{
		{
			//"new_rating": "100",
			"id":         200,
			"new_rating": "40",
			"first_name": "F1_",
			"full_name":  "XX ss",
			"gid":        "gid03",
		},
		{
			"asana_task_gid": "100",
			"new_rating":     "48",
			"id":             201,
			"first_name":     "F2_",
			"last_name":      "L2_",
			"gid":            "gid04",
		},
	}
	//
	result, err := UT.DataEntryUsecase.Handle(biz.Kind_clients, datalist, "gid", nil)
	lib.DPrintln("result: ", result)
	lib.DPrintln(err)
}

func Test_DataEntryUsecase_HandleChangeHistories(t *testing.T) {
	//UT.DataEntryUsecase.HandleChangeHistories("", 1)
}
