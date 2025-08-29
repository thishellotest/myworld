package config_zoho

import "vbc/lib"

var UserMappingConfigs = map[string]string{
	"email": "email",
	//"first_name":    "first_name", // 不同步防止zoho修改影响业务
	//"last_name":     "last_name", // 不同步防止zoho修改影响业务
	"full_name":     "full_name", // 需要同步 与client case :Lead VS , Lead CP 映射关联
	"id":            "gid",
	"Modified_Time": "modified_time",
	"created_time":  "created_time",
}

func UserMappings(row lib.TypeMap) lib.TypeMap {

	if row == nil {
		return nil
	}
	res := make(lib.TypeMap)
	for k, v := range UserMappingConfigs {
		res.Set(v, row.GetString(k))
	}
	return res
}
