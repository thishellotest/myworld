package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type StatementConditionBuzUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	StatementConditionUsecase *StatementConditionUsecase
	TUsecase                  *TUsecase
	DataEntryUsecase          *DataEntryUsecase
}

func NewStatementConditionBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
) *StatementConditionBuzUsecase {
	uc := &StatementConditionBuzUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		StatementConditionUsecase: StatementConditionUsecase,
		TUsecase:                  TUsecase,
		DataEntryUsecase:          DataEntryUsecase,
	}

	return uc
}
func (c *StatementConditionBuzUsecase) DoInitStatementCondition_Deleted(tCase TData) error {

	entity, err := c.StatementConditionUsecase.GetByCond(Eq{"case_id": tCase.Id()})
	if err != nil {
		return err
	}
	if entity != nil {
		return nil
	}
	return c.InitStatementCondition_Deleted(tCase)
}

func (c *StatementConditionBuzUsecase) InitStatementCondition_Deleted(tCase TData) error {

	caseId := tCase.Id()
	statements := tCase.CustomFields.TextValueByNameBasic(FieldName_statements)

	conditions, err := SplitCaseStatements(statements)
	if err != nil {
		return err
	}
	sort := 1000
	for k, v := range conditions {
		destSort := sort + k
		_, err = c.StatementConditionUsecase.Upsert(caseId, v.ConditionValue, v.FrontValue, v.ConditionValue, v.BehindValue, destSort, v.Category)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *StatementConditionBuzUsecase) UpdateCaseStatement(caseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	statement, err := c.GenerateStatementString(tCase.Id())
	if err != nil {
		return err
	}
	if statement != "" {
		data := make(TypeDataEntry)
		data[DataEntry_gid] = tCase.Gid()
		data[FieldName_statements] = statement
		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, data, DataEntry_gid, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateStatementString 前端也有同样的方法，
func (c *StatementConditionBuzUsecase) GenerateStatementString(caseId int32) (statement string, err error) {

	res, err := c.StatementConditionUsecase.AllConditions(caseId)
	if err != nil {
		return "", err
	}
	hasSupplemental := false
	hasNOPRIVATEEXAMS := false
	for _, v := range res {
		if v.Category == StatementCondition_Category_Supplemental {
			if hasSupplemental == false {
				hasSupplemental = true
				statement += "-------Supplemental-------\n\n"
			}
		}
		if v.Category == StatementCondition_Category_NOPRIVATEEXAMS {
			if hasNOPRIVATEEXAMS == false {
				hasNOPRIVATEEXAMS = true
				statement += "-------NO PRIVATE EXAMS-------\n\n"
			}
		}
		statement += v.ToCondition() + "\n\n"
	}
	return strings.TrimSpace(statement), nil
}
