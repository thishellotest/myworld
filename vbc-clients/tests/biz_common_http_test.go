package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_CommonHttpUsecase_t(t *testing.T) {
	body := lib.ToTypeMapByString(`{"uniqid":"Tinnitus","jotform_ids":[{"value":"240905843097159","label":"Back"},{"value":"240917546196162","label":"Hearing Loss and Tinnitus"},{"value":"240915419171152","label":"Knee"}]}`)
	a, err := UT.CommonHttpUsecase.BizSave(biz.UserFacade{}, biz.CommonHttp_CommonType_ConditionQuestionnaires, body)
	lib.DPrintln(a, err)
}
