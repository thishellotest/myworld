package biz

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HttpTUsecase struct {
	log      *log.Helper
	TUsecase *TUsecase
}

func NewHttpTUsecase(logger log.Logger, TUsecase *TUsecase) *HttpTUsecase {

	return &HttpTUsecase{
		log:      log.NewHelper(logger),
		TUsecase: TUsecase,
	}
}

func (c *HttpTUsecase) List(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	data, err := c.BizList("", rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpTUsecase) BizList(kind string, rawData []byte) (lib.TypeMap, error) {

	tList, total, page, pageSize, err := c.DoBizList(kind, rawData)
	if err != nil {
		return nil, err
	}

	data := make(lib.TypeMap)
	data.Set(Fab_TData+"."+Fab_TList, tList)
	data.Set(Fab_TData+"."+Fab_TTotal, int32(total))
	data.Set(Fab_TData+"."+Fab_TPage, page)
	data.Set(Fab_TData+"."+Fab_TPageSize, pageSize)

	return data, nil
}

func (c *HttpTUsecase) DoBizList(kind string, rawData []byte) (tList TDataList, total int64, page int, pageSize int, err error) {

	var request map[string]interface{}
	tListRequest := &TListRequest{}
	if len(rawData) > 0 {
		err = json.Unmarshal(rawData, &request)
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}
	//tListRequest.Page = int(lib.GetFromMapToInt64(request, Fab_TPage))
	//tListRequest.PageSize = int(lib.GetFromMapToInt64(request, Fab_TPageSize))

	tList, err = c.TUsecase.List(kind, nil, tListRequest, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	total, err = c.TUsecase.Total(kind)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return tList, total, page, pageSize, nil
}

func (c *HttpTUsecase) Detail(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	data, err := c.BizList("", rawData)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpTUsecase) BizDetail(kind string, cond Cond) (*TData, error) {

	tData, err := c.TUsecase.Data(kind, cond)
	if err != nil {
		return nil, err
	}
	if tData == nil {
		return nil, nil
	}
	return tData, nil
}
