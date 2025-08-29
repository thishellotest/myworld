package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

type TUsecase struct {
	CommonUsecase      *CommonUsecase
	FieldUsecase       *FieldUsecase
	KindUsecase        *KindUsecase
	FieldOptionUsecase *FieldOptionUsecase
	log                *log.Helper
	TRelaUsecase       *TRelaUsecase
}

func NewTUsecase(CommonUsecase *CommonUsecase, FieldUsecase *FieldUsecase,
	KindUsecase *KindUsecase, FieldOptionUsecase *FieldOptionUsecase, logger log.Logger,
	TRelaUsecase *TRelaUsecase) *TUsecase {
	return &TUsecase{
		CommonUsecase:      CommonUsecase,
		FieldUsecase:       FieldUsecase,
		KindUsecase:        KindUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
		log:                log.NewHelper(logger),
		TRelaUsecase:       TRelaUsecase,
	}
}

func (c *TUsecase) DataById(kind string, id int32) (tData *TData, err error) {
	return c.Data(kind, And(Eq{"id": id}))
}

func (c *TUsecase) DataByGid(kind string, gid string) (tData *TData, err error) {
	return c.Data(kind, And(Eq{"gid": gid}))
}

func (c *TUsecase) ListByCond(kind string, cond Cond) (tList []*TData, err error) {

	//if cond == nil {
	//	return nil, errors.New("cond is nil.")
	//}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	if cond != nil {
		builder.And(cond)
	}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, res, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}
	caches, err := c.GetCaches(fields, res)
	if err != nil {
		return nil, err
	}

	for k, _ := range res {
		tdata := &TData{
			CustomFields: c.GenTFields(&caches, kind, res[k], fields),
			Kind:         kind,
		}
		tList = append(tList, tdata)
	}
	return tList, nil
}

// ListByCondNoRela 此方法不会进行lookup multilookup 关联，specFieldNames为空时所有字段
func (c *TUsecase) ListByCondNoRela(kind string, cond Cond, specFieldNames []string) (tList []*TData, err error) {

	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	if len(specFieldNames) > 0 {
		var destFieldNames []string
		for _, v := range fieldNames {
			if v == DataEntry_gid || lib.InArray(v, specFieldNames) {
				destFieldNames = append(destFieldNames, v)
			}
		}
		fieldNames = destFieldNames
	}

	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	if cond != nil {
		builder.And(cond)
	}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, res, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}
	caches := lib.CacheInit[*TData]()
	for k, _ := range res {
		tdata := &TData{
			CustomFields: c.GenTFields(&caches, kind, res[k], fields),
			Kind:         kind,
		}
		tList = append(tList, tdata)
	}
	return tList, nil
}

func (c *TUsecase) ListByCondCaches_Deleted(caches *lib.Cache[*TData], kind string, cond Cond) (tList []*TData, err error) {

	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}

	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	if cond != nil {
		builder.And(cond)
	}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, res, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}
	for k, _ := range res {
		//lib.DPrintln("ssss___:", len(*caches))
		tdata := &TData{
			CustomFields: c.GenTFields(caches, kind, res[k], fields),
			Kind:         kind,
		}
		tList = append(tList, tdata)
	}
	return tList, nil
}

func (c *TUsecase) DataCaches(caches *lib.Cache[*TData], kind string, cond Cond) (tData *TData, err error) {
	if cond == nil {
		return nil, errors.New("cond is nil.")
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	// 此处使用缓存
	//fields, err := c.FieldUsecase.ListByKind(kind)
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	fields := structField.Records

	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	builder.And(cond)
	//for k, v := range conditions {
	//	eq := Eq{k: v}
	//	builder.And(eq)
	//}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	entity := make(map[string]interface{})
	tx := c.CommonUsecase.DB().Raw(sql).Scan(&entity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if len(entity) == 0 {
		return nil, nil
	}
	tdata := TData{
		CustomFields: c.GenTFields(caches, kind, entity, fields),
		Kind:         kind,
	}
	return &tdata, nil
}

func (c *TUsecase) DataWithOrderBy(kind string, cond Cond, orderBy string) (tData *TData, err error) {

	if cond == nil {
		return nil, errors.New("cond is nil.")
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	// 此处使用缓存
	//fields, err := c.FieldUsecase.ListByKind(kind)
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	fields := structField.Records

	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	builder.And(cond)
	if orderBy != "" {
		builder.OrderBy(orderBy)
	}
	builder.Limit(1)
	//for k, v := range conditions {
	//	eq := Eq{k: v}
	//	builder.And(eq)
	//}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	entity := make(map[string]interface{})
	tx := c.CommonUsecase.DB().Raw(sql).Scan(&entity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if len(entity) == 0 {
		return nil, nil
	}

	caches, err := c.GetCaches(fields, []map[string]interface{}{entity})
	if err != nil {
		return nil, err
	}

	tdata := TData{
		CustomFields: c.GenTFields(&caches, kind, entity, fields),
		Kind:         kind,
	}
	return &tdata, nil
}
func (c *TUsecase) Data(kind string, cond Cond) (tData *TData, err error) {

	//return c.DataWithOrderBy(kind, cond, "")

	if cond == nil {
		return nil, errors.New("cond is nil.")
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	// 此处使用缓存
	//fields, err := c.FieldUsecase.ListByKind(kind)
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	fields := structField.Records

	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}
	builder := Dialect(MYSQL).Select(fieldNames...).From(tableName).Where(Eq{"deleted_at": 0})
	builder.And(cond)
	//for k, v := range conditions {
	//	eq := Eq{k: v}
	//	builder.And(eq)
	//}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return nil, err
	}
	entity := make(map[string]interface{})
	tx := c.CommonUsecase.DB().Raw(sql).Scan(&entity)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if len(entity) == 0 {
		return nil, nil
	}

	caches, err := c.GetCaches(fields, []map[string]interface{}{entity})
	if err != nil {
		return nil, err
	}

	tdata := TData{
		CustomFields: c.GenTFields(&caches, kind, entity, fields),
		Kind:         kind,
	}
	return &tdata, nil
}

func (c *TUsecase) GetCaches(fields TypeFieldList, records []map[string]interface{}) (lib.Cache[*TData], error) {
	var caches lib.Cache[*TData]
	var err error
	if configs.Enable_NewVersionForT {
		caches, err = c.GetRelaCaches(fields, records)
		//lib.DPrintln("caches:", caches)
		if err != nil {
			c.log.Warn(err)
			return nil, err
		}
	} else {
		caches = lib.CacheInit[*TData]()
	}
	return caches, nil
}

func (c *TUsecase) GetRelaCaches(fields TypeFieldList, records []map[string]interface{}) (caches lib.Cache[*TData], err error) {
	relaMap := c.TRelaUsecase.GetRelaMap(fields, records)
	return c.GetRelaCachesByRelaMap(relaMap)

	//caches = lib.CacheInit[*TData]()
	//for k, v := range relaMap {
	//	gids := v.GetGids()
	//	if len(gids) > 0 {
	//		res, err := c.ListByCondNoRela(k, In(DataEntry_gid, gids), []string{v.RelaName})
	//		if err != nil {
	//			return nil, err
	//		}
	//		for k1, v1 := range res {
	//			key := fmt.Sprintf("%s:%s", k, v1.Gid())
	//			caches.Set(key, res[k1])
	//		}
	//	}
	//}
	//return caches, nil
}

func (c *TUsecase) GetRelaCachesByRelaMap(relaMap TRelaMap) (caches lib.Cache[*TData], err error) {
	caches = lib.CacheInit[*TData]()
	for k, v := range relaMap {
		gids := v.GetGids()
		if len(gids) > 0 {
			res, err := c.ListByCondNoRela(k, In(DataEntry_gid, gids), []string{v.RelaName})
			if err != nil {
				return nil, err
			}
			for k1, v1 := range res {
				key := fmt.Sprintf("%s:%s", k, v1.Gid())
				caches.Set(key, res[k1])
			}
		}
	}
	return caches, nil
}

func (c *TUsecase) DataByRawSql(kindEntity KindEntity, sql string) (tData *TData, err error) {

	fields, err := c.FieldUsecase.ListByKind(kindEntity.Kind)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, row, err := lib.SqlRowsToRow(sqlRows)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}

	caches, err := c.GetCaches(fields, []map[string]interface{}{row})
	if err != nil {
		return nil, err
	}

	tData = &TData{
		CustomFields: c.GenTFields(&caches, kindEntity.Kind, row, fields),
		Kind:         kindEntity.Kind,
	}
	return tData, err
}

func (c *TUsecase) Total(kind string) (total int64, err error) {

	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return 0, err
	}

	sql, err := Dialect(MYSQL).Select("count(*) as c").From(tableName).ToBoundSQL()
	if err != nil {
		return 0, err
	}
	return c.CommonUsecase.Count(c.CommonUsecase.DB(), sql)
}

func (c *TUsecase) ListByRawSql(kind string, sql string) (tDataList TDataList, err error) {

	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}

	if err != nil {
		return nil, err
	}
	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, list, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}

	caches, err := c.GetCaches(fields, list)
	if err != nil {
		return nil, err
	}

	for k, _ := range list {
		tdata := TData{
			CustomFields: c.GenTFields(&caches, kind, list[k], fields),
			Kind:         kind,
		}
		tDataList = append(tDataList, tdata)
	}
	return tDataList, err
}

func (c *TUsecase) List(kind string, tUser *TData, request *TListRequest, page int, pageSize int) (tDataList TDataList, err error) {
	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}
	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}

	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}

	sql, err := Dialect(MYSQL).Select(fieldNames...).From(tableName).OrderBy("id desc").
		Limit(page, HandleOffset(page, pageSize)).ToBoundSQL()
	if err != nil {
		return nil, err
	}
	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, list, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}

	caches, err := c.GetCaches(fields, list)
	if err != nil {
		return nil, err
	}

	for k, _ := range list {
		tdata := TData{
			CustomFields: c.GenTFields(&caches, kind, list[k], fields),
			Kind:         kind,
		}
		tDataList = append(tDataList, tdata)
	}
	return tDataList, err
}

func (c *TUsecase) TotalByCond(kindEntity KindEntity, cond Cond) (total int64, err error) {

	tableName, err := c.KindUsecase.CacheTableNameByKind(kindEntity.Kind)
	if err != nil {
		return 0, err
	}

	query := Dialect(MYSQL).Select("count(*) as c").From(tableName)
	query.And(Eq{"deleted_at": 0})
	if cond != nil {
		query.And(cond)
	}
	sql, err := query.ToBoundSQL()
	if err != nil {
		return 0, err
	}
	return c.CommonUsecase.Count(c.CommonUsecase.DB(), sql)
}

// ListByCondWithPaging orderBY: id desc
func (c *TUsecase) ListByCondWithPaging(kindEntity KindEntity, cond Cond, orderBy string, page int, pageSize int) (tDataList TDataList, err error) {
	fields, err := c.FieldUsecase.ListByKind(kindEntity.Kind)
	if err != nil {
		return nil, err
	}
	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}

	tableName, err := c.KindUsecase.CacheTableNameByKind(kindEntity.Kind)
	if err != nil {
		return nil, err
	}

	query := Dialect(MYSQL).Select(fieldNames...).From(tableName)
	query.And(Eq{"deleted_at": 0})
	if orderBy != "" {
		query.OrderBy(orderBy)
	}
	if cond != nil {
		query.And(cond)
	}

	sql, err := query.Limit(pageSize, HandleOffset(page, pageSize)).ToBoundSQL()
	if err != nil {
		return nil, err
	}
	sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, list, err := lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return nil, err
	}

	caches, err := c.GetCaches(fields, list)
	if err != nil {
		return nil, err
	}

	for k, _ := range list {
		tdata := TData{
			CustomFields: c.GenTFields(&caches, kindEntity.Kind, list[k], fields),
			Kind:         kindEntity.Kind,
		}
		tDataList = append(tDataList, tdata)
	}
	return tDataList, err
}

// InterfaceToString nil 会转为 "" 很重要
func InterfaceToString(val interface{}) string {
	return lib.InterfaceToString(val)
}

func (c *TUsecase) DataByGidWithCaches(caches *lib.Cache[*TData], kind string, gid string) *TData {
	cacheKey := fmt.Sprintf("%s:%s", kind, gid)
	entity, exist := caches.Get(cacheKey)

	if configs.Enable_NewVersionForT {
		return entity
	}

	//fmt.Printf("Map variable address: %p\n", &*caches)
	if !exist {
		relaData, err := c.DataCaches(caches, kind, Eq{DataEntry_gid: gid})
		if err != nil {
			c.log.Error(err)
		}
		caches.Set(cacheKey, relaData)
		entity = relaData
	}
	return entity
}

func (c *TUsecase) GenTFieldMultiValues(caches *lib.Cache[*TData], fieldEntity FieldEntity, val string) (multiValues TFieldMultiValues) {

	if val != "" {
		vals := strings.Split(val, ",")
		for _, gid := range vals {
			if gid != "" {
				multiValue := TFieldMultiValue{
					Value: gid,
				}
				entity := c.DataByGidWithCaches(caches, fieldEntity.RelaKind, gid)
				if entity != nil {
					multiValue.Label = entity.CustomFields.DisplayValueByNameBasic(fieldEntity.RelaName)
				}
				multiValues = append(multiValues, multiValue)
			}
		}
	}
	return
}

func (c *TUsecase) GenTField(caches *lib.Cache[*TData], kind string, fieldEntity FieldEntity, val string) TField {

	tField := TField{
		Name:         fieldEntity.FieldName,
		Type:         fieldEntity.FieldType,
		DisplayValue: &val,
		Value:        val,
	}
	if fieldEntity.FieldType == FieldType_decimal {
		if val != "" {
			num, _ := strconv.ParseFloat(val, 64)
			p := message.NewPrinter(language.English)
			tV := p.Sprintf("%.2f", num)
			// todo:lgl 将来可以指定小数的长度
			tField.DisplayValue = &tV
		}
		stringVal := val
		tField.TextValue = &stringVal
	} else if fieldEntity.FieldType == FieldType_timestamp {
		str := ""
		if len(val) != 0 && val != "0" {
			t, _ := strconv.ParseInt(val, 0, 64)
			str = time.Unix(t, 0).Format(time.RFC3339)
			numberVal := lib.InterfaceToInt32(val)
			tField.NumberValue = &numberVal
		} else {
			tField.Value = ""
		}
		tField.DisplayValue = &str
	} else if IsNumberFieldType(fieldEntity.FieldType) {
		if val != "" {
			numberVal := lib.InterfaceToInt32(val)
			tField.NumberValue = &numberVal
		} else {
			tField.Value = ""
		}
	} else if fieldEntity.FieldType == FieldType_dropdown {
		stringVal := val
		tField.DisplayValue = to.Ptr("")
		if stringVal != "" {
			optionStruct, err := c.FieldOptionUsecase.CacheStructByKind(kind)
			if err != nil {
				c.log.Error(err)
			} else {
				if optionStruct != nil {
					temp := optionStruct.AllByFieldName(fieldEntity).GetByValue(stringVal)
					if temp != nil {
						tField.DisplayValue = to.Ptr(temp.OptionLabel)
					}
				}
			}
		}
		tField.TextValue = &stringVal
	} else if fieldEntity.FieldType == FieldType_multidropdown {

		tField.DisplayValue = to.Ptr("")
		//tField.Value = ""
		stringVal := val
		tField.TextValue = &stringVal

		if val != "" {
			var multiValues TFieldMultiValues
			vals := strings.Split(val, ",")

			optionStruct, err := c.FieldOptionUsecase.CacheStructByKind(kind)
			if err != nil {
				c.log.Error(err)
			}
			for _, optionOneValue := range vals {
				if optionOneValue != "" {
					multiValue := TFieldMultiValue{
						Value: optionOneValue,
					}
					if optionStruct != nil {
						temp := optionStruct.AllByFieldName(fieldEntity).GetByValue(optionOneValue)
						if temp != nil {
							multiValue.Label = temp.OptionLabel
						}
					}

					multiValues = append(multiValues, multiValue)
				}
			}
			tField.MultiValues = multiValues
		}

	} else if fieldEntity.FieldType == FieldType_lookup {
		stringVal := val
		tField.DisplayValue = to.Ptr("")
		if stringVal != "" {
			entity := c.DataByGidWithCaches(caches, fieldEntity.RelaKind, stringVal)
			if entity != nil {
				tField.DisplayValue = to.Ptr(entity.CustomFields.DisplayValueByNameBasic(fieldEntity.RelaName))
			}
		}
		tField.TextValue = &stringVal
	} else if fieldEntity.FieldType == FieldType_multilookup {
		tField.DisplayValue = to.Ptr("")
		//tField.Value = ""
		stringVal := val
		tField.TextValue = &stringVal
		if val != "" {
			tField.MultiValues = c.GenTFieldMultiValues(caches, fieldEntity, val)
		}
	} else if fieldEntity.FieldType == FieldType_date {
		stringVal := val
		tField.TextValue = &stringVal
		if fieldEntity.FieldName == DataEntry_sys__due_date {
			tFieldExtendForSysDueDate, err := GenTFieldExtendForSysDueDate(val)
			if err != nil {
				c.log.Error("GenTFieldExtendForSysDueDate:", val, " : ", err)
			}
			tField.Extend = tFieldExtendForSysDueDate
		}
	} else {
		stringVal := val
		tField.TextValue = &stringVal
	}
	return tField
}

func (c *TUsecase) DoFormula(kindEntity KindEntity, tData *TData) {
	for k, v := range tData.CustomFields {
		//if v.Type == FieldType_formula {
		if kindEntity.Kind == Kind_client_cases && v.Name == DataEntry_sys__itf_formula {
			itfDate := tData.CustomFields.TextValueByNameBasic(FieldName_itf_expiration)
			val := itfDate
			if val != "" {
				currentTime := time.Now().In(configs.GetVBCDefaultLocation())
				currentTime, _ = time.ParseInLocation(time.DateOnly, currentTime.Format(time.DateOnly), configs.GetVBCDefaultLocation())
				itfTime, _ := time.ParseInLocation(time.DateOnly, val, configs.GetVBCDefaultLocation())
				daysRemaining := int64(math.Ceil(itfTime.Sub(currentTime).Hours() / 24))
				daysRemainingStr := strconv.FormatInt(daysRemaining, 10)
				tData.CustomFields[k].DisplayValue = &daysRemainingStr
				tData.CustomFields[k].TextValue = &daysRemainingStr
				tData.CustomFields[k].Value = daysRemainingStr
				entity, _ := GenTFieldExtendForSysItfFormula(daysRemainingStr)
				tData.CustomFields[k].Extend = entity
			}
		} else if kindEntity.Kind == Kind_client_cases && v.Name == FieldName_itf_expiration {

			tData.CustomFields[k].Extend = GenItfExpirationExtend(tData)

			//itfDate := tData.CustomFields.TextValueByNameBasic(FieldName_itf_expiration)
			//val := itfDate
			//if val != "" {
			//	currentTime := time.Now().In(lib.GetVBCDefaultLocation())
			//	itfTime, _ := time.ParseInLocation(time.DateOnly, val, lib.GetVBCDefaultLocation())
			//	daysRemaining := int64(math.Ceil(itfTime.Sub(currentTime).Hours() / 24))
			//	daysRemainingStr := strconv.FormatInt(daysRemaining, 10)
			//
			//	entity, _ := GenTFieldExtendForSysItfFormula(daysRemainingStr)
			//	entity.Value = itfDate
			//	entity.Label = tData.CustomFields.DisplayValueByNameBasic(FieldName_itf_expiration)
			//	tData.CustomFields[k].Extend = entity
			//}

		}
		//}
	}
}

func (c *TUsecase) GenTFields(caches *lib.Cache[*TData], kind string, dbRow map[string]interface{}, fields TypeFieldList) (tFields TFields) {

	for k, v := range fields {
		val := lib.InterfaceToString(dbRow[v.FieldName])

		// 新版本开始 如果遇到风险，请回退到旧版本
		tField := c.GenTField(caches, kind, *fields[k], val)
		// 新版本结束

		// 旧版本开始
		/*
			tField := TField{
				Name:         v.FieldName,
				Type:         v.FieldType,
				DisplayValue: &val,
				Value:        val,
			}
			if v.FieldType == FieldType_timestamp {
				str := ""
				if len(val) != 0 && val != "0" {
					t, _ := strconv.ParseInt(val, 0, 64)
					str = time.Unix(t, 0).Format(time.RFC3339)
					numberVal := lib.InterfaceToInt32(dbRow[v.FieldName])
					tField.NumberValue = &numberVal
				} else {
					tField.Value = ""
				}
				tField.DisplayValue = &str
			} else if IsNumberFieldType(v.FieldType) {
				if dbRow[v.FieldName] != nil {
					numberVal := lib.InterfaceToInt32(dbRow[v.FieldName])
					tField.NumberValue = &numberVal
				} else {
					tField.Value = ""
				}
			} else if v.FieldType == FieldType_dropdown {
				stringVal := InterfaceToString(dbRow[v.FieldName])
				tField.DisplayValue = to.Ptr("")
				if stringVal != "" {
					optionStruct, err := c.FieldOptionUsecase.CacheStructByKind(kind)
					if err != nil {
						c.log.Error(err)
					} else {
						if optionStruct != nil {
							temp := optionStruct.AllByFieldName(v.FieldName).GetByValue(stringVal)
							if temp != nil {
								tField.DisplayValue = to.Ptr(temp.OptionLabel)
							}
						}
					}
				}
				tField.TextValue = &stringVal
			} else if v.FieldType == FieldType_lookup {
				stringVal := InterfaceToString(dbRow[v.FieldName])
				tField.DisplayValue = to.Ptr("")
				if stringVal != "" {
					cacheKey := fmt.Sprintf("%s:%s", kind, stringVal)
					entity, exist := caches.Get(cacheKey)
					if !exist {
						relaData, err := c.DataByGid(v.RelaKind, stringVal)
						if err != nil {
							c.log.Error(err)
						}
						caches.Set(cacheKey, relaData)
						entity = relaData
					}
					if entity != nil {
						tField.DisplayValue = to.Ptr(entity.CustomFields.DisplayValueByNameBasic(v.RelaName))
					}
				}
				tField.TextValue = &stringVal
			} else {
				stringVal := InterfaceToString(dbRow[v.FieldName])
				tField.TextValue = &stringVal
			}
		*/
		// 旧版本结束

		tFields = append(tFields, tField)
	}
	return
}
