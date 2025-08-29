package lib

import (
	"fmt"
	"testing"
)

func Test_SqlValueBackslash(t *testing.T) {

	a := "\"aa'\\aa"
	c := SqlValueBackslash(a)
	fmt.Println(c)
}
