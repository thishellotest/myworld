package biz

import (
	"time"

	"vbc/internal/conf"
	"vbc/lib/builder"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	MonitoredEmailsTasksFieldName_inbox_email   = "inbox_email"
	MonitoredEmailsTasksFieldName_uid           = "uid"
	MonitoredEmailsTasksFieldName_status        = "status"
	MonitoredEmailsTasksFieldName_attempt_count = "attempt_count"
	MaxAttemptCount                             = 3  // Max 3 attempts
	InProgressTimeoutMinutes                    = 10 // Timeout for in_progress tasks in minutes
)

// MonitoredEmailsTasksEntity represents the monitored_emails_tasks table
type MonitoredEmailsTasksEntity struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement"`
	InboxEmail   string     `gorm:"size:320;not null;index:uk_inbox_email;index:uk_mail_uid,priority:1;index:idx_pending,priority:1"`
	UID          uint64     `gorm:"not null;index:uk_mail_uid,priority:2,unique"`
	Subject      string     `gorm:"type:text"`
	EmailDate    *time.Time `gorm:"type:datetime(6);index"`
	Status       string     `gorm:"type:enum('received','in_progress','forwarded','skipped_no_match','error','signed');default:'received';not null;index:idx_pending,priority:2"`
	AttemptCount int        `gorm:"default:0;not null"`
	ErrorMessage string     `gorm:"type:text"`
	Data         string     `gorm:"type:json;default:null"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;type:datetime(6)"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;type:datetime(6)"`
}

func (MonitoredEmailsTasksEntity) TableName() string {
	return "monitored_emails_tasks"
}

// MonitoredEmailsTasksUsecase handles operations for monitored_emails_tasks table
type MonitoredEmailsTasksUsecase struct {
	log           *log.Helper
	conf          *conf.Data
	CommonUsecase *CommonUsecase
	DBUsecase[MonitoredEmailsTasksEntity]
}

func NewMonitoredEmailsTasksUsecase(
	logger log.Logger,
	conf *conf.Data,
	CommonUsecase *CommonUsecase,
) *MonitoredEmailsTasksUsecase {
	uc := &MonitoredEmailsTasksUsecase{
		log:           log.NewHelper(logger),
		CommonUsecase: CommonUsecase,
		conf:          conf,
	}
	uc.DBUsecase.DB = CommonUsecase.DB()
	return uc
}

// RegisterTask attempts to register a new email task, returns whether it's a new email
func (c *MonitoredEmailsTasksUsecase) RegisterTask(email IncomingEmail, inboxEmail string) (bool, error) {
	entity := MonitoredEmailsTasksEntity{
		InboxEmail: inboxEmail,
		UID:        uint64(email.UID),
		Subject:    email.Subject,
		EmailDate:  &email.Date,
		Status:     "received",
	}

	// Try to insert, ignore if duplicate exists
	result := c.CommonUsecase.DB().Create(&entity)
	if result.Error != nil {
		// Check if it's a duplicate key error
		if isDuplicateKeyError(result.Error) {
			return false, nil // Not a new email, but no error
		}
		return false, result.Error
	}

	// If affected rows > 0, it was a new insert
	return result.RowsAffected > 0, nil
}

func (c *MonitoredEmailsTasksUsecase) DeleteTask(inboxEmail string, uid uint32) (bool, error) {
	result := c.CommonUsecase.DB().Where("inbox_email = ? AND uid = ?", inboxEmail, uint64(uid)).Delete(&MonitoredEmailsTasksEntity{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// GetPendingTasks retrieves pending tasks (received, error, or stale in_progress status) for an inbox
func (c *MonitoredEmailsTasksUsecase) GetPendingTasks(inboxEmail string, limit int) ([]*MonitoredEmailsTasksEntity, error) {
	// Define timeout for in_progress tasks
	inProgressTimeout := time.Now().Add(-InProgressTimeoutMinutes * time.Minute)

	cond := builder.And(
		builder.Eq{MonitoredEmailsTasksFieldName_inbox_email: inboxEmail},
		builder.Or(
			// Normal pending tasks
			builder.In(MonitoredEmailsTasksFieldName_status, []string{"received", "error"}),
			// Recover stale in_progress tasks (older than 10 minutes)
			builder.And(
				builder.Eq{MonitoredEmailsTasksFieldName_status: "in_progress"},
				builder.Lt{"updated_at": inProgressTimeout},
			),
		),
		builder.Lt{MonitoredEmailsTasksFieldName_attempt_count: MaxAttemptCount},
	)
	return c.AllByCondWithOrderBy(cond, "id ASC", limit)
}

// GetByUID retrieves a task by inbox email and UID
func (c *MonitoredEmailsTasksUsecase) GetByUID(inboxEmail string, uid uint32) (*MonitoredEmailsTasksEntity, error) {
	cond := builder.And(
		builder.Eq{MonitoredEmailsTasksFieldName_inbox_email: inboxEmail},
		builder.Eq{MonitoredEmailsTasksFieldName_uid: uint64(uid)},
	)
	return c.GetByCond(cond)
}

// UpdateStatus updates the status and optionally error message of a task
func (c *MonitoredEmailsTasksUsecase) UpdateStatus(id uint64, status string, errorMsg string) error {
	updates := map[string]interface{}{
		MonitoredEmailsTasksFieldName_status: status,
	}
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	} else if status == "forwarded" {
		// Clear error message when status is forwarded
		updates["error_message"] = ""
	}
	return c.UpdatesByCond(updates, builder.Eq{"id": id})
}

// IncrementAttempt increments the attempt count for a task
func (c *MonitoredEmailsTasksUsecase) IncrementAttempt(id uint64) error {
	return c.CommonUsecase.DB().Model(&MonitoredEmailsTasksEntity{}).
		Where("id = ?", id).
		UpdateColumn("attempt_count", c.CommonUsecase.DB().Raw("attempt_count + 1")).Error
}

// SetInProgress marks a task as in_progress and increments attempt count
func (c *MonitoredEmailsTasksUsecase) SetInProgress(id uint64) error {
	updates := map[string]interface{}{
		MonitoredEmailsTasksFieldName_status: "in_progress",
	}
	err := c.UpdatesByCond(updates, builder.Eq{"id": id})
	if err != nil {
		return err
	}
	return c.IncrementAttempt(id)
}

// UpdateData updates the data field for a task
func (c *MonitoredEmailsTasksUsecase) UpdateData(id uint64, data string) error {
	updates := map[string]interface{}{
		"data": data,
	}
	return c.UpdatesByCond(updates, builder.Eq{"id": id})
}

// GetForwardedTasksOlderThanConfigSeconds retrieves tasks with 'forwarded' status older than configured seconds
func (c *MonitoredEmailsTasksUsecase) GetForwardedTasksOlderThanConfigSeconds() ([]*MonitoredEmailsTasksEntity, error) {
	// Get re-forward seconds from config, default to 24 hours (86400 seconds) if not configured
	reForwardSeconds := int32(86400) // 24 hours in seconds
	if c.conf != nil && c.conf.EmailMonitor != nil && c.conf.EmailMonitor.ReForwardSeconds > 0 {
		reForwardSeconds = c.conf.EmailMonitor.ReForwardSeconds
	}
	
	// Calculate the time threshold
	timeThreshold := time.Now().Add(-time.Duration(reForwardSeconds) * time.Second)
	
	cond := builder.And(
		builder.Eq{MonitoredEmailsTasksFieldName_status: "forwarded"},
		builder.Lt{"updated_at": timeThreshold},
	)
	
	return c.AllByCondWithOrderBy(cond, "id ASC", 0)
}

// isDuplicateKeyError checks if the error is a duplicate key constraint error
// This is a simple implementation - you might need to adjust based on your database driver
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL duplicate entry error patterns
	return contains(errStr, "Duplicate entry") ||
		contains(errStr, "duplicate key") ||
		contains(errStr, "UNIQUE constraint failed")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOfSubstring(s, substr) >= 0)))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
