package tests

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/internal/conf"
	"vbc/lib"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var ConfData *conf.Data
var UT *UnittestApp

// go build -ldflags "-X main.Version=x.y.z"
var (
	id, _   = os.Hostname()
	Name    string
	Version string

	NACOS_HOST string
	NACOS_PORT int
	NACOS_NS   string
	param      string

	sqldb *gorm.DB
)

func init() {

}

type UnittestApp struct {
	Conf          *conf.Data
	CommonUsecase *biz.CommonUsecase
	UserUsecase   *biz.UserUsecase
	ZapLogger     *zap.Logger
	AsanaUsecase  *biz.AsanaUsecase
	//CustomerUsecase      *biz.CustomerUsecase
	FieldUsecase                *biz.FieldUsecase
	FieldOptionUsecase          *biz.FieldOptionUsecase
	DataEntryUsecase            *biz.DataEntryUsecase
	RedisUsecase                *biz.RedisUsecase
	SyncAsanaTaskUsecase        *biz.SyncAsanaTaskUsecase
	TUsecase                    *biz.TUsecase
	HttpTUsecase                *biz.HttpTUsecase
	MapUsecase                  *biz.MapUsecase
	TaskUsecase                 *biz.TaskUsecase
	TaskCreateUsecase           *biz.TaskCreateUsecase
	MailUsecase                 *biz.MailUsecase
	ChangeHistoryUseacse        *biz.ChangeHistoryUseacse
	HttpWebhookUsecase          *biz.HttpWebhookUsecase
	Oauth2ClientUsecase         *biz.Oauth2ClientUsecase
	Oauth2TokenUsecase          *biz.Oauth2TokenUsecase
	EnvelopeStatusChangeUsecase *biz.EnvelopeStatusChangeUsecase
	DocuSignUsecase             *biz.DocuSignUsecase
	BoxUsecase                  *biz.BoxUsecase
	Dag                         *biz.Dag
	AdobeSignUsecase            *biz.AdobeSignUsecase
	ClientAgreementUsecase      *biz.ClientAgreementUsecase
	AdobesignSyncTaskUsecase    *biz.AdobesignSyncTaskUsecase
	FormResponseUsecase         *biz.FormResponseUsecase
	RemindUsecase               *biz.RemindUsecase
	BehaviorUsecase             *biz.BehaviorUsecase
	BoxcontractUsecase          *biz.BoxcontractUsecase
	BoxWebhookLogUsecase        *biz.BoxWebhookLogUsecase
	ClientEnvelopeUsecase       *biz.ClientEnvelopeUsecase
	RollpoingUsecase            *biz.RollpoingUsecase
	RollpoingJobUsecase         *biz.RollpoingJobUsecase
	UniqueCodeGeneratorUsecase  *biz.UniqueCodeGeneratorUsecase
	HttpManualUsecase           *biz.HttpManualUsecase
	TaskFailureLogJobUsecase    *biz.TaskFailureLogJobUsecase
	//GoogleSheetUseacse                 *biz.GoogleSheetUseacse
	WebsiteUsecase       *biz.WebsiteUsecase
	WebhookLogJobUsecase *biz.WebhookLogJobUsecase
	//GoogleSheetSyncTaskUsecase         *biz.GoogleSheetSyncTaskUsecase
	TaskFailureLogUsecase                  *biz.TaskFailureLogUsecase
	ClientUsecase                          *biz.ClientUsecase
	AccessControlWorkUsecase               *biz.AccessControlWorkUsecase
	HttpAccessControl                      *biz.HttpAccessControl
	XeroUsecase                            *biz.XeroUsecase
	XeroInvoiceUsecase                     *biz.XeroInvoiceUsecase
	MaCongratsEmailUsecase                 *biz.MaCongratsEmailUsecase
	ZohoUsecase                            *biz.ZohoUsecase
	ZohobuzUsecase                         *biz.ZohobuzUsecase
	ZohoDealScanJobUsecase                 *biz.ZohoDealScanJobUsecase
	ZohoContactScanJobUsecase              *biz.ZohoContactScanJobUsecase
	AsanaMigrateUsecase                    *biz.AsanaMigrateUsecase
	UsageStatsUsecase                      *biz.UsageStatsUsecase
	ActionOnceUsecase                      *biz.ActionOnceUsecase
	ClientCaseContractBasicDataUsecase     *biz.ClientCaseContractBasicDataUsecase
	FeeUsecase                             *biz.FeeUsecase
	ZohoTaskScanJobUsecase                 *biz.ZohoTaskScanJobUsecase
	ClientTaskUsecase                      *biz.ClientTaskUsecase
	ClientCaseUsecase                      *biz.ClientCaseUsecase
	BoxbuzUsecase                          *biz.BoxbuzUsecase
	DbqsUsecase                            *biz.DbqsUsecase
	DataComboUsecase                       *biz.DataComboUsecase
	MiscUsecase                            *biz.MiscUsecase
	GoogleDriveUsecase                     *biz.GoogleDriveUsecase
	PdfcpuUsecase                          *biz.PdfcpuUsecase
	RecordReviewUsecase                    *biz.RecordReviewUsecase
	RecordReviewJobUsecase                 *biz.RecordReviewJobUsecase
	BoxdcUsecase                           *biz.BoxdcUsecase
	ReminderEventUsecase                   *biz.ReminderEventUsecase
	ReminderEventsJobUsecase               *biz.ReminderEventsJobUsecase
	BoxcUsecase                            *biz.BoxcUsecase
	CounterUsecase                         *biz.CounterUsecase
	CounterbuzUsecase                      *biz.CounterbuzUsecase
	ZohobuzTaskUsecase                     *biz.ZohobuzTaskUsecase
	GoogleDrivebuzUsecase                  *biz.GoogleDrivebuzUsecase
	PricingVersionUsecase                  *biz.PricingVersionUsecase
	AzstorageUsecase                       *biz.AzstorageUsecase
	AzopenaiUsecase                        *biz.AzopenaiUsecase
	BlobbuzUsecase                         *biz.BlobbuzUsecase
	BlobUsecase                            *biz.BlobUsecase
	BlobSliceJobUsecase                    *biz.BlobSliceJobUsecase
	BlobSliceUsecase                       *biz.BlobSliceUsecase
	AzcognitiveUsecase                     *biz.AzcognitiveUsecase
	HaiUsecase                             *biz.HaiUsecase
	HaReportTaskJobUsecase                 *biz.HaReportTaskJobUsecase
	HaReportTaskUsecase                    *biz.HaReportTaskUsecase
	PdfUsecase                             *biz.PdfUsecase
	HaReportPdfUsecase                     *biz.HaReportPdfUsecase
	HaReportPageUsecase                    *biz.HaReportPageUsecase
	HaReportTasksBuzUsecase                *biz.HaReportTasksBuzUsecase
	BlobJobUsecase                         *biz.BlobJobUsecase
	HaReportPageJobUsecase                 *biz.HaReportPageJobUsecase
	ContractReminderUsecase                *biz.ContractReminderUsecase
	DialpadUsecase                         *biz.DialpadUsecase
	DialpadbuzUsecase                      *biz.DialpadbuzUsecase
	CaseWithoutTaskUsecase                 *biz.CaseWithoutTaskUsecase
	CronUsecase                            *biz.CronUsecase
	CronTriggerUsecase                     *biz.CronTriggerUsecase
	SendsmsClientTasksConditionUsecase     *biz.SendsmsClientTasksConditionUsecase
	StageTransUsecase                      *biz.StageTransUsecase
	HttpOauth2Usecase                      *biz.HttpOauth2Usecase
	AppTokenUsecase                        *biz.AppTokenUsecase
	ZoomUsecase                            *biz.ZoomUsecase
	ZoomTokenUsecase                       *biz.ZoomTokenUsecase
	ZoombuzUsecase                         *biz.ZoombuzUsecase
	ZoomRecordingFileUsecase               *biz.ZoomRecordingFileUsecase
	ZoomRecordingFileJobUsecase            *biz.ZoomRecordingFileJobUsecase
	CronTriggerCreateUsecase               *biz.CronTriggerCreateUsecase
	ZoomMeetingUsecase                     *biz.ZoomMeetingUsecase
	Awsclaude3Usecase                      *biz.Awsclaude3Usecase
	FabUsecase                             *biz.FabUsecase
	HttpSettingsUsecase                    *biz.HttpSettingsUsecase
	ClientTaskBuzUsecase                   *biz.ClientTaskBuzUsecase
	ItfexpirationUsecase                   *biz.ItfexpirationUsecase
	LogUsecase                             *biz.LogUsecase
	ZoomMeetingSmsNoticeJobUsecase         *biz.ZoomMeetingSmsNoticeJobUsecase
	TimezoneUsecase                        *biz.TimezoneUsecase
	BUsaStateUsecase                       *biz.BUsaStateUsecase
	LogInfoUsecase                         *biz.LogInfoUsecase
	ManUsecase                             *biz.ManUsecase
	JotformUsecase                         *biz.JotformUsecase
	JotformbuzUsecase                      *biz.JotformbuzUsecase
	ClientTaskJobUsecase                   *biz.ClientTaskJobUsecase
	LongMapUsecase                         *biz.LongMapUsecase
	JotformSubmissionUsecase               *biz.JotformSubmissionUsecase
	QuestionnairesUsecase                  *biz.QuestionnairesUsecase
	AiStatementUsecase                     *biz.AiStatementUsecase
	AiTaskUsecase                          *biz.AiTaskUsecase
	AiTaskJobUsecase                       *biz.AiTaskJobUsecase
	MetadataUsecase                        *biz.MetadataUsecase
	ConditionUsecase                       *biz.ConditionUsecase
	ConditionbuzUsecase                    *biz.ConditionbuzUsecase
	GopdfUsecase                           *biz.GopdfUsecase
	MetadataHttpUsecase                    *biz.MetadataHttpUsecase
	RecordbuzUsecase                       *biz.RecordbuzUsecase
	RecordHttpUsecase                      *biz.RecordHttpUsecase
	SettingSectionFieldUsecase             *biz.SettingSectionFieldUsecase
	KindUsecase                            *biz.KindUsecase
	TimelineUsecase                        *biz.TimelineUsecase
	NotesUsecase                           *biz.NotesUsecase
	TimelinesbuzUsecase                    *biz.TimelinesbuzUsecase
	RatingPaymentUsecase                   *biz.RatingPaymentUsecase
	MailFeeContentUsecase                  *biz.MailFeeContentUsecase
	BUsecase                               *biz.BUsecase
	FieldPermissionUsecase                 *biz.FieldPermissionUsecase
	FieldValidatorUsecase                  *biz.FieldValidatorUsecase
	ManualThingstoknowUsecase              *biz.ManualThingstoknowUsecase
	TimezonesUsecase                       *biz.TimezonesUsecase
	WordUsecase                            *biz.WordUsecase
	ZohoDealScan2JobUsecase                *biz.ZohoDealScan2JobUsecase
	SettingHttpUsecase                     *biz.SettingHttpUsecase
	ClientTaskHandleWhatGidJobUsecase      *biz.ClientTaskHandleWhatGidJobUsecase
	ClientTaskHandleWhoGidJobUsecase       *biz.ClientTaskHandleWhoGidJobUsecase
	QueueUsecase                           *biz.QueueUsecase
	RecordbuzSearchUsecase                 *biz.RecordbuzSearchUsecase
	TaskHttpUsecase                        *biz.TaskHttpUsecase
	ReissueTriggerStrRequestPendingUsecase *biz.ReissueTriggerStrRequestPendingUsecase
	AutomaticTaskCreationUsecase           *biz.AutomaticTaskCreationUsecase
	AutomaticUpdateDueDateUsecase          *biz.AutomaticUpdateDueDateUsecase
	NotificationHttpUsecase                *biz.NotificationHttpUsecase
	ZohoCollaboratorUsecase                *biz.ZohoCollaboratorUsecase
	UnsubscribesbuzUsecase                 *biz.UnsubscribesbuzUsecase
	UnsubscribesUsecase                    *biz.UnsubscribesUsecase
	DialpadbuzInternalUsecase              *biz.DialpadbuzInternalUsecase
	UnsubscribesHttpUsecase                *biz.UnsubscribesHttpUsecase
	ConditionSourcebuzUsecase              *biz.ConditionSourcebuzUsecase
	AiPromptUsecase                        *biz.AiPromptUsecase
	ConditionRelaAiUsecase                 *biz.ConditionRelaAiUsecase
	ConditionLogAiUsecase                  *biz.ConditionLogAiUsecase
	ConditionHttpUsecase                   *biz.ConditionHttpUsecase
	ConditionCategoryUsecase               *biz.ConditionCategoryUsecase
	FilterbuzUsecase                       *biz.FilterbuzUsecase
	RecordLogbuzUsecase                    *biz.RecordLogbuzUsecase
	ZohoNoteScanJobUsecase                 *biz.ZohoNoteScanJobUsecase
	InvokeLogUsecase                       *biz.InvokeLogUsecase
	InvokeLogJobUsecase                    *biz.InvokeLogJobUsecase
	ClientCaseSyncbuzUsecase               *biz.ClientCaseSyncbuzUsecase
	ChangeHistoryNodelayJobUseacse         *biz.ChangeHistoryNodelayJobUseacse
	DueDateUsecase                         *biz.DueDateUsecase
	VbcAIUsecase                           *biz.VbcAIUsecase
	AiUsecase                              *biz.AiUsecase
	AiHttpUsecase                          *biz.AiHttpUsecase
	DocEmailUsecase                        *biz.DocEmailUsecase
	MedicalDbqCostUsecase                  *biz.MedicalDbqCostUsecase
	ClientCasebuzUsecase                   *biz.ClientCasebuzUsecase
	QuestionnairesbuzUsecase               *biz.QuestionnairesbuzUsecase
	CommonHttpUsecase                      *biz.CommonHttpUsecase
	OptionUsecase                          *biz.OptionUsecase
	StatementUsecase                       *biz.StatementUsecase
	UserHttpUsecase                        *biz.UserHttpUsecase
	ClientNameChangeJobUsecase             *biz.ClientNameChangeJobUsecase
	ResourceUsecase                        *biz.ResourceUsecase
	AdobepdfUsecase                        *biz.AdobepdfUsecase
	RoleUsecase                            *biz.RoleUsecase
	LeadVSChangeUsecase                    *biz.LeadVSChangeUsecase
	ChangeHisUsecase                       *biz.ChangeHisUsecase
	AiTaskbuzUsecase                       *biz.AiTaskbuzUsecase
	CacheLogUsecase                        *biz.CacheLogUsecase
	ZoomUploadBoxUsecase                   *biz.ZoomUploadBoxUsecase
	WordbuzUsecase                         *biz.WordbuzUsecase
	PsbuzUsecase                           *biz.PsbuzUsecase
	AiResultUsecase                        *biz.AiResultUsecase
	HttpBlobUsecase                        *biz.HttpBlobUsecase
	WebpUsecase                            *biz.WebpUsecase
	FilebuzUsecase                         *biz.FilebuzUsecase
	ContractbuzUsecase                     *biz.ContractbuzUsecase
	VbcDataVerifyUsecase                   *biz.VbcDataVerifyUsecase
	ManualUsecase                          *biz.ManualUsecase
	LeadConversionSummaryBuzUsecase        *biz.LeadConversionSummaryBuzUsecase
	ReferrerLogUsecase                     *biz.ReferrerLogUsecase
	StatemtUsecase                         *biz.StatemtUsecase
	LeadsUsecase                           *biz.LeadsUsecase
	ClientReviewUsecase                    *biz.ClientReviewUsecase
	ClientReviewBuzUsecase                 *biz.ClientReviewBuzUsecase
	StatementConditionBuzUsecase           *biz.StatementConditionBuzUsecase
	StatementConditionUsecase              *biz.StatementConditionUsecase
	StatementCommentBuzUsecase             *biz.StatementCommentBuzUsecase
	BoxUserBuzUsecase                      *biz.BoxUserBuzUsecase
	BoxCollaborationBuzUsecase             *biz.BoxCollaborationBuzUsecase
	AssistantUsecase                       *biz.AssistantUsecase
	AiAssistantJobBuzUsecase               *biz.AiAssistantJobBuzUsecase
	AttorneybuzUsecase                     *biz.AttorneybuzUsecase
	AttorneyUsecase                        *biz.AttorneyUsecase
	NotesbuzUsecase                        *biz.NotesbuzUsecase
	SendVa2122aUsecase                     *biz.SendVa2122aUsecase
	GlobalEventBusBuzUsecase               *biz.GlobalEventBusBuzUsecase
	CmdUsecase                             *biz.CmdUsecase
	PersonalWebformUsecase                 *biz.PersonalWebformUsecase
	MonitoredEmailsUsecase                 *biz.MonitoredEmailsUsecase
	MonitoredEmailsJobUsecase              *biz.MonitoredEmailsJobUsecase
	MonitoredEmailsTasksUsecase            *biz.MonitoredEmailsTasksUsecase
	ExportUsecase                          *biz.ExportUsecase
	CollaboratorbuzUsecase                 *biz.CollaboratorbuzUsecase
	UserbuzUsecase                         *biz.UserbuzUsecase
}

func newApp(logger log.Logger,
	Conf *conf.Data,
	CommonUsecase *biz.CommonUsecase,
	UserUsecase *biz.UserUsecase,
	ZapLogger *zap.Logger,
	AsanaUsecase *biz.AsanaUsecase,
	//CustomerUsecase *biz.CustomerUsecase,
	FieldUsecase *biz.FieldUsecase,
	FieldOptionUsecase *biz.FieldOptionUsecase,
	DataEntryUsecase *biz.DataEntryUsecase,
	RedisUsecase *biz.RedisUsecase,
	SyncAsanaTaskUsecase *biz.SyncAsanaTaskUsecase,
	TUsecase *biz.TUsecase,
	HttpTUsecase *biz.HttpTUsecase,
	MapUsecase *biz.MapUsecase,
	TaskUsecase *biz.TaskUsecase,
	TaskCreateUsecase *biz.TaskCreateUsecase,
	MailUsecase *biz.MailUsecase,
	ChangeHistoryUseacse *biz.ChangeHistoryUseacse,
	HttpWebhookUsecase *biz.HttpWebhookUsecase,
	Oauth2ClientUsecase *biz.Oauth2ClientUsecase,
	Oauth2TokenUsecase *biz.Oauth2TokenUsecase,
	EnvelopeStatusChangeUsecase *biz.EnvelopeStatusChangeUsecase,
	DocuSignUsecase *biz.DocuSignUsecase,
	BoxUsecase *biz.BoxUsecase,
	Dag *biz.Dag,
	AdobeSignUsecase *biz.AdobeSignUsecase,
	ClientAgreementUsecase *biz.ClientAgreementUsecase,
	AdobesignSyncTaskUsecase *biz.AdobesignSyncTaskUsecase,
	FormResponseUsecase *biz.FormResponseUsecase,
	RemindUsecase *biz.RemindUsecase,
	BehaviorUsecase *biz.BehaviorUsecase,
	BoxcontractUsecase *biz.BoxcontractUsecase,
	BoxWebhookLogUsecase *biz.BoxWebhookLogUsecase,
	ClientEnvelopeUsecase *biz.ClientEnvelopeUsecase,
	RollpoingUsecase *biz.RollpoingUsecase,
	RollpoingJobUsecase *biz.RollpoingJobUsecase,
	UniqueCodeGeneratorUsecase *biz.UniqueCodeGeneratorUsecase,
	HttpManualUsecase *biz.HttpManualUsecase,
	TaskFailureLogJobUsecase *biz.TaskFailureLogJobUsecase,
	//GoogleSheetUseacse *biz.GoogleSheetUseacse,
	WebsiteUsecase *biz.WebsiteUsecase,
	WebhookLogJobUsecase *biz.WebhookLogJobUsecase,
	//GoogleSheetSyncTaskUsecase *biz.GoogleSheetSyncTaskUsecase,
	TaskFailureLogUsecase *biz.TaskFailureLogUsecase,
	ClientUsecase *biz.ClientUsecase,
	AccessControlWorkUsecase *biz.AccessControlWorkUsecase,
	HttpAccessControl *biz.HttpAccessControl,
	XeroUsecase *biz.XeroUsecase,
	XeroInvoiceUsecase *biz.XeroInvoiceUsecase,
	MaCongratsEmailUsecase *biz.MaCongratsEmailUsecase,
	ZohoUsecase *biz.ZohoUsecase,
	ZohobuzUsecase *biz.ZohobuzUsecase,
	ZohoDealScanJobUsecase *biz.ZohoDealScanJobUsecase,
	ZohoContactScanJobUsecase *biz.ZohoContactScanJobUsecase,
	AsanaMigrateUsecase *biz.AsanaMigrateUsecase,
	UsageStatsUsecase *biz.UsageStatsUsecase,
	ActionOnceUsecase *biz.ActionOnceUsecase,
	ClientCaseContractBasicDataUsecase *biz.ClientCaseContractBasicDataUsecase,
	FeeUsecase *biz.FeeUsecase,
	ZohoTaskScanJobUsecase *biz.ZohoTaskScanJobUsecase,
	ClientTaskUsecase *biz.ClientTaskUsecase,
	ClientCaseUsecase *biz.ClientCaseUsecase,
	BoxbuzUsecase *biz.BoxbuzUsecase,
	DbqsUsecase *biz.DbqsUsecase,
	DataComboUsecase *biz.DataComboUsecase,
	MiscUsecase *biz.MiscUsecase,
	GoogleDriveUsecase *biz.GoogleDriveUsecase,
	PdfcpuUsecase *biz.PdfcpuUsecase,
	RecordReviewUsecase *biz.RecordReviewUsecase,
	RecordReviewJobUsecase *biz.RecordReviewJobUsecase,
	BoxdcUsecase *biz.BoxdcUsecase,
	ReminderEventUsecase *biz.ReminderEventUsecase,
	ReminderEventsJobUsecase *biz.ReminderEventsJobUsecase,
	BoxcUsecase *biz.BoxcUsecase,
	CounterUsecase *biz.CounterUsecase,
	CounterbuzUsecase *biz.CounterbuzUsecase,
	ZohobuzTaskUsecase *biz.ZohobuzTaskUsecase,
	GoogleDrivebuzUsecase *biz.GoogleDrivebuzUsecase,
	PricingVersionUsecase *biz.PricingVersionUsecase,
	AzstorageUsecase *biz.AzstorageUsecase,
	AzopenaiUsecase *biz.AzopenaiUsecase,
	BlobbuzUsecase *biz.BlobbuzUsecase,
	BlobUsecase *biz.BlobUsecase,
	BlobSliceJobUsecase *biz.BlobSliceJobUsecase,
	BlobSliceUsecase *biz.BlobSliceUsecase,
	AzcognitiveUsecase *biz.AzcognitiveUsecase,
	HaiUsecase *biz.HaiUsecase,
	HaReportTaskJobUsecase *biz.HaReportTaskJobUsecase,
	HaReportTaskUsecase *biz.HaReportTaskUsecase,
	PdfUsecase *biz.PdfUsecase,
	HaReportPdfUsecase *biz.HaReportPdfUsecase,
	HaReportPageUsecase *biz.HaReportPageUsecase,
	HaReportTasksBuzUsecase *biz.HaReportTasksBuzUsecase,
	BlobJobUsecase *biz.BlobJobUsecase,
	HaReportPageJobUsecase *biz.HaReportPageJobUsecase,
	ContractReminderUsecase *biz.ContractReminderUsecase,
	DialpadUsecase *biz.DialpadUsecase,
	DialpadbuzUsecase *biz.DialpadbuzUsecase,
	CaseWithoutTaskUsecase *biz.CaseWithoutTaskUsecase,
	CronUsecase *biz.CronUsecase,
	CronTriggerUsecase *biz.CronTriggerUsecase,
	SendsmsClientTasksConditionUsecase *biz.SendsmsClientTasksConditionUsecase,
	StageTransUsecase *biz.StageTransUsecase,
	HttpOauth2Usecase *biz.HttpOauth2Usecase,
	AppTokenUsecase *biz.AppTokenUsecase,
	ZoomUsecase *biz.ZoomUsecase,
	ZoomTokenUsecase *biz.ZoomTokenUsecase,
	ZoombuzUsecase *biz.ZoombuzUsecase,
	ZoomRecordingFileUsecase *biz.ZoomRecordingFileUsecase,
	ZoomRecordingFileJobUsecase *biz.ZoomRecordingFileJobUsecase,
	CronTriggerCreateUsecase *biz.CronTriggerCreateUsecase,
	ZoomMeetingUsecase *biz.ZoomMeetingUsecase,
	Awsclaude3Usecase *biz.Awsclaude3Usecase,
	FabUsecase *biz.FabUsecase,
	HttpSettingsUsecase *biz.HttpSettingsUsecase,
	ClientTaskBuzUsecase *biz.ClientTaskBuzUsecase,
	ItfexpirationUsecase *biz.ItfexpirationUsecase,
	LogUsecase *biz.LogUsecase,
	ZoomMeetingSmsNoticeJobUsecase *biz.ZoomMeetingSmsNoticeJobUsecase,
	TimezoneUsecase *biz.TimezoneUsecase,
	BUsaStateUsecase *biz.BUsaStateUsecase,
	LogInfoUsecase *biz.LogInfoUsecase,
	ManUsecase *biz.ManUsecase,
	JotformUsecase *biz.JotformUsecase,
	JotformbuzUsecase *biz.JotformbuzUsecase,
	ClientTaskJobUsecase *biz.ClientTaskJobUsecase,
	LongMapUsecase *biz.LongMapUsecase,
	JotformSubmissionUsecase *biz.JotformSubmissionUsecase,
	QuestionnairesUsecase *biz.QuestionnairesUsecase,
	AiStatementUsecase *biz.AiStatementUsecase,
	AiTaskUsecase *biz.AiTaskUsecase,
	AiTaskJobUsecase *biz.AiTaskJobUsecase,
	MetadataUsecase *biz.MetadataUsecase,
	ConditionUsecase *biz.ConditionUsecase,
	ConditionbuzUsecase *biz.ConditionbuzUsecase,
	GopdfUsecase *biz.GopdfUsecase,
	MetadataHttpUsecase *biz.MetadataHttpUsecase,
	RecordbuzUsecase *biz.RecordbuzUsecase,
	RecordHttpUsecase *biz.RecordHttpUsecase,
	SettingSectionFieldUsecase *biz.SettingSectionFieldUsecase,
	KindUsecase *biz.KindUsecase,
	TimelineUsecase *biz.TimelineUsecase,
	NotesUsecase *biz.NotesUsecase,
	TimelinesbuzUsecase *biz.TimelinesbuzUsecase,
	RatingPaymentUsecase *biz.RatingPaymentUsecase,
	MailFeeContentUsecase *biz.MailFeeContentUsecase,
	BUsecase *biz.BUsecase,
	FieldPermissionUsecase *biz.FieldPermissionUsecase,
	FieldValidatorUsecase *biz.FieldValidatorUsecase,
	ManualThingstoknowUsecase *biz.ManualThingstoknowUsecase,
	TimezonesUsecase *biz.TimezonesUsecase,
	WordUsecase *biz.WordUsecase,
	ZohoDealScan2JobUsecase *biz.ZohoDealScan2JobUsecase,
	SettingHttpUsecase *biz.SettingHttpUsecase,
	ClientTaskHandleWhatGidJobUsecase *biz.ClientTaskHandleWhatGidJobUsecase,
	ClientTaskHandleWhoGidJobUsecase *biz.ClientTaskHandleWhoGidJobUsecase,
	QueueUsecase *biz.QueueUsecase,
	RecordbuzSearchUsecase *biz.RecordbuzSearchUsecase,
	TaskHttpUsecase *biz.TaskHttpUsecase,
	ReissueTriggerStrRequestPendingUsecase *biz.ReissueTriggerStrRequestPendingUsecase,
	AutomaticTaskCreationUsecase *biz.AutomaticTaskCreationUsecase,
	AutomaticUpdateDueDateUsecase *biz.AutomaticUpdateDueDateUsecase,
	NotificationHttpUsecase *biz.NotificationHttpUsecase,
	ZohoCollaboratorUsecase *biz.ZohoCollaboratorUsecase,
	UnsubscribesbuzUsecase *biz.UnsubscribesbuzUsecase,
	UnsubscribesUsecase *biz.UnsubscribesUsecase,
	DialpadbuzInternalUsecase *biz.DialpadbuzInternalUsecase,
	GlobalInjectUsecase *biz.GlobalInjectUsecase,
	UnsubscribesHttpUsecase *biz.UnsubscribesHttpUsecase,
	ConditionSourcebuzUsecase *biz.ConditionSourcebuzUsecase,
	AiPromptUsecase *biz.AiPromptUsecase,
	ConditionRelaAiUsecase *biz.ConditionRelaAiUsecase,
	ConditionLogAiUsecase *biz.ConditionLogAiUsecase,
	ConditionHttpUsecase *biz.ConditionHttpUsecase,
	ConditionCategoryUsecase *biz.ConditionCategoryUsecase,
	FilterbuzUsecase *biz.FilterbuzUsecase,
	RecordLogbuzUsecase *biz.RecordLogbuzUsecase,
	ZohoNoteScanJobUsecase *biz.ZohoNoteScanJobUsecase,
	InvokeLogUsecase *biz.InvokeLogUsecase,
	InvokeLogJobUsecase *biz.InvokeLogJobUsecase,
	ClientCaseSyncbuzUsecase *biz.ClientCaseSyncbuzUsecase,
	ChangeHistoryNodelayJobUseacse *biz.ChangeHistoryNodelayJobUseacse,
	DueDateUsecase *biz.DueDateUsecase,
	VbcAIUsecase *biz.VbcAIUsecase,
	AiUsecase *biz.AiUsecase,
	AiHttpUsecase *biz.AiHttpUsecase,
	DocEmailUsecase *biz.DocEmailUsecase,
	MedicalDbqCostUsecase *biz.MedicalDbqCostUsecase,
	ClientCasebuzUsecase *biz.ClientCasebuzUsecase,
	QuestionnairesbuzUsecase *biz.QuestionnairesbuzUsecase,
	CommonHttpUsecase *biz.CommonHttpUsecase,
	OptionUsecase *biz.OptionUsecase,
	StatementUsecase *biz.StatementUsecase,
	UserHttpUsecase *biz.UserHttpUsecase,
	ClientNameChangeJobUsecase *biz.ClientNameChangeJobUsecase,
	ResourceUsecase *biz.ResourceUsecase,
	AdobepdfUsecase *biz.AdobepdfUsecase,
	RoleUsecase *biz.RoleUsecase,
	LeadVSChangeUsecase *biz.LeadVSChangeUsecase,
	ChangeHisUsecase *biz.ChangeHisUsecase,
	AiTaskbuzUsecase *biz.AiTaskbuzUsecase,
	CacheLogUsecase *biz.CacheLogUsecase,
	ZoomUploadBoxUsecase *biz.ZoomUploadBoxUsecase,
	WordbuzUsecase *biz.WordbuzUsecase,
	PsbuzUsecase *biz.PsbuzUsecase,
	AiResultUsecase *biz.AiResultUsecase,
	HttpBlobUsecase *biz.HttpBlobUsecase,
	WebpUsecase *biz.WebpUsecase,
	FilebuzUsecase *biz.FilebuzUsecase,
	ContractbuzUsecase *biz.ContractbuzUsecase,
	VbcDataVerifyUsecase *biz.VbcDataVerifyUsecase,
	ManualUsecase *biz.ManualUsecase,
	LeadConversionSummaryBuzUsecase *biz.LeadConversionSummaryBuzUsecase,
	ReferrerLogUsecase *biz.ReferrerLogUsecase,
	StatemtUsecase *biz.StatemtUsecase,
	LeadsUsecase *biz.LeadsUsecase,
	ClientReviewUsecase *biz.ClientReviewUsecase,
	ClientReviewBuzUsecase *biz.ClientReviewBuzUsecase,
	StatementConditionBuzUsecase *biz.StatementConditionBuzUsecase,
	StatementConditionUsecase *biz.StatementConditionUsecase,
	StatementCommentBuzUsecase *biz.StatementCommentBuzUsecase,
	BoxUserBuzUsecase *biz.BoxUserBuzUsecase,
	BoxCollaborationBuzUsecase *biz.BoxCollaborationBuzUsecase,
	AssistantUsecase *biz.AssistantUsecase,
	AiAssistantJobBuzUsecase *biz.AiAssistantJobBuzUsecase,
	AttorneybuzUsecase *biz.AttorneybuzUsecase,
	AttorneyUsecase *biz.AttorneyUsecase,
	NotesbuzUsecase *biz.NotesbuzUsecase,
	SendVa2122aUsecase *biz.SendVa2122aUsecase,
	GlobalEventBusBuzUsecase *biz.GlobalEventBusBuzUsecase,
	CmdUsecase *biz.CmdUsecase,
	PersonalWebformUsecase *biz.PersonalWebformUsecase,
	MonitoredEmailsUsecase *biz.MonitoredEmailsUsecase,
	MonitoredEmailsJobUsecase *biz.MonitoredEmailsJobUsecase,
	MonitoredEmailsTasksUsecase *biz.MonitoredEmailsTasksUsecase,
	ExportUsecase *biz.ExportUsecase,
	CollaboratorbuzUsecase *biz.CollaboratorbuzUsecase,
	UserbuzUsecase *biz.UserbuzUsecase,
) *UnittestApp {

	return &UnittestApp{
		Conf:          Conf,
		CommonUsecase: CommonUsecase,
		ZapLogger:     ZapLogger,
		UserUsecase:   UserUsecase,
		AsanaUsecase:  AsanaUsecase,
		//CustomerUsecase:      CustomerUsecase,
		FieldUsecase:                FieldUsecase,
		FieldOptionUsecase:          FieldOptionUsecase,
		DataEntryUsecase:            DataEntryUsecase,
		RedisUsecase:                RedisUsecase,
		SyncAsanaTaskUsecase:        SyncAsanaTaskUsecase,
		TUsecase:                    TUsecase,
		HttpTUsecase:                HttpTUsecase,
		MapUsecase:                  MapUsecase,
		TaskUsecase:                 TaskUsecase,
		TaskCreateUsecase:           TaskCreateUsecase,
		MailUsecase:                 MailUsecase,
		ChangeHistoryUseacse:        ChangeHistoryUseacse,
		HttpWebhookUsecase:          HttpWebhookUsecase,
		Oauth2ClientUsecase:         Oauth2ClientUsecase,
		Oauth2TokenUsecase:          Oauth2TokenUsecase,
		EnvelopeStatusChangeUsecase: EnvelopeStatusChangeUsecase,
		DocuSignUsecase:             DocuSignUsecase,
		BoxUsecase:                  BoxUsecase,
		Dag:                         Dag,
		AdobeSignUsecase:            AdobeSignUsecase,
		ClientAgreementUsecase:      ClientAgreementUsecase,
		AdobesignSyncTaskUsecase:    AdobesignSyncTaskUsecase,
		FormResponseUsecase:         FormResponseUsecase,
		RemindUsecase:               RemindUsecase,
		BehaviorUsecase:             BehaviorUsecase,
		BoxcontractUsecase:          BoxcontractUsecase,
		BoxWebhookLogUsecase:        BoxWebhookLogUsecase,
		ClientEnvelopeUsecase:       ClientEnvelopeUsecase,
		RollpoingUsecase:            RollpoingUsecase,
		RollpoingJobUsecase:         RollpoingJobUsecase,
		UniqueCodeGeneratorUsecase:  UniqueCodeGeneratorUsecase,
		HttpManualUsecase:           HttpManualUsecase,
		TaskFailureLogJobUsecase:    TaskFailureLogJobUsecase,
		//GoogleSheetUseacse:                 GoogleSheetUseacse,
		WebsiteUsecase:       WebsiteUsecase,
		WebhookLogJobUsecase: WebhookLogJobUsecase,
		//GoogleSheetSyncTaskUsecase:         GoogleSheetSyncTaskUsecase,
		TaskFailureLogUsecase:                  TaskFailureLogUsecase,
		ClientUsecase:                          ClientUsecase,
		AccessControlWorkUsecase:               AccessControlWorkUsecase,
		HttpAccessControl:                      HttpAccessControl,
		XeroUsecase:                            XeroUsecase,
		XeroInvoiceUsecase:                     XeroInvoiceUsecase,
		MaCongratsEmailUsecase:                 MaCongratsEmailUsecase,
		ZohoUsecase:                            ZohoUsecase,
		ZohobuzUsecase:                         ZohobuzUsecase,
		ZohoDealScanJobUsecase:                 ZohoDealScanJobUsecase,
		ZohoContactScanJobUsecase:              ZohoContactScanJobUsecase,
		AsanaMigrateUsecase:                    AsanaMigrateUsecase,
		UsageStatsUsecase:                      UsageStatsUsecase,
		ActionOnceUsecase:                      ActionOnceUsecase,
		ClientCaseContractBasicDataUsecase:     ClientCaseContractBasicDataUsecase,
		FeeUsecase:                             FeeUsecase,
		ZohoTaskScanJobUsecase:                 ZohoTaskScanJobUsecase,
		ClientTaskUsecase:                      ClientTaskUsecase,
		ClientCaseUsecase:                      ClientCaseUsecase,
		BoxbuzUsecase:                          BoxbuzUsecase,
		DbqsUsecase:                            DbqsUsecase,
		DataComboUsecase:                       DataComboUsecase,
		MiscUsecase:                            MiscUsecase,
		GoogleDriveUsecase:                     GoogleDriveUsecase,
		PdfcpuUsecase:                          PdfcpuUsecase,
		RecordReviewUsecase:                    RecordReviewUsecase,
		RecordReviewJobUsecase:                 RecordReviewJobUsecase,
		BoxdcUsecase:                           BoxdcUsecase,
		ReminderEventUsecase:                   ReminderEventUsecase,
		ReminderEventsJobUsecase:               ReminderEventsJobUsecase,
		BoxcUsecase:                            BoxcUsecase,
		CounterUsecase:                         CounterUsecase,
		CounterbuzUsecase:                      CounterbuzUsecase,
		ZohobuzTaskUsecase:                     ZohobuzTaskUsecase,
		GoogleDrivebuzUsecase:                  GoogleDrivebuzUsecase,
		PricingVersionUsecase:                  PricingVersionUsecase,
		AzstorageUsecase:                       AzstorageUsecase,
		AzopenaiUsecase:                        AzopenaiUsecase,
		BlobbuzUsecase:                         BlobbuzUsecase,
		BlobUsecase:                            BlobUsecase,
		BlobSliceJobUsecase:                    BlobSliceJobUsecase,
		BlobSliceUsecase:                       BlobSliceUsecase,
		AzcognitiveUsecase:                     AzcognitiveUsecase,
		HaiUsecase:                             HaiUsecase,
		HaReportTaskJobUsecase:                 HaReportTaskJobUsecase,
		HaReportTaskUsecase:                    HaReportTaskUsecase,
		PdfUsecase:                             PdfUsecase,
		HaReportPdfUsecase:                     HaReportPdfUsecase,
		HaReportPageUsecase:                    HaReportPageUsecase,
		HaReportTasksBuzUsecase:                HaReportTasksBuzUsecase,
		BlobJobUsecase:                         BlobJobUsecase,
		HaReportPageJobUsecase:                 HaReportPageJobUsecase,
		ContractReminderUsecase:                ContractReminderUsecase,
		DialpadUsecase:                         DialpadUsecase,
		DialpadbuzUsecase:                      DialpadbuzUsecase,
		CaseWithoutTaskUsecase:                 CaseWithoutTaskUsecase,
		CronUsecase:                            CronUsecase,
		CronTriggerUsecase:                     CronTriggerUsecase,
		SendsmsClientTasksConditionUsecase:     SendsmsClientTasksConditionUsecase,
		StageTransUsecase:                      StageTransUsecase,
		HttpOauth2Usecase:                      HttpOauth2Usecase,
		AppTokenUsecase:                        AppTokenUsecase,
		ZoomUsecase:                            ZoomUsecase,
		ZoomTokenUsecase:                       ZoomTokenUsecase,
		ZoombuzUsecase:                         ZoombuzUsecase,
		ZoomRecordingFileUsecase:               ZoomRecordingFileUsecase,
		ZoomRecordingFileJobUsecase:            ZoomRecordingFileJobUsecase,
		CronTriggerCreateUsecase:               CronTriggerCreateUsecase,
		ZoomMeetingUsecase:                     ZoomMeetingUsecase,
		Awsclaude3Usecase:                      Awsclaude3Usecase,
		FabUsecase:                             FabUsecase,
		HttpSettingsUsecase:                    HttpSettingsUsecase,
		ClientTaskBuzUsecase:                   ClientTaskBuzUsecase,
		ItfexpirationUsecase:                   ItfexpirationUsecase,
		LogUsecase:                             LogUsecase,
		ZoomMeetingSmsNoticeJobUsecase:         ZoomMeetingSmsNoticeJobUsecase,
		TimezoneUsecase:                        TimezoneUsecase,
		BUsaStateUsecase:                       BUsaStateUsecase,
		LogInfoUsecase:                         LogInfoUsecase,
		ManUsecase:                             ManUsecase,
		JotformUsecase:                         JotformUsecase,
		JotformbuzUsecase:                      JotformbuzUsecase,
		ClientTaskJobUsecase:                   ClientTaskJobUsecase,
		LongMapUsecase:                         LongMapUsecase,
		JotformSubmissionUsecase:               JotformSubmissionUsecase,
		QuestionnairesUsecase:                  QuestionnairesUsecase,
		AiStatementUsecase:                     AiStatementUsecase,
		AiTaskUsecase:                          AiTaskUsecase,
		AiTaskJobUsecase:                       AiTaskJobUsecase,
		MetadataUsecase:                        MetadataUsecase,
		ConditionUsecase:                       ConditionUsecase,
		ConditionbuzUsecase:                    ConditionbuzUsecase,
		GopdfUsecase:                           GopdfUsecase,
		MetadataHttpUsecase:                    MetadataHttpUsecase,
		RecordbuzUsecase:                       RecordbuzUsecase,
		RecordHttpUsecase:                      RecordHttpUsecase,
		SettingSectionFieldUsecase:             SettingSectionFieldUsecase,
		KindUsecase:                            KindUsecase,
		TimelineUsecase:                        TimelineUsecase,
		NotesUsecase:                           NotesUsecase,
		TimelinesbuzUsecase:                    TimelinesbuzUsecase,
		RatingPaymentUsecase:                   RatingPaymentUsecase,
		MailFeeContentUsecase:                  MailFeeContentUsecase,
		BUsecase:                               BUsecase,
		FieldPermissionUsecase:                 FieldPermissionUsecase,
		FieldValidatorUsecase:                  FieldValidatorUsecase,
		ManualThingstoknowUsecase:              ManualThingstoknowUsecase,
		TimezonesUsecase:                       TimezonesUsecase,
		WordUsecase:                            WordUsecase,
		ZohoDealScan2JobUsecase:                ZohoDealScan2JobUsecase,
		SettingHttpUsecase:                     SettingHttpUsecase,
		ClientTaskHandleWhatGidJobUsecase:      ClientTaskHandleWhatGidJobUsecase,
		ClientTaskHandleWhoGidJobUsecase:       ClientTaskHandleWhoGidJobUsecase,
		QueueUsecase:                           QueueUsecase,
		RecordbuzSearchUsecase:                 RecordbuzSearchUsecase,
		TaskHttpUsecase:                        TaskHttpUsecase,
		ReissueTriggerStrRequestPendingUsecase: ReissueTriggerStrRequestPendingUsecase,
		AutomaticTaskCreationUsecase:           AutomaticTaskCreationUsecase,
		AutomaticUpdateDueDateUsecase:          AutomaticUpdateDueDateUsecase,
		NotificationHttpUsecase:                NotificationHttpUsecase,
		ZohoCollaboratorUsecase:                ZohoCollaboratorUsecase,
		UnsubscribesbuzUsecase:                 UnsubscribesbuzUsecase,
		UnsubscribesUsecase:                    UnsubscribesUsecase,
		DialpadbuzInternalUsecase:              DialpadbuzInternalUsecase,
		UnsubscribesHttpUsecase:                UnsubscribesHttpUsecase,
		ConditionSourcebuzUsecase:              ConditionSourcebuzUsecase,
		AiPromptUsecase:                        AiPromptUsecase,
		ConditionRelaAiUsecase:                 ConditionRelaAiUsecase,
		ConditionLogAiUsecase:                  ConditionLogAiUsecase,
		ConditionHttpUsecase:                   ConditionHttpUsecase,
		ConditionCategoryUsecase:               ConditionCategoryUsecase,
		FilterbuzUsecase:                       FilterbuzUsecase,
		RecordLogbuzUsecase:                    RecordLogbuzUsecase,
		ZohoNoteScanJobUsecase:                 ZohoNoteScanJobUsecase,
		InvokeLogUsecase:                       InvokeLogUsecase,
		InvokeLogJobUsecase:                    InvokeLogJobUsecase,
		ClientCaseSyncbuzUsecase:               ClientCaseSyncbuzUsecase,
		ChangeHistoryNodelayJobUseacse:         ChangeHistoryNodelayJobUseacse,
		DueDateUsecase:                         DueDateUsecase,
		VbcAIUsecase:                           VbcAIUsecase,
		AiUsecase:                              AiUsecase,
		AiHttpUsecase:                          AiHttpUsecase,
		DocEmailUsecase:                        DocEmailUsecase,
		MedicalDbqCostUsecase:                  MedicalDbqCostUsecase,
		ClientCasebuzUsecase:                   ClientCasebuzUsecase,
		QuestionnairesbuzUsecase:               QuestionnairesbuzUsecase,
		CommonHttpUsecase:                      CommonHttpUsecase,
		OptionUsecase:                          OptionUsecase,
		StatementUsecase:                       StatementUsecase,
		UserHttpUsecase:                        UserHttpUsecase,
		ClientNameChangeJobUsecase:             ClientNameChangeJobUsecase,
		ResourceUsecase:                        ResourceUsecase,
		AdobepdfUsecase:                        AdobepdfUsecase,
		RoleUsecase:                            RoleUsecase,
		LeadVSChangeUsecase:                    LeadVSChangeUsecase,
		ChangeHisUsecase:                       ChangeHisUsecase,
		AiTaskbuzUsecase:                       AiTaskbuzUsecase,
		CacheLogUsecase:                        CacheLogUsecase,
		ZoomUploadBoxUsecase:                   ZoomUploadBoxUsecase,
		WordbuzUsecase:                         WordbuzUsecase,
		PsbuzUsecase:                           PsbuzUsecase,
		AiResultUsecase:                        AiResultUsecase,
		HttpBlobUsecase:                        HttpBlobUsecase,
		WebpUsecase:                            WebpUsecase,
		FilebuzUsecase:                         FilebuzUsecase,
		ContractbuzUsecase:                     ContractbuzUsecase,
		VbcDataVerifyUsecase:                   VbcDataVerifyUsecase,
		ManualUsecase:                          ManualUsecase,
		LeadConversionSummaryBuzUsecase:        LeadConversionSummaryBuzUsecase,
		ReferrerLogUsecase:                     ReferrerLogUsecase,
		StatemtUsecase:                         StatemtUsecase,
		LeadsUsecase:                           LeadsUsecase,
		ClientReviewUsecase:                    ClientReviewUsecase,
		ClientReviewBuzUsecase:                 ClientReviewBuzUsecase,
		StatementConditionBuzUsecase:           StatementConditionBuzUsecase,
		StatementConditionUsecase:              StatementConditionUsecase,
		StatementCommentBuzUsecase:             StatementCommentBuzUsecase,
		BoxUserBuzUsecase:                      BoxUserBuzUsecase,
		BoxCollaborationBuzUsecase:             BoxCollaborationBuzUsecase,
		AssistantUsecase:                       AssistantUsecase,
		AiAssistantJobBuzUsecase:               AiAssistantJobBuzUsecase,
		AttorneybuzUsecase:                     AttorneybuzUsecase,
		AttorneyUsecase:                        AttorneyUsecase,
		NotesbuzUsecase:                        NotesbuzUsecase,
		SendVa2122aUsecase:                     SendVa2122aUsecase,
		GlobalEventBusBuzUsecase:               GlobalEventBusBuzUsecase,
		CmdUsecase:                             CmdUsecase,
		PersonalWebformUsecase:                 PersonalWebformUsecase,
		MonitoredEmailsUsecase:                 MonitoredEmailsUsecase,
		MonitoredEmailsJobUsecase:              MonitoredEmailsJobUsecase,
		MonitoredEmailsTasksUsecase:            MonitoredEmailsTasksUsecase,
		ExportUsecase:                          ExportUsecase,
		CollaboratorbuzUsecase:                 CollaboratorbuzUsecase,
		UserbuzUsecase:                         UserbuzUsecase,
	}
}

func AppMain() *UnittestApp {
	os.Setenv("ENV", "DEV")
	flag.Parse()
	configs.InitApp(configs.App_UnitTest)
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"service.id", id,
		//"service.name", Name,
		//"service.version", Version,
		//"trace_id", tracing.TraceID(),
		//"span_id", tracing.SpanID(),
	)
	flagconf := "/configs/config_dev.yaml"
	a := os.Getenv("VBCONFIG_TEST")
	if len(a) > 0 {
		flagconf = a
	}
	var useProd bool

	useProd = false

	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	if useProd {
		flagconf = "/configs/config_prod.yaml"
		envPath := filepath.Join(baseDir, ".env_prod")
		err := godotenv.Load(envPath)
		if err != nil {
			panic("Error loading .env_prod file")
		}
	} else {
		envPath := filepath.Join(baseDir, ".env_dev")
		err := godotenv.Load(envPath)
		if err != nil {
			log.Fatal("Error loading .env_dev file")
		}
	}
	fmt.Println(flagconf)
	flagconf = lib.GetSourceRootPath() + flagconf
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, i map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, i)
		}),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}
	fmt.Println(c)
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	ConfData = bc.Data
	fmt.Println(ConfData)
	app, _, _ := initApp(bc.Data, logger)
	UT = app
	return app
}
