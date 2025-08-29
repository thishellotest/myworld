package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	. "vbc/lib/builder"
)

type AttorneybuzUsecase struct {
	log             *log.Helper
	conf            *conf.Data
	CommonUsecase   *CommonUsecase
	MapUsecase      *MapUsecase
	AttorneyUsecase *AttorneyUsecase
}

func NewAttorneybuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	MapUsecase *MapUsecase,
	AttorneyUsecase *AttorneyUsecase,
) *AttorneybuzUsecase {
	uc := &AttorneybuzUsecase{
		log:             log.NewHelper(logger),
		CommonUsecase:   CommonUsecase,
		conf:            conf,
		MapUsecase:      MapUsecase,
		AttorneyUsecase: AttorneyUsecase,
	}

	return uc
}

func (c *AttorneybuzUsecase) CurrentAttorneyIndex() (int, error) {
	key := MapKeyCurrentAttorneyIndex()
	a, err := c.MapUsecase.GetForInt(key)
	if err != nil {
		return 0, err
	}
	return int(a), nil
}

func (c *AttorneybuzUsecase) GetAnAttorney() (*AttorneyEntity, error) {

	res, err := c.AttorneyUsecase.AllByCond(Eq{"deleted_at": 0, "status": 1})
	if err != nil {
		return nil, err
	}
	index, err := c.CurrentAttorneyIndex()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	index += 1
	if index >= len(res) {
		index = 0
	}
	key := MapKeyCurrentAttorneyIndex()
	c.MapUsecase.SetInt(key, index)
	return res[index], nil
}
