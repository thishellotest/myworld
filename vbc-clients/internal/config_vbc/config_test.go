package config_vbc

import (
	"encoding/json"
	"testing"
	"vbc/lib"
)

func Test_FeeDefine(t *testing.T) {
	config := make(map[string]interface{})
	config["BoxSignTpl"] = BoxSignTpl
	config["FeeDefine"] = FeeDefine
	a, er := json.Marshal(config)
	lib.DPrintln(string(a))
	lib.DPrintln(er)
}
