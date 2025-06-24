package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetDiff returns the diff of unstaged changes
func GetDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetCommitDiff returns the diff of a specific commit
func GetCommitDiff(commitHash string) (string, error) {
	cmd := exec.Command("git", "show", commitHash)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git show failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
} 