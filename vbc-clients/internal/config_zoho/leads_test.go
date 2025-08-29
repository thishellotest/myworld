package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_LeadsLayout_SectionApiNames(t *testing.T) {
	sectionApiNames := LeadsLayout().SectionApiNames()
	lib.DPrintln(sectionApiNames)
}

func Test_LeadsLayout_LeadFieldInfos(t *testing.T) {
	res := LeadsLayout().LeadFieldInfos()
	lib.DPrintln(res)
}

func Test_LeadApiNames(t *testing.T) {
	res := LeadsLayout().LeadApiNames()
	lib.DPrintln(res)
}
