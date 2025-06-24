# AIG - AI-Powered Git Assistant

[![Go](https://img.shields.io/badge/Go-1.24.2+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/tarantino19/aig)](https://github.com/tarantino19/aig/releases)

**AIG** (AI Git) is a powerful command-line tool that enhances your Git workflow with AI-powered features. Generate intelligent commit messages, code reviews, and project summaries using Google Gemini or OpenAI.

## ✨ Features

- 🤖 **AI-Powered Commit Messages**: Generate conventional, descriptive commit messages from staged changes
- 📊 **Code Reviews**: Get intelligent feedback on your staged changes
- 📝 **Project Summaries**: Create release notes and changelogs from commit history
- 🎯 **Conventional Commits**: Automatic adherence to conventional commit standards
- 🔧 **Flexible Configuration**: Support for multiple AI providers (Gemini, OpenAI)
- 🚀 **Cross-Platform**: Works on Linux, macOS, and Windows
- 💾 **Offline Fallback**: Graceful handling when AI services are unavailable

## 🚀 Quick Start

### Installation

#### Option 1: Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/tarantino19/aig/releases).

#### Option 2: Install Script

```bash
# Clone and install
git clone https://github.com/tarantino19/aig.git
cd aig
make build
./install.sh
```

#### Option 3: Build from Source

```bash
# Requires Go 1.24.2+
git clone https://github.com/tarantino19/aig.git
cd aig
make install
```

### Initial Setup

1. **Get an API Key**:

   - **Gemini**: Get your key from [Google AI Studio](https://makersuite.google.com/app/apikey)
   - **OpenAI**: Get your key from [OpenAI Platform](https://platform.openai.com/api-keys)

2. **Configure AIG**:

   ```bash
   # Set Gemini API key (recommended)
   aig config set ai.api_key YOUR_GEMINI_API_KEY

   # Or set OpenAI API key
   aig config set ai.provider openai
   aig config set ai.api_key YOUR_OPENAI_API_KEY
   ```

3. **Verify Installation**:
   ```bash
   aig --version
   aig --help
   ```

## 📖 Usage

### Generate Commit Messages

```bash
# Stage your changes
git add .

# Generate and create commit
aig commit

# Non-interactive mode
aig commit --interactive=false

# Specify commit type and scope
aig commit --type feat --scope auth

# Preview without committing
aig commit --dry-run

# Commit and push
aig commit --push
```

### Code Reviews

```bash
# Review staged changes
aig review --staged

# Review specific files
aig review file1.go file2.go

# Get detailed feedback
aig review --verbose
```

### Generate Summaries

```bash
# Summarize last 10 commits
aig summary

# Custom range
aig summary --number 20
aig summary --from "2024-01-01" --to "2024-02-01"

# Generate changelog
aig summary --changelog --output markdown
```

### Configuration Management

```bash
# View current configuration
aig config list

# Set values
aig config set ai.provider gemini
aig config set ai.model gemini-1.5-flash
aig config set ai.temperature 0.7

# Reset to defaults
aig config reset
```

## ⚙️ Configuration

AIG can be configured through multiple methods (in order of precedence):

1. **Command-line flags**
2. **Environment variables**
3. **Configuration file** (`~/.config/aig/config.yaml`)
4. **Default values**

### Environment Variables

```bash
# API Keys
export AIG_GEMINI_API_KEY="your-gemini-key"
export AIG_OPENAI_API_KEY="your-openai-key"

# AI Configuration
export AIG_AI_PROVIDER="gemini"
export AIG_AI_MODEL="gemini-1.5-flash"
export AIG_AI_TEMPERATURE="0.7"
export AIG_AI_MAX_TOKENS="2048"
```

### Configuration File

Create `~/.config/aig/config.yaml`:

```yaml
ai:
 provider: 'gemini'
 api_key: 'your-api-key'
 model: 'gemini-1.5-flash'
 temperature: 0.7
 max_tokens: 2048

git:
 auto_stage: false
 conventional_commits: true

ui:
 interactive: true
 colors: true
```

## 🔧 Development

### Prerequisites

- Go 1.24.2 or later
- Git

### Build from Source

```bash
# Clone repository
git clone https://github.com/tarantino19/aig.git
cd aig

# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Build for all platforms
make build-all
```

### Project Structure

```
aig/
├── cmd/aig/           # Main CLI application
├── internal/
│   ├── ai/           # AI provider implementations
│   ├── commands/     # CLI commands
│   ├── config/       # Configuration management
│   ├── git/          # Git operations
│   └── ui/           # User interface components
├── pkg/
│   └── prompts/      # AI prompt templates
├── build/            # Build artifacts
├── Makefile          # Build automation
└── install.sh        # Installation script
```

### Available Make Commands

```bash
make build      # Build the binary
make build-all  # Build for multiple platforms
make install    # Install system-wide
make test       # Run tests
make fmt        # Format code
make lint       # Lint code
make clean      # Clean build artifacts
make help       # Show all commands
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`aig commit` 😉)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📖 [Documentation](https://github.com/tarantino19/aig/wiki)
- 🐛 [Issue Tracker](https://github.com/tarantino19/aig/issues)
- 💬 [Discussions](https://github.com/tarantino19/aig/discussions)

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for interactive UI
- [go-git](https://github.com/go-git/go-git) for Git operations
- [Viper](https://github.com/spf13/viper) for configuration management

---

Made with ❤️ by the AIG team. Happy coding! 🚀
