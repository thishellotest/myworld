package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"sync"
	"vbc/configs"
	"vbc/lib"
	. "vbc/lib/builder"
)

const (
	FieldType_formula = "formula" // itf

	FieldType_text      = "text"
	FieldType_text_url  = "url"
	FieldType_multitext = "multitext"
	FieldType_email     = "email"
	FieldType_tel       = "tel"

	FieldType_dropdown      = "dropdown"
	FieldType_multidropdown = "multidropdown"
	FieldType_timestamp     = "timestamp"
	FieldType_decimal       = "decimal"
	FieldType_number        = "number"
	FieldType_date          = "date"
	FieldType_lookup        = "lookup"
	FieldType_multilookup   = "multilookup"
	FieldType_switch        = "switch"
)

const (
	ContractSource_VBC = "VBC"
	ContractSource_AM  = "AM"

	NewCaseDefaultContractSource = ContractSource_AM
)

const (
	CommonFieldName_common_jotform_ids = "common_jotform_ids"

	FileName_asana_task_gid = "asana_task_gid"
	FileName_asana_user_gid = "asana_user_gid"

	FieldName_statements      = "statements"
	FieldName_attorney_uniqid = "attorney_uniqid"

	FieldName_ContractSource           = "contract_source"
	FieldName_branch                   = "branch"
	FieldName_stages                   = "stages"
	FieldName_current_rating           = "current_rating"
	FieldName_effective_current_rating = "effective_current_rating"
	FieldName_new_rating               = "new_rating"
	FieldName_first_name               = "first_name"
	FieldName_middle_name              = "middle_name"
	FieldName_full_name                = "full_name"
	FieldName_last_name                = "last_name"
	FieldName_biz_deleted_at           = "biz_deleted_at"
	FieldName_active_duty              = "active_duty"
	FieldName_primary_vs               = "primary_vs"
	FieldName_primary_cp               = "primary_cp"
	FieldName_lead_co                  = "lead_co"
	FieldName_support_cp               = "support_cp"
	FieldName_case_files_folder        = "case_files_folder"
	FieldName_data_collection_folder   = "data_collection_folder"

	FieldName_personal_statement_manager  = "personal_statement_manager"
	FieldName_personal_statement_password = "personal_statement_password"

	FieldName_uniqcode       = "uniqcode"
	FieldName_timezone_id    = "timezone_id"
	FieldName_pending_claims = "pending_claims"

	FileName_client_cases_gid = "gid"
	FieldName_gid             = "gid"
	FieldName_updated_at      = "updated_at"
	FieldName_created_at      = "created_at"
	FieldName_deleted_at      = "deleted_at"

	FieldName_collaborators = "collaborators"

	FileName_user_gid = "gid"

	Client_FileName_gid = "gid"

	FieldName_source       = "source"
	FieldName_referrer     = "referrer"
	FieldName_referrer_gid = "referrer_gid"
	FieldName_user_gid     = "user_gid"
	FieldName_email        = "email"
	FieldName_phone        = "phone"
	FieldName_ssn          = "ssn"
	FieldName_dob          = "dob"
	FieldName_state        = "state"
	FieldName_city         = "city"
	FieldName_address      = "address"
	FieldName_apt_number   = "apt_number"
	FieldName_zip_code     = "zip_code"

	FieldName_place_of_birth_city           = "place_of_birth_city"
	FieldName_place_of_birth_state_province = "place_of_birth_state_province"
	FieldName_place_of_birth_country        = "place_of_birth_country"
	FieldName_current_occupation            = "current_occupation"

	FieldName_pricing_version   = "pricing_version"
	FieldName_s_pricing_version = "s_pricing_version"

	FieldName_client_gid = "client_gid"

	FieldName_deal_name = "deal_name"

	FieldName_is_primary_case = "is_primary_case"
	FieldName_amount          = "amount"
	FieldName_description     = "description"
	FieldName_itf_expiration  = "itf_expiration"

	FieldName_am_invoice_amount       = "am_invoice_amount"
	FieldName_personal_statement_type = "personal_statement_type"
)

const (
	Personal_statement_type_WorddocumentWithinBox = "Word document within Box"
	Personal_statement_type_Webform               = "Webform"
)

const (
	Is_primary_case_YES = 1
	Is_primary_case_NO  = 0
)

func IsNumberFieldType(fieldType string) bool {
	if fieldType == FieldType_number || fieldType == FieldType_timestamp || fieldType == FieldType_switch {
		return true
	}
	return false
}

func IsFloatFieldType(fieldType string) bool {
	if fieldType == FieldType_decimal {
		return true
	}
	return false
}

const (
	Field_OnlyReadAndWrite_Yes = 1

	Field_IsCondition_Yes = 1
	Field_IsCondition_No  = 0

	Field_IsCollaborator_Yes = 1

	Field_NoManageable_Yes = 1 // 不允许在字段列表中管理，例：gid, id, created_at, updated_at

	Field_IsRequired_Yes = 1

	Field_IsEnableColor_Yes = 1
	Field_IsEnableColor_No  = 0

	FieldPermissionLevel_OnlyReadAndWrite = 1 // 只允许读和写权限
	FieldPermissionLevel_OnlyRead         = 2 // 只允行读权限  在内部写数据

	Field_NoTimelines_Yes = 1 // 此字段不需要日志

	Field_SortType_Sort              = 0 // field options sort field value (Entered order)
	Field_SortType_AlphabeticalOrder = 1 // Alphabetical order
)

func IsTextField(fieldType string) bool {
	if fieldType == FieldType_text ||
		fieldType == FieldType_text_url ||
		fieldType == FieldType_multitext ||
		fieldType == FieldType_email ||
		fieldType == FieldType_tel {
		return true
	}
	return false
}

type FieldEntity struct {
	ID           int32 `gorm:"primaryKey"`
	Kind         string
	FieldName    string
	FieldLabel   string
	FieldType    string
	RelaKind     string
	RelaName     string
	IsCondition  int
	NoManageable int // 不允许在字段列表中管理，例：gid, id, created_at, updated_at
	//OnlyReadAndWrite     int
	FieldPermissionLevel int
	IsCollaborator       int // 一个Kind 只允许一个， 要求必须设置为：field_type-multilookup; rela_kind-users; rela_name-full_name;
	IsRequired           int
	Dependence           string
	OptionGroupName      string // 有值时，使用 global_sets 的options
	IsEnableColor        int
	Tooltip              string
	ValidRegular         string // 有值时，说明需要验证
	NoTimelines          int
	SortType             int
	DeletedAt            int64
}

func (c *FieldEntity) GetIsRequired() bool {
	if c.IsRequired == Field_IsRequired_Yes {
		return true
	}
	return false
}

func (c *FieldEntity) GetIsEnableColor() bool {
	if c.IsEnableColor == Field_IsEnableColor_Yes {
		return true
	}
	return false
}

// modified_time zoho的更新时间
// created_time zoho的创建时间
var DoNoDisplayFieldNames = []string{DataEntry_gid, DataEntry_Incr_id_name, "uniqcode", "biz_deleted_at", "modified_time", "created_time", TaskFieldName_re_kind, TaskFieldName_se_module}

// 用户不能设置是否显示的列
var DoNoDisplayColumns = []string{}

// 必须返回到前端的字段列表
var MustReturnFieldNamesForRecords = []string{DataEntry_gid, DataEntry_sys__due_date}

var OnlyCRMFieldNames = []string{FieldName_collaborators, FieldName_middle_name, DataEntry_sys__due_date}

var DefaultDisplayFieldNamesForRecords = []string{FieldName_deal_name,
	FieldName_dob,
	FieldName_active_duty,
	FieldName_email,
	FieldName_full_name,
	FieldName_phone,
	FieldName_branch,
	FieldName_stages,
	DataEntry_sys__itf_formula,
	FieldName_itf_expiration,
	FieldName_user_gid,
	TaskFieldName_what_id_gid,
	TaskFieldName_due_date,
	TaskFieldName_status,
	TaskFieldName_priority,
	TaskFieldName_subject,
	DataEntry_updated_at,
}

var DefaultDisplayFieldNamesForSorts = []string{
	DataEntry_sys__itf_formula,
	FieldName_itf_expiration,
	FieldName_deal_name,
	FieldName_full_name,
	FieldName_stages,
	FieldName_branch,
	FieldName_user_gid,
	FieldName_email,
	FieldName_phone,
	FieldName_dob,
	FieldName_active_duty,
	DataEntry_updated_at,
}

func (c *FieldEntity) IsNoDisplayForUser() bool {
	if lib.InArray(c.FieldName, DoNoDisplayFieldNames) {
		return true
	}
	return false
}

func (c *FieldEntity) IsNoDisplayColumnsForUser() bool {
	if lib.InArray(c.FieldName, DoNoDisplayFieldNames) {
		return true
	}
	if lib.InArray(c.FieldName, DoNoDisplayColumns) {
		return true
	}
	return false
}

func (FieldEntity) TableName() string {
	return "fields"
}

const Pipelines_default = ""
const Pipelines_am = "am"
const Pipelines_vbc = "vbc"

func (c *FieldEntity) FieldToApi(FieldOptionUsecase *FieldOptionUsecase, log *log.Helper, tUser *TData) (fabField FabField) {

	//data := make(lib.TypeMap)
	//data.Set("module", c.Kind)
	//data.Set("field_name", c.FieldName)
	//data.Set("field_label", c.FieldLabel)
	//data.Set("field_type", c.FieldType)
	//data.Set("rela_module", c.RelaKind)
	//data.Set("rela_name", c.RelaName)

	fabField.Module = KindConvertToModule(c.Kind)
	fabField.FieldName = c.FieldName
	fabField.FieldLabel = c.FieldLabel
	fabField.FieldType = c.FieldType
	fabField.IsRequired = c.GetIsRequired()
	fabField.IsEnableColor = c.GetIsEnableColor()
	fabField.Tooltip = c.Tooltip
	fabField.CanManageColumn = true
	if c.FieldName == DataEntry_sys__due_date {
		fabField.CanManageColumn = false
	}
	if c.FieldType == FieldType_lookup {
		fabField.RelaModule = KindConvertToModule(c.RelaKind)
		fabField.RelaName = c.RelaName
	}

	//hasTestPermisstion := false
	//if tUser != nil {
	//	if tUser.Gid() != vbc_config.User_Dev_gid ||
	//		tUser.Gid() != vbc_config.User_Yannan_gid ||
	//		tUser.Gid() != vbc_config.User_Edward_gid ||
	//		tUser.Gid() != vbc_config.User_Lili_gid {
	//		hasTestPermisstion = true
	//	}
	//}

	if FieldOptionUsecase != nil {
		if c.FieldType == FieldType_dropdown || c.FieldType == FieldType_multidropdown {
			optionStruct, err := FieldOptionUsecase.CacheStructByKind(c.Kind)
			if err != nil {
				log.Error(err)
			} else {
				if optionStruct != nil {
					var destOptions []FabFieldOption
					optionList := optionStruct.AllByFieldName(*c)
					for _, v := range optionList {
						//if c.FieldName == FieldName_stages {
						//
						//	if pipelines == Pipelines_am {
						//		if !hasTestPermisstion {
						//			continue
						//		}
						//		if v.Pipelines != Pipelines_am {
						//			continue
						//		}
						//	} else if pipelines == Pipelines_vbc {
						//		if v.Pipelines != Pipelines_vbc {
						//			continue
						//		}
						//	} else if pipelines == Pipelines_default {
						//		if !hasTestPermisstion {
						//			if v.OptionValue == vbc_config.Stages_InformationIntake || v.OptionValue == vbc_config.Stages_ContractPending {
						//				continue
						//			}
						//		}
						//		if v.Pipelines == Pipelines_vbc || (hasTestPermisstion && (v.OptionValue == vbc_config.Stages_InformationIntake || v.OptionValue == vbc_config.Stages_ContractPending)) {
						//
						//		} else {
						//			continue
						//		}
						//	}
						//}
						destOptions = append(destOptions, v.FieldOptionToApi())
					}
					fabField.Options = destOptions
				}
			}
		}
	}
	return fabField
}

// DefaultValue 字段的默认值
func (c *FieldEntity) DefaultValue() interface{} {
	if c.FieldType == FieldType_number || c.FieldType == FieldType_decimal || c.FieldType == FieldType_switch {
		return nil
	} else if c.FieldType == FieldType_timestamp {
		return 0
	} else {
		return ""
	}
}

// TransToCorrectValueFormat 把输入的值转为正确的格式
func (c *FieldEntity) TransToCorrectValueFormat(val interface{}) interface{} {

	if c.FieldType == FieldType_number || c.FieldType == FieldType_decimal || c.FieldType == FieldType_switch {
		if val == nil || val == "" {
			return nil
		}
	}

	a := lib.InterfaceToString(val)
	if c.FieldType == FieldType_number || c.FieldType == FieldType_timestamp || c.FieldType == FieldType_switch {
		r, _ := strconv.ParseInt(a, 10, 32)
		return r
	} else if c.FieldType == FieldType_decimal {
		// 不能强行转float，会丢失精度
		return a
	} else {
		return a
	}
}

type FieldUsecase struct {
	CommonUsecase  *CommonUsecase
	GoCacheUsecase *GoCacheUsecase
	DBUsecase[FieldEntity]
	log *log.Helper
}

func NewFieldUsecase(logger log.Logger, CommonUsecase *CommonUsecase, GoCacheUsecase *GoCacheUsecase) *FieldUsecase {

	uc := &FieldUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		GoCacheUsecase: GoCacheUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *FieldUsecase) GetByFieldName(kind string, fieldName string) (*FieldEntity, error) {
	return c.GetByCond(Eq{"kind": kind, "field_name": fieldName, "deleted_at": 0})
}

const DataEntry_Incr_id_name = "id"
const DataEntry_gid = "gid"
const DataEntry_created_at = "created_at"
const DataEntry_updated_at = "updated_at"
const DataEntry_created_by = "created_by"
const DataEntry_modified_by = "modified_by"
const DataEntry_biz_deleted_at = "biz_deleted_at"

const DataEntry_sys__due_date = "sys__due_date"
const DataEntry_sys__itf_formula = "sys__itf_formula"

// DataEntry_user_gid 是一个普通的module都有这个字段，说明记录属于谁
const DataEntry_user_gid = "user_gid"

const DataEntry_deleted_at = "deleted_at"

func (c *FieldUsecase) KindCustomCondition(kind string) (res TypeFieldList, err error) {
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_Incr_id_name,
		FieldType:  FieldType_number,
		FieldLabel: "ID",
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_condition_name,
		FieldType:  FieldType_text,
		FieldLabel: "Condition Name",
		IsRequired: 1,
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_type,
		FieldType:  FieldType_text,
		FieldLabel: "Type",
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_condition_category_id,
		FieldType:  FieldType_lookup,
		FieldLabel: "Category Name",
	})

	res = append(res, &FieldEntity{
		Kind:         kind,
		FieldName:    DataEntry_created_by,
		FieldType:    FieldType_lookup,
		FieldLabel:   "Created By",
		RelaKind:     Kind_users,
		RelaName:     UserFieldName_fullname,
		NoManageable: Field_NoManageable_Yes,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_modified_by,
		FieldType:  FieldType_lookup,
		FieldLabel: "Modified By",
		RelaKind:   Kind_users,
		RelaName:   UserFieldName_fullname,
	})
	return
}

func (c *FieldUsecase) KindCustomConditionSecondary(kind string) (res TypeFieldList, err error) {
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_Incr_id_name,
		FieldType:  FieldType_number,
		FieldLabel: "ID",
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_condition_name,
		FieldType:  FieldType_text,
		FieldLabel: "Secondary Condition Name",
		IsRequired: 1,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_secondary_type,
		FieldType:  FieldType_dropdown,
		FieldLabel: "Type",
		IsRequired: 1,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  ConditionFieldName_type,
		FieldType:  FieldType_text,
		FieldLabel: "Type",
	})

	res = append(res, &FieldEntity{
		Kind:         kind,
		FieldName:    DataEntry_created_by,
		FieldType:    FieldType_lookup,
		FieldLabel:   "Created By",
		RelaKind:     Kind_users,
		RelaName:     UserFieldName_fullname,
		NoManageable: Field_NoManageable_Yes,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_modified_by,
		FieldType:  FieldType_lookup,
		FieldLabel: "Modified By",
		RelaKind:   Kind_users,
		RelaName:   UserFieldName_fullname,
	})
	return
}

func (c *FieldUsecase) KindCustomFilter(kind string) (res TypeFieldList, err error) {
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_Incr_id_name,
		FieldType:  FieldType_number,
		FieldLabel: "ID",
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  Filter_FieldName_filter_name,
		FieldType:  FieldType_text,
		FieldLabel: "Filter Name",
		IsRequired: 1,
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  Filter_FieldName_content,
		FieldType:  FieldType_text,
		FieldLabel: "Content",
		IsRequired: 1,
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  Filter_FieldName_table_type,
		FieldType:  FieldType_text,
		FieldLabel: "Table Type",
	})

	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  "user_gid",
		FieldType:  FieldType_lookup,
		FieldLabel: "User",
		RelaKind:   Kind_users,
		RelaName:   UserFieldName_fullname,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  "kind",
		FieldType:  FieldType_text,
		FieldLabel: "Kind",
	})

	res = c.HandleCustomCommonFields(kind, res)
	return
}
func (c *FieldUsecase) HandleCustomCommonFields(kind string, res TypeFieldList) TypeFieldList {

	res = append(res, &FieldEntity{
		Kind:         kind,
		FieldName:    DataEntry_created_by,
		FieldType:    FieldType_lookup,
		FieldLabel:   "Created By",
		RelaKind:     Kind_users,
		RelaName:     UserFieldName_fullname,
		NoManageable: Field_NoManageable_Yes,
	})
	res = append(res, &FieldEntity{
		Kind:       kind,
		FieldName:  DataEntry_modified_by,
		FieldType:  FieldType_lookup,
		FieldLabel: "Modified By",
		RelaKind:   Kind_users,
		RelaName:   UserFieldName_fullname,
	})

	res = append(res, &FieldEntity{
		Kind:         kind,
		FieldName:    DataEntry_biz_deleted_at,
		FieldType:    FieldType_timestamp,
		FieldLabel:   "Biz Deleted At",
		NoManageable: Field_NoManageable_Yes,
	})
	return res
}

func (c *FieldUsecase) ListByKind(kind string) (res TypeFieldList, err error) {
	err = c.CommonUsecase.DB().Where("deleted_at=0 and kind=?", kind).Find(&res).Error
	if err == nil {

		if kind == Kind_Custom_Condition {
			return c.KindCustomCondition(kind)
		} else if kind == Kind_Custom_ConditionSecondary {
			return c.KindCustomConditionSecondary(kind)
		} else if kind == Kind_Custom_Filter {
			return c.KindCustomFilter(kind)
		}

		// 后续所有系统字段应该都在这里处理
		res = append(res, &FieldEntity{
			Kind:       kind,
			FieldName:  DataEntry_Incr_id_name,
			FieldType:  FieldType_number,
			FieldLabel: "ID",
		})
		res = append(res, &FieldEntity{
			Kind:       kind,
			FieldName:  DataEntry_created_at,
			FieldType:  FieldType_timestamp,
			FieldLabel: "Created Time",
		})
		res = append(res, &FieldEntity{
			Kind:         kind,
			FieldName:    DataEntry_updated_at,
			FieldType:    FieldType_timestamp,
			FieldLabel:   "Modified Time",
			NoManageable: Field_NoManageable_Yes,
		})

		res = append(res, &FieldEntity{
			Kind:         kind,
			FieldName:    DataEntry_created_by,
			FieldType:    FieldType_lookup,
			FieldLabel:   "Created By",
			RelaKind:     Kind_users,
			RelaName:     UserFieldName_fullname,
			NoManageable: Field_NoManageable_Yes,
		})
		res = append(res, &FieldEntity{
			Kind:       kind,
			FieldName:  DataEntry_modified_by,
			FieldType:  FieldType_lookup,
			FieldLabel: "Modified By",
			RelaKind:   Kind_users,
			RelaName:   UserFieldName_fullname,
		})

		if kind != Kind_client_cases && kind != Kind_clients && kind != Kind_email_tpls { // todo:lgl 此处hack
			res = append(res, &FieldEntity{
				Kind:         kind,
				FieldName:    DataEntry_gid,
				FieldType:    FieldType_text,
				FieldLabel:   "Gid",
				NoManageable: Field_NoManageable_Yes,
			})
		}
		if kind == Kind_client_cases {
			res = append(res, &FieldEntity{
				Kind:         kind,
				FieldName:    DataEntry_sys__due_date,
				FieldType:    FieldType_date,
				FieldLabel:   "Due Date",
				NoManageable: Field_NoManageable_Yes,
				//NoTimelines:  Field_NoTimelines_Yes,
			})
			res = append(res, &FieldEntity{
				Kind:         kind,
				FieldName:    DataEntry_sys__itf_formula,
				FieldType:    FieldType_formula,
				FieldLabel:   "Days to ITF Expiration",
				NoManageable: Field_NoManageable_Yes,
				//NoTimelines:  Field_NoTimelines_Yes,
			})

		}
		if kind == Kind_users {
			res = append(res, &FieldEntity{
				Kind:         kind,
				FieldName:    DataEntry_biz_deleted_at,
				FieldType:    FieldType_timestamp,
				FieldLabel:   "Biz Deleted At",
				NoManageable: Field_NoManageable_Yes,
			})
		}
	}
	return
}

func (c *FieldUsecase) StructByKind(kind string) (*TypeFieldStruct, error) {
	list, err := c.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	a := &TypeFieldStruct{}
	a.Init(kind, list)
	return a, nil
}

func (c *FieldUsecase) CacheStructByKind(kind string) (*TypeFieldStruct, error) {
	key := fmt.Sprintf("%s%s", GOCACHE_PREFIX_field, kind)
	res, found := GoCacheGet[*TypeFieldStruct](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.StructByKind(kind)
	if err != nil {
		return nil, err
	}
	GoCacheSet(c.GoCacheUsecase, key, res, configs.CacheExpiredDurationDefault)
	return res, err
}

type TypeFieldList []*FieldEntity

func (c TypeFieldList) ToFieldNames() (r []string) {
	for _, v := range c {
		if v.FieldType == FieldType_formula {
			continue
		}
		r = append(r, v.FieldName)
	}
	return
}

type FieldDependList []FieldDepend

func (c FieldDependList) GetByFieldName(fieldName string) *FieldDepend {
	for k, v := range c {
		if v.FieldName == fieldName {
			r := c[k]
			return &r
		}
	}
	return nil
}

// GetFieldNameIsDependent 获取此字段被依赖的信息
func (c FieldDependList) GetFieldNameIsDependent(fieldName string) (r FieldDependList) {
	for k, v := range c {
		for _, v1 := range v.FieldDependItemList {
			if v1.FieldName == fieldName {
				r = append(r, c[k])
				break
			}
		}
	}
	return r
}

// GetFieldNamesIsDependent 获取这些字段被依赖的信息，如果重复就去重
func (c FieldDependList) GetFieldNamesIsDependent(fieldNames []string) (r FieldDependList) {

	res := make(map[string]FieldDepend)

	for _, fieldName := range fieldNames {
		for k, v := range c {
			for _, v1 := range v.FieldDependItemList {
				if v1.FieldName == fieldName {
					res[v.FieldName] = c[k]
					break
				}
			}
		}
	}
	for k, _ := range res {
		r = append(r, res[k])
	}
	return r
}

type FieldDepend struct {
	FieldName           string
	FieldDependItemList FieldDependItemList
}

type FieldDependItemList []FieldDependItem

type FieldDependItem struct {
	FieldName string
}

type TypeFieldStruct struct {
	Kind                string
	Records             TypeFieldList
	FieldNameIdx        map[string]*FieldEntity
	fieldDependList     FieldDependList
	onceFieldDependList sync.Once
}

func (c *TypeFieldStruct) GetFieldDependList() FieldDependList {
	c.onceFieldDependList.Do(func() {
		for _, v := range c.Records {
			if v.Dependence != "" {
				dependences := strings.Split(v.Dependence, "|")

				fieldDepend := FieldDepend{
					FieldName: v.FieldName,
				}
				for _, v1 := range dependences {
					fieldDependItem := FieldDependItem{
						FieldName: v1,
					}
					fieldDepend.FieldDependItemList = append(fieldDepend.FieldDependItemList, fieldDependItem)
				}
				c.fieldDependList = append(c.fieldDependList, fieldDepend)
			}
		}
	})
	return c.fieldDependList
}

func (c *TypeFieldStruct) Init(kind string, list TypeFieldList) {
	c.FieldNameIdx = make(map[string]*FieldEntity)
	c.Kind = kind
	c.Records = list
	for k, v := range list {
		c.FieldNameIdx[v.FieldName] = list[k]
	}
}

func (c *TypeFieldStruct) GetByFieldName(fieldName string) *FieldEntity {
	if _, ok := c.FieldNameIdx[fieldName]; ok {
		return c.FieldNameIdx[fieldName]
	}
	return nil
}

func (c *TypeFieldStruct) GetCollaborator() *FieldEntity {

	for k, v := range c.Records {
		if v.IsCollaborator == Field_IsCollaborator_Yes {
			return c.Records[k]
		}
	}
	return nil
}

func (c *TypeFieldStruct) FilterFieldName(fieldNames TypeFieldNameMaps) TypeFieldNameMaps {
	r := make(map[string]bool)
	for k, _ := range fieldNames {
		a := c.GetByFieldName(k)
		if a != nil {
			r[k] = true
		}
	}
	return r
}

type TypeFieldNameMaps map[string]bool

func (c TypeFieldNameMaps) ToSqlSelect() (selects []string) {
	if c != nil && len(c) > 0 {
		for k, _ := range c {
			selects = append(selects, k)
		}
	}
	return
}

func (c TypeFieldNameMaps) DeleteFieldName(fieldName string) {
	if c != nil {
		if _, ok := c[fieldName]; ok {
			delete(c, fieldName)
		}
	}
}
