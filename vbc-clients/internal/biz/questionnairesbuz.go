package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type QuestionnairesbuzUsecase struct {
	log                      *log.Helper
	conf                     *conf.Data
	CommonUsecase            *CommonUsecase
	QuestionnairesUsecase    *QuestionnairesUsecase
	RelasLogUsecase          *RelasLogUsecase
	JotformSubmissionUsecase *JotformSubmissionUsecase
	FeeUsecase               *FeeUsecase
}

func NewQuestionnairesbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	QuestionnairesUsecase *QuestionnairesUsecase,
	RelasLogUsecase *RelasLogUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	FeeUsecase *FeeUsecase,
) *QuestionnairesbuzUsecase {
	uc := &QuestionnairesbuzUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		QuestionnairesUsecase:    QuestionnairesUsecase,
		RelasLogUsecase:          RelasLogUsecase,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		FeeUsecase:               FeeUsecase,
	}

	return uc
}

func (c *QuestionnairesbuzUsecase) GetJotformSubmissionsForGenStatementNew(tCase *TData, StatementCondition StatementCondition) (intakeSubmission *JotformSubmissionEntity, others []*JotformSubmissionEntity, err error) {

	if tCase == nil {
		return nil, nil, errors.New("tCase is nil")
	}
	//conditionUniqid := InterfaceToString(tCase.Id()) + ":" + StatementCondition.ConditionValue
	conditionUniqid := InterfaceToString(StatementCondition.StatementConditionId)
	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	submissions, err := c.JotformSubmissionUsecase.AllByUniqcodeAndConditionUniqid(conditionUniqid)
	if err != nil {
		return nil, nil, err
	}

	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return nil, nil, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	intakeSubmission, err = c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return nil, nil, err
	}
	if intakeSubmission == nil {
		return nil, nil, errors.New("intakeSubmission is nil")
	}
	return intakeSubmission, submissions, nil
}

func (c *QuestionnairesbuzUsecase) GetJotformSubmissionsForGenStatement(tCase *TData, conditionEntity *ConditionEntity) (intakeSubmission *JotformSubmissionEntity, others []*JotformSubmissionEntity, err error) {

	if tCase == nil {
		return nil, nil, errors.New("tCase is nil")
	}
	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	if conditionEntity == nil {
		return nil, nil, errors.New("conditionEntity is nil")
	}
	relas, err := c.RelasLogUsecase.ConditionRelas(InterfaceToString(conditionEntity.ID))
	if err != nil {
		return nil, nil, err
	}
	if len(relas) == 0 {
		return nil, nil, errors.New("Please associate Condition with Jotform")
	}
	var formIds []string
	for _, v := range relas {
		formIds = append(formIds, v.TargetId)
	}

	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return nil, nil, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	intakeSubmission, err = c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return nil, nil, err
	}
	if intakeSubmission == nil {
		return nil, nil, errors.New("intakeSubmission is nil")
	}
	sql := fmt.Sprintf(`select * from jotform_submissions where id in (
select max(id) as id from jotform_submissions
 where form_id in (%s) and uniqcode='%s' GROUP BY form_id
 )`, "'"+strings.Join(formIds, "','")+"'", uniqcode)

	submissions, err := c.JotformSubmissionUsecase.AllByRawSql(sql)
	if err != nil {
		return nil, nil, err
	}
	if len(submissions) != len(formIds) {
		return nil, nil, errors.New("The associated Jotform has no corresponding data")
	}
	return intakeSubmission, submissions, nil
}

func (c *QuestionnairesbuzUsecase) AllForCondition(conditionId int32) ([]*QuestionnairesEntity, error) {
	sql := fmt.Sprintf("select q.* from questionnaires q inner join relas_log l on l.target_id=q.jotform_form_id and l.type=\"%s\" and l.source_id=\"%d\" where l.deleted_at=0 and q.deleted_at=0 order by q.base_title ",
		RelasLog_Type_condition_2_jotform, conditionId)
	return c.QuestionnairesUsecase.AllByRawSql(sql)
}

func (c *QuestionnairesbuzUsecase) Manual() error {

	for _, v := range QuestionnairesListConfigs {
		a, err := c.QuestionnairesUsecase.GetByCond(Eq{"jotform_form_id": v.FormId})
		if err != nil {
			return err
		}
		if a == nil {
			a = &QuestionnairesEntity{
				JotformFormId: v.FormId,
				Title:         v.Title,
				BaseTitle:     v.BaseTitle,
				JsonData: InterfaceToString(map[string]interface{}{
					"file_names": v.FileNames,
				}),
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}
			err = c.CommonUsecase.DB().Save(&a).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
