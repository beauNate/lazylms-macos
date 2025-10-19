# lazylms-macos

[![Release](https://img.shields.io/github/v/release/beauNate/lazylms-macos?style=flat-square)](https://github.com/beauNate/lazylms-macos/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beauNate/lazylms-macos?style=flat-square)](https://go.dev/)
[![Tests](https://img.shields.io/github/actions/workflow/status/beauNate/lazylms-macos/test.yml?branch=main&label=tests&style=flat-square)](https://github.com/beauNate/lazylms-macos/actions/workflows/test.yml)

<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/12b94da0-588e-48bc-aab1-a55ec786cc2e" />

A minimal TUI client for LM Studio. Currently in **BETA** state and supports only basic functionalities.

This is a personal hobby project which utilizies LM Studio REST API and the `lms-cli`

Check it out LM Studio here: https://lmstudio.ai/ and https://lmstudio.ai/docs/cli

## Installation

### Homebrew (macOS)

```bash
brew tap beauNate/lazylms-macos
brew install lazylms-macos
```

### DMG Installer (macOS)

Download the appropriate DMG file for your Mac from the [latest release](https://github.com/beauNate/lazylms-macos/releases/latest):

- **For Apple Silicon (M1/M2/M3/M4)**: `lazylms-macos_v2.0.0_darwin_arm64.dmg`
- **For Intel Macs**: `lazylms-macos_v2.0.0_darwin_amd64.dmg`

Double-click the DMG file to open it, then drag the application to your Applications folder.

### Binary Download (macOS)

Download the binary for your architecture from the [latest release](https://github.com/beauNate/lazylms-macos/releases/latest):

- **For Apple Silicon**: `lazylms-macos_darwin_arm64.tar.gz`
- **For Intel Macs**: `lazylms-macos_darwin_amd64.tar.gz`

Extract and make it executable:

```bash
tar -xzf lazylms-macos_darwin_*.tar.gz
chmod +x lazylms-macos
./lazylms-macos
```

Optionally, move it to your PATH:

```bash
sudo mv lazylms-macos /usr/local/bin/
```

### Docker (via GitHub Container Registry)

```bash
# Pull and run the latest version
docker run --rm -it ghcr.io/beaunate/lazylms-macos:latest

# Or specific version
docker run --rm -it ghcr.io/beaunate/lazylms-macos:v2.0.0
```

## Usage

```bash
lazylms-macos
```

<img width="600" height="433" alt="image" src="https://github.com/user-attachments/assets/9bcdcdfd-92f7-4704-939d-c8a0a5aef26a" />

### Options

```
--host     LM Studio host (default: localhost)
--port     LM Studio port (default: 1234)
```

## Requirements

- Running LM Studio instance

## Development

### Building from source

```bash
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos
go build -o lazylms-macos ./cmd/lazylms
```

### Running tests

```bash
go test -v ./...
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Release Workflow**: Automatically builds macOS binaries (arm64/amd64), creates DMG installers, and publishes Docker images to GitHub Container Registry on new version tags
- **Test Workflow**: Runs tests on Ubuntu and macOS with multiple Go versions on every push and pull request
- **Homebrew Workflow**: Automatically updates the Homebrew tap formula when a new release is published

## License

MIT License - see [LICENSE](LICENSE) file for details
