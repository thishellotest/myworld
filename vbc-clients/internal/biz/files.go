package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"strings"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

//
//CREATE TABLE `files` (
//`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
//`case_id` int(11) NOT NULL DEFAULT '0',
//`dest_folder_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '目标文件夹（如：DC/Record Reciew）',
//`type` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'folder/file',
//`source_id` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件唯一ID',
//`source_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`source_created_at` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`source_updated_at` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`source_trashed_at` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`source_purged_at` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
//`created_at` int(11) NOT NULL DEFAULT '0',
//`updated_at` int(11) NOT NULL DEFAULT '0',
//`biz_deleted_at` int(11) NOT NULL DEFAULT '0',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='files';

type FileEntity struct {
	ID                int32 `gorm:"primaryKey"`
	CaseId            int32
	ParentFolderId    string
	DestFolderId      string
	Type              string
	SourceId          string
	SourceName        string
	FileVersionId     string
	SourceCreatedAt   string
	SourceModifiedAt  string
	SourceTrashedAt   string
	SourcePurgedAt    string
	ContentCreatedAt  string
	ContentModifiedAt string
	CreatedAt         int64
	UpdatedAt         int64
	BizDeletedAt      int64
}

func (c *FileEntity) CanReview() bool {
	n := 4
	if len(c.SourceName) < n {
		n = len(c.SourceName)
	}
	last := c.SourceName[len(c.SourceName)-n:]
	if strings.ToLower(last) == ".pdf" {
		return true
	}
	return false
}

func (FileEntity) TableName() string {
	return "files"
}

type FileUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[FileEntity]
}

func NewFileUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *FileUsecase {
	uc := &FileUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *FileUsecase) HandleEntriesFolder(entries lib.TypeList, caseId int32, destFolderId string) error {

	entriesLen := len(entries)

	isOk := false
	var prevFolderEntity *FileEntity
	for i := 0; i < entriesLen; i++ {
		id := entries[i].GetString("id")
		name := entries[i].GetString("name")
		if id == destFolderId {
			isOk = true
		}
		if !isOk {
			continue
		}
		entity, err := c.GetByCond(Eq{"source_id": id, "type": "folder"})
		if err != nil {
			return err
		}
		if entity == nil {
			entity = &FileEntity{
				CaseId: caseId,
				//DestFolderId: destFolderId,
				Type:     "folder",
				SourceId: id,
				//SourceName: name,
				CreatedAt: time.Now().Unix(),
			}

		}
		if prevFolderEntity != nil {
			entity.ParentFolderId = prevFolderEntity.SourceId
		} else {
			entity.ParentFolderId = ""
		}
		entity.SourceName = name
		entity.DestFolderId = destFolderId

		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return err
		}
		prevFolderEntity = entity
	}
	return nil
}
