package tests

import (
	"testing"
	"vbc/lib"
)

func Test_AiStatementUsecase_GenStatementByCondition(t *testing.T) {
	//UT.AiStatementUsecase.GenStatementByCondition("")
}

func Test_AiStatementUsecase_SplitCaseDescription(t *testing.T) {

	notes := `Service:
1983 - 1984 (Army Active Enlisted)
1987 (Commissioned Officer)
1988-1992 (Army Officer)
1992-2002 (Army NG Reserve)
2011-2024 (Army NG Reserve)

Current:
N/A

New:
70 - PTSD
50 - Sleep apnea
30 - Bilateral Plantar Fasciitis
20* - Bilateral knee pain with limitation of flexion and extension
20* - Left shoulder pain AC joint with limitation of flexion and abduction
20* - Low back pain with right leg radiculopathy
10 - Chronic sinus congestion
10 - TBI
10 - Tinnitus
0 - Hearing loss

2nd wave:
30 - GERD
10 - Irritable Bowel Syndrome`
	a, err := UT.AiStatementUsecase.SplitCaseDescription(notes)
	lib.DPrintln(a, err)
}
