package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_FieldValidatorUsecase_CacheStructByKind(t *testing.T) {
	r, err := UT.FieldValidatorUsecase.CacheStructByKind(biz.Kind_client_cases)
	lib.DPrintln(err)
	a := r.GetByFieldName("stages", "Record Review")
	lib.DPrintln(a)
}

func Test_FieldValidatorUsecase_CacheFieldValidatorCenter(t *testing.T) {
	//center, err := UT.FieldValidatorUsecase.CacheFieldValidatorCenter(biz.Kind_client_cases)
	//center.Verify("")
	//a := r.GetByFieldName("stages", "Record Review")
	//lib.DPrintln(a)
}
