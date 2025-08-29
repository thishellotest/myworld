package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"net/http"
	"strings"
	"time"
	"vbc/configs"
	"vbc/lib"
	"vbc/lib/static"
)

func customMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		//if tr, ok := transport.FromServerContext(ctx); ok {
		//fmt.Println("operation:", tr.Operation())
		//}
		reply, err = handler(ctx, req)
		return
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		if strings.Index(origin, "x.com") >= 0 ||
			strings.Index(origin, "localhost") >= 0 ||
			strings.Index(origin, "20.190.193.236") >= 0 ||
			strings.Index(origin, "vbctest-d4fbb39e6b9f.herokuapp.com") >= 0 ||
			strings.Index(origin, "vetbenefitscenter.com") >= 0 ||
			strings.Index(origin, "vercel.app") >= 0 ||
			strings.Index(origin, "veteranbenefitscenter.com") >= 0 ||
			strings.Index(origin, "augustusmiles.com") >= 0 ||
			strings.Index(origin, "v0.dev") >= 0 ||
			strings.Index(origin, "vusercontent.net") >= 0 { // vusercontent.net v0的priview网站
			c.Header("Access-Control-Allow-Origin", origin)
			//c.Header("Access-Control-Allow-Origin", "https://ailive-test2.magics.plus") // 使用*，不通过，需要指定域名
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 需要加入这个，否则不通过
		}

		c.Next()
	}
}

func NewGinRegisterHTTPServer(ApiUsecase *ApiUsecase,
	CommonUsecase *CommonUsecase,
	logger log.Logger,
	HttpWebhookUsecase *HttpWebhookUsecase,
	HttpTUsecase *HttpTUsecase,
	HttpOauth2Usecase *HttpOauth2Usecase,
	HttpApiUsecase *HttpApiUsecase,
	HttpManualUsecase *HttpManualUsecase,
	HttpAccessControl *HttpAccessControl,
	ZohobuzUsecase *ZohobuzUsecase,
	AsanaMigrateUsecase *AsanaMigrateUsecase,
	ZohoinfoSyncUsecase *ZohoinfoSyncUsecase,
	ActionOnceUsecase *ActionOnceUsecase,
	ClientCaseContractBasicDataUsecase *ClientCaseContractBasicDataUsecase,
	HaReportTasksBuzUsecase *HaReportTasksBuzUsecase,
	LoginUsecase *LoginUsecase,
	JWTUsecase *JWTUsecase,
	HttpUserUsecase *HttpUserUsecase,
	HttpBlobUsecase *HttpBlobUsecase,
	ZoomRecordingFileUsecase *ZoomRecordingFileUsecase,
	ZoomRecordingFileJobUsecase *ZoomRecordingFileJobUsecase,
	Awsclaude3Usecase *Awsclaude3Usecase,
	HttpSettingsUsecase *HttpSettingsUsecase,
	MiscUsecase *MiscUsecase,
	JotformbuzUsecase *JotformbuzUsecase,
	QuestionnairesUsecase *QuestionnairesUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	MetadataUsecase *MetadataUsecase,
	ConditionbuzUsecase *ConditionbuzUsecase,
	RecordHttpUsecase *RecordHttpUsecase,
	MetadataHttpUsecase *MetadataHttpUsecase,
	SettingHttpUsecase *SettingHttpUsecase,
	NotesHttpUsecase *NotesHttpUsecase,
	TaskHttpUsecase *TaskHttpUsecase,
	FlowHttpUsecase *FlowHttpUsecase,
	NotificationHttpUsecase *NotificationHttpUsecase,
	EventstreamHttpUsecase *EventstreamHttpUsecase,
	ZohoCollaboratorUsecase *ZohoCollaboratorUsecase,
	UnsubscribesHttpUsecase *UnsubscribesHttpUsecase,
	ConditionSourcebuzUsecase *ConditionSourcebuzUsecase,
	AiHttpUsecase *AiHttpUsecase,
	AiUsecase *AiUsecase,
	ConditionHttpUsecase *ConditionHttpUsecase,
	ZohoNoteScanJobUsecase *ZohoNoteScanJobUsecase,
	DueDateUsecase *DueDateUsecase,
	WebsocketUsecase *WebsocketUsecase,
	JotformSubmissionUsecase *JotformSubmissionUsecase,
	CommonHttpUsecase *CommonHttpUsecase,
	ManualThingstoknowUsecase *ManualThingstoknowUsecase,
	CmdUsecase *CmdUsecase,
	UserHttpUsecase *UserHttpUsecase,
	MgmtHttpUsecase *MgmtHttpUsecase,
	ZoomUploadBoxUsecase *ZoomUploadBoxUsecase,
	PsHttpUsecase *PsHttpUsecase,
	PdfGoFitzUsecase *PdfGoFitzUsecase,
	RemindUsecase *RemindUsecase,
	ContractHttpUsecase *ContractHttpUsecase,
	VbcDataVerifyUsecase *VbcDataVerifyUsecase,
	MonitorUsecase *MonitorUsecase,
	LeadConversionSummaryBuzUsecase *LeadConversionSummaryBuzUsecase,
	AiTaskbuzUsecase *AiTaskbuzUsecase,
	AiAssistantJobBuzUsecase *AiAssistantJobBuzUsecase,
	BoxWebhookLogUsecase *BoxWebhookLogUsecase,
	ExportUsecase *ExportUsecase,
	PersonalWebformUsecase *PersonalWebformUsecase,
) *gin.Engine {
	router := gin.Default()
	router.Use(kgin.Middlewares(recovery.Recovery(), customMiddleware))
	router.Use(Cors())

	if configs.IsDev() {
		router.Use(static.Serve("/", static.LocalFile("./front/prod", false)))
		router.Use(static.Serve("/assets", static.LocalFile("./front/assets", false)))
	} else if configs.AppEnv() == configs.ENV_PROD {
		if configs.IsJobTypeQA() {
			router.Use(static.Serve("/", static.LocalFile("/app/front/qa", false)))
		} else {
			router.Use(static.Serve("/", static.LocalFile("/app/front/prod", false)))
		}
		router.Use(static.Serve("/assets", static.LocalFile("/app/front/assets", false)))
	} else if configs.AppEnv() == configs.ENV_TEST {
		router.Use(static.Serve("/", static.LocalFile("/app/front/test", false)))
		router.Use(static.Serve("/assets", static.LocalFile("/app/front/assets", false)))
	}

	log := log.NewHelper(logger)
	rootGrp := router.Group("/api")
	rootGrp.GET("export", ExportUsecase.Http)
	rootGrp.GET("ws", WebsocketUsecase.HttpHandleWS)
	rootGrp.GET("ws/send", WebsocketUsecase.HttpSendMessage)
	rootGrp.POST("AccessControl/Tasks", HttpAccessControl.Tasks)
	rootGrp.POST("AccessControl/CarryOut", HttpAccessControl.CarryOut)
	rootGrp.GET("_manual_/HandleOnce", func(c *gin.Context) {
		go func() {
			var err error
			//err = PersonalWebformUsecase.ManualHistoryData()
			//err = BoxWebhookLogUsecase.CrontabEveryOneHourHandleQuestionnaireDownloads()

			//AiTaskbuzUsecase.TestGenUpdatePSFromStatementCondition()

			//err = MiscUsecase.UpdateAll()
			//err = LeadConversionSummaryBuzUsecase.ManualAll()
			//MonitorUsecase.DoMonitorVSUsers()
			//VbcDataVerifyUsecase.VerifyContract()
			//err = ZohobuzUsecase.HandleAllMan()
			//if err != nil {
			//	c.Error(err)
			//}
			//err = RemindUsecase.HandleCreateTaskForITFExpirations()
			//PdfGoFitzUsecase.TestFitz()
			//err := CmdUsecase.RunBackup()

			// 更新：Exam Day Checklist v2.7.pdf
			//err = ManualThingstoknowUsecase.HandleUploadNewThingsToKnowFileAllCases()
			//err := DueDateUsecase.SyncDueDate()
			//err := ConditionbuzUsecase.HandleAllCondition()
			//err := JotformSubmissionUsecase.ManualHandleFormId()

			//err := ZoomUploadBoxUsecase.TestUpload()

			//id := c.Query("id")
			//if id != "" {
			//	log.Debug("HandleExecUseChunks test:", id)
			//	zoomRecording, _ := ZoomRecordingFileUsecase.GetByCond(builder.Eq{"id": id})
			//	err = ZoomRecordingFileJobUsecase.HandleExecUseChunks(context.TODO(), zoomRecording)
			//	log.Debug("HandleExecUseChunks err:", err)
			//}

			if err != nil {
				log.Error(err)
			}
		}()
	})

	rootGrp.GET("_manual_/test_env", func(c *gin.Context) {

	})

	rootGrp.GET("login", LoginUsecase.HttpLogin)
	rootGrp.POST("/unsubscribes/list", JWTUsecase.JWTAuthMiddleware(), UnsubscribesHttpUsecase.List)
	rootGrp.POST("/unsubscribes/change-status", JWTUsecase.JWTAuthMiddleware(), UnsubscribesHttpUsecase.ChangeStatus)
	rootGrp.POST("/unsubscribes/save", JWTUsecase.JWTAuthMiddleware(), UnsubscribesHttpUsecase.Save)
	rootGrp.POST("/unsubscribes/delete", JWTUsecase.JWTAuthMiddleware(), UnsubscribesHttpUsecase.Delete)

	conditionGrp := rootGrp.Group("condition")
	conditionGrp.POST("list", JWTUsecase.JWTAuthMiddleware(), ConditionHttpUsecase.List)
	conditionGrp.POST("sources", JWTUsecase.JWTAuthMiddleware(), ConditionHttpUsecase.Sources)
	conditionGrp.POST("delete", JWTUsecase.JWTAuthMiddleware(), ConditionHttpUsecase.Delete)

	man := rootGrp.Group("man")
	man.GET("MiscUsecase_UpdateAll", JWTUsecase.JWTAuthMiddleware(), func(c *gin.Context) {
		go func() {
			//MiscUsecase.UpdateAll()
		}()
	})
	man.GET("HandleQuestionnairesJotformHistory", JWTUsecase.JWTAuthMiddleware(), func(c *gin.Context) {
		err := JotformbuzUsecase.HandleQuestionnairesJotformHistory()
		if err != nil {
			c.Error(err)
		}
	})

	mgmtGrp := rootGrp.Group("mgmt")
	mgmtGrp.POST("init", JWTUsecase.JWTAuthMiddleware(), MgmtHttpUsecase.Init)
	mgmtGrp.POST("users/sync-dailpad", JWTUsecase.JWTAuthMiddleware(), UserHttpUsecase.SyncDailpad)
	mgmtGrp.POST("users/verify-email-outbox", JWTUsecase.JWTAuthMiddleware(), UserHttpUsecase.VerifyEmailOutbox)
	mgmtGrp.POST("contract/get", JWTUsecase.JWTAuthMiddleware(), ContractHttpUsecase.Get)
	mgmtGrp.POST("contract/save", JWTUsecase.JWTAuthMiddleware(), ContractHttpUsecase.Save)
	mgmtGrp.POST("contract/list", JWTUsecase.JWTAuthMiddleware(), ContractHttpUsecase.List)

	flowGrp := rootGrp.Group("flow")
	flowGrp.POST("stages/trigger", JWTUsecase.JWTAuthMiddleware(), FlowHttpUsecase.StagesTrigger)

	settingsGrp := rootGrp.Group("settings")
	settingsGrp.POST("custom_view/:module_name", JWTUsecase.JWTAuthMiddleware(), SettingHttpUsecase.CustomView)
	settingsGrp.POST("custom_view/:module_name/change_sort", JWTUsecase.JWTAuthMiddleware(), SettingHttpUsecase.ChangeSort)
	settingsGrp.POST("custom_view/:module_name/change_fields", JWTUsecase.JWTAuthMiddleware(), SettingHttpUsecase.ChangeFields)
	settingsGrp.POST("custom_view/:module_name/change_columnwidth", JWTUsecase.JWTAuthMiddleware(), SettingHttpUsecase.ChangeColumnwidth)

	notificationGrp := rootGrp.Group("notification")
	notificationGrp.POST("info", JWTUsecase.JWTAuthMiddleware(), NotificationHttpUsecase.Info)
	notificationGrp.POST("list", JWTUsecase.JWTAuthMiddleware(), NotificationHttpUsecase.List)
	notificationGrp.POST(":gid/read", JWTUsecase.JWTAuthMiddleware(), NotificationHttpUsecase.Read)

	userGrp := rootGrp.Group("user")
	userGrp.POST("sign_in", HttpUserUsecase.SignIn)
	userGrp.GET("info", JWTUsecase.JWTAuthMiddleware(), HttpUserUsecase.Info)

	blobGrp := rootGrp.Group("blob")
	blobGrp.POST("create_task", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.CreateTask)
	blobGrp.POST("task_list", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.TaskList)
	blobGrp.POST("detail", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.Detail)
	blobGrp.POST("comments", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.Comments)
	blobGrp.POST("save_comment", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.Save)
	blobGrp.POST("delete", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.Delete)
	blobGrp.POST("record-review-detail", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.RecordReviewDetail)
	blobGrp.POST("slice-join-ocr", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.SliceJoinOcr)
	//blobGrp.POST("record-review-tasks-progress", HttpBlobUsecase.RecordReviewTasksProgress)
	blobGrp.POST("record-review-files", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.RecordReviewFiles)
	blobGrp.POST("ocr-detail", JWTUsecase.JWTAuthMiddleware(), HttpBlobUsecase.OcrDetail)
	questionGrp := rootGrp.Group("question")
	questionGrp.POST("list", JWTUsecase.JWTAuthMiddleware(), JotformbuzUsecase.HttpQuestionList)

	questionnairesGrp := rootGrp.Group("questionnaires")
	questionnairesGrp.POST("list", JWTUsecase.JWTAuthMiddleware(), QuestionnairesUsecase.HttpList)

	casesGrp := rootGrp.Group("cases")
	casesGrp.POST("detail-by-id", JWTUsecase.JWTAuthMiddleware(), ClientCaseUsecase.HttpDetailById)
	casesGrp.POST("claims_info", JWTUsecase.JWTAuthMiddleware(), ClientCaseUsecase.HttpClaimsInfo)
	casesGrp.POST("detail", JWTUsecase.JWTAuthMiddleware(), ClientCaseUsecase.HttpDetail)
	casesGrp.POST("save", JWTUsecase.JWTAuthMiddleware(), ClientCaseUsecase.HttpSave)

	eventstreamGrp := rootGrp.Group("eventstream")
	eventstreamGrp.POST("handle/:jwt", EventstreamHttpUsecase.Handle)

	psGrp := rootGrp.Group("ps")
	psGrp.POST("generate-document", JWTUsecase.JWTAuthMiddleware(), PsHttpUsecase.GenerateDocument)
	psGrp.POST("handle-update-statement", JWTUsecase.JWTAuthMiddleware(), PsHttpUsecase.HandleUpdateStatement)

	aiGrp := rootGrp.Group("ai")
	aiGrp.POST("tasks", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.Tasks)
	aiGrp.POST("task-handle", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TaskHandle)
	aiGrp.POST("task-launch", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TaskLaunch)
	aiGrp.POST("task-renew", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TaskRenew)
	aiGrp.POST("task-delete", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TaskDelete)
	aiGrp.POST("task-result", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TaskResult)
	aiGrp.POST("claude3", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.Claude3)
	aiGrp.POST("test-ai", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.TestAi)
	aiGrp.POST("HandleOnceConditionSourceWithAi", JWTUsecase.JWTAuthMiddleware(), AiHttpUsecase.HandleOnceConditionSourceWithAi)
	aiGrp.POST("assistant/create/job", JWTUsecase.JWTAuthMiddleware(), AiAssistantJobBuzUsecase.HttpCreate)
	aiGrp.POST("assistant/get/job-status", JWTUsecase.JWTAuthMiddleware(), AiAssistantJobBuzUsecase.HttpGetJobStatus)
	aiGrp.GET("assistant/get/job-detail", JWTUsecase.JWTAuthMiddleware(), AiAssistantJobBuzUsecase.HttpGetJobDetail)
	aiGrp.POST("assistant/apply/job", JWTUsecase.JWTAuthMiddleware(), AiAssistantJobBuzUsecase.HttpApplyJob)
	aiGrp.POST("assistant/clear/job", JWTUsecase.JWTAuthMiddleware(), AiAssistantJobBuzUsecase.HttpClearJob)

	metadataGrp := rootGrp.Group("metadata")
	//metadataGrp.POST("basicdata", JWTUsecase.JWTAuthMiddleware(), MetadataUsecase.Basicdata)
	metadataGrp.GET("basicdata", JWTUsecase.JWTAuthMiddleware(), MetadataUsecase.Basicdata)
	metadataGrp.POST("conditions", JWTUsecase.JWTAuthMiddleware(), MetadataUsecase.HttpConditions)
	metadataGrp.POST("fields/:module_name", JWTUsecase.JWTAuthMiddleware(), MetadataHttpUsecase.Fields) // 搜索使用
	metadataGrp.POST("options/:module_name/:field_name", JWTUsecase.JWTAuthMiddleware(), MetadataHttpUsecase.Options)

	commonGrp := rootGrp.Group("common")
	commonGrp.POST(":common_type/get", JWTUsecase.JWTAuthMiddleware(), CommonHttpUsecase.Get)
	commonGrp.POST(":common_type/save", JWTUsecase.JWTAuthMiddleware(), CommonHttpUsecase.Save)

	rootGrp.POST("records/:module_name", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.List)
	rootGrp.POST("record/:module_name/edit/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Edit)     // 修改页面
	rootGrp.POST("record/:module_name/detail/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Detail) // 详情页面
	rootGrp.POST("record/:module_name/layout", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Layout)      // 创建页面
	rootGrp.POST("record/form/column", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.FormColumn)
	rootGrp.POST("record/:module_name/delete/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Delete) // 详情页面
	rootGrp.POST("record/:module_name/related/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Related)
	rootGrp.POST("record/:module_name/timelines/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Timelines)
	rootGrp.POST("record/:module_name/save/:gid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Save)
	rootGrp.POST("record/:module_name/create", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.Save)
	rootGrp.POST("records/related/client/:clientGid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.RelatedClient)
	rootGrp.POST("records/medical-cost/save/:caseGid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.MedicalCost)
	rootGrp.GET("records/statement/detail/:caseGid", RecordHttpUsecase.StatementDetail) // 为什么需要从Post改为Get，因为出现错误：H18 - Server Request Interrupted，
	rootGrp.GET("records/statement/detail-test/:caseGid", RecordHttpUsecase.StatementDetailTest)
	rootGrp.GET("records/statement/revert-version/:caseGid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.StatementRevertVersion)
	rootGrp.POST("records/statement/detail-versions/:caseGid", RecordHttpUsecase.StatementDetailVersions)                                                // jwt或密码访问，所以不能使用：JWTAuthMiddleware
	rootGrp.POST("records/statement/save/:caseGid", RecordHttpUsecase.StatementSave)                                                                     // jwt或密码访问，所以不能使用：JWTAuthMiddleware
	rootGrp.POST("records/statement/verify-password/:caseGid", RecordHttpUsecase.StatementVerifyPassword)                                                // jwt或密码访问，所以不能使用：JWTAuthMiddleware
	rootGrp.POST("records/statement/comment/save/:caseGid", RecordHttpUsecase.StatementCommentSave)                                                      // jwt或密码访问
	rootGrp.POST("records/statement/comment/mark-complete/:caseGid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.StatementCommentMarkComplete)     // jwt访问
	rootGrp.POST("records/statement/comment/unmark-complete/:caseGid", JWTUsecase.JWTAuthMiddleware(), RecordHttpUsecase.StatementCommentUnmarkComplete) // jwt访问
	rootGrp.GET("records/statement/comment/list/:caseGid", RecordHttpUsecase.StatementCommentList)
	rootGrp.POST("records/statement/comment/delete/:caseGid", RecordHttpUsecase.StatementCommentDelete)
	rootGrp.POST("records/statement/comment/submit-for-review/:caseGid", RecordHttpUsecase.StatementCommentSubmitForReview) // 密码访问

	rootGrp.POST("leads/save", RecordHttpUsecase.LeadsSave)
	rootGrp.GET("client/reviews", RecordHttpUsecase.ClientReviews)

	rootGrp.POST("task/:module_name/:kind_gid/list", JWTUsecase.JWTAuthMiddleware(), TaskHttpUsecase.List)
	rootGrp.POST("task/:module_name/complete/:gid", JWTUsecase.JWTAuthMiddleware(), TaskHttpUsecase.Complete)

	rootGrp.POST("notes/:module_name/:kind_gid/save", JWTUsecase.JWTAuthMiddleware(), NotesHttpUsecase.Save)
	rootGrp.POST("notes/:module_name/:kind_gid", JWTUsecase.JWTAuthMiddleware(), NotesHttpUsecase.List)
	rootGrp.POST("notes/delete", JWTUsecase.JWTAuthMiddleware(), NotesHttpUsecase.Delete)

	oauth2Grp := router.Group("/oauth2")
	oauth2Grp.GET("callback", HttpOauth2Usecase.Callback)
	oauth2Grp.GET("auth-list", HttpOauth2Usecase.AuthList)

	routeWebhook := router.Group("/webhook")
	routeWebhook.POST("receive_asana", HttpWebhookUsecase.ReceiveAsana)
	routeWebhook.POST("receive_box", HttpWebhookUsecase.ReceiveBox)
	routeWebhook.POST("post_source", func(c *gin.Context) {

		c.Header("x-xero-signature", "rWhT6E2d6SLIWQRFc6c4ao+Xq+3MneVF5RP2SVzdraa69wupsUFH+zQeSKzurUFayc1wjT9joHbSGhFXllh9Tw==")

		XAdobeSignClientId := c.Request.Header.Get("X-AdobeSign-ClientId")
		c.Header("X-AdobeSign-ClientId", XAdobeSignClientId)
		xHookSecret := c.Request.Header.Get("X-Hook-Secret")
		c.Request.Header.Set("X-Hook-Secret", xHookSecret)

		from := c.Query("from")
		webhookLog := WebhookLogEntity{
			From:      from,
			Query:     lib.InterfaceToString(c.Request.URL.Query()),
			Headers:   lib.InterfaceToString(c.Request.Header),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		if from == WebhookLog_From_jotform {
			formID := c.PostForm("formID")
			submissionID := c.PostForm("submissionID")
			neatBody := make(lib.TypeMap)
			neatBody.Set("formID", formID)
			neatBody.Set("submissionID", submissionID)
			webhookLog.Body = InterfaceToString(c.Request.Form)
			webhookLog.NeatBody = InterfaceToString(neatBody)
		} else {
			body, err := c.GetRawData()
			if err != nil {
				log.Error(err)
			}
			webhookLog.Body = string(body)
		}

		CommonUsecase.DB().Create(&webhookLog)
	})
	routeWebhook.GET("post_source", func(c *gin.Context) {
		c.Header("x-xero-signature", "rWhT6E2d6SLIWQRFc6c4ao+Xq+3MneVF5RP2SVzdraa69wupsUFH+zQeSKzurUFayc1wjT9joHbSGhFXllh9Tw==")
		XAdobeSignClientId := c.Request.Header.Get("X-AdobeSign-ClientId")
		c.Header("X-AdobeSign-ClientId", XAdobeSignClientId)
		xHookSecret := c.Request.Header.Get("X-Hook-Secret")
		c.Request.Header.Set("X-Hook-Secret", xHookSecret)

		body, err := c.GetRawData()
		if err != nil {
			log.Error(err)
		}
		webhookLog := WebhookLogEntity{
			Headers: lib.InterfaceToString(c.Request.Header),
			Body:    string(body),
		}
		CommonUsecase.DB().Create(&webhookLog)
	})
	routeWebhook.POST("receive_adobesign", HttpWebhookUsecase.ReceiveAdobeSign)
	routeWebhook.GET("receive_adobesign", func(c *gin.Context) {
		XAdobeSignClientId := c.Request.Header.Get("X-AdobeSign-ClientId")
		c.Header("X-AdobeSign-ClientId", XAdobeSignClientId)
	})

	routeWebhook.POST("form_responses", HttpWebhookUsecase.FormResponses)

	return router
}
