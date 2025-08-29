package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

type AccountHttpUsecase struct {
	log              *log.Helper
	conf             *conf.Data
	JWTUsecase       *JWTUsecase
	FilterbuzUsecase *FilterbuzUsecase
}

func NewAccountHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	FilterbuzUsecase *FilterbuzUsecase) *AccountHttpUsecase {
	return &AccountHttpUsecase{
		log:              log.NewHelper(logger),
		conf:             conf,
		JWTUsecase:       JWTUsecase,
		FilterbuzUsecase: FilterbuzUsecase,
	}
}

func (c *AccountHttpUsecase) FilterSave(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	moduleName := ctx.Param("module_name")
	content := body.GetString("content")
	filterName := body.GetString("filter_name")
	data, err := c.FilterbuzUsecase.BizFilterSave(userFacade.Gid(), ModuleConvertToKind(moduleName), content, filterName)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AccountHttpUsecase) FilterList(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	moduleName := ctx.Param("module_name")
	data, err := c.FilterbuzUsecase.BizFilterList(userFacade.Gid(), ModuleConvertToKind(moduleName), "")
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AccountHttpUsecase) FilterDelete(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	// 通过路由获取的
	filterId := ctx.Param("filter_id")
	data, err := c.BizFilterDelete(userFacade, lib.StringToInt32(filterId))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *AccountHttpUsecase) BizFilterDelete(userFacade UserFacade, filterId int32) (lib.TypeMap, error) {
	err := c.FilterbuzUsecase.FilterDelete(userFacade.Gid(), []int32{filterId})
	return nil, err
}
