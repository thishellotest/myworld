package utils

import (
	"testing"
	"time"
	"vbc/configs"
	"vbc/lib"
)

func Test_CalIntervalDayTime(t *testing.T) {

	aa, err := CalIntervalDayTime(time.Now(), 1, "08:00", *configs.GetVBCDefaultLocation())
	lib.DPrintln(time.Now().Format(time.RFC3339))
	lib.DPrintln(aa, err)
}

func Test_CalDelayDayTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2024-10-22T01:01:00-07:00")
	a := CalDelayDayTime(t1, *configs.GetVBCDefaultLocation())
	r := a.Format(time.RFC3339)
	lib.DPrintln("r:", r)
}
