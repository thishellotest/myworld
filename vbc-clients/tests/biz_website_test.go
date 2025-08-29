package tests

import (
	"testing"
	"vbc/lib"
)

func Test_WebsiteUsecase_BizSyncToZohoOrVBCRM(t *testing.T) {
	//err := UT.WebsiteUsecase.BizSyncToZohoOrVBCRM("Jerry", "Reed", "j.reednfl@yahoo.com", "3216841000", "SC", "")
	//lib.DPrintln(err)
	//
	err := UT.WebsiteUsecase.BizSyncToZohoOrVBCRM("Austin", "Smith", "01AustinHaley10@gmail.com", "8087405202", "HI", "", "", "", "")
	lib.DPrintln(err)
}

func Test_WebsiteUsecase_SyncToZohoOrVBCRM(t *testing.T) {
	formData := `{"data":{"formName":"VBC Contact Form","field:comp-lozg1qa5":"18600503374","field:comp-lozdvdku1":"I already have a rating, but I think I could be underrated.","submissionTime":"2024-02-05T02:27:27.993Z","context":{"metaSiteId":"3130ca58-a60b-481b-9a48-1a0c9155304c","activationId":"9467050c-53aa-4c8d-9886-6b860bfa1049"},"_context":{"activation":{"id":"9467050c-53aa-4c8d-9886-6b860bfa1049"},"configuration":{"id":"0de7ae89-9781-46aa-a0c0-df74cc10a1fe"},"app":{"id":"14ce1214-b278-a7e4-1373-00cebd1bef7c"},"action":{"id":"00cb17d8-0448-423e-0c73-57d4c506dc45"},"trigger":{"key":"wix_forms-form_submit"}},"contact":{"name":{"first":"TestGary","last":"TestLiao"},"email":"lialing@foxmail.com","locale":"zh","phones":[{"tag":"UNTAGGED","formattedPhone":"18600503374","id":"3d0c071b-0b70-4595-b73d-4f42394aa147","primary":true,"phone":"18600503374"}],"emails":[{"id":"c85f5669-54e6-468a-8626-5a036e29fb71","tag":"MAIN","email":"lialing@foxmail.com","primary":true}],"phone":"18600503374"},"submissionId":"f2538757-c446-4a3a-85eb-af2226a3406f","field:comp-lozdvdjr":"lialing@foxmail.com","contactId":"acb6b1c0-fb0c-4b11-ae8e-9e3dc391df31","field:comp-lozdvdjo2":"TestLiao","field:comp-lozegck4":"Checked","field:comp-loze9sgy":"AL","field:comp-lozdvdje":"TestGary","formId":"comp-lozdvdj82"}}`
	err := UT.WebsiteUsecase.SyncToZohoOrVBCRM(formData)
	lib.DPrintln(err)
}
