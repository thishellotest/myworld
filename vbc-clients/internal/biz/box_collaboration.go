package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
)

/*
CREATE TABLE `box_collaborations` (

	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`case_id` int(11) NOT NULL DEFAULT '0',
	`box_folder_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`box_user_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`created_at` int(11) NOT NULL DEFAULT '0',
	`updated_at` int(11) NOT NULL DEFAULT '0',
	`deleted_at` int(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`),
	KEY `idx` (`box_folder_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='box_collaborations';
*/

const (
	Box_collaboration_ow         = "ow"         // Client Case Owner
	Box_collaboration_vs         = "vs"         // Lead VS
	Box_collaboration_cp         = "cp"         // Lead CP
	Box_collaboration_support_cp = "support_cp" // Support CP
	Box_collaboration_lead_co    = "co"         // Lead CO
)

type BoxCollaborationEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	CaseId             int32
	BoxFolderId        string
	PermissionSource   string
	UserGid            string
	BoxUserId          string
	BoxCollaborationId string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (BoxCollaborationEntity) TableName() string {
	return "box_collaborations"
}

type BoxCollaborationUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[BoxCollaborationEntity]
}

func NewBoxCollaborationUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *BoxCollaborationUsecase {
	uc := &BoxCollaborationUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}
