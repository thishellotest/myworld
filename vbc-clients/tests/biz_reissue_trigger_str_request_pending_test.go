package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ReissueTriggerStrRequestPendingUsecase_ExistTask(t *testing.T) {
	a, err := UT.ReissueTriggerStrRequestPendingUsecase.ExistTask(5374)
	lib.DPrintln(a, err)
}

func Test_ReissueTriggerStrRequestPendingUsecase_LastTask(t *testing.T) {
	e, err := UT.ReissueTriggerStrRequestPendingUsecase.LastTask(66)
	lib.DPrintln(e, err)
}

func Test_ReissueTriggerStrRequestPendingUsecase_Handle(t *testing.T) {
	err := UT.ReissueTriggerStrRequestPendingUsecase.Handle()
	lib.DPrintln(err)
}
func Test_ReissueTriggerStrRequestPendingUsecase_DoOne(t *testing.T) {
	err := UT.ReissueTriggerStrRequestPendingUsecase.DoOne(42)
	lib.DPrintln(err)
}
