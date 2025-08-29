package tests

import (
	"fmt"
	"testing"
	"time"
)

func Test_MonitoredEmailsUsecase_Create(t *testing.T) {
	email := "xxx@gmail.com"
	password := "xxx"
	description := "test description"
	uc := UT.MonitoredEmailsUsecase

	// Ensure clean state
	_ = uc.Delete(email)

	// Create
	if err := uc.Create(email, password, description); err != nil {
		t.Fatalf("Create error: %v", err)
	}
}

func Test_MonitoredEmailsUsecase_CRUD(t *testing.T) {
	uc := UT.MonitoredEmailsUsecase
	if uc == nil {
		t.Fatal("MonitoredEmailsUsecase is nil (wire generation may be required)")
	}

	email := fmt.Sprintf("monitored_%d@example.com", time.Now().UnixNano())
	description := "test email for monitored_emails CRUD"
	password := "123456"

	// Ensure clean state
	_ = uc.Delete(email)

	// Create
	if err := uc.Create(email, password, description); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	// Get
	got, err := uc.GetByEmailAddress(email)
	if err != nil {
		t.Fatalf("GetByEmailAddress error: %v", err)
	}
	if got == nil {
		t.Fatalf("GetByEmailAddress returned nil record for %s", email)
	}
	if got.EmailAddress != email {
		t.Fatalf("expected email %s, got %s", email, got.EmailAddress)
	}
	if !got.IsActive {
		t.Fatalf("expected IsActive true after create")
	}

	// ListActive should contain our email now
	list, err := uc.ListActive()
	if err != nil {
		t.Fatalf("ListActive error: %v", err)
	}
	found := false
	for _, v := range list {
		if v != nil && v.EmailAddress == email {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("ListActive should include created email %s", email)
	}

	// UpdateActive -> false
	if err := uc.UpdateActive(email, false); err != nil {
		t.Fatalf("UpdateActive(false) error: %v", err)
	}
	got, err = uc.GetByEmailAddress(email)
	if err != nil {
		t.Fatalf("Get after UpdateActive(false) error: %v", err)
	}
	if got == nil || got.IsActive {
		t.Fatalf("expected IsActive false after UpdateActive(false)")
	}

	// Ensure not present in ListActive now
	list, err = uc.ListActive()
	if err != nil {
		t.Fatalf("ListActive after deactivation error: %v", err)
	}
	for _, v := range list {
		if v != nil && v.EmailAddress == email {
			t.Fatalf("deactivated email %s should not be in ListActive", email)
		}
	}

	// UpdateActive -> true
	if err := uc.UpdateActive(email, true); err != nil {
		t.Fatalf("UpdateActive(true) error: %v", err)
	}

	// Cleanup
	if err := uc.Delete(email); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}
