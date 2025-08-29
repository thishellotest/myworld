package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_HaReportTasksBuzUsecase_CreateTask(t *testing.T) {

	// 1561231365497 PROD
	// 1580306003745 PROD 46页： 有空白页
	// 1584545143190 Test Local
	// 1584698450081 Test Local OnePage
	// 1584878283919 ：Test Local 46页： 有空白页
	err := UT.HaReportTasksBuzUsecase.CreateTask(context.TODO(), "3158719e92e846059f2336bdd810f449")
	lib.DPrintln(err)
}
