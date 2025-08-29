package biz

import (
	"testing"
	"vbc/tests"
)

var UT *tests.UnittestApp

func TestMain(m *testing.M) {
	tests.AppMain()
	UT = tests.UT
	m.Run()
}
