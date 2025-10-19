package rendering

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// MessageType represents the type of a chat message
type MessageType string

const (
	MessageTypeUser MessageType = "user"
	MessageTypeAI   MessageType = "ai"
)

// ContentType represents the type of content in a segment
type ContentType string

const (
	ContentTypeOutput    ContentType = "output"
	ContentTypeReasoning ContentType = "reasoning"
)

// ContentSegment represents a segment of content with a specific type
type ContentSegment struct {
	Text string
	Type ContentType
}

// WrapText wraps text to the given width
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	var lines []string
	var current string
	for _, word := range words {
		if len(current) > 0 && len(current)+len(word)+1 > width {
			lines = append(lines, current)
			current = word
		} else {
			if len(current) > 0 {
				current += " " + word
			} else {
				current = word
			}
		}
	}
	if len(current) > 0 {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}

// ChatMessage represents a structured chat message
type ChatMessage struct {
	Type     MessageType      `json:"type"`     // MessageTypeUser or MessageTypeAI
	Author   string           `json:"author"`   // "You" or model name
	Content  string           `json:"content"`  // The actual message content (for simple messages)
	Segments []ContentSegment `json:"segments"` // Segments with mixed content types
}

// renderMarkdown renders markdown content with proper styling
func renderMarkdown(content string, width int, logChan chan string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return content, nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return WrapText(content, width), err
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return WrapText(content, width), err
	}

	return strings.TrimRight(rendered, "\n"), nil
}

// RenderChatMessage renders a chat message, applying markdown only to AI responses
func RenderChatMessage(message ChatMessage, width int, logChan chan string) string {
	if message.Type == MessageTypeUser {
		// User messages don't get markdown rendering
		userPrefix := lipgloss.NewStyle().
			Foreground(styles.ColorGreen).
			Bold(true).
			Render(message.Author + ":")
		return userPrefix + " " + WrapText(message.Content, width-len(message.Author)-2) + "\n"
	}

	// Style the model name prefix
	styledPrefix := lipgloss.NewStyle().
		Foreground(styles.ColorPurple).
		Bold(true).
		Render(message.Author + ":")

	// AI messages - check if we have segments or simple content
	if len(message.Segments) > 0 {
		// Render mixed content with segments
		rendered := RenderMixedContent(message.Segments, width, logChan)
		return styledPrefix + rendered + "\n"
	}

	// Fallback to simple content rendering
	rendered, err := renderMarkdown(message.Content, width, logChan)
	if err != nil {
		rendered = WrapText(message.Content, width)
	}

	return styledPrefix + rendered + "\n"
}

// RenderAIMessage renders an AI message with model name styling and markdown
func RenderAIMessage(modelName, content string, width int, logChan chan string) string {
	// Render markdown for the content
	rendered, err := renderMarkdown(content, width, logChan)
	if err != nil {
		rendered = WrapText(content, width)
	}

	// Style the model name prefix with better visibility
	styledPrefix := lipgloss.NewStyle().
		Foreground(styles.ColorPurple).
		Bold(true).
		Render(modelName + ":")

	return styledPrefix + rendered + "\n"
}

// RenderMixedContent renders content with mixed output and reasoning segments
func RenderMixedContent(segments []ContentSegment, width int, logChan chan string) string {
	var result strings.Builder

	// Add newline after prefix if first segment is reasoning
	if len(segments) > 0 && segments[0].Type == ContentTypeReasoning {
		result.WriteString("\n")
	}

	for i, segment := range segments {
		switch segment.Type {
		case ContentTypeReasoning:
			// Render reasoning text with dimmed style and proper wrapping
			wrapped := WrapText(segment.Text, width-2)
			reasoningStyle := lipgloss.NewStyle().
				Foreground(styles.ColorReasoning).
				Italic(true)
			result.WriteString(reasoningStyle.Render(wrapped))
			// Add newline after reasoning if followed by output
			if i < len(segments)-1 && segments[i+1].Type == ContentTypeOutput {
				result.WriteString("\n")
			}
		case ContentTypeOutput:
			// Render output text with markdown
			rendered, err := renderMarkdown(segment.Text, width, logChan)
			if err != nil {
				rendered = WrapText(segment.Text, width)
			}
			result.WriteString(rendered)
		}
	}

	return result.String()
}
