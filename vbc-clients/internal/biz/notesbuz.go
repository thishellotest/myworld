package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
	"vbc/lib/uuid"
)

type NotesbuzUsecase struct {
	log                 *log.Helper
	conf                *conf.Data
	CommonUsecase       *CommonUsecase
	EventBus            *EventBus
	NotificationUsecase *NotificationUsecase
	UserUsecase         *UserUsecase
	TUsecase            *TUsecase
	ClientCaseUsecase   *ClientCaseUsecase
}

func NewNotesbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	EventBus *EventBus,
	NotificationUsecase *NotificationUsecase,
	UserUsecase *UserUsecase,
	TUsecase *TUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
) *NotesbuzUsecase {
	uc := &NotesbuzUsecase{
		log:                 log.NewHelper(logger),
		CommonUsecase:       CommonUsecase,
		conf:                conf,
		EventBus:            EventBus,
		NotificationUsecase: NotificationUsecase,
		UserUsecase:         UserUsecase,
		TUsecase:            TUsecase,
		ClientCaseUsecase:   ClientCaseUsecase,
	}

	uc.EventBus.Subscribe(EventBus_AfterInsertData, uc.HandleAfterInsertData)
	uc.EventBus.Subscribe(EventBus_AfterHandleUpdate, uc.HandleAfterHandleUpdate)

	return uc
}

func (c *NotesbuzUsecase) HandleAfterInsertData(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, modifiedBy string) {
	if kindEntity.Kind == Kind_notes && recognizeFieldName == DataEntry_gid {
		for _, v := range sourceData {
			gid := InterfaceToString(v[DataEntry_gid])
			content := InterfaceToString(v[Notes_FieldName_content])
			if gid != "" && content != "" {
				c.HandleNoteNotification(gid, content)
			}
		}
	}
}

func (c *NotesbuzUsecase) HandleAfterHandleUpdate(kindEntity KindEntity,
	structField *TypeFieldStruct,
	recognizeFieldName string,
	dataEntryOperResult DataEntryOperResult,
	sourceData TypeDataEntryList, modifiedBy string) {

	if kindEntity.Kind == Kind_notes && recognizeFieldName == DataEntry_gid {
		destData := make(map[string]string)
		for gid, v := range dataEntryOperResult {
			if v.IsUpdated {
				for k1, v1 := range v.DataEntryModifyDataMap {
					if k1 == Notes_FieldName_content {
						newVal := v1.GetNewVal(FieldType_text)
						destData[gid] = newVal
						break
					}
				}
			}
		}
		for k, v := range destData {
			err := c.HandleNoteNotification(k, v)
			if err != nil {
				c.log.Error(err)
			}
		}
	}
}

func (c *NotesbuzUsecase) HandleNoteNotification(noteGid string, content string) error {

	tNote, _ := c.TUsecase.DataByGid(Kind_notes, noteGid)
	if tNote == nil {
		c.log.Error("tNote is nil")
		return nil
	}

	result := NotificationTextExtractContext(content, 20, 30)
	for _, v := range result {

		user, _ := c.UserUsecase.GetByGid(v.UserGid)
		if user != nil {
			modifiedBy := tNote.CustomFields.TextValueByNameBasic(DataEntry_modified_by)

			entity := NotificationEntity{
				Gid:         uuid.UuidWithoutStrike(),
				FromType:    Notification_FromType_Notes,
				FromGid:     noteGid,
				ReceiverGid: user.Gid(),
				SenderGid:   modifiedBy,
				Content:     v.ToText(),
				Unread:      1,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			}
			err := c.CommonUsecase.DB().Save(&entity).Error
			if err != nil {
				c.log.Error(err, InterfaceToString(entity))
			}
		}
	}
	return nil
}

func (c *NotesbuzUsecase) HandlePWNotification(caseGid string) error {

	tCase, _ := c.TUsecase.DataByGid(Kind_client_cases, caseGid)
	if tCase == nil {
		c.log.Error("tCase is nil")
		return nil
	}
	leadVs := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
	leadCp := tCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp)

	leadsVsUser, _ := c.UserUsecase.GetByFullName(leadVs)
	leadCpUser, _ := c.UserUsecase.GetByFullName(leadCp)

	userGids := make(map[string]bool)
	if leadsVsUser != nil {
		userGids[leadsVsUser.Gid()] = true
	}
	if leadCpUser != nil {
		userGids[leadCpUser.Gid()] = true
	}
	userGid := tCase.CustomFields.TextValueByNameBasic(FieldName_user_gid)
	if userGid != "" {
		userGids[userGid] = true
	}
	if len(userGids) == 0 {
		c.log.Error("No user who needs to be notified was found")
		return nil
	}
	content := fmt.Sprintf("\"%s\" has submitted Statement Feedback. Please review it.",
		tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name))

	userGidArr := lib.MapKeys(userGids)
	users, err := c.TUsecase.ListByCond(Kind_users, And(In("gid", userGidArr), Eq{"deleted_at": 0}))
	if err != nil {
		c.log.Error(err)
		return nil
	}

	for _, v := range users {

		entity := NotificationEntity{
			Gid:         uuid.UuidWithoutStrike(),
			FromType:    Notification_FromType_PW,
			FromGid:     caseGid,
			ReceiverGid: v.Gid(),
			SenderGid:   "",
			Content:     content,
			Unread:      1,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}
		err := c.CommonUsecase.DB().Save(&entity).Error
		if err != nil {
			c.log.Error(err, InterfaceToString(entity))
		}

	}
	return nil
}
