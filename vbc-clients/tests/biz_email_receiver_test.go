package tests

import (
	"fmt"
	"testing"
	"time"
	"vbc/internal/biz"
)

const (
	EmailAddress  = "xxx@gmail.com"
	EmailPassword = "xxx"
)

func Test_FetchRecentEmails(t *testing.T) {
	maxCount := uint32(1)
	emails, err := biz.FetchRecentEmails(EmailAddress, EmailPassword, maxCount)
	if err != nil {
		t.Fatalf("FetchInboxEmail error: %v", err)
	}
	for _, email := range emails {
		fmt.Printf("UID: %d. Subject: %s\n", email.UID, email.Subject)
		fmt.Printf("HTMLBody: %s\n", email.HTMLBody)
		fmt.Printf("Body: %s\n", email.Body)
		fmt.Printf("Subject: %s\n", email.Subject)
	}
}

func Test_FetchEmailsByUIDRange(t *testing.T) {
	emails, _hasMore, err := biz.FetchEmailsByUIDRange(EmailAddress, EmailPassword, 3176, 10)
	if err != nil {
		t.Fatalf("FetchInboxEmail error: %v", err)
	}
	fmt.Printf("hasMore: %v. count: %d\n", _hasMore, len(emails))
	for _, email := range emails {
		fmt.Printf("UID: %d\n", email.UID)
	}
}

func Test_FetchEmailByUID(t *testing.T) {
	email, err := biz.FetchEmailByUID(EmailAddress, EmailPassword, 3174)
	if err != nil {
		t.Fatalf("FetchInboxEmail error: %v", err)
	}
	fmt.Printf("UID: %d\n", email.UID)
}

func TestIncomingEmail_Structure(t *testing.T) {
	now := time.Now()
	email := biz.IncomingEmail{
		UID:       12345,
		MessageID: "test@example.com",
		From:      "sender@example.com",
		To:        "recipient@example.com",
		Subject:   "Test Subject - John Doe",
		Body:      "Test Body Content",
		Raw:       []byte("raw email content"),
		Date:      now,
	}

	// Test all fields are properly set
	if email.UID != 12345 {
		t.Errorf("Expected UID 12345, got %d", email.UID)
	}

	if email.MessageID != "test@example.com" {
		t.Errorf("Expected MessageID 'test@example.com', got %s", email.MessageID)
	}

	if email.From != "sender@example.com" {
		t.Errorf("Expected From 'sender@example.com', got %s", email.From)
	}

	if email.To != "recipient@example.com" {
		t.Errorf("Expected To 'recipient@example.com', got %s", email.To)
	}

	if email.Subject != "Test Subject - John Doe" {
		t.Errorf("Expected Subject 'Test Subject - John Doe', got %s", email.Subject)
	}

	if email.Body != "Test Body Content" {
		t.Errorf("Expected Body 'Test Body Content', got %s", email.Body)
	}

	if string(email.Raw) != "raw email content" {
		t.Errorf("Expected Raw 'raw email content', got %s", string(email.Raw))
	}

	if !email.Date.Equal(now) {
		t.Errorf("Expected Date %v, got %v", now, email.Date)
	}
}

func TestIncomingEmail_EmptyFields(t *testing.T) {
	email := biz.IncomingEmail{
		UID: 1,
		// Other fields are empty/zero values
	}

	if email.UID != 1 {
		t.Errorf("Expected UID 1, got %d", email.UID)
	}

	if email.MessageID != "" {
		t.Errorf("Expected empty MessageID, got %s", email.MessageID)
	}

	if email.From != "" {
		t.Errorf("Expected empty From, got %s", email.From)
	}

	if email.To != "" {
		t.Errorf("Expected empty To, got %s", email.To)
	}

	if email.Subject != "" {
		t.Errorf("Expected empty Subject, got %s", email.Subject)
	}

	if email.Body != "" {
		t.Errorf("Expected empty Body, got %s", email.Body)
	}

	if len(email.Raw) != 0 {
		t.Errorf("Expected empty Raw, got %d bytes", len(email.Raw))
	}

	if !email.Date.IsZero() {
		t.Errorf("Expected zero Date, got %v", email.Date)
	}
}

func TestIncomingEmail_UIDHandling(t *testing.T) {
	testCases := []struct {
		name string
		uid  uint32
	}{
		{"Small UID", 1},
		{"Medium UID", 12345},
		{"Large UID", 4294967295}, // Max uint32
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email := biz.IncomingEmail{
				UID: tc.uid,
			}

			if email.UID != tc.uid {
				t.Errorf("Expected UID %d, got %d", tc.uid, email.UID)
			}
		})
	}
}

func TestIncomingEmail_DateHandling(t *testing.T) {
	testCases := []struct {
		name string
		date time.Time
	}{
		{"Current time", time.Now()},
		{"Past time", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
		{"Future time", time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)},
		{"Zero time", time.Time{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email := biz.IncomingEmail{
				UID:  1,
				Date: tc.date,
			}

			if !email.Date.Equal(tc.date) {
				t.Errorf("Expected Date %v, got %v", tc.date, email.Date)
			}

			// Test zero time handling
			if tc.name == "Zero time" && !email.Date.IsZero() {
				t.Errorf("Expected zero time, got %v", email.Date)
			}
		})
	}
}

func TestIncomingEmail_RawDataHandling(t *testing.T) {
	testCases := []struct {
		name string
		raw  []byte
	}{
		{"Empty raw", []byte{}},
		{"Small raw", []byte("small")},
		{"Large raw", make([]byte, 1024)}, // 1KB
		{"Binary raw", []byte{0x00, 0x01, 0x02, 0xFF}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			email := biz.IncomingEmail{
				UID: 1,
				Raw: tc.raw,
			}

			if len(email.Raw) != len(tc.raw) {
				t.Errorf("Expected Raw length %d, got %d", len(tc.raw), len(email.Raw))
			}

			// Compare byte by byte for non-empty arrays
			if len(tc.raw) > 0 {
				for i, b := range tc.raw {
					if email.Raw[i] != b {
						t.Errorf("Expected Raw[%d] = %d, got %d", i, b, email.Raw[i])
						break
					}
				}
			}
		})
	}
}
