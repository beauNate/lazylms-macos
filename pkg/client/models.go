package client

import (
	"fmt"
)

func (c *Client) GetActiveModel() (string, error) {
	models, err := c.GetLoadedModels()
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no models are currently loaded")
	}

	return models[0].Identifier, nil
}

func (c *Client) GetActiveModelByID(modelID string) (string, error) {
	models, err := c.GetLoadedModels()
	if err != nil {
		return "", err
	}

	for _, model := range models {
		if model.Identifier == modelID {

			return model.Identifier, nil
		}
	}

	if len(models) > 0 {
		return models[0].Identifier, nil
	}

	return "", fmt.Errorf("no models are currently loaded")
}
