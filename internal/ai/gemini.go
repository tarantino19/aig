package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/tarantino19/aig/pkg/prompts"
	"google.golang.org/api/option"
)

// GeminiProvider implements the Provider interface using Google's Gemini AI
type GeminiProvider struct {
	client      *genai.Client
	model       *genai.GenerativeModel
	temperature float32
	maxTokens   int32
}

// NewGeminiProvider creates a new Gemini AI provider
func NewGeminiProvider(apiKey, modelName string, temperature float64, maxTokens int) (*GeminiProvider, error) {
	ctx := context.Background()
	
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	
	model := client.GenerativeModel(modelName)
	model.SetTemperature(float32(temperature))
	model.SetMaxOutputTokens(int32(maxTokens))
	
	// Configure safety settings to be less restrictive for code review
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}
	
	return &GeminiProvider{
		client:      client,
		model:       model,
		temperature: float32(temperature),
		maxTokens:   int32(maxTokens),
	}, nil
}

// GenerateCommitMessage generates a commit message from a git diff
func (g *GeminiProvider) GenerateCommitMessage(ctx context.Context, diff string, options CommitOptions) (*CommitMessage, error) {
	prompt := prompts.GetCommitMessagePrompt(diff, options.Type, options.Scope, options.Conventional)
	
	resp, err := g.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate commit message: %w", err)
	}
	
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
	// Extract text from response
	text := extractTextFromResponse(resp)
	
	// Parse the commit message
	commitMsg := parseCommitMessage(text, options.Conventional)
	
	return commitMsg, nil
}

// GenerateSummary generates a summary of commits
func (g *GeminiProvider) GenerateSummary(ctx context.Context, commits []Commit, options SummaryOptions) (*Summary, error) {
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
	
	resp, err := g.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
	text := extractTextFromResponse(resp)
	
	// Try to parse as JSON first for structured response
	var summary Summary
	if err := json.Unmarshal([]byte(text), &summary); err != nil {
		// Fallback to text parsing
		summary = parseSummaryText(text, commits, options)
	}
	
	return &summary, nil
}

// ReviewCode performs a code review on the given diff
func (g *GeminiProvider) ReviewCode(ctx context.Context, diff string, options ReviewOptions) (*Review, error) {
	prompt := prompts.GetReviewPrompt(diff, options.FocusAreas, options.Security, options.Performance)
	
	resp, err := g.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to review code: %w", err)
	}
	
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
		text := extractTextFromResponse(resp)
	fmt.Println("--- Raw AI Response Start ---")
	fmt.Println(text)
	fmt.Println("--- Raw AI Response End ---")

	// Parse the review response
	review := parseReviewResponse(text, options)
	
	return review, nil
}

// GeneratePRDescription generates a PR description from branch analysis
func (g *GeminiProvider) GeneratePRDescription(ctx context.Context, analysis PRAnalysis) (*PRDescriptionAI, error) {
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
	
	resp, err := g.generateWithRetry(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PR description: %w", err)
	}
	
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
	text := extractTextFromResponse(resp)
	
	// Try to parse as JSON first
	var prDesc PRDescriptionAI
	if err := json.Unmarshal([]byte(text), &prDesc); err != nil {
		// Fallback to text parsing if JSON fails
		return parsePRDescriptionFromText(text), nil
	}
	
	return &prDesc, nil
}

// generateWithRetry implements exponential backoff retry logic for rate limiting
func (g *GeminiProvider) generateWithRetry(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	maxRetries := 3
	baseDelay := time.Second
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			// Check if it's a rate limiting error
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") {
				if attempt == maxRetries {
					return nil, fmt.Errorf("rate limit exceeded after %d attempts. Please try again later or check your Gemini API quota at https://ai.google.dev/gemini-api/docs/rate-limits", maxRetries+1)
				}
				
				// Calculate exponential backoff delay
				delay := time.Duration(math.Pow(2, float64(attempt))) * baseDelay
				fmt.Printf("Rate limit hit, retrying in %v... (attempt %d/%d)\n", delay, attempt+1, maxRetries+1)
				
				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			
			// For other errors, return immediately
			return nil, err
		}
		
		return resp, nil
	}
	
	return nil, fmt.Errorf("failed after %d retries", maxRetries+1)
}

// Close closes the Gemini client
func (g *GeminiProvider) Close() error {
	return g.client.Close()
}

// Helper functions

func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	var texts []string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					texts = append(texts, string(text))
				}
			}
		}
	}
	return strings.Join(texts, "\n")
}

func parseCommitMessage(text string, conventional bool) *CommitMessage {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) == 0 {
		return &CommitMessage{FullMessage: text}
	}
	
	commitMsg := &CommitMessage{}
	
	// Parse conventional commit format
	if conventional && strings.Contains(lines[0], ":") {
		parts := strings.SplitN(lines[0], ":", 2)
		if len(parts) == 2 {
			// Extract type and scope
			typeAndScope := strings.TrimSpace(parts[0])
			if strings.Contains(typeAndScope, "(") && strings.Contains(typeAndScope, ")") {
				typeParts := strings.SplitN(typeAndScope, "(", 2)
				commitMsg.Type = strings.TrimSpace(typeParts[0])
				commitMsg.Scope = strings.TrimSuffix(strings.TrimSpace(typeParts[1]), ")")
			} else {
				commitMsg.Type = typeAndScope
			}
			commitMsg.Subject = strings.TrimSpace(parts[1])
		}
	} else {
		// Non-conventional format
		commitMsg.Subject = lines[0]
	}
	
	// Extract body and footer
	if len(lines) > 2 {
		bodyLines := []string{}
		footerStart := -1
		
		for i := 2; i < len(lines); i++ {
			line := lines[i]
			// Check for footer patterns
			if strings.HasPrefix(line, "BREAKING CHANGE:") ||
				strings.HasPrefix(line, "Fixes #") ||
				strings.HasPrefix(line, "Closes #") ||
				strings.HasPrefix(line, "Resolves #") {
				footerStart = i
				break
			}
			bodyLines = append(bodyLines, line)
		}
		
		commitMsg.Body = strings.TrimSpace(strings.Join(bodyLines, "\n"))
		
		if footerStart > 0 {
			footerLines := lines[footerStart:]
			commitMsg.Footer = strings.TrimSpace(strings.Join(footerLines, "\n"))
		}
	}
	
	// Build full message
	fullParts := []string{}
	if commitMsg.Type != "" {
		if commitMsg.Scope != "" {
			fullParts = append(fullParts, fmt.Sprintf("%s(%s): %s", commitMsg.Type, commitMsg.Scope, commitMsg.Subject))
		} else {
			fullParts = append(fullParts, fmt.Sprintf("%s: %s", commitMsg.Type, commitMsg.Subject))
		}
	} else {
		fullParts = append(fullParts, commitMsg.Subject)
	}
	
	if commitMsg.Body != "" {
		fullParts = append(fullParts, "", commitMsg.Body)
	}
	
	if commitMsg.Footer != "" {
		fullParts = append(fullParts, "", commitMsg.Footer)
	}
	
	commitMsg.FullMessage = strings.Join(fullParts, "\n")
	
	return commitMsg
}

func parseSummaryText(text string, commits []Commit, options SummaryOptions) Summary {
	// Simple text parsing for summary
	lines := strings.Split(text, "\n")
	
	summary := Summary{
		Groups: make(map[string][]CommitSummary),
	}
	
	// Extract title (first non-empty line)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			summary.Title = line
			break
		}
	}
	
	// Set description as the full text for now
	summary.Description = text
	
	// If markdown format requested, format accordingly
	if options.Format == "markdown" {
		summary.Markdown = text
	}
	
	return summary
}

func parseReviewResponse(text string, options ReviewOptions) *Review {
	review := &Review{
		Summary:       "",
		Issues:        []Issue{},
		Suggestions:   []Suggestion{},
		SecurityRisks: []SecurityRisk{},
		Performance:   []PerformanceIssue{},
	}

	sections := make(map[string]string)
	currentSection := ""
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "## ") {
			currentSection = strings.ToLower(strings.TrimPrefix(line, "## "))
			if currentSection == "security risks" {
				currentSection = "security"
			} else if currentSection == "performance issues" {
				currentSection = "performance"
			}
			sections[currentSection] = "" // Initialize section content
		} else if currentSection != "" {
			sections[currentSection] += line + "\n"
		}
	}

	if summaryContent, ok := sections["summary"]; ok {
		review.Summary = strings.TrimSpace(summaryContent)
	}

	if issuesContent, ok := sections["issues"]; ok {
		for _, line := range strings.Split(issuesContent, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•")) {
				content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"), "•"))
				review.Issues = append(review.Issues, Issue{
					Severity:    "medium", // Default severity
					Type:        "general",  // Default type
					Description: content,
				})
			}
		}
	}

	if suggestionsContent, ok := sections["suggestions"]; ok {
		for _, line := range strings.Split(suggestionsContent, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•")) {
				content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"), "•"))
				review.Suggestions = append(review.Suggestions, Suggestion{
					Type:        "general", // Default type
					Description: content,
				})
			}
		}
	}

	if securityContent, ok := sections["security risks"]; ok {
		for _, line := range strings.Split(securityContent, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•")) {
				content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"), "•"))
				review.SecurityRisks = append(review.SecurityRisks, SecurityRisk{
					Severity:    "medium", // Default severity
					Description: content,
				})
			}
		}
	}

	if performanceContent, ok := sections["performance issues"]; ok {
		for _, line := range strings.Split(performanceContent, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•")) {
				content := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"), "•"))
				review.Performance = append(review.Performance, PerformanceIssue{
					Type:        "general", // Default type
					Description: content,
				})
			}
		}
	}

	return review
}

func parsePRDescriptionFromText(text string) *PRDescriptionAI {
	// Simple fallback parsing when JSON fails
	lines := strings.Split(text, "\n")
	
	prDesc := &PRDescriptionAI{
		Title:           "Generated PR Title",
		Summary:         "",
		Changes:         []string{},
		Testing:         "",
		BreakingChanges: []string{},
	}
	
	// Extract first non-empty line as title
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			prDesc.Title = line
			break
		}
	}
	
	// Use the full text as summary
	prDesc.Summary = strings.TrimSpace(text)
	
	// Simple change detection
	if strings.Contains(text, "add") || strings.Contains(text, "new") {
		prDesc.Changes = append(prDesc.Changes, "Added new functionality")
	}
	if strings.Contains(text, "fix") || strings.Contains(text, "bug") {
		prDesc.Changes = append(prDesc.Changes, "Fixed bugs")
	}
	if strings.Contains(text, "update") || strings.Contains(text, "modify") {
		prDesc.Changes = append(prDesc.Changes, "Updated existing features")
	}
	
	if len(prDesc.Changes) == 0 {
		prDesc.Changes = append(prDesc.Changes, "Made various improvements")
	}
	
	prDesc.Testing = "Please test the changes manually"
	
	return prDesc
} 