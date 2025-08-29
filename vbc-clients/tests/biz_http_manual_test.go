package tests

import (
	gin2 "github.com/gin-gonic/gin"
	"testing"
	"vbc/lib"
)

func Test_http_manual(t *testing.T) {
	gin := &gin2.Context{}
	UT.HttpManualUsecase.ReplenishClientUniqcode(gin)
}

func Test_HttpManualUsecase_HandleSyncAllClientFromAsana(t *testing.T) {
	err := UT.HttpManualUsecase.HandleSyncAllClientFromAsana("")
	lib.DPrintln(err)
}
