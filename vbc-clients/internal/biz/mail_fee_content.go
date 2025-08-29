package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type MailFeeContentUsecase struct {
	log                  *log.Helper
	CommonUsecase        *CommonUsecase
	conf                 *conf.Data
	RatingPaymentUsecase *RatingPaymentUsecase
}

func NewMailFeeContentUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	RatingPaymentUsecase *RatingPaymentUsecase) *MailFeeContentUsecase {
	uc := &MailFeeContentUsecase{
		log:                  log.NewHelper(logger),
		CommonUsecase:        CommonUsecase,
		conf:                 conf,
		RatingPaymentUsecase: RatingPaymentUsecase,
	}

	return uc
}

func (c *MailFeeContentUsecase) GetCurrentEvaluation(currentRating int) (int, error) {
	currentRatingPayments, err := c.RatingPaymentUsecase.CurrentRatingPayments()
	if err != nil {
		c.log.Error(err)
		return 0, err
	}
	var ratingPaymentEntity *RatingPaymentEntity
	for k, v := range currentRatingPayments {
		if v.Rating == currentRating {
			ratingPaymentEntity = currentRatingPayments[k]
		}
	}
	if ratingPaymentEntity == nil {
		return 0, errors.New("ratingPaymentEntity is nil")
	}
	return ratingPaymentEntity.GetDollar(), nil
}

// GenContent startRating 只会是-1 50 70 90 100
func (c *MailFeeContentUsecase) GenContent(startRating int) (content string, err error) {

	currentRatingPayments, err := c.RatingPaymentUsecase.CurrentRatingPayments()
	if err != nil {
		c.log.Error(err)
		return "", err
	}

	var destRating []int
	if startRating == 50 || startRating == -1 {
		destRating = []int{50, 70, 90, 100}
	} else if startRating == 70 {
		destRating = []int{70, 90, 100}
	} else if startRating == 90 {
		destRating = []int{90, 100}
	} else if startRating == 100 {
		destRating = []int{100}
	} else {
		return "", errors.New("startRating is wrong")
	}
	for _, v := range destRating {
		ratingPayment := currentRatingPayments.GetByRating(v)
		if ratingPayment == nil {
			return "", errors.New("ratingPayment is nil")
		}
		/*
			 <li style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">50% evaluation: $1,075/month ($12,900/year, $129,000/10 years, $258,000/20 years)</li>
			<li style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">70% evaluation: $1,716/month ($20,592/year, $205,920/10 years, $411,840/20 years)</li>
			<li style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">90% evaluation: $2,241/month ($26,892/year, $268,920/10 years, $537,840/20 years)</li>
			<li style="font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;">100% evaluation: $3,737/month ($44,844/year, $448,440/10 years, $896,880/20 years)</li>
		*/
		monthPayment := ratingPayment.GetDollar()
		yearPayment := monthPayment * 12
		year10Payment := yearPayment * 10
		year20Payment := yearPayment * 20

		content += fmt.Sprintf("<li style=\"font-family:'Open Sans','Times New Roman',Arial;font-size:14px;line-height:20px;\">%d%% evaluation: $%s/month ($%s/year, $%s/10 years, $%s/20 years)</li>",
			v,
			lib.NumberEnglishPrinter(int64(monthPayment)),
			lib.NumberEnglishPrinter(int64(yearPayment)),
			lib.NumberEnglishPrinter(int64(year10Payment)),
			lib.NumberEnglishPrinter(int64(year20Payment)))

	}
	return content, nil
}
