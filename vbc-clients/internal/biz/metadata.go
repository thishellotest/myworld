package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type MetadataUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	ConditionUsecase *ConditionUsecase
	JWTUsecase       *JWTUsecase
	TUsecase         *TUsecase
	FieldbuzUsecase  *FieldbuzUsecase
}

func NewMetadataUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ConditionUsecase *ConditionUsecase,
	JWTUsecase *JWTUsecase,
	TUsecase *TUsecase,
	FieldbuzUsecase *FieldbuzUsecase) *MetadataUsecase {
	uc := &MetadataUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		ConditionUsecase: ConditionUsecase,
		JWTUsecase:       JWTUsecase,
		TUsecase:         TUsecase,
		FieldbuzUsecase:  FieldbuzUsecase,
	}

	return uc
}

func (c *MetadataUsecase) Basicdata(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizBasicdata(userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type MetadataHttpBasicdataResponse struct {
	Users      []FabUser  `json:"users"`
	CaseFields []FabField `json:"case_fields"`
}

func (c *MetadataUsecase) BizBasicdata(userFacade UserFacade) (lib.TypeMap, error) {

	data := make(lib.TypeMap)
	var metadataHttpBasicdataResponse MetadataHttpBasicdataResponse

	users, err := c.TUsecase.ListByCond(Kind_users, Eq{DataEntry_biz_deleted_at: 0})
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		metadataHttpBasicdataResponse.Users = append(metadataHttpBasicdataResponse.Users, FabUser{
			Gid:      v.Gid(),
			FullName: v.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
			Email:    v.CustomFields.TextValueByNameBasic(UserFieldName_email),
		})
	}

	caseFields, err := c.FieldbuzUsecase.FabFieldsForBasicdata(Kind_client_cases, &userFacade)
	if err != nil {
		return nil, err
	}
	metadataHttpBasicdataResponse.CaseFields = caseFields
	data.Set("data", metadataHttpBasicdataResponse)

	return data, nil
}

func (c *MetadataUsecase) HttpConditions(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpConditions()
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

type ClaimCondition struct {
	Condition string `json:"condition"`
}

func (c *MetadataUsecase) BizHttpConditions() (lib.TypeMap, error) {
	data := make(lib.TypeMap)

	var conditions []ClaimCondition

	records, err := c.ConditionUsecase.AllByCondWithOrderBy(Eq{"deleted_at": 0}, "condition_name", 50)
	if err != nil {
		return nil, err
	}
	for _, v := range records {
		conditions = append(conditions, ClaimCondition{
			Condition: v.ConditionName,
		})
	}

	data.Set("conditions", conditions)
	return data, nil
}
