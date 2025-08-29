package lib

import "testing"

func Test_files(t *testing.T) {
	a1, aa := FileExt("aac/..ssaaa.pdf", false)
	DPrintln(a1, aa)
	a1, aa = FileExt("aaapdf", false)
	DPrintln(a1, aa)
	a1, aa = FileExt("aaapdf.doc.pdf", false)
	DPrintln(a1, aa)
	a1, aa = FileExt(".aaapdf.doc", false)
	DPrintln(a1, aa)
}

func Test_TrimHiddenCharacter(t *testing.T) {
	aaa := TrimHiddenCharacter("aa 我是中文\n\ta.\\ */a cc")
	DPrintln(aaa)
}

func Test_TrimCharacterForFileName(t *testing.T) {
	aaa := TrimCharacterForFileName("aa 我是中文\n\ta.\\ */a cc")
	DPrintln(aaa)
}
