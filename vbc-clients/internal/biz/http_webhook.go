package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/lib"
)

type HttpWebhookUsecase struct {
	log                  *log.Helper
	SyncAsanaTaskUsecase *SyncAsanaTaskUsecase
	CommonUsecase        *CommonUsecase
}

func NewHttpWebhookUsecase(logger log.Logger, SyncAsanaTaskUsecase *SyncAsanaTaskUsecase, CommonUsecase *CommonUsecase) *HttpWebhookUsecase {
	return &HttpWebhookUsecase{
		log:                  log.NewHelper(logger),
		SyncAsanaTaskUsecase: SyncAsanaTaskUsecase,
		CommonUsecase:        CommonUsecase,
	}
}

// ReceiveAsana 接收webhook
func (c *HttpWebhookUsecase) ReceiveAsana(ctx *gin.Context) {

	//return
	a, err := ctx.GetRawData()
	if err != nil {
		c.log.Error(err)
	} else {
		r := lib.StringToT[*AsanaWebhookVo](string(a))
		if r.IsOk() {
			for _, v := range r.Unwrap().Events {
				if v.IsTaskWebhook() {
					err = c.SyncAsanaTaskUsecase.LPushSyncTaskQueue(ctx, v.Resource.Gid)
					if err != nil {
						c.log.Error(err)
					}
					//if v.User.Gid != "" {
					//	err = c.SyncAsanaTaskUsecase.UserLPushSyncTaskQueue(ctx, v.User.Gid)
					//	if err != nil {
					//		c.log.Error(err)
					//	}
					//}
				}
			}
		} else {
			c.log.Error(r.Err())
		}
		c.AsyncAsanaAllActivityLog(a)
	}

	xHookSecret := ctx.Request.Header.Get("X-Hook-Secret")
	//ctx.Request.Header.Set("X-Hook-Secret", xHookSecret)
	ctx.Header("X-Hook-Secret", xHookSecret)
	ctx.JSON(200, gin.H{
		"msg": "ok",
	})
}

func (c *HttpWebhookUsecase) AsyncAsanaAllActivityLog(rawData []byte) {
	go func() {
		typeMap := lib.ToTypeMapByString(string(rawData))
		typeList := lib.ToTypeList(typeMap["events"])
		var entities []AsanaAllActivityLogEntity
		for _, v := range typeList {
			entities = append(entities, AsanaAllActivityLogEntity{
				UserGid:                       v.GetString("user.gid"),
				UserResourceType:              v.GetString("user.resource_type"),
				AsanaCreatedAt:                v.GetString("created_at"),
				Action:                        v.GetString("action"),
				ResourceGid:                   v.GetString("resource.gid"),
				ResourceType:                  v.GetString("resource.resource_type"),
				ResourceSubtype:               v.GetString("resource.resource_subtype"),
				ParentGid:                     v.GetString("parent.gid"),
				ParentResourceType:            v.GetString("parent.resource_type"),
				ParentResourceSubtype:         v.GetString("parent.resource_subtype"),
				ChangeField:                   v.GetString("change.field"),
				ChangeAction:                  v.GetString("change.action"),
				ChangeNewValueGid:             v.GetString("change.new_value.gid"),
				ChangeNewValueResourceType:    v.GetString("change.new_value.resource_type"),
				ChangeNewValueResourceSubtype: v.GetString("change.new_value.resource_subtype"),
			})
		}
		if len(entities) > 0 {
			err := c.CommonUsecase.DB().Create(entities).Error
			if err != nil {
				c.log.Error(err)
			}
		}
	}()
}

func (c *HttpWebhookUsecase) ReceiveBox(ctx *gin.Context) {

	body, err := ctx.GetRawData()
	if err != nil {
		c.log.Error(err)
	} else {
		webhookLog := BoxWebhookLogEntity{
			Remarks:   ctx.Request.Method,
			Headers:   lib.InterfaceToString(ctx.Request.Header),
			Query:     lib.InterfaceToString(ctx.Request.URL.Query()),
			Body:      string(body),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		c.CommonUsecase.DB().Create(&webhookLog)
	}
	ctx.JSON(200, gin.H{
		"msg": "ok",
	})
}

// ReceiveAdobeSign 接收webhook
func (c *HttpWebhookUsecase) ReceiveAdobeSign(ctx *gin.Context) {

	a, err := ctx.GetRawData()
	if err != nil {
		c.log.Error(err)
	} else {
		r := lib.ToTypeMapByString(string(a))
		err = c.SaveAdobeWebhookEvent(r)
		if err != nil {
			c.log.Error(err)
		}
	}
	ctx.JSON(200, gin.H{
		"msg": "ok",
	})
}

// FormResponses 接收表单信息
func (c *HttpWebhookUsecase) FormResponses(ctx *gin.Context) {

	body, err := ctx.GetRawData()
	if err != nil {
		c.log.Error(err)
	} else {
		sourceForm := ctx.Query("source_form")

		entity := FormResponseEntity{
			SourceForm: sourceForm,
			Data:       string(body),
			CreatedAt:  time.Now().Unix(),
		}
		err := c.CommonUsecase.DB().Create(&entity).Error
		if err != nil {
			c.log.Error(err)
		}
	}
	ctx.JSON(200, gin.H{
		"msg": "ok",
	})
}

func (c *HttpWebhookUsecase) SaveAdobeWebhookEvent(row lib.TypeMap) error {

	entity := AdobeWebhookEventEntity{
		WebhookId:            row.GetString("webhookId"),
		WebhookName:          row.GetString("webhookName"),
		Event:                row.GetString("event"),
		EventDate:            row.GetString("eventDate"),
		EventResourceType:    row.GetString("eventResourceType"),
		ParticipantRole:      row.GetString("participantRole"),
		ParticipantUserId:    row.GetString("participantUserId"),
		ParticipantUserEmail: row.GetString("participantUserEmail"),
		ActingUserId:         row.GetString("actingUserId"),
		ActingUserEmail:      row.GetString("actingUserEmail"),
		InitiatingUserId:     row.GetString("initiatingUserId"),
		InitiatingUserEmail:  row.GetString("initiatingUserEmail"),
		AgreementId:          row.GetString("agreement.id"),
		AgreementName:        row.GetString("agreement.name"),
		AgreementStatus:      row.GetString("agreement.status"),
		CreatedAt:            time.Now().Unix(),
	}
	return c.CommonUsecase.DB().Create(&entity).Error
}
