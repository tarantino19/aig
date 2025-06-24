# AI Git CLI Tool - Plan of Action

## Project Overview

**Name**: `ai-git` (or `aig` for short)  
**Language**: Go (100%)  
**Purpose**: AI-powered Git development assistant providing intelligent commit messages, summaries, and code reviews  
**UI Framework**: Charm libraries (Bubble Tea, Lipgloss, Glamour)

## Core Features

### 1. Git Commit Message Generator

- Analyzes staged changes using `git diff --cached`
- Generates conventional commit messages following best practices
- Supports different commit types (feat, fix, docs, style, refactor, test, chore)
- Interactive mode for editing generated messages

### 2. Git Commit Summary

- Summarizes recent commits in a branch or date range
- Generates release notes or changelog entries
- Groups commits by type and impact
- Markdown output support

### 3. AI Powered Code Review Assistant

- Reviews staged or committed changes
- Identifies potential issues, bugs, or improvements
- Suggests refactoring opportunities
- Security and performance considerations

## Technical Architecture

### Core Components

```
ai-git/
├── cmd/
│   └── aig/
│       └── main.go          # Entry point
├── internal/
│   ├── ai/
│   │   ├── client.go        # AI provider interface
│   │   └── gemini.go        # Gemini AI implementation
│   ├── git/
│   │   ├── diff.go          # Git diff operations
│   │   ├── commit.go        # Commit operations
│   │   └── log.go           # Git log parsing
│   ├── commands/
│   │   ├── commit.go        # Commit message generator
│   │   ├── summary.go       # Commit summary
│   │   └── review.go        # Code review
│   ├── config/
│   │   └── config.go        # Configuration management
│   └── ui/
│       ├── styles.go        # Lipgloss styles
│       └── components.go    # Bubble Tea components
├── pkg/
│   └── prompts/
│       └── templates.go     # AI prompt templates
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .goreleaser.yml         # For building releases
```

### Key Dependencies

```go
// Charm Libraries
github.com/charmbracelet/bubbles
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
github.com/charmbracelet/glamour

// CLI Framework
github.com/spf13/cobra
github.com/spf13/viper

// Git Operations
github.com/go-git/go-git/v5

// AI Client
github.com/google/generative-ai-go/genai

// Utilities
github.com/joho/godotenv
github.com/fatih/color
```

## Command Structure

### Main Commands

```bash
# Generate commit message for staged changes
aig commit
aig c                    # shorthand

# Generate commit summary
aig summary [flags]
aig s                    # shorthand

# Review code changes
aig review [flags]
aig r                    # shorthand

# Configuration
aig config set <key> <value>
aig config get <key>
aig config list

# Help
aig help
aig --version
```

### Command Options

#### Commit Command

```bash
aig commit [flags]
  -t, --type string      Commit type (feat|fix|docs|style|refactor|test|chore)
  -s, --scope string     Commit scope
  -i, --interactive      Interactive mode for editing
  -c, --conventional     Force conventional commit format
  -p, --push            Auto-push after commit
  --dry-run             Show what would be committed
```

#### Summary Command

```bash
aig summary [flags]
  -n, --number int       Number of commits to summarize (default 10)
  -b, --branch string    Target branch (default: current)
  -f, --from string      Start date/commit
  -t, --to string        End date/commit
  -o, --output string    Output format (text|markdown|json)
  -g, --group            Group by commit type
  --changelog           Generate changelog format
```

#### Review Command

```bash
aig review [flags]
  -s, --staged          Review staged changes only
  -c, --commit string   Review specific commit
  -r, --range string    Review commit range
  -f, --files string    Review specific files (glob pattern)
  -v, --verbose         Detailed review output
  --security           Focus on security issues
  --performance        Focus on performance issues
```

## Implementation Phases

### Phase 1: Foundation (Week 1)

- [ ] Project structure setup
- [ ] Basic CLI framework with Cobra
- [ ] Configuration management with Viper
- [ ] Git operations wrapper using go-git
- [ ] Basic UI components with Lipgloss

### Phase 2: AI Integration (Week 2)

- [ ] AI provider interface design
- [ ] Gemini AI client implementation
- [ ] Prompt template system optimized for Gemini
- [ ] API key management and security
- [ ] Error handling for Gemini API responses

### Phase 3: Commit Message Generator (Week 3)

- [ ] Git diff parsing and analysis
- [ ] Commit message generation logic
- [ ] Interactive editing with Bubble Tea
- [ ] Conventional commit format support
- [ ] Testing and refinement

### Phase 4: Commit Summary (Week 4)

- [ ] Git log parsing
- [ ] Summary generation algorithms
- [ ] Grouping and categorization
- [ ] Multiple output formats
- [ ] Changelog generation

### Phase 5: Code Review Assistant (Week 5-6)

- [ ] Code analysis framework
- [ ] Review prompt engineering
- [ ] Issue categorization (bugs, security, performance)
- [ ] Suggestion formatting
- [ ] Interactive review mode

### Phase 6: Polish & Release (Week 7)

- [ ] Comprehensive error handling
- [ ] Performance optimization
- [ ] Documentation
- [ ] Installation scripts
- [ ] CI/CD setup with GitHub Actions
- [ ] Release automation with GoReleaser

## Configuration

### Config File Location

- `~/.config/aig/config.yaml` (Linux/macOS)
- `%APPDATA%\aig\config.yaml` (Windows)

### Configuration Options

```yaml
# AI Provider Settings
ai:
 provider: gemini
 api_key: ${AIG_GEMINI_API_KEY} # Environment variable
 model: gemini-1.5-pro # or gemini-1.5-flash for faster responses
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
```

## Prompt Engineering

### Commit Message Template

```
Analyze the following git diff and generate a concise, conventional commit message.

Rules:
1. Use conventional commit format: <type>(<scope>): <subject>
2. Types: feat, fix, docs, style, refactor, test, chore
3. Subject line max 50 characters
4. Use present tense ("add" not "added")
5. No period at the end

Diff:
{diff}

Generate the commit message:
```

### Code Review Template

```
Review the following code changes for:
1. Potential bugs or errors
2. Security vulnerabilities
3. Performance issues
4. Code style and best practices
5. Suggestions for improvement

Focus areas: {focus_areas}

Code changes:
{diff}

Provide a structured review with specific line references where applicable.
```

## Error Handling

- Graceful degradation when AI service is unavailable
- Local fallback for basic operations
- Clear error messages with suggested actions
- Retry logic with exponential backoff
- Offline mode capabilities

## Security Considerations

- API keys stored securely (never in code)
- Environment variable support
- Keyring integration for secure storage
- No sensitive data in logs
- Git hooks for preventing API key commits

## Testing Strategy

- Unit tests for all core functions
- Integration tests for Git operations
- Mock AI providers for testing
- CLI command tests
- Prompt effectiveness testing
- Performance benchmarks

## Documentation

- Comprehensive README with examples
- Man pages for each command
- Interactive help with examples
- Video tutorials
- Contributing guidelines

## Future Enhancements

- Support for more AI providers (OpenAI, Anthropic, Llama)
- Git hooks integration
- Team collaboration features
- Custom prompt templates
- Plugin system
- Web UI dashboard
- IDE integrations
- Multi-language support
