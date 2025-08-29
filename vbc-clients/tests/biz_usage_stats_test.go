package tests

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_aa(t *testing.T) {
	err := UT.UsageStatsUsecase.Stat("a", time.Now(), 2)
	lib.DPrintln(err)
}
