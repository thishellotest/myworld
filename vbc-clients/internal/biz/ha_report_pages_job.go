package biz

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"vbc/internal/conf"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type HaReportPageJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[HaReportPageEntity]
	BaseHandleCustom[HaReportPageEntity]
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
	DataEntryUsecase    *DataEntryUsecase
	BlobCommentUsecase  *BlobCommentUsecase
}

func NewHaReportPageJobUsecase(logger log.Logger,
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
	HaReportTaskUsecase *HaReportTaskUsecase,
	DataEntryUsecase *DataEntryUsecase,
	BlobCommentUsecase *BlobCommentUsecase) *HaReportPageJobUsecase {
	uc := &HaReportPageJobUsecase{
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
		DataEntryUsecase:    DataEntryUsecase,
		BlobCommentUsecase:  BlobCommentUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandleCustom.Log = log.NewHelper(logger)
	uc.BaseHandleCustom.DB = CommonUsecase.DB()

	return uc
}

func (c *HaReportPageJobUsecase) WaitingTasks(ctx context.Context) (*sql.Rows, error) {

	sql := fmt.Sprintf(`select * from ha_report_pages where handle_status=%d`, HandleStatus_waiting)
	return c.CommonUsecase.DB().Raw(sql).Rows()

}

func (c *HaReportPageJobUsecase) Handle(ctx context.Context, task *HaReportPageEntity) error {

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
		c.log.Error("HaReportPageJobUsecase: Handle: ", task.ID, " : ", err.Error())
		return err
	}

	// 判断任务是否全部完成
	whetherExistEntity, err := c.HaReportPageUsecase.GetByCond(Eq{
		"handle_status": HandleStatus_waiting,
		"ha_report_gid": task.HaReportGid})
	if err != nil {
		return err
	}
	if whetherExistEntity == nil {

		haReport, err := c.HaReportTaskUsecase.GetByCond(Eq{"gid": task.HaReportGid})
		if err != nil {
			c.log.Error(err)
			return err
		}
		if haReport == nil {
			return errors.New("haReport is nil")
		}

		row := make(TypeDataEntry)
		row["gid"] = haReport.BlobGid
		row[BlobFieldName_status] = Blob_Status_ready
		_, err = c.DataEntryUsecase.UpdateOne(Kind_blobs, row, FieldName_gid, nil)
		if err != nil {
			c.log.Error(err)
			return err
		}

		// 已经全部完成了
		//err = c.HaReportTaskUsecase.CommonUsecase.DB().Model(&HaReportTaskEntity{}).
		//	Where("gid=?", task.HaReportGid).
		//	Updates(map[string]interface{}{
		//		"handle_status": HaReportTask_HandleStatus_WaitingCreatingPdf,
		//		"updated_at":    time.Now().Unix()}).Error
	}

	return err
}

func (c *HaReportPageJobUsecase) HandleExec(ctx context.Context, task *HaReportPageEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}

	blobSlice, err := c.BlobSliceUsecase.GetByCond(Eq{"gid": task.BlobSliceGid, "deleted_at": 0})
	if err != nil {
		return err
	}
	if blobSlice == nil {
		return errors.New("blobSlice is nil")
	}

	aiReport, err := c.GetAiReport(ctx, blobSlice)
	if err != nil {
		c.log.Error(err)
		return err
	}

	task.AiReport = InterfaceToString(aiReport)
	task.UpdatedAt = time.Now().Unix()

	haReport, err := c.HaReportTaskUsecase.GetByCond(Eq{"gid": task.HaReportGid})
	if err != nil {
		c.log.Error(err)
		return err
	}
	if haReport == nil {
		return errors.New("haReport is nil")
	}

	tBlob, err := c.BlobUsecase.GetByGid(haReport.BlobGid)
	if err != nil {
		c.log.Error(err)
		return err
	}
	if tBlob == nil {
		return errors.New("tBlob is nil")
	}

	for _, v := range aiReport.Conditions {
		if v != "" {

			page, _ := strconv.ParseInt(blobSlice.SliceId, 10, 32)
			entity := BlobCommentEntity{
				Gid:         uuid.UuidWithoutStrike(),
				BlobGid:     haReport.BlobGid,
				Content:     v,
				Page:        int32(page),
				Type:        BlobComment_Type_OnlyText,
				UserGid:     tBlob.CustomFields.TextValueByNameBasic(BlobFieldName_user_gid),
				HaReportGid: haReport.Gid,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			}
			err = c.BlobCommentUsecase.CommonUsecase.DB().Save(&entity).Error
			if err != nil {
				c.log.Error(err)
				return err
			}
		}
	}

	return nil
}

func (c *HaReportPageJobUsecase) GetAiReport(ctx context.Context, blobSliceEntity *BlobSliceEntity) (getDiseaseNamesByMedicalTextWithAIResponse GetDiseaseNamesByMedicalTextWithAIResponse, err error) {
	if blobSliceEntity == nil {
		return getDiseaseNamesByMedicalTextWithAIResponse, errors.New("haReportPageEntity is nil")
	}

	ocrResult, err := c.BlobSliceUsecase.GetOcrResult(ctx, blobSliceEntity)
	if err != nil {
		return getDiseaseNamesByMedicalTextWithAIResponse, err
	}
	content := ocrResult.GetContent()
	if content == "" {
		return getDiseaseNamesByMedicalTextWithAIResponse, nil
	}
	fromUniqKey := fmt.Sprintf("%s:%s", blobSliceEntity.Gid, time.Now().Format(time.RFC3339))
	getDiseaseNamesByMedicalTextWithAIResponse, err = c.HaiUsecase.GetDiseaseNamesByMedicalTextWithAI(ctx, content, fromUniqKey)
	if err != nil {
		return getDiseaseNamesByMedicalTextWithAIResponse, err
	}
	return getDiseaseNamesByMedicalTextWithAIResponse, nil
}
