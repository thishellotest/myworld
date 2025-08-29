package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type BoxUserBuzUsecase struct {
	log            *log.Helper
	conf           *conf.Data
	CommonUsecase  *CommonUsecase
	BoxUserUsecase *BoxUserUsecase
	BoxUsecase     *BoxUsecase
}

func NewBoxUserBuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	BoxUserUsecase *BoxUserUsecase,
	BoxUsecase *BoxUsecase,
) *BoxUserBuzUsecase {
	uc := &BoxUserBuzUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		BoxUserUsecase: BoxUserUsecase,
		BoxUsecase:     BoxUsecase,
	}

	return uc
}

func (c *BoxUserBuzUsecase) ScanBoxUser() error {
	res, _, err := c.BoxUsecase.Users()
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("res is nil")
	}
	data := lib.ToTypeMapByString(*res)
	entries := data.GetTypeList("entries")
	for k, _ := range entries {
		err := c.BoxUserUsecase.Upsert(entries[k])
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}
