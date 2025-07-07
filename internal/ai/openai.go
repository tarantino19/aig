package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/tarantino19/aig/pkg/prompts"
)

// OpenAIProvider implements the Provider interface using OpenAI's GPT models
type OpenAIProvider struct {
	client      *openai.Client
	model       string
	temperature float32
	maxTokens   int
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, modelName string, temperature float64, maxTokens int) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	
	client := openai.NewClient(apiKey)
	
	// Default to gpt-4o-mini if no model specified
	if modelName == "" {
		modelName = openai.GPT4oMini
	}
	
	return &OpenAIProvider{
		client:      client,
		model:       modelName,
		temperature: float32(temperature),
		maxTokens:   maxTokens,
	}, nil
}

// GenerateCommitMessage generates a commit message from a git diff
func (o *OpenAIProvider) GenerateCommitMessage(ctx context.Context, diff string, options CommitOptions) (*CommitMessage, error) {
	prompt := prompts.GetCommitMessagePrompt(diff, options.Type, options.Scope, options.Conventional)
	
	response, err := o.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}
	
	// Parse the commit message
	commitMsg := parseCommitMessage(response, options.Conventional)
	
	return commitMsg, nil
}

// GenerateSummary generates a summary of commits
func (o *OpenAIProvider) GenerateSummary(ctx context.Context, commits []Commit, options SummaryOptions) (*Summary, error) {
	// Convert ai.Commit to prompts.Commit
	promptCommits := make([]prompts.Commit, len(commits))
	for i, c := range commits {
		promptCommits[i] = prompts.Commit{
			Hash:    c.Hash,
			Author:  c.Author,
			Date:    c.Date,
			Message: c.Message,
		}
	}
	
	prompt := prompts.GetSummaryPrompt(promptCommits, options.GroupByType, options.Changelog)
	
	response, err := o.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	
	// Try to parse as JSON first for structured response
	var summary Summary
	if err := json.Unmarshal([]byte(response), &summary); err != nil {
		// Fallback to text parsing
		summary = parseSummaryText(response, commits, options)
	}
	
	return &summary, nil
}

// ReviewCode performs a code review on the given diff
func (o *OpenAIProvider) ReviewCode(ctx context.Context, diff string, options ReviewOptions) (*Review, error) {
	prompt := prompts.GetReviewPrompt(diff, options.FocusAreas, options.Security, options.Performance)
	
	response, err := o.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to review code: %w", err)
	}
	
	// Parse the review response
	review := parseReviewResponse(response, options)
	
	return review, nil
}

// GeneratePRDescription generates a PR description from branch analysis
func (o *OpenAIProvider) GeneratePRDescription(ctx context.Context, analysis PRAnalysis) (*PRDescriptionAI, error) {
	// Convert ai.Commit to prompts.Commit
	promptCommits := make([]prompts.Commit, len(analysis.Commits))
	for i, c := range analysis.Commits {
		promptCommits[i] = prompts.Commit{
			Hash:    c.Hash,
			Author:  c.Author,
			Date:    c.Date,
			Message: c.Message,
		}
	}
	
	prompt := prompts.GetPRDescriptionPrompt(
		analysis.CurrentBranch,
		analysis.TargetBranch,
		analysis.Diff,
		promptCommits,
		analysis.IssueNumbers,
		analysis.Platform,
	)
	
	response, err := o.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PR description: %w", err)
	}
	
	// Try to parse as JSON first
	var prDesc PRDescriptionAI
	if err := json.Unmarshal([]byte(response), &prDesc); err != nil {
		// Fallback to text parsing if JSON fails
		return parsePRDescriptionFromText(response), nil
	}
	
	return &prDesc, nil
}

// generateWithRetry implements exponential backoff retry logic for rate limiting
func (o *OpenAIProvider) generateWithRetry(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:       o.model,
		Temperature: o.temperature,
		MaxTokens:   o.maxTokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	
	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		// Check if it's specifically an OpenAI API error
		if apiErr, ok := err.(*openai.APIError); ok {
			// Only retry on actual rate limit errors (HTTP 429)
			if apiErr.HTTPStatusCode == 429 {
				return "", fmt.Errorf("rate limit exceeded. Please check your OpenAI API quota and billing at https://platform.openai.com/usage")
			}
		}
		
		// For other errors, return immediately with better error message
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}
	
	return resp.Choices[0].Message.Content, nil
}

// Close closes the OpenAI client (no-op for OpenAI client)
func (o *OpenAIProvider) Close() error {
	// OpenAI client doesn't need explicit closing
	return nil
}

// Note: parseSummaryText and parseReviewResponse are defined in gemini.go and shared 