package config_vbc

import (
	"fmt"
	"testing"
)

func Test_a(t *testing.T) {
	a := GetAsanaCustomFields()
	c := a.GetByName("Agent Orange Exposure")
	s := c.GetEnumGidByName("Yes")
	fmt.Println(s)
}
