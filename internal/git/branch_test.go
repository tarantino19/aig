package git

import (
	"testing"
)

func TestExtractCommitDetails(t *testing.T) {
	tests := []struct {
		branchName   string
		expectedType string
		expectedTicket string
	}{
		{"feature/1234-new-login", "feat", "1234"},
		{"fix/5678-fix-bug", "fix", "5678"},
		{"bugfix/8765-another-bug", "fix", "8765"},
		{"feat/4321-add-feature", "feat", "4321"},
		{"hotfix/9999-critical-issue", "fix", "9999"},
		{"release/v1.0", "", ""},
		{"no-ticket-feat", "feat", ""},
		{"12345-fix-something", "fix", "12345"},
		{"feature/1234-20250620-new-login", "feat", "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.branchName, func(t *testing.T) {
			commitType, ticketNumber := ExtractCommitDetails(tt.branchName)
			if commitType != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, commitType)
			}
			if ticketNumber != tt.expectedTicket {
				t.Errorf("expected ticket %q, got %q", tt.expectedTicket, ticketNumber)
			}
		})
	}
}
