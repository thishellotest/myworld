package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"vbc/internal/conf"
	"vbc/internal/config_box"
	"vbc/lib"
	. "vbc/lib/builder"
)

var boxWebhookStr = `{
    "type": "webhook_event",
    "id": "f214d32f-6e36-4550-873b-2f8e86a0966d",
    "created_at": "2024-06-06T01:19:50-07:00",
    "trigger": "FILE.UPLOADED",
    "webhook": {
        "id": "2802339849",
        "type": "webhook"
    },
    "created_by": {
        "type": "user",
        "id": "30888625898",
        "name": "VBC Team",
        "login": "info@vetbenefitscenter.com"
    },
    "source": {
        "id": "1552371904380",
        "type": "file",
        "file_version": {
            "type": "file_version",
            "id": "1705561482780",
            "sha1": "da39a3ee5e6b4b0d3255bfef95601890afd80709"
        },
        "sequence_id": "0",
        "etag": "0",
        "sha1": "da39a3ee5e6b4b0d3255bfef95601890afd80709",
        "name": "\u7b14\u8bb0\u672a\u8bbe\u7f6e\u6807\u9898 2024-06-06 16.19.49.boxnote",
        "description": "",
        "size": 0,
        "path_collection": {
            "total_count": 7,
            "entries": [{
                "type": "folder",
                "id": "0",
                "sequence_id": null,
                "etag": null,
                "name": "All Files"
            }, {
                "type": "folder",
                "id": "241183180615",
                "sequence_id": "4",
                "etag": "4",
                "name": "VBC Engineering Team"
            }, {
                "type": "folder",
                "id": "247457173873",
                "sequence_id": "1",
                "etag": "1",
                "name": "Testing"
            }, {
                "type": "folder",
                "id": "255166311971",
                "sequence_id": "1",
                "etag": "1",
                "name": "Test Clients"
            }, {
                "type": "folder",
                "id": "264686374897",
                "sequence_id": "2",
                "etag": "2",
                "name": "[PROD]VBC - TestLiao, TestGary #5076"
            }, {
                "type": "folder",
                "id": "264686394097",
                "sequence_id": "0",
                "etag": "0",
                "name": "VA Medical Records"
            }, {
                "type": "folder",
                "id": "268258622342",
                "sequence_id": "0",
                "etag": "0",
                "name": "leve1_folder"
            }]
        },
        "created_at": "2024-06-06T01:19:50-07:00",
        "modified_at": "2024-06-06T01:19:50-07:00",
        "trashed_at": null,
        "purged_at": null,
        "content_created_at": "2024-06-06T01:19:50-07:00",
        "content_modified_at": "2024-06-06T01:19:50-07:00",
        "created_by": {
            "type": "user",
            "id": "30888625898",
            "name": "VBC Team",
            "login": "info@vetbenefitscenter.com"
        },
        "modified_by": {
            "type": "user",
            "id": "30888625898",
            "name": "VBC Team",
            "login": "info@vetbenefitscenter.com"
        },
        "owned_by": {
            "type": "user",
            "id": "30690179025",
            "name": "Yannan Wang",
            "login": "ywang@vetbenefitscenter.com"
        },
        "shared_link": null,
        "parent": {
            "type": "folder",
            "id": "268258622342",
            "sequence_id": "0",
            "etag": "0",
            "name": "leve1_folder"
        },
        "item_status": "active"
    },
    "additional_info": []
}`

func aaa() {
}

type RecordReviewUsecase struct {
	log                    *log.Helper
	CommonUsecase          *CommonUsecase
	conf                   *conf.Data
	MapUsecase             *MapUsecase
	RecordReviewJobUsecase *RecordReviewJobUsecase
	LogUsecase             *LogUsecase
}

func NewRecordReviewUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	MapUsecase *MapUsecase,
	RecordReviewJobUsecase *RecordReviewJobUsecase,
	LogUsecase *LogUsecase) *RecordReviewUsecase {
	uc := &RecordReviewUsecase{
		log:                    log.NewHelper(logger),
		CommonUsecase:          CommonUsecase,
		conf:                   conf,
		MapUsecase:             MapUsecase,
		RecordReviewJobUsecase: RecordReviewJobUsecase,
		LogUsecase:             LogUsecase,
	}
	return uc
}

// Process 处理webhook数所
// needHandle：true 需要处理，放入队列；false: 不需要处理；
func (c *RecordReviewUsecase) Process(ctx context.Context, boxWebhookTypeMap lib.TypeMap, sourceFromId int32) (needHandle bool, err error) {

	if c.conf.Box.SyncRecordReviewVersion == config_box.SyncRecordReviewVersionV1 {
		return c.ProcessV1(ctx, boxWebhookTypeMap, sourceFromId)
	}

	if boxWebhookTypeMap == nil {
		return false, errors.New("boxWebhookTypeMap is nil")
	}
	if boxWebhookTypeMap.GetString("source.type") != string(config_box.BoxResType_folder) &&
		boxWebhookTypeMap.GetString("source.type") != string(config_box.BoxResType_file) {
		return false, nil
	}
	trigger := boxWebhookTypeMap.GetString("trigger")
	c.log.Info("RecordReviewUsecase_Process trigger:", trigger)
	// 更新也使用FILE.UPLOADED
	if trigger != "FILE.UPLOADED" && trigger != "FILE.RENAMED" &&
		trigger != "FOLDER.CREATED" && trigger != "FOLDER.RENAMED" &&
		trigger != "FOLDER.MOVED" && trigger != "FOLDER.COPIED" {
		return false, nil
	}

	entries := boxWebhookTypeMap.GetTypeList("source.path_collection.entries")
	belongsClientFolder := false
	indexEntries := 0
	for i := 0; i < len(entries); i++ {
		id := entries[i].GetString("id")
		// 255166311971: 为测试文件夹Test Clients
		if id == c.conf.Box.ClientFolderStructureParentId || id == c.conf.Box.ClientFolderStructureParentIdV2 || id == "255166311971" {
			belongsClientFolder = true
			indexEntries = i
			break
		}
	}
	c.log.Info("RecordReviewUsecase_Process belongsClientFolder:", belongsClientFolder)
	if !belongsClientFolder {
		return false, nil
	}
	if (indexEntries + 2) >= len(entries) {
		return false, nil
	}

	clientFolder := entries[indexEntries+1]
	clientFolderFirstSubFolder := entries[indexEntries+2]

	firstSubFolderName := clientFolderFirstSubFolder.GetString("name")

	c.log.Info(firstSubFolderName)
	if firstSubFolderName != config_box.FolderName_PrivateMedicalRecords &&
		firstSubFolderName != config_box.FolderName_VAMedicalRecords &&
		firstSubFolderName != config_box.FolderName_ServiceTreatmentRecords {
		return false, nil
	}
	folderMap, err := c.MapUsecase.GetByCond(And(Eq{"mval": clientFolder.Get("id")}, Like{"mkey", "ClientBoxFolderId:%"}))
	if err != nil {
		return false, err
	}
	if folderMap == nil {
		return false, errors.New("Process: folderMap is nil")
	}

	mkeyRes := strings.Split(folderMap.Mkey, ":")
	clientCaseId := lib.InterfaceToInt32(mkeyRes[1])

	if clientCaseId != 5076 {
		c.log.Info("RecordReviewUsecase_Process 2:非测试帐号5076")
		return false, nil
	}

	// 判断此人是否进行了文件拷贝
	key := MapKeyCopyRecordReviewFiles(clientCaseId)
	val, err := c.MapUsecase.GetForString(key)
	if err != nil {
		return false, err
	}
	if val == "" {
		c.log.Info("RecordReviewUsecase_Process  CopyRecordReviewFiles:false")
		return false, nil
	}

	recordReviewParams := &RecordReviewParams{
		ClientCaseId:       clientCaseId,
		FirstSubFolderId:   clientFolderFirstSubFolder.GetString("id"),
		FirstSubFolderName: firstSubFolderName,
		//DestId:             clientFolderDestFolder.GetString("id"),
		//DestType:           config_box.BoxResType(clientFolderDestFolder.GetString("type")),
		//DestName:           clientFolderDestFolder.GetString("name"),
	}

	if (indexEntries + 3) < len(entries) { // 说明是文件夹
		clientFolderDestFolder := entries[indexEntries+3]
		recordReviewParams.DestId = clientFolderDestFolder.GetString("id")
		recordReviewParams.DestName = clientFolderDestFolder.GetString("name")
		recordReviewParams.DestType = config_box.BoxResType(clientFolderDestFolder.GetString("type"))
	} else {
		recordReviewParams.DestId = boxWebhookTypeMap.GetString("source.id")
		recordReviewParams.DestName = boxWebhookTypeMap.GetString("source.name")
		recordReviewParams.DestType = config_box.BoxResType_file
	}

	//
	//lib.DPrintln("recordReviewParams:", recordReviewParams)
	recordReviewParamsBytes, _ := json.Marshal(recordReviewParams)

	uniqueKey := fmt.Sprintf("%d:%s:%s", recordReviewParams.ClientCaseId, recordReviewParams.DestId, recordReviewParams.DestType)
	customTaskParams := CustomTaskParams{
		UniqueKey: uniqueKey,
		Params:    string(recordReviewParamsBytes),
	}
	//lib.DPrintln("RecordReviewUsecase_Process LPushCustomTaskQueue:", customTaskParams)
	err = c.RecordReviewJobUsecase.LPushCustomTaskQueue(ctx, customTaskParams)
	er := c.LogUsecase.SaveLog(clientCaseId, Log_FromType_RecordReviewLPushRedisQueue, map[string]interface{}{
		"customTaskParams": customTaskParams,
	})
	if er != nil {
		c.log.Error(er)
	}

	//lib.DPrintln(folderMap, err)
	//lib.DPrintln(clientFolder)
	//lib.DPrintln(entries[indexEntries+2])
	//lib.DPrintln(clientFolderDestFolder)

	if err != nil {
		return false, err
	}

	return true, nil
}

// ProcessV1 处理webhook数所
// needHandle：true 需要处理，放入队列；false: 不需要处理；
func (c *RecordReviewUsecase) ProcessV1(ctx context.Context, boxWebhookTypeMap lib.TypeMap, sourceFromId int32) (needHandle bool, err error) {
	if boxWebhookTypeMap == nil {
		return false, errors.New("boxWebhookTypeMap is nil")
	}
	if boxWebhookTypeMap.GetString("source.type") != string(config_box.BoxResType_file) {
		return false, nil
	}
	trigger := boxWebhookTypeMap.GetString("trigger")
	lib.DPrintln("RecordReviewUsecase_Process trigger:", trigger)
	// 更新也使用FILE.UPLOADED 新版本是考虑文件的修改
	if trigger != "FILE.UPLOADED" && trigger != "FILE.RENAMED" {
		return false, nil
	}

	entries := boxWebhookTypeMap.GetTypeList("source.path_collection.entries")
	belongsClientFolder := false
	indexEntries := 0
	for i := 0; i < len(entries); i++ {
		id := entries[i].GetString("id")
		// 255166311971: 为测试文件夹Test Clients
		if id == c.conf.Box.ClientFolderStructureParentId || id == c.conf.Box.ClientFolderStructureParentIdV2 || id == "255166311971" {
			belongsClientFolder = true
			indexEntries = i
			break
		}
	}
	lib.DPrintln("RecordReviewUsecase_Process belongsClientFolder:", belongsClientFolder)
	if !belongsClientFolder {
		return false, nil
	}
	if (indexEntries + 2) >= len(entries) {
		return false, nil
	}

	clientFolder := entries[indexEntries+1]
	clientFolderFirstSubFolder := entries[indexEntries+2]

	firstSubFolderName := clientFolderFirstSubFolder.GetString("name")

	lib.DPrintln(firstSubFolderName)
	if firstSubFolderName != config_box.FolderName_PrivateMedicalRecords &&
		firstSubFolderName != config_box.FolderName_VAMedicalRecords &&
		firstSubFolderName != config_box.FolderName_ServiceTreatmentRecords &&
		strings.Index(firstSubFolderName, "New Evidence") == -1 {
		return false, nil
	}
	folderMap, err := c.MapUsecase.GetByCond(And(Eq{"mval": clientFolder.Get("id")}, Like{"mkey", "ClientBoxFolderId:%"}))
	if err != nil {
		return false, err
	}
	if folderMap == nil {
		return false, errors.New("Process: folderMap is nil")
	}

	mkeyRes := strings.Split(folderMap.Mkey, ":")
	clientCaseId := lib.InterfaceToInt32(mkeyRes[1])

	//if clientCaseId != 5076 {
	//	lib.DPrintln("RecordReviewUsecase_Process 2:非测试帐号5076")
	//	return false, nil
	//}

	recordReviewParams := &RecordReviewParams{
		ClientCaseId:       clientCaseId,
		FirstSubFolderId:   clientFolderFirstSubFolder.GetString("id"),
		FirstSubFolderName: firstSubFolderName,
		SourceFromId:       sourceFromId,
	}

	recordReviewParams.DestId = boxWebhookTypeMap.GetString("source.id")
	recordReviewParams.DestName = boxWebhookTypeMap.GetString("source.name")
	recordReviewParams.DestType = config_box.BoxResType_file

	c.log.Debug("recordReviewParams:", recordReviewParams)
	recordReviewParamsBytes, _ := json.Marshal(recordReviewParams)

	uniqueKey := fmt.Sprintf("%d:%s:%s", recordReviewParams.ClientCaseId, recordReviewParams.DestId, recordReviewParams.DestType)
	customTaskParams := CustomTaskParams{
		UniqueKey: uniqueKey,
		Params:    string(recordReviewParamsBytes),
	}
	c.log.Debug("RecordReviewUsecase_Process LPushCustomTaskQueue:", customTaskParams)
	err = c.RecordReviewJobUsecase.LPushCustomTaskQueue(ctx, customTaskParams)
	er := c.LogUsecase.SaveLog(clientCaseId, Log_FromType_RecordReviewLPushRedisQueue, map[string]interface{}{
		"customTaskParams": customTaskParams,
	})
	if er != nil {
		c.log.Error(er)
	}

	//lib.DPrintln(folderMap, err)
	//lib.DPrintln(clientFolder)
	//lib.DPrintln(entries[indexEntries+2])
	//lib.DPrintln(clientFolderDestFolder)

	if err != nil {
		return false, err
	}

	return true, nil
}
