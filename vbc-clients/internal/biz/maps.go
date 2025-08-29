package biz

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
	"vbc/lib"
)

/*

CREATE TABLE `maps` (
  `mkey` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
  `mval` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '',
  UNIQUE KEY `uniq_k` (`mkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

*/

const (
	Map_BuildAuto                       = "BuildAuto:"
	Map_AdobeWebhookEvent_divide        = "AdobeWebhookEvent_divide"
	Map_EnvelopeStatusChange_divide     = "EnvelopeStatusChange_divide"
	Map_Change_histories_divide         = "Change_histories_divide"
	Map_Change_histories_divide_nodelay = "Change_histories_divide_nodelay"
	Map_mail_FeeScheduleCommunication   = "MailFeeScheduleCommunication:"
	Map_CreateEnvelope                  = "CreateEnvelope:"
	//Map_CreateFolderInBox             = "CreateFolderInBox:"
	Map_docusignTpl                    = "docusignTpl:"
	Map_adobesignTpl                   = "adobesignTpl:"
	Map_boxsignTpl                     = "boxsignTpl:"
	Map_ClientContractBoxFolderId      = "ClientContractBoxFolderId:"
	Map_HandleCreateFolderInBoxAndMail = "HandleCreateFolderInBoxAndMail:"

	Map_ClientBoxFolderIdParentId                       = "ClientBoxFolderIdParentId:"
	Map_ClientBoxFolderId                               = "ClientBoxFolderId:"
	Map_ClientMiscThingsToKnowCPExamFileId              = "ClientMiscThingsToKnowCPExamFileId:"
	Map_BuildAuto_BoxDC_RecordReviewFolderId            = Map_BuildAuto + "BoxDC_RecordReviewFolderId:"
	Map_BuildAuto_BoxDC_QuestionnairesFolderId          = Map_BuildAuto + "BoxDC_QuestionnairesFolderId:"
	Map_BuildAuto_BoxDC_PrivateExamsFolderId            = Map_BuildAuto + "BoxDC_PrivateExamsFolderId:"
	Map_BuildAuto_BoxDC_PersonalStatementsFolderId      = Map_BuildAuto + "BoxDC_PersonalStatementsFolderId:"
	Map_BuildAuto_BoxDC_ReleaseOfInformationForm_FileId = Map_BuildAuto + "BoxDC_ReleaseOfInformationForm_FileId:"
	Map_BuildAuto_BoxDC_PatientPaymentForm_FileId       = Map_BuildAuto + "BoxDC_PatientPaymentForm_FileId:"
	Map_BuildAuto_BoxDC_ClaimsAnalysisFolderId          = Map_BuildAuto + "BoxDC_ClaimsAnalysisFolderId:"

	Map_BuildAuto_BoxC_PrivateMedicalRecordsFolderId   = Map_BuildAuto + "BoxC_PrivateMedicalRecordsFolderId:"
	Map_BuildAuto_BoxC_VAMedicalRecordsFolderId        = Map_BuildAuto + "BoxC_VAMedicalRecordsFolderId:"
	Map_BuildAuto_BoxC_ServiceTreatmentRecordsFolderId = Map_BuildAuto + "BoxC_ServiceTreatmentRecordsFolderId:"
	Map_BuildAuto_BoxC_MiscFolderId                    = Map_BuildAuto + "BoxC_MiscFolderId:"
	Map_BuildAuto_BoxC_DD214FolderId                   = Map_BuildAuto + "BoxC_DD214FolderId:"
	Map_BuildAuto_BoxC_DisabilityRatingListFolderId    = Map_BuildAuto + "BoxC_DisabilityRatingListFolderId:"
	Map_BuildAuto_BoxC_RatingDecisionLettersFolderId   = Map_BuildAuto + "BoxC_RatingDecisionLettersFolderId:"
	Map_BuildAuto_BoxC_PersonalStatementsFolderId      = Map_BuildAuto + "BoxC_PersonalStatementsFolderId:"

	Map_EnvelopeDocuments                 = "EnvelopeDocuments:"
	Map_XeroContactId                     = "XeroContactId:"
	Map_XeroInvoiceId                     = "XeroInvoiceId:"
	Map_MaCongratsEmail                   = "MaCongratsEmail:"
	Map_ZohoDealHandleLastModifiedTime    = "ZohoDealHandleLastModifiedTime"    // Zoho deal处理的最后时间
	Map_ZohoDealHandleLastModifiedTime2   = "ZohoDealHandleLastModifiedTime2"   // Zoho deal处理的最后时间
	Map_ZohoContactHandleLastModifiedTime = "ZohoContactHandleLastModifiedTime" // Zoho contact处理的最后时间
	Map_ZohoTaskHandleLastModifiedTime    = "ZohoTaskHandleLastModifiedTime"    // Zoho task处理的最后时间
	Map_ZohoNoteHandleLastModifiedTime    = "ZohoNoteHandleLastModifiedTime"    // Zoho note处理的最后时间
	Map_ActionOnce                        = "ActionOnce:"

	// 合同来源基础数据, 快照，用户生成invoice
	Map_ClientCaseContractBasicData = "ClientCaseContractBasicData:"
	//Map_ClientCaseContractPricingVersion = "ClientCaseContractPricingVersion:"

	MapGoogleDrivePaymentFolderId = "GoogleDrivePaymentFolderId:"
	MapGoogleDrivePsychFolderId   = "GoogleDrivePsychFolderId:"
	MapGoogleDriveGeneralFolderId = "GoogleDriveGeneralFolderId:"

	MapCaseWithoutTaskFlag = "CaseWithoutTaskFlag:"

	MapMeetingSmsNotice = "MeetingSmsNotice:"
	MapITFClientTask    = "ITFClientTask:"

	MapHandleQuestionnairesJotformHistory = "HandleQuestionnairesJotformHistory:"

	MapCustomViewSortBy      = "CustomViewSortBy:"
	MapCustomViewColumns     = "CustomViewColumns:"
	MapCustomViewColumnwidth = "CustomViewColumnwidth:"

	MapMedicalDbqCost                       = "MedicalDbqCost:"
	MapLeadVSChangeLog                      = "LeadVSChangeLog:"
	MapUpcomingContactInformation           = "UpcomingContactInformation:"
	MapCopyPersonalStatementsDoc            = "CopyPersonalStatementsDoc:"
	MapEmailMiniDBQsDrafts                  = "EmailMiniDBQsDrafts:"
	MapYourRecordsReviewProcessHasBegun     = "YourRecordsReviewProcessHasBegun:"
	MapPleaseScheduleYourDoctorAppointments = "PleaseScheduleYourDoctorAppointments:"

	MapPersonalStatementPassword        = "PersonalStatementPassword:"
	MapPersonalStatementSubmitForReview = "PersonalStatementSubmitForReview:"
	//MapPersonalStatementUsePW           = "PersonalStatementUsePW:"

	MapAmInformationIntake = "AmInformationIntake:"
	MapAmContractPending   = "AmContractPending:"

	MapCurrentAttorneyIndex = "CurrentAttorneyIndex"

	MapQuestionnairesFileTime = "QuestionnairesFileTime:"

	MapHasInitStatementsEdit = "HasInitStatementsEdit:"

	MapPSCurrentDocEmailAiResultId = "PSCurrentDocEmailAiResultId:"

	MapClientCaseAmSignedAgreementBoxFolderId = "AmSignedAgreementBoxFolderId:"
	MapClientCaseAmSignedVA2122aBoxFolderId   = "AmSignedVA2122aBoxFolderId:"
	MapClientCaseAmSignedAgreementBoxFileId   = "AmSignedAgreementBoxFileId:"
	MapClientCaseAmSignedVA2122aBoxFileId     = "AmSignedVA2122aBoxFileId:"

	MapMoving2122aFileId = "MapMoving2122aFile:"

	MapUseVBCActiveCasesFolder = "UseVBCActiveCasesFolder:"
	MapSyncInitCase            = "SyncInitCase:"
)

type LeadVSChangeLogVo struct {
	PreviousVSUserGid string
	NewVSUserGid      string
}

func MapKeySyncInitCase(caseId int32) string {
	return fmt.Sprintf("%s%d", MapSyncInitCase, caseId)
}

func MapKeyUseVBCActiveCasesFolder(caseId int32) string {
	return fmt.Sprintf("%s%d", MapUseVBCActiveCasesFolder, caseId)
}

func MapKeyMoving2122aFileId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapMoving2122aFileId, caseId)
}

func MapKeyClientCaseAmSignedAgreementBoxFileId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapClientCaseAmSignedAgreementBoxFileId, caseId)
}

func MapKeyClientCaseAmSignedVA2122aBoxFileId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapClientCaseAmSignedVA2122aBoxFileId, caseId)
}

func MapKeyClientCaseAmSignedAgreementBoxFolderId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapClientCaseAmSignedAgreementBoxFolderId, caseId)
}

func MapKeyClientCaseAmSignedVA2122aBoxFolderId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapClientCaseAmSignedVA2122aBoxFolderId, caseId)
}

func MapKeyPSCurrentDocEmailAiResultId(caseId int32) string {
	return fmt.Sprintf("%s%d", MapPSCurrentDocEmailAiResultId, caseId)
}

func MapKeyHasInitStatementsEdit(caseId int32) string {
	return fmt.Sprintf("%s%d", MapHasInitStatementsEdit, caseId)
}

func MapKeyQuestionnairesFileTime(submissionId string) string {
	return fmt.Sprintf("%s%s", MapQuestionnairesFileTime, submissionId)
}

func MapKeyCurrentAttorneyIndex() string {
	return MapCurrentAttorneyIndex
}

func MapKeyAmContractPending(caseId int32) string {
	return fmt.Sprintf("%s%d", MapAmContractPending, caseId)
}

func MapKeyAmInformationIntake(caseId int32) string {
	return fmt.Sprintf("%s%d", MapAmInformationIntake, caseId)
}

//
//func MapKeyPersonalStatementUsePW(caseId int32) string {
//	return fmt.Sprintf("%s%d", MapPersonalStatementUsePW, caseId)
//}

func MapKeyPersonalStatementSubmitForReview(caseId int32) string {
	return fmt.Sprintf("%s%d", MapPersonalStatementSubmitForReview, caseId)
}

func MapKeyPersonalStatementPassword(caseId int32) string {
	return fmt.Sprintf("%s%d", MapPersonalStatementPassword, caseId)
}

func MapKeyClientCaseContractBasicData(caseId int32) string {
	return fmt.Sprintf("%s%d", Map_ClientCaseContractBasicData, caseId)
}

func MapKeyYourRecordsReviewProcessHasBegun(caseId int32) string {
	return fmt.Sprintf("%s%d", MapYourRecordsReviewProcessHasBegun, caseId)
}

func MapKeyPleaseScheduleYourDoctorAppointments(caseId int32) string {
	return fmt.Sprintf("%s%d", MapPleaseScheduleYourDoctorAppointments, caseId)
}

func MapKeyPersonalStatementsReadyforYourReview(caseId int32) string {
	return fmt.Sprintf("PersonalStatementsReadyforYourReview:%d", caseId)
}

func MapKeyPleaseReviewYourPersonalStatementsinSharedFolder(caseId int32) string {
	return fmt.Sprintf("PleaseReviewYourPersonalStatementsinSharedFolder:%d", caseId)
}

func MapKeyEmailMiniDBQsDrafts(caseId int32) string {
	return fmt.Sprintf("%s%d", MapEmailMiniDBQsDrafts, caseId)
}

func MapKeyCopyPersonalStatementsDoc(caseId int32) string {
	return fmt.Sprintf("%s%d", MapCopyPersonalStatementsDoc, caseId)
}

func MapKeyUpcomingContactInformation(caseId int32) string {
	return fmt.Sprintf("%s%d", MapUpcomingContactInformation, caseId)
}

func MapKeyLeadVSChangeLog(caseId int32) string {
	return fmt.Sprintf("%s%d", MapLeadVSChangeLog, caseId)
}

func MapKeyMedicalDbqCost(caseId int32) string {
	return fmt.Sprintf("%s:%d", MapMedicalDbqCost, caseId)
}

func MapKeyCustomView(userGid string, kind string, tableType string) string {
	if tableType != "" {
		tableType = ":" + tableType
	}
	return fmt.Sprintf("%s%s:%s%s", MapCustomViewSortBy, userGid, kind, tableType)
}

func MapKeyCustomViewColumns(userGid string, kind string, tableType string) string {
	if tableType != "" {
		tableType = ":" + tableType
	}
	return fmt.Sprintf("%s%s:%s%s", MapCustomViewColumns, userGid, kind, tableType)
}

func MapKeyCustomViewColumnwidth(userGid string, kind string, tableType string) string {
	if tableType != "" {
		tableType = ":" + tableType
	}
	return fmt.Sprintf("%s%s:%s%s", MapCustomViewColumnwidth, userGid, kind, tableType)
}

func MapKeyMapITFClientTask(caseId int32, days int, itf string) string {
	return fmt.Sprintf("%s%s:%d:%d", MapITFClientTask, itf, days, caseId)
}

func MapKeyCaseWithoutTaskFlag(caseId int32, stage string) string {
	return fmt.Sprintf("%s%d:%s", MapCaseWithoutTaskFlag, caseId, stage)
}

//	func MapKeyClientCaseContractPricingVersion(clientCaseId int32) string {
//		return fmt.Sprintf("%s%d", Map_ClientCaseContractPricingVersion, clientCaseId)
//	}
func MapKeyGoogleDrivePaymentFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", MapGoogleDrivePaymentFolderId, clientCaseId)
}
func MapKeyGoogleDrivePsychFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", MapGoogleDrivePsychFolderId, clientCaseId)
}
func MapKeyGoogleDriveGeneralFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", MapGoogleDriveGeneralFolderId, clientCaseId)
}

func MapKeyCopyRecordReviewFiles(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "CopyRecordReviewFiles", clientCaseId)
}

func MapKeyClientMiscThingsToKnowCPExamFileId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_ClientMiscThingsToKnowCPExamFileId, clientCaseId)
}

func MapKeyBuildAutoBoxDCRecordReviewFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_RecordReviewFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxDCQuestionnairesFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_QuestionnairesFolderId, clientCaseId)
}

func MapKeyBuildAutoBoxDCReleaseOfInformationFormFileId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_ReleaseOfInformationForm_FileId, clientCaseId)
}
func MapKeyBuildAutoBoxDCPatientPaymentFormFileId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_PatientPaymentForm_FileId, clientCaseId)
}

func MapKeyBuildAutoBoxDCPrivateExamsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_PrivateExamsFolderId, clientCaseId)
}

func MapKeyBuildAutoBoxDCClaimsAnalysisFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_ClaimsAnalysisFolderId, clientCaseId)
}

func MapKeyBuildAutoBoxDCPersonalStatementsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxDC_PersonalStatementsFolderId, clientCaseId)
}

func MapKeyBuildAutoBoxCPrivateMedicalRecordsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_PrivateMedicalRecordsFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCVAMedicalRecordsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_VAMedicalRecordsFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCServiceTreatmentRecordsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_ServiceTreatmentRecordsFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCMiscFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_MiscFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCDD214FolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_DD214FolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCDisabilityRatingListFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_DisabilityRatingListFolderId, clientCaseId)
}
func MapKeyBuildAutoBoxCRatingDecisionLettersFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_RatingDecisionLettersFolderId, clientCaseId)
}

func MapKeyBuildAutoBoxCPersonalStatementsFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_BuildAuto_BoxC_PersonalStatementsFolderId, clientCaseId)
}

func MapKeyClientBoxFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%d", Map_ClientBoxFolderId, clientCaseId)
}

func MapKeyDataCollectionFolderId(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "DataCollectionFolderId", clientCaseId)
}

func MapKeyPersonalStatementsFile(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "PersonalStatementsFile", clientCaseId)
}

func MapKeyClaimsAnalysisFile(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "ClaimsAnalysis", clientCaseId)
}

func MapKeyDocEmailFile(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "DocEmailFile", clientCaseId)
}

func MapKeyCopyDocEmailFile(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "CopyDocEmailFile", clientCaseId)
}

func MapKeyDoCopyReadPriorToYourDoctorVisitFile(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "CopyReadPriorToYourDoctorVisitFile", clientCaseId)
}

func MapKeyMedicalTeamForms(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "MedicalTeamForms", clientCaseId)
}

func MapKeyMedicalTeamFormsReminderEmail(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "MedicalTeamFormsReminderEmail", clientCaseId)
}

func MapKeyCaseRelaBox(uniqueKey string, clientCaseId int32) string {
	return fmt.Sprintf("CaseRelaBox:%s:%d", uniqueKey, clientCaseId)
}

/*
func MapKeyReleaseOfInformation(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "ReleaseOfInformation", clientCaseId)
}
func MapKeyPatientPaymentForm(clientCaseId int32) string {
	return fmt.Sprintf("%s%s:%d", Map_ActionOnce, "PatientPaymentForm", clientCaseId)
}*/

type MapEntity struct {
	Mkey      string
	Mval      string
	CreatedAt int64
	UpdatedAt int64
}

func (MapEntity) TableName() string {
	return "maps"
}

type MapUsecase struct {
	CommonUsecase *CommonUsecase
	DBUsecase[MapEntity]
}

func NewMapUsecase(CommonUsecase *CommonUsecase) *MapUsecase {

	uc := &MapUsecase{
		CommonUsecase: CommonUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *MapUsecase) SetInt(mkey string, mval int) error {
	return c.Set(mkey, InterfaceToString(mval))
}

func (c *MapUsecase) Set(mkey string, mval string) error {
	currentTime := time.Now()
	mval = lib.SqlBindValue(mval)
	sql := fmt.Sprintf("INSERT INTO %s (mkey, mval ,created_at, updated_at) VALUES(\"%s\", %s, %d, %d ) ON DUPLICATE KEY UPDATE updated_at=values(updated_at),mval=values(mval)",
		MapEntity{}.TableName(), mkey, mval, currentTime.Unix(), currentTime.Unix())
	return c.CommonUsecase.DB().Exec(sql).Error
}

func (c *MapUsecase) GetForString(mkey string) (string, error) {
	var entity MapEntity
	err := c.CommonUsecase.DB().Where("mkey=?", mkey).
		Take(&entity).
		Error
	if err == nil {
		return entity.Mval, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return "", err
}

func (c *MapUsecase) GetForInt(mkey string) (int32, error) {
	val, err := c.GetForString(mkey)
	if err != nil {
		return 0, err
	}
	i, _ := strconv.ParseInt(val, 10, 32)
	return int32(i), nil
}
