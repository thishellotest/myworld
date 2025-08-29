package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_HttpBlobUsecase_BizCreateTask(t *testing.T) {
	tUser, _ := UT.TUsecase.DataById(biz.Kind_users, 4)
	a, err := UT.HttpBlobUsecase.BizCreateTask(context.TODO(), *tUser, "1560495055079")
	lib.DPrintln(a, err)
}

func Test_HttpBlobUsecase_BizRecordReviewDetail(t *testing.T) {
	user, _ := UT.UserUsecase.GetUserFacadeById(4)
	a, err := UT.HttpBlobUsecase.BizRecordReviewDetail(context.TODO(), *user, "dc87883221884396ba374be83617d756", nil)
	lib.DPrintln(a, err)
}

func Test_HttpBlobUsecase_BizRecordReviewFiles(t *testing.T) {
	user, _ := UT.UserUsecase.GetUserFacadeById(4)
	a, err := UT.HttpBlobUsecase.BizRecordReviewFiles(context.TODO(), *user, "d1fbcc1328424c3699057dd71f14e970")
	lib.DPrintln(a, err)
}
