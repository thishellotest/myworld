package tests

import (
	"context"
	"testing"
	"vbc/lib"
)

func Test_ZohoCollaboratorUsecase_HandleClientCases(t *testing.T) {
	err := UT.ZohoCollaboratorUsecase.HandleClientCases(context.TODO())
	lib.DPrintln(err)
}

func Test_ZohoCollaboratorUsecase_BizHandleClientCases(t *testing.T) {
	err := UT.ZohoCollaboratorUsecase.BizHandleClientCases("6159272000014605020")
	lib.DPrintln(err)
}

func Test_ZohoCollaboratorUsecase_BizHandleClients(t *testing.T) {
	err := UT.ZohoCollaboratorUsecase.BizHandleClients("6159272000012737026")
	lib.DPrintln(err)
}

func Test_ZohoCollaboratorUsecase_HandleClients(t *testing.T) {
	err := UT.ZohoCollaboratorUsecase.HandleClients(context.TODO())
	lib.DPrintln(err)
}
