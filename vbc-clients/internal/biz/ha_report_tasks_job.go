package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"os"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HaReportTaskJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[HaReportTaskEntity]
	BaseHandleCustom[HaReportTaskEntity]
	AzcognitiveUsecase  *AzcognitiveUsecase
	BlobUsecase         *BlobUsecase
	BlobSliceUsecase    *BlobSliceUsecase
	AzstorageUsecase    *AzstorageUsecase
	HaReportPageUsecase *HaReportPageUsecase
	HaiUsecase          *HaiUsecase
	HaReportPdfUsecase  *HaReportPdfUsecase
	TUsecase            *TUsecase
	BoxUsecase          *BoxUsecase
	DataComboUsecase    *DataComboUsecase
	HaReportTaskUsecase *HaReportTaskUsecase
}

func NewHaReportTaskJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AzcognitiveUsecase *AzcognitiveUsecase,
	BlobUsecase *BlobUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	AzstorageUsecase *AzstorageUsecase,
	HaReportPageUsecase *HaReportPageUsecase,
	HaiUsecase *HaiUsecase,
	HaReportPdfUsecase *HaReportPdfUsecase,
	TUsecase *TUsecase,
	BoxUsecase *BoxUsecase,
	DataComboUsecase *DataComboUsecase,
	HaReportTaskUsecase *HaReportTaskUsecase) *HaReportTaskJobUsecase {
	uc := &HaReportTaskJobUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		AzcognitiveUsecase:  AzcognitiveUsecase,
		BlobUsecase:         BlobUsecase,
		BlobSliceUsecase:    BlobSliceUsecase,
		AzstorageUsecase:    AzstorageUsecase,
		HaReportPageUsecase: HaReportPageUsecase,
		HaiUsecase:          HaiUsecase,
		HaReportPdfUsecase:  HaReportPdfUsecase,
		TUsecase:            TUsecase,
		BoxUsecase:          BoxUsecase,
		DataComboUsecase:    DataComboUsecase,
		HaReportTaskUsecase: HaReportTaskUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

//func (c *HaReportTaskJobUsecase) WaitingTasksByCreatingPdf(ctx context.Context) (*sql.Rows, error) {
//
//	sql := fmt.Sprintf(`select ha_report_tasks.* from ha_report_tasks
//where ha_report_tasks.handle_status=%d and ha_report_tasks.deleted_at=0`, HaReportTask_HandleStatus_WaitingCreatingPdf)
//	return c.CommonUsecase.DB().Raw(sql).Rows()
//
//}

//
//func (c *HaReportTaskJobUsecase) HandleByCreatingPdf(ctx context.Context, task *HaReportTaskEntity) error {
//
//	if task == nil {
//		return errors.New("task is nil.")
//	}
//	err := c.HandleExecByCreatingPdf(ctx, task)
//	task.HandleStatus = HaReportTask_HandleStatus_CreatingPdfDone
//	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
//	if err != nil {
//		task.HandleResult = HandleResult_failure
//		task.AppendHandleResultDetail(err.Error())
//	} else {
//		task.HandleResult = HandleResult_ok
//	}
//	err = c.CommonUsecase.DB().Save(task).Error
//	if err != nil {
//		c.log.Error("HaReportTaskJobUsecase: HandleByCreatingPdf: ", task.ID, " : ", err.Error())
//	}
//	return err
//}

//
//func (c *HaReportTaskJobUsecase) HandleExecByCreatingPdf(ctx context.Context, task *HaReportTaskEntity) error {
//
//	if task == nil {
//		return errors.New("task is nil")
//	}
//
//	tCase, err := c.TUsecase.DataById(Kind_client_cases, task.ClientCaseId)
//	if err != nil {
//		return err
//	}
//	if tCase == nil {
//		return errors.New("tCase is nil")
//	}
//
//	blobEntity, err := c.BlobUsecase.GetByGid(task.BlobGid)
//	if err != nil {
//		return err
//	}
//	if blobEntity == nil {
//		return errors.New("blobEntity is nil")
//	}
//
//	// 开记生成pdf报告
//	pdfFilePath, err := c.HaReportPdfUsecase.CreateHaReportPdf(ctx, task, tCase)
//	if err != nil {
//		return err
//	}
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	if pdfFilePath == "" {
//		return errors.New("pdfFilePath is empty")
//	}
//	defer func() {
//		os.Remove(pdfFilePath)
//	}()
//
//	boxBlobFileId, _ := GetUniqblobFileInfo(blobEntity.CustomFields.TextValueByNameBasic("uniqblob"))
//
//	boxFileId, err := c.UploadPdfReportToBox(ctx, task, tCase, blobEntity.CustomFields.TextValueByNameBasic("name"), boxBlobFileId, pdfFilePath)
//	if err != nil {
//		return err
//	}
//	task.BoxFileId = boxFileId
//	return c.CommonUsecase.DB().Save(&task).Error
//}

func (c *HaReportTaskJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	sql := fmt.Sprintf(`select ha_report_tasks.* from ha_report_tasks 
inner join blobs on blobs.gid=ha_report_tasks.blob_gid
where ha_report_tasks.handle_status=%d and blobs.status=%s and blobs.deleted_at=0 and ha_report_tasks.deleted_at=0`, HandleStatus_waiting, Blob_Status_ready)
	return c.CommonUsecase.DB().Raw(sql).Rows()

}

func (c *HaReportTaskJobUsecase) Handle(ctx context.Context, task *HaReportTaskEntity) error {

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
		c.log.Error("HaReportTaskJobUsecase: Handle: ", task.ID, " : ", err.Error())
	}

	return err
}

func (c *HaReportTaskJobUsecase) HandleExec(ctx context.Context, task *HaReportTaskEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	slices, err := c.BlobSliceUsecase.AllByCond(Eq{"blob_gid": task.BlobGid, "deleted_at": 0})
	if err != nil {
		return err
	}
	for _, v := range slices {
		haPage, err := c.HaReportPageUsecase.GetByCond(Eq{"blob_slice_gid": v.Gid, "ha_report_gid": task.Gid})
		if err != nil {
			return err
		}
		if haPage == nil {
			haPage = &HaReportPageEntity{
				HaReportGid:  task.Gid,
				BlobSliceGid: v.Gid,
				CreatedAt:    time.Now().Unix(),
				UpdatedAt:    time.Now().Unix(),
			}
			err = c.CommonUsecase.DB().Save(&haPage).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//
//func (c *HaReportTaskJobUsecase) HandleExecBackup(ctx context.Context, task *HaReportTaskEntity) error {
//
//	if task == nil {
//		return errors.New("task is nil")
//	}
//	slices, err := c.BlobSliceUsecase.AllByCond(Eq{"blob_gid": task.BlobGid, "deleted_at": 0})
//	if err != nil {
//		return err
//	}
//	for k, _ := range slices {
//
//		// 判断任务是否已经关闭，已经关闭了，不再执行后续任务
//		haReportTask, err := c.HaReportTaskUsecase.GetByCond(Eq{"gid": task.Gid, "deleted_at": 0})
//		if err != nil {
//			c.log.Error(err)
//			return err
//		}
//		if haReportTask == nil {
//			return errors.New("haReportTask is nil")
//		}
//		if haReportTask.HandleStatus == HandleStatus_done {
//			return errors.New("HandleStatus is HandleStatus_done")
//		}
//
//		err = c.HandleAiReport(ctx, task, slices[k])
//		if err != nil {
//			c.log.Error(err)
//			return err
//		}
//	}
//
//	tCase, err := c.TUsecase.DataById(Kind_client_cases, task.ClientCaseId)
//	if err != nil {
//		return err
//	}
//	if tCase == nil {
//		return errors.New("tCase is nil")
//	}
//
//	blobEntity, err := c.BlobUsecase.GetByGid(task.BlobGid)
//	if err != nil {
//		return err
//	}
//	if blobEntity == nil {
//		return errors.New("blobEntity is nil")
//	}
//
//	//tClient, _, err := c.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
//	//if err != nil {
//	//	return err
//	//}
//	//if tClient == nil {
//	//	return errors.New("tClient is nil")
//	//}
//
//	// 开记生成pdf报告
//	pdfFilePath, err := c.HaReportPdfUsecase.CreateHaReportPdf(ctx, task, tCase)
//	if err != nil {
//		return err
//	}
//	if err != nil {
//		c.log.Error(err)
//		return err
//	}
//	if pdfFilePath == "" {
//		return errors.New("pdfFilePath is empty")
//	}
//	defer func() {
//		os.Remove(pdfFilePath)
//	}()
//
//	boxBlobFileId, _ := GetUniqblobFileInfo(blobEntity.CustomFields.TextValueByNameBasic("uniqblob"))
//
//	boxFileId, err := c.UploadPdfReportToBox(ctx, task, tCase, blobEntity.CustomFields.TextValueByNameBasic("name"), boxBlobFileId, pdfFilePath)
//	if err != nil {
//		return err
//	}
//	task.BoxFileId = boxFileId
//	return c.CommonUsecase.DB().Save(&task).Error
//}

func (c *HaReportTaskJobUsecase) UploadPdfReportToBox(ctx context.Context,
	task *HaReportTaskEntity,
	tCase *TData,
	sourcePdfName string,
	fileId string,
	pdfFilePath string) (boxFileId string, err error) {

	if task == nil {
		return "", errors.New("task is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	if pdfFilePath == "" {
		return "", errors.New("pdfFilePath is empty")
	}
	folderId := "274235574426"
	if configs.IsProd() {
		folderId = "274040007842"
	}
	sourcePdfPrefix, _ := lib.FileExt(sourcePdfName, false)

	fileName := fmt.Sprintf("%d#%s(%s)_%d.pdf", tCase.CustomFields.NumberValueByNameBasic("id"), sourcePdfPrefix, fileId, task.ID)
	pdfFilePathFile, err := os.Open(pdfFilePath)
	if err != nil {
		return "", err
	}
	defer func() {
		pdfFilePathFile.Close()
	}()

	boxFileId, err = c.BoxUsecase.UploadFile(folderId, pdfFilePathFile, fileName)
	if err != nil {
		return "", err
	}
	return
}

func (c *HaReportTaskJobUsecase) HandleAiReport(ctx context.Context, task *HaReportTaskEntity, blobSliceEntity *BlobSliceEntity) error {

	isFinish, err := c.HaReportPageUsecase.IsFinish(task.Gid, blobSliceEntity.Gid)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if !isFinish {
		aiReport, err := c.GetAiReport(ctx, blobSliceEntity)
		if err != nil {
			c.log.Error(err)
			return err
		}

		haReportPage := &HaReportPageEntity{
			HaReportGid:  task.Gid,
			BlobSliceGid: blobSliceEntity.Gid,
			AiReport:     aiReport,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		err = c.HaReportPageUsecase.CommonUsecase.DB().Save(&haReportPage).Error
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	return nil
}

func (c *HaReportTaskJobUsecase) GetAiReport(ctx context.Context, blobSliceEntity *BlobSliceEntity) (string, error) {
	if blobSliceEntity == nil {
		return "", errors.New("haReportPageEntity is nil")
	}

	ocrResult, err := c.BlobSliceUsecase.GetOcrResult(ctx, blobSliceEntity)
	if err != nil {
		return "", err
	}
	content := ocrResult.GetContent()
	if content == "" {
		return "", nil
	}
	fromUniqKey := fmt.Sprintf("%s:%s", blobSliceEntity.Gid, time.Now().Format(time.RFC3339))
	res, err := c.HaiUsecase.GetDiseaseNamesByMedicalTextWithAI(ctx, content, fromUniqKey)
	if err != nil {
		return "", err
	}
	return InterfaceToString(res), nil
}
