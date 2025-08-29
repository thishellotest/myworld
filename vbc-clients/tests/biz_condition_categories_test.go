package tests

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
	"vbc/lib"
)

func Test_ConditionCategoryUsecase_Upsert(t *testing.T) {
	a, err := UT.ConditionCategoryUsecase.Upsert("aa2")
	lib.DPrintln(a, err)
}

func Test_ConditionCategoryUsecase_HandleExcel(t *testing.T) {
	f, err := excelize.OpenFile("./excel/VBC Conditions.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	//cell, err := f.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for k, colCell := range row {
			if k == 0 && colCell != "" && colCell != "Total Category (Jotform)" {
				fmt.Println(colCell)
				UT.ConditionCategoryUsecase.Upsert(colCell)
			}
		}
	}
}
