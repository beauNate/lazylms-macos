package tui

import (
	"time"

	"github.com/Rugz007/lazylms/pkg/client"
)

// Message types for the TUI
type tickMsg time.Time

type estimateTickMsg time.Time

type animationTickMsg time.Time

type emptyMsg struct{}

type statusMsg bool

type logMsg string

type modelsMsg struct {
	downloaded []client.LMSDownloadedListItem
	loaded     []client.LMSLoadedListItem
}

type estimatesMsg []client.LMSDownloadedListItem

type logListenerMsg struct{}

type streamChunkMsg struct {
	content     string
	contentType string // "output" or "reasoning"
}

type streamCompleteMsg struct{}
type nextViewMsg string
