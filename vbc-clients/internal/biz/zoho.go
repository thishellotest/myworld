package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_zoho"
	"vbc/lib"
)

type ZohoUsecase struct {
	log                *log.Helper
	CommonUsecase      *CommonUsecase
	conf               *conf.Data
	Oauth2TokenUsecase *Oauth2TokenUsecase
	UsageStatsUsecase  *UsageStatsUsecase
}

func NewZohoUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	Oauth2TokenUsecase *Oauth2TokenUsecase,
	UsageStatsUsecase *UsageStatsUsecase) *ZohoUsecase {
	uc := &ZohoUsecase{
		log:                log.NewHelper(logger),
		CommonUsecase:      CommonUsecase,
		conf:               conf,
		Oauth2TokenUsecase: Oauth2TokenUsecase,
		UsageStatsUsecase:  UsageStatsUsecase,
	}

	return uc
}

func (c *ZohoUsecase) Headers() (map[string]string, error) {

	token, err := c.Oauth2TokenUsecase.GetAccessToken(Oauth2_AppId_zohocrm)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Authorization": "Zoho-oauthtoken " + token}, nil
}

func (c *ZohoUsecase) UsersCount() error {

	c.UsageStatsUsecase.Stat("ZohoUsecase_UsersCount", time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/crm/v6/users/actions/count", c.conf.Zoho.ApiUrl)
	a, _, err := lib.Request("GET", api, nil, headers)
	lib.DPrintln(err)
	lib.DPrintln(a)

	return nil
}

func (c *ZohoUsecase) SettingsModules() error {

	c.UsageStatsUsecase.Stat("ZohoUsecase_SettingsModules", time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/crm/v6/settings/modules", c.conf.Zoho.ApiUrl)
	a, _, err := lib.Request("GET", api, nil, headers)
	lib.DPrintln(err)
	lib.DPrintln(a)

	return nil
}

func (c *ZohoUsecase) SettingsLayouts(moduleApiName string) error {

	c.UsageStatsUsecase.Stat("ZohoUsecase_SettingsLayouts_"+moduleApiName, time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/crm/v6/settings/layouts?module=%s", c.conf.Zoho.ApiUrl, moduleApiName)
	a, _, err := lib.Request("GET", api, nil, headers)
	lib.DPrintln(err)
	lib.DPrintln(a)

	return nil
}

// GetDeal 返回：{"Active_Duty":"No","Agent_Orange_Exposure":null,"Amount":13000,"Amyotrophic_Lateral_Sclerosis_ALS":null,"Atomic_Veterans_and_Radiation_Exposure":null,"Branch":"Marine Corps","Burn_Pits_and_Other_Airborne_Hazards":null,"C_File_Submitted":null,"Case_Files_Folder":"https://veteranbenefitscenter.app.box.com/folder/257067522190","City":"city 0413_9","Contact_Form":null,"Contact_Name":{"id":"6159272000001008204","name":"TestFn LnN"},"Created_By":{"email":"glliao@vetbenefitscenter.com","id":"6159272000000453669","name":"Engineering Team"},"Created_Time":"2024-04-09T22:05:21+08:00","Current_Rating":0,"DD214":null,"Date_of_Birth":"1990-12-08","Deal_Name":"TestFn LnN-0","Description":"Service:\n\n\nCurrent:\n\n\nNew:","Disability_Rating_List_Screenshot":null,"Effective_Current_Rating":0,"Email":"lialing@foxmail.com","Gulf_War_Illness":null,"ITF_Expiration":null,"Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun":null,"Last_Activity_Time":"2024-04-13T18:21:04+08:00","Lead_Conversion_Time":null,"Lead_Source":"Partner","Modified_By":{"email":"glliao@vetbenefitscenter.com","id":"6159272000000453669","name":"Engineering Team"},"Modified_Time":"2024-04-13T18:21:04+08:00","New_Rating":null,"Overall_Sales_Duration":null,"Owner":{"email":"glliao@vetbenefitscenter.com","id":"6159272000000453669","name":"Engineering Team"},"Phone":"041-041-0418","Pipeline":"VBC Clients","Rating_Decision_Letters":null,"Referring_Person":null,"Retired":"No","SSN":"041-04-0419","STRs":null,"Sales_Cycle_Duration":null,"Stage":"1. Fee Schedule and Contract","State":"Delaware","Street_Address":"addr 0413_9","TDIU":null,"Tag":[],"Zip_Code":"20009","id":"6159272000001184003"}
func (c *ZohoUsecase) GetDeal(gid string) (lib.TypeMap, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_GetDeal", time.Now(), 1)

	fields := config_zoho.DealLayout().DealApiNames()
	params := make(url.Values)
	params.Add("ids", gid)
	r, err := c.GetRecords(config_zoho.Deals, fields, params)
	if err != nil {
		return nil, err
	}
	if r != nil {
		data := r.GetTypeList("data")
		if len(data) > 0 {
			return data[0], nil
		}
	}
	return nil, nil
}

func (c *ZohoUsecase) GetContact(gid string) (lib.TypeMap, error) {
	c.UsageStatsUsecase.Stat("ZohoUsecase_GetContact", time.Now(), 1)
	fields := config_zoho.ContactLayout().ContactApiNames()
	params := make(url.Values)
	params.Add("ids", gid)
	r, err := c.GetRecords(config_zoho.Contacts, fields, params)
	if err != nil {
		return nil, err
	}
	if r != nil {
		data := r.GetTypeList("data")
		if len(data) > 0 {
			return data[0], nil
		}
	}
	return nil, nil
}

func (c *ZohoUsecase) GetLead(gid string) (lib.TypeMap, error) {
	c.UsageStatsUsecase.Stat("ZohoUsecase_GetLead", time.Now(), 1)
	fields := config_zoho.LeadsLayout().LeadApiNames()
	params := make(url.Values)
	params.Add("ids", gid)
	r, err := c.GetRecords(config_zoho.Leads, fields, params)
	if err != nil {
		return nil, err
	}
	if r != nil {
		data := r.GetTypeList("data")
		if len(data) > 0 {
			return data[0], nil
		}
	}
	return nil, nil
}

// GetASpecificRecord 此方法，可以返回mutli user
// https://help.zoho.com/portal/en/community/topic/kaizen-87-multi-user-lookup-in-api-v4
// You can use Get Records API to fetch details of your multi-user lookup field. Note that it is returned in the response only when you fetch a specific record. Let us see how to fetch details of a record.
func (c *ZohoUsecase) GetASpecificRecord(moduleApiName string, fieldApiNames []string, gid string) (lib.TypeMap, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_GetASpecificRecord_"+moduleApiName, time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/crm/v6/%s/%s", c.conf.Zoho.ApiUrl, moduleApiName, gid)
	query := make(url.Values)

	fieldsString := strings.Join(fieldApiNames, ",")
	query.Add("fields", fieldsString)

	records, _, err := lib.RequestGet(api, query, headers)
	if err != nil {
		return nil, lib.ErrorWrap(err, records)
	}
	if records == nil {
		return nil, errors.New("records is nil")
	}
	return lib.ToTypeMapByString(*records), nil
}
func (c *ZohoUsecase) GetRecords(moduleApiName string, fieldApiNames []string, queryParams url.Values) (lib.TypeMap, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_GetRecords_"+moduleApiName, time.Now(), 1)

	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/crm/v6/%s", c.conf.Zoho.ApiUrl, moduleApiName)
	query := make(url.Values)

	fieldsString := strings.Join(fieldApiNames, ",")
	query.Add("fields", fieldsString)
	// 必须指定字段， 否则只返回id
	//query.Add("fields", "Last_Name,Email,Record_Status__s,Converted__s,Converted_Date_Time")

	//query.Add("fields", "*")
	//if converted != "" {
	//	query.Add("converted", converted) // both / true /  false
	//}
	if queryParams != nil {
		for k, _ := range queryParams {
			for _, v1 := range queryParams[k] {
				query.Add(k, v1)
			}
		}
	}
	//query.Add("per_page", InterfaceToString(200))
	//query.Add("page", InterfaceToString(1))
	query.Add("sort_by", "Modified_Time")
	query.Add("sort_order", "desc")

	//query.Add("next_page_token", "c2e52874eaf459516778d639c4886e2e55d102bfdce9240618cf0a53e917f51257620f1e34d00536f2dce4b6224bdcf981f168f56017679f06e9fddb09259176a9c6d86dc5d8a38db5650306989563eeb5958661c9cc701265f7a570ab646be530eb62556123ddab18e87343218f9b590650a766623742652901e5b3153abdd5")
	records, _, err := lib.RequestGet(api, query, headers)
	if err != nil {
		return nil, lib.ErrorWrap(err, records)
	}
	if records == nil {
		return nil, errors.New("records is nil")
	}
	return lib.ToTypeMapByString(*records), nil
}

func (c *ZohoUsecase) GetShareRecords(moduleApiName string, recordId string) (lib.TypeMap, error) {

	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	/*

		curl "https://www.zohoapis.com/crm/v6/Contacts/3652397000000649013/actions/share"
		-X GET
		-H "Authorization: Zoho-oauthtoken 1000.8cb99dxxxxxxxxxxxxx9be93.9b8xxxxxxxxxxxxxxxf"

	*/
	api := fmt.Sprintf("%s/crm/v6/%s/%s/actions/share", c.conf.Zoho.ApiUrl, moduleApiName, recordId)

	records, _, err := lib.RequestGet(api, nil, headers)
	if err != nil {
		return nil, lib.ErrorWrap(err, records)
	}
	if records == nil {
		return nil, errors.New("GetShareRecords: records is nil")
	}
	return lib.ToTypeMapByString(*records), nil
}

func (c *ZohoUsecase) GetRelatedRecords(moduleApiName string, recordId string, relatedListApiNames []string) (lib.TypeMap, error) {

	//c.UsageStatsUsecase.Stat("ZohoUsecase_GetRecords_"+moduleApiName, time.Now(), 1)

	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/crm/v6/%s/%s/Users", c.conf.Zoho.ApiUrl, moduleApiName, recordId)
	query := make(url.Values)
	relatedListApiNameStr := strings.Join(relatedListApiNames, ",")
	query.Add("fields", relatedListApiNameStr)

	records, _, err := lib.RequestGet(api, query, headers)
	if err != nil {
		return nil, lib.ErrorWrap(err, records)
	}
	if records == nil {
		return nil, errors.New("records is nil")
	}
	return lib.ToTypeMapByString(*records), nil
}

func (c *ZohoUsecase) ChangeDealAmount(gid string, amount string) error {

	c.UsageStatsUsecase.Stat("ZohoUsecase_ChangeDealAmount", time.Now(), 1)

	record := make(lib.TypeMap)
	record.Set("id", gid)
	record.Set("Amount", amount)
	gid, r, err := c.PutRecordV1(config_zoho.Deals, record)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("ChangeDealAmount: r is nil")
	}
	if r.GetString("code") != "SUCCESS" {
		return errors.New("ChangeDealAmount: error code: " + r.GetString("code"))
	}
	return nil
}

func (c *ZohoUsecase) ChangeDealV1(gid string, params lib.TypeMap) error {
	c.UsageStatsUsecase.Stat("ZohoUsecase_ChangeDeal", time.Now(), 1)
	record := make(lib.TypeMap)
	record.Set("id", gid)

	for k, v := range params {
		record.Set(k, v)
	}
	gid, r, err := c.PutRecordV1(config_zoho.Deals, record)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("ChangeDeal: r is nil")
	}
	if r.GetString("code") != "SUCCESS" {
		return errors.New("ChangeDeal: error code: " + r.GetString("code"))
	}
	return nil
}

func (c *ZohoUsecase) PutRecordV1(moduleApiName string, record lib.TypeMap) (gid string, row lib.TypeMap, err error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_PutRecord_"+moduleApiName, time.Now(), 1)
	records := make(lib.TypeMap)
	var data lib.TypeList
	data = append(data, record)
	records.Set("data", data)

	response, err := c.PutRecordsV1(moduleApiName, records)
	return c.TidyRecordResponse(response, err)
}

func (c *ZohoUsecase) PostRecordV1(moduleApiName string, record lib.TypeMap) (gid string, row lib.TypeMap, err error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_PostRecordV1_"+moduleApiName, time.Now(), 1)
	records := make(lib.TypeMap)
	var data lib.TypeList
	data = append(data, record)
	records.Set("data", data)

	response, err := c.PostRecordsV1(moduleApiName, records)
	return c.TidyRecordResponse(response, err)
}

func (c *ZohoUsecase) TidyRecordResponse(response *string, err error) (string, lib.TypeMap, error) {

	if err != nil {
		return "", nil, lib.ErrorWrap(err, response)
	}
	if response == nil {
		return "", nil, errors.New("response is nil")
	}
	responseMap := lib.ToTypeMapByString(*response)
	list := responseMap.GetTypeList("data")
	if len(list) != 1 {
		return "", nil, errors.New("data is wrong.")
	}
	if list[0].GetString("code") != "SUCCESS" {
		return "", nil, errors.New("code error:" + list[0].GetString("code"))
	}
	gid := list[0].GetString("details.id")
	if gid == "" {
		return "", nil, errors.New("Response gid is wrong.")
	}
	return gid, list[0], nil
}

func (c *ZohoUsecase) CreateRecord(moduleApiName string, record lib.TypeMap) (gid string, row lib.TypeMap, err error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_CreateRecord_"+moduleApiName, time.Now(), 1)
	records := make(lib.TypeMap)
	var data lib.TypeList
	data = append(data, record)
	records.Set("data", data)

	response, err := c.CreateRecords(moduleApiName, records)
	return c.TidyRecordResponse(response, err)
}

func (c *ZohoUsecase) CreateRecords(moduleApiName string, records lib.TypeMap) (*string, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_CreateRecords_"+moduleApiName, time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/crm/v6/%s", c.conf.Zoho.ApiUrl, moduleApiName)
	res, _, err := lib.HTTPJsonWithHeaders(http.MethodPost, api, records.ToBytes(), headers)
	return res, err
}

func (c *ZohoUsecase) PutRecordsV1(moduleApiName string, records lib.TypeMap) (*string, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_PutRecords_"+moduleApiName, time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/crm/v6/%s", c.conf.Zoho.ApiUrl, moduleApiName)
	res, _, err := lib.HTTPJsonWithHeaders(http.MethodPut, api, records.ToBytes(), headers)
	return res, err
}

func (c *ZohoUsecase) PostRecordsV1(moduleApiName string, records lib.TypeMap) (*string, error) {

	c.UsageStatsUsecase.Stat("ZohoUsecase_PostRecordsV1_"+moduleApiName, time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}
	api := fmt.Sprintf("%s/crm/v6/%s", c.conf.Zoho.ApiUrl, moduleApiName)
	res, _, err := lib.HTTPJsonWithHeaders(http.MethodPost, api, records.ToBytes(), headers)
	return res, err
}

func (c *ZohoUsecase) Watch(input string) error {

	c.UsageStatsUsecase.Stat("ZohoUsecase_Watch", time.Now(), 1)
	headers, err := c.Headers()
	if err != nil {
		return err
	}

	api := fmt.Sprintf("%s/crm/v6/actions/watch", c.conf.Zoho.ApiUrl)
	a, _, err := lib.Request("POST", api, []byte(input), headers)
	lib.DPrintln(err)
	lib.DPrintln(a)

	return nil
}

func (c *ZohoUsecase) Users() (lib.TypeMap, error) {
	headers, err := c.Headers()
	if err != nil {
		return nil, err
	}

	api := fmt.Sprintf("%s/crm/v6/users?type=AllUsers", c.conf.Zoho.ApiUrl)
	records, _, err := lib.Request("GET", api, nil, headers)

	if err != nil {
		return nil, lib.ErrorWrap(err, records)
	}
	if records == nil {
		return nil, errors.New("records si nil")
	}
	return lib.ToTypeMapByString(*records), nil
}
