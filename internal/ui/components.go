package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
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