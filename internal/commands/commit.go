package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tarantino19/aig/internal/ai"
	"github.com/tarantino19/aig/internal/config"
	"github.com/tarantino19/aig/internal/git"
	"github.com/tarantino19/aig/internal/ui"
)

var (
	commitType      string
	commitScope     string
	interactive     bool
	conventional    bool
	push           bool
	dryRun         bool
)

// NewCommitCmd creates the commit command
func NewCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commit",
		Aliases: []string{"c"},
		Short:   "Generate AI-powered commit message for staged changes",
		Long: `Analyzes your staged changes and generates an intelligent commit message
following best practices and conventional commit format.`,
		RunE: runCommit,
	}

	cmd.Flags().StringVarP(&commitType, "type", "t", "", "Commit type (feat|fix|docs|style|refactor|test|chore)")
	cmd.Flags().StringVarP(&commitScope, "scope", "s", "", "Commit scope")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Interactive mode for editing")
	cmd.Flags().BoolVarP(&conventional, "conventional", "c", true, "Force conventional commit format")
	cmd.Flags().BoolVarP(&push, "push", "p", false, "Auto-push after commit")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be committed")

	return cmd
}

func runCommit(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if API key is configured
	if cfg.AI.APIKey == "" || cfg.AI.APIKey == "your-gemini-api-key-here" || cfg.AI.APIKey == "your-openai-api-key-here" {
		ui.ShowError(fmt.Errorf("%s API key not configured", strings.Title(cfg.AI.Provider)))
		ui.ShowInfo("Please set your API key in one of these ways:")
		
		switch cfg.AI.Provider {
		case "openai":
			ui.ShowInfo("1. Set environment variable: export AIG_OPENAI_API_KEY=your-key")
			ui.ShowInfo("2. Set environment variable: export OPENAI_API_KEY=your-key")
			ui.ShowInfo("3. Edit .env file and add: AIG_OPENAI_API_KEY=your-key")
			ui.ShowInfo("4. Use: aig config set ai.api_key your-key")
			ui.ShowInfo("\nGet your API key from: https://platform.openai.com/api-keys")
		case "gemini":
			ui.ShowInfo("1. Set environment variable: export AIG_GEMINI_API_KEY=your-key")
			ui.ShowInfo("2. Set environment variable: export GEMINI_API_KEY=your-key")
			ui.ShowInfo("3. Edit .env file and add: AIG_GEMINI_API_KEY=your-key")
			ui.ShowInfo("4. Use: aig config set ai.api_key your-key")
			ui.ShowInfo("\nGet your API key from: https://makersuite.google.com/app/apikey")
		}
		return nil
	}

	// Get current branch name
	branchName, err := git.GetCurrentBranch()
	if err != nil {
		ui.ShowWarning("Could not get current branch name, proceeding without it.")
	}

	// Extract commit details from branch name
	extractedType, ticketNumber := git.ExtractCommitDetails(branchName)
	if extractedType != "" {
		commitType = extractedType
	}

	// Get staged changes
	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %w", err)
	}

	if diff == "" {
		ui.ShowWarning("No staged changes found. Stage your changes with 'git add' first")
		return nil
	}

	// Show what will be committed in dry-run mode
	if dryRun {
		ui.ShowDryRun(diff)
		return nil
	}

	// Create AI provider using the factory
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

	// Generate commit message using AI
	ui.ShowInfo(fmt.Sprintf("🤖 Analyzing staged changes with %s...", strings.Title(cfg.AI.Provider)))
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commitMsg, err := provider.GenerateCommitMessage(ctx, diff, ai.CommitOptions{
		Type:         commitType,
		Scope:        commitScope,
		Conventional: conventional,
	})
	if err != nil {
		// Check if it's specifically a rate limit or quota error
		errorStr := strings.ToLower(err.Error())
		if strings.Contains(errorStr, "rate limit exceeded") || 
		   strings.Contains(errorStr, "quota exceeded") ||
		   strings.Contains(errorStr, "insufficient_quota") ||
		   strings.Contains(errorStr, "429") {
			ui.ShowWarning("⚠️  API quota exceeded. Falling back to manual mode...")
			
			// Provide a fallback manual commit message
			fallbackMsg := generateFallbackCommitMessage(diff, ai.CommitOptions{
				Type:         commitType,
				Scope:        commitScope,
				Conventional: conventional,
			})
			
			ui.ShowInfo("📝 Generated fallback commit message:")
			ui.ShowCommitMessage(fallbackMsg.Type, fallbackMsg.Scope, fallbackMsg.Subject)
			
			if interactive {
				// Allow user to edit the message
				fmt.Print("\n✏️  Edit commit message? (y/N): ")
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
					fmt.Print("Enter commit message: ")
					var customMsg string
					fmt.Scanln(&customMsg)
					if customMsg != "" {
						fallbackMsg.FullMessage = customMsg
						fallbackMsg.Subject = customMsg
					}
				}
			}
			
			commitMsg = fallbackMsg
		} else {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}
	}

	if ticketNumber != "" {
		commitMsg.Subject = fmt.Sprintf("%s-%s", ticketNumber, commitMsg.Subject)
	}

	// Display the generated commit message
	ui.ShowCommitMessage(commitMsg.Type, commitMsg.Scope, commitMsg.Subject)
	
	if commitMsg.Body != "" {
		fmt.Printf("\nBody:\n%s\n", commitMsg.Body)
	}
	
	if commitMsg.Footer != "" {
		fmt.Printf("\nFooter:\n%s\n", commitMsg.Footer)
	}

	// In interactive mode, ask for confirmation
	if interactive {
		fmt.Print("\nUse this commit message? [Y/n]: ")
		var response string
		fmt.Scanln(&response)
		
		if response != "" && response != "y" && response != "Y" {
			ui.ShowInfo("Commit cancelled")
			return nil
		}
	}

	// Create the commit
	if err := git.CreateCommit(commitMsg.FullMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	ui.ShowSuccess("Commit created successfully!")

	// Auto-push if requested
	if push {
		ui.ShowInfo("Pushing to remote...")
		if err := git.Push(); err != nil {
			ui.ShowWarning(fmt.Sprintf("Failed to push: %v", err))
		} else {
			ui.ShowSuccess("Pushed to remote successfully!")
		}
	}

	return nil
}

// generateFallbackCommitMessage creates a simple commit message when AI is unavailable
func generateFallbackCommitMessage(diff string, options ai.CommitOptions) *ai.CommitMessage {
	// Simple analysis of the diff
	lines := strings.Split(diff, "\n")
	
	var addedFiles, modifiedFiles, deletedFiles []string
	var addedLines, deletedLines int
	
	for _, line := range lines {
		if strings.HasPrefix(line, "+++") && !strings.Contains(line, "/dev/null") {
			filename := strings.TrimPrefix(line, "+++ b/")
			if !contains(addedFiles, filename) && !contains(modifiedFiles, filename) {
				if strings.Contains(diff, "new file mode") {
					addedFiles = append(addedFiles, filename)
				} else {
					modifiedFiles = append(modifiedFiles, filename)
				}
			}
		} else if strings.HasPrefix(line, "---") && strings.Contains(line, "/dev/null") {
			filename := strings.TrimPrefix(line, "--- a/")
			if !contains(deletedFiles, filename) {
				deletedFiles = append(deletedFiles, filename)
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletedLines++
		}
	}
	
	// Generate commit message based on changes
	var commitType, subject string
	
	if options.Type != "" {
		commitType = options.Type
	} else {
		// Infer type from changes
		if len(addedFiles) > 0 {
			commitType = "feat"
		} else if len(deletedFiles) > 0 {
			commitType = "chore"
		} else {
			commitType = "fix"
		}
	}
	
	// Generate subject
	if len(addedFiles) > 0 {
		if len(addedFiles) == 1 {
			subject = fmt.Sprintf("add %s", filepath.Base(addedFiles[0]))
		} else {
			subject = fmt.Sprintf("add %d files", len(addedFiles))
		}
	} else if len(deletedFiles) > 0 {
		if len(deletedFiles) == 1 {
			subject = fmt.Sprintf("remove %s", filepath.Base(deletedFiles[0]))
		} else {
			subject = fmt.Sprintf("remove %d files", len(deletedFiles))
		}
	} else if len(modifiedFiles) > 0 {
		if len(modifiedFiles) == 1 {
			subject = fmt.Sprintf("update %s", filepath.Base(modifiedFiles[0]))
		} else {
			subject = fmt.Sprintf("update %d files", len(modifiedFiles))
		}
	} else {
		subject = "update code"
	}
	
	// Build full message
	var fullMessage string
	if options.Conventional {
		if options.Scope != "" {
			fullMessage = fmt.Sprintf("%s(%s): %s", commitType, options.Scope, subject)
		} else {
			fullMessage = fmt.Sprintf("%s: %s", commitType, subject)
		}
	} else {
		fullMessage = strings.Title(subject)
	}
	
	return &ai.CommitMessage{
		Type:        commitType,
		Scope:       options.Scope,
		Subject:     subject,
		FullMessage: fullMessage,
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 