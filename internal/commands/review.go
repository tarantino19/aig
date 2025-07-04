package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tarantino19/aig/internal/ai"
	"github.com/tarantino19/aig/internal/config"
	"github.com/tarantino19/aig/internal/git"
	"github.com/tarantino19/aig/internal/ui"
)

var (
	reviewStaged      bool
	reviewCommit      string
	reviewRange       string
	reviewBranch      string
	reviewFiles       string
	reviewVerbose     bool
	reviewSecurity    bool
	reviewPerformance bool
)

// NewReviewCmd creates the review command
func NewReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "review",
		Aliases: []string{"r"},
		Short:   "Get AI-powered code review for changes",
		Long: `Analyzes code changes and provides intelligent feedback on
potential issues, improvements, and best practices.`,
		RunE: runReview,
	}

	cmd.Flags().BoolVarP(&reviewStaged, "staged", "s", false, "Review staged changes only")
	cmd.Flags().StringVarP(&reviewCommit, "commit", "c", "", "Review specific commit")
	cmd.Flags().StringVarP(&reviewRange, "range", "r", "", "Review commit range")
	cmd.Flags().StringVarP(&reviewBranch, "branch", "b", "", "Review changes against a specific branch")
	cmd.Flags().StringVarP(&reviewFiles, "files", "f", "", "Review specific files (glob pattern)")
	cmd.Flags().BoolVarP(&reviewVerbose, "verbose", "v", false, "Detailed review output")
	cmd.Flags().BoolVar(&reviewSecurity, "security", false, "Focus on security issues")
	cmd.Flags().BoolVar(&reviewPerformance, "performance", false, "Focus on performance issues")

	return cmd
}

func runReview(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get diff based on flags
	var diff string
	switch {
	case reviewStaged:
		diff, err = git.GetStagedDiff()
		ui.ShowInfo("Reviewing staged changes...")
	case reviewCommit != "":
		diff, err = git.GetCommitDiff(reviewCommit)
		ui.ShowInfo(fmt.Sprintf("Reviewing commit %s...", reviewCommit))
	case reviewRange != "":
		diff, err = git.GetCommitRangeDiff(reviewRange)
		ui.ShowInfo(fmt.Sprintf("Reviewing commit range %s...", reviewRange))
	case reviewBranch != "":
		diff, err = git.GetBranchDiff(reviewBranch)
		ui.ShowInfo(fmt.Sprintf("Reviewing changes against branch %s...", reviewBranch))
	default:
		// Default to unstaged changes
		diff, err = git.GetDiff()
		ui.ShowInfo("Reviewing unstaged changes...")
	}

	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("no changes found to review")
	}

	// Show diff preview if verbose
	if reviewVerbose {
		ui.ShowDiff(truncateString(diff, 500))
	}

	// Initialize AI provider
	aiProvider, err := ai.NewProvider(ai.ProviderConfig{
		Provider:    cfg.AI.Provider,
		APIKey:      cfg.AI.APIKey,
		Model:       cfg.AI.Model,
		Temperature: cfg.AI.Temperature,
		MaxTokens:   cfg.AI.MaxTokens,
	})
	if err != nil {
		return fmt.Errorf("failed to create AI provider: %w", err)
	}
	defer func() {
		if err := aiProvider.Close(); err != nil {
			log.Printf("Error closing AI provider: %v", err)
		}
	}()

	ui.ShowInfo("Sending diff to AI for review...")

	reviewOptions := ai.ReviewOptions{
		FocusAreas:  cfg.Review.FocusAreas,
		Verbose:     reviewVerbose,
		Security:    reviewSecurity,
		Performance: reviewPerformance,
	}

	review, err := aiProvider.ReviewCode(cmd.Context(), diff, reviewOptions)
	if err != nil {
		return fmt.Errorf("failed to get code review: %w", err)
	}

	ui.ShowReview(review)

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
} 