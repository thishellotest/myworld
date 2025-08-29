package lib

import "testing"

func Test_GeneratePassword(t *testing.T) {
	a := GeneratePassword(16)
	DPrintln(a)
}
