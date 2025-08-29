package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
)

type BlobJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	BaseHandleT[TData]
	AzcognitiveUsecase *AzcognitiveUsecase
	BlobUsecase        *BlobUsecase
	BlobSliceUsecase   *BlobSliceUsecase
	AzstorageUsecase   *AzstorageUsecase
	BlobbuzUsecase     *BlobbuzUsecase
	BoxUsecase         *BoxUsecase
	TUsecase           *TUsecase
	DataEntryUsecase   *DataEntryUsecase
}

func NewBlobJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AzcognitiveUsecase *AzcognitiveUsecase,
	BlobUsecase *BlobUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	AzstorageUsecase *AzstorageUsecase,
	BlobbuzUsecase *BlobbuzUsecase,
	BoxUsecase *BoxUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase) *BlobJobUsecase {
	uc := &BlobJobUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		AzcognitiveUsecase: AzcognitiveUsecase,
		BlobUsecase:        BlobUsecase,
		BlobSliceUsecase:   BlobSliceUsecase,
		AzstorageUsecase:   AzstorageUsecase,
		BlobbuzUsecase:     BlobbuzUsecase,
		BoxUsecase:         BoxUsecase,
		TUsecase:           TUsecase,
		DataEntryUsecase:   DataEntryUsecase,
	}

	uc.BaseHandleT.Log = log.NewHelper(logger)

	return uc
}

func (c *BlobJobUsecase) WaitingTasks(ctx context.Context) ([]TData, error) {

	sql := `select * from blobs where blobs.handle_status=0 and deleted_at=0`
	list, err := c.TUsecase.ListByRawSql(Kind_blobs, sql)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *BlobJobUsecase) Handle(ctx context.Context, task TData) error {
	err := c.HandleExec(ctx, &task)

	update := make(TypeDataEntry)
	update[FieldName_gid] = task.CustomFields.TextValueByNameBasic("gid")
	update["handle_status"] = HandleStatus_done
	update[FieldName_updated_at] = time.Now().Unix()

	if err != nil {
		update["handle_result"] = HandleResult_failure
		update["handle_result_detail"] = AppendHandleResultDetail(&task, err)
	} else {
		update["handle_result"] = HandleResult_ok
	}
	_, err = c.DataEntryUsecase.UpdateOne(Kind_blobs, update, "gid", nil)
	if err != nil {
		c.log.Error("BlobJobUsecase: Handle: ", task.CustomFields.NumberValueByNameBasic("id"), " : ", err.Error())
	}
	return err
}

func (c *BlobJobUsecase) HandleExec(ctx context.Context, tBlob *TData) error {
	if tBlob == nil {
		return errors.New("tBlob is nil")
	}
	fileId, versionId := GetUniqblobFileInfo(tBlob.CustomFields.TextValueByNameBasic("uniqblob"))
	return c.BlobbuzUsecase.HandleBlobSlices(ctx, tBlob, fileId, versionId)
}
