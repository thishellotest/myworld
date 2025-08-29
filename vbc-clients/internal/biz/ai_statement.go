package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
)

type AiStatementUsecase struct {
	log               *log.Helper
	CommonUsecase     *CommonUsecase
	conf              *conf.Data
	LongMapUsecase    *LongMapUsecase
	Awsclaude3Usecase *Awsclaude3Usecase
	TUsecase          *TUsecase
}

func NewAiStatementUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	LongMapUsecase *LongMapUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	TUsecase *TUsecase) *AiStatementUsecase {
	uc := &AiStatementUsecase{
		log:               log.NewHelper(logger),
		CommonUsecase:     CommonUsecase,
		conf:              conf,
		LongMapUsecase:    LongMapUsecase,
		Awsclaude3Usecase: Awsclaude3Usecase,
		TUsecase:          TUsecase,
	}

	return uc
}

func (c *AiStatementUsecase) SplitCaseDescription(description string) (conditions []StatementCondition, err error) {
	strArr := strings.Split(description, "\n")
	var res []string
	for k, _ := range strArr {
		t := strings.TrimSpace(strArr[k])
		if t != "" {
			res = append(res, t)
		}
	}
	isOk := false
	for _, v := range res {
		if v == "New:" {
			isOk = true
			continue
		}
		if isOk {
			v := strings.TrimSpace(v)
			temp := strings.Split(v, " - ")
			if len(temp) == 2 {
				conditions = append(conditions, StatementCondition{
					//OriginValue: v,
					FrontValue:  temp[0],
					BehindValue: temp[1],
				})
			}
		}
	}
	return
}

func (c *AiStatementUsecase) GenStatementByCondition(input AiTaskInputStatement, tCase *TData) (string, error) {

	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	systemConfig, err := c.LongMapUsecase.GetForString(LongMapKey_sck1)
	if err != nil {
		return "", err
	}
	lib.DPrintln(systemConfig)

	//tCase.CustomFields.TextValueByNameBasic()

	//c.Awsclaude3Usecase.AskV2()

	return "", nil
}
