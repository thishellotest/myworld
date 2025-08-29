package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
	//. "vbc/lib/builder"
)

type RecordReviewParams struct {
	ClientCaseId       int32
	FirstSubFolderId   string
	FirstSubFolderName string //  config_box.FolderName_PrivateMedicalRecords config_box.FolderName_VAMedicalRecords  config_box.FolderName_ServiceTreatmentRecords

	// {"etag":"0","id":"268258622342","name":"leve1_folder","sequence_id":"0","type":"folder"}
	// 通过 ClientCaseId DestId DestType, 合并任务做一次处理
	DestId       string                // 268258622342
	DestType     config_box.BoxResType // folder file
	DestName     string                // leve1_folder
	SourceFromId int32                 // 事件来源ID
}

type RecordReviewJobUsecase struct {
	conf             *conf.Data
	CommonUsecase    *CommonUsecase
	log              *log.Helper
	DataEntryUsecase *DataEntryUsecase
	CustomTask
	DataComboUsecase     *DataComboUsecase
	TUsecase             *TUsecase
	BoxdcUsecase         *BoxdcUsecase
	BoxbuzUsecase        *BoxbuzUsecase
	BoxUsecase           *BoxUsecase
	LogUsecase           *LogUsecase
	ReminderEventUsecase *ReminderEventUsecase
	MapUsecase           *MapUsecase
	BoxcUsecase          *BoxcUsecase
	ClientCaseUsecase    *ClientCaseUsecase
	CounterbuzUsecase    *CounterbuzUsecase
}

func NewRecordReviewJobUsecase(CommonUsecase *CommonUsecase,
	logger log.Logger,
	DataEntryUsecase *DataEntryUsecase,
	DataComboUsecase *DataComboUsecase,
	TUsecase *TUsecase,
	BoxdcUsecase *BoxdcUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	LogUsecase *LogUsecase,
	ReminderEventUsecase *ReminderEventUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	BoxcUsecase *BoxcUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	CounterbuzUsecase *CounterbuzUsecase,
) *RecordReviewJobUsecase {

	recordReviewJobUsecase := &RecordReviewJobUsecase{
		CommonUsecase:        CommonUsecase,
		log:                  log.NewHelper(logger),
		DataEntryUsecase:     DataEntryUsecase,
		DataComboUsecase:     DataComboUsecase,
		TUsecase:             TUsecase,
		BoxdcUsecase:         BoxdcUsecase,
		BoxbuzUsecase:        BoxbuzUsecase,
		BoxUsecase:           BoxUsecase,
		LogUsecase:           LogUsecase,
		ReminderEventUsecase: ReminderEventUsecase,
		conf:                 conf,
		MapUsecase:           MapUsecase,
		BoxcUsecase:          BoxcUsecase,
		ClientCaseUsecase:    ClientCaseUsecase,
		CounterbuzUsecase:    CounterbuzUsecase,
	}
	recordReviewJobUsecase.CustomTask.RedisQueue = Redis_record_review_tasks_queue
	recordReviewJobUsecase.CustomTask.RedisProcessing = Redis_record_review_tasks_processing
	recordReviewJobUsecase.CustomTask.RedisClient = CommonUsecase.RedisClient()
	recordReviewJobUsecase.CustomTask.Log = log.NewHelper(logger)
	recordReviewJobUsecase.CustomTask.Handle = recordReviewJobUsecase.HandleTask
	recordReviewJobUsecase.CustomTask.MaxBatchLimit = 100000
	// todo:lgl需要改为3分钟到5分钟，测试时，可以设置为30秒
	recordReviewJobUsecase.CustomTask.WindowSeconds = 5 * 60

	return recordReviewJobUsecase
}

func RecordReviewVersion() string {
	return "v" + time.Now().In(configs.VBCDefaultLocation).Format("2006-01-02_15:04:05")
}

func (c *RecordReviewJobUsecase) TidyData(recordReviewParams RecordReviewParams) (
	needHandle bool,
	lastName string,
	err error) {

	var httpCode int
	var res lib.TypeMap
	if recordReviewParams.DestType == config_box.BoxResType_folder {
		res, httpCode, err = c.BoxUsecase.GetFolderInfoForTypeMap(recordReviewParams.DestId)
	} else {
		res, httpCode, err = c.BoxUsecase.GetFileInfoForTypeMap(recordReviewParams.DestId)
	}
	if err != nil {
		return false, "", err
	}
	if httpCode == 404 {
		lib.DPrintln("RecordReviewJobUsecase: TidyData 404")
		return false, "", nil
	}
	if httpCode != 200 {
		lib.DPrintln("RecordReviewJobUsecase: TidyData no 200:", httpCode)
		return false, "", errors.New("HttpCode error: " + InterfaceToString(httpCode))
	}
	if res.GetString("parent.id") != recordReviewParams.FirstSubFolderId {
		c.log.Info("RecordReviewJobUsecase: TidyData, parent.id unequal to target folder id")
		return false, "", nil
	}

	return true, res.GetString("name"), nil
}

func (c *RecordReviewJobUsecase) TidyDataV1(recordReviewParams RecordReviewParams) (
	needHandle bool,
	res lib.TypeMap,
	err error) {

	var httpCode int
	if recordReviewParams.DestType == config_box.BoxResType_folder {
		res, httpCode, err = c.BoxUsecase.GetFolderInfoForTypeMap(recordReviewParams.DestId)
	} else {
		res, httpCode, err = c.BoxUsecase.GetFileInfoForTypeMap(recordReviewParams.DestId)
	}
	if err != nil {
		return false, res, err
	}
	if httpCode == 404 {
		lib.DPrintln("RecordReviewJobUsecase: TidyData 404")
		return false, res, nil
	}
	if httpCode != 200 {
		lib.DPrintln("RecordReviewJobUsecase: TidyData no 200:", httpCode)
		return false, res, errors.New("HttpCode error: " + InterfaceToString(httpCode))
	}

	return true, res, nil
}

func (c *RecordReviewJobUsecase) HandleTask(ctx context.Context, customTaskParams CustomTaskParams) error {
	err := c.BizHandleTask(ctx, customTaskParams)

	if err != nil {
		c.LogUsecase.SaveLog(0, "RecordReviewJobUsecase:HandleTask", map[string]interface{}{
			"customTaskParams": customTaskParams,
			"err":              err.Error(),
		})
	}
	return err
}

// BizHandleTask 处理单个task
func (c *RecordReviewJobUsecase) BizHandleTask(ctx context.Context, customTaskParams CustomTaskParams) error {

	if c.conf.Box.SyncRecordReviewVersion == config_box.SyncRecordReviewVersionV1 {
		return c.BizHandleTaskV1(ctx, customTaskParams)
	}

	lib.DPrintln("RecordReviewJobUsecase customTaskParams:", customTaskParams)

	var recordReviewParams RecordReviewParams
	err := json.Unmarshal([]byte(customTaskParams.Params), &recordReviewParams)
	if err != nil {
		return err
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, recordReviewParams.ClientCaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	tCaseId := tCase.CustomFields.NumberValueByNameBasic("id")
	folderId, err := c.BoxdcUsecase.RecordReviewFirstSubFolderByName(recordReviewParams.FirstSubFolderName, tCase)
	if err != nil {
		return err
	}
	subItems, err := c.BoxUsecase.ListItemsInFolderFormat(folderId)
	if err != nil {
		return err
	}

	// 此处
	// - 需要获取最新的名称：recordReviewParams.DestName， 因为客户有可能会改名
	// - 需要确认是否在正确的层级，因为客户有可能会移动文件夹或目录
	needHandle, newDestName, err := c.TidyData(recordReviewParams)
	if err != nil {
		return err
	}
	if !needHandle {
		return nil
	}
	recordReviewParams.DestName = newDestName

	existBoxRes := false
	for _, v := range subItems {
		if v.GetString("type") == string(recordReviewParams.DestType) &&
			v.GetString("name") == string(recordReviewParams.DestName) {
			existBoxRes = true
			break
		}
	}
	newName := ""
	if existBoxRes {
		newName = lib.FormatNameWithVersion(recordReviewParams.DestName, RecordReviewVersion())
	} else {
		newName = recordReviewParams.DestName
	}

	if recordReviewParams.DestType == config_box.BoxResType_folder {
		resId, _, err := c.BoxUsecase.CopyFolder(recordReviewParams.DestId, newName, folderId)
		if err != nil {
			return err
		}
		lib.DPrintln("RecordReviewJobUsecase CopyFolder: resId:", resId)
		eventData := &ReminderClientUpdateFilesEventVo{
			Items: []*ReminderClientUpdateFilesEventVoItem{
				{
					BoxResId:   resId,
					BoxResType: config_box.BoxResType_folder,
					BoxResName: newName,
				},
			},
		}
		err = c.ReminderEventUsecase.AddClientUpdateFilesEvent(tCaseId, eventData)
		if err != nil {
			return err
		}
	} else {
		resFieldId, err := c.BoxUsecase.CopyFileNewFileNameReturnFileId(recordReviewParams.DestId, newName, folderId)
		if err != nil {
			return err
		}
		lib.DPrintln("RecordReviewJobUsecase CopyFileNewFileName:", resFieldId)
		eventData := &ReminderClientUpdateFilesEventVo{
			Items: []*ReminderClientUpdateFilesEventVoItem{
				{
					BoxResId:   resFieldId,
					BoxResType: config_box.BoxResType_file,
					BoxResName: newName,
				},
			},
		}
		err = c.ReminderEventUsecase.AddClientUpdateFilesEvent(tCaseId, eventData)
		if err != nil {
			return err
		}
	}

	return nil
}

// BizHandleTaskV1 处理单个task
func (c *RecordReviewJobUsecase) BizHandleTaskV1(ctx context.Context, customTaskParams CustomTaskParams) error {

	lib.DPrintln("RecordReviewJobUsecase customTaskParams:", customTaskParams)

	var recordReviewParams RecordReviewParams
	err := json.Unmarshal([]byte(customTaskParams.Params), &recordReviewParams)
	if err != nil {
		return err
	}
	er := c.LogUsecase.SaveLog(recordReviewParams.ClientCaseId, Log_FromType_RecordReviewBizHandleTask, map[string]interface{}{
		"customTaskParams": customTaskParams,
	})
	if er != nil {
		c.log.Error(er)
	}

	tCase, err := c.TUsecase.DataById(Kind_client_cases, recordReviewParams.ClientCaseId)
	if err != nil {
		return err
	}
	if tCase == nil {
		return errors.New("tCase is nil")
	}

	hasLimit, err := c.CounterbuzUsecase.ClientUploadFilesLimit(tCase.CustomFields.NumberValueByNameBasic("id"), 500)
	if err != nil {
		c.log.Error(err)
		return err
	}

	if hasLimit {
		customTaskParamsBytes, _ := json.Marshal(customTaskParams)
		c.log.Error("RecordReviewJob_Sync_Limit: ", string(customTaskParamsBytes), " hasLimit: ", hasLimit)
		return nil
	}

	err = c.CounterbuzUsecase.ClientUploadFilesStat(tCase.CustomFields.NumberValueByNameBasic("id"))
	if err != nil {
		c.log.Error(err)
		return err
	}

	var mainCaseId int32

	// 此处
	// - 需要获取最新的名称：recordReviewParams.DestName， 因为客户有可能会改名
	// - 需要确认是否在正确的层级，因为客户有可能会移动文件夹或目录
	needHandle, resTypeMap, err := c.TidyDataV1(recordReviewParams)
	if err != nil {
		c.log.Error(err, " : ", InterfaceToString(recordReviewParams))
		return err
	}
	if !needHandle {
		return nil
	}

	entries := resTypeMap.GetTypeList("path_collection.entries")

	belongsClientFolder := false
	indexEntries := 0
	for i := 0; i < len(entries); i++ {
		id := entries[i].GetString("id")
		//lib.DPrintln("___:", id)
		// 255166311971: 为测试文件夹Test Clients
		if id == c.conf.Box.ClientFolderStructureParentId || id == c.conf.Box.ClientFolderStructureParentIdV2 || id == "255166311971" || id == "241109085470" {
			belongsClientFolder = true
			indexEntries = i
			break
		}
	}
	if !belongsClientFolder {
		c.log.Info("belongsClientFolder:", belongsClientFolder)
		return nil
	}
	if (indexEntries + 2) >= len(entries) {
		return nil
	}

	//clientFolder := entries[indexEntries+1]
	clientFolderFirstSubFolder := entries[indexEntries+2]
	lib.DPrintln(clientFolderFirstSubFolder)

	firstSubFolderName := clientFolderFirstSubFolder.GetString("name")
	if firstSubFolderName != config_box.FolderName_PrivateMedicalRecords &&
		firstSubFolderName != config_box.FolderName_VAMedicalRecords &&
		firstSubFolderName != config_box.FolderName_ServiceTreatmentRecords &&
		strings.Index(firstSubFolderName, "New Evidence") == -1 {
		c.log.Info("firstSubFolderName:", firstSubFolderName)
		return nil
	}

	var folderId string

	if strings.Index(firstSubFolderName, "New Evidence") >= 0 {

		NewEvidenceArr := strings.Split(firstSubFolderName, "#")
		if len(NewEvidenceArr) != 2 {
			return errors.New("NewEvidenceArr length is error : " + InterfaceToString(customTaskParams))
		}
		secondCaseId := lib.InterfaceToInt32(NewEvidenceArr[1])

		// 判断此人是否进行了文件拷贝
		key := MapKeyCopyRecordReviewFiles(secondCaseId)
		val, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if val == "" {
			c.log.Debug("CopyRecordReviewFiles:false", " secondCaseId: ", secondCaseId)
			return nil
		}

		secondCase, err := c.TUsecase.DataById(Kind_client_cases, secondCaseId)
		if err != nil {
			return err
		}
		if secondCase == nil {
			return errors.New("secondCase is nil")
		}
		mainCaseId = secondCase.CustomFields.NumberValueByNameBasic("id")
		dcRecordReviewFolderId, err := c.BoxbuzUsecase.DCRecordReviewFolderId(secondCase)
		if err != nil {
			return err
		}
		if dcRecordReviewFolderId == "" {
			return errors.New("dcRecordReviewFolderId is empty")
		}
		folderId = dcRecordReviewFolderId

		//if secondCase == nil {
		//	return errors.New("secondCase is nil")
		//}
		//
		//newEvidenceFolderId, err := c.BoxcUsecase.GetNewEvidenceFolderId(tCase, secondCase)
		//if err != nil {
		//	return err
		//}
		//folderId = newEvidenceFolderId
	} else {

		//currentCase, err := c.ClientCaseUsecase.CurrentCaseInProgress(tCase.CustomFields.TextValueByNameBasic(FieldName_client_gid))
		//if err != nil {
		//	return err
		//}
		currentCase := tCase
		caseId := currentCase.CustomFields.NumberValueByNameBasic("id")

		// 判断此人是否进行了文件拷贝
		key := MapKeyCopyRecordReviewFiles(caseId)
		val, err := c.MapUsecase.GetForString(key)
		if err != nil {
			return err
		}
		if val == "" {
			c.log.Debug("CopyRecordReviewFiles:false", " caseId: ", caseId)
			return nil
		}

		if currentCase == nil {
			c.log.Info("currentCase is nil")
			return nil
		}
		folderId, err = c.BoxdcUsecase.RecordReviewFirstSubFolderByName(firstSubFolderName, currentCase)
		if err != nil {
			return err
		}
		if folderId == "" {
			return errors.New("firstSubFolderName: folderId is empty " +
				InterfaceToString(currentCase.CustomFields.NumberValueByNameBasic("id")) +
				" firstSubFolderName: " + firstSubFolderName)
		}

		mainCaseId = caseId
	}

	sourceName := resTypeMap.GetString("name")
	sourceId := resTypeMap.GetString("id")
	var path string
	var destParentId string
	if (indexEntries + 3) >= len(entries) { // 此处说明要第一级子目录的文件
		destParentId = folderId
	} else {
		destParentId, path, err = c.BoxbuzUsecase.CreateFolderByEntries(indexEntries+3, resTypeMap.GetTypeList("path_collection.entries"), folderId)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}
	path = firstSubFolderName + "/" + path
	lib.DPrintln(path)
	fileNameExist, newFileId, err := c.BoxbuzUsecase.CopyFileToFolderNoCover(destParentId, sourceId, sourceName)
	if err != nil {
		return err
	}

	if fileNameExist {
		// 文件存在需要记录日志，后续分析处理
		c.log.Info("RecordReviewSyncExistsSameFileName fileNameExist:", fileNameExist, " path: ", path, " mainCaseId: ", mainCaseId, " sourceName: ", sourceName, " sourceId: ", sourceId)
		c.LogUsecase.SaveLog(mainCaseId, Log_FromType_RecordReviewSyncExistsSameFileName, map[string]interface{}{
			"tCaseId":    mainCaseId,
			"path":       path,
			"sourceId":   sourceId,
			"sourceName": sourceName,
		})
		return nil
	}

	lib.DPrintln(fileNameExist, newFileId, err)

	eventData := &ReminderClientUpdateFilesEventVo{
		Items: []*ReminderClientUpdateFilesEventVoItem{
			{
				BoxResId:       newFileId,
				BoxResType:     config_box.BoxResType_file,
				BoxResName:     sourceName,
				SourceBoxResId: sourceId,
				SourceBoxPath:  "[Clients]/" + path,
				BoxPath:        "[Data Collection]/Record Review/" + path,
			},
		},
	}
	err = c.ReminderEventUsecase.AddClientUpdateFilesEvent(mainCaseId, eventData)
	if err != nil {
		return err
	}

	return nil
}
