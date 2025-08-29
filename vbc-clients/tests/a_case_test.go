package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/to"
)

func Test_aa1(t *testing.T) {
	tTpl, err := UT.TUsecase.Data(biz.Kind_email_tpls, Eq{"tpl": "PersonalStatementsReadyforYourReview", "sub_id": 0})
	if err != nil {
		panic(err)
	}
	tTpl.CustomFields.SetTextValueByName("body", to.Ptr("sss"))
	cc := tTpl.CustomFields.TextValueByNameBasic("body")
	lib.DPrintln(cc)
}

func Test_OCR_JOB_RUN(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(1)
	ctx := context.TODO()

	err := UT.HaReportTaskJobUsecase.RunHandleCustomJob(ctx, 2, 5*time.Second,
		UT.HaReportTaskJobUsecase.WaitingTasks,
		UT.HaReportTaskJobUsecase.Handle)
	if err != nil {
		panic(err)
	}
	err = UT.BlobJobUsecase.RunHandleCustomJob(ctx, 2, 5*time.Second,
		UT.BlobJobUsecase.WaitingTasks,
		UT.BlobJobUsecase.Handle)
	if err != nil {
		panic(err)
	}

	err = UT.BlobSliceJobUsecase.RunHandleCustomJob(ctx, 2, 5*time.Second,
		UT.BlobSliceJobUsecase.WaitingTasks,
		UT.BlobSliceJobUsecase.Handle)
	if err != nil {
		panic(err)
	}
	wait.Wait()
}

func Test_StreetClean(t *testing.T) {
	street := "918 Dr. Martin Luther      King Jr    Blvd,   Suite 100,    Unit 3107,aa,"
	lib.DPrintln(street)
	aa := biz.StreetClean(street)
	lib.DPrintln(string(aa))
}

func Test_StringToStandardHeaderRevisionAiResult(t *testing.T) {
	str := `{
  "Special Notes": "SERVICE CONNECTION: My treatment records reflect that I have a diagnosis of Hypertension, a potentially devastating condition that threatens my cardiovascular health daily. Additionally, I served in Southwest Asia during the Persian Gulf War Era, enduring harsh environmental exposures and toxic hazards, and I am entitled to the application of presumptive provisions under the PACT Act. All legal requirements for establishing service connection for Hypertension have been met; service connection for such disease is warranted to address this life-altering condition that continues to impact every aspect of my health and well-being.",
  
  "Introduction": "I am respectfully requesting Veteran Affairs benefits for my condition of Hypertension which began during my service.\n\nI served in the United States Air Force from 2002 to 2006 as a Occupation 01. During this service period, I developed Hypertension that continues to affect my daily life and ability to work.",
  
  "Onset and Service Connection": "My Hypertension condition first developed during my active duty service. During regular medical check-ups, my blood pressure readings were consistently elevated, leading to my initial diagnosis.\n\nThe high-stress environment of military service, combined with irregular sleep schedules, contributed significantly to the development of my condition. While serving, I noticed symptoms including headaches, dizziness, and occasional shortness of breath, which medical personnel attributed to my elevated blood pressure.",
  
  "Current Symptoms Severity and Frequency": "My Hypertension has worsened in both frequency and severity since my initial diagnosis. I experience persistent headaches, dizziness, and fatigue that affect my daily functioning. These symptoms occur multiple times per week and have increased in intensity over time.\n\nI often feel lightheaded when standing up quickly, and experience episodes of blurred vision when my blood pressure spikes. During particularly stressful periods, I can physically feel my heart racing and pounding in my chest, which causes significant anxiety and further elevates my blood pressure in a vicious cycle.",
  
  "Medication": "I am currently taking Amlodipine to manage my Hypertension. Despite following my medication regimen strictly, I continue to experience breakthrough symptoms that impact my quality of life. My medication helps keep my blood pressure from reaching dangerous levels, but does not fully alleviate my symptoms or prevent flare-ups during periods of stress or physical exertion.",
  
  "Impact on Daily Life": "Hypertension has significantly impacted my daily life.\n\nI have lost interest in activities I once enjoyed due to persistent fatigue and concern about triggering blood pressure spikes. Physical activities that I previously engaged in without issue now cause me to feel unwell and require extended recovery time.\n\nI must constantly monitor my diet, avoiding many foods I once enjoyed. My social life has diminished as I often need to cancel plans when experiencing severe symptoms. Family members have expressed concern about my health, adding emotional stress to my physical condition.",
  
  "Professional Impact": "My Hypertension condition has negatively impacted my professional life. I experience difficulty concentrating during periods of elevated blood pressure, which affects my job performance. I have had to take sick days when symptoms are particularly severe, impacting my reliability at work. The fatigue associated with my condition makes it challenging to maintain productivity throughout the workday. My employer has noted concerns about my health, creating additional stress about job security. The combination of physical symptoms and medication side effects has limited my career advancement opportunities.",
  
  "Nexus Between Service and Current Condition": "The onset of my Hypertension occurred during my military service, where the combination of high-stress environments, irregular schedules, and demanding physical requirements created conditions conducive to developing Hypertension. The initial symptoms I experienced during service have persisted and worsened since separation. Military medical records document my elevated blood pressure readings during service, establishing a clear connection between my current condition and my time in the military. The progression of my symptoms follows the typical pattern of service-connected Hypertension, with gradual worsening over time despite treatment.",
  
  "Request": "I respectfully request your thorough consideration of my application for benefits to help me access the necessary medical care, treatments, and resources to effectively manage my condition. Your support and understanding are crucial as I strive to regain a sense of normalcy, improve my overall well-being, and rebuild my ability to participate fully in both my personal and professional life."
}`
	cc := biz.StringToStandardHeaderRevisionAiResult(str)
	lib.DPrintln(cc)
}
