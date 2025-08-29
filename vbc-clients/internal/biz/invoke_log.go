package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

/*
CREATE TABLE `invoke_log` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`invoke_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`record_uniq_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`record_modified_time` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`record_created_time` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	`content` mediumtext COLLATE utf8mb4_unicode_ci,
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	`deleted_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	KEY `idx` (`record_uniq_id`(191))

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='invoke_log';
*/

type InvokeLogEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	RecordUniqId       string
	RecordModifiedTime string
	RecordCreatedTime  string
	Content            string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
	HandleStatus       int
	HandleResult       int
	NextAt             int64
	HandleResultDetail string
}

func (c *InvokeLogEntity) GetContentMap() lib.TypeMap {
	return lib.ToTypeMapByString(c.Content)
}

func (c *InvokeLogEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

func (InvokeLogEntity) TableName() string {
	return "invoke_log"
}

type InvokeLogUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[InvokeLogEntity]
}

func NewInvokeLogUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *InvokeLogUsecase {
	uc := &InvokeLogUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *InvokeLogUsecase) Upsert(RecordUniqId string, RecordModifiedTime string, RecordCreatedTime string, content lib.TypeMap) (*InvokeLogEntity, error) {
	var entity *InvokeLogEntity
	isNew := false
	needUpdate := false
	var err error
	if RecordUniqId == "" {
		isNew = true
	} else {
		entity, err = c.GetByCond(Eq{"record_uniq_id": RecordUniqId})
		if err != nil {
			return nil, err
		}
		if entity == nil {
			isNew = true
		} else {
			entity.UpdatedAt = time.Now().Unix()
			newContent := InterfaceToString(content)
			if entity.Content != newContent {
				entity.Content = newContent
				needUpdate = true
				entity.RecordModifiedTime = RecordModifiedTime
				entity.RecordCreatedTime = RecordCreatedTime
			}
		}
	}
	if isNew {
		entity = &InvokeLogEntity{
			RecordUniqId:       RecordUniqId,
			RecordModifiedTime: RecordModifiedTime,
			RecordCreatedTime:  RecordCreatedTime,
			Content:            InterfaceToString(content),
			CreatedAt:          time.Now().Unix(),
			UpdatedAt:          time.Now().Unix(),
		}
	}
	if isNew || needUpdate {
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return nil, err
		}
	}
	return entity, nil
}
