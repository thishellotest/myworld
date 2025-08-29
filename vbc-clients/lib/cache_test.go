package lib

import "testing"

func Test_abc(t *testing.T) {
	aaa := CacheInit[*string]()
	aaa.Set("abc", nil)
	e, exists := aaa.Get("abc")
	DPrintln(exists, e)
}
