package biz

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/go-kratos/kratos/v2/log"
	"image/jpeg"
	"os"
	"path/filepath"
	"vbc/internal/conf"
	"vbc/lib"
)

func TestCgoFitz() {
	doc, err := fitz.New("../pdf/tmp/0_full.pdf")
	if err != nil {
		panic(err)
	}

	defer doc.Close()
}

type PdfGoFitzUsecase struct {
	log             *log.Helper
	conf            *conf.Data
	CommonUsecase   *CommonUsecase
	ResourceUsecase *ResourceUsecase
}

func NewPdfGoFitzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ResourceUsecase *ResourceUsecase,
) *PdfGoFitzUsecase {
	uc := &PdfGoFitzUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		ResourceUsecase: ResourceUsecase,
	}

	return uc
}

// PdfToImages 图片从1开始
func (c *PdfGoFitzUsecase) PdfToImages(filePath string, dir string, jpgNamePrefix string) error {

	doc, err := fitz.New(filePath)
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s_%d.jpg", jpgNamePrefix, n+1)))
		if err != nil {
			return err
		}
		defer f.Close()

		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PdfGoFitzUsecase) TestFitz() {

	file := c.ResourceUsecase.ResPath() + "/Untitled document.pdf"

	c.log.Debug("TestFitz:", file)
	// "../pdf/tmp/0_full.pdf"

	doc, err := fitz.New(file)
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	//tmpDir, err := os.MkdirTemp(os.TempDir(), "fitz")
	//if err != nil {
	//	panic(err)
	//}
	tmpDir := "/tmp"
	c.log.Debug("TestFitz:", tmpDir)
	//fmt.Println(tmpDir)

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		lib.DPrintln(n, doc.NumPage())
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(tmpDir, fmt.Sprintf("test%03d.jpg", n)))
		if err != nil {
			c.log.Debug("TestFitz:err", err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			c.log.Debug("TestFitz:err", err)
		}
		f.Close()
	}
}
