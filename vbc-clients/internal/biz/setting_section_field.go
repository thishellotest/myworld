package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type SettingSectionVo struct {
	Sections []SettingSection `json:"sections"`
}

type SettingSection struct {
	SectionLabel string                `json:"section_label"`
	SectionName  string                `json:"section_name"`
	Left         SettingSectionElement `json:"left"`
	Right        SettingSectionElement `json:"right"`
}

type SettingSectionElement struct {
	Fields []SettingSectionField `json:"fields"`
}

type SettingSectionField struct {
	FieldName string `json:"field_name"`
}

func (c *SettingSectionField) ToDetailField(typeFieldStruct *TypeFieldStruct, tData *TData, StatementUsecase *StatementUsecase) (detailField *DetailField, err error) {
	if typeFieldStruct == nil {
		return nil, errors.New("typeFieldStruct is nil")
	}
	//if tData == nil {
	//	return nil, errors.New("tData is nil")
	//}

	field := typeFieldStruct.GetByFieldName(c.FieldName)
	if field == nil {
		return nil, errors.New("field is nil")
	}

	detailField = &DetailField{}
	detailField.FieldLabel = field.FieldLabel
	detailField.FieldName = field.FieldName
	detailField.FieldType = field.FieldType
	if tData != nil {
		if false && field.IsCondition == Field_IsCondition_Yes { // todo:lgl 暂不支持此方法
			detailField.Value = CaseClaimsDivide(tData.CustomFields.TextValueByNameBasic(field.FieldName))
		} else {
			detailField.Value = tData.CustomFields.ValueByName(field.FieldName)
			if field.Kind == Kind_users && field.FieldName == UserFieldName_MailPassword {
				val, _ := DecryptSensitive(InterfaceToString(detailField.Value))
				detailField.Value = val
			} else if field.FieldName == FieldName_itf_expiration {
				itfExpiration := tData.CustomFields.GetByName(FieldName_itf_expiration)
				if itfExpiration != nil {
					detailField.Extend = itfExpiration.Extend
				}
			} else if field.FieldName == FieldName_statements {
				//statements := tData.CustomFields.TextValueByNameBasic(FieldName_statements)
				//conditions, err := SplitCaseStatements(statements)
				//if err != nil {
				//	lib.DPrintln(err)
				//}
				//detailField.Extend = conditions

				if tData != nil {
					conditions, err := StatementUsecase.GetCaseStatementExtend(*tData)
					if err != nil {
						lib.DPrintln(err)
					}
					detailField.Extend = conditions
				}

			}
		}
	}
	return detailField, err
}

var FormColumnAccountInfo = `
{
	"sections": [
	{
		"section_label": "User Information",
		"section_name": "user_information",
		"left": {
			"fields": [
			
			{"field_name": "timezone_id"}
			]
		}
	}
	]
}`

var FormColumnCondition = `
{
	"sections": [
	{
		"section_label": "",
		"section_name": "information",
		"left": {
			"fields": [
			{"field_name": "` + ConditionFieldName_condition_name + `"},
			{"field_name": "` + ConditionFieldName_condition_category_id + `"}
			]
		}
	}
	]
}`

var FormColumnConditionSecondary = `
{
	"sections": [
	{
		"section_label": "",
		"section_name": "information",
		"left": {
			"fields": [
{"field_name": "` + ConditionFieldName_secondary_type + `"},
			{"field_name": "` + ConditionFieldName_condition_name + `"}
			
			]
		}
	}
	]
}`

var FormColumnFilter = `
{
	"sections": [
	{
		"section_label": "",
		"section_name": "information",
		"left": {
			"fields": [
{"field_name": "` + Filter_FieldName_filter_name + `"}
			
			
			]
		}
	}
	]
}`

var Form_mgmt_user = `{
	"sections": [
	{
		"section_label": "Account Details",
		"section_name": "Account Details",
		"left": {
			"fields": [
{"field_name": "` + UserFieldName_first_name + `"},
{"field_name": "` + UserFieldName_last_name + `"},
{"field_name": "` + UserFieldName_title + `"},
{"field_name": "` + UserFieldName_email + `"},
{"field_name": "` + UserFieldName_mobile + `"},
{"field_name": "` + User_FieldName_timezone_id + `"},
{"field_name": "` + User_FieldName_profile_gid + `"},
{"field_name": "` + UserFieldName_role_gid + `"},
{"field_name": "` + UserFieldName_gender + `"},
{"field_name": "` + UserFieldName_permissions + `"}
			]
		}
	},
{
		"section_label": "Email Outbox Settings",
		"section_name": "Email Outbox Settings",
		"left": {
			"fields": [
{"field_name": "` + UserFieldName_MailSender + `"},
{"field_name": "` + UserFieldName_MailPassword + `"}
			]
		}
	},
{
		"section_label": "Dialpad Settings",
		"section_name": "Dialpad Settings",
		"left": {
			"fields": [
{"field_name": "` + UserFieldName_dialpad_userid + `"},
{"field_name": "` + UserFieldName_dialpad_phonenumber + `"}
			]
		}
	}
	]
}`

var Form_mgmt_attorney = `{
	"sections": [
	{
		"section_label": "Details",
		"section_name": "Details",
		"left": {
			"fields": [
{"field_name": "first_name"},
{"field_name": "last_name"},
{"field_name": "accreditation_date"},
{"field_name": "accreditation_number"},
{"field_name": "province"},
{"field_name": "city"},
{"field_name": "street"},
{"field_name": "zip_code"},
{"field_name": "company_name"},
{"field_name": "email"},
{"field_name": "ro_email"}


			]
		}
	}
	]
}`

/*
{"field_name": "first_name"},
			{"field_name": "last_name"},
			{"field_name": "email"},
*/

// DefaultClientCaseSectionConfig 此配置文件 后续放到数据库
var DefaultClientCaseSectionConfig = `
{
	"sections": [
	{
		"section_label": "Client Case Information",
		"section_name": "Client_Case_Information",
		"left": {
			"fields": [
			{"field_name": "user_gid"},
			{"field_name": "primary_vs"},
			{"field_name": "lead_co"},
			{"field_name": "deal_name"},
			{"field_name": "current_rating"},
			{"field_name": "retired"},
			{"field_name": "branch"},
			{"field_name": "source"},
			{"field_name": "created_by"},
			{"field_name": "client_gid"},
			{"field_name": "data_collection_folder"},
			{"field_name": "contract_source"},
			{"field_name": "attorney_uniqid"}]
		},
		"right": {
			"fields": [{"field_name":"collaborators"},
				{"field_name":"primary_cp"},
				{"field_name":"support_cp"},
				{"field_name":"stages"},
				{"field_name":"` + DataEntry_sys__due_date + `"},
				{"field_name":"effective_current_rating"},
				
				{"field_name":"referrer"},
				{"field_name":"new_rating"},
				{"field_name":"final_rating"},
				{"field_name": "modified_by"},
				{"field_name":"amount"},
				{"field_name":"am_invoice_amount"},
				{"field_name":"case_files_folder"},
				{"field_name":"pricing_version"}]
		}
	}

	,{
		"section_label": "Client Information",
		"section_name": "client_information",
		"left": {
			"fields": [
{"field_name": "email"}, 
{"field_name": "dob"}, 
{"field_name": "address"},
{"field_name": "apt_number"},
{"field_name": "timezone_id"},
{"field_name": "place_of_birth_city"},
{"field_name": "place_of_birth_country"}
]
		},
		"right": {
			"fields": [
{"field_name": "ssn"},
{"field_name": "phone"},
{"field_name": "state"},
{"field_name": "city"},
{"field_name": "zip_code"},
{"field_name": "place_of_birth_state_province"},
{"field_name": "current_occupation"}
]
		}
	}

,{
		"section_label": "Presumptive Information",
		"section_name": "presumptive_information",
		"left": {
			"fields": [
{"field_name": "year_entering_service"},
{"field_name": "answer_to_presumptive_questions"},
{"field_name": "agent_orange"},
{"field_name": "amyotrophic"},
{"field_name": "atomic_veterans"},
{"field_name": "burn_pits"},
{"field_name": "gulf_war"},
{"field_name": "illness_due"}
]
		},
		"right": {
			"fields": [
{"field_name": "year_separate_from_service"}
]
		}
	}

	,{
		"section_label": "Service Information",
		"section_name": "service_information",
		"left": {
			"fields": [
{"field_name": "service_years"}
]
		},
		"right": {
			"fields": [
{"field_name": "occupation_during_service"}
]
		}
	}
	,{
		"section_label": "Claims Information",
		"section_name": "claims_information",
		"left": {
			"fields": [
{"field_name": "contact_form"},
{"field_name": "21_22a_submitted"},
{"field_name": "foia"},
{"field_name": "sf_180"},
{"field_name": "disability_rating"},
{"field_name": "strs"},
{"field_name": "va_records"},
{"field_name": "va_healthcare"},
{"field_name": "tinnitus_nexus"},
{"field_name": "service_connections"}, 
{"field_name": "personal_statement_type"},
{"field_name": "personal_statement_manager"},
{"field_name": "personal_statement_password"},
{"field_name": "statements"}
]
		},
		"right": {
			"fields": [
{"field_name": "itf_expiration"},
{"field_name": "dd214"},
{"field_name": "rating_decision"},
{"field_name": "benefits_summary_letter"},
{"field_name": "private_records"},
{"field_name": "tdiu"},
{"field_name": "private_exams_needed"},
{"field_name": "previous_denials"}, 
{"field_name": "claims_next_round"},
{"field_name": "claims_online"},
{"field_name": "claims_supplemental"},
{"field_name": "claims_hlr"},
{"field_name": "description"}

]
		}
	}

]
}
`

var DefaultClientSectionConfig = `
{
	"sections": [
	{
		"section_label": "Client Information",
		"section_name": "Client_Information",
		"left": {
			"fields": [
			{"field_name": "user_gid"},
			{"field_name": "email"},
			{"field_name":"current_rating"},
			{"field_name": "retired"},
			{"field_name": "branch"}, 
			{"field_name": "dob"},
			{"field_name": "created_by"}, 
			{"field_name": "first_name"}, 
			{"field_name": "middle_name"},
			{"field_name": "source"},
			{"field_name": "client_type"},
			{"field_name": "referrer_gid"}
]
		},
		"right": {
			"fields": [
			{"field_name": "collaborators"}, 
			{"field_name": "phone"},
			{"field_name": "effective_current_rating"},
			{"field_name": "initial_claim_filing"},
			{"field_name": "pending_claims"},
			{"field_name": "ssn"},
			{"field_name": "referrer"},
			{"field_name": "modified_by"},
			{"field_name": "last_name"},
			{"field_name": "current_occupation"}]
		}
	}
	,{
		"section_label": "Address Information",
		"section_name": "address_information",
		"left": {
			"fields": [
			{"field_name": "address"}, 
			{"field_name": "apt_number"}, 
			{"field_name": "timezone_id"},
			{"field_name": "place_of_birth_city"}, 
			{"field_name": "place_of_birth_country"}]
		},
		"right": {
			"fields": [
			{"field_name": "state"}, 
			{"field_name": "city"},
			{"field_name": "zip_code"},
			{"field_name": "place_of_birth_state_province"}
			
]
		}
	}]
}
`

var DefaultTaskSectionConfig = `
{
	"sections": [
	{
		"section_label": "Task Information",
		"section_name": "Task_Information",
		"left": {
			"fields": [
			{"field_name": "user_gid"},
			{"field_name": "who_id_gid"},
			{"field_name":"subject"},
			{"field_name": "due_date"},
			{"field_name": "what_id_gid"}, 
			{"field_name": "status"},
			{"field_name": "priority"}]
		},
		"right": {
			"fields": [
			]
		}
	}]
}
`

type SettingSectionFieldUsecase struct {
	log                    *log.Helper
	CommonUsecase          *CommonUsecase
	conf                   *conf.Data
	FieldPermissionUsecase *FieldPermissionUsecase
}

func NewSettingSectionFieldUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldPermissionUsecase *FieldPermissionUsecase) *SettingSectionFieldUsecase {
	uc := &SettingSectionFieldUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		FieldPermissionUsecase: FieldPermissionUsecase,
	}

	return uc
}

func (c *SettingSectionFieldUsecase) GetSettingSectionsByKindForFormColumn(kind string, tProfile TData, recordHttpFormColumnRequestType string) (settingSectionVo SettingSectionVo, err error) {

	if kind != Kind_attorneys && kind != Kind_users && !IsCustomKind(kind) {
		return settingSectionVo, errors.New("GetSettingSectionsByKindForFormColumn: kind does not support")
	}
	var sectionVo SettingSectionVo
	if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_Account {
		err = json.Unmarshal([]byte(FormColumnAccountInfo), &sectionVo)
	} else if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_FormCondition {
		err = json.Unmarshal([]byte(FormColumnCondition), &sectionVo)
	} else if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_FormConditionSecondary {
		err = json.Unmarshal([]byte(FormColumnConditionSecondary), &sectionVo)
	} else if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_FormFilter {
		err = json.Unmarshal([]byte(FormColumnFilter), &sectionVo)
	} else if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_Form_mgmt_user {
		err = json.Unmarshal([]byte(Form_mgmt_user), &sectionVo)
	} else if recordHttpFormColumnRequestType == RecordHttpFormColumnRequest_Type_Form_mgmt_attorney {
		err = json.Unmarshal([]byte(Form_mgmt_attorney), &sectionVo)
	}

	// 进行字段权限过过虑
	fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, tProfile.Gid())
	if err != nil {
		return settingSectionVo, err
	}
	for _, v := range sectionVo.Sections {

		var settingSection SettingSection
		settingSection.SectionLabel = v.SectionLabel
		settingSection.SectionName = v.SectionName
		isOk := false

		for k1, v1 := range v.Left.Fields {
			//if kind == Kind_clients && from == Record_From_detail { // 特殊处理
			//	if v1.FieldName == "first_name" || v1.FieldName == "middle_name" {
			//		continue
			//	}
			//	if v1.FieldName == "last_name" {
			//		v1.FieldName = "full_name"
			//		v.Left.Fields[k1].FieldName = "full_name"
			//	}
			//}
			if kind == Kind_users && v1.FieldName == User_FieldName_profile_gid {
				if !IsAdminProfile(&tProfile) {
					continue
				}
			}

			fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v1.FieldName)
			if err != nil {
				return settingSectionVo, err
			}
			if fieldPermissionVo.CanShow() {
				isOk = true
				settingSection.Left.Fields = append(settingSection.Left.Fields, v.Left.Fields[k1])
			}
		}
		for k1, v1 := range v.Right.Fields {

			//if kind == Kind_clients && from == Record_From_detail { // 特殊处理
			//	if v1.FieldName == "first_name" || v1.FieldName == "middle_name" {
			//		continue
			//	}
			//	if v1.FieldName == "last_name" {
			//		v1.FieldName = "full_name"
			//		v.Right.Fields[k1].FieldName = "full_name"
			//	}
			//}
			if kind == Kind_users && v1.FieldName == User_FieldName_profile_gid {
				if !IsAdminProfile(&tProfile) {
					continue
				}
			}
			fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v1.FieldName)
			if err != nil {
				return settingSectionVo, err
			}
			if fieldPermissionVo.CanShow() {
				isOk = true
				settingSection.Right.Fields = append(settingSection.Right.Fields, v.Right.Fields[k1])
			}
		}
		if isOk {
			settingSectionVo.Sections = append(settingSectionVo.Sections, settingSection)
		}
	}
	return settingSectionVo, err
}

func (c *SettingSectionFieldUsecase) GetSettingSectionsByKind(kind string, profileGid string, from string) (settingSectionVo SettingSectionVo, err error) {

	if kind != Kind_client_cases && kind != Kind_clients && kind != Kind_client_tasks {
		return settingSectionVo, errors.New("GetSettingSectionsByKind: kind does not support")
	}
	var sectionVo SettingSectionVo
	// todo:lgl 后续此部分，通过kind，放到数据库管理
	if kind == Kind_client_cases {
		err = json.Unmarshal([]byte(DefaultClientCaseSectionConfig), &sectionVo)
	} else if kind == Kind_clients {
		err = json.Unmarshal([]byte(DefaultClientSectionConfig), &sectionVo)
	} else if kind == Kind_client_tasks {
		err = json.Unmarshal([]byte(DefaultTaskSectionConfig), &sectionVo)
	}

	// 进行字段权限过过虑
	fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, profileGid)
	if err != nil {
		return settingSectionVo, err
	}
	for _, v := range sectionVo.Sections {

		var settingSection SettingSection
		settingSection.SectionLabel = v.SectionLabel
		settingSection.SectionName = v.SectionName
		isOk := false

		for k1, v1 := range v.Left.Fields {

			if kind == Kind_clients && from == Record_From_detail { // 特殊处理
				if v1.FieldName == "first_name" || v1.FieldName == "middle_name" {
					continue
				}
				if v1.FieldName == "last_name" {
					v1.FieldName = "full_name"
					v.Left.Fields[k1].FieldName = "full_name"
				}
			}
			if kind == Kind_client_cases {
				// FieldName_statements 使用Statements Editor，不允许快捷添加了
				if from == Record_From_create || from == Record_From_edit {
					if v1.FieldName == DataEntry_sys__due_date || v1.FieldName == FieldName_statements || v1.FieldName == FieldName_ContractSource {
						continue
					}
				}
				if from == Record_From_edit {
					if v1.FieldName == FieldName_client_gid {
						continue
					}
				}
			}

			fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v1.FieldName)
			if err != nil {
				return settingSectionVo, err
			}
			if fieldPermissionVo.CanShow() {
				isOk = true
				settingSection.Left.Fields = append(settingSection.Left.Fields, v.Left.Fields[k1])
			}
		}
		for k1, v1 := range v.Right.Fields {

			if kind == Kind_clients && from == Record_From_detail { // 特殊处理
				if v1.FieldName == "first_name" || v1.FieldName == "middle_name" {
					continue
				}
				if v1.FieldName == "last_name" {
					v1.FieldName = "full_name"
					v.Right.Fields[k1].FieldName = "full_name"
				}
			}
			if kind == Kind_client_cases {
				// FieldName_statements 使用Statements Editor，不允许快捷添加了
				if from == Record_From_create || from == Record_From_edit {
					if v1.FieldName == DataEntry_sys__due_date || v1.FieldName == FieldName_statements || v1.FieldName == FieldName_ContractSource {
						continue
					}
				}
				if from == Record_From_edit {
					if v1.FieldName == FieldName_client_gid {
						continue
					}
				}
			}

			fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v1.FieldName)
			if err != nil {
				return settingSectionVo, err
			}
			if fieldPermissionVo.CanShow() {
				isOk = true
				settingSection.Right.Fields = append(settingSection.Right.Fields, v.Right.Fields[k1])
			}
		}
		if isOk {
			settingSectionVo.Sections = append(settingSectionVo.Sections, settingSection)
		}
	}

	return settingSectionVo, err
}
