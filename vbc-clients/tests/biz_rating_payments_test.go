package tests

import (
	"testing"
	"vbc/lib"
)

func Test_RatingPaymentUsecase_CurrentRatingPayments(t *testing.T) {
	a, err := UT.RatingPaymentUsecase.CurrentRatingPayments()
	for _, v := range a {
		lib.DPrintln(v.Rating, v.Payment, v.EffectiveDate, "GetDollar: ", v.GetDollar())
	}
	te := a.GetByRating(10)
	lib.DPrintln("aaa:", te)
	lib.DPrintln(a, err)
}
