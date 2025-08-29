package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_QueueUsecase_PushClientTaskHandleWhoGidJobTasks(t *testing.T) {
	err := UT.QueueUsecase.PushClientTaskHandleWhoGidJobTasks(context.TODO(), []string{"6159272000005519042", "scc", "s3"})
	lib.DPrintln(err)
}

func Test_QueueUsecase_PushClientTaskHandleWhatGidJobTasks(t *testing.T) {
	err := UT.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), []string{"6159272000009972111", "scc", "s3"})
	lib.DPrintln(err)
}

func Test_QueueUsecase_PushClientTaskHandleWhatGidJobTasks_All(t *testing.T) {

	//caches := lib.CacheInit[*biz.TData]()
	//res, _ := UT.TUsecase.ListByCondCaches(&caches, biz.Kind_client_cases,
	//	builder.And(builder.Eq{"1": 1}, builder.In("id", []int{5360, 5399})))
	//res, _ := UT.TUsecase.ListByCondCaches(&caches, biz.Kind_client_cases,
	//	builder.And(builder.Eq{"1": 1}))

	kindEntity, _ := UT.KindUsecase.GetByKind((biz.Kind_client_cases))

	res, _ := UT.TUsecase.ListByCondWithPaging(*kindEntity,
		builder.And(builder.Eq{"1": 1}), "", 1, 10)

	for _, v := range res {
		err := UT.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), []string{v.Gid()})
		lib.DPrintln(err)
		lib.DPrintln(len(v.CustomFields))
	}
}

func Test_QueueUsecase_PushClientNameChangeJobTasks(t *testing.T) {
	UT.QueueUsecase.PushClientNameChangeJobTasks(context.TODO(), []string{"aaa", "bb", "aaa", "bb", "sssc"})
}
