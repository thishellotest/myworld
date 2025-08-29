package biz

type ResponseUser struct {
	Gid  string `json:"gid"`
	Name string `json:"name"` // 一般用于显示
}

type ResponseRelatedRecord struct {
	ModuleName  string `json:"module_name"`
	ModuleLabel string `json:"module_label"`
	Gid         string `json:"gid"`
	Name        string `json:"name"` // 一般用于显示
}
