package docx

import (
	"baliance.com/gooxml/document"
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_aaa(t *testing.T) {
	doc, err := document.Open("/tmp/PersonalStatements_ABagarry_5373.docx")
	if err != nil {
		panic(err)
	}

	lib.DPrintln("sss:", len(doc.Paragraphs()))
	//for _, para := range doc.Paragraphs() {
	//
	//}
}

func Test_readDocxText(t *testing.T) {
	text, err := biz.ReadDocxText("/tmp/PersonalStatements_ABagarry_5373_v1.docx")
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
