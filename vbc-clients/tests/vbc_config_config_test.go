package tests

import (
	"fmt"
	"testing"
	"vbc/internal/config_vbc"
)

func Test_Config(t *testing.T) {
	c := config_vbc.FeeDefine.Charge(0, 100)
	fmt.Println(c)
}
