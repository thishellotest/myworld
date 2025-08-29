package config_zoho

import (
	"fmt"
	"testing"
)

func Test_ZohoDealFieldNameByVbcFieldName(t *testing.T) {
	a := ZohoDealFieldNameByVbcFieldName("address")
	fmt.Println(a)
}

func Test_ZohoContactFieldNameByVbcFieldName(t *testing.T) {
	a := ZohoContactFieldNameByVbcFieldName("address")
	fmt.Println(a)
}
