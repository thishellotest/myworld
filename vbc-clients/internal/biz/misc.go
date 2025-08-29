package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"os"
	"strings"
	"sync"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
	. "vbc/lib/builder"
)

type MiscUsecase struct {
	log                                    *log.Helper
	CommonUsecase                          *CommonUsecase
	conf                                   *conf.Data
	miscThingsToKnowCPExamLock             sync.Mutex
	handleRemoveMiscThingsToKnowCPExamLock sync.Mutex
	TUsecase                               *TUsecase
	MapUsecase                             *MapUsecase
	BoxUsecase                             *BoxUsecase
	BoxbuzUsecase                          *BoxbuzUsecase
	LogUsecase                             *LogUsecase
	PrimaryUsecase                         *PrimaryUsecase
	AttorneyUsecase                        *AttorneyUsecase
	DataComboUsecase                       *DataComboUsecase
}

func NewMiscUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TUsecase *TUsecase,
	MapUsecase *MapUsecase,
	BoxUsecase *BoxUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	LogUsecase *LogUsecase,
	PrimaryUsecase *PrimaryUsecase,
	AttorneyUsecase *AttorneyUsecase,
	DataComboUsecase *DataComboUsecase) *MiscUsecase {
	uc := &MiscUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		TUsecase:         TUsecase,
		MapUsecase:       MapUsecase,
		BoxUsecase:       BoxUsecase,
		BoxbuzUsecase:    BoxbuzUsecase,
		LogUsecase:       LogUsecase,
		PrimaryUsecase:   PrimaryUsecase,
		AttorneyUsecase:  AttorneyUsecase,
		DataComboUsecase: DataComboUsecase,
	}

	return uc
}

func (c *MiscUsecase) Gen2122aFileNameForMisc(tCase TData, tClient TData) (string, error) {

	attorneyUniqid := tCase.CustomFields.TextValueByNameBasic(FieldName_attorney_uniqid)
	if attorneyUniqid == "" {
		return "", errors.New("attorneyUniqid is empty")
	}
	attorney, err := c.AttorneyUsecase.GetByGid(attorneyUniqid)
	if err != nil {
		return "", err
	}
	if attorney == nil {
		return "", errors.New("attorney is nil")
	}
	attorneyName := fmt.Sprintf("%s%s", attorney.FirstName, attorney.LastName)

	return fmt.Sprintf("%s - 21-22a - Appointment as Claimant Representative for %s.pdf", attorneyName, tClient.CustomFields.TextValueByNameBasic(FieldName_full_name)), nil
}

func (c *MiscUsecase) HandleMoving2122aFile(caseId int32) (boxFileId string, err error) {
	key := MapKeyMoving2122aFileId(caseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return "", err
	}
	if val == "" {
		destBoxFileId, err := c.DoHandleMoving2122aFile(caseId)
		if err != nil {
			return "", err
		}
		c.MapUsecase.Set(key, destBoxFileId)
		return destBoxFileId, nil
	} else {
		return val, nil
	}

}

func (c *MiscUsecase) Delete2122aFile(caseId int32) (r string, err error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return "", err
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}

	CMiscFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
		MapKeyBuildAutoBoxCMiscFolderId(tCase.Id()),
		tCase)
	if err != nil {
		return "", err
	}
	if CMiscFolderId == "" {
		return "", errors.New("CMiscFolderId is empty")
	}
	aFileName, err := c.Gen2122aFileNameForMisc(*tCase, *tClient)
	if err != nil {
		return "", err
	}

	res, err := c.BoxUsecase.ListItemsInFolderFormat(CMiscFolderId)
	if err != nil {
		return "", err
	}
	for _, v := range res {
		if v.GetString("name") == aFileName {
			_, err = c.BoxUsecase.DeleteFile(v.GetString("id"))
			return "", err
		}
	}
	return "", nil
}

func (c *MiscUsecase) DoHandleMoving2122aFile(caseId int32) (destBoxFileId string, err error) {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, caseId)
	if err != nil {
		return "", err
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}

	tClient, _, _ := c.DataComboUsecase.ClientWithCase(*tCase)
	if tClient == nil {
		return "", errors.New("tClient is nil")
	}

	CMiscFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
		MapKeyBuildAutoBoxCMiscFolderId(tCase.Id()),
		tCase)
	if err != nil {
		return "", err
	}
	if CMiscFolderId == "" {
		return "", errors.New("CMiscFolderId is empty")
	}
	aFileName, err := c.Gen2122aFileNameForMisc(*tCase, *tClient)
	if err != nil {
		return "", err
	}
	key := MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId)
	boxFileId, _ := c.MapUsecase.GetForString(key)
	if boxFileId == "" {
		return "", errors.New("boxFileId is empty")
	}
	lib.DPrintln("21-22a filename:", aFileName)
	destBoxFileId, err = c.BoxUsecase.CopyFileNewFileNameReturnFileId(boxFileId, aFileName, CMiscFolderId)
	if err != nil {
		return "", err
	}
	return destBoxFileId, nil
}

func (c *MiscUsecase) HandleRemoveMiscThingsToKnowCPExam(clientCaseId int32) error {

	// 后续发现性能低下，可以把此处的锁关闭
	c.handleRemoveMiscThingsToKnowCPExamLock.Lock()
	defer c.handleRemoveMiscThingsToKnowCPExamLock.Unlock()

	tCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	// 此处需要获取primaryCase
	primaryCase, _, err := c.PrimaryUsecase.GetPrimaryCase(tCase)
	if err != nil {
		return err
	}

	key := MapKeyClientMiscThingsToKnowCPExamFileId(primaryCase.Id())
	// 判断文件是否存在
	clientMiscThingsToKnowCPExamFileId, err := c.MapUsecase.GetForString(key)
	c.log.Info("HandleRemoveMiscThingsToKnowCPExam clientCaseId: ", clientCaseId, " PrimaryCaseId: ", primaryCase.Id(), " clientMiscThingsToKnowCPExamFileId: ", clientMiscThingsToKnowCPExamFileId)
	if err != nil {
		return err
	}

	var hasMapValue = false
	if clientMiscThingsToKnowCPExamFileId == "" {

		CMiscFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
			MapKeyBuildAutoBoxCMiscFolderId(primaryCase.Id()),
			primaryCase)
		if err != nil {
			return err
		}
		if CMiscFolderId == "" {
			return errors.New("HandleMiscThingsToKnowCPExam: CMiscFolderId is empty")
		}

		clientMiscThingsToKnowCPExamFileId, err = c.ClientThingsToKnowExamFileIdByApi(CMiscFolderId)
		if err != nil {
			return err
		}
	} else {
		hasMapValue = true
	}
	if clientMiscThingsToKnowCPExamFileId != "" {
		_, err = c.BoxUsecase.DeleteFile(clientMiscThingsToKnowCPExamFileId)
		if err != nil {
			return err
		}
		if hasMapValue {
			err = c.MapUsecase.Set(key, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *MiscUsecase) HandleMiscThingsToKnowCPExam(clientCaseId int32) error {

	tCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	// 此处需要获取primaryCase
	primaryCase, _, err := c.PrimaryUsecase.GetPrimaryCase(tCase)
	if err != nil {
		return err
	}

	key := MapKeyClientMiscThingsToKnowCPExamFileId(primaryCase.Id())
	// 判断文件是否存在
	clientMiscThingsToKnowCPExamFileId, err := c.MapUsecase.GetForString(key)
	c.log.Info("HandleMiscThingsToKnowCPExam clientCaseId: ", clientCaseId, " PrimaryCase: ", primaryCase.Id(), " clientMiscThingsToKnowCPExamFileId: ", clientMiscThingsToKnowCPExamFileId)
	if err != nil {
		return err
	}
	if clientMiscThingsToKnowCPExamFileId != "" { // The file already exists.
		return nil
	}

	// 后续发现性能低下，可以把此处的锁关闭
	c.miscThingsToKnowCPExamLock.Lock()
	defer c.miscThingsToKnowCPExamLock.Unlock()

	//clientBoxFolderId, err := c.MapUsecase.GetForString(MapKeyClientBoxFolderId(clientCaseId))
	//if err != nil {
	//	return err
	//}
	//if clientBoxFolderId == "" {
	//	return errors.New("HandleMiscThingsToKnowCPExam: clientBoxFolderId is nil")
	//}

	CMiscFolderId, err := c.BoxbuzUsecase.GetClientSubFolderId(
		MapKeyBuildAutoBoxCMiscFolderId(primaryCase.Id()),
		primaryCase)
	if err != nil {
		return err
	}
	if CMiscFolderId == "" {
		return errors.New("HandleMiscThingsToKnowCPExam: CMiscFolderId is empty")
	}
	thingsToKnowExamFileId, err := c.ClientThingsToKnowExamFileId(CMiscFolderId)
	if err != nil {
		return err
	}
	if thingsToKnowExamFileId == "" {
		return errors.New("thingsToKnowExamFileId is empty")
	}
	err = c.MapUsecase.Set(key, thingsToKnowExamFileId)

	return err
}

func GetHowToGuide(conf *conf.Data) (path string, name string) {
	var resPath string
	name = "How-to-Guide v7.3.pdf"
	if configs.IsDev() {
		resPath = "/Users/garyliao/code/vbc-clients/resource"
	} else {
		resPath = conf.ResourcePath
	}
	return resPath + "/" + name, name
}

func (c *MiscUsecase) UpdateAll() error {

	lib.DPrintln("MiscUsecase_UpdateAll:running", time.Now().Format(time.RFC3339))

	list, err := c.TUsecase.ListByCond(Kind_client_cases, And(Eq{"deleted_at": 0, "biz_deleted_at": 0}, Gt{"id": 280}))
	if err != nil {
		lib.DPrintln(err)
	}
	for k, v := range list {
		key := MapKeyClientBoxFolderId(v.Id())
		a, _ := c.MapUsecase.GetForString(key)
		if a != "" {
			err := c.UpdateGuideForClientCase(list[k])
			if err != nil {
				c.LogUsecase.SaveLog(v.Id(), "MiscUsecase_UpdateAll_Error", map[string]interface{}{})
				lib.DPrintln("MiscUsecase_UpdateAll:updated error:", v.Id())
			} else {
				c.LogUsecase.SaveLog(v.Id(), "MiscUsecase_UpdateAll_Ok", map[string]interface{}{})
				lib.DPrintln("MiscUsecase_UpdateAll:updated ok:", v.Id())
			}
		}
	}
	return nil
}

func (c *MiscUsecase) UpdateGuideForClientCase(tCase *TData) error {

	path, newGuideName := GetHowToGuide(c.conf)

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

	guideFileInfo, err := c.GuideFileByApi(CMiscFolderId)
	if err != nil {
		return err
	}
	currentName := guideFileInfo.GetString("name")
	if currentName == newGuideName {
		c.log.Info(InterfaceToString(clientCaseId) + " It's already been a new guide file")
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		c.log.Error("clientCaseId:", clientCaseId, " ", err)
		return err
	}
	defer file.Close()
	_, err = c.BoxUsecase.UploadFileVersionWithNewFileName(guideFileInfo.GetString("id"), file, newGuideName)
	if err != nil {
		c.log.Error("clientCaseId:", clientCaseId, " ", err)
		return err
	}
	err = c.LogUsecase.SaveLog(clientCaseId, "UpdateGuideForClientCase", nil)
	if err != nil {
		c.log.Error("clientCaseId:", clientCaseId, " ", err)
		return err
	}

	//c.BoxUsecase.UploadFileVersion()
	return nil
}

// GuideFileByApi {"etag":"0","file_version":{"id":"1793231125933","sha1":"8e24b54e8987d75d92bde9409bf1a8baf7506865","type":"file_version"},"id":"1630836585133","name":"How-to-Guide v5.3.pdf","sequence_id":"0","sha1":"8e24b54e8987d75d92bde9409bf1a8baf7506865","type":"file"}
func (c *MiscUsecase) GuideFileByApi(CMiscFolderId string) (fileInfo lib.TypeMap, err error) {

	CMiscFolderSubs, err := c.BoxUsecase.ListItemsInFolderFormat(CMiscFolderId)
	if err != nil {
		return nil, err
	}
	for k, v := range CMiscFolderSubs {
		if v.GetString("type") == "file" &&
			(strings.Index(v.GetString("name"), config_box.FileName_HowToGuide_Prefix) >= 0) {
			return CMiscFolderSubs[k], nil
		}
	}
	return fileInfo, nil
}

func (c *MiscUsecase) ClientThingsToKnowExamFileIdByApi(CMiscFolderId string) (thingsToKnowExamFileId string, err error) {

	CMiscFolderSubs, err := c.BoxUsecase.ListItemsInFolderFormat(CMiscFolderId)
	if err != nil {
		return "", err
	}
	for _, v := range CMiscFolderSubs {
		if v.GetString("type") == "file" &&
			(strings.Index(v.GetString("name"), config_box.FileName_ThingsToKnowExam_Prefix) >= 0) {
			thingsToKnowExamFileId = v.GetString("id")
			break
		}
	}
	return thingsToKnowExamFileId, nil
}

func (c *MiscUsecase) ClientThingsToKnowExamFileId(CMiscFolderId string) (thingsToKnowExamFileId string, err error) {

	thingsToKnowExamFileId, err = c.ClientThingsToKnowExamFileIdByApi(CMiscFolderId)
	if err != nil {
		return "", err
	}
	if thingsToKnowExamFileId == "" { // file copies to the Misc Folder
		thingsToKnowExamFileId, err = c.CopyThingsToKnowExamFile(CMiscFolderId)
	}
	return
}

func (c *MiscUsecase) CopyThingsToKnowExamFile(CMiscFolderId string) (thingsToKnowExamFileId string, err error) {
	return c.BoxUsecase.CopyFileNewFileNameReturnFileId(config_box.ThingsToKnowExamTplFileId, config_box.FileName_ThingsToKnowExam, CMiscFolderId)
}
