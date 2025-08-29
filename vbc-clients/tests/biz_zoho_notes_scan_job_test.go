package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ZohoNoteScanJobUsecase_BizRunJob(t *testing.T) {
	err := UT.ZohoNoteScanJobUsecase.BizRunJob()
	lib.DPrintln(err)
}

func Test_ZohoNoteScanJobUsecase_BatchHandle(t *testing.T) {
	var lastModifiedTime string
	err := UT.ZohoNoteScanJobUsecase.BatchHandle(&lastModifiedTime, "", 1)
	if err != nil {
		lib.DPrintln(err)
	}
}

func Test_ZohoNoteScanJobUsecase_SyncAll(t *testing.T) {
	err := UT.ZohoNoteScanJobUsecase.SyncAll("633f9b19a278e0cc844e2bb2610496561491aca5d6820ae37ae18213323e50f8ff8b479009ddd22531254abf869de8ba9480931e3a57443370edbce50446cfdb5d14035b5204197e435c2e068da3e9f80882047943f06fd8c6b6e4d2d47b125da6382db46aff7e440426b4d91b679b4944b937c578f1c64f60de2669258ac226a652887bfc775983f58bae462916cf79b979694f8d0df8f9b5fedad0ac0c1788", 1)
	lib.DPrintln(err)
}
