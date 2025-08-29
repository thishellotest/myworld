package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
)

type BoxWebhookLogEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	Remarks            string
	Headers            string
	Query              string
	Body               string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (BoxWebhookLogEntity) TableName() string {
	return "box_webhook_log"
}

func (c *BoxWebhookLogEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

type BoxWebhookLogUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[BoxWebhookLogEntity]
	BaseHandle[BoxWebhookLogEntity]
	TaskCreateUsecase        *TaskCreateUsecase
	ClientEnvelopeUsecase    *ClientEnvelopeUsecase
	BehaviorUsecase          *BehaviorUsecase
	TUsecase                 *TUsecase
	BoxbuzUsecase            *BoxbuzUsecase
	BoxUsecase               *BoxUsecase
	TaskFailureLogUsecase    *TaskFailureLogUsecase
	RecordReviewUsecase      *RecordReviewUsecase
	JotformbuzUsecase        *JotformbuzUsecase
	JotformSubmissionUsecase *JotformSubmissionUsecase
	FilebuzUsecase           *FilebuzUsecase
	MapUsecase               *MapUsecase
}

func NewBoxWebhookLogUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	TaskCreateUsecase *TaskCreateUsecase,
	ClientEnvelopeUsecase *ClientEnvelopeUsecase,
	BehaviorUsecase *BehaviorUsecase,
	TUsecase *TUsecase,
	BoxbuzUsecase *BoxbuzUsecase,
	BoxUsecase *BoxUsecase,
	TaskFailureLogUsecase *TaskFailureLogUsecase,
	RecordReviewUsecase *RecordReviewUsecase,
	JotformbuzUsecase *JotformbuzUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	FilebuzUsecase *FilebuzUsecase,
	MapUsecase *MapUsecase) *BoxWebhookLogUsecase {
	uc := &BoxWebhookLogUsecase{
		log:                      log.NewHelper(logger),
		CommonUsecase:            CommonUsecase,
		conf:                     conf,
		TaskCreateUsecase:        TaskCreateUsecase,
		ClientEnvelopeUsecase:    ClientEnvelopeUsecase,
		BehaviorUsecase:          BehaviorUsecase,
		TUsecase:                 TUsecase,
		BoxbuzUsecase:            BoxbuzUsecase,
		BoxUsecase:               BoxUsecase,
		TaskFailureLogUsecase:    TaskFailureLogUsecase,
		RecordReviewUsecase:      RecordReviewUsecase,
		JotformbuzUsecase:        JotformbuzUsecase,
		JotformSubmissionUsecase: JotformSubmissionUsecase,
		FilebuzUsecase:           FilebuzUsecase,
		MapUsecase:               MapUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()

	uc.BaseHandle.Log = log.NewHelper(logger)
	uc.TableName = BoxWebhookLogEntity{}.TableName()
	uc.BaseHandle.DB = CommonUsecase.DB()
	uc.BaseHandle.Handle = uc.Handle

	return uc
}

func (c *BoxWebhookLogUsecase) Handle(ctx context.Context, task *BoxWebhookLogEntity) error {

	if task == nil {
		return errors.New("task is nil.")
	}
	err := c.HandleExec(ctx, task)
	task.HandleStatus = HandleStatus_done
	task.UpdatedAt = time.Now().Unix() // 解决修改无更新有一次sql的问题
	if err != nil {
		task.HandleResult = HandleResult_failure
		task.AppendHandleResultDetail(err.Error())

		// 报警处理
		err = c.TaskFailureLogUsecase.Add(TaskType_BoxWebhookLog, 0,
			map[string]interface{}{
				"BoxWebhookLogId": task.ID,
				"err":             err.Error(),
			})

		if err != nil {
			c.log.Error(err)
		}

	} else {
		task.HandleResult = HandleResult_ok
	}
	return c.CommonUsecase.DB().Save(task).Error
}

const BoxQuestionnaireDownloadsFolderId = "258821922844"

// SIGN_REQUEST.COMPLETED,
// SIGN_REQUEST.DECLINED,
// SIGN_REQUEST.EXPIRED,
// SIGN_REQUEST.SIGNER_EMAIL_BOUNCED
// 28
// HandleExec {"type":"webhook_event","id":"2b4ded3a-2296-4d04-a8b7-f86005ced427","created_at":"2024-01-29T23:32:33-08:00","trigger":"SIGN_REQUEST.COMPLETED","webhook":{"id":"2383852458","type":"webhook"},"created_by":{"type":"user","id":"30888625898","name":"VBC Team","login":"team@vetbenefitscenter.com"},"source":{"id":"1429599030083","type":"file","file_version":{"type":"file_version","id":"1567474367100","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35"},"sequence_id":"2","etag":"2","sha1":"9d83c850476f1f6a2c7226d8d798d498d3d03c35","name":"Your Veteran Benefits Center Contract (3).pdf","description":"","size":452587,"path_collection":{"total_count":3,"entries":[{"type":"folder","id":"0","sequence_id":null,"etag":null,"name":"All Files"},{"type":"folder","id":"241183180615","sequence_id":"4","etag":"4","name":"VBC Engineering Team"},{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"}]},"created_at":"2024-01-29T23:25:43-08:00","modified_at":"2024-01-29T23:32:28-08:00","trashed_at":null,"purged_at":null,"content_created_at":"2024-01-29T23:25:43-08:00","content_modified_at":"2024-01-29T23:32:28-08:00","created_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"modified_by":{"type":"user","id":"16371441643","name":"Box Sign","login":"AutomationUser_1519487_GBsgja6E9G@boxdevedition.com"},"owned_by":{"type":"user","id":"30690179025","name":"Yannan Wang","login":"ywang@vetbenefitscenter.com"},"shared_link":null,"parent":{"type":"folder","id":"246205309773","sequence_id":"1","etag":"1","name":"Test Box Sign Requests"},"item_status":"active"},"additional_info":{"sign_request_id":"b1685f64-b285-4786-abc0-174931e7be0e","signer_emails":["team@vetbenefitscenter.com","lialing@foxmail.com","liaogling@gmail.com"],"external_id":null}}
func (c *BoxWebhookLogUsecase) HandleExec(ctx context.Context, task *BoxWebhookLogEntity) error {

	if task == nil {
		return errors.New("task is nil")
	}
	bodyMap := lib.ToTypeMapByString(task.Body)
	trigger := bodyMap.GetString("trigger")
	if trigger == "SIGN_REQUEST.COMPLETED" {
		signRequestId := bodyMap.GetString("additional_info.sign_request_id")
		if signRequestId == "" {
			return errors.New("signRequestId is empty.")
		}

		clientEnvelope, err := c.ClientEnvelopeUsecase.GetByEnvelopeId(EsignVendor_box, signRequestId)
		if err != nil {
			return err
		}
		if clientEnvelope == nil {
			return errors.New("clientEnvelope is nil")
		}
		BehaviorType, err := ClientEnvelopeTypeToBehaviorType(clientEnvelope.Type)
		if err != nil {
			return err
		}
		clientEnvelope.IsSigned = ClientEnvelope_IsSigned_Yes
		clientEnvelope.UpdatedAt = time.Now().Unix()
		er := c.CommonUsecase.DB().Save(&clientEnvelope).Error
		if er != nil {
			c.log.Error(er)
		}
		return c.BehaviorUsecase.Add(clientEnvelope.ClientId, BehaviorType, time.Now(), "")
		//return c.TaskCreateUsecase.CreateBoxCreateFolderForNewClientTask(&CreateBoxCreateFolderForNewClientTask{
		//	ClientId: clientEnvelope.ClientId,
		//})
	} else if trigger == "FOLDER.CREATED" {
		parentId := bodyMap.GetString("source.parent.id")
		// 258821922844：是DataCollection Backup

		if parentId == BoxQuestionnaireDownloadsFolderId && !configs.EnableNewJotformName {
			if bodyMap.GetString("source.type") == "folder" {
				sourceId := bodyMap.GetString("source.id")
				folderName := bodyMap.GetString("source.name")

				if strings.Index(folderName, "#") >= 0 {
					folderArr := strings.Split(folderName, "#")
					if len(folderArr) != 2 {
						return errors.New("folderName: " + folderName + " format is wrong.")
					}
					clientCaseIdStr := strings.TrimSpace(folderArr[1])
					newFolderName := strings.TrimSpace(folderArr[0])
					clientCaseIdInt64, _ := strconv.ParseInt(clientCaseIdStr, 10, 32)
					clientCaseId := int32(clientCaseIdInt64)

					tCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
					if err != nil {
						return err
					}
					if tCase == nil {
						return errors.New("tCase is nil.")
					}
					QuestionnairesFolderId, err := c.BoxbuzUsecase.GetDCSubFolderId(MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId), tCase)
					if err != nil {
						return err
					}
					if QuestionnairesFolderId == "" {
						return errors.New("QuestionnairesFolderId is empty.")
					}

					// 处理文件名是否冲突的问题
					needMove, err := c.HandleQuestionnairesBuz(sourceId, newFolderName, QuestionnairesFolderId, tCase)
					if err != nil {
						return err
					}
					if needMove {
						_, err = c.BoxUsecase.MoveFolderName(sourceId, newFolderName, QuestionnairesFolderId)
						if err != nil {
							return err
						} else {
							// 处理jotform数据入DB
							items, err := c.BoxUsecase.ListItemsInFolderFormat(sourceId)
							if err != nil {
								return err
							}
							for _, v := range items {
								name := v.GetString("name")
								c.log.Info("HandleQuestionnairesJotform:", name)
								err := c.HandleQuestionnairesJotform(name, tCase)
								if err != nil {
									c.log.Error(err, InterfaceToString(v), " clientCaseId:", clientCaseId)
								}
							}
						}
					}
				}
			}
		}
	} else if trigger == "FILE.UPLOADED" {
		if configs.EnableNewJotformName {
			err := c.HandleNewJotformName(bodyMap)
			if err != nil {
				c.log.Error(err, " BoxWebhookLogId: ", task.ID)
			}
		}
	}
	// 只有创建文件夹，没什么意义
	if trigger == "FOLDER.MOVED" || trigger == "FOLDER.COPIED" || trigger == "FILE.UPLOADED" || trigger == "FILE.RENAMED" {
		needHandle, err := c.RecordReviewUsecase.Process(ctx, bodyMap, task.ID)
		c.log.Info("RecordReviewUsecase_Process 1: ", needHandle, err)
	}

	err := c.FilebuzUsecase.DCRecordReviewFileHandle(bodyMap)
	if err != nil {
		c.log.Error(err, " id: ", task.ID)
	}
	return nil
}

func GetJotformIdFromBoxFileName(boxFileName string) (submissionId string, err error) {
	strs := strings.Split(boxFileName, ".")
	if len(strs) != 2 {
		return "", errors.New("GetJotformIdFromBoxFileName error : " + boxFileName)
	}
	if strs[1] != "pdf" {
		return "", errors.New("GetJotformIdFromBoxFileName error " + boxFileName)
	}
	return strs[0], nil
}

func (c *BoxWebhookLogUsecase) HandleNewJotformName(bodyMap lib.TypeMap) error {
	boxFileId := bodyMap.GetString("source.id")
	boxFileName := bodyMap.GetString("source.name")
	entries := bodyMap.GetTypeList("source.path_collection.entries")
	if len(entries) < 2 {
		return nil
	}
	//"modified_at": "2025-06-04T10:41:58-07:00",
	boxFileModifiedAt := bodyMap.GetString("source.modified_at")
	layout := time.RFC3339
	//str := "2025-06-04T10:41:58-07:00"
	boxFileModifiedTime, err := time.Parse(layout, boxFileModifiedAt)
	if err != nil {
		return err
	}

	maybeQuestionnaireDownloads := entries[len(entries)-2]
	if maybeQuestionnaireDownloads.GetString("type") == "folder" &&
		maybeQuestionnaireDownloads.GetString("id") == BoxQuestionnaireDownloadsFolderId {
		boxJotformFolder := entries[len(entries)-1]
		folderName := boxJotformFolder.GetString("name")
		folderArr := strings.Split(folderName, "#")
		if len(folderArr) != 2 {
			return errors.New("folderName: " + folderName + " format is wrong.")
		}
		clientCaseIdStr := strings.TrimSpace(folderArr[1])
		clientCaseIdInt64, _ := strconv.ParseInt(clientCaseIdStr, 10, 32)
		clientCaseId := int32(clientCaseIdInt64)
		tCase, err := c.TUsecase.DataById(Kind_client_cases, clientCaseId)
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil.")
		}

		var submissionId string
		// 特殊处理解决jotform pdf没有ID的问题
		if strings.Index(folderName, "Muscle Injuries") >= 0 ||
			strings.Index(folderName, "Chronic Fatigue Syndrome") >= 0 ||
			strings.Index(folderName, "Stomach and Duodenal Conditions") >= 0 {
			submissionId = GetJotformSubmissionsIdFromFolderName(folderName)
		} else {
			submissionId, err = GetJotformIdFromBoxFileName(boxFileName)
			if err != nil {
				return err
			}
		}
		err = c.JotformbuzUsecase.HandleSubmission(submissionId, tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode), "")
		if err != nil {
			return err
		}
		err = c.HandleNewJotformNameStep2Copy(boxFileId, boxFileModifiedTime, submissionId, tCase)
		if err != nil {
			return err
		}

		// 判断文件夹是否为空， 为空时直接删除
		boxJotformFolderId := boxJotformFolder.GetString("id")
		destFolderEntities, err := c.BoxUsecase.ListItemsInFolderFormat(boxJotformFolderId)
		if err != nil {
			return err
		}
		if len(destFolderEntities) == 0 {
			_, err = c.BoxUsecase.DeleteFolder(boxJotformFolderId, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetJotformSubmissionsIdFromFolderName(folderName string) (jotformSubmissionsId string) {
	//text := "Test1 TestL Muscle Injuries Increase 2025-03-14 19:48:10 -6178160905016488647# 5511"

	text := folderName
	// 找到 '#' 的索引
	idx := strings.LastIndex(text, "#")
	if idx == -1 {
		//fmt.Println("No '#' found")
		return
	}

	// 从 '#' 左侧开始截取
	leftPart := text[:idx]

	// 按空格拆分
	words := strings.Split(leftPart, "-")

	// 获取最后一个数字
	if len(words) > 0 {
		//fmt.Println("Extracted number:", words[len(words)-1])
		return strings.TrimSpace(words[len(words)-1])
	}
	return
}

func (c *BoxWebhookLogUsecase) HandleNewJotformNameStep2Copy(boxFileId string, boxFileModifiedTime time.Time, jotformSubmissionId string, tCase *TData) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	if boxFileId == "" {
		return errors.New("boxFileId is empty")
	}
	if jotformSubmissionId == "" {
		return errors.New("jotformSubmissionId is empty")
	}
	clientCaseId := tCase.Id()
	var destFolderId string
	var err error

	jotformSubmisstionEntity, err := c.JotformSubmissionUsecase.GetLatestFormInfo(jotformSubmissionId)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if jotformSubmisstionEntity.FormId == QuestionnairesUpdateQuestionnaire_FormId {
		destFolderId, err = c.BoxbuzUsecase.DCGetOrMakeUpdateQuestionnairesFolderId(*tCase)
		if err != nil {
			return err
		}
	} else {
		destFolderId, err = c.BoxbuzUsecase.GetDCSubFolderId(MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId), tCase)
		if err != nil {
			return err
		}
	}
	if destFolderId == "" {
		return errors.New("destFolderId is empty.")
	}

	newFileName, err := GenJotformNewFileNameForBox(jotformSubmisstionEntity, tCase)
	if err != nil {
		return err
	}
	c.log.Info("newFileName:", newFileName)

	destFolderEntities, err := c.BoxUsecase.ListItemsInFolderFormat(destFolderId)
	if err != nil {
		return err
	}

	existBoxFileId := ""
	for _, v1 := range destFolderEntities {
		if v1.GetString("name") == newFileName && v1.GetString("type") == "file" {
			existBoxFileId = v1.GetString("id")
			break
		}
	}
	key := MapKeyQuestionnairesFileTime(jotformSubmissionId)
	if existBoxFileId != "" {
		boxFileModifiedTimeUnix := boxFileModifiedTime.Unix()
		currentBoxTimeUnix, err := c.MapUsecase.GetForInt(key)
		if err != nil {
			return err
		}
		if currentBoxTimeUnix == 0 || boxFileModifiedTimeUnix > int64(currentBoxTimeUnix) {
			fileReader, err := c.BoxUsecase.DownloadFile(boxFileId, "")
			if err != nil {
				return errors.New("BoxUsecase.DownloadFile: " + err.Error())
			}
			if fileReader == nil {
				return errors.New("BoxUsecase.DownloadFile: fileReader is nil")
			}
			defer fileReader.Close()
			// 上传新版本
			_, err = c.BoxUsecase.UploadFileVersion(existBoxFileId, fileReader)
			if err != nil {
				return errors.New("BoxUsecase.UploadFileVersion: " + err.Error())
			}
			c.MapUsecase.Set(key, InterfaceToString(boxFileModifiedTime.Unix()))
			c.BoxUsecase.DeleteFile(boxFileId)
		} else {
			c.BoxUsecase.DeleteFile(boxFileId)
		}
	} else {
		_, err = c.BoxUsecase.MoveFileWithNewName(boxFileId, destFolderId, newFileName)
		if err != nil {
			return err
		}
		c.MapUsecase.Set(key, InterfaceToString(boxFileModifiedTime.Unix()))
	}

	return nil
}

func GenJotformNewFileNameForBox(jotformSubmisstionEntity *JotformSubmissionEntity, tCase *TData) (newFileName string, err error) {
	if jotformSubmisstionEntity == nil {
		return "", errors.New("jotformSubmisstionEntity is nil")
	}
	if tCase == nil {
		return "", errors.New("tCase is nil")
	}
	questionnairesItem := GetQuestionnairesItemByFormId(jotformSubmisstionEntity.FormId)
	if questionnairesItem == nil {
		return "", errors.New("questionnairesItem is nil")
	}
	//if len(questionnairesItem.FileNames) == 0 {
	//	return "", errors.New("questionnairesItem.FileNames is not config : " + jotformSubmisstionEntity.FormId)
	//}
	notesMap := lib.ToTypeMapByString(jotformSubmisstionEntity.Notes)
	newFileName = fmt.Sprintf("%d-%s", tCase.Id(), questionnairesItem.BaseTitle)
	for _, v := range questionnairesItem.FileNames {
		name := notesMap.GetString(v)
		name = BoxFileNameFilter(name)
		name = lib.TruncateStringWithRune(name, 80)
		newFileName += "-" + name
	}
	newFileName += "-" + jotformSubmisstionEntity.SubmissionId + ".pdf"
	newFileName = BoxFileNameFilter(newFileName)
	return newFileName, nil
}

func GenJotformNewFileNameForAI(jotformSubmisstionEntity *JotformSubmissionEntity) (newFileName string, err error) {
	if jotformSubmisstionEntity == nil {
		return "", errors.New("jotformSubmisstionEntity is nil")
	}
	questionnairesItem := GetQuestionnairesItemByFormId(jotformSubmisstionEntity.FormId)
	if questionnairesItem == nil {
		return "", errors.New("questionnairesItem is nil")
	}
	//if len(questionnairesItem.FileNames) == 0 {
	//	return "", errors.New("questionnairesItem.FileNames is not config : " + jotformSubmisstionEntity.FormId)
	//}
	notesMap := lib.ToTypeMapByString(jotformSubmisstionEntity.Notes)
	newFileName = fmt.Sprintf("%s", questionnairesItem.BaseTitle)
	for _, v := range questionnairesItem.FileNames {
		newFileName += "-" + notesMap.GetString(v)
	}
	newFileName += "-" + jotformSubmisstionEntity.SubmissionId + ".pdf"
	newFileName = BoxFileNameFilter(newFileName)
	return newFileName, nil
}

func BoxFileNameFilter(fieldName string) string {

	fieldName = strings.ReplaceAll(fieldName, "\n", "")
	fieldName = strings.ReplaceAll(fieldName, "\t", "")
	fieldName = strings.ReplaceAll(fieldName, "\r", "")
	fieldName = strings.ReplaceAll(fieldName, "<", "")
	fieldName = strings.ReplaceAll(fieldName, ">", "")
	fieldName = strings.ReplaceAll(fieldName, "\"", "")
	fieldName = strings.ReplaceAll(fieldName, ":", "")
	fieldName = strings.ReplaceAll(fieldName, "|", "")
	fieldName = strings.ReplaceAll(fieldName, "/", "")
	fieldName = strings.ReplaceAll(fieldName, "\\", "")
	fieldName = strings.ReplaceAll(fieldName, "?", "")
	fieldName = strings.ReplaceAll(fieldName, "*", "")
	return fieldName
}

func (c *BoxWebhookLogUsecase) HandleQuestionnairesJotform(pdfName string, tCase *TData) error {
	if tCase == nil {
		return errors.New("tCase is nil")
	}
	submissionId, _ := lib.FileExt(pdfName, false)
	if submissionId == "" {
		return errors.New("submissionId is empty")
	}
	return c.JotformbuzUsecase.HandleSubmission(submissionId, tCase.CustomFields.TextValueByNameBasic(FieldName_uniqcode), "")
	//return nil
}

// HandleQuestionnairesBuz 只支持文件
func (c *BoxWebhookLogUsecase) HandleQuestionnairesBuz(srcFolderId string, srcFolderName string,
	QuestionnairesFolderId string, tCase *TData) (needMove bool, err error) {

	if tCase == nil {
		return false, errors.New("tCase is nil")
	}
	// 判断文件夹名是否重复
	sameId, err := c.BoxbuzUsecase.SameNameFolderOrFile("folder", srcFolderName, QuestionnairesFolderId)
	if err != nil {
		return false, err
	}
	if sameId == "" { // 不重复走移动操作
		return true, nil
	}
	srcFolderFiles, err := c.BoxUsecase.ListItemsInFolderFormat(srcFolderId)
	if err != nil {
		return false, err
	}

	destFolderEntities, err := c.BoxUsecase.ListItemsInFolderFormat(sameId)
	if err != nil {
		return false, err
	}

	for _, v := range srcFolderFiles {
		time.Sleep(1 * time.Second)
		v.GetString("id")
		v.GetString("name")
		v.GetString("type")

		isConflict := false
		for _, v1 := range destFolderEntities {
			if v1.GetString("name") == v.GetString("name") &&
				v1.GetString("type") == v.GetString("type") { // file or folder conflict
				if v.GetString("type") == "file" { //
					isConflict = true

					fileReader, err := c.BoxUsecase.DownloadFile(v.GetString("id"), "")
					if err != nil {
						return false, errors.New("BoxUsecase.DownloadFile: " + err.Error())
					}
					if fileReader == nil {
						return false, errors.New("BoxUsecase.DownloadFile: fileReader is nil")
					}
					defer fileReader.Close()
					// 上传新版本
					_, err = c.BoxUsecase.UploadFileVersion(v1.GetString("id"), fileReader)
					if err != nil {
						return false, errors.New("BoxUsecase.UploadFileVersion: " + err.Error())
					} else {
						pdfName := v.GetString("name")
						c.log.Info("HandleQuestionnairesJotform 1:", " pdfName:", pdfName, " caseId: ", tCase.Id())
						err = c.HandleQuestionnairesJotform(pdfName, tCase)
						if err != nil {
							c.log.Error("HandleQuestionnairesJotform error1:", pdfName, " caseId: ", tCase.Id())
						}
					}
					break
				}
			}
		}
		if !isConflict && v.GetString("type") == "file" { // 移动操作
			_, err := c.BoxUsecase.MoveFile(v.GetString("id"), sameId)
			if err != nil {
				return false, err
			} else {
				pdfName := v.GetString("name")
				c.log.Info("HandleQuestionnairesJotform 2:", " pdfName:", pdfName, " caseId: ", tCase.Id())
				err = c.HandleQuestionnairesJotform(pdfName, tCase)
				if err != nil {
					c.log.Error("HandleQuestionnairesJotform error2:", pdfName, " caseId: ", tCase.Id())
				}
			}
		}
	}

	// 处理完华后，把冲突的文件夹删除
	_, err = c.BoxUsecase.DeleteFolder(srcFolderId, true)
	if err != nil {
		return false, errors.New("BoxUsecase.DeleteFolder: " + err.Error())
	}

	return false, nil
}

func (c *BoxWebhookLogUsecase) CrontabEveryOneHourHandleQuestionnaireDownloads() error {
	c.log.Info("CrontabEveryOneHourHandleQuestionnaireDownloads")
	res, err := c.BoxUsecase.ListItemsInFolderFormat(BoxQuestionnaireDownloadsFolderId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	c.log.Info("CrontabEveryOneHourHandleQuestionnaireDownloads", len(res))
	for _, v := range res {
		if v.GetString("type") == "folder" {
			temp, err := c.BoxUsecase.ListItemsInFolderFormat(v.GetString("id"))
			if err != nil {
				c.log.Error(err)
			} else {
				for _, v1 := range temp {
					fileInfo, _, err := c.BoxUsecase.GetFileInfo(v1.GetString("id"))
					if err != nil {
						c.log.Error(err)
						continue
					}
					if fileInfo == nil {
						continue
					}
					fileInfoMap := lib.ToTypeMapByString(*fileInfo)
					newFileInfoMap := make(lib.TypeMap)
					newFileInfoMap.Set("source", fileInfoMap)
					err = c.HandleNewJotformName(newFileInfoMap)
					if err != nil {
						c.log.Error(err)
					}
				}
			}
		}
	}
	return nil
}
