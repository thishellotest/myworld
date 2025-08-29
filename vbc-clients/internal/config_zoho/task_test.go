package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_TaskLayout_SectionApiNames(t *testing.T) {
	sectionApiNames := TaskLayout().SectionApiNames()
	lib.DPrintln(sectionApiNames)
}

func Test_TaskLayout_TaskFieldInfos(t *testing.T) {
	DealFieldInfos := TaskLayout().TaskFieldInfos()
	lib.DPrintln(DealFieldInfos)
}

func Test_TaskLayout_TaskApiNames(t *testing.T) {
	res := TaskLayout().TaskApiNames()
	lib.DPrintln(res)
}

/*
Not Started
Deferred
In Progress
Completed
Waiting for input
*/
func Test_TaskLayout_FieldByApiName(t *testing.T) {
	a := TaskLayout()
	c := a.FieldByApiName("Status")
	pickListValues := c.GetTypeList("pick_list_values")
	for _, v := range pickListValues {
		lib.DPrintln(v.GetString("display_value"))
	}
}
