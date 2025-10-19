package client

import "time"

const (
	DefaultLogChannelSize    = 100
	DefaultStreamChannelSize = 1000
	DefaultMaxLogLines       = 1000

	DefaultHTTPTimeout    = 5 * time.Minute
	DefaultMaxRetries     = 3
	DefaultRetryBaseDelay = 1 * time.Second

	DefaultLMStudioPort   = "1234"
	DefaultLMStudioHost   = "localhost"
	DefaultLMStudioScheme = "http"

	MaxSystemMessageLength = 10000
	MaxChatMessageLength   = 50000
	MaxConversationLength  = 100

	DefaultTickInterval      = 5 * time.Second
	DefaultEstimateInterval  = 5 * time.Minute
	DefaultAnimationInterval = 80 * time.Millisecond
	DefaultLogCheckInterval  = 100 * time.Millisecond

	// Validation constants
	MinPortNumber     = 1
	MaxPortNumber     = 65535
	MaxHostnameLength = 253
	MaxModelIDLength  = 256

	OpenAIEndpointSuffix     = "/v1"
	DefaultEmbeddingModelKey = "text-embedding-nomic-embed-text-v1.5"
	EstimateSuccessMessage   = "Estimate: This model may be loaded based on your resource guardrails settings."

	HTTPStatusOK            = 200
	HTTPStatusBadRequest    = 400
	HTTPStatusInternalError = 500
)

var ValidSchemes = []string{"http", "https"}

type ClientConfig struct {
	Host           string
	Port           string
	Scheme         string
	HTTPTimeout    time.Duration
	MaxRetries     int
	LogChannelSize int
}

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Host:           DefaultLMStudioHost,
		Port:           DefaultLMStudioPort,
		Scheme:         DefaultLMStudioScheme,
		HTTPTimeout:    DefaultHTTPTimeout,
		MaxRetries:     DefaultMaxRetries,
		LogChannelSize: DefaultLogChannelSize,
	}
}

func (c ClientConfig) GetFullURL() string {
	return c.Scheme + "://" + c.Host + ":" + c.Port
}

func (c ClientConfig) GetAPIURL() string {
	return c.GetFullURL() + OpenAIEndpointSuffix
}
