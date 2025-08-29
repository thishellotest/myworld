package tests

import (
	"net/url"
	"testing"
	"time"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

func Test_ZohoUsecase_UsersCount(t *testing.T) {
	err := UT.ZohoUsecase.UsersCount()
	lib.DPrintln(err)
}

func Test_ZohoUsecase_SettingsModules(t *testing.T) {
	err := UT.ZohoUsecase.SettingsModules()
	lib.DPrintln(err)
}

func Test_ZohoUsecase_SettingsLayouts(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts("Leads")
	lib.DPrintln(err)
}

func Test_ZohoUsecase_SettingsLayouts_Notes(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts(config_zoho.Notes)
	lib.DPrintln(err)
}

func Test_ZohoUsecase_SettingsLayouts_Tasks(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts(config_zoho.Tasks)
	lib.DPrintln(err)
}

// Contacts:Clients
func Test_ZohoUsecase_Contacts(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts(config_zoho.Contacts)
	lib.DPrintln(err)
}

func Test_ZohoUsecase_Deals(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts("Deals")
	lib.DPrintln(err)
}

func Test_ZohoUsecase_Tasks(t *testing.T) {
	err := UT.ZohoUsecase.SettingsLayouts("Tasks")
	lib.DPrintln(err)
}

func Test_ZohoUsecase_GetRecords_Notes(t *testing.T) {
	fields := config_zoho.NotesLayout().NoteApiNames()
	params := make(url.Values)
	//params.Add("page", "20")
	//params.Add("per_page", "100")
	params.Add("page_token", "7dbf61cdffa375ef6698d83aa54439f629e25ed75a07f7dce03305ad5e17274d4c9796b2c321be5ab96bd8f00e928fe02a225d034e35b340bef58dffec0a55e8a579b60f73358fc01b80631197cc97f75a07c51c24883b7a9e27b6e1164a7536da474b3b3bb9bed9e7bda4356c2ba26e79ef568f9545bbc2c1c6d97d806cf4ac45ce8eb606bf89312d0cefed7d342a0c139510e635f2e20f316e73115030aeff")
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Notes, fields, params)
	lib.DPrintln(err)
	data := r.GetTypeList("data")
	lib.DPrintln(r)
	lib.DPrintln("sss:", len(data))
	//for _, v := range data {
	//	lib.DPrintln(v.GetString("Note_Content"))
	//	lib.DPrintln(v.GetString("Modified_Time"))
	//	//break
	//}
}

func Test_ZohoUsecase_GetRecords_Leads(t *testing.T) {
	fields := config_zoho.LeadsLayout().LeadApiNames()
	params := make(url.Values)
	//params.Add("converted", "both")
	/*
		To get the list of converted records. Default value is false.
		Possible values: true - get only converted records, false - get only non-converted records, both - get all records.
	*/
	// 需要在列表显示的leads，把这里设置为false
	params.Add("converted", "false")
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Leads, fields, params)
	lib.DPrintln(r, err)
}

func Test_ZohoUsecase_GetLead(t *testing.T) {
	r, err := UT.ZohoUsecase.GetLead("6159272000001347005")
	lib.DPrintln(r, err)
}

func Test_Leads_CreateRecord(t *testing.T) {
	record := make(lib.TypeMap)
	record.Set("First_Name", "TestF")
	record.Set("Last_Name", "TestL")
	record.Set("Email", "liaogling@gmail.com")
	record.Set("State", "California")
	record.Set("Mobile", "123-123-1234")
	record.Set("Description", "Objective: I already have a rating, but I think I could be underrated.\n\nCurrent:\n\nNew:\n\n")
	record.Set("Lead_Source", config_vbc.Source_Website)

	// 6159272000001347005
	leadGid, row, err := UT.ZohoUsecase.CreateRecord(config_zoho.Leads, record)
	lib.DPrintln("err:", err)
	lib.DPrintln(leadGid)
	lib.DPrintln(row)
}

func Test_Leads_PutRecord(t *testing.T) {
	record := make(lib.TypeMap)
	record.Set("First_Name", "TestF")
	record.Set("Last_Name", "TestL")
	record.Set("Email", "liaogling@gmail.com")
	record.Set("State", "California")
	record.Set("Mobile", "123-123-1234")
	record.Set("Description", "Objective: I already have a rating, but I think I could be underrated.\n\nService:\n\nCurrent:\n\nNew:\n\n")
	record.Set("Lead_Source", config_vbc.Source_Website)
	record.Set("id", "6159272000001347005")
	// 6159272000001347005
	leadGid, row, err := UT.ZohoUsecase.PutRecordV1(config_zoho.Leads, record)
	lib.DPrintln("err:", err)
	lib.DPrintln(leadGid)
	lib.DPrintln(row)
}

func Test_ZohoUsecase_GetRecords_Deals(t *testing.T) {
	//fields := config_zoho.DealLayout().FieldApiNamesByApiName(config_zoho.Deal_Sections_ApiName_ClientCaseInformation)
	fields := config_zoho.DealLayout().DealApiNames()
	lib.DPrintln("fields::", fields)
	params := make(url.Values)
	params.Add("page", "1")
	params.Add("per_page", "2")
	//params.Add("ids", "6159272000006133775")
	//return
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	lib.DPrintln(r, err)

}

func Test_ZohoUsecase_GetRecords_Tasks(t *testing.T) {
	//fields := config_zoho.DealLayout().FieldApiNamesByApiName(config_zoho.Deal_Sections_ApiName_ClientCaseInformation)
	fields := config_zoho.TaskLayout().TaskApiNames()
	lib.DPrintln(fields)
	//return
	params := make(url.Values)
	params.Add("page", "1")
	params.Add("per_page", "2")
	params.Add("ids", "6159272000007377003")
	//params.Add("ids", "6159272000003674030")

	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Tasks, fields, params)
	lib.DPrintln(r, err)

}

func Test_ZohoUsecase_PutRecord_Tasks(t *testing.T) {

	record := make(lib.TypeMap)
	record.Set("id", "6159272000001873015")
	record.Set("Status", config_zoho.ClientTaskStatus_Completed)

	gid, r, err := UT.ZohoUsecase.PutRecordV1(config_zoho.Tasks, record)
	lib.DPrintln(gid, r, err)
}

func Test_ZohoUsecase_PostRecordV1_Tasks(t *testing.T) {

	record := make(lib.TypeMap)
	record.Set("Status", config_zoho.ClientTaskStatus_Completed)
	record.Set("Owner", "6159272000000453669")
	record.Set("Created_By", "6159272000000453669")
	record.Set("Modified_By", "6159272000000453669")
	record.Set("Who_Id", "6159272000005519042")
	record.Set("What_Id", "6159272000007376001")
	record.Set("Priority", "High")
	record.Set("Status", "Not Started")
	record.Set("Subject", "ITF 80 days1")
	record.Set("$se_module", "Deals")
	record.Set("Due_Date", time.Now().Format(time.DateOnly))
	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")
	//record.Set("", "")

	gid, r, err := UT.ZohoUsecase.PostRecordV1(config_zoho.Tasks, record)
	lib.DPrintln(gid, r, err)
}

func Test_ZohoUsecase_GetRecords_Deals1(t *testing.T) {
	fields := config_zoho.DealLayout().DealApiNames()
	//fields := []string{"Owner", "Collaborators", "Deal_Name", "Stage"}
	params := make(url.Values)
	params.Add("ids", "6159272000007553001")
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	lib.DPrintln(r, err)
}

func Test_ZohoUsecase_GetRecords_Deals2(t *testing.T) {
	fields := config_zoho.DealLayout().DealApiNames2()
	//fields := []string{"Owner", "Collaborators", "Deal_Name", "Stage"}
	//params := make(url.Values)
	//params.Add("ids", "6159272000007553001")
	params := make(url.Values)
	params.Add("page", "1")
	params.Add("per_page", "100")
	params.Add("ids", "6159272000013512017")
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	lib.DPrintln(r, err)
}

func Test_ZohoUsecase_GetRelatedRecords(t *testing.T) {
	records, err := UT.ZohoUsecase.GetRelatedRecords(config_zoho.Deals, "6159272000001184003", []string{"Collaborators"})
	lib.DPrintln(records)
	lib.DPrintln(err)
}

func Test_ZohoUsecase_GetShareRecords(t *testing.T) {
	records, err := UT.ZohoUsecase.GetShareRecords(config_zoho.Deals, "6159272000001184003")
	lib.DPrintln(records)
	lib.DPrintln(err)
}

func Test_ZohoUsecase_GetDeal(t *testing.T) {
	e, rr := UT.ZohoUsecase.GetDeal("6159272000008327001")
	lib.DPrintln(e, rr)
}

func Test_ZohoUsecase_GetContact(t *testing.T) {
	e, rr := UT.ZohoUsecase.GetContact("6159272000000708001")
	lib.DPrintln(e, rr)
}

func Test_ZohoUsecase_GetRecords_Contacts(t *testing.T) {
	fields := config_zoho.ContactLayout().ContactApiNames()
	params := make(url.Values)
	//params.Add("ids", "6159272000000708051") // 已经删除的数据，如果数据都删除了，就会返回205状态码
	//params.Add("ids", "6159272000000708001")
	r, err := UT.ZohoUsecase.GetRecords(config_zoho.Contacts, fields, params)
	lib.DPrintln(r, err)
}

func Test_ZohoUsecase_Watch(t *testing.T) {
	input := `{
    "watch": [
        {
            "channel_id": "10000",
            "events": [
                "Deals.all"
            ],
            "notification_condition": [
                {
                    "type": "field_selection",
                    "module": {
                        "api_name": "Deals",
                        "id": "554023000000000131"
                    },
                    "field_selection": {
                        "group_operator": "or",
                        "group": [
                            {
                                "field": {
                                    "api_name": "Stage",
                                    "id": "554023000000000525"
                                }
                            },
                            {
                                "group_operator": "and",
                                "group": [
                                    {
                                        "field": {
                                            "api_name": "Account_Name",
                                            "id": "554023000000000523"
                                        }
                                    },
                                    {
                                        "field": {
                                            "api_name": "Lead_Source",
                                            "id": "554023000000000535"
                                        }
                                    }
                                ]
                            }
                        ]
                    }
                }
            ],
            "token": "deals.all.notif",
            "return_affected_field_values": true,
            "notify_url": "https:///webhook/post_source?from=zoho"
        }
    ]
}`
	err := UT.ZohoUsecase.Watch(input)
	lib.DPrintln(err)

	// 返回结果：
	//{"watch":[{"code":"SUCCESS","details":{"events":[{"channel_expiry":"2024-03-10T22:16:03+08:00","resource_uri":"https://www.zohoapis.com/crm/v2/Deals","resource_id":"6159272000000002181","resource_name":"Deals","channel_id":"10000"}]},"message":"Successfully subscribed for actions-watch of the given module","status":"success"}]}

}

func Test_ZohoUsecase_Users(t *testing.T) {
	r, err := UT.ZohoUsecase.Users()
	lib.DPrintln(r, err)

}

func Test_ZohoUsecase_sync_Users(t *testing.T) {

	err := UT.ZohobuzUsecase.BizSyncUsers()
	lib.DPrintln(err)

	//r, err := UT.ZohoUsecase.Users()
	//users := r.GetTypeList("users")
	//UT.ZohobuzUsecase.SyncUsers(users)
	//lib.DPrintln(r, err)
}

func Test_ZohoUsecase_PutRecord(t *testing.T) {

	record := make(lib.TypeMap)
	record.Set("id", "6159272000002701086")
	record.Set("Private_Exams_Needed", "TBI (str)\n\nLow Back Pain with bilateral radiculopathy (str)")
	record.Set("Pricing_Version", "v2024")

	//body := make(lib.TypeMap)
	//var data lib.TypeList
	//data = append(data, record)
	//body.Set("data", data)

	gid, r, err := UT.ZohoUsecase.PutRecordV1(config_zoho.Deals, record)
	lib.DPrintln(gid, r, err)
}

func Test_ZohoUsecase_PutRecord_Amount(t *testing.T) {

	record := make(lib.TypeMap)
	record.Set("id", "6159272000005255012")
	//record.Set("Amount", 1.24)
	record.Set("Contact_Form", "Yes")

	//body := make(lib.TypeMap)
	//var data lib.TypeList
	//data = append(data, record)
	//body.Set("data", data)

	gid, r, err := UT.ZohoUsecase.PutRecordV1(config_zoho.Deals, record)
	lib.DPrintln(gid, r, err)
}

func Test_ZohoUsecase_PutRecord_Contact(t *testing.T) {

	record := make(lib.TypeMap)
	record.Set("id", "6159272000000708001")
	record.Set("First_Name", "Ling1")
	record.Set("Last_Name", "Li1")

	//body := make(lib.TypeMap)
	//var data lib.TypeList
	//data = append(data, record)
	//body.Set("data", data)

	gid, r, err := UT.ZohoUsecase.PutRecordV1(config_zoho.Contacts, record)
	lib.DPrintln(gid, r, err)
}

func Test_ZohoUsecase_GetASpecificRecord(t *testing.T) {

	fields := config_zoho.DealLayout().DealApiNames()
	fields = []string{"Collaborators"}
	r, err := UT.ZohoUsecase.GetASpecificRecord(config_zoho.Deals, fields, "6159272000012996001")
	lib.DPrintln(r, err)

	//fields := []string{"Collaborators"}
	//params := make(url.Values)
	//params.Add("ids", "6159272000012996001")
	//r, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	//lib.DPrintln(r, err)
}

func Test_ZohoUsecase_GetASpecificRecord_Contacts(t *testing.T) {

	fields := config_zoho.DealLayout().DealApiNames()
	fields = []string{"Collaborators"}
	r, err := UT.ZohoUsecase.GetASpecificRecord(config_zoho.Contacts, fields, "6159272000012737026")
	lib.DPrintln(r, err)

	//fields := []string{"Collaborators"}
	//params := make(url.Values)
	//params.Add("ids", "6159272000012996001")
	//r, err := UT.ZohoUsecase.GetRecords(config_zoho.Deals, fields, params)
	//lib.DPrintln(r, err)
}
