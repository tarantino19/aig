package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")
	warningColor   = lipgloss.Color("#F59E0B")
	infoColor      = lipgloss.Color("#3B82F6")
	mutedColor     = lipgloss.Color("#6B7280")
	
	// Base styles
	baseStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1)
	
	// Title styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1)
	
	// Header styles
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(secondaryColor).
		MarginTop(1).
		MarginBottom(1)
	
	// Error style
	errorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)
	
	// Success style
	successStyle = lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true)
	
	// Info style
	infoStyle = lipgloss.NewStyle().
		Foreground(infoColor)
	
	// Warning style
	warningStyle = lipgloss.NewStyle().
		Foreground(warningColor)
	
	// Muted style
	mutedStyle = lipgloss.NewStyle().
		Foreground(mutedColor)
	
	// Box styles
	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)
	
	// Code block style
	codeBlockStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#1F2937")).
		Foreground(lipgloss.Color("#E5E7EB")).
		Padding(1).
		MarginTop(1).
		MarginBottom(1)
	
	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
		Foreground(primaryColor)
	
	// Prompt style
	promptStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)
)

// Commit type colors
var commitTypeColors = map[string]lipgloss.Color{
	"feat":     lipgloss.Color("#10B981"),
	"fix":      lipgloss.Color("#EF4444"),
	"docs":     lipgloss.Color("#3B82F6"),
	"style":    lipgloss.Color("#8B5CF6"),
	"refactor": lipgloss.Color("#F59E0B"),
	"test":     lipgloss.Color("#EC4899"),
	"chore":    lipgloss.Color("#6B7280"),
	"perf":     lipgloss.Color("#F97316"),
	"ci":       lipgloss.Color("#06B6D4"),
	"build":    lipgloss.Color("#84CC16"),
}

// GetCommitTypeStyle returns a style for a commit type
func GetCommitTypeStyle(commitType string) lipgloss.Style {
	color, ok := commitTypeColors[commitType]
	if !ok {
		color = mutedColor
	}
	return lipgloss.NewStyle().Foreground(color).Bold(true)
} 