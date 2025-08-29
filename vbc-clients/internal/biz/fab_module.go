package biz

const ModuleCases = "cases"
const ModuleTasks = "tasks"

func ModuleConvertToKind(moduleName string) string {
	if moduleName == ModuleCases {
		return ""
	} else if moduleName == ModuleTasks {
		return Kind_client_tasks
	}
	return moduleName
}

func KindConvertToModule(kind string) string {
	if kind == "" {
		return ModuleCases
	} else if kind == Kind_client_tasks {
		return ModuleTasks
	}
	return kind
}
