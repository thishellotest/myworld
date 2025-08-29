package tests

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_RecordReviewJobUsecase_HandleTask(t *testing.T) {

	str := `{
	"UniqueKey": "5076268258622342folder",
	"Params": "{\"ClientCaseId\":5004,\"FirstSubFolderId\":\"264686394097\",\"FirstSubFolderName\":\"VA Medical Records\",\"DestId\":\"268258622342\",\"DestType\":\"folder\",\"DestName\":\"leve1_folder\"}"
}`

	//	str = `{
	//	"UniqueKey": "5076268258622342folder",
	//	"Params": "{\"ClientCaseId\":5004,\"FirstSubFolderId\":\"264686394097\",\"FirstSubFolderName\":\"VA Medical Records\",\"DestId\":\"1534417478116\",\"DestType\":\"file\",\"DestName\":\"VA Medical Records\"}"
	//}`

	var customTaskParams biz.CustomTaskParams
	err := json.Unmarshal([]byte(str), &customTaskParams)
	if err != nil {
		panic(err)
	}
	err = UT.RecordReviewJobUsecase.HandleTask(context.TODO(), customTaskParams)
	lib.DPrintln(err)
}

func Test_RecordReviewJobUsecase_LPushCustomTaskQueue(t *testing.T) {
	a := biz.CustomTaskParams{
		UniqueKey: "44",
	}
	err := UT.RecordReviewJobUsecase.LPushCustomTaskQueue(context.TODO(), a, a)
	lib.DPrintln(err)
}

func Test_RecordReviewJobUsecase_RunCustomTaskJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)

	err := UT.RecordReviewJobUsecase.RunCustomTaskJob(context.TODO())
	lib.DPrintln(err)

	wait.Wait()
}

func Test_RecordReviewJobUsecase_BizHandleTaskV1(t *testing.T) {

	str := `{
	"UniqueKey": "5076:1559005155765:file",
	"Params": "{\"ClientCaseId\":5076,\"FirstSubFolderId\":\"264686394097\",\"FirstSubFolderName\":\"VA Medical Records\",\"DestId\":\"1555593308983\",\"DestType\":\"file\",\"DestName\":\"b.pdf\",\"SourceFromId\":5538}"
}`
	str = `{"UniqueKey":"5814:1954487526959:file","Params":"{\"ClientCaseId\":5814,\"FirstSubFolderId\":\"333843805321\",\"FirstSubFolderName\":\"Service Treatment Records\",\"DestId\":\"1954487526959\",\"DestType\":\"file\",\"DestName\":\"111.jpeg\",\"SourceFromId\":174401}"}`
	var customTaskParams biz.CustomTaskParams
	err := json.Unmarshal([]byte(str), &customTaskParams)
	if err != nil {
		panic(err)
	}
	err = UT.RecordReviewJobUsecase.BizHandleTaskV1(context.TODO(), customTaskParams)
	lib.DPrintln(err)
}
