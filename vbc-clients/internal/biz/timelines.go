package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"vbc/internal/conf"
	"vbc/lib/uuid"
)

const (
	Timeline_FieldName_kind             = "kind"
	Timeline_FieldName_kind_gid         = "kind_gid"
	Timeline_FieldName_related_kind     = "related_kind"
	Timeline_FieldName_related_kind_gid = "related_kind_gid"
	Timeline_FieldName_action           = "action"
	Timeline_FieldName_notes            = "notes"
)

type TimelineForNotes struct {
	Content string `json:"content"`
}

type TimelineFieldHistoryNotes struct {
	FieldHistory TimelineFieldHistory `json:"field_history"`
}

func (c *TimelineFieldHistoryNotes) ToString() string {
	return InterfaceToString(c)
}

type TimelineFieldHistory []TimelineFieldHistoryItem
type TimelineFieldHistoryItem struct {
	FieldName string `json:"field_name"`
	NewValue  string `json:"new_value"`
	OldValue  string `json:"old_value"`
}

const (
	Timeline_action_added    = "added"
	Timeline_action_updated  = "updated"
	Timeline_action_deleted  = "deleted"
	Timeline_action_restored = "restored"
)

type TimelineUsecase struct {
	log              *log.Helper
	CommonUsecase    *CommonUsecase
	conf             *conf.Data
	DataEntryUsecase *DataEntryUsecase
}

func NewTimelineUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	DataEntryUsecase *DataEntryUsecase) *TimelineUsecase {
	uc := &TimelineUsecase{
		log:              log.NewHelper(logger),
		CommonUsecase:    CommonUsecase,
		conf:             conf,
		DataEntryUsecase: DataEntryUsecase,
	}

	return uc
}

/*
CREATE TABLE `timelines` (
	`id` int unsigned NOT NULL AUTO_INCREMENT,
	`gid` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'timelines gid',
	`kind` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '属于那个module',
	`kind_gid` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '属于谁的',
	`action` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
	`related_kind` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '关联那个module',
	`related_kind_gid` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '关联那个',
	`notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'json',
	`created_by` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'users gid',
	`modified_by` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'users gid',
	`created_at` int NOT NULL DEFAULT '0',
	`updated_at` int NOT NULL DEFAULT '0',
	`deleted_at` int NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='timelines';
*/

// Create 创建
func (c *TimelineUsecase) Create(Kind string, kindGid string, action string, relatedKind string, relatedKindGid string,
	notes string, operUser *TData) (gid string, err error) {

	var createdBy, modifiedBy string
	if operUser != nil {
		createdBy = operUser.Gid()
		modifiedBy = createdBy
	}

	gid = uuid.UuidWithoutStrike()
	data := make(TypeDataEntry)
	data["kind"] = Kind
	data["kind_gid"] = kindGid
	data["action"] = action
	data["related_kind"] = relatedKind
	data["related_kind_gid"] = relatedKindGid
	data["notes"] = notes
	data["created_by"] = createdBy
	data["modified_by"] = modifiedBy
	data["gid"] = gid

	_, err = c.DataEntryUsecase.HandleOne(Kind_timelines, data, "gid", operUser)
	if err != nil {
		return "", err
	}
	return gid, nil
}
