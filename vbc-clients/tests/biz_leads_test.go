package tests

import (
	"testing"
	"vbc/lib"
)

func Test_LeadsUsecase_BizLeadsSave(t *testing.T) {

	str := `{
	"firstName": "f1",
	"lastName": "l1",
	"email": "aaa@qq.com",
	"phone": "123",
	"state": "Alabama"
}`
	a, err := UT.LeadsUsecase.BizLeadsSave([]byte(str))
	lib.DPrintln(a, err)
}
