package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

func Test_JotformbuzUsecase_HandleSubmission2(t *testing.T) {
	err := UT.JotformbuzUsecase.HandleSubmission("6206799971121513298", "", "")
	lib.DPrintln(err)
}

func Test_JotformbuzUsecase_HandleSubmission(t *testing.T) {
	err := UT.JotformbuzUsecase.HandleSubmission("6277002835013497426", "", "251865711410149")
	lib.DPrintln(err)
}

func Test_JotformbuzUsecase_HandleSubmission1(t *testing.T) {

	res, _ := UT.JotformSubmissionUsecase.AllByCond(builder.In("submission_id", "6037031974216948122", "6036971094217009889"))
	for _, v := range res {
		notes := lib.ToTypeMapByString(v.Notes)

		lib.DPrintln("", notes.GetString("content.answers"))
	}
}

func Test_JotformbuzUsecase_HandleQuestionnairesJotform(t *testing.T) {
	UT.JotformbuzUsecase.HandleQuestionnairesJotformHistory()
}

func Test_JotformbuzUsecase_GetSubmissionInfo(t *testing.T) {
	str := `{"content":{"answers":{"1":{"name":"heading","order":"1","text":"Eye Questionnaire","type":"control_head"},"11":{"name":"brieflyDescribe","order":"13","text":"Briefly describe where and how you first encountered this condition: ","type":"control_textarea"},"111":{"answer":["Worse in frequency","Worse in severity"],"name":"howWould","order":"22","prettyFormat":"Worse in frequency; Worse in severity","text":"How would you characterize your condition since your first encounter (check all that apply)?","type":"control_checkbox"},"145":{"name":"doYou145","order":"24","text":"Do you currently take medications for this condition?","type":"control_radio"},"198":{"name":"ifYes198","order":"28","text":"If yes, specify the condition causing incapacitating episodes:","type":"control_textarea"},"199":{"name":"veryImportant","order":"23","text":"VERY IMPORTANT: History of Surgery – Provide type of surgery, indicate if VA did surgery or private doctor, and date:","type":"control_textarea"},"2":{"name":"submit2","order":"41","text":"Submit","type":"control_button"},"205":{"answer":"claiming Value","name":"whatIs","order":"7","text":"What is the exact condition you are claiming? ","type":"control_textarea"},"206":{"name":"whatYear","order":"10","text":"What year did your condition begin?","type":"control_textbox"},"208":{"name":"doesYour208","order":"37","text":"Does this condition affect your ability to work?","type":"control_radio"},"209":{"name":"didThe209","order":"8","text":"Did the condition begin during active duty service?","type":"control_radio"},"210":{"name":"doYou","order":"26","text":"Do you have scarring or disfigurement attributable to this specified condition?","type":"control_radio"},"211":{"name":"duringThe","order":"27","text":"During the past 12 months, do you have any incapacitating episodes attributable to this specified condition?","type":"control_radio"},"212":{"name":"howDoes","order":"32","text":"How does this condition affect your daily life? ","type":"control_textarea"},"213":{"name":"indicateThe","order":"29","text":"Indicate the number of documented medical visits for treatment of eye condition over the past 12 months:","type":"control_radio"},"214":{"name":"indicateThe214","order":"30","text":"Indicate the type of intervention that occurred during the incapacitating episode. (Check all that apply)","type":"control_checkbox"},"215":{"name":"ifYes215","order":"31","text":"If yes to any of the above, please list name of medication or describe:","type":"control_textarea"},"216":{"answer":["New condition found in Service Treatment Records (STR)"],"name":"checkAll","order":"5","prettyFormat":"New condition found in Service Treatment Records (STR)","text":"Check all that apply:","type":"control_checkbox"},"217":{"cfname":"Date Picker","name":"typeA","order":"20","selectedField":"52934dbf3be147110a000030","static":"No","text":"When did you most recently see a doctor regarding this medical issue?","type":"control_widget"},"218":{"name":"howDoes218","order":"33","text":"How does this condition affect your daily life?","type":"control_checkbox"},"219":{"name":"howDoes219","order":"36","text":"How does this condition affect your sleep?","type":"control_checkbox"},"22":{"name":"ifYes","order":"17","text":"If yes, what did they do for you? ","type":"control_textarea"},"220":{"name":"doesThis","order":"35","text":"Does this condition cause sleep disturbance?","type":"control_radio"},"221":{"name":"howDoes221","order":"38","text":"How does this condition affect your ability to work?","type":"control_checkbox"},"222":{"name":"whatUnit","order":"11","text":"What unit were you assigned to?","type":"control_textbox"},"223":{"name":"whereWere","order":"12","text":"Where were you located?","type":"control_textbox"},"224":{"name":"ifYes224","order":"15","text":"If yes, what did they do for you?","type":"control_checkbox"},"225":{"name":"ifNo","order":"16","text":"If no, why didn't you seek treatment?","type":"control_checkbox"},"226":{"name":"ifCondition","order":"6","text":"If condition found in STR, what is the name of the condition in STR?","type":"control_textarea"},"24":{"name":"didYou24","order":"14","text":"Did you seek treatment when the condition began?","type":"control_radio"},"26":{"name":"ifYes26","order":"25","text":"If yes, please list all medications.","type":"control_textarea"},"32":{"name":"doYou32","order":"21","text":"Do you currently have an eye condition?","type":"control_radio"},"34":{"answer":"Yes","name":"haveYou","order":"34","text":"Have you lost interest in activities you once enjoyed?","type":"control_radio"},"35":{"name":"ifYes35","order":"39","text":"If yes, explain in as much detail as possible: ","type":"control_textarea"},"39":{"answer":{"first":"Test1","last":"TestL"},"name":"fullName","order":"2","prettyFormat":"Test1 TestL","sublabels":"{\"prefix\":\"Prefix\",\"first\":\"First Name\",\"middle\":\"Middle Name\",\"last\":\"Last Name\",\"suffix\":\"Suffix\"}","text":"Full Name","type":"control_fullname"},"40":{"name":"pageBreak","order":"18","text":"Page Break","type":"control_pagebreak"},"43":{"name":"currentSituation","order":"19","text":"Current Situation","type":"control_head"},"44":{"name":"additionalComments","order":"40","text":"Additional comments:","type":"control_textarea"},"45":{"answer":"5511","name":"uniqueId","order":"3","text":"Unique ID","type":"control_textbox"},"52":{"answer":"New","name":"isThis","order":"4","text":"Is this for a new condition or increase?","type":"control_radio"},"88":{"name":"whichSide","order":"9","text":"Which side is affected? (skip if not needed)","type":"control_radio"}},"created_at":"2025-03-13 06:30:18","flag":"0","form_id":"241499180510152","id":"6176710185019648226","ip":"46.3.240.105","new":"1","notes":"","status":"ACTIVE","updated_at":null},"duration":"195.61ms","info":null,"limit-left":99997,"message":"success","responseCode":200}`
	typeMap := lib.ToTypeMapByString(str)
	lib.DPrintln("typeMap:", typeMap)
	formId, a, f, c, err := biz.GetSubmissionInfo(typeMap)
	lib.DPrintln(formId, a, f, c, err)
}
