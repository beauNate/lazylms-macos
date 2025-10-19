package tui

import (
	"time"

	"github.com/Rugz007/lazylms/pkg/client"
)

const (
	// UI-specific constants that don't belong in client config
	ChatInputCharLimit   = 500
	SystemInputCharLimit = 1000

	// UI refresh intervals
	LogCheckInterval = 100 * time.Millisecond
)

// Use client configuration constants for shared values
var (
	MaxLogLines             = client.DefaultMaxLogLines
	StreamChannelBufferSize = client.DefaultStreamChannelSize
	EstimateTickInterval    = client.DefaultEstimateInterval
	AnimationTickInterval   = client.DefaultAnimationInterval
)

// UI Strings
const (
	ChatInputPlaceholder     = "Type your message and press Enter..."
	SystemInputPlaceholder   = "Enter system prompt..."
	NoModelsPlaceholder      = "No models loaded - load a model first"
	SelectModelPlaceholder   = "Select a model first (press Enter on loaded model)"
	ChatWithModelPlaceholder = "Chat with %s (Enter to send, arrows to scroll)"
	GeneratingPlaceholder    = "Generating response..."
	WelcomeMessage           = "  Press 4 to start chatting\n  Press Tab to navigate panels\n  Press h for help\n  Press Ctrl+S for system prompt\n\n  Ready to chat with your AI models!"
)
