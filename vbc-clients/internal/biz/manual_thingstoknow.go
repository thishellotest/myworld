package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"os"
	"strings"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/internal/config_vbc"
	"vbc/lib"
	"vbc/lib/builder"
)

func GetThingsToKnow(conf *conf.Data) (path string, name string) {
	var resPath string
	name = config_box.FileName_ThingsToKnowExam
	if configs.IsDev() {
		resPath = "/Users/garyliao/code/vbc-clients/resource"
	} else {
		resPath = conf.ResourcePath
	}
	return resPath + "/" + name, name
}

type ManualThingstoknowUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	TUsecase      *TUsecase
	BoxbuzUsecase *BoxbuzUsecase
	BoxUsecase    *BoxUsecase
	LogUsecase    *LogUsecase
}

func NewManualThingstoknowUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	LogUsecase *LogUsecase) *ManualThingstoknowUsecase {
	uc := &ManualThingstoknowUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
		BoxbuzUsecase: BoxbuzUsecase,
		BoxUsecase:    BoxUsecase,
		LogUsecase:    LogUsecase,
	}

	return uc
}

func (c *ManualThingstoknowUsecase) DestClientCases() ([]*TData, error) {

	return c.TUsecase.ListByCond(Kind_client_cases, builder.And(builder.In(FieldName_stages,
		config_vbc.Stages_MedicalTeamExamsScheduled,
		config_vbc.Stages_MedicalTeamCallVet,
		config_vbc.Stages_DBQ_Completed,
		config_vbc.Stages_FileClaims_Draft,
		config_vbc.Stages_FileClaims,
		config_vbc.Stages_VerifyEvidenceReceived,
		config_vbc.Stages_AwaitingDecision,
	),
		//builder.In("id", 94, 75, 61, 58, 55, 5495, 5465, 5434, 5355),
		builder.Eq{"biz_deleted_at": 0}))
}

func (c *ManualThingstoknowUsecase) HandleUploadNewThingsToKnowFileAllCases() error {
	cases, err := c.DestClientCases()
	if err != nil {
		return err
	}
	for k, v := range cases {
		c.log.Info("ManualThingstoknowUsecase CaseId: ", v.Id())
		c.LogUsecase.SaveLog(v.Id(), "UploadNewThingsToKnowFileAllCasesLog", nil)
		err := c.UploadNewThingsToKnowFile(cases[k])
		if err != nil {
			c.LogUsecase.SaveLog(v.Id(), "UploadNewThingsToKnowFileAllCasesLogError", map[string]interface{}{
				"error": err.Error(),
			})
		}
		//break
	}
	return nil
}

func (c *ManualThingstoknowUsecase) UploadNewThingsToKnowFile(tCase *TData) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	clientCaseId := tCase.Id()
	CMiscFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
		MapKeyBuildAutoBoxCMiscFolderId(clientCaseId),
		tCase)
	if err != nil {
		c.log.Error("clientCaseId:", clientCaseId, " ", err)
		return err
	}
	if CMiscFolderId == "" {
		c.log.Error("clientCaseId:", clientCaseId, " CMiscFolderId is empty")
		return errors.New("CMiscFolderId is empty")
	}
	thingsToKnowFileInfo, err := c.ThingsToKnowFileByApi(CMiscFolderId)
	if err != nil {
		return err
	}
	if thingsToKnowFileInfo == nil {
		return errors.New("thingsToKnowFileInfo is nil")
	}
	name := thingsToKnowFileInfo.GetString("name")
	if name != config_box.FileName_ThingsToKnowExam {
		newFilePath, newFileName := GetThingsToKnow(c.conf)

		file, err := os.Open(newFilePath)
		if err != nil {
			c.log.Error("clientCaseId:", clientCaseId, " ", err)
			return err
		}
		defer file.Close()
		_, err = c.BoxUsecase.UploadFileVersionWithNewFileName(thingsToKnowFileInfo.GetString("id"), file, newFileName)
		if err != nil {
			c.log.Error("clientCaseId:", clientCaseId, " ", err)
			return err
		}
		c.log.Info("UploadFileVersionWithNewFileName : ", tCase.Id(), " ok")
		//err = c.LogUsecase.SaveLog(clientCaseId, "UpdateGuideForClientCase", nil)
		//if err != nil {
		//	c.log.Error("clientCaseId:", clientCaseId, " ", err)
		//	return err
		//}

	} else {
		c.log.Info("UploadFileVersionWithNewFileName : ", tCase.Id(), " exists")
	}
	return nil
}

// ThingsToKnowFileByApi {"etag":"0","file_version":{"id":"1859167527731","sha1":"0bdc5e6a7adbb7f88c869cd09bf6a39248243163","type":"file_version"},"id":"1689345042131","name":"Things to know before your exam v2.1.pdf","sequence_id":"0","sha1":"0bdc5e6a7adbb7f88c869cd09bf6a39248243163","type":"file"}
func (c *ManualThingstoknowUsecase) ThingsToKnowFileByApi(CMiscFolderId string) (fileInfo lib.TypeMap, err error) {

	CMiscFolderSubs, err := c.BoxUsecase.ListItemsInFolderFormat(CMiscFolderId)
	if err != nil {
		return nil, err
	}
	for k, v := range CMiscFolderSubs {
		if v.GetString("type") == "file" &&
			(strings.Index(v.GetString("name"), config_box.FileName_ThingsToKnowExam_Prefix) >= 0) {
			return CMiscFolderSubs[k], nil
		}
	}
	return fileInfo, nil
}
