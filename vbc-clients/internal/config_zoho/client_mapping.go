package config_zoho

import (
	"vbc/lib"
)

const Zoho_Date_of_Birth = "Date_of_Birth"

var ClientMappingConfigs = map[string]string{
	"Owner.id":                      "user_gid",
	"Email":                         "email",
	Zoho_Date_of_Birth:              "dob",
	"Mobile":                        "phone",
	"SSN":                           "ssn",
	"Mailing_City":                  "city",
	"Mailing_State":                 "state",
	"Mailing_Street":                "address",
	"Mailing_Zip":                   "zip_code",
	"First_Name":                    "first_name",
	"Last_Name":                     "last_name",
	"Current_Rating":                "current_rating",
	"Branch":                        "branch",
	"id":                            "gid",
	"Referring_Person":              "referrer",
	"Retired":                       "retired",
	"Effective_Current_Rating":      "effective_current_rating",
	"Lead_Source":                   "source",
	"Modified_Time":                 "modified_time",
	"Created_Time":                  "created_time",
	"Place_of_Birth_City":           "place_of_birth_city",
	"Place_of_Birth_State_Province": "place_of_birth_state_province",
	"Place_of_Birth_Country":        "place_of_birth_country",
	"Active_Duty":                   "active_duty",
	"Current_Occupation":            "current_occupation",
}

func ZohoContactVbcFieldNameByZohoFieldName(zohoFieldName string) string {
	for k, v := range ClientMappingConfigs {
		if k == zohoFieldName {
			return v
		}
	}
	return ""
}

func ZohoContactFieldNameByVbcFieldName(vbcFieldName string) string {
	for k, v := range ClientMappingConfigs {
		if v == vbcFieldName {
			return k
		}
	}
	return ""
}

func ClientMappings(row lib.TypeMap) lib.TypeMap {
	if row == nil {
		return nil
	}
	res := make(lib.TypeMap)
	for k, v := range ClientMappingConfigs {
		//if k == Zoho_Date_of_Birth {
		//	val := row.GetString(k)
		//	if val != "" {
		//		t, _ := lib.TimeString(val, "2006-01-02")
		//		val = lib.InterfaceToString(t.Unix())
		//	}
		//	res.Set(v, val)
		//} else {
		if row.Get(k) != nil {
			res.Set(v, row.GetString(k))
		}
		//}
	}
	return res
}
