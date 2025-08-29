package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

const (
	HaReportTask_HandleStatus_WaitingCreatingPdf = 3
	HaReportTask_HandleStatus_CreatingPdfDone    = 4
)

type HaReportTaskEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	Gid                string
	BlobGid            string
	ReportBlobname     string
	BoxFileId          string // 报告上传到box的ID
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (HaReportTaskEntity) TableName() string {
	return "ha_report_tasks"
}

func (c *HaReportTaskEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

type HaReportTaskUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[HaReportTaskEntity]
}

func NewHaReportTaskUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *HaReportTaskUsecase {
	uc := &HaReportTaskUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
