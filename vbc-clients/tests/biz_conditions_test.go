package tests

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_ConditionUsecase_ConditionUpsert(t *testing.T) {
	a, err := UT.ConditionUsecase.ConditionUpsert("abc")
	lib.DPrintln(a, err)
}

func Test_ConditionUsecase_Upsert(t *testing.T) {
	err := UT.ConditionUsecase.Upsert("abc")
	lib.DPrintln(err)
}

func Test_ConditionUsecase_Import(t *testing.T) {

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
		categoryName := ""
		conditionName := ""
		for k, colCell := range row {
			if k == 1 && colCell != "" && colCell != "Category (Jotform)" {
				categoryName = colCell
			} else if k == 2 && colCell != "" && colCell != "Condition (Claim)" {
				conditionName = colCell
			}
		}
		if categoryName != "" && conditionName != "" {

			categoryEntity, err := UT.ConditionCategoryUsecase.Upsert(categoryName)
			if err != nil {
				panic(err)
			}
			primaryEntity, _, err := UT.ConditionUsecase.UpsertPrimaryCondition(conditionName, categoryEntity.ID)
			if err != nil {
				panic(err)
			}
			UT.ConditionLogAiUsecase.AddLogConditionSourceFromImport(primaryEntity.ID, "VBC Conditions.xlsx", "")
			UT.ConditionLogAiUsecase.AddCategoryOfConditionSource(primaryEntity.ID, categoryEntity.ID, biz.ConditionLogAi_FromType_Import, "VBC Conditions.xlsx", "")
			lib.DPrintln("categoryName:", categoryName, "conditionName:", conditionName)
		}
	}
}
