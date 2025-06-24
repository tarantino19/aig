package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tarantino19/aig/internal/config"
	"github.com/tarantino19/aig/internal/git"
	"github.com/tarantino19/aig/internal/ui"
)

var (
	summaryNumber   int
	summaryBranch   string
	summaryFrom     string
	summaryTo       string
	summaryOutput   string
	summaryGroup    bool
	summaryChangelog bool
)

// NewSummaryCmd creates the summary command
func NewSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "summary",
		Aliases: []string{"s"},
		Short:   "Generate AI-powered summary of commits",
		Long: `Analyzes commit history and generates intelligent summaries,
release notes, or changelog entries.`,
		RunE: runSummary,
	}

	cmd.Flags().IntVarP(&summaryNumber, "number", "n", 10, "Number of commits to summarize")
	cmd.Flags().StringVarP(&summaryBranch, "branch", "b", "", "Target branch (default: current)")
	cmd.Flags().StringVarP(&summaryFrom, "from", "f", "", "Start date/commit")
	cmd.Flags().StringVarP(&summaryTo, "to", "t", "", "End date/commit")
	cmd.Flags().StringVarP(&summaryOutput, "output", "o", "text", "Output format (text|markdown|json)")
	cmd.Flags().BoolVarP(&summaryGroup, "group", "g", false, "Group by commit type")
	cmd.Flags().BoolVar(&summaryChangelog, "changelog", false, "Generate changelog format")

	return cmd
}

func runSummary(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get commits
	commits, err := git.GetCommits(git.CommitOptions{
		Number: summaryNumber,
		Branch: summaryBranch,
		From:   summaryFrom,
		To:     summaryTo,
	})
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("no commits found in the specified range")
	}

	ui.ShowInfo(fmt.Sprintf("Found %d commits to summarize", len(commits)))
	
	// This will be implemented when we complete the AI integration
	fmt.Printf("Configuration loaded: %+v\n", cfg)
	fmt.Printf("Output format: %s\n", summaryOutput)
	
	return fmt.Errorf("AI integration not yet implemented")
} 