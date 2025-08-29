package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type SendsmsClientTasksConditionUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	ClientCaseUsecase *ClientCaseUsecase
	ClientTaskUsecase *ClientTaskUsecase
}

func NewSendsmsClientTasksConditionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ClientCaseUsecase *ClientCaseUsecase,
	ClientTaskUsecase *ClientTaskUsecase) *SendsmsClientTasksConditionUsecase {
	uc := &SendsmsClientTasksConditionUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		ClientCaseUsecase: ClientCaseUsecase,
		ClientTaskUsecase: ClientTaskUsecase,
	}

	return uc
}

func (c *SendsmsClientTasksConditionUsecase) WhetherExistsTask(caseGid string, stages string) (bool, error) {
	lists, err := c.ClientTaskUsecase.TasksByCaseGid(caseGid)
	if err != nil {
		c.log.Error(err)
		return false, err
	}
	//stages := tCase.CustomFields.TextValueByNameBasic(FieldName_stages)
	for _, v := range lists {
		lib.DPrintln(v.CustomFields.TextValueByNameBasic("subject"), stages)
		if config_vbc.JudgeTaskWhetherBelongsStage(v.CustomFields.TextValueByNameBasic("subject"), stages) {
			return true, nil
		}
	}
	return false, nil
}
