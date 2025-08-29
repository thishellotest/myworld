package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"io"
	"math"
	"strconv"
	"vbc/internal/conf"
)

type PdfUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewPdfUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *PdfUsecase {
	uc := &PdfUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	return uc
}

func PdfcpuJsonInfo(info *pdfcpu.PDFInfo, pages types.IntSet) (map[string]model.PageBoundaries, []types.Dim) {
	if len(pages) > 0 {
		pbs := map[string]model.PageBoundaries{}
		for i, pb := range info.PageBoundaries {
			if _, found := pages[i+1]; !found {
				continue
			}
			d := pb.CropBox().Dimensions()
			if pb.Rot%180 != 0 {
				d.Width, d.Height = d.Height, d.Width
			}
			pb.Orientation = "portrait"
			if d.Landscape() {
				pb.Orientation = "landscape"
			}
			if pb.Media != nil {
				pb.Media.Rect = pb.Media.Rect.ConvertToUnit(info.Unit)
				pb.Media.Rect.LL.X = math.Round(pb.Media.Rect.LL.X*100) / 100
				pb.Media.Rect.LL.Y = math.Round(pb.Media.Rect.LL.Y*100) / 100
				pb.Media.Rect.UR.X = math.Round(pb.Media.Rect.UR.X*100) / 100
				pb.Media.Rect.UR.Y = math.Round(pb.Media.Rect.UR.Y*100) / 100
			}
			if pb.Crop != nil {
				pb.Crop.Rect = pb.Crop.Rect.ConvertToUnit(info.Unit)
				pb.Crop.Rect.LL.X = math.Round(pb.Crop.Rect.LL.X*100) / 100
				pb.Crop.Rect.LL.Y = math.Round(pb.Crop.Rect.LL.Y*100) / 100
				pb.Crop.Rect.UR.X = math.Round(pb.Crop.Rect.UR.X*100) / 100
				pb.Crop.Rect.UR.Y = math.Round(pb.Crop.Rect.UR.Y*100) / 100
			}
			if pb.Trim != nil {
				pb.Trim.Rect = pb.Trim.Rect.ConvertToUnit(info.Unit)
				pb.Trim.Rect.LL.X = math.Round(pb.Trim.Rect.LL.X*100) / 100
				pb.Trim.Rect.LL.Y = math.Round(pb.Trim.Rect.LL.Y*100) / 100
				pb.Trim.Rect.UR.X = math.Round(pb.Trim.Rect.UR.X*100) / 100
				pb.Trim.Rect.UR.Y = math.Round(pb.Trim.Rect.UR.Y*100) / 100
			}
			if pb.Bleed != nil {
				pb.Bleed.Rect = pb.Bleed.Rect.ConvertToUnit(info.Unit)
				pb.Bleed.Rect.LL.X = math.Round(pb.Bleed.Rect.LL.X*100) / 100
				pb.Bleed.Rect.LL.Y = math.Round(pb.Bleed.Rect.LL.Y*100) / 100
				pb.Bleed.Rect.UR.X = math.Round(pb.Bleed.Rect.UR.X*100) / 100
				pb.Bleed.Rect.UR.Y = math.Round(pb.Bleed.Rect.UR.Y*100) / 100
			}
			if pb.Art != nil {
				pb.Art.Rect = pb.Art.Rect.ConvertToUnit(info.Unit)
				pb.Art.Rect.LL.X = math.Round(pb.Art.Rect.LL.X*100) / 100
				pb.Art.Rect.LL.Y = math.Round(pb.Art.Rect.LL.Y*100) / 100
				pb.Art.Rect.UR.X = math.Round(pb.Art.Rect.UR.X*100) / 100
				pb.Art.Rect.UR.Y = math.Round(pb.Art.Rect.UR.Y*100) / 100
			}
			pbs[strconv.Itoa(i+1)] = pb
		}
		return pbs, nil
	}

	var dims []types.Dim
	for k, v := range info.PageDimensions {
		if v {
			dc := k.ConvertToUnit(info.Unit)
			dc.Width = math.Round(dc.Width*100) / 100
			dc.Height = math.Round(dc.Height*100) / 100
			dims = append(dims, dc)
		}
	}
	return nil, dims
}

func (c *PdfUsecase) PdfInfo(fileReadSeeker io.ReadSeeker, fileName string) (pdfInfo pdfcpu.PDFInfo, err error) {

	info, err := api.PDFInfo(fileReadSeeker, fileName, nil, nil)
	if err != nil {
		return pdfInfo, err
	}
	pages, err := api.PagesForPageSelection(info.PageCount, nil, false, false)
	if err != nil {
		return pdfInfo, err
	}
	//lib.DPrintln(pages)
	info.Boundaries, info.Dimensions = PdfcpuJsonInfo(info, pages)
	//lib.DPrintln(info)
	//lib.DPrintln(info.PageCount)
	//lib.DPrintln("info.Boundaries:", info.Boundaries)
	//lib.DPrintln("info.Dimensions:", info.Dimensions)
	return *info, nil
}
