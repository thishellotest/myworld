package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_ContactLayout_DealApiNames(t *testing.T) {
	res := ContactLayout().ContactApiNames()
	lib.DPrintln(res)
	lib.DPrintln(len(res))
}
