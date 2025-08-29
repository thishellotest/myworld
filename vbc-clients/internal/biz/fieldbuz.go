package biz

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"sort"
	"vbc/internal/conf"
	"vbc/lib"
)

type FieldbuzUsecase struct {
	log                    *log.Helper
	CommonUsecase          *CommonUsecase
	conf                   *conf.Data
	FieldUsecase           *FieldUsecase
	FieldOptionUsecase     *FieldOptionUsecase
	LongMapUsecase         *LongMapUsecase
	FieldPermissionUsecase *FieldPermissionUsecase
}

func NewFieldbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	LongMapUsecase *LongMapUsecase,
	FieldPermissionUsecase *FieldPermissionUsecase) *FieldbuzUsecase {
	uc := &FieldbuzUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		FieldUsecase:           FieldUsecase,
		FieldOptionUsecase:     FieldOptionUsecase,
		LongMapUsecase:         LongMapUsecase,
		FieldPermissionUsecase: FieldPermissionUsecase,
	}

	return uc
}

// FabColumnwidth 返回用户设置的字段宽度
func (c *FieldbuzUsecase) FabColumnwidth(kind string, userFacade *UserFacade, tableType string) (Columnwidths, error) {
	if userFacade == nil {
		return nil, errors.New("userFacade is nil")
	}
	key := MapKeyCustomViewColumnwidth(userFacade.Gid(), kind, tableType)
	val, _ := c.LongMapUsecase.GetForString(key)
	if val != "" {
		var columnwidthVo ColumnwidthVo
		err := json.Unmarshal([]byte(val), &columnwidthVo)
		if err != nil {
			return nil, err
		}
		return columnwidthVo.Columns, nil
	}

	return nil, nil
}

// FabCustomView 返回用户设置
func (c *FieldbuzUsecase) FabCustomView(kind string, userFacade *UserFacade, tableType string) (fabCustomView FabCustomView, err error) {
	if userFacade == nil {
		return fabCustomView, errors.New("userFacade is nil")
	}
	customViewKey := MapKeyCustomView(userFacade.Gid(), kind, tableType)
	customViewVal, _ := c.LongMapUsecase.GetForString(customViewKey)
	if customViewVal != "" {
		err = json.Unmarshal([]byte(customViewVal), &fabCustomView)
		if err != nil {
			return fabCustomView, err
		}
	}
	return fabCustomView, nil
}

func (c *FieldbuzUsecase) GetChangeFieldsVo(kind string, userFacade *UserFacade, tableType string) ChangeFieldsVo {
	columnsKey := MapKeyCustomViewColumns(userFacade.Gid(), kind, tableType)
	columnsVal, _ := c.LongMapUsecase.GetForString(columnsKey)
	var changeFieldsVo ChangeFieldsVo
	if columnsVal != "" {
		err := json.Unmarshal([]byte(columnsVal), &changeFieldsVo)
		if err != nil {
			c.log.Error(err)
		}
	}
	return changeFieldsVo
}

func (c *FieldbuzUsecase) GetDisplayFieldNameForRecords(kind string, userFacade *UserFacade, tableType string) map[string]bool {
	fieldNames := make(map[string]bool)
	if kind == Kind_users {
		_, checkedFieldNames, _ := c.UsersCheckedAndSortFieldNames()
		checkedFieldNames = append(checkedFieldNames, UserFieldName_role_gid)
		checkedFieldNames = append(checkedFieldNames, User_FieldName_profile_gid)
		for _, v := range checkedFieldNames {
			fieldNames[v] = true
		}
	} else if kind == Kind_attorneys {
		_, checkedFieldNames, _ := c.AttorneysCheckedAndSortFieldNames()
		for _, v := range checkedFieldNames {
			fieldNames[v] = true
		}
	} else {
		vo := c.GetChangeFieldsVo(kind, userFacade, tableType)
		if len(vo.Fields) > 0 {
			for _, v := range vo.Fields {
				if v.Checked {
					fieldNames[v.FieldName] = true
				}
			}
		} else {
			for _, v := range DefaultDisplayFieldNamesForRecords {
				fieldNames[v] = true
			}
		}
	}
	return fieldNames
}

func (c *FieldbuzUsecase) FabFieldsForBasicdata(kind string, userFacade *UserFacade) ([]FabField, error) {

	if userFacade == nil {
		return nil, errors.New("FabFields: userFacade is nil")
	}
	fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
	if err != nil {
		return nil, err
	}

	canShowFieldNames := []string{
		FieldName_stages,
	}

	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	var fabFields []FabField
	for _, v := range fieldStruct.Records {
		if kind == Kind_client_cases { // 此处为定制，临时解决heroku返回数据超限的问题
			if !lib.InArray(v.FieldName, canShowFieldNames) {
				continue
			}
		}
		if v.IsNoDisplayColumnsForUser() {
			continue
		}

		fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
		if err != nil {
			c.log.Error(err)
			continue
		}
		if fieldPermissionVo.CanShow() {

			var tUser *TData
			if userFacade != nil {
				tUser = &userFacade.TData
			}
			a := v.FieldToApi(c.FieldOptionUsecase, c.log, tUser)
			if fieldPermissionVo.CanWrite() {
				a.CanWrite = true
			}
			fabFields = append(fabFields, a)
		}
	}
	return fabFields, nil
}

func (c *FieldbuzUsecase) AttorneysCheckedAndSortFieldNames() (columnwidths Columnwidths, checkedFieldNames []string, sortFieldNames []string) {
	columnwidths = make(Columnwidths)
	columnwidths[AttorneyFieldName_first_name] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_last_name] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_email] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_ro_email] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_status] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_province] = ColumnwidthUnitVo{
		Width: 10,
	}
	columnwidths[AttorneyFieldName_city] = ColumnwidthUnitVo{
		Width: 10,
	}

	checkedFieldNames = []string{
		AttorneyFieldName_first_name,
		AttorneyFieldName_last_name,
		AttorneyFieldName_status,
		AttorneyFieldName_email,
		AttorneyFieldName_ro_email,
		AttorneyFieldName_province,
		AttorneyFieldName_city,
	}
	sortFieldNames = []string{
		AttorneyFieldName_first_name,
		AttorneyFieldName_last_name,
		AttorneyFieldName_status,
		AttorneyFieldName_email,
		AttorneyFieldName_ro_email,
		AttorneyFieldName_province,
		AttorneyFieldName_city,
	}
	return
}

func (c *FieldbuzUsecase) UsersCheckedAndSortFieldNames() (columnwidths Columnwidths, checkedFieldNames []string, sortFieldNames []string) {
	columnwidths = make(Columnwidths)
	columnwidths[UserFieldName_fullname] = ColumnwidthUnitVo{
		Width: 24,
	}
	columnwidths[UserFieldName_email] = ColumnwidthUnitVo{
		Width: 20,
	}
	columnwidths[UserFieldName_status] = ColumnwidthUnitVo{
		Width: 6,
	}
	columnwidths[UserFieldName_mobile] = ColumnwidthUnitVo{
		Width: 15,
	}
	columnwidths[UserFieldName_dialpad_phonenumber] = ColumnwidthUnitVo{
		Width: 15,
	}
	columnwidths[UserFieldName_MailSender] = ColumnwidthUnitVo{
		Width: 20,
	}
	//columnwidths[User_FieldName_profile_gid] = ColumnwidthUnitVo{
	//	Width: 10,
	//}

	checkedFieldNames = []string{
		UserFieldName_fullname,
		UserFieldName_email,
		UserFieldName_status,
		//UserFieldName_role_gid,
		UserFieldName_mobile,
		UserFieldName_dialpad_phonenumber,
		UserFieldName_MailSender,
		//User_FieldName_profile_gid,
	}
	sortFieldNames = []string{
		UserFieldName_fullname,
		UserFieldName_email,
		UserFieldName_status,
		//UserFieldName_role_gid,
		UserFieldName_mobile,
		UserFieldName_dialpad_phonenumber,
		UserFieldName_MailSender,
		//User_FieldName_profile_gid,
	}
	return
}

// FabFields 返回用户可显示的字段（有权限的字段，选中的字段，需要排序）
func (c *FieldbuzUsecase) FabFields(kind string, userFacade *UserFacade, tableType string) ([]FabField, error) {

	if userFacade == nil {
		return nil, errors.New("FabFields: userFacade is nil")
	}
	fieldPermissionCenter, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
	if err != nil {
		return nil, err
	}

	var checkedFieldNames []string
	var sortFieldNames []string

	if kind == Kind_users {
		_, checkedFieldNames, sortFieldNames = c.UsersCheckedAndSortFieldNames()
	} else if kind == Kind_attorneys {
		_, checkedFieldNames, sortFieldNames = c.AttorneysCheckedAndSortFieldNames()
	} else {
		changeFieldsVo := c.GetChangeFieldsVo(kind, userFacade, tableType)
		if len(changeFieldsVo.Fields) > 0 {
			fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
			if err != nil {
				return nil, err
			}
			if fieldStruct == nil {
				return nil, errors.New("fieldStruct is nil")
			}
			for _, v := range changeFieldsVo.Fields {
				fieldEntity := fieldStruct.GetByFieldName(v.FieldName)
				if fieldEntity == nil {
					continue
				}
				fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
				if err != nil {
					c.log.Error(err)
					continue
				}
				if fieldPermissionVo.CanShow() {
					if v.Checked {
						checkedFieldNames = append(checkedFieldNames, v.FieldName)
					}
					sortFieldNames = append(sortFieldNames, v.FieldName)
				}
			}
		}
		// 用户没有选中字段，使用默认字段
		if len(checkedFieldNames) == 0 {
			// todo: 使用kind强制关联字段
			checkedFieldNames = DefaultDisplayFieldNamesForRecords
		}
		// 用户没有设置排序时，使用默认字段
		if len(sortFieldNames) == 0 {
			sortFieldNames = DefaultDisplayFieldNamesForSorts
		}
	}

	//checkedFieldNames = append(checkedFieldNames, mustCheckFieldNames...)
	//checkedFieldNames = lib.RemoveDuplicates(checkedFieldNames)

	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
	if err != nil {
		return nil, err
	}
	var fabFields []FabField
	for _, v := range fieldStruct.Records {

		if v.IsNoDisplayColumnsForUser() {
			continue
		}

		fieldPermissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
		if err != nil {
			c.log.Error(err)
			continue
		}
		if fieldPermissionVo.CanShow() {
			var tUser *TData
			if userFacade != nil {
				tUser = &userFacade.TData
			}
			a := v.FieldToApi(c.FieldOptionUsecase, c.log, tUser)
			if lib.InArray(v.FieldName, checkedFieldNames) {
				a.Checked = true
			}
			if fieldPermissionVo.CanWrite() {
				a.CanWrite = true
			}
			fabFields = append(fabFields, a)
		}
	}

	// 把字段排序
	fabFields = SortFabField(fabFields, sortFieldNames)

	return fabFields, nil
}

func SortFabField(fabFields []FabField, sortFieldNames []string) []FabField {

	if len(sortFieldNames) == 0 {
		return fabFields
	}
	orderMap := make(map[string]int)
	for k, v := range sortFieldNames {
		orderMap[v] = k
	}

	var fabFields01, fabFields02 []FabField
	for k, v := range fabFields {
		if _, ok := orderMap[v.FieldName]; ok {
			fabFields01 = append(fabFields01, fabFields[k])
		} else {
			fabFields02 = append(fabFields02, fabFields[k])
		}
	}
	//lib.DPrintln(fabFields01)
	//lib.DPrintln(fabFields02)

	// 使用 sort.Slice 按照 order 数组中的顺序排序
	sort.Slice(fabFields01, func(i, j int) bool {
		// 根据 orderMap 中记录的索引进行比较
		return orderMap[fabFields01[i].FieldName] < orderMap[fabFields01[j].FieldName]
	})
	//lib.DPrintln(fabFields01)
	//lib.DPrintln("___=== orderMap:", orderMap)
	//lib.DPrintln("fabFields01:", fabFields01)
	//return fabFields01
	fabFields01 = append(fabFields01, fabFields02...)
	return fabFields01
}

func (c *FieldbuzUsecase) FabFieldsForSearchView(kind string, userFacade *UserFacade) ([]FabField, error) {

	//lib.DPrintln("FabFieldsForSearchView kind: ", kind)

	//configDestFieldNames := []string{
	//	"amount",
	//	"updated_at",
	//	"claims_next_round",
	//	"deal_name",
	//	"active_duty",
	//	"effective_current_rating",
	//	"dob",
	//	"user_gid",
	//}
	//var destFieldNames []string

	var fieldPermissionCenter *FieldPermissionCenter

	if userFacade != nil {
		// 权限过虑
		fieldPermissionCenterR, err := c.FieldPermissionUsecase.CacheFieldPermissionCenter(kind, userFacade.ProfileGid())
		if err != nil {
			return nil, err
		}
		fieldPermissionCenter = &fieldPermissionCenterR

		//for _, v := range configDestFieldNames {
		//	permissionVo, err := fieldPermissionCenter.PermissionByFieldName(v)
		//	if err != nil {
		//		return nil, err
		//	}
		//	if permissionVo.CanShow() {
		//		destFieldNames = append(destFieldNames, v)
		//	}
		//}

	}

	fieldStruct, err := c.FieldUsecase.StructByKind(kind)
	if err != nil {
		return nil, err
	}
	var fabFields []FabField
	for _, v := range fieldStruct.Records {
		if v.IsNoDisplayForUser() {
			continue
		}

		isOk := true
		if fieldPermissionCenter != nil {
			permissionVo, err := fieldPermissionCenter.PermissionByFieldName(v.FieldName)
			if err != nil {
				return nil, err
			}
			if !permissionVo.CanShow() {
				isOk = false
			}
		}

		if isOk {
			var tUser *TData
			if userFacade != nil {
				tUser = &userFacade.TData
			}
			a := v.FieldToApi(c.FieldOptionUsecase, c.log, tUser)
			fabFields = append(fabFields, a)
		}
	}
	return fabFields, nil
}
