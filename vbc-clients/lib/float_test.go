package lib

import "testing"

func Test_Float32Sum(t *testing.T) {
	cc := Float32Sum(1.1, 2.44)
	DPrintln(cc)
}

func Test_FloatSum(t *testing.T) {
	cc := FloatSum(1.1, 2.44)
	DPrintln(cc)
}
