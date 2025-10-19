package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func (c *Client) SendMessageStreamWithModel(ctx context.Context, message string, modelID string, callback func(string, string)) error {
	return c.SendMessageStreamWithResponses(ctx, message, modelID, callback)
}

func (c *Client) ClearConversation() {
	c.conversation = make([]openai.ChatCompletionMessage, 0)
	c.ClearResponseHistory()
}

func (c *Client) SetSystemMessage(content string) error {
	content = SanitizeInput(content)

	if err := ValidateSystemMessage(content); err != nil {
		return fmt.Errorf("invalid system message: %w", err)
	}

	if len(c.conversation) > 0 && c.conversation[0].Role == openai.ChatMessageRoleSystem {
		c.conversation = c.conversation[1:]
	}

	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: content,
	}
	c.conversation = append([]openai.ChatCompletionMessage{systemMessage}, c.conversation...)
	c.logger.Info("System message updated: %s", strings.TrimSpace(content)[:min(50, len(strings.TrimSpace(content)))]+"...")
	return nil
}

func (c *Client) SendMessageStreamWithResponses(ctx context.Context, message string, modelID string, callback func(string, string)) error {
	if c.IsClosed() {
		return fmt.Errorf("client is closed")
	}

	message = SanitizeInput(message)

	if err := ValidateChatMessage(message); err != nil {
		return fmt.Errorf("invalid chat message: %w", err)
	}

	if modelID != "" {
		if err := ValidateModelID(modelID); err != nil {
			return fmt.Errorf("invalid model ID: %w", err)
		}
	}

	var activeModel string
	var err error

	if modelID != "" {
		activeModel, err = c.GetActiveModelByID(modelID)
	} else {
		activeModel, err = c.GetActiveModel()
	}

	if err != nil {
		c.logger.Error("Failed to get active model: %v", err)
		return fmt.Errorf("failed to get active model: %w", err)
	}

	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	}
	c.conversation = append(c.conversation, userMessage)

	c.TruncateConversation(MaxConversationLength)
	messages := make([]InputMessage, len(c.conversation))
	for i, msg := range c.conversation {
		messages[i] = InputMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	req := ResponseRequest{
		Model: activeModel,
		Input: messages,
		Store: false,
	}

	return c.SendResponseStream(ctx, req, callback)
}
