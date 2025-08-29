package tests

import (
	"testing"
	"vbc/lib"
)

func Test_CounterUsecase_Stat(t *testing.T) {
	UT.CounterUsecase.Stat("aaa", 2)
	UT.CounterUsecase.Stat("aaa", 1)
}

func Test_CounterUsecase_HasLimit(t *testing.T) {
	r, err := UT.CounterUsecase.HasLimit("aaa", 6)
	lib.DPrintln(r, err)
}
