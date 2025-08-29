package biz

import (
	"testing"
	"vbc/lib"
)

func Test_IsUSAPhone(t *testing.T) {
	a := IsUSAPhone("+1111-111-1111")
	lib.DPrintln(a)
}
