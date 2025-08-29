package biz

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib/uuid"
)

type BlobSliceJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BlobSliceEntity]
	BaseHandleCustom[BlobSliceEntity]
	AzcognitiveUsecase      *AzcognitiveUsecase
	BlobUsecase             *BlobUsecase
	BlobSliceUsecase        *BlobSliceUsecase
	AzstorageUsecase        *AzstorageUsecase
	DataEntryUsecase        *DataEntryUsecase
	HaReportTasksBuzUsecase *HaReportTasksBuzUsecase
}

func NewBlobSliceJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AzcognitiveUsecase *AzcognitiveUsecase,
	BlobUsecase *BlobUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	AzstorageUsecase *AzstorageUsecase,
	DataEntryUsecase *DataEntryUsecase,
	HaReportTasksBuzUsecase *HaReportTasksBuzUsecase) *BlobSliceJobUsecase {
	uc := &BlobSliceJobUsecase{
		log:                     log.NewHelper(logger),
		CommonUsecase:           CommonUsecase,
		conf:                    conf,
		AzcognitiveUsecase:      AzcognitiveUsecase,
		BlobUsecase:             BlobUsecase,
		BlobSliceUsecase:        BlobSliceUsecase,
		AzstorageUsecase:        AzstorageUsecase,
		DataEntryUsecase:        DataEntryUsecase,
		HaReportTasksBuzUsecase: HaReportTasksBuzUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *BlobSliceJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	sql := `select blob_slices.* from blob_slices inner join blobs on blobs.gid=blob_slices.blob_gid
  where blob_slices.deleted_at=0 and blob_slices.handle_status=0 
  and blobs.handle_status=1 
  and blobs.handle_result=0 and blobs.deleted_at=0`

	return c.CommonUsecase.DB().Raw(sql).Rows()

	//return c.CommonUsecase.DB().
	//	Table(BlobSliceEntity{}.TableName()).
	//	Where("handle_status=? ",
	//		HandleStatus_waiting).Rows()

}

func (c *BlobSliceJobUsecase) Handle(ctx context.Context, task *BlobSliceEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.AppendHandleResultDetail(err.Error())
	} else {
		task.HandleResult = HandleResult_ok
	}
	err = c.CommonUsecase.DB().Save(task).Error
	if err != nil {
		c.log.Error("BlobSliceJobUsecase: Handle: ", task.ID, " : ", err.Error())
	}
	isOk, err := c.BlobSliceUsecase.HasFinish(task.BlobGid)
	if err != nil {
		return err
	}
	if isOk {
		dataEntry := make(TypeDataEntry)
		dataEntry[FieldName_gid] = task.BlobGid
		dataEntry[BlobFieldName_status] = Blob_Status_ready
		dataEntry[FieldName_updated_at] = time.Now().Unix()
		_, err = c.DataEntryUsecase.UpdateOne(Kind_blobs, dataEntry, "gid", nil)

		if err != nil {
			c.log.Error(err)
			return err
		}

		err = c.HaReportTasksBuzUsecase.CreateTask(ctx, task.BlobGid)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	return err
}

func (c *BlobSliceJobUsecase) HandleExec(ctx context.Context, task *BlobSliceEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	//blob, err := c.BlobUsecase.GetByGid(task.Gid)
	//if err != nil {
	//	return err
	//}
	//if blob == nil {
	//	return errors.New("blob is nil")
	//}

	storageRes, err := c.AzstorageUsecase.DownloadStream(ctx, task.FileBlobname)
	if err != nil {
		return err
	}
	defer storageRes.Body.Close()

	operationLocation, err := c.AzcognitiveUsecase.PrebuiltRead(storageRes.Body)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return c.HandleOperationLocation(ctx, task, operationLocation)
}

func (c *BlobSliceJobUsecase) HandleOperationLocation(ctx context.Context, task *BlobSliceEntity, operationLocation string) error {

	if task == nil {
		return errors.New("task is nil")
	}

	res, err := c.AzcognitiveUsecase.GetPrebuiltReadResultWithBlock(ctx, operationLocation)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if res == nil {
		return errors.New("GetPrebuiltReadResultWithBlock:res is nil")
	}
	ocrUuid := InterfaceToString(task.SliceId) + "_" + uuid.UuidWithoutStrike()
	blobname := OcrResultBlobname(task.BlobGid, ocrUuid)
	_, err = c.AzstorageUsecase.UploadStream(ctx, blobname, strings.NewReader(InterfaceToString(res)))
	if err != nil {
		c.log.Error(err)
		return err
	}
	vo := BlobSliceBlobnameVo{
		Blobname: blobname,
	}
	task.OcrResultBlobname = InterfaceToString(vo)
	ocrResultVo := OcrResultVo(res)
	task.HasOcr = 1
	task.OcrResult = InterfaceToString(res)
	task.OcrResultContent = ocrResultVo.GetContent()
	task.UpdatedAt = time.Now().Unix()
	err = c.CommonUsecase.DB().Save(&task).Error
	return err
}
