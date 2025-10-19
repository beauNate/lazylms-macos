package client

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	// ErrInvalidInput represents validation errors
	ErrInvalidInput = errors.New("invalid input")

	// Regular expressions for validation
	hostRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s='%s': %s", e.Field, e.Value, e.Message)
}

func ValidatePort(port string) error {
	if strings.TrimSpace(port) == "" {
		return ValidationError{
			Field:   "port",
			Value:   port,
			Message: "port cannot be empty",
		}
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return ValidationError{
			Field:   "port",
			Value:   port,
			Message: "port must be a valid number",
		}
	}

	if portNum < MinPortNumber || portNum > MaxPortNumber {
		return ValidationError{
			Field:   "port",
			Value:   port,
			Message: fmt.Sprintf("port must be between %d and %d", MinPortNumber, MaxPortNumber),
		}
	}

	return nil
}

func ValidateHost(host string) error {
	if strings.TrimSpace(host) == "" {
		return ValidationError{
			Field:   "host",
			Value:   host,
			Message: "host cannot be empty",
		}
	}

	if len(host) > MaxHostnameLength {
		return ValidationError{
			Field:   "host",
			Value:   host,
			Message: fmt.Sprintf("host exceeds maximum length of %d characters", MaxHostnameLength),
		}
	}

	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	if !hostRegex.MatchString(host) {
		return ValidationError{
			Field:   "host",
			Value:   host,
			Message: "host must be a valid hostname or IP address",
		}
	}

	return nil
}

func ValidateScheme(scheme string) error {
	if strings.TrimSpace(scheme) == "" {
		return ValidationError{
			Field:   "scheme",
			Value:   scheme,
			Message: "scheme cannot be empty",
		}
	}

	scheme = strings.ToLower(scheme)
	for _, valid := range ValidSchemes {
		if scheme == valid {
			return nil
		}
	}

	return ValidationError{
		Field:   "scheme",
		Value:   scheme,
		Message: fmt.Sprintf("scheme must be one of: %v", ValidSchemes),
	}
}

func ValidateURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return ValidationError{
			Field:   "url",
			Value:   urlStr,
			Message: "URL cannot be empty",
		}
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ValidationError{
			Field:   "url",
			Value:   urlStr,
			Message: fmt.Sprintf("invalid URL format: %v", err),
		}
	}

	if err := ValidateScheme(parsedURL.Scheme); err != nil {
		return err
	}

	if err := ValidateHost(parsedURL.Hostname()); err != nil {
		return err
	}

	if parsedURL.Port() != "" {
		if err := ValidatePort(parsedURL.Port()); err != nil {
			return err
		}
	}

	return nil
}

func ValidateSystemMessage(content string) error {
	if strings.TrimSpace(content) == "" {
		return ValidationError{
			Field:   "system_message",
			Value:   content,
			Message: "system message cannot be empty",
		}
	}

	if !utf8.ValidString(content) {
		return ValidationError{
			Field:   "system_message",
			Value:   content,
			Message: "system message must be valid UTF-8",
		}
	}

	if len(content) > MaxSystemMessageLength {
		return ValidationError{
			Field:   "system_message",
			Value:   content,
			Message: fmt.Sprintf("system message exceeds maximum length of %d characters", MaxSystemMessageLength),
		}
	}

	return nil
}

func ValidateChatMessage(content string) error {
	if strings.TrimSpace(content) == "" {
		return ValidationError{
			Field:   "chat_message",
			Value:   content,
			Message: "chat message cannot be empty",
		}
	}

	if !utf8.ValidString(content) {
		return ValidationError{
			Field:   "chat_message",
			Value:   content,
			Message: "chat message must be valid UTF-8",
		}
	}

	if len(content) > MaxChatMessageLength {
		return ValidationError{
			Field:   "chat_message",
			Value:   content,
			Message: fmt.Sprintf("chat message exceeds maximum length of %d characters", MaxChatMessageLength),
		}
	}

	return nil
}

func ValidateModelID(modelID string) error {
	if strings.TrimSpace(modelID) == "" {
		return ValidationError{
			Field:   "model_id",
			Value:   modelID,
			Message: "model ID cannot be empty",
		}
	}

	if len(modelID) > MaxModelIDLength {
		return ValidationError{
			Field:   "model_id",
			Value:   modelID,
			Message: fmt.Sprintf("model ID exceeds maximum length of %d characters", MaxModelIDLength),
		}
	}

	return nil
}

func ValidateClientConfig(config ClientConfig) error {
	if err := ValidateHost(config.Host); err != nil {
		return err
	}

	if err := ValidatePort(config.Port); err != nil {
		return err
	}

	if err := ValidateScheme(config.Scheme); err != nil {
		return err
	}

	if config.HTTPTimeout <= 0 {
		return ValidationError{
			Field:   "http_timeout",
			Value:   config.HTTPTimeout.String(),
			Message: "HTTP timeout must be positive",
		}
	}

	if config.MaxRetries < 0 {
		return ValidationError{
			Field:   "max_retries",
			Value:   strconv.Itoa(config.MaxRetries),
			Message: "max retries cannot be negative",
		}
	}

	if config.LogChannelSize <= 0 {
		return ValidationError{
			Field:   "log_channel_size",
			Value:   strconv.Itoa(config.LogChannelSize),
			Message: "log channel size must be positive",
		}
	}

	return nil
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)

	// Remove null bytes and other dangerous characters
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove common command injection patterns
	dangerousPatterns := []string{
		"&&", "||", ";", "|", "`", "$(", "${",
	}
	for _, pattern := range dangerousPatterns {
		input = strings.ReplaceAll(input, pattern, "")
	}

	// Filter control characters except newline and tab
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
