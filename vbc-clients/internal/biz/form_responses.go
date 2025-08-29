package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	. "vbc/lib/builder"
)

const SourceFrom_google_intake_form = "google_intake_form"

type FormResponseEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	SourceForm         string
	Data               string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (FormResponseEntity) TableName() string {
	return "form_responses"
}

type FormResponseUsecase struct {
	CommonUsecase *CommonUsecase
	log           *log.Helper
	DBUsecase[FormResponseEntity]
	BaseHandle[FormResponseEntity]
	TUsecase          *TUsecase
	AsanaUsecase      *AsanaUsecase
	BehaviorUsecase   *BehaviorUsecase
	DataComboUsecase  *DataComboUsecase
	ZohoUsecase       *ZohoUsecase
	ActionOnceUsecase *ActionOnceUsecase
}

func NewFormResponseUsecase(CommonUsecase *CommonUsecase, logger log.Logger, TUsecase *TUsecase,
	AsanaUsecase *AsanaUsecase,
	BehaviorUsecase *BehaviorUsecase,
	DataComboUsecase *DataComboUsecase,
	ZohoUsecase *ZohoUsecase,
	ActionOnceUsecase *ActionOnceUsecase) *FormResponseUsecase {
	uc := &FormResponseUsecase{
		CommonUsecase:     CommonUsecase,
		log:               log.NewHelper(logger),
		TUsecase:          TUsecase,
		AsanaUsecase:      AsanaUsecase,
		BehaviorUsecase:   BehaviorUsecase,
		DataComboUsecase:  DataComboUsecase,
		ZohoUsecase:       ZohoUsecase,
		ActionOnceUsecase: ActionOnceUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandle.DB = CommonUsecase.DB()
	uc.BaseHandle.Log = log.NewHelper(logger)
	uc.BaseHandle.TableName = FormResponseEntity{}.TableName()
	uc.BaseHandle.Handle = uc.Handle

	return uc
}

func (c *FormResponseUsecase) Handle(ctx context.Context, task *FormResponseEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix()
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.HandleResultDetail = err.Error()
	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Save(task).Error
}

/*
{"VBC Case ID":"2466789032","First Name":"TestGary","Last Name":"TestLiao","Street Address":"addr55","City":"city55","State":"Arizona","Zip Code":"10000","Phone Number":"100-100-1000","SSN":"100-00-0055","Date of Birth":"1987-10-22","Overall Rating":"00","Branch of Service ":["Army","Navy","Marine Corps","Air Force","Space Force","Coast Guard","National Oceanic and Atmospheric Administration","Public Health Service","Army National Guard","Air National Guard"],"Have you retired from US Military services?":"No","Are you currently active duty in the US Military?":"No","Military Toxic Exposures":["Agent Orange Exposure","Gulf War Illness","Burn Pits and Other Airborne Hazards","Illness Due to Toxic Drinking Water at Camp Lejeune","\"Atomic Veterans\" and Radiation Exposure","Amyotrophic Lateral Sclerosis (ALS)"]}
*/
func (c *FormResponseUsecase) HandleExec(ctx context.Context, task *FormResponseEntity) error {
	if task == nil {
		return errors.New("task is nil.")
	}
	dataMap := lib.ToTypeMapByString(task.Data)

	VBCCaseID := dataMap.GetString("VBC Case ID")
	if VBCCaseID == "" {
		return errors.New("VBC Case ID is empty")
	}

	firstName := dataMap.GetString("First Name")
	lastName := dataMap.GetString("Last Name")
	if len(firstName) == 0 {
		return errors.New("firstName is empty")
	}
	if len(lastName) == 0 {
		return errors.New("lastName is empty")
	}
	//tClient, err := c.TUsecase.Data(Kind_client_cases, And(Eq{"first_name": firstName}, Eq{"last_name": lastName}))
	tClientCase, err := c.TUsecase.Data(Kind_client_cases, And(Eq{FieldName_uniqcode: VBCCaseID}))
	if err != nil {
		return err
	}
	if tClientCase == nil {
		return errors.New("The tClientCase does not exist via VBC Case ID.")
	}
	clientCaseGid := tClientCase.CustomFields.TextValueByNameBasic("gid")
	if len(clientCaseGid) == 0 {
		return errors.New("clientCaseGid is empty")
	}

	_, tContactFields, err := c.DataComboUsecase.Client(tClientCase.CustomFields.TextValueByNameBasic("client_gid"))
	if err != nil {
		return err
	}
	if tContactFields == nil {
		return errors.New("tContactFields is nil.")
	}
	clientGid := tContactFields.TextValueByNameBasic("gid")
	if clientGid == "" {
		return errors.New("clientGid is empty.")
	}

	contactMap := c.HandeMappingContact(dataMap)
	if len(contactMap) == 0 {
		return errors.New("contactMap Nothing need modify.")
	}
	contactMap.Set("id", clientGid)

	_, _, err = c.ZohoUsecase.PutRecordV1(config_zoho.Contacts, contactMap)
	if err != nil {
		return err
	}
	err = c.HandleExecDeal(clientCaseGid, dataMap)
	if err != nil {
		return err
	}

	err = c.BehaviorUsecase.Add(tClientCase.CustomFields.NumberValueByNameBasic("id"),
		BehaviorType_complete_intake_form, time.Now(), "")
	if err != nil {
		return err
	} else {
		// - Move task from "Getting Started Email" to " Awaiting Client Files"
		err = c.ActionOnceUsecase.StageGettingStartedEmailToAwaitingClientFiles(tClientCase.CustomFields.NumberValueByNameBasic("id"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FormResponseUsecase) HandleExecDeal(clientCaseGid string, dataMap lib.TypeMap) error {

	DealFieldInfos := config_zoho.DealLayout().DealFieldInfos()

	military := lib.InterfaceToTDef[[]string](dataMap.Get("Military Toxic Exposures"), nil)
	destMap := make(lib.TypeMap)
	destMap.Set("id", clientCaseGid)
	destMap.Set("Agent_Orange_Exposure", "NO")
	destMap.Set("Gulf_War_Illness", "NO")
	destMap.Set("Burn_Pits_and_Other_Airborne_Hazards", "NO")
	destMap.Set("Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun", "NO")
	destMap.Set("Atomic_Veterans_and_Radiation_Exposure", "NO")
	destMap.Set("Amyotrophic_Lateral_Sclerosis_ALS", "NO")
	for _, v := range military {
		apiName := config_zoho.ApiNameByFieldLabel(DealFieldInfos, v)
		if apiName != "" {
			destMap.Set(apiName, "YES")
		} else {
			if strings.Index(v, "Illness") >= 0 {
				destMap.Set("Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun", "YES")
			} else if strings.Index(v, "Atomic") >= 0 {
				destMap.Set("Atomic_Veterans_and_Radiation_Exposure", "YES")
			}
		}
	}
	branch := FormResponseBranch(dataMap)
	if branch != "" {
		destMap.Set("Branch", branch)
	}
	if len(dataMap.GetString("Overall Rating")) > 0 {
		destMap.Set("Current_Rating", dataMap.GetInt("Overall Rating"))
	}
	retired := dataMap.GetString("Have you retired from US Military services?")
	if len(retired) > 0 {
		destMap.Set("Retired", retired)
	}
	activeDuty := dataMap.GetString("Are you currently active duty in the US Military?")
	if len(activeDuty) > 0 {
		destMap.Set("Active_Duty", activeDuty)
	}
	destMap.Set("Contact_Form", "Yes")
	_, _, err := c.ZohoUsecase.PutRecordV1(config_zoho.Deals, destMap)
	return err
}

func (c *FormResponseUsecase) HandeMappingContact(dataMap lib.TypeMap) lib.TypeMap {
	destMap := make(lib.TypeMap)
	if len(dataMap.GetString("Street Address")) > 0 {
		destMap.Set("Mailing_Street", dataMap.GetString("Street Address"))
	}
	if len(dataMap.GetString("City")) > 0 {
		destMap.Set("Mailing_City", dataMap.GetString("City"))
	}
	if len(dataMap.GetString("State")) > 0 {
		destMap.Set("Mailing_State", dataMap.GetString("State"))
	}
	if len(dataMap.GetString("Zip Code")) > 0 {
		destMap.Set("Mailing_Zip", dataMap.GetString("Zip Code"))
	}
	if len(dataMap.GetString("Phone Number")) > 0 {
		destMap.Set("Mobile", dataMap.GetString("Phone Number"))
	}
	if len(dataMap.GetString("SSN")) > 0 {
		destMap.Set("SSN", dataMap.GetString("SSN"))
	}
	if len(dataMap.GetString("Date of Birth")) > 0 {
		destMap.Set("Date_of_Birth", dataMap.GetString("Date of Birth"))
	}
	if len(dataMap.GetString("Overall Rating")) > 0 {
		destMap.Set("Current_Rating", dataMap.GetInt("Overall Rating"))
	}
	branch := FormResponseBranch(dataMap)
	if branch != "" {
		destMap.Set("Branch", branch)
	}
	return destMap
}

func FormResponseBranch(dataMap lib.TypeMap) string {
	BranchOfService := dataMap.Get("Branch of Service ")
	BranchOfServiceList := lib.InterfaceToTDef[[]string](BranchOfService, nil)
	if len(BranchOfServiceList) > 0 {
		return FormResponseBranchMappingToZoho(BranchOfServiceList[0])
	}
	return ""
}

func FormResponseBranchMappingToZoho(formBranch string) string {

	if formBranch == "Air National Guard" {
		return "Air NG"
	} else if formBranch == "Army National Guard" {
		return "Army NG"
	} else if formBranch == "Air Force" ||
		formBranch == "Army" ||
		formBranch == "Coast Guard" ||
		formBranch == "Marine Corps" ||
		formBranch == "Navy" ||
		formBranch == "Space Force" {
		return formBranch
	}
	return ""
}

func (c *FormResponseUsecase) HandeMapping(dataMap lib.TypeMap) lib.TypeMap {
	asanaField := config_vbc.GetAsanaCustomFields()
	destMap := make(lib.TypeMap)

	firstNameGid := asanaField.GetByName("First Name").GetGid()
	firstName := dataMap.GetString("First Name")
	if firstNameGid != "" {
		destMap.Set(firstNameGid, firstName)
	}
	lastNameGid := asanaField.GetByName("Last Name").GetGid()
	lastName := dataMap.GetString("Last Name")
	if lastNameGid != "" {
		destMap.Set(lastNameGid, lastName)
	}

	SSNGid := asanaField.GetByName("SSN").GetGid()
	SSN := dataMap.GetString("SSN")
	if SSN != "" && SSNGid != "" {
		destMap.Set(SSNGid, SSN)
	}
	phoneNumberGid := asanaField.GetByName("Phone Number").GetGid()
	phoneNumber := dataMap.GetString("Phone Number")
	if phoneNumber != "" && phoneNumberGid != "" {
		destMap.Set(phoneNumberGid, phoneNumber)
	}
	StreetAddressGid := asanaField.GetByName("Street Address").GetGid()
	StreetAddress := dataMap.GetString("Street Address")
	//fmt.Println("sss StreetAddressGid:", StreetAddressGid, StreetAddress)
	if StreetAddress != "" && StreetAddressGid != "" {
		destMap.Set(StreetAddressGid, StreetAddress)
	}
	DOBGid := asanaField.GetByName("DOB").GetGid()
	dateOfBirth := dataMap.GetString("Date of Birth")
	if len(dateOfBirth) > 0 && DOBGid != "" {
		//dob, _ := time.ParseInLocation(time.DateOnly, dateOfBirth, lib.LoadLocation)
		//fmt.Println(dob.Format(time.DateOnly))
		destMap.Set(DOBGid+".date", dateOfBirth)
	}
	// 1205964024662896 Current Rating
	CurrentRatingGid := asanaField.GetByName("Current Rating").GetGid()
	overallRating := dataMap.GetInt("Overall Rating")
	if CurrentRatingGid != "" {
		destMap.Set(CurrentRatingGid, overallRating)
	}

	// 1206398481016928 Address - State
	AddressStateGid := asanaField.GetByName("Address - State").GetGid()
	State := dataMap.GetString("State")
	//fmt.Println("sss AddressStateGid:", AddressStateGid, State)
	if State != "" && AddressStateGid != "" {
		gid := asanaField.GetByName("Address - State").GetEnumGidByName(State)
		if gid != "" {
			destMap.Set(AddressStateGid, gid)
		}
	}
	// 1206401684235934 Zip Code
	AddressZipCodeGid := asanaField.GetByName("Address - Zip Code").GetGid()
	zipCode := dataMap.GetString("Zip Code")
	if zipCode != "" && AddressZipCodeGid != "" {
		destMap.Set(AddressZipCodeGid, zipCode)
	}

	AddressCityGid := asanaField.GetByName("Address - City").GetGid()
	// 1206401658215277 Address - City
	City := dataMap.GetString("City")
	if City != "" && AddressCityGid != "" {
		destMap.Set(AddressCityGid, City)
	}
	// 1206422409732583 Retired enum
	// 1206422409732584 Yes // 1206422409732585 No
	retired := dataMap.GetString("Have you retired from US Military services?")

	if len(retired) > 0 {
		field := asanaField.GetByName("Retired")
		if field != nil {
			gid := field.GetEnumGidByName(retired)
			if gid != "" {
				destMap.Set(lib.TypeMap(field).GetString("gid"), gid)
			}
		}
		//if retired == "Yes" {
		//	destMap.Set("1206401658215277", "1206422409732584")
		//} else {
		//	destMap.Set("1206401658215277", "1206422409732585")
		//}
	}
	// Military Toxic Exposures
	military := dataMap.Get("Military Toxic Exposures")
	militarys := lib.InterfaceToTDef[[]string](military, nil)

	for _, v := range militarys {
		field := asanaField.GetByName(v)
		if field != nil {
			gid := field.GetEnumGidByName("Yes")
			if gid != "" {
				destMap.Set(lib.TypeMap(field).GetString("gid"), gid)
			}
		}
	}

	BranchofService := dataMap.Get("Branch of Service ")
	BranchofServiceList := lib.InterfaceToTDef[[]string](BranchofService, nil)
	branchGid := lib.TypeMap(asanaField.GetByName("Branch")).GetString("gid")
	if branchGid != "" {
		var enumGid string
		for _, v := range BranchofServiceList {
			enumGid = asanaField.GetByName("Branch").GetEnumGidByName(v)
			if enumGid != "" {
				break
			}
		}
		if enumGid != "" {
			destMap.Set(branchGid, enumGid)
		}
	}

	return destMap
}

type TypeAddressState lib.TypeList

func (c TypeAddressState) GidByName(name string) string {
	for _, v := range c {
		if v.GetString("name") == name {
			return v.GetString("gid")
		}
	}
	return ""
}

//
//func GetTypeAddressState() TypeAddressState {
//	return TypeAddressState(lib.ToTypeListByString(addressState))
//}

const addressState = `[{
				"gid": "1206398481016929",
				"color": "none",
				"enabled": true,
				"name": "Alabama",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016930",
				"color": "none",
				"enabled": true,
				"name": "Alaska",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016931",
				"color": "none",
				"enabled": true,
				"name": "Arizona",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016932",
				"color": "none",
				"enabled": true,
				"name": "Arkansas",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016933",
				"color": "none",
				"enabled": true,
				"name": "American Samoa",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016934",
				"color": "none",
				"enabled": true,
				"name": "California",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016935",
				"color": "none",
				"enabled": true,
				"name": "Colorado",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016936",
				"color": "none",
				"enabled": true,
				"name": "Connecticut",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016937",
				"color": "none",
				"enabled": true,
				"name": "Delaware",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016938",
				"color": "none",
				"enabled": true,
				"name": "District of Columbia",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016939",
				"color": "none",
				"enabled": true,
				"name": "Florida",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016940",
				"color": "none",
				"enabled": true,
				"name": "Georgia",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016941",
				"color": "none",
				"enabled": true,
				"name": "Guam",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016942",
				"color": "none",
				"enabled": true,
				"name": "Hawaii",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016943",
				"color": "none",
				"enabled": true,
				"name": "Idaho",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016944",
				"color": "none",
				"enabled": true,
				"name": "Illinois",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016945",
				"color": "none",
				"enabled": true,
				"name": "Indiana",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016946",
				"color": "none",
				"enabled": true,
				"name": "Iowa",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016947",
				"color": "none",
				"enabled": true,
				"name": "Kansas",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016948",
				"color": "none",
				"enabled": true,
				"name": "Kentucky",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016949",
				"color": "none",
				"enabled": true,
				"name": "Louisiana",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016950",
				"color": "none",
				"enabled": true,
				"name": "Maine",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016951",
				"color": "none",
				"enabled": true,
				"name": "Maryland",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016952",
				"color": "none",
				"enabled": true,
				"name": "Massachusetts",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016953",
				"color": "none",
				"enabled": true,
				"name": "Michigan",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016954",
				"color": "none",
				"enabled": true,
				"name": "Minnesota",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016955",
				"color": "none",
				"enabled": true,
				"name": "Mississippi",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016956",
				"color": "none",
				"enabled": true,
				"name": "Missouri",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016957",
				"color": "none",
				"enabled": true,
				"name": "Montana",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016958",
				"color": "none",
				"enabled": true,
				"name": "Nebraska",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016959",
				"color": "none",
				"enabled": true,
				"name": "Nevada",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016960",
				"color": "none",
				"enabled": true,
				"name": "New Hampshire",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016961",
				"color": "none",
				"enabled": true,
				"name": "New Jersey",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016962",
				"color": "none",
				"enabled": true,
				"name": "New Mexico",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016963",
				"color": "none",
				"enabled": true,
				"name": "New York",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016964",
				"color": "none",
				"enabled": true,
				"name": "North Carolina",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016965",
				"color": "none",
				"enabled": true,
				"name": "North Dakota",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016966",
				"color": "none",
				"enabled": true,
				"name": "Northern Mariana Islands",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016967",
				"color": "none",
				"enabled": true,
				"name": "Ohio",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016968",
				"color": "none",
				"enabled": true,
				"name": "Oklahoma",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016969",
				"color": "none",
				"enabled": true,
				"name": "Oregon",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016970",
				"color": "none",
				"enabled": true,
				"name": "Pennsylvania",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016971",
				"color": "none",
				"enabled": true,
				"name": "Puerto Rico",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016972",
				"color": "none",
				"enabled": true,
				"name": "Rhode Island",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016973",
				"color": "none",
				"enabled": true,
				"name": "South Carolina",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016974",
				"color": "none",
				"enabled": true,
				"name": "South Dakota",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016975",
				"color": "none",
				"enabled": true,
				"name": "Tennessee",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016976",
				"color": "none",
				"enabled": true,
				"name": "Texas",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016977",
				"color": "none",
				"enabled": true,
				"name": "Trust Territories",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016978",
				"color": "none",
				"enabled": true,
				"name": "Utah",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016979",
				"color": "none",
				"enabled": true,
				"name": "Vermont",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016980",
				"color": "none",
				"enabled": true,
				"name": "Virginia",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016981",
				"color": "none",
				"enabled": true,
				"name": "Virgin Islands",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016982",
				"color": "none",
				"enabled": true,
				"name": "Washington",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016983",
				"color": "none",
				"enabled": true,
				"name": "West Virginia",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016984",
				"color": "none",
				"enabled": true,
				"name": "Wisconsin",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016985",
				"color": "none",
				"enabled": true,
				"name": "Wyoming",
				"resource_type": "enum_option"
			}, {
				"gid": "1206398481016986",
				"color": "none",
				"enabled": true,
				"name": "Outside United States",
				"resource_type": "enum_option"
			}]`
