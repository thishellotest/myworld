package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

func Test_Dag_BuzEmail(t *testing.T) {
	var task biz.TaskEntity
	err := UT.CommonUsecase.DB().Where("id=?", 14533).Take(&task).Error
	lib.DPrintln(err)
	err = UT.Dag.BuzEmail(&task)
	lib.DPrintln(err)
}

func Test_Dag_GetEnvelopeDocuments(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=22").Take(&entity)
	err := UT.Dag.GetEnvelopeDocuments(&entity)
	lib.DPrintln(err)
}

func Test_Dag_CreateTask_SaveSignedContractInBox(t *testing.T) {
	a := UT.TaskCreateUsecase.CreateTask(0,
		map[string]interface{}{"envelope_id": "641eafbd-1d03-409b-b092-37219af0ae41", "folder_id": "241927737195"},
		biz.Task_Dag_SaveSignedContractInBox, 0, "", "")
	lib.DPrintln(a)
}

func Test_Dag_SaveSignedContractInBox(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=23").Take(&entity)
	err := UT.Dag.SaveSignedContractInBox(&entity)
	lib.DPrintln(err)
}

func Test_Dag_CreateEnvelopeAndSentFromAdobeSign(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=27").Take(&entity)
	err := UT.Dag.CreateEnvelopeAndSentFromAdobeSign(&entity)
	lib.DPrintln(err)
}

func Test_Dag_SaveSignedContractInBox_Adobe(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=28").Take(&entity)
	err := UT.Dag.SaveSignedContractInBox(&entity)
	lib.DPrintln(err)
}

func Test_Dag_BizCreateBoxFolder(t *testing.T) {
	err := UT.Dag.BizCreateBoxFolder(5813)
	lib.DPrintln(err)
}

func Test_Dag_BizCreateBoxFolder_Fix(t *testing.T) {

	return
	clientCaseId := int32(5209)
	tClientCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, clientCaseId)
	if err != nil {
		panic(err)
	}
	if tClientCase == nil {
		panic("tClientCase is nil.")
	}

	_, tContactFields, err := UT.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		panic(err)
	}
	if tContactFields == nil {
		panic("tContactFields is nil.")
	}

	key := fmt.Sprintf("%s%d", biz.Map_ClientBoxFolderId, clientCaseId)
	boxFolderId, _ := UT.MapUsecase.GetForString(key)
	if boxFolderId == "" {
		panic("boxFolderId")
	}

	lib.DPrintln(boxFolderId)
	//return
	// 创建共享
	email := tContactFields.TextValueByNameBasic("email")
	_, err = UT.BoxUsecase.Collaborations(boxFolderId, email)

	if err != nil {
		panic(err)
	}

	// 更新zoho box文件夹
	deal, err := UT.ZohoUsecase.GetDeal(tClientCase.CustomFields.TextValueByNameBasic("gid"))
	if err != nil {
		panic(err)
	}
	if deal == nil {
		panic("zoho deal is nil.")
	}
	if deal.GetString("Case_Files_Folder") == "" {
		row := make(lib.TypeMap)
		key := fmt.Sprintf("%s%d", biz.Map_ClientBoxFolderId, clientCaseId)
		boxFolderId, _ = UT.MapUsecase.GetForString(key)
		if boxFolderId != "" {
			row.Set("Case_Files_Folder", "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
		}
		if len(row) > 0 {
			row.Set("id", tClientCase.CustomFields.TextValueByNameBasic("gid"))
			UT.ZohoUsecase.PutRecordV1(config_zoho.Deals, row)
		}
		lib.DPrintln("BizCreateBoxFolder Case_Files_Folder: row: ", row)
	}
}

func Test_Dag_HandleNonResponsive(t *testing.T) {
	err := UT.Dag.HandleNonResponsive(5094, false)
	lib.DPrintln(err)
}

func Test_Dag_CronTrigger(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=48").Take(&entity)
	err := UT.Dag.CronTrigger(&entity)
	lib.DPrintln(err)
}

func Test_Dag_CreateEnvelopeAndSentFromBoxWithTemplate(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=7480").Take(&entity)
	err := UT.Dag.CreateEnvelopeAndSentFromBoxWithTemplate(&entity)
	lib.DPrintln(err)
}

func Test_Dag_DoCreateEnvelopeAndSentFromBoxAm(t *testing.T) {

	err := UT.Dag.DoCreateEnvelopeAndSentFromBoxAm(5728)
	lib.DPrintln(err)
}
