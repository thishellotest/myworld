package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
)

type FeeUsecase struct {
	log                   *log.Helper
	CommonUsecase         *CommonUsecase
	conf                  *conf.Data
	ClientCaseUsecase     *ClientCaseUsecase
	PricingVersionUsecase *PricingVersionUsecase
	TUsecase              *TUsecase
	RatingPaymentUsecase  *RatingPaymentUsecase
}

func NewFeeUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ClientCaseUsecase *ClientCaseUsecase,
	PricingVersionUsecase *PricingVersionUsecase,
	TUsecase *TUsecase,
	RatingPaymentUsecase *RatingPaymentUsecase,
) *FeeUsecase {
	uc := &FeeUsecase{
		log:                   log.NewHelper(logger),
		CommonUsecase:         CommonUsecase,
		conf:                  conf,
		ClientCaseUsecase:     ClientCaseUsecase,
		PricingVersionUsecase: PricingVersionUsecase,
		TUsecase:              TUsecase,
		RatingPaymentUsecase:  RatingPaymentUsecase,
	}
	return uc
}

func (c *FeeUsecase) GetIncreaseAmount(currentRating int, newRating int) (increaseAmount int, err error) {

	payments, err := c.RatingPaymentUsecase.CurrentRatingPayments()
	if err != nil {
		return 0, err
	}
	var current *RatingPaymentEntity
	var final *RatingPaymentEntity
	for k, v := range payments {
		if v.Rating == currentRating {
			current = payments[k]
		}
		if v.Rating == newRating {
			final = payments[k]
		}
	}
	if current == nil || final == nil {
		return 0, errors.New("RatingPaymentEntity is nil")
	}
	increaseAmount = final.Payment - current.Payment
	return increaseAmount, nil
}

func (c *FeeUsecase) FeeScheduleCommunicationSubId(tClientCase *TData) (subId int, err error) {
	if tClientCase == nil {
		return 0, errors.New("FeeScheduleCommunicationSubId: tClientCase is nil.")
	}
	fieldData := tClientCase.CustomFields
	if fieldData.TextValueByNameBasic("active_duty") == "Yes" {
		return -1, nil
	}
	return config_vbc.MailSubIdV1(fieldData.NumberValueByNameBasic("effective_current_rating")), nil
}

// VBCFees VBC Fees
func (c *FeeUsecase) VBCFees(tClientCase *TData) ([]config_vbc.FeeVo, error) {

	if tClientCase == nil {
		return nil, errors.New("VbcFees: tClientCase is nil.")
	}
	fieldData := tClientCase.CustomFields
	var index int
	if fieldData.TextValueByNameBasic("active_duty") == "Yes" {
		index = -1
	} else {
		effectiveCurrentRating := fieldData.NumberValueByNameBasic("effective_current_rating")
		index = config_vbc.SignContractIndexV1(effectiveCurrentRating)
	}
	if configs.EnabledDBPricingVersion {
		pricingVersion, err := c.ClientCaseUsecase.GetPricingVersion(tClientCase)
		if err != nil {
			return nil, err
		}
		CurrentVersionConfig, _, err := c.PricingVersionUsecase.ConfigByPricingVersion(pricingVersion)
		if err != nil {
			return nil, err
		}
		val := CurrentVersionConfig.GetByIndex(index)
		c.log.Debug("EnabledDBPricingVersion: ", val)
		return val, nil
	}
	return config_vbc.FeeDefine.GetByIndex(index), nil
}

// NotPrimaryCaseAmount 计算价格
func (c *FeeUsecase) NotPrimaryCaseAmount(tClientCase *TData, newRating int32) (amount int, err error) {
	if tClientCase == nil {
		return 0, errors.New("NotPrimaryCaseAmount: tClientCase is nil")
	}
	if newRating <= 0 {
		return 0, errors.New("NotPrimaryCaseAmount: newRating is zero")
	}
	clientGid := tClientCase.CustomFields.TextValueByNameBasic("client_gid")
	clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
	primaryCase, err := c.ClientCaseUsecase.PrimaryCase(clientGid)

	if err != nil {
		return 0, err
	}
	if primaryCase == nil { // 计算invoice
		return 0, errors.New("NotPrimaryCaseAmount: primaryCase is nil")
	}

	clientCaseContractBasicDataVo, err := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(primaryCase.CustomFields.NumberValueByNameBasic("id"))
	if err != nil {
		return 0, err
	}
	if clientCaseContractBasicDataVo == nil {
		return 0, errors.New("NotPrimaryCaseAmount: clientCaseContractBasicDataVo is nil")
	}
	c.log.Info("NotPrimaryCaseAmount clientCaseContractBasicDataVo:", clientCaseContractBasicDataVo)
	originAmount, err := c.InvoiceAmount(primaryCase.CustomFields.NumberValueByNameBasic("id"), clientCaseContractBasicDataVo.ActiveDuty, clientCaseContractBasicDataVo.EffectiveCurrentRating, newRating)
	if err != nil {
		return 0, err
	}
	if originAmount <= 0 {
		return 0, errors.New(fmt.Sprintf("NotPrimaryCaseAmount: originAmount is %d", originAmount))
	}

	originAmountDecimal := decimal.NewFromInt(int64(originAmount))
	c.log.Info("NotPrimaryCaseAmount originAmountDecimal:", originAmountDecimal.String())

	primaryAmountStr := primaryCase.CustomFields.TextValueByNameBasic(FieldName_amount)
	primaryAmountDecimal, _ := decimal.NewFromString(primaryAmountStr)
	//zeroDecimal := decimal.NewFromInt(0)
	//if primaryAmountDecimal.LessThanOrEqual(zeroDecimal) {
	//	return 0, errors.New("NotPrimaryCaseAmount: primaryAmountDecimal is " + primaryAmountDecimal.String())
	//}
	originAmountDecimal = originAmountDecimal.Sub(primaryAmountDecimal)

	c.log.Info("NotPrimaryCaseAmount primaryAmountDecimal:", primaryAmountDecimal.String())

	siblingCases, err := c.ClientCaseUsecase.NotPrimaryCases(clientGid, clientCaseGid)
	if err != nil {
		return 0, err
	}
	for k, _ := range siblingCases {
		siblingAmountStr := siblingCases[k].CustomFields.TextValueByNameBasic(FieldName_amount)
		siblingAmountDecimal, _ := decimal.NewFromString(siblingAmountStr)
		//if siblingAmountDecimal.LessThanOrEqual(zeroDecimal) {
		//	return 0, errors.New("NotPrimaryCaseAmount: siblingAmountDecimal is " + siblingAmountDecimal.String())
		//}
		c.log.Info("NotPrimaryCaseAmount siblingAmountDecimal:", siblingAmountDecimal.String())
		originAmountDecimal = originAmountDecimal.Sub(siblingAmountDecimal)
	}
	amount = int(originAmountDecimal.Floor().IntPart())
	if amount <= 0 {
		return 0, errors.New("NotPrimaryCaseAmount: amount is " + originAmountDecimal.String())
	}
	return amount, nil
}

// UsePrimaryCaseCalc 是否为primary case计算方法，当isPrimaryCase=false时，一定有primaryCase
func (c *FeeUsecase) UsePrimaryCaseCalc(tClientCase *TData) (isPrimaryCase bool, primaryCase *TData, err error) {
	if tClientCase == nil {
		return false, nil, errors.New("UsePrimaryCaseCalc: tClientCase is nil")
	}
	fieldData := tClientCase.CustomFields
	if fieldData.NumberValueByNameBasic(FieldName_is_primary_case) == Is_primary_case_YES {
		return true, nil, nil
	}
	clientGid := fieldData.TextValueByNameBasic("client_gid")
	//clientCaseGid := fieldData.TextValueByNameBasic("gid")
	primaryCase, err = c.ClientCaseUsecase.PrimaryCase(clientGid)
	if err != nil {
		return false, nil, err
	}
	if primaryCase == nil {
		return true, nil, nil
	}
	return false, primaryCase, nil
}

func HasEnabledPrimaryCase(clientGid string) bool {
	return true

}

// ClientCaseAmount 此方法只能预估使用
func (c *FeeUsecase) ClientCaseAmount(tClientCase *TData) (amount int, err error, noCaseContractBasicDataVo bool) {
	if tClientCase == nil {
		return 0, errors.New("ClientCaseAmount: tClientCase is nil."), false
	}

	// New version
	fieldData := tClientCase.CustomFields
	newRating := fieldData.NumberValueByNameBasic("new_rating")
	if newRating <= 0 {
		newRating = 100
	}

	if HasEnabledPrimaryCase(fieldData.TextValueByNameBasic("client_gid")) {
		usePrimaryCaseCalc, _, err := c.UsePrimaryCaseCalc(tClientCase)
		if err != nil {
			return 0, err, false
		}
		if !usePrimaryCaseCalc {
			amount, err = c.NotPrimaryCaseAmount(tClientCase, newRating)
			return amount, err, false
		}
	}
	clientCaseContractBasicDataVo, _ := c.ClientCaseUsecase.ClientCaseContractBasicDataVoById(fieldData.NumberValueByNameBasic("id"))
	var activeDuty bool
	var effectiveCurrentRating int32
	if clientCaseContractBasicDataVo != nil {
		//c.log.Debug("ClientCaseAmount: clientCaseContractBasicDataVo is not nil", InterfaceToString(clientCaseContractBasicDataVo))
		activeDuty = clientCaseContractBasicDataVo.ActiveDuty
		effectiveCurrentRating = clientCaseContractBasicDataVo.EffectiveCurrentRating
	} else {
		return 0, nil, true
		if fieldData.TextValueByNameBasic("active_duty") == "Yes" {
			activeDuty = true
		}
		effectiveCurrentRating = fieldData.NumberValueByNameBasic(FieldName_effective_current_rating)
		//c.log.Debug("ClientCaseAmount: clientCaseContractBasicDataVo is nil", InterfaceToString(activeDuty), InterfaceToString(effectiveCurrentRating))
	}
	amount, err = c.InvoiceAmount(tClientCase.CustomFields.NumberValueByNameBasic("id"), activeDuty, effectiveCurrentRating, newRating)
	return amount, err, false
	// End new version

	/*
		feeVos, err := c.VBCFees(tClientCase)
		if err != nil {
			return 0, err
		}
		newRating := tClientCase.CustomFields.NumberValueByNameBasic("new_rating")
		if newRating <= 0 {
			newRating = 100
		}

		for _, v := range feeVos {
			if newRating >= int32(v.Rating) {
				amount = v.Fee
			}
		}
		if amount > 0 {
			return amount, nil
		}

		//  此处创建客户真实收费有风险，
		return 0, nil
	*/
}

// InvoiceAmount 计算
func (c *FeeUsecase) InvoiceAmount(caseId int32, activeDuty bool, effectiveCurrentRating int32, newRating int32) (amount int, err error) {

	var index int
	if activeDuty {
		index = -1
	} else {
		index = config_vbc.SignContractIndexV1(effectiveCurrentRating)
	}
	var fees []config_vbc.FeeVo
	if configs.EnabledDBPricingVersion {
		tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
		if err != nil {
			return 0, err
		}
		if tCase == nil {
			return 0, errors.New(InterfaceToString(caseId) + ": tCase is nil")
		}

		pricingVersion, err := c.ClientCaseUsecase.GetPricingVersion(tCase)
		if err != nil {
			return 0, err
		}
		CurrentVersionConfig, _, err := c.PricingVersionUsecase.ConfigByPricingVersion(pricingVersion)
		if err != nil {
			c.log.Error(err)
			return 0, err
		}
		if CurrentVersionConfig == nil {
			return 0, errors.New("CurrentVersionConfig is nil")
		}
		fees = CurrentVersionConfig.GetByIndex(index)
	} else {
		fees = config_vbc.FeeDefine.GetByIndex(index)
	}

	for _, v := range fees {
		if newRating >= int32(v.Rating) {
			amount = v.Fee
		}
	}
	if amount > 0 {
		return amount, nil
	}
	return 0, nil
}

//func (c FeeVoConfigs) Charge(currentRating int, newRating int) int {
//	if _, ok := c[currentRating]; ok {
//		for _, v := range c[currentRating] {
//			if v.Rating == newRating {
//				return v.Fee
//			}
//		}
//	}
//	return 0
//}
