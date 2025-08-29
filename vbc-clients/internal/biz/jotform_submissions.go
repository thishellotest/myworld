package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

const JotformIntakeFormID = "242466899584074"
const JotformAmIntakeFormID = "251865711410149"

type JotformSubmissionEntity struct {
	ID           int32 `gorm:"primaryKey"`
	SubmissionId string
	FormId       string
	Uniqcode     string
	Notes        string
	CreatedAt    int64
	UpdatedAt    int64
}

func (JotformSubmissionEntity) TableName() string {
	return "jotform_submissions"
}

func (c *JotformSubmissionEntity) JotformNewFileNameForBox(tCase *TData) (newFileName string, err error) {
	return GenJotformNewFileNameForBox(c, tCase)
}

func (c *JotformSubmissionEntity) JotformNewFileNameForAI() (newFileName string, err error) {
	return GenJotformNewFileNameForAI(c)
}

type JotformSubmissionUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[JotformSubmissionEntity]
}

func NewJotformSubmissionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *JotformSubmissionUsecase {
	uc := &JotformSubmissionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *JotformSubmissionUsecase) GetLatestIntakeFormInfoByFormId(uniqcodes []string) (*JotformSubmissionEntity, error) {
	return c.GetLatestFormInfoByFormId(QuestionnairesInitialIntake_FormId, uniqcodes)
}

func (c *JotformSubmissionUsecase) GetLatestFormInfoByFormId(formId string, uniqcodes []string) (*JotformSubmissionEntity, error) {
	return c.GetByCondWithOrderBy(And(Eq{"form_id": formId}, In("uniqcode", uniqcodes)), "id desc")
}

func (c *JotformSubmissionUsecase) GetLatestFormInfo(submissionId string) (*JotformSubmissionEntity, error) {
	return c.GetByCondWithOrderBy(Eq{"submission_id": submissionId}, "id desc")
}

func (c *JotformSubmissionUsecase) GetLatestFormInfoWithUniqcode(submissionId string) (*JotformSubmissionEntity, error) {
	return c.GetByCondWithOrderBy(Eq{"submission_id": submissionId}, "id desc")
}

func (c *JotformSubmissionUsecase) AllLatestByUniqcodeExceptIntake(caseUniqcodes []string) (r []*JotformSubmissionEntity, err error) {
	caseUniqcodeStr := "'" + strings.Join(caseUniqcodes, "','") + "'"
	sql := fmt.Sprintf(`SELECT s.*
FROM jotform_submissions s
INNER JOIN questionnaires q ON q.jotform_form_id = s.form_id
WHERE q.deleted_at = 0 
  AND s.uniqcode in (%s)
  AND s.form_id not in ('%s')
  AND q.is_intake = 0
  AND s.id = (
    SELECT MAX(s2.id)
    FROM jotform_submissions s2
    WHERE s2.submission_id = s.submission_id
  ) order by s.id 
  `, caseUniqcodeStr, QuestionnairesUpdateQuestionnaire_FormId)

	c.log.Debug("AllLatestByUniqcodeExceptIntake:sql", strings.ReplaceAll(sql, "\n", ""))

	return c.AllByRawSql(sql)
}

func (c *JotformSubmissionUsecase) AllLatestByUniqcode(caseUniqcode string) ([]*JotformSubmissionEntity, error) {
	sql := fmt.Sprintf(`SELECT s.*
FROM jotform_submissions s
INNER JOIN questionnaires q ON q.jotform_form_id = s.form_id
WHERE q.deleted_at = 0 
  AND s.uniqcode = '%s'
  AND s.id = (
    SELECT MAX(s2.id)
    FROM jotform_submissions s2
    WHERE s2.submission_id = s.submission_id
  ) order by s.id 
  `, caseUniqcode)

	return c.AllByRawSql(sql)
}

func (c *JotformSubmissionUsecase) AllLatestUpdateQuestionnaires(caseUniqcode string) ([]*JotformSubmissionEntity, error) {

	return c.AllLatestByUniqcodeAndFormId(caseUniqcode, QuestionnairesUpdateQuestionnaire_FormId)
}

func (c *JotformSubmissionUsecase) AllLatestByUniqcodeAndFormId(caseUniqcode string, formId string) ([]*JotformSubmissionEntity, error) {
	sql := fmt.Sprintf(`SELECT s.*
FROM jotform_submissions s
INNER JOIN questionnaires q ON q.jotform_form_id = s.form_id
WHERE q.deleted_at = 0 
  AND s.uniqcode = '%s'
  AND s.form_id = '%s'
  AND s.id = (
    SELECT MAX(s2.id)
    FROM jotform_submissions s2
    WHERE s2.submission_id = s.submission_id
  ) order by s.id 
  `, caseUniqcode, formId)

	return c.AllByRawSql(sql)
}

func (c *JotformSubmissionUsecase) AllByUniqcodeAndConditionUniqid(conditionUniqid string) ([]*JotformSubmissionEntity, error) {

	sql := fmt.Sprintf(`SELECT s.*
FROM jotform_submissions s
INNER JOIN questionnaires q ON q.jotform_form_id = s.form_id
inner join relas_log on relas_log.target_id=s.submission_id and relas_log.type='%s' and source_id='%s' and relas_log.deleted_at=0
WHERE q.deleted_at = 0 
  AND q.is_intake = 0
  AND s.id = (
    SELECT MAX(s2.id)
    FROM jotform_submissions s2
    WHERE s2.submission_id = s.submission_id
  ) order by s.id 
  `, RelasLog_Type_condition_2_jotform_new, conditionUniqid)

	return c.AllByRawSql(sql)
}

func (c *JotformSubmissionUsecase) ManualHandleFormId() error {

	res, err := c.AllByCond(Eq{"form_id": ""})
	if err != nil {
		return err
	}
	for k, v := range res {
		if v.FormId == "" {
			notesMap := lib.ToTypeMapByString(v.Notes)
			formId, _, _, _, err := GetSubmissionInfo(notesMap)
			if err != nil {
				return err
			}
			if formId != "" {
				res[k].FormId = formId
				err := c.CommonUsecase.DB().Save(&res[k]).Error
				if err != nil {
					return err
				}
				c.log.Info("Id:", v.ID)
				//break
			}

		}
	}

	return nil
}
