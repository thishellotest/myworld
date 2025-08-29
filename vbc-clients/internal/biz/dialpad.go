package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	"vbc/lib/to"
)

type DialpadUsecase struct {
	log                 *log.Helper
	CommonUsecase       *CommonUsecase
	conf                *conf.Data
	LogUsecase          *LogUsecase
	UnsubscribesUsecase *UnsubscribesUsecase
}

func NewDialpadUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	LogUsecase *LogUsecase,
	UnsubscribesUsecase *UnsubscribesUsecase) *DialpadUsecase {
	uc := &DialpadUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		LogUsecase:          LogUsecase,
		UnsubscribesUsecase: UnsubscribesUsecase,
	}

	return uc
}

func (c *DialpadUsecase) AccessKey() string {
	return configs.EnvDialpadKey()
}

func (c *DialpadUsecase) ApiUrl() string {
	return "https://dialpad.com"
}

// UserList Get Dialpad users
func (c *DialpadUsecase) UserList() (lib.TypeList, error) {
	api := fmt.Sprintf("%s/api/v2/users?apikey=%s", c.ApiUrl(), c.AccessKey())
	res, _, err := lib.HTTPJson("GET", api, nil)
	if err != nil {
		return nil, err
	}
	//lib.DPrintln(code, err)
	//lib.DPrintln(*res)
	data := lib.ToTypeMapByString(*res)
	list := data.GetTypeList("items")
	return list, nil
}

func DialpadGetEmail(res lib.TypeMap) string {
	if res != nil {
		emails := res.GetTypeListInterface("emails")
		if len(emails) > 0 {
			return InterfaceToString(emails[0])
		}
	}
	return ""
}

func DialpadGetPhoneNumber(res lib.TypeMap) string {
	if res != nil {
		phoneNumbers := res.GetTypeListInterface("phone_numbers")
		if len(phoneNumbers) > 0 {
			return InterfaceToString(phoneNumbers[0])
		}
	}
	return ""
}

func DialpadGetId(res lib.TypeMap) string {
	if res != nil {
		return res.GetString("id")
	}
	return ""
}

func (c *DialpadUsecase) UserInfoByEmail(email string) (destEmail, phoneNumber, dialpadId string) {

	res, err := c.UserByEmail(email)
	if err != nil {
		c.log.Error(err)
	}
	return DialpadGetEmail(res), DialpadGetPhoneNumber(res), DialpadGetId(res)

}

// UserByEmail Get Dialpad user
func (c *DialpadUsecase) UserByEmail(email string) (lib.TypeMap, error) {
	list, err := c.UserListByEmail(email)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// UserListByEmail Get Dialpad user
func (c *DialpadUsecase) UserListByEmail(email string) (lib.TypeList, error) {
	api := fmt.Sprintf("%s/api/v2/users?apikey=%s&email=%s", c.ApiUrl(), c.AccessKey(), email)
	res, _, err := lib.HTTPJson("GET", api, nil)
	if err != nil {
		return nil, err
	}
	//lib.DPrintln(code, err)
	//lib.DPrintln(*res)
	data := lib.ToTypeMapByString(*res)
	list := data.GetTypeList("items")
	return list, nil
}

// 错误的格式：{"error":{"code":400,"errors":[{"domain":"global","message":"Must provide text field, media field, or both","reason":"badRequest"}],"message":"Must provide text field, media field, or both"}}
// 正确的格式：{"contact_id":"http://www.google.com/m8/feeds/contacts/ywang@newcitycap.com/base/14bde6ca0d93a7fa","created_date":"2024-07-23T01:31:43.840181","device_type":"public_api","direction":"outbound","from_number":"+13109719619","id":"6497527001595904","message_status":"pending","target_id":"5243348849786880","target_type":"user","text":"Dear Tom, I hope this message finds you well. This is a follow-up ","to_numbers":["+18056604465"],"user_id":"5243348849786880"}
// SendSms phoneNumber: +18056604465 text:短信内容 dialpadUserId:5243348849786880 Yannan： +13109719619
func (c *DialpadUsecase) SendSms(phoneNumber string, text string, dialpadUserId string, caseId int32, fromType string) error {

	if fromType == "" {
		fromType = "Dialpad:SendSms"
	}

	// 此处拦截短信发送
	if configs.Enable_Unsubscribes_SMS {
		canSendSms, err := c.UnsubscribesUsecase.CanSendSms(phoneNumber)
		if err != nil {
			return err
		}
		if !canSendSms {
			er := c.LogUsecase.SaveLog(caseId, fromType+":UnsubFilter", map[string]interface{}{
				"canSendSms":  InterfaceToString(canSendSms),
				"phoneNumber": InterfaceToString(phoneNumber),
				"text":        text,
			})
			if er != nil {
				c.log.Error(er)
			}
			return nil
		}
	}

	api := fmt.Sprintf("%s/api/v2/sms?apikey=%s", c.ApiUrl(), c.AccessKey())
	/*
		{
		  "infer_country_code": false,
		  "to_numbers": [
		    "+18056604465"
		  ],
		  "text": "Dear Tom, I hope this message finds you well. This is a follow-up ",
		  "user_id": 5243348849786880
		}
	*/
	params := make(lib.TypeMap)
	params.Set("infer_country_code", false)
	params.Set("to_numbers", []string{phoneNumber})
	params.Set("text", text)
	params.Set("user_id", dialpadUserId)

	var res *string
	var err error

	// todo:lgl 暂时注解
	if configs.IsDev() {
		res = to.Ptr(`{"contact_id":"http://www.google.com/m8/feeds/contacts/ywang@newcitycap.com/base/14bde6ca0d93a7fa","created_date":"2024-07-23T01:31:43.840181","device_type":"public_api","direction":"outbound","from_number":"+13109719619","id":"6497527001595904","message_status":"pending","target_id":"5243348849786880","target_type":"user","text":"Dear Tom, I hope this message finds you well. This is a follow-up ","to_numbers":["+18056604465"],"user_id":"5243348849786880"}`)
	} else {
		res, _, err = lib.HTTPJson("POST", api, params.ToBytes())
	}

	logContent := ""
	if res != nil {
		logContent = *res
	}
	er := c.LogUsecase.SaveLog(caseId, fromType, map[string]interface{}{
		"params": InterfaceToString(params),
		"res":    logContent,
	})
	if er != nil {
		c.log.Error(er)
	}
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("res is nil")
	}
	return nil
}

// SendSmsNoFilter 直接发送，不经过过虑
func (c *DialpadUsecase) SendSmsNoFilter(phoneNumber string, text string, dialpadUserId string, caseId int32, fromType string) error {

	if fromType == "" {
		fromType = "Dialpad:SendSmsNoFilter"
	}

	api := fmt.Sprintf("%s/api/v2/sms?apikey=%s", c.ApiUrl(), c.AccessKey())
	/*
		{
		  "infer_country_code": false,
		  "to_numbers": [
		    "+18056604465"
		  ],
		  "text": "Dear Tom, I hope this message finds you well. This is a follow-up ",
		  "user_id": 5243348849786880
		}
	*/
	params := make(lib.TypeMap)
	params.Set("infer_country_code", false)
	params.Set("to_numbers", []string{phoneNumber})
	params.Set("text", text)
	params.Set("user_id", dialpadUserId)

	var res *string
	var err error

	// todo:lgl 暂时注解
	if configs.IsDev() {
		res = to.Ptr(`{"contact_id":"http://www.google.com/m8/feeds/contacts/ywang@newcitycap.com/base/14bde6ca0d93a7fa","created_date":"2024-07-23T01:31:43.840181","device_type":"public_api","direction":"outbound","from_number":"+13109719619","id":"6497527001595904","message_status":"pending","target_id":"5243348849786880","target_type":"user","text":"Dear Tom, I hope this message finds you well. This is a follow-up ","to_numbers":["+18056604465"],"user_id":"5243348849786880"}`)
	} else {
		res, _, err = lib.HTTPJson("POST", api, params.ToBytes())
	}

	logContent := ""
	if res != nil {
		logContent = *res
	}
	er := c.LogUsecase.SaveLog(caseId, fromType, map[string]interface{}{
		"params": InterfaceToString(params),
		"res":    logContent,
	})
	if er != nil {
		c.log.Error(er)
	}
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("res is nil")
	}
	return nil
}
