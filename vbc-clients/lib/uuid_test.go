package lib

import (
	"fmt"
	"testing"
)

func Test_google_a(t *testing.T) {
	i := 0
	aaa := make(map[string]bool)
	for {
		i++

		id := UuidNumeric()
		fmt.Println(id, len(id))
		if _, ok := aaa[id]; ok {
			panic("sss")
		} else {
			aaa[id] = true
		}
		if i > 100000 {
			break
		}
	}
}

func Test_UuidNumericTime(t *testing.T) {
	c := UuidNumeric()
	fmt.Println(len(c), c)
}

func Test_RuntimePath(t *testing.T) {
	fmt.Println(RuntimePath())
}
