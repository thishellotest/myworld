package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"vbc/internal/conf"
)

const (
	Notification_FromType_Notes = 1
	Notification_FromType_PW    = 2
)

type NotificationEntity struct {
	ID          int32 `gorm:"primaryKey"`
	Gid         string
	FromType    int
	FromGid     string
	ReceiverGid string
	SenderGid   string
	Content     string
	Unread      int
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
}

func (NotificationEntity) TableName() string {
	return "notifications"
}

type NotificationUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[NotificationEntity]
}

func NewNotificationUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *NotificationUsecase {
	uc := &NotificationUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

func (c *NotificationUsecase) LatestRecords(userGid string, notificationType int32, lastId int32, pageSize int) ([]NotificationEntity, error) {
	var records []NotificationEntity
	var query *gorm.DB
	str := ""
	if notificationType == 1 {
		str = "and unread!=0"
	}
	if lastId == 0 {
		query = c.CommonUsecase.DB().Where("deleted_at=0 and receiver_gid=? "+str, userGid)
	} else {
		query = c.CommonUsecase.DB().Where("deleted_at=0 and receiver_gid=? and id <? "+str, userGid, lastId)
	}
	err := query.Order("id desc").
		Limit(pageSize).
		Find(&records).Error
	return records, err
}
