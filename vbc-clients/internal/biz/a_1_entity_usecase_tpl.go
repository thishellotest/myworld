package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
)

/*


CREATE TABLE `profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `gid` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'gid',
  `created_by` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'users gid',
  `modified_by` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'users gid',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Created At',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Updated At',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT 'Deleted At',
  PRIMARY KEY (`id`),
  KEY `idx_e` (`gid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='profiles';



INSERT INTO `fields` ( `kind`, `field_name`, `field_label`, `field_type`, `rela_kind`, `rela_name`, `is_condition`, `deleted_at`)
VALUES
	( 'profiles', 'is_retain', 'is_retain', 'number', '', '', 0, 0);

INSERT INTO `fields` ( `kind`, `field_name`, `field_label`, `field_type`, `rela_kind`, `rela_name`, `is_condition`, `deleted_at`)
VALUES
	( 'profiles', 'is_admin', 'is_admin', 'number', '', '', 0, 0);


INSERT INTO `kinds` ( `kind`, `tablename`, `label`, `tab_label`, `primary_field_name`, `no_change_history`, `no_timelines`, `deleted_at`)
VALUES
	( 'profiles', 'profiles', 'Profile', 'Profiles', '', 0, 1, 0);


*/

type TTemplateEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (TTemplateEntity) TableName() string {
	return "test_table_name"
}

func (c *TTemplateEntity) AppendHandleResultDetail(str string) {
	if c.HandleResultDetail == "" {
		c.HandleResultDetail = time.Now().Format(time.RFC3339) + " " + str
	} else {
		c.HandleResultDetail += "\r\n" + time.Now().Format(time.RFC3339) + " " + str
	}
}

type TTemplateUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[TTemplateEntity]
}

func NewTTemplateUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *TTemplateUsecase {
	uc := &TTemplateUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
