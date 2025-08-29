package biz

const Kind_common = "common"

const Kind_client_cases = ""
const Kind_clients = "clients"
const Kind_client_tasks = "client_tasks"
const Kind_email_tpls = "email_tpls"
const Kind_users = "users"
const Kind_blobs = "blobs"
const Kind_notes = "notes"
const Kind_timelines = "timelines"
const Kind_attorneys = "attorneys"

const Kind_roles = "roles"
const Kind_profiles = "profiles"

const Kind_global_sets = "global_sets" // 特殊的

const Kind_Custom_Condition = "custom_condition"
const Kind_Custom_ConditionSecondary = "custom_condition_secondary"
const Kind_Custom_Filter = "custom_filter"
const Kind_Custom_ConditionQuestionnaires = "custom_condition_questionnaires"

func IsCustomKind(kind string) bool {
	if kind == Kind_Custom_Condition || kind == Kind_Custom_ConditionSecondary ||
		kind == Kind_Custom_Filter || kind == Kind_Custom_ConditionQuestionnaires {
		return true
	}
	return false
}
