package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_EnvelopeStatusChangeUsecase_GetByCond(t *testing.T) {

	a, err := UT.EnvelopeStatusChangeUsecase.GetByCond(builder.Eq{"envelope_id": 1})
	lib.DPrintln(a, err)

	a, err = UT.EnvelopeStatusChangeUsecase.GetByCond(builder.Eq{"envelope_id": "\\a'\"."})
	lib.DPrintln(a, err, a)
}

func Test_EnvelopeStatusChangeUsecase_Handle(t *testing.T) {
	a, err := UT.EnvelopeStatusChangeUsecase.GetByCond(
		builder.Eq{"envelope_id": "295e31a2-faa6-4c38-8373-f0f25de37e57"})
	lib.DPrintln(a, err)
	err = UT.EnvelopeStatusChangeUsecase.Handle(a)
	fmt.Println(err)
}

//
//func Test_EnvelopeStatusChangeUsecase_HandleCompleted(t *testing.T) {
//	a, err := UT.EnvelopeStatusChangeUsecase.GetByCond(UT.CommonUsecase.DB(),
//		builder.Eq{"envelope_id": "295e31a2-faa6-4c38-8373-f0f25de37e57"})
//	lib.DPrintln(a, err)
//	err = UT.EnvelopeStatusChangeUsecase.HandleCompleted(a)
//	fmt.Println(err)
//}

func Test_EnvelopeStatusChangeUsecase_RunEnvelopeStatusChangeJob(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	UT.EnvelopeStatusChangeUsecase.RunEnvelopeStatusChangeJob(context.TODO())
	wait.Wait()
}
