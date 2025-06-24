package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tarantino19/aig/internal/config"
	"github.com/tarantino19/aig/internal/git"
	"github.com/tarantino19/aig/internal/ui"
)

var (
	reviewStaged      bool
	reviewCommit      string
	reviewRange       string
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
		// TODO: Implement range diff
		return fmt.Errorf("commit range review not yet implemented")
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

	// This will be implemented when we complete the AI integration
	fmt.Printf("Configuration loaded: %+v\n", cfg)
	fmt.Printf("Security focus: %v, Performance focus: %v\n", reviewSecurity, reviewPerformance)
	
	return fmt.Errorf("AI integration not yet implemented")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
} 