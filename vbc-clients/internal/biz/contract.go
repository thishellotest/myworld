package biz

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	ActiveDuty_Yes = "Yes"
	ActiveDuty_No  = "No"
)

type ClientCaseContractBasicDataVo struct {
	ActiveDuty             bool  // 用于生成invoice
	EffectiveCurrentRating int32 // 用于生成invoice
	CurrentRating          int32 // 目前不使用此值，备份
}

func CreateClientCaseContractBasicDataVoByCase(tClientCaseFields TFields) (ClientCaseContractBasicDataVo, error) {
	if tClientCaseFields == nil {
		return ClientCaseContractBasicDataVo{}, errors.New("CreateClientCaseContractBasicDataVoByCase: tClientCaseFields is nil.")
	}
	clientCaseContractBasicDataVo := CreateClientCaseContractBasicDataVo(tClientCaseFields.TextValueByNameBasic(FieldName_active_duty),
		tClientCaseFields.NumberValueByNameBasic(FieldName_effective_current_rating),
		tClientCaseFields.NumberValueByNameBasic(FieldName_current_rating))
	return clientCaseContractBasicDataVo, nil
}

func CreateClientCaseContractBasicDataVo(activeDuty string, EffectiveCurrentRating int32, CurrentRating int32) ClientCaseContractBasicDataVo {
	vo := ClientCaseContractBasicDataVo{}
	if activeDuty == ActiveDuty_Yes {
		vo.ActiveDuty = true
	}
	vo.CurrentRating = CurrentRating
	vo.EffectiveCurrentRating = EffectiveCurrentRating
	return vo
}

type ClientCaseContractBasicDataUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	TUsecase      *TUsecase
	MapUsecase    *MapUsecase
}

func NewClientCaseContractBasicDataUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase) *ClientCaseContractBasicDataUsecase {
	uc := &ClientCaseContractBasicDataUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
		MapUsecase:    MapUsecase,
	}

	return uc
}

func (c *ClientCaseContractBasicDataUsecase) HttpHandleHistory(ctx *gin.Context) {
	reply := CreateReply()
	err := c.BizHttpHandleHistory(ctx.Query("do"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *ClientCaseContractBasicDataUsecase) BizHttpHandleHistory(do string) error {
	if do != "1" {
		return errors.New("params is errors.")
	}

	tList, err := c.TUsecase.ListByCond(Kind_client_cases, Eq{FieldName_biz_deleted_at: 0})
	if err != nil {
		return err
	}
	for k, _ := range tList {
		fields := tList[k].CustomFields
		clientCaseContractBasicDataKey := fmt.Sprintf("%s%d", Map_ClientCaseContractBasicData, fields.NumberValueByNameBasic("id"))
		val, err := c.MapUsecase.GetForString(clientCaseContractBasicDataKey)
		if err != nil {
			return err
		}
		if val == "" {
			clientCaseContractBasicDataVo, err := CreateClientCaseContractBasicDataVoByCase(fields)
			if err != nil {
				return err
			}
			err = c.MapUsecase.Set(clientCaseContractBasicDataKey, InterfaceToString(clientCaseContractBasicDataVo))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
