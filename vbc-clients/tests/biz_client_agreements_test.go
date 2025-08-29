package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_ClientAgreementUsecase_Update(t *testing.T) {
	a, er := UT.AdobeSignUsecase.GetAgreement(context.TODO(), "CBJCHBCAABAAQHOUh6xTo-ifxir7if6T5CgpUndpzuQn")
	lib.DPrintln(er)
	lib.DPrintln(a)
	err := UT.ClientAgreementUsecase.Update(a)
	lib.DPrintln(err)
}
