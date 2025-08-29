package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_BoxFileNameFilter(t *testing.T) {
	fileName := "aa\na|\"aaa*ac:a#$a\ta\ra\\.pdf"
	a := biz.BoxFileNameFilter(fileName)
	lib.DPrintln(fileName)
	lib.DPrintln(a)
}

func Test_BoxWebhookLogUsecase_HandleExec(t *testing.T) {
	task, _ := UT.BoxWebhookLogUsecase.GetByCond(Eq{"id": 42219})
	err := UT.BoxWebhookLogUsecase.HandleExec(context.TODO(), task)
	lib.DPrintln(err)
}

func Test_BoxWebhookLogUsecase_RunHandleJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.BoxWebhookLogUsecase.RunHandleJob(context.TODO())
	wait.Wait()
}

func Test_BoxWebhookLogUsecase_Handle(t *testing.T) {

	task, _ := UT.BoxWebhookLogUsecase.GetByCond(Eq{"id": 5182})
	err := UT.BoxWebhookLogUsecase.Handle(context.TODO(), task)
	lib.DPrintln(err)
}

func Test_BoxWebhookLogUsecase_HandleQuestionnairesJotform(t *testing.T) {

	tCase, _ := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 5182})
	err := UT.BoxWebhookLogUsecase.HandleQuestionnairesJotform("6037031974216948122.pdf", tCase)
	lib.DPrintln(err)
}

func Test_BoxWebhookLogUsecase_HandleQuestionnairesBuz(t *testing.T) {

	needMove, err := UT.BoxWebhookLogUsecase.HandleQuestionnairesBuz("264812519184",
		"TestGary TestLiaoInitial Intake2024-05-19 00:04:51", "264658751993", nil)
	lib.DPrintln(needMove, err)
}

func Test_Fix_QuestionnaireDownload(t *testing.T) {

	res, err := UT.BoxUsecase.ListItemsInFolderFormat("258821922844")
	if err != nil {
		panic(err)
	}
	for _, v := range res {
		if v.GetString("type") == "folder" {
			temp, err := UT.BoxUsecase.ListItemsInFolderFormat(v.GetString("id"))
			if err != nil {
				panic(err)
			}
			for _, v1 := range temp {
				fileInfo, _, err := UT.BoxUsecase.GetFileInfo(v1.GetString("id"))
				if err != nil {
					panic(err)
				}
				if fileInfo == nil {
					panic("is nil")
				}
				fileInfoMap := lib.ToTypeMapByString(*fileInfo)
				newFileInfoMap := make(lib.TypeMap)
				newFileInfoMap.Set("source", fileInfoMap)

				//lib.DPrintln(newFileInfoMap)

				err = UT.BoxWebhookLogUsecase.HandleNewJotformName(newFileInfoMap)

				if err != nil {
					panic(err)
				}
				lib.DPrintln(newFileInfoMap)
			}
		}
		//lib.DPrintln(v)
		//break
		//lib.DPrintln(v.GetString(""))
	}
}

func Test_BoxWebhookLogUsecase_HandleNewJotformNameStep2Copy(t *testing.T) {
	tCase, _ := UT.TUsecase.Data(biz.Kind_client_cases, Eq{"id": 5511})
	err := UT.BoxWebhookLogUsecase.HandleNewJotformNameStep2Copy("1802087108824", time.Now(), "6176710185019648226", tCase)
	lib.DPrintln(err)
}

func Test_BoxWebhookLogUsecase_CrontabEveryOneHourHandleQuestionnaireDownloads(t *testing.T) {
	err := UT.BoxWebhookLogUsecase.CrontabEveryOneHourHandleQuestionnaireDownloads()
	lib.DPrintln(err)
}

func Test_GetJotformSubmissionsIdFromFolderName(t *testing.T) {
	folderName := "Test1 TestL Muscle Injuries Increase 2025-03-14 19:48:10 -6178160905016488647# 5511"
	aa := biz.GetJotformSubmissionsIdFromFolderName(folderName)
	lib.DPrintln(aa)
}
