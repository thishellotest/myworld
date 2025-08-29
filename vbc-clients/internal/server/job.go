package server

import (
	"context"
	"fmt"
	"time"
	"vbc/configs"
	"vbc/internal/biz"
	"vbc/internal/conf"
	"vbc/lib"

	"github.com/go-kratos/kratos/v2/log"
)

var VbcJobManager *lib.JobManager

type Job struct {
	conf                              *conf.Data
	log                               *log.Helper
	SyncAsanaTaskUsecase              *biz.SyncAsanaTaskUsecase
	ChangeHistoryUseacse              *biz.ChangeHistoryUseacse
	TaskUsecase                       *biz.TaskUsecase
	Oauth2TokenUsecase                *biz.Oauth2TokenUsecase
	DocuSignUsecase                   *biz.DocuSignUsecase
	EnvelopeStatusChangeUsecase       *biz.EnvelopeStatusChangeUsecase
	AdobeWebhookEventUsecase          *biz.AdobeWebhookEventUsecase
	AdobesignSyncTaskUsecase          *biz.AdobesignSyncTaskUsecase
	FormResponseUsecase               *biz.FormResponseUsecase
	BoxWebhookLogUsecase              *biz.BoxWebhookLogUsecase
	TaskFailureLogJobUsecase          *biz.TaskFailureLogJobUsecase
	RollpoingJobUsecase               *biz.RollpoingJobUsecase
	WebhookLogJobUsecase              *biz.WebhookLogJobUsecase
	ZohoDealScanJobUsecase            *biz.ZohoDealScanJobUsecase
	ZohoContactScanJobUsecase         *biz.ZohoContactScanJobUsecase
	ZohoTaskScanJobUsecase            *biz.ZohoTaskScanJobUsecase
	CronUsecase                       *biz.CronUsecase
	ZohobuzUsecase                    *biz.ZohobuzUsecase
	RecordReviewJobUsecase            *biz.RecordReviewJobUsecase
	ReminderEventsJobUsecase          *biz.ReminderEventsJobUsecase
	ZohobuzTaskUsecase                *biz.ZohobuzTaskUsecase
	HaReportTaskJobUsecase            *biz.HaReportTaskJobUsecase
	BlobJobUsecase                    *biz.BlobJobUsecase
	BlobSliceJobUsecase               *biz.BlobSliceJobUsecase
	HaReportPageJobUsecase            *biz.HaReportPageJobUsecase
	CaseWithoutTaskUsecase            *biz.CaseWithoutTaskUsecase
	ZoombuzUsecase                    *biz.ZoombuzUsecase
	ZoomRecordingFileJobUsecase       *biz.ZoomRecordingFileJobUsecase
	ItfexpirationUsecase              *biz.ItfexpirationUsecase
	ZoomMeetingSmsNoticeJobUsecase    *biz.ZoomMeetingSmsNoticeJobUsecase
	ClientTaskBuzUsecase              *biz.ClientTaskBuzUsecase
	ZohoDealScan2JobUsecase           *biz.ZohoDealScan2JobUsecase
	ClientTaskHandleWhatGidJobUsecase *biz.ClientTaskHandleWhatGidJobUsecase
	ClientTaskHandleWhoGidJobUsecase  *biz.ClientTaskHandleWhoGidJobUsecase
	AutomaticUpdateDueDateUsecase     *biz.AutomaticUpdateDueDateUsecase
	NotesbuzUsecase                   *biz.NotesbuzUsecase
	ZohoCollaboratorUsecase           *biz.ZohoCollaboratorUsecase
	GlobalInjectUsecase               *biz.GlobalInjectUsecase
	ZohoNoteScanJobUsecase            *biz.ZohoNoteScanJobUsecase
	InvokeLogJobUsecase               *biz.InvokeLogJobUsecase
	ChangeHistoryNodelayJobUseacse    *biz.ChangeHistoryNodelayJobUseacse
	DueDateUsecase                    *biz.DueDateUsecase
	AiTaskJobUsecase                  *biz.AiTaskJobUsecase
	CmdUsecase                        *biz.CmdUsecase
	ClientNameChangeJobUsecase        *biz.ClientNameChangeJobUsecase
	RemindUsecase                     *biz.RemindUsecase
	MonitorUsecase                    *biz.MonitorUsecase
	BoxUserBuzUsecase                 *biz.BoxUserBuzUsecase
	MonitoredEmailsJobUsecase         *biz.MonitoredEmailsJobUsecase
}

func NewJob(conf *conf.Data,
	logger log.Logger,
	SyncAsanaTaskUsecase *biz.SyncAsanaTaskUsecase,
	ChangeHistoryUseacse *biz.ChangeHistoryUseacse,
	TaskUsecase *biz.TaskUsecase,
	Oauth2TokenUsecase *biz.Oauth2TokenUsecase,
	DocuSignUsecase *biz.DocuSignUsecase,
	EnvelopeStatusChangeUsecase *biz.EnvelopeStatusChangeUsecase,
	AdobeWebhookEventUsecase *biz.AdobeWebhookEventUsecase,
	AdobesignSyncTaskUsecase *biz.AdobesignSyncTaskUsecase,
	FormResponseUsecase *biz.FormResponseUsecase,
	BoxWebhookLogUsecase *biz.BoxWebhookLogUsecase,
	TaskFailureLogJobUsecase *biz.TaskFailureLogJobUsecase,
	RollpoingJobUsecase *biz.RollpoingJobUsecase,
	WebhookLogJobUsecase *biz.WebhookLogJobUsecase,
	ZohoDealScanJobUsecase *biz.ZohoDealScanJobUsecase,
	ZohoContactScanJobUsecase *biz.ZohoContactScanJobUsecase,
	ZohoTaskScanJobUsecase *biz.ZohoTaskScanJobUsecase,
	CronUsecase *biz.CronUsecase,
	ZohobuzUsecase *biz.ZohobuzUsecase,
	RecordReviewJobUsecase *biz.RecordReviewJobUsecase,
	ReminderEventsJobUsecase *biz.ReminderEventsJobUsecase,
	ZohobuzTaskUsecase *biz.ZohobuzTaskUsecase,
	HaReportTaskJobUsecase *biz.HaReportTaskJobUsecase,
	BlobJobUsecase *biz.BlobJobUsecase,
	BlobSliceJobUsecase *biz.BlobSliceJobUsecase,
	HaReportPageJobUsecase *biz.HaReportPageJobUsecase,
	CaseWithoutTaskUsecase *biz.CaseWithoutTaskUsecase,
	ZoombuzUsecase *biz.ZoombuzUsecase,
	ZoomRecordingFileJobUsecase *biz.ZoomRecordingFileJobUsecase,
	ItfexpirationUsecase *biz.ItfexpirationUsecase,
	ZoomMeetingSmsNoticeJobUsecase *biz.ZoomMeetingSmsNoticeJobUsecase,
	ClientTaskBuzUsecase *biz.ClientTaskBuzUsecase,
	ZohoDealScan2JobUsecase *biz.ZohoDealScan2JobUsecase,
	ClientTaskHandleWhatGidJobUsecase *biz.ClientTaskHandleWhatGidJobUsecase,
	ClientTaskHandleWhoGidJobUsecase *biz.ClientTaskHandleWhoGidJobUsecase,
	AutomaticUpdateDueDateUsecase *biz.AutomaticUpdateDueDateUsecase,
	NotesbuzUsecase *biz.NotesbuzUsecase,
	ZohoCollaboratorUsecase *biz.ZohoCollaboratorUsecase,
	GlobalInjectUsecase *biz.GlobalInjectUsecase,
	ZohoNoteScanJobUsecase *biz.ZohoNoteScanJobUsecase,
	InvokeLogJobUsecase *biz.InvokeLogJobUsecase,
	ChangeHistoryNodelayJobUseacse *biz.ChangeHistoryNodelayJobUseacse,
	DueDateUsecase *biz.DueDateUsecase,
	AiTaskJobUsecase *biz.AiTaskJobUsecase,
	CmdUsecase *biz.CmdUsecase,
	ClientNameChangeJobUsecase *biz.ClientNameChangeJobUsecase,
	RemindUsecase *biz.RemindUsecase,
	MonitorUsecase *biz.MonitorUsecase,
	BoxUserBuzUsecase *biz.BoxUserBuzUsecase,
	MonitoredEmailsJobUsecase *biz.MonitoredEmailsJobUsecase,

) *Job {

	return &Job{
		conf:                        conf,
		log:                         log.NewHelper(logger),
		SyncAsanaTaskUsecase:        SyncAsanaTaskUsecase,
		ChangeHistoryUseacse:        ChangeHistoryUseacse,
		TaskUsecase:                 TaskUsecase,
		Oauth2TokenUsecase:          Oauth2TokenUsecase,
		DocuSignUsecase:             DocuSignUsecase,
		EnvelopeStatusChangeUsecase: EnvelopeStatusChangeUsecase,
		AdobeWebhookEventUsecase:    AdobeWebhookEventUsecase,
		AdobesignSyncTaskUsecase:    AdobesignSyncTaskUsecase,
		FormResponseUsecase:         FormResponseUsecase,
		BoxWebhookLogUsecase:        BoxWebhookLogUsecase,
		TaskFailureLogJobUsecase:    TaskFailureLogJobUsecase,
		RollpoingJobUsecase:         RollpoingJobUsecase,
		WebhookLogJobUsecase:        WebhookLogJobUsecase,
		//GoogleSheetSyncTaskUsecase:  GoogleSheetSyncTaskUsecase,
		ZohoDealScanJobUsecase:            ZohoDealScanJobUsecase,
		ZohoContactScanJobUsecase:         ZohoContactScanJobUsecase,
		ZohoTaskScanJobUsecase:            ZohoTaskScanJobUsecase,
		CronUsecase:                       CronUsecase,
		ZohobuzUsecase:                    ZohobuzUsecase,
		RecordReviewJobUsecase:            RecordReviewJobUsecase,
		ReminderEventsJobUsecase:          ReminderEventsJobUsecase,
		ZohobuzTaskUsecase:                ZohobuzTaskUsecase,
		HaReportTaskJobUsecase:            HaReportTaskJobUsecase,
		BlobJobUsecase:                    BlobJobUsecase,
		BlobSliceJobUsecase:               BlobSliceJobUsecase,
		HaReportPageJobUsecase:            HaReportPageJobUsecase,
		CaseWithoutTaskUsecase:            CaseWithoutTaskUsecase,
		ZoombuzUsecase:                    ZoombuzUsecase,
		ZoomRecordingFileJobUsecase:       ZoomRecordingFileJobUsecase,
		ItfexpirationUsecase:              ItfexpirationUsecase,
		ZoomMeetingSmsNoticeJobUsecase:    ZoomMeetingSmsNoticeJobUsecase,
		ClientTaskBuzUsecase:              ClientTaskBuzUsecase,
		ZohoDealScan2JobUsecase:           ZohoDealScan2JobUsecase,
		ClientTaskHandleWhatGidJobUsecase: ClientTaskHandleWhatGidJobUsecase,
		ClientTaskHandleWhoGidJobUsecase:  ClientTaskHandleWhoGidJobUsecase,
		AutomaticUpdateDueDateUsecase:     AutomaticUpdateDueDateUsecase,
		NotesbuzUsecase:                   NotesbuzUsecase,
		ZohoCollaboratorUsecase:           ZohoCollaboratorUsecase,
		GlobalInjectUsecase:               GlobalInjectUsecase,
		ZohoNoteScanJobUsecase:            ZohoNoteScanJobUsecase,
		InvokeLogJobUsecase:               InvokeLogJobUsecase,
		ChangeHistoryNodelayJobUseacse:    ChangeHistoryNodelayJobUseacse,
		DueDateUsecase:                    DueDateUsecase,
		AiTaskJobUsecase:                  AiTaskJobUsecase,
		CmdUsecase:                        CmdUsecase,
		ClientNameChangeJobUsecase:        ClientNameChangeJobUsecase,
		RemindUsecase:                     RemindUsecase,
		MonitorUsecase:                    MonitorUsecase,
		BoxUserBuzUsecase:                 BoxUserBuzUsecase,
		MonitoredEmailsJobUsecase:         MonitoredEmailsJobUsecase,
	}
}

func (c *Job) Run(ctx context.Context) error {
	c.log.Info("Job Run: ", " AppJobType: ", configs.AppJobType(), " AppEnv: ", configs.AppEnv())
	var err error

	if configs.IsProd() && configs.IsJobTypeDefault() {
		err = c.Oauth2TokenUsecase.RunRefreshTokenJob(ctx)
		if err != nil {
			panic(err)
		}
		err = c.AiTaskJobUsecase.RunHandleCustomJob(ctx, 1, 3*time.Second,
			c.AiTaskJobUsecase.WaitingTasks,
			c.AiTaskJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 8 * * 1", func() { //
			er := c.RemindUsecase.HandleCreateTaskForITFExpirations()
			if er != nil {
				c.log.Error("HandleCreateTaskForITFExpirations: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 8,9,10 * * *", func() { //
			er := c.ItfexpirationUsecase.HandleITFExpireReminder()
			if er != nil {
				c.log.Error("HandleITFExpireReminder: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		err = c.InvokeLogJobUsecase.RunHandleCustomJob(context.TODO(), 1, 3*time.Second,
			c.InvokeLogJobUsecase.WaitingTasks,
			c.InvokeLogJobUsecase.Handle,
		)
		if err != nil {
			panic(err)
		}

		go func() {
			for {
				fmt.Println("heartbeat")
				time.Sleep(30 * time.Second)
			}
		}()

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 * * * *", func() { // 每小时执行一次-同步users
			er := c.ZoombuzUsecase.ExecuteSyncZoomUsers()
			if er != nil {
				c.log.Error("ExecuteSyncZoomUsers: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 * * * *", func() { // 每小时执行一次-同步users
			er := c.ZoombuzUsecase.ExecuteSyncRecords()
			if er != nil {
				c.log.Error("ExecuteSyncRecords: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 20 * * * *", func() { // 每小时执行一次-同步users
			er := c.BoxWebhookLogUsecase.CrontabEveryOneHourHandleQuestionnaireDownloads()
			if er != nil {
				c.log.Error("CrontabEveryOneHourHandleQuestionnaireDownloads: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 0 * * * *", func() { // 每小时执行一次-删除
			er := c.ZoombuzUsecase.ExecuteDeleteMeetingRecording()
			if er != nil {
				c.log.Error("ExecuteDeleteMeetingRecording: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles */30 * * * *", func() { // 每小时执行一次-删除
			er := c.MonitorUsecase.DoMonitorVSUsers()
			if er != nil {
				c.log.Error("ExecuteDeleteMeetingRecording: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles */3 * * * *", func() { //
			er := c.ZoomRecordingFileJobUsecase.ExecuteCrontabHandleProcessingRecording()
			if er != nil {
				c.log.Error("ExecuteCrontabHandleProcessingRecording: ", er)
			}
		})
		if err != nil {
			panic(err)
		}

		err = c.TaskUsecase.RunTaskJob(ctx)
		if err != nil {
			panic(err)
		}
		err = c.FormResponseUsecase.RunHandleJob(ctx)
		if err != nil {
			panic(err)
		}
		err = c.BoxWebhookLogUsecase.RunHandleJob(ctx)
		if err != nil {
			panic(err)
		}

		err = c.RollpoingJobUsecase.RunHandleCustomJob(context.TODO(), 2, 0,
			c.RollpoingJobUsecase.WaitingTasks,
			c.RollpoingJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		err = c.WebhookLogJobUsecase.RunHandleCustomJob(ctx, 1, 0, c.WebhookLogJobUsecase.WaitingTasks, c.WebhookLogJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		// 暂时关闭
		err = c.TaskFailureLogJobUsecase.RunHandleCustomJob(ctx, 2, 0,
			c.TaskFailureLogJobUsecase.WaitingTasks,
			c.TaskFailureLogJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		err = c.RecordReviewJobUsecase.RunCustomTaskJob(ctx)
		if err != nil {
			panic(err)
		}
		err = c.ReminderEventsJobUsecase.RunHandleJob(ctx)
		if err != nil {
			panic(err)
		}

		err = c.HaReportTaskJobUsecase.RunHandleCustomJob(ctx, 2, 10*time.Second,
			c.HaReportTaskJobUsecase.WaitingTasks,
			c.HaReportTaskJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		//err = c.HaReportTaskJobUsecase.RunHandleCustomJob(ctx, 2, 10*time.Second,
		//	c.HaReportTaskJobUsecase.WaitingTasksByCreatingPdf,
		//	c.HaReportTaskJobUsecase.HandleByCreatingPdf)
		//if err != nil {
		//	panic(err)
		//}

		// todo:lgl 任务先暂停
		//err = c.BlobJobUsecase.RunHandleCustomJob(ctx, 2, 10*time.Second,
		//	c.BlobJobUsecase.WaitingTasks,
		//	c.BlobJobUsecase.Handle)
		//if err != nil {
		//	panic(err)
		//}

		err = c.BlobSliceJobUsecase.RunHandleCustomJob(ctx, 2, 10*time.Second,
			c.BlobSliceJobUsecase.WaitingTasks,
			c.BlobSliceJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		// 并发3，已经会出现请求超限了
		/*
			{
			  "error": {
			    "code": "429",
			    "message": "Requests to the ChatCompletions_Create Operation under Azure OpenAI API version 2024-05-01-preview have exceeded token rate limit of your current OpenAI S0 pricing tier. Please retry after 60 seconds. Please go here: https://aka.ms/oai/quotaincrease if you would like to further increase the default rate limit."
			  }
			}
		*/
		//err = c.HaReportPageJobUsecase.RunHandleCustomJob(context.TODO(), 2, 10*time.Second, c.HaReportPageJobUsecase.WaitingTasks,
		//	c.HaReportPageJobUsecase.Handle)
		//if err != nil {
		//	panic(err)
		//}

		err = c.ZoomRecordingFileJobUsecase.RunHandleCustomJob(context.TODO(),
			1,
			time.Second*10,
			c.ZoomRecordingFileJobUsecase.WaitingTasks,
			c.ZoomRecordingFileJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		err = c.ZoomMeetingSmsNoticeJobUsecase.Run(ctx)
		if err != nil {
			panic(err)
		}

	} else if configs.IsProd() && configs.IsJobTypeLargeMemory() {
		fmt.Println("IsJobTypeLargeMemory job")
		err = c.ZoomRecordingFileJobUsecase.RunHandleCustomJob(context.TODO(),
			1,
			time.Second*60,
			c.ZoomRecordingFileJobUsecase.WaitingTasks,
			c.ZoomRecordingFileJobUsecase.Handle)
		if err != nil {
			panic(err)
		}

		// 每天凌晨 3 点 5 分执行一次任务
		_, err = c.CronUsecase.Cron().AddFunc("CRON_TZ=America/Los_Angeles 40 3 * * *", func() { //
			er := c.CmdUsecase.RunBackup()
			if er != nil {
				c.log.Error("CmdUsecase RunBackup: ", er)
			}
		})
		if err != nil {
			panic(err)
		}
	}

	if (configs.IsProd() && configs.IsJobTypeDefault()) || configs.IsTest() || configs.IsDev() {

		if configs.Enable_Client_Task_ForCRM {
			c.log.Info("ClientTaskHandleWhatGidJobUsecase Running")
			err = c.ClientTaskHandleWhatGidJobUsecase.RunCustomTaskJob(ctx)
			if err != nil {
				panic(err)
			}
			c.log.Info("ClientTaskHandleWhoGidJobUsecase Running")
			err = c.ClientTaskHandleWhoGidJobUsecase.RunCustomTaskJob(ctx)
			if err != nil {
				panic(err)
			}
		}
		err = c.ChangeHistoryUseacse.RunChangeHistoryJob(ctx)
		if err != nil {
			panic(err)
		}

		// Cron: monitored emails listener and auto-forwarding flow (every 10 seconds)
		_, err = c.CronUsecase.Cron().AddFunc("@every 10s", func() {
			er := c.MonitoredEmailsJobUsecase.Handle()
			if er != nil {
				c.log.Error("MonitoredEmailsJobUsecase.Handle: ", er)
			}
		})
		if err != nil {
			panic(err)
		}
	}

	if configs.IsProd() && configs.IsJobTypeDefault() {
		err = c.ClientNameChangeJobUsecase.RunCustomTaskJob(ctx)
		if err != nil {
			panic(err)
		}
	}

	// 任务只能启动一个
	if configs.IsJobTypeDefault() {
		err = c.ChangeHistoryNodelayJobUseacse.RunChangeHistoryNodelayJobJob(ctx)
		if err != nil {
			panic(err)
		}

	}

	return nil
}
