package tests

import (
	"os"
	"strings"
	"testing"
	"vbc/lib"
)

/*
func Test_PostUpload(t *testing.T) {
	url := "https://upload.box.com/api/2.0/files/content"

	values := map[string]*lib.UploadReader{
		"file":       lib.NewUploadReader(mustOpen("1.txt"), "1.txt"), // lets assume its this file
		"attributes": lib.NewUploadReader(strings.NewReader("{\"name\":\"1.txt\", \"parent\":{\"id\":\"241927737195\"}}"), ""),
	}
	err := lib.PostUpload(url, values, map[string]string{"authorization": "Bearer " + "mJ2zfbivJynqJj01PpCCd7H31WhEB7Gw"})
	lib.DPrintln(err)
}
*/

func Test_PostUpload(t *testing.T) {
	url := "https://upload.box.com/api/2.0/files/content"

	values := []*lib.UploadReader{
		lib.NewUploadReader("file", mustOpen("1.txt"), "1.txt"), // lets assume its this file
		lib.NewUploadReader("attributes", strings.NewReader("{\"name\":\"2.txt\", \"parent\":{\"id\":\"264924117433\"}}"), ""),
	}
	r, err := lib.PostUpload(url, values, map[string]string{"authorization": "Bearer " + "R393VVv6ee5csHyQSru94ytD7vn7zyIE"})
	lib.DPrintln(r, err)
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
