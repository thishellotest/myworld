package lib

import "testing"

func Test_ArrayReverse(t *testing.T) {
	//res := []string{}
	var aaa []string
	a := ArrayReverse(aaa)
	DPrintln(a, aaa)
}

func Test_RemoveDuplicates(t *testing.T) {
	dest := RemoveDuplicates([]int{1, 2, 2, 4})
	DPrintln(dest)
}
