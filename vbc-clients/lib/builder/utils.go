package builder

import (
	"strings"
)

func ExprLikeBindValue(val string) string {
	c := strings.ReplaceAll(val, "%", "\\%")
	//lib.DPrintln("ssss__: ", c)
	return c
}
