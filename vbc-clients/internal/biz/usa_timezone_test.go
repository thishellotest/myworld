package biz

import (
	"testing"
	"time"
	"vbc/lib"
)

func Test_UsaTime(t *testing.T) {
	UsaTime()
}

func Test_GetLocationByUsaTimezone(t *testing.T) {
	aa, err := GetLocationByUsaTimezone(USA_TIMEZONE_AK)
	now := time.Now().In(aa).Format(time.DateTime)
	lib.DPrintln(aa, err)
	lib.DPrintln(now)
	lib.DPrintln(time.Now().Format(time.DateTime))
}
