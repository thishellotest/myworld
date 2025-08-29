package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_ContactLayout_SectionApiNames(t *testing.T) {
	sectionApiNames := ContactLayout().SectionApiNames()
	lib.DPrintln(sectionApiNames)
}

func Test_DealLayout_SectionApiNames(t *testing.T) {
	sectionApiNames := DealLayout().SectionApiNames()
	lib.DPrintln(sectionApiNames)
}
