package biz

import (
	"testing"
	"vbc/lib"
)

func Test_TimeDateOnlyToTimestamp(t *testing.T) {
	a := TimeDateOnlyToTimestamp("2020-01-01")
	lib.DPrintln(a)
}
