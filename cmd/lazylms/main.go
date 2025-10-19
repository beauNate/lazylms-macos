package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/beauNate/lazylms-macos/internal/api"
	"github.com/beauNate/lazylms-macos/internal/config"
	"github.com/beauNate/lazylms-macos/internal/security"
	"github.com/beauNate/lazylms-macos/internal/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse command-line flags
	var (
		host       = flag.String("host", "localhost", "LM Studio host")
		port       = flag.Int("port", 1234, "LM Studio port")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("lazylms-macos %s\ncommit: %s\nbuilt: %s\n", version, commit, date)
		os.Exit(0)
	}

	// Validate host and port
	if err := security.ValidateHost(*host); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid host: %v\n", err)
		os.Exit(1)
	}
	if err := security.ValidatePort(*port); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid port: %v\n", err)
		os.Exit(1)
	}

	// Initialize configuration
	cfg := config.New(*host, *port)

	// Initialize API client
	client := api.NewClient(cfg)

	// Initialize Bubble Tea model
	model := ui.NewModel(client, cfg)

	// Run the Bubble Tea program with Mac OS 26 Liquid Glass styling
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
