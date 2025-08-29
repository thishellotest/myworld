package biz

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type MetadataHttpUsecase struct {
	log                      *log.Helper
	CommonUsecase            *CommonUsecase
	conf                     *conf.Data
	FieldUsecase             *FieldUsecase
	FieldOptionUsecase       *FieldOptionUsecase
	FieldbuzUsecase          *FieldbuzUsecase
	TUsecase                 *TUsecase
	JWTUsecase               *JWTUsecase
	ConditionCategoryUsecase *ConditionCategoryUsecase
	OptionUsecase            *OptionUsecase
	RoleUsecase              *RoleUsecase
	UserUsecase              *UserUsecase
	ClientUsecase            *ClientUsecase
}

func NewMetadataHttpUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
	FieldbuzUsecase *FieldbuzUsecase,
	TUsecase *TUsecase,
	JWTUsecase *JWTUsecase,
	ConditionCategoryUsecase *ConditionCategoryUsecase,
	OptionUsecase *OptionUsecase,
	RoleUsecase *RoleUsecase,
	UserUsecase *UserUsecase,
	ClientUsecase *ClientUsecase) *MetadataHttpUsecase {
	uc := &MetadataHttpUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		FieldUsecase:             FieldUsecase,
		FieldOptionUsecase:       FieldOptionUsecase,
		FieldbuzUsecase:          FieldbuzUsecase,
		TUsecase:                 TUsecase,
		JWTUsecase:               JWTUsecase,
		ConditionCategoryUsecase: ConditionCategoryUsecase,
		OptionUsecase:            OptionUsecase,
		RoleUsecase:              RoleUsecase,
		UserUsecase:              UserUsecase,
		ClientUsecase:            ClientUsecase,
	}

	return uc
}

func (c *MetadataHttpUsecase) Fields(ctx *gin.Context) {
	reply := CreateReply()
	moduleName := ctx.Param("module_name")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizFields(ModuleConvertToKind(moduleName), userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *MetadataHttpUsecase) BizFields(kind string, userFacade UserFacade) (lib.TypeMap, error) {

	fabFields, err := c.FieldbuzUsecase.FabFieldsForSearchView(kind, &userFacade)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set("fields", fabFields)

	return data, nil
}

func (c *MetadataHttpUsecase) Options(ctx *gin.Context) {
	reply := CreateReply()
	moduleName := ctx.Param("module_name")
	fieldName := ctx.Param("field_name")
	rawData, _ := ctx.GetRawData()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizOptions(userFacade, ModuleConvertToKind(moduleName), fieldName, rawData, ctx)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type OptionsRequestVo struct {
	Keyword string `json:"keyword"`
}

func (c *MetadataHttpUsecase) BizOptions(userFacade UserFacade, kind string, fieldName string, rawData []byte, ctx *gin.Context) (lib.TypeMap, error) {

	data := make(lib.TypeMap)
	var optionsRequestVo OptionsRequestVo
	json.Unmarshal(rawData, &optionsRequestVo)

	var fabFieldOptions []FabFieldOption
	if kind == Kind_common {
		if fieldName == CommonFieldName_common_jotform_ids {
			valueType := ctx.Query("value_type")
			caseGid := ctx.Query("case_gid")
			tCase, err := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
			if err != nil {
				return nil, err
			}
			fabFieldOptions, err = c.OptionUsecase.JotformIdsOptions(tCase, valueType, optionsRequestVo.Keyword)
			if err != nil {
				return nil, err
			}
		}
	} else if kind == Kind_Custom_Condition {
		if fieldName == ConditionFieldName_condition_category_id {
			var conds []Cond
			if optionsRequestVo.Keyword != "" {
				conds = append(conds, Like{"category_name", optionsRequestVo.Keyword})
			}
			res, err := c.ConditionCategoryUsecase.ListByCondWithPaging(And(conds...), "id desc", 1, 50)
			if err != nil {
				return nil, err
			}
			for _, v := range res {
				fabFieldOptions = append(fabFieldOptions, v.ToFabFieldOption())
			}
			fabFieldOptions = append(fabFieldOptions, FabFieldOption{
				OptionLabel: UndefinedConditionCategory.ConditionCategoryName,
				OptionValue: InterfaceToString(UndefinedConditionCategory.ConditionCategoryId),
			})
		}

	} else {
		fieldStruct, err := c.FieldUsecase.StructByKind(kind)
		if err != nil {
			return nil, err
		}
		if fieldStruct == nil {
			return nil, errors.New("fieldStruct is nil")
		}

		fieldEntity := fieldStruct.GetByFieldName(fieldName)
		if fieldEntity == nil {
			return nil, errors.New("fieldEntity is nil")
		}

		if fieldEntity.FieldType != FieldType_multilookup && fieldEntity.FieldType != FieldType_lookup && fieldEntity.FieldType != FieldType_dropdown {
			return nil, errors.New("FieldType is wrong")
		}

		if fieldEntity.FieldType == FieldType_lookup || fieldEntity.FieldType == FieldType_multilookup {

			var options []*TData
			useNormal := false
			if kind == Kind_users && fieldEntity.FieldName == UserFieldName_role_gid {
				userProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)
				if userProfile == nil {
					return nil, errors.New("Profile is nil")
				}
				if IsAdminProfile(userProfile) {
					useNormal = true
				} else {
					userRole, _ := c.RoleUsecase.GetRole(userFacade.CustomFields.TextValueByNameBasic(UserFieldName_role_gid))
					if userRole == nil {
						return nil, errors.New("userRole is nil")
					}
					options, _ = c.RoleUsecase.ChildrenRoles(userRole)
				}
			} else if kind == Kind_users && fieldEntity.FieldName == User_FieldName_profile_gid {
				userProfile, _ := c.UserUsecase.GetProfile(&userFacade.TData)
				if userProfile == nil {
					return nil, errors.New("Profile is nil")
				}
				if IsAdminProfile(userProfile) {
					useNormal = true
				} else {
					options, _ = c.TUsecase.ListByCond(Kind_profiles, Eq{"is_admin": 0, "deleted_at": 0})
				}
			} else {
				useNormal = true
			}

			if useNormal {
				var conds []Cond
				if optionsRequestVo.Keyword != "" {
					conds = append(conds, Like{fieldEntity.RelaName, optionsRequestVo.Keyword})
				}
				conds = append(conds, Eq{FieldName_biz_deleted_at: 0})
				options, err = c.TUsecase.ListByCond(fieldEntity.RelaKind, And(conds...))
			}

			if err != nil {
				return nil, err
			}

			var clientPipelines ClientPipelines
			if kind == Kind_client_cases && fieldName == FieldName_client_gid {
				var clientGids []string
				for _, v := range options {
					clientGids = append(clientGids, v.Gid())
				}
				clientPipelines, _ = c.ClientUsecase.GetClientsPipelines(clientGids)
			}
			for _, v := range options {
				fabFieldOptions = append(fabFieldOptions, FabFieldOption{
					OptionLabel: v.CustomFields.TextValueByNameBasic(fieldEntity.RelaName),
					OptionValue: v.Gid(),
					Pipelines:   clientPipelines.GetByClientGid(v.Gid()),
				})
			}
		} else if fieldEntity.FieldType == FieldType_dropdown || fieldEntity.FieldType == FieldType_multidropdown {
			structOptions, err := c.FieldOptionUsecase.StructByKind(kind)
			if err != nil {
				return nil, err
			}
			options := structOptions.AllByFieldName(*fieldEntity)
			for _, v := range options {
				fabFieldOptions = append(fabFieldOptions, FabFieldOption{
					OptionLabel: v.OptionLabel,
					OptionValue: v.OptionValue,
				})
			}
		}
	}
	//time.Sleep(5 * time.Second)
	data.Set("options", fabFieldOptions)

	return data, nil
}
