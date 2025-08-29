package tests

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_TimezoneUsecase_AllByCond(t *testing.T) {
	res, err := UT.TimezonesUsecase.AllByCond(nil)
	if err != nil {
		return
	}
	n := time.Now()
	for _, v := range res {
		lib.DPrintln(v.CodeValue)
		a, err := time.LoadLocation(v.CodeValue)
		if err != nil {
			panic(err)
		}
		aa, err := time.Parse(time.RFC3339, n.Format(time.RFC3339))
		if err != nil {
			panic(err)
		}
		aa = aa.In(a)
		lib.DPrintln(aa.Format(time.DateTime))
	}
}
