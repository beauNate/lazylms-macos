package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/Rugz007/lazylms/pkg/client"
	"github.com/Rugz007/lazylms/pkg/tui/rendering"
)

// ResponseBuffer tracks segments during streaming
type ResponseBuffer struct {
	Segments []rendering.ContentSegment
}

func (rb *ResponseBuffer) AddSegment(text string, contentType string) {
	segType := rendering.ContentTypeOutput
	if contentType == "reasoning" {
		segType = rendering.ContentTypeReasoning
	}

	// Merge with last segment if same type
	if len(rb.Segments) > 0 && rb.Segments[len(rb.Segments)-1].Type == segType {
		rb.Segments[len(rb.Segments)-1].Text += text
	} else {
		rb.Segments = append(rb.Segments, rendering.ContentSegment{
			Text: text,
			Type: segType,
		})
	}
}

func (rb *ResponseBuffer) Reset() {
	rb.Segments = nil
}

// Model represents the main TUI model
type Model struct {
	client                  *client.Client
	width                   int
	height                  int
	currentView             string
	status                  bool
	chatMessages            []rendering.ChatMessage
	downloadedModels        []client.LMSDownloadedListItem
	loadedModels            []client.LMSLoadedListItem
	selectedModel           string // Currently selected model for chat
	explicitlySelectedModel string // Explicitly selected model via Enter key
	systemPrompt            string // System prompt for chat
	ctx                     context.Context
	cancel                  context.CancelFunc
	logChan                 chan string // Channel for receiving log messages
	loadedList              list.Model
	downloadedList          list.Model
	logsViewport            viewport.Model
	chatViewport            viewport.Model
	chatInput               textinput.Model
	systemInput             textinput.Model // Input for system prompt
	showHelp                bool
	showSystemPopup         bool            // Whether to show system prompt popup
	hasWelcomeMessage       bool            // Whether the welcome message is still displayed
	animationTime           time.Time       // Current time for animations
	streaming               bool            // Whether we're currently streaming a response
	currentResponse         *ResponseBuffer // Buffer for current streaming response with segments
	// streamChan is used for streaming response chunks from the API to the UI.
	//
	// Synchronization guarantees:
	// - Channel is buffered (size: StreamChannelBufferSize = 1000) to prevent blocking
	// - Ownership: Created and owned by Model, garbage collected on program exit
	// - Writer: Single goroutine in handleChatInputKeys writes chunks, errors, and completion signal
	// - Reader: streamSubscription() tea.Cmd reads chunks and converts them to tea.Msg
	// - Lifecycle: Created in NewModel(), garbage collected when Model is destroyed
	// - Thread-safety: Writes use select with default case to prevent blocking even if buffer is full
	//
	// Write protocol:
	//   1. Stream chunks as they arrive: streamChunkMsg{content, contentType}
	//   2. On error: "ERROR: <error message>" (as plain string for backward compat)
	//   3. On completion: "STREAM_COMPLETE" (as plain string for backward compat)
	streamChan              chan interface{}
	originalStreamingStatus string    // Original status of the model being streamed to
	lastEstimateTime        time.Time // Last time estimates were run
	lastLoadedModelIDs      []string  // IDs of loaded models from last estimate run
	firstLoad               bool      // Whether this is the first load
}

func NewModel(lmsClient *client.Client, logChannel chan string) Model {
	ctx, cancel := context.WithCancel(context.Background())

	streamChan := make(chan interface{}, StreamChannelBufferSize)

	downloadedDelegate := list.NewDefaultDelegate()
	loadedDelegate := list.NewDefaultDelegate()

	// Create custom keymap without Q quit and help hotkeys
	customKeyMap := list.DefaultKeyMap()
	customKeyMap.Quit.SetKeys()
	customKeyMap.ShowFullHelp.SetKeys()
	customKeyMap.CloseFullHelp.SetKeys()

	loadedList := list.New([]list.Item{}, loadedDelegate, 0, 0)
	loadedList.KeyMap = customKeyMap

	downloadedList := list.New([]list.Item{}, downloadedDelegate, 0, 0)
	downloadedList.KeyMap = customKeyMap

	// Initialize viewports
	logsViewport := viewport.New(0, 0)
	logsViewport.SetContent("")

	chatViewport := viewport.New(0, 0)
	chatViewport.SetContent("")

	// Initialize chat input
	chatInput := textinput.New()
	chatInput.Placeholder = ChatInputPlaceholder
	chatInput.CharLimit = ChatInputCharLimit
	chatInput.Width = 50

	// Initialize system prompt input
	systemInput := textinput.New()
	systemInput.Placeholder = SystemInputPlaceholder
	systemInput.CharLimit = SystemInputCharLimit
	systemInput.Width = 50

	// Check if any models are loaded and update placeholder accordingly
	if loadedModels, err := lmsClient.GetLoadedModels(); err == nil && len(loadedModels) == 0 {
		chatInput.Placeholder = NoModelsPlaceholder
	}

	// Create welcoming message
	chatViewport.SetContent(WelcomeMessage)

	return Model{
		client:             lmsClient,
		currentView:        "status",
		ctx:                ctx,
		cancel:             cancel,
		logChan:            logChannel,
		chatMessages:       []rendering.ChatMessage{},
		loadedList:         loadedList,
		downloadedList:     downloadedList,
		logsViewport:       logsViewport,
		chatViewport:       chatViewport,
		chatInput:          chatInput,
		systemInput:        systemInput,
		showHelp:           false,
		showSystemPopup:    false,
		hasWelcomeMessage:  true,
		streaming:          false,
		currentResponse:    &ResponseBuffer{},
		streamChan:         streamChan,
		firstLoad:          true,
		lastEstimateTime:   time.Time{},
		lastLoadedModelIDs: []string{},
	}
}

// Cleanup performs cleanup of TUI resources
func (m *Model) Cleanup() error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}

// IsClosed returns whether the model has been cleaned up
func (m *Model) IsClosed() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}
