package uuid

import (
	satoriuuid "github.com/satori/go.uuid"
	"strings"
)

func UuidWithoutStrike() string {
	a := satoriuuid.NewV4().String()
	return strings.Replace(a, "-", "", -1)
}
