package biz

type EmailLogEntity struct {
	ID         int32 `gorm:"primaryKey"`
	ClientId   int32
	SubId      int32
	Email      string
	SenderMail string
	SenderName string
	TaskId     int32
	Tpl        string
	Subject    string
	Body       string
	CreatedAt  int64
}

// email_log

func (EmailLogEntity) TableName() string {
	return "email_log"
}
