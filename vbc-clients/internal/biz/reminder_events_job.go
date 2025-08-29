package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type ReminderEventsJobUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[TTemplateEntity]
	ReminderEventUsecase *ReminderEventUsecase
	UserUsecase          *UserUsecase
	RemindUsecase        *RemindUsecase
	MailUsecase          *MailUsecase
	TUsecase             *TUsecase
}

func NewReminderEventsJobUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data,
	ReminderEventUsecase *ReminderEventUsecase,
	UserUsecase *UserUsecase,
	RemindUsecase *RemindUsecase,
	MailUsecase *MailUsecase,
	TUsecase *TUsecase) *ReminderEventsJobUsecase {
	uc := &ReminderEventsJobUsecase{
		log:                  log.NewHelper(logger),
		CommonUsecase:        CommonUsecase,
		conf:                 conf,
		ReminderEventUsecase: ReminderEventUsecase,
		UserUsecase:          UserUsecase,
		RemindUsecase:        RemindUsecase,
		MailUsecase:          MailUsecase,
		TUsecase:             TUsecase,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

func (c *ReminderEventsJobUsecase) RunHandleJob(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ReminderEventsJobUsecase RunHandleJob Done")
				return
			default:
				err := c.Handle()
				if err != nil {
					c.log.Error(err)
				}
				time.Sleep(60 * time.Second)
			}
		}
	}()
	return nil
}

func (c *ReminderEventsJobUsecase) Handle() error {
	events, err := c.ReminderEventUsecase.AllByCond(And(Eq{"deleted_at": 0, "handle_status": 0}))
	if err != nil {
		return err
	}
	reminderEventGroup := make(ReminderEventGroup)
	for k, _ := range events {
		reminderEventGroup.Append(events[k])
	}

	for k, _ := range reminderEventGroup {
		reminders := reminderEventGroup[k]
		if len(reminders) > 0 {
			firstEntity := reminders[0]
			if firstEntity.EventType == ReminderEventType_ClientUpdateFiles {

				tCase, err := c.TUsecase.DataById(Kind_client_cases, firstEntity.IncrId)
				if err != nil {
					c.log.Error(err)
					continue
				}
				err = c.ReminderClientUpdateFiles(tCase, reminders)
				if err != nil {
					c.log.Error(err, " ", firstEntity.ID, " ", firstEntity.IncrId)
				}
			} else {
				c.log.Error(firstEntity.EventType + ": firstEntity.EventType does not support")
			}

			err = c.FinishReminder(reminders)
			if err != nil {
				c.log.Error("FinishReminder err:", err)
			}
		}
	}
	return nil
}

func (c *ReminderEventsJobUsecase) FinishReminder(events []*ReminderEventEntity) error {

	var ids []int32
	for _, v := range events {
		ids = append(ids, v.ID)
	}
	if len(ids) <= 0 {
		return nil
	}

	return c.CommonUsecase.DB().Model(&ReminderEventEntity{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"handle_status": 1,
			"updated_at":    time.Now().Unix()}).Error
}

func (c *ReminderEventsJobUsecase) ReminderClientUpdateFiles(tClientCase *TData, events []*ReminderEventEntity) error {

	var items []*ReminderClientUpdateFilesEventVoItem
	for _, v := range events {
		eventVo := v.GetReminderClientUpdateFilesEventVo()
		if eventVo == nil {
			return errors.New("eventVo is nil")
		}
		items = append(items, eventVo.Items...)
	}

	if len(items) < 0 {
		return errors.New("errors is nil")
	}

	if tClientCase == nil {
		return errors.New("tClientCase is nil")
	}
	primaryCPFullName := ""
	primaryVSFullName := ""

	tClientCaseId := tClientCase.CustomFields.NumberValueByNameBasic("id")

	if tClientCase.CustomFields.TextValueByNameBasic("email") == "lialing@foxmail.com" ||
		tClientCase.CustomFields.TextValueByNameBasic("email") == "liaogling@gmail.com" {
		primaryCPFullName = "Engineering Team"
		primaryVSFullName = "Engineering Team"
	} else {
		primaryCPFullName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_primary_cp)
		primaryVSFullName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_primary_vs)
	}

	if primaryCPFullName == "" {
		return errors.New("primaryCPFullName is empty")
	}
	if primaryVSFullName == "" {
		return errors.New("primaryVSFullName is empty")
	}
	email := ""
	if primaryCPFullName == primaryVSFullName {
		tUser, err := c.UserUsecase.GetByFullName(primaryCPFullName)
		if err != nil {
			return err
		}
		if tUser == nil {
			return errors.New("tUser is nil")
		}
		email = tUser.CustomFields.TextValueByNameBasic("email")

	} else {
		tUserCP, err := c.UserUsecase.GetByFullName(primaryCPFullName)
		if err != nil {
			return err
		}
		tUserVS, err := c.UserUsecase.GetByFullName(primaryVSFullName)
		if err != nil {
			return err
		}
		if tUserCP == nil {
			return errors.New("tUserCP is nil")
		}
		if tUserVS == nil {
			return errors.New("tUserVS is nil")
		}
		cpEmail := tUserCP.CustomFields.TextValueByNameBasic("email")
		if cpEmail == "" {
			return errors.New("cpEmail is empty")
		}
		vsEmail := tUserVS.CustomFields.TextValueByNameBasic("email")
		if vsEmail == "" {
			return errors.New("vsEmail is empty")
		}
		email = fmt.Sprintf("%s;%s", cpEmail, vsEmail)
	}

	vo, err := c.RemindUsecase.FollowingUpUploadedDocumentEmailBody(tClientCase, email, items)

	if err != nil {
		return err
	}

	mailServiceConfig := InitMailServiceConfig()
	mailMessage := &MailMessage{
		To:      vo.Email,
		Subject: vo.Subject,
		Body:    vo.Body,
	}
	if !configs.IsDev() {
		if tClientCase.CustomFields.TextValueByNameBasic("email") != "lialing@foxmail.com" {
			mailMessage.Cc = []string{"info@vetbenefitscenter.com"}
		}
	}
	err = c.MailUsecase.SendEmail(mailServiceConfig, mailMessage, "", nil)

	lib.DPrintln("ReminderClientUpdateFiles mailServiceConfig:", mailServiceConfig)
	lib.DPrintln("ReminderClientUpdateFiles mailMessage:", mailMessage)
	c.CommonUsecase.DB().Save(&EmailLogEntity{
		ClientId:   tClientCaseId,
		Email:      vo.Email,
		TaskId:     0,
		Tpl:        "FollowingUpUploadedDocumentEmailBody",
		SubId:      0,
		SenderMail: mailServiceConfig.Username,
		SenderName: mailServiceConfig.Name,
		Subject:    mailMessage.Subject,
		Body:       mailMessage.Body,
	})

	return nil
}

type ReminderEventGroup map[string][]*ReminderEventEntity

func (c ReminderEventGroup) Append(reminder *ReminderEventEntity) {
	if c == nil {
		return
	}
	if reminder != nil {
		groupId := reminder.GroupId()
		if _, ok := c[groupId]; !ok {
			c[groupId] = make([]*ReminderEventEntity, 0)
		}
		c[groupId] = append(c[groupId], reminder)
	}
}
