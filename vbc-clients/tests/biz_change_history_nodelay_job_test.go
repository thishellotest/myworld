package tests

import (
	"context"
	"sync"
	"testing"
	"vbc/lib"
)

func Test_ChangeHistoryNodelayJobUseacse_RunChangeHistoryNodelayJobJob(t *testing.T) {
	var syncwait sync.WaitGroup
	syncwait.Add(1)
	UT.ChangeHistoryNodelayJobUseacse.RunChangeHistoryNodelayJobJob(context.TODO())
	syncwait.Wait()
}

func Test_ChangeHistoryNodelayJobUseacse_GetLastValueByFieldName(t *testing.T) {
	a, b, err := UT.ChangeHistoryNodelayJobUseacse.GetLastValueByFieldName(5004, 0, "email")
	lib.DPrintln(a, b, err)
}
