package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/tarantino19/aig/internal/ai"
)

// ShowDryRun displays what will be committed in dry-run mode
func ShowDryRun(diff string) {
	fmt.Println(titleStyle.Render("🔍 Dry Run - Staged Changes:"))
	fmt.Println(codeBlockStyle.Render(truncateString(diff, 1000)))
	fmt.Println(mutedStyle.Render("\n(This is a preview. Remove --dry-run to generate commit message)"))
}

// ShowError displays an error message
func ShowError(err error) {
	fmt.Println(errorStyle.Render("❌ Error: " + err.Error()))
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	fmt.Println(successStyle.Render("✅ " + message))
}

// ShowInfo displays an info message
func ShowInfo(message string) {
	fmt.Println(infoStyle.Render("ℹ️  " + message))
}

// ShowWarning displays a warning message
func ShowWarning(message string) {
	fmt.Println(warningStyle.Render("⚠️  " + message))
}

// ShowCommitMessage displays a formatted commit message
func ShowCommitMessage(commitType, scope, message string) {
	typeStyle := GetCommitTypeStyle(commitType)
	
	var commit string
	if scope != "" {
		commit = fmt.Sprintf("%s(%s): %s", commitType, scope, message)
	} else {
		commit = fmt.Sprintf("%s: %s", commitType, message)
	}
	
	fmt.Println(headerStyle.Render("Generated Commit Message:"))
	fmt.Println(boxStyle.Render(typeStyle.Render(commitType) + mutedStyle.Render(commit[len(commitType):])))
}

// ShowDiff displays a diff with syntax highlighting
func ShowDiff(diff string) {
	lines := strings.Split(diff, "\n")
	var styledLines []string
	
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	removeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))
	
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+"):
			styledLines = append(styledLines, addStyle.Render(line))
		case strings.HasPrefix(line, "-"):
			styledLines = append(styledLines, removeStyle.Render(line))
		case strings.HasPrefix(line, "@@"):
			styledLines = append(styledLines, metaStyle.Render(line))
		case strings.HasPrefix(line, "diff"):
			styledLines = append(styledLines, metaStyle.Render(line))
		default:
			styledLines = append(styledLines, line)
		}
	}
	
	fmt.Println(strings.Join(styledLines, "\n"))
}

// GetSpinner returns a configured spinner
func GetSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return s
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// FormatList formats a list of items with bullets
func FormatList(items []string) string {
	var formatted []string
	for _, item := range items {
		formatted = append(formatted, fmt.Sprintf("  • %s", item))
	}
	return strings.Join(formatted, "\n")
}

// RenderBox renders content in a styled box
func RenderBox(title, content string) string {
	titleRendered := headerStyle.Render(title)
	contentRendered := boxStyle.Render(content)
	return fmt.Sprintf("%s\n%s", titleRendered, contentRendered)
}

// ShowReview displays a formatted code review
func ShowReview(review *ai.Review) {
	fmt.Println(headerStyle.Render("Code Review Results Completed"))

	if review.Summary != "" {
		fmt.Println(mutedStyle.Render("## Summary"))
		fmt.Println(boxStyle.Render(review.Summary))
	}

	if len(review.Issues) > 0 {
		fmt.Println(errorStyle.Render("## Issues"))
		for _, issue := range review.Issues {
			fmt.Printf("  %s [Severity: %s, Type: %s] %s\n", errorStyle.Render("•"), issue.Severity, issue.Type, issue.Description)
			if issue.Suggestion != "" {
				fmt.Printf("    %s Suggestion: %s\n", mutedStyle.Render("↳"), issue.Suggestion)
			}
		}
	}

	if len(review.Suggestions) > 0 {
		fmt.Println(infoStyle.Render("## Suggestions"))
		for _, suggestion := range review.Suggestions {
			fmt.Printf("  %s [Type: %s] %s\n", infoStyle.Render("•"), suggestion.Type, suggestion.Description)
			if suggestion.Example != "" {
				fmt.Printf("    %s Example: %s\n", mutedStyle.Render("↳"), suggestion.Example)
			}
		}
	}

	if len(review.SecurityRisks) > 0 {
		fmt.Println(errorStyle.Render("## Security Risks"))
		for _, risk := range review.SecurityRisks {
			fmt.Printf("  %s [Severity: %s] %s\n", errorStyle.Render("•"), risk.Severity, risk.Description)
			if risk.Mitigation != "" {
				fmt.Printf("    %s Mitigation: %s\n", mutedStyle.Render("↳"), risk.Mitigation)
			}
		}
	}

	if len(review.Performance) > 0 {
		fmt.Println(warningStyle.Render("## Performance Issues"))
		for _, perf := range review.Performance {
			fmt.Printf("  %s [Type: %s] %s (Impact: %s)\n", warningStyle.Render("•"), perf.Type, perf.Description, perf.Impact)
			if perf.Solution != "" {
				fmt.Printf("    %s Solution: %s\n", mutedStyle.Render("↳"), perf.Solution)
			}
		}
	}
}

// ShowPRDescription displays a formatted PR description
func ShowPRDescription(pr *PRDescription, platform string) {
	fmt.Println(headerStyle.Render("🚀 Generated PR/MR Description"))
	
	// Title
	fmt.Println(titleStyle.Render("Title:"))
	fmt.Println(boxStyle.Render(pr.Title))
	
	// Summary
	if pr.Summary != "" {
		fmt.Println(titleStyle.Render("Summary:"))
		fmt.Println(boxStyle.Render(pr.Summary))
	}
	
	// Changes
	if len(pr.Changes) > 0 {
		fmt.Println(titleStyle.Render("Changes:"))
		changesList := FormatList(pr.Changes)
		fmt.Println(boxStyle.Render(changesList))
	}
	
	// Issue Links
	if len(pr.IssueLinks) > 0 {
		fmt.Println(titleStyle.Render("Related Issues:"))
		issuesList := FormatList(pr.IssueLinks)
		fmt.Println(boxStyle.Render(issuesList))
	}
	
	// Testing Notes
	if pr.TestingNotes != "" {
		fmt.Println(titleStyle.Render("Testing:"))
		fmt.Println(boxStyle.Render(pr.TestingNotes))
	}
	
	// Checklist
	if len(pr.Checklist) > 0 {
		fmt.Println(titleStyle.Render("Checklist:"))
		checklistText := FormatChecklist(pr.Checklist, platform)
		fmt.Println(boxStyle.Render(checklistText))
	}
	
	// Breaking Changes
	if len(pr.BreakingChanges) > 0 {
		fmt.Println(errorStyle.Render("⚠️  Breaking Changes:"))
		breakingList := FormatList(pr.BreakingChanges)
		fmt.Println(boxStyle.Render(breakingList))
	}
	
	// Screenshots reminder
	if pr.Screenshots {
		fmt.Println(warningStyle.Render("📸 Don't forget to add screenshots of UI changes!"))
	}
	
	// Platform-specific formatting
	fmt.Println(mutedStyle.Render(fmt.Sprintf("\n📋 Formatted for %s", strings.Title(platform))))
	
	// Full markdown output
	fmt.Println(titleStyle.Render("Markdown Output:"))
	markdown := FormatPRMarkdown(pr, platform)
	fmt.Println(codeBlockStyle.Render(markdown))
}

// FormatChecklist formats checklist items for display
func FormatChecklist(items []ChecklistItem, platform string) string {
	var formatted []string
	
	// Different platforms use different checkbox syntax
	var checkboxEmpty, checkboxChecked string
	switch platform {
	case "gitlab":
		checkboxEmpty = "- [ ]"
		checkboxChecked = "- [x]"
	case "bitbucket":
		checkboxEmpty = "- [ ]"
		checkboxChecked = "- [x]"
	default: // github
		checkboxEmpty = "- [ ]"
		checkboxChecked = "- [x]"
	}
	
	for _, item := range items {
		if item.Checked {
			formatted = append(formatted, fmt.Sprintf("  %s %s", checkboxChecked, item.Text))
		} else {
			formatted = append(formatted, fmt.Sprintf("  %s %s", checkboxEmpty, item.Text))
		}
	}
	
	return strings.Join(formatted, "\n")
}

// FormatPRMarkdown generates the complete markdown for the PR description
func FormatPRMarkdown(pr *PRDescription, platform string) string {
	var markdown strings.Builder
	
	// Title is handled separately in PR creation
	
	// Summary
	if pr.Summary != "" {
		markdown.WriteString("## Summary\n\n")
		markdown.WriteString(pr.Summary)
		markdown.WriteString("\n\n")
	}
	
	// Changes
	if len(pr.Changes) > 0 {
		markdown.WriteString("## Changes\n\n")
		for _, change := range pr.Changes {
			markdown.WriteString(fmt.Sprintf("- %s\n", change))
		}
		markdown.WriteString("\n")
	}
	
	// Testing
	if pr.TestingNotes != "" {
		markdown.WriteString("## Testing\n\n")
		markdown.WriteString(pr.TestingNotes)
		markdown.WriteString("\n\n")
	}
	
	// Checklist
	if len(pr.Checklist) > 0 {
		markdown.WriteString("## Checklist\n\n")
		for _, item := range pr.Checklist {
			checkbox := "[ ]"
			if item.Checked {
				checkbox = "[x]"
			}
			markdown.WriteString(fmt.Sprintf("- %s %s\n", checkbox, item.Text))
		}
		markdown.WriteString("\n")
	}
	
	// Breaking Changes
	if len(pr.BreakingChanges) > 0 {
		markdown.WriteString("## ⚠️ Breaking Changes\n\n")
		for _, breaking := range pr.BreakingChanges {
			markdown.WriteString(fmt.Sprintf("- %s\n", breaking))
		}
		markdown.WriteString("\n")
	}
	
	// Issue Links (at the end for auto-linking)
	if len(pr.IssueLinks) > 0 {
		markdown.WriteString("## Related Issues\n\n")
		for _, link := range pr.IssueLinks {
			markdown.WriteString(fmt.Sprintf("%s\n", link))
		}
	}
	
	return strings.TrimSpace(markdown.String())
}

// PRDescription represents a generated PR description
type PRDescription struct {
	Title           string
	Summary         string
	Changes         []string
	IssueLinks      []string
	TestingNotes    string
	Checklist       []ChecklistItem
	BreakingChanges []string
	Screenshots     bool
	Platform        string
}

// ChecklistItem represents a checklist item in the PR
type ChecklistItem struct {
	Text    string
	Checked bool
}

 