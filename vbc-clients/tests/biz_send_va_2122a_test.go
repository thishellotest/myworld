package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_SendVa2122aUsecase_HandleSeparateAmContract(t *testing.T) {

	caseId := int32(5842)
	amContractBoxFileId, err := UT.ClientEnvelopeUsecase.AmContractBoxFileId(caseId)
	if err != nil {
		panic(err)
	}
	if amContractBoxFileId == "" {
		panic("amContractBoxFileId is empty")
	}
	lib.DPrintln("Manual amContractBoxFileId: ", amContractBoxFileId)
	//return
	key := biz.MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId)
	val, _ := UT.MapUsecase.GetForString(key)
	if val != "" {
		_, err := UT.BoxUsecase.DeleteFile(val)
		if err != nil {
			lib.DPrintln(err)
		}
	}

	key1 := biz.MapKeyClientCaseAmSignedAgreementBoxFileId(caseId)
	val1, _ := UT.MapUsecase.GetForString(key1)
	if val1 != "" {
		_, err := UT.BoxUsecase.DeleteFile(val1)
		if err != nil {
			lib.DPrintln(err)
		}
	}
	err = UT.SendVa2122aUsecase.HandleSeparateAmContract(caseId, amContractBoxFileId)
	if err != nil {
		panic(err)
	}

	_, err = UT.MiscUsecase.Delete2122aFile(caseId)
	if err != nil {
		panic(err)
	}
	lib.DPrintln("Manual Delete2122aFile: ", err)
	_, err = UT.MiscUsecase.DoHandleMoving2122aFile(caseId)
	if err != nil {
		panic(err)
	}
	lib.DPrintln("Manual DoHandleMoving2122aFile: ", err)

}

func Test_SendVa2122aUsecase_DoHandleMoving2122aFile(t *testing.T) {
	caseId := int32(5885)
	_, err := UT.MiscUsecase.Delete2122aFile(caseId)
	lib.DPrintln("Manual Delete2122aFile: ", err)
	_, err = UT.MiscUsecase.DoHandleMoving2122aFile(caseId)
	lib.DPrintln("Manual DoHandleMoving2122aFile: ", err)
}

/*
select id,deal_name, `attorney_uniqid` from client_cases where biz_deleted_at=0 and deleted_at=0 and contract_source='AM' and stages not in ('Am__Terminated', 'Am__Dormant', "Am__Information Intake", "Am__Incoming Request", "Am__Contract Pending") and deal_name not like '%Test%'
*/
func Test_SendVa2122aUsecase_RunHandleSeparateAmContract(t *testing.T) {
	err := UT.SendVa2122aUsecase.RunHandleSeparateAmContract(5781)
	lib.DPrintln(err)
}

func Test_SendVa2122aUsecase_Download21Pdf(t *testing.T) {
	err := UT.SendVa2122aUsecase.Download21Pdf()
	if err != nil {
		panic(err)
	}
}
