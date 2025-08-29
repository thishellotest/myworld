package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_GetJotformSubmissionIdFromFileName(t *testing.T) {
	// 要处理的字符串列表
	files := []string{
		"5004-Back-New-Back pain with radiculopathy secondary to Right knee strain with tibial stress fractures -597484667-4216669565.pdf",
		"5004-Back-New-Back pain with radiculopathy secondary to Right knee strain with tibial stress fractures -5974846674216669565.pdf",
		"5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf",
		"5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf",
		"5004-Knee-Increase-Right knee strain with tibial stress fractures with limitation of flexion and extension (increase)Bilateral left knee pain secondary to Right knee strain with tibial stress fractures with limitation of flexion and extension (opinion)-5974845864217999778.pdf",
		"5004-Mental Disorders Secondaries-New-Major Depressive Disorder secondary to Right knee strain with tibial stress fractures-5974857144211924747.pdf",
	}

	// 正则表达式：匹配 - 和 .pdf 之间的数字，数字长度不固定
	//re := regexp.MustCompile(`-(\d+)\.pdf`)

	// 遍历文件名列表并提取匹配的数字
	for _, file := range files {
		aaa := biz.GetJotformSubmissionIdFromFileName(file)
		lib.DPrintln(aaa)

		//matches := re.FindStringSubmatch(file)
		//if len(matches) > 0 {
		//	fmt.Println(matches[1]) // 打印匹配到的数字
		//}
	}
}

func Test_GetJsonFromAiResult(t *testing.T) {

	text := `这里有一些文本：
{
  "related_entries": [
    "5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf",
    "5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf"
  ]
}
后面还有一些文本。`
	text = `{
  "state": "North Carolina",
  "city": "Whitsett",
  "timezone": "America/New_York"
}`

	//text = "```json\n{\n  \"related_entries\": [\n    \"5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf\",\n    \"5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf\"\n  ]\n}\n```"

	aaa := biz.GetJsonFromAiResult(text)
	lib.DPrintln(aaa)
	//
	//// 正则表达式：匹配最外层的 JSON 对象（简单假设只有一个 {} 块）
	//re := regexp.MustCompile(`(?s)\{.*?\}`) // (?s) 让 . 匹配换行
	//
	//match := re.FindString(text)
	//if match != "" {
	//	fmt.Println("找到 JSON:")
	//	fmt.Println(match)
	//} else {
	//	fmt.Println("没有找到 JSON。")
	//}

}

//aaa:=`From the list below, identify the entries that are directly or indirectly related to "{{condition}}". A connection can be:
//- the condition is caused by {{condition}} (e.g., "secondary to {{condition}}"),
//- {{condition}} is a mentioned symptom or diagnosis, or {{condition}} is part of the main condition description.
//
//Return only the related entries in valid JSON format like this:
//{
//  "related_entries": [
//    "entry1",
//    "entry2"
//  ]
//}
//
//Here is the data list:
//{{data_list}}`

func Test_VbcAIUsecase_Claude3(t *testing.T) {
	r, i, err := UT.VbcAIUsecase.Claude3("", `From the list below, identify the entries that are directly or indirectly related to "Tinnitus". A connection can be:
- the condition is caused by Tinnitus (e.g., "secondary to Tinnitus"),
- Tinnitus is a mentioned symptom or diagnosis,
- or Tinnitus is part of the main condition description.

Return only the related entries in valid JSON format like this:
{
  "related_entries": [
    "entry1",
    "entry2"
  ]
}

Here is the data list:
5004-Back-New-Back pain with radiculopathy secondary to Right knee strain with tibial stress fractures -5974846674216669565.pdf  
5004-Headaches and Migraines-New-Headaches secondary to tinnitus -5974825784219522597.pdf  
5004-Hearing Loss and Tinnitus-New-Tinnitus-5974825894217185001.pdf  
5004-Knee-Increase-Right knee strain with tibial stress fractures with limitation of flexion and extension (increase)Bilateral left knee pain secondary to Right knee strain with tibial stress fractures with limitation of flexion and extension (opinion)-5974845864217999778.pdf  
5004-Mental Disorders Secondaries-New-Major Depressive Disorder secondary to Right knee strain with tibial stress fractures-5974857144211924747.pdf`, "")
	lib.DPrintln(r, i, err)
}

func Test_VbcAIUsecase_Claude333(t *testing.T) {
	str := "California, Los Angeles"
	str = "Pennsylvania, Warrington"
	r, i, err := UT.VbcAIUsecase.Claude3(`Given a U.S. state and city, return the corresponding time zone in JSON format. Only respond with one of the following valid IANA time zones:
"Pacific/Honolulu"
"America/Anchorage"
"America/Chicago"
"America/Denver"
"America/Los_Angeles"
"America/New_York"
Respond in the following JSON format:
{
  "state": "<State>",
  "city": "<City>",
  "timezone": "<IANA Time Zone>"
}
`, str, "")
	lib.DPrintln(r, i, err)
}

func Test_VbcAIUsecase_Claude4(t *testing.T) {
	str := "Given the condition: \"Back pain secondary to bilateral pes planus\"\n\nFrom the list below, return only the entries that are directly or explicitly related to this condition. A related entry must meet **at least one** of the following criteria:\n\n1. The entry's title **exactly matches** the given condition or includes it as a primary diagnosis (e.g., \"New-[condition]\" or \"[condition]-aggravated by service\").\n2. The entry **explicitly states** the given condition as a diagnosis or symptom within its title (not inferred or implied by a shared term).\n3. The condition is the **main subject** of the entry, rather than a contributing factor to another diagnosis.\n\nDo **not** return entries that:\n- Only contain partial phrases or terms (e.g., \"bilateral pes planus\" alone is not sufficient),\n- Refer to other conditions that merely include components of the given condition,\n- Are indirectly related via other diagnoses.\n\nReturn the results in the following JSON format:\n{\n  \"related_entries\": [\n    \"entry1\",\n    \"entry2\"\n  ]\n}\n\nHere is the data list:\nHearing Loss and Tinnitus-New-Tinnitus-6169292782894013022.pdf  \nHypertension-New-Hypertension-6169303632898518752.pdf  \nFoot--Bilateral pes planus aggravated by service-6169310762892987299.pdf  \nBack-New-Back pain secondary to bilateral pes planus-6169316572895335368.pdf\nMental Disorders-Increase-Insomnia Disorder-6169340152892747999.pdf  \nMale Reproductive Organ Conditions-New-Erectile  dysfunction secondary to insomnia disorder-6169352952897678956.pdf"
	r, i, err := UT.VbcAIUsecase.Claude3(``, str, "")
	lib.DPrintln(r, i, err)
}

func Test_VbcAIUsecase_Claude5(t *testing.T) {
	//str := "Given the condition: \"Degenerative Arthritis; right knee, patellofemoral chondromalacia\"\n\nFrom the list below, return only the entries that are directly or explicitly related to this condition. A related entry must meet **at least one** of the following criteria:\n\n1. The entry's title **exactly matches** the given condition or includes it as a primary diagnosis (e.g., \"New-[condition]\" or \"[condition]-aggravated by service\").\n2. The entry **explicitly states** the given condition as a diagnosis or symptom within its title (not inferred or implied by a shared term).\n3. The condition is the **main subject** of the entry, rather than a contributing factor to another diagnosis.\n\nDo **not** return entries that:\n- Only contain partial phrases or terms (e.g., \"bilateral pes planus\" alone is not sufficient),\n- Refer to other conditions that merely include components of the given condition,\n- Are indirectly related via other diagnoses.\n\nReturn the results in the following JSON format:\n{\n  \"related_entries\": [\n    \"entry1\",\n    \"entry2\"\n  ]\n}\n\nHere is the data list:\nHearing Loss and Tinnitus-New-Tinnitus-6194989904217511966.pdf\nHeadaches and Migraines-Increase-Tension headache-6194996724217746641.pdf\nEar-New-Dizziness-6195001744215515843.pdf\nSinus-New-Allergic rhinitishay fever-6195006264215396564.pdf\nAnkle-New-Left ankle sprain-6195014234219988921.pdf\nKnee-Increase-Traumatic Arthritis; sp two ACL reconstructions, left knee-6195024684219386014.pdf\nKnee-Increase-Degenerative Arthritis; right knee, patellofemoral chondromalacia-6195032454212156629.pdf\nScar-Increase-Painful scars, right and left knees-6195037934213042032.pdf\nBack-New-Low back pain with radiculopathy in right lower extremity secondary to degenerative Arthritis; right knee, patellofemoral chondromalacia and traumatic Arthritis; sp two ACL reconstructions, left knee-6195044314219896456.pdf\nPost-Traumatic Stress Disorder (PTSD)--6195050814212857188.pdf\nMale Reproductive Organ Conditions-New-Erectile dysfunction secondary to PTSD-6195055014215176355.pdf"
	//r, i, err := UT.VbcAIUsecase.Claude3(``, str, "")
	//lib.DPrintln(r, i, err)

	str := "Given the condition: \"Painful scars, right and left knees\"\n\nFrom the list below, return only the entries that are directly or explicitly related to this condition. A related entry must meet **at least one** of the following criteria:\n\n1. The entry's title **exactly matches** the given condition or includes it as a primary diagnosis (e.g., \"New-[condition]\" or \"[condition]-aggravated by service\").\n2. The entry **explicitly states** the given condition as a diagnosis or symptom within its title (not inferred or implied by a shared term).\n3. The condition is the **main subject** of the entry, rather than a contributing factor to another diagnosis.\n\nDo **not** return entries that:\n- Only contain partial phrases or terms (e.g., \"bilateral pes planus\" alone is not sufficient),\n- Refer to other conditions that merely include components of the given condition,\n- Are indirectly related via other diagnoses.\n\nReturn the results in the following JSON format:\n{\n  \"related_entries\": [\n    \"entry1\",\n    \"entry2\"\n  ]\n}\n\nHere is the data list:\nHearing Loss and Tinnitus-New-Tinnitus-6194989904217511966.pdf\nHeadaches and Migraines-Increase-Tension headache-6194996724217746641.pdf\nEar-New-Dizziness-6195001744215515843.pdf\nSinus-New-Allergic rhinitishay fever-6195006264215396564.pdf\nAnkle-New-Left ankle sprain-6195014234219988921.pdf\nKnee-Increase-Traumatic Arthritis; sp two ACL reconstructions, left knee-6195024684219386014.pdf\nKnee-Increase-Degenerative Arthritis; right knee, patellofemoral chondromalacia-6195032454212156629.pdf\nScar-Increase-Painful scars, right and left knees-6195037934213042032.pdf\nBack-New-Low back pain with radiculopathy in right lower extremity secondary to degenerative Arthritis; right knee, patellofemoral chondromalacia and traumatic Arthritis; sp two ACL reconstructions, left knee-6195044314219896456.pdf\nPost-Traumatic Stress Disorder (PTSD)--6195050814212857188.pdf\nMale Reproductive Organ Conditions-New-Erectile dysfunction secondary to PTSD-6195055014215176355.pdf"
	r, i, err := UT.VbcAIUsecase.Claude3(``, str, "")
	lib.DPrintln(r, i, err)

}

func Test_VbcAIUsecase_AssociateJotformGetDataList(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5451)
	uniqcode := tCase.CustomFields.TextValueByNameBasic(biz.FieldName_uniqcode)
	jotformSubmissions, _ := UT.JotformSubmissionUsecase.AllLatestByUniqcodeExceptIntake([]string{uniqcode})
	dataList, err := biz.AssociateJotformGetDataList(jotformSubmissions)
	lib.DPrintln(dataList)
	lib.DPrintln(err)
	//str := "Hearing Loss and Tinnitus-New-Tinnitus-6194989904217511966.pdf\nHeadaches and Migraines-Increase-Tension headache-6194996724217746641.pdf\nEar-New-Dizziness-6195001744215515843.pdf\nSinus-New-Allergic rhinitishay fever-6195006264215396564.pdf\nAnkle-New-Left ankle sprain-6195014234219988921.pdf\nKnee-Increase-Traumatic Arthritis; sp two ACL reconstructions, left knee-6195024684219386014.pdf\nKnee-Increase-Degenerative Arthritis; right knee, patellofemoral chondromalacia-6195032454212156629.pdf\nScar-Increase-Painful scars, right and left knees-6195037934213042032.pdf\nBack-New-Low back pain with radiculopathy in right lower extremity secondary to degenerative Arthritis; right knee, patellofemoral chondromalacia and traumatic Arthritis; sp two ACL reconstructions, left knee-6195044314219896456.pdf\nPost-Traumatic Stress Disorder (PTSD)--6195050814212857188.pdf\nMale Reproductive Organ Conditions-New-Erectile dysfunction secondary to PTSD-6195055014215176355.pdf"
}
