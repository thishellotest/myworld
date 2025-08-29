package tests

import (
	"context"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_Awsclaude3Usecase_Ask_test(t *testing.T) {
	a, _, err := UT.Awsclaude3Usecase.Ask(context.TODO(), "", "hi", "Test", "")
	lib.DPrintln(a)
	lib.DPrintln(err)
}
func Test_Awsclaude3Usecase_Ask(t *testing.T) {

	systemConfig := `# Role
You are a comprehensive medical expert proficient in analyzing complex medical data, capable of extracting valuable information from vast datasets. You are also adept at efficiently extracting and analyzing all disease names, providing comprehensive medical information to patients and healthcare professionals.

## Skills
- Primary Care: Family medicine, general diagnostics, and health check-ups
- Internal Medicine: Cardiovascular diseases, gastrointestinal disorders, endocrine, and metabolic diseases
- Surgery: General surgery, minimally invasive surgery, and trauma surgery
- Obstetrics and Gynecology: Gynecological diseases, prenatal and postnatal care, reproductive health
- Dermatology: Skin diseases, skincare, and cosmetic treatments
- Cardiology: Coronary artery disease, hypertension, and cardiac rehabilitation
- Neurology: Neurological disorders, stroke, and epilepsy
- Emergency Medicine: Emergency treatment, trauma care, and urgent surgery
- Mental Health: Psychological counseling, mental disorders, and mental health management
- Data Analysis: Medical data mining, disease name extraction, and medical information organization

## Action
- Based on the patient's description, extract all disease names and output them in JSON format

## Constraints
- The case of the found disease name must match exactly the content provided
- Ensure your result contains only one valid JSON format

## Format
- The corresponding JSON key is: Conditions`
	text := `.\n' NAVMED 6150 3 (REV 7.72) FRONT (MODIFIED)\nHEALTH RECORD\nSICK CALL TREATMENT RECORD\nDATE\nNAME OF TREATING FACILITY. COMPLAINT, TREATMENT ADMINISTERED. SIGNATURE AND GRADE/RATE OF PERSON ADMINISTERING TREATMENT\n4CCTTL\nMISS DETROIT (ACC. 4)\nHEADACHE,RUNNY NOSE SORE THROAT REYAM TeP. 99-2; THROAT CLEAR -\nRX; ACTIFED 2 T TIO - TyLENOL #3X TI GID .\n54 11 acum 5N\n.\n-.\n22 oct 74\nUSS Detail AOE 4\n/0 Nausea, vanutri\n.. 5. Vanuitil x q today has felt Nauseated for the past week aldo loss of Appetite for\n. : the past week. It was Recruit at San Diego\nrate it Callers # 5 + 8 dienas\nAugust -left there on Aug 12 1974 + arrived\nNon t outbreak of barnabitis. Has not been -\ntreated for marasite infection since lite Jal\nNouses queaus to coincide with rough weather\n0- Tem. 97: No salerad una interes\n. Abd - No larato galenomical BSN. No mais tardemen. V/A - WN2\n. (A) Prote motion unkness. v. wird gastroenteritis.\n(p) Compagne tale now + 6h. #3\nDramamine # 15 min pausea .\nE-a Castle Stue\nSEX\nMALE Malayer\nRACE\nGRADE. RATING, OR POSITION\nORGANIZATION UNIT\nCOMPONENT OR BRANCH USN\nSERVICE. DEPT. OR AGENCY\nDOD\nSR\nNTC SDIEGO\nPATIENT 'S LAST NAME . FIRST NAME . MIDDLE NAME\nFANQUILUT REMIGIO DIMALANTA 75 C168\nDATE OF BIRTH ( DAY-MONTH . YEAR)\nIDENTIFICATION NO\n01-21-49\n20/ 571299855/\n21Jan44\nSICK CALL TREATMENT RECORD NAVMED 6150/3\n3\nIME\n(3F)`

	res, err := UT.Awsclaude3Usecase.GetContentByAsk(context.TODO(), systemConfig, text)
	lib.DPrintln("+++", res, "___")
	lib.DPrintln(err)
}

func Test_Awsclaude3Usecase_GenStatement(t *testing.T) {

	//systemConfig, _ := UT.LongMapUsecase.GetForString("prompt1_0")
	//
	//a, _ := UT.AiPromptUsecase.GetByPromptKey("prompt3_0")
	//systemConfig = a.Prompt
	promptKey := "prompt3_0"
	//promptKey = "prompt1_0"

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	statementCondition := lib.StringToTDef(`{"origin_value":"10 - Tinnitus (new)","condition_value":"Tinnitus","front_value":"10 -","behind_value":"(new)"}`, biz.StatementCondition{})
	res, _, aiResultId, err := UT.AiTaskbuzUsecase.GenStatement(context.TODO(), tCase, statementCondition, "", promptKey, "")
	lib.DPrintln("+++", res, "___")
	lib.DPrintln(err, aiResultId)
}

func Test_Awsclaude3Usecase_GenStatementTest(t *testing.T) {

	text := ``
	//systemConfig, _ := UT.LongMapUsecase.GetForString("prompt1_0")
	//
	//a, _ := UT.AiPromptUsecase.GetByPromptKey("prompt3_0")
	//systemConfig = a.Prompt
	promptKey := "prompt3_0"
	promptKey = "prompt1_0"
	res, aiResultId, err := UT.Awsclaude3Usecase.GenStatementTest(context.TODO(), text, promptKey)
	lib.DPrintln("+++", res, "___")
	lib.DPrintln(err, aiResultId)
}
