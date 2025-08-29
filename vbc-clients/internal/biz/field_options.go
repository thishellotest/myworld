package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"sort"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

/*
更新field_options

UPDATE `field_options`
SET `option_sort` = 1000+CAST(SUBSTRING_INDEX(`option_label`, '.', 1) AS UNSIGNED)
WHERE `field_name` = 'stages';

UPDATE `field_options`
SET `option_label_alias` = CONCAT(

	LPAD(CAST(SUBSTRING_INDEX(`option_label`, '.', 1) AS UNSIGNED), 2, '0'),
	'.',
	SUBSTRING_INDEX(`option_label`, '.', -1)

)
WHERE `field_name` = 'stages';
*/
type FieldOptionEntity struct {
	ID          int32 `gorm:"primaryKey"`
	Kind        string
	Pipelines   string
	FieldName   string
	OptionValue string
	OptionLabel string
	OptionColor string
	OptionSort  int
	DeletedAt   int64
}

func (c *FieldOptionEntity) GetKey() string {
	return c.Kind + ":" + c.FieldName
}

func (FieldOptionEntity) TableName() string {
	return "field_options"
}

func (c *FieldOptionEntity) FieldOptionToApi() (fabFieldOption FabFieldOption) {
	fabFieldOption.OptionValue = c.OptionValue
	fabFieldOption.OptionLabel = c.OptionLabel
	fabFieldOption.OptionColor = c.OptionColor
	fabFieldOption.Pipelines = c.Pipelines
	return
}

type FieldOptionUsecase struct {
	CommonUsecase  *CommonUsecase
	GoCacheUsecase *GoCacheUsecase
	log            *log.Helper
	conf           *conf.Data
	DBUsecase[FieldOptionEntity]
}

func NewFieldOptionUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase, GoCacheUsecase *GoCacheUsecase) *FieldOptionUsecase {

	uc := &FieldOptionUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		GoCacheUsecase: GoCacheUsecase,
		conf:           conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *FieldOptionUsecase) ListByKind(kind string) (list TypeFieldOptionList, err error) {
	err = c.CommonUsecase.DB().Where("(kind=? or (kind=?)) and deleted_at=0", kind, Kind_global_sets).Find(&list).Error
	return
}

func (c *FieldOptionUsecase) GetByFieldName(kind string, fieldName string, value string) (*FieldOptionEntity, error) {
	return c.GetByCond(Eq{"kind": kind, "field_name": fieldName, "option_value": value, "deleted_at": 0})
}

func (c *FieldOptionUsecase) GetByEntity(entity FieldEntity, value string) (*FieldOptionEntity, error) {
	if entity.OptionGroupName != "" {
		return c.GetByCond(Eq{"kind": Kind_global_sets, "field_name": entity.OptionGroupName, "option_value": value, "deleted_at": 0})
	} else {
		return c.GetByCond(Eq{"kind": entity.Kind, "field_name": entity.FieldName, "option_value": value, "deleted_at": 0})
	}
}

func (c *FieldOptionUsecase) GetByEntityAndLabel(entity FieldEntity, label string) (*FieldOptionEntity, error) {
	if entity.OptionGroupName != "" {
		return c.GetByCond(Eq{"kind": Kind_global_sets, "field_name": entity.OptionGroupName, "option_label": label, "deleted_at": 0})
	} else {
		return c.GetByCond(Eq{"kind": entity.Kind, "field_name": entity.FieldName, "option_label": label, "deleted_at": 0})
	}
}

func (c *FieldOptionUsecase) StructByKind(kind string) (*TypeFieldOptionStruct, error) {
	list, err := c.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	res := &TypeFieldOptionStruct{}
	res.Init(kind, list)
	return res, nil
}

func (c *FieldOptionUsecase) CacheStructByKind(kind string) (*TypeFieldOptionStruct, error) {
	key := fmt.Sprintf("%s%s", GOCACHE_PREFIX_field_option, kind)
	res, found := GoCacheGet[*TypeFieldOptionStruct](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.StructByKind(kind)
	if err != nil {
		return nil, err
	}
	GoCacheSet[*TypeFieldOptionStruct](c.GoCacheUsecase, key, res, configs.CacheExpiredDuration5Seconds)
	return res, nil
}

type TypeFieldOptionList []*FieldOptionEntity

func (c TypeFieldOptionList) GetByLabel(optionLabel string) *FieldOptionEntity {

	for k, v := range c {
		if v.OptionLabel == optionLabel {
			return c[k]
		}
	}
	return nil
}

func (c TypeFieldOptionList) GetByValue(optionValue string) *FieldOptionEntity {

	for k, v := range c {
		if v.OptionValue == optionValue {
			return c[k]
		}
	}
	return nil
}

type TypeFieldOptionStruct struct {
	Kind         string
	List         TypeFieldOptionList
	FieldNameIdx map[string]TypeFieldOptionList
}

func (c *TypeFieldOptionStruct) Init(kind string, list TypeFieldOptionList) {
	c.FieldNameIdx = make(map[string]TypeFieldOptionList)
	c.Kind = kind
	c.List = list
	for k, v := range list {
		c.FieldNameIdx[v.GetKey()] = append(c.FieldNameIdx[v.GetKey()], list[k])
	}
}

func (c *TypeFieldOptionStruct) AllByFieldName(field FieldEntity) TypeFieldOptionList {

	var key string
	if field.OptionGroupName != "" {
		key = fmt.Sprintf("%s:%s", Kind_global_sets, field.OptionGroupName)
	} else {
		key = fmt.Sprintf("%s:%s", field.Kind, field.FieldName)
	}

	if _, ok := c.FieldNameIdx[key]; ok {
		optionList := c.FieldNameIdx[key]

		if field.SortType == Field_SortType_Sort {
			sort.Slice(optionList, func(i, j int) bool {
				return optionList[i].OptionSort < optionList[j].OptionSort
			})
		} else {
			sort.Slice(optionList, func(i, j int) bool {
				return optionList[i].OptionLabel < optionList[j].OptionLabel
			})
		}
		return optionList
	}
	return nil
}
