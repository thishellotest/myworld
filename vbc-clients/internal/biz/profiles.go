package biz

const (
	Profile_IsAdmin_Yes = 1

	Profile_FieldName_is_admin = "is_admin"

	Profile_Standard_Gid    = "441f19d51858417cb948cc286ef1b585"
	Profile_TeamLeaders_gid = "c695eb2366ff4796a7806feb6b591740"
)

func HaveAllDataPermissions(kind string, profileGid string) bool {
	if profileGid == Profile_TeamLeaders_gid && (kind == Kind_client_cases || kind == Kind_clients) {
		return true
	}
	return false
}

func IsAdminProfile(tProfile *TData) bool {
	if tProfile.CustomFields.NumberValueByNameBasic(Profile_FieldName_is_admin) == Profile_IsAdmin_Yes {
		return true
	}
	return false
}
