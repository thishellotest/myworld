package biz

import (
	"errors"
	"github.com/gin-gonic/gin"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HttpAccessControl struct {
	AccessControlWorkUsecase            *AccessControlWorkUsecase
	AccessControlWorkPayloadTaskUsecase *AccessControlWorkPayloadTaskUsecase
}

func NewHttpAccessControl(AccessControlWorkUsecase *AccessControlWorkUsecase,
	AccessControlWorkPayloadTaskUsecase *AccessControlWorkPayloadTaskUsecase) *HttpAccessControl {
	return &HttpAccessControl{
		AccessControlWorkUsecase:            AccessControlWorkUsecase,
		AccessControlWorkPayloadTaskUsecase: AccessControlWorkPayloadTaskUsecase,
	}
}

func (c *HttpAccessControl) Tasks(ctx *gin.Context) {
	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	data, err := c.HandleTasks(body.GetString("token"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply["data"] = data
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpAccessControl) HandleTasks(token string) (interface{}, error) {
	if len(token) == 0 {
		return nil, errors.New("token is empty.")
	}
	work, err := c.AccessControlWorkUsecase.GetByCond(Eq{"deleted_at": 0, "token": token})
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("work is nil.")
	}

	return c.AccessControlWorkUsecase.AccessControlWorkPayload(work)
}

func (c *HttpAccessControl) CarryOut(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	err := c.HandeCarryOut(body.GetString("token"), body.GetInt("index"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *HttpAccessControl) HandeCarryOut(token string, index int32) error {

	work, err := c.AccessControlWorkUsecase.GetByCond(Eq{"deleted_at": 0, "token": token})
	if err != nil {
		return err
	}
	if work == nil {
		return errors.New("work is nil.")
	}
	err = c.AccessControlWorkUsecase.VerifyAccess(work)
	if err != nil {
		return err
	}

	payload, err := work.GetPayload()
	if err != nil {
		return err
	}

	taskPayload := payload.GetByIndex(index)
	if taskPayload == nil {
		return errors.New("taskPayload is nil.")
	}
	return c.AccessControlWorkPayloadTaskUsecase.Handle(work, taskPayload)
}
