package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/internal/biz"
	"vbc/lib"
	"vbc/lib/builder"
)

/*
https://crm.zoho.com/crm/org847391426/tab/Potentials/6159272000005460893
Timothy Fortson-10#5182

Service:
1996-2004

Current:
10 - Bilateral tinnitus
0 - Bilateral sensorineural hearing loss

New - Online:
70 - Mood Disorder secondary to bilateral tinnitus and bilateral sensorineural hearing loss (new)

50 - Migraine headaches secondary to tinnitus (opinion)

30 - GERD secondary to chronic NSAID use related to right knee and bilateral ankle sprains (opinion)

30 - Vertigo secondary to bilateral tinnitus and bilateral sensorineural hearing loss (opinion)

20* - Back pain secondary to right knee medial meniscal tear and left ankle inversion sprain (opinion)

New - Supplemental:
20* - Right knee medial meniscal tear with limitation of flexion and extension (str, opinion, previous denial)

20* - Left knee pain with limitation of flexion and extension secondary to right knee medial meniscal tear (opinion, previous denial)

10 - Left ankle inversion sprain  (str, opinion, previous denial)

10 - Right ankle sprain  (str, opinion, previous denial)
*/
func Test_AiTaskUsecase_CreateTask(t *testing.T) {
	entity, err := UT.AiTaskUsecase.CreateTask(biz.AiTaskFromType_statement,
		"70 - Mood Disorder secondary to bilateral tinnitus and bilateral sensorineural hearing loss (new)",
		5182, 0, biz.AiTaskInputStatement{
			SubmissionId: []string{
				"6036951604218420113", "6036973064211091711"},
			Condition: "Mood Disorder secondary to bilateral tinnitus and bilateral sensorineural hearing loss (new)",
		}, "", nil, nil)
	lib.DPrintln(entity, err)
}

func Test_AiTaskUsecase_CreateGenerateDocEmail(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5431)
	entity, err := UT.AiTaskUsecase.CreateGenerateDocEmail(tCase)
	lib.DPrintln(entity)
	lib.DPrintln(err)
}

func Test_AiTaskJobUsecase_RunHandleCustomJob(t *testing.T) {
	var wait sync.WaitGroup
	err := UT.AiTaskJobUsecase.RunHandleCustomJob(context.TODO(), 1, 3*time.Second,
		UT.AiTaskJobUsecase.WaitingTasks,
		UT.AiTaskJobUsecase.Handle)
	if err != nil {
		panic(err)
	}
	wait.Add(1)
	wait.Wait()
}

func Test_AiTaskJobUsecase_DoHandleVeteranSummary(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5662)
	aiResultId, err := UT.AiTaskbuzUsecase.DoHandleVeteranSummary(context.TODO(), tCase)
	lib.DPrintln(aiResultId, err)
	a, err := UT.AiResultUsecase.GetByCond(builder.Eq{"id": aiResultId})
	lib.DPrintln(err)
	if a != nil {
		lib.DPrintln(a.ParseResult)
	}
}
