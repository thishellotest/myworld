package biz

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/goccy/go-json"
	"time"
	"vbc/internal/conf"
	"vbc/internal/config_box"
)

// 客户更新了文件:VA Medical Records/Service Treatment Records/Private Medical Records
const ReminderEventType_ClientUpdateFiles = "ClientUpdateFiles"

// Reminder_IncrType client cases
const Reminder_IncrType = ""

type ReminderEventEntity struct {
	ID           int32 `gorm:"primaryKey"`
	HandleStatus int
	IncrType     string
	IncrId       int32
	EventType    string
	EventData    string
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    int64
}

func (ReminderEventEntity) TableName() string {
	return "reminder_events"
}

// GroupId 把事件分为一组，统一触发
func (c *ReminderEventEntity) GroupId() string {
	return fmt.Sprintf("%s:%d:%s", c.IncrType, c.IncrId, c.EventType)
}

func (c *ReminderEventEntity) GetReminderClientUpdateFilesEventVo() *ReminderClientUpdateFilesEventVo {
	var res ReminderClientUpdateFilesEventVo
	err := json.Unmarshal([]byte(c.EventData), &res)
	if err != nil {
		return nil
	}
	return &res
}

type ReminderEventUsecase struct {
	log           *log.Helper
	CommonUsecase *CommonUsecase
	conf          *conf.Data
	DBUsecase[ReminderEventEntity]
}

func NewReminderEventUsecase(logger log.Logger,
	CommonUsecase *CommonUsecase,
	conf *conf.Data) *ReminderEventUsecase {
	uc := &ReminderEventUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}

	uc.DBUsecase.DB = CommonUsecase.DB()

	return uc
}

type ReminderClientUpdateFilesEventVo struct {
	Items []*ReminderClientUpdateFilesEventVoItem
}

type ReminderClientUpdateFilesEventVoItem struct {
	BoxResName     string
	BoxResId       string
	BoxPath        string
	SourceBoxResId string
	SourceBoxPath  string
	BoxResType     config_box.BoxResType
}

func (c *ReminderClientUpdateFilesEventVoItem) GetUrl() string {
	prefixUrl := ""
	if c.BoxResType == config_box.BoxResType_file {
		prefixUrl = "https://veteranbenefitscenter.app.box.com/file/"
	} else {
		prefixUrl = "https://veteranbenefitscenter.app.box.com/folder/"
	}
	return prefixUrl + c.BoxResId
}

func (c *ReminderEventUsecase) AddClientUpdateFilesEvent(IncrId int32, EventData *ReminderClientUpdateFilesEventVo) error {
	currentTime := time.Now()
	event := &ReminderEventEntity{
		IncrType:  Reminder_IncrType,
		IncrId:    IncrId,
		EventType: ReminderEventType_ClientUpdateFiles,
		EventData: InterfaceToString(EventData),
		CreatedAt: currentTime.Unix(),
		UpdatedAt: currentTime.Unix(),
	}
	return c.CommonUsecase.DB().Create(&event).Error
}
