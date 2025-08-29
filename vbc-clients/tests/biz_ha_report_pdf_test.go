package tests

import (
	"os"
	"testing"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_HaReportPdfUsecase_ReportPage(t *testing.T) {

	pdf := "./tmp/1_4636f3657f8043b1bf25d3fefc8bc990.pdf"
	//pdf = "./tmp/STR Full_2.pdf"
	//pdf = "./tmp/1_4_1.pdf"
	//pdf = "./tmp/1_4.pdf"
	file, err := os.Open(pdf)

	if err != nil {
		panic(err)
	}

	entity, _ := UT.HaReportPageUsecase.GetByCond(Eq{"id": 82})

	cc, err := UT.HaReportPdfUsecase.ReportPage(file, entity, "Augustus Ivan IV Goodson-70#5124", 231)
	lib.DPrintln(cc, err)
}

//func Test_HaReportPdfUsecase_CreateHaReportPdf(t *testing.T) {
//
//	entity, _ := UT.HaReportTaskUsecase.GetByCond(Eq{"gid": "g1"})
//	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, entity.ClientCaseId)
//	aa, err := UT.HaReportPdfUsecase.CreateHaReportPdf(context.TODO(), entity, tCase)
//	lib.DPrintln(aa, err)
//}
