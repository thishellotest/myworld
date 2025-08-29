package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_KindUsecase_GetByKind(t *testing.T) {
	e, err := UT.KindUsecase.GetByKind(biz.Kind_client_cases)
	lib.DPrintln(e, err)
}
