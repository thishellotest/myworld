package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

const TDefaultPageSize = 10
const TMaxPageSize = 1000

type TListRequest struct {
	Filter    TFilterVo `json:"filter"`
	TableType string    `json:"table_type"`
}

const (
	TFilterVo_Operator_AND = "AND"
	TFilterVo_Operator_OR  = "OR"
)

type TFilterVo struct {
	Operator string          `json:"operator"`
	Group    TListConditions `json:"group"`
	//Page     int
	//PageSize int
}

type TListConditions []TListCondition

type TListCondition struct {
	Comparator string `json:"comparator"`
	Field      TListConditionField
	Value      []interface{}
}

type TListConditionValueOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func (c *TListCondition) FirstValue() string {
	if len(c.Value) > 0 {
		return InterfaceToString(c.Value[0])
	}
	return ""
}

func (c *TListCondition) ValueForOptions() []TListConditionValueOption {
	if len(c.Value) > 0 {
		return lib.InterfaceToTDef[[]TListConditionValueOption](c.Value, nil)
	}
	return nil
}

func (c *TListCondition) SecondValue() string {
	if len(c.Value) > 1 {
		return InterfaceToString(c.Value[1])
	}
	return ""
}

type TListConditionField struct {
	FieldName string `json:"field_name"`
}

type TFilterUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	FieldUsecase       *FieldUsecase
	FieldOptionUsecase *FieldOptionUsecase
}

func NewTFilterUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase) *TFilterUsecase {
	uc := &TFilterUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		FieldUsecase:       FieldUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
	}
	return uc
}

func (c *TFilterUsecase) Do(kind string, tListRequest TListRequest, timezoneId string, aliasTableName string) (destCond Cond, err error) {

	fieldStruct, err := c.FieldUsecase.StructByKind(kind)
	if err != nil {
		return nil, err
	}
	if fieldStruct == nil {
		return nil, errors.New("fieldStruct is nil")
	}

	var conds []Cond
	for k, _ := range tListRequest.Filter.Group {
		cond, err := c.GenCond(fieldStruct, tListRequest.Filter.Group[k], timezoneId, aliasTableName)
		if err != nil {
			return nil, err
		}
		if cond != nil {
			conds = append(conds, cond)
		}
	}
	if len(conds) == 0 {
		return nil, nil
	}
	if tListRequest.Filter.Operator == TFilterVo_Operator_OR {
		return Or(conds...), nil
	} else {
		return And(conds...), nil
	}
}

func (c *TFilterUsecase) GenCond(fieldStruct *TypeFieldStruct, tListCondition TListCondition, timezoneId string, aliasTableName string) (Cond, error) {

	if fieldStruct == nil {
		return nil, errors.New("fieldStruct is nil")
	}

	fieldEntity := fieldStruct.GetByFieldName(tListCondition.Field.FieldName)
	if fieldEntity == nil {
		return nil, nil
	}
	if fieldEntity.FieldType == FieldType_text || fieldEntity.FieldType == FieldType_multitext ||
		fieldEntity.FieldType == FieldType_email ||
		fieldEntity.FieldType == FieldType_text_url ||
		fieldEntity.FieldType == FieldType_tel {
		return c.GenCondText(*fieldEntity, tListCondition, aliasTableName)
	} else if fieldEntity.FieldType == FieldType_decimal || fieldEntity.FieldType == FieldType_number {
		return c.GenCondDecimal(*fieldEntity, tListCondition, aliasTableName)
	} else if fieldEntity.FieldType == FieldType_lookup || fieldEntity.FieldType == FieldType_dropdown {
		return c.GenCondDropdown(*fieldEntity, tListCondition, aliasTableName)
	} else if fieldEntity.FieldType == FieldType_multilookup {
		return c.GenCondMultiDropdown(*fieldEntity, tListCondition, aliasTableName)
	} else if fieldEntity.FieldType == FieldType_timestamp {
		return c.GenCondTimestamp(*fieldEntity, tListCondition, timezoneId, aliasTableName)
	} else if fieldEntity.FieldType == FieldType_date || fieldEntity.FieldName == DataEntry_sys__itf_formula {
		return c.GenCondDate(*fieldEntity, tListCondition, aliasTableName)
	}
	return nil, nil
}

const (
	Comparator_eq           = "eq"
	Comparator_neq          = "neq"
	Comparator_is_empty     = "is_empty"
	Comparator_is_not_empty = "is_not_empty"
	Comparator_contains     = "contains"
	Comparator_not_contains = "not_contains"
	Comparator_lt           = "lt"
	Comparator_lte          = "lte"
	Comparator_gt           = "gt"
	Comparator_gte          = "gte"
	Comparator_between      = "between"
	Comparator_not_between  = "not_between"
)

func (c *TFilterUsecase) GenCondDecimal(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {
	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else if tListCondition.Comparator == Comparator_eq {
		return Eq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_lt {
		return Lt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_lte {
		return Lte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_gt {
		return Gt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_gte {
		return Gte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_between {
		return And(Gte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()},
			Lte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.SecondValue()}), nil
	} else if tListCondition.Comparator == Comparator_not_between {
		return Or(Lt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.FirstValue()},
			Gt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): tListCondition.SecondValue()}), nil
	} else {
		return nil, errors.New("GenCondDecimal Comparator is wrong")
	}
}

func ItfFormulaValueToItfExpirationValue(leftDays int32) string {
	currentTime := time.Now().In(configs.GetVBCDefaultLocation())
	currentTime = currentTime.AddDate(0, 0, int(leftDays))
	return currentTime.Format(time.DateOnly)
}

func (c *TFilterUsecase) GenCondDate(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {

	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else {

		fieldName := fieldEntity.FieldName
		firstValue := tListCondition.FirstValue()
		if fieldEntity.FieldName == DataEntry_sys__itf_formula {
			fieldName = FieldName_itf_expiration
			firstValue = ItfFormulaValueToItfExpirationValue(lib.StringToInt32(firstValue))
		}
		secondValue := ""
		if tListCondition.Comparator == Comparator_between || tListCondition.Comparator == Comparator_not_between {
			secondValue = tListCondition.SecondValue()
			secondValue = ItfFormulaValueToItfExpirationValue(lib.StringToInt32(secondValue))
		}

		if tListCondition.Comparator == Comparator_eq {
			return Eq{TidyTableFieldForSql(fieldName, aliasTableName): firstValue}, nil
		} else if tListCondition.Comparator == Comparator_lt {
			return Lt{TidyTableFieldForSql(fieldName, aliasTableName): firstValue}, nil
		} else if tListCondition.Comparator == Comparator_lte {
			return Lte{TidyTableFieldForSql(fieldName, aliasTableName): firstValue}, nil
		} else if tListCondition.Comparator == Comparator_gt {
			return Gt{TidyTableFieldForSql(fieldName, aliasTableName): firstValue}, nil
		} else if tListCondition.Comparator == Comparator_gte {
			return Gte{TidyTableFieldForSql(fieldName, aliasTableName): firstValue}, nil
		} else if tListCondition.Comparator == Comparator_between {
			return And(Gte{TidyTableFieldForSql(fieldName, aliasTableName): firstValue},
				Lte{TidyTableFieldForSql(fieldName, aliasTableName): secondValue}), nil
		} else if tListCondition.Comparator == Comparator_not_between {
			return Or(Lt{TidyTableFieldForSql(fieldName, aliasTableName): firstValue},
				Gt{TidyTableFieldForSql(fieldName, aliasTableName): secondValue}), nil
		} else {
			return nil, errors.New("GenCondDecimal Comparator is wrong")
		}
	}
}

func TransGenCondTimestamp(timezoneId string, value string) (time2 time.Time, err error) {
	if value == "" {
		return time2, nil
	}
	location := configs.GetVBCDefaultLocation()
	if timezoneId != "" {
		location, err = time.LoadLocation(timezoneId)
		if err != nil {
			return time2, err
		}
	}
	time2, err = time.ParseInLocation(time.DateOnly, value, location)
	return
}

func (c *TFilterUsecase) GenCondTimestamp(fieldEntity FieldEntity, tListCondition TListCondition, timezoneId string, aliasTableName string) (Cond, error) {
	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else if tListCondition.Comparator == Comparator_lt {
		time, err := TransGenCondTimestamp(timezoneId, tListCondition.FirstValue())
		if err != nil {
			return nil, err
		}
		return Lt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time.Unix()}, nil
	} else if tListCondition.Comparator == Comparator_gt {
		time, err := TransGenCondTimestamp(timezoneId, tListCondition.FirstValue())
		if err != nil {
			return nil, err
		}
		time = time.AddDate(0, 0, 1)
		return Gte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time.Unix()}, nil
	} else if tListCondition.Comparator == Comparator_between {
		time, err := TransGenCondTimestamp(timezoneId, tListCondition.FirstValue())
		if err != nil {
			return nil, err
		}
		time2, err := TransGenCondTimestamp(timezoneId, tListCondition.SecondValue())
		if err != nil {
			return nil, err
		}
		time2 = time2.AddDate(0, 0, 1)
		return And(Gte{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time.Unix()}, Lt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time2.Unix()}), nil
	} else if tListCondition.Comparator == Comparator_not_between {

		time, err := TransGenCondTimestamp(timezoneId, tListCondition.FirstValue())
		if err != nil {
			return nil, err
		}
		time2, err := TransGenCondTimestamp(timezoneId, tListCondition.SecondValue())
		if err != nil {
			return nil, err
		}
		time2 = time2.AddDate(0, 0, 1)
		// todo:lgl or可能有问题
		return Or(Lt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time.Unix()}, Gt{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): time2.Unix()}), nil
	} else {
		return nil, errors.New("GenCondDecimal Comparator is wrong")
	}
}

func (c *TFilterUsecase) GenCondDropdown(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {
	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else if tListCondition.Comparator == Comparator_eq {

		valueOptions := tListCondition.ValueForOptions()
		if len(valueOptions) > 0 {
			var values []string
			for _, v := range valueOptions {
				values = append(values, v.Value)
			}
			return In(TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName), values), nil
		}
		return nil, nil

	} else if tListCondition.Comparator == Comparator_neq {
		valueOptions := tListCondition.ValueForOptions()
		if len(valueOptions) > 0 {
			var values []string
			for _, v := range valueOptions {
				values = append(values, v.Value)
			}
			return NotIn(TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName), values), nil
		}
		return nil, nil
	} else {
		return nil, errors.New("GenCondDropdown Comparator is wrong")
	}
}

func (c *TFilterUsecase) GenCondMultiDropdown(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {
	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else if tListCondition.Comparator == Comparator_eq {

		valueOptions := tListCondition.ValueForOptions()
		if len(valueOptions) > 0 {
			var conds []Cond
			for _, v := range valueOptions {
				conds = append(conds, Like{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName), "," + v.Value + ","})
			}
			return Or(conds...), nil
		}
		return nil, nil

	} else if tListCondition.Comparator == Comparator_neq {
		valueOptions := tListCondition.ValueForOptions()
		if len(valueOptions) > 0 {
			var conds []Cond
			for _, v := range valueOptions {
				conds = append(conds, Expr(TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName)+" not like ? ", "%,"+ExprLikeBindValue(v.Value)+",%"))
			}
			return And(conds...), nil
		}
		return nil, nil
	} else {
		return nil, errors.New("GenCondDropdown Comparator is wrong")
	}
}

func (c *TFilterUsecase) GenCondText(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {
	if tListCondition.Comparator == Comparator_is_empty || tListCondition.Comparator == Comparator_is_not_empty {
		return c.HandleComparatorEmpty(fieldEntity, tListCondition, aliasTableName)
	} else if tListCondition.Comparator == Comparator_eq {
		return Eq{fieldEntity.FieldName: tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_neq {
		return Neq{fieldEntity.FieldName: tListCondition.FirstValue()}, nil
	} else if tListCondition.Comparator == Comparator_contains {
		return Like{fieldEntity.FieldName, ExprLikeBindValue(tListCondition.FirstValue())}, nil
	} else if tListCondition.Comparator == Comparator_not_contains {
		return Expr(fieldEntity.FieldName+" not like ? ", "%"+ExprLikeBindValue(tListCondition.FirstValue())+"%"), nil
		//return Expr(fieldEntity.FieldName + " not like '" + "%" + lib.SqlBindValue(tListCondition.FirstValue()) + "%" + "'"), nil
	} else {
		return nil, errors.New("GenCondText Comparator is wrong")
	}
}

func (c *TFilterUsecase) HandleComparatorEmpty(fieldEntity FieldEntity, tListCondition TListCondition, aliasTableName string) (Cond, error) {
	if fieldEntity.FieldType == FieldType_text ||
		fieldEntity.FieldType == FieldType_multitext ||
		fieldEntity.FieldType == FieldType_email ||
		fieldEntity.FieldType == FieldType_text_url ||
		fieldEntity.FieldType == FieldType_tel {
		if tListCondition.Comparator == Comparator_is_empty {
			return Eq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): ""}, nil
		} else {
			return Neq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): ""}, nil
		}
	} else if fieldEntity.FieldType == FieldType_decimal || fieldEntity.FieldType == FieldType_number {
		if tListCondition.Comparator == Comparator_is_empty {
			return IsNull{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName)}, nil
		} else {
			return NotNull{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName)}, nil
		}
	} else if fieldEntity.FieldType == FieldType_dropdown ||
		fieldEntity.FieldType == FieldType_lookup ||
		fieldEntity.FieldType == FieldType_multilookup ||
		fieldEntity.FieldType == FieldType_date {
		if tListCondition.Comparator == Comparator_is_empty {
			return Eq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): ""}, nil
		} else {
			return Neq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): ""}, nil
		}
	} else if fieldEntity.FieldType == FieldType_timestamp {
		if tListCondition.Comparator == Comparator_is_empty {
			return Eq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): 0}, nil
		} else {
			return Neq{TidyTableFieldForSql(fieldEntity.FieldName, aliasTableName): 0}, nil
		}
	} else if fieldEntity.FieldType == FieldType_formula {
		if fieldEntity.FieldName == DataEntry_sys__itf_formula {
			if tListCondition.Comparator == Comparator_is_empty {
				return Eq{TidyTableFieldForSql(FieldName_itf_expiration, aliasTableName): ""}, nil
			} else {
				return Neq{TidyTableFieldForSql(FieldName_itf_expiration, aliasTableName): ""}, nil
			}
		}
	}
	return nil, nil
}
