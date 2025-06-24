package prompts

import (
	"fmt"
	"strings"
)

// Commit represents a git commit for prompts
type Commit struct {
	Hash    string
	Author  string
	Date    string
	Message string
}

// GetCommitMessagePrompt returns the prompt for generating commit messages
func GetCommitMessagePrompt(diff, commitType, scope string, conventional bool) string {
	var prompt strings.Builder
	
	prompt.WriteString("Analyze the following git diff and generate a concise, conventional commit message.\n\n")
	
	prompt.WriteString("Rules:\n")
	if conventional {
		prompt.WriteString("1. Use conventional commit format: <type>(<scope>): <subject>\n")
		prompt.WriteString("2. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build\n")
		prompt.WriteString("3. Subject line max 50 characters\n")
		prompt.WriteString("4. Use present tense (\"add\" not \"added\")\n")
		prompt.WriteString("5. No period at the end of subject\n")
		prompt.WriteString("6. Include body if changes are complex (wrap at 72 chars)\n")
		prompt.WriteString("7. Include footer for breaking changes or issue references\n")
	} else {
		prompt.WriteString("1. Subject line max 50 characters\n")
		prompt.WriteString("2. Use imperative mood (\"Add feature\" not \"Added feature\")\n")
		prompt.WriteString("3. Capitalize the subject line\n")
		prompt.WriteString("4. No period at the end\n")
		prompt.WriteString("5. Include body if needed (wrap at 72 chars)\n")
	}
	
	if commitType != "" {
		prompt.WriteString(fmt.Sprintf("\nCommit type must be: %s\n", commitType))
	}
	
	if scope != "" {
		prompt.WriteString(fmt.Sprintf("Scope must be: %s\n", scope))
	}
	
	prompt.WriteString("\nDiff:\n")
	prompt.WriteString("```\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n```\n\n")
	
	prompt.WriteString("Generate the commit message (respond with ONLY the commit message, no explanations):")
	
	return prompt.String()
}

// GetSummaryPrompt returns the prompt for generating commit summaries
func GetSummaryPrompt(commits []Commit, groupByType, changelog bool) string {
	var prompt strings.Builder
	
	prompt.WriteString("Summarize the following git commits ")
	if changelog {
		prompt.WriteString("in changelog format.\n\n")
		prompt.WriteString("Format the output as a proper changelog entry with:\n")
		prompt.WriteString("- Version header\n")
		prompt.WriteString("- Date\n")
		prompt.WriteString("- Grouped changes by type (Features, Bug Fixes, etc.)\n")
		prompt.WriteString("- Clear, user-facing descriptions\n\n")
	} else {
		prompt.WriteString("in a clear, concise manner.\n\n")
		if groupByType {
			prompt.WriteString("Group commits by their type (feat, fix, docs, etc.).\n")
		}
	}
	
	prompt.WriteString("Commits:\n\n")
	for _, commit := range commits {
		prompt.WriteString(fmt.Sprintf("- %s: %s\n", commit.Hash[:7], commit.Message))
	}
	
	prompt.WriteString("\nGenerate the summary:")
	
	return prompt.String()
}

// GetReviewPrompt returns the prompt for code review
func GetReviewPrompt(diff string, focusAreas []string, security, performance bool) string {
	var prompt strings.Builder
	
	prompt.WriteString("Review the following code changes and provide constructive feedback.\n\n")
	
	prompt.WriteString("Focus on:\n")
	prompt.WriteString("1. Potential bugs or errors\n")
	prompt.WriteString("2. Code quality and best practices\n")
	prompt.WriteString("3. Readability and maintainability\n")
	
	if security {
		prompt.WriteString("4. Security vulnerabilities (PRIORITY)\n")
	}
	
	if performance {
		prompt.WriteString("5. Performance issues and optimization opportunities (PRIORITY)\n")
	}
	
	if len(focusAreas) > 0 {
		prompt.WriteString(fmt.Sprintf("6. Specific areas: %s\n", strings.Join(focusAreas, ", ")))
	}
	
	prompt.WriteString("\nProvide:\n")
	prompt.WriteString("- Summary of the changes\n")
	prompt.WriteString("- List of issues found (if any)\n")
	prompt.WriteString("- Suggestions for improvement\n")
	if security {
		prompt.WriteString("- Security risks and mitigations\n")
	}
	if performance {
		prompt.WriteString("- Performance concerns and solutions\n")
	}
	
	prompt.WriteString("\nCode changes:\n")
	prompt.WriteString("```diff\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n```\n\n")
	
	prompt.WriteString("Provide a structured review:")
	
	return prompt.String()
} 