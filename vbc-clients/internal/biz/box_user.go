package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
CREATE TABLE `box_users` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`box_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`user_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`login` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`box_created_at` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`box_modified_at` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`language` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`timezone` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`space_amount` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`space_used` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`max_upload_size` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`job_title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`phone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`address` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`avatar_url` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`notification_email` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='box_users';
*/
type BoxUserEntity struct {
	ID                int32 `gorm:"primaryKey"`
	BoxUserId         string
	UserType          string
	Name              string
	Login             string
	BoxCreatedAt      string
	BoxModifiedAt     string
	Language          string
	Timezone          string
	SpaceAmount       string
	SpaceUsed         string
	MaxUploadSize     string
	JobTitle          string
	Phone             string
	Status            string
	Address           string
	AvatarUrl         string
	NotificationEmail string
	CreatedAt         int64
	UpdatedAt         int64
}

func (BoxUserEntity) TableName() string {
	return "box_users"
}

type BoxUserUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[BoxUserEntity]
}

func NewBoxUserUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *BoxUserUsecase {
	uc := &BoxUserUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

/*
Upsert

	{
			"type": "user",
			"id": "40712863180",
			"name": "Debra Harris",
			"login": "dharris@vetbenefitscenter.com",
			"created_at": "2025-03-10T09:16:17-07:00",
			"modified_at": "2025-06-21T07:36:05-07:00",
			"language": "en",
			"timezone": "America\/Los_Angeles",
			"space_amount": 10737418240,
			"space_used": 0,
			"max_upload_size": 53687091200,
			"status": "active",
			"job_title": "",
			"phone": "",
			"address": "",
			"avatar_url": "https:\/\/veteranbenefitscenter.app.box.com\/api\/avatar\/large\/40712863180",
			"notification_email": null
		}
*/
func (c *BoxUserUsecase) Upsert(data lib.TypeMap) error {

	boxUserId := data.GetString("id")
	if boxUserId == "" {
		return errors.New("boxUserId is empty")
	}
	entity, err := c.GetByCond(Eq{"box_user_id": boxUserId})
	if err != nil {
		return err
	}
	if entity == nil {
		entity = &BoxUserEntity{
			BoxUserId: boxUserId,
			CreatedAt: time.Now().Unix(),
		}
	}
	entity.UserType = data.GetString("type")
	entity.Name = data.GetString("name")
	entity.Login = data.GetString("login")
	entity.BoxCreatedAt = data.GetString("created_at")
	entity.BoxModifiedAt = data.GetString("modified_at")
	entity.Language = data.GetString("language")
	entity.Timezone = data.GetString("timezone")
	entity.SpaceAmount = data.GetString("space_amount")
	entity.SpaceUsed = data.GetString("space_used")
	entity.MaxUploadSize = data.GetString("max_upload_size")
	entity.Status = data.GetString("status")
	entity.JobTitle = data.GetString("job_title")
	entity.Phone = data.GetString("phone")
	entity.Address = data.GetString("address")
	entity.AvatarUrl = data.GetString("avatar_url")
	entity.NotificationEmail = data.GetString("notification_email")
	return c.CommonUsecase.DB().Save(&entity).Error
}
