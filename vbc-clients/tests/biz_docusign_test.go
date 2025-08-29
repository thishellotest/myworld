package tests

import (
	"context"
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/esign/v2.1/envelopes"
	"vbc/lib/esign/v2.1/model"
)

func Test_DocuSignUsecase_HandleEnvelopeChangeStatus(t *testing.T) {
	err := UT.DocuSignUsecase.HandleEnvelopeChangeStatus()
	lib.DPrintln(err)
}

func Test_Create_Envolope_from_a_template(t *testing.T) {

	token, err := UT.Oauth2TokenUsecase.GetByAppId(biz.Oauth2_AppId_docusign)
	fmt.Println(err)
	if token == nil {
	}
	cred := biz.NewDocuSignCredential(token, UT.Conf.Docusign.AppAccountId)
	srv := envelopes.New(cred)

	templateRole := model.TemplateRole{}
	templateRole.Name = "张三"
	templateRole.RoleName = "Client"
	templateRole.Email = "gengling.liao@hotmail.com"

	templateRole1 := model.TemplateRole{}
	templateRole1.Name = "Gary"
	templateRole1.RoleName = "Agent"
	templateRole1.Email = "liaogling@gmail.com"

	//EnvelopeDefinition
	envlopeDef := model.EnvelopeDefinition{}
	envlopeDef.TemplateRoles = append(envlopeDef.TemplateRoles, templateRole, templateRole1)
	envlopeDef.Status = "sent"
	envlopeDef.TemplateID = "e57d1c5a-650f-4446-b9d8-315d93350bbb"

	res, err := srv.Create(&envlopeDef).Do(context.TODO())
	lib.DPrintln(err)
	lib.DPrintln(res)
}

func Test_DocuSignUsecase_ContractTemplateIdByDB(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	a, err := UT.DocuSignUsecase.ContractTemplateIdOld(tCase)
	lib.DPrintln("1:", err)

	_, b, pricingVersion, err := UT.DocuSignUsecase.ContractTemplateIdByDB(tCase)
	lib.DPrintln("2:", err, pricingVersion)
	if a == b {
		lib.DPrintln(a)
		lib.DPrintln("ok")
	} else {
		panic("ss")
	}
}
