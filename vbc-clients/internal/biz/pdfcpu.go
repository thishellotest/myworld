package biz

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/uuid"
)

type PdfcpuUsecase struct {
	log             *log.Helper
	CommonUsecase   *CommonUsecase
	conf            *conf.Data
	ResourceUsecase *ResourceUsecase
}

func NewPdfcpuUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ResourceUsecase *ResourceUsecase) *PdfcpuUsecase {
	uc := &PdfcpuUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		ResourceUsecase: ResourceUsecase,
	}
	return uc
}

type SplitPdfAndCombineConf struct {
	PageBegin int // Begin with page 1
	PageEnd   int
}

func (c *PdfcpuUsecase) SplitPdfAndCombine(sourcePdf io.ReadSeeker, conf []*SplitPdfAndCombineConf) (result []io.Reader, err error) {

	pdfUuid := uuid.UuidWithoutStrike()
	if sourcePdf == nil {
		return nil, errors.New("sourcePdf is nil")
	}
	if conf == nil {
		return nil, errors.New("conf is nil")
	}

	tempDir, err := os.MkdirTemp(configs.GetAppRuntimePath(), pdfUuid)
	if err != nil {
		return nil, err
	}
	lib.DPrintln("tempDir:", tempDir)
	defer func() {
		os.RemoveAll(tempDir)
	}()

	err = pdfcpu.Split(sourcePdf, tempDir, pdfUuid, 1, nil)
	if err != nil {
		return nil, err
	}
	lib.DPrintln("conf:", conf)

	for _, v := range conf {
		//file, err := os.Create("../pdf/combine2.pdf")
		//err = pdfcpu.Merge("", []string{"../pdf/new_c_2.pdf", "../pdf/new_c_3.pdf"}, file, nil, false)
		//lib.DPrintln(err)

		var pdfs []string
		for i := v.PageBegin; i <= v.PageEnd; i++ {
			subPdf := fmt.Sprintf("%s/%s_%d.pdf", tempDir, pdfUuid, i)
			exists, err := lib.PathExists(subPdf)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
				//return nil, errors.New(subPdf + " does not exist. ")
			}
			pdfs = append(pdfs, subPdf)
		}
		if len(pdfs) == 0 {
			return nil, errors.New("pdfs length is 0")
		}

		fileBytes := make([]byte, 100)
		file := bytes.NewBuffer(fileBytes)
		err = pdfcpu.Merge("", pdfs, file, nil, false)
		if err != nil {
			return nil, err
		}
		result = append(result, file)
	}

	return result, nil
}

func (c *PdfcpuUsecase) FormLists() {

	pdf := c.ResourceUsecase.ResPath() + "/vba-21-4138-are.pdf"
	pdfFile, err := os.Open(pdf)
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()
	fields, err := pdfcpu.FormFields(pdfFile, nil)
	lib.DPrintln(fields)
	lib.DPrintln(err)
}

func (c *PdfcpuUsecase) Fonts() {

	pdf := c.ResourceUsecase.ResPath() + "/vba-21-4138-are-v1.pdf"
	pdfFile, err := os.Open(pdf)
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()
	err = pdfcpu.ExtractFonts(pdfFile, "/tmp", "b", nil, nil)
	lib.DPrintln(err)
}

func (c *PdfcpuUsecase) FillForm() {

	inPdf := c.ResourceUsecase.ResPath() + "/vba-21-4138-are-v1.pdf"
	outPdf := "/tmp/outpdf.pdf"
	inFileJSON := `{
		"forms":[
			{"textfield":[
				{
					"name":"form1[0].#subform[0].TelephoneNumber_LastFourNumbers[0]",
					"value":"1234",
					"locked": false
				}
			]}
		]}`
	//inFileJSON = `{}`
	inPdfFile, err := os.Open(inPdf)

	f2, err := os.Create(outPdf)
	rd := strings.NewReader(inFileJSON)

	// 打印当前的 UserConfigDir
	configDir, _ := os.UserConfigDir()
	fmt.Println("UserConfigDir:", configDir)

	config := model.NewDefaultConfiguration()
	//config.SetUnit("aaa")
	err = pdfcpu.FillForm(inPdfFile, rd, f2, config)
	lib.DPrintln(err)
}
