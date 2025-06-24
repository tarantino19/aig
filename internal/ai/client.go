package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI providers
type Provider interface {
	// GenerateCommitMessage generates a commit message from a git diff
	GenerateCommitMessage(ctx context.Context, diff string, options CommitOptions) (*CommitMessage, error)
	
	// GenerateSummary generates a summary of commits
	GenerateSummary(ctx context.Context, commits []Commit, options SummaryOptions) (*Summary, error)
	
	// ReviewCode performs a code review on the given diff
	ReviewCode(ctx context.Context, diff string, options ReviewOptions) (*Review, error)
	
	// Close closes the provider connection
	Close() error
}

// ProviderConfig holds configuration for creating AI providers
type ProviderConfig struct {
	Provider    string
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
}

// NewProvider creates a new AI provider based on the configuration
func NewProvider(config ProviderConfig) (Provider, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIProvider(config.APIKey, config.Model, config.Temperature, config.MaxTokens)
	case "gemini":
		return NewGeminiProvider(config.APIKey, config.Model, config.Temperature, config.MaxTokens)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s. Supported providers: openai, gemini", config.Provider)
	}
}

// CommitOptions contains options for commit message generation
type CommitOptions struct {
	Type         string
	Scope        string
	Conventional bool
}

// CommitMessage represents a generated commit message
type CommitMessage struct {
	Type        string
	Scope       string
	Subject     string
	Body        string
	Footer      string
	FullMessage string
}

// Commit represents a git commit
type Commit struct {
	Hash    string
	Author  string
	Date    string
	Message string
}

// SummaryOptions contains options for summary generation
type SummaryOptions struct {
	GroupByType bool
	Format      string // text, markdown, json
	Changelog   bool
}

// Summary represents a generated summary
type Summary struct {
	Title       string
	Description string
	Groups      map[string][]CommitSummary
	Markdown    string
}

// CommitSummary represents a summarized commit
type CommitSummary struct {
	Hash    string
	Type    string
	Scope   string
	Subject string
}

// ReviewOptions contains options for code review
type ReviewOptions struct {
	FocusAreas []string
	Verbose    bool
	Security   bool
	Performance bool
}

// Review represents a code review result
type Review struct {
	Summary      string
	Issues       []Issue
	Suggestions  []Suggestion
	SecurityRisks []SecurityRisk
	Performance  []PerformanceIssue
}

// Issue represents a code issue
type Issue struct {
	Severity    string // high, medium, low
	Type        string // bug, style, logic
	File        string
	Line        int
	Description string
	Suggestion  string
}

// Suggestion represents a code improvement suggestion
type Suggestion struct {
	Type        string // refactor, optimization, clarity
	File        string
	Line        int
	Description string
	Example     string
}

// SecurityRisk represents a security issue
type SecurityRisk struct {
	Severity    string
	Type        string
	Description string
	Mitigation  string
}

// PerformanceIssue represents a performance issue
type PerformanceIssue struct {
	Type        string
	Description string
	Impact      string
	Solution    string
} 