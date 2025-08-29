package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type BUsaStateEntity struct {
	ID          int32 `gorm:"primaryKey"`
	EnglishName string
	Nano        string
	Tzone       string
}

func (BUsaStateEntity) TableName() string {
	return "b_usa_states"
}

type BUsaStateUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BUsaStateEntity]
	GoCacheUsecase *GoCacheUsecase
}

func NewBUsaStateUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	GoCacheUsecase *GoCacheUsecase) *BUsaStateUsecase {
	uc := &BUsaStateUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		GoCacheUsecase: GoCacheUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *BUsaStateUsecase) GetTimeLocationByUsaState(state string) (ut string, loc *time.Location, err error) {

	states, err := c.CacheAll()
	if err != nil {
		return "", nil, err
	}
	for _, v := range states {
		if strings.ToLower(v.EnglishName) == strings.ToLower(state) {
			ut = v.Tzone
			break
		}
	}
	loc, err = GetLocationByUsaTimezone(ut)
	return ut, loc, err
}

func (c *BUsaStateUsecase) CacheAll() ([]*BUsaStateEntity, error) {
	key := fmt.Sprintf("%s", "b_usa_states")
	res, found := GoCacheGet[[]*BUsaStateEntity](c.GoCacheUsecase, key)
	if found {
		return res, nil
	}
	var err error
	res, err = c.AllByCond(Neq{"id": 0})
	if err != nil {
		return nil, err
	}
	GoCacheSet(c.GoCacheUsecase, key, res, configs.CacheExpiredDurationDefault)
	return res, err
}
