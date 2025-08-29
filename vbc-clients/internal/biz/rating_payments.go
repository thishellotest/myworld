package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type RatingPaymentEntity struct {
	ID            int32 `gorm:"primaryKey"`
	EffectiveDate string
	Rating        int
	Payment       int
}

func (RatingPaymentEntity) TableName() string {
	return "rating_payments"
}

// GetDollar 把美分转为美元,采用四舍五入方法
func (c *RatingPaymentEntity) GetDollar() int {
	if c.Payment <= 0 {
		return 0
	}
	a := math.Floor(float64(c.Payment) / 100)
	return int(a)
}

func CentToDollar(cent float32) int {
	a := math.Floor(float64(cent) / 100)
	return int(a)
}

type RatingPaymentUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[RatingPaymentEntity]
}

func NewRatingPaymentUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *RatingPaymentUsecase {
	uc := &RatingPaymentUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type RatingPaymentList []*RatingPaymentEntity

func (c RatingPaymentList) GetByRating(rating int) *RatingPaymentEntity {
	for k, v := range c {
		if v.Rating == rating {
			return c[k]
		}
	}
	return nil
}

func (c *RatingPaymentUsecase) CurrentRatingPayments() (RatingPaymentList, error) {
	currentDate := time.Now().In(configs.GetVBCDefaultLocation()).Format(time.DateOnly)
	//currentDate = "2024-12-01" // test
	return c.AllByCondWithOrderBy(Lte{"effective_date": currentDate}, "effective_date desc,rating ", 11)
}
