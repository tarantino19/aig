package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommitOptions contains options for retrieving commits
type CommitOptions struct {
	Number int
	Branch string
	From   string
	To     string
}

// Commit represents a git commit
type Commit struct {
	Hash    string
	Author  string
	Date    string
	Message string
}

// GetCommits retrieves commits based on the provided options
func GetCommits(opts CommitOptions) ([]Commit, error) {
	args := []string{"log", "--pretty=format:%H|%an|%ad|%s", "--date=short"}
	
	if opts.Number > 0 {
		args = append(args, fmt.Sprintf("-n%d", opts.Number))
	}
	
	if opts.Branch != "" {
		args = append(args, opts.Branch)
	}
	
	if opts.From != "" && opts.To != "" {
		args = append(args, fmt.Sprintf("%s..%s", opts.From, opts.To))
	} else if opts.From != "" {
		args = append(args, fmt.Sprintf("%s..HEAD", opts.From))
	}
	
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w, stderr: %s", err, stderr.String())
	}
	
	return parseCommits(out.String()), nil
}

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

func parseCommits(output string) []Commit {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]Commit, 0, len(lines))
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "|", 4)
		if len(parts) == 4 {
			commits = append(commits, Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Date:    parts[2],
				Message: parts[3],
			})
		}
	}
	
	return commits
} 