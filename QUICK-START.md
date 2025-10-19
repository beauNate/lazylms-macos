# 🚀 Quick Start Guide

Get up and running with **lazylms-macos** in 5 minutes!

---

## Prerequisites

Before you begin, ensure you have:

1. ✅ **macOS 13.0+** (Ventura or later)
2. ✅ **Go 1.21+** installed ([Download](https://golang.org/dl/))
3. ✅ **LM Studio** installed and running ([Download](https://lmstudio.ai/))
4. ✅ At least one model loaded in LM Studio

---

## Step 1: Install LM Studio

If you haven't already:

1. Download LM Studio from [https://lmstudio.ai/](https://lmstudio.ai/)
2. Install and launch the application
3. Download and load a model (e.g., `Llama 3.3`, `Mistral`, `Phi`)
4. Start the local server:
   - Click **→ Local Server** in the left sidebar
   - Click **Start Server**
   - Note the port (default: `1234`)

---

## Step 2: Install lazylms-macos

### Option A: Automated Installation (Recommended)

```bash
# Clone the repository
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos

# Run the installation script
chmod +x scripts/install.sh
./scripts/install.sh
```

The script will:
- ✅ Check for Go installation
- ✅ Build the binary
- ✅ Install to `/usr/local/bin`
- ✅ Create config directory

### Option B: Manual Installation

```bash
# Clone the repository
git clone https://github.com/beauNate/lazylms-macos.git
cd lazylms-macos

# Build
make build

# Install (may require sudo)
make install
```

---

## Step 3: Run lazylms-macos

Simply run:

```bash
lazylms-macos
```

You should see the **Liquid Glass UI** appear:

```
╭─────────────────────────╮
│ 🔮 lazylms-macos      │
╰─────────────────────────╯

Connected to localhost:1234

Available Models:

  → llama-3.3-70b-instruct
    mistral-7b-instruct-v0.2
    phi-3-mini-4k-instruct

↑/↓: navigate • enter: select • q: quit
```

---

## Step 4: Navigate and Select a Model

### Keyboard Controls

| Key | Action |
|-----|--------|
| ↑ / `k` | Move up |
| ↓ / `j` | Move down |
| `Enter` | Select model |
| `Esc` | Go back |
| `q` / `Ctrl+C` | Quit |

### Selecting a Model

1. Use arrow keys or `j`/`k` to navigate
2. Press `Enter` to select a model
3. The model will be highlighted and ready for interaction

---

## Step 5: Start Chatting (Coming Soon)

> 🚧 **Note**: Full chat functionality is in active development. Current version focuses on model listing and selection.

Stay tuned for:
- 💬 Real-time chat interface
- 📝 Conversation history
- 💾 Save/load sessions
- ⚡ Streaming responses

---

## Troubleshooting

### Problem: "Connection refused"

**Solution**: Ensure LM Studio server is running

```bash
# Check if LM Studio is listening on port 1234
lsof -i :1234
```

### Problem: "No models found"

**Solution**: Load at least one model in LM Studio

1. Open LM Studio
2. Go to **Search** tab
3. Download a model
4. Go to **→ Local Server** tab
5. Load the model
6. Start the server

### Problem: "Permission denied" during installation

**Solution**: Use `sudo` for installation

```bash
sudo make install
```

### Problem: Binary not found after installation

**Solution**: Add `/usr/local/bin` to your PATH

```bash
# For Zsh (default on macOS)
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# For Bash
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bash_profile
source ~/.bash_profile
```

---

## Advanced Usage

### Custom Host/Port

If LM Studio is running on a different port:

```bash
lazylms-macos --port 8080
```

### Check Version

```bash
lazylms-macos --version
```

### Help

```bash
lazylms-macos --help
```

---

## What's Next?

- 📖 Read the full [README](README.md) for detailed documentation
- 🛠️ Explore the [Makefile](Makefile) for build commands
- 🔒 Learn about [Security features](README.md#-security)
- 🎨 Customize the [Liquid Glass UI](README.md#-liquid-glass-ui)
- 🐛 Report issues on [GitHub](https://github.com/beauNate/lazylms-macos/issues)
- 💬 Join discussions on [GitHub Discussions](https://github.com/beauNate/lazylms-macos/discussions)

---

## Quick Reference Card

```
╭────────────────────────────────────────────────────╮
│  lazylms-macos Quick Reference                       │
├────────────────────────────────────────────────────┤
│  Navigation:                                          │
│    ↑/↓ or j/k  : Move up/down                       │
│    Enter      : Select model                         │
│    Esc        : Go back                              │
│    q / Ctrl+C : Quit                                 │
├────────────────────────────────────────────────────┤
│  Commands:                                            │
│    lazylms-macos           : Run (default)           │
│    lazylms-macos --port N  : Custom port             │
│    lazylms-macos --version : Show version            │
│    lazylms-macos --help    : Show help               │
├────────────────────────────────────────────────────┤
│  Build Commands:                                      │
│    make build    : Build binary                      │
│    make install  : Install to /usr/local/bin         │
│    make clean    : Remove build artifacts            │
│    make test     : Run tests                         │
│    make app      : Create .app bundle                │
╰────────────────────────────────────────────────────╯
```

---

**Happy Chatting! 🚀**

For more help, see the [README](README.md) or [open an issue](https://github.com/beauNate/lazylms-macos/issues/new).
