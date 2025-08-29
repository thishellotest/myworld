package biz

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type NotificationHttpUsecase struct {
	log                    *log.Helper
	conf                   *conf.Data
	JWTUsecase             *JWTUsecase
	NotificationUsecase    *NotificationUsecase
	NotificationbuzUsecase *NotificationbuzUsecase
}

func NewNotificationHttpUsecase(logger log.Logger,
	conf *conf.Data,
	JWTUsecase *JWTUsecase,
	NotificationUsecase *NotificationUsecase,
	NotificationbuzUsecase *NotificationbuzUsecase) *NotificationHttpUsecase {
	return &NotificationHttpUsecase{
		log:                    log.NewHelper(logger),
		conf:                   conf,
		JWTUsecase:             JWTUsecase,
		NotificationUsecase:    NotificationUsecase,
		NotificationbuzUsecase: NotificationbuzUsecase,
	}
}

func (c *NotificationHttpUsecase) List(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	rawData, _ := ctx.GetRawData()
	body := lib.ToTypeMapByString(string(rawData))
	notificationType := body.GetInt("notification_type")
	data, err := c.BizList(userFacade, body.GetInt("id"), notificationType)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *NotificationHttpUsecase) Info(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)
	data, err := c.BizInfo(userFacade)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *NotificationHttpUsecase) BizInfo(userFacade UserFacade) (lib.TypeMap, error) {

	var count int64
	err := c.NotificationUsecase.CommonUsecase.DB().
		Model(&NotificationEntity{}).
		Where("deleted_at=0 and receiver_gid=? and unread>0", userFacade.Gid()).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	data.Set("data.count", int32(count))
	return data, nil
}

func (c *NotificationHttpUsecase) Read(ctx *gin.Context) {
	reply := CreateReply()
	userFacade, _ := c.JWTUsecase.JWTUserFacade(ctx)

	gid := ctx.Param("gid")
	data, err := c.BizRead(userFacade, gid)
	if err != nil {
		reply.CommonError(err)
	} else {
		reply.Merge(data)
		reply.Success()
	}
	ctx.JSON(200, reply)
}

func (c *NotificationHttpUsecase) BizRead(userFacade UserFacade, gid string) (lib.TypeMap, error) {
	notif, err := c.NotificationUsecase.GetByCond(Eq{
		"deleted_at":   0,
		"gid":          gid,
		"receiver_gid": userFacade.Gid(),
	})
	if err != nil {
		return nil, err
	}
	if notif == nil {
		return nil, errors.New("Parameter incorrect")
	}
	notif.Unread = 0
	notif.UpdatedAt = time.Now().Unix()
	err = c.NotificationUsecase.CommonUsecase.DB().Save(&notif).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type NotificationHttpListResponse struct {
	Records NotificationList `json:"records"`
}

type NotificationList []NotificationItem
type NotificationItem struct {
	Id               int32   `json:"id"`
	Gid              string  `json:"gid"`
	Title            string  `json:"title"`        // 例：Notes
	TriggerUser      FabUser `json:"trigger_user"` // 发送人
	NotificationTime int32   `json:"notification_time"`
	Content          string  `json:"content"` // 内容
	Unread           int     `json:"unread"`
	Url              string  `json:"url"`
	OpenNewWindow    bool    `json:"open_new_window"`
}

func (c *NotificationHttpUsecase) BizList(userFacade UserFacade, lastId int32, notificationType int32) (lib.TypeMap, error) {

	pageSize := 10
	records, err := c.NotificationUsecase.LatestRecords(userFacade.Gid(), notificationType, lastId, pageSize)
	if err != nil {
		return nil, err
	}

	items, err := c.NotificationbuzUsecase.Tidy(records)
	if err != nil {
		return nil, err
	}
	data := make(lib.TypeMap)
	var notificationHttpListResponse NotificationHttpListResponse
	notificationHttpListResponse.Records = items
	data.Set("notifications", notificationHttpListResponse)
	data.Set("page_size", pageSize)
	return data, nil
}
