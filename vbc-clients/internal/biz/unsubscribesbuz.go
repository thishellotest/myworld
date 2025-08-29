package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"strings"
	"time"
	"vbc/configs"
	"vbc/internal/conf"
	"vbc/lib"
	. "vbc/lib/builder"
)

type UnsubscribesbuzUsecase struct {
	log                       *log.Helper
	conf                      *conf.Data
	CommonUsecase             *CommonUsecase
	UnsubscribesUsecase       *UnsubscribesUsecase
	DialpadbuzInternalUsecase *DialpadbuzInternalUsecase
	ClientUsecase             *ClientUsecase
	DialpadbuzUsecase         *DialpadbuzUsecase
	ClientCaseUsecase         *ClientCaseUsecase
	LogUsecase                *LogUsecase
	DialpadUsecase            *DialpadUsecase
}

func NewUnsubscribesbuzUsecase(logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
	UnsubscribesUsecase *UnsubscribesUsecase,
	DialpadbuzInternalUsecase *DialpadbuzInternalUsecase,
	ClientUsecase *ClientUsecase,
	DialpadbuzUsecase *DialpadbuzUsecase,
	ClientCaseUsecase *ClientCaseUsecase,
	LogUsecase *LogUsecase,
	DialpadUsecase *DialpadUsecase,
) *UnsubscribesbuzUsecase {
	uc := &UnsubscribesbuzUsecase{
		log:                       log.NewHelper(logger),
		CommonUsecase:             CommonUsecase,
		conf:                      conf,
		UnsubscribesUsecase:       UnsubscribesUsecase,
		DialpadbuzInternalUsecase: DialpadbuzInternalUsecase,
		ClientUsecase:             ClientUsecase,
		DialpadbuzUsecase:         DialpadbuzUsecase,
		ClientCaseUsecase:         ClientCaseUsecase,
		LogUsecase:                LogUsecase,
		DialpadUsecase:            DialpadUsecase,
	}

	return uc
}

const (
	SMS_Text_action_stop   = "stop"   // 停止接收短信
	SMS_Text_action_unstop = "unstop" // 允许接收短信
	SMS_Text_action_help   = "help"   // 需要帮助，提醒vs，联系客户
)

func GetSMSTextAction(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	if strings.Index(text, SMS_Text_action_stop) == 0 {
		return SMS_Text_action_stop
	} else if strings.Index(text, SMS_Text_action_unstop) == 0 {
		return SMS_Text_action_unstop
	} else if strings.Index(text, SMS_Text_action_help) == 0 {
		return SMS_Text_action_help
	}
	return ""
}

// HandleFromDialpadWebhookEvent {"contact":{"id":6230136767135744,"name":"Test - VBC","phone_number":"+16192000000"},"created_date":1738808741123,"direction":"inbound","event_timestamp":1738808741527,"from_number":"+16192788886","id":5594630430703616,"is_internal":false,"message_delivery_result":null,"message_status":"pending","mms":false,"mms_url":null,"sender_id":null,"target":{"id":4693435745845248,"name":"Edward Bunting Jr.","phone_number":"(619) 800-0000","type":"user"},"text":"sTOp","text_content":"sTOp","to_number":["+16198005543"]}
func (c *UnsubscribesbuzUsecase) HandleFromDialpadWebhookEvent(data lib.TypeMap) error {

	if data == nil {
		return errors.New("data is nil")
	}
	contactPhoneNumber := data.GetString("contact.phone_number")
	fromId := data.GetString("id")
	if data.GetString("direction") == "inbound" && contactPhoneNumber != "" && fromId != "" {
		text := data.GetString("text")
		action := GetSMSTextAction(text)

		if action == SMS_Text_action_unstop || action == SMS_Text_action_stop {
			return c.Upsert(fromId, contactPhoneNumber, action)
		} else if action == SMS_Text_action_help {
			return c.HandleActionHelp(contactPhoneNumber)
		}
	}

	return nil
}

const (
	SysDialpadUserid = "5243348849786880"
	SysDialpadName   = "System"
)

func (c *UnsubscribesbuzUsecase) HandleActionHelp(phone string) error {

	leadVsTUser, err := c.ClientCaseUsecase.GetLeadVSByPhone(phone)
	if err != nil {
		return err
	}
	if leadVsTUser == nil {
		return errors.New("leadVsTUser is nil")
	}
	vsDialpadPhoneNumber := leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_dialpad_phonenumber)
	if vsDialpadPhoneNumber == "" {
		return errors.New("vsDialpadPhoneNumber is empty")
	}

	vsEmail := leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_email)

	tCase, err := c.ClientCaseUsecase.GetByPhone(phone)
	if err != nil {
		return err
	}
	clientName := ""
	if tCase != nil {
		clientName = tCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	text := c.DialpadbuzInternalUsecase.ClientSMSReplyAlertText(clientName, phone, time.Now(), SMS_Text_action_help)

	if vsEmail != "" {
		er := c.DialpadbuzInternalUsecase.CreateTaskClientSMSReplyAlertForEmail(clientName, phone, time.Now(), SMS_Text_action_help, vsEmail)
		if er != nil {
			c.log.Error(er)
		}
	}
	c.LogUsecase.SaveLog(0, "Unsubscribesbuz:HandleActionHelp", map[string]interface{}{
		"ReceivePhone":        vsDialpadPhoneNumber,
		"ReceiveText":         text,
		"ReceiveName":         leadVsTUser.CustomFields.TextValueByNameBasic(UserFieldName_fullname),
		"SenderDialpadUserid": SysDialpadUserid,
		"SenderUserFullName":  SysDialpadName,
	})

	if configs.Enable_SMS_New_Version_Debug {
		return nil
	}

	err = c.DialpadUsecase.SendSmsNoFilter(vsDialpadPhoneNumber, text, SysDialpadUserid, 0, "Unsubscribesbuz:HandleActionHelp")
	if err != nil {
		c.log.Error(err)
	}

	return nil
}

func (c *UnsubscribesbuzUsecase) Upsert(latestFromId string, contactPhoneNumber string, action string) error {

	status := Unsubscribes_Status_No
	if action == SMS_Text_action_stop {
		status = Unsubscribes_Status_Yes
	}
	if status != Unsubscribes_Status_No && status != Unsubscribes_Status_Yes {
		return errors.New("status is error")
	}
	entity, err := c.UnsubscribesUsecase.GetByCond(Eq{"contact_phone_number": contactPhoneNumber})
	if err != nil {
		return err
	}

	if entity == nil {
		entity = &UnsubscribesEntity{
			CreatedAt:          time.Now().Unix(),
			ContactPhoneNumber: contactPhoneNumber,
		}

	} else {
		if entity.Status == status {
			return nil
		}

	}

	entity.UpdatedAt = time.Now().Unix()
	entity.Status = status
	entity.LatestFromId = latestFromId
	entity.BizDeletedAt = 0
	err = c.CommonUsecase.DB().Save(&entity).Error
	if err != nil {
		return err
	}

	er := c.NotifyAdmin(contactPhoneNumber, action)
	if er != nil {
		c.log.Error(er)
	}

	if action == SMS_Text_action_stop {
		err := c.DialpadbuzUsecase.HandleAfterActionStop(contactPhoneNumber)
		if err != nil {
			c.log.Error(err, "contactPhoneNumber:", contactPhoneNumber)
		}
	}

	return nil
}

// NotifyAdmin contactPhoneNumber格式有 +1
func (c *UnsubscribesbuzUsecase) NotifyAdmin(contactPhoneNumber string, action string) error {
	tClientCase, err := c.ClientCaseUsecase.GetByPhone(contactPhoneNumber)
	if err != nil {
		c.log.Error(err)
	}
	clientName := ""
	if tClientCase != nil {
		clientName = tClientCase.CustomFields.TextValueByNameBasic(FieldName_deal_name)
	}
	er := c.DialpadbuzInternalUsecase.ClientSMSReplyAlertForAdmin(clientName, contactPhoneNumber, time.Now(), action)
	if er != nil {
		c.log.Error(er)
	}
	er = c.DialpadbuzInternalUsecase.HandleClientSMSReplyAlertTextForEmailAdmin(clientName, contactPhoneNumber, time.Now(), action)
	if er != nil {
		c.log.Error(er)
	}
	return nil
}
