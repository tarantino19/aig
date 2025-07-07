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
	
	prompt.WriteString("Provide a structured review using markdown headers (e.g., ## Summary, ## Issues, ## Suggestions, ## Security Risks, ## Performance Issues):")
	
	return prompt.String()
}

// GetPRDescriptionPrompt returns the prompt for generating PR descriptions
func GetPRDescriptionPrompt(currentBranch, targetBranch, diff string, commits []Commit, issueNumbers []string, platform string) string {
	var prompt strings.Builder
	
	prompt.WriteString("Generate a comprehensive Pull Request description based on the following information.\n\n")
	
	prompt.WriteString("Branch Information:\n")
	prompt.WriteString(fmt.Sprintf("- Current Branch: %s\n", currentBranch))
	prompt.WriteString(fmt.Sprintf("- Target Branch: %s\n", targetBranch))
	prompt.WriteString(fmt.Sprintf("- Platform: %s\n", platform))
	
	if len(issueNumbers) > 0 {
		prompt.WriteString(fmt.Sprintf("- Related Issues: %s\n", strings.Join(issueNumbers, ", ")))
	}
	
	prompt.WriteString("\nCommits in this branch:\n")
	for i, commit := range commits {
		if i >= 10 { // Limit to first 10 commits for prompt size
			prompt.WriteString(fmt.Sprintf("... and %d more commits\n", len(commits)-10))
			break
		}
		prompt.WriteString(fmt.Sprintf("- %s: %s\n", commit.Hash[:7], commit.Message))
	}
	
	prompt.WriteString("\nGenerate a PR description with the following structure:\n")
	prompt.WriteString("1. **Title**: Concise, descriptive title (50 chars max)\n")
	prompt.WriteString("2. **Summary**: Brief overview of what this PR accomplishes\n")
	prompt.WriteString("3. **Changes**: Bullet points of key changes made\n")
	prompt.WriteString("4. **Testing**: How the changes should be tested\n")
	prompt.WriteString("5. **Breaking Changes**: Any breaking changes (if applicable)\n\n")
	
	prompt.WriteString("Requirements:\n")
	prompt.WriteString("- Use clear, professional language\n")
	prompt.WriteString("- Focus on business value and impact\n")
	prompt.WriteString("- Include technical details where relevant\n")
	prompt.WriteString("- Mention any dependencies or requirements\n")
	
	switch platform {
	case "gitlab":
		prompt.WriteString("- Use GitLab-specific formatting\n")
		prompt.WriteString("- Use 'Closes #issue' for issue linking\n")
	case "bitbucket":
		prompt.WriteString("- Use Bitbucket-specific formatting\n")
		prompt.WriteString("- Use 'Fixes #issue' for issue linking\n")
	default: // github
		prompt.WriteString("- Use GitHub-specific formatting\n")
		prompt.WriteString("- Use 'Fixes #issue' for issue linking\n")
	}
	
	prompt.WriteString("\nCode changes:\n")
	prompt.WriteString("```diff\n")
	// Truncate diff if too long to fit in prompt
	if len(diff) > 8000 {
		prompt.WriteString(diff[:8000])
		prompt.WriteString("\n... (diff truncated for brevity)")
	} else {
		prompt.WriteString(diff)
	}
	prompt.WriteString("\n```\n\n")
	
	prompt.WriteString("Respond with a JSON object containing:\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"title\": \"PR title\",\n")
	prompt.WriteString("  \"summary\": \"Brief summary paragraph\",\n")
	prompt.WriteString("  \"changes\": [\"change 1\", \"change 2\", ...],\n")
	prompt.WriteString("  \"testing\": \"Testing instructions\",\n")
	prompt.WriteString("  \"breaking_changes\": [\"breaking change 1\", ...] // empty array if none\n")
	prompt.WriteString("}\n")
	
	return prompt.String()
} 