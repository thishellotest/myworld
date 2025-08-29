package tests

import (
	"os"
	"testing"
	"time"
	"vbc/internal/biz"
)

func Test_GopdfUsecase_CreateContractAm2(t *testing.T) {

	caseId := int32(5814)
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, caseId)
	if err != nil {
		return
	}
	if tCase == nil {
		return
	}

	tClient, _, err := UT.DataComboUsecase.ClientWithCase(*tCase)
	if err != nil {
		return
	}
	if tClient == nil {
		return
	}
	contractVetVo := biz.GenContractVetVo(*tClient, *tCase)

	attorney, err := UT.AttorneyUsecase.GetByGid(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_attorney_uniqid))
	if err != nil {
		panic(err)
		return
	}
	if attorney == nil {
		return
	}
	contractAttorneyVo := attorney.ToContractAttorneyVo()
	contractTime := time.Now()
	signFileBytes, err := UT.GopdfUsecase.CreateContractAm(contractTime, contractVetVo, contractAttorneyVo)
	if err != nil {
		panic(err)
	}
	a, err := os.Create("a1.pdf")
	if err != nil {
		panic(err)
	}
	a.Write(signFileBytes)
	a.Close()
}

func Test_GopdfUsecase_CreatePersonalStatementsPDFForAiV1(t *testing.T) {
	caseGid := "6159272000001066012"
	const statementConditionId int32 = 4

	pdfBytes, fileName, err := UT.ExportUsecase.BizHttpStatement(caseGid, statementConditionId)
	if err != nil {
		t.Fatalf("BizHttpAllStatements failed: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("PDF bytes are empty")
	}

	err = os.MkdirAll(".data.local", 0755)
	if err != nil {
		t.Fatalf("Failed to create .data.local directory: %v", err)
	}

	filePath := ".data.local/" + fileName
	err = os.WriteFile(filePath, pdfBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write PDF to %s: %v", filePath, err)
	}

	t.Logf("Successfully saved PDF (%d bytes) to %s", len(pdfBytes), filePath)
}
