package tests

import (
	"context"
	"fmt"
	"testing"
)

func Test_SyncAsanaTaskUsecase_InitSyncAsanaTask(t *testing.T) {

	err := UT.SyncAsanaTaskUsecase.InitSyncAsanaTask(context.TODO())
	fmt.Println(err)
	//UT.CommonUsecase.RedisClient().LPush(context.TODO(), "aa", aa)
}

func Test_SyncAsanaTaskUsecase_GetTasks(t *testing.T) {
	a := UT.SyncAsanaTaskUsecase.GetTasks(context.TODO())
	fmt.Println(a)
}

func Test_SyncAsanaTaskUsecase_LPushSyncTaskQueue(t *testing.T) {
	//aa := []string{"c1", "c2"}
	UT.SyncAsanaTaskUsecase.LPushSyncTaskQueue(context.TODO(), "c3")
}

func Test_SyncAsanaTaskUsecase_FinishSyncTask(t *testing.T) {
	UT.SyncAsanaTaskUsecase.FinishSyncTask(context.TODO())
}

func Test_SyncAsanaTaskUsecase_SyncTask(t *testing.T) {
	// 1206178230625366
	// 1206234446219801
	// 1206511875930215 生产
	// 1206398481017098 没有删除
	// 1206398481017067
	err := UT.SyncAsanaTaskUsecase.SyncTask("1206398481017067")
	fmt.Println(err)
}

func Test_SyncAsanaTaskUsecase_SyncUser(t *testing.T) {
	// 1205444097333494 1206230291638946
	err := UT.SyncAsanaTaskUsecase.SyncUser("1206230291638946")
	fmt.Println(err)
}
