package tests

import (
	"testing"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/internal/utils"
	"vbc/lib"
)

func Test_CalDelayDayTime(t *testing.T) {
	aa := utils.CalDelayDayTime(time.Now(), *configs.GetVBCDefaultLocation())
	lib.DPrintln(aa.Format(time.RFC3339))
}

func Test_CalIntervalDayTime(t *testing.T) {
	//utils.CalIntervalDayTime()
	executeTime := time.Now().In(configs.GetVBCDefaultLocation())
	hour := executeTime.Hour()
	minute := executeTime.Minute()
	lib.DPrintln(hour)
	lib.DPrintln(minute)

}
func Test_GetCaseTimeLocation(t *testing.T) {
	a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	b := biz.GetCaseTimeLocation(a, UT.CommonUsecase.Log)
	lib.DPrintln(b.String())
}

func Test_EncryptSensitive(t *testing.T) {
	a, err := biz.EncryptSensitive("abc")
	lib.DPrintln(a, err)
}

func Test_DecryptSensitive(t *testing.T) {
	a, err := biz.DecryptSensitive("UwhH")
	lib.DPrintln(a, err)
}
