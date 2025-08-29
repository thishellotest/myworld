package biz

import (
	"testing"
	"vbc/lib"
)

func Test_TimeToVBCDisplay(t *testing.T) {
	c, err := TimeToVBCDisplay("2024-03-29T23:06:04+08:00")
	lib.DPrintln(c, err)
}
