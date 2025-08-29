package tests

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_BUsaStateUsecase_GetTimeLocationByUsaState(t *testing.T) {
	a, loc, err := UT.BUsaStateUsecase.GetTimeLocationByUsaState("Hawaii")
	now := time.Now().In(loc).Format(time.DateTime)
	lib.DPrintln(a, loc, err, now)
}
