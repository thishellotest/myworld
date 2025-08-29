package biz

import (
	"archive/zip"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"io"
	"os"
	"regexp"
	"strings"
	"vbc/internal/conf"
	"vbc/lib"
)

type WordUsecase struct {
	log             *log.Helper
	CommonUsecase   *CommonUsecase
	conf            *conf.Data
	ResourceUsecase *ResourceUsecase
}

func NewWordUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ResourceUsecase *ResourceUsecase) *WordUsecase {
	uc := &WordUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		ResourceUsecase: ResourceUsecase,
	}
	return uc
}

func (c *WordUsecase) PersonalStatementsTpl() string {
	return c.ResourceUsecase.ResPath() + "/Personal_Statements_tpl.docx"
}

func (c *WordUsecase) DocEmailTpl() string {
	return c.ResourceUsecase.ResPath() + "/DocEmail_tpl.docx"
}

func (c *WordUsecase) DoPersonalStatementsWord(tCase *TData) (io.Reader, error) {
	if tCase == nil {
		return nil, errors.New("tCase is nil")
	}
	dealName := tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)

	statements := tCase.CustomFields.TextValueByNameBasic("statements")
	//claimsSupplemental := tCase.CustomFields.TextValueByNameBasic("claims_supplemental")

	cliams := FormatClaimsInfo(statements)
	//t1 := FormatClaimsInfo(claimsSupplemental)
	//cliams = append(cliams, t1...)

	//lib.DPrintln(dealName)
	//lib.DPrintln(cliams)
	return c.CreatePersonalStatementsWord(dealName, cliams)
	//if err != nil {
	//	return err
	//}

	//file, err := os.Create("/tmp/output.docx")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//
	//// 使用 io.Copy 将数据从 reader 写入文件
	//_, err = io.Copy(file, reader)

}

const (
	WordLine_Type_Normal = "normal"
	WordLine_Type_List   = "list"
)

type WordLineList []WordLine

func (c WordLineList) ToString() (r string) {
	i := 0
	for _, v := range c {
		if v.Type == WordLine_Type_Normal {
			if r == "" {
				r = v.Value
			} else {
				r += "\n\n" + v.Value
			}
		} else if v.Type == WordLine_Type_List {
			i++
			if r == "" {
				r = "    " + InterfaceToString(i) + ". " + v.Value
			} else {
				r += "\n\n    " + InterfaceToString(i) + ". " + v.Value
			}
		}
	}
	return r
}

type WordLine struct {
	Type  string
	Value string
}

var WordDocEmailTop = WordLineList{
	{Type: WordLine_Type_Normal, Value: "Subject: Several Health Issues to Discuss"},
	{Type: WordLine_Type_Normal, Value: "Dear Dr. [DOCTOR'S NAME],"},
	{Type: WordLine_Type_Normal, Value: "I hope you're doing well."},
	{Type: WordLine_Type_Normal, Value: "[OPENER 1 - Recent doctor visit] I forgot to mention several important health issues during our last visit. Is this something you can help with, or do I need to schedule an appointment?"},
	{Type: WordLine_Type_Normal, Value: "[OPENER 2 - Upcoming doctor visit] I wanted to mention these health issues prior to our appointment on [DATE]."},
	{Type: WordLine_Type_Normal, Value: "[OPENER 3 - Have not been to the doctor] I have been experiencing several health issues lately. Is this something you can help with, or do I need to schedule an appointment?"},
	{Type: WordLine_Type_Normal, Value: "Here's a rundown of my current health concerns:"},
}

var WordDocEmailBottom = WordLineList{
	{Type: WordLine_Type_Normal, Value: "All of these issues are making daily life extremely challenging. My sleep is severely disrupted, I'm frequently in pain or discomfort, and my mental health is suffering. These conditions are affecting my ability to work effectively and maintain healthy relationships with my family."},
	{Type: WordLine_Type_Normal, Value: "Is this something we can discuss at my next visit, or should I schedule a separate appointment to go over all of this?"},
	{Type: WordLine_Type_Normal, Value: "Thank you for your help and understanding."},
	{Type: WordLine_Type_Normal, Value: "Sincerely,"},
}

func (c *WordUsecase) CreateDocEmailWord(wordLineList WordLineList) (io.Reader, error) {

	doc, err := document.OpenTemplate(c.DocEmailTpl())
	if err != nil {
		return nil, err
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	//for _, s := range doc.Styles.Styles() {
	//	fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	//}
	//return

	Normal := "36"
	//SeptalLine := "33"
	ListParagraph := "35"

	for _, v := range wordLineList {
		if v.Type == WordLine_Type_Normal {
			para := doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText(v.Value)
		} else if v.Type == WordLine_Type_List {
			para := doc.AddParagraph()
			para.SetStyle(ListParagraph)
			para.AddRun().AddText(v.Value)
		}
	}

	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("Dear Dr. [DOCTOR'S NAME],")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("I hope you're doing well.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("[OPENER 1 - Recent doctor visit] I forgot to mention several important health issues during our last visit. Is this something you can help with, or do I need to schedule an appointment?")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("[OPENER 2 - Upcoming doctor visit] I wanted to mention these health issues prior to our appointment on [DATE].")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("[OPENER 3 - Have not been to the doctor] I have been experiencing several health issues lately. Is this something you can help with, or do I need to schedule an appointment?")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("Here's a rundown of my current health concerns:")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Nonintractable episodic headaches: I suffer from severe headaches 3-4 days per month, lasting up to 6 hours each. These are often accompanied by sensitivity to light and sound.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Bilateral knee pain: I experience pain and limited range of motion in both knees, with cracking and popping sounds when I move them. This makes walking and standing difficult.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("GERD: I experience reflux 3-4 times a week and regurgitation about twice a week, along with nausea and occasional vomiting.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Obstructive Sleep Apnea: Despite using a CPAP machine, I still struggle with chronic fatigue and poor quality sleep.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Depression and Anxiety: I'm experiencing worsening symptoms including intrusive thoughts, memory issues, and difficulty concentrating.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Traumatic Brain Injury effects: I continue to experience cognitive difficulties related to a TBI I suffered years ago.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(ListParagraph)
	//para.AddRun().AddText("Low Back Pain: I have persistent lower back pain with numbness, tingling, and shooting pain down both legs.")
	//
	//for k, v := range WordDocEmailBottom {
	//	para = doc.AddParagraph()
	//	para.SetStyle(Normal)
	//	para.AddRun().AddText("All of these issues are making daily life extremely challenging. My sleep is severely disrupted, I'm frequently in pain or discomfort, and my mental health is suffering. These conditions are affecting my ability to work effectively and maintain healthy relationships with my family.")
	//
	//}
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("Is this something we can discuss at my next visit, or should I schedule a separate appointment to go over all of this?")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("Thank you for your help and understanding.")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("Sincerely, ")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText(FullName)

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)

	writer := &bytes.Buffer{}
	err = doc.Save(writer)
	if err != nil {
		return nil, err
	}
	//doc.SaveToFile("/tmp/use-template2.docx")
	return writer, nil
}

// CreatePersonalStatementsWordForAiV1 statementConditionId=0时，全部； 非0时，指定
func (c *WordUsecase) CreatePersonalStatementsWordForAiV1(dealName string, statementDetail StatementDetail, statementConditionId int32) (io.Reader, error) {

	doc, err := document.OpenTemplate(c.PersonalStatementsTpl())
	if err != nil {
		return nil, err
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	//for _, s := range doc.Styles.Styles() {
	//	fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	//}
	//return

	Normal := "1"
	SeptalLine := "33"
	//ListParagraph := "34"

	para := doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText(dealName)

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	var statementBaseInfoList StatementBaseInfoList
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Full Name",
		Value: statementDetail.BaseInfo.FullName,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Unique ID",
		Value: InterfaceToString(statementDetail.CaseId),
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Branch of Service",
		Value: statementDetail.BaseInfo.BranchOfService,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Years of Service",
		Value: statementDetail.BaseInfo.YearsOfService,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Retired from service",
		Value: statementDetail.BaseInfo.RetiredFromService,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Marital Status",
		Value: statementDetail.BaseInfo.MaritalStatus,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Children",
		Value: statementDetail.BaseInfo.Children,
	})
	statementBaseInfoList = append(statementBaseInfoList, StatementBaseInfo{
		Label: "Occupation in service",
		Value: statementDetail.BaseInfo.OccupationInService,
	})

	for _, v := range statementBaseInfoList {
		//text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(text)
	}

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	for k, v := range statementDetail.Statements {

		if v.IsEmptyResult() {
			continue
		}

		if statementConditionId != 0 {
			if v.StatementCondition.StatementConditionId != statementConditionId {
				continue
			}
		}

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText("Name of Disability/Condition: " + v.StatementCondition.ConditionValue)

		//para = doc.AddParagraph()
		//para.SetStyle(Normal)
		//para.AddRun().AddText("")

		for _, v1 := range v.Rows {

			if v1.SectionType == Statemt_Section_CurrentTreatmentFacility ||
				v1.SectionType == Statemt_Section_CurrentMedication {
				para = doc.AddParagraph()
				para.SetStyle(Normal)
				para.AddRun().AddText(GetSectionTitleFromSectionType(v1.SectionType) + ": " + v1.Body)
			}
		}

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText("")

		for _, v1 := range v.Rows {

			if v1.SectionType != Statemt_Section_CurrentTreatmentFacility &&
				v1.SectionType != Statemt_Section_CurrentMedication {

				if v1.SectionType == Statemt_Section_SpecialNotes || v1.SectionType == Statemt_Section_IntroductionParagraph {

					if v1.Body != "" {
						lines := strings.Split(v1.Body, "\n")
						for _, line := range lines {
							line := strings.TrimSpace(line)
							if line == "" {
								continue
							}
							para = doc.AddParagraph()
							para.SetStyle(Normal)
							para.AddRun().AddText(line)
							para = doc.AddParagraph()
							para.SetStyle(Normal)
							para.AddRun().AddText("")
						}
					}
				} else {
					if v1.Body == "" {
						continue
					}

					para = doc.AddParagraph()
					para.SetStyle(Normal)
					para.AddRun().AddText(v1.Title + ":")
					para = doc.AddParagraph()
					para.SetStyle(Normal)
					para.AddRun().AddText("")

					if v1.Body != "" {
						lines := strings.Split(v1.Body, "\n")
						for lineKey, line := range lines {
							line := strings.TrimSpace(line)
							if line == "" {
								continue
							}
							para = doc.AddParagraph()
							para.SetStyle(Normal)
							para.AddRun().AddText(line)

							if v1.SectionType == Statemt_Section_Request && lineKey == len(lines)-1 {

							} else {
								para = doc.AddParagraph()
								para.SetStyle(Normal)
								para.AddRun().AddText("")
							}
						}
					}

				}
			}
		}

		if k != len(statementDetail.Statements)-1 && statementConditionId == 0 {
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
			para = doc.AddParagraph()
			para.SetStyle(SeptalLine)
			para.AddRun().AddText("")

			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
			//para = doc.AddParagraph()
			//para.SetStyle(Normal)
			//para.AddRun().AddText("")
		}
	}

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)

	writer := &bytes.Buffer{}
	err = doc.Save(writer)
	if err != nil {
		return nil, err
	}
	//doc.SaveToFile("/tmp/use-template2.docx")
	return writer, nil
}

func (c *WordUsecase) CreatePersonalStatementsWordForAi(dealName string, statementBaseInfoList StatementBaseInfoList, parseResults []string) (io.Reader, error) {

	doc, err := document.OpenTemplate(c.PersonalStatementsTpl())
	if err != nil {
		return nil, err
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	//for _, s := range doc.Styles.Styles() {
	//	fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	//}
	//return

	Normal := "1"
	SeptalLine := "33"
	//ListParagraph := "34"

	para := doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText(dealName)

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	for _, v := range statementBaseInfoList {
		//text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(text)
	}

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	for i, v := range parseResults {

		lines := strings.Split(v, "\n")

		for _, v1 := range lines {
			isTitle := false
			isTopTitle := false
			if strings.Index(v1, "## ") == 0 {
				isTitle = true
			}
			if strings.Index(v1, "# Name of Disability") == 0 {
				isTopTitle = true
			}
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			if isTitle {
				v1 = strings.TrimPrefix(v1, "## ")
			}
			v1 = strings.TrimPrefix(v1, "# ")
			para.AddRun().AddText(v1)
			if isTopTitle {
				para = doc.AddParagraph()
				para.SetStyle(Normal)
				para.AddRun().AddText("")
			}
		}
		if i != len(parseResults)-1 {
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
			para = doc.AddParagraph()
			para.SetStyle(SeptalLine)
			para.AddRun().AddText("")

			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
		}
	}

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)

	writer := &bytes.Buffer{}
	err = doc.Save(writer)
	if err != nil {
		return nil, err
	}
	//doc.SaveToFile("/tmp/use-template2.docx")
	return writer, nil
}

type PersonalStatementsForCreateWordVo struct {
	//DealName string
	//Conditions []string //在Client的Personal Statements_ABagarry_5373.docx没有此项
	BaseInfo []string
}

type PersonalStatementsVo struct {
	DealName   string
	Conditions []string //在Client的Personal Statements_ABagarry_5373.docx没有此项
	BaseInfo   []string
	Statements []PersonalStatementsVoStatement
}

type PersonalStatementOneVo struct {
	DealName  string
	BaseInfo  []string
	Statement PersonalStatementsVoStatement
}

func (c *PersonalStatementOneVo) ToParseUpdateStatementVo() (ParseUpdateStatementVo, error) {
	text := c.ToText()
	return ParseUpdateStatementResultSplitLine(text)
}

func (c *PersonalStatementOneVo) ToTextForAi() (text string, err error) {

	vo, err := c.ToParseUpdateStatementVo()
	if err != nil {
		return "", err
	}
	baseInfoStr := strings.Join(c.BaseInfo, "\n")
	text = baseInfoStr
	text += "\n\n" + vo.ToOneConditionText()
	return text, nil
}

func (c *PersonalStatementOneVo) ToText() (text string) {
	for _, v := range c.Statement {
		if text == "" {
			text = v
		} else {
			text += "\n" + v
		}
	}
	return text
}

type PersonalStatementsVoStatement []string

func (c PersonalStatementsVoStatement) ToText() (text string) {
	for _, v := range c {
		if text == "" {
			text = v
		} else {
			text += "\n" + v
		}
	}
	return text
}

func IsCurrentTreatment(line string) bool {
	if strings.Index(strings.ToLower(line), strings.ToLower("Current Treatment")) >= 0 {
		return true
	}
	return false
}
func IsCurrentMedication(line string) bool {
	if strings.Index(strings.ToLower(line), strings.ToLower("Current Medication")) >= 0 {
		return true
	}
	return false
}

func IsNameOfDisability(line string) bool {
	if strings.Index(strings.ToLower(line), strings.ToLower("Name of Disability")) >= 0 {
		return true
	}
	return false
}

func (c *PersonalStatementsVo) ToStatements() (listParseUpdateStatementVo ListParseUpdateStatementVo, err error) {
	//text := c.ToText()

	for _, v := range c.Statements {
		text := v.ToText()
		parseUpdateStatementVo, err := ParseUpdateStatementResultSplitLine(text)
		if err != nil {
			return nil, err
		}
		listParseUpdateStatementVo = append(listParseUpdateStatementVo, parseUpdateStatementVo)
	}
	return
}

func (c *PersonalStatementsVo) ToText() (text string) {

	text = c.DealName + "\n\n"
	//for _, v := range c.Conditions {
	//	text += v + "\n"
	//}
	//text += "\n"
	//for _, v := range c.BaseInfo {
	//	text += v + "\n"
	//}
	//text += "\n"
	for _, v := range c.Statements {
		for _, v1 := range v {
			if IsCurrentTreatment(v1) ||
				IsCurrentMedication(v1) ||
				IsNameOfDisability(v1) {
				text += v1 + "\n"
			}
		}
		text += "\n"
	}
	return text
}

func ReadDocxText(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var content string
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			bytes, _ := io.ReadAll(rc)
			type Text struct {
				Text string `xml:"t"`
			}
			var doc struct {
				Paragraphs []struct {
					Runs []Text `xml:"r"`
				} `xml:"body>p"`
			}
			xml.Unmarshal(bytes, &doc)

			for _, p := range doc.Paragraphs {
				for _, r := range p.Runs {
					content += r.Text
				}
				content += "\n"
			}
		}
	}
	return content, nil
}

func ParseUpdateStatementResultSplitLine(text string) (parseUpdateStatementVo ParseUpdateStatementVo, err error) {

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		//lib.DPrintln("line: ", line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if IsNameOfDisability(key) {
				parseUpdateStatementVo.NameOfDisability.Key = key
				parseUpdateStatementVo.NameOfDisability.Value = value
			} else if IsCurrentTreatment(key) {
				parseUpdateStatementVo.CurrentTreatment.Key = key
				parseUpdateStatementVo.CurrentTreatment.Value = value
			} else if IsCurrentMedication(key) {
				parseUpdateStatementVo.CurrentMedication.Key = key
				parseUpdateStatementVo.CurrentMedication.Value = value
			} else {
				parseUpdateStatementVo.Statements = append(parseUpdateStatementVo.Statements, line)
			}
		} else {
			parseUpdateStatementVo.Statements = append(parseUpdateStatementVo.Statements, line)
		}
	}
	if len(parseUpdateStatementVo.NameOfDisability.Key) == 0 {
		return parseUpdateStatementVo, errors.New("There was an error in parsing the Name of Disability/Condition")
	}
	if len(parseUpdateStatementVo.CurrentTreatment.Key) == 0 {
		return parseUpdateStatementVo, errors.New("There was an error in parsing the Current Treatment Facility")
	}
	if len(parseUpdateStatementVo.CurrentMedication.Key) == 0 {
		return parseUpdateStatementVo, errors.New("There was an error in parsing the Current Medication")
	}
	return
}

type ListParseUpdateStatementVo []ParseUpdateStatementVo

type ParseUpdateStatementVo struct {
	NameOfDisability struct {
		Key   string
		Value string
	}
	CurrentTreatment struct {
		Key   string
		Value string
	}
	CurrentMedication struct {
		Key   string
		Value string
	}
	Statements []string
}

func (c *ParseUpdateStatementVo) ToOneConditionText() string {
	text := fmt.Sprintf("%s: %s", c.NameOfDisability.Key, c.NameOfDisability.Value)
	text += fmt.Sprintf("\n%s: %s", c.CurrentTreatment.Key, c.CurrentTreatment.Value)
	text += fmt.Sprintf("\n%s: %s", c.CurrentMedication.Key, c.CurrentMedication.Value)
	text += "\n"
	for _, v := range c.Statements {
		text += "\n" + v
	}
	return text
}

func ParseUpdateStatementResult(parseResult string) (listParseUpdateStatementVo ListParseUpdateStatementVo, err error) {
	re := regexp.MustCompile(`(?s)(Name of Disability/Condition:.*?Current Medication:.*?(?:\n|$))`)
	matches := re.FindAllString(parseResult, -1)

	for _, match := range matches {
		//lib.DPrintln("matches:", match)
		parseUpdateStatementVo, err := ParseUpdateStatementResultSplitLine(match)
		if err != nil {
			return nil, err
		}
		listParseUpdateStatementVo = append(listParseUpdateStatementVo, parseUpdateStatementVo)
	}
	return
}

func SplitUpdatePersonalStatementsAiResult(text string) (ParseUpdateStatementVo, error) {

	lines := strings.Split(text, "\n")
	isOk := false
	destText := ""
	for _, v := range lines {
		v := strings.TrimSpace(v)
		if v != "" {
			if strings.Index(v, "Name of Disability") >= 0 {
				isOk = true
			}
			if isOk {
				if destText == "" {
					destText = v
				} else {
					destText += "\n" + v
				}
			}
		}
	}
	return ParseUpdateStatementResultSplitLine(destText)
}

func SplitPersonalStatementsString(text string) (PersonalStatementsVo, error) {

	//lib.DPrintln(text)
	//return PersonalStatementsVo{}, nil
	lines := strings.Split(text, "\n")
	var personalStatementsVo PersonalStatementsVo
	var temp []string
	for _, v := range lines {
		v := strings.TrimSpace(v)
		if v != "" {
			//if personalStatementsVo.DealName == "" {
			//	personalStatementsVo.DealName = v
			//} else {
			//if len(personalStatementsVo.Conditions) == 0 {
			//	if strings.Index(v, "Full Name") >= 0 {
			//		//aa := make([]string, len(temp))
			//		aa := append([]string(nil), temp...)
			//		//copy(aa, temp)
			//		personalStatementsVo.Conditions = aa
			//		temp = nil
			//		temp = append(temp, v)
			//	} else {
			//		temp = append(temp, v)
			//	}
			//} else

			if len(personalStatementsVo.BaseInfo) == 0 {
				if strings.Index(v, "Name of Disability") >= 0 {
					//aa := make([]string, len(temp))
					//copy(aa, temp)
					aa := append([]string(nil), temp...)
					personalStatementsVo.BaseInfo = aa
					temp = nil
					temp = append(temp, v)
				} else {
					temp = append(temp, v)
				}
			} else {
				if strings.Index(v, "Name of Disability") >= 0 {
					aa := make([]string, len(temp))
					copy(aa, temp)
					personalStatementsVo.Statements = append(personalStatementsVo.Statements, aa)
					temp = nil
					temp = append(temp, v)
				} else {
					temp = append(temp, v)
				}
			}
			//}
		}
		//lib.DPrintln(v)
	}
	aa := make([]string, len(temp))
	copy(aa, temp)
	personalStatementsVo.Statements = append(personalStatementsVo.Statements, aa)

	var err error
	if len(personalStatementsVo.Statements) == 0 {
		err = errors.New("There was an error in parsing Statements")
	} else if len(personalStatementsVo.BaseInfo) == 0 {
		err = errors.New("There was an error in parsing BaseInfo")
	}
	//else if len(personalStatementsVo.Conditions) == 0 {
	//	err = errors.New("There was an error in parsing Conditions")
	//}
	return personalStatementsVo, err
}

func (c *WordUsecase) OpenWord() {
	file := "/tmp/Personal Statements_ABagarry_5373.docx"
	doc, err := document.Open(file)
	if err != nil {
		c.log.Error("error opening document: %s", err)
		return
	}

	//批注
	for _, docfile := range doc.DocBase.ExtraFiles {
		if docfile.ZipPath != "word/comments.xml" {
			continue
		}

		file, err := os.Open(docfile.DiskPath)
		if err != nil {
			continue
		}
		defer file.Close()
		f, err := file.Stat()
		if err != nil {
			continue
		}
		size := f.Size()
		var fileinfo []byte = make([]byte, size)
		_, err = file.Read(fileinfo)
		if err != nil {
			continue
		}
		//实际应该解析<w:t>中的数据
		fmt.Println(string(fileinfo))
	}

	//书签
	for _, bookmark := range doc.Bookmarks() {
		bookname := bookmark.Name()
		if len(bookname) == 0 {
			continue
		}
		fmt.Println(bookmark.Name())
	}

	//页眉
	for _, head := range doc.Headers() {
		var text string
		for _, para := range head.Paragraphs() {
			for _, run := range para.Runs() {
				text += run.Text()
			}
		}
		if len(text) == 0 {
			continue
		}
		fmt.Println(text)
	}

	//页脚
	for _, footer := range doc.Footers() {
		for _, para := range footer.Paragraphs() {
			var text string
			for _, run := range para.Runs() {
				text += run.Text()
			}
			if len(text) == 0 {
				continue
			}
			fmt.Println(text)
		}
	}

	lib.DPrintln("doc.Paragraphs() -- :", len(doc.Paragraphs()))
	//doc.Paragraphs()得到包含文档所有的段落的切片
	for _, para := range doc.Paragraphs() {
		var text string
		//run为每个段落相同格式的文字组成的片段
		for _, run := range para.Runs() {
			text += run.Text()
		}
		if len(text) == 0 {
			continue
		}
		//打印一段
		fmt.Println(text)
	}

	//获取表格中的文本
	for _, table := range doc.Tables() {
		for _, run := range table.Rows() {
			for _, cell := range run.Cells() {
				var text string
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						//fmt.Print("\t-----------第", j, "格式片段-------------")
						text += run.Text()
					}
				}
				if len(text) == 0 {
					continue
				}
				fmt.Println(text)
			}
		}
	}

}

func (c *WordUsecase) CreateUpdatePersonalStatementsWord(personalStatementsForCreateWordVo PersonalStatementsForCreateWordVo, listParseUpdateStatementVo ListParseUpdateStatementVo) (io.Reader, error) {

	doc, err := document.OpenTemplate(c.PersonalStatementsTpl())
	if err != nil {
		return nil, err
	}

	Normal := "1"
	SeptalLine := "33"
	ListParagraph := "34"
	divide1 := "36"

	c.log.Info(ListParagraph, divide1)

	para := doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText(personalStatementsVo.DealName)
	//
	//para = doc.AddParagraph()
	//para.SetStyle(SeptalLine)
	//para.AddRun().AddText("")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("")

	//for i, v := range personalStatementsVo.Conditions {
	//	para = doc.AddParagraph()
	//	para.SetStyle(ListParagraph)
	//	para.AddRun().AddText(v)
	//
	//	if i != len(personalStatementsVo.Conditions)-1 {
	//		para = doc.AddParagraph()
	//		para.SetStyle(Normal)
	//		para.AddRun().AddText("")
	//	}
	//}

	//para = doc.AddParagraph()
	//para.SetStyle(SeptalLine)
	//para.AddRun().AddText("")
	//
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("")

	for _, v := range personalStatementsForCreateWordVo.BaseInfo {
		//text := fmt.Sprintf("• %s: %s", v.Label, v.Value)
		text := v
		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(text)
	}

	//para = doc.AddParagraph()
	//para.SetStyle(SeptalLine)
	//para.AddRun().AddText("")

	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("")
	//para = doc.AddParagraph()
	//para.SetStyle(Normal)
	//para.AddRun().AddText("")

	for i, v := range listParseUpdateStatementVo {

		para = doc.AddParagraph()
		para.SetStyle(SeptalLine)
		para.AddRun().AddText("")
		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText("")

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(fmt.Sprintf("%s: %s", v.NameOfDisability.Key, v.NameOfDisability.Value))

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(fmt.Sprintf("%s: %s", v.CurrentTreatment.Key, v.CurrentTreatment.Value))

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText(fmt.Sprintf("%s: %s", v.CurrentMedication.Key, v.CurrentMedication.Value))

		para = doc.AddParagraph()
		para.SetStyle(Normal)
		para.AddRun().AddText("")

		for k1, v1 := range v.Statements {
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText(v1)

			if k1 != len(v.Statements)-1 {
				para = doc.AddParagraph()
				para.SetStyle(Normal)
				para.AddRun().AddText("")
			}
		}
		if i != len(listParseUpdateStatementVo)-1 {
			//para = doc.AddParagraph()
			//para.SetStyle(Normal)
			//para.AddRun().AddText("")
			//para = doc.AddParagraph()
			//para.SetStyle(Normal)
			//para.AddRun().AddText("")
		}
	}

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)

	writer := &bytes.Buffer{}
	err = doc.Save(writer)
	if err != nil {
		return nil, err
	}
	//doc.SaveToFile("/tmp/use-template2.docx")
	return writer, nil

}

func (c *WordUsecase) CreatePersonalStatementsWord(dealName string, cliams []string) (io.Reader, error) {

	doc, err := document.OpenTemplate(c.PersonalStatementsTpl())
	if err != nil {
		return nil, err
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	//for _, s := range doc.Styles.Styles() {
	//	fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	//}
	//return

	Normal := "1"
	SeptalLine := "33"
	ListParagraph := "34"

	para := doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText(dealName)

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	for i, v := range cliams {
		para = doc.AddParagraph()
		para.SetStyle(ListParagraph)
		para.AddRun().AddText(v)

		if i != len(cliams)-1 {
			para = doc.AddParagraph()
			para.SetStyle(Normal)
			para.AddRun().AddText("")
		}
	}

	para = doc.AddParagraph()
	para.SetStyle(SeptalLine)
	para.AddRun().AddText("")

	para = doc.AddParagraph()
	para.SetStyle(Normal)
	para.AddRun().AddText("")

	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch,
		measurement.Inch, measurement.Inch, measurement.Inch, 0)

	writer := &bytes.Buffer{}
	err = doc.Save(writer)
	if err != nil {
		return nil, err
	}
	//doc.SaveToFile("/tmp/use-template2.docx")
	return writer, nil
}

func FormatClaimsInfo(val string) (r []string) {
	res := strings.Split(val, "\n")
	for _, v := range res {
		v := strings.TrimSpace(v)
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}
