package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

type MgmtHttpUsecase struct {
	log         *log.Helper
	conf        *conf.Data
	JWTUsecase  *JWTUsecase
	MenuUsecase *MenuUsecase
}

func NewMgmtHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	MenuUsecase *MenuUsecase) *MgmtHttpUsecase {
	return &MgmtHttpUsecase{
		log:         log.NewHelper(logger),
		conf:        conf,
		JWTUsecase:  JWTUsecase,
		MenuUsecase: MenuUsecase,
	}
}

func (c *MgmtHttpUsecase) Init(ctx *gin.Context) {
	reply := CreateReply()
	//rawData, _ := ctx.GetRawData()
	//body := lib.ToTypeMapByString(string(rawData))
	// 通过路由获取的
	//moduleName := ctx.Param("module_name")
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizInit(userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *MgmtHttpUsecase) BizInit(userFacade UserFacade) (lib.TypeMap, error) {
	menu, err := c.MenuUsecase.GetMenu(userFacade)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("menu", menu)
	return data, nil
}
