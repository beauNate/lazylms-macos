package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Allowed LMS commands (allowlist)
var allowedCommands = map[string]bool{
	"status": true,
	"load":   true,
	"unload": true,
	"ls":     true,
	"ps":     true,
}

func (c *Client) RunLMSCommand(command []string) (string, error) {
	if len(command) == 0 {
		return "", ValidationError{
			Field:   "command",
			Value:   "",
			Message: "command cannot be empty",
		}
	}

	// Validate command against allowlist
	if !allowedCommands[command[0]] {
		return "", ValidationError{
			Field:   "command",
			Value:   command[0],
			Message: fmt.Sprintf("command not allowed: %s", command[0]),
		}
	}

	// Validate all arguments don't contain dangerous characters
	for i, arg := range command {
		if strings.ContainsAny(arg, "&|;`$\x00") {
			return "", ValidationError{
				Field:   fmt.Sprintf("command[%d]", i),
				Value:   arg,
				Message: "argument contains dangerous characters",
			}
		}
	}

	if c.IsClosed() {
		return "", fmt.Errorf("client is closed")
	}
	if _, err := exec.LookPath("lms"); err != nil {
		c.logger.Error("lms command not found in PATH: %v", err)
		return "", fmt.Errorf("lms command not found: %w", err)
	}

	cmd := exec.Command("lms", command...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		if errOut.Len() > 0 {
			return "", fmt.Errorf("Failed to run lms command: %w, stderr: %s", err, errOut.String())
		}
		return "", fmt.Errorf("Failed to run lms command: %w", err)
	}

	output := out.String()
	if len(strings.TrimSpace(output)) == 0 && errOut.Len() > 0 {
		output = errOut.String()
	}

	return strings.TrimSpace(output), nil
}

func (c *Client) GetStatus() (Status, error) {
	cmd := []string{"status"}
	output, err := c.RunLMSCommand(cmd)

	if err != nil {
		return StatusUnknown, fmt.Errorf("Failed to run status %w", err)
	}

	if strings.Contains(output, "ON") {
		return StatusOn, nil
	}
	return StatusOff, nil
}

func (c *Client) GetModelEstimate(identifier string) (Status, error) {
	cmd := []string{"load", "--exact", "--estimate-only", identifier}
	output, err := c.RunLMSCommand(cmd)
	if err != nil {
		return StatusUnknown, fmt.Errorf("Failed to run lms load --estimate-only for %s: %w", identifier, err)
	}
	if strings.Contains(output, EstimateSuccessMessage) {
		return StatusAvailable, nil
	}
	return StatusUnavailable, nil
}

func (c *Client) GetDownloadedModelsWithoutEstimates() ([]LMSDownloadedListItem, error) {
	cmd := []string{"ls", "--json"}
	output, err := c.RunLMSCommand(cmd)

	if err != nil {
		return []LMSDownloadedListItem{}, fmt.Errorf("Failed to run lms ls: %w", err)
	}

	var downloadedList []LMSDownloadedListItem
	if err := json.Unmarshal([]byte(output), &downloadedList); err != nil {
		return []LMSDownloadedListItem{}, fmt.Errorf("Failed to parse JSON output: %w", err)
	}

	var filteredList []LMSDownloadedListItem
	for _, item := range downloadedList {
		if item.ModelKey != DefaultEmbeddingModelKey {
			filteredList = append(filteredList, item)
		}
	}

	for i := range filteredList {
		filteredList[i].CanLoad = true
	}

	return filteredList, nil
}

func (c *Client) GetDownloadedModels() ([]LMSDownloadedListItem, error) {
	if c.IsClosed() {
		return nil, fmt.Errorf("client is closed")
	}
	cmd := []string{"ls", "--json"}
	output, err := c.RunLMSCommand(cmd)

	if err != nil {
		return []LMSDownloadedListItem{}, fmt.Errorf("Failed to run lms ls: %w", err)
	}

	var downloadedList []LMSDownloadedListItem
	if err := json.Unmarshal([]byte(output), &downloadedList); err != nil {
		return []LMSDownloadedListItem{}, fmt.Errorf("Failed to parse JSON output: %w", err)
	}

	var filteredList []LMSDownloadedListItem
	for _, item := range downloadedList {
		if item.ModelKey != DefaultEmbeddingModelKey {
			filteredList = append(filteredList, item)
		}
	}
	downloadedList = filteredList

	for i, element := range downloadedList {
		status, err := c.GetModelEstimate(element.Path)
		if err != nil {
			return []LMSDownloadedListItem{}, fmt.Errorf("Failed to run estimation is listing: %w", err)
		}
		downloadedList[i].CanLoad = (status == StatusAvailable)
	}
	return downloadedList, nil
}

func (c *Client) GetLoadedModels() ([]LMSLoadedListItem, error) {
	if c.IsClosed() {
		return nil, fmt.Errorf("client is closed")
	}
	cmd := []string{"ps", "--json"}
	output, err := c.RunLMSCommand(cmd)

	if err != nil {
		return []LMSLoadedListItem{}, fmt.Errorf("Failed to run lms ls: %w", err)
	}

	var downloadedList []LMSLoadedListItem
	if err := json.Unmarshal([]byte(output), &downloadedList); err != nil {
		return []LMSLoadedListItem{}, fmt.Errorf("Failed to parse JSON output: %w", err)
	}

	return downloadedList, nil
}

func (c *Client) LoadModel(modelID string) error {
	if err := ValidateModelID(modelID); err != nil {
		return fmt.Errorf("invalid model ID: %w", err)
	}

	if c.IsClosed() {
		return fmt.Errorf("client is closed")
	}

	_, err := c.RunLMSCommand([]string{"load", modelID})
	if err != nil {
		c.logger.Error("Failed to load model %s: %v", modelID, err)
		return fmt.Errorf("failed to load model: %w", err)
	}
	return nil
}

func (c *Client) UnloadModel(modelKey string) error {
	if err := ValidateModelID(modelKey); err != nil {
		return fmt.Errorf("invalid model key: %w", err)
	}

	if c.IsClosed() {
		return fmt.Errorf("client is closed")
	}

	_, err := c.RunLMSCommand([]string{"unload", modelKey})
	if err != nil {
		c.logger.Error("Failed to unload model %s: %v", modelKey, err)
		return fmt.Errorf("failed to unload model: %w", err)
	}
	c.logger.Info("Successfully unloaded model: %s", modelKey)
	return nil
}

func (c *Client) UnloadAllModels() error {
	if c.IsClosed() {
		return fmt.Errorf("client is closed")
	}
	_, err := c.RunLMSCommand([]string{"unload", "--all"})
	if err == nil {
		c.logger.Info("Unloaded all models using --all flag")
		return nil
	}

	loadedModels, err := c.GetLoadedModels()
	if err != nil {
		c.logger.Error("Failed to get loaded models for unload all: %v", err)
		return err
	}

	if len(loadedModels) == 0 {
		c.logger.Info("No models to unload")
		return nil
	}

	for _, model := range loadedModels {
		if err := c.UnloadModel(model.Identifier); err != nil {
			c.logger.Error("Failed to unload model %s during unload all: %v", model.Identifier, err)
			return err
		}
	}

	c.logger.Info("Unloaded all models individually")
	return nil
}
