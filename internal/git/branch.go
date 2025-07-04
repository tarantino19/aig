package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GetCurrentBranch returns the current git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git branch failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// ExtractCommitDetails extracts the commit type and ticket number from the branch name.
func ExtractCommitDetails(branchName string) (string, string) {
	branchName = strings.ToLower(branchName)
	var commitType, ticketNumber string

	// Remove date patterns (YYYYMMDD) from branch name
	dateRegex := regexp.MustCompile(`\d{8}`)
	branchName = dateRegex.ReplaceAllString(branchName, "")

	// Extract ticket number
	re := regexp.MustCompile(`(\d{4,5})`)
	matches := re.FindStringSubmatch(branchName)
	if len(matches) > 1 {
		ticketNumber = matches[1]
	}

	// Determine commit type
	if strings.Contains(branchName, "fix") || strings.Contains(branchName, "bugfix") {
		commitType = "fix"
	} else if strings.Contains(branchName, "feature") || strings.Contains(branchName, "feat") {
		commitType = "feat"
	}

	return commitType, ticketNumber
}
