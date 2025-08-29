package biz

import (
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/internal/config_box"
)

type BoxdcUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	BoxbuzUsecase *BoxbuzUsecase
}

func NewBoxdcUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	BoxbuzUsecase *BoxbuzUsecase) *BoxdcUsecase {
	uc := &BoxdcUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		BoxbuzUsecase: BoxbuzUsecase,
	}

	return uc
}

// RecordReviewFirstSubFolderByName 获取Data Collection/RecordReview 第一层文件夹ID
func (c *BoxdcUsecase) RecordReviewFirstSubFolderByName(name string, tCase *TData) (folderId string, err error) {

	if tCase == nil {
		return "", errors.New("RecordReviewFirstSubFolderByName: tCase is nil")
	}
	clientCaseId := tCase.CustomFields.NumberValueByNameBasic("id")
	dcRecordReviewFolderId, err := c.BoxbuzUsecase.DCRecordReviewFolderId(tCase)
	if err != nil {
		return "", err
	}
	if dcRecordReviewFolderId == "" {
		return "", errors.New("dcRecordReviewFolderId is empty")
	}

	if name == config_box.FolderName_PrivateMedicalRecords {
		return c.BoxbuzUsecase.GetBoxResId(dcRecordReviewFolderId,
			config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_RV_PrivateMedicalRecords_Folder),
			clientCaseId)
	} else if name == config_box.FolderName_VAMedicalRecords {
		return c.BoxbuzUsecase.GetBoxResId(dcRecordReviewFolderId,
			config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_RV_VAMedicalRecords_Folder),
			clientCaseId)
	} else if name == config_box.FolderName_ServiceTreatmentRecords {
		return c.BoxbuzUsecase.GetBoxResId(dcRecordReviewFolderId,
			config_box.GetCaseRelaBoxVo(config_box.CaseRelaBox_DC_RV_ServiceTreatmentRecords_Folder),
			clientCaseId)
	} else {
		return "", errors.New(name + " does not support")
	}
}
