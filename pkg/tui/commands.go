package tui

import (
	"fmt"
	"time"

	"github.com/Rugz007/lazylms/pkg/client"
	tea "github.com/charmbracelet/bubbletea"
)

func tickCmd() tea.Cmd {
	return tea.Tick(client.DefaultTickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func estimateTickCmd() tea.Cmd {
	return tea.Tick(EstimateTickInterval, func(t time.Time) tea.Msg {
		return estimateTickMsg(t)
	})
}

func animationTickCmd() tea.Cmd {
	return tea.Tick(AnimationTickInterval, func(t time.Time) tea.Msg {
		return animationTickMsg(t)
	})
}

// nextViewCmd cycles to the next view
func (m Model) nextViewCmd() tea.Cmd {
	return func() tea.Msg {
		views := []string{"status", "loaded", "downloaded", "chat", "logs"}
		for i, view := range views {
			if m.currentView == view {
				nextView := views[(i+1)%len(views)]
				newModel := m
				newModel.currentView = nextView
				if nextView == "chat" {
					newModel.chatInput.Focus()
				} else {
					newModel.chatInput.Blur()
				}
				return nextViewMsg(nextView)
			}
		}
		return nil
	}
}

func (m Model) updateStatusCmd() tea.Cmd {
	return func() tea.Msg {
		status, err := m.client.GetStatus()
		if err != nil {
			return statusMsg(false)
		}
		return statusMsg(status == client.StatusOn)
	}
}

func (m Model) updateModelsCmd() tea.Cmd {
	return func() tea.Msg {
		downloaded, err := m.client.GetDownloadedModelsWithoutEstimates()
		if err != nil {
			// Log the error for debugging
			if m.client != nil && m.client.GetLogger() != nil {
				m.client.GetLogger().Error("Failed to get downloaded models: %v", err)
			}
			downloaded = []client.LMSDownloadedListItem{}
		}

		// Preserve existing CanLoad status from current models
		existingCanLoad := make(map[string]bool)
		for _, model := range m.downloadedModels {
			existingCanLoad[model.ModelKey] = model.CanLoad
		}

		// Apply existing CanLoad status to new list
		for i, model := range downloaded {
			if canLoad, exists := existingCanLoad[model.ModelKey]; exists {
				downloaded[i].CanLoad = canLoad
			}
		}

		loaded, err := m.client.GetLoadedModels()
		if err != nil {
			// Log the error for debugging
			if m.client != nil && m.client.GetLogger() != nil {
				m.client.GetLogger().Error("Failed to get loaded models: %v", err)
			}
			loaded = []client.LMSLoadedListItem{}
		}

		return modelsMsg{
			downloaded: downloaded,
			loaded:     loaded,
		}
	}
}

// shouldRunEstimates determines if estimates should be run based on loaded model changes
func (m Model) shouldRunEstimates(loadedModels []client.LMSLoadedListItem) bool {
	if m.firstLoad {
		return true
	}

	currentLoadedIDs := extractModelIDs(loadedModels)

	if len(currentLoadedIDs) != len(m.lastLoadedModelIDs) {
		return true
	}

	for i, id := range currentLoadedIDs {
		if i >= len(m.lastLoadedModelIDs) || m.lastLoadedModelIDs[i] != id {
			return true
		}
	}

	return false
}

func (m Model) updateEstimatesCmd() tea.Cmd {
	return func() tea.Msg {
		downloaded, err := m.client.GetDownloadedModels()
		if err != nil {
			if m.client != nil && m.client.GetLogger() != nil {
				m.client.GetLogger().Error("Failed to get downloaded models with estimates: %v", err)
			}
			downloaded = []client.LMSDownloadedListItem{}
		}

		return estimatesMsg(downloaded)
	}
}

func (m Model) startLogListening() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return logMsg("LazyLMS TUI started")
		},
		m.logListenerCmd(),
	)
}

// logListenerCmd creates a command that checks for log messages
func (m Model) logListenerCmd() tea.Cmd {
	return tea.Tick(LogCheckInterval, func(t time.Time) tea.Msg {
		return logListenerMsg{}
	})
}

func (m Model) handleDownloadedModelSelection() tea.Cmd {
	if m.downloadedList.Items() == nil || len(m.downloadedList.Items()) == 0 {
		return nil
	}

	selectedItem := m.downloadedList.SelectedItem()
	if selectedItem == nil {
		return nil
	}

	modelItem, ok := selectedItem.(downloadedModelItem)
	if !ok {
		if m.client != nil && m.client.GetLogger() != nil {
			m.client.GetLogger().Warn("Selected item is not a modelItem: %v", selectedItem)
		}
		return nil
	}

	if modelItem.model.ModelKey == "" {
		if m.client != nil && m.client.GetLogger() != nil {
			m.client.GetLogger().Warn("Selected model has empty ID")
		}
		return nil
	}
	m.client.GetLogger().Info("Loading model: %s", modelItem.model.ModelKey)
	return m.loadModelCmd(modelItem.model.ModelKey)
}

func (m Model) loadModelCmd(modelID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.LoadModel(modelID)
		if err != nil {
			return logMsg(fmt.Sprintf("Failed to load model %s: %v", modelID, err))
		}
		return logMsg(fmt.Sprintf("Loaded model: %s", modelID))
	}
}

func (m Model) unloadModelCmd(modelKey string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.UnloadModel(modelKey)
		if err != nil {
			return logMsg(fmt.Sprintf("Failed to unload model %s: %v", modelKey, err))
		}
		return nil
	}
}

func (m Model) unloadAllModelsCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.client.UnloadAllModels()
		if err != nil {
			return logMsg(fmt.Sprintf("Failed to unload all models: %v", err))
		}
		return nil
	}
}

// streamSubscription creates a subscription to read streaming chunks
func (m Model) streamSubscription() tea.Cmd {
	return func() tea.Msg {
		msg := <-m.streamChan
		// Handle both typed messages and legacy string messages
		switch v := msg.(type) {
		case streamChunkMsg:
			return v
		case string:
			// Legacy string messages (error and completion)
			return streamChunkMsg{content: v, contentType: "output"}
		default:
			return streamChunkMsg{content: "", contentType: "output"}
		}
	}
}
