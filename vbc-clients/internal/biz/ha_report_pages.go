package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type HaReportPageEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	HaReportGid        string
	BlobSliceGid       string
	AiReport           string
	CreatedAt          int64
	UpdatedAt          int64
}

func (c *HaReportPageEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

func (c *HaReportPageEntity) GetAiReportToPdfText() (text string) {

	aa := lib.ToTypeMapByString(c.AiReport)
	diseaseNames := lib.InterfaceToTDef[[]string](aa.Get("diseaseNames"), nil)
	for _, v := range diseaseNames {
		if text == "" {
			text += v
		} else {
			text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v
			//text += "\n" + v

		}
	}
	return text
}

func (HaReportPageEntity) TableName() string {
	return "ha_report_pages"
}

type HaReportPageUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[HaReportPageEntity]
}

func NewHaReportPageUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *HaReportPageUsecase {
	uc := &HaReportPageUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *HaReportPageUsecase) IsFinish(haReportGid string, sliceGid string) (bool, error) {

	r, err := c.GetByCond(Eq{"ha_report_gid": haReportGid, "blob_slice_gid": sliceGid})
	if err != nil {
		return false, err
	}
	if r == nil {
		return false, nil
	}
	return true, nil
}
