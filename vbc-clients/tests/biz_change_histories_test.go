package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ChangeHistoryUseacse_GenTask(t *testing.T) {
	var entity biz.ChangeHistoryEntity
	UT.CommonUsecase.DB().Where("id=?", 89).Take(&entity)
	UT.ChangeHistoryUseacse.GenTask(&entity)
}

func Test_ChangeHistoryUseacse_HandleCreateFolderInBox(t *testing.T) {

	var entity biz.ChangeHistoryEntity
	UT.CommonUsecase.DB().Where("id=?", 12).Take(&entity)

	//UT.ChangeHistoryUseacse.HandleCreateFolderInBox(&entity)
}

//
//func Test_ChangeHistoryUseacse_CreateFolderInBox(t *testing.T) {
//
//	s, err := UT.ChangeHistoryUseacse.CreateFolderInBox(9)
//	lib.DPrintln(s)
//	lib.DPrintln(err)
//}

func Test_ChangeHistoryUseacse_HandleEnvelope(t *testing.T) {
	var entity biz.ChangeHistoryEntity
	UT.CommonUsecase.DB().Where("id=?", 55).Take(&entity)
	err := UT.ChangeHistoryUseacse.HandleEnvelope(&entity, "")
	fmt.Println(err)
}

func Test_HandleInitClientCaseChangeHistory(t *testing.T) {
	//err := UT.ChangeHistoryUseacse.HandleInitClientCaseChangeHistory(5086)
	//lib.DPrintln(err)
}

func Test_ChangeHistoryUseacse_HandlePrimaryCase(t *testing.T) {
	err := UT.ChangeHistoryUseacse.HandlePrimaryCase(5369)
	lib.DPrintln(err)
}

func Test_ChangeHistoryUseacse_HandleCreateFolderInBoxAndMail(t *testing.T) {
	var entity biz.ChangeHistoryEntity
	UT.CommonUsecase.DB().Where("id=?", 151554).Take(&entity)
	err := UT.ChangeHistoryUseacse.HandleCreateFolderInBoxAndMail(&entity)
	lib.DPrintln(err)
}

func Test_ChangeHistoryUseacse_TriggerGettingStartedEmail(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	err := UT.ChangeHistoryUseacse.TriggerGettingStartedEmail(*tCase)
	lib.DPrintln(err)
}
