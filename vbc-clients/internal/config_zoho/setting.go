package config_zoho

import (
	"vbc/lib"
)

type ZohoLayout struct {
	module    string
	source    string
	sourceMap lib.TypeMap
	sections  lib.TypeList
}

func NotesLayout() *ZohoLayout {
	return newZohoLayout(notesLayout, Notes)
}

func LeadsLayout() *ZohoLayout {
	return newZohoLayout(leadsLayout, Leads)
}

func ContactLayout() *ZohoLayout {
	return newZohoLayout(contactLayout, Contacts)
}

func DealLayout() *ZohoLayout {
	return newZohoLayout(dealLayout, Deals)
}

func newZohoLayout(source string, module string) *ZohoLayout {

	sourceMap := lib.ToTypeMapByString(source)
	layouts := lib.ToTypeList(sourceMap.Get("layouts"))
	sections := lib.ToTypeList(layouts[0].Get("sections"))

	return &ZohoLayout{
		module:    module,
		source:    source,
		sourceMap: sourceMap,
		sections:  sections,
	}
}

const (
	Lead_Sections_ApiName_Lead_Information        = "Lead Information"
	Lead_Sections_ApiName_Address_Information     = "Address Information"
	Lead_Sections_ApiName_Description_Information = "Description Information"

	Contact_Sections_ApiName_Client_Image            = "Client Image" // 暂不使用
	Contact_Sections_ApiName_Client_Information      = "Client Information"
	Contact_Sections_ApiName_Address_Information     = "Address Information"
	Contact_Sections_ApiName_Description_Information = "Description Information" // 暂不使用

	Deal_Sections_ApiName_ClientCaseImage         = "Client Case Image" // 暂不使用
	Deal_Sections_ApiName_ClientCaseInformation   = "Client Case Information"
	Deal_Sections_ApiName_ClientInformation       = "Client Information"
	Deal_Sections_ApiName_Description_Information = "Claims Information" // 暂不使用
	Deal_Sections_ApiName_Presumptive_Information = "Presumptive Information"
	Deal_Sections_ApiName_Service_Information     = "Service Information"

	Task_Sections_ApiName_TaskInformation        = "Task Information"
	Task_Sections_ApiName_DescriptionInformation = "Description Information"

	Notes_Sections_ApiName_NotesInformation = "Note Information"
)

func (c *ZohoLayout) SectionApiNames() (r []string) {

	for _, v := range c.sections {
		r = append(r, v.GetString("api_name"))
	}
	return
}

// FieldsByApiName 此处为字段分组，例： Contacts_Sections_ApiName_Client_Information
func (c *ZohoLayout) FieldsByApiName(sectionApiName string) lib.TypeList {
	for k, v := range c.sections {
		if v.GetString("api_name") == sectionApiName {
			return lib.ToTypeList(c.sections[k].Get("fields"))
		}
	}
	return nil
}

func (c *ZohoLayout) FieldApiNamesByApiName(sectionApiName string) (r []string) {
	fields := c.FieldsByApiName(sectionApiName)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	return
}

func (c *ZohoLayout) FieldByApiName(apiName string) lib.TypeMap {
	var fields lib.TypeList
	a1 := c.FieldsByApiName(Deal_Sections_ApiName_Description_Information)
	a2 := c.FieldsByApiName(Deal_Sections_ApiName_ClientCaseInformation)
	a2_1 := c.FieldsByApiName(Deal_Sections_ApiName_ClientInformation)
	a3 := c.FieldsByApiName(Contact_Sections_ApiName_Client_Information)
	a4 := c.FieldsByApiName(Contact_Sections_ApiName_Address_Information)
	fields = append(fields, a1...)
	fields = append(fields, a2...)
	fields = append(fields, a2_1...)
	fields = append(fields, a3...)
	fields = append(fields, a4...)

	if c.module == Tasks {
		t1 := c.FieldsByApiName(Task_Sections_ApiName_TaskInformation)
		t2 := c.FieldsByApiName(Task_Sections_ApiName_DescriptionInformation)
		fields = append(fields, t1...)
		fields = append(fields, t2...)
	}

	for k, v := range fields {
		if v.GetString("api_name") == apiName {
			return fields[k]
		}
	}
	return nil

}

func (c *ZohoLayout) FieldInfoByApiName(sectionApiName string) (r lib.TypeMap) {
	fields := c.FieldsByApiName(sectionApiName)
	r = make(lib.TypeMap)
	for _, v := range fields {
		r.Set(v.GetString("api_name"), v.GetString("field_label"))
	}
	return
}

func (c *ZohoLayout) LeadFieldInfos() (r lib.TypeMap) {
	d1 := c.FieldInfoByApiName(Lead_Sections_ApiName_Lead_Information)
	d2 := c.FieldInfoByApiName(Lead_Sections_ApiName_Address_Information)
	d3 := c.FieldInfoByApiName(Lead_Sections_ApiName_Description_Information)
	return lib.TypeMapMerge(d1, d2, d3)
}

func (c *ZohoLayout) LeadApiNames() (r []string) {
	fields := c.FieldsByApiName(Lead_Sections_ApiName_Lead_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Lead_Sections_ApiName_Address_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Lead_Sections_ApiName_Description_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	var res []string
	for _, v := range r {
		//if v != "Locked__s" && v != "Change_Log_Time__s" { //  解决字段超50个的问题
		res = append(res, v)
		//}
	}

	return res
}

func (c *ZohoLayout) NoteApiNames() (r []string) {
	fields := c.FieldsByApiName(Notes_Sections_ApiName_NotesInformation)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	var res []string
	for _, v := range r {
		res = append(res, v)
	}
	return res
}

func (c *ZohoLayout) DealFieldInfos() (r lib.TypeMap) {
	d1 := c.FieldInfoByApiName(Deal_Sections_ApiName_ClientCaseInformation)
	d2 := c.FieldInfoByApiName(Deal_Sections_ApiName_Description_Information)
	d3 := c.FieldInfoByApiName(Deal_Sections_ApiName_ClientInformation)
	d4 := c.FieldInfoByApiName(Deal_Sections_ApiName_Presumptive_Information)
	return lib.TypeMapMerge(d1, d2, d3, d4)
}

func ApiNameByFieldLabel(fieldInfos lib.TypeMap, fieldLabel string) string {
	if fieldInfos != nil {
		for k, v := range fieldInfos {
			if lib.InterfaceToString(v) == fieldLabel {
				return k
			}
		}
	}
	return ""
}

func (c *ZohoLayout) PrintDealApiNames() {
	fields := c.FieldsByApiName(Deal_Sections_ApiName_ClientCaseInformation)
	var list lib.TypeList
	for _, v := range fields {
		if v.GetString("api_name") == Zoho_Collaborators {
			lib.DPrintln(Zoho_Collaborators+":", v)
			lib.DPrintln("\n\n\n\n")
		}
		//lib.DPrintln(v)
		list = append(list, map[string]interface{}{
			"api_name":  v.GetString("api_name"),
			"api_label": v.GetString("field_label"),
		})
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_Description_Information)
	for _, v := range fields {
		list = append(list, map[string]interface{}{
			"api_name":  v.GetString("api_name"),
			"api_label": v.GetString("field_label"),
		})
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_ClientInformation)
	for _, v := range fields {
		list = append(list, map[string]interface{}{
			"api_name":  v.GetString("api_name"),
			"api_label": v.GetString("field_label"),
		})
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_Presumptive_Information)
	for _, v := range fields {
		list = append(list, map[string]interface{}{
			"api_name":  v.GetString("api_name"),
			"api_label": v.GetString("field_label"),
		})
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_Service_Information)
	for _, v := range fields {
		list = append(list, map[string]interface{}{
			"api_name":  v.GetString("api_name"),
			"api_label": v.GetString("field_label"),
		})
	}

	lib.DPrintln(list)
}

func (c *ZohoLayout) DealApiNames() (r []string) {
	fields := c.FieldsByApiName(Deal_Sections_ApiName_ClientCaseInformation)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_Description_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_ClientInformation)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Deal_Sections_ApiName_Presumptive_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}

	//lib.DPrintln("sss:", len(r))
	//lib.DPrintln(r)

	var res []string
	for _, v := range r {
		if v != "Locked__s" && v != "Change_Log_Time__s" && v != "Tag" &&
			v != "Sales_Cycle_Duration" &&
			v != "Last_Activity_Time" && v != Zoho_Collaborators &&
			v != "Lead_Conversion_Time" && v != "Overall_Sales_Duration" &&
			v != "Benefits_Summary_Letter" &&
			v != "Modified_By" && v != "Pipeline" && v != "Tinnitus" && v != "Record_Status__s" && v != "VA_Records" && v != "Answer_to_Presumptive_Questions" &&
			v != "C_File_Submitted" &&
			v != "Referring_Person" && v != "VA_Healthcare_Registration" && v != "Private_Records" &&
			v != "Disability_Rating_List_Screenshot" && v != "Rating_Decision_Letters" && v != "DD214" &&
			v != "TDIU" && v != "STRs" && v != "Lead_CO" && v != "Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun" && v != "Primary_VS" &&
			v != Zoho_Statements && v != Zoho_SF_180_Submitted && v != Zoho_SF_180 && v != Zoho_FOIA { //  解决字段超50个的问题
			res = append(res, v)
		}
	}
	return res
}

// DealApiNames2  暂不能同步
func (c *ZohoLayout) DealApiNames2() (r []string) {
	return []string{"Referring_Person", "Lead_CO", "Modified_Time",
		"Answer_to_Presumptive_Questions",
		"Illness_Due_to_Toxic_Drinking_Water_at_Camp_Lejeun",
		"Service_Years",
		"Occupation_during_Service",
		"C_File_Submitted",
		"DD214",
		"Disability_Rating_List_Screenshot",
		"Rating_Decision_Letters",
		"STRs",
		"Benefits_Summary_Letter",
		"VA_Records",
		"Private_Records",
		"VA_Healthcare_Registration",
		"TDIU",
		"Tinnitus_Nexus",
		"Primary_VS",
		Zoho_Statements,
		Zoho_SF_180_Submitted,
		Zoho_SF_180,
		Zoho_FOIA}
}

func (c *ZohoLayout) ContactFieldInfos() (r lib.TypeMap) {
	d1 := c.FieldInfoByApiName(Contact_Sections_ApiName_Client_Information)
	d2 := c.FieldInfoByApiName(Contact_Sections_ApiName_Address_Information)
	return lib.TypeMapMerge(d1, d2)
}

func (c *ZohoLayout) ContactApiNames() (r []string) {

	fields := c.FieldsByApiName(Contact_Sections_ApiName_Address_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	fields = c.FieldsByApiName(Contact_Sections_ApiName_Client_Information)
	for _, v := range fields {
		r = append(r, v.GetString("api_name"))
	}
	return
}

//func ContactsLayout()
