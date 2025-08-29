package biz

import (
	"testing"
	"vbc/lib"
)

func Test_CaseClaimsTypeV2(t *testing.T) {
	str := "Bilateral hip pain with limitation of flexion and extension secondary to right ankle strain (str, opinion)"
	a, b := CaseClaimsTypeV2(str)
	lib.DPrintln(a)
	lib.DPrintln(b)
}

func Test_CaseClaimsDivide(t *testing.T) {

	val := `70 - Major Depressive Disorder (new,opinion)

50 - Headaches (str)

30 - GERD secondary to chronic NSAID use related to headaches (opinion)

20* - Bilateral hip pain with limitation of flexion and extension secondary to right ankle strain (opinion)

20 - Diabetes Mellitus Type II (Agent Orange)

10 - Hypertension secondary to Diabetes Mellitus Type II (opinion)

10 - Eczema (str)

10 - Dermatitis (str)

10 - Sinusitis (str)

10 - Tinnitus (new)

10 - Right ankle sprain (str)

10 - Urinary frequency secondary to Diabetes (opinion)

0 - ED secondary to Major Depressive Disorder (new)`

	info := CaseClaimsDivideV2(val)
	lib.DPrintln(info)
}

func TestCaseClaimsType(t *testing.T) {
	a, ty := CaseClaimsType("ED secondary to Major Depressive Disorder (new)")
	lib.DPrintln(a, ty)
}
