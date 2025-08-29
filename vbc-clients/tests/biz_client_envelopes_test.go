package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ClientEnvelopeUsecase_GetByEnvelopeId(t *testing.T) {
	e, err := UT.ClientEnvelopeUsecase.GetByEnvelopeId("box", "ab0991e8-2404-44d8-b179-670feb859633")
	lib.DPrintln(e.BoxContactFileId(), err)
}

func Test_ClientEnvelopeUsecase_GetBoxSignByCaseId(t *testing.T) {
	e, err := UT.ClientEnvelopeUsecase.GetBoxSignByCaseId(5728, biz.Type_AmContract)
	lib.DPrintln(e, err)
}

func Test_ClientEnvelopeUsecase_ContractDateOn(t *testing.T) {
	a, err := UT.ClientEnvelopeUsecase.ContractDateOn(5004, false)
	lib.DPrintln(a, err)
}

func Test_ClientEnvelopeUsecase_AmContractBoxFileId(t *testing.T) {
	boxFileId, err := UT.ClientEnvelopeUsecase.AmContractBoxFileId(5809)
	lib.DPrintln(boxFileId, err)
}
