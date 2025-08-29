package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
)

func Test_MonitoredEmailsJobUsecase_Handle(t *testing.T) {
	UT.MonitoredEmailsJobUsecase.Handle()
}

func Test_MonitoredEmailsJobDev_Handle(t *testing.T) {
	cc := UT.MonitoredEmailsJobUsecase.GetForwardCcEmails()
	if len(cc) > 0 {
		fmt.Printf("cc: %v\n", cc)
	} else {
		fmt.Printf("No cc\n")
	}
}

func Test_FindNameInSubject(t *testing.T) {
	tests := []struct {
		name     string
		subject  string
		expected string
	}{
		{
			name:     "normal case",
			subject:  "New Client Case #123 - John Smith",
			expected: "John Smith",
		},
		{
			name:     "multiple separators",
			subject:  "Case #123 - Status Update - Jane Doe",
			expected: "Jane Doe",
		},
		{
			name:     "no separator",
			subject:  "Simple Subject Line",
			expected: "",
		},
		{
			name:     "empty subject",
			subject:  "",
			expected: "",
		},
		{
			name:     "only separator",
			subject:  " - ",
			expected: "",
		},
		{
			name:     "separator at end",
			subject:  "Case #123 - ",
			expected: "",
		},
		{
			name:     "multiple spaces after separator",
			subject:  "Case #123 -    Robert Johnson   ",
			expected: "Robert Johnson",
		},
		{
			name:     "multiple -",
			subject:  "转发：[GitHub] Your Dependabot alerts for the week of Aug 5 - Aug 12 - Lin Chen ",
			expected: "Lin Chen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := biz.FindNameInSubject(tt.subject)
			if got != tt.expected {
				t.Errorf("FindNameInSubject(%q) = %q, want %q", tt.subject, got, tt.expected)
			}
		})
	}
}

func Test_NormalizeRecipientList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "single email",
			input:    "john@example.com",
			expected: "john@example.com",
		},
		{
			name:     "single email with spaces",
			input:    "  john@example.com  ",
			expected: "john@example.com",
		},
		{
			name:     "comma separated emails",
			input:    "john@example.com,jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "comma separated with spaces",
			input:    "john@example.com, jane@example.com, bob@example.com",
			expected: "john@example.com;jane@example.com;bob@example.com",
		},
		{
			name:     "semicolon separated emails",
			input:    "john@example.com;jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "semicolon separated with spaces",
			input:    "john@example.com; jane@example.com; bob@example.com",
			expected: "john@example.com;jane@example.com;bob@example.com",
		},
		{
			name:     "mixed separators with commas",
			input:    "john@example.com,jane@example.com;bob@example.com",
			expected: "john@example.com;jane@example.com;bob@example.com",
		},
		{
			name:     "empty entries with commas",
			input:    "john@example.com,,jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "empty entries with semicolons",
			input:    "john@example.com;;jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "trailing comma",
			input:    "john@example.com,jane@example.com,",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "trailing semicolon",
			input:    "john@example.com;jane@example.com;",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "leading comma",
			input:    ",john@example.com,jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "leading semicolon",
			input:    ";john@example.com;jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "multiple consecutive commas",
			input:    "john@example.com,,,jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "multiple consecutive semicolons",
			input:    "john@example.com;;;jane@example.com",
			expected: "john@example.com;jane@example.com",
		},
		{
			name:     "complex whitespace and separators",
			input:    "  john@example.com  ,  , jane@example.com  ;  ;  bob@example.com  ",
			expected: "john@example.com;jane@example.com;bob@example.com",
		},
		{
			name:     "only semicolons no commas normalize",
			input:    "john@example.com; ; jane@example.com;;bob@example.com;",
			expected: "john@example.com;jane@example.com;bob@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := biz.NormalizeRecipientList(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeRecipientList(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
