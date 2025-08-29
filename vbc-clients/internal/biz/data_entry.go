package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/lib"
	//. "vbc/lib/builder"
	"vbc/lib/uuid"
)

const (
	TableName_client       = "clients"
	TableName_client_cases = "client_cases"
	TableName_client_tasks = "client_tasks"
)

type DataEntryUsecase struct {
	log                        *log.Helper
	CommonUsecase              *CommonUsecase
	FieldUsecase               *FieldUsecase
	KindUsecase                *KindUsecase
	UniqueCodeGeneratorUsecase *UniqueCodeGeneratorUsecase
	EventBus                   *EventBus
}

func NewDataEntryUsecase(logger log.Logger, CommonUsecase *CommonUsecase,
	FieldUsecase *FieldUsecase, KindUsecase *KindUsecase,
	UniqueCodeGeneratorUsecase *UniqueCodeGeneratorUsecase,
	EventBus *EventBus) *DataEntryUsecase {

	uc := &DataEntryUsecase{
		log:                        log.NewHelper(logger),
		CommonUsecase:              CommonUsecase,
		FieldUsecase:               FieldUsecase,
		KindUsecase:                KindUsecase,
		UniqueCodeGeneratorUsecase: UniqueCodeGeneratorUsecase,
		EventBus:                   EventBus,
	}
	return uc
}

// TypeDataEntry 数据库真实值
// dropdown: option_value
// lookup: gid
// multilookup: ,gid,gid,gid,
// timestamp: 时间戳（秒）
// date:  2022-12-12 格式
type TypeDataEntry lib.TypeMap

func (c TypeDataEntry) AllFieldNames() (res TypeFieldNameMaps) {
	res = make(TypeFieldNameMaps)
	for k, _ := range c {
		res[k] = true
	}
	return res
}

type TypeDataEntryList []TypeDataEntry

func (c TypeDataEntryList) Get(key string, value string) TypeDataEntry {
	for k, v := range c {
		if _, ok := v[key]; ok {
			if InterfaceToString(v[key]) == value {
				return c[k]
			}
		}
	}
	return nil
}

func (c TypeDataEntryList) AllFieldNames() (res TypeFieldNameMaps) {
	res = make(TypeFieldNameMaps)
	for k, _ := range c {
		for k1, _ := range c[k] {
			res[k1] = true
		}
	}
	return res
}

func (c TypeDataEntryList) RecognizeValues(recognizeFieldName string) (r []string) {
	res := make(map[string]bool)
	for k, _ := range c {
		for k1, _ := range c[k] {
			if k1 == recognizeFieldName {
				val := lib.InterfaceToString(c[k][k1])
				res[val] = true
			}
		}
	}
	for k, _ := range res {
		r = append(r, k)
	}
	return r
}

// RetainTypeDataEntryList 只保留可用的数据
func RetainTypeDataEntryList(dataList TypeDataEntryList, enableFieldNames TypeFieldNameMaps) (res TypeDataEntryList) {
	for k, _ := range dataList {
		for k1, _ := range dataList[k] {
			if _, ok := enableFieldNames[k1]; !ok {
				delete(dataList[k], k1)
			}
		}
		if len(dataList[k]) > 0 {
			res = append(res, dataList[k])
		}
	}
	return res
}

// RetainTypeDataEntry 只保留可用的数据
func RetainTypeDataEntry(row TypeDataEntry, enableFieldNames TypeFieldNameMaps) (res TypeDataEntry) {
	for k1, _ := range row {
		if _, ok := enableFieldNames[k1]; !ok {
			delete(row, k1)
		}
	}
	return row
}

func (c *DataEntryUsecase) HandleOne(kind string, data TypeDataEntry, recognizeFieldName string, operUser *TData) (dataEntryOperResult DataEntryOperResult, err error) {
	var dataList TypeDataEntryList
	dataList = append(dataList, data)
	return c.Handle(kind, dataList, recognizeFieldName, operUser)

}

func (c *DataEntryUsecase) Delete(kind KindEntity, ids []interface{}, recognizeFieldName string) error {
	if len(ids) <= 0 {
		return nil
	}
	tableName := kind.Tablename
	//sql, args, err := Update(Eq{"deleted_at": time.Now().Unix()}).From(tableName).Where(In(recognizeFieldName, ids)).ToSQL()
	//if err != nil {
	//	return err
	//}
	err := c.CommonUsecase.DB().Table(tableName).Where(recognizeFieldName+" IN ?", ids).Updates(map[string]interface{}{
		"deleted_at": time.Now().Unix(),
	}).Error
	return err
}

func DataEntryOperResultCombine(result DataEntryOperResult, newResult DataEntryOperResult) DataEntryOperResult {
	if result == nil {
		result = make(DataEntryOperResult)
	}
	for k, v := range newResult {
		if _, ok := result[k]; ok {
			oriItem := result[k]

			if !oriItem.IsNewRecord {
				oriItem.IsNewRecord = v.IsNewRecord
			}
			if !oriItem.IsUpdated {
				oriItem.IsUpdated = v.IsUpdated
			}
		} else {
			result[k] = v
		}
	}
	return result
}

// Handle 处理数据入库（可以处理插入和更新）：说明进入此方法的数据确保  recognizeFieldName对应在dataList必段唯一。否则产生多条数据
func (c *DataEntryUsecase) Handle(kind string, dataList TypeDataEntryList, recognizeFieldName string, operUser *TData) (dataEntryOperResult DataEntryOperResult, err error) {

	dataEntryOperResult = make(DataEntryOperResult)
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return dataEntryOperResult, err
	}
	kindList, err := c.KindUsecase.CacheList()
	if err != nil {
		return dataEntryOperResult, err
	}
	tableName, err := kindList.TableNameByKind(kind)
	if err != nil {
		return dataEntryOperResult, err
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return dataEntryOperResult, err
	}
	if kindEntity == nil {
		return dataEntryOperResult, errors.New("kindEntity does not exist.")
	}

	handleUpdateResult, newData, err := c.HandleUpdate(*kindEntity, dataList, recognizeFieldName, structField, tableName, operUser)
	if err != nil {
		c.log.Error(err)
	}
	dataEntryOperResult = DataEntryOperResultCombine(dataEntryOperResult, handleUpdateResult)

	if len(newData) > 0 {
		if kindEntity.Kind == Kind_clients || kindEntity.Kind == Kind_client_cases {
			// 新版开始 此版本为解决client cases 初始化时，没有change_history无法触发自动
			var newDataRecognizeFieldOnly TypeDataEntryList
			for k, _ := range newData {
				typeDataEntry := make(TypeDataEntry)
				typeDataEntry[recognizeFieldName] = newData[k][recognizeFieldName]
				newDataRecognizeFieldOnly = append(newDataRecognizeFieldOnly, typeDataEntry)
			}
			insertDataResult, err := c.InsertData(*kindEntity, tableName, newDataRecognizeFieldOnly, structField, recognizeFieldName, operUser)
			if err != nil {
				c.log.Error(err)
			}
			dataEntryOperResult = DataEntryOperResultCombine(dataEntryOperResult, insertDataResult)

			handleResult, _, err := c.HandleUpdate(*kindEntity, newData, recognizeFieldName, structField, tableName, operUser)
			if err != nil {
				c.log.Error(err)
			}
			dataEntryOperResult = DataEntryOperResultCombine(dataEntryOperResult, handleResult)
			// 新版本结束
		} else {
			//if kindEntity.Kind == Kind_notes || kindEntity.Kind == Kind_timelines {
			// 以前版本开始 此版本需要一次入库，才能产生正确的timelines
			insertDataResult, err := c.InsertData(*kindEntity, tableName, newData, structField, recognizeFieldName, operUser)
			if err != nil {
				c.log.Error(err)
			}
			dataEntryOperResult = DataEntryOperResultCombine(dataEntryOperResult, insertDataResult)
			// 以前版本结束
			//}
		}
	}
	return dataEntryOperResult, nil
}

func (c *DataEntryUsecase) HandleDependFieldSelect(structField *TypeFieldStruct, destSelectColumns []string) ([]string, error) {
	if structField == nil {
		return nil, errors.New("HandleDependFieldSelect: structField is nil")
	}
	res := structField.GetFieldDependList().GetFieldNamesIsDependent(destSelectColumns)
	for _, v := range res {
		for _, v1 := range v.FieldDependItemList {
			destSelectColumns = append(destSelectColumns, v1.FieldName)
		}
	}
	return destSelectColumns, nil
}

// UpdateOne 只能处理更新 已经兼容after方法
func (c *DataEntryUsecase) UpdateOne(kind string, row TypeDataEntry, recognizeFieldName string, operUser *TData) (dataEntryOperItem DataEntryOperItem, err error) {

	fieldNames := row.AllFieldNames()
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return dataEntryOperItem, err
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return dataEntryOperItem, err
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return dataEntryOperItem, err
	}
	if kindEntity == nil {
		return dataEntryOperItem, errors.New("kindEntity does not exist.")
	}

	fieldNames = structField.FilterFieldName(fieldNames)
	if len(fieldNames) == 0 {
		return dataEntryOperItem, errors.New("fieldNames都不可用")
	}
	// 只保留可用数据
	row = RetainTypeDataEntry(row, fieldNames)
	if len(row) == 0 {
		return dataEntryOperItem, errors.New("row数据都不可用")
	}

	if _, ok := row[recognizeFieldName]; !ok {
		return dataEntryOperItem, errors.New("recognizeFieldName字段在原始数据无值:" + recognizeFieldName)
	}

	selectColumns := fieldNames.ToSqlSelect()

	var dbRowsList DbRowsList

	destSelectColumns := selectColumns
	destSelectColumns = append(destSelectColumns, DataEntry_Incr_id_name) // 加入id更新使用
	destSelectColumns = append(destSelectColumns, DataEntry_gid)          // 加入gid业务使用
	destSelectColumns = append(destSelectColumns, DataEntry_modified_by)

	if configs.EnabledDataEntryDependField {
		destSelectColumns, err = c.HandleDependFieldSelect(structField, destSelectColumns)
		if err != nil {
			return dataEntryOperItem, err
		}
	}
	destSelectColumns = lib.RemoveDuplicates(destSelectColumns)

	sqlRows, err := c.CommonUsecase.DB().Table(tableName).Select(destSelectColumns).
		Where(recognizeFieldName+" in ? and deleted_at=0", []interface{}{row[recognizeFieldName]}).
		Rows()
	if err != nil {
		return dataEntryOperItem, err
	}
	if sqlRows != nil {
		defer sqlRows.Close()
	}
	_, dbRowsList, err = lib.SqlRowsTrans(sqlRows)
	if err != nil {
		return dataEntryOperItem, err
	}

	// 开始入库操作
	//var newData TypeDataEntryList
	recoVal := lib.InterfaceToString(row[recognizeFieldName])
	oldRow := dbRowsList.GetByFieldNameValue(recoVal, recognizeFieldName)
	if oldRow == nil { // 新增数据
		return dataEntryOperItem, errors.New("row does not exist")
	} else { // 修改数据
		isUpdated, err, dataEntryModifyDataMap := c.UpdateModifyData(*kindEntity, tableName, oldRow, row, structField, recognizeFieldName, operUser)
		if err != nil {
			c.log.Error(err)
			return dataEntryOperItem, err
		}
		dataEntryOperItem.IsUpdated = isUpdated
		dataEntryOperItem.DataEntryModifyDataMap = dataEntryModifyDataMap
	}

	if recognizeFieldName == DataEntry_gid {
		dataEntryOperResult := make(DataEntryOperResult)
		dataEntryOperResult[recoVal] = dataEntryOperItem
		var dataList TypeDataEntryList
		dataList = append(dataList, row)
		err = c.AfterHandleUpdate(*kindEntity, structField, recognizeFieldName, dataEntryOperResult, dataList, operUser)
	}
	if err != nil {
		c.log.Error(err)
	}

	return dataEntryOperItem, nil
}

// DataEntryOperResult string是指RecognizeValue 记录即可以是插入且更新，
type DataEntryOperResult map[string]DataEntryOperItem

func (c DataEntryOperResult) GetOne() (key string, item DataEntryOperItem) {
	for k, _ := range c {
		key = k
		item = c[k]
		return
	}
	return
}

func (c DataEntryOperResult) GetByValue(recognizeValue string) (item DataEntryOperItem) {
	if _, ok := c[recognizeValue]; ok {
		item = c[recognizeValue]
	}
	return item
}

type DataEntryOperItem struct {
	//RecognizeValue string // 此值与GetByValue，会有一些歧义，所以不使用
	IsUpdated              bool
	IsNewRecord            bool
	DataEntryModifyDataMap DataEntryModifyDataMap
}

func (c *DataEntryUsecase) AfterHandleUpdate(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, operUser *TData) error {

	//if structField == nil {
	//	c.log.Error("structField is nil")
	//}
	//if kindEntity.Kind == Kind_client_tasks {
	//	if recognizeFieldName == DataEntry_gid && len(dataEntryOperResult) > 0 {
	//		c.log.Info("AfterHandleUpdate:1 dataEntryOperResult:", InterfaceToString(dataEntryOperResult))
	//		var whatGids []string
	//		var whoGids []string
	//		for gid, v := range dataEntryOperResult {
	//			if v.IsUpdated {
	//				for k1, v1 := range v.DataEntryModifyDataMap {
	//					if k1 == TaskFieldName_what_id_gid {
	//						newVal := v1.GetNewVal(FieldType_text)
	//						oldNew := v1.GetOldVal(FieldType_text)
	//						if newVal != "" {
	//							whatGids = append(whatGids, newVal)
	//						}
	//						if oldNew != "" {
	//							whatGids = append(whatGids, oldNew)
	//						}
	//					} else if k1 == TaskFieldName_who_id_gid {
	//						newVal := v1.GetNewVal(FieldType_text)
	//						oldNew := v1.GetOldVal(FieldType_text)
	//						if newVal != "" {
	//							whoGids = append(whoGids, newVal)
	//						}
	//						if oldNew != "" {
	//							whoGids = append(whoGids, oldNew)
	//						}
	//					} else if k1 == TaskFieldName_due_date {
	//						err := c.QueueUsecase.TaskDueDateChange([]string{gid})
	//						if err != nil {
	//							c.log.Error(err)
	//						}
	//					}
	//				}
	//			}
	//		}
	//		c.log.Info("AfterHandleUpdate:2 whatGids", whatGids)
	//		c.log.Info("AfterHandleUpdate:2 whoGids", whoGids)
	//		err := c.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), whatGids)
	//		if err != nil {
	//			c.log.Error(err, "AfterHandleUpdate whatGids:", whatGids)
	//		}
	//		err = c.QueueUsecase.PushClientTaskHandleWhoGidJobTasks(context.TODO(), whoGids)
	//		if err != nil {
	//			c.log.Error(err, "AfterHandleUpdate whoGids:", whoGids)
	//		}
	//	}
	//}

	modifiedBy := ""
	if operUser != nil {
		modifiedBy = operUser.Gid()
	}
	c.EventBus.Publish(EventBus_AfterHandleUpdate, kindEntity, structField, recognizeFieldName, dataEntryOperResult, sourceData, modifiedBy)
	return nil
}

// HandleUpdate dataList 不能包含系统保留的字段：create_at updated_at crated_by modified_by deleted_at
func (c *DataEntryUsecase) HandleUpdate(kindEntity KindEntity, dataList TypeDataEntryList,
	recognizeFieldName string, structField *TypeFieldStruct,
	tableName string,
	operUser *TData) (dataEntryOperResult DataEntryOperResult, newData TypeDataEntryList, err error) {

	dataEntryOperResult = make(DataEntryOperResult)

	fieldNames := dataList.AllFieldNames()

	fieldNames = structField.FilterFieldName(fieldNames)
	if len(fieldNames) == 0 {
		c.log.Info("fieldNames都不可用")
		return dataEntryOperResult, nil, nil
	}

	// 只保留可用数据
	dataList = RetainTypeDataEntryList(dataList, fieldNames)
	if len(dataList) == 0 {
		return dataEntryOperResult, nil, errors.New("dataList数据都不可用")
	}

	// 验证数据合法性
	if recognizeFieldName != DataEntry_Incr_id_name {
		for _, v := range dataList {
			if _, ok := v[recognizeFieldName]; !ok {
				return dataEntryOperResult, nil, errors.New("recognizeFieldName字段在原始数据无值:" + recognizeFieldName)
			}
		}
	}

	selectColumns := fieldNames.ToSqlSelect()
	values := dataList.RecognizeValues(recognizeFieldName)

	var dbRowsList DbRowsList
	if len(values) > 0 {
		destSelectColumns := selectColumns
		destSelectColumns = append(destSelectColumns, DataEntry_Incr_id_name) // 加入id更新使用

		if !IsCustomKind(kindEntity.Kind) {
			destSelectColumns = append(destSelectColumns, DataEntry_gid) // 加入gid业务使用
		}
		destSelectColumns = append(destSelectColumns, DataEntry_modified_by)
		if configs.EnabledDataEntryDependField {
			destSelectColumns, err = c.HandleDependFieldSelect(structField, destSelectColumns)
			if err != nil {
				return dataEntryOperResult, nil, err
			}
		}
		destSelectColumns = lib.RemoveDuplicates(destSelectColumns)

		sqlRows, err := c.CommonUsecase.DB().Table(tableName).Select(destSelectColumns).
			Where(recognizeFieldName+" in ? and deleted_at=0", values).
			Rows()
		if err != nil {
			return dataEntryOperResult, nil, err
		}
		if sqlRows != nil {
			defer sqlRows.Close()
		}
		_, dbRowsList, _ = lib.SqlRowsTrans(sqlRows)
	}

	// 开始入库操作
	//var newData TypeDataEntryList
	for k, v := range dataList {
		recoVal := lib.InterfaceToString(v[recognizeFieldName])
		oldRow := dbRowsList.GetByFieldNameValue(recoVal, recognizeFieldName)
		if oldRow == nil { // 新增数据
			newData = append(newData, dataList[k])
		} else { // 修改数据
			isUpdated, err, dataEntryModifyDataMap := c.UpdateModifyData(kindEntity, tableName, oldRow, dataList[k], structField, recognizeFieldName, operUser)
			if err != nil {
				c.log.Error(err)
			}
			if isUpdated {
				row := lib.TypeMap(v)
				recognizeVal := row.GetString(recognizeFieldName)
				if recognizeVal != "" {
					if _, ok := dataEntryOperResult[recognizeVal]; ok {
						item := dataEntryOperResult[recognizeVal]
						item.IsUpdated = true
						item.DataEntryModifyDataMap = dataEntryModifyDataMap
						dataEntryOperResult[recognizeVal] = item
					} else {
						dataEntryOperResult[recognizeVal] = DataEntryOperItem{
							IsUpdated:              true,
							DataEntryModifyDataMap: dataEntryModifyDataMap,
						}
					}
				}
			}
		}
	}

	err = c.AfterHandleUpdate(kindEntity, structField, recognizeFieldName, dataEntryOperResult, dataList, operUser)
	if err != nil {
		c.log.Error(err)
	}
	return dataEntryOperResult, newData, nil
}

// InsertOne 此方法过期，不要使用， 不能兼容after方法
func (c *DataEntryUsecase) InsertOne(kind string, newRow TypeDataEntry, operUser *TData) (gid string, err error) {
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	uniqGid := uuid.UuidWithoutStrike()
	newRow["gid"] = uniqGid
	err = c.InsertOneData(tableName, newRow, structField, operUser)
	if err != nil {
		return "", err
	}
	return uniqGid, nil
}

// InsertOneData 插入一条新数据 不要使用， 不能兼容after方法
func (c *DataEntryUsecase) InsertOneData(tableName string, newRow TypeDataEntry, structField *TypeFieldStruct, operUser *TData) error {

	if newRow == nil {
		return errors.New("newRow is nil")
	}
	allFieldNames := newRow.AllFieldNames()
	allFieldNames.DeleteFieldName(DataEntry_Incr_id_name)
	allFieldNames.DeleteFieldName("deleted_at")

	if len(allFieldNames) == 0 {
		return errors.New("无字段数据")
	}
	listFields := allFieldNames.ToSqlSelect()
	fieldStr := "`" + strings.Join(listFields, "`,`") + "`" + ",created_at,updated_at"
	if tableName == TableName_client_cases {
		fieldStr += "," + FieldName_uniqcode
	}
	if operUser != nil {
		fieldStr += fmt.Sprintf(",%s,%s", DataEntry_created_by, DataEntry_modified_by)
	}
	sql := fmt.Sprintf("insert into %s (%s) values ", tableName, fieldStr)
	var values []string
	v := newRow
	s := "("
	var sVals []string
	for _, fileName := range listFields {
		field := structField.GetByFieldName(fileName)
		if field == nil {
			return errors.New("字段：" + fileName + "不存在")
		}
		if _, ok := v[fileName]; ok {
			afterFormatVal := field.TransToCorrectValueFormat(v[fileName])
			sVals = append(sVals, fmt.Sprintf("%s", lib.SqlBindValue(afterFormatVal)))
		} else {
			sVals = append(sVals, fmt.Sprintf("%s", lib.SqlBindValue(field.DefaultValue())))
		}
	}
	s += strings.Join(sVals, ",")
	s += fmt.Sprintf(",%d,%d", time.Now().Unix(), time.Now().Unix())
	if tableName == TableName_client_cases {
		uniqCode, _ := c.UniqueCodeGeneratorUsecase.GenUuid(UniqueCodeGenerator_Type_ClientUniqCode, 0)
		s += "," + lib.SqlBindValue(uniqCode)
	}
	if operUser != nil {
		s += "," + lib.SqlBindValue(operUser.Gid())
		s += "," + lib.SqlBindValue(operUser.Gid())
	}

	s += ")"
	values = append(values, s)

	sql += strings.Join(values, ",")

	return c.CommonUsecase.DB().Exec(sql).Error
}

// AfterInsertData todo:lgl
func (c *DataEntryUsecase) AfterInsertData(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList,
	operUser *TData) error {

	/*
		if kindEntity.Kind == Kind_client_tasks {
			if recognizeFieldName == DataEntry_gid {
				c.log.Info("AfterInsertData:1")
				var whatGids []string
				var whoGids []string
				for rowUniqId, v := range dataEntryOperResult {
					if v.IsNewRecord {
						row := lib.TypeMap(sourceData.Get(DataEntry_gid, rowUniqId))
						if row != nil {
							whatGid := row.GetString(TaskFieldName_what_id_gid)
							whoGid := row.GetString(TaskFieldName_who_id_gid)
							if whatGid != "" {
								whatGids = append(whatGids, whatGid)
							}
							if whoGid != "" {
								whoGids = append(whoGids, whoGid)
							}
						}
					}
				}
				c.log.Info("AfterInsertData:2 whatGids", whatGids)
				c.log.Info("AfterInsertData:2 whoGids", whoGids)
				err := c.QueueUsecase.PushClientTaskHandleWhatGidJobTasks(context.TODO(), whatGids)
				if err != nil {
					c.log.Error(err, "whatGids:", whatGids)
				}
				err = c.QueueUsecase.PushClientTaskHandleWhoGidJobTasks(context.TODO(), whoGids)
				if err != nil {
					c.log.Error(err, "whoGids:", whatGids)
				}
			}
		}*/
	modifiedBy := ""
	if operUser != nil {
		modifiedBy = operUser.Gid()
	}
	c.EventBus.Publish(EventBus_AfterInsertData, kindEntity, structField, recognizeFieldName, dataEntryOperResult, sourceData, modifiedBy)
	return nil
}

// InsertData 插入新数据
func (c *DataEntryUsecase) InsertData(kindEntity KindEntity, tableName string, newData TypeDataEntryList,
	structField *TypeFieldStruct, recognizeFieldName string, operUser *TData) (dataEntryOperResult DataEntryOperResult, err error) {

	dataEntryOperResult = make(DataEntryOperResult)

	if len(newData) == 0 {
		return dataEntryOperResult, nil
	}
	allFieldNames := newData.AllFieldNames()
	if recognizeFieldName != DataEntry_Incr_id_name {
		// 需要删除系统id
		allFieldNames.DeleteFieldName(DataEntry_Incr_id_name)
	}
	allFieldNames.DeleteFieldName(DataEntry_created_at)
	allFieldNames.DeleteFieldName(DataEntry_updated_at)
	allFieldNames.DeleteFieldName(DataEntry_created_by)
	allFieldNames.DeleteFieldName(DataEntry_modified_by)

	if len(allFieldNames) == 0 {
		return dataEntryOperResult, errors.New("无字段数据")
	}
	listFields := allFieldNames.ToSqlSelect()

	var fieldDependList FieldDependList
	if configs.EnabledDataEntryDependField {
		fieldDependList = structField.GetFieldDependList().GetFieldNamesIsDependent(listFields)
		for _, v := range fieldDependList {
			allFieldNames[v.FieldName] = true
		}
		listFields = allFieldNames.ToSqlSelect()
	}

	fieldStr := "`" + strings.Join(listFields, "`,`") + "`" + ",created_at,updated_at"
	if tableName == TableName_client_cases {
		fieldStr += "," + FieldName_uniqcode
	}
	if operUser != nil {
		fieldStr += fmt.Sprintf(",%s,%s", DataEntry_created_by, DataEntry_modified_by)
	}
	sql := fmt.Sprintf("insert into %s (%s) values ", tableName, fieldStr)
	var values []string
	destDataEntryOperResult := make(DataEntryOperResult)

	for _, v := range newData {
		s := "("
		var sVals []string
		for _, fileName := range listFields {
			field := structField.GetByFieldName(fileName)
			if field == nil {
				return dataEntryOperResult, errors.New("字段：" + fileName + "不存在")
			}

			if configs.EnabledDataEntryDependField {
				fieldDepend := fieldDependList.GetByFieldName(fileName)
				if fieldDepend != nil {
					var dependValues []interface{}
					for _, v2 := range fieldDepend.FieldDependItemList {
						dependFd := structField.GetByFieldName(v2.FieldName)
						if dependFd == nil {
							return nil, errors.New("dependFd is nil")
						}
						dependValues = append(dependValues, dependFd.TransToCorrectValueFormat(v[v2.FieldName]))
					}

					var newDependValue string
					for _, v3 := range dependValues {
						if v3 != nil {
							if newDependValue == "" {
								newDependValue = InterfaceToString(v3)
							} else {
								newDependValue += " " + InterfaceToString(v3)
							}
						}
					}
					if newDependValue != "" {
						v[fileName] = newDependValue
					}
				}
			}

			if _, ok := v[fileName]; ok {
				afterFormatVal := field.TransToCorrectValueFormat(v[fileName])
				sVals = append(sVals, fmt.Sprintf("%s", lib.SqlBindValue(afterFormatVal)))
			} else {
				sVals = append(sVals, fmt.Sprintf("%s", lib.SqlBindValue(field.DefaultValue())))
			}
		}
		s += strings.Join(sVals, ",")
		s += fmt.Sprintf(",%d,%d", time.Now().Unix(), time.Now().Unix())
		if tableName == TableName_client_cases {
			uniqCode, _ := c.UniqueCodeGeneratorUsecase.GenUuid(UniqueCodeGenerator_Type_ClientUniqCode, 0)
			s += "," + lib.SqlBindValue(uniqCode)
		}
		if operUser != nil {
			s += "," + lib.SqlBindValue(operUser.Gid())
			s += "," + lib.SqlBindValue(operUser.Gid())
		}
		s += ")"
		values = append(values, s)
		row := lib.TypeMap(v)
		recognizeValue := row.GetString(recognizeFieldName)
		if recognizeValue != "" {
			destDataEntryOperResult[recognizeValue] = DataEntryOperItem{
				//RecognizeValue: recognizeValue,
				IsNewRecord: true,
			}
		}
	}

	sql += strings.Join(values, ",")

	var inertLastId int64
	// todo:lgl 先行验证，后续全部放开
	if IsCustomKind(kindEntity.Kind) {
		c.CommonUsecase.DB().Transaction(func(tx *gorm.DB) error {
			err = tx.Exec(sql).Error
			if err != nil {
				return err
			}

			var lastIds []int64
			err = tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &lastIds).Error
			if err != nil {
				return err
			}
			if len(lastIds) == 0 {
				return errors.New("LAST_INSERT_ID error")
			}
			inertLastId = lastIds[0]
			//lastInsertSql := `select LAST_INSERT_ID()`
			//var lastInertId int64
			//err = tx.Exec(lastInsertSql).Scan(&lastInertId).Error // the request id always 0
			//if err != nil {
			//	return err
			//}
			//c.log.Info("lastInsertSql lastInsertSql lastInertId:", lastIds)
			return nil
		})
	} else {
		err = c.CommonUsecase.DB().Exec(sql).Error
	}
	if recognizeFieldName == DataEntry_Incr_id_name {
		destDataEntryOperResult[InterfaceToString(inertLastId)] = DataEntryOperItem{
			//RecognizeValue: recognizeValue,
			IsNewRecord: true,
		}
	}

	if err != nil {
		return dataEntryOperResult, err
	}
	dataEntryOperResult = destDataEntryOperResult

	err = c.AfterInsertData(kindEntity, structField, recognizeFieldName, dataEntryOperResult, newData, operUser)
	if err != nil {
		return dataEntryOperResult, err
	}

	err = c.HandleTimelinesForAdded(kindEntity, newData, operUser)
	if err != nil {
		return dataEntryOperResult, err
	}

	err = c.HandleTimelinesRelaKindForAdded(kindEntity, newData, operUser)
	if err != nil {
		return dataEntryOperResult, err
	}

	return dataEntryOperResult, nil
}

func (c *DataEntryUsecase) HandleTimelinesRelaKindForAdded(kindEntity KindEntity, newData TypeDataEntryList,
	operUser *TData) error {
	if kindEntity.Kind == Kind_notes {
		userGid := ""
		if operUser != nil {
			userGid = operUser.Gid()
		}
		var destDataList TypeDataEntryList
		for _, v := range newData {
			row := lib.TypeMap(v)

			gid := row.GetString(DataEntry_gid)
			kindGid := row.GetString(Notes_FieldName_kind_gid)
			kind := row.GetString(Notes_FieldName_kind)
			content := row.GetString(Notes_FieldName_content)
			if kindGid != "" && gid != "" {
				var timelineForNotes TimelineForNotes
				timelineForNotes.Content = content

				timelineGid := uuid.UuidWithoutStrike()
				timelineData := make(TypeDataEntry)
				timelineData[Timeline_FieldName_kind] = kindEntity.Kind
				timelineData[Timeline_FieldName_kind_gid] = gid
				timelineData[Timeline_FieldName_related_kind] = kind
				timelineData[Timeline_FieldName_related_kind_gid] = kindGid
				timelineData[Timeline_FieldName_action] = Timeline_action_added
				timelineData[Timeline_FieldName_notes] = InterfaceToString(timelineForNotes)
				timelineData[DataEntry_created_by] = userGid
				timelineData[DataEntry_modified_by] = userGid
				timelineData[DataEntry_gid] = timelineGid
				destDataList = append(destDataList, timelineData)
			}
		}

		if len(destDataList) > 0 {
			_, err := c.Handle(Kind_timelines, destDataList, DataEntry_gid, operUser)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *DataEntryUsecase) HandleTimelinesForAdded(kindEntity KindEntity, newData TypeDataEntryList,
	operUser *TData) error {

	if kindEntity.NoTimelines == NoTimelines_Yes {
		return nil
	}
	userGid := ""
	if operUser != nil {
		userGid = operUser.Gid()
	}
	var destDataList TypeDataEntryList

	for _, v := range newData {
		if gid, ok := v[DataEntry_gid]; ok && gid != "" {
			timelineGid := uuid.UuidWithoutStrike()
			timelineData := make(TypeDataEntry)
			timelineData["kind"] = kindEntity.Kind
			timelineData["kind_gid"] = gid
			timelineData["action"] = Timeline_action_added
			timelineData["notes"] = ""
			timelineData["created_by"] = userGid
			timelineData["modified_by"] = userGid
			timelineData["gid"] = timelineGid
			destDataList = append(destDataList, timelineData)
		}
	}

	if len(destDataList) > 0 {
		_, err := c.Handle(Kind_timelines, destDataList, DataEntry_gid, operUser)
		if err != nil {
			return err
		}
	}
	return nil
}

type DataEntryModifyDataMap map[string]ModifyDataVo

// HasModifyDataVo 字段是否有日志
func (c DataEntryModifyDataMap) HasModifyDataVo(fieldName string) bool {
	if _, ok := c[fieldName]; ok {
		return true
	}
	return false
}

type ModifyDataVo struct {
	OldVal interface{}
	NewVal interface{}
}

func (c *ModifyDataVo) GetOldVal(fieldType string) string {
	return ChangeHistoryValueFormat(fieldType, c.OldVal)
}

func (c *ModifyDataVo) GetNewVal(fieldType string) string {
	return ChangeHistoryValueFormat(fieldType, c.NewVal)
}

func (c *DataEntryUsecase) HandleFieldDepend(structField *TypeFieldStruct, destRow map[string]ModifyDataVo, dbRow map[string]interface{}) (map[string]ModifyDataVo, error) {

	if destRow == nil && len(destRow) == 0 {
		return nil, nil
	}
	if structField == nil {
		return nil, errors.New("HandleFieldDepend: structField is nil")
	}
	var newValueFieldNames []string
	for k, _ := range destRow {
		newValueFieldNames = append(newValueFieldNames, k)
	}
	dependInfos := structField.GetFieldDependList().GetFieldNamesIsDependent(newValueFieldNames)
	dependRow := make(map[string]ModifyDataVo)
	for _, v := range dependInfos {
		fieldEntity := structField.GetByFieldName(v.FieldName)
		if fieldEntity == nil {
			return nil, errors.New("HandleFieldDepend: fieldEntity is nil")
		}
		var values []string
		for _, v1 := range v.FieldDependItemList {
			val, err := c.GetLastValueByFieldName(structField, v1.FieldName, destRow, dbRow)
			if err != nil {
				return nil, err
			}
			if val != nil {
				values = append(values, InterfaceToString(val))
			}
		}
		var newValueStr string
		for _, v3 := range values {
			if newValueStr == "" {
				if v3 != "" {
					newValueStr = v3
				}
			} else {
				if v3 != "" {
					newValueStr += " " + v3
				}
			}
		}
		newValue := fieldEntity.TransToCorrectValueFormat(newValueStr)
		var oldValue interface{}
		if _, ok := dbRow[v.FieldName]; ok {
			oldValue = fieldEntity.TransToCorrectValueFormat(dbRow[v.FieldName])
		}
		if InterfaceToString(newValue) != InterfaceToString(oldValue) {
			dependRow[v.FieldName] = ModifyDataVo{
				OldVal: oldValue,
				NewVal: newValue,
			}
		}
	}

	for k, _ := range dependRow {
		destRow[k] = dependRow[k]
	}
	return destRow, nil
}

func (c *DataEntryUsecase) GetLastValueByFieldName(structField *TypeFieldStruct, fieldName string, destRow map[string]ModifyDataVo, dbRow map[string]interface{}) (value interface{}, err error) {
	if structField == nil {
		return nil, errors.New("GetLastValueByFieldName: structField is nil")
	}
	fieldEntity := structField.GetByFieldName(fieldName)
	if fieldEntity == nil {
		return nil, errors.New("GetLastValueByFieldName: fieldEntity is nil")
	}
	isOk := false
	if destRow != nil {
		if _, ok := destRow[fieldName]; ok {
			isOk = true
			value = destRow[fieldName].NewVal
		}
	}
	if !isOk {
		if _, ok := dbRow[fieldName]; ok {
			value = dbRow[fieldName]
		}
	}
	value = fieldEntity.TransToCorrectValueFormat(value)
	return value, nil
}

// UpdateModifyData 修改数据
func (c *DataEntryUsecase) UpdateModifyData(kindEntity KindEntity,
	tableName string,
	dbRow map[string]interface{},
	newRow map[string]interface{},
	structField *TypeFieldStruct,
	recognizeFieldName string,
	operUser *TData,
) (isUpdated bool, err error, dataEntryModifyDataMap DataEntryModifyDataMap) {

	destRow := make(DataEntryModifyDataMap)
	for k, v := range newRow {
		field := structField.GetByFieldName(k)
		if field == nil {
			c.log.Error("field不存在:" + k)
			continue
		}
		v := field.TransToCorrectValueFormat(v)
		if _, ok := dbRow[k]; ok {
			if field.FieldType == FieldType_decimal {

				if dbRow[k] == nil && v == nil {
					continue
				}

				o1, _ := decimal.NewFromString(InterfaceToString(dbRow[k]))
				if v == nil {
					destRow[k] = ModifyDataVo{
						OldVal: o1.String(),
						NewVal: v,
					}
				} else if dbRow[k] == nil {
					destRow[k] = ModifyDataVo{
						OldVal: nil,
						NewVal: v,
					}
				} else {
					n1, _ := decimal.NewFromString(InterfaceToString(v))
					if !(o1.Equal(n1)) {
						destRow[k] = ModifyDataVo{
							OldVal: o1.String(),
							NewVal: n1.String(),
						}
					}
				}
			} else if field.FieldType == FieldType_number || field.FieldType == FieldType_switch {
				if v == nil && dbRow[k] == nil {
					continue
				}
				if v == nil {
					destRow[k] = ModifyDataVo{
						OldVal: dbRow[k],
						NewVal: v,
					}
				} else if dbRow[k] == nil {
					destRow[k] = ModifyDataVo{
						OldVal: nil,
						NewVal: v,
					}
				} else {
					if lib.InterfaceToInt32(dbRow[k]) != lib.InterfaceToInt32(v) {
						destRow[k] = ModifyDataVo{
							OldVal: dbRow[k],
							NewVal: v,
						}
					}
				}
			} else if lib.InterfaceToString(dbRow[k]) != lib.InterfaceToString(v) {
				destRow[k] = ModifyDataVo{
					OldVal: dbRow[k],
					NewVal: v,
				}
			}

		} else {
			if v != nil {
				destRow[k] = ModifyDataVo{
					OldVal: nil,
					NewVal: v,
				}
			}
		}
	}

	if configs.EnabledDataEntryDependField {
		// todo:lgl 此处很关键，需要测试, 发现有问题立刻注释
		destRow, err = c.HandleFieldDepend(structField, destRow, dbRow)
		if err != nil {
			return false, err, nil
		}
	}

	if len(destRow) > 0 {

		updateData := make(map[string]interface{})
		for k, v := range destRow {
			updateData[k] = v.NewVal
		}
		updateData["updated_at"] = time.Now().Unix()
		if operUser != nil {
			updateData[DataEntry_modified_by] = operUser.Gid()
		}
		if _, ok := dbRow[DataEntry_Incr_id_name]; ok {
			err := c.CommonUsecase.DB().Table(tableName).Where("id=?", dbRow[DataEntry_Incr_id_name]).
				Updates(updateData).Error
			if err != nil {
				c.log.Error(err)
				return false, err, nil
			} else {
				err := c.HandleChangeHistories(kindEntity, dbRow[DataEntry_Incr_id_name], destRow, structField)
				if err != nil {
					c.log.Error("HandleChangeHistories:", err, kindEntity.Kind, dbRow[DataEntry_Incr_id_name])
					return true, err, destRow
				}
				err = c.HandleTimelines(kindEntity, InterfaceToString(dbRow[DataEntry_gid]), destRow, structField, operUser)
				if err != nil {
					c.log.Error("HandleTimelines:", err, kindEntity.Kind, dbRow[DataEntry_Incr_id_name])
					return true, err, destRow
				}
			}
			// 此处返回已更新
			return true, nil, destRow
		} else {
			c.log.Error("DataEntry_Incr_id_name不存在")
		}
	}
	return false, nil, nil
}

func (c *DataEntryUsecase) HandleTimelines(kindEntity KindEntity, gid string,
	data map[string]ModifyDataVo,
	structField *TypeFieldStruct, operUser *TData) error {

	if kindEntity.NoTimelines == NoTimelines_Yes {
		return nil
	}

	if data == nil {
		return nil
	}

	var fieldHistory TimelineFieldHistory
	for fieldName, modifyDataVo := range data {
		field := structField.GetByFieldName(fieldName)

		if field.NoTimelines == Field_NoTimelines_Yes {
			continue
		}

		oldVal := field.DefaultValue()
		if modifyDataVo.OldVal != nil {
			oldVal = modifyDataVo.OldVal
		}
		newVal := modifyDataVo.NewVal

		destOldVal := ChangeHistoryValueFormat(field.FieldType, oldVal)
		destNewVal := ChangeHistoryValueFormat(field.FieldType, newVal)

		fieldHistory = append(fieldHistory, TimelineFieldHistoryItem{
			FieldName: fieldName,
			NewValue:  destNewVal,
			OldValue:  destOldVal,
		})
	}

	if len(fieldHistory) != 0 {

		var timelineFieldHistoryNotes TimelineFieldHistoryNotes
		timelineFieldHistoryNotes.FieldHistory = fieldHistory
		userGid := ""
		if operUser != nil {
			userGid = operUser.Gid()
		}

		timelineGid := uuid.UuidWithoutStrike()
		timelineData := make(TypeDataEntry)
		timelineData["kind"] = kindEntity.Kind
		timelineData["kind_gid"] = gid
		timelineData["action"] = Timeline_action_updated
		timelineData["notes"] = timelineFieldHistoryNotes.ToString()
		timelineData["created_by"] = userGid
		timelineData["modified_by"] = userGid
		timelineData["gid"] = timelineGid

		_, err := c.HandleOne(Kind_timelines, timelineData, DataEntry_gid, operUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DataEntryUsecase) HandleChangeHistories(kindEntity KindEntity, incrId interface{},
	data map[string]ModifyDataVo,
	structField *TypeFieldStruct) error {

	if kindEntity.NoChangeHistory == NoChangeHistory_Yes {
		return nil
	}

	if data == nil {
		return nil
	}

	incrIdInt, err := strconv.ParseInt(lib.InterfaceToString(incrId), 10, 32)
	if err != nil {
		c.log.Error(err)
	}

	var logs []ChangeHistoryEntity
	for fieldName, modifyDataVo := range data {
		field := structField.GetByFieldName(fieldName)
		oldVal := field.DefaultValue()
		if modifyDataVo.OldVal != nil {
			oldVal = modifyDataVo.OldVal
		}
		newVal := modifyDataVo.NewVal
		e := ChangeHistoryEntity{
			Kind:      kindEntity.Kind,
			IncrId:    int32(incrIdInt),
			FieldName: field.FieldName,
			OldValue:  ChangeHistoryValueFormat(field.FieldType, oldVal),
			NewValue:  ChangeHistoryValueFormat(field.FieldType, newVal),
			CreatedAt: time.Now().Unix(),
		}
		logs = append(logs, e)
	}
	if len(logs) == 0 {
		return nil
	}
	return c.CommonUsecase.DB().Create(logs).Error
}

type DbRowsList []map[string]interface{}

func (c DbRowsList) GetByFieldNameValue(fieldNameValue string, fieldName string) map[string]interface{} {
	for k, _ := range c {
		for k1, v1 := range c[k] {
			if k1 == fieldName {
				if lib.InterfaceToString(v1) == fieldNameValue {
					return c[k]
				}
			}
		}
	}
	return nil
}
