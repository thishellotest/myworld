package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

const (
	Option_VS_Color = "#FEA36A"
	Option_CP_Color = "#5CB3FD"
)

type OptionUsecase struct {
	log                      *log.Helper
	conf                     *conf.Data
	CommonUsecase            *CommonUsecase
	QuestionnairesUsecase    *QuestionnairesUsecase
	JotformSubmissionUsecase *JotformSubmissionUsecase
	FeeUsecase               *FeeUsecase
}

func NewOptionUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	QuestionnairesUsecase *QuestionnairesUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	FeeUsecase *FeeUsecase,
) *OptionUsecase {
	uc := &OptionUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		QuestionnairesUsecase:    QuestionnairesUsecase,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		FeeUsecase:               FeeUsecase,
	}

	return uc
}

func (c *OptionUsecase) NewVersionJotformIdsOptions(tCase *TData, keyword string) (r []FabFieldOption, err error) {

	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return nil, err
	}
	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	if !isPrimaryCase {
		if primaryCase == nil {
			return nil, errors.New("primaryCase is nil")
		}
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}
	uniqcodeStr := "'" + strings.Join(uniqcodes, "','") + "'"

	sql := fmt.Sprintf(`SELECT s.*
FROM jotform_submissions s
INNER JOIN questionnaires q ON q.jotform_form_id = s.form_id
WHERE q.deleted_at = 0 
  AND s.uniqcode in (%s) 
  AND q.is_intake = 0
  AND s.id = (
    SELECT MAX(s2.id)
    FROM jotform_submissions s2
    WHERE s2.submission_id = s.submission_id
  ) order by s.id `, uniqcodeStr)
	//c.log.Debug("NewVersionJotformIdsOptions sql: ", sql)
	res, err := c.JotformSubmissionUsecase.AllByRawSql(sql)
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		var newFileName string
		if v.Uniqcode == uniqcode {
			newFileName, err = v.JotformNewFileNameForBox(tCase)
		} else {
			newFileName, err = v.JotformNewFileNameForBox(primaryCase)
		}
		if err != nil {
			return nil, err
		}

		if keyword != "" {
			if strings.Index(strings.ToLower(newFileName), strings.ToLower(keyword)) < 0 {
				continue
			}
		}

		fabFieldOption := FabFieldOption{
			OptionLabel: newFileName,
			OptionValue: v.SubmissionId,
		}
		r = append(r, fabFieldOption)
	}
	return r, nil
}

func (c *OptionUsecase) JotformIdsOptions(tCase *TData, valueType string, keyword string) (r []FabFieldOption, err error) {

	if configs.NewPSGen {
		return c.NewVersionJotformIdsOptions(tCase, keyword)
	}

	var items []*QuestionnairesEntity
	if valueType == "1" {
		if tCase != nil {
			sql := fmt.Sprintf("select * from questionnaires where jotform_form_id in (select q.jotform_form_id from questionnaires q inner join jotform_submissions j  on j.form_id=q.`jotform_form_id` and j.uniqcode=\"%s\" and q.is_intake=0 and q.deleted_at=0) order by base_title", tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
			items, err = c.QuestionnairesUsecase.AllByRawSql(sql)
			if err != nil {
				return nil, err
			}
		}
	} else {
		items, err = c.QuestionnairesUsecase.AllByCondWithOrderBy(Eq{"is_intake": 0, "deleted_at": 0}, "base_title asc", 50)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range items {
		fabFieldOption := FabFieldOption{
			OptionLabel: v.BaseTitle,
			OptionValue: v.JotformFormId,
		}
		r = append(r, fabFieldOption)
	}
	return
}
