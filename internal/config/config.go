package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	AI      AIConfig      `mapstructure:"ai"`
	Git     GitConfig     `mapstructure:"git"`
	UI      UIConfig      `mapstructure:"ui"`
	Review  ReviewConfig  `mapstructure:"review"`
}

// AIConfig holds AI provider settings
type AIConfig struct {
	Provider    string  `mapstructure:"provider"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// GitConfig holds git-related settings
type GitConfig struct {
	AutoStage      bool   `mapstructure:"auto_stage"`
	DefaultBranch  string `mapstructure:"default_branch"`
	CommitTemplate string `mapstructure:"commit_template"`
}

// UIConfig holds UI settings
type UIConfig struct {
	Theme   string `mapstructure:"theme"`
	Emoji   bool   `mapstructure:"emoji"`
	Color   bool   `mapstructure:"color"`
	Spinner string `mapstructure:"spinner"`
}

// ReviewConfig holds code review settings
type ReviewConfig struct {
	IncludePatterns []string `mapstructure:"include_patterns"`
	ExcludePatterns []string `mapstructure:"exclude_patterns"`
	FocusAreas      []string `mapstructure:"focus_areas"`
}

// Load loads the configuration from file and environment
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Set default config path
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	
	// Initialize viper
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	
	// Set defaults
	setDefaults()
	
	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvPrefix("AIG")
	
	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// Create default config if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(configPath); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			// Re-read the created config
			if err := viper.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read created config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}
	
	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Override provider from environment if set
	if provider := os.Getenv("AIG_AI_PROVIDER"); provider != "" {
		cfg.AI.Provider = provider
	}
	
	// Override model from environment if set
	if model := os.Getenv("AIG_AI_MODEL"); model != "" {
		cfg.AI.Model = model
	}
	
	// Override API key from environment based on provider
	switch cfg.AI.Provider {
	case "openai":
		if apiKey := os.Getenv("AIG_OPENAI_API_KEY"); apiKey != "" {
			cfg.AI.APIKey = apiKey
		} else if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
			cfg.AI.APIKey = apiKey
		}
	case "gemini":
		if apiKey := os.Getenv("AIG_GEMINI_API_KEY"); apiKey != "" {
			cfg.AI.APIKey = apiKey
		} else if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
			cfg.AI.APIKey = apiKey
		}
	}
	
	return &cfg, nil
}

func setDefaults() {
	// AI defaults - now defaulting to OpenAI
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4o-mini")
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 2000)
	
	// Git defaults
	viper.SetDefault("git.auto_stage", false)
	viper.SetDefault("git.default_branch", "main")
	viper.SetDefault("git.commit_template", "conventional")
	
	// UI defaults
	viper.SetDefault("ui.theme", "dark")
	viper.SetDefault("ui.emoji", true)
	viper.SetDefault("ui.color", true)
	viper.SetDefault("ui.spinner", "dots")
	
	// Review defaults
	viper.SetDefault("review.include_patterns", []string{"*.go", "*.js", "*.py"})
	viper.SetDefault("review.exclude_patterns", []string{"*_test.go", "vendor/*"})
	viper.SetDefault("review.focus_areas", []string{"security", "performance", "best_practices"})
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(homeDir, ".config", "aig")
	return configDir, nil
}

func createDefaultConfig(configPath string) error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Write default config
	defaultConfig := `# AI Git Configuration

# AI Provider Settings
ai:
  provider: openai # openai or gemini
  api_key: ${AIG_OPENAI_API_KEY} # Environment variable
  model: gpt-4o-mini # OpenAI: gpt-4o-mini, gpt-4o, gpt-3.5-turbo | Gemini: gemini-1.5-pro, gemini-1.5-flash
  temperature: 0.7
  max_tokens: 2000

# Git Settings
git:
  auto_stage: false
  default_branch: main
  commit_template: conventional # or custom

# UI Settings
ui:
  theme: dark # or light, auto
  emoji: true
  color: true
  spinner: dots

# Review Settings
review:
  include_patterns:
    - '*.go'
    - '*.js'
    - '*.py'
  exclude_patterns:
    - '*_test.go'
    - 'vendor/*'
  focus_areas:
    - security
    - performance
    - best_practices
`
	
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}
	
	return nil
} 