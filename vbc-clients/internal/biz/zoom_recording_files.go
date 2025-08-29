package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
	"vbc/internal/conf"
)

const HandleResult_HandleProcessing_Recording_Missing = 3 // 记录不见了(zoom的问题)
const HandleResult_HandleProcessing_Error = 4             // 出现异常

type ZoomRecordingFileEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	RecordingFileId    string
	MeetingUuid        string
	RecordingStart     string
	RecordingEnd       string
	FileType           string
	FileExtension      string
	FileSize           string
	PlayUrl            string
	DownloadUrl        string
	Status             string
	RecordingType      string
	BoxResId           string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (c *ZoomRecordingFileEntity) GetFileSize() int64 {
	v, _ := strconv.ParseInt(c.FileSize, 10, 32)
	return v
}

func (ZoomRecordingFileEntity) TableName() string {
	return "zoom_recording_files"
}

func (c *ZoomRecordingFileEntity) FileName() string {
	return fmt.Sprintf("%s.%s", c.RecordingType, strings.ToLower(c.FileExtension))
}

func (c *ZoomRecordingFileEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

type ZoomRecordingFileUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ZoomRecordingFileEntity]
}

func NewZoomRecordingFileUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ZoomRecordingFileUsecase {
	uc := &ZoomRecordingFileUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
