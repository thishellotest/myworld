package config_zoho

import "vbc/lib"

var TaskMappingConfigs = map[string]string{
	"$se_module":     "se_module",
	"Created_By.id":  "created_by",
	"Modified_By.id": "modified_by",
	"Owner.id":       "user_gid",
	"Priority":       "priority",
	"Status":         "status",
	"Subject":        "subject",
	"What_Id.id":     "what_id_gid",
	"Who_Id.id":      "who_id_gid",
	"id":             "gid",
	"Modified_Time":  "modified_time",
	"Created_Time":   "created_time",
	"Due_Date":       "due_date",
}

func TasksMappings(row lib.TypeMap) lib.TypeMap {
	if row == nil {
		return nil
	}
	res := make(lib.TypeMap)
	for k, v := range TaskMappingConfigs {
		res.Set(v, row.GetString(k))
	}
	return res
}
