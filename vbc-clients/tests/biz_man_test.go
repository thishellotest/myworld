package tests

import "testing"

func Test_ManUsecase_HandleHistoryCreateEnvelope(t *testing.T) {
	err := UT.ManUsecase.HandleHistoryCreateEnvelope()
	if err != nil {
		panic(err)
	}
}
