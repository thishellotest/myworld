package biz

import (
	evtBus "github.com/asaskevich/EventBus"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

type GlobalEventBus struct {
	Bus evtBus.Bus
}

func NewGlobalEventBus() *GlobalEventBus {
	bus := evtBus.New()

	return &GlobalEventBus{
		Bus: bus,
	}
}

const (
	GlobalEventBus_AfterHandleCompleteBoxSign    = "AfterHandleCompleteBoxSign"
	GlobalEventBus_AfterHandleCompleteAmContract = "AfterHandleCompleteAmContract"
)

type GlobalEventBusBuzUsecase struct {
	log                      *log.Helper
	conf                     *conf.Data
	CommonUsecase            *CommonUsecase
	TUsecase                 *TUsecase
	DialpadbuzUsecase        *DialpadbuzUsecase
	GlobalEventBus           *GlobalEventBus
	CronTriggerCreateUsecase *CronTriggerCreateUsecase
	SendVa2122aUsecase       *SendVa2122aUsecase
	ClientEnvelopeUsecase    *ClientEnvelopeUsecase
	MiscUsecase              *MiscUsecase
}

func NewGlobalEventBusBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	DialpadbuzUsecase *DialpadbuzUsecase,
	GlobalEventBus *GlobalEventBus,
	CronTriggerCreateUsecase *CronTriggerCreateUsecase,
	SendVa2122aUsecase *SendVa2122aUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	MiscUsecase *MiscUsecase,
) *GlobalEventBusBuzUsecase {
	uc := &GlobalEventBusBuzUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		TUsecase:                 TUsecase,
		DialpadbuzUsecase:        DialpadbuzUsecase,
		GlobalEventBus:           GlobalEventBus,
		CronTriggerCreateUsecase: CronTriggerCreateUsecase,
		SendVa2122aUsecase:       SendVa2122aUsecase,
		ClientEnvelopeUsecase:    ClientEnvelopeUsecase,
		MiscUsecase:              MiscUsecase,
	}
	GlobalEventBus.Bus.Subscribe(GlobalEventBus_AfterHandleCompleteBoxSign, uc.AfterHandleCompleteBoxSign)
	GlobalEventBus.Bus.Subscribe(GlobalEventBus_AfterHandleCompleteAmContract, uc.AfterHandleCompleteAmContract)
	return uc
}

//
//func (c *GlobalEventBusBuzUsecase) RunHandleSeparateAmContract(caseId int32) error {
//
//	amContractBoxFieldId, err := c.ClientEnvelopeUsecase.AmContractBoxFileId(caseId)
//	if err != nil {
//		return err
//	}
//	if amContractBoxFieldId == "" {
//		return errors.New("HandleSeparateAmContract amContractBoxFieldId is empty")
//	}
//	err = c.SendVa2122aUsecase.HandleSeparateAmContract(caseId, amContractBoxFieldId)
//	if err != nil {
//		c.log.Error(err, " caseId: ", caseId)
//		return err
//	}
//	return nil
//}

func (c *GlobalEventBusBuzUsecase) AfterHandleCompleteAmContract(clientCaseGid string) {

	tClientCase, er := c.TUsecase.DataByGid(Kind_client_cases, clientCaseGid)
	if er != nil {
		c.log.Error(er)
	} else if tClientCase != nil {

		usaPhone, er := FormatUSAPhoneHandle(tClientCase.CustomFields.TextValueByNameBasic(FieldName_phone))
		c.log.Info("AfterHandleCompleteBoxSign:FormatUSAPhoneHandle usaPhone:", usaPhone)
		if er != nil {
			c.log.Error(er, "AfterHandleCompleteBoxSign:FormatUSAPhoneHandle ID:", tClientCase.Id())
		} else if usaPhone != "" {
			er = c.CronTriggerCreateUsecase.CreateAfterSignedContract(tClientCase)
			if er != nil {
				c.log.Error(er, "AfterHandleCompleteBoxSign:HandleAfterSignedContract ID:", tClientCase.Id())
			}
		}
	}
	lib.DPrintln("AfterHandleCompleteBoxSign:clientCaseGid:", clientCaseGid)
}

func (c *GlobalEventBusBuzUsecase) AfterHandleCompleteBoxSign(clientCaseGid string) {

	tClientCase, er := c.TUsecase.DataByGid(Kind_client_cases, clientCaseGid)
	if er != nil {
		c.log.Error(er)
	} else if tClientCase != nil {
		usaPhone, er := FormatUSAPhoneHandle(tClientCase.CustomFields.TextValueByNameBasic(FieldName_phone))
		c.log.Info("AfterHandleCompleteBoxSign:FormatUSAPhoneHandle usaPhone:", usaPhone)
		if er != nil {
			c.log.Error(er, "AfterHandleCompleteBoxSign:FormatUSAPhoneHandle ID:", tClientCase.Id())
		} else if usaPhone != "" {

			er = c.CronTriggerCreateUsecase.CreateAfterSignedContract(tClientCase)
			//er := c.DialpadbuzUsecase.HandleAfterSignedContract(usaPhone)

			if er != nil {
				c.log.Error(er, "AfterHandleCompleteBoxSign:HandleAfterSignedContract ID:", tClientCase.Id())
			}
		}
	}

	lib.DPrintln("AfterHandleCompleteBoxSign:clientCaseGid:", clientCaseGid)

}
