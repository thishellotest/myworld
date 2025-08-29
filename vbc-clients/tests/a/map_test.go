package a

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Fun(aa *map[string]bool) {
	(*aa)["aaa"] = true
}

func Test_aa(t *testing.T) {
	m := make(map[string]bool)
	Fun(&m)
	lib.DPrintln(m)
}

func Test_NotificationTextExtractContext(t *testing.T) {
	text := "3@[aaa]我是前 面@[Yannan Wang](6159272000000453001)我是123456789 @[Engineering Team](6159272000000453669) jjcc @[Yannan Wang](6159272000000453001) 1@[Engineering Team1](61592720000004536691)12345"
	result := biz.NotificationTextExtractContext(text, 15, 5)
	for _, r := range result {
		fmt.Println(r, "+++++")
	}
}
