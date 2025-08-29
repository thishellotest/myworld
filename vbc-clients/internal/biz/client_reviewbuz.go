package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ClientReviewBuzUsecase struct {
	log                 *log.Helper
	conf                *conf.Data
	CommonUsecase       *CommonUsecase
	ClientReviewUsecase *ClientReviewUsecase
}

func NewClientReviewBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	ClientReviewUsecase *ClientReviewUsecase,
) *ClientReviewBuzUsecase {
	uc := &ClientReviewBuzUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		ClientReviewUsecase: ClientReviewUsecase,
	}

	return uc
}

func (c *ClientReviewBuzUsecase) BizClientReviews() (lib.TypeMap, error) {

	records, err := c.ClientReviewUsecase.AllByCondWithOrderBy(Eq{"communication_rating": 5,
		"status":                ClientReview_Status_Enable,
		"allow_testimonial_use": ClientReview_AllowTestimonialUse_Yes}, "sort desc", 50)
	if err != nil {
		return nil, err
	}
	var res []ClientReviewVo
	for _, v := range records {
		vo := v.ToClientReviewVo()
		res = append(res, vo)
	}
	data := make(lib.TypeMap)
	data.Set("records", res)
	return data, nil
}
