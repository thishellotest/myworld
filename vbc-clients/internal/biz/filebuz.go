package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type FilebuzUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	TUsecase      *TUsecase
	FileUsecase   *FileUsecase
}

func NewFilebuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	TUsecase *TUsecase,
	FileUsecase *FileUsecase,
) *FilebuzUsecase {
	uc := &FilebuzUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
		TUsecase:      TUsecase,
		FileUsecase:   FileUsecase,
	}

	return uc
}

func (c *FilebuzUsecase) DCRecordReviewFileHandle(webhookBodyMap lib.TypeMap) error {

	trigger := webhookBodyMap.GetString("trigger")

	if trigger == "FILE.TRASHED" {
		fileEntity, err := c.FileUsecase.GetByCond(Eq{"source_id": webhookBodyMap.GetString("source.id"), "type": "file"})
		if err != nil {
			return err
		}
		if fileEntity != nil {
			fileEntity.BizDeletedAt = time.Now().Unix()
			return c.CommonUsecase.DB().Save(&fileEntity).Error
		}
		return nil
	} else if trigger == "FOLDER.TRASHED" {
		entity, err := c.FileUsecase.GetByCond(Eq{"source_id": webhookBodyMap.GetString("source.id"), "type": "folder"})
		if err != nil {
			return err
		}
		if entity != nil {
			entity.BizDeletedAt = time.Now().Unix()
			return c.CommonUsecase.DB().Save(&entity).Error
		}
		return nil
	} else if trigger == "FOLDER.RESTORED" {

		entity, err := c.FileUsecase.GetByCond(Eq{"source_id": webhookBodyMap.GetString("source.id"), "type": "folder"})
		if err != nil {
			return err
		}
		if entity != nil {
			entity.BizDeletedAt = 0
			err = c.CommonUsecase.DB().Save(&entity).Error
			if err != nil {
				return err
			}
			entries := webhookBodyMap.GetTypeList("source.path_collection.entries")
			return c.FileUsecase.HandleEntriesFolder(entries, entity.CaseId, "0")
		}
		return nil
	}

	entries := webhookBodyMap.GetTypeList("source.path_collection.entries")
	if len(entries) < 5 {
		return nil
	}
	if entries[2].GetString("id") == configs.DCFolderId {
		recordReviewFolderName := entries[4].GetString("name")
		if recordReviewFolderName != "Record Review" {
			return nil
		}
		dcFolderName := entries[3].GetString("name")
		dcFolderNameArr := strings.Split(dcFolderName, "#")
		if len(dcFolderNameArr) != 2 {
			return errors.New("dcFolderNameArr is wrong : " + dcFolderName)
		}

		caseId, err := strconv.ParseInt(strings.TrimSpace(dcFolderNameArr[1]), 10, 32)
		if err != nil {
			return err
		}
		tCase, err := c.TUsecase.DataById(Kind_client_cases, int32(caseId))
		if err != nil {
			return err
		}
		if tCase == nil {
			return errors.New("tCase is nil : " + dcFolderName)
		}
		return c.HandleFile(tCase.Id(), entries[4].GetString("id"), webhookBodyMap)
	}
	return nil
}

func (c *FilebuzUsecase) HandleFile(caseId int32, destFolderId string, webHookMap lib.TypeMap) error {

	trigger := webHookMap.GetString("trigger")
	if trigger == "FILE.RENAMED" || trigger == "FILE.UPLOADED" || trigger == "FILE.RESTORED" {

		return c.HandleSourceFile(caseId, destFolderId, webHookMap.GetTypeMap("source"))
	}
	return nil
}

func (c *FilebuzUsecase) HandleSourceFile(caseId int32, destFolderId string, sourceMap lib.TypeMap) error {

	resType := sourceMap.GetString("type")
	resId := sourceMap.GetString("id")
	if resType == "file" {
		name := sourceMap.GetString("name")
		fileVersionId := sourceMap.GetString("source.file_version.id")
		entity, err := c.FileUsecase.GetByCond(Eq{"source_id": resId, "type": "file"})
		if err != nil {
			return err
		}
		if entity == nil {
			entity = &FileEntity{
				CaseId:       caseId,
				DestFolderId: destFolderId,
				Type:         "file",
				SourceId:     resId,
				SourceName:   name,

				CreatedAt: time.Now().Unix(),
			}
		}
		entity.FileVersionId = fileVersionId
		entity.ParentFolderId = sourceMap.GetString("parent.id")
		entity.BizDeletedAt = 0

		/*
			"created_at": "2025-05-11T17:58:26-07:00",
				"modified_at": "2025-05-11T18:00:55-07:00",
				"trashed_at": null,
				"purged_at": null,
				"content_created_at": "2025-05-09T03:42:23-07:00",
				"content_modified_at": "2025-05-09T03:42:23-07:00",
		*/
		if sourceMap.GetString("created_at") != "" {
			entity.SourceCreatedAt = sourceMap.GetString("created_at")
		}
		if sourceMap.GetString("modified_at") != "" {
			entity.SourceModifiedAt = sourceMap.GetString("modified_at")
		}
		if sourceMap.GetString("trashed_at") != "" {
			entity.SourceTrashedAt = sourceMap.GetString("trashed_at")
		}
		if sourceMap.GetString("purged_at") != "" {
			entity.SourcePurgedAt = sourceMap.GetString("purged_at")
		}
		if sourceMap.GetString("content_created_at") != "" {
			entity.ContentCreatedAt = sourceMap.GetString("content_created_at")
		}
		if sourceMap.GetString("content_modified_at") != "" {
			entity.ContentModifiedAt = sourceMap.GetString("content_modified_at")
		}
		err = c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			return err
		}
		entries := sourceMap.GetTypeList("path_collection.entries")
		err = c.FileUsecase.HandleEntriesFolder(entries, caseId, "0")
		if err != nil {
			return err
		}
	}
	return nil
}
