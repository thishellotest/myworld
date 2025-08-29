package tests

import (
	"fmt"
	"io"
	"os"
	"testing"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_PdfcpuUsecase_SplitPdf(t *testing.T) {
	sourcePdf, err := os.Open("./res/c.pdf")
	if err != nil {
		panic(err)
	}
	defer sourcePdf.Close()
	result, err := UT.PdfcpuUsecase.SplitPdfAndCombine(sourcePdf, []*biz.SplitPdfAndCombineConf{
		{
			PageBegin: 1,
			PageEnd:   2,
		},
		{
			PageBegin: 2,
			PageEnd:   2,
		},
		{
			PageBegin: 1,
			PageEnd:   3,
		},
	})

	lib.DPrintln(err, len(result))
	for k, v := range result {
		file := fmt.Sprintf("./tmp/%d.pdf", k)
		fw, err := os.Create(file)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(fw, v)
		if err != nil {
			panic(err)
		}
	}
}

func Test_sss(t *testing.T) {
	a, err := os.MkdirTemp(configs.GetAppRuntimePath()+"", "aaac")
	lib.DPrintln(a, err)
}

func Test_PdfcpuUsecase_FormLists(t *testing.T) {
	UT.PdfcpuUsecase.FormLists()
}

func Test_PdfcpuUsecase_Fonts(t *testing.T) {
	UT.PdfcpuUsecase.Fonts()
}

func Test_PdfcpuUsecase_FillForm(t *testing.T) {
	UT.PdfcpuUsecase.FillForm()
}
