package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tarantino19/aig/internal/ai"
	"github.com/tarantino19/aig/internal/config"
	"github.com/tarantino19/aig/internal/git"
	"github.com/tarantino19/aig/internal/ui"
)

var (
	prTargetBranch string
	prPlatform     string
	prTemplate     string
	prDraft        bool
	prInteractive  bool
	prCopyToClipboard bool
)

// NewPRCmd creates the PR command
func NewPRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pr",
		Aliases: []string{"pull-request", "merge-request", "mr"},
		Short:   "Generate AI-powered PR/MR description",
		Long: `Analyzes your branch changes against the target branch and generates
an intelligent PR/MR description with summary, changes, issue links, and checklists.`,
		RunE: runPR,
	}

	cmd.Flags().StringVarP(&prTargetBranch, "target", "t", "main", "Target branch for comparison")
	cmd.Flags().StringVarP(&prPlatform, "platform", "p", "github", "Platform (github|gitlab|bitbucket)")
	cmd.Flags().StringVar(&prTemplate, "template", "standard", "Template type (standard|minimal|detailed)")
	cmd.Flags().BoolVarP(&prDraft, "draft", "d", false, "Generate draft PR description")
	cmd.Flags().BoolVarP(&prInteractive, "interactive", "i", true, "Interactive mode for editing")
	cmd.Flags().BoolVarP(&prCopyToClipboard, "copy", "c", false, "Copy description to clipboard")

	return cmd
}

func runPR(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if API key is configured
	if cfg.AI.APIKey == "" || cfg.AI.APIKey == "your-gemini-api-key-here" || cfg.AI.APIKey == "your-openai-api-key-here" {
		ui.ShowError(fmt.Errorf("%s API key not configured", strings.Title(cfg.AI.Provider)))
		ui.ShowInfo("Please configure your API key first using 'aig config set ai.api_key YOUR_KEY'")
		return nil
	}

	// Get current branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch == prTargetBranch {
		return fmt.Errorf("current branch (%s) is the same as target branch (%s)", currentBranch, prTargetBranch)
	}

	ui.ShowInfo(fmt.Sprintf("🔍 Analyzing changes from %s to %s...", prTargetBranch, currentBranch))

	// Get branch diff
	diff, err := git.GetBranchDiff(prTargetBranch)
	if err != nil {
		return fmt.Errorf("failed to get branch diff: %w", err)
	}

	if diff == "" {
		ui.ShowWarning("No differences found between current branch and target branch")
		return nil
	}

	// Get commits in current branch that are not in target
	commits, err := git.GetCommits(git.CommitOptions{
		Branch: fmt.Sprintf("%s..%s", prTargetBranch, currentBranch),
		Number: 50, // Limit to last 50 commits
	})
	if err != nil {
		return fmt.Errorf("failed to get branch commits: %w", err)
	}

	ui.ShowInfo(fmt.Sprintf("📊 Found %d commits and analyzing diff...", len(commits)))

	// Extract issue numbers from branch name and commits
	issueNumbers := extractIssueNumbers(currentBranch, commits)

	// Create AI provider
	provider, err := ai.NewProvider(ai.ProviderConfig{
		Provider:    cfg.AI.Provider,
		APIKey:      cfg.AI.APIKey,
		Model:       cfg.AI.Model,
		Temperature: cfg.AI.Temperature,
		MaxTokens:   cfg.AI.MaxTokens,
	})
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}
	defer provider.Close()

	// Generate PR description
	ui.ShowInfo(fmt.Sprintf("🤖 Generating PR description with %s...", strings.Title(cfg.AI.Provider)))
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	prDescription, err := generatePRDescription(ctx, provider, PRAnalysis{
		CurrentBranch: currentBranch,
		TargetBranch:  prTargetBranch,
		Diff:          diff,
		Commits:       commits,
		IssueNumbers:  issueNumbers,
		Platform:      prPlatform,
		Template:      prTemplate,
		IsDraft:       prDraft,
	})
	if err != nil {
		return fmt.Errorf("failed to generate PR description: %w", err)
	}

	// Display the generated PR description
	ui.ShowPRDescription(prDescription, prPlatform)

	// Interactive editing if requested
	if prInteractive {
		fmt.Print("\n✏️  Edit PR description? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			// TODO: Implement interactive editing using Bubble Tea
			ui.ShowInfo("Interactive editing feature coming soon! For now, copy and edit manually.")
		}
	}

	// Copy to clipboard if requested
	if prCopyToClipboard {
		// TODO: Implement clipboard functionality
		ui.ShowInfo("📋 Clipboard functionality coming soon! For now, copy manually from above.")
	}

	ui.ShowSuccess("PR description generated successfully!")
	ui.ShowInfo(fmt.Sprintf("💡 Tip: Use 'aig pr --platform %s' to format for different platforms", 
		map[string]string{"github": "gitlab", "gitlab": "bitbucket", "bitbucket": "github"}[prPlatform]))

	return nil
}

// PRAnalysis contains all the data needed for PR description generation
type PRAnalysis struct {
	CurrentBranch string
	TargetBranch  string
	Diff          string
	Commits       []git.Commit
	IssueNumbers  []string
	Platform      string
	Template      string
	IsDraft       bool
}

// PRDescription represents a generated PR description
type PRDescription struct {
	Title         string
	Summary       string
	Changes       []string
	IssueLinks    []string
	TestingNotes  string
	Checklist     []ChecklistItem
	BreakingChanges []string
	Screenshots   bool
	Platform      string
}

// ChecklistItem represents a checklist item in the PR
type ChecklistItem struct {
	Text    string
	Checked bool
}

func generatePRDescription(ctx context.Context, provider ai.Provider, analysis PRAnalysis) (*PRDescription, error) {
	// For now, we'll use the existing AI interface. Later we can extend it for PR-specific generation
	commitOptions := ai.CommitOptions{
		Type:         "",
		Scope:        "",
		Conventional: false,
	}

	// Generate a comprehensive commit message that we'll transform into PR description
	commitMsg, err := provider.GenerateCommitMessage(ctx, analysis.Diff, commitOptions)
	if err != nil {
		return nil, err
	}

	// Transform the commit message into a PR description
	prDesc := &PRDescription{
		Platform: analysis.Platform,
	}

	// Extract title from commit subject or generate from branch name
	if commitMsg.Subject != "" {
		prDesc.Title = capitalizeFirst(commitMsg.Subject)
	} else {
		prDesc.Title = generateTitleFromBranch(analysis.CurrentBranch)
	}

	// Use commit body as base summary, or generate from commits
	if commitMsg.Body != "" {
		prDesc.Summary = commitMsg.Body
	} else {
		prDesc.Summary = generateSummaryFromCommits(analysis.Commits)
	}

	// Analyze changes and generate sections
	prDesc.Changes = analyzeChanges(analysis.Diff)
	prDesc.IssueLinks = formatIssueLinks(analysis.IssueNumbers, analysis.Platform)
	prDesc.TestingNotes = generateTestingNotes(analysis.Diff)
	prDesc.Checklist = generateChecklist(analysis.Diff, analysis.Commits)
	prDesc.BreakingChanges = detectBreakingChanges(analysis.Commits, analysis.Diff)
	prDesc.Screenshots = needsScreenshots(analysis.Diff)

	return prDesc, nil
}

func extractIssueNumbers(branchName string, commits []git.Commit) []string {
	issueSet := make(map[string]bool)
	
	// Extract from branch name
	branchIssues := extractIssuesFromText(branchName)
	for _, issue := range branchIssues {
		issueSet[issue] = true
	}
	
	// Extract from commit messages
	for _, commit := range commits {
		commitIssues := extractIssuesFromText(commit.Message)
		for _, issue := range commitIssues {
			issueSet[issue] = true
		}
	}
	
	// Convert set to slice
	var issues []string
	for issue := range issueSet {
		issues = append(issues, issue)
	}
	
	return issues
}

func extractIssuesFromText(text string) []string {
	// Match various issue number patterns: #123, fixes #123, closes #456, resolves #789
	patterns := []string{
		`#(\d+)`,
		`(?i)(?:fix|fixes|fixed|close|closes|closed|resolve|resolves|resolved)\s+#(\d+)`,
		`(?i)(?:fix|fixes|fixed|close|closes|closed|resolve|resolves|resolved)\s+(\d+)`,
	}
	
	var issues []string
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				issues = append(issues, match[1])
			}
		}
	}
	
	return issues
}

func generateTitleFromBranch(branchName string) string {
	// Remove prefixes and clean up branch name
	title := branchName
	prefixes := []string{"feature/", "feat/", "fix/", "bugfix/", "hotfix/", "chore/", "docs/"}
	
	for _, prefix := range prefixes {
		if strings.HasPrefix(title, prefix) {
			title = strings.TrimPrefix(title, prefix)
			break
		}
	}
	
	// Remove issue numbers and dates
	title = regexp.MustCompile(`\d{4,5}-?`).ReplaceAllString(title, "")
	title = regexp.MustCompile(`\d{8}`).ReplaceAllString(title, "")
	
	// Replace hyphens/underscores with spaces and capitalize
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.TrimSpace(title)
	
	return capitalizeFirst(title)
}

func generateSummaryFromCommits(commits []git.Commit) string {
	if len(commits) == 0 {
		return "No commits found in this branch."
	}
	
	if len(commits) == 1 {
		return fmt.Sprintf("This PR contains a single commit: %s", commits[0].Message)
	}
	
	return fmt.Sprintf("This PR contains %d commits with various changes and improvements.", len(commits))
}

func analyzeChanges(diff string) []string {
	changes := []string{}
	
	if strings.Contains(diff, "+++ /dev/null") {
		changes = append(changes, "🗑️ Removed files")
	}
	if strings.Contains(diff, "--- /dev/null") {
		changes = append(changes, "📄 Added new files")
	}
	if strings.Contains(diff, "package.json") || strings.Contains(diff, "go.mod") || strings.Contains(diff, "requirements.txt") {
		changes = append(changes, "📦 Updated dependencies")
	}
	if strings.Contains(diff, "_test.") || strings.Contains(diff, ".test.") || strings.Contains(diff, "spec.") {
		changes = append(changes, "🧪 Updated tests")
	}
	if strings.Contains(diff, "README") || strings.Contains(diff, ".md") {
		changes = append(changes, "📚 Updated documentation")
	}
	if strings.Contains(diff, ".css") || strings.Contains(diff, ".scss") || strings.Contains(diff, "style") {
		changes = append(changes, "🎨 Updated styles")
	}
	
	if len(changes) == 0 {
		changes = append(changes, "🔧 Modified existing functionality")
	}
	
	return changes
}

func formatIssueLinks(issueNumbers []string, platform string) []string {
	if len(issueNumbers) == 0 {
		return nil
	}
	
	links := make([]string, len(issueNumbers))
	for i, issue := range issueNumbers {
		switch platform {
		case "gitlab":
			links[i] = fmt.Sprintf("Closes #%s", issue)
		case "bitbucket":
			links[i] = fmt.Sprintf("Fixes #%s", issue)
		default: // github
			links[i] = fmt.Sprintf("Fixes #%s", issue)
		}
	}
	
	return links
}

func generateTestingNotes(diff string) string {
	if strings.Contains(diff, "_test.") || strings.Contains(diff, ".test.") {
		return "✅ Tests have been updated to cover the changes"
	}
	
	if strings.Contains(diff, "package.json") || strings.Contains(diff, "go.mod") {
		return "🔄 Run tests after installing new dependencies"
	}
	
	return "🧪 Manual testing recommended for the modified functionality"
}

func generateChecklist(diff string, commits []git.Commit) []ChecklistItem {
	checklist := []ChecklistItem{
		{Text: "Code follows project style guidelines", Checked: false},
		{Text: "Self-review of code has been performed", Checked: false},
	}
	
	if strings.Contains(diff, "_test.") || strings.Contains(diff, ".test.") {
		checklist = append(checklist, ChecklistItem{Text: "Tests pass locally", Checked: false})
	} else {
		checklist = append(checklist, ChecklistItem{Text: "Tests have been added/updated", Checked: false})
	}
	
	if strings.Contains(diff, "README") || strings.Contains(diff, ".md") {
		checklist = append(checklist, ChecklistItem{Text: "Documentation has been updated", Checked: true})
	} else {
		checklist = append(checklist, ChecklistItem{Text: "Documentation updated if needed", Checked: false})
	}
	
	return checklist
}

func detectBreakingChanges(commits []git.Commit, diff string) []string {
	var breaking []string
	
	for _, commit := range commits {
		if strings.Contains(strings.ToLower(commit.Message), "breaking change") ||
		   strings.Contains(strings.ToLower(commit.Message), "breaking:") {
			breaking = append(breaking, commit.Message)
		}
	}
	
	// Simple heuristics for detecting potential breaking changes
	if strings.Contains(diff, "public interface") || strings.Contains(diff, "export") {
		breaking = append(breaking, "Modified public interfaces - review for compatibility")
	}
	
	return breaking
}

func needsScreenshots(diff string) bool {
	// Check for UI-related changes
	uiPatterns := []string{".css", ".scss", ".html", ".jsx", ".tsx", ".vue", "component", "style", "ui/", "frontend/"}
	
	for _, pattern := range uiPatterns {
		if strings.Contains(strings.ToLower(diff), pattern) {
			return true
		}
	}
	
	return false
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
} 