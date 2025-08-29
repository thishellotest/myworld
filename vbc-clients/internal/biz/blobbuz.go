package biz

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type BlobbuzUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	BlobUsecase      *BlobUsecase
	BoxUsecase       *BoxUsecase
	PdfcpuUsecase    *PdfcpuUsecase
	BlobSliceUsecase *BlobSliceUsecase
	AzstorageUsecase *AzstorageUsecase
	MapUsecase       *MapUsecase
	TUsecase         *TUsecase
	DataEntryUsecase *DataEntryUsecase
	PdfGoFitzUsecase *PdfGoFitzUsecase
	WebpUsecase      *WebpUsecase
}

func NewBlobbuzUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BlobUsecase *BlobUsecase,
	BoxUsecase *BoxUsecase,
	PdfcpuUsecase *PdfcpuUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	AzstorageUsecase *AzstorageUsecase,
	MapUsecase *MapUsecase,
	TUsecase *TUsecase,
	DataEntryUsecase *DataEntryUsecase,
	PdfGoFitzUsecase *PdfGoFitzUsecase,
	WebpUsecase *WebpUsecase) *BlobbuzUsecase {
	uc := &BlobbuzUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		BlobUsecase:      BlobUsecase,
		BoxUsecase:       BoxUsecase,
		PdfcpuUsecase:    PdfcpuUsecase,
		BlobSliceUsecase: BlobSliceUsecase,
		AzstorageUsecase: AzstorageUsecase,
		MapUsecase:       MapUsecase,
		TUsecase:         TUsecase,
		DataEntryUsecase: DataEntryUsecase,
		PdfGoFitzUsecase: PdfGoFitzUsecase,
		WebpUsecase:      WebpUsecase,
	}

	return uc
}

func GenUniqblob(fileId, fileVersionId string) string {
	return fmt.Sprintf("%s_%s", fileId, fileVersionId)
}

func (c *BlobbuzUsecase) GetCaseByBoxFileInfo(fileInfo lib.TypeMap) (*TData, error) {

	entries := fileInfo.GetTypeList("path_collection.entries")
	if len(entries) < 4 {
		return nil, errors.New("The client cannot be identified correctly")
	}

	clientFolderId := entries[3].GetString("id")
	if entries[1].GetString("id") == "241183180615" { // 测试帐号
		if len(entries) < 5 {
			return nil, errors.New("The client cannot be identified correctly:1")
		}
		clientFolderId = entries[4].GetString("id")
	}
	if clientFolderId == "" {
		return nil, errors.New("clientFolderId is wrong")
	}
	folderMap, err := c.MapUsecase.GetByCond(And(Eq{"mval": clientFolderId}, Like{"mkey", "ClientBoxFolderId:%"}))
	if err != nil {
		return nil, err
	}
	if folderMap == nil {
		return nil, errors.New("No clients found")
	}
	mkeys := strings.Split(folderMap.Mkey, ":")
	if len(mkeys) != 2 {
		return nil, errors.New("No clients found:1")
	}
	tCase, err := c.TUsecase.DataById(Kind_client_cases, lib.InterfaceToInt32(mkeys[1]))
	if err != nil {
		return nil, err
	}
	if tCase == nil {
		return nil, errors.New("No clients found:2")
	}
	return tCase, nil
}

func (c *BlobbuzUsecase) HandleBoxFile(ctx context.Context, fileInfo lib.TypeMap, caseGid string, userGid string) (tBlob *TData, err error) {

	fileId := fileInfo.GetString("id")
	versionId := fileInfo.GetString("file_version.id")
	if versionId == "" {
		return nil, errors.New(fileId + " : versionId is empty:")
	}
	fileName := fileInfo.GetString("name")
	_, fileSuffixName := lib.FileExt(fileName, true)
	if fileSuffixName != BlobType_pdf {
		return nil, errors.New(fileId + " : Does not support")
	}

	uniqblob := GenUniqblob(fileId, versionId)

	blob, err := c.BlobUsecase.GetByUniqblob(uniqblob)
	if err != nil {
		return nil, err
	}
	if blob != nil {
		c.log.Info("blob exists")
		return nil, nil
	}

	row := make(TypeDataEntry)
	row[BlobFieldName_blob_type] = BlobType_pdf
	row[BlobFieldName_uniqblob] = uniqblob
	row[BlobFieldName_case_gid] = caseGid
	row[BlobFieldName_user_gid] = userGid

	gid, err := c.DataEntryUsecase.InsertOne(Kind_blobs, row, nil)
	if err != nil {
		return nil, err
	}
	return c.TUsecase.DataByGid(Kind_blobs, gid)
}

// HandleBlobSlices fileInfo.GetString("id"), fileInfo.GetString("file_version.id")
func (c *BlobbuzUsecase) HandleBlobSlices(ctx context.Context, blob *TData, boxFileId string, boxVersionId string) error {
	if blob == nil {
		return errors.New("tBlob is nil")
	}
	fileReader, err := c.BoxUsecase.DownloadFile(boxFileId, boxVersionId)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	c.log.Debug("star reader : ", boxFileId)
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return err
	}
	c.log.Debug("end reader : ", boxFileId)
	bytesReader := bytes.NewReader(fileBytes)

	pdfUuid := uuid.UuidWithoutStrike()

	tempDir, err := os.MkdirTemp(configs.GetAppRuntimePath(), pdfUuid)
	if err != nil {
		return err
	}

	sourcePdf := fmt.Sprintf("%s/%s_source.pdf", tempDir, pdfUuid)
	c.log.Debug("sourcePdf:", sourcePdf)
	sourceFile, err := os.Create(sourcePdf)
	if err != nil {
		return err
	}
	defer func() {
		sourceFile.Close()
	}()

	_, err = io.Copy(sourceFile, bytesReader)
	if err != nil {
		return err
	}
	c.log.Debug(tempDir)
	defer func() {
		os.RemoveAll(tempDir)
	}()

	err = pdfcpu.Split(bytesReader, tempDir, pdfUuid, 1, nil)
	if err != nil {
		return err
	}

	err = c.PdfGoFitzUsecase.PdfToImages(sourcePdf, tempDir, pdfUuid+"_img")
	if err != nil {
		return err
	}

	maxPage := 1000
	filePath := fmt.Sprintf("%s/%s_%d.pdf", tempDir, pdfUuid, maxPage+1)
	exceedMaxPage, err := lib.PathExists(filePath)
	if err != nil {
		return err
	}
	if exceedMaxPage {
		return errors.New("The maximum support is 1000 pages")
	}

	originReader := bytes.NewReader(fileBytes)

	blobUuid := uuid.UuidWithoutStrike()
	blobName := FileBlobname(blob.Gid(),
		blob.CustomFields.TextValueByNameBasic(BlobFieldName_blob_type),
		blobUuid)
	_, err = c.AzstorageUsecase.UploadStream(ctx, blobName, originReader)
	if err != nil {
		c.log.Error(err, "blob:", blob.Gid())
		return err
	}
	blobRow := make(TypeDataEntry)
	blobRow[FieldName_gid] = blob.Gid()
	blobRow[BlobFieldName_file_blobname] = blobName
	_, err = c.DataEntryUsecase.UpdateOne(Kind_blobs, blobRow, FieldName_gid, nil)
	if err != nil {
		c.log.Error(err, "blob.Gid: ", blob.Gid())
		return err
	}
	// 最多处理1000页
	for i := 1; i <= maxPage; i++ {
		filePath = fmt.Sprintf("%s/%s_%d.pdf", tempDir, pdfUuid, i)
		ok, err := lib.PathExists(filePath)
		if err != nil {
			return err
		}
		c.log.Debug("filePath:", filePath)
		imgPath := fmt.Sprintf("%s/%s_img_%d.jpg", tempDir, pdfUuid, i)
		ok, err = lib.PathExists(imgPath)
		if err != nil {
			return err
		}
		c.log.Debug("imgPath:", imgPath)
		if ok {
			sliceFileBytes, err := os.ReadFile(filePath)
			if err != nil {
				c.log.Error(err)
				return err
			}

			webpPath := imgPath + ".webp"
			width, height, err := c.WebpUsecase.JpgToWebp(imgPath, webpPath)
			if err != nil {
				return err
			}

			imgPathBtyes, err := os.ReadFile(imgPath)
			if err != nil {
				c.log.Error(err)
				return err
			}

			webpPathBtyes, err := os.ReadFile(webpPath)
			if err != nil {
				c.log.Error(err)
				return err
			}

			sliceUuid := InterfaceToString(i) + "_" + uuid.UuidWithoutStrike()
			blobName = FileBlobname(blob.Gid(),
				blob.CustomFields.TextValueByNameBasic(BlobFieldName_blob_type),
				sliceUuid)
			fileBlobNameJpg := FileBlobnameJpg(blob.Gid(), sliceUuid)
			fileBlobNameWebp := fileBlobNameJpg + ".webp"
			//blobName := fmt.Sprintf("blobs/%s/%d.pdf", pdfUuid, i)
			_, err = c.AzstorageUsecase.UploadStream(ctx, blobName, strings.NewReader(string(sliceFileBytes)))
			if err != nil {
				c.log.Error(err)
				return err
			}
			_, err = c.AzstorageUsecase.UploadStream(ctx, fileBlobNameJpg, strings.NewReader(string(imgPathBtyes)))
			if err != nil {
				c.log.Error(err)
				return err
			}
			_, err = c.AzstorageUsecase.UploadStream(ctx, fileBlobNameWebp, strings.NewReader(string(webpPathBtyes)))
			if err != nil {
				c.log.Error(err)
				return err
			}
			sliceGid := uuid.UuidWithoutStrike()

			FileBlobnameJpgVo := BlobSliceBlobnameVo{
				Blobname: fileBlobNameJpg,
				Width:    width,
				Height:   height,
			}
			fileBlobNameWebpVo := BlobSliceBlobnameVo{
				Blobname: fileBlobNameWebp,
				Width:    width,
				Height:   height,
			}

			c.log.Debug("FileBlobnameJpgVo.Blobname:", FileBlobnameJpgVo.Blobname)
			c.log.Debug("fileBlobNameWebpVo.Blobname:", fileBlobNameWebpVo.Blobname)

			entity := &BlobSliceEntity{
				Gid:              sliceGid,
				HandleStatus:     BlobSlice_HandleStatus_wait_operation,
				BlobGid:          blob.CustomFields.TextValueByNameBasic("gid"),
				SliceId:          InterfaceToString(i),
				FileBlobname:     blobName,
				FileBlobnameJpg:  InterfaceToString(FileBlobnameJpgVo),
				FileBlobnameWebp: InterfaceToString(fileBlobNameWebpVo),
				CreatedAt:        time.Now().Unix(),
				UpdatedAt:        time.Now().Unix(),
			}
			err = c.BlobSliceUsecase.CommonUsecase.DB().Create(entity).Error
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (c *BlobbuzUsecase) BizSliceJoinOcr(ctx context.Context, userFacade UserFacade, blobSliceGid string) (lib.TypeMap, error) {

	sliceEntity, _ := c.BlobSliceUsecase.GetByCond(Eq{"gid": blobSliceGid})
	if sliceEntity == nil {
		return nil, errors.New("sliceEntity is nil")
	}
	if sliceEntity.HandleStatus == BlobSlice_HandleStatus_wait_operation {
		sliceEntity.HandleStatus = BlobSlice_HandleStatus_waiting // 零值没有保存
		err := c.CommonUsecase.DB().Select("HandleStatus").Save(&sliceEntity).Error
		if err != nil {
			return nil, err
		}
		row := sliceEntity.ToApi(c.CommonUsecase, c.log, c.AzstorageUsecase)
		data := make(lib.TypeMap)
		data.Set("data", row)
		return data, nil
	}
	return nil, errors.New("Not allowed to operate")
}
