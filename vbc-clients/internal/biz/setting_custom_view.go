package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
)

type SettingCustomViewUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	FieldbuzUsecase  *FieldbuzUsecase
	LongMapUsecase   *LongMapUsecase
	KindUsecase      *KindUsecase
	FilterbuzUsecase *FilterbuzUsecase
}

func NewSettingCustomViewUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldbuzUsecase *FieldbuzUsecase,
	LongMapUsecase *LongMapUsecase,
	KindUsecase *KindUsecase,
	FilterbuzUsecase *FilterbuzUsecase) *SettingCustomViewUsecase {
	uc := &SettingCustomViewUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		FieldbuzUsecase:  FieldbuzUsecase,
		LongMapUsecase:   LongMapUsecase,
		KindUsecase:      KindUsecase,
		FilterbuzUsecase: FilterbuzUsecase,
	}
	return uc
}

func (c *SettingCustomViewUsecase) Get(kind string, userFacade UserFacade, tableType string) (fabSettingCustomView FabSettingCustomView, err error) {

	var fabCustomView FabCustomView
	if kind != Kind_users {
		fabCustomView, err = c.FieldbuzUsecase.FabCustomView(kind, &userFacade, tableType)
		if err != nil {
			c.log.Error(err)
			//return fabSettingCustomView, err
		}
	}
	fabSettingCustomView.CustomView = fabCustomView

	fabFields, err := c.FieldbuzUsecase.FabFields(kind, &userFacade, tableType)
	if err != nil {
		return fabSettingCustomView, err
	}
	fabSettingCustomView.Fields = fabFields

	var columnwidths Columnwidths
	if kind == Kind_users {
		columnwidths, _, _ = c.FieldbuzUsecase.UsersCheckedAndSortFieldNames()
	} else {
		columnwidths, err = c.FieldbuzUsecase.FabColumnwidth(kind, &userFacade, tableType)
		if err != nil {
			return fabSettingCustomView, err
		}
	}
	fabSettingCustomView.Columns = columnwidths

	searchFabFields, err := c.FieldbuzUsecase.FabFieldsForSearchView(kind, &userFacade)
	if err != nil {
		return fabSettingCustomView, err
	}
	fabSettingCustomView.SearchFields = searchFabFields

	kindEntity, err := c.KindUsecase.GetByKind(kind)
	if err != nil {
		return fabSettingCustomView, err
	}
	if kindEntity == nil {
		return fabSettingCustomView, errors.New("kindEntity is nil")
	}

	filters, err := c.FilterbuzUsecase.FilterList(userFacade.Gid(), kind, tableType)
	if err != nil {
		return fabSettingCustomView, err
	}
	fabSettingCustomView.Filters = filters

	fabSettingCustomView.Module = FabModule{
		ModuleName: KindConvertToModule(kind),
		Label:      kindEntity.Label,
		TabLabel:   kindEntity.TabLabel,
	}

	return fabSettingCustomView, nil
}

func (c *SettingCustomViewUsecase) ChangeSort(kind string, userFacade UserFacade, fabCustomView FabCustomView) error {
	key := MapKeyCustomView(userFacade.Gid(), kind, fabCustomView.TableType)
	c.LongMapUsecase.Set(key, InterfaceToString(fabCustomView))
	return nil
}

func (c *SettingCustomViewUsecase) ChangeFields(kind string, userFacade UserFacade, changeFieldsVo ChangeFieldsVo) error {
	key := MapKeyCustomViewColumns(userFacade.Gid(), kind, changeFieldsVo.TableType)
	c.LongMapUsecase.Set(key, InterfaceToString(changeFieldsVo))
	return nil
}

type ColumnwidthVo struct {
	Columns   Columnwidths `json:"columns"`
	TableType string       `json:"table_type"`
}

type Columnwidths map[string]ColumnwidthUnitVo

type ColumnwidthUnitVo struct {
	Width float32 `json:"width"`
}

func (c *SettingCustomViewUsecase) ChangeColumnwidth(kind string, userFacade UserFacade, columnwidthVo ColumnwidthVo) error {
	key := MapKeyCustomViewColumnwidth(userFacade.Gid(), kind, columnwidthVo.TableType)
	c.LongMapUsecase.Set(key, InterfaceToString(columnwidthVo))
	return nil
}
