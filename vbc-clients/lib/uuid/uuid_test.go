package uuid

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
	"vbc/lib"
)

func Test_uuid(t *testing.T) {
	a := uuid.NewV4().String()
	fmt.Println(a, len(a))
}

func Test_UuidWithoutStrike(t *testing.T) {
	aa := UuidWithoutStrike()
	lib.DPrintln(aa)
}
