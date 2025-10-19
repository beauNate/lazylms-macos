package security

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	// Maximum allowed prompt length (100KB)
	maxPromptLength = 100 * 1024
	// Maximum allowed model name length
	maxModelNameLength = 256
)

var (
	// Valid model name pattern: alphanumeric, dash, underscore, dot, colon, slash
	modelNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._:/-]+$`)
)

// ValidateHost validates that the host is not an external address
// and prevents Server-Side Request Forgery (SSRF) attacks
func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Only allow localhost and 127.0.0.1 for security
	allowedHosts := map[string]bool{
		"localhost":       true,
		"127.0.0.1":       true,
		"::1":             true,
		"0.0.0.0":         true,
	}

	if !allowedHosts[strings.ToLower(host)] {
		// Check if it's a local IP address
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return fmt.Errorf("host must be localhost or loopback address for security (got: %s)", host)
		}
	}

	return nil
}

// ValidatePort validates that the port is within valid range
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535 (got: %d)", port)
	}

	// Warn about privileged ports (< 1024) but don't block
	if port < 1024 {
		// This is just informational - we'll allow it but it's unusual
	}

	return nil
}

// ValidateModelName validates the model name to prevent injection attacks
func ValidateModelName(modelName string) error {
	if modelName == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	if len(modelName) > maxModelNameLength {
		return fmt.Errorf("model name exceeds maximum length of %d characters", maxModelNameLength)
	}

	if !modelNameRegex.MatchString(modelName) {
		return fmt.Errorf("model name contains invalid characters")
	}

	// Check for path traversal attempts
	if strings.Contains(modelName, "..") {
		return fmt.Errorf("model name contains invalid path traversal")
	}

	return nil
}

// ValidatePromptLength validates that the prompt is within acceptable size limits
func ValidatePromptLength(prompt string) error {
	if len(prompt) == 0 {
		return fmt.Errorf("prompt cannot be empty")
	}

	if len(prompt) > maxPromptLength {
		return fmt.Errorf("prompt exceeds maximum length of %d bytes", maxPromptLength)
	}

	return nil
}

// SanitizeInput performs basic sanitization on user input
// This is a defense-in-depth measure
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim excessive whitespace
	input = strings.TrimSpace(input)

	return input
}
