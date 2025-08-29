package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	. "vbc/lib/builder"
)

type HttpManualUsecase struct {
	CommonUsecase              *CommonUsecase
	TUsecase                   *TUsecase
	ClientUsecase              *ClientUsecase
	log                        *log.Helper
	UniqueCodeGeneratorUsecase *UniqueCodeGeneratorUsecase
	SyncAsanaTaskUsecase       *SyncAsanaTaskUsecase
	//GoogleSheetSyncTaskUsecase *GoogleSheetSyncTaskUsecase
}

func NewHttpManualUsecase(CommonUsecase *CommonUsecase,
	TUsecase *TUsecase, ClientUsecase *ClientUsecase, logger log.Logger,
	UniqueCodeGeneratorUsecase *UniqueCodeGeneratorUsecase,
	SyncAsanaTaskUsecase *SyncAsanaTaskUsecase,
// GoogleSheetSyncTaskUsecase *GoogleSheetSyncTaskUsecase,
) *HttpManualUsecase {
	return &HttpManualUsecase{
		CommonUsecase:              CommonUsecase,
		TUsecase:                   TUsecase,
		ClientUsecase:              ClientUsecase,
		log:                        log.NewHelper(logger),
		UniqueCodeGeneratorUsecase: UniqueCodeGeneratorUsecase,
		SyncAsanaTaskUsecase:       SyncAsanaTaskUsecase,
		//GoogleSheetSyncTaskUsecase: GoogleSheetSyncTaskUsecase,
	}
}

func (c *HttpManualUsecase) ReplenishClientUniqcode(ctx *gin.Context) {

	res, err := c.ClientUsecase.AllByCond(Eq{FieldName_uniqcode: ""})
	if err != nil {
		c.log.Error(err)
	}
	for k, _ := range res {
		Uniqcode, _ := c.UniqueCodeGeneratorUsecase.GenUuid(UniqueCodeGenerator_Type_ClientUniqCode, 0)
		res[k].Uniqcode = Uniqcode
		err = c.CommonUsecase.DB().Save(res[k]).Error
		if err != nil {
			c.log.Error(err)
		}
	}
	ctx.JSON(0, "ok")
}

func (c *HttpManualUsecase) SyncAllClientFromAsana(ctx *gin.Context) {

	taskGid := ctx.Query("gid")

	err := c.HandleSyncAllClientFromAsana(taskGid)
	msg := "ok"
	if err != nil {
		msg = err.Error()
	}
	ctx.JSON(0, msg)
}

func (c *HttpManualUsecase) HandleSyncAllClientFromAsana(taskGid string) error {

	if taskGid != "" {
		err := c.SyncAsanaTaskUsecase.LPushSyncTaskQueue(context.TODO(), taskGid)
		if err != nil {
			c.log.Error(err)
		}
		return nil
	}
	res, err := c.TUsecase.ListByCond(Kind_client_cases, Neq{"assignee_gid": ""})
	if err != nil {
		return err
	}
	for _, v := range res {
		gid := v.CustomFields.TextValueByNameBasic("asana_task_gid")
		err = c.SyncAsanaTaskUsecase.LPushSyncTaskQueue(context.TODO(), gid)
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}

func (c *HttpManualUsecase) SyncDashboardGoogleSheet(ctx *gin.Context) {

	err := c.HandleSyncDashboardGoogleSheet()
	msg := "ok"
	if err != nil {
		msg = err.Error()
	}
	ctx.JSON(0, msg)
}

func (c *HttpManualUsecase) HandleSyncDashboardGoogleSheet() error {
	//return c.GoogleSheetSyncTaskUsecase.LPushSyncTaskQueue(context.TODO(), "1")
	return nil
}
