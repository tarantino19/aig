package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tarantino19/aig/internal/ui"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage aig configuration",
		Long:  `View and modify aig configuration settings.`,
	}

	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigPathCmd())

	return cmd
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			// Set the value
			viper.Set(key, value)

			// Write config
			if err := viper.WriteConfig(); err != nil {
				// If config file doesn't exist, create it
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					configPath := viper.ConfigFileUsed()
					if configPath == "" {
						homeDir, _ := os.UserHomeDir()
						configPath = filepath.Join(homeDir, ".config", "aig", "config.yaml")
					}
					
					// Create directory if needed
					configDir := filepath.Dir(configPath)
					if err := os.MkdirAll(configDir, 0755); err != nil {
						return fmt.Errorf("failed to create config directory: %w", err)
					}
					
					// Write new config file
					if err := viper.WriteConfigAs(configPath); err != nil {
						return fmt.Errorf("failed to write config: %w", err)
					}
				} else {
					return fmt.Errorf("failed to write config: %w", err)
				}
			}

			ui.ShowSuccess(fmt.Sprintf("Set %s = %s", key, value))
			return nil
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := viper.Get(key)

			if value == nil {
				return fmt.Errorf("configuration key '%s' not found", key)
			}

			fmt.Printf("%s = %v\n", key, value)
			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings := viper.AllSettings()

			if len(settings) == 0 {
				ui.ShowInfo("No configuration values set")
				return nil
			}

			ui.ShowInfo("Current configuration:")
			printSettings(settings, "")
			return nil
		},
	}
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show configuration file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.ConfigFileUsed()
			if configFile == "" {
				homeDir, _ := os.UserHomeDir()
				configFile = filepath.Join(homeDir, ".config", "aig", "config.yaml")
				ui.ShowInfo(fmt.Sprintf("Default config path: %s", configFile))
			} else {
				ui.ShowInfo(fmt.Sprintf("Config file: %s", configFile))
			}
			return nil
		},
	}
}

// Helper function to recursively print settings
func printSettings(settings map[string]interface{}, prefix string) {
	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			printSettings(v, fullKey)
		default:
			fmt.Printf("  %s = %v\n", fullKey, value)
		}
	}
} 