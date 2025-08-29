package tests

import (
	"strconv"
	"strings"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_MiscUsecase_ClientThingsToKnowExamFileId(t *testing.T) {
	thingsToKnowExamFileId, err := UT.MiscUsecase.ClientThingsToKnowExamFileId("257067526990")
	lib.DPrintln(thingsToKnowExamFileId, err)
}

func Test_MiscUsecase_HandleMiscThingsToKnowCPExam(t *testing.T) {
	err := UT.MiscUsecase.HandleMiscThingsToKnowCPExam(5465)
	lib.DPrintln(err)
}

func Test_MiscUsecase_HandleRemoveMiscThingsToKnowCPExam(t *testing.T) {
	err := UT.MiscUsecase.HandleRemoveMiscThingsToKnowCPExam(5492)
	lib.DPrintln(err)
}

func Test_MiscUsecase_GuideFileByApi(t *testing.T) {
	info, err := UT.MiscUsecase.GuideFileByApi("281846152333")
	lib.DPrintln(err)
	lib.DPrintln(info)
}

func Test_MiscUsecase_UpdateGuideForClientCase(t *testing.T) {
	tCase, err := UT.TUsecase.DataById(biz.Kind_client_cases, 5217)
	lib.DPrintln(err)
	err = UT.MiscUsecase.UpdateGuideForClientCase(tCase)
	lib.DPrintln(err)
}

func Test_MiscUsecase_UpdateGuideForClientCase_UpdateAll(t *testing.T) {

	list, err := UT.TUsecase.ListByCond(biz.Kind_client_cases, builder.Eq{"deleted_at": 0, "biz_deleted_at": 0})
	if err != nil {
		lib.DPrintln(err)
	}
	for k, v := range list {
		key := biz.MapKeyClientBoxFolderId(v.Id())
		a, _ := UT.MapUsecase.GetForString(key)
		if a != "" {
			err := UT.MiscUsecase.UpdateGuideForClientCase(list[k])
			if err != nil {
				UT.LogUsecase.SaveLog(v.Id(), "MiscUsecase_UpdateAll_Error", map[string]interface{}{})
				lib.DPrintln("updated error:", v.Id())
			} else {
				UT.LogUsecase.SaveLog(v.Id(), "MiscUsecase_UpdateAll_Ok", map[string]interface{}{})
				lib.DPrintln("updated ok:", v.Id())
			}
		}
	}
}

func Test_MiscUsecase_HandleMoving2122aFile(t *testing.T) {

	str := "5778,5781,5784,5788,5789,5791,5794,5810,5815,5816,5818,5819,5820,5821,5822"
	strs := strings.Split(str, ",")

	for _, v := range strs {
		a, _ := strconv.ParseInt(v, 10, 32)
		_, err := UT.MiscUsecase.HandleMoving2122aFile(int32(a))
		if err != nil {
			lib.DPrintln(err, "caseID:", v)
		}

	}
}

func Test_MiscUsecase_DoHandleMoving2122aFile(t *testing.T) {

	a, err := UT.MiscUsecase.DoHandleMoving2122aFile(5829)
	lib.DPrintln(a, err)
}
