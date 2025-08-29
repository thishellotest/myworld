package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"vbc/internal/conf"
	"vbc/internal/config_box"
)

type BoxcUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	BoxbuzUsecase *BoxbuzUsecase
}

func NewboxcUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxbuzUsecase *BoxbuzUsecase) *BoxcUsecase {
	uc := &BoxcUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		BoxbuzUsecase: BoxbuzUsecase,
	}

	return uc
}

func (c *BoxcUsecase) GetNewEvidenceFolderId(primaryCase *TData, tCase *TData) (newEvidenceFolderId string, err error) {
	if primaryCase == nil {
		return "", errors.New("primaryCase is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	primaryCaseClientFolderId, err := c.BoxbuzUsecase.GetClientBoxFolderId(primaryCase)
	if err != nil {
		return "", err
	}
	boxResId, err := c.BoxbuzUsecase.GetBoxResId(primaryCaseClientFolderId,
		config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_C_New_Evidence_Folder),
		tCase.CustomFields.NumberValueByNameBasic("id"))
	if err != nil {
		return "", err
	}
	return boxResId, nil
}
