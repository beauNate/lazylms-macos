# lazylms-macos

[![Release](https://img.shields.io/github/v/release/beauNate/lazylms-macos?style=flat-square)](https://github.com/beauNate/lazylms-macos/releases/latest) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT) [![Go Version](https://img.shields.io/github/go-mod/go-version/beauNate/lazylms-macos?style=flat-square)](https://go.dev/) [![Tests](https://img.shields.io/github/actions/workflow/status/beauNate/lazylms-macos/test.yml?branch=main&label=tests&style=flat-square)](https://github.com/beauNate/lazylms-macos/actions/workflows/test.yml) [![Homebrew](https://img.shields.io/badge/homebrew-available-orange?style=flat-square)](https://github.com/beauNate/homebrew-lazylms-macos)

<p align="center">
<img src="docs/demo.gif" alt="LazyLMS macOS Demo" width="100%"/>
</p>

> A minimal TUI (Text User Interface) client for LM Studio. This is a personal hobby project that utilizes the LM Studio REST API and provides an elegant terminal interface for interacting with local language models.
>
> ⚠️ **Currently in BETA**: Supports core functionalities with active development ongoing.

**Attribution**: This project was originally inspired by and adopted from the excellent work in [Rugz007/lazylms](https://github.com/Rugz007/lazylms), with gratitude to its authors and contributors. We've forked and evolved it for a Mac‑native TUI experience while keeping the spirit of the original.

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

#### Option 3: Go Install

```bash
go install github.com/beauNate/lazylms-macos@latest
```

#### Option 4: Build from Source

```bash
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos
make build
./bin/lazylms-macos
```

## 💻 Usage

### Starting the Application

```bash
lazylms-macos
```

### Keyboard Shortcuts

#### Global

- `Ctrl+C` / `q` - Quit application
- `Ctrl+N` - New chat session
- `Ctrl+S` - Select model
- `Ctrl+H` - View command history
- `Tab` - Cycle through UI elements

#### Chat View

- `Enter` - Send message
- `Esc` - Clear input / cancel operation
- `↑` / `↓` - Navigate chat history
- `Ctrl+L` - Clear screen
- `j` / `k` - Scroll messages (Vim-style)

#### Model Selection

- `↑` / `↓` or `j` / `k` - Navigate model list
- `Enter` - Select model
- `Esc` - Cancel

## ⚙️ Configuration

### Default Configuration Path

```
~/.config/lazylms-macos/config.yaml
```

### Configuration Options

```yaml
lm_studio:
  base_url: "http://localhost:1234"
  timeout: 30s

ui:
  theme: "default"
  show_timestamps: true
  max_history_lines: 1000

api:
  temperature: 0.7
  max_tokens: 2048
  stream: true
```

### Environment Variables

- `LAZYLMS_LM_STUDIO_URL` - Override LM Studio base URL
- `LAZYLMS_CONFIG_PATH` - Custom configuration file path

## 🛠️ Development

### Prerequisites

- Go 1.22 or higher
- Make (optional, for build automation)

### Building

```bash
# Clone the repository
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos

# Install dependencies
go mod download

# Build
make build

# Or without make
go build -o bin/lazylms-macos .
```

### Running Tests

```bash
make test

# Or without make
go test ./...
```

### Linting

```bash
make lint

# Or without make
golangci-lint run
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

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

## 🙏 Acknowledgments & Attribution

**Original Project**: This project is a fork and evolution of [Rugz007/lazylms](https://github.com/Rugz007/lazylms), an excellent terminal-based LLM client. We are deeply grateful to [Rugz007](https://github.com/Rugz007) and all contributors to the original project for laying the foundation and inspiration for this work. Their vision of a simple, elegant TUI for local language models made this macOS-focused adaptation possible. Thank you for your creativity, dedication, and for sharing your work with the open-source community! 🎉

We also want to thank:

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
  Made with ❤️ by [beauNate](https://github.com/beauNate)
</p>
