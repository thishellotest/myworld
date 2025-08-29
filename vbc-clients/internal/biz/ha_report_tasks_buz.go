package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type HaReportTasksBuzUsecase struct {
	log                 *log.Helper
	CommonUsecase       *CommonUsecase
	conf                *conf.Data
	BoxUsecase          *BoxUsecase
	MapUsecase          *MapUsecase
	TUsecase            *TUsecase
	BlobUsecase         *BlobUsecase
	BlobbuzUsecase      *BlobbuzUsecase
	HaReportTaskUsecase *HaReportTaskUsecase
}

func NewHaReportTasksBuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxUsecase *BoxUsecase,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	BlobUsecase *BlobUsecase,
	BlobbuzUsecase *BlobbuzUsecase,
	HaReportTaskUsecase *HaReportTaskUsecase) *HaReportTasksBuzUsecase {
	uc := &HaReportTasksBuzUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		BoxUsecase:          BoxUsecase,
		MapUsecase:          MapUsecase,
		TUsecase:            TUsecase,
		BlobUsecase:         BlobUsecase,
		BlobbuzUsecase:      BlobbuzUsecase,
		HaReportTaskUsecase: HaReportTaskUsecase,
	}

	return uc
}

func (c *HaReportTasksBuzUsecase) CreateTask(ctx context.Context, blobGid string) error {

	tBlob, err := c.BlobUsecase.GetByGid(blobGid)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tBlob == nil {
		return errors.New("tBlob is nil")
	}

	// 判断任务是否存在，存在不重复创建任务
	exists, err := c.HaReportTaskUsecase.GetByCond(Eq{"blob_gid": tBlob.CustomFields.TextValueByNameBasic("gid"),
		"handle_status": HandleStatus_waiting,
		"deleted_at":    0})
	if err != nil {
		return err
	}
	if exists != nil {
		return errors.New("The task has been created")
	}

	task := &HaReportTaskEntity{
		Gid:     uuid.UuidWithoutStrike(),
		BlobGid: tBlob.CustomFields.TextValueByNameBasic(FieldName_gid),
	}
	return c.CommonUsecase.DB().Create(&task).Error
}

func (c *HaReportTasksBuzUsecase) deprecated_CreateTask(ctx context.Context, BoxFileId string) error {

	fileInfo, _, err := c.BoxUsecase.GetFileInfoForTypeMap(BoxFileId)
	lib.DPrintln(fileInfo)
	if err != nil {
		return err
	}

	tCase, err := c.BlobbuzUsecase.GetCaseByBoxFileInfo(fileInfo)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	lib.DPrintln(tCase, err)

	fileVersionID := fileInfo.GetString("file_version.id")
	fileName := fileInfo.GetString("name")

	_, suffix := lib.FileExt(fileName, true)
	if suffix != BlobType_pdf {
		return errors.New("The file format is not supported")
	}

	uniqblob := GenUniqblob(BoxFileId, fileVersionID)
	tBlob, err := c.TUsecase.Data(Kind_blobs, Eq{"uniqblob": uniqblob, "deleted_at": 0})
	if err != nil {
		return err
	}
	if tBlob == nil {
		tBlob, err = c.BlobbuzUsecase.HandleBoxFile(ctx,
			fileInfo,
			tCase.CustomFields.TextValueByNameBasic("gid"), "")
		if err != nil {
			return err
		}
	}
	if tBlob == nil {
		return errors.New("tBlob is nil")
	}

	// 判断任务是否存在，存在不重复创建任务
	exists, err := c.HaReportTaskUsecase.GetByCond(Eq{"blob_gid": tBlob.CustomFields.TextValueByNameBasic("gid"),
		"handle_status": HandleStatus_waiting,
		"deleted_at":    0})
	if err != nil {
		return err
	}
	if exists != nil {
		return errors.New("The task has been created")
	}

	task := &HaReportTaskEntity{
		Gid:     uuid.UuidWithoutStrike(),
		BlobGid: tBlob.CustomFields.TextValueByNameBasic("gid"),
	}
	return c.CommonUsecase.DB().Create(&task).Error
	//lib.DPrintln(task)
	//return nil
}

func (c *HaReportTasksBuzUsecase) deprecated_HttpCreateTask(ctx *gin.Context) {

	reply := CreateReply()
	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	err := c.deprecated_CreateTask(ctx, body.GetString("file_id"))
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Success()
	}
	ctx.JSON(200, reply)
}
