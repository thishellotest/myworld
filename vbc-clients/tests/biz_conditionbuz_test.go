package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ConditionbuzUsecase_HandleAllCondition(t *testing.T) {
	err := UT.ConditionbuzUsecase.HandleAllCondition()
	lib.DPrintln(err)
}

func Test_ConditionbuzUsecase_HandleParseConditionResultFromAi(t *testing.T) {
	str := `[{"PrimaryCondition":"Left ankle pain","DirectSecondaryConditions":[],"AggravationConditions":[],"SourceData":"Left ankle pain (str, opinion)"},{"PrimaryCondition":"Depression","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"70 - Depression secondary to tinnitus and hearing loss"},{"PrimaryCondition":"Vertigo","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"60 - Vertigo secondary to tinnitus and hearing loss (opinion)"},{"PrimaryCondition":"Headaches","DirectSecondaryConditions":["tinnitus"],"AggravationConditions":[],"SourceData":"50 - Headaches secondary to tinnitus (opinion)"},{"PrimaryCondition":"Sleep apnea","DirectSecondaryConditions":["tinnitus","hearing loss"],"AggravationConditions":[],"SourceData":"50 - Sleep apnea secondary to tinnitus and hearing loss (opinion)"},{"PrimaryCondition":"Vertigo","DirectSecondaryConditions":[],"AggravationConditions":[],"SourceData":"30 - Vertigo"},{"PrimaryCondition":"Right knee patellar tendonitis","DirectSecondaryConditions":[],"AggravationConditions":["limitation of flexion","limitation of extension"],"SourceData":"20x - Right knee patellar tendonitis with limitation of flexion and extension (str, opinion)"}]`
	err := UT.ConditionbuzUsecase.HandleParseConditionResultFromAi(str, "test1_prompt_key")
	lib.DPrintln(err)
}
