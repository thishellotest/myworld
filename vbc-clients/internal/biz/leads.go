package biz

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/lib"
)

type LeadsUsecase struct {
	log            *log.Helper
	conf           *conf.Data
	CommonUsecase  *CommonUsecase
	LogUsecase     *LogUsecase
	WebsiteUsecase *WebsiteUsecase
}

func NewLeadsUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	LogUsecase *LogUsecase,
	WebsiteUsecase *WebsiteUsecase,
) *LeadsUsecase {
	uc := &LeadsUsecase{
		log:            log.NewHelper(logger),
		CommonUsecase:  CommonUsecase,
		conf:           conf,
		LogUsecase:     LogUsecase,
		WebsiteUsecase: WebsiteUsecase,
	}

	return uc
}

type BizLeadsSaveVo struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	State       string `json:"state"`
	LeadSource  string `json:"leadSource"`
	Branch      string `json:"branch"`
	Description string `json:"description"`
}

func (c *LeadsUsecase) BizLeadsSave(raws []byte) (lib.TypeMap, error) {
	var bizLeadsSaveVo BizLeadsSaveVo
	err := json.Unmarshal(raws, &bizLeadsSaveVo)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	if bizLeadsSaveVo.FirstName == "" || bizLeadsSaveVo.Email == "" || bizLeadsSaveVo.Phone == "" {
		return nil, errors.New("Parameter error")
	}
	if !lib.IsValidEmail(bizLeadsSaveVo.Email) {
		c.log.Error(bizLeadsSaveVo.Email)
		return nil, errors.New(bizLeadsSaveVo.Email + " is wrong")
	}
	c.LogUsecase.SaveLog(0, "BizLeadsSave", map[string]string{
		"vo": string(raws),
	})
	err = c.WebsiteUsecase.BizSyncToZohoOrVBCRM(bizLeadsSaveVo.FirstName, bizLeadsSaveVo.LastName, bizLeadsSaveVo.Email, bizLeadsSaveVo.Phone, "", bizLeadsSaveVo.State, bizLeadsSaveVo.Description, bizLeadsSaveVo.LeadSource, bizLeadsSaveVo.Branch)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
