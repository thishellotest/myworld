package tests

import (
	"net/url"
	"testing"
	"vbc/internal/biz"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_DefaultStageTrans(t *testing.T) {
	stage := "16. Medical Team - Private Exams Submitted"

	aa := biz.DefaultStageTrans(stage)
	lib.DPrintln(aa)
}

func Test_ZohobuzUsecase_SyncClients(t *testing.T) {

	fields := config_zoho.ContactLayout().ContactApiNames()
	//fields := config_zoho.ContactLayout().FieldApiNamesByApiName(config_zoho.Contact_Sections_ApiName_Client_Information)
	params := make(url.Values)
	//params.Add("ids", "6159272000000516059")
	params.Add("ids", "6159272000000708001")
	params.Add("ids", "6159272000000708025")
	records, err := UT.ZohoUsecase.GetRecords(config_zoho.Contacts, fields, params)
	lib.DPrintln(err)

	listMaps := records.GetTypeList("data")
	lib.DPrintln(listMaps)
	err = UT.ZohobuzUsecase.SyncClients(listMaps)
	lib.DPrintln(err)

	//contactRecords := `{"data":[{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"TestClientAndCaseFN@a.com","Salutation":null,"Last_Activity_Time":"2024-03-16T10:39:25+08:00","First_Name":"TestClientAndCaseFN","Full_Name":"TestClientAndCaseFN TestClientAndCaseLN","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":"Air Force","id":"6159272000000516059","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":"Referring PersonVal","Modified_Time":"2024-03-10T18:36:14+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":"Yes","Created_Time":"2024-03-10T18:36:14+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":"138","SSN":null,"Effective_Current_Rating":10,"Last_Name":"TestClientAndCaseLN","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null},{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"JustClientFN1@a.com","Salutation":null,"Last_Activity_Time":"2024-03-10T18:28:35+08:00","First_Name":"JustClientFN1","Full_Name":"JustClientFN1 JustClientLN1","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":null,"id":"6159272000000516012","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":null,"Modified_Time":"2024-03-10T18:28:35+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":null,"Created_Time":"2024-03-10T18:28:35+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":null,"SSN":null,"Effective_Current_Rating":10,"Last_Name":"JustClientLN1","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null}],"info":{"per_page":5,"next_page_token":null,"count":2,"sort_by":"id","page":1,"previous_page_token":null,"page_token_expiry":null,"sort_order":"desc","more_records":false}}`
	//lib.DPrintln(contactRecords)
}

func Test_ZohobuzUsecase_SyncClientCases(t *testing.T) {

	//fields := config_zoho.DealLayout().FieldApiNamesByApiName(config_zoho.Deal_Sections_ApiName_ClientCaseInformation)
	fields := config_zoho.DealLayout().DealApiNames()
	params := make(url.Values)
	//params.Add("ids", "6159272000000516059")
	params.Add("ids", "6159272000009972111")
	//params.Add("ids", "6159272000000708006")
	records, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	lib.DPrintln(err)

	listMaps := records.GetTypeList("data")
	lib.DPrintln(listMaps)
	err = UT.ZohobuzUsecase.SyncClientCases(listMaps, "")
	lib.DPrintln(err)

	//contactRecords := `{"data":[{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"TestClientAndCaseFN@a.com","Salutation":null,"Last_Activity_Time":"2024-03-16T10:39:25+08:00","First_Name":"TestClientAndCaseFN","Full_Name":"TestClientAndCaseFN TestClientAndCaseLN","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":"Air Force","id":"6159272000000516059","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":"Referring PersonVal","Modified_Time":"2024-03-10T18:36:14+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":"Yes","Created_Time":"2024-03-10T18:36:14+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":"138","SSN":null,"Effective_Current_Rating":10,"Last_Name":"TestClientAndCaseLN","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null},{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"JustClientFN1@a.com","Salutation":null,"Last_Activity_Time":"2024-03-10T18:28:35+08:00","First_Name":"JustClientFN1","Full_Name":"JustClientFN1 JustClientLN1","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":null,"id":"6159272000000516012","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":null,"Modified_Time":"2024-03-10T18:28:35+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":null,"Created_Time":"2024-03-10T18:28:35+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":null,"SSN":null,"Effective_Current_Rating":10,"Last_Name":"JustClientLN1","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null}],"info":{"per_page":5,"next_page_token":null,"count":2,"sort_by":"id","page":1,"previous_page_token":null,"page_token_expiry":null,"sort_order":"desc","more_records":false}}`
	//lib.DPrintln(contactRecords)
}

func Test_ZohobuzUsecase_SyncTasks(t *testing.T) {

	//fields := config_zoho.DealLayout().FieldApiNamesByApiName(config_zoho.Deal_Sections_ApiName_ClientCaseInformation)
	fields := config_zoho.TaskLayout().TaskApiNames()
	params := make(url.Values)
	//params.Add("ids", "6159272000000516059")
	//params.Add("ids", "6159272000001873015")
	params.Add("ids", "6159272000001600615")
	records, err := UT.ZohoUsecase.GetRecords(config_zoho.Tasks, fields, params)
	lib.DPrintln(err)
	lib.DPrintln(records)
	//return
	listMaps := records.GetTypeList("data")
	lib.DPrintln(listMaps)
	err = UT.ZohobuzUsecase.SyncTasks(listMaps)
	lib.DPrintln(err)

	//contactRecords := `{"data":[{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"TestClientAndCaseFN@a.com","Salutation":null,"Last_Activity_Time":"2024-03-16T10:39:25+08:00","First_Name":"TestClientAndCaseFN","Full_Name":"TestClientAndCaseFN TestClientAndCaseLN","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":"Air Force","id":"6159272000000516059","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":"Referring PersonVal","Modified_Time":"2024-03-10T18:36:14+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":"Yes","Created_Time":"2024-03-10T18:36:14+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":"138","SSN":null,"Effective_Current_Rating":10,"Last_Name":"TestClientAndCaseLN","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null},{"Owner":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Email":"JustClientFN1@a.com","Salutation":null,"Last_Activity_Time":"2024-03-10T18:28:35+08:00","First_Name":"JustClientFN1","Full_Name":"JustClientFN1 JustClientLN1","Modified_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Unsubscribed_Mode":null,"Current_Rating":10,"Branch":null,"id":"6159272000000516012","Email_Opt_Out":false,"Data_Source":"Converted","Referring_Person":null,"Modified_Time":"2024-03-10T18:28:35+08:00","Date_of_Birth":null,"Enrich_Status__s":null,"Retired":null,"Created_Time":"2024-03-10T18:28:35+08:00","Unsubscribed_Time":null,"Change_Log_Time__s":null,"Mobile":null,"SSN":null,"Effective_Current_Rating":10,"Last_Name":"JustClientLN1","Locked__s":false,"Lead_Source":"Manual","Created_By":{"name":"Gary Liao","id":"6159272000000453669","email":"glliao@vetbenefitscenter.com"},"Tag":[],"Last_Enriched_Time__s":null}],"info":{"per_page":5,"next_page_token":null,"count":2,"sort_by":"id","page":1,"previous_page_token":null,"page_token_expiry":null,"sort_order":"desc","more_records":false}}`
	//lib.DPrintln(contactRecords)
}

func Test_ZohobuzUsecase_SyncClientOne(t *testing.T) {

	contactRowMap := lib.ToTypeMapByString(`{
	"Branch": "National Oceanic and Atmospheric Administration",
	"Change_Log_Time__s": null,
	"Created_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Created_Time": "2024-03-23T18:58:48+08:00",
	"Current_Rating": 100,
	"Data_Source": "Manual",
	"Date_of_Birth": "1990-11-11",
	"Effective_Current_Rating": 20,
	"Email": "lialing@foxmail.com",
	"Email_Opt_Out": false,
	"Enrich_Status__s": null,
	"First_Name": "Shi",
	"Full_Name": "Shi Li",
	"Last_Activity_Time": "2024-03-30T18:13:35+08:00",
	"Last_Enriched_Time__s": null,
	"Last_Name": "Li",
	"Lead_Source": null,
	"Mailing_City": "city 01",
	"Mailing_Country": "USA",
	"Mailing_State": "Georgia",
	"Mailing_Street": "address new val01",
	"Mailing_Zip": "60000",
	"Mobile": "600-000-0000",
	"Modified_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Modified_Time": "2024-03-30T18:13:28+08:00",
	"Owner": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Referring_Person": "Referring Person val11",
	"Retired": "Yes",
	"SSN": "600-00-0000",
	"Salutation": null,
	"Tag": [],
	"Unsubscribed_Mode": null,
	"Unsubscribed_Time": null,
	"Vendor_Name": null,
	"id": "6159272000000708001"
}`)
	err := UT.ZohobuzUsecase.SyncClientOne(contactRowMap)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SyncClientCaseOne(t *testing.T) {
	// "Amount": null,
	// "Amount": 100.11,
	rowMap := lib.ToTypeMapByString(`{
	"Agent_Orange_Exposure": "Yes",
	"Amount": 100.11,
	"Amyotrophic_Lateral_Sclerosis_ALS": "Yes",
	"Approval_Date": "2024-03-29",
	"Atomic_Veterans_and_Radiation_Exposure": "Yes",
	"Branch": "National Oceanic and Atmospheric Administration",
	"Burn_Pits_and_Other_Airborne_Hazards": "Yes",
	"C_File_Submitted": "N/A",
	"Change_Log_Time__s": null,
	"Contact_Form": "Yes",
	"Contact_Name": {
		"id": "6159272000000708001",
		"name": "Shi Li"
	},
	"Created_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Created_Time": "2024-03-29T23:06:04+08:00",
	"Current_Rating": 80,
	"DD214": "Yes",
	"Deal_Name": "Shi Li@Case1",
	"Description": "Description val",
	"Disability_Rating_List_Screenshot": "No",
	"Effective_Current_Rating": 90,
	"Filing_Date": "2024-03-28",
	"Gulf_War_Illness": "Yes",
	"ITF_Expiration": "2024-03-30",
	"Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun": "Yes",
	"Last_Activity_Time": "2024-03-30T16:41:59+08:00",
	"Lead_Conversion_Time": null,
	"Lead_Source": null,
	"Modified_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Modified_Time": "2024-03-30T16:41:59+08:00",
	"New_Rating": null,
	"Next_Step": "a",
	"Overall_Sales_Duration": null,
	"Owner": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Pipeline": "VBC Clients",
	"Rating_Decision_Letters": "Yes",
	"Referring_Person": "Referring Person val",
	"Retired": "Yes",
	"STRs": "Yes",
	"Sales_Cycle_Duration": null,
	"Stage": "Getting Started Email",
	"TDIU": "Yes",
	"Tag": [],
	"id": "6159272000000820046",
	"Description":"Service:\n\n\nCurrent:\n\n\nNew:"
}`)
	err := UT.ZohobuzUsecase.SyncClientCaseOne(rowMap)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SyncClientCaseOne_with_null_value(t *testing.T) {
	rowMap := lib.ToTypeMapByString(`{
	"Agent_Orange_Exposure": null,
	"Amount": null,
	"Amyotrophic_Lateral_Sclerosis_ALS": null,
	"Approval_Date": null,
	"Atomic_Veterans_and_Radiation_Exposure": null,
	"Branch": null,
	"Burn_Pits_and_Other_Airborne_Hazards": null,
	"C_File_Submitted": null,
	"Change_Log_Time__s": null,
	"Contact_Form": null,
	"Contact_Name": {
		"id": "6159272000000878001",
		"name": "TestForZohoOwner TestForZohoLn"
	},
	"Created_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Created_Time": "2024-04-01T22:54:51+08:00",
	"Current_Rating": 0,
	"DD214": null,
	"Deal_Name": "TestForZohoLn TestForZohoOwner",
	"Description": null,
	"Disability_Rating_List_Screenshot": null,
	"Effective_Current_Rating": 0,
	"Filing_Date": null,
	"Gulf_War_Illness": null,
	"ITF_Expiration": null,
	"Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun": null,
	"Last_Activity_Time": "2024-04-01T23:22:42+08:00",
	"Lead_Conversion_Time": null,
	"Lead_Source": "Manual",
	"Modified_By": {
		"email": "glliao@vetbenefitscenter.com",
		"id": "6159272000000453669",
		"name": "Gary Liao"
	},
	"Modified_Time": "2024-04-01T23:22:42+08:00",
	"New_Rating": null,
	"Overall_Sales_Duration": null,
	"Owner": {
		"email": "ebunting@vetbenefitscenter.com",
		"id": "6159272000000453640",
		"name": "Edward Bunting"
	},
	"Pipeline": "VBC Clients",
	"Rating_Decision_Letters": null,
	"Referring_Person": null,
	"Retired": "No",
	"STRs": null,
	"Sales_Cycle_Duration": null,
	"Stage": "Fee Schedule and Contract",
	"TDIU": null,
	"Tag": [],
	"id": "6159272000000881012"
}`)
	err := UT.ZohobuzUsecase.SyncClientCaseOne(rowMap)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SyncUserOne(t *testing.T) {
	rowMap := lib.ToTypeMapByString(`{
		"$current_shift": null,
		"$next_shift": null,
		"$shift_effective_from": null,
		"Isonline": true,
		"Modified_By": {
			"id": "6159272000000453001",
			"name": "Yannan Wang"
		},
		"Modified_Time": "2024-03-17T09:34:41+08:00",
		"Source__s": null,
		"alias": null,
		"category": "regular_user",
		"city": null,
		"confirm": true,
		"country": null,
		"country_locale": "en_US",
		"created_by": {
			"id": "6159272000000453001",
			"name": "Yannan Wang"
		},
		"created_time": "2024-03-04T00:30:07+08:00",
		"customize_info": {
			"bc_view": null,
			"notes_desc": null,
			"show_detail_view": true,
			"show_home": false,
			"show_right_panel": null,
			"unpin_recent_item": null
		},
		"date_format": "MMM d, yyyy",
		"decimal_separator": "Period",
		"default_tab_group": "0",
		"dob": null,
		"email": "glliao@vetbenefitscenter.com",
		"fax": null,
		"first_name": "Gary",
		"full_name": "Gary Liao",
		"id": "6159272000000453669",
		"language": "en_US",
		"last_name": "Liao",
		"locale": "en_US",
		"microsoft": false,
		"mobile": null,
		"name_format__s": "Salutation,First Name,Last Name",
		"number_separator": "Comma",
		"offset": 28800000,
		"personal_account": false,
		"phone": null,
		"profile": {
			"id": "6159272000000026011",
			"name": "Administrator"
		},
		"role": {
			"id": "6159272000000453682",
			"name": "Engineering Team Leader"
		},
		"sandboxDeveloper": false,
		"signature": null,
		"sort_order_preference__s": "First Name,Last Name",
		"state": null,
		"status": "active",
		"status_reason__s": null,
		"street": null,
		"theme": {
			"background": "#F3F0EB",
			"new_background": null,
			"normal_tab": {
				"background": "#222222",
				"font_color": "#FFFFFF"
			},
			"screen": "fixed",
			"selected_tab": {
				"background": "#222222",
				"font_color": "#FFFFFF"
			},
			"type": "default"
		},
		"time_format": "hh:mm a",
		"time_zone": "Asia/Hong_Kong",
		"website": null,
		"zip": null,
		"zuid": "847408674"
	}`)
	err := UT.ZohobuzUsecase.SyncUserOne(rowMap)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_BizSyncUsers(t *testing.T) {
	err := UT.ZohobuzUsecase.BizSyncUsers()
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_HandleZohoDelete(t *testing.T) {
	err := UT.ZohobuzUsecase.HandleZohoDelete([]*biz.ClientCaseEntity{
		{
			Gid: "6159272000000708060",
		},
		{
			Gid: "6159272000000708006",
		},
		{
			Gid: "6159272000000708005",
		},
	})
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SyncDealsDeletes(t *testing.T) {
	err := UT.ZohobuzUsecase.SyncDealsDeletes()
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SYNC_HandleAmount(t *testing.T) {

	aaa, err := UT.TUsecase.ListByCond(biz.Kind_client_cases, And(Eq{"id": "163", "biz_deleted_at": 0, "deleted_at": 0}))

	lib.DPrintln(err)
	for _, v := range aaa {
		caseId := v.CustomFields.NumberValueByNameBasic("id")
		er := UT.ZohobuzUsecase.HandleAmount(caseId)
		if er != nil {
			panic(er)
		}
	}

	//err := UT.ZohobuzUsecase.HandleAmount(59)
	//lib.DPrintln(err)
}

func Test_ZohobuzUsecase_HandleAllMan(t *testing.T) {
	err := UT.ZohobuzUsecase.HandleAllMan()
	lib.DPrintln(err)
}
func Test_ZohobuzUsecase_HandleAmount(t *testing.T) {
	err := UT.ZohobuzUsecase.HandleAmount(5885)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_HandleClientCaseName(t *testing.T) {
	err := UT.ZohobuzUsecase.HandleClientCaseName(5004)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_BizHttpHandleClientCaseName(t *testing.T) {
	clientCaseIds := "5004"
	err := UT.ZohobuzUsecase.BizHttpHandleClientCaseName(clientCaseIds)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_BizHttpHandleNotes(t *testing.T) {
	clientCaseIds := "5004"
	err := UT.ZohobuzUsecase.BizHttpHandleNotes(clientCaseIds)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_HandleSyncZohoPricingVersion(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5076)
	err := UT.ZohobuzUsecase.HandleSyncZohoPricingVersion(tCase)
	lib.DPrintln(err)
}

func Test_ZohobuzUsecase_SyncClientsDeletes(t *testing.T) {
	UT.ZohobuzUsecase.SyncClientsDeletes()
}

func Test_ZohobuzUsecase_HandleAmountForAm(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5728)
	err := UT.ZohobuzUsecase.HandleAmountForAm(*tCase)
	lib.DPrintln(err)
}
