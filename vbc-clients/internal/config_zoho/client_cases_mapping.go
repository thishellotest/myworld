package config_zoho

const (
	Zoho_Claims_Online    = "Claims_Online"
	Zoho_Statements       = "Statements"
	Zoho_SF_180_Submitted = "SF_180_Submitted"
	Zoho_Collaborators    = "Collaborators"

	Zoho_SF_180 = "SF_180"
	Zoho_FOIA   = "FOIA"
)

var ClientCasesMappingConfigs = map[string]string{
	"Deal_Name":                "deal_name",
	"Owner.id":                 "user_gid",
	"Current_Rating":           "current_rating",
	"Branch":                   "branch",
	"Stage":                    "stages",
	"id":                       "gid",
	"Retired":                  "retired",
	"New_Rating":               "new_rating",
	"Contact_Name.id":          "client_gid",
	"Effective_Current_Rating": "effective_current_rating",
	"Modified_Time":            "modified_time",
	"Created_Time":             "created_time",
	"Lead_Source":              "source",
	//"Pipeline":                               "pipeline",
	"Contact_Form":   "contact_form",
	"ITF_Expiration": "itf_expiration",
	"Filing_Date":    "itf_date",
	//"Approval_Date":                          "claim_decision_date",
	"Agent_Orange_Exposure":                  "agent_orange",
	"Burn_Pits_and_Other_Airborne_Hazards":   "burn_pits",
	"Gulf_War_Illness":                       "gulf_war",
	"Atomic_Veterans_and_Radiation_Exposure": "atomic_veterans",
	"Amyotrophic_Lateral_Sclerosis_ALS":      "amyotrophic",
	"Amount":                                 "amount",
	"Active_Duty":                            "active_duty",
	"Description":                            "description",
	"Email":                                  "email",
	"Phone":                                  "phone",
	"SSN":                                    "ssn",
	"Date_of_Birth":                          "dob",
	"State":                                  "state",
	"City":                                   "city",
	"Street_Address":                         "address",
	"Zip_Code":                               "zip_code",
	"Case_Files_Folder":                      "case_files_folder",
	"Primary_CP":                             "primary_cp",
	"Private_Exams_Needed":                   "private_exams_needed",
	"Data_Collection_Folder":                 "data_collection_folder",
	"Pricing_Version":                        "pricing_version",
	"Year_Entering_Service":                  "year_entering_service",
	"Year_Separate_from_Service":             "year_separate_from_service",
	"Service_Connections":                    "service_connections",
	"Previous_Denials":                       "previous_denials",
	Zoho_Claims_Online:                       "claims_online",
	"Claims_Next_Round":                      "claims_next_round",
	"Claims_Supplemental":                    "claims_supplemental",
	"Place_of_Birth_City":                    "place_of_birth_city",
	"Place_of_Birth_State_Province":          "place_of_birth_state_province",
	"Place_of_Birth_Country":                 "place_of_birth_country",
	"Current_Occupation":                     "current_occupation",
}

var ClientCasesMappingConfigs2 = map[string]string{
	// 此处开始使用任务2处理
	"id":            "gid",
	"Modified_Time": "modified_time",

	"Lead_CO":                         "lead_co",
	"Referring_Person":                "referrer",
	"Answer_to_Presumptive_Questions": "answer_to_presumptive_questions",
	"Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun": "illness_due",
	"Service_Years":                     "service_years",
	"Occupation_during_Service":         "occupation_during_service",
	"C_File_Submitted":                  "record_request_submitted",
	"DD214":                             "dd214",
	"Disability_Rating_List_Screenshot": "disability_rating",
	"Rating_Decision_Letters":           "rating_decision",
	"STRs":                              "strs",
	"Benefits_Summary_Letter":           "benefits_summary_letter",
	"VA_Records":                        "va_records",
	"Private_Records":                   "private_records",
	"VA_Healthcare_Registration":        "va_healthcare",
	"TDIU":                              "tdiu",
	"Tinnitus_Nexus":                    "tinnitus_nexus",
	"Primary_VS":                        "primary_vs",
	Zoho_Statements:                     "statements",
	Zoho_SF_180_Submitted:               "sf_180_submitted",
	Zoho_FOIA:                           "foia",
	Zoho_SF_180:                         "sf_180",
}

func ZohoDealVbcFieldNameByZohoFieldName(zohoFieldName string) string {
	for k, v := range ClientCasesMappingConfigs {
		if k == zohoFieldName {
			return v
		}
	}
	for k, v := range ClientCasesMappingConfigs2 {
		if k == zohoFieldName {
			return v
		}
	}
	return ""
}

func ZohoDealFieldNameByVbcFieldName(vbcFieldName string) string {
	for k, v := range ClientCasesMappingConfigs {
		if v == vbcFieldName {
			return k
		}
	}
	for k, v := range ClientCasesMappingConfigs2 {
		if v == vbcFieldName {
			return k
		}
	}
	return ""
}
