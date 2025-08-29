package biz

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"vbc/configs"
	"vbc/lib"
)

type TField struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	DisplayValue *string           `json:"display_value"`
	Value        string            `json:"value"` // 数据库存储的值， 把NumberValue与TextValue都放在这里
	NumberValue  *int32            `json:"number_value"`
	TextValue    *string           `json:"text_value"`   // DecimalValue
	MultiValues  TFieldMultiValues `json:"multi_values"` // multilookup 使用
	Extend       interface{}       `json:"extend"`
}

const TFieldExtendForSysDueDate_Badge_red = "red"
const TFieldExtendForSysDueDate_Badge_green = "green"
const TFieldExtendForSysDueDate_Badge_warn = "warn"

type TFieldExtendForSysDueDate struct {
	Badge     string `json:"badge"` //
	FontColor string `json:"font_color"`
	BgColor   string `json:"bg_color"`
	Label     string `json:"label"`
	Value     string `json:"value"`
}

func SysDueDateLabel(dueDate time.Time) string {
	dueDate = dueDate.In(configs.GetVBCDefaultLocation())
	currentDate := time.Now().In(configs.GetVBCDefaultLocation())
	if dueDate.Format("2006") == currentDate.Format("2006") { // 同一年
		return dueDate.Format(configs.TimeFormatDateThisYear)
	} else {
		return dueDate.Format(configs.TimeFormatDate)
	}
}

type SysItfFormulaColorConfig struct {
	LeftDays int
	Extend   TFieldExtendForSysDueDate
}

var DefaultSysItfFormulaColorConfig = SysItfFormulaColorConfig{
	Extend: TFieldExtendForSysDueDate{
		FontColor: "#008000E6",
		BgColor:   "#008000E6",
	},
}
var SysItfFormulaColorConfigs = []SysItfFormulaColorConfig{ // 注意必须保证LeftDays的顺序
	{
		LeftDays: 0, //  **Expired** | Red | #FF0000
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#FF0000",
			BgColor:   "#FF0000E6",
		},
	}, {
		LeftDays: 3, // | **3 days or less** | Orange-Red | #FF4500
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#FF0000", //
			BgColor:   "#FF4500E6",
		},
	}, {
		LeftDays: 10, // | **10 days or less** | Orange |
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#FF4500",
			//BgColor:   "#FFA500",
		},
	}, {
		LeftDays: 20, //| **20 days or less** | Amber | #FFBF00
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#FFA500",
			BgColor:   "#FFBF00E6",
		},
	}, {
		LeftDays: 30, // **30 days or less** | Yellow | #FFFF00
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#FFA500", // FFD700
			BgColor:   "#FFD700E6",
		},
	}, {
		LeftDays: 60, //| **60 days or less** | Light Green | #90EE90
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#D4AF37", // 2E8B57
			BgColor:   "#90EE90E6",
		},
	}, {
		LeftDays: 90, // **90 days or less** | Green | #008000
		Extend: TFieldExtendForSysDueDate{
			FontColor: "#008000",
			BgColor:   "#008000E6",
		},
	},
}

func GetSysItfFormulaColorConfig(val int) TFieldExtendForSysDueDate {
	for k, v := range SysItfFormulaColorConfigs {
		if val <= v.LeftDays {
			return SysItfFormulaColorConfigs[k].Extend
		}
	}
	return DefaultSysItfFormulaColorConfig.Extend
}

// GenTFieldExtendForSysItfFormula val e.g.: "" "0" "-1"
func GenTFieldExtendForSysItfFormula(val string) (*TFieldExtendForSysDueDate, error) {
	if val == "" {
		return nil, nil
	}

	newVal, _ := strconv.ParseInt(val, 10, 32)

	extend := GetSysItfFormulaColorConfig(int(newVal))
	extend.Value = val
	extend.Label = val
	return &extend, nil
}

// GenTFieldExtendForSysDueDate val e.g.: 2025-01-03
func GenTFieldExtendForSysDueDate(val string) (*TFieldExtendForSysDueDate, error) {
	if val == "" {
		return nil, nil
	}
	currentTime := time.Now()
	currentTime = currentTime.In(configs.GetVBCDefaultLocation())
	currentDate := currentTime.Format(time.DateOnly)
	if currentDate == val {
		return &TFieldExtendForSysDueDate{
			Badge:     TFieldExtendForSysDueDate_Badge_warn,
			FontColor: "#C46A0F",
			BgColor:   "#FDEDDD",
			Label:     "Today",
			Value:     val,
		}, nil
	} else {
		dueDate, err := time.ParseInLocation(time.DateOnly, val, configs.GetVBCDefaultLocation())
		if err != nil {
			return nil, err
		}
		currentDataTime, err := time.ParseInLocation(time.DateOnly, currentDate, configs.GetVBCDefaultLocation())
		if err != nil {
			return nil, err
		}
		if dueDate.Before(currentDataTime) {
			return &TFieldExtendForSysDueDate{
				Badge:     TFieldExtendForSysDueDate_Badge_red,
				FontColor: "#FF5D5A",
				BgColor:   "#FFECEC",
				Label:     SysDueDateLabel(dueDate),
				Value:     val,
			}, nil
		} else {
			return &TFieldExtendForSysDueDate{
				Badge:     TFieldExtendForSysDueDate_Badge_green,
				FontColor: "#12AA67",
				BgColor:   "#DDF8EC",
				Label:     SysDueDateLabel(dueDate),
				Value:     val,
			}, nil
		}
	}
}

type TFieldMultiValues []TFieldMultiValue

type TFieldMultiValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type TFields []TField
type TData struct {
	CustomFields   TFields           `json:"custom_fields"`
	CacheRelaKinds map[string]*TData `json:"-"`
	Kind           string            `json:"-"`
}

func (c *TData) HandleTDataTimezone(TimezonesUsecase *TimezonesUsecase, fieldStruct TypeFieldStruct, userFacade *UserFacade) error {
	if TimezonesUsecase == nil {
		return errors.New("HandleTDataTimezone: TimezonesUsecase is nil")
	}
	timezonesEntity, err := userFacade.GetTimezonesEntity(TimezonesUsecase)
	if err != nil {
		return err
	}
	if timezonesEntity == nil {
		return errors.New("HandleTDataTimezone: timezonesEntity is nil")
	}

	for k, v := range c.CustomFields {
		fieldEntity := fieldStruct.GetByFieldName(v.Name)
		if fieldEntity == nil {
			return errors.New(v.Name + ": HandleTDataTimezone: fieldEntity is nil")
		}
		if fieldEntity.FieldType == FieldType_date {
			displayValue := v.DisplayValue
			if displayValue != nil && *displayValue != "" {
				timeParse, _ := time.Parse(time.DateOnly, *displayValue)
				newDisplayValue := timeParse.Format("Jan 02, 2006")
				c.CustomFields[k].DisplayValue = &newDisplayValue
			}
		} else if fieldEntity.FieldType == FieldType_timestamp {
			displayValue := v.DisplayValue
			if displayValue != nil && *displayValue != "" {
				timeParse, _ := time.Parse(time.RFC3339, *displayValue)
				la, err := time.LoadLocation(timezonesEntity.CodeValue)
				if err != nil {
					return err
				}
				newDisplayValue := timeParse.In(la).Format("Jan 02, 2006 03:04 PM")
				//2024-12-09T07:05:21Z
				c.CustomFields[k].DisplayValue = &newDisplayValue
			}
		}
	}
	return nil
}

func (c *TData) RelaData(bUsecase *BUsecase, fieldName string) (*TData, error) {
	if bUsecase == nil {
		return nil, errors.New("RelaData: bUsecase is nil")
	}
	fieldStruct, err := bUsecase.FieldUsecase.CacheStructByKind(c.Kind)
	if err != nil {
		return nil, err
	}
	if fieldStruct == nil {
		return nil, errors.New("RelaData: fieldStruct is nil")
	}
	fieldEntity := fieldStruct.GetByFieldName(fieldName)
	if fieldEntity == nil {
		return nil, errors.New("RelaData: fieldEntity is nil")
	}

	if fieldEntity.RelaName == "" {
		return nil, errors.New(fieldEntity.FieldName + ":The field is not rela field")
	}
	if c.CacheRelaKinds == nil {
		c.CacheRelaKinds = make(map[string]*TData)
	}
	fieldVal := c.CustomFields.TextValueByNameBasic(fieldEntity.FieldName)
	if fieldVal == "" {
		return nil, nil
	}

	key := fmt.Sprintf("%s:%s", fieldEntity.RelaKind, fieldVal)
	if _, ok := c.CacheRelaKinds[key]; ok {
		return c.CacheRelaKinds[key], nil
	}
	tRelaData, err := bUsecase.TUsecase.DataByGid(fieldEntity.RelaKind, fieldVal)
	if err != nil {
		return nil, err
	}
	c.CacheRelaKinds[key] = tRelaData
	return tRelaData, nil
}

func (c *TData) Gid() string {
	return c.CustomFields.TextValueByNameBasic("gid")
}
func (c *TData) Id() int32 {
	return c.CustomFields.NumberValueByNameBasic("id")
}

func (c *TData) CreatedAt() int32 {
	return c.CustomFields.NumberValueByNameBasic(DataEntry_created_at)
}

func (c *TData) UpdatedAt() int32 {
	return c.CustomFields.NumberValueByNameBasic(DataEntry_updated_at)
}

func (c TFields) DisplayValueByName(name string) *string {
	for k, v := range c {
		if v.Name == name {
			return c[k].DisplayValue
		}
	}
	return nil
}

func (c TFields) DisplayValueByNameBasic(name string) string {
	a := c.DisplayValueByName(name)
	if a != nil {
		return *a
	}
	return ""
}

func (c TFields) NumberValueByName(name string) *int32 {
	for k, v := range c {
		if v.Name == name {
			return c[k].NumberValue
		}
	}
	return nil
}

// NumberValueByNameBasic id是number所有可以获取到
func (c TFields) NumberValueByNameBasic(name string) int32 {
	a := c.NumberValueByName(name)
	if a != nil {
		return *a
	}
	return 0
}

func (c TFields) TextValueByName(name string) *string {
	for k, v := range c {
		if v.Name == name {
			return c[k].TextValue
		}
	}
	return nil
}

func (c TFields) TFieldMultiValuesByName(name string) TFieldMultiValues {
	for k, v := range c {
		if v.Name == name {
			return c[k].MultiValues
		}
	}
	return nil
}

func (c TFields) ValueByName(name string) any {
	for k, v := range c {
		if v.Name == name {
			if v.Type == FieldType_dropdown || v.Type == FieldType_lookup {
				return c._valueOpByName(name)
			} else if v.Type == FieldType_multilookup || v.Type == FieldType_multidropdown {
				return c.TFieldMultiValuesByName(name)
			}
			return c[k].Value
		}
	}
	return ""
}

func (c TFields) GetByName(name string) *TField {
	for k, v := range c {
		if v.Name == name {
			return &c[k]
		}
	}
	return nil
}

type ValueOp struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (c TFields) _valueOpByName(name string) *ValueOp {
	label := c.DisplayValueByNameBasic(name)
	value := c.TextValueByNameBasic(name)
	if value != "" {
		return &ValueOp{
			Label: label,
			Value: value,
		}
	}

	return nil
}

// TextValueByNameBasic id是number此方法不能获取
func (c TFields) TextValueByNameBasic(name string) string {
	a := c.TextValueByName(name)
	if a != nil {
		return *a
	}
	return ""
}

func (c TFields) ToResponseUser(name string) *ResponseUser {
	val := c.TextValueByNameBasic(name)
	if val != "" {
		displayValue := c.DisplayValueByNameBasic(name)
		return &ResponseUser{
			Gid:  val,
			Name: displayValue,
		}
	}
	return nil
}

func (c TFields) ToResponseRelatedRecord(kindEntity KindEntity, name string, tData *TData) *ResponseRelatedRecord {
	val := c.TextValueByNameBasic(name)
	if val != "" {
		displayValue := ""
		if tData != nil {
			displayValue = tData.CustomFields.TextValueByNameBasic(kindEntity.PrimaryFieldName)
		}
		return &ResponseRelatedRecord{
			ModuleName:  kindEntity.ModuleName(),
			ModuleLabel: kindEntity.ModuleLabel(),
			Gid:         val,
			Name:        displayValue,
		}
	}
	return nil
}

func (c TFields) SetTextValueByName(name string, val *string) *string {
	for k, v := range c {
		if v.Name == name {
			c[k].TextValue = val
		}
	}
	return nil
}

func (c TFields) SetNumberValueByName(name string, val *int32) *string {
	for k, v := range c {
		if v.Name == name {
			c[k].NumberValue = val
		}
	}
	return nil
}

func (c TFields) ToDisplayMaps() lib.TypeMap {
	r := make(map[string]interface{})
	for _, v := range c {
		r[v.Name] = v.DisplayValue
	}
	return r
}

func (c TFields) ToApiMap() lib.TypeMap {
	r := make(map[string]interface{})
	for k, v := range c {
		bytes, _ := json.Marshal(c[k])
		var maps map[string]interface{}
		json.Unmarshal(bytes, &maps)
		r[v.Name] = maps
	}
	return r
}

func (c TFields) ToMaps() lib.TypeMap {
	r := make(map[string]interface{})
	for _, v := range c {
		if IsNumberFieldType(v.Type) {
			r[v.Name] = v.NumberValue
		} else if IsFloatFieldType(v.Type) {
			r[v.Name] = v.TextValue
		} else if v.Type == FieldType_multilookup {
			r[v.Name] = v.Value
		} else {
			r[v.Name] = v.TextValue
		}
	}
	return r
}

type TDataList []TData
