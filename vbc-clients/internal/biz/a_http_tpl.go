package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
)

type HttpTpl struct {
	log        *log.Helper
	conf       *conf.Data
	JWTUsecase *JWTUsecase
}

func NewHttpTpl(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase) *HttpTpl {
	return &HttpTpl{
		log:        log.NewHelper(logger),
		conf:       conf,
		JWTUsecase: JWTUsecase,
	}
}

func (c *HttpTpl) HttpFun(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))

	// 通过路由获取的
	//moduleName := ctx.Param("module_name")

	// tUser, _ := c.JWTUsecase.JWTUser(ctx)

	data, err := c.BizHttpFun(body.GetString("token"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpTpl) BizHttpFun(str string) (lib.TypeMap, error) {

	data := make(lib.TypeMap)
	data.Set("data.val", "aaa")
	return data, nil
}
