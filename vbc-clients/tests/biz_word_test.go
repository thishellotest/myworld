package tests

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_WordUsecase_CreatePersonalStatementsWordForAi(t *testing.T) {
	dealName := "MarkDean Ronduen-30#5431"
	parseResults := []string{`# Name of Disability/Condition: Tinnitus
Current Treatment Facility: None provided
Current Medication: None provided

I am respectfully requesting Veteran Affairs benefits for my condition of Tinnitus. I served in the United States Navy from April 2022 to June 2022. During my time in service, I developed this condition that continues to affect my daily life and ability to work.

SERVICE CONNECTION: My treatment records reflect that I have a diagnosis of Tinnitus. All legal requirements for establishing service connection for Tinnitus have been met; service connection for such disease is warranted.

## Onset and Service Connection:
My tinnitus began during my time in Boot Camp with the Navy in 2022. While in boot camp, I was exposed to loud shouting in my ears from the Drill Instructors who were only inches away from my ears. I was also exposed to the loud sounds of guns being shot on the firing range. We were provided foam ear plugs but they did not seal well in my ears and provided inadequate protection. I did not have tinnitus before entering military service. The condition only started after going to the firing range in Boot Camp. The ringing began immediately after exposure to these loud noises, and I have experienced it continuously since that time.

## Current Symptoms, Severity and Frequency:
Since my initial exposure to loud noise during my military service, my condition has worsened both in frequency and severity. I now experience tinnitus daily, with the ringing lasting all day. There is no relief from this constant ringing. The sound is often so loud that it interferes with my ability to hear conversations and focus on tasks at hand. On particularly bad days, the ringing is so intense that it causes me to develop headaches, which further impacts my ability to function normally.

## Medication:
I do not currently take any prescribed medication specifically for my tinnitus, as I have been told there is no medication that can eliminate the ringing. I sometimes take over-the-counter pain relievers when the tinnitus triggers headaches, but these provide only minimal and temporary relief.

## Impact on Daily Life:
The ringing in my ears makes it difficult to hear others when I am trying to speak with them. I have to constantly ask them to repeat themselves, which is frustrating and embarrassing. After asking them several times to repeat themselves, I sometimes give up on the conversation altogether. Going out with family to dinner in a busy restaurant causes me anxiety because I know I will have difficulty communicating with everyone. The constant ringing also affects my ability to enjoy quiet activities like reading or watching television. The persistent noise in my ears is mentally exhausting and often leaves me irritable and withdrawn from social situations.

## Professional Impact:
The ringing in my ears makes it nearly impossible to hear the alarms or beeping on heavy equipment like a forklift. This creates a safety hazard in many work environments. The ringing in my ears makes it difficult to concentrate and focus on tasks, reducing my productivity and job performance. In meetings or training sessions, I often miss important information due to the interference from the tinnitus. This has limited my employment options and made it challenging to maintain steady employment.

## Nexus Between Service and Current Condition:
There is a clear connection between my current tinnitus condition and my military service. Prior to entering the Navy, I had no issues with ringing in my ears. The onset of my tinnitus directly coincided with my exposure to loud noises during Boot Camp, specifically the firearms training on the range and the constant loud yelling from drill instructors. The condition began during my service and has persisted and worsened since that time.

## Request:
I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.`}

	parseResults = append(parseResults, parseResults[0])

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5004)
	tClient, _, _ := UT.DataComboUsecase.Client(tCase.CustomFields.TextValueByNameBasic(biz.FieldName_client_gid))
	r, err := UT.StatementUsecase.GetStatementBaseInfo(tCase, tClient)
	if err != nil {
		panic(err)
	}

	aa, err := UT.WordUsecase.CreatePersonalStatementsWordForAi(dealName, r, parseResults)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("test2.docx")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, aa)
	if err != nil {
		panic(err)
	}
}

func Test_WordUsecase_CreateDocEmailWord(t *testing.T) {
	//a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5431)
	//lib.DPrintln(a)

	var wordLineList biz.WordLineList
	wordLineList = append(wordLineList, biz.WordDocEmailTop...)

	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Nonintractable episodic headaches: I suffer from severe headaches 3-4 days per month, lasting up to 6 hours each. These are often accompanied by sensitivity to light and sound.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Bilateral knee pain: I experience pain and limited range of motion in both knees, with cracking and popping sounds when I move them. This makes walking and standing difficult.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "GERD: I experience reflux 3-4 times a week and regurgitation about twice a week, along with nausea and occasional vomiting.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Obstructive Sleep Apnea: Despite using a CPAP machine, I still struggle with chronic fatigue and poor quality sleep.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Depression and Anxiety: I'm experiencing worsening symptoms including intrusive thoughts, memory issues, and difficulty concentrating.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Traumatic Brain Injury effects: I continue to experience cognitive difficulties related to a TBI I suffered years ago.",
	})
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_List,
		Value: "Low Back Pain: I have persistent lower back pain with numbness, tingling, and shooting pain down both legs.",
	})

	wordLineList = append(wordLineList, biz.WordDocEmailBottom...)
	wordLineList = append(wordLineList, biz.WordLine{
		Type:  biz.WordLine_Type_Normal,
		Value: "James Stuart",
	})

	lib.DPrintln(wordLineList)
	aa, err := UT.WordUsecase.CreateDocEmailWord(wordLineList)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("test2.docx")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, aa)
	if err != nil {
		panic(err)
	}
}

func Test_WordUsecase_DoPersonalStatementsWord(t *testing.T) {
	a, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5431)
	aa, err := UT.WordUsecase.DoPersonalStatementsWord(a)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("test.docx")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, aa)
	if err != nil {
		panic(err)
	}
}

func Test_WordUsecase_OpenWord(t *testing.T) {
	UT.WordUsecase.OpenWord()
}

func Test_SplitPersonalStatementsString(t *testing.T) {
	a, err := biz.SplitPersonalStatementsString(PersonalStatementsString)

	text := a.ToText()
	lib.DPrintln(text)
	return
	lib.DPrintln(err)
	for _, v := range a.Statements {
		for _, v1 := range v {
			if strings.Index(v1, "Name of Disability") >= 0 {
				lib.DPrintln(v1)
			}
		}
	}
}

func Test_aaa222(t *testing.T) {

	// 正则：从 "Name of Disability/Condition:" 开始，到 "Current Medication:" 后面一行结束
	re := regexp.MustCompile(`(?s)(Name of Disability/Condition:.*?Current Medication:.*?(?:\n|$))`)

	matches := re.FindAllString(PersonalStatementsStringAiResult, -1)

	for i, match := range matches {
		fmt.Printf("Match %d:\n%s\n", i+1, match)
	}
}

func Test_SplitUpdatePersonalStatementsAiResult(t *testing.T) {
	str := `I've processed the update for Alexander Bagarry's VA statement. Below is the updated statement with the requested changes:

• Full Name: Alexander Bagarry
• Unique ID: 5373
• Branch of Service: Navy
• Years of Service: 1994-1998
• Retired from service: No
• Deployments: Persian Gulf (USS Vandergrift FFG-48 1995, USS Shiloh CG-67 1997), Operation Desert Strike
• Marital Status: Married (2007-Present)
• Children: 2 young children
• Occupation in service: Torpedoman's Mate (TM)

Name of Disability/Condition: Hemorrhoids
Current Treatment Facility: Sharp
Current Medication: Linzess 72 mcg, Infrared coagulation

SERVICE CONNECTION: My service treatment records reflect that I suffered from External Hemorrhoids while on active duty. All legal requirements for establishing service connection for Hemorrhoids has been met; service connection for such disease is warranted.
I am respectfully requesting Veteran Affairs benefits for my condition of hemorrhoids. This condition has been a significant challenge in my life since its onset during my active-duty service in 1996. Throughout my service in the Navy as a Torpedoman's Mate from 1994 to 1998, I faced circumstances that have contributed to this condition. The effects of hemorrhoids have become increasingly debilitating, affecting not just my physical well-being but also my ability to maintain a stable and fulfilling personal and professional life.
While deployed aboard the USS Vandergrift to the Persian Gulf and later during a port call in Hong Kong in 1996, I ate local food that I believe was contaminated. I became severely ill for several days, experiencing vomiting, extreme diarrhea, and straining to defecate. I developed pain in my rectum and went to medical, where they documented that I had hemorrhoids. I did not have hemorrhoids before this deployment, and since then, my condition has progressively worsened.
Currently, my hemorrhoids have worsened both in frequency and severity since their onset. My symptoms are moderate in nature, characterized by large or thrombotic, irreducible hemorrhoids with excessive redundant tissue and frequent recurrences. I experience persistent bleeding along with severe pain, swelling, and itching. I have undergone radiation treatment to remove them, but the surgeon would not perform surgery due to their large size.
This condition severely affects my daily life. I experience chronic pain and discomfort that affects my overall well-being, and I have difficulty sitting for extended periods. The persistent symptoms have caused me to lose interest in activities I once enjoyed.
My sleep is frequently disturbed due to discomfort and pain, leading to fatigue and irritability during the day. This chronic sleep disturbance compounds the challenges I face in my daily life.
The condition has significantly impacted my ability to work. I have difficulty with jobs requiring prolonged sitting and physical exertion. These limitations affect my productivity and career opportunities.
I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life.`
	a, err := biz.SplitUpdatePersonalStatementsAiResult(str)
	lib.DPrintln(err)
	lib.DPrintln(a)
}
