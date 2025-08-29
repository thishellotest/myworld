package tests

import (
	"testing"
	"vbc/lib"
)

func Test_HttpAccessControl_HandleTasks(t *testing.T) {
	UT.HttpAccessControl.HandleTasks("57140384499e47dcaf11daa4eba408bd")
}

func Test_HttpAccessControl_HandeCarryOut(t *testing.T) {
	err := UT.HttpAccessControl.HandeCarryOut("57140384499e47dcaf11daa4eba408bd", 0)
	lib.DPrintln(err)
}
