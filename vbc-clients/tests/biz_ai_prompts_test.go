package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_AiPromptUsecase_GetAiInfoByPromptKey(t *testing.T) {

	a, b, err := UT.AiPromptUsecase.GetAiInfoByPromptKey("multiline_condition_parser", lib.TypeMap{
		"text": "Right shoulder strain with nerve damage and limitation of flexion and abduction.\n20* - Right ankle with limitation of plantar flexion and dorsiflexion secondary to Infrapatellar tendinitis right knee with limitation of flexion and extension (opinion) (BVA)",
	})
	lib.DPrintln(a)
	lib.DPrintln(b)
	lib.DPrintln(err)
}

func Test_AiPromptUsecase_GetAiInfoByPromptKey1(t *testing.T) {

	a, b, err := UT.AiPromptUsecase.GetAiInfoByPromptKey(biz.Prompt_associate_jotform, lib.TypeMap{
		"condition": "Back pain secondary to bilateral pes planus",
		"data_list": "Hearing Loss and Tinnitus-New-Tinnitus-6169292782894013022.pdf\nHypertension-New-Hypertension-6169303632898518752.pdf\nFoot--Bilateral pes planus aggravated by service-6169310762892987299.pdf\nBack-New-Back pain secondary to bilateral pes planus-6169316572895335368.pdf\nMental Disorders-Increase-Insomnia Disorder-6169340152892747999.pdf\nMale Reproductive Organ Conditions-New-Erectile  dysfunction secondary to insomnia disorder-6169352952897678956.pdf",
	})
	lib.DPrintln(a)
	lib.DPrintln(b)
	lib.DPrintln(err)
}
