# lazylms-macos

[![Release](https://img.shields.io/github/v/release/beauNate/lazylms-macos?style=flat-square)](https://github.com/beauNate/lazylms-macos/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beauNate/lazylms-macos?style=flat-square)](https://go.dev/)
[![Tests](https://img.shields.io/github/actions/workflow/status/beauNate/lazylms-macos/test.yml?branch=main&label=tests&style=flat-square)](https://github.com/beauNate/lazylms-macos/actions/workflows/test.yml)
[![Homebrew](https://img.shields.io/badge/homebrew-available-orange?style=flat-square)](https://github.com/beauNate/homebrew-lazylms-macos)

<p align="center">
  <img src="docs/demo.gif" alt="LazyLMS macOS Demo" width="100%"/>
</p>

A minimal TUI (Text User Interface) client for LM Studio. This is a personal hobby project that utilizes the LM Studio REST API and provides an elegant terminal interface for interacting with local language models.

> ⚠️ **Currently in BETA**: Supports core functionalities with active development ongoing.

## ✨ Features

- 🖥️ **Beautiful TUI**: Clean, minimal terminal interface built with Bubble Tea
- 🚀 **Fast & Lightweight**: Native Go performance with minimal resource footprint
- 🔌 **LM Studio Integration**: Direct integration with LM Studio's REST API
- 💬 **Interactive Chat**: Real-time streaming responses from local models
- 🎨 **Syntax Highlighting**: Code block rendering with language detection
- ⌨️ **Keyboard Shortcuts**: Vim-inspired navigation and controls
- 🔄 **Model Switching**: Easy model selection and management
- 📝 **Session History**: Persistent conversation history
- 🎯 **Cross-Platform**: Native support for macOS (Apple Silicon & Intel)

## 🚀 Quick Start

### Prerequisites

1. **LM Studio** installed and running ([Download here](https://lmstudio.ai/))
2. At least one model loaded in LM Studio
3. LM Studio server started (default: `http://localhost:1234`)

### Installation

#### Option 1: Homebrew (Recommended)

```bash
brew tap beauNate/lazylms-macos
brew install lazylms-macos
lazylms-macos
```

#### Option 2: DMG Installer

Download from [latest release](https://github.com/beauNate/lazylms-macos/releases/latest):

- **Apple Silicon (M1/M2/M3/M4)**: `lazylms-macos_v2.0.0_darwin_arm64.dmg`
- **Intel Macs**: `lazylms-macos_v2.0.0_darwin_amd64.dmg`

Double-click the DMG and drag to Applications folder.

#### Option 3: Binary Download

```bash
# Download for your architecture
tar -xzf lazylms-macos_darwin_*.tar.gz
chmod +x lazylms-macos
./lazylms-macos
```

## 📖 Usage

### Basic Commands

```bash
# Start with default settings
lazylms-macos

# Specify custom LM Studio server
lazylms-macos --server http://localhost:8080

# Show version
lazylms-macos --version

# Show help
lazylms-macos --help
```

### Interface Overview

- **Top Bar**: Current model and connection status
- **Chat Area**: Conversation history with syntax highlighting
- **Input Box**: Type your prompts here
- **Status Bar**: Keyboard shortcuts and help

## ⌨️ Keyboard Shortcuts

<details>
<summary><b>View all shortcuts</b></summary>

### Navigation
- `↑/↓` or `j/k` - Scroll chat history
- `Page Up/Down` - Fast scroll
- `Home/End` - Jump to start/end
- `Tab` - Switch focus (chat/input)

### Chat Controls
- `Enter` - Send message
- `Ctrl+C` - Cancel current generation
- `Ctrl+L` - Clear chat history
- `Ctrl+R` - Reload/refresh

### Model Management
- `Ctrl+M` - Open model selector
- `Ctrl+S` - Show model settings

### Application
- `Ctrl+Q` or `Esc` - Quit application
- `Ctrl+H` - Toggle help panel
- `?` - Show shortcuts overlay

### Text Editing (Input)
- `Ctrl+A` - Move to start of line
- `Ctrl+E` - Move to end of line
- `Ctrl+K` - Delete to end of line
- `Ctrl+U` - Delete entire line
- `Ctrl+W` - Delete word backwards

</details>

## 🔧 Troubleshooting

### Connection Issues

**Problem**: "Cannot connect to LM Studio"

- Ensure LM Studio is running
- Check the server is started (LM Studio > Developer > Start Server)
- Verify port is `1234` (default) or use `--server` flag
- Check firewall settings

### Model Not Loading

**Problem**: "No models available"

- Load a model in LM Studio first
- Ensure the model is fully downloaded
- Try restarting LM Studio server

### Slow Response Times

- Check model size vs available RAM
- Review LM Studio GPU acceleration settings
- Consider using a smaller/faster model
- Monitor CPU/GPU usage in Activity Monitor

### Display Issues

- Ensure terminal supports 256 colors
- Try resizing terminal window
- Update to latest version: `brew upgrade lazylms-macos`

## 🛠️ Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos

# Install dependencies
go mod download

# Build
go build -o lazylms-macos .

# Run
./lazylms-macos
```

### Project Structure

```
lazylms-macos/
├── cmd/           # Command-line interface
├── internal/      # Internal packages
│   ├── api/      # LM Studio API client
│   ├── ui/       # TUI components
│   └── config/   # Configuration management
├── docs/          # Documentation and assets
└── scripts/       # Build and release scripts
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/api
```

### Code Style

This project follows standard Go conventions:
- `gofmt` for formatting
- `golangci-lint` for linting
- Conventional Commits for commit messages

## 🗺️ Roadmap

- [ ] **v2.1.0**: Multi-model comparison view
- [ ] **v2.2.0**: Conversation export (JSON, Markdown)
- [ ] **v2.3.0**: Custom prompt templates
- [ ] **v2.4.0**: Plugin system
- [ ] **v3.0.0**: Cross-platform support (Linux, Windows)

### Completed
- [x] Basic chat interface
- [x] Streaming responses
- [x] Model selection
- [x] Syntax highlighting
- [x] Keyboard shortcuts
- [x] Homebrew distribution
- [x] DMG installer
- [x] CI/CD pipeline

## 🔄 CI/CD

This project uses GitHub Actions for continuous integration and delivery:

- **Testing**: Automated tests on every push and PR
- **Building**: Multi-architecture builds (arm64, amd64)
- **Releasing**: Automated GitHub Releases with GoReleaser
- **Distribution**: 
  - Homebrew tap updates
  - DMG creation and notarization
  - Binary archives
  - Docker images (if configured)

### Workflows

- `.github/workflows/test.yml` - Run tests and linting
- `.github/workflows/release.yml` - Build and publish releases
- `.github/workflows/homebrew.yml` - Update Homebrew formula

## 🤝 Contributing

Contributions are welcome! This is a hobby project, but PRs and issues are appreciated.

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Test additions or modifications
- `chore:` - Maintenance tasks

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [LM Studio](https://lmstudio.ai/) - For the excellent local LLM platform
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the TUI framework
- [Glamour](https://github.com/charmbracelet/glamour) - For Markdown rendering
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For terminal styling

## 🔗 Links

- [LM Studio Documentation](https://lmstudio.ai/docs)
- [LM Studio CLI Documentation](https://lmstudio.ai/docs/cli)
- [Issue Tracker](https://github.com/beauNate/lazylms-macos/issues)
- [Discussions](https://github.com/beauNate/lazylms-macos/discussions)

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/beauNate">beauNate</a>
</p>
