package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type RecordbuzUsecase struct {
	log                         *log.Helper
	CommonUsecase               *CommonUsecase
	conf                        *conf.Data
	FieldUsecase                *FieldUsecase
	KindUsecase                 *KindUsecase
	TUsecase                    *TUsecase
	TFilterUsecase              *TFilterUsecase
	SettingCustomViewUsecase    *SettingCustomViewUsecase
	PermissionDataFilterUsecase *PermissionDataFilterUsecase
	BUsecase                    *BUsecase
	FieldPermissionUsecase      *FieldPermissionUsecase
	FieldbuzUsecase             *FieldbuzUsecase
	RecordbuzSearchUsecase      *RecordbuzSearchUsecase
}

func NewRecordbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	KindUsecase *KindUsecase,
	TUsecase *TUsecase,
	TFilterUsecase *TFilterUsecase,
	SettingCustomViewUsecase *SettingCustomViewUsecase,
	PermissionDataFilterUsecase *PermissionDataFilterUsecase,
	BUsecase *BUsecase,
	FieldPermissionUsecase *FieldPermissionUsecase,
	FieldbuzUsecase *FieldbuzUsecase,
	RecordbuzSearchUsecase *RecordbuzSearchUsecase) *RecordbuzUsecase {
	uc := &RecordbuzUsecase{
		log:                         log.NewHelper(logger),
		CommonUsecase:               CommonUsecase,
		conf:                        conf,
		FieldUsecase:                FieldUsecase,
		KindUsecase:                 KindUsecase,
		TUsecase:                    TUsecase,
		TFilterUsecase:              TFilterUsecase,
		SettingCustomViewUsecase:    SettingCustomViewUsecase,
		PermissionDataFilterUsecase: PermissionDataFilterUsecase,
		BUsecase:                    BUsecase,
		FieldPermissionUsecase:      FieldPermissionUsecase,
		FieldbuzUsecase:             FieldbuzUsecase,
		RecordbuzSearchUsecase:      RecordbuzSearchUsecase,
	}

	return uc
}

func (c *RecordbuzUsecase) Filter(kind string, request *TListRequest, query *Builder, timezoneId string) error {
	if request != nil {
		conds, err := c.TFilterUsecase.Do(kind, *request, timezoneId, "")
		if err != nil {
			return err
		}
		if conds != nil {
			query.And(conds)
		}
	}
	query.And(Eq{"biz_deleted_at": 0})
	query.And(Eq{"deleted_at": 0})
	return nil
}

func (c *RecordbuzUsecase) FilterDataPermission(kindEntity KindEntity, userFacade *UserFacade, query *Builder, normalUserOnlyOwner bool) error {

	cond, err := c.FilterDataPermissionCond(kindEntity, userFacade, normalUserOnlyOwner)
	if err != nil {
		return err
	}
	if cond != nil {
		query.And(cond)
	}
	return nil
}

func (c *RecordbuzUsecase) FilterDataPermissionCond(kindEntity KindEntity, userFacade *UserFacade, normalUserOnlyOwner bool) (Cond, error) {

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
			conds = append(conds, Eq{TidyTableFieldForSql(DataEntry_user_gid, ""): userFacade.Gid()})
			if len(conds) > 0 {
				return And(conds...), nil
			}
		}
		return nil, nil
	} else {
		cond, err := c.PermissionDataFilterUsecase.Filter(kindEntity, tProfile, userFacade, "", normalUserOnlyOwner)
		if err != nil {
			return nil, err
		}
		return cond, nil
	}
	return nil, nil
}

// List specificFieldNames如果有值，就不使用配置信息了 normalUserOnlyOwner: true-只要owner数据
func (c *RecordbuzUsecase) List(kind string, userFacade *UserFacade, request *TListRequest, page int, pageSize int, specificFieldNames []string, normalUserOnlyOwner bool) (tDataList TDataList, err error) {
	fields, err := c.FieldUsecase.ListByKind(kind)
	if err != nil {
		return nil, err
	}

	timezoneId := ""
	if userFacade != nil {
		timezoneId = userFacade.TimezoneId()

		var newFields TypeFieldList
		// 进行字段权限过过虑
		fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
		if err != nil {
			return nil, err
		}

		var displayFieldNames map[string]bool
		if len(specificFieldNames) > 0 {
			displayFieldNames = make(map[string]bool)
			for _, v := range specificFieldNames {
				displayFieldNames[v] = true
			}
		} else {
			displayFieldNames = c.FieldbuzUsecase.GetDisplayFieldNameForRecords(kind, userFacade, request.TableType)
		}

		for k, v := range fields {
			fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
			if err != nil {
				return nil, err
			}
			if fieldPermissionVo.CanShow() {
				isOk := false
				if lib.InArray(v.FieldName, MustReturnFieldNamesForRecords) {
					isOk = true
				} else {
					if _, ok := displayFieldNames[v.FieldName]; ok {
						isOk = true
					}
				}
				if isOk {
					newFields = append(newFields, fields[k])
				}
			}
		}
		fields = newFields
	}

	fieldNames := fields.ToFieldNames()
	if len(fieldNames) == 0 {
		return nil, errors.New("No fields available")
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return nil, err
	}
	if kindEntity == nil {
		return nil, errors.New("kindEntity is nil")
	}
	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return nil, err
	}
	var fabSettingCustomView *FabSettingCustomView
	if userFacade != nil {
		scv, err := c.SettingCustomViewUsecase.Get(kind, *userFacade, request.TableType)
		if err != nil {
			return nil, err
		}
		fabSettingCustomView = &scv
	}

	isNewVersion := true
	var sql string
	structField, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}

	if isNewVersion {
		// 注意这里与计算total需要一起修改
		recordbuzSearchCls := c.RecordbuzSearchUsecase.NewRecordbuzSearchCls(normalUserOnlyOwner)
		var recordFieldSortList RecordFieldSortList
		if fabSettingCustomView != nil {
			if fabSettingCustomView.CustomView.IsCustomSortBy() {
				for _, v := range fabSettingCustomView.CustomView.SortBys {
					destFieldName := v.FieldName
					if destFieldName == DataEntry_sys__itf_formula {
						destFieldName = FieldName_itf_expiration
					}
					fieldEntity := structField.GetByFieldName(destFieldName)
					if fieldEntity != nil {
						recordFieldSortList = append(recordFieldSortList, RecordFieldSort{
							SortFieldEntity: *fieldEntity,
							SortOrder:       v.SortOrder,
						})
					}
				}
			}
		}
		builder, err := recordbuzSearchCls.GenBuilder(userFacade, *kindEntity, request, recordFieldSortList, timezoneId)
		if err != nil {
			return nil, err
		}
		sql, err = builder.Limit(pageSize, HandleOffset(page, pageSize)).ToBoundSQL()
		if err != nil {
			return nil, err
		}

		//lib.DPrintln("sql:", sql)
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
		var gids []string
		for _, v := range list {
			gids = append(gids, InterfaceToString(v["gid"]))
		}
		//lib.DPrintln("gids:", gids)
		if len(gids) > 0 {
			query := Dialect(MYSQL).Select(fieldNames...).From(tableName)
			query.And(In("gid", gids))
			query.OrderBy(fmt.Sprintf("FIELD(gid, '%s')", strings.Join(gids, "','")))
			sql, err = query.ToBoundSQL()
			if err != nil {
				return nil, err
			}
		}
	}

	c.log.Info("sql: ", sql)
	var list []map[string]interface{}

	if sql != "" {
		sqlRows, err := c.CommonUsecase.DB().Raw(sql).Rows()
		if err != nil {
			return nil, err
		}
		if sqlRows != nil {
			defer sqlRows.Close()
		}
		_, list, err = lib.SqlRowsTrans(sqlRows)
		if err != nil {
			return nil, err
		}
	}

	caches, err := c.TUsecase.GetCaches(fields, list)
	if err != nil {
		return nil, err
	}
	for k, _ := range list {
		tdata := TData{
			CustomFields: c.TUsecase.GenTFields(&caches, kind, list[k], fields),
		}
		// 处理formula
		c.TUsecase.DoFormula(*kindEntity, &tdata)
		tDataList = append(tDataList, tdata)
	}
	return tDataList, err
}

func (c *RecordbuzUsecase) Total(kind string, userFacade *UserFacade, request *TListRequest, normalUserOnlyOwner bool) (total int64, err error) {

	tableName, err := c.KindUsecase.CacheTableNameByKind(kind)
	if err != nil {
		return 0, err
	}

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return 0, err
	}
	if kindEntity == nil {
		return 0, errors.New("kindEntity is nil")
	}

	timezoneId := ""
	if userFacade != nil {
		timezoneId = userFacade.TimezoneId()
	}
	query := Dialect(MYSQL).Select("count(*) as c").From(tableName)

	err = c.Filter(kind, request, query, timezoneId)
	if err != nil {
		return 0, err
	}
	err = c.FilterDataPermission(*kindEntity, userFacade, query, normalUserOnlyOwner)
	if err != nil {
		return 0, err
	}

	sql, err := query.ToBoundSQL()
	if err != nil {
		return 0, err
	}

	return c.CommonUsecase.Count(c.CommonUsecase.DB(), sql)
}

func (c *RecordbuzUsecase) GetZohoValueByDbValue(dbValue string, field FieldEntity, fieldOptionStruct TypeFieldOptionStruct) string {
	if field.FieldType == FieldType_dropdown {
		fieldOption := fieldOptionStruct.AllByFieldName(field).GetByValue(dbValue)
		if fieldOption != nil {
			return fieldOption.OptionLabel
		}
	}
	return dbValue
}
