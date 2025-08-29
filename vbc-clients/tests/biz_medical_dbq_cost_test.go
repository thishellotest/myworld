package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_MedicalDbqCostCalculator(t *testing.T) {
	a := biz.MedicalDbqCostCalculator(true, 10, 10)
	lib.DPrintln(a)
}

func Test_MedicalDbqCostUsecase_GetMedicalDbqCost(t *testing.T) {
	costInfo := UT.MedicalDbqCostUsecase.GetMedicalDbqCost(32)
	lib.DPrintln(costInfo)
}
