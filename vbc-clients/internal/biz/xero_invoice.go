package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type XeroInvoiceUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	MapUsecase       *MapUsecase
	XeroUsecase      *XeroUsecase
	TUsecase         *TUsecase
	LogUsecase       *LogUsecase
	DataComboUsecase *DataComboUsecase
	BehaviorUsecase  *BehaviorUsecase
}

func NewXeroInvoiceUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	XeroUsecase *XeroUsecase,
	TUsecase *TUsecase,
	LogUsecase *LogUsecase,
	DataComboUsecase *DataComboUsecase,
	BehaviorUsecase *BehaviorUsecase,
) *XeroInvoiceUsecase {
	uc := &XeroInvoiceUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		MapUsecase:       MapUsecase,
		XeroUsecase:      XeroUsecase,
		TUsecase:         TUsecase,
		LogUsecase:       LogUsecase,
		DataComboUsecase: DataComboUsecase,
		BehaviorUsecase:  BehaviorUsecase,
	}

	return uc
}

func (c *XeroInvoiceUsecase) HandleInvoice(clientCaseId int32) error {
	key := fmt.Sprintf("%s%d", Map_XeroInvoiceId, clientCaseId)
	invoiceId, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if invoiceId != "" { // 已经生成帐单号
		lib.DPrintln("HandleInvoice: already exists invoice.", invoiceId)
		return nil
	}
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}

	if tClientCase == nil {
		return errors.New("tClientCase is nil.")
	}
	clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")

	_, tContactField, err := c.DataComboUsecase.Client(clientGid)
	if err != nil {
		return err
	}
	if tContactField == nil {
		return errors.New("tContactField is nil.")
	}

	email := tContactField.TextValueByNameBasic(FieldName_email)
	if email == "" {
		lib.DPrintln("HandleInvoice email is empty")
		return nil
	}
	firstName := tContactField.TextValueByNameBasic(FieldName_first_name)
	if firstName == "" {
		return errors.New("firstName is empty.")
	}
	if tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_AwaitingPayment {
		lib.DPrintln("HandleInvoice FieldName_stages:", tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages))
		return nil
	}
	if tClientCase.CustomFields.NumberValueByNameBasic(FieldName_new_rating) <= 0 {
		lib.DPrintln("HandleInvoice FieldName_new_rating:", tClientCase.CustomFields.TextValueByNameBasic(FieldName_new_rating))
		return nil
	}

	InvoiceID, InvoiceNumber, err := c.XeroUsecase.BizCreateInvoice(tClientCase)
	if err != nil {
		return err
	}
	er := c.LogUsecase.SaveLog(clientCaseId, Log_FromType_Xero_CreateInvoice, map[string]interface{}{
		"InvoiceID":     InvoiceID,
		"InvoiceNumber": InvoiceNumber,
	})
	if er != nil {
		c.log.Error(er)
	}
	er = c.BehaviorUsecase.BehaviorForCreateInvoice(clientCaseId, time.Now(), "")
	if er != nil {
		c.log.Error(er)
	}
	return c.MapUsecase.Set(key, InvoiceID)
}

func (c *XeroInvoiceUsecase) HandleAmInvoice(clientCaseId int32) error {
	key := fmt.Sprintf("%s%d", Map_XeroInvoiceId, clientCaseId)
	invoiceId, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return err
	}
	if invoiceId != "" { // 已经生成帐单号
		lib.DPrintln("HandleInvoice: already exists invoice.", invoiceId)
		return nil
	}
	tClientCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}

	if tClientCase == nil {
		return errors.New("tClientCase is nil.")
	}
	clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")

	_, tContactField, err := c.DataComboUsecase.Client(clientGid)
	if err != nil {
		return err
	}
	if tContactField == nil {
		return errors.New("tContactField is nil.")
	}

	email := tContactField.TextValueByNameBasic(FieldName_email)
	if email == "" {
		lib.DPrintln("HandleInvoice email is empty")
		return nil
	}
	firstName := tContactField.TextValueByNameBasic(FieldName_first_name)
	if firstName == "" {
		return errors.New("firstName is empty.")
	}
	if tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages) != config_vbc.Stages_AmAwaitingPayment {
		lib.DPrintln("HandleInvoice FieldName_stages:", tClientCase.CustomFields.TextValueByNameBasic(FieldName_stages))
		return nil
	}
	//if tClientCase.CustomFields.NumberValueByNameBasic(FieldName_new_rating) <= 0 {
	//	lib.DPrintln("HandleInvoice FieldName_new_rating:", tClientCase.CustomFields.TextValueByNameBasic(FieldName_new_rating))
	//	return nil
	//}

	InvoiceID, InvoiceNumber, err := c.XeroUsecase.BizAmCreateInvoice(tClientCase)
	if err != nil {
		return err
	}
	er := c.LogUsecase.SaveLog(clientCaseId, Log_FromType_Xero_AmCreateInvoice, map[string]interface{}{
		"InvoiceID":     InvoiceID,
		"InvoiceNumber": InvoiceNumber,
	})
	if er != nil {
		c.log.Error(er)
	}
	er = c.BehaviorUsecase.BehaviorForAmCreateInvoice(clientCaseId, time.Now(), "")
	if er != nil {
		c.log.Error(er)
	}
	return c.MapUsecase.Set(key, InvoiceID)
}
