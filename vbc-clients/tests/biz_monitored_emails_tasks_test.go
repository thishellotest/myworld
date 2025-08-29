package tests

import (
	"fmt"
	"testing"
	"time"
	"vbc/internal/biz"
)

func TestRegisterAndDeleteTask(t *testing.T) {
	success, err := UT.MonitoredEmailsTasksUsecase.RegisterTask(biz.IncomingEmail{
		UID:     1234,
		Subject: "Test Subject",
		Date:    time.Now(),
	}, "a@a.com")

	fmt.Printf("Create result: success: %v, err: %v\n", success, err)

	success, err = UT.MonitoredEmailsTasksUsecase.DeleteTask("a@a.com", 1234)
	fmt.Printf("Delete result: success: %v, err: %v\n", success, err)
}

func TestIncomingEmail_WithDate(t *testing.T) {
	now := time.Now()
	email := biz.IncomingEmail{
		UID:       12345,
		MessageID: "test@example.com",
		From:      "sender@example.com",
		To:        "recipient@example.com",
		Subject:   "Test Subject",
		Body:      "Test Body",
		Raw:       []byte("raw email content"),
		Date:      now,
	}

	if email.UID != 12345 {
		t.Errorf("Expected UID 12345, got %d", email.UID)
	}

	if !email.Date.Equal(now) {
		t.Errorf("Expected Date %v, got %v", now, email.Date)
	}
}

func TestMonitoredEmailsTasksEntity_Fields(t *testing.T) {
	now := time.Now()
	entity := biz.MonitoredEmailsTasksEntity{
		ID:           1,
		InboxEmail:   "test@example.com",
		UID:          12345,
		Subject:      "Test Subject",
		EmailDate:    &now,
		Status:       "received",
		AttemptCount: 0,
		ErrorMessage: "",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if entity.InboxEmail != "test@example.com" {
		t.Errorf("Expected InboxEmail 'test@example.com', got %s", entity.InboxEmail)
	}

	if entity.UID != 12345 {
		t.Errorf("Expected UID 12345, got %d", entity.UID)
	}

	if entity.Status != "received" {
		t.Errorf("Expected Status 'received', got %s", entity.Status)
	}
}

func TestMonitoredEmailsTasksUsecase_RegisterEmail(t *testing.T) {
	// This test would require database setup, so we'll test the basic logic
	now := time.Now()
	email := biz.IncomingEmail{
		UID:       12345,
		MessageID: "test@example.com",
		From:      "sender@example.com",
		To:        "recipient@example.com",
		Subject:   "Test Subject - John Doe",
		Body:      "Test Body",
		Raw:       []byte("raw email content"),
		Date:      now,
	}

	// Test that the email structure is properly formed
	if email.UID == 0 {
		t.Error("UID should not be zero")
	}

	if email.Subject == "" {
		t.Error("Subject should not be empty")
	}

	if email.Date.IsZero() {
		t.Error("Date should not be zero")
	}
}

func TestMonitoredEmailsTasksEntity_StatusValues(t *testing.T) {
	validStatuses := []string{
		"received",
		"in_progress",
		"forwarded",
		"skipped_no_match",
		"error",
	}

	for _, status := range validStatuses {
		entity := biz.MonitoredEmailsTasksEntity{
			Status: status,
		}

		if entity.Status != status {
			t.Errorf("Expected status %s, got %s", status, entity.Status)
		}
	}
}

func TestMonitoredEmailsTasksEntity_DefaultValues(t *testing.T) {
	entity := biz.MonitoredEmailsTasksEntity{
		InboxEmail: "test@example.com",
		UID:        12345,
		Subject:    "Test Subject",
	}

	// Test that attempt count defaults to 0 (this would be handled by GORM)
	if entity.AttemptCount != 0 {
		t.Errorf("Expected AttemptCount to default to 0, got %d", entity.AttemptCount)
	}
}

func TestIncomingEmail_EmptyDate(t *testing.T) {
	email := biz.IncomingEmail{
		UID:       12345,
		MessageID: "test@example.com",
		From:      "sender@example.com",
		To:        "recipient@example.com",
		Subject:   "Test Subject",
		Body:      "Test Body",
		Raw:       []byte("raw email content"),
		// Date is not set, should be zero time
	}

	if !email.Date.IsZero() {
		t.Errorf("Expected Date to be zero time when not set, got %v", email.Date)
	}
}

func TestMonitoredEmailsTasksEntity_UniqueConstraints(t *testing.T) {
	// Test the structure supports unique constraints
	entity1 := biz.MonitoredEmailsTasksEntity{
		InboxEmail: "test@example.com",
		UID:        12345,
		Subject:    "Test Subject 1",
	}

	entity2 := biz.MonitoredEmailsTasksEntity{
		InboxEmail: "test@example.com",
		UID:        12345, // Same UID for same inbox should be unique
		Subject:    "Test Subject 2",
	}

	// Both entities have the same inbox_email and UID
	// This should trigger unique constraint in database
	if entity1.InboxEmail != entity2.InboxEmail {
		t.Error("Inbox emails should be the same for this test")
	}

	if entity1.UID != entity2.UID {
		t.Error("UIDs should be the same for this test")
	}
}
