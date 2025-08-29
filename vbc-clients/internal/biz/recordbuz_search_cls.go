package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	. "vbc/lib/builder"
)

type RecordbuzSearchCls struct {
	alias                       int
	TFilterUsecase              *TFilterUsecase
	KindUsecase                 *KindUsecase
	BUsecase                    *BUsecase
	PermissionDataFilterUsecase *PermissionDataFilterUsecase
	TUsecase                    *TUsecase
	NormalUserOnlyOwner         bool
}

func CreateRecordbuzSearchCls(TFilterUsecase *TFilterUsecase,
	KindUsecase *KindUsecase,
	BUsecase *BUsecase,
	PermissionDataFilterUsecase *PermissionDataFilterUsecase,
	TUsecase *TUsecase,
	NormalUserOnlyOwner bool) *RecordbuzSearchCls {
	return &RecordbuzSearchCls{
		TFilterUsecase:              TFilterUsecase,
		KindUsecase:                 KindUsecase,
		BUsecase:                    BUsecase,
		PermissionDataFilterUsecase: PermissionDataFilterUsecase,
		TUsecase:                    TUsecase,
		NormalUserOnlyOwner:         NormalUserOnlyOwner,
	}
}

func TidyTableFieldForSql(fieldName string, aliasTableName string) string {
	if aliasTableName == "" {
		return fieldName
	}
	return fmt.Sprintf("%s.%s", aliasTableName, fieldName)
}

func (c *RecordbuzSearchCls) GenAliasTableName() string {
	c.alias += 1
	return fmt.Sprintf("t%d", c.alias)
}

// GenBuilder (list) sortFieldEntity 排序字段， 如果没有就算了 timezoneId：时区，为空时：默认时区：America/Los_Angeles
func (c *RecordbuzSearchCls) GenBuilder(userFacade *UserFacade, kindEntity KindEntity, request *TListRequest, recordFieldSortList RecordFieldSortList, timezoneId string) (builder *Builder, err error) {

	recordTableName, err := c.KindUsecase.CacheTableNameByKind(kindEntity.Kind)
	if err != nil {
		return nil, err
	}

	driveAliasTableName := c.GenAliasTableName()
	query := Dialect(MYSQL).Select(fmt.Sprintf("%s.gid", driveAliasTableName)).
		From(recordTableName, driveAliasTableName)

	err = c.Filter(kindEntity.Kind, request, query, timezoneId, driveAliasTableName)
	if err != nil {
		return nil, err
	}
	err = c.FilterDataPermission(kindEntity, userFacade, query, driveAliasTableName)
	if err != nil {
		return nil, err
	}
	//err = c.HandleSort(query, kindEntity, sortFieldEntity, sortOrder, driveAliasTableName)

	//todo:lgl begin test code
	//var fieldNames = []string{"sys__due_date", "stages", "user_gid"}
	//caseStruct, _ := c.BUsecase.FieldUsecase.StructByKind(Kind_client_cases)
	//var recordFieldSortList RecordFieldSortList

	//for _, v := range fieldNames {
	//	fieldEntity := caseStruct.GetByFieldName(v)
	//	recordFieldSortList = append(recordFieldSortList, RecordFieldSort{
	//		SortFieldEntity: *fieldEntity,
	//		SortOrder:       Fab_SortOrder_desc,
	//	})
	//}
	err = c.HandleSortV2(query, kindEntity, recordFieldSortList, driveAliasTableName)
	// todo:lgl end test code

	if err != nil {
		return nil, err
	}

	return query, nil
}

func (c *RecordbuzSearchCls) HasDeletePermission(userFacade UserFacade, userProfile TData, kindEntity KindEntity, gid string) (hasDeletePermission bool, err error) {
	if IsAdminProfile(&userProfile) {
		return true, nil
	}
	return false, nil
}

// HasPermissionRow 判断数据是否有权限
func (c *RecordbuzSearchCls) HasPermissionRow(userFacade UserFacade, kindEntity KindEntity, gid string) (hasPermission bool, err error) {
	builder, err := c.GenBuilderRow(&userFacade, kindEntity, Eq{DataEntry_gid: gid})
	if err != nil {
		return false, err
	}
	sql, err := builder.ToBoundSQL()
	if err != nil {
		return false, err
	}
	tData, err := c.TUsecase.DataByRawSql(kindEntity, sql)
	if err != nil {
		return false, err
	}
	if tData != nil {
		return true, nil
	}
	return false, nil
}

// GenBuilderRow (row) 主要验证是否有此数据权限
func (c *RecordbuzSearchCls) GenBuilderRow(userFacade *UserFacade, kindEntity KindEntity, cond Cond) (builder *Builder, err error) {

	recordTableName, err := c.KindUsecase.CacheTableNameByKind(kindEntity.Kind)
	if err != nil {
		return nil, err
	}
	if cond == nil {
		return nil, errors.New("cond is nil")
	}

	driveAliasTableName := c.GenAliasTableName()
	query := Dialect(MYSQL).Select(fmt.Sprintf("%s.gid", driveAliasTableName)).
		From(recordTableName, driveAliasTableName)

	query.And(Eq{TidyTableFieldForSql(DataEntry_biz_deleted_at, driveAliasTableName): 0})
	query.And(Eq{TidyTableFieldForSql(DataEntry_deleted_at, driveAliasTableName): 0})
	query.And(cond)

	err = c.FilterDataPermission(kindEntity, userFacade, query, driveAliasTableName)
	if err != nil {
		return nil, err
	}
	return query, nil
}

type RecordFieldSortList []RecordFieldSort

type RecordFieldSort struct {
	SortFieldEntity FieldEntity
	SortOrder       string
}

func (c *RecordbuzSearchCls) HandleSortV2(query *Builder, kindEntity KindEntity, recordFieldSortList RecordFieldSortList, driveAliasTableName string) error {
	if len(recordFieldSortList) == 0 {
		query.OrderBy(TidyTableFieldForSql("id", driveAliasTableName) + " desc")
		return nil
	}

	var orderbys []string
	relaTableNameMap := make(map[string]string)
	for _, v := range recordFieldSortList {
		sortFieldEntity := v.SortFieldEntity
		sortOrder := v.SortOrder
		var st string
		if sortFieldEntity.FieldType == FieldType_lookup {
			relaTableName, err := c.KindUsecase.CacheTableNameByKind(sortFieldEntity.RelaKind)
			if err != nil {
				return err
			}
			var lookupAliasTableName string
			if _, ok := relaTableNameMap[relaTableName]; !ok {
				lookupAliasTableName = c.GenAliasTableName()
				var conds []Cond
				conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName), TidyTableFieldForSql("gid", lookupAliasTableName))))
				query.LeftJoin(fmt.Sprintf("%s as %s", relaTableName, lookupAliasTableName), And(conds...))
			} else {
				lookupAliasTableName = relaTableNameMap[relaTableName]
			}

			st = fmt.Sprintf("%s %s", TidyTableFieldForSql(sortFieldEntity.RelaName, lookupAliasTableName), sortOrder)
		} else if sortFieldEntity.FieldType == FieldType_dropdown {
			dropDownAliasTableName := c.GenAliasTableName()
			var conds []Cond

			if sortFieldEntity.OptionGroupName == "" {
				conds = append(conds, Eq{TidyTableFieldForSql("kind", dropDownAliasTableName): kindEntity.Kind})
				conds = append(conds, Eq{TidyTableFieldForSql("field_name", dropDownAliasTableName): sortFieldEntity.FieldName})
				conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql("option_value", dropDownAliasTableName), TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName))))
			} else {
				conds = append(conds, Eq{TidyTableFieldForSql("kind", dropDownAliasTableName): Kind_global_sets})
				conds = append(conds, Eq{TidyTableFieldForSql("field_name", dropDownAliasTableName): sortFieldEntity.OptionGroupName})
				conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql("option_value", dropDownAliasTableName), TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName))))
			}

			//query.Select(fmt.Sprintf("%s.gid,%s.option_label,%s.option_sort", driveAliasTableName, dropDownAliasTableName, dropDownAliasTableName))
			query.LeftJoin(fmt.Sprintf("%s as %s", FieldOptionEntity{}.TableName(), dropDownAliasTableName), And(conds...))
			if sortFieldEntity.SortType == Field_SortType_Sort {
				st = fmt.Sprintf("%s %s", TidyTableFieldForSql("option_sort", dropDownAliasTableName), sortOrder)
			} else {
				st = fmt.Sprintf("%s %s", TidyTableFieldForSql("option_label", dropDownAliasTableName), sortOrder)
			}
		} else {
			st = fmt.Sprintf("%s %s", TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName), sortOrder)
		}
		orderbys = append(orderbys, st)
	}

	ob := strings.Join(orderbys, ",")
	query.OrderBy(ob)

	sql, _ := query.ToBoundSQL()
	c.BUsecase.log.Info("sql:", sql)

	return nil
}

func (c *RecordbuzSearchCls) HandleSort(query *Builder, kindEntity KindEntity, sortFieldEntity *FieldEntity, sortOrder string, driveAliasTableName string) error {
	if sortFieldEntity == nil {
		query.OrderBy(TidyTableFieldForSql("id", driveAliasTableName) + " desc")
		return nil
	}
	if sortFieldEntity.FieldType == FieldType_lookup {

		relaTableName, err := c.KindUsecase.CacheTableNameByKind(sortFieldEntity.RelaKind)
		if err != nil {
			return err
		}
		lookupAliasTableName := c.GenAliasTableName()
		var conds []Cond
		conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName), TidyTableFieldForSql("gid", lookupAliasTableName))))
		query.LeftJoin(fmt.Sprintf("%s as %s", relaTableName, lookupAliasTableName), And(conds...))
		query.OrderBy(fmt.Sprintf("%s %s", TidyTableFieldForSql(sortFieldEntity.RelaName, lookupAliasTableName), sortOrder))
	} else if sortFieldEntity.FieldType == FieldType_dropdown {
		dropDownAliasTableName := c.GenAliasTableName()
		var conds []Cond

		if sortFieldEntity.OptionGroupName == "" {
			conds = append(conds, Eq{TidyTableFieldForSql("kind", dropDownAliasTableName): kindEntity.Kind})
			conds = append(conds, Eq{TidyTableFieldForSql("field_name", dropDownAliasTableName): sortFieldEntity.FieldName})
			conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql("option_value", dropDownAliasTableName), TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName))))
		} else {
			conds = append(conds, Eq{TidyTableFieldForSql("kind", dropDownAliasTableName): Kind_global_sets})
			conds = append(conds, Eq{TidyTableFieldForSql("field_name", dropDownAliasTableName): sortFieldEntity.OptionGroupName})
			conds = append(conds, Expr(fmt.Sprintf("%s=%s", TidyTableFieldForSql("option_value", dropDownAliasTableName), TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName))))
		}

		//query.Select(fmt.Sprintf("%s.gid,%s.option_label,%s.option_sort", driveAliasTableName, dropDownAliasTableName, dropDownAliasTableName))
		query.LeftJoin(fmt.Sprintf("%s as %s", FieldOptionEntity{}.TableName(), dropDownAliasTableName), And(conds...))
		if sortFieldEntity.SortType == Field_SortType_Sort {
			query.OrderBy(fmt.Sprintf("%s %s", TidyTableFieldForSql("option_sort", dropDownAliasTableName), sortOrder))
		} else {
			query.OrderBy(fmt.Sprintf("%s %s", TidyTableFieldForSql("option_label", dropDownAliasTableName), sortOrder))
		}
	} else {
		query.OrderBy(fmt.Sprintf("%s %s", TidyTableFieldForSql(sortFieldEntity.FieldName, driveAliasTableName), sortOrder))
	}
	return nil
}

func (c *RecordbuzSearchCls) Filter(kind string, request *TListRequest, query *Builder, timezoneId string, aliasTableName string) error {
	if request != nil {
		conds, err := c.TFilterUsecase.Do(kind, *request, timezoneId, aliasTableName)
		if err != nil {
			return err
		}
		if conds != nil {
			query.And(conds)
		}
	}
	query.And(Eq{TidyTableFieldForSql(DataEntry_biz_deleted_at, aliasTableName): 0})
	query.And(Eq{TidyTableFieldForSql(DataEntry_deleted_at, aliasTableName): 0})
	return nil
}

func (c *RecordbuzSearchCls) FilterDataPermission(kindEntity KindEntity, userFacade *UserFacade, query *Builder, aliasTableName string) error {

	cond, err := c.FilterDataPermissionCond(kindEntity, userFacade, aliasTableName, c.NormalUserOnlyOwner)
	if err != nil {
		return err
	}
	if cond != nil {
		query.And(cond)
	}
	return nil
}

// FilterDataPermissionCond normalUserOnlyOwner 所有用户都使用owner数据
func (c *RecordbuzSearchCls) FilterDataPermissionCond(kindEntity KindEntity, userFacade *UserFacade, aliasTableName string, normalUserOnlyOwner bool) (Cond, error) {

	if userFacade == nil {
		return nil, nil
	}

	tProfile, err := userFacade.RelaData(c.BUsecase, User_FieldName_profile_gid)
	if err != nil {
		return nil, err
	}
	if tProfile == nil {
		return nil, errors.New("tProfile is nil")
	}
	if IsAdminProfile(tProfile) || HaveAllDataPermissions(kindEntity.Kind, tProfile.Gid()) {
		if normalUserOnlyOwner {
			var conds []Cond
			conds = append(conds, Eq{TidyTableFieldForSql(DataEntry_user_gid, aliasTableName): userFacade.Gid()})
			if len(conds) > 0 {
				return And(conds...), nil
			}
		}
		return nil, nil
	} else {
		cond, err := c.PermissionDataFilterUsecase.Filter(kindEntity, tProfile, userFacade, aliasTableName, normalUserOnlyOwner)
		if err != nil {
			return nil, err
		}
		return cond, nil
	}
	return nil, nil
}
