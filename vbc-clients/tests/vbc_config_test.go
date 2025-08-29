package tests

import (
	"fmt"
	"testing"
	"vbc/internal/config_vbc"
)

func Test_ca(t *testing.T) {
	a := config_vbc.StateConfigs.FullNameByShort("WA")
	fmt.Println(a)
}
