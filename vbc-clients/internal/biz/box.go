package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
	"vbc/lib/to"
	"vbc/lib/uuid"
)

type BoxUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
	UsageStatsUsecase  *UsageStatsUsecase
}

func NewBoxUsecase(logger log.Logger, CommonUsecase *CommonUsecase, conf *conf.Data, Oauth2TokenUsecase *Oauth2TokenUsecase, UsageStatsUsecase *UsageStatsUsecase) *BoxUsecase {
	return &BoxUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
		UsageStatsUsecase:  UsageStatsUsecase,
	}

}

func (c *BoxUsecase) Token() (string, error) {
	e, err := c.Oauth2TokenUsecase.GetByAppId(Oauth2_AppId_box)
	if err != nil {
		return "", err
	}
	if e == nil {
		return "", errors.New("BoxUsecase token is nil.")
	}
	return e.AccessToken, nil
}

func (c *BoxUsecase) ListBoxSignTemplates() (res *string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/sign_templates", c.conf.Box.ApiUrl)
	res, _, err = lib.HTTPJsonWithHeaders("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return
}

func (c *BoxUsecase) GetSignRequest(signRequestId string) (response *http.Response, err error) {
	// sign_requests
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests/%s", c.conf.Box.ApiUrl, signRequestId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "GetSignRequest"), time.Now(), 1)
	return lib.RequestDo("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
}

// 成功返回：{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Your Veteran Benefits Center Contract","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"ebunting@vetbenefitscenter.com","role":"signer","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"97a6f0e3-a238-4fb2-95fd-37c5136dd266","prefill_tags":[{"document_tag_id":"clientName","text_value":"TestF TestL","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"TestF TestL","checkbox_value":null,"date_value":null},{"document_tag_id":"VSNameSign","text_value":"Edward Bunting Jr.","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"276350317208","etag":"0","type":"folder","sequence_id":"0","name":"TestL, TestF #5164"},"name":"Agreement for Consulting Services.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"cancelled","sign_files":{"files":[{"id":"1599331175954","etag":"0","type":"file","sequence_id":"0","name":"Agreement for Consulting Services.pdf","sha1":"74311f7a04c32b632829f9f5eca4fc9589a1c56d","file_version":{"id":"1758060655154","type":"file_version","sha1":"74311f7a04c32b632829f9f5eca4fc9589a1c56d"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"8a19aa3c-689b-4e46-a4ef-cc2bf318144d","external_system_name":null}
// 合同ID不存在：返回404
// 合同已经取消了： 返回400
func (c *BoxUsecase) CancelSignRequest(signRequestId string) (res *string, err error) {
	// sign_requests
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests/%s/cancel", c.conf.Box.ApiUrl, signRequestId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CancelSignRequest"), time.Now(), 1)

	res, _, err = lib.Request("POST", api, nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil, err
	}
	if res != nil {
	}
	return res, nil
}

func (c *BoxUsecase) ListSignRequest() (res *string, err error) {
	// sign_requests
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "ListSignRequest"), time.Now(), 1)
	res, _, err = lib.HTTPJsonWithHeaders("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return
}

// https://box.dev/reference/post-sign-requests/#param-template_id
// 返回参数：{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Your Veteran Benefits Center Contract From API","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (team@vetbenefitscenter.com)","are_reminders_enabled":false,"signers":[{"email":"team@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null},{"email":"glliao@vetbenefitscenter.com","role":"signer","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null}],"id":"5a36e3d9-05aa-4f19-bc7d-1c2ff92fc32e","prefill_tags":[{"document_tag_id":"clientName","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"Shi Li","checkbox_value":null,"date_value":null},{"document_tag_id":"VSNameSign","text_value":"Gary Liao","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"246205309773","etag":"1","type":"folder","sequence_id":"1","name":"Test Box Sign Requests"},"name":"Your Veteran Benefits Center Contract.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1426975036434","etag":"0","type":"file","sequence_id":"0","name":"Your Veteran Benefits Center Contract.pdf","sha1":"52aed5e732c43468dc6a0575109cb84a8f09a173","file_version":{"id":"1564573626834","type":"file_version","sha1":"52aed5e732c43468dc6a0575109cb84a8f09a173"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"a0a2df1f-3fd4-42f2-8881-6d6b1df55ea2"}
// folderId: 合同存在的位置  示例:246205309773
func (c *BoxUsecase) SignRequests(createEnvelopeTaskInput *CreateEnvelopeTaskInput, folderId string) (res *string, contractUniqId string, err error) {

	if createEnvelopeTaskInput == nil {
		return nil, "", errors.New("SignRequests:SignRequests is nil")
	}
	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if configs.IsProd() {
		params.Set("email_subject", "Your Veteran Benefits Center Contract")
	} else {
		params.Set("email_subject", "Your Veteran Benefits Center Contract")
	}
	params.Set("are_reminders_enabled", true)
	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": createEnvelopeTaskInput.ClientEmail,
			"order": 1,
		},
		{
			"role":  "signer",
			"email": createEnvelopeTaskInput.AgentEmail,
			"order": 2,
		},
	})
	params.Set("prefill_tags", lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      GenFullName(createEnvelopeTaskInput.ClientFirstName, createEnvelopeTaskInput.ClientLastName),
		},
		{
			"document_tag_id": "ClientNameSign",
			"text_value":      GenFullName(createEnvelopeTaskInput.ClientFirstName, createEnvelopeTaskInput.ClientLastName),
		},
		{
			"document_tag_id": "VSNameSign",
			"text_value":      GenFullName(createEnvelopeTaskInput.AgentFirstName, createEnvelopeTaskInput.AgentLastName),
		},
	})
	params.Set("template_id", createEnvelopeTaskInput.TemplateId) // a0a2df1f-3fd4-42f2-8881-6d6b1df55ea2
	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})
	res, _, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "SignRequests"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_CreateContract, logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

func (c *BoxUsecase) SignRequestsWithoutTemplateAm(ClientEmail string, AgentEmail string, attorneyFullName string, folderId string, pdfSourceBoxFileId string) (res *string, contractUniqId string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)

	params.Set("email_subject", "Your VA Representation Agreement with August Miles - "+attorneyFullName)

	params.Set("are_reminders_enabled", true)
	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": ClientEmail,
			"order": 1,
		},
		{
			"role":  "signer",
			"email": AgentEmail,
			"order": 2,
		},
	})
	params.Set("source_files", lib.TypeList{lib.TypeMap{
		"type": "file",
		"id":   pdfSourceBoxFileId,
	}})
	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})
	res, _, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "SignRequestsWithoutTemplateAm"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_CreateContractAm, logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

func (c *BoxUsecase) SignRequestsWithoutTemplate(createEnvelopeTaskInput *CreateEnvelopeTaskInput, folderId string, pdfSourceBoxFileId string) (res *string, contractUniqId string, err error) {

	if createEnvelopeTaskInput == nil {
		return nil, "", errors.New("SignRequests:SignRequests is nil")
	}
	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if configs.IsProd() {
		params.Set("email_subject", "Your Veteran Benefits Center Contract")
	} else {
		params.Set("email_subject", "Your Veteran Benefits Center Contract")
	}
	params.Set("are_reminders_enabled", true)
	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": createEnvelopeTaskInput.ClientEmail,
			"order": 1,
		},
		{
			"role":  "signer",
			"email": createEnvelopeTaskInput.AgentEmail,
			"order": 2,
		},
	})
	params.Set("source_files", lib.TypeList{lib.TypeMap{
		"type": "file",
		"id":   pdfSourceBoxFileId,
	}})
	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})
	res, _, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "SignRequests"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_CreateContract, logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

/*
func (c *BoxUsecase) PatientPaymentFormSignRequests(folderId string, signerEmail string, copyEmail string, fullName string) (res *string, contractUniqId string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if lib.IsProd() {
		params.Set("email_subject", "Patient Payment Form")
	} else {
		params.Set("email_subject", "Patient Payment Form")
	}
	params.Set("are_reminders_enabled", true)
	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": signerEmail, // lialing@foxmail.com
			"order": 1,
		},
		{
			"role":  "final_copy_reader",
			"email": copyEmail, // liaogling@gmail.com
			"order": 2,
		},
	})

	params.Set("prefill_tags", lib.TypeList{
		{
			"document_tag_id": "clientName",
			"text_value":      fullName,
		},
		{
			"document_tag_id": "ClientNameSign",
			"text_value":      fullName,
		},
	})
	params.Set("template_id", "e282a4b2-be3d-4970-978b-2b112abcff0a") // a0a2df1f-3fd4-42f2-8881-6d6b1df55ea2
	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})

	//if lib.DebugMedTeamFormBoxSign {
	//	testRes := `{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Release of Information","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"ab0991e8-2404-44d8-b179-670feb859633","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257082839680","etag":"0","type":"folder","sequence_id":"0","name":"Li, Shi#5040"},"name":"Release of Information.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1541669723989","etag":"0","type":"file","sequence_id":"0","name":"Release of Information.pdf","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","file_version":{"id":"1693470407989","type":"file_version","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"c4c50315-7e5a-44af-9d5b-ff0e4ddbad48","external_system_name":null}`
	//	return to.Ptr(testRes), "ab0991e8-2404-44d8-b179-670feb859633", nil
	//}
	res, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "PatientPaymentFormSignRequests"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_CreatePatientPaymentFormContract, logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

*/

// ReleaseOfInformationSignRequests
// 测试数据： https://veteranbenefitscenter.app.box.com/file/1494561604477
// https://veteranbenefitscenter.app.box.com/folder/257590681761
// ab0991e8-2404-44d8-b179-670feb859633 {"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Release of Information","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"ab0991e8-2404-44d8-b179-670feb859633","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257082839680","etag":"0","type":"folder","sequence_id":"0","name":"Li, Shi#5040"},"name":"Release of Information.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1541669723989","etag":"0","type":"file","sequence_id":"0","name":"Release of Information.pdf","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","file_version":{"id":"1693470407989","type":"file_version","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"c4c50315-7e5a-44af-9d5b-ff0e4ddbad48","external_system_name":null}
/*
func (c *BoxUsecase) ReleaseOfInformationSignRequests(folderId string, signerEmail string, copyEmail string, prefillTags lib.TypeList) (res *string, contractUniqId string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if lib.IsProd() {
		params.Set("email_subject", "Release of Information")
	} else {
		params.Set("email_subject", "Release of Information")
	}
	params.Set("are_reminders_enabled", true)

	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": signerEmail,
			"order": 1,
		},
		{
			"role":  "final_copy_reader",
			"email": copyEmail,
			"order": 2,
		},
	})

	/*
		params.Set("signers", lib.TypeList{
			{
				"role":  "signer",
				"email": "lialing@foxmail.com",
				"order": 1,
			},
			{
				"role":  "final_copy_reader",
				"email": "liaogling@gmail.com",
				"order": 2,
			},
		})
*/

//params.Set("prefill_tags", prefillTags)
/*params.Set("prefill_tags", lib.TypeList{
	{
		"document_tag_id": "clientName",
		"text_value":      GenFullName("Gary", "Liao"),
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
		"text_value":      "1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)",
	},
})*/
//	params.Set("template_id", "c4c50315-7e5a-44af-9d5b-ff0e4ddbad48") // a0a2df1f-3fd4-42f2-8881-6d6b1df55ea2
//	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})
//
//	if lib.DebugMedTeamFormBoxSign {
//		//
//		testRes := `{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Release of Information","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"ab0991e8-2404-44d8-b179-670feb859633","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257082839680","etag":"0","type":"folder","sequence_id":"0","name":"Li, Shi#5040"},"name":"Release of Information.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1541669723989","etag":"0","type":"file","sequence_id":"0","name":"Release of Information.pdf","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","file_version":{"id":"1693470407989","type":"file_version","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"c4c50315-7e5a-44af-9d5b-ff0e4ddbad48","external_system_name":null}`
//		return to.Ptr(testRes), "ab0991e8-2404-44d8-b179-670feb859633", nil
//	}
//	res, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
//
//	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "ReleaseOfInformationSignRequests"), time.Now(), 1)
//
//	logNotes := make(lib.TypeMap)
//	logNotes.Set("params", params)
//	logNotes.Set("error", err)
//	logNotes.Set("res", res)
//
//	log := GenLog(0, Log_FromType_Box_ReleaseOfInformationSignRequests, logNotes.ToString())
//	er := c.CommonUsecase.DB().Save(&log).Error
//	if er != nil {
//		c.log.Error(er)
//	}
//
//	if err != nil {
//		return nil, "", err
//	}
//	if res == nil {
//		return nil, "", errors.New("res is nil")
//	}
//	resMap := lib.ToTypeMapByString(*res)
//	if resMap.GetString("id") == "" {
//		return nil, "", errors.New("res is error: " + *res)
//	}
//	return res, resMap.GetString("id"), nil
//}*/

// MedicalTeamFormsSignRequests
// 测试数据： https://veteranbenefitscenter.app.box.com/file/1550028092926
// https://veteranbenefitscenter.app.box.com/folder/257590681761
// ab0991e8-2404-44d8-b179-670feb859633 {"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"Release of Information","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"ab0991e8-2404-44d8-b179-670feb859633","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257082839680","etag":"0","type":"folder","sequence_id":"0","name":"Li, Shi#5040"},"name":"Release of Information.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1541669723989","etag":"0","type":"file","sequence_id":"0","name":"Release of Information.pdf","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","file_version":{"id":"1693470407989","type":"file_version","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"c4c50315-7e5a-44af-9d5b-ff0e4ddbad48","external_system_name":null}
func (c *BoxUsecase) MedicalTeamFormsSignRequests(folderId string, signerEmail string, copyEmail string, prefillTags lib.TypeList, version string, boxSignTmpId string) (res *string, contractUniqId string, err error) {

	if boxSignTmpId == "" {
		return nil, "", errors.New("boxSignTmpId is empty")
	}
	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if configs.IsProd() {
		params.Set("email_subject", "VBC: Medical Team Forms")
	} else {
		params.Set("email_subject", "VBC: Medical Team Forms")
	}
	params.Set("are_reminders_enabled", true)

	params.Set("signers", lib.TypeList{
		{
			"role":  "signer",
			"email": signerEmail,
			"order": 1,
		},
		{
			"role":  "final_copy_reader",
			"email": copyEmail,
			"order": 2,
		},
	})

	params.Set("prefill_tags", prefillTags)
	//if version == MedicalTeamFormsV2 {
	//	params.Set("template_id", "2aac3aed-358f-491a-a459-1b9d51bcc154")
	//} else {
	//	params.Set("template_id", "aaaf350c-9f4f-437c-a8ce-ba5602b360b3")
	//}

	params.Set("template_id", boxSignTmpId)

	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})

	//if lib.DebugMedTeamFormBoxSign {
	//	//
	//	testRes := `{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"VBC: Medical Team Forms","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"0328b7ed-d475-44a6-aaa1-dafcc97baf07","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257590681761","etag":"0","type":"folder","sequence_id":"0","name":"Demo, test#5007"},"name":"Medical Team Forms.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1550028092926","etag":"0","type":"file","sequence_id":"0","name":"Medical Team Forms.pdf","sha1":"7eb34c27b1dd17f7ed40466e250f3566b8604081","file_version":{"id":"1702881668926","type":"file_version","sha1":"7eb34c27b1dd17f7ed40466e250f3566b8604081"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"aaaf350c-9f4f-437c-a8ce-ba5602b360b3","external_system_name":null}`
	//	return to.Ptr(testRes), "0328b7ed-d475-44a6-aaa1-dafcc97baf07", nil
	//}
	res, _, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "ReleaseOfInformationSignRequests"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_MedicalTeamFormsSignRequests, logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

func (c *BoxUsecase) MedicalTeamFormsSignRequestsWithoutTemplate(folderId string, signerEmail string, copyEmail string, copyEmail2 string, pdfSourceBoxFileId string) (res *string, contractUniqId string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, "", err
	}
	api := fmt.Sprintf("%s/2.0/sign_requests", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	if configs.IsProd() {
		params.Set("email_subject", "VBC: Medical Team Forms")
	} else {
		params.Set("email_subject", "VBC: Medical Team Forms")
	}
	params.Set("are_reminders_enabled", true)

	if copyEmail2 == copyEmail {
		copyEmail2 = ""
	}

	paramsList := lib.TypeList{
		{
			"role":  "signer",
			"email": signerEmail,
			"order": 1,
		},
		{
			"role":  "final_copy_reader",
			"email": copyEmail,
			"order": 2,
		},
	}
	if copyEmail2 != "" {
		paramsList = append(paramsList, lib.TypeMap{
			"role":  "final_copy_reader",
			"email": copyEmail2,
			"order": 3,
		})
	}

	params.Set("signers", paramsList)
	//params.Set("prefill_tags", prefillTags)
	//if version == MedicalTeamFormsV2 {
	//	params.Set("template_id", "2aac3aed-358f-491a-a459-1b9d51bcc154")
	//} else {
	//	params.Set("template_id", "aaaf350c-9f4f-437c-a8ce-ba5602b360b3")
	//}
	//pdfSourceBoxFileId
	params.Set("source_files", lib.TypeList{lib.TypeMap{
		"type": "file",
		"id":   pdfSourceBoxFileId,
	}})
	//params.Set("template_id", boxSignTmpId)
	params.Set("parent_folder", lib.TypeMap{"type": "folder", "id": folderId})
	//return

	//if lib.DebugMedTeamFormBoxSign {
	//	//
	//	testRes := `{"is_document_preparation_needed":false,"redirect_url":null,"declined_redirect_url":null,"are_text_signatures_enabled":true,"signature_color":null,"is_phone_verification_required_to_view":false,"email_subject":"VBC: Medical Team Forms","email_message":"Please sign this document.\n\nKind regards,\n\nVBC Team (info@vetbenefitscenter.com)","are_reminders_enabled":true,"signers":[{"email":"info@vetbenefitscenter.com","role":"final_copy_reader","is_in_person":false,"order":0,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"lialing@foxmail.com","role":"signer","is_in_person":false,"order":1,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false},{"email":"liaogling@gmail.com","role":"final_copy_reader","is_in_person":false,"order":2,"verification_phone_number":null,"embed_url_external_user_id":null,"redirect_url":null,"declined_redirect_url":null,"login_required":false,"has_viewed_document":false,"signer_decision":null,"signer_group_id":null,"inputs":[],"embed_url":null,"iframeable_embed_url":null,"suppress_notifications":false}],"id":"0328b7ed-d475-44a6-aaa1-dafcc97baf07","prefill_tags":[{"document_tag_id":"clientName","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"ClientNameSign","text_value":"Gary Liao","checkbox_value":null,"date_value":null},{"document_tag_id":"dob","text_value":"Jun. 8, 1984","checkbox_value":null,"date_value":null},{"document_tag_id":"ssn","text_value":"573-79-7392","checkbox_value":null,"date_value":null},{"document_tag_id":"phone","text_value":"619-948-5488","checkbox_value":null,"date_value":null},{"document_tag_id":"email","text_value":"lialing@foxmail.com","checkbox_value":null,"date_value":null},{"document_tag_id":"address","text_value":"000 W 9th Avenue","checkbox_value":null,"date_value":null},{"document_tag_id":"location","text_value":"Escondido, CA 92000","checkbox_value":null,"date_value":null},{"document_tag_id":"itf","text_value":"Mar. 1, 2025","checkbox_value":null,"date_value":null},{"document_tag_id":"privateExamsNeeded","text_value":"1_Morning headaches (increase)\nRight knee sprain (str)\nLeft knee pain secondary to right knee sprain (opinion)\nErectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n2_Morning headaches (increase)\n4_Right knee sprain (str)\n5_Left knee pain secondary to right knee sprain (opinion)\n6_Erectile dysfunction secondary to Major depression disorder with anxious distress (opinion)\n7_Morning headaches (increase)\n8_Right knee sprain (str)\n9_Left knee pain secondary to right knee sprain (opinion)","checkbox_value":null,"date_value":null}],"days_valid":0,"prepare_url":null,"source_files":[],"parent_folder":{"id":"257590681761","etag":"0","type":"folder","sequence_id":"0","name":"Demo, test#5007"},"name":"Medical Team Forms.pdf","external_id":null,"type":"sign-request","signing_log":null,"status":"created","sign_files":{"files":[{"id":"1550028092926","etag":"0","type":"file","sequence_id":"0","name":"Medical Team Forms.pdf","sha1":"7eb34c27b1dd17f7ed40466e250f3566b8604081","file_version":{"id":"1702881668926","type":"file_version","sha1":"7eb34c27b1dd17f7ed40466e250f3566b8604081"}}],"is_ready_for_download":true},"auto_expire_at":null,"template_id":"aaaf350c-9f4f-437c-a8ce-ba5602b360b3","external_system_name":null}`
	//	return to.Ptr(testRes), "0328b7ed-d475-44a6-aaa1-dafcc97baf07", nil
	//}
	res, _, err = lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "ReleaseOfInformationSignRequests"), time.Now(), 1)

	logNotes := make(lib.TypeMap)
	logNotes.Set("params", params)
	logNotes.Set("error", err)
	logNotes.Set("res", res)

	log := GenLog(0, Log_FromType_Box_MedicalTeamFormsSignRequests+"V2", logNotes.ToString())
	er := c.CommonUsecase.DB().Save(&log).Error
	if er != nil {
		c.log.Error(er)
	}

	if err != nil {
		return nil, "", err
	}
	if res == nil {
		return nil, "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	if resMap.GetString("id") == "" {
		return nil, "", errors.New("res is error: " + *res)
	}
	return res, resMap.GetString("id"), nil
}

const (
	MedicalTeamFormsV2 = "v2"
)

func (c *BoxUsecase) CreateFolder(destFolderName string, destParentId string) (boxFolderId string, err error) {

	token, err := c.Token()
	if err != nil {
		return "", err
	}

	api := fmt.Sprintf("%s/2.0/folders", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	params.Set("name", destFolderName)
	params.Set("parent.id", destParentId)
	r, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CreateFolder"), time.Now(), 1)
	if err != nil {
		c.log.Error(err, InterfaceToString(r), ":", destFolderName, ":", destParentId)
		return "", err
	}
	if r == nil {
		return "", errors.New("r is nil")
	}
	a := lib.ToTypeMapByString(*r)

	return a.GetString("id"), nil
}

/*
	func (c *AsanaUsecase) UpdateFolderName(gid string, customFields lib.TypeMap) (*string, error) {
		url := c.conf.Asana.GetApiHost()
		url = fmt.Sprintf("%s/tasks/%s", url, gid)

		destTypeMap := make(lib.TypeMap)
		destTypeMap.Set("data.custom_fields", customFields)

		r, er := lib.HTTPJsonWithHeaders(http.MethodPut, url, destTypeMap.ToBytes(), map[string]string{
			"authorization": "Bearer " + c.conf.Asana.Pat,
		})
		return r, er
	}
*/
func (c *BoxUsecase) UpdateFolderName(folderId string, newFolderName string) (res *string, err error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/folders/%s", c.conf.Box.ApiUrl, folderId)
	params := make(lib.TypeMap)
	params.Set("name", newFolderName)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "UpdateFolderName"), time.Now(), 1)
	res, _, err = lib.HTTPJsonWithHeaders(http.MethodPut, api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	return
}

func (c *BoxUsecase) MoveFolderName(folderId string, newFolderName string, parentFolderId string) (res *string, err error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/folders/%s", c.conf.Box.ApiUrl, folderId)
	params := make(lib.TypeMap)
	params.Set("name", newFolderName)
	params.Set("parent.id", parentFolderId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "MoveFolderName"), time.Now(), 1)
	res, _, err = lib.HTTPJsonWithHeaders(http.MethodPut, api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	return
}

func (c *BoxUsecase) CopyFolder(sourceFolderId string, destFolderName string, destParentId string) (boxFolderId string, httpCode int, err error) {

	token, err := c.Token()
	if err != nil {
		return "", 0, err
	}

	api := fmt.Sprintf("%s/2.0/folders/%s/copy", c.conf.Box.ApiUrl, sourceFolderId)
	params := make(lib.TypeMap)
	params.Set("name", destFolderName)
	params.Set("parent.id", destParentId)
	r, httpCode, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CopyFolder"), time.Now(), 1)
	if err != nil {
		c.log.Error(err, r)
		return "", httpCode, err
	}
	if r == nil {
		return "", httpCode, errors.New("r is nil")
	}
	a := lib.ToTypeMapByString(*r)

	return a.GetString("id"), httpCode, nil
}

func (c *BoxUsecase) CopyFileNewFileNameReturnFileId(fileId string, newFileName string, destFolderId string) (resultFileId string, err error) {
	res, err := c.CopyFileNewFileName(fileId, newFileName, destFolderId)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("res is nil")
	}
	resMap := lib.ToTypeMapByString(*res)
	resultFileId = resMap.GetString("id")
	if resultFileId == "" {
		return "", errors.New("resultFileId is empty")
	}
	return
}
func (c *BoxUsecase) CopyFileNewFileName(fileId string, newFileName string, destFolderId string) (*string, error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s/copy", c.conf.Box.ApiUrl, fileId)
	params := make(lib.TypeMap)
	params.Set("name", newFileName)
	params.Set("parent.id", destFolderId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CopyFile"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

func (c *BoxUsecase) CopyFile(fileId string, destFolderId string) (*string, int, error) {

	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s/copy", c.conf.Box.ApiUrl, fileId)
	params := make(lib.TypeMap)
	params.Set("parent.id", destFolderId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CopyFile"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	if err == nil {
		return res, config_box.HttpCode_200, nil
	}
	if res != nil {
		resMap := lib.ToTypeMapByString(*res)
		return res, int(resMap.GetInt("status")), err
	}
	return res, 0, err
}

func (c *BoxUsecase) MoveFile(fileId string, destFolderId string) (*string, error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s", c.conf.Box.ApiUrl, fileId)
	params := make(lib.TypeMap)
	params.Set("parent.id", destFolderId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "MoveFile"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("PUT", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

func (c *BoxUsecase) MoveFileWithNewName(fileId string, destFolderId string, newName string) (*string, error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s", c.conf.Box.ApiUrl, fileId)
	params := make(lib.TypeMap)
	params.Set("name", newName)
	params.Set("parent.id", destFolderId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "MoveFile"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("PUT", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

// CollaborationsByBoxUserId 把员工加入协作 返回好像是：协作的ID：64086181594
func (c *BoxUsecase) CollaborationsByBoxUserId(folderId string, boxUserId string) (id string, httpCode int, err error) {

	token, err := c.Token()
	if err != nil {
		return "", 0, err
	}

	api := fmt.Sprintf("%s/2.0/collaborations", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	params.Set("item.type", "folder")
	params.Set("item.id", folderId)
	params.Set("accessible_by.type", "user")
	params.Set("accessible_by.id", boxUserId)
	params.Set("role", "editor")
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "Collaborations"), time.Now(), 1)
	r, httpCode, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		str := ""
		if r != nil {
			str = *r
		}
		c.log.Error(err, str)
		return "", 0, err
	}
	if r == nil {
		return "", 0, errors.New("r is nil")
	}
	a := lib.ToTypeMapByString(*r)

	/*
		{
		    "type": "error",
		    "status": 400,
		    "code": "user_already_collaborator",
		    "help_url": "http://developers.box.com/docs/#errors",
		    "message": "User is already a collaborator",
		    "request_id": "vz6y9hhlpvgj1ek0"
		}
	*/

	return a.GetString("id"), httpCode, nil
}

func (c *BoxUsecase) Collaborations(folderId string, email string) (id string, err error) {

	token, err := c.Token()
	if err != nil {
		return "", err
	}

	api := fmt.Sprintf("%s/2.0/collaborations", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	params.Set("item.type", "folder")
	params.Set("item.id", folderId)
	params.Set("accessible_by.type", "user")
	params.Set("accessible_by.login", email)
	params.Set("role", "viewer uploader")
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "Collaborations"), time.Now(), 1)
	r, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		str := ""
		if r != nil {
			str = *r
		}
		c.log.Error(err, str)
		return "", err
	}
	if r == nil {
		return "", errors.New("r is nil")
	}
	a := lib.ToTypeMapByString(*r)

	/*
		{
		    "type": "error",
		    "status": 400,
		    "code": "user_already_collaborator",
		    "help_url": "http://developers.box.com/docs/#errors",
		    "message": "User is already a collaborator",
		    "request_id": "vz6y9hhlpvgj1ek0"
		}
	*/

	return a.GetString("id"), nil
}

func (c *BoxUsecase) UploadFile(folderId string, reader io.Reader, fileName string) (fileId string, err error) {

	token, err := c.Token()
	if err != nil {
		return "", err
	}
	url := c.conf.Box.UploadUrl + "/api/2.0/files/content"
	attributes := make(lib.TypeMap)
	attributes.Set("name", fileName)
	attributes.Set("parent.id", folderId)

	values := []*lib.UploadReader{
		lib.NewUploadReader("attributes", strings.NewReader(InterfaceToString(attributes)), ""),
		lib.NewUploadReader("file", reader, fileName), // lets assume its this file
	}
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "UploadFile"), time.Now(), 1)
	res, err := lib.PostUpload(url, values, map[string]string{"authorization": "Bearer " + token})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("UploadFile: res is nil")
	}
	typeMap := lib.ToTypeMapByString(*res)
	entries := typeMap.GetTypeList("entries")
	if len(entries) > 0 {
		return entries[0].GetString("id"), nil
	} else {
		return "", errors.New("UploadFile: entries is wrong")
	}
}

func (c *BoxUsecase) WebhooksList() (*string, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/webhooks", c.conf.Box.ApiUrl)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "WebhooksList"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

/*
curl -i -X POST "https://api.box.com/2.0/webhooks" \
     -H "authorization: Bearer <ACCESS_TOKEN>" \
     -H "content-type: application/json" \
     -d '{
       "target": {
         "id": "21322",
         "type": "file"
       },
       "address": "https://example.com/webhooks",
       "triggers": [
         "FILE.PREVIEWED"
       ]
     }'
*/

func (c *BoxUsecase) CreateWebhooks(folderId string, email string) (id string, err error) {

	token, err := c.Token()
	if err != nil {
		return "", err
	}

	api := fmt.Sprintf("%s/2.0/collaborations", c.conf.Box.ApiUrl)
	params := make(lib.TypeMap)
	params.Set("item.type", "folder")
	params.Set("item.id", folderId)
	params.Set("accessible_by.type", "user")
	params.Set("accessible_by.login", email)
	params.Set("role", "uploader")
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "CreateWebhooks"), time.Now(), 1)
	r, _, err := lib.HTTPJsonWithHeaders("POST", api, params.ToBytes(), map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		c.log.Error(err, r)
		return "", err
	}
	if r == nil {
		return "", errors.New("r is nil")
	}
	a := lib.ToTypeMapByString(*r)

	/*
		{
		    "type": "error",
		    "status": 400,
		    "code": "user_already_collaborator",
		    "help_url": "http://developers.box.com/docs/#errors",
		    "message": "User is already a collaborator",
		    "request_id": "vz6y9hhlpvgj1ek0"
		}
	*/

	return a.GetString("id"), nil
}

// ListItemsInFolderFormat
// [{"etag":"2","id":"263407033477","name":"Abutin, Niko Ralphluis","sequence_id":"2","type":"folder"},{"etag":"1","id":"263407629244","name":"Acuario, Edralin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409690475","name":"Albrecht, Keith Richard","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407432521","name":"Alcantar, Francisco Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409699906","name":"Alexander, Keith David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409220477","name":"Alexander, Troy Don","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409383635","name":"Allen, Robert Joseph","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408142981","name":"Ancho, Romulo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409340623","name":"Anderson, Jilleah","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408260360","name":"Anderson, Trever Shaw","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409604748","name":"Andrews, Jamaal","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085591","name":"Angulo, Don Clark","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408638849","name":"Arellano, Hector Gibram","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409608463","name":"Ayala, Joe","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407544468","name":"Ayuyao, Bernard Elijah","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407775301","name":"Bailey, Bernard","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408325137","name":"Baker, Jeffrey Stephen Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407895271","name":"Balavram, Jason Sy","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408186009","name":"Baldemeca, Noel Genetia","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407002935","name":"Ballesteros, Judeasar Galapon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408572372","name":"Banuex, Noemi","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409027699","name":"Barajas, Emery Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407850270","name":"Barbosa, Brendan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407864388","name":"Barnes, Santana Venique","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408646454","name":"Barnett, Jovan Eric","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408332367","name":"Battle, Rodney Terez","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408192682","name":"Bautista, Alexander Clement","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407518222","name":"Becerra, Jonathan Contreras","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408711489","name":"Beeler, Rebecca #5019","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408505251","name":"Blaine, Christopher Charles","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408101803","name":"Boden, Boyce Robert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408584289","name":"Bolino, Louis Alto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408994347","name":"Briggs, Scott Christopher","sequence_id":"1","type":"folder"},{"etag":"0","id":"263966723631","name":"Brooks, Roy #5069","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408480938","name":"Brown, Terence Lewis","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408941892","name":"Burgess, Dominique Farrell","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407893513","name":"Camacho, Pedro","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408432830","name":"Campbell, Dailyn","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408809561","name":"Canseco, Manuel Valiente","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408946196","name":"Carr, Christopher","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407168436","name":"Carreon, Lovelito Flores","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408315542","name":"Carrillo, Adrian","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407400559","name":"Carter, Gabrielle Alexandria","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409575210","name":"Castelluccio, Dillon Thomas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409090018","name":"Castillo, Jacinto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408282334","name":"Castro, Roman Lucas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408953729","name":"Cavaliere, Robynn","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408315294","name":"Chacon, Guy #5052","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409193805","name":"Claire, Seth Alexander","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409058831","name":"Cobian Jr., Jose","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409335804","name":"Coley, Alvin Lorenzo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407511190","name":"Collins, Samantha Rose","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408706463","name":"Cook, Robert #5020","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408775705","name":"D'Alessandro, James Scott","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408622377","name":"De Leon, Cecilia","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409573530","name":"DelPrete, Desiree","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408560027","name":"Demps, Kelvin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409486296","name":"Deocampo, Teddy Baysan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408562555","name":"Devine, Justin Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409697401","name":"DiBenedetto, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408598259","name":"Dickerson, Ramon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408170942","name":"Dickey, James #5005","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408262134","name":"Dishmon, Varian Dione","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407789701","name":"Dodd, Brent","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408883616","name":"Dukes, Spencer Patrick","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409395451","name":"Dunkin, John Steven","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408255849","name":"Edwards, Addison Chase","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407679255","name":"FaisonLanier, Terrica","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408082641","name":"Fannin, Quentin Dekote","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407698819","name":"Farmer, Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407825270","name":"Flores, Jacinth Aaron","sequence_id":"1","type":"folder"},{"etag":"2","id":"263409236508","name":"Fowler, Casey #345","sequence_id":"2","type":"folder"},{"etag":"1","id":"263409256573","name":"Fowler, Jerry Joseph","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407897776","name":"Franco, Alec Robert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409003856","name":"Galac, Cesar","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408776073","name":"Garcia, Rey #5041","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408716059","name":"Garlejo, Jason Doctolero","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409049595","name":"Gilmore, Earl Glenn Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408759623","name":"Gonzales, Beau Matthew","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407379970","name":"Goodson, Augustus Ivan IV","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409028000","name":"Green, Antionette","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407516742","name":"Green, Donnell","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409561453","name":"Green, Sandra","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407261628","name":"Griffin, Major Pete III","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407014935","name":"Haley, George Walter Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407237574","name":"Harris, Debra Lee","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408234122","name":"Herron, Leslie Rhea","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409282194","name":"Hiers, Tommy Lamar Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408977216","name":"Ho, Duyet #5043","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407717756","name":"Houston, Keith Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407411164","name":"Howard, Nicole Marie","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408903510","name":"Huerta Jr., Edward David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408925455","name":"Huerta Jr., Edward David #96","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407497275","name":"Huerta Sr., Edward David #260","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407386390","name":"Hutchinson, Derek","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408325009","name":"Ibarrondo, William","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408646453","name":"Ibarrondo, William Basean","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407665390","name":"Inzer, Russell Dustin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408617842","name":"Jacob, Melvin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408202281","name":"James, Dillon Randall","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409584864","name":"Johnson, Christopher #5009","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409280743","name":"Johnson, Jermaine","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409239761","name":"Johnson, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407513317","name":"Johnson, Robin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408383080","name":"Jones, Cyril Evans Jr.","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407472517","name":"Jones, Naurice","sequence_id":"1","type":"folder"},{"etag":"0","id":"264597201436","name":"Kallmeyer, Michael #5072","sequence_id":"0","type":"folder"},{"etag":"1","id":"263409011236","name":"Keith, Kristopher #5025","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407636399","name":"Keller, Tony Jordan","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407378937","name":"Kennedy, Payton Gary","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407672521","name":"Kubas III, William Philip","sequence_id":"1","type":"folder"},{"etag":"0","id":"264171116583","name":"Kunkowski, James #5070","sequence_id":"0","type":"folder"},{"etag":"1","id":"263407746455","name":"Lane, Michael John","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408048756","name":"Lang, Shanise Eileen","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407759064","name":"Lastrella, Amando Llorin","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407079735","name":"Laxa, Eduardo Dulu","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407578819","name":"Le, Khanh Si","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409409666","name":"Lepe, Christopher","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407774875","name":"Liggans, Nyjerus Lavondai Onijar","sequence_id":"1","type":"folder"},{"etag":"0","id":"264108189026","name":"Long, Jonathan #5068","sequence_id":"0","type":"folder"},{"etag":"0","id":"263965342942","name":"Luna, Larzen #355","sequence_id":"0","type":"folder"},{"etag":"0","id":"264296809567","name":"Maldonado, Jimmy #5060","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408704030","name":"Mangra, Robbi #5023","sequence_id":"1","type":"folder"},{"etag":"2","id":"263409358994","name":"Marcial, Janry #351","sequence_id":"2","type":"folder"},{"etag":"1","id":"263408360692","name":"Mendez, Marco Antonio","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407919167","name":"Montoya, Mario","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408174156","name":"Montoya, Mary","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407556962","name":"Moreno, Michael Joey","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408493117","name":"Morris, Lee Roy","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407593367","name":"Murrell, Brashaad","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408197832","name":"Myers, Luis Anthony","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408754259","name":"Nellis, Lawrence #5024","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085859","name":"Netemeyer, Aaron","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409335873","name":"Newman, James Wesley","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407708055","name":"Olivetti, Anthony Ryan Borja","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409189404","name":"Orias, Ricardo","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407676838","name":"Padua, Michael Daniel","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409107750","name":"Peca, Jason","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407622264","name":"Perez, Joaquin Xavier","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409369154","name":"Perry, Fred Douglas","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408536697","name":"Petit-Frere, Alexandre Freud","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408613185","name":"Pharnes, Eric Dwayne","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408872130","name":"Pierre, Gilbert","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408335098","name":"Prado, Martius Oris","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409357404","name":"Pratko, Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407645792","name":"Provasek, Jared #5046","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407539990","name":"Ralat, Carlos Alberto","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407549727","name":"Reynoso, Algis #5036","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407439423","name":"Rios, Jonathan David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408483369","name":"Rivera, Louie","sequence_id":"1","type":"folder"},{"etag":"2","id":"263408128360","name":"Rosales, Juan #66","sequence_id":"2","type":"folder"},{"etag":"1","id":"263409671231","name":"Rutledge, Ronnie #5042","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408531552","name":"Salb, Austin Reid","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407996012","name":"Santillan-Mondaca, Kristy #5050","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407479863","name":"Sayles, Matthew Evans","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408282344","name":"Serrano, Ronald","sequence_id":"1","type":"folder"},{"etag":"0","id":"263964486190","name":"Sese, Maria #5002","sequence_id":"0","type":"folder"},{"etag":"1","id":"263407777243","name":"Shelrud, Cierra Kay","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408291548","name":"Sida, Andrew Michael","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409143090","name":"Slater, Jamie","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408173982","name":"Smith, Andrew","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409109595","name":"Smith, Andrew #5011","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409138194","name":"Smith, Austin Cole","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408466241","name":"Smith, Christopher Manuel #5026","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408730987","name":"Smith, Max","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409028559","name":"Smith, Rhyheime","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408486215","name":"Smith, Zane Eugene","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407441338","name":"Smolinski, Donald Jay","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409085790","name":"Stacks, Taylor","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408843374","name":"Stewart, Robert #5038","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408879881","name":"Stuart, James Francis","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409189108","name":"Summer, Aaron Michael Blake","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408183176","name":"Sutton, Derrick","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408877126","name":"Tanquilut, Remigio Dimalanta","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408462333","name":"Taylor, Bobbee Nykole","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408290941","name":"Terrell, Conley II","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408989699","name":"TestLi, TestShi #5057","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408579201","name":"TestLi, TestShi #5058","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407556960","name":"TestLiao, TestGary #5061","sequence_id":"1","type":"folder"},{"etag":"0","id":"263470914421","name":"TestLiao, TestGary #5064","sequence_id":"0","type":"folder"},{"etag":"1","id":"263409010923","name":"Thompson, Jason Scott","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408644030","name":"Thrower, Tony Lee","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408560526","name":"Tran, Danny Minh","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408893699","name":"Tran, Joanne Perea","sequence_id":"1","type":"folder"},{"etag":"1","id":"263409387704","name":"Tran, Kenny #5045","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408661172","name":"Valdez, Jacob","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408233836","name":"Valdez, Joshua Raul","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407744765","name":"Valli, Matthew Lawrence","sequence_id":"1","type":"folder"},{"etag":"2","id":"263408943728","name":"Valli, Ronald #391","sequence_id":"2","type":"folder"},{"etag":"1","id":"263408636453","name":"Vargas, Gabrian","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408008709","name":"Velasquez, Jose David","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407730033","name":"Walker, Ronald Stanley","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408847871","name":"Warren, Olga","sequence_id":"1","type":"folder"},{"etag":"0","id":"264219349997","name":"Warren, Olga #5071","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408262984","name":"Watts, Patrick Levon","sequence_id":"1","type":"folder"},{"etag":"1","id":"263407622121","name":"Webster, Craig","sequence_id":"1","type":"folder"},{"etag":"1","id":"263408430635","name":"West, Melissa Theresa","sequence_id":"1","type":"folder"},{"etag":"0","id":"263723175862","name":"Westerveld, Michael #5017","sequence_id":"0","type":"folder"},{"etag":"1","id":"263408356225","name":"White, Michael #5015","sequence_id":"1","type":"folder"},{"etag":"0","id":"263570678121","name":"Wirth, David #5029","sequence_id":"0","type":"folder"}]
func (c *BoxUsecase) ListItemsInFolderFormat(folderId string) (lib.TypeList, error) {
	res, err := c.ListItemsInFolder(folderId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("ListItemsInFolderFormat: res is nil.")
	}
	resMap := lib.ToTypeMapByString(*res)
	return resMap.GetTypeList("entries"), nil
}

func (c *BoxUsecase) ListItemsInFolder(folderId string) (res *string, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/2.0/folders/%s/items", c.conf.Box.ApiUrl, folderId)
	queryParams := make(url.Values)
	queryParams.Add("limit", "1000")
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "ListItemsInFolder"), time.Now(), 1)
	r, _, err := lib.RequestGet(api, queryParams, map[string]string{"Authorization": "Bearer " + token})
	return r, err
}

func (c *BoxUsecase) DeleteFolder(folderId string, recursive bool) (*string, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/folders/%s", c.conf.Box.ApiUrl, folderId)
	if recursive {
		api = api + "?recursive=true"
	}
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "DeleteFolder"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("DELETE", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

func (c *BoxUsecase) DeleteFile(fileId string) (*string, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/files/%s", c.conf.Box.ApiUrl, fileId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "DeleteFile"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("DELETE", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return res, err
}

// UploadFileVersion 上传文件，产生新版本
func (c *BoxUsecase) UploadFileVersion(fileId string, fileReader io.Reader) (*string, error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/api/2.0/files/%s/content", c.conf.Box.UploadUrl, fileId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "UploadFileVersion"), time.Now(), 1)

	return lib.PostUpload(api, []*lib.UploadReader{
		lib.NewUploadReader("file", fileReader, "a.pdf"),
		//lib.NewUploadReader("attributes", strings.NewReader("{\"name\":\"a_newname.pdf\", \"parent\":{\"id\":\"264924117433\"}}"), ""),
		//lib.NewUploadReader("attributes", strings.NewReader("{\"name\":\"b.pdf\"}"), ""), // 改名为b.pdf
	}, map[string]string{"Authorization": "Bearer " + token})
}

// UploadFileVersionWithNewFileName 上传文件，产生新版本
func (c *BoxUsecase) UploadFileVersionWithNewFileName(fileId string, fileReader io.Reader, newFileName string) (*string, error) {

	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/api/2.0/files/%s/content", c.conf.Box.UploadUrl, fileId)
	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "UploadFileVersion"), time.Now(), 1)

	return lib.PostUpload(api, []*lib.UploadReader{
		lib.NewUploadReader("file", fileReader, "a.pdf"),
		//lib.NewUploadReader("attributes", strings.NewReader("{\"name\":\"a_newname.pdf\", \"parent\":{\"id\":\"264924117433\"}}"), ""),
		lib.NewUploadReader("attributes", strings.NewReader("{\"name\":\""+newFileName+"\"}"), ""), // 改名为b.pdf
	}, map[string]string{"Authorization": "Bearer " + token})
}

// DownloadFile 下载文件
func (c *BoxUsecase) DownloadFile(fileId string, fileVersion string) (io.ReadCloser, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/files/%s/content", c.conf.Box.ApiUrl, fileId)

	queryParams := make(url.Values)
	if fileVersion != "" {
		queryParams.Add("version", fileVersion)
	}
	if len(queryParams) > 0 {
		api += "?" + queryParams.Encode()
	}
	//return nil, nil
	response, err := lib.RequestDoTimeout("GET", api, nil, map[string]string{"Authorization": "Bearer " + token}, time.Minute*10)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return response.Body, nil
	}
	return nil, errors.New("DownloadFile: response.StatusCode: " + InterfaceToString(response.StatusCode))
}

func (c *BoxUsecase) GetFileInfoForTypeMap(fileId string) (res lib.TypeMap, httpCode int, err error) {
	data, httpCode, err := c.GetFileInfo(fileId)
	if err != nil {
		return nil, httpCode, err
	}
	if data == nil {
		return nil, httpCode, errors.New("data is nil")
	}
	res = lib.ToTypeMapByString(*data)
	return res, httpCode, nil
}

// GetFileInfo 获取文件的信息
// {"type":"file","id":"1554567330874","file_version":{"type":"file_version","id":"1708008340474","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"},"sequence_id":"1","etag":"1","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","name":"b.pdf","description":"","size":182138,"path_collection":{"total_count":7,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"247457173873","sequence_id":"1","etag":"1","name":"Testing"},{"type":"folder","id":"255166311971","sequence_id":"1","etag":"1","name":"Test Clients"},{"type":"folder","id":"264686374897","sequence_id":"2","etag":"2","name":"[PROD]VBC - TestLiao, TestGary #5076"},{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"},{"type":"folder","id":"268258622342","sequence_id":"0","etag":"0","name":"leve1_folder"}]},"created_at":"2024-06-08T03:41:17-07:00","modified_at":"2024-06-08T03:41:17-07:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-05-18T19:35:02-07:00","content_modified_at":"2024-05-18T19:35:02-07:00","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"268258622342","sequence_id":"0","etag":"0","name":"leve1_folder"},"item_status":"active"}
func (c *BoxUsecase) GetFileInfo(fileId string) (res *string, httpCode int, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s", c.conf.Box.ApiUrl, fileId)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "GetFileInfo"), time.Now(), 1)
	//return lib.RequestGet(api, nil, map[string]string{"Authorization": "Bearer " + token})

	return lib.Request("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
}

// GetFileVersions 获取文件版本
// {"type":"file","id":"1554567330874","file_version":{"type":"file_version","id":"1708008340474","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68"},"sequence_id":"1","etag":"1","sha1":"58bb20dfae95e8640fd6606a36c575796e45ab68","name":"b.pdf","description":"","size":182138,"path_collection":{"total_count":7,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"247457173873","sequence_id":"1","etag":"1","name":"Testing"},{"type":"folder","id":"255166311971","sequence_id":"1","etag":"1","name":"Test Clients"},{"type":"folder","id":"264686374897","sequence_id":"2","etag":"2","name":"[PROD]VBC - TestLiao, TestGary #5076"},{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"},{"type":"folder","id":"268258622342","sequence_id":"0","etag":"0","name":"leve1_folder"}]},"created_at":"2024-06-08T03:41:17-07:00","modified_at":"2024-06-08T03:41:17-07:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-05-18T19:35:02-07:00","content_modified_at":"2024-05-18T19:35:02-07:00","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"268258622342","sequence_id":"0","etag":"0","name":"leve1_folder"},"item_status":"active"}
func (c *BoxUsecase) GetFileVersions(fileId string) (res *string, httpCode int, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}

	api := fmt.Sprintf("%s/2.0/files/%s/versions", c.conf.Box.ApiUrl, fileId)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "GetFileVersions"), time.Now(), 1)
	//return lib.RequestGet(api, nil, map[string]string{"Authorization": "Bearer " + token})

	response, err := lib.RequestDo("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil, 0, err
	}
	if response == nil {
		return nil, 0, errors.New("response is nil")
	}
	body, _ := io.ReadAll(response.Body)
	return to.Ptr(string(body)), response.StatusCode, nil
}

func (c *BoxUsecase) GetFolderInfoForTypeMap(folderId string) (res lib.TypeMap, httpCode int, err error) {
	data, httpCode, err := c.GetFolderInfo(folderId)
	if err != nil {
		return nil, httpCode, err
	}
	if data == nil {
		return nil, httpCode, errors.New("data is nil")
	}
	res = lib.ToTypeMapByString(*data)
	return res, httpCode, nil
}

// GetFolderInfo 获取文件夹的信息
// {"type":"folder","id":"268906213262","sequence_id":"2","etag":"2","name":"FolderTest","created_at":"2024-06-08T03:33:44-07:00","modified_at":"2024-06-08T03:45:48-07:00","description":"","size":162709,"path_collection":{"total_count":6,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"247457173873","sequence_id":"1","etag":"1","name":"Testing"},{"type":"folder","id":"255166311971","sequence_id":"1","etag":"1","name":"Test Clients"},{"type":"folder","id":"264686374897","sequence_id":"2","etag":"2","name":"[PROD]VBC - TestLiao, TestGary #5076"},{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"}]},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"modified_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"trashed_at":null,"purged_at":null,"content_created_at":"2024-06-08T03:33:44-07:00","content_modified_at":"2024-06-08T03:45:48-07:00","owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"folder_upload_email":null,"parent":{"type":"folder","id":"264686394097","sequence_id":"0","etag":"0","name":"VA Medical Records"},"item_status":"active","item_collection":{"total_count":1,"entries":[{"type":"folder","id":"268905519956","sequence_id":"0","etag":"0","name":"FolderTestSub"}],"offset":0,"limit":100,"order":[{"by":"type","direction":"ASC"},{"by":"name","direction":"ASC"}]}}
func (c *BoxUsecase) GetFolderInfo(folderId string) (res *string, httpCode int, err error) {

	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}

	api := fmt.Sprintf("%s/2.0/folders/%s", c.conf.Box.ApiUrl, folderId)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "GetFolderInfo"), time.Now(), 1)

	response, err := lib.RequestDo("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil, 0, err
	}
	if response == nil {
		return nil, 0, errors.New("response is nil")
	}
	body, _ := io.ReadAll(response.Body)
	return to.Ptr(string(body)), response.StatusCode, nil
}

func (c *BoxUsecase) DownloadToLocal(boxFileId string, suffix string) (fileName string, path string, err error) {
	path = configs.GetTempDir()
	fileName = fmt.Sprintf("%s.%s", uuid.UuidWithoutStrike(), suffix)
	filenamePath := path + "/" + fileName
	c.log.Info("filenamePath: ", filenamePath)
	reader, err := c.DownloadFile(boxFileId, "")
	if err != nil {
		return "", "", err
	}
	defer reader.Close()

	file, err := os.Create(filenamePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", "", err
	}
	return fileName, path, nil
}

func (c *BoxUsecase) Users() (res *string, httpCode int, err error) {
	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}
	api := fmt.Sprintf("%s/2.0/users", c.conf.Box.ApiUrl)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "Users"), time.Now(), 1)

	response, err := lib.RequestDo("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil, 0, err
	}
	if response == nil {
		return nil, 0, errors.New("response is nil")
	}
	body, _ := io.ReadAll(response.Body)
	return to.Ptr(string(body)), response.StatusCode, nil
}

// FolderCollaborations {"total_count":4,"entries":[{"type":"collaboration","id":"63772909822","created_by":null,"created_at":"2025-06-10T18:55:44-07:00","modified_at":"2025-06-10T18:55:44-07:00","expires_at":null,"status":"accepted","accessible_by":{"type":"group","id":"23718777558","name":"VBC Partners","group_type":"managed_group"},"invite_email":null,"role":"co-owner","acknowledged_at":"2025-06-10T18:55:44-07:00","item":null,"is_access_only":false,"app_item":null},{"type":"collaboration","id":"63772942826","created_by":null,"created_at":"2025-06-10T19:00:01-07:00","modified_at":"2025-06-10T19:00:10-07:00","expires_at":null,"status":"accepted","accessible_by":{"type":"user","id":"30690469672","name":"Edward Bunting","login":"ebunting@vetbenefitscenter.com","is_active":true},"invite_email":null,"role":"co-owner","acknowledged_at":"2025-06-10T19:00:01-07:00","item":null,"is_access_only":false,"app_item":null},{"type":"collaboration","id":"63773147751","created_by":null,"created_at":"2025-06-10T18:39:24-07:00","modified_at":"2025-06-10T18:39:24-07:00","expires_at":null,"status":"accepted","accessible_by":{"type":"group","id":"23718553499","name":"Executive Team","group_type":"managed_group"},"invite_email":null,"role":"editor","acknowledged_at":"2025-06-10T18:39:24-07:00","item":{"type":"folder","id":"324555049746","sequence_id":"1","etag":"1","name":"VBC Active Cases"},"is_access_only":false,"app_item":null},{"type":"collaboration","id":"64086181594","created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"info@vetbenefitscenter.com"},"created_at":"2025-06-23T02:17:55-07:00","modified_at":"2025-06-23T02:17:55-07:00","expires_at":null,"status":"accepted","accessible_by":{"type":"user","id":"41426608287","name":"Benigno Decena","login":"bdecena@vetbenefitscenter.com","is_active":true},"invite_email":null,"role":"editor","acknowledged_at":"2025-06-23T02:17:55-07:00","item":{"type":"folder","id":"327431320500","sequence_id":"0","etag":"0","name":"Test"},"is_access_only":false,"app_item":null}]}
func (c *BoxUsecase) FolderCollaborations(folderId string) (res *string, httpCode int, err error) {
	token, err := c.Token()
	if err != nil {
		return nil, 0, err
	}
	api := fmt.Sprintf("%s/2.0/folders/%s/collaborations", c.conf.Box.ApiUrl, folderId)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "FolderCollaborations"), time.Now(), 1)

	response, err := lib.RequestDo("GET", api, nil, map[string]string{"Authorization": "Bearer " + token})
	if err != nil {
		return nil, 0, err
	}
	if response == nil {
		return nil, 0, errors.New("response is nil")
	}
	body, _ := io.ReadAll(response.Body)
	return to.Ptr(string(body)), response.StatusCode, nil
}

func (c *BoxUsecase) DeleteCollaborations(collaborationId string) (*string, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/2.0/collaborations/%s", c.conf.Box.ApiUrl, collaborationId)

	c.UsageStatsUsecase.Stat(UsageTypeValue(UsageType_PREFIX_BOX, "DeleteCollaborations"), time.Now(), 1)
	res, _, err := lib.HTTPJsonWithHeaders("DELETE", api, nil, map[string]string{"Authorization": "Bearer " + token})
	return res, err
}
