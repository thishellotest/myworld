package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

func GetBlobsParentName() string {
	if configs.IsProd() {
		return "blobs"
	} else {
		return "blobs_dev"
	}
}

func FileBlobname(blobGid string, blobType string, uuid string) string {
	blobName := fmt.Sprintf("%s/%s/%s.%s", GetBlobsParentName(), blobGid, uuid, blobType)
	return blobName
}

func FileBlobnameJpg(blobGid string, uuid string) string {
	blobName := fmt.Sprintf("%s/%s/%s_jpg.jpg", GetBlobsParentName(), blobGid, uuid)
	return blobName
}

func OcrResultBlobname(blobGid string, uuid string) string {
	blobName := fmt.Sprintf("%s/%s/%s.json", GetBlobsParentName(), blobGid, uuid)
	return blobName
}

const (
	BlobSlice_HandleStatus_waiting        = 0 // ocr进行中
	BlobSlice_HandleStatus_done           = 1 // ocr完成
	BlobSlice_HandleStatus_wait_operation = 3 // 等待用户操作， 成为 BlobSlice_HandleStatus_waiting
)

type BlobSliceEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	Gid                string
	BlobGid            string
	SliceId            string
	FileBlobname       string
	FileBlobnameJpg    string
	FileBlobnameWebp   string
	OcrResultBlobname  string
	HasOcr             int
	OcrResult          string
	OcrResultContent   string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

type BlobSliceBlobnameVo struct {
	Blobname   string `json:"blobname"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	ExpiryTime int64  `json:"expiry_time"`
	Url        string `json:"url"`
}

func (c *BlobSliceEntity) GetFileBlobnameJpgVo() (fileBlobnameJpgVo BlobSliceBlobnameVo) {
	return lib.StringToTDef(c.FileBlobnameJpg, fileBlobnameJpgVo)
}

func (c *BlobSliceEntity) GetFileBlobnameWebpVo() (fileBlobnameWebpVo BlobSliceBlobnameVo) {
	return lib.StringToTDef(c.FileBlobnameWebp, fileBlobnameWebpVo)
}

func (c *BlobSliceEntity) GetOCRJsonVo() (ocrJsonVo BlobSliceBlobnameVo) {
	return lib.StringToTDef(c.OcrResultBlobname, ocrJsonVo)
}

func (BlobSliceEntity) TableName() string {
	return "blob_slices"
}

func (c *BlobSliceEntity) GetBlobNameUrlVo(CommonUsecase *CommonUsecase, log *log.Helper, AzstorageUsecase *AzstorageUsecase) (fileBlobnameJpgVo BlobSliceBlobnameVo, fileBlobnameWebpVo BlobSliceBlobnameVo, ocrJsonVo BlobSliceBlobnameVo) {

	fileBlobnameJpgVo = c.GetFileBlobnameJpgVo()
	fileBlobnameWebpVo = c.GetFileBlobnameWebpVo()
	ocrJsonVo = c.GetOCRJsonVo()

	currentTime := time.Now().Unix()
	needUpdate := false
	if fileBlobnameJpgVo.Blobname != "" && fileBlobnameJpgVo.ExpiryTime <= currentTime { // 过期了
		newExpiryTime := time.Now().UTC().Add(20 * 48 * time.Hour)
		realNewExpiryTime := time.Now().UTC().Add(21 * 48 * time.Hour)
		url, err := AzstorageUsecase.SasReadUrl(fileBlobnameJpgVo.Blobname, &realNewExpiryTime)
		if err != nil {
			log.Error(err)
		} else {
			fileBlobnameJpgVo.Url = url
			fileBlobnameJpgVo.ExpiryTime = newExpiryTime.Unix()
			needUpdate = true
		}
	}
	if fileBlobnameWebpVo.Blobname != "" && fileBlobnameWebpVo.ExpiryTime <= currentTime { // 过期了
		newExpiryTime := time.Now().UTC().Add(20 * 48 * time.Hour)
		realNewExpiryTime := time.Now().UTC().Add(21 * 48 * time.Hour)
		url, err := AzstorageUsecase.SasReadUrl(fileBlobnameWebpVo.Blobname, &realNewExpiryTime)
		if err != nil {
			log.Error(err)
		} else {
			fileBlobnameWebpVo.Url = url
			fileBlobnameWebpVo.ExpiryTime = newExpiryTime.Unix()
			needUpdate = true
		}
	}
	if ocrJsonVo.Blobname != "" && ocrJsonVo.ExpiryTime <= currentTime { // 过期了
		newExpiryTime := time.Now().UTC().Add(20 * 48 * time.Hour)
		realNewExpiryTime := time.Now().UTC().Add(21 * 48 * time.Hour)
		url, err := AzstorageUsecase.SasReadUrl(ocrJsonVo.Blobname, &realNewExpiryTime)
		if err != nil {
			log.Error(err)
		} else {
			ocrJsonVo.Url = url
			ocrJsonVo.ExpiryTime = newExpiryTime.Unix()
			needUpdate = true
		}
	}
	if needUpdate {
		c.FileBlobnameWebp = InterfaceToString(fileBlobnameWebpVo)
		c.FileBlobnameJpg = InterfaceToString(fileBlobnameJpgVo)
		c.OcrResultBlobname = InterfaceToString(ocrJsonVo)
		err := CommonUsecase.DB().Select("FileBlobnameWebp", "FileBlobnameJpg").Save(&c).Error
		if err != nil {
			log.Error(err)
		}
	}
	return
}

func (c *BlobSliceEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

var (
	BlobSlice_ToApi_Columns = []string{"slice_id", "file_blobname_webp", "gid", "has_ocr", "handle_status", "handle_result"}
)

func (c *BlobSliceEntity) ToApi(CommonUsecase *CommonUsecase, log *log.Helper, AzstorageUsecase *AzstorageUsecase) lib.TypeMap {
	data := make(lib.TypeMap)
	_, webpVo, _ := c.GetBlobNameUrlVo(CommonUsecase, log, AzstorageUsecase)
	//data.Set("jpg", fileBlobnameJpgVo)
	data.Set("page", c.SliceId)
	data.Set("webp", webpVo)
	data.Set("has_ocr", c.HasOcr)
	data.Set("gid", c.Gid)
	if c.HandleStatus == BlobSlice_HandleStatus_done && c.HandleResult == 0 {
		data.Set("progress_status", BlobSlice_HandleStatus_done) // 处理完成
	} else if c.HandleStatus == BlobSlice_HandleStatus_wait_operation {
		data.Set("progress_status", BlobSlice_HandleStatus_wait_operation) // 处理用户操作完成
	} else {
		data.Set("progress_status", BlobSlice_HandleStatus_waiting) // 正在处理
	}
	return data
}

type BlobSliceUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BlobSliceEntity]
	AzstorageUsecase *AzstorageUsecase
}

func NewBlobSliceUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AzstorageUsecase *AzstorageUsecase) *BlobSliceUsecase {
	uc := &BlobSliceUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		AzstorageUsecase: AzstorageUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

// HasFinish 判断是OCR是否完成
func (c *BlobSliceUsecase) HasFinish(blobGid string) (bool, error) {

	if blobGid == "" {
		return false, errors.New("blobGid is empty")
	}
	a, err := c.GetByCond(Eq{"deleted_at": 0,
		"handle_status": HandleStatus_waiting,
		"blob_gid":      blobGid,
		"handle_result": HandleResult_ok})
	if err != nil {
		return false, err
	}
	if a == nil {
		return true, nil
	}
	return false, nil
}

func (c *BlobSliceUsecase) GetByGid(gid string) (*BlobSliceEntity, error) {
	return c.GetByCond(Eq{"gid": gid})
}

func (c *BlobSliceUsecase) GetOcrResultByGid(ctx context.Context, gid string) (OcrResultVo, error) {
	a, err := c.GetByCond(Eq{"gid": gid})
	if err != nil {
		return nil, err
	}
	return c.GetOcrResult(ctx, a)
}

func (c *BlobSliceUsecase) GetOcrResult(ctx context.Context, blobSliceEntity *BlobSliceEntity) (OcrResultVo, error) {
	if blobSliceEntity == nil {
		return nil, errors.New("blobSliceEntity is nil")
	}
	if blobSliceEntity.OcrResultBlobname == "" {
		return nil, errors.New(blobSliceEntity.Gid + " :blobSliceEntity.OcrResultBlobname is empty")
	}
	DownloadStreamResponse, err := c.AzstorageUsecase.DownloadStream(ctx, blobSliceEntity.OcrResultBlobname)
	if err != nil {
		return nil, err
	}
	defer DownloadStreamResponse.Body.Close()
	ocrStr, err := io.ReadAll(DownloadStreamResponse.Body)
	if err != nil {
		return nil, err
	}
	ocrResultVo := lib.ToTypeMapByString(string(ocrStr))
	return OcrResultVo(ocrResultVo), nil
}

type OcrResultVo lib.TypeMap

func (c OcrResultVo) GetContent() string {
	return lib.TypeMap(c).GetString("analyzeResult.content")
}
