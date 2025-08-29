package biz

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

type SettingHttpUsecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	JWTUsecase               *JWTUsecase
	SettingCustomViewUsecase *SettingCustomViewUsecase
	FieldbuzUsecase          *FieldbuzUsecase
}

func NewSettingHttpUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	SettingCustomViewUsecase *SettingCustomViewUsecase,
	FieldbuzUsecase *FieldbuzUsecase) *SettingHttpUsecase {
	uc := &SettingHttpUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		JWTUsecase:               JWTUsecase,
		SettingCustomViewUsecase: SettingCustomViewUsecase,
		FieldbuzUsecase:          FieldbuzUsecase,
	}

	return uc
}

const (
	SettingHttpCustomViewRequest_TableType_upcoming = "upcoming"
	SettingHttpCustomViewRequest_TableType_overdue  = "overdue"
	SettingHttpCustomViewRequest_TableType_ongoing  = "ongoing"
	SettingHttpCustomViewRequest_TableType_Default  = ""
	SettingHttpCustomViewRequest_TableType_Search   = "search"
)

type SettingHttpCustomViewRequest struct {
	TableType string `json:"table_type"`
}

func (c *SettingHttpUsecase) CustomView(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	var settingHttpCustomViewRequest SettingHttpCustomViewRequest
	json.Unmarshal(rawData, &settingHttpCustomViewRequest)

	moduleName := ctx.Param("module_name")

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizCustomView(ModuleConvertToKind(moduleName), userFacade, settingHttpCustomViewRequest)

	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *SettingHttpUsecase) BizCustomView(kind string, userFacade UserFacade, settingHttpCustomViewRequest SettingHttpCustomViewRequest) (lib.TypeMap, error) {
	fabCustomView, err := c.SettingCustomViewUsecase.Get(kind, userFacade, settingHttpCustomViewRequest.TableType)
	if err != nil {
		return nil, err
	}
	if kind == Kind_client_cases && userFacade.Gid() == config_vbc.User_Dev_gid { // 这是方便开发帐号测试使用
		for k, v := range fabCustomView.Fields {
			if v.FieldName == FieldName_primary_vs ||
				v.FieldName == FieldName_primary_cp ||
				v.FieldName == FieldName_lead_co ||
				v.FieldName == FieldName_support_cp {
				fabCustomView.Fields[k].Options = append(fabCustomView.Fields[k].Options, FabFieldOption{
					OptionValue: "Yannan Wang",
					OptionLabel: "Yannan Wang",
				})
				fabCustomView.Fields[k].Options = append(fabCustomView.Fields[k].Options, FabFieldOption{
					OptionValue: "Engineering Team",
					OptionLabel: "Engineering Team",
				})
			}
		}

	}

	return lib.ToTypeMapByString(InterfaceToString(fabCustomView)), nil
}

func (c *SettingHttpUsecase) ChangeSort(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	moduleName := ctx.Param("module_name")

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	//lib.DPrintln(tUser)
	data, err := c.BizChangeSort(ModuleConvertToKind(moduleName), userFacade, rawData)

	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *SettingHttpUsecase) BizChangeSort(kind string, userFacade UserFacade, rawData []byte) (lib.TypeMap, error) {

	var fabCustomView FabCustomView
	err := json.Unmarshal(rawData, &fabCustomView)
	if err != nil {
		return nil, err
	}

	err = c.SettingCustomViewUsecase.ChangeSort(kind, userFacade, fabCustomView)
	if err != nil {
		return nil, err
	}

	destFabCustomView, err := c.FieldbuzUsecase.FabCustomView(kind, &userFacade, fabCustomView.TableType)
	if err != nil {
		return nil, err
	}

	typeMaps := make(lib.TypeMap)
	typeMaps.Set("custom_view", destFabCustomView)

	return typeMaps, nil

}

type ChangeFieldsVo struct {
	Fields    []ChangeFieldVo `json:"fields"`
	TableType string          `json:"table_type"`
}

type ChangeFieldVo struct {
	FieldName string `json:"field_name"`
	Checked   bool   `json:"checked"`
}

func (c *SettingHttpUsecase) ChangeFields(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	moduleName := ctx.Param("module_name")

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizChangeFields(ModuleConvertToKind(moduleName), userFacade, rawData)

	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *SettingHttpUsecase) BizChangeFields(kind string, userFacade UserFacade, rawData []byte) (lib.TypeMap, error) {

	var changeFieldsVo ChangeFieldsVo
	err := json.Unmarshal(rawData, &changeFieldsVo)
	if err != nil {
		return nil, err
	}

	err = c.SettingCustomViewUsecase.ChangeFields(kind, userFacade, changeFieldsVo)
	if err != nil {
		return nil, err
	}

	fabFields, err := c.FieldbuzUsecase.FabFields(kind, &userFacade, changeFieldsVo.TableType)
	if err != nil {
		return nil, err
	}
	typeMaps := make(lib.TypeMap)
	typeMaps.Set("fields", fabFields)

	return typeMaps, nil
}

func (c *SettingHttpUsecase) ChangeColumnwidth(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	moduleName := ctx.Param("module_name")

	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizChangeColumnwidth(ModuleConvertToKind(moduleName), userFacade, rawData)

	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *SettingHttpUsecase) BizChangeColumnwidth(kind string, userFacade UserFacade, rawData []byte) (lib.TypeMap, error) {

	var columnwidthVo ColumnwidthVo
	err := json.Unmarshal(rawData, &columnwidthVo)
	if err != nil {
		return nil, err
	}

	err = c.SettingCustomViewUsecase.ChangeColumnwidth(kind, userFacade, columnwidthVo)
	if err != nil {
		return nil, err
	}

	columnwidths, err := c.FieldbuzUsecase.FabColumnwidth(kind, &userFacade, columnwidthVo.TableType)
	if err != nil {
		return nil, err
	}

	typeMaps := make(lib.TypeMap)
	typeMaps.Set("columns", columnwidths)

	return typeMaps, nil
}
