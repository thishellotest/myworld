package tests

import (
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_BehaviorUsecase_Add(t *testing.T) {
	UT.BehaviorUsecase.Add(0, "aa", time.Now(), "")
}

func Test_BehaviorUsecase_Add_1(t *testing.T) {
	UT.BehaviorUsecase.Add(44, biz.BehaviorType_complete_fee_schedule_contract, time.Now(), "")
}

func Test_BehaviorUsecase_MedicalTeamFormsContractSentAt(t *testing.T) {

	//UT.BehaviorUsecase.Add(5004, biz.BehaviorType_sent_medical_team_forms_contract, time.Now(), "")

	a := UT.BehaviorUsecase.MedicalTeamFormsContractSentAt(5004)
	lib.DPrintln(a)
}

//	func Test_BehaviorUsecase_Add_2(t *testing.T) {
//		UT.BehaviorUsecase.Add(5052, biz.BehaviorType_complete_release_of_information_contract, time.Now(), "")
//	}
//
//	func Test_BehaviorUsecase_Add_3(t *testing.T) {
//		UT.BehaviorUsecase.Add(5052, biz.BehaviorType_complete_patient_payment_form_contract, time.Now(), "")
//	}
func Test_BehaviorUsecase_HandleCompleteBoxSign(t *testing.T) {
	err := UT.BehaviorUsecase.HandleCompleteBoxSign(biz.BehaviorType_complete_fee_schedule_contract, 5460)
	if err != nil {
		panic(err)
	}
}

func Test_BehaviorUsecase_BehaviorForCreateInvoice(t *testing.T) {
	er := UT.BehaviorUsecase.BehaviorForCreateInvoice(5511, time.Now(), "")
	lib.DPrintln(er)
}
