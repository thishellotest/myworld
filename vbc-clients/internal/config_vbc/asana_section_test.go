package config_vbc

import (
	"fmt"
	"testing"
)

func Test_AsanaSections_GetSectionGidByName(t *testing.T) {
	gid := GetAsanaSections().GetSectionGidByName(AsanaSection_GETTING_STARTED_EMAIL)
	fmt.Println(gid)
}
