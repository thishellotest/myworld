package biz

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
	"vbc/internal/config_vbc"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

const Statemt_Section_CurrentTreatmentFacility = "CurrentTreatmentFacility"
const Statemt_Section_CurrentMedication = "CurrentMedication"
const Statemt_Section_SpecialNotes = "SpecialNotes"
const Statemt_Section_IntroductionParagraph = "IntroductionParagraph"
const Statemt_Section_OnsetAndServiceConnection = "OnsetAndServiceConnection"
const Statemt_Section_CurrentSymptomsSeverityFrequency = "CurrentSymptomsSeverityFrequency"
const Statemt_Section_Medication = "Medication"
const Statemt_Section_ImpactOnDailyLife = "ImpactOnDailyLife"
const Statemt_Section_ProfessionalImpact = "ProfessionalImpact"
const Statemt_Section_NexusBetweenSC = "NexusBetweenSC"
const Statemt_Section_Request = "Request"

func CompareListStatementDetail(a1, b1 StatementDetail) (isEqual bool) {

	if a1.BaseInfo.YearsOfService != b1.BaseInfo.YearsOfService {
		return false
	}
	if a1.BaseInfo.RetiredFromService != b1.BaseInfo.RetiredFromService {
		return false
	}
	if a1.BaseInfo.Deployments != b1.BaseInfo.Deployments {
		return false
	}
	if a1.BaseInfo.MaritalStatus != b1.BaseInfo.MaritalStatus {
		return false
	}
	if a1.BaseInfo.Children != b1.BaseInfo.Children {
		return false
	}
	if a1.BaseInfo.OccupationInService != b1.BaseInfo.OccupationInService {
		return false
	}
	if a1.BaseInfo.BranchOfService != b1.BaseInfo.BranchOfService {
		return false
	}

	a := a1.Statements
	b := b1.Statements
	if len(a) != len(b) {
		return false
	}

	for k, _ := range a {
		if a[k].StatementCondition.ConditionValue != b[k].StatementCondition.ConditionValue {
			return false
		}
		if a[k].StatementCondition.FrontValue != b[k].StatementCondition.FrontValue {
			return false
		}
		if a[k].StatementCondition.BehindValue != b[k].StatementCondition.BehindValue {
			return false
		}
		if a[k].StatementCondition.Category != b[k].StatementCondition.Category {
			return false
		}
		//if a[k].CurrentMedication != b[k].CurrentMedication {
		//	//lib.DPrintln(1)
		//	return false
		//}
		//if a[k].CurrentTreatmentFacility != b[k].CurrentTreatmentFacility {
		//	//lib.DPrintln(2)
		//	return false
		//}
		//if a[k].StatementCondition.ConditionValue != b[k].StatementCondition.ConditionValue {
		//	//lib.DPrintln(3)
		//	return false
		//}
		if len(a[k].Rows) != len(b[k].Rows) {
			//lib.DPrintln(4)
			return false
		}
		for k1, _ := range a[k].Rows {
			if a[k].Rows[k1].Body != b[k].Rows[k1].Body {
				return false
			}
		}
	}
	return true
}

type StatementDetail struct {
	CaseId     int32                   `json:"case_id"`
	Gid        string                  `json:"gid"`
	Version    int32                   `json:"version"`
	BaseInfo   StatementDetailBaseInfo `json:"base_info"`
	Statements ListStatementDetail     `json:"statements"`
	Versions   []StatemtVersionVo      `json:"versions"`
}

func FormatStatementDetail(statementDetail StatementDetail) StatementDetail {

	for k, v := range statementDetail.Statements {
		for k1, v1 := range v.Rows {
			newBody := ""
			lines := strings.Split(v1.Body, "\n")
			for _, v2 := range lines {
				v2 := strings.TrimSpace(v2)
				if v2 == "" {
					continue
				}
				if newBody == "" {
					newBody += v2
				} else {
					newBody += "\n\n" + v2
				}

			}
			statementDetail.Statements[k].Rows[k1].Body = newBody
		}
	}

	return statementDetail
}

type ListStatementDetail []StatementDetailVo

type StatementDetailVo struct {
	StatementCondition StatementCondition `json:"statement_condition"`
	//CurrentTreatmentFacility string                   `json:"current_treatment_facility"`
	//CurrentMedication        string                   `json:"current_medication"`
	Rows ListStatementDetailVoRow `json:"rows"`
}

func (c *StatementDetailVo) IsEmptyResult() bool {
	isEmpty := true
	for _, v := range c.Rows {
		if v.Body != "" {
			isEmpty = false
		}
	}
	return isEmpty
}

type StatementDetailBaseInfo struct {
	FullName            string `json:"full_name"`
	CaseId              int32  `json:"case_id"`
	BranchOfService     string `json:"branch_of_service"`
	YearsOfService      string `json:"years_of_service"`
	RetiredFromService  string `json:"retired_from_service"`
	Deployments         string `json:"deployments"`
	MaritalStatus       string `json:"marital_status"`
	Children            string `json:"children"`
	OccupationInService string `json:"occupation_in_service"`
}

func (c *StatementDetailVo) ToStatemtEntity() StatemtEntity {
	var statemtEntity StatemtEntity

	//statemtEntity.CurrentTreatmentFacility = c.CurrentTreatmentFacility
	//statemtEntity.CurrentMedication = c.CurrentMedication
	statemtEntity.ConditionUniqid = InterfaceToString(c.StatementCondition.StatementConditionId)
	statemtEntity.FrontValue = c.StatementCondition.FrontValue
	statemtEntity.ConditionValue = c.StatementCondition.ConditionValue
	statemtEntity.BehindValue = c.StatementCondition.BehindValue
	statemtEntity.Category = c.StatementCondition.Category

	for _, v := range c.Rows {
		if v.SectionType == Statemt_Section_CurrentTreatmentFacility {
			statemtEntity.CurrentTreatmentFacility = v.Body
		} else if v.SectionType == Statemt_Section_CurrentMedication {
			statemtEntity.CurrentMedication = v.Body
		} else if v.SectionType == Statemt_Section_SpecialNotes {
			statemtEntity.SpecialNotes = v.Body
		} else if v.SectionType == Statemt_Section_IntroductionParagraph {
			statemtEntity.IntroductionParagraph = v.Body
		} else if v.SectionType == Statemt_Section_OnsetAndServiceConnection {
			statemtEntity.OnsetAndServiceConnection = v.Body
		} else if v.SectionType == Statemt_Section_CurrentSymptomsSeverityFrequency {
			statemtEntity.CurrentSymptomsSeverityFrequency = v.Body
		} else if v.SectionType == Statemt_Section_Medication {
			statemtEntity.Medication = v.Body
		} else if v.SectionType == Statemt_Section_ImpactOnDailyLife {
			statemtEntity.ImpactOnDailyLife = v.Body
		} else if v.SectionType == Statemt_Section_ProfessionalImpact {
			statemtEntity.ProfessionalImpact = v.Body
		} else if v.SectionType == Statemt_Section_NexusBetweenSC {
			statemtEntity.NexusBetweenSC = v.Body
		} else if v.SectionType == Statemt_Section_Request {
			statemtEntity.Request = v.Body
		}
	}

	return statemtEntity
}

type ListStatementDetailVoRow []StatementDetailVoRow

func (c ListStatementDetailVoRow) GetCurrentTreatmentFacilityAndCurrentMedication() (CurrentTreatmentFacility string, CurrentMedication string) {
	for _, v := range c {
		if v.SectionType == Statemt_Section_CurrentTreatmentFacility {
			CurrentTreatmentFacility = v.Body
		} else if v.SectionType == Statemt_Section_CurrentMedication {
			CurrentMedication = v.Body
		}
	}
	return
}

func (c ListStatementDetailVoRow) ToStringForStandardHeaderRevision() string {
	text := ""
	for _, v := range c {
		if v.IsStatementRow() {
			if text == "" {
				text += v.Title + ":\n"
				text += v.Body
			} else {
				text += "\n\n" + v.Title + ":\n"
				text += v.Body
			}
		}
	}
	return text
}

type StatementDetailVoRow struct {
	Title       string `json:"title"`
	SectionType string `json:"section_type"`
	Body        string `json:"body"`
}

func IsStatementRow(SectionType string) bool {

	if SectionType == Statemt_Section_SpecialNotes ||
		SectionType == Statemt_Section_IntroductionParagraph ||
		SectionType == Statemt_Section_OnsetAndServiceConnection ||
		SectionType == Statemt_Section_CurrentSymptomsSeverityFrequency ||
		SectionType == Statemt_Section_Medication ||
		SectionType == Statemt_Section_ImpactOnDailyLife ||
		SectionType == Statemt_Section_ProfessionalImpact ||
		SectionType == Statemt_Section_NexusBetweenSC ||
		SectionType == Statemt_Section_Request {
		return true
	}
	return false
}
func (c *StatementDetailVoRow) IsStatementRow() bool {
	return IsStatementRow(c.SectionType)
}

func GetSectionTitleFromSectionType(SectionType string) string {
	for _, v := range ListStatementDetailVoRowConfig {
		if v.SectionType == SectionType {
			return v.Title
		}
	}
	return ""
}

var ListStatementDetailVoRowConfig = ListStatementDetailVoRow{
	{
		Title:       "Current Treatment Facility",
		SectionType: Statemt_Section_CurrentTreatmentFacility,
	},
	{
		Title:       "Current Medication",
		SectionType: Statemt_Section_CurrentMedication,
	},
	{
		Title:       "Special Notes",
		SectionType: Statemt_Section_SpecialNotes,
	},
	{
		Title:       "Introduction",
		SectionType: Statemt_Section_IntroductionParagraph,
	},
	{
		Title:       "Onset and Service Connection",
		SectionType: Statemt_Section_OnsetAndServiceConnection,
	},
	{
		Title:       "Current Symptoms, Severity and Frequency",
		SectionType: Statemt_Section_CurrentSymptomsSeverityFrequency,
	},
	{
		Title:       "Medication",
		SectionType: Statemt_Section_Medication,
	},
	{
		Title:       "Impact on Daily Life",
		SectionType: Statemt_Section_ImpactOnDailyLife,
	},
	{
		Title:       "Professional Impact",
		SectionType: Statemt_Section_ProfessionalImpact,
	},
	{
		Title:       "Nexus Between Service and Current Condition",
		SectionType: Statemt_Section_NexusBetweenSC,
	},
	{
		Title:       "Request",
		SectionType: Statemt_Section_Request,
	},
}

func StatementDetailVoRowBody(sectionType string, stateEntity StatemtEntity) string {
	if sectionType == Statemt_Section_CurrentTreatmentFacility {
		return stateEntity.CurrentTreatmentFacility
	}
	if sectionType == Statemt_Section_CurrentMedication {
		return stateEntity.CurrentMedication
	}
	if sectionType == Statemt_Section_SpecialNotes {
		return stateEntity.SpecialNotes
	}
	if sectionType == Statemt_Section_IntroductionParagraph {
		return stateEntity.IntroductionParagraph
	}
	if sectionType == Statemt_Section_OnsetAndServiceConnection {
		return stateEntity.OnsetAndServiceConnection
	}
	if sectionType == Statemt_Section_CurrentSymptomsSeverityFrequency {
		return stateEntity.CurrentSymptomsSeverityFrequency
	}
	if sectionType == Statemt_Section_Medication {
		return stateEntity.Medication
	}
	if sectionType == Statemt_Section_ImpactOnDailyLife {
		return stateEntity.ImpactOnDailyLife
	}
	if sectionType == Statemt_Section_ProfessionalImpact {
		return stateEntity.ProfessionalImpact
	}
	if sectionType == Statemt_Section_NexusBetweenSC {
		return stateEntity.NexusBetweenSC
	}
	if sectionType == Statemt_Section_Request {
		return stateEntity.Request
	}
	return ""
}

func (c *StatementUsecase) BizStatementRevertVersion(userFacade UserFacade, caseGid string, versionId int32) (lib.TypeMap, error) {

	tCase, err := c.RecordbuzUsecase.VerifyDataPermission(caseGid, userFacade.TData)
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("The data does not exist or there is no permission")
	}
	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return nil, errors.New("tClient is nil")
	}

	statementDetail, err := c.GetListStatementDetail(false, *tClient, *tCase, versionId)
	if err != nil {
		return nil, err
	}
	//newVersion, err := c.StatemtUsecase.ObtainAvailableVersionID(tCase.Id())
	//if err != nil {
	//	return nil, err
	//}
	bytes, err := json.Marshal(statementDetail)
	if err != nil {
		return nil, err
	}
	return c.BizStatementSave(false, &userFacade.TData, caseGid, bytes)
}

func (c *StatementUsecase) BizStatementVerifyPassword(caseId int32, password string) (isOk bool, err error) {

	c.log.Info("BizStatementVerifyPassword:", caseId, "-", password)
	realPassword, err := c.PersonalStatementPassword(caseId)
	if err != nil {
		return false, err
	}
	c.log.Info("BizStatementVerifyPassword:", caseId, "-", password, "-", realPassword, "-")
	if password == realPassword {
		return true, nil
	}
	return false, nil
}

type BizStatementDetailVersionsRequest struct {
	NewVersion int32 `json:"new_version"`
	OldVersion int32 `json:"old_version"`
}

func (c *StatementUsecase) BizStatementDetailVersions(usePasswordAccess bool, tUser *TData, caseGid string, raws []byte) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementDetail: tClient is nil")
	}
	var bizStatementDetailVersionsRequest BizStatementDetailVersionsRequest
	bizStatementDetailVersionsRequest = lib.BytesToTDef(raws, bizStatementDetailVersionsRequest)
	if bizStatementDetailVersionsRequest.NewVersion <= 0 || bizStatementDetailVersionsRequest.OldVersion <= 0 {
		return nil, errors.New("The parameters are incorrect")
	}

	newStatementDetail, err := c.GetListStatementDetail(usePasswordAccess, *tClient, *tCase, bizStatementDetailVersionsRequest.NewVersion)
	if err != nil {
		return nil, err
	}

	oldStatementDetail, err := c.GetListStatementDetail(usePasswordAccess, *tClient, *tCase, bizStatementDetailVersionsRequest.OldVersion)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set("new_data", newStatementDetail)
	data.Set("old_data", oldStatementDetail)
	return data, nil
}

func HasCocoAssistant(usePasswordAccess bool, tUser *TData) bool {
	if usePasswordAccess {
		return false
	}
	return true
	if tUser != nil {
		gid := tUser.Gid()
		if gid == config_vbc.User_Dev_gid ||
			gid == config_vbc.User_Lili_gid ||
			gid == config_vbc.User_Edward_gid ||
			gid == config_vbc.User_Yannan_gid {
			return true
		}
	}
	return false
}

func (c *StatementUsecase) BizStatementDetail(usePasswordAccess bool, tUser *TData, caseGid string) (lib.TypeMap, error) {

	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementDetail: tClient is nil")
	}
	statementDetail, err := c.GetListStatementDetail(usePasswordAccess, *tClient, *tCase, 0)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set("data", statementDetail)
	data.Set("use_password_access", usePasswordAccess)
	data.Set("has_coco_assistant", HasCocoAssistant(usePasswordAccess, tUser))

	if usePasswordAccess {
		data.Set("doc_email_body", "")
	} else {
		docEmailBody, _ := c.DocEmailUsecase.DocEmailResultTextByCase(*tCase, *tClient)
		data.Set("doc_email_body", docEmailBody)
	}

	a, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(tCase.Id())
	if err != nil {
		c.log.Error(err)
	}
	data.Set("has_client_PW_active", a)
	return data, nil
}

func (c *StatementUsecase) GetStatementDetailBaseInfo(statemtEntity *StatemtEntity, tClient TData, tCase TData) (statementDetailBaseInfo StatementDetailBaseInfo, err error) {

	uniqcode := tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode)
	var uniqcodes []string
	uniqcodes = append(uniqcodes, uniqcode)
	isPrimaryCase, primaryCase, err := c.FeeUsecase.UsePrimaryCaseCalc(&tCase)
	if err != nil {
		return statementDetailBaseInfo, err
	}
	if !isPrimaryCase {
		uniqcodes = append(uniqcodes, primaryCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode))
	}

	statementDetailBaseInfo.FullName = tClient.CustomFields.TextValueByNameBasic(FieldName_full_name)
	statementDetailBaseInfo.CaseId = tCase.Id()
	statementDetailBaseInfo.BranchOfService = InterfaceToString(tCase.CustomFields.DisplayValueByName(FieldName_branch))

	intakeSubmission, err := c.JotformSubmissionUsecase.GetLatestIntakeFormInfoByFormId(uniqcodes)
	if err != nil {
		return statementDetailBaseInfo, err
	}
	if intakeSubmission != nil {
		notesInfo := lib.ToTypeMapByString(intakeSubmission.Notes)
		fearsOfService := StatementYearsOfServiceFormat(notesInfo.GetString("content.answers.210.answer"))
		statementDetailBaseInfo.YearsOfService = fearsOfService
		statementDetailBaseInfo.RetiredFromService = notesInfo.GetString("content.answers.211.answer")
		maritalStatus := ""
		if notesInfo.GetString("content.answers.213.answer") == "No" {
			maritalStatus = "Single"
		} else {
			maritalStatus = "Married"
		}
		statementDetailBaseInfo.MaritalStatus = maritalStatus
		statementDetailBaseInfo.Children = notesInfo.GetString("content.answers.218.answer")
		statementDetailBaseInfo.OccupationInService = notesInfo.GetString("content.answers.207.answer")
	}
	if statemtEntity != nil {
		//Deployments 可以通过AI任务初始化获取
		statementDetailBaseInfo.YearsOfService = statemtEntity.YearsOfService
		statementDetailBaseInfo.RetiredFromService = statemtEntity.RetiredFromService
		statementDetailBaseInfo.MaritalStatus = statemtEntity.MaritalStatus
		statementDetailBaseInfo.Children = statemtEntity.Children
		statementDetailBaseInfo.OccupationInService = statemtEntity.OccupationInService
		statementDetailBaseInfo.Deployments = statemtEntity.Deployments
	}

	return statementDetailBaseInfo, nil
}

func (c *StatementUsecase) GetListStatementDetail(usePasswordAccess bool, tClient TData, tCase TData, version int32) (statementDetail StatementDetail, err error) {

	//statements := tCase.CustomFields.TextValueByNameBasic(FieldName_statements)
	statementConditions, err := c.StatementConditionUsecase.AllConditions(tCase.Id())
	if err != nil {
		return statementDetail, err
	}
	statementDetail.CaseId = tCase.Id()
	statementDetail.Gid = tCase.Gid()

	//list, err := SplitCaseStatements(statements)
	//if err != nil {
	//	return statementDetail, err
	//}

	var listStatemtEntity ListStatemtEntity
	if version == 0 {
		listStatemtEntity, err = c.StatemtUsecase.AllLatestStatements(tCase.Id())
	} else {
		listStatemtEntity, err = c.StatemtUsecase.AllStatementsByVersion(tCase.Id(), version)
	}
	if err != nil {
		return statementDetail, err
	}

	var statemtEntity *StatemtEntity
	if len(listStatemtEntity) > 0 {
		statemtEntity = listStatemtEntity[0]
		statementDetail.Version = statemtEntity.Versions
	}
	statementDetailBaseInfo, err := c.GetStatementDetailBaseInfo(statemtEntity, tClient, tCase)
	if err != nil {
		return statementDetail, err
	}
	statementDetail.BaseInfo = statementDetailBaseInfo

	if !usePasswordAccess {
		versions, err := c.StatemtUsecase.Versions(tCase.Id())
		if err != nil {
			return statementDetail, err
		}
		statementDetail.Versions = versions
	}
	for _, v := range statementConditions {

		copyConfig := make(ListStatementDetailVoRow, len(ListStatementDetailVoRowConfig))
		copy(copyConfig, ListStatementDetailVoRowConfig)

		vo := StatementDetailVo{
			Rows: copyConfig,
		}
		vo.StatementCondition = v.ToStatementCondition()
		stateEntity := listStatemtEntity.GetByConditionUniqid(InterfaceToString(vo.StatementCondition.StatementConditionId))
		//lib.DPrintln(stateEntity.ID, stateEntity.ConditionUniqid)
		if stateEntity != nil {
			vo.StatementCondition.FrontValue = stateEntity.FrontValue
			//vo.StatementCondition.ConditionValue = stateEntity.ConditionValue
			vo.StatementCondition.BehindValue = stateEntity.BehindValue
			vo.StatementCondition.Category = stateEntity.Category
			//lib.DPrintln("*stateEntity.Request:", stateEntity.Request)
			for k1, v1 := range vo.Rows {
				vo.Rows[k1].Body = StatementDetailVoRowBody(v1.SectionType, *stateEntity)
			}
			//vo.CurrentTreatmentFacility = stateEntity.CurrentTreatmentFacility
			//vo.CurrentMedication = stateEntity.CurrentMedication
		}

		statementDetail.Statements = append(statementDetail.Statements, vo)
	}
	return
}

func SortListStatementDetail(listStatementDetail ListStatementDetail) (r ListStatementDetail) {

	var dest [][]StatementDetailVo
	for _, v := range StatementConditionCategoryOrder {
		var temp []StatementDetailVo
		for k1, v1 := range listStatementDetail {
			if v1.StatementCondition.Category == v {
				temp = append(temp, listStatementDetail[k1])
			}
		}
		if len(temp) > 0 {
			dest = append(dest, temp)
		}
	}

	for k, _ := range dest {
		r = append(r, dest[k]...)
	}
	lib.DPrintln("SortListStatementDetail: ", r)
	return r
}

func (c *StatementUsecase) BizStatementSave(usePasswordAccess bool, tUser *TData, caseGid string, rawData []byte) (lib.TypeMap, error) {
	if usePasswordAccess {
		return nil, errors.New("Not allowed to operate")
	}
	if caseGid == "" {
		return nil, errors.New("Parameters Incorrect")
	}
	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
	if err != nil {
		return nil, err
	}
	if tClient == nil {
		return nil, errors.New("BizStatementSave: tClient is nil")
	}

	var statementDetail StatementDetail
	err = json.Unmarshal(rawData, &statementDetail)
	if err != nil {
		return nil, err
	}
	statementDetail = FormatStatementDetail(statementDetail)

	oldStatementDetail, err := c.GetListStatementDetail(usePasswordAccess, *tClient, *tCase, 0)
	if err != nil {
		return nil, err
	}
	var newStatementDetail StatementDetail
	isModified := false
	if !CompareListStatementDetail(statementDetail, oldStatementDetail) { // 按key值比较完毕
		isModified = true
		newVerions, err := c.StatemtUsecase.ObtainAvailableVersionID(tCase.Id())
		if err != nil {
			return nil, err
		}
		// 需要重新排序更新了
		sortListStatementDetail := SortListStatementDetail(statementDetail.Statements)

		sort := 1000
		for k, v := range sortListStatementDetail {

			statementConditionEntity, err := c.StatementConditionUsecase.GetByCond(Eq{"id": v.StatementCondition.StatementConditionId})
			if err != nil {
				return nil, err
			}
			if statementConditionEntity == nil {
				return nil, errors.New(InterfaceToString(v.StatementCondition.StatementConditionId) + ":statementConditionEntity is nil")
			}
			statementConditionEntity.Sort = sort + k
			statementConditionEntity.ConditionValue = v.StatementCondition.ConditionValue
			statementConditionEntity.FrontValue = v.StatementCondition.FrontValue
			statementConditionEntity.BehindValue = v.StatementCondition.BehindValue
			statementConditionEntity.Category = v.StatementCondition.Category
			err = c.CommonUsecase.DB().Save(&statementConditionEntity).Error
			if err != nil {
				return nil, err
			}

			statemtEntity := v.ToStatemtEntity()
			statemtEntity.CaseId = tCase.Id()
			statemtEntity.Versions = newVerions
			statemtEntity.CreatedAt = time.Now().Unix()
			statemtEntity.UpdatedAt = time.Now().Unix()
			statemtEntity.YearsOfService = statementDetail.BaseInfo.YearsOfService
			statemtEntity.Children = statementDetail.BaseInfo.Children
			statemtEntity.Deployments = statementDetail.BaseInfo.Deployments
			statemtEntity.MaritalStatus = statementDetail.BaseInfo.MaritalStatus
			statemtEntity.RetiredFromService = statementDetail.BaseInfo.RetiredFromService
			statemtEntity.OccupationInService = statementDetail.BaseInfo.OccupationInService
			if !usePasswordAccess {
				if tUser != nil {
					statemtEntity.ModifiedBy = tUser.Gid()
				}
			}
			err = c.CommonUsecase.DB().Save(&statemtEntity).Error
			if err != nil {
				return nil, err
			}
		}
		newStatementDetail, err = c.GetListStatementDetail(usePasswordAccess, *tClient, *tCase, 0)
		if err != nil {
			return nil, err
		}

		go func() {
			usePw, err := c.PersonalWebformUsecase.IsUseNewPersonalWebForm(tCase.Id())
			if err != nil {
				c.log.Error(err)
			} else {
				if !usePw {
					er := c.GenerateDocument(*tCase, *tClient)
					if er != nil {
						c.log.Error(er)
					}
				}
			}
		}()

		err = c.StatementConditionBuzUsecase.UpdateCaseStatement(tCase.Id())
		if err != nil {
			c.log.Error(err, " caseId: ", tCase.Id())
		}
	} else {
		newStatementDetail = oldStatementDetail
	}

	data := make(lib.TypeMap)
	data.Set("data", newStatementDetail)
	data.Set("is_modified", isModified)
	return data, nil
}

func (c *StatementUsecase) PersonalStatementPassword(caseId int32) (password string, err error) {
	if caseId <= 0 {
		return "", errors.New("The case does not exist")
	}
	key := MapKeyPersonalStatementPassword(caseId)
	password, err = c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if password == "" {
		password = lib.GenerateSafePassword(8)
		err = c.MapUsecase.Set(key, password)
		if err != nil {
			return "", err
		}
	}
	return password, nil
}

func (c *StatementUsecase) GetUpdatePSTextForAiParamWithAssistant(tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, sectionType string) (referenceContent string, statement string, err error) {

	statementDetail, err := c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return "", "", err
	}
	var statementDetailVo *StatementDetailVo
	for k, v := range statementDetail.Statements {
		if v.StatementCondition.StatementConditionId == statementConditionEntity.ID {
			statementDetailVo = to.Ptr(statementDetail.Statements[k])
			break
		}
	}
	if statementDetailVo == nil {
		return "", "", errors.New("The corresponding statement was not found")
	}

	isOk := false
	body := ""

	for _, v := range statementDetailVo.Rows {
		if v.SectionType == sectionType {
			isOk = true
			body = v.Body
		}
	}
	if isOk {
		// start test code
		//CurrentTreatmentFacility = "Doctors Community Hospital"
		//CurrentMedication = "Flexeril, Naproxen, Ibuprofen"
		// end test code

		return referenceContent, body, nil
	}

	return "", "", nil

}

func (c *StatementUsecase) GetUpdatePSTextForAiParamWithMedication(tClient TData, tCase TData, statementConditionEntity StatementConditionEntity, sectionType string) (referenceContent string, statement string, err error) {

	statementDetail, err := c.GetListStatementDetail(false, tClient, tCase, 0)
	if err != nil {
		return "", "", err
	}
	var statementDetailVo *StatementDetailVo
	for k, v := range statementDetail.Statements {
		if v.StatementCondition.StatementConditionId == statementConditionEntity.ID {
			statementDetailVo = to.Ptr(statementDetail.Statements[k])
			break
		}
	}
	if statementDetailVo == nil {
		return "", "", errors.New("The corresponding statement was not found")
	}

	isOk := false
	body := ""
	CurrentTreatmentFacility := ""
	CurrentMedication := ""
	for _, v := range statementDetailVo.Rows {
		if v.SectionType == sectionType {
			isOk = true
			body = v.Body
			//body += "Current Treatment Facility: "+statementDetailVo.
			//return v.Body, nil
		} else if v.SectionType == Statemt_Section_CurrentTreatmentFacility {
			CurrentTreatmentFacility = v.Body
		} else if v.SectionType == Statemt_Section_CurrentMedication {
			CurrentMedication = v.Body
		}
	}
	if isOk {
		// start test code
		//CurrentTreatmentFacility = "Doctors Community Hospital"
		//CurrentMedication = "Flexeril, Naproxen, Ibuprofen"
		// end test code

		referenceContent = fmt.Sprintf("Current Treatment Facility: %s\nCurrent Medication: %s", CurrentTreatmentFacility, CurrentMedication)
		return referenceContent, body, nil
	}

	return "", "", nil

}

type ParseAiStatementConditionVo struct {
	NameOfDisabilityCondition        string `json:"name_of_disability_condition"`
	CurrentTreatmentFacility         string `json:"current_treatment_facility"`
	CurrentMedication                string `json:"current_medication"`
	SpecialNotes                     string `json:"special_notes"`
	IntroductionParagraph            string `json:"introduction_paragraph"`
	OnsetAndServiceConnection        string `json:"onset_and_service_connection"`
	CurrentSymptomsSeverityFrequency string `json:"current_symptoms_severity_frequency"`
	Medication                       string `json:"medication"`
	ImpactOnDailyLife                string `json:"impact_on_daily_life"`
	ProfessionalImpact               string `json:"professional_impact"`
	NexusBetweenSC                   string `json:"nexus_between_sc"`
	Request                          string `json:"request"`
}

type SourceParseAiStatementConditionVo struct {
	NameOfDisabilityCondition        string   `json:"name_of_disability_condition"`
	CurrentTreatmentFacility         string   `json:"current_treatment_facility"`
	CurrentMedication                string   `json:"current_medication"`
	SpecialNotes                     []string `json:"special_notes"`
	IntroductionParagraph            []string `json:"introduction_paragraph"`
	OnsetAndServiceConnection        []string `json:"onset_and_service_connection"`
	CurrentSymptomsSeverityFrequency []string `json:"current_symptoms_severity_frequency"`
	Medication                       []string `json:"medication"`
	ImpactOnDailyLife                []string `json:"impact_on_daily_life"`
	ProfessionalImpact               []string `json:"professional_impact"`
	NexusBetweenSC                   []string `json:"nexus_between_sc"`
	Request                          []string `json:"request"`
}

func IsEmptyResultForStatement(str string) bool {
	str = strings.ToLower(str)
	if str == "" {
		return true
	}
	if strings.Index(str, "not specified") >= 0 || strings.Index(str, "not provided") >= 0 {
		return true
	}
	return false
}

func ParseAiStatementCondition(parseResult string) (result ParseAiStatementConditionVo) {

	//re := regexp.MustCompile(`(?m)^#{1,2}\s*`)
	//// (?m) 开启多行模式，这样 ^ 代表每一行的开头。
	//// 接着是 1 或 2 个 #（#{1,2}）
	//// 后面也允许有空格（\s*）
	//parseResult = re.ReplaceAllString(parseResult, "")
	var vo SourceParseAiStatementConditionVo

	isSpecialNotes := false
	isIntroductionParagraph := false
	isOnsetAndServiceConnection := false
	isCurrentSymptomsSeverityFrequency := false
	isMedication := false
	isImpactOnDailyLife := false
	isProfessionalImpact := false
	isNexusBetweenSC := false
	isRequest := false

	stateToDefault := func() {
		isSpecialNotes = false
		isIntroductionParagraph = false
		isOnsetAndServiceConnection = false
		isCurrentSymptomsSeverityFrequency = false
		isMedication = false
		isImpactOnDailyLife = false
		isProfessionalImpact = false
		isNexusBetweenSC = false
		isRequest = false
	}

	arr := strings.Split(parseResult, "\n")
	for _, v := range arr {
		v := strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if vo.NameOfDisabilityCondition == "" && strings.Index(v, "Name of Disability") >= 0 {
			vo.NameOfDisabilityCondition = v

		} else if vo.CurrentTreatmentFacility == "" && strings.Index(v, "Current Treatment Facility") >= 0 {

			if !IsEmptyResultForStatement(v) {
				vo.CurrentTreatmentFacility = v
			}

		} else if vo.CurrentMedication == "" && strings.Index(v, "Current Medication") >= 0 {
			if !IsEmptyResultForStatement(v) {
				vo.CurrentMedication = v
			}
		} else {
			if strings.Index(v, "SERVICE CONNECTION") >= 0 {
				stateToDefault()
				isSpecialNotes = true
			} else if strings.Index(v, "I am respectfully requesting") >= 0 {
				stateToDefault()
				isIntroductionParagraph = true
			} else if strings.Index(v, "Onset and Service Connection") >= 0 {
				stateToDefault()
				isOnsetAndServiceConnection = true
			} else if strings.Index(v, "Current Symptoms") >= 0 {
				stateToDefault()
				isCurrentSymptomsSeverityFrequency = true
			} else if strings.Index(v, "Medication:") >= 0 {
				stateToDefault()
				isMedication = true
			} else if strings.Index(v, "Impact on Daily Life") >= 0 {
				stateToDefault()
				isImpactOnDailyLife = true
			} else if strings.Index(v, "Professional Impact") >= 0 {
				stateToDefault()
				isProfessionalImpact = true
			} else if strings.Index(v, "Nexus Between Service") >= 0 {
				stateToDefault()
				isNexusBetweenSC = true
			} else if strings.Index(v, "Request:") >= 0 {
				stateToDefault()
				isRequest = true
			}

			if isSpecialNotes {
				vo.SpecialNotes = append(vo.SpecialNotes, v)
			}
			if isIntroductionParagraph {
				vo.IntroductionParagraph = append(vo.IntroductionParagraph, v)
			}
			if isOnsetAndServiceConnection {
				vo.OnsetAndServiceConnection = append(vo.OnsetAndServiceConnection, v)
			}
			if isCurrentSymptomsSeverityFrequency {
				vo.CurrentSymptomsSeverityFrequency = append(vo.CurrentSymptomsSeverityFrequency, v)
			}
			if isMedication {
				vo.Medication = append(vo.Medication, v)
			}
			if isImpactOnDailyLife {
				vo.ImpactOnDailyLife = append(vo.ImpactOnDailyLife, v)
			}
			if isProfessionalImpact {
				vo.ProfessionalImpact = append(vo.ProfessionalImpact, v)
			}
			if isNexusBetweenSC {
				vo.NexusBetweenSC = append(vo.NexusBetweenSC, v)
			}
			if isRequest {
				vo.Request = append(vo.Request, v)
			}
		}
	}

	return SourceParseAiStatementConditionVoToParseAiStatementConditionVo(vo)
}

func SourceParseAiStatementConditionVoToParseAiStatementConditionVo(vo SourceParseAiStatementConditionVo) (result ParseAiStatementConditionVo) {
	cc := strings.Split(vo.NameOfDisabilityCondition, ":")
	if len(cc) > 1 {
		result.NameOfDisabilityCondition = strings.TrimSpace(cc[1])
	}
	cc = strings.Split(vo.CurrentTreatmentFacility, ":")
	if len(cc) > 1 {
		result.CurrentTreatmentFacility = strings.TrimSpace(cc[1])
	}
	cc = strings.Split(vo.CurrentMedication, ":")
	if len(cc) > 1 {
		result.CurrentMedication = strings.TrimSpace(cc[1])
	}
	if len(vo.SpecialNotes) > 0 {
		result.SpecialNotes = strings.Join(vo.SpecialNotes, "\n")
	}
	if len(vo.IntroductionParagraph) > 0 {
		result.IntroductionParagraph = strings.Join(vo.IntroductionParagraph, "\n")
	}
	if len(vo.OnsetAndServiceConnection) > 1 {
		result.OnsetAndServiceConnection = strings.Join(vo.OnsetAndServiceConnection[1:], "\n")
	}
	if len(vo.CurrentSymptomsSeverityFrequency) > 1 {
		result.CurrentSymptomsSeverityFrequency = strings.Join(vo.CurrentSymptomsSeverityFrequency[1:], "\n")
	}
	if len(vo.Medication) > 1 {
		result.Medication = strings.Join(vo.Medication[1:], "\n")
	}
	if len(vo.ImpactOnDailyLife) > 1 {
		result.ImpactOnDailyLife = strings.Join(vo.ImpactOnDailyLife[1:], "\n")
	}
	if len(vo.ProfessionalImpact) > 1 {
		result.ProfessionalImpact = strings.Join(vo.ProfessionalImpact[1:], "\n")
	}
	if len(vo.NexusBetweenSC) > 1 {
		result.NexusBetweenSC = strings.Join(vo.NexusBetweenSC[1:], "\n")
	}
	if len(vo.Request) > 1 {
		result.Request = strings.Join(vo.Request[1:], "\n")
	}
	return
}

func ParseAiVeteranSummary(parseResult string) (baseInfo StatementDetailBaseInfo) {
	arr := strings.Split(parseResult, "\n")
	getVal := func(val string) string {
		res := strings.Split(val, ":")
		if len(res) > 1 {
			return strings.TrimSpace(res[1])
		}
		return ""
	}
	for _, v := range arr {
		v = strings.ReplaceAll(v, "** ", "")
		if baseInfo.FullName == "" && strings.Index(strings.ToLower(v), strings.ToLower("Full Name:")) >= 0 {
			baseInfo.FullName = getVal(v)
		} else if baseInfo.CaseId == 0 && strings.Index(strings.ToLower(v), strings.ToLower("Unique ID")) >= 0 {
			str := getVal(v)
			val, _ := strconv.ParseInt(str, 10, 32)
			baseInfo.CaseId = int32(val)
		} else if baseInfo.BranchOfService == "" && strings.Index(strings.ToLower(v), strings.ToLower("Branch of Service")) >= 0 {
			baseInfo.BranchOfService = getVal(v)
		} else if baseInfo.YearsOfService == "" && strings.Index(strings.ToLower(v), strings.ToLower("Years of Service")) >= 0 {
			baseInfo.YearsOfService = getVal(v)
		} else if baseInfo.RetiredFromService == "" && strings.Index(strings.ToLower(v), strings.ToLower("Retirement status")) >= 0 {
			baseInfo.RetiredFromService = getVal(v)
		} else if baseInfo.Deployments == "" && strings.Index(strings.ToLower(v), strings.ToLower("Deployments")) >= 0 {
			baseInfo.Deployments = getVal(v)
		} else if baseInfo.MaritalStatus == "" && strings.Index(strings.ToLower(v), strings.ToLower("Marital Status")) >= 0 {
			baseInfo.MaritalStatus = getVal(v)
		} else if baseInfo.Children == "" && strings.Index(strings.ToLower(v), strings.ToLower("Children:")) >= 0 {
			baseInfo.Children = getVal(v)
		} else if baseInfo.OccupationInService == "" && strings.Index(strings.ToLower(v), strings.ToLower("Occupation in service")) >= 0 {
			baseInfo.OccupationInService = getVal(v)
		}
	}
	return
}
