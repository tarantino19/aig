package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CreateCommit creates a git commit with the given message
func CreateCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git commit failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// Push pushes commits to the remote repository
func Push() error {
	cmd := exec.Command("git", "push")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git push failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// IsRepoClean checks if the repository has no uncommitted changes
func IsRepoClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("git status failed: %w, stderr: %s", err, stderr.String())
	}

	return out.String() == "", nil
}

// HasStagedChanges checks if there are staged changes
func HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	
	// git diff --cached --quiet returns 0 if no staged changes, 1 if there are changes
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil // There are staged changes
		}
		return false, fmt.Errorf("git diff failed: %w", err)
	}
	
	return false, nil // No staged changes
} 