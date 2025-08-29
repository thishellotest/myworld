package config_vbc

import (
	"vbc/configs"
	"vbc/lib"
)

const (
	Asana_Field_Source = "Source"
)

type AsanaCustomFields lib.TypeMap

func (c AsanaCustomFields) CustomFields() lib.TypeList {
	return lib.ToTypeList(lib.TypeMap(c).Get("data.custom_fields"))
}

func (c AsanaCustomFields) GetByName(name string) AsanaCustomField {
	res := c.CustomFields()
	for k, v := range res {
		if v.GetString("name") == name {
			return AsanaCustomField(res[k])
		}
	}
	return nil
}

type AsanaCustomField lib.TypeMap

func (c AsanaCustomField) ToTypeMap() lib.TypeMap {
	return lib.TypeMap(c)
}

func (c AsanaCustomField) GetGid() string {
	return c.ToTypeMap().GetString("gid")
}

func (c AsanaCustomField) GetEnumOptions() lib.TypeList {
	t := lib.TypeMap(c)
	r := t.Get("enum_options")
	return lib.ToTypeList(r)
}

func (c AsanaCustomField) GetEnumGidByName(name string) string {
	res := c.GetEnumOptions()
	for _, v := range res {
		if v.GetString("name") == name {
			return v.GetString("gid")
		}
	}
	return ""
}

func GetAsanaCustomFields() AsanaCustomFields {
	return AsanaCustomFields(lib.ToTypeMapByString(getAsanaStr()))
}

func getAsanaStr() string {
	if configs.AppEnv() == configs.ENV_PROD {
		return asanaStrProd
	}
	return asanaStr
}

// const asanaStrProd = `{"data":{"gid":"1206343093612429","actual_time_minutes":null,"assignee":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"assignee_status":"inbox","completed":false,"completed_at":null,"created_at":"2024-01-15T06:20:22.368Z","custom_fields":[{"gid":"1206184385949353","enabled":true,"name":"First Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Rizalino","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Rizalino"},{"gid":"1206184460895324","enabled":true,"name":"Last Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Arrabis","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Arrabis"},{"gid":"1206479747533741","enabled":true,"name":"Email","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747533745","enabled":true,"name":"Phone Number","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747466575","enabled":true,"name":"Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553545","enabled":true,"name":"Effective Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206422409732583","enabled":true,"name":"Retired","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206422409732584","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206422409732585","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1205468334222326","enabled":true,"name":"Branch","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205468334222327","color":"blue","enabled":true,"name":"Navy","resource_type":"enum_option"},{"gid":"1205468334222328","color":"yellow-green","enabled":true,"name":"Army","resource_type":"enum_option"},{"gid":"1205468334222329","color":"blue-green","enabled":true,"name":"Air Force","resource_type":"enum_option"},{"gid":"1205468334222330","color":"red","enabled":true,"name":"Marine Corps","resource_type":"enum_option"},{"gid":"1205468334222331","color":"yellow-green","enabled":true,"name":"Army NG","resource_type":"enum_option"},{"gid":"1205960749182653","color":"orange","enabled":true,"name":"Coast Guard","resource_type":"enum_option"}]},{"gid":"1206479747466584","enabled":true,"name":"New Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747466593","enabled":true,"name":"ITF Expiration","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1206479747501015","enabled":true,"name":"Contact Form","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501016","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"}]},{"gid":"1206479747501021","enabled":true,"name":"C-File Submitted","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501022","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501023","color":"cool-gray","enabled":true,"name":"N/A","resource_type":"enum_option"},{"gid":"1206479747501024","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501031","enabled":true,"name":"DD214","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501032","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501033","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206478961495299","enabled":true,"name":"Disability Rating List Screenshot","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206478961495300","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206478961495301","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501039","enabled":true,"name":"Rating Decision Letters","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501040","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501041","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533715","enabled":true,"name":"STRs","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533716","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533717","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533724","enabled":true,"name":"TDIU","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533725","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533726","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533732","enabled":true,"name":"Item ID (auto generated)","description":"","precision":2,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553535","enabled":true,"name":"SSN","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553539","enabled":true,"name":"DOB","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1205512827064319","enabled":true,"name":"Street Address","description":"Street Address","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206401658215277","enabled":true,"name":"Address - City","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553554","enabled":true,"name":"Address - State","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747553555","color":"none","enabled":true,"name":"Alabama","resource_type":"enum_option"},{"gid":"1206479747553556","color":"none","enabled":true,"name":"Alaska","resource_type":"enum_option"},{"gid":"1206479747553557","color":"none","enabled":true,"name":"Arizona","resource_type":"enum_option"},{"gid":"1206479747553558","color":"none","enabled":true,"name":"Arkansas","resource_type":"enum_option"},{"gid":"1206479747553559","color":"none","enabled":true,"name":"American Samoa","resource_type":"enum_option"},{"gid":"1206479747553560","color":"none","enabled":true,"name":"California","resource_type":"enum_option"},{"gid":"1206479747553561","color":"none","enabled":true,"name":"Colorado","resource_type":"enum_option"},{"gid":"1206479747553562","color":"none","enabled":true,"name":"Connecticut","resource_type":"enum_option"},{"gid":"1206479747553563","color":"none","enabled":true,"name":"Delaware","resource_type":"enum_option"},{"gid":"1206479747553564","color":"none","enabled":true,"name":"District of Columbia","resource_type":"enum_option"},{"gid":"1206479747604120","color":"none","enabled":true,"name":"Florida","resource_type":"enum_option"},{"gid":"1206479747604121","color":"none","enabled":true,"name":"Georgia","resource_type":"enum_option"},{"gid":"1206479747604122","color":"none","enabled":true,"name":"Guam","resource_type":"enum_option"},{"gid":"1206479747604123","color":"none","enabled":true,"name":"Hawaii","resource_type":"enum_option"},{"gid":"1206479747604124","color":"none","enabled":true,"name":"Idaho","resource_type":"enum_option"},{"gid":"1206479747604125","color":"none","enabled":true,"name":"Illinois","resource_type":"enum_option"},{"gid":"1206479747604126","color":"none","enabled":true,"name":"Indiana","resource_type":"enum_option"},{"gid":"1206479747604127","color":"none","enabled":true,"name":"Iowa","resource_type":"enum_option"},{"gid":"1206479747604128","color":"none","enabled":true,"name":"Kansas","resource_type":"enum_option"},{"gid":"1206479747604129","color":"none","enabled":true,"name":"Kentucky","resource_type":"enum_option"},{"gid":"1206479747604130","color":"none","enabled":true,"name":"Louisiana","resource_type":"enum_option"},{"gid":"1206479747604131","color":"none","enabled":true,"name":"Maine","resource_type":"enum_option"},{"gid":"1206479747604132","color":"none","enabled":true,"name":"Maryland","resource_type":"enum_option"},{"gid":"1206479747604133","color":"none","enabled":true,"name":"Massachusetts","resource_type":"enum_option"},{"gid":"1206479747604134","color":"none","enabled":true,"name":"Michigan","resource_type":"enum_option"},{"gid":"1206479747604135","color":"none","enabled":true,"name":"Minnesota","resource_type":"enum_option"},{"gid":"1206479747604136","color":"none","enabled":true,"name":"Mississippi","resource_type":"enum_option"},{"gid":"1206479747604137","color":"none","enabled":true,"name":"Missouri","resource_type":"enum_option"},{"gid":"1206479747604138","color":"none","enabled":true,"name":"Montana","resource_type":"enum_option"},{"gid":"1206479747604139","color":"none","enabled":true,"name":"Nebraska","resource_type":"enum_option"},{"gid":"1206479747604140","color":"none","enabled":true,"name":"Nevada","resource_type":"enum_option"},{"gid":"1206479747604141","color":"none","enabled":true,"name":"New Hampshire","resource_type":"enum_option"},{"gid":"1206479747604142","color":"none","enabled":true,"name":"New Jersey","resource_type":"enum_option"},{"gid":"1206479747604143","color":"none","enabled":true,"name":"New Mexico","resource_type":"enum_option"},{"gid":"1206479747604144","color":"none","enabled":true,"name":"New York","resource_type":"enum_option"},{"gid":"1206479747604145","color":"none","enabled":true,"name":"North Carolina","resource_type":"enum_option"},{"gid":"1206479747604146","color":"none","enabled":true,"name":"North Dakota","resource_type":"enum_option"},{"gid":"1206479747604147","color":"none","enabled":true,"name":"Northern Mariana Islands","resource_type":"enum_option"},{"gid":"1206479747604148","color":"none","enabled":true,"name":"Ohio","resource_type":"enum_option"},{"gid":"1206479747604149","color":"none","enabled":true,"name":"Oklahoma","resource_type":"enum_option"},{"gid":"1206479747604150","color":"none","enabled":true,"name":"Oregon","resource_type":"enum_option"},{"gid":"1206479747604151","color":"none","enabled":true,"name":"Pennsylvania","resource_type":"enum_option"},{"gid":"1206479747639138","color":"none","enabled":true,"name":"Puerto Rico","resource_type":"enum_option"},{"gid":"1206479747639139","color":"none","enabled":true,"name":"Rhode Island","resource_type":"enum_option"},{"gid":"1206479747639140","color":"none","enabled":true,"name":"South Carolina","resource_type":"enum_option"},{"gid":"1206479747639141","color":"none","enabled":true,"name":"South Dakota","resource_type":"enum_option"},{"gid":"1206479747639142","color":"none","enabled":true,"name":"Tennessee","resource_type":"enum_option"},{"gid":"1206479747639143","color":"none","enabled":true,"name":"Texas","resource_type":"enum_option"},{"gid":"1206479747639144","color":"none","enabled":true,"name":"Trust Territories","resource_type":"enum_option"},{"gid":"1206479747639145","color":"none","enabled":true,"name":"Utah","resource_type":"enum_option"},{"gid":"1206479747639146","color":"none","enabled":true,"name":"Vermont","resource_type":"enum_option"},{"gid":"1206479747639147","color":"none","enabled":true,"name":"Virginia","resource_type":"enum_option"},{"gid":"1206479747639148","color":"none","enabled":true,"name":"Virgin Islands","resource_type":"enum_option"},{"gid":"1206479747639149","color":"none","enabled":true,"name":"Washington","resource_type":"enum_option"},{"gid":"1206479747639150","color":"none","enabled":true,"name":"West Virginia","resource_type":"enum_option"},{"gid":"1206479747639151","color":"none","enabled":true,"name":"Wisconsin","resource_type":"enum_option"},{"gid":"1206479747639152","color":"none","enabled":true,"name":"Wyoming","resource_type":"enum_option"},{"gid":"1206479747639153","color":"none","enabled":true,"name":"Outside United States","resource_type":"enum_option"}]},{"gid":"1206401684235934","enabled":true,"name":"Address - Zip Code","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1205769148648333","enabled":true,"name":"Stages","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206235855089869","color":"purple","enabled":true,"name":"Incoming Request","resource_type":"enum_option"},{"gid":"1206235855089870","color":"indigo","enabled":true,"name":"Fee Schedule and Contract","resource_type":"enum_option"},{"gid":"1206235855089871","color":"red","enabled":true,"name":"Getting Started Email","resource_type":"enum_option"},{"gid":"1205769148648382","color":"aqua","enabled":true,"name":"Update Client Intake Info","resource_type":"enum_option"},{"gid":"1206089580019064","color":"blue-green","enabled":true,"name":"Awaiting C-File","resource_type":"enum_option"},{"gid":"1205769148648383","color":"blue-green","enabled":true,"name":"Record Review","resource_type":"enum_option"},{"gid":"1205769148648385","color":"yellow-green","enabled":true,"name":"Schedule Call","resource_type":"enum_option"},{"gid":"1205769148648386","color":"yellow","enabled":true,"name":"Statement Drafts","resource_type":"enum_option"},{"gid":"1205769148648387","color":"yellow-orange","enabled":true,"name":"Statements Finalized","resource_type":"enum_option"},{"gid":"1205769148648388","color":"pink","enabled":true,"name":"Current Treatment","resource_type":"enum_option"},{"gid":"1205773364138828","color":"hot-pink","enabled":true,"name":"Mini-DBQs","resource_type":"enum_option"},{"gid":"1205773364138829","color":"magenta","enabled":true,"name":"Nexus Letters","resource_type":"enum_option"},{"gid":"1205773364138830","color":"purple","enabled":true,"name":"Medical Team","resource_type":"enum_option"},{"gid":"1205773364138831","color":"indigo","enabled":true,"name":"File Claims (New or Supplemental)","resource_type":"enum_option"},{"gid":"1205773364138832","color":"blue","enabled":true,"name":"Verify Evidence Received","resource_type":"enum_option"},{"gid":"1205773364138833","color":"red","enabled":true,"name":"Awaiting Decision","resource_type":"enum_option"},{"gid":"1206122140407998","color":"green","enabled":true,"name":"Awaiting Payment (Update Rating Board, Send Invoice and Clothing Email)","resource_type":"enum_option"},{"gid":"1206235855089872","color":"orange","enabled":true,"name":"Completed","resource_type":"enum_option"}]},{"gid":"1205769148648348","enabled":true,"name":"Status","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Overdue","resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205769148648349","color":"red","enabled":true,"name":"Overdue","resource_type":"enum_option"}],"enum_value":{"color":"red","enabled":true,"gid":"1205769148648349","name":"Overdue","resource_type":"enum_option"}},{"gid":"1205964025396864","enabled":true,"name":"Agent Orange Exposure","description":"","created_by":{"gid":"5609879243379","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205964025396865","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401859348576","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206398481016988","enabled":true,"name":"Burn Pits and Other Airborne Hazards","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206398481016989","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206398481016990","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858650700","enabled":true,"name":"Gulf War Illness","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858650701","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858650702","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401856356610","enabled":true,"name":"Illness Due to Toxic Drinking Water at Camp Lejeune","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401856356611","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401856356612","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858770487","enabled":true,"name":"\"Atomic Veterans\" and Radiation Exposure","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858770488","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858770489","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858924307","enabled":true,"name":"Amyotrophic Lateral Sclerosis (ALS)","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858924308","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858924309","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]}],"due_at":null,"due_on":"2024-01-26","followers":[{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"}],"hearted":false,"hearts":[],"liked":false,"likes":[],"memberships":[{"project":{"gid":"1206472580135542","name":"VBC Client List - Automated","resource_type":"project"},"section":{"gid":"1206472580135543","name":"INCOMING REQUEST","resource_type":"section"}}],"modified_at":"2024-01-30T20:16:09.984Z","name":"Arrabis, Rizalino","notes":"Service:\n\n\nCurrent:\n\n\nNew:\n100 - tdiu","num_hearts":0,"num_likes":0,"parent":null,"permalink_url":"https://app.asana.com/0/1206472580135542/1206343093612429","projects":[{"gid":"1206472580135542","name":"VBC Client List - Automated","resource_type":"project"}],"resource_type":"task","start_at":null,"start_on":null,"tags":[],"resource_subtype":"default_task","workspace":{"gid":"1205444246629009","name":"Veteran Benefits Center LLC","resource_type":"workspace"}}}`
// const asanaStrProd = `{"data":{"gid":"1206343093612429","actual_time_minutes":null,"assignee":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"assignee_status":"inbox","completed":false,"completed_at":null,"created_at":"2024-01-15T06:20:22.368Z","custom_fields":[{"gid":"1206629760929727","enabled":true,"name":"Source","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Manual","resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206629760929729","color":"blue","enabled":true,"name":"Manual","resource_type":"enum_option"},{"gid":"1206629760929728","color":"blue-green","enabled":true,"name":"Website","resource_type":"enum_option"}],"enum_value":{"color":"blue","enabled":true,"gid":"1206629760929729","name":"Manual","resource_type":"enum_option"}},{"gid":"1206184385949353","enabled":true,"name":"First Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Rizalino","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Rizalino"},{"gid":"1206184460895324","enabled":true,"name":"Last Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Arrabis","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Arrabis"},{"gid":"1206539090084664","enabled":true,"name":"Referrer","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":"Jorge Balares","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Jorge Balares"},{"gid":"1206479747533741","enabled":true,"name":"Email","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747533745","enabled":true,"name":"Phone Number","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747466575","enabled":true,"name":"Current Rating","description":"","number_value":50,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":"50","resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553545","enabled":true,"name":"Effective Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206422409732583","enabled":true,"name":"Retired","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206422409732584","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206422409732585","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1205468334222326","enabled":true,"name":"Branch","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205468334222327","color":"blue","enabled":true,"name":"Navy","resource_type":"enum_option"},{"gid":"1205468334222328","color":"yellow-green","enabled":true,"name":"Army","resource_type":"enum_option"},{"gid":"1205468334222329","color":"blue-green","enabled":true,"name":"Air Force","resource_type":"enum_option"},{"gid":"1205468334222330","color":"red","enabled":true,"name":"Marine Corps","resource_type":"enum_option"},{"gid":"1205468334222331","color":"yellow-green","enabled":true,"name":"Army NG","resource_type":"enum_option"},{"gid":"1205960749182653","color":"orange","enabled":true,"name":"Coast Guard","resource_type":"enum_option"}]},{"gid":"1206479747466584","enabled":true,"name":"New Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747466593","enabled":true,"name":"ITF Expiration","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1206479747501015","enabled":true,"name":"Contact Form","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501016","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206640523877079","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501021","enabled":true,"name":"C-File Submitted","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501022","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501023","color":"cool-gray","enabled":true,"name":"N/A","resource_type":"enum_option"},{"gid":"1206479747501024","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501031","enabled":true,"name":"DD214","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501032","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501033","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206478961495299","enabled":true,"name":"Disability Rating List Screenshot","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206478961495300","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206478961495301","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501039","enabled":true,"name":"Rating Decision Letters","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501040","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501041","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533715","enabled":true,"name":"STRs","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533716","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533717","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533724","enabled":true,"name":"TDIU","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533725","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533726","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533732","enabled":true,"name":"Item ID (auto generated)","description":"","precision":2,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553535","enabled":true,"name":"SSN","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553539","enabled":true,"name":"DOB","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1205512827064319","enabled":true,"name":"Street Address","description":"Street Address","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206401658215277","enabled":true,"name":"Address - City","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553554","enabled":true,"name":"Address - State","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747553555","color":"none","enabled":true,"name":"Alabama","resource_type":"enum_option"},{"gid":"1206479747553556","color":"none","enabled":true,"name":"Alaska","resource_type":"enum_option"},{"gid":"1206479747553557","color":"none","enabled":true,"name":"Arizona","resource_type":"enum_option"},{"gid":"1206479747553558","color":"none","enabled":true,"name":"Arkansas","resource_type":"enum_option"},{"gid":"1206479747553559","color":"none","enabled":true,"name":"American Samoa","resource_type":"enum_option"},{"gid":"1206479747553560","color":"none","enabled":true,"name":"California","resource_type":"enum_option"},{"gid":"1206479747553561","color":"none","enabled":true,"name":"Colorado","resource_type":"enum_option"},{"gid":"1206479747553562","color":"none","enabled":true,"name":"Connecticut","resource_type":"enum_option"},{"gid":"1206479747553563","color":"none","enabled":true,"name":"Delaware","resource_type":"enum_option"},{"gid":"1206479747553564","color":"none","enabled":true,"name":"District of Columbia","resource_type":"enum_option"},{"gid":"1206479747604120","color":"none","enabled":true,"name":"Florida","resource_type":"enum_option"},{"gid":"1206479747604121","color":"none","enabled":true,"name":"Georgia","resource_type":"enum_option"},{"gid":"1206479747604122","color":"none","enabled":true,"name":"Guam","resource_type":"enum_option"},{"gid":"1206479747604123","color":"none","enabled":true,"name":"Hawaii","resource_type":"enum_option"},{"gid":"1206479747604124","color":"none","enabled":true,"name":"Idaho","resource_type":"enum_option"},{"gid":"1206479747604125","color":"none","enabled":true,"name":"Illinois","resource_type":"enum_option"},{"gid":"1206479747604126","color":"none","enabled":true,"name":"Indiana","resource_type":"enum_option"},{"gid":"1206479747604127","color":"none","enabled":true,"name":"Iowa","resource_type":"enum_option"},{"gid":"1206479747604128","color":"none","enabled":true,"name":"Kansas","resource_type":"enum_option"},{"gid":"1206479747604129","color":"none","enabled":true,"name":"Kentucky","resource_type":"enum_option"},{"gid":"1206479747604130","color":"none","enabled":true,"name":"Louisiana","resource_type":"enum_option"},{"gid":"1206479747604131","color":"none","enabled":true,"name":"Maine","resource_type":"enum_option"},{"gid":"1206479747604132","color":"none","enabled":true,"name":"Maryland","resource_type":"enum_option"},{"gid":"1206479747604133","color":"none","enabled":true,"name":"Massachusetts","resource_type":"enum_option"},{"gid":"1206479747604134","color":"none","enabled":true,"name":"Michigan","resource_type":"enum_option"},{"gid":"1206479747604135","color":"none","enabled":true,"name":"Minnesota","resource_type":"enum_option"},{"gid":"1206479747604136","color":"none","enabled":true,"name":"Mississippi","resource_type":"enum_option"},{"gid":"1206479747604137","color":"none","enabled":true,"name":"Missouri","resource_type":"enum_option"},{"gid":"1206479747604138","color":"none","enabled":true,"name":"Montana","resource_type":"enum_option"},{"gid":"1206479747604139","color":"none","enabled":true,"name":"Nebraska","resource_type":"enum_option"},{"gid":"1206479747604140","color":"none","enabled":true,"name":"Nevada","resource_type":"enum_option"},{"gid":"1206479747604141","color":"none","enabled":true,"name":"New Hampshire","resource_type":"enum_option"},{"gid":"1206479747604142","color":"none","enabled":true,"name":"New Jersey","resource_type":"enum_option"},{"gid":"1206479747604143","color":"none","enabled":true,"name":"New Mexico","resource_type":"enum_option"},{"gid":"1206479747604144","color":"none","enabled":true,"name":"New York","resource_type":"enum_option"},{"gid":"1206479747604145","color":"none","enabled":true,"name":"North Carolina","resource_type":"enum_option"},{"gid":"1206479747604146","color":"none","enabled":true,"name":"North Dakota","resource_type":"enum_option"},{"gid":"1206479747604147","color":"none","enabled":true,"name":"Northern Mariana Islands","resource_type":"enum_option"},{"gid":"1206479747604148","color":"none","enabled":true,"name":"Ohio","resource_type":"enum_option"},{"gid":"1206479747604149","color":"none","enabled":true,"name":"Oklahoma","resource_type":"enum_option"},{"gid":"1206479747604150","color":"none","enabled":true,"name":"Oregon","resource_type":"enum_option"},{"gid":"1206479747604151","color":"none","enabled":true,"name":"Pennsylvania","resource_type":"enum_option"},{"gid":"1206479747639138","color":"none","enabled":true,"name":"Puerto Rico","resource_type":"enum_option"},{"gid":"1206479747639139","color":"none","enabled":true,"name":"Rhode Island","resource_type":"enum_option"},{"gid":"1206479747639140","color":"none","enabled":true,"name":"South Carolina","resource_type":"enum_option"},{"gid":"1206479747639141","color":"none","enabled":true,"name":"South Dakota","resource_type":"enum_option"},{"gid":"1206479747639142","color":"none","enabled":true,"name":"Tennessee","resource_type":"enum_option"},{"gid":"1206479747639143","color":"none","enabled":true,"name":"Texas","resource_type":"enum_option"},{"gid":"1206479747639144","color":"none","enabled":true,"name":"Trust Territories","resource_type":"enum_option"},{"gid":"1206479747639145","color":"none","enabled":true,"name":"Utah","resource_type":"enum_option"},{"gid":"1206479747639146","color":"none","enabled":true,"name":"Vermont","resource_type":"enum_option"},{"gid":"1206479747639147","color":"none","enabled":true,"name":"Virginia","resource_type":"enum_option"},{"gid":"1206479747639148","color":"none","enabled":true,"name":"Virgin Islands","resource_type":"enum_option"},{"gid":"1206479747639149","color":"none","enabled":true,"name":"Washington","resource_type":"enum_option"},{"gid":"1206479747639150","color":"none","enabled":true,"name":"West Virginia","resource_type":"enum_option"},{"gid":"1206479747639151","color":"none","enabled":true,"name":"Wisconsin","resource_type":"enum_option"},{"gid":"1206479747639152","color":"none","enabled":true,"name":"Wyoming","resource_type":"enum_option"},{"gid":"1206479747639153","color":"none","enabled":true,"name":"Outside United States","resource_type":"enum_option"}]},{"gid":"1206401684235934","enabled":true,"name":"Address - Zip Code","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1205769148648333","enabled":true,"name":"Stages","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206235855089869","color":"purple","enabled":true,"name":"Incoming Request","resource_type":"enum_option"},{"gid":"1206235855089870","color":"indigo","enabled":true,"name":"Fee Schedule and Contract","resource_type":"enum_option"},{"gid":"1206235855089871","color":"red","enabled":true,"name":"Getting Started Email","resource_type":"enum_option"},{"gid":"1205769148648382","color":"aqua","enabled":true,"name":"Update Client Intake Info","resource_type":"enum_option"},{"gid":"1206089580019064","color":"blue-green","enabled":true,"name":"Awaiting C-File","resource_type":"enum_option"},{"gid":"1205769148648383","color":"blue-green","enabled":true,"name":"Record Review","resource_type":"enum_option"},{"gid":"1205769148648385","color":"yellow-green","enabled":true,"name":"Schedule Call","resource_type":"enum_option"},{"gid":"1205769148648386","color":"yellow","enabled":true,"name":"Statement Drafts","resource_type":"enum_option"},{"gid":"1205769148648387","color":"yellow-orange","enabled":true,"name":"Statements Finalized","resource_type":"enum_option"},{"gid":"1205769148648388","color":"pink","enabled":true,"name":"Current Treatment","resource_type":"enum_option"},{"gid":"1205773364138828","color":"hot-pink","enabled":true,"name":"Mini-DBQs","resource_type":"enum_option"},{"gid":"1205773364138829","color":"magenta","enabled":true,"name":"Nexus Letters","resource_type":"enum_option"},{"gid":"1205773364138830","color":"purple","enabled":true,"name":"Medical Team","resource_type":"enum_option"},{"gid":"1205773364138831","color":"indigo","enabled":true,"name":"File Claims (New or Supplemental)","resource_type":"enum_option"},{"gid":"1205773364138832","color":"blue","enabled":true,"name":"Verify Evidence Received","resource_type":"enum_option"},{"gid":"1205773364138833","color":"red","enabled":true,"name":"Awaiting Decision","resource_type":"enum_option"},{"gid":"1206122140407998","color":"green","enabled":true,"name":"Awaiting Payment (Update Rating Board, Send Invoice and Clothing Email)","resource_type":"enum_option"},{"gid":"1206235855089872","color":"orange","enabled":true,"name":"Completed","resource_type":"enum_option"}]},{"gid":"1205964025396864","enabled":true,"name":"Agent Orange Exposure","description":"","created_by":{"gid":"5609879243379","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205964025396865","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401859348576","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206398481016988","enabled":true,"name":"Burn Pits and Other Airborne Hazards","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206398481016989","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206398481016990","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858650700","enabled":true,"name":"Gulf War Illness","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858650701","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858650702","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401856356610","enabled":true,"name":"Illness Due to Toxic Drinking Water at Camp Lejeune","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401856356611","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401856356612","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858770487","enabled":true,"name":"\"Atomic Veterans\" and Radiation Exposure","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858770488","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858770489","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858924307","enabled":true,"name":"Amyotrophic Lateral Sclerosis (ALS)","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858924308","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858924309","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]}],"due_at":null,"due_on":"2024-03-29","followers":[{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"}],"hearted":false,"hearts":[],"liked":false,"likes":[],"memberships":[{"project":{"gid":"1206472580135542","name":"VBC Client List - Automated","resource_type":"project"},"section":{"gid":"1206472580135543","name":"INCOMING REQUEST","resource_type":"section"}}],"modified_at":"2024-02-22T18:22:01.998Z","name":"Arrabis, Rizalino","notes":"Service:\n\n\nCurrent:\n\n\nNew:\n100 - tdiu","num_hearts":0,"num_likes":0,"parent":null,"permalink_url":"https://app.asana.com/0/1206472580135542/1206343093612429","projects":[{"gid":"1206472580135542","name":"VBC Client List - Automated","resource_type":"project"}],"resource_type":"task","start_at":null,"start_on":null,"tags":[],"resource_subtype":"default_task","workspace":{"gid":"1205444246629009","name":"Veteran Benefits Center LLC","resource_type":"workspace"}}}`
const asanaStrProd = `{"data":{"gid":"1206343093612429","actual_time_minutes":null,"assignee":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"assignee_status":"inbox","completed":false,"completed_at":null,"created_at":"2024-01-15T06:20:22.368Z","custom_fields":[{"gid":"1206629760929727","enabled":true,"name":"Source","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Manual","resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206629760929729","color":"blue","enabled":true,"name":"Manual","resource_type":"enum_option"},{"gid":"1206629760929728","color":"blue-green","enabled":true,"name":"Website","resource_type":"enum_option"}],"enum_value":{"color":"blue","enabled":true,"gid":"1206629760929729","name":"Manual","resource_type":"enum_option"}},{"gid":"1206184385949353","enabled":true,"name":"First Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Rizalino","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Rizalino"},{"gid":"1206184460895324","enabled":true,"name":"Last Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Arrabis","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Arrabis"},{"gid":"1206539090084664","enabled":true,"name":"Referrer","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":"Jorge Balares","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Jorge Balares"},{"gid":"1206479747533741","enabled":true,"name":"Email","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747533745","enabled":true,"name":"Phone Number","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747466575","enabled":true,"name":"Current Rating","description":"","number_value":50,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":"50","resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553545","enabled":true,"name":"Effective Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206422409732583","enabled":true,"name":"Retired","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206422409732584","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206422409732585","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1205468334222326","enabled":true,"name":"Branch","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205468334222327","color":"blue","enabled":true,"name":"Navy","resource_type":"enum_option"},{"gid":"1205468334222328","color":"yellow-green","enabled":true,"name":"Army","resource_type":"enum_option"},{"gid":"1205468334222329","color":"blue-green","enabled":true,"name":"Air Force","resource_type":"enum_option"},{"gid":"1205468334222330","color":"red","enabled":true,"name":"Marine Corps","resource_type":"enum_option"},{"gid":"1205468334222331","color":"yellow-green","enabled":true,"name":"Army NG","resource_type":"enum_option"},{"gid":"1205960749182653","color":"orange","enabled":true,"name":"Coast Guard","resource_type":"enum_option"}]},{"gid":"1206479747466584","enabled":true,"name":"New Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747466593","enabled":true,"name":"ITF Expiration","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1206479747501015","enabled":true,"name":"Contact Form","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501016","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206640523877079","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501021","enabled":true,"name":"C-File Submitted","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501022","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501023","color":"cool-gray","enabled":true,"name":"N/A","resource_type":"enum_option"},{"gid":"1206479747501024","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501031","enabled":true,"name":"DD214","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501032","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501033","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206478961495299","enabled":true,"name":"Disability Rating List Screenshot","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206478961495300","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206478961495301","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501039","enabled":true,"name":"Rating Decision Letters","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501040","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501041","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533715","enabled":true,"name":"STRs","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533716","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533717","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533724","enabled":true,"name":"TDIU","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533725","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533726","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533732","enabled":true,"name":"Item ID (auto generated)","description":"","precision":2,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553535","enabled":true,"name":"SSN","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553539","enabled":true,"name":"DOB","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1205512827064319","enabled":true,"name":"Street Address","description":"Street Address","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206401658215277","enabled":true,"name":"Address - City","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553554","enabled":true,"name":"Address - State","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747553555","color":"none","enabled":true,"name":"Alabama","resource_type":"enum_option"},{"gid":"1206479747553556","color":"none","enabled":true,"name":"Alaska","resource_type":"enum_option"},{"gid":"1206479747553557","color":"none","enabled":true,"name":"Arizona","resource_type":"enum_option"},{"gid":"1206479747553558","color":"none","enabled":true,"name":"Arkansas","resource_type":"enum_option"},{"gid":"1206479747553559","color":"none","enabled":true,"name":"American Samoa","resource_type":"enum_option"},{"gid":"1206479747553560","color":"none","enabled":true,"name":"California","resource_type":"enum_option"},{"gid":"1206479747553561","color":"none","enabled":true,"name":"Colorado","resource_type":"enum_option"},{"gid":"1206479747553562","color":"none","enabled":true,"name":"Connecticut","resource_type":"enum_option"},{"gid":"1206479747553563","color":"none","enabled":true,"name":"Delaware","resource_type":"enum_option"},{"gid":"1206479747553564","color":"none","enabled":true,"name":"District of Columbia","resource_type":"enum_option"},{"gid":"1206479747604120","color":"none","enabled":true,"name":"Florida","resource_type":"enum_option"},{"gid":"1206479747604121","color":"none","enabled":true,"name":"Georgia","resource_type":"enum_option"},{"gid":"1206479747604122","color":"none","enabled":true,"name":"Guam","resource_type":"enum_option"},{"gid":"1206479747604123","color":"none","enabled":true,"name":"Hawaii","resource_type":"enum_option"},{"gid":"1206479747604124","color":"none","enabled":true,"name":"Idaho","resource_type":"enum_option"},{"gid":"1206479747604125","color":"none","enabled":true,"name":"Illinois","resource_type":"enum_option"},{"gid":"1206479747604126","color":"none","enabled":true,"name":"Indiana","resource_type":"enum_option"},{"gid":"1206479747604127","color":"none","enabled":true,"name":"Iowa","resource_type":"enum_option"},{"gid":"1206479747604128","color":"none","enabled":true,"name":"Kansas","resource_type":"enum_option"},{"gid":"1206479747604129","color":"none","enabled":true,"name":"Kentucky","resource_type":"enum_option"},{"gid":"1206479747604130","color":"none","enabled":true,"name":"Louisiana","resource_type":"enum_option"},{"gid":"1206479747604131","color":"none","enabled":true,"name":"Maine","resource_type":"enum_option"},{"gid":"1206479747604132","color":"none","enabled":true,"name":"Maryland","resource_type":"enum_option"},{"gid":"1206479747604133","color":"none","enabled":true,"name":"Massachusetts","resource_type":"enum_option"},{"gid":"1206479747604134","color":"none","enabled":true,"name":"Michigan","resource_type":"enum_option"},{"gid":"1206479747604135","color":"none","enabled":true,"name":"Minnesota","resource_type":"enum_option"},{"gid":"1206479747604136","color":"none","enabled":true,"name":"Mississippi","resource_type":"enum_option"},{"gid":"1206479747604137","color":"none","enabled":true,"name":"Missouri","resource_type":"enum_option"},{"gid":"1206479747604138","color":"none","enabled":true,"name":"Montana","resource_type":"enum_option"},{"gid":"1206479747604139","color":"none","enabled":true,"name":"Nebraska","resource_type":"enum_option"},{"gid":"1206479747604140","color":"none","enabled":true,"name":"Nevada","resource_type":"enum_option"},{"gid":"1206479747604141","color":"none","enabled":true,"name":"New Hampshire","resource_type":"enum_option"},{"gid":"1206479747604142","color":"none","enabled":true,"name":"New Jersey","resource_type":"enum_option"},{"gid":"1206479747604143","color":"none","enabled":true,"name":"New Mexico","resource_type":"enum_option"},{"gid":"1206479747604144","color":"none","enabled":true,"name":"New York","resource_type":"enum_option"},{"gid":"1206479747604145","color":"none","enabled":true,"name":"North Carolina","resource_type":"enum_option"},{"gid":"1206479747604146","color":"none","enabled":true,"name":"North Dakota","resource_type":"enum_option"},{"gid":"1206479747604147","color":"none","enabled":true,"name":"Northern Mariana Islands","resource_type":"enum_option"},{"gid":"1206479747604148","color":"none","enabled":true,"name":"Ohio","resource_type":"enum_option"},{"gid":"1206479747604149","color":"none","enabled":true,"name":"Oklahoma","resource_type":"enum_option"},{"gid":"1206479747604150","color":"none","enabled":true,"name":"Oregon","resource_type":"enum_option"},{"gid":"1206479747604151","color":"none","enabled":true,"name":"Pennsylvania","resource_type":"enum_option"},{"gid":"1206479747639138","color":"none","enabled":true,"name":"Puerto Rico","resource_type":"enum_option"},{"gid":"1206479747639139","color":"none","enabled":true,"name":"Rhode Island","resource_type":"enum_option"},{"gid":"1206479747639140","color":"none","enabled":true,"name":"South Carolina","resource_type":"enum_option"},{"gid":"1206479747639141","color":"none","enabled":true,"name":"South Dakota","resource_type":"enum_option"},{"gid":"1206479747639142","color":"none","enabled":true,"name":"Tennessee","resource_type":"enum_option"},{"gid":"1206479747639143","color":"none","enabled":true,"name":"Texas","resource_type":"enum_option"},{"gid":"1206479747639144","color":"none","enabled":true,"name":"Trust Territories","resource_type":"enum_option"},{"gid":"1206479747639145","color":"none","enabled":true,"name":"Utah","resource_type":"enum_option"},{"gid":"1206479747639146","color":"none","enabled":true,"name":"Vermont","resource_type":"enum_option"},{"gid":"1206479747639147","color":"none","enabled":true,"name":"Virginia","resource_type":"enum_option"},{"gid":"1206479747639148","color":"none","enabled":true,"name":"Virgin Islands","resource_type":"enum_option"},{"gid":"1206479747639149","color":"none","enabled":true,"name":"Washington","resource_type":"enum_option"},{"gid":"1206479747639150","color":"none","enabled":true,"name":"West Virginia","resource_type":"enum_option"},{"gid":"1206479747639151","color":"none","enabled":true,"name":"Wisconsin","resource_type":"enum_option"},{"gid":"1206479747639152","color":"none","enabled":true,"name":"Wyoming","resource_type":"enum_option"},{"gid":"1206479747639153","color":"none","enabled":true,"name":"Outside United States","resource_type":"enum_option"}]},{"gid":"1206401684235934","enabled":true,"name":"Address - Zip Code","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1205769148648333","enabled":true,"name":"Stages","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206235855089869","color":"purple","enabled":true,"name":"Incoming Request","resource_type":"enum_option"},{"gid":"1206235855089870","color":"indigo","enabled":true,"name":"Fee Schedule and Contract","resource_type":"enum_option"},{"gid":"1206235855089871","color":"red","enabled":true,"name":"Getting Started Email","resource_type":"enum_option"},{"gid":"1205769148648382","color":"aqua","enabled":true,"name":"Awaiting Client Records","resource_type":"enum_option"},{"gid":"1206089580019064","color":"blue-green","enabled":true,"name":"Awaiting C-File","resource_type":"enum_option"},{"gid":"1205769148648383","color":"blue-green","enabled":true,"name":"Record Review","resource_type":"enum_option"},{"gid":"1205769148648385","color":"yellow-green","enabled":true,"name":"Schedule Call","resource_type":"enum_option"},{"gid":"1205769148648386","color":"yellow","enabled":true,"name":"Statement Drafts","resource_type":"enum_option"},{"gid":"1205769148648387","color":"yellow-orange","enabled":true,"name":"Statements Finalized","resource_type":"enum_option"},{"gid":"1205769148648388","color":"pink","enabled":true,"name":"Current Treatment","resource_type":"enum_option"},{"gid":"1205773364138828","color":"hot-pink","enabled":true,"name":"Mini-DBQs","resource_type":"enum_option"},{"gid":"1205773364138829","color":"magenta","enabled":true,"name":"Nexus Letters","resource_type":"enum_option"},{"gid":"1205773364138830","color":"purple","enabled":true,"name":"Medical Team","resource_type":"enum_option"},{"gid":"1205773364138831","color":"indigo","enabled":true,"name":"File Claims (New or Supplemental)","resource_type":"enum_option"},{"gid":"1205773364138832","color":"blue","enabled":true,"name":"Verify Evidence Received","resource_type":"enum_option"},{"gid":"1205773364138833","color":"red","enabled":true,"name":"Awaiting Decision","resource_type":"enum_option"},{"gid":"1206122140407998","color":"green","enabled":true,"name":"Awaiting Payment (Update Rating Board, Send Invoice and Clothing Email)","resource_type":"enum_option"},{"gid":"1206235855089872","color":"orange","enabled":true,"name":"Completed","resource_type":"enum_option"}]},{"gid":"1205964025396864","enabled":true,"name":"Agent Orange Exposure","description":"","created_by":{"gid":"5609879243379","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205964025396865","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401859348576","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206398481016988","enabled":true,"name":"Burn Pits and Other Airborne Hazards","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206398481016989","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206398481016990","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858650700","enabled":true,"name":"Gulf War Illness","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858650701","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858650702","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401856356610","enabled":true,"name":"Illness Due to Toxic Drinking Water at Camp Lejeune","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401856356611","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401856356612","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858770487","enabled":true,"name":"\"Atomic Veterans\" and Radiation Exposure","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858770488","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858770489","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858924307","enabled":true,"name":"Amyotrophic Lateral Sclerosis (ALS)","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858924308","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858924309","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]}],"due_at":null,"due_on":"2024-03-29","followers":[{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"}],"hearted":false,"hearts":[],"liked":false,"likes":[],"memberships":[{"project":{"gid":"1206472580135542","name":"VBC Client List - Automated - VS","resource_type":"project"},"section":{"gid":"1206472580135543","name":"INCOMING REQUEST","resource_type":"section"}}],"modified_at":"2024-02-22T18:22:01.998Z","name":"Arrabis, Rizalino","notes":"Service:\n\n\nCurrent:\n\n\nNew:\n100 - tdiu","num_hearts":0,"num_likes":0,"parent":null,"permalink_url":"https://app.asana.com/0/1206472580135542/1206343093612429","projects":[{"gid":"1206472580135542","name":"VBC Client List - Automated - VS","resource_type":"project"}],"resource_type":"task","start_at":null,"start_on":null,"tags":[],"resource_subtype":"default_task","workspace":{"gid":"1205444246629009","name":"Veteran Benefits Center LLC","resource_type":"workspace"}}}`

const asanaStr = `{"data":{"gid":"1206398481017098","actual_time_minutes":null,"assignee":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"assignee_status":"inbox","completed":false,"completed_at":null,"created_at":"2024-01-23T20:38:33.814Z","custom_fields":[{"gid":"1206629760929727","enabled":true,"name":"Source","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206629760929729","color":"blue","enabled":true,"name":"Manual","resource_type":"enum_option"},{"gid":"1206629760929728","color":"blue-green","enabled":true,"name":"Website","resource_type":"enum_option"}]},{"gid":"1206184385949353","enabled":true,"name":"First Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Testing","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"Testing"},{"gid":"1206184460895324","enabled":true,"name":"Last Name","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"90","resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text","text_value":"90"},{"gid":"1206539090084664","enabled":true,"name":"Referrer","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747533741","enabled":true,"name":"Email","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747533745","enabled":true,"name":"Phone Number","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747466575","enabled":true,"name":"Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553545","enabled":true,"name":"Effective Current Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206422409732583","enabled":true,"name":"Retired","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206422409732584","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206422409732585","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1205468334222326","enabled":true,"name":"Branch","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205468334222327","color":"blue","enabled":true,"name":"Navy","resource_type":"enum_option"},{"gid":"1205468334222328","color":"yellow-green","enabled":true,"name":"Army","resource_type":"enum_option"},{"gid":"1205468334222329","color":"blue-green","enabled":true,"name":"Air Force","resource_type":"enum_option"},{"gid":"1205468334222330","color":"red","enabled":true,"name":"Marine Corps","resource_type":"enum_option"},{"gid":"1205468334222331","color":"yellow-green","enabled":true,"name":"Army NG","resource_type":"enum_option"},{"gid":"1205960749182653","color":"orange","enabled":true,"name":"Coast Guard","resource_type":"enum_option"}]},{"gid":"1206479747466584","enabled":true,"name":"New Rating","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747466593","enabled":true,"name":"ITF Expiration","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1206479747501015","enabled":true,"name":"Contact Form","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501016","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206640523877079","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501021","enabled":true,"name":"C-File Submitted","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501022","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501023","color":"cool-gray","enabled":true,"name":"N/A","resource_type":"enum_option"},{"gid":"1206479747501024","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501031","enabled":true,"name":"DD214","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501032","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501033","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206478961495299","enabled":true,"name":"Disability Rating List Screenshot","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206478961495300","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206478961495301","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747501039","enabled":true,"name":"Rating Decision Letters","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747501040","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747501041","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533715","enabled":true,"name":"STRs","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533716","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533717","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533724","enabled":true,"name":"TDIU","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747533725","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206479747533726","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206479747533732","enabled":true,"name":"Item ID (auto generated)","description":"","precision":2,"created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"number","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"number"},{"gid":"1206479747553535","enabled":true,"name":"SSN","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553539","enabled":true,"name":"DOB","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"date","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"date"},{"gid":"1205512827064319","enabled":true,"name":"Street Address","description":"Street Address","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206401658215277","enabled":true,"name":"Address - City","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1206479747553554","enabled":true,"name":"Address - State","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206479747553555","color":"none","enabled":true,"name":"Alabama","resource_type":"enum_option"},{"gid":"1206479747553556","color":"none","enabled":true,"name":"Alaska","resource_type":"enum_option"},{"gid":"1206479747553557","color":"none","enabled":true,"name":"Arizona","resource_type":"enum_option"},{"gid":"1206479747553558","color":"none","enabled":true,"name":"Arkansas","resource_type":"enum_option"},{"gid":"1206479747553559","color":"none","enabled":true,"name":"American Samoa","resource_type":"enum_option"},{"gid":"1206479747553560","color":"none","enabled":true,"name":"California","resource_type":"enum_option"},{"gid":"1206479747553561","color":"none","enabled":true,"name":"Colorado","resource_type":"enum_option"},{"gid":"1206479747553562","color":"none","enabled":true,"name":"Connecticut","resource_type":"enum_option"},{"gid":"1206479747553563","color":"none","enabled":true,"name":"Delaware","resource_type":"enum_option"},{"gid":"1206479747553564","color":"none","enabled":true,"name":"District of Columbia","resource_type":"enum_option"},{"gid":"1206479747604120","color":"none","enabled":true,"name":"Florida","resource_type":"enum_option"},{"gid":"1206479747604121","color":"none","enabled":true,"name":"Georgia","resource_type":"enum_option"},{"gid":"1206479747604122","color":"none","enabled":true,"name":"Guam","resource_type":"enum_option"},{"gid":"1206479747604123","color":"none","enabled":true,"name":"Hawaii","resource_type":"enum_option"},{"gid":"1206479747604124","color":"none","enabled":true,"name":"Idaho","resource_type":"enum_option"},{"gid":"1206479747604125","color":"none","enabled":true,"name":"Illinois","resource_type":"enum_option"},{"gid":"1206479747604126","color":"none","enabled":true,"name":"Indiana","resource_type":"enum_option"},{"gid":"1206479747604127","color":"none","enabled":true,"name":"Iowa","resource_type":"enum_option"},{"gid":"1206479747604128","color":"none","enabled":true,"name":"Kansas","resource_type":"enum_option"},{"gid":"1206479747604129","color":"none","enabled":true,"name":"Kentucky","resource_type":"enum_option"},{"gid":"1206479747604130","color":"none","enabled":true,"name":"Louisiana","resource_type":"enum_option"},{"gid":"1206479747604131","color":"none","enabled":true,"name":"Maine","resource_type":"enum_option"},{"gid":"1206479747604132","color":"none","enabled":true,"name":"Maryland","resource_type":"enum_option"},{"gid":"1206479747604133","color":"none","enabled":true,"name":"Massachusetts","resource_type":"enum_option"},{"gid":"1206479747604134","color":"none","enabled":true,"name":"Michigan","resource_type":"enum_option"},{"gid":"1206479747604135","color":"none","enabled":true,"name":"Minnesota","resource_type":"enum_option"},{"gid":"1206479747604136","color":"none","enabled":true,"name":"Mississippi","resource_type":"enum_option"},{"gid":"1206479747604137","color":"none","enabled":true,"name":"Missouri","resource_type":"enum_option"},{"gid":"1206479747604138","color":"none","enabled":true,"name":"Montana","resource_type":"enum_option"},{"gid":"1206479747604139","color":"none","enabled":true,"name":"Nebraska","resource_type":"enum_option"},{"gid":"1206479747604140","color":"none","enabled":true,"name":"Nevada","resource_type":"enum_option"},{"gid":"1206479747604141","color":"none","enabled":true,"name":"New Hampshire","resource_type":"enum_option"},{"gid":"1206479747604142","color":"none","enabled":true,"name":"New Jersey","resource_type":"enum_option"},{"gid":"1206479747604143","color":"none","enabled":true,"name":"New Mexico","resource_type":"enum_option"},{"gid":"1206479747604144","color":"none","enabled":true,"name":"New York","resource_type":"enum_option"},{"gid":"1206479747604145","color":"none","enabled":true,"name":"North Carolina","resource_type":"enum_option"},{"gid":"1206479747604146","color":"none","enabled":true,"name":"North Dakota","resource_type":"enum_option"},{"gid":"1206479747604147","color":"none","enabled":true,"name":"Northern Mariana Islands","resource_type":"enum_option"},{"gid":"1206479747604148","color":"none","enabled":true,"name":"Ohio","resource_type":"enum_option"},{"gid":"1206479747604149","color":"none","enabled":true,"name":"Oklahoma","resource_type":"enum_option"},{"gid":"1206479747604150","color":"none","enabled":true,"name":"Oregon","resource_type":"enum_option"},{"gid":"1206479747604151","color":"none","enabled":true,"name":"Pennsylvania","resource_type":"enum_option"},{"gid":"1206479747639138","color":"none","enabled":true,"name":"Puerto Rico","resource_type":"enum_option"},{"gid":"1206479747639139","color":"none","enabled":true,"name":"Rhode Island","resource_type":"enum_option"},{"gid":"1206479747639140","color":"none","enabled":true,"name":"South Carolina","resource_type":"enum_option"},{"gid":"1206479747639141","color":"none","enabled":true,"name":"South Dakota","resource_type":"enum_option"},{"gid":"1206479747639142","color":"none","enabled":true,"name":"Tennessee","resource_type":"enum_option"},{"gid":"1206479747639143","color":"none","enabled":true,"name":"Texas","resource_type":"enum_option"},{"gid":"1206479747639144","color":"none","enabled":true,"name":"Trust Territories","resource_type":"enum_option"},{"gid":"1206479747639145","color":"none","enabled":true,"name":"Utah","resource_type":"enum_option"},{"gid":"1206479747639146","color":"none","enabled":true,"name":"Vermont","resource_type":"enum_option"},{"gid":"1206479747639147","color":"none","enabled":true,"name":"Virginia","resource_type":"enum_option"},{"gid":"1206479747639148","color":"none","enabled":true,"name":"Virgin Islands","resource_type":"enum_option"},{"gid":"1206479747639149","color":"none","enabled":true,"name":"Washington","resource_type":"enum_option"},{"gid":"1206479747639150","color":"none","enabled":true,"name":"West Virginia","resource_type":"enum_option"},{"gid":"1206479747639151","color":"none","enabled":true,"name":"Wisconsin","resource_type":"enum_option"},{"gid":"1206479747639152","color":"none","enabled":true,"name":"Wyoming","resource_type":"enum_option"},{"gid":"1206479747639153","color":"none","enabled":true,"name":"Outside United States","resource_type":"enum_option"}]},{"gid":"1206401684235934","enabled":true,"name":"Address - Zip Code","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"text","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"text"},{"gid":"1205769148648333","enabled":true,"name":"Stages","description":"","created_by":{"gid":"1205444097333494","name":"Edward Bunting","resource_type":"user"},"display_value":"Getting Started Email","resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206235855089869","color":"purple","enabled":true,"name":"Incoming Request","resource_type":"enum_option"},{"gid":"1206235855089870","color":"indigo","enabled":true,"name":"Fee Schedule and Contract","resource_type":"enum_option"},{"gid":"1206235855089871","color":"red","enabled":true,"name":"Getting Started Email","resource_type":"enum_option"},{"gid":"1205769148648382","color":"aqua","enabled":true,"name":"Awaiting Client Records","resource_type":"enum_option"},{"gid":"1206089580019064","color":"blue-green","enabled":true,"name":"Awaiting C-File","resource_type":"enum_option"},{"gid":"1205769148648383","color":"blue-green","enabled":true,"name":"Record Review","resource_type":"enum_option"},{"gid":"1205769148648385","color":"yellow-green","enabled":true,"name":"Schedule Call","resource_type":"enum_option"},{"gid":"1205769148648386","color":"yellow","enabled":true,"name":"Statement Drafts","resource_type":"enum_option"},{"gid":"1205769148648387","color":"yellow-orange","enabled":true,"name":"Statements Finalized","resource_type":"enum_option"},{"gid":"1205769148648388","color":"pink","enabled":true,"name":"Current Treatment","resource_type":"enum_option"},{"gid":"1205773364138828","color":"hot-pink","enabled":true,"name":"Mini-DBQs","resource_type":"enum_option"},{"gid":"1205773364138829","color":"magenta","enabled":true,"name":"Nexus Letters","resource_type":"enum_option"},{"gid":"1205773364138830","color":"purple","enabled":true,"name":"Medical Team","resource_type":"enum_option"},{"gid":"1205773364138831","color":"indigo","enabled":true,"name":"File Claims (New or Supplemental)","resource_type":"enum_option"},{"gid":"1205773364138832","color":"blue","enabled":true,"name":"Verify Evidence Received","resource_type":"enum_option"},{"gid":"1205773364138833","color":"red","enabled":true,"name":"Awaiting Decision","resource_type":"enum_option"},{"gid":"1206122140407998","color":"green","enabled":true,"name":"Awaiting Payment (Update Rating Board, Send Invoice and Clothing Email)","resource_type":"enum_option"},{"gid":"1206235855089872","color":"orange","enabled":true,"name":"Completed","resource_type":"enum_option"}],"enum_value":{"color":"red","enabled":true,"gid":"1206235855089871","name":"Getting Started Email","resource_type":"enum_option"}},{"gid":"1205964025396864","enabled":true,"name":"Agent Orange Exposure","description":"","created_by":{"gid":"5609879243379","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1205964025396865","color":"yellow-green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401859348576","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206398481016988","enabled":true,"name":"Burn Pits and Other Airborne Hazards","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206398481016989","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206398481016990","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858650700","enabled":true,"name":"Gulf War Illness","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858650701","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858650702","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401856356610","enabled":true,"name":"Illness Due to Toxic Drinking Water at Camp Lejeune","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401856356611","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401856356612","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858770487","enabled":true,"name":"\"Atomic Veterans\" and Radiation Exposure","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858770488","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858770489","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]},{"gid":"1206401858924307","enabled":true,"name":"Amyotrophic Lateral Sclerosis (ALS)","description":"","created_by":{"gid":"1205826312615618","name":"Yannan Wang","resource_type":"user"},"display_value":null,"resource_subtype":"enum","resource_type":"custom_field","is_formula_field":false,"is_value_read_only":false,"type":"enum","enum_options":[{"gid":"1206401858924308","color":"green","enabled":true,"name":"Yes","resource_type":"enum_option"},{"gid":"1206401858924309","color":"red","enabled":true,"name":"No","resource_type":"enum_option"}]}],"due_at":null,"due_on":null,"followers":[],"hearted":false,"hearts":[],"liked":false,"likes":[],"memberships":[{"project":{"gid":"1205962480830589","name":"VBC Client List for Testing - VS","resource_type":"project"},"section":{"gid":"1205962480830599","name":"GETTING STARTED EMAIL","resource_type":"section"}}],"modified_at":"2024-02-24T10:27:20.276Z","name":"Testing 90","notes":"","num_hearts":0,"num_likes":0,"parent":null,"permalink_url":"https://app.asana.com/0/1205962480830589/1206398481017098","projects":[{"gid":"1205962480830589","name":"VBC Client List for Testing - VS","resource_type":"project"}],"resource_type":"task","start_at":null,"start_on":null,"tags":[],"resource_subtype":"default_task","workspace":{"gid":"1205444246629009","name":"Veteran Benefits Center LLC","resource_type":"workspace"}}}`

//const asanaStr = `{
//	"data": {
//		"gid": "1206237015703342",
//		"actual_time_minutes": null,
//		"assignee": {
//			"gid": "1205826312615618",
//			"name": "Yannan Wang",
//			"resource_type": "user"
//		},
//		"assignee_status": "inbox",
//		"completed": false,
//		"completed_at": null,
//		"created_at": "2023-12-25T14:25:45.166Z",
//		"custom_fields": [{
//			"gid": "1206184385949353",
//			"enabled": true,
//			"name": "First Name",
//			"description": "",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": "Shi",
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text",
//			"text_value": "Shi"
//		}, {
//			"gid": "1206184460895324",
//			"enabled": true,
//			"name": "Last Name",
//			"description": "",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": "Li",
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text",
//			"text_value": "Li"
//		}, {
//			"gid": "1205468334222326",
//			"enabled": true,
//			"name": "Branch",
//			"description": "",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": "Coast Guard",
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205468334222327",
//				"color": "blue",
//				"enabled": true,
//				"name": "Navy",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205468334222328",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Army",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205468334222329",
//				"color": "blue-green",
//				"enabled": true,
//				"name": "Air Force",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205468334222330",
//				"color": "red",
//				"enabled": true,
//				"name": "Marine Corps",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205468334222331",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Army NG",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205960749182653",
//				"color": "orange",
//				"enabled": true,
//				"name": "Coast Guard",
//				"resource_type": "enum_option"
//			}],
//			"enum_value": {
//				"color": "orange",
//				"enabled": true,
//				"gid": "1205960749182653",
//				"name": "Coast Guard",
//				"resource_type": "enum_option"
//			}
//		}, {
//			"gid": "1205964024662896",
//			"enabled": true,
//			"name": "Current Rating",
//			"description": "",
//			"number_value": 0,
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "0",
//			"resource_subtype": "number",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "number"
//		}, {
//			"gid": "1206241124341932",
//			"enabled": true,
//			"name": "Effective Current Rating",
//			"description": "",
//			"number_value": 0,
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "0",
//			"resource_subtype": "number",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "number"
//		}, {
//			"gid": "1205964024662905",
//			"enabled": true,
//			"name": "New Rating",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "number",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "number"
//		}, {
//			"gid": "1205964025385462",
//			"enabled": true,
//			"name": "ITF Expiration",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "date",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "date"
//		}, {
//			"gid": "1205964025385466",
//			"enabled": true,
//			"name": "Contact Form",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025385467",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025385472",
//			"enabled": true,
//			"name": "C-File Submitted",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025385473",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025385474",
//				"color": "cool-gray",
//				"enabled": true,
//				"name": "N/A",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025385475",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025385482",
//			"enabled": true,
//			"name": "DD214",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025385483",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025385484",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025385490",
//			"enabled": true,
//			"name": "Rating Decision Letters",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025385491",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025396850",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025396856",
//			"enabled": true,
//			"name": "STRs",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025396857",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025396858",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025396870",
//			"enabled": true,
//			"name": "TDIU",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025396871",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205964025396872",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025396878",
//			"enabled": true,
//			"name": "Item ID (auto generated)",
//			"description": "",
//			"precision": 2,
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "number",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "number"
//		}, {
//			"gid": "1205964025409303",
//			"enabled": true,
//			"name": "Email",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "yannanwang@gmail.com",
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text",
//			"text_value": "yannanwang@gmail.com"
//		}, {
//			"gid": "1205964025409307",
//			"enabled": true,
//			"name": "Phone Number",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "118-776-6677",
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text",
//			"text_value": "118-776-6677"
//		}, {
//			"gid": "1205964025409311",
//			"enabled": true,
//			"name": "SSN",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "111-21-8888",
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text",
//			"text_value": "111-21-8888"
//		}, {
//			"gid": "1205964025409315",
//			"enabled": true,
//			"name": "DOB",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": "1987-12-22T00:00:00.000Z",
//			"resource_subtype": "date",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "date",
//			"date_value": {
//				"date": "1987-12-22",
//				"date_time": null
//			}
//		}, {
//			"gid": "1205512827064319",
//			"enabled": true,
//			"name": "Street Address",
//			"description": "Street Address",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text"
//		}, {
//			"gid": "1206401658215277",
//			"enabled": true,
//			"name": "Address - City",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text"
//		}, {
//			"gid": "1206398481016928",
//			"enabled": true,
//			"name": "Address - State",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206398481016929",
//				"color": "none",
//				"enabled": true,
//				"name": "Alabama",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016930",
//				"color": "none",
//				"enabled": true,
//				"name": "Alaska",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016931",
//				"color": "none",
//				"enabled": true,
//				"name": "Arizona",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016932",
//				"color": "none",
//				"enabled": true,
//				"name": "Arkansas",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016933",
//				"color": "none",
//				"enabled": true,
//				"name": "American Samoa",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016934",
//				"color": "none",
//				"enabled": true,
//				"name": "California",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016935",
//				"color": "none",
//				"enabled": true,
//				"name": "Colorado",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016936",
//				"color": "none",
//				"enabled": true,
//				"name": "Connecticut",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016937",
//				"color": "none",
//				"enabled": true,
//				"name": "Delaware",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016938",
//				"color": "none",
//				"enabled": true,
//				"name": "District of Columbia",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016939",
//				"color": "none",
//				"enabled": true,
//				"name": "Florida",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016940",
//				"color": "none",
//				"enabled": true,
//				"name": "Georgia",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016941",
//				"color": "none",
//				"enabled": true,
//				"name": "Guam",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016942",
//				"color": "none",
//				"enabled": true,
//				"name": "Hawaii",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016943",
//				"color": "none",
//				"enabled": true,
//				"name": "Idaho",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016944",
//				"color": "none",
//				"enabled": true,
//				"name": "Illinois",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016945",
//				"color": "none",
//				"enabled": true,
//				"name": "Indiana",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016946",
//				"color": "none",
//				"enabled": true,
//				"name": "Iowa",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016947",
//				"color": "none",
//				"enabled": true,
//				"name": "Kansas",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016948",
//				"color": "none",
//				"enabled": true,
//				"name": "Kentucky",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016949",
//				"color": "none",
//				"enabled": true,
//				"name": "Louisiana",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016950",
//				"color": "none",
//				"enabled": true,
//				"name": "Maine",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016951",
//				"color": "none",
//				"enabled": true,
//				"name": "Maryland",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016952",
//				"color": "none",
//				"enabled": true,
//				"name": "Massachusetts",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016953",
//				"color": "none",
//				"enabled": true,
//				"name": "Michigan",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016954",
//				"color": "none",
//				"enabled": true,
//				"name": "Minnesota",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016955",
//				"color": "none",
//				"enabled": true,
//				"name": "Mississippi",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016956",
//				"color": "none",
//				"enabled": true,
//				"name": "Missouri",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016957",
//				"color": "none",
//				"enabled": true,
//				"name": "Montana",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016958",
//				"color": "none",
//				"enabled": true,
//				"name": "Nebraska",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016959",
//				"color": "none",
//				"enabled": true,
//				"name": "Nevada",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016960",
//				"color": "none",
//				"enabled": true,
//				"name": "New Hampshire",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016961",
//				"color": "none",
//				"enabled": true,
//				"name": "New Jersey",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016962",
//				"color": "none",
//				"enabled": true,
//				"name": "New Mexico",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016963",
//				"color": "none",
//				"enabled": true,
//				"name": "New York",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016964",
//				"color": "none",
//				"enabled": true,
//				"name": "North Carolina",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016965",
//				"color": "none",
//				"enabled": true,
//				"name": "North Dakota",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016966",
//				"color": "none",
//				"enabled": true,
//				"name": "Northern Mariana Islands",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016967",
//				"color": "none",
//				"enabled": true,
//				"name": "Ohio",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016968",
//				"color": "none",
//				"enabled": true,
//				"name": "Oklahoma",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016969",
//				"color": "none",
//				"enabled": true,
//				"name": "Oregon",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016970",
//				"color": "none",
//				"enabled": true,
//				"name": "Pennsylvania",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016971",
//				"color": "none",
//				"enabled": true,
//				"name": "Puerto Rico",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016972",
//				"color": "none",
//				"enabled": true,
//				"name": "Rhode Island",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016973",
//				"color": "none",
//				"enabled": true,
//				"name": "South Carolina",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016974",
//				"color": "none",
//				"enabled": true,
//				"name": "South Dakota",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016975",
//				"color": "none",
//				"enabled": true,
//				"name": "Tennessee",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016976",
//				"color": "none",
//				"enabled": true,
//				"name": "Texas",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016977",
//				"color": "none",
//				"enabled": true,
//				"name": "Trust Territories",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016978",
//				"color": "none",
//				"enabled": true,
//				"name": "Utah",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016979",
//				"color": "none",
//				"enabled": true,
//				"name": "Vermont",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016980",
//				"color": "none",
//				"enabled": true,
//				"name": "Virginia",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016981",
//				"color": "none",
//				"enabled": true,
//				"name": "Virgin Islands",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016982",
//				"color": "none",
//				"enabled": true,
//				"name": "Washington",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016983",
//				"color": "none",
//				"enabled": true,
//				"name": "West Virginia",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016984",
//				"color": "none",
//				"enabled": true,
//				"name": "Wisconsin",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016985",
//				"color": "none",
//				"enabled": true,
//				"name": "Wyoming",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016986",
//				"color": "none",
//				"enabled": true,
//				"name": "Outside United States",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206401684235934",
//			"enabled": true,
//			"name": "Address - Zip Code",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "text",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "text"
//		}, {
//			"gid": "1205769148648333",
//			"enabled": true,
//			"name": "Stages",
//			"description": "",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": "Fee Schedule and Contract",
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206235855089869",
//				"color": "purple",
//				"enabled": true,
//				"name": "Incoming Request",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206235855089870",
//				"color": "indigo",
//				"enabled": true,
//				"name": "Fee Schedule and Contract",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206235855089871",
//				"color": "red",
//				"enabled": true,
//				"name": "Getting Started Email",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648382",
//				"color": "aqua",
//				"enabled": true,
//				"name": "Update Client Intake Info",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206089580019064",
//				"color": "blue-green",
//				"enabled": true,
//				"name": "Awaiting C-File",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648383",
//				"color": "blue-green",
//				"enabled": true,
//				"name": "Record Review",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648385",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Schedule Call",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648386",
//				"color": "yellow",
//				"enabled": true,
//				"name": "Statement Drafts",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648387",
//				"color": "yellow-orange",
//				"enabled": true,
//				"name": "Statements Finalized",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205769148648388",
//				"color": "pink",
//				"enabled": true,
//				"name": "Current Treatment",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138828",
//				"color": "hot-pink",
//				"enabled": true,
//				"name": "Mini-DBQs",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138829",
//				"color": "magenta",
//				"enabled": true,
//				"name": "Nexus Letters",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138830",
//				"color": "purple",
//				"enabled": true,
//				"name": "Medical Team",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138831",
//				"color": "indigo",
//				"enabled": true,
//				"name": "File Claims (New or Supplemental)",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138832",
//				"color": "blue",
//				"enabled": true,
//				"name": "Verify Evidence Received",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773364138833",
//				"color": "red",
//				"enabled": true,
//				"name": "Awaiting Decision",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1205773365925299",
//				"color": "aqua",
//				"enabled": true,
//				"name": "Update Rating Board, Send Invoice and Clothing Email",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206122140407998",
//				"color": "green",
//				"enabled": true,
//				"name": "Awaiting Payment",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206235855089872",
//				"color": "orange",
//				"enabled": true,
//				"name": "Completed",
//				"resource_type": "enum_option"
//			}],
//			"enum_value": {
//				"color": "indigo",
//				"enabled": true,
//				"gid": "1206235855089870",
//				"name": "Fee Schedule and Contract",
//				"resource_type": "enum_option"
//			}
//		}, {
//			"gid": "1205769148648348",
//			"enabled": true,
//			"name": "Status",
//			"description": "",
//			"created_by": {
//				"gid": "1205444097333494",
//				"name": "Edward Bunting",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205769148648349",
//				"color": "red",
//				"enabled": true,
//				"name": "Overdue",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1205964025396864",
//			"enabled": true,
//			"name": "Agent Orange Exposure",
//			"description": "",
//			"created_by": {
//				"gid": "5609879243379",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1205964025396865",
//				"color": "yellow-green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206401859348576",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206398481016988",
//			"enabled": true,
//			"name": "Burn Pits and Other Airborne Hazards",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206398481016989",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206398481016990",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206401858650700",
//			"enabled": true,
//			"name": "Gulf War Illness",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206401858650701",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206401858650702",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206401856356610",
//			"enabled": true,
//			"name": "Illness Due to Toxic Drinking Water at Camp Lejeune",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206401856356611",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206401856356612",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206401858770487",
//			"enabled": true,
//			"name": "\"Atomic Veterans\" and Radiation Exposure",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206401858770488",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206401858770489",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206401858924307",
//			"enabled": true,
//			"name": "Amyotrophic Lateral Sclerosis (ALS)",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206401858924308",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206401858924309",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}, {
//			"gid": "1206422409732583",
//			"enabled": true,
//			"name": "Retired",
//			"description": "",
//			"created_by": {
//				"gid": "1205826312615618",
//				"name": "Yannan Wang",
//				"resource_type": "user"
//			},
//			"display_value": null,
//			"resource_subtype": "enum",
//			"resource_type": "custom_field",
//			"is_formula_field": false,
//			"is_value_read_only": false,
//			"type": "enum",
//			"enum_options": [{
//				"gid": "1206422409732584",
//				"color": "green",
//				"enabled": true,
//				"name": "Yes",
//				"resource_type": "enum_option"
//			}, {
//				"gid": "1206422409732585",
//				"color": "red",
//				"enabled": true,
//				"name": "No",
//				"resource_type": "enum_option"
//			}]
//		}],
//		"due_at": null,
//		"due_on": null,
//		"followers": [{
//			"gid": "1206230291638946",
//			"name": "liaogling",
//			"resource_type": "user"
//		}, {
//			"gid": "1205826312615618",
//			"name": "Yannan Wang",
//			"resource_type": "user"
//		}],
//		"hearted": false,
//		"hearts": [],
//		"liked": false,
//		"likes": [],
//		"memberships": [{
//			"project": {
//				"gid": "1205962480830589",
//				"name": "VBC Client List for Testing",
//				"resource_type": "project"
//			},
//			"section": {
//				"gid": "1206176340563332",
//				"name": "FEE SCHEDULE AND CONTRACT",
//				"resource_type": "section"
//			}
//		}],
//		"modified_at": "2024-01-23T19:22:36.003Z",
//		"name": "Li Shi",
//		"notes": "",
//		"num_hearts": 0,
//		"num_likes": 0,
//		"parent": null,
//		"permalink_url": "https://app.asana.com/0/1205962480830589/1206237015703342",
//		"projects": [{
//			"gid": "1205962480830589",
//			"name": "VBC Client List for Testing",
//			"resource_type": "project"
//		}],
//		"resource_type": "task",
//		"start_at": null,
//		"start_on": null,
//		"tags": [],
//		"resource_subtype": "default_task",
//		"workspace": {
//			"gid": "1205444246629009",
//			"name": "Veteran Benefits Center LLC",
//			"resource_type": "workspace"
//		}
//	}
//}`
