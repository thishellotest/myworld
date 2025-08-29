package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
)

/*
CREATE TABLE `adobe_webhook_events` (
  `id` int NOT NULL DEFAULT '0',
  `webhook_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `webhook_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `event` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `event_date` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `event_resource_type` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `participant_role` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `participant_user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `participant_user_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `acting_user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `acting_user_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `initiating_user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `initiating_user_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `agreement_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `agreement_name` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `agreement_status` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` int NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
*/

type AdobeWebhookEventEntity struct {
	ID                   int32 `gorm:"primaryKey"`
	WebhookId            string
	WebhookName          string
	Event                string
	EventDate            string
	EventResourceType    string
	ParticipantRole      string
	ParticipantUserId    string
	ParticipantUserEmail string
	ActingUserId         string
	ActingUserEmail      string
	InitiatingUserId     string
	InitiatingUserEmail  string
	AgreementId          string
	AgreementName        string
	AgreementStatus      string
	CreatedAt            int64
}

func (AdobeWebhookEventEntity) TableName() string {
	return "adobe_webhook_events"
}

type AdobeWebhookEventUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	MapUsecase    *MapUsecase
	DBUsecase[AdobeWebhookEventEntity]
	AdobesignSyncTaskUsecase *AdobesignSyncTaskUsecase
}

func NewAdobeWebhookEventUsecase(logger log.Logger, conf *conf.Data, CommonUsecase *CommonUsecase, MapUsecase *MapUsecase,
	AdobesignSyncTaskUsecase *AdobesignSyncTaskUsecase) *AdobeWebhookEventUsecase {
	uc := &AdobeWebhookEventUsecase{
		log:                      log.NewHelper(logger),
		conf:                     conf,
		CommonUsecase:            CommonUsecase,
		MapUsecase:               MapUsecase,
		AdobesignSyncTaskUsecase: AdobesignSyncTaskUsecase,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *AdobeWebhookEventUsecase) RunAdobeWebhookEventJob(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("AdobeWebhookEventUsecase:RunAdobeWebhookEventJob:Done")
				return
			default:
				divideId, err := c.MapUsecase.GetForInt(Map_AdobeWebhookEvent_divide)
				if err != nil {
					c.log.Error(err)
				} else {
					sqlRows, err := c.CommonUsecase.DB().Table(AdobeWebhookEventEntity{}.TableName()).
						Where("id>?",
							divideId).Rows()
					if err != nil {
						c.log.Error(err)
					} else {
						if sqlRows != nil {
							newDivideId := int32(0)
							for sqlRows.Next() {
								var entity AdobeWebhookEventEntity
								err = c.CommonUsecase.DB().ScanRows(sqlRows, &entity)
								if err != nil {
									c.log.Error(err)
								} else {
									newDivideId = entity.ID
									err = c.Handle(&entity)
									if err != nil {
										c.log.Error(err)
									}
									c.MapUsecase.Set(Map_AdobeWebhookEvent_divide, lib.InterfaceToString(newDivideId))
								}
							}
							err = sqlRows.Close()
							if err != nil {
								c.log.Error(err)
							}
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return nil
}

//const ESign_type_adobe = "adobe"

// Handle 处理事件
func (c *AdobeWebhookEventUsecase) Handle(entity *AdobeWebhookEventEntity) error {
	if entity == nil {
		return errors.New("EnvelopeStatusChangeEntity is nil.")
	}
	if entity.AgreementStatus == "SIGNED" && entity.EventResourceType == "agreement" {
		//err := c.HandleCompleted(entity)
		typeMap := make(lib.TypeMap)
		typeMap.Set("AdobeWebhookEventId", entity.ID)
		typeMap.Set("type", Sign_type_adobe)
		err := c.CommonUsecase.DB().Save(&TaskEntity{
			IncrId:    entity.ID,
			TaskInput: typeMap.ToString(),
			Event:     Task_Dag_BoxCreateClientContracts,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}).Error

		if err != nil {
			a := GenLog(entity.ID, Log_FromType_AdobeSign_completed_failed, err.Error())
			err = c.CommonUsecase.DB().Save(a).Error
			if err != nil {
				c.log.Error(err)
			}
		}
	}
	if entity.EventResourceType == "agreement" {
		err := c.AdobesignSyncTaskUsecase.LPushSyncTaskQueue(context.Background(), entity.AgreementId)
		if err != nil {
			c.log.Error(err)
		}
	}
	return nil
}
