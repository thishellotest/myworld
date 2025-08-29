package biz

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
	. "vbc/lib/builder"
)

type StatemtEntity struct {
	ID                               int32 `gorm:"primaryKey"`
	CaseId                           int32
	ConditionUniqid                  string
	FrontValue                       string
	ConditionValue                   string
	BehindValue                      string
	Category                         string
	YearsOfService                   string
	RetiredFromService               string
	Deployments                      string
	MaritalStatus                    string
	Children                         string
	OccupationInService              string
	CurrentTreatmentFacility         string
	CurrentMedication                string
	SpecialNotes                     string // 可能为空由AI决定，如果有AI会返回：类似这部分 SERVICE CONNECTION: My service treatment records reflect that I suffered from Chronic Neck and Low back Pain while on active duty. All legal requirements for establishing service connection for Neck pain with right upper extremity radiculopathy have been met; service connection for such disease is warranted.
	IntroductionParagraph            string // 类似I am respectfully requesting Veteran Affairs benefits for my condition of neck pain with right upper extremity radiculopathy. This condition has been a significant challenge in my life since its onset during my active-duty service in 2008. Throughout my service in the Air Force as a Special Operations Pilot from 2007 to 2018, I faced circumstances that have contributed to this condition. The effects of my neck pain have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
	OnsetAndServiceConnection        string
	CurrentSymptomsSeverityFrequency string
	Medication                       string
	ImpactOnDailyLife                string
	ProfessionalImpact               string
	NexusBetweenSC                   string
	Request                          string
	Versions                         int32 // 数字越大版本越新
	ModifiedBy                       string
	CreatedAt                        int64
	UpdatedAt                        int64
	DeletedAt                        int64
}

func (StatemtEntity) TableName() string {
	return "statements"
}

type StatemtUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[StatemtEntity]
}

func NewStatemtUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *StatemtUsecase {
	uc := &StatemtUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type ListStatemtEntity []*StatemtEntity

func (c ListStatemtEntity) GetByConditionUniqid(conditionUniqid string) *StatemtEntity {
	for k, v := range c {
		if v.ConditionUniqid == conditionUniqid {
			return c[k]
		}
	}
	return nil
}

func (c *StatemtUsecase) AllStatementsByVersion(caseId int32, version int32) (ListStatemtEntity, error) {
	sql := fmt.Sprintf(`SELECT t1.*
FROM statements t1
JOIN (
    SELECT condition_uniqid, MAX(versions) AS max_versions
    FROM statements
    WHERE case_id=%d and deleted_at=0 and versions=%d
    GROUP BY condition_uniqid
) t2
ON t1.condition_uniqid = t2.condition_uniqid AND t1.versions = t2.max_versions
where t1.case_id=%d and t1.deleted_at=0 
`, caseId, version, caseId)
	return c.AllByRawSql(sql)
}

func (c *StatemtUsecase) AllLatestStatements(caseId int32) (ListStatemtEntity, error) {
	sql := fmt.Sprintf(`SELECT t1.*
FROM statements t1
JOIN (
    SELECT condition_uniqid, MAX(versions) AS max_versions
    FROM statements
    WHERE case_id=%d and deleted_at=0
    GROUP BY condition_uniqid
) t2
ON t1.condition_uniqid = t2.condition_uniqid AND t1.versions = t2.max_versions
where t1.case_id=%d and t1.deleted_at=0
`, caseId, caseId)
	return c.AllByRawSql(sql)
}

type StatemtVersionVo struct {
	Version   int32 `json:"version"`
	CreatedAt int32 `json:"created_at"`
}

func (c *StatemtUsecase) Versions(caseId int32) (versions []StatemtVersionVo, err error) {
	sql := fmt.Sprintf(`select versions,min(created_at) as  created_at from statements where case_id=%d group by versions order by versions desc limit 50`, caseId)
	res, err := c.AllByRawSql(sql)
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		vo := StatemtVersionVo{
			Version:   v.Versions,
			CreatedAt: int32(v.CreatedAt),
		}
		versions = append(versions, vo)
	}
	return
}

func (c *StatemtUsecase) GetLatestStatement(caseId int32, conditionUniqid string) (*StatemtEntity, error) {
	return c.GetByCondWithOrderBy(Eq{"case_id": caseId,
		"condition_uniqid": conditionUniqid,
		"deleted_at":       0}, "versions desc")
}

func (c *StatemtUsecase) ObtainAvailableVersionID(caseId int32) (versions int32, err error) {

	statementEntity, err := c.GetByCondWithOrderBy(Eq{"case_id": caseId,
		"deleted_at": 0}, "versions desc")
	if err != nil {
		return 0, err
	}
	if statementEntity == nil {
		return 1, nil
	}
	return statementEntity.Versions + 1, nil
}

type StatementConditionListApi []StatementConditionApi

type StatementConditionList []StatementCondition

const (
	StatementCondition_Category_General        = "General"          // 这个唯一值不能改
	StatementCondition_Category_Supplemental   = "Supplemental"     // 这个唯一值不能改
	StatementCondition_Category_NOPRIVATEEXAMS = "NO PRIVATE EXAMS" // 这个唯一值不能改
)

var StatementConditionCategoryOrder = []string{StatementCondition_Category_General, StatementCondition_Category_Supplemental, StatementCondition_Category_NOPRIVATEEXAMS}

type StatementCondition struct {
	StatementConditionId int32 `json:"id"`
	//StatementConditionUuid string `json:"uuid"`
	Sort int `json:"sort"`
	//OriginValue    string `json:"origin_value"`
	ConditionValue string `json:"condition_value"` // 此值忽略大小写，去数据库匹配，新建等
	FrontValue     string `json:"front_value"`
	BehindValue    string `json:"behind_value"`
	Category       string `json:"category"`
}

type StatementConditionApi struct {
	StatementConditionId string `json:"id"`
	//StatementConditionUuid string `json:"uuid"`
	Sort int `json:"sort"`
	//OriginValue    string `json:"origin_value"`
	ConditionValue string `json:"condition"` // 此值忽略大小写，去数据库匹配，新建等
	FrontValue     string `json:"rating"`
	BehindValue    string `json:"association"`
	Category       string `json:"category"`
}

func (c *StatementCondition) ToOriginValue() (val string) {
	if c.FrontValue != "" {
		val += c.FrontValue + " - "
	}
	val += c.ConditionValue
	if c.BehindValue != "" {
		val += " " + c.BehindValue
	}
	return val
}

//func (c *StatementCondition) ConditionUuid(caseId int32) string {
//	return fmt.Sprintf("%d:%s", caseId, c.ConditionValue)
//}

func SplitCaseStatements(statements string) (conditions StatementConditionList, err error) {
	strArr := strings.Split(statements, "\n")
	var res []string
	for k, _ := range strArr {
		t := strings.TrimSpace(strArr[k])
		if t != "" {
			res = append(res, t)
		}
	}
	category := StatementCondition_Category_General
	for _, v := range res {
		a := strings.TrimSpace(v)

		b := strings.ToLower(a)
		if strings.Index(b, "supplemental") >= 0 {
			category = StatementCondition_Category_Supplemental
			continue
		} else if strings.Index(b, "no private exams") >= 0 {
			category = StatementCondition_Category_NOPRIVATEEXAMS
			continue
		}

		// 匹配 "-" 之前的内容
		re1 := regexp.MustCompile(`^(.*?)\s*-`)
		// re1 := regexp.MustCompile(`^(.*?)\s*[-–]`) 暂不兼容此–
		match1 := re1.FindStringSubmatch(a)
		if len(match1) > 1 {
			//fmt.Println("匹配到的前半部分:", match1[1])
		} else {
			return nil, errors.New(fmt.Sprintf("The format of \"%s\" is incorrect ", v))
		}
		a = strings.ReplaceAll(a, match1[0], "")
		frontValue := match1[0]

		// 匹配最后一个 ")" 之前的 "(" 到最后一个 ")" 之间的内容

		last := strings.LastIndex(a, ")")
		var condition string
		var behindValue string
		if last >= 0 && last == (len(a)-1) {
			//re2 := regexp.MustCompile(`\(([^()]*)\)`) // 匹配所有括号内的内容
			re2 := regexp.MustCompile(`\(([^()]*)\)[^()]*$`)
			match2 := re2.FindAllStringSubmatch(a, -1)
			if len(match2) > 0 {
				match3 := match2[len(match2)-1]
				condition = strings.ReplaceAll(a, match3[0], "")
				behindValue = match3[0]
			}
		} else {
			condition = a
		}
		statementCondition := StatementCondition{
			//OriginValue:    v,
			ConditionValue: strings.TrimSpace(condition),
			FrontValue:     strings.ReplaceAll(strings.TrimSpace(frontValue), " -", ""),
			BehindValue:    strings.TrimSpace(strings.Trim(strings.TrimSpace(behindValue), "()")),
			Category:       category,
		}
		conditions = append(conditions, statementCondition)
	}
	return
}

type StatementUsecase struct {
	log                          *log.Helper
	conf                         *conf.Data
	CommonUsecase                *CommonUsecase
	RelasLogUsecase              *RelasLogUsecase
	ConditionUsecase             *ConditionUsecase
	JotformSubmissionUsecase     *JotformSubmissionUsecase
	AiTaskUsecase                *AiTaskUsecase
	WordUsecase                  *WordUsecase
	AiResultUsecase              *AiResultUsecase
	BoxbuzUsecase                *BoxbuzUsecase
	BoxUsecase                   *BoxUsecase
	FeeUsecase                   *FeeUsecase
	TUsecase                     *TUsecase
	StatemtUsecase               *StatemtUsecase
	DataComboUsecase             *DataComboUsecase
	MapUsecase                   *MapUsecase
	StatementConditionUsecase    *StatementConditionUsecase
	StatementConditionBuzUsecase *StatementConditionBuzUsecase
	DocEmailUsecase              *DocEmailUsecase
	AiAssistantJobUsecase        *AiAssistantJobUsecase
	PersonalWebformUsecase       *PersonalWebformUsecase
	RecordbuzUsecase             *RecordbuzUsecase
}

func NewStatementUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	RelasLogUsecase *RelasLogUsecase,
	ConditionUsecase *ConditionUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	AiTaskUsecase *AiTaskUsecase,
	WordUsecase *WordUsecase,
	AiResultUsecase *AiResultUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	FeeUsecase *FeeUsecase,
	TUsecase *TUsecase,
	StatemtUsecase *StatemtUsecase,
	DataComboUsecase *DataComboUsecase,
	MapUsecase *MapUsecase,
	StatementConditionUsecase *StatementConditionUsecase,
	StatementConditionBuzUsecase *StatementConditionBuzUsecase,
	DocEmailUsecase *DocEmailUsecase,
	AiAssistantJobUsecase *AiAssistantJobUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
	RecordbuzUsecase *RecordbuzUsecase,
) *StatementUsecase {
	uc := &StatementUsecase{
		log:                          log.NewHelper(logger),
		CommonUsecase:                CommonUsecase,
		conf:                         conf,
		RelasLogUsecase:              RelasLogUsecase,
		ConditionUsecase:             ConditionUsecase,
		JotformSubmissionUsecase:     JotformSubmissionUsecase,
		AiTaskUsecase:                AiTaskUsecase,
		WordUsecase:                  WordUsecase,
		AiResultUsecase:              AiResultUsecase,
		BoxbuzUsecase:                BoxbuzUsecase,
		BoxUsecase:                   BoxUsecase,
		FeeUsecase:                   FeeUsecase,
		TUsecase:                     TUsecase,
		StatemtUsecase:               StatemtUsecase,
		DataComboUsecase:             DataComboUsecase,
		MapUsecase:                   MapUsecase,
		StatementConditionUsecase:    StatementConditionUsecase,
		StatementConditionBuzUsecase: StatementConditionBuzUsecase,
		DocEmailUsecase:              DocEmailUsecase,
		AiAssistantJobUsecase:        AiAssistantJobUsecase,
		PersonalWebformUsecase:       PersonalWebformUsecase,
		RecordbuzUsecase:             RecordbuzUsecase,
	}

	return uc
}

func (c *StatementUsecase) ManualSyncStatement(caseGid string) error {

	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	caseId := tCase.Id()
	hasInit, err := c.HasInitStatementsEdit(caseId)
	if err != nil {
		return err
	}
	if hasInit {
		c.log.Info("ManualSyncStatement caseId: " + InterfaceToString(caseId) + " " + tCase.Gid() + " already had initialized")
		return nil
	}
	statements := tCase.CustomFields.TextValueByNameBasic(FieldName_statements)
	statementConditionList, err := SplitCaseStatements(statements)
	if err != nil {
		return err
	}
	var statementConditionListApi StatementConditionListApi
	for k, v := range statementConditionList {
		statementConditionListApi = append(statementConditionListApi, StatementConditionApi{
			StatementConditionId: InterfaceToString(v.StatementConditionId),
			Sort:                 1000 + k,
			ConditionValue:       v.ConditionValue,
			FrontValue:           v.FrontValue,
			BehindValue:          v.BehindValue,
			Category:             v.Category,
		})
	}
	if len(statementConditionListApi) <= 0 {
		c.log.Info("statementConditionListApi is empty")
		return nil
	}

	res, err := json.Marshal(&statementConditionListApi)
	if err != nil {
		return err
	}
	newStatement, err := c.SaveCaseStatement(caseGid, string(res))
	if err != nil {
		return err
	}
	lib.DPrintln(newStatement)
	return nil
}

func (c *StatementUsecase) SaveCaseStatement(caseGid string, newStatements string) (statements string, err error) {

	var statementConditionList StatementConditionListApi
	err = json.Unmarshal([]byte(newStatements), &statementConditionList)
	if err != nil {
		return "", err
	}
	tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if err != nil {
		return "", err
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	caseId := tCase.Id()
	hasInit, err := c.HasInitStatementsEdit(caseId)
	if err != nil {
		return "", err
	}
	var needReserveIds []int32
	//lib.DPrintln(statementConditionList)
	//return nil
	if !hasInit {
		for k, v := range statementConditionList {
			entity, err := c.StatementConditionUsecase.GetByCond(Eq{"case_id": caseId, "uuid": v.ConditionValue, "deleted_at": 0})
			if err != nil {
				return "", err
			}
			if entity == nil {
				entity = &StatementConditionEntity{
					CaseId:    caseId,
					Uuid:      v.ConditionValue,
					CreatedAt: time.Now().Unix(),
				}
			}
			entity.FrontValue = v.FrontValue
			entity.ConditionValue = v.ConditionValue
			entity.BehindValue = v.BehindValue
			entity.Category = v.Category
			entity.Sort = 1000 + k
			entity.UpdatedAt = time.Now().Unix()
			err = c.CommonUsecase.DB().Save(&entity).Error
			if err != nil {
				return "", err
			}
			needReserveIds = append(needReserveIds, entity.ID)
		}
		err = c.CaseInitStatementsEdit(caseId)
		if err != nil {
			return "", err
		}
	} else {
		for k, v := range statementConditionList {
			var entity *StatementConditionEntity
			isNew := false
			if strings.Index(v.StatementConditionId, "new_") == 0 {
				isNew = true
			} else {
				entity, err = c.StatementConditionUsecase.GetByCond(Eq{"case_id": caseId, "id": v.StatementConditionId, "deleted_at": 0})
				if err != nil {
					return "", err
				}
				if entity == nil {
					return "", errors.New("The StatementCondition was not found")
				}
			}

			if isNew {
				entity = &StatementConditionEntity{
					CaseId:    caseId,
					Uuid:      v.ConditionValue,
					CreatedAt: time.Now().Unix(),
				}
			}
			if entity != nil {
				entity.FrontValue = v.FrontValue
				entity.ConditionValue = v.ConditionValue
				entity.BehindValue = v.BehindValue
				entity.Category = v.Category
				entity.Sort = 1000 + k
				entity.UpdatedAt = time.Now().Unix()
				err = c.CommonUsecase.DB().Save(&entity).Error
				if err != nil {
					return "", err
				}
				needReserveIds = append(needReserveIds, entity.ID)
			}
		}
	}

	err = c.StatementConditionUsecase.UpdatesByCond(map[string]interface{}{
		"deleted_at": time.Now().Unix(),
	}, And(Eq{"deleted_at": 0, "case_id": caseId}, NotIn("id", needReserveIds)))

	if err != nil {
		return "", err
	}

	text, err := c.GetCaseStatementByStatementCondition(caseId)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (c *StatementUsecase) GetCaseStatementByStatementCondition(caseId int32) (string, error) {
	records, err := c.StatementConditionUsecase.AllConditions(caseId)
	if err != nil {
		return "", err
	}

	text := ""
	flag := make(map[string]bool)
	for _, v := range records {
		if v.Category != StatementCondition_Category_General {
			if _, ok := flag[v.Category]; !ok {
				flag[v.Category] = true
				text += "-------" + v.Category + "-------\n\n"
			}
		}
		text += v.ToCondition() + "\n\n"
	}
	text = strings.TrimSpace(text)

	return text, nil
}

func (c *StatementUsecase) GetCaseStatementExtend(tCase TData) (conditions StatementConditionList, err error) {

	hasInitStatementsEdit, err := c.HasInitStatementsEdit(tCase.Id())
	if err != nil {
		return nil, err
	}
	if hasInitStatementsEdit {
		records, err := c.StatementConditionUsecase.AllConditions(tCase.Id())
		if err != nil {
			return nil, err
		}
		for _, v := range records {
			conditions = append(conditions, v.ToStatementCondition())
		}
	} else {
		statements := tCase.CustomFields.TextValueByNameBasic(FieldName_statements)
		conditions, err = SplitCaseStatements(statements)
	}

	return conditions, err
}

func (c *StatementUsecase) HasInitStatementsEdit(caseId int32) (hasInitStatementsEdit bool, err error) {

	key := MapKeyHasInitStatementsEdit(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val != "" {
		return true, nil
	}
	return false, nil
}

func (c *StatementUsecase) CaseInitStatementsEdit(caseId int32) error {
	key := MapKeyHasInitStatementsEdit(caseId)
	return c.MapUsecase.Set(key, "1")
}

func (c *StatementUsecase) GetJotFormInfo(statementCondition StatementCondition) error {

	conditionEntity, err := c.ConditionUsecase.ConditionUpsert(statementCondition.ConditionValue)
	if err != nil {
		return err
	}
	if conditionEntity == nil {
		return errors.New("conditionEntity is nil")
	}
	relasLog, err := c.RelasLogUsecase.ConditionRelas(InterfaceToString(conditionEntity.ID))
	if err != nil {
		return err
	}
	lib.DPrintln(relasLog)

	return nil
}

type StatementBaseInfoList []StatementBaseInfo
type StatementBaseInfo struct {
	Label string
	Value string
}

func (c *StatementUsecase) HandleStatementToBox(tCase *TData, tClient *TData, veteranSummary VeteranSummaryVo) error {

	if tCase == nil {
		return errors.New("tCase is nil")
	}
	if tClient == nil {
		return errors.New("tClient is nil")
	}

	existEntity, err := c.AiTaskUsecase.GetByCond(Eq{"deleted_at": 0,
		"from_type":     AiTaskFromType_statement,
		"case_id":       tCase.Id(),
		"handle_status": 0,
	})
	if err != nil {
		return err
	}
	if existEntity != nil { // 还有任务没有完成不允许生成文档
		return nil
	}

	er := c.GenerateNewStatementVersion(*tCase, *tClient)
	if er != nil {
		c.log.Error("GenerateNewStatementVersion:", er, " caseId: ", tCase.Id())
	}

	return c.GenerateDocument(*tCase, *tClient)
}

// statementDetail StatementDetail

func (c *StatementUsecase) UpdateOneConditionStatement(tCase TData, tClient TData, statementConditionEntity StatementConditionEntity, parseAiStatementConditionVo ParseAiStatementConditionVo) error {
	statementDetail, err := c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return err
	}

	isAiTaskOk := false
	for k, v := range statementDetail.Statements {

		if v.StatementCondition.StatementConditionId == statementConditionEntity.ID {

			//parseAiStatementConditionVo := ParseAiStatementCondition(ParseResult)

			for k1, v1 := range statementDetail.Statements[k].Rows {
				if v1.SectionType == Statemt_Section_CurrentTreatmentFacility {
					if parseAiStatementConditionVo.CurrentTreatmentFacility != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentTreatmentFacility
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_CurrentMedication {
					if parseAiStatementConditionVo.CurrentMedication != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentMedication
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_SpecialNotes {
					if parseAiStatementConditionVo.SpecialNotes != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.SpecialNotes
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_IntroductionParagraph {
					if parseAiStatementConditionVo.IntroductionParagraph != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.IntroductionParagraph
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_OnsetAndServiceConnection {
					if parseAiStatementConditionVo.OnsetAndServiceConnection != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.OnsetAndServiceConnection
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_CurrentSymptomsSeverityFrequency {
					if parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_Medication {
					if parseAiStatementConditionVo.Medication != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Medication
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_ImpactOnDailyLife {
					if parseAiStatementConditionVo.ImpactOnDailyLife != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ImpactOnDailyLife
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_ProfessionalImpact {
					if parseAiStatementConditionVo.ProfessionalImpact != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ProfessionalImpact
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_NexusBetweenSC {
					if parseAiStatementConditionVo.NexusBetweenSC != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.NexusBetweenSC
						isAiTaskOk = true
					}
				} else if v1.SectionType == Statemt_Section_Request {
					if parseAiStatementConditionVo.Request != "" {
						statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Request
						isAiTaskOk = true
					}
				}
			}

		}
	}
	raws, err := json.Marshal(&statementDetail)
	if err != nil {
		return err
	}
	if isAiTaskOk {
		_, err = c.BizStatementSave(false, nil, tCase.Gid(), raws)
	}
	return err
}

func (c *StatementUsecase) UpdateOneSectionStatement(tCase TData, tClient TData, statementConditionEntity StatementConditionEntity, sectionType string, sectionStatement string) error {
	statementDetail, err := c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return err
	}

	for k, v := range statementDetail.Statements {

		if v.StatementCondition.StatementConditionId == statementConditionEntity.ID {
			for k1, v1 := range statementDetail.Statements[k].Rows {
				if v1.SectionType == sectionType {
					if statementDetail.Statements[k].Rows[k1].Body != sectionStatement {
						statementDetail.Statements[k].Rows[k1].Body = sectionStatement
					}
				}
			}
		}
	}
	raws, err := json.Marshal(&statementDetail)
	if err != nil {
		return err
	}
	_, err = c.BizStatementSave(false, nil, tCase.Gid(), raws)
	return err
}

func (c *StatementUsecase) GetVeteranSummaryForVeteranSummary(tCase TData) (isOk bool, isToPsform bool, veteranSummary VeteranSummaryVo, aiTaskId int32) {
	aiTaskEntity, _ := c.AiTaskUsecase.GetVeteranSummary(&tCase)
	if aiTaskEntity == nil {
		return false, false, veteranSummary, 0
	}

	aiResult, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTaskEntity.CurrentResultId})
	if err != nil {
		c.log.Error(err)
		return false, false, veteranSummary, 0
	}
	if aiResult == nil {
		c.log.Error("aiResult is nil")
		return false, false, veteranSummary, 0
	}

	str := GetJsonFromAiResultForAssistant(aiResult.ParseResult)
	veteranSummary = VeteranSummaryJsonToVo(str)

	if aiTaskEntity.ToPsform == AiTask_ToPsform_Yes {
		isToPsform = true
	}
	return true, isToPsform, veteranSummary, aiTaskEntity.ID
}

func (c *StatementUsecase) GetNewStatementForWebForm(tCase TData, tClient TData) (statementDetail StatementDetail, useAiTaskIds []int32, useAssistantIds []int32, err error) {
	statementDetail, err = c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return statementDetail, useAiTaskIds, useAssistantIds, err
	}
	isOk, isToPsform, veteranSummary, veteranSummaryAiTaskId := c.GetVeteranSummaryForVeteranSummary(tCase)

	if isOk && !isToPsform {
		useAiTaskIds = append(useAiTaskIds, veteranSummaryAiTaskId)
		if veteranSummary.YearsOfService != "" {
			statementDetail.BaseInfo.YearsOfService = veteranSummary.YearsOfService
		}
		if veteranSummary.RetirementStatus != "" {
			statementDetail.BaseInfo.RetiredFromService = veteranSummary.RetirementStatus
		}
		if veteranSummary.Deployments != "" {
			statementDetail.BaseInfo.Deployments = veteranSummary.Deployments
		}
		if veteranSummary.MaritalStatus != "" {
			statementDetail.BaseInfo.MaritalStatus = veteranSummary.MaritalStatus
		}
		if veteranSummary.Children != "" {
			statementDetail.BaseInfo.Children = veteranSummary.Children
		}
		if veteranSummary.OccupationInService != "" {
			statementDetail.BaseInfo.OccupationInService = veteranSummary.OccupationInService
		}
	}

	for k, v := range statementDetail.Statements {

		jobUuid := GenStatementConditionJobUuid(tCase.Id(), v.StatementCondition.StatementConditionId)
		aJobEntity, err := c.AiAssistantJobUsecase.GetByCond(Eq{"job_uuid": jobUuid})
		if err != nil {
			return statementDetail, useAiTaskIds, useAssistantIds, err
		}
		if aJobEntity == nil {
			continue
		}
		aiTaskEntity, err := c.AiTaskUsecase.GetByCond(Eq{"id": aJobEntity.AiTaskId})

		if err != nil {
			return statementDetail, useAiTaskIds, useAssistantIds, err
		}
		if aiTaskEntity == nil {
			continue
		}
		if aiTaskEntity.ToPsform == AiTask_ToPsform_Yes {
			continue
		}
		aiResultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTaskEntity.CurrentResultId})
		if err != nil {
			return statementDetail, useAiTaskIds, useAssistantIds, err
		}
		if aiResultEntity == nil {
			c.log.Error("aiResultEntity is nil")
			continue
		}

		parseAiStatementConditionVo := ParseAiStatementCondition(aiResultEntity.ParseResult)

		isAiTaskOk := false
		for k1, v1 := range statementDetail.Statements[k].Rows {
			if v1.SectionType == Statemt_Section_CurrentTreatmentFacility {
				if parseAiStatementConditionVo.CurrentTreatmentFacility != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentTreatmentFacility
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_CurrentMedication {
				if parseAiStatementConditionVo.CurrentMedication != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentMedication
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_SpecialNotes {
				if parseAiStatementConditionVo.SpecialNotes != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.SpecialNotes
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_IntroductionParagraph {
				if parseAiStatementConditionVo.IntroductionParagraph != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.IntroductionParagraph
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_OnsetAndServiceConnection {
				if parseAiStatementConditionVo.OnsetAndServiceConnection != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.OnsetAndServiceConnection
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_CurrentSymptomsSeverityFrequency {
				if parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_Medication {
				if parseAiStatementConditionVo.Medication != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Medication
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_ImpactOnDailyLife {
				if parseAiStatementConditionVo.ImpactOnDailyLife != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ImpactOnDailyLife
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_ProfessionalImpact {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ProfessionalImpact
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_NexusBetweenSC {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.NexusBetweenSC
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_Request {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Request
					isAiTaskOk = true
				}
			}
		}
		if isAiTaskOk {

			useAiTaskIds = append(useAiTaskIds, aiTaskEntity.ID)
			useAssistantIds = append(useAssistantIds, aJobEntity.ID)
		}
	}
	return statementDetail, useAiTaskIds, useAssistantIds, nil
}

// GetNewStatement ps-gen 使用，可能需要下架
func (c *StatementUsecase) GetNewStatement(tCase TData, tClient TData) (statementDetail StatementDetail, useAiTaskIds []int32, err error) {
	statementDetail, err = c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return statementDetail, useAiTaskIds, err
	}
	isOk, isToPsform, veteranSummary, veteranSummaryAiTaskId := c.GetVeteranSummaryForVeteranSummary(tCase)

	if isOk && !isToPsform {
		useAiTaskIds = append(useAiTaskIds, veteranSummaryAiTaskId)
		if veteranSummary.YearsOfService != "" {
			statementDetail.BaseInfo.YearsOfService = veteranSummary.YearsOfService
		}
		if veteranSummary.RetirementStatus != "" {
			statementDetail.BaseInfo.RetiredFromService = veteranSummary.RetirementStatus
		}
		if veteranSummary.Deployments != "" {
			statementDetail.BaseInfo.Deployments = veteranSummary.Deployments
		}
		if veteranSummary.MaritalStatus != "" {
			statementDetail.BaseInfo.MaritalStatus = veteranSummary.MaritalStatus
		}
		if veteranSummary.Children != "" {
			statementDetail.BaseInfo.Children = veteranSummary.Children
		}
		if veteranSummary.OccupationInService != "" {
			statementDetail.BaseInfo.OccupationInService = veteranSummary.OccupationInService
		}
	}

	for k, v := range statementDetail.Statements {

		aiTaskEntity, err := c.AiTaskUsecase.GetByCond(Eq{"deleted_at": 0,
			"from_type":     AiTaskFromType_statement,
			"case_id":       tCase.Id(),
			"task_uniqcode": v.StatementCondition.StatementConditionId,
			"handle_status": 1,
			"handle_result": 0,
		})
		if err != nil {
			return statementDetail, useAiTaskIds, err
		}
		if aiTaskEntity == nil {
			continue
		}
		if aiTaskEntity.ToPsform == AiTask_ToPsform_Yes {
			continue
		}
		aiResultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": aiTaskEntity.CurrentResultId})
		if err != nil {
			return statementDetail, useAiTaskIds, err
		}
		if aiResultEntity == nil {
			c.log.Error("aiResultEntity is nil")
			continue
		}

		parseAiStatementConditionVo := ParseAiStatementCondition(aiResultEntity.ParseResult)

		isAiTaskOk := false
		for k1, v1 := range statementDetail.Statements[k].Rows {
			if v1.SectionType == Statemt_Section_CurrentTreatmentFacility {
				if parseAiStatementConditionVo.CurrentTreatmentFacility != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentTreatmentFacility
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_CurrentMedication {
				if parseAiStatementConditionVo.CurrentMedication != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentMedication
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_SpecialNotes {
				if parseAiStatementConditionVo.SpecialNotes != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.SpecialNotes
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_IntroductionParagraph {
				if parseAiStatementConditionVo.IntroductionParagraph != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.IntroductionParagraph
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_OnsetAndServiceConnection {
				if parseAiStatementConditionVo.OnsetAndServiceConnection != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.OnsetAndServiceConnection
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_CurrentSymptomsSeverityFrequency {
				if parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.CurrentSymptomsSeverityFrequency
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_Medication {
				if parseAiStatementConditionVo.Medication != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Medication
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_ImpactOnDailyLife {
				if parseAiStatementConditionVo.ImpactOnDailyLife != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ImpactOnDailyLife
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_ProfessionalImpact {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.ProfessionalImpact
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_NexusBetweenSC {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.NexusBetweenSC
					isAiTaskOk = true
				}
			} else if v1.SectionType == Statemt_Section_Request {
				if parseAiStatementConditionVo.ProfessionalImpact != "" {
					statementDetail.Statements[k].Rows[k1].Body = parseAiStatementConditionVo.Request
					isAiTaskOk = true
				}
			}
		}
		if isAiTaskOk {
			useAiTaskIds = append(useAiTaskIds, aiTaskEntity.ID)
		}
	}
	return statementDetail, useAiTaskIds, nil
}

func (c *StatementUsecase) GenerateNewStatementVersionForWebForm(tCase TData, tClient TData) error {

	statementDetail, useAiTaskIds, useAssistantIds, err := c.GetNewStatementForWebForm(tCase, tClient)
	if err != nil {
		return err
	}
	raws, err := json.Marshal(&statementDetail)
	if err != nil {
		return err
	}
	_, err = c.BizStatementSave(false, nil, tCase.Gid(), raws)
	if err == nil { // 成功生成，把useAiTaskIds设置为使用
		if len(useAiTaskIds) > 0 {
			er := c.AiTaskUsecase.UpdatesByCond(map[string]interface{}{"to_psform": AiTask_ToPsform_Yes}, In("id", useAiTaskIds))

			if er != nil {
				c.log.Error(er)
			}
		}
		if len(useAssistantIds) > 0 {
			er := c.AiAssistantJobUsecase.UpdatesByCond(map[string]interface{}{"job_status": AiAssistantJob_JobStatus_Normal}, In("id", useAssistantIds))

			if er != nil {
				c.log.Error(er)
			}
		}
	}
	return err
}

func (c *StatementUsecase) GenerateNewStatementVersion(tCase TData, tClient TData) error {

	statementDetail, useAiTaskIds, err := c.GetNewStatement(tCase, tClient)
	if err != nil {
		return err
	}
	raws, err := json.Marshal(&statementDetail)
	if err != nil {
		return err
	}
	_, err = c.BizStatementSave(false, nil, tCase.Gid(), raws)
	if err == nil { // 成功生成，把useAiTaskIds设置为使用
		if len(useAiTaskIds) > 0 {
			er := c.AiTaskUsecase.UpdatesByCond(map[string]interface{}{"to_psform": AiTask_ToPsform_Yes}, In("id", useAiTaskIds))
			if er != nil {
				c.log.Error(er)
			}
		}
	}
	return err
}

func (c *StatementUsecase) GenerateDocument(tCase TData, tClient TData) error {
	dealName := tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	var wordReader io.Reader
	isNewVersion := true
	if isNewVersion {
		//statementDetail, err := c.GetNewStatement(tCase, tClient, veteranSummary)
		statementDetail, err := c.GetListStatementDetail(false, tClient, tCase, 0)
		if err != nil {
			return err
		}
		wordReader, err = c.WordUsecase.CreatePersonalStatementsWordForAiV1(dealName, statementDetail, 0)
		if err != nil {
			return err
		}
	} else {

		conditions, err := c.StatementConditionUsecase.AllConditions(tCase.Id())

		//conditions, err := SplitCaseStatements(tCase.CustomFields.TextValueByNameBasic(FieldName_statements))
		if err != nil {
			return err
		}
		if len(conditions) == 0 {
			return errors.New("conditions is empty")
		}

		type Record struct {
			Condition    StatementCondition
			AiTaskEntity AiTaskEntity
		}

		var records []Record
		for _, v := range conditions {
			a, err := c.AiTaskUsecase.GetByCond(Eq{"deleted_at": 0,
				"from_type":     AiTaskFromType_statement,
				"case_id":       tCase.Id(),
				"task_uniqcode": v.ID,
				"handle_status": 1,
				"handle_result": 0,
			})
			if err != nil {
				return err
			}
			if a == nil {
				continue
			}
			records = append(records, Record{
				Condition:    v.ToStatementCondition(),
				AiTaskEntity: *a,
			})
		}

		var parseResults []string

		for _, v := range records {
			aiResultEntity, err := c.AiResultUsecase.GetByCond(Eq{"id": v.AiTaskEntity.CurrentResultId})
			if err != nil {
				return err
			}
			if aiResultEntity == nil {
				return errors.New("aiResultEntity is nil")
			}
			parseResults = append(parseResults, aiResultEntity.ParseResult)

		}

		statementBaseInfoList, err := c.GetStatementBaseInfo(&tCase, &tClient)
		if err != nil {
			return err
		}
		wordReader, err = c.WordUsecase.CreatePersonalStatementsWordForAi(dealName, statementBaseInfoList, parseResults)
		if err != nil {
			return err
		}
	}

	dCPersonalStatementsFolderId, aiStatementFileName, boxFileId, err := c.DocStatementBoxFileId(&tClient, &tCase)
	if err != nil {
		return err
	}
	if boxFileId == "" {
		boxFileId, err = c.BoxUsecase.UploadFile(dCPersonalStatementsFolderId, wordReader, aiStatementFileName)
		if err != nil {
			return err
		}
	} else {
		_, err = c.BoxUsecase.UploadFileVersion(boxFileId, wordReader)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *StatementUsecase) DocUpdateStatementBoxFileId(tClient TData, tCase TData) (dCPersonalStatementsFolderId string, updateStatementFileName string, boxFileId string, err error) {

	updateStatementFileName = GenUpdatePersonalStatementsFileName(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err = c.BoxbuzUsecase.DCPersonalStatementsFolderId(&tCase)
	if err != nil {
		return "", "", "", err
	}
	if dCPersonalStatementsFolderId == "" {
		return "", "", "", errors.New("dCPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(dCPersonalStatementsFolderId)
	if err != nil {
		return "", "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == updateStatementFileName {
				boxFileId = resId
				break
			}
		}
	}
	return
}

func (c *StatementUsecase) DocClientPSSourceBoxFileId(tClient TData, tCase TData) (dCPersonalStatementsFolderId string, ClientPSSourceFileName string, boxFileId string, err error) {

	ClientPSSourceFileName = GenPersonalStatementsFileNameForUpdateStatement(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err = c.BoxbuzUsecase.DCPersonalStatementsFolderId(&tCase)
	if err != nil {
		return "", "", "", err
	}
	if dCPersonalStatementsFolderId == "" {
		return "", "", "", errors.New("dCPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(dCPersonalStatementsFolderId)
	if err != nil {
		return "", "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == ClientPSSourceFileName {
				boxFileId = resId
				break
			}
		}
	}
	return
}

func (c *StatementUsecase) DocStatementBoxFileId(tClient *TData, tCase *TData) (dCPersonalStatementsFolderId string, aiStatementFileName string, boxFileId string, err error) {

	if tClient == nil {
		return "", "", "", errors.New("tClient is nil")
	}
	if tCase == nil {
		return "", "", "", errors.New("tCase is nil")
	}

	aiStatementFileName = GenPersonalStatementsFileNameAiAuto(tClient.CustomFields.TextValueByNameBasic(FieldName_first_name), tClient.CustomFields.TextValueByNameBasic(FieldName_last_name), tCase.Id())

	dCPersonalStatementsFolderId, err = c.BoxbuzUsecase.DCPersonalStatementsFolderId(tCase)
	if err != nil {
		return "", "", "", err
	}
	if dCPersonalStatementsFolderId == "" {
		return "", "", "", errors.New("dCPersonalStatementsFolderId is empty")
	}
	resItems, err := c.BoxUsecase.ListItemsInFolderFormat(dCPersonalStatementsFolderId)
	if err != nil {
		return "", "", "", err
	}
	for _, v := range resItems {
		resId := v.GetString("id")
		resType := v.GetString("type")
		resName := v.GetString("name")
		if resType == string(config_box.BoxResType_file) {
			if resName == aiStatementFileName {
				boxFileId = resId
				break
			}
		}
	}
	return
}

func (c *StatementUsecase) GetStatementBaseInfo(tCase *TData, tClient *TData) (result StatementBaseInfoList, err error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}
	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)

	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(tCase)
	if err != nil {
		return nil, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	intakeSubmission, err := c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return nil, err
	}
	if intakeSubmission == nil {
		return nil, errors.New("GetStatementBaseInfo: intakeSubmission is nil")
	}
	result = append(result, StatementBaseInfo{
		Label: "Full Name",
		Value: tClient.CustomFields.TextValueByNameBasic(FieldName_full_name),
	})
	result = append(result, StatementBaseInfo{
		Label: "Unique ID",
		Value: InterfaceToString(tCase.Id()),
	})
	result = append(result, StatementBaseInfo{
		Label: "Branch of Service",
		Value: InterfaceToString(tCase.CustomFields.DisplayValueByName(FieldName_branch)),
	})

	notesInfo := lib.ToTypeMapByString(intakeSubmission.Notes)
	fearsOfService := StatementYearsOfServiceFormat(notesInfo.GetString("content.answers.210.answer"))
	result = append(result, StatementBaseInfo{
		Label: "Years of Service",
		Value: fearsOfService,
	})
	result = append(result, StatementBaseInfo{
		Label: "Retired from service",
		Value: notesInfo.GetString("content.answers.211.answer"),
	})
	maritalStatus := ""
	if notesInfo.GetString("content.answers.213.answer") == "No" {
		maritalStatus = "Single"
	} else {
		maritalStatus = "Married"
	}
	result = append(result, StatementBaseInfo{
		Label: "Marital Status",
		Value: maritalStatus,
	})
	result = append(result, StatementBaseInfo{
		Label: "Children",
		Value: notesInfo.GetString("content.answers.218.answer"),
	})
	result = append(result, StatementBaseInfo{
		Label: "Occupation in service",
		Value: notesInfo.GetString("content.answers.207.answer"),
	})

	// todo:lgl Deployments

	return
}

func StatementYearsOfServiceFormat(str string) (r string) {
	aa := strings.Split(str, "-")
	for _, v := range aa {
		cc := strings.Split(v, "/")
		if len(cc) == 2 {
			if r == "" {
				r = cc[1]
			} else {
				r += "-" + cc[1]
			}
		}
	}
	return
}

func IsNewStatementVersion(tCase TData) bool {
	psmUrl := tCase.CustomFields.TextValueByNameBasic(FieldName_personal_statement_manager)
	if psmUrl != "" {
		return true
	}
	return false
}
