package tests

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_box"
	"vbc/lib"
)

func Test_BoxUsecase_CreateFolder(t *testing.T) {
	c, er := UT.BoxUsecase.CreateFolder("aaa bb#11", "241182558605")
	lib.DPrintln(c, er)
}

func Test_boxsign_templateTl(t *testing.T) {
	aa := lib.ToTypeMapByString(config_box.BoxSignTemplateSourceJson)
	entries := lib.ToTypeList(aa.Get("entries"))
	for _, v := range entries {
		re := regexp.MustCompile(`.*\((.*)\).*`)
		match := re.FindStringSubmatch(v.GetString("name"))
		key := fmt.Sprintf("%s%s", biz.Map_boxsignTpl, match[1])
		UT.MapUsecase.Set(key, v.GetString("id"))
	}
}

func Test_BoxUsecase_UpdateFolderName(t *testing.T) {
	res, err := UT.BoxUsecase.UpdateFolderName("247252076715", "VBC - BoxB10, BoxB_1")
	lib.DPrintln(err)
	lib.DPrintln(res)
}

func Test_BoxUsecase_MoveFolderName(t *testing.T) {
	res, err := UT.BoxUsecase.MoveFolderName("265036612209", "NewName", "264686095074")
	lib.DPrintln(err)
	lib.DPrintln(res)
}

func Test_BoxUsecase_CopyFolder(t *testing.T) {
	id1, httpCode, err := UT.BoxUsecase.CopyFolder("267683828516", "abc", "271859979627")
	lib.DPrintln(id1, httpCode)
	lib.DPrintln(err)
}

func Test_BoxUsecase_ListBoxSignTemplates(t *testing.T) {
	res, err := UT.BoxUsecase.ListBoxSignTemplates()
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_SignRequests(t *testing.T) {
	task := &biz.CreateEnvelopeTaskInput{}

	str := `{"agentEmail":"liaogling@gmail.com","agentFirstName":"Gary","agentLastName":"Liao","clientEmail":"18891706@qq.com","clientFirstName":"FirstTest","clientLastName":"LastTest","signType":"box","templateId":"86df0f08-497c-4960-9121-afbcd4f4251b"}`
	task, err := lib.StringToTE[*biz.CreateEnvelopeTaskInput](str, nil)
	lib.DPrintln(err)

	res, cId, err := UT.BoxUsecase.SignRequests(task, "249619374866")
	lib.DPrintln(cId)
	lib.DPrintln(res, err)
	// {"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Your Veteran Benefits Center Contract From API","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (team@vetbenefitscenter.com)","are_reminders_enabled":false,"signers":[{"email":"team@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null},{"email":"liaogling@gmail.com","role":"signer","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null}],"id":"b1685f64-b285-4786-abc0-174931e7be0e","prefill_tags":[{"document_tag_id":"clientName","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"VSNameSign","text_value":"Gary Liao","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"246205309773","etag":"1","type":"folder","sequence_id":"1","name":"Test Box Sign Requests"},"name":"Your Veteran Benefits Center Contract (3).pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1429599030083","etag":"0","type":"file","sequence_id":"0","name":"Your Veteran Benefits Center Contract (3).pdf","sha1":"52aed5e732c43468dc6a0575109cb84a8f09a173","file_version":{"id":"1567482284483","type":"file_version","sha1":"52aed5e732c43468dc6a0575109cb84a8f09a173"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"a0a2df1f-3fd4-42f2-8881-6d6b1df55ea2"}

	// webhook 28: {"type":"webhook_event","id":"2b4ded3a-2296-4d04-a8b7-f86005ced427","created_at":"2024-01-29T23:32:33-08:00","trigger":"SIGN_REQUEST.COMPLETED","webhook":{"id":"2383852458","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"team@vetbenefitscenter.com"},"source":{"id":"1429599030083","type":"file","file_version":{"type":"file_version","id":"1567474367100","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35"},"sequence_id":"2","etag":"2","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35","name":"Your Veteran Benefits Center Contract (3).pdf","description":"","size":452587,"path_collection":{"total_count":3,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"}]},"created_at":"2024-01-29T23:25:43-08:00","modified_at":"2024-01-29T23:32:28-08:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-01-29T23:25:43-08:00","content_modified_at":"2024-01-29T23:32:28-08:00","created_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"modified_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"},"item_status":"active"},"additional_info":{"sign_request_id":"b1685f64-b285-4786-abc0-174931e7be0e","signer_emails":["team@vetbenefitscenter.com","lialing@foxmail.com","liaogling@gmail.com"],"external_id":null}}
	// webhook 27: {"type":"webhook_event","id":"2b4ded3a-2296-4d04-a8b7-f86005ced427","created_at":"2024-01-29T23:32:33-08:00","trigger":"SIGN_REQUEST.COMPLETED","webhook":{"id":"2383696467","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"team@vetbenefitscenter.com"},"source":{"id":"1429599030083","type":"file","file_version":{"type":"file_version","id":"1567474367100","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35"},"sequence_id":"2","etag":"2","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35","name":"Your Veteran Benefits Center Contract (3).pdf","description":"","size":452587,"path_collection":{"total_count":3,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"}]},"created_at":"2024-01-29T23:25:43-08:00","modified_at":"2024-01-29T23:32:28-08:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-01-29T23:25:43-08:00","content_modified_at":"2024-01-29T23:32:28-08:00","created_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"modified_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"},"item_status":"active"},"additional_info":{"sign_request_id":"b1685f64-b285-4786-abc0-174931e7be0e","signer_emails":["team@vetbenefitscenter.com","lialing@foxmail.com","liaogling@gmail.com"],"external_id":null}}
}

//
//func Test_BoxUsecase_PatientPaymentFormSignRequests(t *testing.T) {
//	//task := &biz.CreateEnvelopeTaskInput{}
//	//
//	//str := `{"agentEmail":"liaogling@gmail.com","agentFirstName":"Gary","agentLastName":"Liao","clientEmail":"18891706@qq.com","clientFirstName":"FirstTest","clientLastName":"LastTest","signType":"box","templateId":"86df0f08-497c-4960-9121-afbcd4f4251b"}`
//	//task, err := lib.StringToTE[*biz.CreateEnvelopeTaskInput](str, nil)
//	//lib.DPrintln(err)
//
//	res, cId, err := UT.BoxUsecase.PatientPaymentFormSignRequests("257082839680", "", "", "")
//	lib.DPrintln(res, cId, err)
//}

//func Test_BoxUsecase_ReleaseOfInformationSignRequests(t *testing.T) {
//	//task := &biz.CreateEnvelopeTaskInput{}
//	//
//	//str := `{"agentEmail":"liaogling@gmail.com","agentFirstName":"Gary","agentLastName":"Liao","clientEmail":"18891706@qq.com","clientFirstName":"FirstTest","clientLastName":"LastTest","signType":"box","templateId":"86df0f08-497c-4960-9121-afbcd4f4251b"}`
//	//task, err := lib.StringToTE[*biz.CreateEnvelopeTaskInput](str, nil)
//	//lib.DPrintln(err)
//
//	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5005)
//	client, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
//	prefillTags, err := UT.DbqsUsecase.ReleaseOfInformationPrefillTags(tCase, client)
//	if err != nil {
//		panic(err)
//	}
//	res, cId, err := UT.BoxUsecase.ReleaseOfInformationSignRequests("257082839680",
//		"lialing@foxmail.com", "liaogling@gmail.com", prefillTags)
//	lib.DPrintln(res, cId, err)
//}

func Test_BoxUsecase_MedicalTeamFormsSignRequests(t *testing.T) {

	prefill := lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      biz.GenFullName("Gary", "Liao"),
		},
		{
			"document_tag_id": "ClientNameSign",
			"text_value":      biz.GenFullName("Gary", "Liao"),
		},
		{
			"document_tag_id": "dob",
			"text_value":      "Jun. 8, 1984",
		},
		{
			"document_tag_id": "ssn",
			"text_value":      "573-79-7392",
		},
		{
			"document_tag_id": "phone",
			"text_value":      "619-948-5488",
		},
		{
			"document_tag_id": "email",
			"text_value":      "lialing@foxmail.com",
		},
		{
			"document_tag_id": "address",
			"text_value":      "000 W 9th Avenue",
		},
		{
			"document_tag_id": "location",
			"text_value":      "Escondido, CA 92000",
		},
		{
			"document_tag_id": "itf",
			"text_value":      "Mar. 1, 2025",
		},
		{
			"document_tag_id": "privateExamsNeeded",
			//"text_value":      "1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)\n        ",
			//"text_value": "Migraines (increase)\nLumbosacral strain w/ bilateral LE radiculopathy (increase)\nCervical strain (increase)\n——————————————————————————————————————————————————————————————",
			"text_value": "Migraines (increase)\nLumbosacral strain w/ bilateral LE radiculopathy (increase)\nCervical strain (increase)\n__________________________________________________________________",
			//"text_value": "Migraines (increase)\nLumbosacral strain w/ bilateral LE radiculopathy (increase)\nCervical strain (increase)\n",
		},
	}

	// lialing@foxmail.com
	// yannanwang@gmail.com
	boxSignTplId := ""
	a, b, er := UT.BoxUsecase.MedicalTeamFormsSignRequests("257590681761", "lialing@foxmail.com", "liaogling@gmail.com", prefill, "", boxSignTplId)
	lib.DPrintln(a)
	lib.DPrintln(b, er)
}

func Test_BoxUsecase_MedicalTeamFormsSignRequestsV2(t *testing.T) {

	// lialing@foxmail.com
	// yannanwang@gmail.com
	a, b, er := UT.BoxUsecase.MedicalTeamFormsSignRequestsWithoutTemplate("255166374931", "lialing@foxmail.com", "liaogling@gmail.com", "gengling.liao@hotmail.com", "1690115325489")
	lib.DPrintln(a)
	lib.DPrintln(b, er)
}

func Test_BoxUsecase_WebhooksList(t *testing.T) {
	res, err := UT.BoxUsecase.WebhooksList()
	lib.DPrintln(res, err)
}

/*
// 9bc35773-d528-41b0-9e19-71c5a777a04a 一个人签了
status:
converting, // 中间状态
created, // 中间状态
sent, // 中间状态
viewed, // #中间状态
signed, // 完成签属
cancelled, // 最后状态
declined, // 最后状态
error_converting, // 最后状态
error_sending,  // 最后状态
expired, // 最后状态
finalizing, // 最后状态
error_finalizing // 最后状态
*/
// {"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Your Veteran Benefits Center Contract From API","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (team@vetbenefitscenter.com)","are_reminders_enabled":false,"signers":[{"email":"team@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":true,"signer_decision":{"type":"signed","finalized_at":"2024-01-31T06:09:27.608Z","additional_info":null},"signer_group_id":null,"inputs":[{"document_tag_id":null,"text_value":"","checkbox_value":null,"date_value":null,"type":"signature","content_type":"signature","page_index":5},{"document_tag_id":"ClientNameSign","text_value":"Shi Li","checkbox_value":null,"date_value":null,"type":"text","content_type":"full_name","page_index":5},{"document_tag_id":null,"text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null,"type":"text","content_type":"email","page_index":5},{"document_tag_id":null,"text_value":"Jan 31, 2024","checkbox_value":null,"date_value":"2024-01-31","type":"date","content_type":"date","page_index":5},{"document_tag_id":null,"text_value":"Jan 31, 2024","checkbox_value":null,"date_value":"2024-01-31","type":"date","content_type":"date","page_index":0},{"document_tag_id":null,"text_value":"VT","checkbox_value":null,"date_value":null,"type":"signature","content_type":"initial","page_index":1}],"embed_url":null,"iframeable_embed_url":null},{"email":"liaogling@gmail.com","role":"signer","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null}],"id":"9bc35773-d528-41b0-9e19-71c5a777a04a","prefill_tags":[{"document_tag_id":"clientName","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"VSNameSign","text_value":"Gary Liao","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"246205309773","etag":"1","type":"folder","sequence_id":"1","name":"Test Box Sign Requests"},"name":"Your Veteran Benefits Center Contract(0) (1).pdf","external_id":null,"type":"sign-request","signing_log":{"id":"1430702826577","etag":"0","type":"file","sequence_id":"0","name":"Your Veteran Benefits Center Contract(0) (1) Signing Log.pdf","sha1":"012010c68a48d1911e249edcb24ddf50faeab4dd","file_version":{"id":"1568744082577","type":"file_version","sha1":"012010c68a48d1911e249edcb24ddf50faeab4dd"}},"status":"viewed","sign_files":{"files":[{"id":"1429944822737","etag":"1","type":"file","sequence_id":"1","name":"Your Veteran Benefits Center Contract(0) (1).pdf","sha1":"5224729774656b7fd0399fb4b7a03ea1f0effa2f","file_version":{"id":"1567872045137","type":"file_version","sha1":"52aed5e732c43468dc6a0575109cb84a8f09a173"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"2f1d4e50-3778-4345-adef-2165f30af5fe"}
func Test_BoxUsecase_GetSignRequest(t *testing.T) {
	// b1685f64-b285-4786-abc0-174931e7be0e // 已经签属
	// 9bc35773-d528-41b0-9e19-71c5a777a04a // 一个人签了
	res, err := UT.BoxUsecase.GetSignRequest("5ae49336-eb39-43dd-8b05-00cf4ae8998e")
	a, er := io.ReadAll(res.Body)
	lib.DPrintln(string(a), "___", er, res.StatusCode, res, err, "===")
	lib.DPrintln("+++")
}

func Test_BoxUsecase_ListSignRequest(t *testing.T) {
	// 已经签属？
	// b1685f64-b285-4786-abc0-174931e7be0e
	res, err := UT.BoxUsecase.ListSignRequest()
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_ListItemsInFolder(t *testing.T) {
	// 已经签属？
	// b1685f64-b285-4786-abc0-174931e7be0e
	res, err := UT.BoxUsecase.ListItemsInFolder("264686394097")
	lib.DPrintln(res, err)
	//res, err = UT.BoxUsecase.ListItemsInFolder("263406592476")
	//lib.DPrintln(res, err)
}

func Test_BoxUsecase_ListItemsInFolderFormat(t *testing.T) {
	// 已经签属？
	// b1685f64-b285-4786-abc0-174931e7be0e
	res, err := UT.BoxUsecase.ListItemsInFolderFormat("263406803830")
	lib.DPrintln(res, err)
	//res, err = UT.BoxUsecase.ListItemsInFolder("263406592476")
	//lib.DPrintln(res, err)
}

func Test_aaaa(t *testing.T) {
	sourceEntries := `{"total_count":35,"entries":[{"type":"folder","id":"263252423092","sequence_id":"0","etag":"0","name":"Alcantar, Francisco"},{"type":"folder","id":"261849118832","sequence_id":"0","etag":"0","name":"Alexander, Troy"},{"type":"folder","id":"263424822143","sequence_id":"0","etag":"0","name":"Allen, Robert"},{"type":"folder","id":"262831893453","sequence_id":"0","etag":"0","name":"Angulo, Don"},{"type":"folder","id":"260944375133","sequence_id":"0","etag":"0","name":"Ayuyao, Bernard"},{"type":"folder","id":"261514026402","sequence_id":"0","etag":"0","name":"Burgess, Dominique"},{"type":"folder","id":"262836135891","sequence_id":"0","etag":"0","name":"Camacho, Pedro"},{"type":"folder","id":"260937542087","sequence_id":"1","etag":"1","name":"Carrillo, Adrian"},{"type":"folder","id":"260939600452","sequence_id":"0","etag":"0","name":"Castillo, Jacinto"},{"type":"folder","id":"261717551273","sequence_id":"0","etag":"0","name":"De Leon, Cecelia"},{"type":"folder","id":"260937927949","sequence_id":"0","etag":"0","name":"Devine, Justin"},{"type":"folder","id":"260938403992","sequence_id":"0","etag":"0","name":"Dodd, Brent"},{"type":"folder","id":"260943786958","sequence_id":"0","etag":"0","name":"Galac, Cesar"},{"type":"folder","id":"260939356241","sequence_id":"0","etag":"0","name":"Houston, Keith"},{"type":"folder","id":"260940515467","sequence_id":"1","etag":"1","name":"Johnson, Christopher"},{"type":"folder","id":"260939834336","sequence_id":"0","etag":"0","name":"Johnson, Robin"},{"type":"folder","id":"260940914337","sequence_id":"1","etag":"1","name":"Jones, Cyril"},{"type":"folder","id":"262601317810","sequence_id":"0","etag":"0","name":"Keith, Kristopher"},{"type":"folder","id":"260942313293","sequence_id":"0","etag":"0","name":"Kennedy, Payton"},{"type":"folder","id":"261815812608","sequence_id":"0","etag":"0","name":"Lastrella, Amando"},{"type":"folder","id":"260939008536","sequence_id":"1","etag":"1","name":"Netemeyer, Aaron"},{"type":"folder","id":"262564917168","sequence_id":"0","etag":"0","name":"Orias, Ricardo"},{"type":"folder","id":"260942586693","sequence_id":"0","etag":"0","name":"Perez, Joaquin"},{"type":"folder","id":"260939281799","sequence_id":"0","etag":"0","name":"Perry, Fred"},{"type":"folder","id":"263250198764","sequence_id":"0","etag":"0","name":"Petit-Frere, Alexandre"},{"type":"folder","id":"261957561306","sequence_id":"0","etag":"0","name":"Ralat, Carlos"},{"type":"folder","id":"261683043974","sequence_id":"0","etag":"0","name":"Reynoso, Algis"},{"type":"folder","id":"261821085474","sequence_id":"0","etag":"0","name":"Shelrud, Cierra"},{"type":"folder","id":"261632097331","sequence_id":"0","etag":"0","name":"Sida, Andrew"},{"type":"folder","id":"261719322823","sequence_id":"1","etag":"1","name":"Smith, Andrew#5011"},{"type":"folder","id":"261962365453","sequence_id":"0","etag":"0","name":"Smith, Christopher"},{"type":"folder","id":"260940454986","sequence_id":"0","etag":"0","name":"Smith, Zane"},{"type":"folder","id":"263257386442","sequence_id":"0","etag":"0","name":"Smith,Andrew"},{"type":"folder","id":"260937572816","sequence_id":"1","etag":"1","name":"Stuart, James"},{"type":"folder","id":"260939913450","sequence_id":"0","etag":"0","name":"Tanquilut, Remigio"}],"offset":0,"limit":1000,"order":[{"by":"type","direction":"ASC"},{"by":"name","direction":"ASC"}]}`
	newEntries := `{"total_count":190,"entries":[{"type":"folder","id":"263407033477","sequence_id":"2","etag":"2","name":"Abutin, Niko Ralphluis"},{"type":"folder","id":"263407629244","sequence_id":"1","etag":"1","name":"Acuario, Edralin"},{"type":"folder","id":"263409690475","sequence_id":"1","etag":"1","name":"Albrecht, Keith Richard"},{"type":"folder","id":"263407432521","sequence_id":"1","etag":"1","name":"Alcantar, Francisco Jr."},{"type":"folder","id":"263409699906","sequence_id":"1","etag":"1","name":"Alexander, Keith David"},{"type":"folder","id":"263409220477","sequence_id":"1","etag":"1","name":"Alexander, Troy Don"},{"type":"folder","id":"263409383635","sequence_id":"1","etag":"1","name":"Allen, Robert Joseph"},{"type":"folder","id":"263408142981","sequence_id":"1","etag":"1","name":"Ancho, Romulo"},{"type":"folder","id":"263409340623","sequence_id":"1","etag":"1","name":"Anderson, Jilleah"},{"type":"folder","id":"263408260360","sequence_id":"1","etag":"1","name":"Anderson, Trever Shaw"},{"type":"folder","id":"263409604748","sequence_id":"1","etag":"1","name":"Andrews, Jamaal"},{"type":"folder","id":"263409085591","sequence_id":"1","etag":"1","name":"Angulo, Don Clark"},{"type":"folder","id":"263408638849","sequence_id":"1","etag":"1","name":"Arellano, Hector Gibram"},{"type":"folder","id":"263409608463","sequence_id":"1","etag":"1","name":"Ayala, Joe"},{"type":"folder","id":"263407544468","sequence_id":"1","etag":"1","name":"Ayuyao, Bernard Elijah"},{"type":"folder","id":"263407775301","sequence_id":"1","etag":"1","name":"Bailey, Bernard"},{"type":"folder","id":"263408325137","sequence_id":"1","etag":"1","name":"Baker, Jeffrey Stephen Jr."},{"type":"folder","id":"263407895271","sequence_id":"1","etag":"1","name":"Balavram, Jason Sy"},{"type":"folder","id":"263408186009","sequence_id":"1","etag":"1","name":"Baldemeca, Noel Genetia"},{"type":"folder","id":"263407002935","sequence_id":"1","etag":"1","name":"Ballesteros, Judeasar Galapon"},{"type":"folder","id":"263408572372","sequence_id":"1","etag":"1","name":"Banuex, Noemi"},{"type":"folder","id":"263409027699","sequence_id":"1","etag":"1","name":"Barajas, Emery Michael"},{"type":"folder","id":"263407850270","sequence_id":"1","etag":"1","name":"Barbosa, Brendan"},{"type":"folder","id":"263407864388","sequence_id":"1","etag":"1","name":"Barnes, Santana Venique"},{"type":"folder","id":"263408646454","sequence_id":"1","etag":"1","name":"Barnett, Jovan Eric"},{"type":"folder","id":"263408332367","sequence_id":"1","etag":"1","name":"Battle, Rodney Terez"},{"type":"folder","id":"263408192682","sequence_id":"1","etag":"1","name":"Bautista, Alexander Clement"},{"type":"folder","id":"263407518222","sequence_id":"1","etag":"1","name":"Becerra, Jonathan Contreras"},{"type":"folder","id":"263408711489","sequence_id":"1","etag":"1","name":"Beeler, Rebecca #5019"},{"type":"folder","id":"263408505251","sequence_id":"1","etag":"1","name":"Blaine, Christopher Charles"},{"type":"folder","id":"263408101803","sequence_id":"1","etag":"1","name":"Boden, Boyce Robert"},{"type":"folder","id":"263408584289","sequence_id":"1","etag":"1","name":"Bolino, Louis Alto"},{"type":"folder","id":"263408994347","sequence_id":"1","etag":"1","name":"Briggs, Scott Christopher"},{"type":"folder","id":"263408480938","sequence_id":"1","etag":"1","name":"Brown, Terence Lewis"},{"type":"folder","id":"263408941892","sequence_id":"1","etag":"1","name":"Burgess, Dominique Farrell"},{"type":"folder","id":"263407893513","sequence_id":"1","etag":"1","name":"Camacho, Pedro"},{"type":"folder","id":"263408432830","sequence_id":"1","etag":"1","name":"Campbell, Dailyn"},{"type":"folder","id":"263408809561","sequence_id":"1","etag":"1","name":"Canseco, Manuel Valiente"},{"type":"folder","id":"263408946196","sequence_id":"1","etag":"1","name":"Carr, Christopher"},{"type":"folder","id":"263407168436","sequence_id":"1","etag":"1","name":"Carreon, Lovelito Flores"},{"type":"folder","id":"263408315542","sequence_id":"1","etag":"1","name":"Carrillo, Adrian"},{"type":"folder","id":"263407400559","sequence_id":"1","etag":"1","name":"Carter, Gabrielle Alexandria"},{"type":"folder","id":"263409575210","sequence_id":"1","etag":"1","name":"Castelluccio, Dillon Thomas"},{"type":"folder","id":"263409090018","sequence_id":"1","etag":"1","name":"Castillo, Jacinto"},{"type":"folder","id":"263408282334","sequence_id":"1","etag":"1","name":"Castro, Roman Lucas"},{"type":"folder","id":"263408953729","sequence_id":"1","etag":"1","name":"Cavaliere, Robynn"},{"type":"folder","id":"263408315294","sequence_id":"1","etag":"1","name":"Chacon, Guy #5052"},{"type":"folder","id":"263409193805","sequence_id":"1","etag":"1","name":"Claire, Seth Alexander"},{"type":"folder","id":"263409058831","sequence_id":"1","etag":"1","name":"Cobian Jr., Jose"},{"type":"folder","id":"263409335804","sequence_id":"1","etag":"1","name":"Coley, Alvin Lorenzo"},{"type":"folder","id":"263407511190","sequence_id":"1","etag":"1","name":"Collins, Samantha Rose"},{"type":"folder","id":"263408706463","sequence_id":"1","etag":"1","name":"Cook, Robert #5020"},{"type":"folder","id":"263408775705","sequence_id":"1","etag":"1","name":"D'Alessandro, James Scott"},{"type":"folder","id":"263408622377","sequence_id":"1","etag":"1","name":"De Leon, Cecilia"},{"type":"folder","id":"263409573530","sequence_id":"1","etag":"1","name":"DelPrete, Desiree"},{"type":"folder","id":"263408560027","sequence_id":"1","etag":"1","name":"Demps, Kelvin"},{"type":"folder","id":"263409486296","sequence_id":"1","etag":"1","name":"Deocampo, Teddy Baysan"},{"type":"folder","id":"263408562555","sequence_id":"1","etag":"1","name":"Devine, Justin Michael"},{"type":"folder","id":"263409697401","sequence_id":"1","etag":"1","name":"DiBenedetto, Michael"},{"type":"folder","id":"263408598259","sequence_id":"1","etag":"1","name":"Dickerson, Ramon"},{"type":"folder","id":"263408170942","sequence_id":"1","etag":"1","name":"Dickey, James #5005"},{"type":"folder","id":"263408262134","sequence_id":"1","etag":"1","name":"Dishmon, Varian Dione"},{"type":"folder","id":"263407789701","sequence_id":"1","etag":"1","name":"Dodd, Brent"},{"type":"folder","id":"263408883616","sequence_id":"1","etag":"1","name":"Dukes, Spencer Patrick"},{"type":"folder","id":"263409395451","sequence_id":"1","etag":"1","name":"Dunkin, John Steven"},{"type":"folder","id":"263408255849","sequence_id":"1","etag":"1","name":"Edwards, Addison Chase"},{"type":"folder","id":"263407679255","sequence_id":"1","etag":"1","name":"FaisonLanier, Terrica"},{"type":"folder","id":"263408082641","sequence_id":"1","etag":"1","name":"Fannin, Quentin Dekote"},{"type":"folder","id":"263407698819","sequence_id":"1","etag":"1","name":"Farmer, Anthony"},{"type":"folder","id":"263407825270","sequence_id":"1","etag":"1","name":"Flores, Jacinth Aaron"},{"type":"folder","id":"263409236508","sequence_id":"1","etag":"1","name":"Fowler, Casey"},{"type":"folder","id":"263409256573","sequence_id":"1","etag":"1","name":"Fowler, Jerry Joseph"},{"type":"folder","id":"263407897776","sequence_id":"1","etag":"1","name":"Franco, Alec Robert"},{"type":"folder","id":"263409003856","sequence_id":"1","etag":"1","name":"Galac, Cesar"},{"type":"folder","id":"263408776073","sequence_id":"1","etag":"1","name":"Garcia, Rey #5041"},{"type":"folder","id":"263408716059","sequence_id":"1","etag":"1","name":"Garlejo, Jason Doctolero"},{"type":"folder","id":"263409049595","sequence_id":"1","etag":"1","name":"Gilmore, Earl Glenn Jr."},{"type":"folder","id":"263408759623","sequence_id":"1","etag":"1","name":"Gonzales, Beau Matthew"},{"type":"folder","id":"263407379970","sequence_id":"1","etag":"1","name":"Goodson, Augustus Ivan IV"},{"type":"folder","id":"263409028000","sequence_id":"1","etag":"1","name":"Green, Antionette"},{"type":"folder","id":"263407516742","sequence_id":"1","etag":"1","name":"Green, Donnell"},{"type":"folder","id":"263409561453","sequence_id":"1","etag":"1","name":"Green, Sandra"},{"type":"folder","id":"263407261628","sequence_id":"1","etag":"1","name":"Griffin, Major Pete III"},{"type":"folder","id":"263407014935","sequence_id":"1","etag":"1","name":"Haley, George Walter Jr."},{"type":"folder","id":"263407237574","sequence_id":"1","etag":"1","name":"Harris, Debra Lee"},{"type":"folder","id":"263408234122","sequence_id":"1","etag":"1","name":"Herron, Leslie Rhea"},{"type":"folder","id":"263409282194","sequence_id":"1","etag":"1","name":"Hiers, Tommy Lamar Jr."},{"type":"folder","id":"263408977216","sequence_id":"1","etag":"1","name":"Ho, Duyet #5043"},{"type":"folder","id":"263407717756","sequence_id":"1","etag":"1","name":"Houston, Keith Anthony"},{"type":"folder","id":"263407411164","sequence_id":"1","etag":"1","name":"Howard, Nicole Marie"},{"type":"folder","id":"263408903510","sequence_id":"1","etag":"1","name":"Huerta Jr., Edward David"},{"type":"folder","id":"263408925455","sequence_id":"1","etag":"1","name":"Huerta Jr., Edward David #96"},{"type":"folder","id":"263407497275","sequence_id":"1","etag":"1","name":"Huerta Sr., Edward David #260"},{"type":"folder","id":"263407386390","sequence_id":"1","etag":"1","name":"Hutchinson, Derek"},{"type":"folder","id":"263408325009","sequence_id":"1","etag":"1","name":"Ibarrondo, William"},{"type":"folder","id":"263408646453","sequence_id":"1","etag":"1","name":"Ibarrondo, William Basean"},{"type":"folder","id":"263407665390","sequence_id":"1","etag":"1","name":"Inzer, Russell Dustin"},{"type":"folder","id":"263408617842","sequence_id":"1","etag":"1","name":"Jacob, Melvin"},{"type":"folder","id":"263408202281","sequence_id":"1","etag":"1","name":"James, Dillon Randall"},{"type":"folder","id":"263409584864","sequence_id":"1","etag":"1","name":"Johnson, Christopher #5009"},{"type":"folder","id":"263409280743","sequence_id":"1","etag":"1","name":"Johnson, Jermaine"},{"type":"folder","id":"263409239761","sequence_id":"1","etag":"1","name":"Johnson, Michael"},{"type":"folder","id":"263407513317","sequence_id":"1","etag":"1","name":"Johnson, Robin"},{"type":"folder","id":"263408383080","sequence_id":"1","etag":"1","name":"Jones, Cyril Evans Jr."},{"type":"folder","id":"263407472517","sequence_id":"1","etag":"1","name":"Jones, Naurice"},{"type":"folder","id":"263409011236","sequence_id":"1","etag":"1","name":"Keith, Kristopher #5025"},{"type":"folder","id":"263407636399","sequence_id":"1","etag":"1","name":"Keller, Tony Jordan"},{"type":"folder","id":"263407378937","sequence_id":"1","etag":"1","name":"Kennedy, Payton Gary"},{"type":"folder","id":"263407672521","sequence_id":"1","etag":"1","name":"Kubas III, William Philip"},{"type":"folder","id":"263407746455","sequence_id":"1","etag":"1","name":"Lane, Michael John"},{"type":"folder","id":"263408048756","sequence_id":"1","etag":"1","name":"Lang, Shanise Eileen"},{"type":"folder","id":"263407759064","sequence_id":"1","etag":"1","name":"Lastrella, Amando Llorin"},{"type":"folder","id":"263407079735","sequence_id":"1","etag":"1","name":"Laxa, Eduardo Dulu"},{"type":"folder","id":"263407578819","sequence_id":"1","etag":"1","name":"Le, Khanh Si"},{"type":"folder","id":"263409409666","sequence_id":"1","etag":"1","name":"Lepe, Christopher"},{"type":"folder","id":"263407774875","sequence_id":"1","etag":"1","name":"Liggans, Nyjerus Lavondai Onijar"},{"type":"folder","id":"263408704030","sequence_id":"1","etag":"1","name":"Mangra, Robbi #5023"},{"type":"folder","id":"263409358994","sequence_id":"1","etag":"1","name":"Marcial, Janry"},{"type":"folder","id":"263408360692","sequence_id":"1","etag":"1","name":"Mendez, Marco Antonio"},{"type":"folder","id":"263407919167","sequence_id":"1","etag":"1","name":"Montoya, Mario"},{"type":"folder","id":"263408174156","sequence_id":"1","etag":"1","name":"Montoya, Mary"},{"type":"folder","id":"263407556962","sequence_id":"1","etag":"1","name":"Moreno, Michael Joey"},{"type":"folder","id":"263408493117","sequence_id":"1","etag":"1","name":"Morris, Lee Roy"},{"type":"folder","id":"263407593367","sequence_id":"1","etag":"1","name":"Murrell, Brashaad"},{"type":"folder","id":"263408197832","sequence_id":"1","etag":"1","name":"Myers, Luis Anthony"},{"type":"folder","id":"263408754259","sequence_id":"1","etag":"1","name":"Nellis, Lawrence #5024"},{"type":"folder","id":"263409085859","sequence_id":"1","etag":"1","name":"Netemeyer, Aaron"},{"type":"folder","id":"263409335873","sequence_id":"1","etag":"1","name":"Newman, James Wesley"},{"type":"folder","id":"263407708055","sequence_id":"1","etag":"1","name":"Olivetti, Anthony Ryan Borja"},{"type":"folder","id":"263409189404","sequence_id":"1","etag":"1","name":"Orias, Ricardo"},{"type":"folder","id":"263407676838","sequence_id":"1","etag":"1","name":"Padua, Michael Daniel"},{"type":"folder","id":"263409107750","sequence_id":"1","etag":"1","name":"Peca, Jason"},{"type":"folder","id":"263407622264","sequence_id":"1","etag":"1","name":"Perez, Joaquin Xavier"},{"type":"folder","id":"263409369154","sequence_id":"1","etag":"1","name":"Perry, Fred Douglas"},{"type":"folder","id":"263408536697","sequence_id":"1","etag":"1","name":"Petit-Frere, Alexandre Freud"},{"type":"folder","id":"263408613185","sequence_id":"1","etag":"1","name":"Pharnes, Eric Dwayne"},{"type":"folder","id":"263408872130","sequence_id":"1","etag":"1","name":"Pierre, Gilbert"},{"type":"folder","id":"263408335098","sequence_id":"1","etag":"1","name":"Prado, Martius Oris"},{"type":"folder","id":"263409357404","sequence_id":"1","etag":"1","name":"Pratko, Michael"},{"type":"folder","id":"263407645792","sequence_id":"1","etag":"1","name":"Provasek, Jared #5046"},{"type":"folder","id":"263407539990","sequence_id":"1","etag":"1","name":"Ralat, Carlos Alberto"},{"type":"folder","id":"263407549727","sequence_id":"1","etag":"1","name":"Reynoso, Algis #5036"},{"type":"folder","id":"263407439423","sequence_id":"1","etag":"1","name":"Rios, Jonathan David"},{"type":"folder","id":"263408483369","sequence_id":"1","etag":"1","name":"Rivera, Louie"},{"type":"folder","id":"263408128360","sequence_id":"1","etag":"1","name":"Rosales, Juan"},{"type":"folder","id":"263409671231","sequence_id":"1","etag":"1","name":"Rutledge, Ronnie #5042"},{"type":"folder","id":"263408531552","sequence_id":"1","etag":"1","name":"Salb, Austin Reid"},{"type":"folder","id":"263407996012","sequence_id":"1","etag":"1","name":"Santillan-Mondaca, Kristy #5050"},{"type":"folder","id":"263407479863","sequence_id":"1","etag":"1","name":"Sayles, Matthew Evans"},{"type":"folder","id":"263408282344","sequence_id":"1","etag":"1","name":"Serrano, Ronald"},{"type":"folder","id":"263407777243","sequence_id":"1","etag":"1","name":"Shelrud, Cierra Kay"},{"type":"folder","id":"263408291548","sequence_id":"1","etag":"1","name":"Sida, Andrew Michael"},{"type":"folder","id":"263409143090","sequence_id":"1","etag":"1","name":"Slater, Jamie"},{"type":"folder","id":"263408173982","sequence_id":"1","etag":"1","name":"Smith, Andrew"},{"type":"folder","id":"263409109595","sequence_id":"1","etag":"1","name":"Smith, Andrew #5011"},{"type":"folder","id":"263409138194","sequence_id":"1","etag":"1","name":"Smith, Austin Cole"},{"type":"folder","id":"263408466241","sequence_id":"1","etag":"1","name":"Smith, Christopher Manuel #5026"},{"type":"folder","id":"263408730987","sequence_id":"1","etag":"1","name":"Smith, Max"},{"type":"folder","id":"263409028559","sequence_id":"1","etag":"1","name":"Smith, Rhyheime"},{"type":"folder","id":"263408486215","sequence_id":"1","etag":"1","name":"Smith, Zane Eugene"},{"type":"folder","id":"263407441338","sequence_id":"1","etag":"1","name":"Smolinski, Donald Jay"},{"type":"folder","id":"263409085790","sequence_id":"1","etag":"1","name":"Stacks, Taylor"},{"type":"folder","id":"263408843374","sequence_id":"1","etag":"1","name":"Stewart, Robert #5038"},{"type":"folder","id":"263408879881","sequence_id":"1","etag":"1","name":"Stuart, James Francis"},{"type":"folder","id":"263409189108","sequence_id":"1","etag":"1","name":"Summer, Aaron Michael Blake"},{"type":"folder","id":"263408183176","sequence_id":"1","etag":"1","name":"Sutton, Derrick"},{"type":"folder","id":"263408877126","sequence_id":"1","etag":"1","name":"Tanquilut, Remigio Dimalanta"},{"type":"folder","id":"263408462333","sequence_id":"1","etag":"1","name":"Taylor, Bobbee Nykole"},{"type":"folder","id":"263408290941","sequence_id":"1","etag":"1","name":"Terrell, Conley II"},{"type":"folder","id":"263408989699","sequence_id":"1","etag":"1","name":"TestLi, TestShi #5057"},{"type":"folder","id":"263408579201","sequence_id":"1","etag":"1","name":"TestLi, TestShi #5058"},{"type":"folder","id":"263407556960","sequence_id":"1","etag":"1","name":"TestLiao, TestGary #5061"},{"type":"folder","id":"263470914421","sequence_id":"0","etag":"0","name":"TestLiao, TestGary #5064"},{"type":"folder","id":"263409010923","sequence_id":"1","etag":"1","name":"Thompson, Jason Scott"},{"type":"folder","id":"263408644030","sequence_id":"1","etag":"1","name":"Thrower, Tony Lee"},{"type":"folder","id":"263408560526","sequence_id":"1","etag":"1","name":"Tran, Danny Minh"},{"type":"folder","id":"263408893699","sequence_id":"1","etag":"1","name":"Tran, Joanne Perea"},{"type":"folder","id":"263409387704","sequence_id":"1","etag":"1","name":"Tran, Kenny #5045"},{"type":"folder","id":"263408661172","sequence_id":"1","etag":"1","name":"Valdez, Jacob"},{"type":"folder","id":"263408233836","sequence_id":"1","etag":"1","name":"Valdez, Joshua Raul"},{"type":"folder","id":"263407744765","sequence_id":"1","etag":"1","name":"Valli, Matthew Lawrence"},{"type":"folder","id":"263408943728","sequence_id":"1","etag":"1","name":"Valli, Ronald"},{"type":"folder","id":"263408636453","sequence_id":"1","etag":"1","name":"Vargas, Gabrian"},{"type":"folder","id":"263408008709","sequence_id":"1","etag":"1","name":"Velasquez, Jose David"},{"type":"folder","id":"263407730033","sequence_id":"1","etag":"1","name":"Walker, Ronald Stanley"},{"type":"folder","id":"263408847871","sequence_id":"1","etag":"1","name":"Warren, Olga"},{"type":"folder","id":"263408262984","sequence_id":"1","etag":"1","name":"Watts, Patrick Levon"},{"type":"folder","id":"263407622121","sequence_id":"1","etag":"1","name":"Webster, Craig"},{"type":"folder","id":"263408430635","sequence_id":"1","etag":"1","name":"West, Melissa Theresa"},{"type":"folder","id":"263408356225","sequence_id":"1","etag":"1","name":"White, Michael #5015"}],"offset":0,"limit":1000,"order":[{"by":"type","direction":"ASC"},{"by":"name","direction":"ASC"}]}`

	a := lib.ToTypeMapByString(sourceEntries)
	b := lib.ToTypeMapByString(newEntries)
	aEntries := a.GetTypeList("entries")
	bEntries := b.GetTypeList("entries")

	for _, v := range aEntries {
		cc := v.GetString("name")
		//lib.DPrintln(cc)
		destFolderId := ""
		for _, v1 := range bEntries {
			nName := v1.GetString("name")
			if nName == cc {
				destFolderId = v1.GetString("id")
				break
			}

			ss := strings.Split(cc, ",")
			tt := strings.TrimSpace(ss[1])

			ss1 := strings.Split(nName, ",")
			tt1 := strings.TrimSpace(ss1[1])

			if tt == tt1 {
				//fmt.Println("++++", tt, tt1)
				destFolderId = v1.GetString("id")
				break
			}

		}
		//lib.DPrintln("ccc:", destFolderId, cc)
		if destFolderId == "" {
			lib.DPrintln("err:", v.GetString("id"))
		} else {
			fmt.Println(v.GetString("id"), destFolderId)
		}
	}

	lib.DPrintln(aEntries)
	lib.DPrintln(bEntries)
}

func Test_BoxUsecase_DeleteFolder(t *testing.T) {
	res, err := UT.BoxUsecase.DeleteFolder("263466623853", false)
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_CopyFileNewFileNameReturnFileId(t *testing.T) {
	fileId, err := UT.BoxUsecase.CopyFileNewFileNameReturnFileId("1535378103870", config_box.FileName_ThingsToKnowExam, "267683828516")
	lib.DPrintln(fileId, err)
}

func Test_BoxUsecase_CopyFileNewFileName(t *testing.T) {
	res, err := UT.BoxUsecase.CopyFileNewFileName("1535378103870", config_box.FileName_ThingsToKnowExam, "267683828516")
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_CopyFile(t *testing.T) {
	res, code, err := UT.BoxUsecase.CopyFile("1527958777860", "263406592476")
	lib.DPrintln("code:", code)
	lib.DPrintln("err:", err)
	lib.DPrintln(res)
}

func Test_BoxUsecase_MoveFile(t *testing.T) {
	res, err := UT.BoxUsecase.MoveFile("1541224873274", "264658751993")
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_UploadFileVersion(t *testing.T) {

	fileOp, err := os.Open("res/b.pdf")
	if err != nil {
		panic(err)
	}
	defer fileOp.Close()

	res, err := UT.BoxUsecase.UploadFileVersion("1541221137574", fileOp)
	lib.DPrintln(res, err)
}

func Test_BoxUsecase_DownloadFile(t *testing.T) {

	fileOp, err := os.Open("res/b.pdf")
	if err != nil {
		panic(err)
	}
	defer fileOp.Close()

	// 版本与页面URL对应
	res, err := UT.BoxUsecase.DownloadFile("1541221137574", "1734145792144")
	if res != nil {
		defer res.Close()
	}
	lib.DPrintln("res:", res, "err:", err)
	file, err := os.Create("tmp/bbb_new2.pdf")
	a, err := io.Copy(file, res)
	lib.DPrintln("a:", a, err)
}

func Test_BoxUsecase_DeleteFile(t *testing.T) {
	aa, err := UT.BoxUsecase.DeleteFile("1547834286959")
	lib.DPrintln(aa, err)
}

func Test_BoxUsecase_GetFileInfoForTypeMap(t *testing.T) {
	r, httpCode, err := UT.BoxUsecase.GetFileInfoForTypeMap("1554567330874")
	lib.DPrintln(httpCode)
	lib.DPrintln(r, err)
}

func Test_BoxUsecase_GetFileInfo(t *testing.T) {
	r, httpCode, err := UT.BoxUsecase.GetFileInfo("1956511490228")
	lib.DPrintln(r, httpCode, err)
}

func Test_BoxUsecase_GetFileVersions(t *testing.T) {
	r, httpCode, err := UT.BoxUsecase.GetFileVersions("1541221137574")
	lib.DPrintln(r, httpCode, err)
}

// "created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T03:45:48-07:00"
// "created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T04:03:24-07:00" 上传文件后，有修改，可能有延时
// "created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T04:04:18-07:00" 目录创建文件夹修改时间有变化
// "created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T04:06:20-07:00" 子子目录创建文件夹后，修改时间也有变化
// "created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T04:08:32-07:00"
// 通过文件夹修改时间，有风险：暂时不通过此方法
func Test_BoxUsecase_GetFolderInfo(t *testing.T) {
	r, httpCode, err := UT.BoxUsecase.GetFolderInfo("268906213262")
	lib.DPrintln(httpCode)
	lib.DPrintln(r, err)
}

func Test_BoxUsecase_CancelSignRequest(t *testing.T) {
	// 97a6f0e3-a238-4fb2-95fd-37c5136dd266
	res, err := UT.BoxUsecase.CancelSignRequest("97a6f0e3-a238-4fb2-95fd-37c5136dd266")
	lib.DPrintln(res)
	lib.DPrintln(err)
}

func Test_BoxUsecase_UploadFileVersionWithNewFileName(t *testing.T) {

	file, err := os.Open("/Users/garyliao/Desktop/res/pdf-test.pdf")
	if err != nil {
		lib.DPrintln(err)
		return
	}
	defer file.Close()
	aa, err := UT.BoxUsecase.UploadFileVersionWithNewFileName("1639619613514", file, "diffSize v2.pdf")
	lib.DPrintln("err:", err)
	lib.DPrintln(aa)
}

func Test_BoxUsecase_DownloadToLocal(t *testing.T) {
	a, b, err := UT.BoxUsecase.DownloadToLocal("1757895743473", "docx")
	lib.DPrintln(a, b, err)
}

func Test_BoxUsecase_Collaborations(t *testing.T) {
	b, err := UT.BoxUsecase.Collaborations("319249984263", "lialing@foxmail.com")
	lib.DPrintln(b, err)
}

func Test_BoxUsecase_Collaborations1(t *testing.T) {

	//Jeremy Rosas-20#5563
	//Dalton Shipley-70#5557
	//Conrado Recasas-40#5597
	//Eugenio Colon-0#5598
	//Matthew Honc-0#5587
	//Charles Webb-60#5593
	//Rashard Marshall-60#5590
	//Melissa Kapoi-30#5589

	// 5563 手动 Jeremy Rosas-20#5563
	// 5557 ok Dalton Shipley-70#5557
	// 5597 ok  Conrado Recasas-40#5597
	// 5598 ok Eugenio Colon-0#5598
	// 5587 ok Matthew Honc-0#5587
	// 5593 ok Charles Webb-60#5593
	// 5590 ok Rashard Marshall-60#5590
	// 5589 ok Melissa Kapoi-30#5589
	//
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5589)
	_, tContactFields, err := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		panic(err)
	}
	email := tContactFields.TextValueByNameBasic("email")
	if email == "" {
		panic("Email does not exists.")
	}

	clientBoxFolderId, err := UT.BoxbuzUsecase.GetClientBoxFolderId(tCase)
	if err != nil {
		panic(err)
	}
	if clientBoxFolderId == "" {
		panic("clientBoxFolderId is empty")
	}

	lib.DPrintln("clientBoxFolderId: ", clientBoxFolderId, email)
	b, err := UT.BoxUsecase.Collaborations(clientBoxFolderId, email)
	if err != nil {
		panic(err)
	}
	lib.DPrintln(b, err)
	UpdateCaseFolder(tCase)
}

func Test_BoxUsecase_Collaborations2(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5598)
	UpdateCaseFolder(tCase)

}
func UpdateCaseFolder(tCase *biz.TData) {
	//tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5598)
	//_, tContactFields, err := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic("client_gid"))
	//if err != nil {
	//	panic(err)
	//}
	//email := tContactFields.TextValueByNameBasic("email")
	//if email == "" {
	//	panic("Email does not exists.")
	//}

	clientBoxFolderId, err := UT.BoxbuzUsecase.GetClientBoxFolderId(tCase)
	if err != nil {
		panic(err)
	}
	if clientBoxFolderId == "" {
		panic("clientBoxFolderId is empty")
	}

	//lib.DPrintln("clientBoxFolderId: ", clientBoxFolderId, email)

	caseFileFolderValue := tCase.CustomFields.TextValueByNameBasic(biz.FieldName_case_files_folder)
	clientCaseId := tCase.Id()
	if caseFileFolderValue == "" {
		row := make(lib.TypeMap)
		key := fmt.Sprintf("%s%d", biz.Map_ClientBoxFolderId, clientCaseId)
		boxFolderId, _ := UT.MapUsecase.GetForString(key)
		if boxFolderId != "" {
			row.Set(biz.FieldName_case_files_folder, "https://veteranbenefitscenter.app.box.com/folder/"+boxFolderId)
		}
		if len(row) > 0 {
			row.Set(biz.DataEntry_gid, tCase.Gid())
			_, err = UT.DataEntryUsecase.HandleOne(biz.Kind_client_cases, biz.TypeDataEntry(row), biz.DataEntry_gid, nil)
		}
	}
}

func Test_BoxUsecase_Users(t *testing.T) {
	a, b, err := UT.BoxUsecase.Users()
	lib.DPrintln(a, b, err)
}

func Test_BoxUsecase_CollaborationsByBoxUserId(t *testing.T) {
	// [INFO ts=2025-06-24T08:09:27Z biz_box_test.go:485] 64110671893 201
	// 39217319801
	a, code, err := UT.BoxUsecase.CollaborationsByBoxUserId("327431320500", "39217319801")
	lib.DPrintln(a, code, err)
	// 返回 a:64086181594
}

func Test_BoxUsecase_FolderCollaborations(t *testing.T) {
	// 64086181594
	a, b, err := UT.BoxUsecase.FolderCollaborations("327431320500")
	lib.DPrintln(a, b, err)
}

func Test_BoxUsecase_DeleteCollaborations(t *testing.T) {
	a, err := UT.BoxUsecase.DeleteCollaborations("64112567378")
	lib.DPrintln(a, err)
}
