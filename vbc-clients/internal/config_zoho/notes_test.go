package config_zoho

import (
	"testing"
	"vbc/lib"
)

func Test_NoteLayout_NoteApiNames(t *testing.T) {
	res := NotesLayout()
	res.LeadFieldInfos()
	lib.DPrintln(res.NoteApiNames())
	//lib.DPrintln(res.sections)
}
