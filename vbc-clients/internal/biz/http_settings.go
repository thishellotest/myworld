package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

type HttpSettingsUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	JWTUsecase         *JWTUsecase
	UserUsecase        *UserUsecase
	AppTokenUsecase    *AppTokenUsecase
	FieldUsecase       *FieldUsecase
	FieldOptionUsecase *FieldOptionUsecase
}

func NewHttpSettingsUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	UserUsecase *UserUsecase,
	AppTokenUsecase *AppTokenUsecase,
	FieldUsecase *FieldUsecase,
	FieldOptionUsecase *FieldOptionUsecase,
) *HttpSettingsUsecase {
	uc := &HttpSettingsUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		JWTUsecase:         JWTUsecase,
		UserUsecase:        UserUsecase,
		AppTokenUsecase:    AppTokenUsecase,
		FieldUsecase:       FieldUsecase,
		FieldOptionUsecase: FieldOptionUsecase,
	}

	return uc
}

func (c *HttpSettingsUsecase) HttpFields(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.BizHttpFields(body.GetString("module"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpSettingsUsecase) BizHttpFields(module string) (lib.TypeMap, error) {

	//fieldList, err := c.FieldUsecase.ListByKind(module)
	//if err != nil {
	//	c.log.Error(err)
	//	return nil, err
	//}
	//var fields lib.TypeList
	//for _, v := range fieldList {
	//	fields = append(fields, v.FieldToApi(c.FieldOptionUsecase, c.log))
	//}
	//data := make(lib.TypeMap)
	//data.Set("fields", fields)

	return nil, nil
}
