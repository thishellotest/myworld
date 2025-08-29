package biz

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

type AsanaUsecase struct {
	CommonUsecase *CommonUsecase
	conf          *conf.Data
}

func NewAsanaUsecase(CommonUsecase *CommonUsecase, conf *conf.Data) *AsanaUsecase {
	return &AsanaUsecase{
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
}

func (c *AsanaUsecase) GetATask(gid string) (vo *AsanaGetATaskVo, isDel bool, err error) {
	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/tasks/%s", url, gid)
	response, err := lib.RequestDo(http.MethodGet, url, nil, map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	if err != nil {
		return nil, false, err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, false, err
	}
	a := string(bodyBytes)
	bodyMap := lib.ToTypeMapByString(a)

	if (response.StatusCode == 403 || response.StatusCode == 404) && bodyMap.GetString("errors") != "" {
		return nil, true, nil
	} else if response.StatusCode == 200 {
		r := lib.StringToT[*AsanaGetATaskVo](a)
		if r.IsOk() {
			return r.Unwrap(), false, nil
		} else {
			return nil, false, r.Err()
		}
	} else {
		return nil, false, errors.New(InterfaceToString(response.StatusCode) + " | " + a)
	}
}

func (c *AsanaUsecase) CreateATask(customFields lib.TypeMap, firstName string, lastName string, desc string) (*string, error) {
	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/tasks", url)

	destTypeMap := make(lib.TypeMap)
	destTypeMap.Set("data.name", ""+fmt.Sprintf("%s, %s", lastName, firstName))
	destTypeMap.Set("data.projects", []string{c.conf.Asana.GetProjectGid()})
	destTypeMap.Set("data.custom_fields", customFields)
	destTypeMap.Set("data.html_notes", "<body>Objective: "+desc+"</body>")

	r, _, er := lib.HTTPJsonWithHeaders(http.MethodPost, url, destTypeMap.ToBytes(), map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	return r, er
}

func (c *AsanaUsecase) PutATask(gid string, customFields lib.TypeMap, assigneeSectionGid string) (*string, error) {
	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/tasks/%s", url, gid)

	destTypeMap := make(lib.TypeMap)
	destTypeMap.Set("data.custom_fields", customFields)
	if assigneeSectionGid != "" {
		destTypeMap.Set("data.assignee_section", assigneeSectionGid)
	}
	// Create-only
	//destTypeMap.Set("data.projects", []string{c.conf.Asana.ProjectGidCp, c.conf.Asana.ProjectGid})
	r, _, er := lib.HTTPJsonWithHeaders(http.MethodPut, url, destTypeMap.ToBytes(), map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	return r, er
}

func (c *AsanaUsecase) AddAProjectToATask(taskGid string, projectGid string) {

	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/tasks/%s/addProject", url, taskGid)
	destTypeMap := make(lib.TypeMap)
	destTypeMap.Set("data.project", projectGid)

	lib.HTTPJsonWithHeaders(http.MethodPost, url, destTypeMap.ToBytes(), map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
}

func (c *AsanaUsecase) RemoveAProjectToATask(taskGid string, projectGid string) {

	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/tasks/%s/removeProject", url, taskGid)
	destTypeMap := make(lib.TypeMap)
	destTypeMap.Set("data.project", projectGid)

	lib.HTTPJsonWithHeaders(http.MethodPost, url, destTypeMap.ToBytes(), map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
}

func (c *AsanaUsecase) GetAUser(gid string) (lib.TypeMap, error) {
	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/users/%s", url, gid)
	a, _, err := lib.HTTPJsonWithHeaders(http.MethodGet, url, nil, map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	if err != nil {
		return nil, err
	}
	if a != nil {
		r := lib.StringToT[lib.TypeMap](*a)
		if r.IsOk() {
			return r.Unwrap(), nil
		} else {
			return nil, r.Err()
		}
	}
	return nil, err
}

func (c *AsanaUsecase) ListWebhooks() (lib.TypeMap, error) {
	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/webhooks?workspace=1205444246629009", url)
	//params := make(lib.TypeMap)
	//params.Set("workspace", "1205444246629009")
	a, _, err := lib.HTTPJsonWithHeaders(http.MethodGet, url, nil, map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	if err != nil {
		return nil, err
	}
	if a != nil {
		return lib.StringToTE[lib.TypeMap](*a, nil)
	}
	return nil, err
}

//
//func (c *AsanaUsecase) ChangeStagesToGettingStartedEmail(taskGid string) error {
//
//	taskInfo, _, err := c.GetATask(taskGid)
//	if err != nil {
//		return err
//	}
//	if taskInfo == nil {
//		return errors.New("taskInfo is nil.")
//	}
//	stagesFieldValue := taskInfo.Data.GetCustomFieldByName("Stages")
//	if stagesFieldValue == nil {
//		return errors.New("stagesFieldValue is nil")
//	}
//	if stagesFieldValue.DisplayValue != vbc_config.Stages_FeeScheduleandContract {
//		return errors.New("Stages is not Stages_FeeScheduleandContract.")
//	}
//
//	field := vbc_config.GetAsanaCustomFields()
//	fGid := field.GetByName("Stages").GetGid()
//	eGid := field.GetByName("Stages").GetEnumGidByName(vbc_config.Stages_GettingStartedEmail)
//	if fGid == "" || eGid == "" {
//		return errors.New("fGid or eGid is empty.")
//	}
//	customFields := make(lib.TypeMap)
//	customFields.Set(fGid, eGid)
//	sectionGid := vbc_config.GetAsanaSections().GetSectionGidByName(vbc_config.AsanaSection_GETTING_STARTED_EMAIL)
//	_, er := c.PutATask(taskGid, customFields, sectionGid)
//	return er
//}

func (c *AsanaUsecase) GetAllProjects(AsanaWorkspaceGid string) {

	url := c.conf.Asana.GetApiHost()
	url = fmt.Sprintf("%s/workspaces/%s/projects", url, AsanaWorkspaceGid)
	a, _, err := lib.HTTPJsonWithHeaders(http.MethodGet, url, nil, map[string]string{
		"authorization": "Bearer " + c.conf.Asana.Pat,
	})
	lib.DPrintln(a, err)
}

type AsanaGetATaskVo struct {
	Data AsanaGetATaskData `json:"data"`
}

func (c *AsanaGetATaskData) GetCustomFieldByName(name string) *CustomFields {

	for k, v := range c.CustomFields {
		if v.Name == name {
			t := c.CustomFields[k]
			return &t
		}
	}
	return nil
}

func AsanaTaskMapping() FieldMapping {
	return map[string]string{
		"First Name":                           "first_name",
		"Last Name":                            "last_name",
		"Current Rating":                       "current_rating",
		"Effective Current Rating":             "effective_current_rating",
		"Branch":                               "branch",
		"New Rating":                           "new_rating",
		"Email":                                "email",
		"Phone Number":                         "phone",
		"SSN":                                  "ssn", // custom_field_gid: 1205964025409311   google form: SSN
		"DOB":                                  "dob",
		"Stages":                               "stages",
		"Street Address":                       "address",
		"Address - Zip Code":                   "zip_code",
		"Address - State":                      "address_state",
		"Address - City":                       "city",
		"Retired":                              "retired",
		"Agent Orange Exposure":                "agent_orange",
		"Gulf War Illness":                     "gulf_war",
		"Burn Pits and Other Airborne Hazards": "burn_pits",
		"Illness Due to Toxic Drinking Water at Camp Lejeune": "illness_due",
		"\"Atomic Veterans\" and Radiation Exposure":          "atomic_veterans",
		"Amyotrophic Lateral Sclerosis (ALS)":                 "amyotrophic",
		"Disability Rating List Screenshot":                   "disability_rating",
		"Referrer":                                            "referrer",
		"Description":                                         "description",
		"Source":                                              "source",
	}
}

type FieldMapping map[string]string

func (c FieldMapping) FieldName(remoteName string) *string {
	if c != nil {
		if _, ok := c[remoteName]; ok {
			v := c[remoteName]
			return &v
		}
	}
	return nil
}

func (c *AsanaGetATaskVo) ToDataEntry() TypeDataEntry {
	typeDataEntry := make(TypeDataEntry)
	typeDataEntry["asana_task_gid"] = c.Data.Gid
	typeDataEntry["asana_projects"] = c.Data.AsanaProjects()
	typeDataEntry["assignee_gid"] = ""
	typeDataEntry["assignee_name"] = ""
	if c.Data.Assignee != nil {
		assignee := lib.TypeMap(c.Data.Assignee.(map[string]interface{}))
		if assignee.Get("gid") != nil {
			typeDataEntry["assignee_gid"] = lib.InterfaceToString(assignee.Get("gid"))
			typeDataEntry["assignee_name"] = lib.InterfaceToString(assignee.Get("name"))
		}
	}

	mappings := AsanaTaskMapping()
	for _, v := range c.Data.CustomFields {
		v := v
		if fieldName := mappings.FieldName(v.Name); fieldName != nil {
			typeDataEntry[*fieldName] = v.GetValueForDataEntry()
		}
	}
	typeDataEntry["task_name"] = c.Data.Name

	return typeDataEntry
}

type CreatedBy struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type EnumOptions struct {
	Gid          string `json:"gid"`
	Color        string `json:"color"`
	Enabled      bool   `json:"enabled"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type CustomFields struct {
	Gid             string        `json:"gid"`
	Enabled         bool          `json:"enabled"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	NumberValue     interface{}   `json:"number_value,omitempty"`
	Precision       int           `json:"precision,omitempty"`
	CreatedBy       CreatedBy     `json:"created_by"`
	DisplayValue    interface{}   `json:"display_value"`
	ResourceSubtype string        `json:"resource_subtype"`
	ResourceType    string        `json:"resource_type"`
	IsFormulaField  bool          `json:"is_formula_field"`
	IsValueReadOnly bool          `json:"is_value_read_only"`
	Type            string        `json:"type"`
	TextValue       string        `json:"text_value,omitempty"`
	EnumOptions     []EnumOptions `json:"enum_options,omitempty"`
	EnumValue       interface{}   `json:"enum_value,omitempty"`
	DateValue       interface{}   `json:"date_value,omitempty"`
}

func (c *CustomFields) GetValueForDataEntry() interface{} {
	if c.ResourceSubtype == "number" {
		return c.NumberValue
	} else if c.ResourceSubtype == "enum" {
		return c.DisplayValue
	} else if c.ResourceSubtype == "date" {
		// 2023-12-27T00:00:00.000Z
		if c.DisplayValue != nil {
			t := lib.InterfaceToString(c.DisplayValue)
			a, err := lib.TimeParse(t)
			if err != nil {
			} else {
				return a.Unix()
			}
		}
		return 0
	}
	return c.TextValue
}

type Followers struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type Project struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type Section struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type Memberships struct {
	Project Project `json:"project"`
	Section Section `json:"section"`
}
type Projects struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type Workspace struct {
	Gid          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}
type AsanaGetATaskData struct {
	Gid               string         `json:"gid"`
	ActualTimeMinutes interface{}    `json:"actual_time_minutes"`
	Assignee          interface{}    `json:"assignee"`
	AssigneeStatus    string         `json:"assignee_status"`
	Completed         bool           `json:"completed"`
	CompletedAt       interface{}    `json:"completed_at"`
	CreatedAt         time.Time      `json:"created_at"`
	CustomFields      []CustomFields `json:"custom_fields"`
	DueAt             interface{}    `json:"due_at"`
	DueOn             interface{}    `json:"due_on"`
	Followers         []Followers    `json:"followers"`
	Hearted           bool           `json:"hearted"`
	Hearts            []interface{}  `json:"hearts"`
	Liked             bool           `json:"liked"`
	Likes             []interface{}  `json:"likes"`
	Memberships       []Memberships  `json:"memberships"`
	ModifiedAt        time.Time      `json:"modified_at"`
	Name              string         `json:"name"`
	Notes             string         `json:"notes"`
	NumHearts         int            `json:"num_hearts"`
	NumLikes          int            `json:"num_likes"`
	Parent            interface{}    `json:"parent"`
	PermalinkURL      string         `json:"permalink_url"`
	Projects          []Projects     `json:"projects"`
	ResourceType      string         `json:"resource_type"`
	StartAt           interface{}    `json:"start_at"`
	StartOn           interface{}    `json:"start_on"`
	Tags              []interface{}  `json:"tags"`
	ResourceSubtype   string         `json:"resource_subtype"`
	Workspace         Workspace      `json:"workspace"`
}

func (c *AsanaGetATaskData) AsanaProjects() string {
	r := ""
	if len(c.Projects) > 0 {
		r = ","
		for _, v := range c.Projects {
			r += v.Gid + ","
		}
	}
	return r
}
