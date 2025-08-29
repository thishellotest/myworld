package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_vbc"
	"vbc/lib"
)

func Test_PricingVersionUsecase_CurrentVersion(t *testing.T) {
	a, err := UT.PricingVersionUsecase.CurrentVersion()
	lib.DPrintln(err)
	lib.DPrintln(a)
}
func Test_PricingVersionUsecase_CurrentVersionConfig(t *testing.T) {
	c, _, err := UT.PricingVersionUsecase.CurrentVersionConfig()
	lib.DPrintln(err)
	index := 90
	dd := c.GetByIndex(index)
	lib.DPrintln(c.GetByIndex(index))
	ca := config_vbc.FeeDefine.GetByIndex(index)
	if biz.InterfaceToString(dd) == biz.InterfaceToString(ca) {
		lib.DPrintln("ok")
	} else {
		panic("error")
	}
}

func Test_PricingVersionUsecase_CurrentVersionConfig1(t *testing.T) {
	c, _, err := UT.PricingVersionUsecase.CurrentVersionConfig()
	lib.DPrintln(err)
	res, err := c.GetBoxSignTpl("90")
	lib.DPrintln(res, err)
}

func Test_PricingVersionUsecase_ConfigByPricingVersion(t *testing.T) {
	c, _, err := UT.PricingVersionUsecase.ConfigByPricingVersion(biz.DefaultPricingVersion)
	lib.DPrintln(err)
	res, err := c.GetBoxSignTpl("10")
	lib.DPrintln(res, err)
}
