package config_vbc

import (
	"vbc/configs"
)

func DialpadWebhookSecret() string {
	return configs.EnvDialpadWebhookSecret()
}

// CaseMappingClient  key: client case field name;  value: client field name
var CaseMappingClient = map[string]string{
	"email":                         "email",
	"ssn":                           "ssn",
	"dob":                           "dob",
	"phone":                         "phone",
	"address":                       "address",
	"city":                          "city",
	"state":                         "state",
	"zip_code":                      "zip_code",
	"place_of_birth_city":           "place_of_birth_city",
	"place_of_birth_state_province": "place_of_birth_state_province",
	"place_of_birth_country":        "place_of_birth_country",
	"current_occupation":            "current_occupation",
	"timezone_id":                   "timezone_id",
	"apt_number":                    "apt_number",
}

func SyncFieldNamesForClient() (r []string) {
	for _, v := range CaseMappingClient {
		r = append(r, v)
	}
	return
}

func SyncFieldNamesForCase() (r []string) {
	for k, _ := range CaseMappingClient {
		r = append(r, k)
	}
	return
}

func GetSyncFieldNameByCaseForClient(caseFieldName string) (clientFieldName string) {
	for k, v := range CaseMappingClient {
		if k == caseFieldName {
			return v
		}
	}
	return ""
}

func GetSyncFieldNameByClientForCase(clientFieldName string) (caseFieldName string) {
	for k, v := range CaseMappingClient {
		if v == clientFieldName {
			return k
		}
	}
	return ""
}

// MailSubIdV1 mail template select
func MailSubIdV1(currentRating int32) int {
	if currentRating < 50 {
		return 50
	} else if currentRating < 70 {
		return 70
	} else if currentRating < 90 {
		return 90
	} else {
		return 100
	}
}

// SignContractIndexV1 计算合同模板的索引 和 邮件
func SignContractIndexV1(effectiveCurrentRating int32) int {
	if effectiveCurrentRating < 10 {
		return 0
	} else if effectiveCurrentRating < 20 {
		return 10
	} else if effectiveCurrentRating < 30 {
		return 20
	} else if effectiveCurrentRating < 40 {
		return 30
	} else if effectiveCurrentRating < 50 {
		return 40
	} else if effectiveCurrentRating < 60 {
		return 50
	} else if effectiveCurrentRating < 70 {
		return 60
	} else if effectiveCurrentRating < 80 {
		return 70
	} else if effectiveCurrentRating < 90 {
		return 80
	} else {
		return 90
	}

}

type FeeVo struct {
	Rating int
	Fee    int
}

type FeeVoConfigs map[int][]FeeVo

func (c FeeVoConfigs) GetByIndex(index int) []FeeVo {
	if c == nil {
		return nil
	}
	if _, ok := c[index]; ok {
		return c[index]
	}
	return nil
}

func (c FeeVoConfigs) Charge(currentRating int, newRating int) int {
	if _, ok := c[currentRating]; ok {
		for _, v := range c[currentRating] {
			if v.Rating == newRating {
				return v.Fee
			}
		}
	}
	return 0
}

// FeeDefine 此配置信息，必须保持顺序：50 70 90 100， 这样才能计算80的情况
var FeeDefine = FeeVoConfigs{
	-1: { // active duty
		{Rating: 70, Fee: 2000},
		{Rating: 90, Fee: 4000},
		{Rating: 100, Fee: 9000}},
	0: {
		{Rating: 50, Fee: 4000},
		{Rating: 70, Fee: 6000},
		{Rating: 90, Fee: 8000},
		{Rating: 100, Fee: 13000}},
	10: {
		{Rating: 50, Fee: 3000},
		{Rating: 70, Fee: 5000},
		{Rating: 90, Fee: 7000},
		{Rating: 100, Fee: 12000}},
	20: {
		{Rating: 50, Fee: 3000},
		{Rating: 70, Fee: 5000},
		{Rating: 90, Fee: 7000},
		{Rating: 100, Fee: 12000}},
	30: {
		{Rating: 50, Fee: 2000},
		{Rating: 70, Fee: 4000},
		{Rating: 90, Fee: 6000},
		{Rating: 100, Fee: 11000}},
	40: {
		{Rating: 50, Fee: 1000},
		{Rating: 70, Fee: 3000},
		{Rating: 90, Fee: 5000},
		{Rating: 100, Fee: 10000}},
	50: {
		{Rating: 70, Fee: 2000},
		{Rating: 90, Fee: 4000},
		{Rating: 100, Fee: 9000}},
	60: {
		{Rating: 70, Fee: 1000},
		{Rating: 90, Fee: 3000},
		{Rating: 100, Fee: 8000}},
	70: {
		{Rating: 90, Fee: 2000},
		{Rating: 100, Fee: 7000}},
	80: {
		{Rating: 90, Fee: 1000},
		{Rating: 100, Fee: 6000}},
	90: {{Rating: 100, Fee: 5000}},
}

var BoxSignTpl = map[string]string{
	"-1": "8a19aa3c-689b-4e46-a4ef-cc2bf318144d",
	"0":  "3befc364-7faf-417d-90cc-c56106b73fcd",
	"10": "86df0f08-497c-4960-9121-afbcd4f4251b",
	"20": "99be37cf-da2b-44fe-ac1d-fd88d7314290",
	"30": "0e442e60-2dae-48ff-84e3-2d71fa5eb477",
	"40": "3e2e553a-d493-43a1-b603-7752dba74abc",
	"50": "3fd34073-ea89-4e68-b215-a4f76a4f3a25",
	"60": "0d862f95-b60e-4a61-8c55-d7d962f64d27",
	"70": "599f3078-6e28-4a16-b1b0-1ac834220c1d",
	"80": "757aef47-93a3-42d2-9a7e-8b19ac9550aa",
	"90": "31807a9a-70f7-4b16-9492-10c0b274375a",
}
