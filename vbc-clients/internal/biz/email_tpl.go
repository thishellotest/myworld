package biz

import (
	"bytes"
	"sync"
	"text/template"
	"vbc/configs"
	"vbc/lib"
)

func GlobTpl() *template.Template {
	var once sync.Once
	var tmpl *template.Template
	once.Do(
		func() {
			funcMap := template.FuncMap{
				// The name "inc" is what the function will be called in the template text.
				"inc": func(i int) int {
					return i + 1
				},
			}
			if configs.IsProd() {
				tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("/app/templates/email/*.html"))
			} else {
				tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("/Users/garyliao/code/vbc-clients/templates/email/*.html"))
			}

		})
	return tmpl
}

func FollowingUpSignMedicalTeamFormsEmailBody(pageTitle string, fullName string, items lib.TypeList) (string, error) {

	typeMap := make(lib.TypeMap)
	typeMap.Set("PageTitle", pageTitle)
	typeMap.Set("FullName", fullName)
	typeMap.Set("Items", items)
	/*typeMap.Set("Items", lib.TypeList{{
		"Label": "Subject",
		"Value": "Following Up with Clients to Sign Medical Team Forms",
	}, {
		"Label": "Client Case Name",
		"Value": "TestFN TestLN-80#5093",
	}, {
		"Label": "Contract Sent On",
		"Value": "Sat, 18 May 2024 03:59 PM",
	}, {
		"Label": "Client Case URL",
		"Value": "<a target=\"_blank\" href=\"https://crm.zoho.com/crm/org847391426/tab/Potentials/6159272000003416077\">https://crm.zoho.com/crm/org847391426/tab/Potentials/6159272000003416077</a>",
	}})*/

	body := make([]byte, 0)
	buffer := bytes.NewBuffer(body)
	err := GlobTpl().ExecuteTemplate(buffer, "remind.html", typeMap)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func FollowingUpUploadedDocumentEmailBody(pageTitle string, fullName string, clientCaseName string, updateFiles []*ReminderClientUpdateFilesEventVoItem) (string, error) {

	typeMap := make(lib.TypeMap)
	typeMap.Set("PageTitle", pageTitle)
	typeMap.Set("FullName", fullName)
	typeMap.Set("ClientCaseName", clientCaseName)
	typeMap.Set("UpdateFiles", updateFiles)
	/*typeMap.Set("Items", lib.TypeList{{
		"Label": "Subject",
		"Value": "Following Up with Clients to Sign Medical Team Forms",
	}, {
		"Label": "Client Case Name",
		"Value": "TestFN TestLN-80#5093",
	}, {
		"Label": "Contract Sent On",
		"Value": "Sat, 18 May 2024 03:59 PM",
	}, {
		"Label": "Client Case URL",
		"Value": "<a target=\"_blank\" href=\"https://crm.zoho.com/crm/org847391426/tab/Potentials/6159272000003416077\">https://crm.zoho.com/crm/org847391426/tab/Potentials/6159272000003416077</a>",
	}})*/

	body := make([]byte, 0)
	buffer := bytes.NewBuffer(body)
	err := GlobTpl().ExecuteTemplate(buffer, "remind_newfiles.html", typeMap)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func CaseWithoutTasksEmailBody(pageTitle string, result []*CaseWithoutTaskVo) (string, error) {

	//lib.DPrintln(result)

	typeMap := make(lib.TypeMap)
	typeMap.Set("PageTitle", pageTitle)
	typeMap.Set("Items", result)
	body := make([]byte, 0)
	buffer := bytes.NewBuffer(body)
	err := GlobTpl().ExecuteTemplate(buffer, "case_without_tasks.html", typeMap)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
