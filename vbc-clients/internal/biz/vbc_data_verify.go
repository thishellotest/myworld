package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strings"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

type VbcDataVerifyUsecase struct {
	log               *log.Helper
	conf              *conf.Data
	CommonUsecase     *CommonUsecase
	TUsecase          *TUsecase
	FeeUsecase        *FeeUsecase
	ClientCaseUsecase *ClientCaseUsecase
	DataComboUsecase  *DataComboUsecase
	ZohobuzUsecase    *ZohobuzUsecase
}

func NewVbcDataVerifyUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	FeeUsecase *FeeUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	DataComboUsecase *DataComboUsecase,
	ZohobuzUsecase *ZohobuzUsecase,
) *VbcDataVerifyUsecase {
	uc := &VbcDataVerifyUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		TUsecase:          TUsecase,
		FeeUsecase:        FeeUsecase,
		ClientCaseUsecase: ClientCaseUsecase,
		DataComboUsecase:  DataComboUsecase,
		ZohobuzUsecase:    ZohobuzUsecase,
	}

	return uc
}

func (c *VbcDataVerifyUsecase) VerifyContract() error {
	cases, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_biz_deleted_at: 0})
	if err != nil {
		return err
	}

	for k, v := range cases {
		err := c.VerifyContractOne(cases[k])
		if err != nil {
			c.log.Debug("VerifyContractError:", err, " case: ", v.Gid())
		} else {
			c.log.Debug("VerifyContractOK:", err, " case: ", v.Gid())
		}
	}
	return nil
}

func (c *VbcDataVerifyUsecase) VerifyContractOne(v *TData) error {

	isPrimaryCase, _, err := c.FeeUsecase.UsePrimaryCaseCalc(v)
	if err != nil {
		return err
	}
	if isPrimaryCase {
		//key := MapKeyClientCaseContractBasicData(v.Id())

		stages := v.CustomFields.TextValueByNameBasic(FieldName_stages)

		if stages == config_vbc.Stages_IncomingRequest ||
			stages == config_vbc.Stages_Completed ||
			stages == config_vbc.Stages_Dormant ||
			stages == config_vbc.Stages_Terminated {
			return nil
		}

		contractVo, err := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(v.Id())
		if err != nil {
			return err
		}
		if contractVo == nil {
			return errors.New("contractVo is nil: " + InterfaceToString(v.Id()))
		}

		if strings.Index(v.CustomFields.TextValueByNameBasic(FieldName_deal_name), "Test") >= 0 {
			return nil
		}

		newAmount, err, noCaseContractBasicDataVo := c.FeeUsecase.ClientCaseAmount(v)
		if err != nil {
			return err
		}
		if noCaseContractBasicDataVo {
			return errors.New("noCaseContractBasicDataVo: " + InterfaceToString(v.Id()))
		}

		newAmountDecimal := decimal.NewFromInt(int64(newAmount))
		amountDecimal, _ := decimal.NewFromString(v.CustomFields.TextValueByNameBasic(FieldName_amount))
		if !amountDecimal.Equal(newAmountDecimal) {
			return errors.New("amountDecimal neq newAmountDecimal: " + v.CustomFields.TextValueByNameBasic(FieldName_stages) + " amountDecimal: " + amountDecimal.String() + " newAmountDecimal: " + newAmountDecimal.String() + " caseId: " + InterfaceToString(v.Id()))
		}

		dealName := v.CustomFields.TextValueByNameBasic(FieldName_deal_name)

		client, _, _ := c.DataComboUsecase.Client(v.CustomFields.TextValueByNameBasic(FieldName_client_gid))

		clientCaseName := ClientCaseNameByCase(client.CustomFields.TextValueByNameBasic("first_name"),
			client.CustomFields.TextValueByNameBasic("last_name"),
			*v,
		)
		if clientCaseName != dealName {

			er := c.ZohobuzUsecase.HandleClientCaseName(v.Id())
			if er != nil {
				c.log.Error(er, " id: ", v.Id())
			}
			return errors.New("clientCaseName clientCaseName: " + clientCaseName + " : " + dealName + " : " + InterfaceToString(v.Id()))
		}
		//if v.CustomFields.TextValueByNameBasic(FieldName_active_duty) == ActiveDuty_Yes && !contractVo.ActiveDuty {
		//	return errors.New("ActiveDuty different: " + InterfaceToString(v.Id()))
		//}
		if v.CustomFields.NumberValueByNameBasic(FieldName_effective_current_rating) != contractVo.EffectiveCurrentRating {
			return errors.New("EffectiveCurrentRating different: " + InterfaceToString(v.Id()))
		}
	}
	return nil
}
