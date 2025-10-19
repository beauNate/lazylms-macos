package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	config         ClientConfig
	logger         *Logger
	openAIClient   *openai.Client
	conversation   []openai.ChatCompletionMessage
	httpClient     *http.Client
	lastResponseID *string
	mu             sync.Mutex
	cancelled      atomic.Bool
	cancelChan     chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewClientWithConfig(ctx context.Context, config ClientConfig, logChannel chan string) (*Client, error) {
	if err := ValidateClientConfig(config); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)

	openaiConfig := openai.DefaultConfig("dummy-key") // LM Studio doesn't need real API key
	openaiConfig.BaseURL = config.GetAPIURL()         // LM Studio OpenAI-compatible endpoint

	openaiClient := openai.NewClientWithConfig(openaiConfig)

	// HTTP client for /v1/responses endpoint with proper connection pooling
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	httpClient := &http.Client{
		Timeout:   config.HTTPTimeout,
		Transport: transport,
	}

	if _, err := exec.LookPath("lms"); err != nil {
		cancel()
		return nil, fmt.Errorf("lms command not found: %w", err)
	}

	client := &Client{
		config:         config,
		logger:         NewLogger(logChannel),
		openAIClient:   openaiClient,
		conversation:   make([]openai.ChatCompletionMessage, 0),
		httpClient:     httpClient,
		lastResponseID: nil,
		cancelChan:     make(chan struct{}, 1),
		ctx:            ctx,
		cancel:         cancel,
	}

	return client, nil
}

func (c *Client) ClearResponseHistory() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastResponseID = nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() ClientConfig {
	return c.config
}

// GetLogger returns the client logger
func (c *Client) GetLogger() *Logger {
	return c.logger
}

// GetConversation returns a copy of the conversation history
func (c *Client) GetConversation() []openai.ChatCompletionMessage {
	return append([]openai.ChatCompletionMessage{}, c.conversation...)
}

// GetLastResponseID returns the last response ID
func (c *Client) GetLastResponseID() *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lastResponseID
}

func (c *Client) TruncateConversation(maxLength int) {
	if maxLength <= 0 {
		maxLength = MaxConversationLength
	}

	if len(c.conversation) > maxLength {
		c.conversation = c.conversation[len(c.conversation)-maxLength:]
		c.ClearResponseHistory()
		c.logger.Info("Truncated conversation history to %d messages", maxLength)
	}
}

func (c *Client) CancelRequest() {
	c.cancelled.Store(true)
	select {
	case c.cancelChan <- struct{}{}:
	default:
	}
}

func (c *Client) isCancelled() bool {
	return c.cancelled.Load()
}

func (c *Client) IsCancelled() bool {
	return c.cancelled.Load()
}

func (c *Client) AddAssistantMessage(content string) {
	if content == "" {
		return
	}
	assistantMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	}
	c.conversation = append(c.conversation, assistantMessage)
	c.logger.Info("Added assistant response to conversation (%d chars)", len(content))
}

func (c *Client) Cleanup() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *Client) Context() context.Context {
	return c.ctx
}

func (c *Client) IsClosed() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}
