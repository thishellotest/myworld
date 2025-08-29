package tests

import (
	"errors"
	"testing"
	"vbc/lib"
)

func Test_AsanaMigrateUsecase_BizHttpMigrateOne(t *testing.T) {
	err := UT.AsanaMigrateUsecase.BizHttpMigrateOne("1206911825266122")
	lib.DPrintln(err)
}

func GetAsanaRow(asanaClientGid string) (lib.TypeMap, error) {

	//asanaClientGid := "1206911825266122"
	_, list, err := UT.AsanaMigrateUsecase.Data(asanaClientGid)
	if err != nil {
		lib.DPrintln(err)
		return nil, err
	}
	if len(list) == 0 {
		lib.DPrintln("length is 0")
		return nil, errors.New("length is 0")
	}
	for _, v := range list {
		v := lib.TypeMap(v)
		return v, nil
	}
	return nil, nil
}

func Test_AsanaMigrateUsecase_HandleZohoDeal(t *testing.T) {

	asanaClientGid := "1206911825266122"
	asanaClientRow, err := GetAsanaRow(asanaClientGid)
	if err != nil {
		lib.DPrintln(err)
		return
	}

	clientId, clientGid, err := UT.AsanaMigrateUsecase.HandleZohoContact(asanaClientRow.GetString("email"), asanaClientRow)
	lib.DPrintln(clientId, clientGid, err)
	if err != nil {
		lib.DPrintln(err)
		return
	}
	lib.DPrintln(clientId, clientGid)

	clientCaseId, clientCaseGid, err := UT.AsanaMigrateUsecase.HandleZohoDeal(clientGid, asanaClientRow)
	lib.DPrintln(clientCaseId, clientCaseGid, err)
}

func Test_AsanaMigrateUsecase_HandleZohoContact(t *testing.T) {

	asanaClientGid := "1206911825266122"
	asanaClientRow, err := GetAsanaRow(asanaClientGid)
	if err != nil {
		lib.DPrintln(err)
		return
	}
	clientId, clientGid, err := UT.AsanaMigrateUsecase.HandleZohoContact(asanaClientRow.GetString("email"), asanaClientRow)
	lib.DPrintln(clientId, clientGid, err)
}

func Test_AsanaMigrateUsecase_HandleMaps(t *testing.T) {
	err := UT.AsanaMigrateUsecase.HandleMaps(12)
	lib.DPrintln(err)
}

func Test_AsanaMigrateUsecase_HandleBehaviors(t *testing.T) {
	err := UT.AsanaMigrateUsecase.HandleBehaviors(12)
	lib.DPrintln(err)
}
