package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

type DbqsUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	UserUsecase   *UserUsecase
}

func NewDbqsUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	UserUsecase *UserUsecase) *DbqsUsecase {
	uc := &DbqsUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		UserUsecase:   UserUsecase,
	}
	return uc
}

func (c *DbqsUsecase) BizLeadVSEmail(tClientCase *TData) (email string, err error) {

	email, err = c.LeadVSEmail(tClientCase)
	if err != nil {
		return "", err
	}
	if email == "" {
		email = "info@vetbenefitscenter.com"
	}
	return
}
func (c *DbqsUsecase) LeadVSEmail(tClientCase *TData) (email string, err error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	var primaryVSFullName string
	if tClientCase.CustomFields.TextValueByNameBasic(FieldName_email) == "lialing@foxmail.com" {
		primaryVSFullName = "Engineering Team"
	} else {
		primaryVSFullName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
		if primaryVSFullName == "" {
			return "", nil
		}
	}
	tUser, err := c.UserUsecase.GetByFullName(primaryVSFullName)
	if err != nil {
		return "", err
	}
	if tUser == nil {
		return "", nil
	}
	return tUser.CustomFields.TextValueByNameBasic("email"), nil
}

func (c *DbqsUsecase) BizLeadCPEmail(tClientCase *TData) (email string, err error) {

	email, err = c.LeadCPEmail(tClientCase)
	if err != nil {
		return "", err
	}
	if email == "" {
		email = "info@vetbenefitscenter.com"
	}
	return
}

func (c *DbqsUsecase) LeadCPEmail(tClientCase *TData) (email string, err error) {
	if tClientCase == nil {
		return "", errors.New("tClientCase is nil")
	}
	var primaryVSFullName string
	if tClientCase.CustomFields.TextValueByNameBasic(FieldName_email) == "lialing@foxmail.com" {
		primaryVSFullName = "Engineering Team"
	} else {
		primaryVSFullName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp)
		if primaryVSFullName == "" {
			return "", nil
		}
	}
	tUser, err := c.UserUsecase.GetByFullName(primaryVSFullName)
	if err != nil {
		return "", err
	}
	if tUser == nil {
		return "", nil
	}
	return tUser.CustomFields.TextValueByNameBasic("email"), nil
}

func FormatDate(dbStr string) (string, error) {
	if dbStr == "" {
		return "", nil
	}
	a, err := time.ParseInLocation(time.DateOnly, dbStr, configs.LoadLocation)
	if err != nil {
		return "", err
	}
	return a.Format("Jan. 2, 2006"), nil
}

func FormatLocation(city string, state string, zipCode string) string {
	return fmt.Sprintf("%s, %s %s", city, state, zipCode)
}

func FormatPrivateExamsNeeded(str string) string {
	strArr := strings.Split(str, "\n")
	res := ""
	for k, _ := range strArr {
		t := strings.TrimSpace(strArr[k])
		if t != "" {
			res += t + "\n"
		}
	}
	if res != "" {
		// 加一些空格处理，文字被截一半的情况
		res += "__________________________________________________________________"
	}
	return res
}

func FormatPrivateExamsNeededV2(str string) (res []string) {
	strArr := strings.Split(str, "\n")
	for k, _ := range strArr {
		t := strings.TrimSpace(strArr[k])
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func (c *DbqsUsecase) ReleaseOfInformationPrefillTags(tCase *TData, tClient *TData) (prefillTags lib.TypeList, err error) {
	if tCase == nil {
		return nil, errors.New("ReleaseOfInformationPrefillTags: tCase is nil.")
	}
	if tClient == nil {
		return nil, errors.New("ReleaseOfInformationPrefillTags: tClient is nil.")
	}
	clientFields := tClient.CustomFields
	caseFields := tCase.CustomFields
	dob, err := FormatDate(clientFields.TextValueByNameBasic("dob"))
	if err != nil {
		return nil, err
	}
	itf, err := FormatDate(clientFields.TextValueByNameBasic("itf_expiration"))
	if err != nil {
		return nil, err
	}
	privateExamsNeeded := FormatPrivateExamsNeeded(caseFields.TextValueByNameBasic("private_exams_needed"))
	prefillTags = lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      GenFullName(clientFields.TextValueByNameBasic(FieldName_first_name), clientFields.TextValueByNameBasic(FieldName_last_name)),
		},
		{
			"document_tag_id": "dob",
			"text_value":      dob, //"Jun. 8, 1984",
		},
		{
			"document_tag_id": "ssn",
			"text_value":      caseFields.TextValueByNameBasic("ssn"),
		},
		{
			"document_tag_id": "phone",
			"text_value":      caseFields.TextValueByNameBasic("phone"),
		},
		{
			"document_tag_id": "email",
			"text_value":      caseFields.TextValueByNameBasic("email"),
		},
		{
			"document_tag_id": "address",
			"text_value":      caseFields.TextValueByNameBasic("address"),
		},
		{
			"document_tag_id": "location",
			"text_value": FormatLocation(caseFields.TextValueByNameBasic("city"),
				caseFields.TextValueByNameBasic("state"),
				caseFields.TextValueByNameBasic("zip_code")), //"Escondido, CA 92000",
		},
		{
			"document_tag_id": "itf",
			"text_value":      itf,
		},
		{
			"document_tag_id": "privateExamsNeeded",
			"text_value":      privateExamsNeeded, //"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)",
		},
	}
	return
}

func (c *DbqsUsecase) MedicalTeamFormsPrefillTags(tCase *TData, tClient *TData) (prefillTags lib.TypeList, boxSignTplId string, err error) {
	if tCase == nil {
		return nil, "", errors.New("MedicalTeamFormsPrefillTags: tCase is nil.")
	}
	if tClient == nil {
		return nil, "", errors.New("MedicalTeamFormsPrefillTags: tClient is nil.")
	}
	clientFields := tClient.CustomFields
	caseFields := tCase.CustomFields
	dob, err := FormatDate(clientFields.TextValueByNameBasic("dob"))
	if err != nil {
		return nil, "", err
	}
	itf, err := FormatDate(caseFields.TextValueByNameBasic("itf_expiration"))
	if err != nil {
		return nil, "", err
	}
	privateExamsNeeded := FormatPrivateExamsNeeded(caseFields.TextValueByNameBasic("private_exams_needed"))
	fullName := GenFullName(clientFields.TextValueByNameBasic(FieldName_first_name), clientFields.TextValueByNameBasic(FieldName_last_name))

	prefillTags = lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      fullName,
		},
		{
			"document_tag_id": "ClientNameSign",
			"text_value":      fullName,
		},
		{
			"document_tag_id": "dob",
			"text_value":      dob, //"Jun. 8, 1984",
		},
		{
			"document_tag_id": "ssn",
			"text_value":      caseFields.TextValueByNameBasic("ssn"),
		},
		{
			"document_tag_id": "phone",
			"text_value":      caseFields.TextValueByNameBasic("phone"),
		},
		{
			"document_tag_id": "email",
			"text_value":      caseFields.TextValueByNameBasic("email"),
		},
		{
			"document_tag_id": "address",
			"text_value":      caseFields.TextValueByNameBasic("address"),
		},
		{
			"document_tag_id": "location",
			"text_value": FormatLocation(caseFields.TextValueByNameBasic("city"),
				caseFields.TextValueByNameBasic("state"),
				caseFields.TextValueByNameBasic("zip_code")), //"Escondido, CA 92000",
		},
		{
			"document_tag_id": "itf",
			"text_value":      itf,
		},
		{
			"document_tag_id": "privateExamsNeeded",
			"text_value":      privateExamsNeeded, //"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)",
		},
	}
	boxSignTplId = "aaaf350c-9f4f-437c-a8ce-ba5602b360b3"
	return
}

func (c *DbqsUsecase) MedicalTeamFormsPrefillTagsV2(tCase *TData, tClient *TData) (prefillTags lib.TypeList, boxSignTplId string, err error) {
	if tCase == nil {
		return nil, "", errors.New("MedicalTeamFormsPrefillTags: tCase is nil.")
	}
	if tClient == nil {
		return nil, "", errors.New("MedicalTeamFormsPrefillTags: tClient is nil.")
	}
	clientFields := tClient.CustomFields
	caseFields := tCase.CustomFields
	dob, err := FormatDate(clientFields.TextValueByNameBasic("dob"))
	if err != nil {
		return nil, "", err
	}
	itf, err := FormatDate(caseFields.TextValueByNameBasic("itf_expiration"))
	if err != nil {
		return nil, "", err
	}
	privateExamsNeededs := FormatPrivateExamsNeededV2(caseFields.TextValueByNameBasic("private_exams_needed"))
	fullName := GenFullName(clientFields.TextValueByNameBasic(FieldName_first_name), clientFields.TextValueByNameBasic(FieldName_last_name))

	prefillTags = lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      fullName,
		},
		{
			"document_tag_id": "ClientNameSign",
			"text_value":      fullName,
		},
		{
			"document_tag_id": "dob",
			"text_value":      dob, //"Jun. 8, 1984",
		},
		{
			"document_tag_id": "ssn",
			"text_value":      caseFields.TextValueByNameBasic("ssn"),
		},
		{
			"document_tag_id": "phone",
			"text_value":      caseFields.TextValueByNameBasic("phone"),
		},
		{
			"document_tag_id": "email",
			"text_value":      caseFields.TextValueByNameBasic("email"),
		},
		{
			"document_tag_id": "address",
			"text_value":      caseFields.TextValueByNameBasic("address"),
		},
		{
			"document_tag_id": "location",
			"text_value": FormatLocation(caseFields.TextValueByNameBasic("city"),
				caseFields.TextValueByNameBasic("state"),
				caseFields.TextValueByNameBasic("zip_code")), //"Escondido, CA 92000",
		},
		{
			"document_tag_id": "itf",
			"text_value":      itf,
		},
		//{
		//	"document_tag_id": "privateExamsNeeded",
		//	"text_value":      privateExamsNeeded, //"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)",
		//},
	}
	i := 0
	boxSignTplId = "2aac3aed-358f-491a-a459-1b9d51bcc154"
	for _, v := range privateExamsNeededs {
		i++
		if i > 8 {
			boxSignTplId = "ab50be2f-f5fc-4362-b7bd-8a6f0b1e5e73"
		}
		if i > 20 {
			break
		}
		prefillTags = append(prefillTags, map[string]interface{}{
			"document_tag_id": fmt.Sprintf("privateExamsNeeded%d", i),
			"text_value":      v,
		})
	}

	return
}

func (c *DbqsUsecase) MedicalTeamFormsPrefillTagsWithoutTemplate(tCase *TData, tClient *TData) (createMedicalTeamFormVo CreateMedicalTeamFormVo, err error) {
	if tCase == nil {
		return createMedicalTeamFormVo, errors.New("MedicalTeamFormsPrefillTags: tCase is nil.")
	}
	if tClient == nil {
		return createMedicalTeamFormVo, errors.New("MedicalTeamFormsPrefillTags: tClient is nil.")
	}
	clientFields := tClient.CustomFields
	caseFields := tCase.CustomFields
	dob, err := FormatDate(clientFields.TextValueByNameBasic("dob"))
	if err != nil {
		return createMedicalTeamFormVo, err
	}
	itf, err := FormatDate(caseFields.TextValueByNameBasic("itf_expiration"))
	if err != nil {
		return createMedicalTeamFormVo, err
	}
	privateExamsNeededs := FormatPrivateExamsNeededV2(caseFields.TextValueByNameBasic("private_exams_needed"))
	fullName := GenFullName(clientFields.TextValueByNameBasic(FieldName_first_name), clientFields.TextValueByNameBasic(FieldName_last_name))

	createMedicalTeamFormVo.ClientName = fullName
	createMedicalTeamFormVo.Address = caseFields.TextValueByNameBasic("address")
	createMedicalTeamFormVo.Location = FormatLocation(caseFields.TextValueByNameBasic("city"),
		caseFields.TextValueByNameBasic("state"),
		caseFields.TextValueByNameBasic("zip_code"))
	createMedicalTeamFormVo.Dob = dob
	createMedicalTeamFormVo.Ssn = caseFields.TextValueByNameBasic("ssn")
	createMedicalTeamFormVo.Phone = caseFields.TextValueByNameBasic("phone")
	createMedicalTeamFormVo.Email = caseFields.TextValueByNameBasic("email")
	createMedicalTeamFormVo.Itf = itf
	createMedicalTeamFormVo.PrivateExamsNeededS = privateExamsNeededs

	return createMedicalTeamFormVo, nil
}
