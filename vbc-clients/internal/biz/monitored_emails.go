package biz

import (
	"time"

	"vbc/internal/conf"
	"vbc/lib/builder"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	MonitoredEmailFieldName_email_address   = "email_address"
	MonitoredEmailFieldName_description     = "description"
	MonitoredEmailFieldName_is_active       = "is_active"
	MonitoredEmailFieldName_last_seen_uid   = "last_seen_uid"
	MonitoredEmailFieldName_last_scanned_at = "last_scanned_at"
)

type MonitoredEmailEntity struct {
	ID            uint32     `gorm:"primaryKey;autoIncrement"`
	EmailAddress  string     `gorm:"uniqueIndex;size:255;not null"`
	EmailPassword string     `gorm:"size:255;not null"`
	Description   string     `gorm:"type:text"`
	IsActive      bool       `gorm:"default:true;not null"`
	LastSeenUID   uint64     `gorm:"default:0;not null"`
	LastScannedAt *time.Time `gorm:"type:datetime(6)"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;type:datetime(6)"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;type:datetime(6)"`
}

func (MonitoredEmailEntity) TableName() string {
	return "monitored_emails"
}

type MonitoredEmailsUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[MonitoredEmailEntity]
}

func NewMonitoredEmailsUsecase(
	logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *MonitoredEmailsUsecase {
	uc := &MonitoredEmailsUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

// GetByEmailAddress returns one by exact email address.
func (c *MonitoredEmailsUsecase) GetByEmailAddress(email string) (*MonitoredEmailEntity, error) {
	return c.GetByCond(builder.Eq{MonitoredEmailFieldName_email_address: email})
}

// ListActive returns all active monitored emails.
func (c *MonitoredEmailsUsecase) ListActive() ([]*MonitoredEmailEntity, error) {
	return c.AllByCond(builder.Eq{MonitoredEmailFieldName_is_active: true})
}

// Create inserts a new monitored email record. `is_active` defaults to true.
// Pass empty string for `description` if not provided.
func (c *MonitoredEmailsUsecase) Create(emailAddress string, password string, description string) error {
	entity := MonitoredEmailEntity{
		EmailAddress:  emailAddress,
		EmailPassword: password,
		Description:   description,
		IsActive:      true,
	}
	return c.CommonUsecase.DB().Create(&entity).Error
}

// UpdateActive updates the is_active flag by email address.
func (c *MonitoredEmailsUsecase) UpdateActive(emailAddress string, isActive bool) error {
	return c.UpdatesByCond(map[string]interface{}{
		MonitoredEmailFieldName_is_active: isActive,
	}, builder.Eq{MonitoredEmailFieldName_email_address: emailAddress})
}

// Delete removes a monitored email by email address.
func (c *MonitoredEmailsUsecase) Delete(emailAddress string) error {
	return c.CommonUsecase.DB().Where("email_address = ?", emailAddress).Delete(&MonitoredEmailEntity{}).Error
}

// UpdateLastSeenUID updates the last_seen_uid field for a monitored email.
func (c *MonitoredEmailsUsecase) UpdateLastSeenUID(id uint32, uid uint64) error {
	return c.UpdatesByCond(map[string]interface{}{
		MonitoredEmailFieldName_last_seen_uid: uid,
	}, builder.Eq{"id": id})
}

// UpdateLastScannedAt updates the last_scanned_at field for a monitored email.
func (c *MonitoredEmailsUsecase) UpdateLastScannedAt(id uint32) error {
	now := time.Now()
	return c.UpdatesByCond(map[string]interface{}{
		MonitoredEmailFieldName_last_scanned_at: &now,
	}, builder.Eq{"id": id})
}
