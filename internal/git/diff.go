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
		return "", fmt.Errorf("git show %s failed: %w, stderr: %s", commitHash, err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetCommitRangeDiff returns the diff of a specific commit range
func GetCommitRangeDiff(commitRange string) (string, error) {
	cmd := exec.Command("git", "diff", commitRange)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff %s failed: %w, stderr: %s", commitRange, err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

// GetBranchDiff returns the diff between the current branch and a specified branch
func GetBranchDiff(branchName string) (string, error) {
	cmd := exec.Command("git", "diff", branchName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff %s failed: %w, stderr: %s", branchName, err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
} 