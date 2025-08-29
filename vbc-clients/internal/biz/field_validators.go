package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"regexp"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

const (
	FieldValidator_IsRequired_Yes = 1
)

type FieldValidatorEntity struct {
	ID              int32 `gorm:"primaryKey"`
	FieldKind       string
	FieldName       string
	FieldValue      string
	DependFieldName string
	IsRequired      int
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       int64
}

func (FieldValidatorEntity) TableName() string {
	return "field_validators"
}

type FieldValidatorUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[FieldValidatorEntity]
	GoCacheUsecase *GoCacheUsecase
	FieldUsecase   *FieldUsecase
}

func NewFieldValidatorUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	GoCacheUsecase *GoCacheUsecase,
	FieldUsecase *FieldUsecase) *FieldValidatorUsecase {
	uc := &FieldValidatorUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		GoCacheUsecase: GoCacheUsecase,
		FieldUsecase:   FieldUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type TypeFieldValidatorList []FieldValidatorEntity

type TypeFieldValidatorStruct struct {
	Kind string
	List TypeFieldValidatorList
}

func (c *TypeFieldValidatorStruct) GetByFieldName(fieldName string, fieldValue string) (list TypeFieldValidatorList) {
	for k, v := range c.List {
		if v.FieldName == fieldName && v.FieldValue == fieldValue {
			list = append(list, c.List[k])
		}
	}
	return list
}

func (c *TypeFieldValidatorStruct) Init(kind string, list TypeFieldValidatorList) {
	c.Kind = kind
	c.List = list
}

func (c *FieldValidatorUsecase) ListByKind(kind string) (list TypeFieldValidatorList, err error) {
	err = c.CommonUsecase.DB().Where("field_kind=? and deleted_at=0", kind).
		Find(&list).Error
	return
}

func (c *FieldValidatorUsecase) StructByKind(kind string) (*TypeFieldValidatorStruct, error) {
	list, err := c.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	res := &TypeFieldValidatorStruct{}
	res.Init(kind, list)
	return res, nil
}

func (c *FieldValidatorUsecase) CacheStructByKind(kind string) (*TypeFieldValidatorStruct, error) {
	key := fmt.Sprintf("%s%s", GOCACHE_PREFIX_field_validator, kind)
	res, found := GoCacheGet[*TypeFieldValidatorStruct](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.StructByKind(kind)
	if err != nil {
		return nil, err
	}
	GoCacheSet[*TypeFieldValidatorStruct](c.GoCacheUsecase, key, res, configs.CacheExpiredDuration5Seconds)
	return res, nil
}

func (c *FieldValidatorUsecase) CacheFieldValidatorCenter(kind string) (fieldValidatorCenter FieldValidatorCenter, err error) {
	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
	fieldValidatorCenter.Kind = kind
	if err != nil {
		return fieldValidatorCenter, err
	}
	if fieldStruct == nil {
		return fieldValidatorCenter, errors.New("fieldStruct is nil")
	}
	fieldValidatorStruct, err := c.CacheStructByKind(kind)
	if err != nil {
		return fieldValidatorCenter, err
	}
	if fieldValidatorStruct == nil {
		return fieldValidatorCenter, errors.New("fieldValidatorStruct is nil")
	}
	fieldValidatorCenter.FieldStruct = *fieldStruct
	fieldValidatorCenter.FieldValidatorStruct = *fieldValidatorStruct
	return fieldValidatorCenter, nil
}

type FieldValidatorCenter struct {
	Kind                 string
	FieldStruct          TypeFieldStruct
	FieldValidatorStruct TypeFieldValidatorStruct
}

type VerifyFailureResultList []VerifyFailureResultItem

// RemoveDuplicateResult 去除重复的结果
func (c VerifyFailureResultList) RemoveDuplicateResult() (r VerifyFailureResultList) {
	flag := make(map[string]bool)
	for k, v := range c {
		key := v.ModuleName + v.FieldName
		if _, ok := flag[key]; !ok {
			flag[key] = true
			r = append(r, c[k])
		}
	}
	return r
}

type VerifyFailureResultItem struct {
	ModuleName string `json:"module_name"`
	FieldName  string `json:"field_name"`
	Message    string `json:"message"`
}

// Verify dataEntry:数据的真实值, isVerified:说明此字的情况
func (c *FieldValidatorCenter) Verify(fieldName string, fieldValue string, dataEntry TypeDataEntry) (isVerified bool, message string, verifyFailureResultList VerifyFailureResultList, err error) {

	//lib.DPrintln("sss:", dataEntry)
	fieldEntity := c.FieldStruct.GetByFieldName(fieldName)
	if fieldEntity == nil {
		return false, "", nil, errors.New("Verify: fieldEntity is nil")
	}
	if fieldEntity.IsRequired == Field_IsRequired_Yes {
		if fieldValue == "" {
			item := VerifyFailureResultItem{
				ModuleName: KindConvertToModule(c.Kind),
				FieldName:  fieldName,
				Message:    fieldEntity.FieldLabel + " cannot be empty",
			}
			verifyFailureResultList = append(verifyFailureResultList, item)
			return false, item.Message, verifyFailureResultList, nil
		}
	}
	if fieldEntity.FieldType == FieldType_email {
		if fieldValue != "" && !lib.IsValidEmail(fieldValue) {
			item := VerifyFailureResultItem{
				ModuleName: KindConvertToModule(c.Kind),
				FieldName:  fieldName,
				Message:    fieldEntity.FieldLabel + " is not the correct email format",
			}
			verifyFailureResultList = append(verifyFailureResultList, item)
			return false, item.Message, verifyFailureResultList, nil
		}
	}
	if fieldEntity.FieldType == FieldType_text_url {
		if fieldValue != "" && !lib.IsValidURL(fieldValue) {
			item := VerifyFailureResultItem{
				ModuleName: KindConvertToModule(c.Kind),
				FieldName:  fieldName,
				Message:    fieldEntity.FieldLabel + " is not the correct URL format",
			}
			verifyFailureResultList = append(verifyFailureResultList, item)
			return false, item.Message, verifyFailureResultList, nil
		}
	}

	if fieldEntity.ValidRegular != "" && fieldValue != "" {
		emailRegex := fieldEntity.ValidRegular
		re := regexp.MustCompile(emailRegex)
		if !re.MatchString(fieldValue) {
			if fieldValue != "" && !lib.IsValidURL(fieldValue) {
				item := VerifyFailureResultItem{
					ModuleName: KindConvertToModule(c.Kind),
					FieldName:  fieldName,
					Message:    fieldEntity.FieldLabel + " is not verified by regular expression \"" + fieldEntity.ValidRegular + "\"",
				}
				verifyFailureResultList = append(verifyFailureResultList, item)
				return false, item.Message, verifyFailureResultList, nil
			}
		}
	}

	fieldValidatorList := c.FieldValidatorStruct.GetByFieldName(fieldName, fieldValue)

	isVerifyFailure := false
	for _, v := range fieldValidatorList {
		if v.IsRequired == FieldValidator_IsRequired_Yes {
			dataEntryMap := lib.TypeMap(dataEntry)
			dependFieldValue := dataEntryMap.GetString(v.DependFieldName)

			fieldEntity1 := c.FieldStruct.GetByFieldName(v.DependFieldName)

			if dependFieldValue == "" {
				isVerifyFailure = true
				item := VerifyFailureResultItem{
					ModuleName: KindConvertToModule(c.Kind),
					FieldName:  v.DependFieldName,
					Message:    "Requirements of " + fieldEntity.FieldLabel + ": " + fieldEntity1.FieldLabel + " cannot be empty",
				}
				verifyFailureResultList = append(verifyFailureResultList, item)
			}
		}
	}
	if isVerifyFailure {
		item := VerifyFailureResultItem{
			ModuleName: KindConvertToModule(c.Kind),
			FieldName:  fieldName,
			Message:    "Dependent fields cannot be empty",
		}
		verifyFailureResultList = append(verifyFailureResultList, item)
		return false, item.Message, verifyFailureResultList, nil
	}
	return true, "", verifyFailureResultList, nil
}
