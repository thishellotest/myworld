package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ContractReminderUsecase_FirstReminder(t *testing.T) {
	err := UT.ContractReminderUsecase.FirstReminder(5004)
	lib.DPrintln(err)
}

func Test_ContractReminderUsecase_AmIntakeFormReminderFirstReminder(t *testing.T) {
	err := UT.ContractReminderUsecase.AmIntakeFormReminderFirstReminder(5746)
	lib.DPrintln(err)
}

//
//func Test_ContractReminderUsecase_SecondReminder(t *testing.T) {
//	err := UT.ContractReminderUsecase.SecondReminder(5004)
//	lib.DPrintln(err)
//}
//
//func Test_ContractReminderUsecase_ThirdReminder(t *testing.T) {
//	err := UT.ContractReminderUsecase.ThirdReminder(5004)
//	lib.DPrintln(err)
//}
//
//func Test_ContractReminderUsecase_FourthReminder(t *testing.T) {
//	err := UT.ContractReminderUsecase.FourthReminder(5004)
//	lib.DPrintln(err)
//}
//
//func Test_ContractReminderUsecase_CreateNonResponsiveTask(t *testing.T) {
//	err := UT.ContractReminderUsecase.CreateNonResponsiveTask(5004)
//	lib.DPrintln(err)
//}
