package biz

const (
	WebhookLog_From_website = "website"
	WebhookLog_From_jotform = "jotform"
	WebhookLog_From_dialpad = "dialpad"
)

type WebhookLogEntity struct {
	ID                 int32 `gorm:"primaryKey"`
	From               string
	HandleStatus       int
	HandleResult       int
	HandleResultDetail string
	Remarks            string
	Headers            string
	Query              string
	Body               string
	NeatBody           string
	CreatedAt          int64
	UpdatedAt          int64
	DeletedAt          int64
}

func (WebhookLogEntity) TableName() string {
	return "webhook_log"
}
