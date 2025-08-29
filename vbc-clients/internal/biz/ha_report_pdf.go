package biz

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"io"
	"os"
	"time"
	"vbc/configs"
	"vbc/internal/conf"

	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type HaReportPdfUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TTemplateEntity]
	PdfUsecase          *PdfUsecase
	BlobSliceUsecase    *BlobSliceUsecase
	HaReportPageUsecase *HaReportPageUsecase
	HaReportTaskUsecase *HaReportTaskUsecase
	TUsecase            *TUsecase
	AzstorageUsecase    *AzstorageUsecase
	BlobUsecase         *BlobUsecase
}

func NewHaReportPdfUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	PdfUsecase *PdfUsecase,
	BlobSliceUsecase *BlobSliceUsecase,
	HaReportPageUsecase *HaReportPageUsecase,
	HaReportTaskUsecase *HaReportTaskUsecase,
	TUsecase *TUsecase,
	AzstorageUsecase *AzstorageUsecase,
	BlobUsecase *BlobUsecase) *HaReportPdfUsecase {
	uc := &HaReportPdfUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		PdfUsecase:          PdfUsecase,
		BlobSliceUsecase:    BlobSliceUsecase,
		HaReportPageUsecase: HaReportPageUsecase,
		HaReportTaskUsecase: HaReportTaskUsecase,
		TUsecase:            TUsecase,
		AzstorageUsecase:    AzstorageUsecase,
		BlobUsecase:         BlobUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *HaReportPdfUsecase) HandleHaReportPdf(ctx context.Context, haReportGid string) error {

	haReport, err := c.HaReportTaskUsecase.GetByCond(Eq{"gid": haReportGid})
	if err != nil {
		return err
	}
	if haReport == nil {
		return errors.New("haReport is nil")
	}

	tBlob, err := c.BlobUsecase.GetByGid(haReport.BlobGid)
	if err != nil {
		return err
	}
	if tBlob == nil {
		return errors.New("tBlob is nil")
	}
	tCase, err := c.TUsecase.Data(Kind_client_cases, Eq{"gid": tBlob.CustomFields.TextValueByNameBasic(BlobFieldName_case_gid)})
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	filePath, err := c.CreateHaReportPdf(ctx, haReport, tCase)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(filePath)
	}()

	return nil
}

func (c *HaReportPdfUsecase) CreateHaReportPdf(ctx context.Context, haReport *HaReportTaskEntity, tCase *TData) (pdfFilePath string, err error) {

	if haReport == nil {
		return "", errors.New("haReport is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	countSql := fmt.Sprintf("select count(*) c from blob_slices where blob_gid='%s' and deleted_at=0", haReport.BlobGid)
	total, err := c.BlobSliceUsecase.CommonUsecase.Count(c.CommonUsecase.DB(), countSql)
	if err != nil {
		c.log.Error(err)
		return "", err
	}

	sql := fmt.Sprintf(`select ha_report_pages.* from ha_report_pages 
inner join blob_slices on ha_report_pages.blob_slice_gid=blob_slices.gid
where ha_report_gid='%s' and blob_slices.deleted_at=0
order by convert(blob_slices.slice_id , signed)`, haReport.Gid)
	sqlRows, err := c.HaReportPageUsecase.CommonUsecase.DB().Raw(sql).Rows()
	if err != nil {
		return "", err
	}
	defer func() {
		sqlRows.Close()
	}()
	var pages []*HaReportPageEntity
	pages, err = lib.SqlRowsToEntities[HaReportPageEntity](c.CommonUsecase.DB(), sqlRows)
	if err != nil {
		return "", err
	}

	clientCaseName := tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)

	type TStruct struct {
		PdfFilePath        string
		HaReportPageEntity *HaReportPageEntity
		BlobSliceEntity    *BlobSliceEntity
	}
	var pdfs []TStruct
	for k, v := range pages {
		blobSliceEntity, err := c.BlobSliceUsecase.GetByGid(v.BlobSliceGid)
		if err != nil {
			return "", err
		}
		if blobSliceEntity == nil {
			return "", errors.New("blobSliceEnity is nil")
		}

		diseaseNames := v.GetAiReportToPdfText()
		if diseaseNames == "" {
			continue
		}

		var fileBytes *bytes.Reader
		if false && configs.IsDev() {
			// begin test code
			file, _ := os.Open("./tmp/STR Full_2.pdf")
			defer func() {
				file.Close()
			}()
			blobSliceBytes, err := io.ReadAll(file)
			if err != nil {
				c.log.Error("STR Full_2.pdf Dev: ", err)
				return "", err
			}
			fileBytes = bytes.NewReader(blobSliceBytes)
			// end test code
		} else {
			DownloadStreamResponse, err := c.AzstorageUsecase.DownloadStream(ctx, blobSliceEntity.FileBlobname)
			if err != nil {
				c.log.Error(err)
				return "", err
			}
			defer DownloadStreamResponse.Body.Close()

			blobSliceBytes, err := io.ReadAll(DownloadStreamResponse.Body)
			if err != nil {
				c.log.Error("DownloadStreamResponse: ", err)
				return "", err
			}
			fileBytes = bytes.NewReader(blobSliceBytes)
		}

		pagePdfFilePath, err := c.ReportPage(fileBytes, pages[k], clientCaseName, int(total))
		if err != nil {
			return "", err
		}
		pdfs = append(pdfs, TStruct{
			PdfFilePath:        pagePdfFilePath,
			HaReportPageEntity: pages[k],
			BlobSliceEntity:    blobSliceEntity,
		})
	}
	defer func() {
		for _, v := range pdfs {
			os.Remove(v.PdfFilePath)
		}
	}()

	pdfFilePath = uuid.UuidWithoutStrike() + ".pdf"

	file, err := os.Create(pdfFilePath)
	if err != nil {
		return "", err
	}
	defer func() {
		file.Close()
	}()
	var tempFiles []string
	for _, v := range pdfs {
		tempFiles = append(tempFiles, v.PdfFilePath)
	}
	err = api.Merge("", tempFiles, file, nil, false)
	if err != nil {
		return "", err
	}
	return pdfFilePath, nil
}

func (c *HaReportPdfUsecase) ReportPage(medicalPagePdf io.ReadSeeker, haReportPageEntity *HaReportPageEntity, clientCaseName string, totalPage int) (pdfFilePath string, err error) {

	if haReportPageEntity == nil {
		return "", errors.New("haReportPageEntity is nil")
	}

	tmpPdf := uuid.UuidWithoutStrike() + ".pdf"
	file, err := os.Create(tmpPdf)
	if err != nil {
		return "", err
	}
	defer func() {
		file.Close()
		os.Remove(tmpPdf)
	}()

	_, err = io.Copy(file, medicalPagePdf)
	if err != nil {
		return "", err
	}

	images, err := ToImages(tmpPdf, ".")
	if err != nil {
		return "", err
	}
	if len(images) < 1 {
		return "", errors.New("images length is wrong")
	}

	defer func() {
		os.Remove(images[0])
	}()
	blobSliceEntity, err := c.BlobSliceUsecase.GetByGid(haReportPageEntity.BlobSliceGid)
	if err != nil {
		return "", err
	}
	if blobSliceEntity == nil {
		return "", errors.New("blobSliceEntity is nil")
	}
	pageOfSource := fmt.Sprintf("%s / %d", blobSliceEntity.SliceId, totalPage)
	dimCustom := types.Dim{
		Width:  500,
		Height: 700,
	}
	gopdf, err := c.CreateGoPdf(dimCustom, images[0], haReportPageEntity.GetAiReportToPdfText(), clientCaseName, pageOfSource)
	if err != nil {
		return "", err
	}

	fileName := uuid.UuidWithoutStrike() + ".pdf"
	err = gopdf.WritePdf(fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
	//defer func() {
	//	os.Remove(fileName)
	//}()
	//
	//desc := "scale:1 abs, pos:tl, rot:0" // v1
	//desc = "scale:0.5 abs, rot:0, pos:tl"
	//wm, err := pdfcpu.ParsePDFWatermarkDetails(tmpPdf+":1", desc, false, types.POINTS)
	//if err != nil {
	//	c.log.Error(err)
	//	return "", err
	//}
	//pdfFilePath = uuid.UuidWithoutStrike() + ".pdf"
	//cmd := cli.AddWatermarksCommand(fileName, pdfFilePath, nil, wm, nil)
	//if _, err := cli.Process(cmd); err != nil {
	//	return "", err
	//}
	//return pdfFilePath, nil
}

func (c *HaReportPdfUsecase) ReportPageBackup(medicalPagePdf io.ReadSeeker, haReportPageEntity *HaReportPageEntity, clientCaseName string, totalPage int) (pdfFilePath string, err error) {

	if haReportPageEntity == nil {
		return "", errors.New("haReportPageEntity is nil")
	}

	tmpPdf := uuid.UuidWithoutStrike() + ".pdf"
	file, err := os.Create(tmpPdf)
	if err != nil {
		return "", err
	}
	defer func() {
		file.Close()
		//os.Remove(tmpPdf)
	}()

	_, err = io.Copy(file, medicalPagePdf)
	if err != nil {
		return "", err
	}

	pdfInfo, err := c.PdfUsecase.PdfInfo(medicalPagePdf, "medicalPagePdf.pdf")
	c.log.Debug("pdfInfo.Dimensions: ", pdfInfo.Dimensions)
	if err != nil {
		return "", err
	}
	if len(pdfInfo.Dimensions) <= 0 {
		return "", errors.New("pdfInfo.Dimensions length is wrong")
	}
	//

	blobSliceEntity, err := c.BlobSliceUsecase.GetByGid(haReportPageEntity.BlobSliceGid)
	if err != nil {
		return "", err
	}
	if blobSliceEntity == nil {
		return "", errors.New("blobSliceEntity is nil")
	}
	pageOfSource := fmt.Sprintf("%s / %d", blobSliceEntity.SliceId, totalPage)
	gopdf, err := c.CreateGoPdf(pdfInfo.Dimensions[0], "", haReportPageEntity.GetAiReportToPdfText(), clientCaseName, pageOfSource)
	if err != nil {
		return "", err
	}

	fileName := time.Now().String() + ".pdf"

	err = gopdf.WritePdf(fileName)

	if err != nil {
		return "", err
	}
	defer func() {
		os.Remove(fileName)
	}()

	desc := "scale:1 abs, pos:tl, rot:0" // v1
	desc = "scale:0.5 abs, rot:0, pos:tl"
	wm, err := pdfcpu.ParsePDFWatermarkDetails(tmpPdf+":1", desc, false, types.POINTS)
	if err != nil {
		c.log.Error(err)
		return "", err
	}
	pdfFilePath = uuid.UuidWithoutStrike() + ".pdf"
	cmd := cli.AddWatermarksCommand(fileName, pdfFilePath, nil, wm, nil)
	if _, err := cli.Process(cmd); err != nil {
		return "", err
	}
	return pdfFilePath, nil
}

const (
	HaReportPdf_Side_Width  = float64(200)
	HaReportPdf_Side_Margin = float64(10)
	HaReportPdf_Font_Size   = 16
)

func (c *HaReportPdfUsecase) CreateGoPdf(dim types.Dim, imagesFilePath string, text string, clientCaseName string, pageOfSource string) (*gopdf.GoPdf, error) {
	pdf := &gopdf.GoPdf{}
	size := *gopdf.PageSizeA4
	size.W = lib.FloatSum(dim.Width + HaReportPdf_Side_Width)
	size.H = dim.Height
	pdf.Start(gopdf.Config{
		Unit:     gopdf.UnitPT,
		PageSize: size,
	})

	ttfName := "Arial"
	err := pdf.AddTTFFont(ttfName, c.conf.ResourcePath+"/ttf/Arial.ttf")
	if err != nil {
		return nil, err
	}
	err = pdf.SetFont(ttfName, "", 14)
	if err != nil {
		return nil, err
	}

	pdf.AddPage()

	imageWidth, imagesHeight, err := lib.GetImageDimensions(imagesFilePath)
	if err != nil {
		return nil, err
	}

	destImageWidth, destImageHeight := lib.CalDimensions(float64(imageWidth), float64(imagesHeight), dim.Width, dim.Height)

	err = pdf.Image(imagesFilePath, 0, 0, &gopdf.Rect{
		W: destImageWidth,
		H: destImageHeight,
	})
	if err != nil {
		return nil, err
	}

	lineWidth := float64(2)
	// 	43, 123, 191
	pdf.SetStrokeColor(43, 123, 191)
	pdf.SetLineWidth(lineWidth)
	pdf.Line(dim.Width, 0, dim.Width, size.H)

	err = pdf.SetFontSize(HaReportPdf_Font_Size)
	if err != nil {
		return nil, err
	}

	/*alignOption := gopdf.CellOption{
		Align:  gopdf.Center | gopdf.Middle,
		Border: gopdf.Left | gopdf.Right | gopdf.Bottom | gopdf.Top}

	// https://blog.csdn.net/weixin_43881017/article/details/112849522

	pdf.SetX(dim.Width)
	pdf.SetY(10)
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1)
	SourcePageNumber := &gopdf.Rect{
		W: 200,
		H: 50,
	}
	pdf.CellWithOption(SourcePageNumber, "Source page number", alignOption)*/

	bgColor := &PdfDrawColor{
		R: 43,
		G: 123,
		B: 191,
	}

	fontSizeTitle := float64(14)
	fontSizeContent := float64(11)
	sideBeginX := dim.Width + lineWidth/2

	PdfDrawText(pdf, sideBeginX, 20, HaReportPdf_Side_Width, 40,
		"Client Case Name", PdfDrawConfig{
			BgColor:  bgColor,
			FontSize: fontSizeTitle,
			FontColor: &PdfDrawColor{
				R: 255,
				G: 255,
				B: 255,
			},
		})

	PdfDrawText(pdf, sideBeginX, pdf.GetY(), HaReportPdf_Side_Width, 40,
		clientCaseName, PdfDrawConfig{
			FontSize: fontSizeContent,
			FontColor: &PdfDrawColor{
				R: 0,
				G: 0,
				B: 0,
			},
		})
	PdfDrawText(pdf, sideBeginX, pdf.GetY(), HaReportPdf_Side_Width, 40,
		"Page of The Source", PdfDrawConfig{
			BgColor:  bgColor,
			FontSize: fontSizeTitle,
			FontColor: &PdfDrawColor{
				R: 255,
				G: 255,
				B: 255,
			},
		})
	PdfDrawText(pdf, sideBeginX, pdf.GetY(), HaReportPdf_Side_Width, 40,
		pageOfSource, PdfDrawConfig{
			FontSize: fontSizeContent,
			FontColor: &PdfDrawColor{
				R: 0,
				G: 0,
				B: 0,
			},
		})
	PdfDrawText(pdf, sideBeginX, pdf.GetY(), HaReportPdf_Side_Width, 40,
		"Disease Names", PdfDrawConfig{
			BgColor:  bgColor,
			FontSize: fontSizeTitle,
			FontColor: &PdfDrawColor{
				R: 255,
				G: 255,
				B: 255,
			},
		})

	PdfDrawAutoscalingHeightText(pdf,
		sideBeginX, pdf.GetY(), HaReportPdf_Side_Width,
		text,
		PdfDrawConfig{
			FontSize: fontSizeContent,
			Padding: PdfDrawPadding{
				Left:   10,
				Top:    10,
				Right:  10,
				Bottom: 10,
			},
		})

	return pdf, err
}

type PdfDrawColor struct {
	R uint8
	G uint8
	B uint8
}

type PdfDrawPadding struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type PdfDrawBorder struct {
	Size  float64
	Color PdfDrawColor
}

type PdfDrawConfig struct {
	BgColor   *PdfDrawColor
	FontColor *PdfDrawColor
	FontSize  float64
	//CellOption *gopdf.CellOption
	TextAlign int
	Padding   PdfDrawPadding
	Border    *PdfDrawBorder
}

func PdfDrawUpdateFontSize(pdf *gopdf.GoPdf, pdfDrawConfig PdfDrawConfig) {
	if pdfDrawConfig.FontSize > 0.0 {
		pdf.SetFontSize(pdfDrawConfig.FontSize)
	} else {
		pdf.SetFontSize(16)
	}
}

func PdfDrawAutoscalingHeightText(pdf *gopdf.GoPdf, x float64, y float64, width float64, text string, pdfDrawConfig PdfDrawConfig) (boxHeight float64, err error) {

	originX := x
	originY := y
	originWidth := width
	PdfDrawUpdateFontSize(pdf, pdfDrawConfig)

	borderSize := float64(0)
	if pdfDrawConfig.Border == nil {
		pdf.SetStrokeColor(0, 0, 0)
		pdf.SetLineWidth(0)
	} else {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.SetStrokeColor(pdfDrawConfig.Border.Color.R, pdfDrawConfig.Border.Color.G, pdfDrawConfig.Border.Color.B)
		pdf.SetLineWidth(pdfDrawConfig.Border.Size)
		borderSize = pdfDrawConfig.Border.Size
	}

	//
	//pdf.MultiCell(&gopdf.Rect{
	//	W: width,
	//	H: 50,
	//}, "fdsafsadfsfsdafsdafdscsafdsfdasfasfdsaf fdsafsadfsfsdafsdafdscsafdsfdasfasfdsaf fdsafsadfsfsdafsdafdscsafdsfdasfasfdsaf fdsafsadfsfsdafsdafdscsafdsfdasfasfdsaf")
	//
	//return 0, nil

	width = width - pdfDrawConfig.Padding.Left - pdfDrawConfig.Padding.Right - borderSize*2

	texts, err := pdf.SplitText(text, width) // remove left and right margins
	if err != nil {
		return 0.0, err
	}

	//var x float64 = HaReportPdf_Side_Margin
	//var y float64 = HaReportPdf_Side_Margin + 20 + pdf.GetY() + DiseaseNamesRect.H
	//pdf.SetXY(x, y)

	defaultHeight := float64(20)
	//y += defaultHeight

	height := originY + float64(len(texts))*defaultHeight
	if pdfDrawConfig.BgColor != nil {
		pdf.SetFillColor(pdfDrawConfig.BgColor.R, pdfDrawConfig.BgColor.G, pdfDrawConfig.BgColor.B)
		// D, F, DF, FD
		// F: 只有背景
		// DF：有边框和背景
		err = pdf.Rectangle(x, originY, x+originWidth, height+borderSize, "F", 0, 0)
		if err != nil {
			fmt.Println(err)
		}
		//
	}

	pdf.SetTextColor(0, 0, 0)
	if pdfDrawConfig.FontColor != nil {
		pdf.SetTextColor(pdfDrawConfig.FontColor.R, pdfDrawConfig.FontColor.G, pdfDrawConfig.FontColor.B)
	}

	var alignOption *gopdf.CellOption
	//if pdfDrawConfig != nil {
	//	alignOption = pdfDrawConfig.CellOption
	//} else {
	alignOption = &gopdf.CellOption{
		Align: pdfDrawConfig.TextAlign,
		//Border: gopdf.Left | gopdf.Right | gopdf.Bottom | gopdf.Top, // 开启后有边框
	}
	//}
	// gopdf.Center | gopdf.Middle
	if pdfDrawConfig.TextAlign == 0 {
		alignOption.Align = gopdf.Left
	}

	//lib.DPrintln("sss borderSize:", borderSize)
	y = y + pdfDrawConfig.Padding.Top + borderSize

	pdf.SetMarginBottom(20)
	pdf.SetMarginTop(20)
	for _, t := range texts {
		pdf.SetNewY(y, defaultHeight) // More than 1 page will be automatically paginated. 20 is a user-defined row height, not a fixed value.
		pdf.SetX(x + pdfDrawConfig.Padding.Left + borderSize)

		// 新的一页 todo:page size
		//if pdf.GetY()+defaultHeight > pdf.-pdf.MarginBottom() {
		//
		//}

		y = pdf.GetY()
		err = pdf.CellWithOption(&gopdf.Rect{W: width,
			H: defaultHeight,
		}, t, *alignOption)
		//err = pdf.Text(t)
		if err != nil {
			continue
		}
		y += defaultHeight
	}
	pdf.SetX(x - pdfDrawConfig.Padding.Left)
	pdf.SetY(y + pdfDrawConfig.Padding.Bottom + borderSize)

	if pdfDrawConfig.Border != nil {
		var points = []gopdf.Point{
			{
				X: originX + pdfDrawConfig.Border.Size/2,
				Y: originY + pdfDrawConfig.Border.Size/2,
			},

			{
				X: originX + originWidth - pdfDrawConfig.Border.Size/2,
				Y: originY + pdfDrawConfig.Border.Size/2,
			},
			{
				X: originX + originWidth - pdfDrawConfig.Border.Size/2,
				Y: y + pdfDrawConfig.Padding.Bottom + pdfDrawConfig.Border.Size/2,
			},
			{
				X: originX + pdfDrawConfig.Border.Size/2,
				Y: y + pdfDrawConfig.Padding.Bottom + pdfDrawConfig.Border.Size/2,
			},
		}
		pdf.SetStrokeColor(pdfDrawConfig.Border.Color.R, pdfDrawConfig.Border.Color.G, pdfDrawConfig.Border.Color.B)
		pdf.SetLineWidth(pdfDrawConfig.Border.Size)
		pdf.Polygon(points, "")
	}

	return y, nil
}

func PdfDrawText(pdf *gopdf.GoPdf, x float64, y float64, width float64, height float64, text string, pdfDrawConfig PdfDrawConfig) (boxHeight float64) {

	originX := x
	originY := y
	originWidth := width
	originHeight := height
	//var alignOption *gopdf.CellOption
	//alignOption = &gopdf.CellOption{
	//	//Align: gopdf.Center | gopdf.Middle,
	//	//Border: gopdf.Left | gopdf.Right | gopdf.Bottom | gopdf.Top,
	//}
	alignOption := &gopdf.CellOption{
		Align: pdfDrawConfig.TextAlign,
		//Border: gopdf.Left | gopdf.Right | gopdf.Bottom | gopdf.Top,
	}
	if pdfDrawConfig.TextAlign == 0 {
		alignOption.Align = gopdf.Center | gopdf.Middle
	}

	pdf.SetX(x)
	pdf.SetY(y)
	borderSize := float64(0)
	if pdfDrawConfig.Border == nil {
		pdf.SetStrokeColor(0, 0, 0)
		pdf.SetLineWidth(0)
	} else {
		pdf.SetStrokeColor(pdfDrawConfig.Border.Color.R, pdfDrawConfig.Border.Color.G, pdfDrawConfig.Border.Color.B)
		pdf.SetLineWidth(pdfDrawConfig.Border.Size)
		borderSize = pdfDrawConfig.Border.Size
	}

	// https://blog.csdn.net/weixin_43881017/article/details/112849522
	PdfDrawUpdateFontSize(pdf, pdfDrawConfig)

	// 设置背景
	if pdfDrawConfig.BgColor != nil {
		pdf.SetFillColor(pdfDrawConfig.BgColor.R, pdfDrawConfig.BgColor.G, pdfDrawConfig.BgColor.B)
		err := pdf.Rectangle(x+borderSize/2, y+borderSize/2, x+width-borderSize/2, y+height-borderSize/2, "F", 0, 0)
		if err != nil {
			fmt.Println(err)
		}
	}

	pdf.SetTextColor(0, 0, 0)
	if pdfDrawConfig.FontColor != nil {
		pdf.SetTextColor(pdfDrawConfig.FontColor.R, pdfDrawConfig.FontColor.G, pdfDrawConfig.FontColor.B)
	}

	pdf.SetX(x + pdfDrawConfig.Padding.Left + borderSize)
	pdf.SetY(y + pdfDrawConfig.Padding.Top + borderSize)
	rect := &gopdf.Rect{
		W: width - pdfDrawConfig.Padding.Left - pdfDrawConfig.Padding.Right - borderSize*2,
		H: height - pdfDrawConfig.Padding.Top - pdfDrawConfig.Padding.Bottom - borderSize*2,
	}
	pdf.CellWithOption(rect, text, *alignOption)

	if pdfDrawConfig.Border != nil {
		var points = []gopdf.Point{
			{
				X: originX + borderSize/2,
				Y: originY + borderSize/2,
			},

			{
				X: originX + originWidth - borderSize/2,
				Y: originY + borderSize/2,
			},
			{
				X: originX + originWidth - borderSize/2,
				Y: originY + originHeight - borderSize/2,
			},
			{
				X: originX + borderSize/2,
				Y: originY + originHeight - borderSize/2,
			},
		}
		pdf.SetStrokeColor(pdfDrawConfig.Border.Color.R, pdfDrawConfig.Border.Color.G, pdfDrawConfig.Border.Color.B)
		pdf.SetLineWidth(borderSize)
		pdf.Polygon(points, "")
	}

	pdf.SetY(originY + height)
	pdf.SetX(originX)

	return height
}
