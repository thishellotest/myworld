package biz

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	"vbc/internal/config_zoho"
	"vbc/lib"
	"vbc/lib/uuid"
)

type WebsiteUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	AsanaUsecase     *AsanaUsecase
	ZohoUsecase      *ZohoUsecase
	DataEntryUsecase *DataEntryUsecase
}

func NewWebsiteUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	AsanaUsecase *AsanaUsecase,
	ZohoUsecase *ZohoUsecase,
	DataEntryUsecase *DataEntryUsecase) *WebsiteUsecase {
	uc := &WebsiteUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		AsanaUsecase:     AsanaUsecase,
		ZohoUsecase:      ZohoUsecase,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

func (c *WebsiteUsecase) SyncToZohoOrVBCRM(formData string) error {
	data := lib.ToTypeMapByString(formData)

	//lastName := data.GetString("data.field:comp-lozdvdjo2")
	//firstName := data.GetString("data.field:comp-lozdvdje")
	//shortState := data.GetString("data.field:comp-loze9sgy")
	//email := data.GetString("data.field:comp-lozdvdjr")
	//phone := data.GetString("data.field:comp-lozg1qa5")
	//fullState := vbc_config.StateConfigs.FullNameByShort(shortState)
	//desc := data.GetString("data.field:comp-lozdvdku1")

	lastName := data.GetString("data.contact.name.last")
	firstName := data.GetString("data.contact.name.first")
	email := data.GetString("data.contact.email")
	phone := data.GetString("data.contact.phone")
	var shortState string
	var desc string
	wixFormId := data.GetString("formId")
	if wixFormId == "comp-lozdvdj82" { // VBC Contact Form
		shortState = data.GetString("data.field:comp-loze9sgy")
		desc = data.GetString("data.field:comp-lozdvdku1")
	} else if wixFormId == "comp-m5zftakk3" { // VBC Contact Form 3
		shortState = data.GetString("data.field:comp-m5zftaks8")
		formLastName := data.GetString("data.field:comp-m5zftakq6")
		formFirstName := data.GetString("data.field:comp-m5zftakn6")
		formEmail := data.GetString("data.field:comp-m5zftakr6")
		formPhone := data.GetString("data.field:comp-m5zftaks1")
		if formLastName != "" {
			lastName = formLastName
		}
		if formFirstName != "" {
			firstName = formFirstName
		}
		if formEmail != "" {
			email = formEmail
		}
		if formPhone != "" {
			phone = formPhone
		}
	}
	return c.BizSyncToZohoOrVBCRM(firstName, lastName, email, phone, shortState, "", desc, "", "")
}

func (c *WebsiteUsecase) BizSyncToZohoOrVBCRM(firstName, lastName, email, phone, shortState, fullState, desc, leadSource string, branch string) error {

	if fullState == "" {
		fullState = config_vbc.StateConfigs.FullNameByShort(shortState)
	}

	if email == "" && phone == "" {
		return errors.New("Email and Phone is empty.")
	}
	zohoDescription := fmt.Sprintf("Objective: %s\n\nService:\n\nCurrent:\n\nNew:\n\n", desc)

	firstName = lib.Capitalize(strings.ToLower(firstName))
	lastName = lib.Capitalize(strings.ToLower(lastName))
	if leadSource == "" {
		leadSource = config_vbc.Source_Website
	}

	record := make(lib.TypeMap)
	record.Set("First_Name", firstName)
	record.Set("Last_Name", lastName)
	record.Set("Email", email)
	record.Set("State", fullState)
	record.Set("Description", zohoDescription)
	record.Set("Lead_Source", leadSource)

	if configs.StoppedZoho {
		userGid := ""
		if configs.IsProd() && email != "liaogling@gmail.com" && email != "18891706@qq.com" && email != "lialing@foxmail.com" {
			userGid = config_vbc.User_Edward_gid
		} else {
			userGid = config_vbc.User_Dev_gid
		}
		clientGid := uuid.UuidWithoutStrike()
		clientData := make(TypeDataEntry)
		clientData[DataEntry_gid] = clientGid
		clientData[DataEntry_user_gid] = userGid
		clientData[FieldName_first_name] = firstName
		clientData[FieldName_last_name] = lastName
		clientData[FieldName_email] = email
		clientData[FieldName_phone] = phone
		clientData[FieldName_state] = fullState
		clientData[FieldName_source] = leadSource
		clientData[FieldName_branch] = branch

		_, err := c.DataEntryUsecase.HandleOne(Kind_clients, clientData, DataEntry_gid, nil)
		if err != nil {
			return err
		}

		caseGid := uuid.UuidWithoutStrike()
		caseData := make(TypeDataEntry)
		caseData[DataEntry_gid] = caseGid
		caseData[DataEntry_user_gid] = userGid
		caseData[FieldName_client_gid] = clientGid
		caseData[FieldName_description] = zohoDescription
		caseData[FieldName_email] = email
		caseData[FieldName_phone] = phone
		caseData[FieldName_state] = fullState
		caseData[FieldName_source] = leadSource
		caseData[FieldName_branch] = branch
		caseData[FieldName_stages] = config_vbc.Stages_AmIncomingRequest
		caseData[FieldName_ContractSource] = ContractSource_AM

		_, err = c.DataEntryUsecase.HandleOne(Kind_client_cases, caseData, DataEntry_gid, nil)
		if err != nil {
			return err
		}
	} else {
		record.Set("Mobile", phone)
		if configs.IsProd() && email != "liaogling@gmail.com" && email != "18891706@qq.com" && email != "lialing@foxmail.com" {
			record.Set("Owner.id", config_vbc.User_Edward_gid)
		} else {
			record.Set("Owner.id", config_vbc.User_Dev_gid)
		}

		// 6159272000001347005
		// todo:lgl stopzoho
		leadGid, _, err := c.ZohoUsecase.CreateRecord(config_zoho.Leads, record)
		if err != nil {
			return err
		}
		if leadGid == "" {
			return errors.New("SyncToZoho: leadGid is empty.")
		}
	}
	return nil
}

// SyncToAsana_Deprecated 不再使用
func (c *WebsiteUsecase) SyncToAsana_Deprecated(formData string) error {
	data := lib.ToTypeMapByString(formData)
	firstName := data.GetString("data.field:comp-lozdvdjo2")
	lastName := data.GetString("data.field:comp-lozdvdje")
	shortState := data.GetString("data.field:comp-loze9sgy")
	email := data.GetString("data.field:comp-lozdvdjr")
	phone := data.GetString("data.contact.phone")
	fullState := config_vbc.StateConfigs.FullNameByShort(shortState)

	desc := data.GetString("data.field:comp-lozdvdku1")

	if email == "" && phone == "" {
		return errors.New("Email and Phone is empty.")
	}

	field := config_vbc.GetAsanaCustomFields()
	firstNameGid := field.GetByName("First Name").GetGid()
	lastNameGid := field.GetByName("Last Name").GetGid()
	emailGid := field.GetByName("Email").GetGid()
	phoneNumberGid := field.GetByName("Phone Number").GetGid()
	stateGid := field.GetByName("Address - State").GetGid()
	stateEnumGid := field.GetByName("Address - State").GetEnumGidByName(fullState)

	sourceGid := field.GetByName("Source").GetGid()
	websiteEnumGid := field.GetByName("Source").GetEnumGidByName(config_vbc.Source_Website)

	customFields := make(lib.TypeMap)
	customFields.Set(firstNameGid, firstName)
	customFields.Set(lastNameGid, lastName)
	customFields.Set(emailGid, email)
	customFields.Set(phoneNumberGid, phone)
	if stateEnumGid != "" {
		customFields.Set(stateGid, stateEnumGid)
	}
	customFields.Set(sourceGid, websiteEnumGid)
	_, err := c.AsanaUsecase.CreateATask(customFields, firstName, lastName, desc)
	return err
}
