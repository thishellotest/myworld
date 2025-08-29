package tests

import (
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_DocEmailUsecase_HandleDocEmailToBox(t *testing.T) {
	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5431)
	err := UT.DocEmailUsecase.HandleDocEmailToBox(tCase)
	lib.DPrintln(err)
}

func Test_DocEmailExtractHealthIssues(t *testing.T) {
	text := `# Doctor's Email

Subject: Several Health Issues to Discuss

Dear Dr. [DOCTOR'S NAME],

I hope you're doing well. 

I forgot to mention several important health issues during our last visit. Is this something you can help with, or do I need to schedule an appointment?

I wanted to mention these health issues prior to our appointment on [DATE].

I have been experiencing several health issues lately. Is this something you can help with, or do I need to schedule an appointment?

Here's a rundown of my current health concerns:

1. Insomnia Disorder: My sleep issues have worsened significantly. I'm getting very little quality sleep each night, and it's affecting every aspect of my life.

2. Bilateral Pes Planus: The pain in my feet has become increasingly severe, making it difficult to stand for even short periods.

3. Back Pain: I'm experiencing persistent and worsening pain in my lower back that seems connected to my foot problems.

4. Hypertension: My blood pressure readings have been consistently high despite medication.

5. Tinnitus: I've been experiencing constant ringing in my ears that's extremely distracting and affecting my concentration.
6. Erectile Dysfunction: I've been having ongoing issues that appear to be related to my sleep medication and overall poor sleep quality.

All of these issues are making daily life extremely challenging. My sleep is severely disrupted, I'm frequently in pain or discomfort, and my mental health is suffering. These conditions are affecting my ability to work effectively and maintain healthy relationships with my family.

Is this something we can discuss at my next visit, or should I schedule a separate appointment to go over all of this?

Thank you for your help and understanding.

Sincerely,
MarkDean Ronduen`
	aa := biz.DocEmailExtractHealthIssues(text)
	lib.DPrintln(aa)
}

func Test_DocEmailUsecase_DocEmailResultTextByCase(t *testing.T) {

	tCase, _ := UT.TUsecase.DataById(biz.Kind_client_cases, 5511)
	tClient, _, _ := UT.DataComboUsecase.ClientWithCase(*tCase)
	a, err := UT.DocEmailUsecase.DocEmailResultTextByCase(*tCase, *tClient)
	lib.DPrintln(a, err)
}
