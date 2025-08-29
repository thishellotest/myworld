package tests

import (
	"testing"
	"vbc/lib"
)

func Test_ClientReviewBuzUsecase_BizClientReviews(t *testing.T) {
	res, err := UT.ClientReviewBuzUsecase.BizClientReviews()
	lib.DPrintln(res, err)
}
