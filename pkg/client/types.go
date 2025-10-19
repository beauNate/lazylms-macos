package client

// Status represents the status of a resource or operation
type Status int

const (
	StatusUnknown Status = iota
	StatusOn
	StatusOff
	StatusAvailable
	StatusUnavailable
)

// Model represents an LM Studio model
type Model struct {
	ID    string `json:"id"`
	State string `json:"state"` // "loaded" or "not-loaded" or "loading"
}

// ModelListResponse represents the response from the LM Studio API
type ModelListResponse struct {
	Data []Model `json:"data"`
}

type Quantization struct {
	Name string `json:"name"`
	Bits int    `json:"bits"`
}

type LMSDownloadedListItem struct {
	Type             string       `json:"type"`
	ModelKey         string       `json:"modelKey"`
	Format           string       `json:"format"`
	DisplayName      string       `json:"displayName"`
	Publisher        string       `json:"publisher"`
	Path             string       `json:"path"`
	SizeBytes        int64        `json:"sizeBytes"`
	Architecture     string       `json:"architecture"`
	Quantization     Quantization `json:"quantization"`
	MaxContextLength int          `json:"maxContextLength"`
	CanLoad          bool
}

type LMSLoadedListItem struct {
	Type              string       `json:"type"`
	ModelKey          string       `json:"modelKey"`
	Format            string       `json:"format"`
	DisplayName       string       `json:"displayName"`
	Publisher         string       `json:"publisher"`
	Path              string       `json:"path"`
	SizeBytes         int64        `json:"sizeBytes"`
	Architecture      string       `json:"architecture"`
	Quantization      Quantization `json:"quantization"`
	Identifier        string       `json:"identifier"`
	TtlMs             *int64       `json:"ttlMs"`
	LastUsedTime      int64        `json:"lastUsedTime"`
	Vision            bool         `json:"vision"`
	TrainedForToolUse bool         `json:"trainedForToolUse"`
	MaxContextLength  int          `json:"maxContextLength"`
	ContextLength     int          `json:"contextLength"`
	Status            string       `json:"status"`
	Queued            int          `json:"queued"`
}

// Response API structs for /v1/responses endpoint

type ReasoningConfig struct {
	Effort string `json:"effort"` // "low", "medium", "high"
}

type Tool struct {
	Type string `json:"type"`
}

type InputMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseRequest struct {
	Model              string           `json:"model"`
	Input              []InputMessage   `json:"input"`
	PreviousResponseID *string          `json:"previous_response_id,omitempty"`
	Reasoning          *ReasoningConfig `json:"reasoning,omitempty"`
	Tools              []Tool           `json:"tools,omitempty"`
	Stream             bool             `json:"stream,omitempty"`
	Store              bool             `json:"store"`
}

type OutputContent struct {
	Text string `json:"text,omitempty"`
}

type OutputItem struct {
	Type    string        `json:"type"`
	Content OutputContent `json:"content,omitempty"`
}

type ResponseResponse struct {
	ID     string       `json:"id"`
	Output []OutputItem `json:"output"`
}

// SSE Event structs for streaming
type SSEEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	ID    string `json:"id,omitempty"`
}

// Specific event data structs
type ResponseCreatedEvent struct {
	ResponseID string `json:"response_id"`
}

type ResponseOutputTextDeltaEvent struct {
	Delta string `json:"delta"`
}

type ResponseCompletedEvent struct {
	ResponseID string `json:"response_id"`
}
