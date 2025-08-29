package tests

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_TimeTz(t *testing.T) {
	LoadLocation := time.FixedZone("CST", 8*3600)
	time.Local = LoadLocation
	a, err := lib.TimeParse("2023-12-27T00:00:00.000Z")
	lib.DPrintln(err)
	a = a.In(LoadLocation)
	lib.DPrintln("sss", a.Format("2006-01-02 15:04:05"), a.Local())
}
