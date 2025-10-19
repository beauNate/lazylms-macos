package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (c *Client) SendResponseStream(ctx context.Context, req ResponseRequest, callback func(string, string)) error {
	c.cancelled.Store(false)
	return c.sendResponseStreamWithRetry(ctx, req, callback, 3)
}

func (c *Client) sendResponseStreamWithRetry(ctx context.Context, req ResponseRequest, callback func(string, string), maxRetries int) error {
	req.Stream = true

	c.mu.Lock()
	if c.lastResponseID != nil {
		req.PreviousResponseID = c.lastResponseID
	}
	c.mu.Unlock()

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := DefaultRetryBaseDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(backoff)
			c.logger.Info("Retrying streaming request (attempt %d/%d) after %v", attempt+1, maxRetries+1, backoff)
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.GetFullURL()+"/v1/responses", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "text/event-stream")
		httpReq.Header.Set("Cache-Control", "no-cache")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("http request: %w", err)
			if attempt < maxRetries {
				continue
			}
			return lastErr
		}
		defer resp.Body.Close()

		if resp.StatusCode != HTTPStatusOK {
			// Don't retry on client errors (4xx), but retry on server errors (5xx)
			if resp.StatusCode >= HTTPStatusInternalError && attempt < maxRetries {
				lastErr = fmt.Errorf("api error: %s", resp.Status)
				continue
			}
			return fmt.Errorf("api error: %s", resp.Status)
		}

		err = c.parseSSEStream(resp.Body, callback)
		if err != nil {
			return fmt.Errorf("parse SSE stream: %w", err)
		}

		return nil
	}

	return lastErr
}

func (c *Client) parseSSEStream(body io.Reader, callback func(string, string)) error {
	scanner := bufio.NewScanner(body)
	var eventData strings.Builder
	var eventType string
	var responseID string

	for scanner.Scan() {
		if c.isCancelled() {
			return fmt.Errorf("request cancelled by user")
		}

		line := scanner.Text()

		if line == "" {
			if eventData.Len() > 0 {
				eventStr := eventData.String()
				eventData.Reset()

				if err := c.processSSEEvent(eventType, eventStr, &responseID, callback); err != nil {
					return fmt.Errorf("process SSE event: %w", err)
				}
				eventType = ""
			}
			continue
		}

		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimPrefix(line, "event: ")
		}

		// Accumulate event data
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			eventData.WriteString(data)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	if responseID != "" {
		c.mu.Lock()
		c.lastResponseID = &responseID
		c.mu.Unlock()
	}

	return nil
}

func (c *Client) processSSEEvent(eventType string, eventData string, responseID *string, callback func(string, string)) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(eventData), &data); err != nil {
		return fmt.Errorf("unmarshal event data: %w", err)
	}

	switch eventType {
	case "response.created":
		if id, ok := data["response_id"].(string); ok {
			*responseID = id
		}
	case "response.output_text.delta":
		if delta, ok := data["delta"].(string); ok {
			callback(delta, "output")
		}
	case "response.reasoning_text.delta":
		if delta, ok := data["delta"].(string); ok {
			callback(delta, "reasoning")
		}
	case "response.completed":
		// Streaming complete, no action needed
	}

	return nil
}
