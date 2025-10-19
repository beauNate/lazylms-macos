package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/Rugz007/lazylms/pkg/client"
	"github.com/Rugz007/lazylms/pkg/tui"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	config := client.DefaultClientConfig()

	app := &cli.App{
		Name:  "lazylms",
		Usage: "TUI client for LM Studio",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       config.Host,
				Usage:       "LM Studio host address",
				Destination: &config.Host,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       config.Port,
				Usage:       "LM Studio port",
				Destination: &config.Port,
			},
			&cli.StringFlag{
				Name:        "scheme",
				Value:       config.Scheme,
				Usage:       "Scheme (http or https)",
				Destination: &config.Scheme,
			},
			&cli.DurationFlag{
				Name:        "timeout",
				Value:       config.HTTPTimeout,
				Usage:       "HTTP request timeout",
				Destination: &config.HTTPTimeout,
			},
			&cli.IntFlag{
				Name:        "max-retries",
				Value:       config.MaxRetries,
				Usage:       "Maximum number of HTTP retries",
				Destination: &config.MaxRetries,
			},
		},
		Action: func(c *cli.Context) error {
			if err := client.ValidateClientConfig(config); err != nil {
				return fmt.Errorf("invalid configuration: %w", err)
			}
			return runApp(sugar, config)
		},
	}

	if err := app.Run(os.Args); err != nil {
		sugar.Fatalf("Application failed: %v", err)
	}
}

func runApp(sugar *zap.SugaredLogger, config client.ClientConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logChannel := make(chan string, config.LogChannelSize)

	lmsClient, err := client.NewClientWithConfig(ctx, config, logChannel)
	if err != nil {
		return fmt.Errorf("failed to create LM Studio client: %w", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	cleanup := func() error {
		var cleanupErr error
		if err := lmsClient.Cleanup(); err != nil {
			sugar.Warnf("Client cleanup encountered issues: %v", err)
			cleanupErr = err
		}
		return cleanupErr
	}

	go func() {
		<-signalChan
		if err := cleanup(); err != nil {
			sugar.Errorf("Cleanup failed on signal: %v", err)
		}
		cancel()
		os.Exit(0)
	}()

	model := tui.NewModel(lmsClient, logChannel)

	defer model.Cleanup()

	p := tea.NewProgram(model, tea.WithInputTTY(), tea.WithContext(ctx))

	_, err = p.Run()
	if cleanupErr := cleanup(); cleanupErr != nil {
		if err != nil {
			return fmt.Errorf("program error: %w, cleanup error: %v", err, cleanupErr)
		}
		return fmt.Errorf("cleanup error: %w", cleanupErr)
	}
	return err
}
