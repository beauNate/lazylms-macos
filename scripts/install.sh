#!/bin/bash
# lazylms-macos installation script
# Installs the binary and sets up the environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}"
echo "╔═══════════════════════════════════════════════════╗"
echo "║   lazylms-macos Installation Script              ║"
echo "║   Mac OS 26 Liquid Glass UI                      ║"
echo "╚═══════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo -e "${RED}Error: This script is designed for macOS only${NC}"
    exit 1
fi

# Check for required dependencies
echo -e "${YELLOW}Checking dependencies...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed${NC}"

# Build the binary
echo -e "${YELLOW}Building lazylms-macos...${NC}"
make build

if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful${NC}"

# Install the binary
echo -e "${YELLOW}Installing to /usr/local/bin...${NC}"

if [ ! -w /usr/local/bin ]; then
    echo -e "${YELLOW}Requesting administrator privileges...${NC}"
    sudo make install
else
    make install
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}Installation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Installation complete${NC}"

# Create config directory
echo -e "${YELLOW}Setting up configuration directory...${NC}"
mkdir -p ~/.config/lazylms-macos

echo -e "${GREEN}✓ Configuration directory created${NC}"

# Verify installation
if command -v lazylms-macos &> /dev/null; then
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════════════════╗"
    echo "║   Installation Complete!                          ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo "Run 'lazylms-macos' to start the application"
    echo "Run 'lazylms-macos --help' for usage information"
else
    echo -e "${RED}Warning: Installation completed but binary not found in PATH${NC}"
    echo "You may need to add /usr/local/bin to your PATH"
fi

echo ""
echo -e "${YELLOW}Note: Make sure LM Studio is running before starting lazylms-macos${NC}"
echo "Download LM Studio from: https://lmstudio.ai/"
echo ""
