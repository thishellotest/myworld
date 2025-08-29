package tests

import (
	"testing"
	"vbc/lib"
)

func Test_MaCongratsEmailUsecase_HandleInputTask(t *testing.T) {
	err := UT.MaCongratsEmailUsecase.HandleInputTask(8)
	lib.DPrintln(err)
}
