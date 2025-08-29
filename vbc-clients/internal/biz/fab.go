package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

const Fab_TData = "data"
const Fab_TList = "list"
const Fab_TRecords = "records"
const Fab_TTotal = "total"
const Fab_TPage = "page"
const Fab_TPageSize = "page_size"
const Fab_HasMore = "has_more"

const (
	Fab_User = "user"
	Fab_Case = "case"
	Fab_Blob = "blob"
)

type FabUser struct {
	Gid      string `json:"gid"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type FabModule struct {
	ModuleName string `json:"module_name"`
	TabLabel   string `json:"tab_label"`
	Label      string `json:"label"`
}

type FabField struct {
	Module          string           `json:"module"`
	FieldName       string           `json:"field_name"`
	FieldLabel      string           `json:"field_label"`
	FieldType       string           `json:"field_type"`
	RelaModule      string           `json:"rela_module"`
	RelaName        string           `json:"rela_name"`
	Options         []FabFieldOption `json:"options"`
	Checked         bool             `json:"checked"` // 是否选中
	CanWrite        bool             `json:"can_write"`
	CanManageColumn bool             `json:"can_show_manage_column"` // 用户能否在manage column管理
	IsRequired      bool             `json:"is_required"`
	IsEnableColor   bool             `json:"is_enable_color"`
	Tooltip         string           `json:"tooltip"`
}

type FabFieldOption struct {
	OptionLabel string `json:"option_label"`
	OptionValue string `json:"option_value"`
	OptionColor string `json:"option_color"`
	Pipelines   string `json:"pipelines"`
}

type FabFieldOptionNew struct { // 与前端格式保持一致，后续都使用New
	OptionLabel string `json:"label"`
	OptionValue string `json:"value"`
	OptionColor string `json:"color"`
}

type FabSettingCustomView struct {
	CustomView   FabCustomView `json:"custom_view"`
	Fields       []FabField    `json:"fields"`
	SearchFields []FabField    `json:"search_fields"`
	Columns      Columnwidths  `json:"columns"`
	Module       FabModule     `json:"module"`
	Filters      []FilterVo    `json:"filters"`
}

type FabCustomView struct {
	TableType string                `json:"table_type"`
	SortBys   []FabCustomViewSortBy `json:"sort_bys"`
}

func (c *FabCustomView) IsCustomSortBy() bool {
	if len(c.SortBys) > 0 {
		return true
	}
	return false
}

type FabCustomViewSortBy struct {
	FieldName string `json:"field_name"`
	SortOrder string `json:"sort_order"`
}

const Fab_SortOrder_asc = "asc"
const Fab_SortOrder_desc = "desc"

func GetFabByKind(kind string) (string, error) {
	if kind == Kind_users {
		return Fab_User, nil
	} else if kind == Kind_blobs {
		return Fab_Blob, nil
	} else if kind == Kind_client_cases {
		return Fab_Case, nil
	}
	return "", errors.New("GetFabByKind is wrong")
}

type FabUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	FieldUsecase       *FieldUsecase
	FieldOptionUsecase *FieldOptionUsecase
}

func NewFabUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase) *FabUsecase {
	uc := &FabUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		FieldUsecase:       FieldUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
	}
	return uc
}

//
//func (c *FabUsecase) ToApiMap(kind string, tData *TData) (lib.TypeMap, error) {
//	if tData == nil {
//		return nil, errors.New("ToApiMap: tData is nil")
//	}
//
//	fieldStruct, err := c.FieldUsecase.CacheStructByKind(kind)
//	if err != nil {
//		return nil, err
//	}
//	optionStruct, err:=c.FieldOptionUsecase.CacheStructByKind(kind)
//	if err!=nil{
//		return nil, err
//	}
//	for k, v := range tData.CustomFields {
//		fieldEntity := fieldStruct.GetByFieldName(v.Name)
//		if fieldEntity == nil {
//			return nil, errors.New("ToApiMap: fieldEntity is nil : " + v.Name)
//		}
//		optionStruct.AllByFieldName(v.Name).GetByValue().
//
//	}
//
//	return tData.CustomFields.ToApiMap()
//}
