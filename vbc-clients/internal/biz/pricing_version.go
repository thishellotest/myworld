package biz

import (
	"encoding/json"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/internal/config_vbc"
	. "vbc/lib/builder"
)

const (
	DefaultPricingVersion = "v20240206" // 初始价格版本
)

/*


v20241010(new):
{"BoxSignTpl":{"-1":"7a57d923-cdab-4e73-b558-a11a5e592212","0":"7fbb8d6b-ba95-4b71-8743-3163f4d973fc","10":"5c6e0607-f890-46fd-8109-5f8f3f19cd2f","20":"499aabe0-44a1-492b-8510-6f1c7341f186","30":"d9af6642-d875-4c9c-a3a0-707c4380b7f1","40":"373a5a54-b927-4114-bb4d-3f79a3876c12","50":"2745fc19-7cf2-49b2-95d4-6456dfb9e265","60":"8c164d3f-3271-4344-954e-7f1c36fc790c","70":"74319944-988e-41c9-bd88-69b63704efe4","80":"84cc5854-c5bc-4e1d-9762-15a46cfbba4d","90":"84feaede-1b7f-4df4-aa16-4f9310e8908f"},"FeeDefine":{"-1":[{"Rating":70,"Fee":3000},{"Rating":90,"Fee":5000},{"Rating":100,"Fee":10000}],"0":[{"Rating":50,"Fee":5000},{"Rating":70,"Fee":7000},{"Rating":90,"Fee":9000},{"Rating":100,"Fee":14000}],"10":[{"Rating":50,"Fee":4000},{"Rating":70,"Fee":6000},{"Rating":90,"Fee":8000},{"Rating":100,"Fee":13000}],"20":[{"Rating":50,"Fee":4000},{"Rating":70,"Fee":6000},{"Rating":90,"Fee":8000},{"Rating":100,"Fee":13000}],"30":[{"Rating":50,"Fee":3000},{"Rating":70,"Fee":5000},{"Rating":90,"Fee":7000},{"Rating":100,"Fee":12000}],"40":[{"Rating":50,"Fee":2000},{"Rating":70,"Fee":4000},{"Rating":90,"Fee":6000},{"Rating":100,"Fee":11000}],"50":[{"Rating":70,"Fee":3000},{"Rating":90,"Fee":5000},{"Rating":100,"Fee":10000}],"60":[{"Rating":70,"Fee":2000},{"Rating":90,"Fee":4000},{"Rating":100,"Fee":9000}],"70":[{"Rating":90,"Fee":3000},{"Rating":100,"Fee":8000}],"80":[{"Rating":90,"Fee":2000},{"Rating":100,"Fee":7000}],"90":[{"Rating":100,"Fee":6000}]}}


v20241010:
{"BoxSignTpl":{"-1":"27f3a7f8-e6ec-4638-b290-76e2d604b1b6","0":"a1d0c7d0-6e2d-44a0-9a8a-63475819af36","10":"4aabbd65-e665-49b2-8f50-a3509b7157c1","20":"02373279-a746-46a4-8ccf-1c9ef4e04ec6","30":"d16064bf-8734-4ad0-ac73-4762de36cb01","40":"d1703bd3-8712-44d3-9e1f-b318e4217baf","50":"cad5d179-3b1f-4fd8-9ade-f90dcefc75b3","60":"bdf8af96-7d62-4cb1-ac8d-084020fa6c85","70":"d47262c3-fabf-40c5-96e3-3655fcd0291d","80":"37daabb1-d526-40bf-a501-d36153daa321","90":"b19826c2-8f56-46c9-9126-7b75a73e0451"},"FeeDefine":{"-1":[{"Rating":70,"Fee":3000},{"Rating":90,"Fee":5000},{"Rating":100,"Fee":10000}],"0":[{"Rating":50,"Fee":5000},{"Rating":70,"Fee":7000},{"Rating":90,"Fee":9000},{"Rating":100,"Fee":14000}],"10":[{"Rating":50,"Fee":4000},{"Rating":70,"Fee":6000},{"Rating":90,"Fee":8000},{"Rating":100,"Fee":13000}],"20":[{"Rating":50,"Fee":4000},{"Rating":70,"Fee":6000},{"Rating":90,"Fee":8000},{"Rating":100,"Fee":13000}],"30":[{"Rating":50,"Fee":3000},{"Rating":70,"Fee":5000},{"Rating":90,"Fee":7000},{"Rating":100,"Fee":12000}],"40":[{"Rating":50,"Fee":2000},{"Rating":70,"Fee":4000},{"Rating":90,"Fee":6000},{"Rating":100,"Fee":11000}],"50":[{"Rating":70,"Fee":3000},{"Rating":90,"Fee":5000},{"Rating":100,"Fee":10000}],"60":[{"Rating":70,"Fee":2000},{"Rating":90,"Fee":4000},{"Rating":100,"Fee":9000}],"70":[{"Rating":90,"Fee":3000},{"Rating":100,"Fee":8000}],"80":[{"Rating":90,"Fee":2000},{"Rating":100,"Fee":7000}],"90":[{"Rating":100,"Fee":6000}]}}

*/

type PricingVersionConfig struct {
	BoxSignTpl map[string]string
	FeeDefine  config_vbc.FeeVoConfigs
}

func (c *PricingVersionConfig) GetBoxSignTpl(index string) (string, error) {
	if c.BoxSignTpl == nil {
		return "", errors.New("BoxSignTpl is nil")
	}
	if a, ok := c.BoxSignTpl[index]; ok {
		return a, nil
	} else {
		return "", errors.New("BoxSignTpl is empty from index")
	}
}

func (c *PricingVersionConfig) GetByIndex(index int) []config_vbc.FeeVo {
	if c.FeeDefine == nil {
		return nil
	}
	if _, ok := c.FeeDefine[index]; ok {
		return c.FeeDefine[index]
	}
	return nil
}

func (c *PricingVersionConfig) Charge(currentRating int, newRating int) int {
	if _, ok := c.FeeDefine[currentRating]; ok {
		for _, v := range c.FeeDefine[currentRating] {
			if v.Rating == newRating {
				return v.Fee
			}
		}
	}
	return 0
}

func (c *PricingVersionConfig) Info() {

}

const (
	IsCurrentVersion_Yes = 1
	IsCurrentVersion_No  = 0

	PricingVersion_Disabled_Yes = 1
	PricingVersion_Disabled_No  = 0
)

type PricingVersionEntity struct {
	ID               int32 `gorm:"primaryKey"`
	Version          string
	Disabled         int
	IsCurrentVersion int
	VersionConfig    string
	CreatedAt        int64
	UpdatedAt        int64
}

func (PricingVersionEntity) TableName() string {
	return "pricing_version"
}

func (c *PricingVersionEntity) GetVersionConfig() (*PricingVersionConfig, error) {
	if c.VersionConfig == "" {
		return nil, errors.New("VersionConfig is empty")
	}
	var config PricingVersionConfig
	err := json.Unmarshal([]byte(c.VersionConfig), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type PricingVersionUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[PricingVersionEntity]
}

func NewPricingVersionUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *PricingVersionUsecase {
	uc := &PricingVersionUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *PricingVersionUsecase) CurrentVersionConfig() (*PricingVersionConfig, *PricingVersionEntity, error) {
	entity, err := c.DBUsecase.GetByCond(Eq{"disabled": PricingVersion_Disabled_No, "is_current_version": IsCurrentVersion_Yes})
	if err != nil {
		return nil, nil, err
	}
	if entity == nil {
		return nil, nil, errors.New("PricingVersion Entity is nil")
	}
	config, err := entity.GetVersionConfig()
	if err != nil {
		c.log.Error(err)
		return nil, nil, err
	}
	return config, entity, nil
}

func (c *PricingVersionUsecase) CurrentVersion() (*PricingVersionEntity, error) {
	return c.DBUsecase.GetByCond(Eq{"disabled": PricingVersion_Disabled_No, "is_current_version": IsCurrentVersion_Yes})
}

func (c *PricingVersionUsecase) ConfigByPricingVersion(pricingVersion string) (*PricingVersionConfig, *PricingVersionEntity, error) {
	entity, err := c.DBUsecase.GetByCond(Eq{"version": pricingVersion})
	if err != nil {
		return nil, nil, err
	}
	if entity == nil {
		return nil, nil, errors.New("PricingVersion entity is nil: " + pricingVersion)
	}
	config, err := entity.GetVersionConfig()
	if err != nil {
		c.log.Error(err)
		return nil, nil, err
	}
	return config, entity, nil
}
