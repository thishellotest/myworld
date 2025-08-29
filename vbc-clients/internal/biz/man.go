package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"os"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
MailFeeScheduleCommunication:5250:gengling.liao@hotmail.com
CreateEnvelope:5250:gengling.liao@hotmail.com
*/

type ManUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	TUsecase      *TUsecase
	MapUsecase    *MapUsecase
}

func NewManUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase) *ManUsecase {
	uc := &ManUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
		MapUsecase:    MapUsecase,
	}

	return uc
}

func (c *ManUsecase) HandleHistoryCreateEnvelope() error {
	return nil
	res, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{"biz_deleted_at": 0, "deleted_at": 0})
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fmt.Sprintf("/tmp/aaa_%d.log", time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, v := range res {

		caseId := v.Id()
		email := v.CustomFields.TextValueByNameBasic("email")
		if email != "" {
			lib.DPrintln("ok:", caseId)
			feeScheduleCommunication := fmt.Sprintf("%s%d", Map_mail_FeeScheduleCommunication, caseId)
			createEnvelope := fmt.Sprintf("%s%d", Map_CreateEnvelope, caseId)
			newFeeScheduleCommunication := fmt.Sprintf("%s%d:%s", Map_mail_FeeScheduleCommunication, caseId, email)
			newCreateEnvelope := fmt.Sprintf("%s%d:%s", Map_CreateEnvelope, caseId, email)
			a, er := c.MapUsecase.GetForString(feeScheduleCommunication)
			if er != nil {
				panic(er)
			}
			if a != "" {
				c.MapUsecase.Set(newFeeScheduleCommunication, "1")
			}

			a, er = c.MapUsecase.GetForString(createEnvelope)
			if er != nil {
				panic(er)
			}
			if a != "" {
				c.MapUsecase.Set(newCreateEnvelope, "1")
			}
		} else {
			io.WriteString(f, email+"\n")
		}
		//break
	}
	return nil
}
