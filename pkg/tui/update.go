package tui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/client"
	"github.com/Rugz007/lazylms/pkg/tui/keybindings"
	"github.com/Rugz007/lazylms/pkg/tui/rendering"
	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// extractModelIDs extracts identifiers from loaded models
func extractModelIDs(models []client.LMSLoadedListItem) []string {
	ids := make([]string, len(models))
	for i, model := range models {
		ids[i] = model.Identifier
	}
	return ids
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	globalKeyMap := keybindings.DefaultGlobalKeyMap()
	viewKeyMap := keybindings.DefaultViewKeyMap()
	chatKeyMap := keybindings.DefaultChatKeyMap()
	listKeyMap := keybindings.DefaultListKeyMap()

	switch msg := msg.(type) {
	case emptyMsg:
	// Do nothing
	case tea.KeyMsg:
		return m.handleKeyMsg(msg, globalKeyMap, viewKeyMap, chatKeyMap, listKeyMap)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listWidth := m.width/4 - 4
		listHeight := (m.height - 12) / 2
		m.loadedList.SetSize(listWidth-2, listHeight-2)
		m.downloadedList.SetSize(listWidth-2, listHeight-2)

		rightColumnWidth := 3 * m.width / 4
		footerHeight := 1
		availableHeight := m.height - footerHeight

		// Always reserve space for input: chat (60%), input (10%), logs (30%)
		chatHeight := 6 * availableHeight / 10
		inputHeight := availableHeight / 10
		logsHeight := 3 * availableHeight / 10

		m.chatViewport.Width = rightColumnWidth - 4
		m.chatViewport.Height = chatHeight - 2
		m.logsViewport.Width = rightColumnWidth - 4
		m.logsViewport.Height = logsHeight - 3

		m.chatInput.Width = rightColumnWidth - 4
		_ = inputHeight

	case tickMsg:
		return m, tea.Batch(
			tickCmd(),
			m.updateStatusCmd(),
			m.updateModelsCmd(),
		)

	case estimateTickMsg:
		// Check if loaded models have changed since last estimate
		currentLoadedIDs := extractModelIDs(m.loadedModels)

		// Compare with last known loaded model IDs
		loadedModelsChanged := len(currentLoadedIDs) != len(m.lastLoadedModelIDs)
		if !loadedModelsChanged {
			for i, id := range currentLoadedIDs {
				if i >= len(m.lastLoadedModelIDs) || m.lastLoadedModelIDs[i] != id {
					loadedModelsChanged = true
					break
				}
			}
		}

		if m.firstLoad || loadedModelsChanged {
			m.firstLoad = false
			m.lastEstimateTime = time.Now()
			m.lastLoadedModelIDs = currentLoadedIDs
			return m, tea.Batch(
				estimateTickCmd(),
				m.updateEstimatesCmd(),
			)
		}

		return m, estimateTickCmd()

	case animationTickMsg:
		m.animationTime = time.Time(msg)
		return m, animationTickCmd()

	case statusMsg:
		m.status = bool(msg)

	case modelsMsg:
		previousLoadedModels := m.loadedModels
		m.downloadedModels = msg.downloaded
		m.loadedModels = msg.loaded

		loadedItems := make([]list.Item, len(msg.loaded))
		for i, model := range msg.loaded {
			loadedItems[i] = loadedModelItem{model: model}
		}
		m.loadedList.SetItems(loadedItems)

		downloadedItems := make([]list.Item, len(msg.downloaded))
		for i, model := range msg.downloaded {
			downloadedItems[i] = downloadedModelItem{model: model}
		}
		m.downloadedList.SetItems(downloadedItems)

		// Auto-select most recently loaded model
		if len(msg.loaded) > 0 && len(previousLoadedModels) < len(msg.loaded) {
			// Find newly loaded model(s)
			previousIdentifiers := make(map[string]bool)
			for _, model := range previousLoadedModels {
				previousIdentifiers[model.Identifier] = true
			}

			// Find the first new model and auto-select it
			for _, model := range msg.loaded {
				if !previousIdentifiers[model.Identifier] {
					m.explicitlySelectedModel = model.Identifier
					m.selectedModel = model.Identifier
					break
				}
			}
		}

		// Check if explicitly selected model is still available
		if m.explicitlySelectedModel != "" {
			found := false
			for _, model := range msg.loaded {
				if model.Identifier == m.explicitlySelectedModel {
					found = true
					break
				}
			}
			if !found {
				// Selected model was unloaded, clear selection
				m.explicitlySelectedModel = ""
				m.selectedModel = ""
			}
		}

		// Check if we should run estimates due to loaded model changes
		shouldRunEstimates := m.shouldRunEstimates(msg.loaded)
		if shouldRunEstimates {
			m.firstLoad = false
			m.lastEstimateTime = time.Now()
			currentLoadedIDs := make([]string, len(msg.loaded))
			for i, model := range msg.loaded {
				currentLoadedIDs[i] = model.Identifier
			}
			m.lastLoadedModelIDs = currentLoadedIDs

			return m, m.updateEstimatesCmd()
		}

		if len(msg.loaded) == 0 {
			m.chatInput.Placeholder = NoModelsPlaceholder
		} else if m.selectedModel != "" {
			m.chatInput.Placeholder = fmt.Sprintf(ChatWithModelPlaceholder, m.selectedModel)
		} else {
			m.chatInput.Placeholder = SelectModelPlaceholder
		}

	case estimatesMsg:
		estimateMap := make(map[string]bool)
		for _, model := range msg {
			estimateMap[model.ModelKey] = model.CanLoad
		}

		for i, model := range m.downloadedModels {
			if canLoad, exists := estimateMap[model.ModelKey]; exists {
				m.downloadedModels[i].CanLoad = canLoad
			}
		}

		// Update list items with new CanLoad status
		downloadedItems := make([]list.Item, len(m.downloadedModels))
		for i, model := range m.downloadedModels {
			downloadedItems[i] = downloadedModelItem{model: model}
		}
		m.downloadedList.SetItems(downloadedItems)

	case logMsg:
		// Add log message and keep only last MaxLogLines messages
		logLines := strings.Split(m.logsViewport.View(), "\n")
		if len(logLines) == 1 && logLines[0] == "" {
			logLines = []string{}
		}
		// Wrap the new log message to viewport width
		wrappedMsg := rendering.WrapText(string(msg), m.logsViewport.Width-2)
		logLines = append(logLines, wrappedMsg)
		if len(logLines) > MaxLogLines {
			logLines = logLines[len(logLines)-MaxLogLines:]
		}
		m.logsViewport.SetContent(strings.Join(logLines, "\n"))
		m.logsViewport.GotoBottom()

	case streamChunkMsg:
		chunk := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(msg.content, "")

		if chunk == "STREAM_COMPLETE" {
			m.streaming = false

			// Restore the original status of the selected model
			if m.originalStreamingStatus != "" {
				for i, model := range m.loadedModels {
					if model.Identifier == m.selectedModel {
						m.loadedModels[i].Status = m.originalStreamingStatus
						m.originalStreamingStatus = ""
						break
					}
				}

				// Update the loaded list items to reflect the change
				loadedItems := make([]list.Item, len(m.loadedModels))
				for i, model := range m.loadedModels {
					loadedItems[i] = loadedModelItem{model: model}
				}
				m.loadedList.SetItems(loadedItems)
			}

			aiMsg := rendering.ChatMessage{
				Type:     rendering.MessageTypeAI,
				Author:   m.selectedModel,
				Segments: m.currentResponse.Segments,
			}
			m.chatMessages = append(m.chatMessages, aiMsg)

			// Add response to client conversation for context
			// Build content for assistant message from segments
			var responseContent strings.Builder
			for _, seg := range m.currentResponse.Segments {
				responseContent.WriteString(seg.Text)
			}
			m.client.AddAssistantMessage(responseContent.String())

			m.currentResponse.Reset()
			m.chatInput.Placeholder = fmt.Sprintf(ChatWithModelPlaceholder, m.selectedModel)
			// Update viewport - render all messages with markdown
			var renderedMessages []string
			for _, msg := range m.chatMessages {
				renderedMessages = append(renderedMessages, rendering.RenderChatMessage(msg, m.chatViewport.Width, m.logChan))
			}
			content := strings.Join(renderedMessages, "\n")
			m.chatViewport.SetContent(content)
			m.chatViewport.GotoBottom()
			return m, nil
		} else if strings.HasPrefix(chunk, "ERROR:") {
			m.streaming = false

			// Restore the original status of the selected model
			if m.originalStreamingStatus != "" {
				for i, model := range m.loadedModels {
					if model.Identifier == m.selectedModel {
						m.loadedModels[i].Status = m.originalStreamingStatus
						m.originalStreamingStatus = ""
						break
					}
				}

				loadedItems := make([]list.Item, len(m.loadedModels))
				for i, model := range m.loadedModels {
					loadedItems[i] = loadedModelItem{model: model}
				}
				m.loadedList.SetItems(loadedItems)
			}

			errorMsg := strings.TrimPrefix(chunk, "ERROR:")
			m.chatInput.Placeholder = fmt.Sprintf(ChatWithModelPlaceholder, m.selectedModel)
			return m, tea.Cmd(func() tea.Msg { return logMsg(fmt.Sprintf("Streaming error: %s", errorMsg)) })
		} else {
			// Regular chunk - add segment
			m.currentResponse.AddSegment(chunk, msg.contentType)
			// Update chat viewport with partial response - render all messages with markdown
			var renderedMessages []string
			for _, chatMsg := range m.chatMessages {
				renderedMessages = append(renderedMessages, rendering.RenderChatMessage(chatMsg, m.chatViewport.Width, m.logChan))
			}
			content := strings.Join(renderedMessages, "\n")
			if m.streaming {
				// Render streaming message with segments
				styledPrefix := lipgloss.NewStyle().
					Foreground(styles.ColorPurple).
					Bold(true).
					Render(m.selectedModel + ":")
				rendered := rendering.RenderMixedContent(m.currentResponse.Segments, m.chatViewport.Width, m.logChan)
				content += "\n" + styledPrefix + rendered
			}
			m.chatViewport.SetContent(content)
			m.chatViewport.GotoBottom()
			return m, m.streamSubscription() // Continue listening for more chunks
		}

	case logListenerMsg:
		// Check for new log messages from client.Logger
		select {
		case logLine := <-m.logChan:
			return m, tea.Batch(m.logListenerCmd(), func() tea.Msg { return logMsg(logLine) })
		default:
			return m, m.logListenerCmd()
		}
	case nextViewMsg:
		m.currentView = string(msg)
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg, globalKeyMap keybindings.GlobalKeyMap, viewKeyMap keybindings.ViewKeyMap, chatKeyMap keybindings.ChatKeyMap, listKeyMap keybindings.ListKeyMap) (tea.Model, tea.Cmd) {
	if m.currentView == "chat" && m.chatInput.Focused() {
		return m.handleChatInputKeys(msg, globalKeyMap, chatKeyMap)
	}

	// If system popup is open and input is focused, let the input handle all keys first
	if m.showSystemPopup && m.systemInput.Focused() {
		return m.handleSystemInputKeys(msg, globalKeyMap)
	}

	switch msg.String() {
	case globalKeyMap.Quit.Keys()[0]:
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case globalKeyMap.NextView.Keys()[0]:
		return m, m.nextViewCmd()
	case viewKeyMap.Status.Keys()[0]:
		m.currentView = "status"
		m.chatInput.Blur()
		return m, nil
	case viewKeyMap.Loaded.Keys()[0]:
		m.currentView = "loaded"
		m.chatInput.Blur()
		return m, nil
	case viewKeyMap.Downloaded.Keys()[0]:
		m.currentView = "downloaded"
		m.chatInput.Blur()
		return m, nil
	case viewKeyMap.Chat.Keys()[0]:
		m.currentView = "chat"
		m.chatInput.Focus()
		return m, nil
	case viewKeyMap.Logs.Keys()[0]:
		m.currentView = "logs"
		m.chatInput.Blur()
		return m, nil
	case globalKeyMap.SystemPrompt.Keys()[0]:
		// Toggle system prompt popup
		m.showSystemPopup = !m.showSystemPopup
		if m.showSystemPopup {
			m.systemInput.SetValue(m.systemPrompt)
			m.systemInput.Focus()
		} else {
			m.systemInput.Blur()
		}
		return m, nil
	case globalKeyMap.Help.Keys()[0], globalKeyMap.Help.Keys()[1]:
		m.showHelp = !m.showHelp
		return m, nil
	case globalKeyMap.ClearChat.Keys()[0]:
		m.chatMessages = []rendering.ChatMessage{}
		m.client.ClearConversation()
		m.chatViewport.SetContent("")
		m.chatViewport.GotoTop()
		return m, tea.Cmd(func() tea.Msg { return logMsg("Chat cleared") })
	case "esc":
		if m.showHelp {
			m.showHelp = false
		} else if m.showSystemPopup {
			m.showSystemPopup = false
			m.systemInput.Blur()
		}
		return m, nil

	case listKeyMap.Unload.Keys()[0]:
		if m.currentView == "loaded" && len(m.loadedModels) > 0 {
			selectedIndex := m.loadedList.Index()
			if selectedIndex < len(m.loadedModels) {
				modelID := m.loadedModels[selectedIndex].Identifier
				if modelID == m.explicitlySelectedModel {
					m.explicitlySelectedModel = ""
					m.selectedModel = ""
				}
				return m, m.unloadModelCmd(modelID)
			}
		}
	case listKeyMap.UnloadAll.Keys()[0]:
		if m.currentView == "loaded" {
			m.explicitlySelectedModel = ""
			m.selectedModel = ""
			return m, m.unloadAllModelsCmd()
		}
	case "enter":
		if m.currentView == "downloaded" {
			return m, m.handleDownloadedModelSelection()
		} else if m.currentView == "loaded" && len(m.loadedModels) > 0 {
			selectedIndex := m.loadedList.Index()
			if selectedIndex < len(m.loadedModels) {
				m.explicitlySelectedModel = m.loadedModels[selectedIndex].Identifier
				m.selectedModel = m.explicitlySelectedModel
				return m, tea.Cmd(func() tea.Msg { return logMsg(fmt.Sprintf("Selected model: %s", m.selectedModel)) })
			}
		} else if m.currentView == "system" {
			// Set system prompt
			systemPrompt := m.systemInput.Value()
			if systemPrompt != "" {
				m.systemPrompt = systemPrompt
				m.client.SetSystemMessage(systemPrompt)
				m.systemInput.SetValue("")
				return m, nil
			}
		}
	default:
		switch m.currentView {
		case "loaded":
			var cmd tea.Cmd
			m.loadedList, cmd = m.loadedList.Update(msg)
			return m, cmd
		case "downloaded":
			var cmd tea.Cmd
			m.downloadedList, cmd = m.downloadedList.Update(msg)
			return m, cmd
		case "logs":
			return m.handleLogsKeys(msg, chatKeyMap)
		case "chat":
			return m.handleChatViewportKeys(msg, chatKeyMap)
		}
	}

	return m, nil
}

func (m Model) handleChatInputKeys(msg tea.KeyMsg, globalKeyMap keybindings.GlobalKeyMap, chatKeyMap keybindings.ChatKeyMap) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case chatKeyMap.ScrollUp.Keys()[0]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.ScrollDown.Keys()[0]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.PageUp.Keys()[0], chatKeyMap.PageUp.Keys()[1]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.PageDown.Keys()[0], chatKeyMap.PageDown.Keys()[1]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.Home.Keys()[0]:
		m.chatViewport.GotoTop()
		return m, nil
	case chatKeyMap.End.Keys()[0]:
		m.chatViewport.GotoBottom()
		return m, nil
	case globalKeyMap.NextView.Keys()[0]:
		return m, m.nextViewCmd()
	case globalKeyMap.ClearChat.Keys()[0]:
		m.chatMessages = []rendering.ChatMessage{}
		m.client.ClearConversation()
		m.chatViewport.SetContent("")
		m.chatViewport.GotoTop()
		return m, tea.Cmd(func() tea.Msg { return logMsg("Chat cleared") })
	case "ctrl+c":
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case chatKeyMap.SendMessage.Keys()[0]:
		// Don't allow sending messages while streaming
		if m.streaming {
			return m, nil
		}
		message := m.chatInput.Value()
		if message != "" {
			if m.selectedModel == "" {
				return m, tea.Cmd(func() tea.Msg {
					return logMsg("No model selected. Please select a model first using 's' key in loaded models view.")
				})
			}

			m.chatInput.SetValue("")
			// Clear welcome message if present
			if m.hasWelcomeMessage {
				m.hasWelcomeMessage = false
			}
			userMsg := rendering.ChatMessage{
				Type:    rendering.MessageTypeUser,
				Author:  "You",
				Content: message,
			}
			m.chatMessages = append(m.chatMessages, userMsg)
			// Update viewport with new message - render all messages with markdown
			var renderedMessages []string
			for _, msg := range m.chatMessages {
				renderedMessages = append(renderedMessages, rendering.RenderChatMessage(msg, m.chatViewport.Width, m.logChan))
			}
			content := strings.Join(renderedMessages, "\n")
			m.chatViewport.SetContent(content)
			m.chatViewport.GotoBottom()
			// Start streaming response
			m.streaming = true
			m.chatInput.Placeholder = GeneratingPlaceholder
			m.currentResponse.Reset()

			// Update the selected model's status to "generating"
			for i, model := range m.loadedModels {
				if model.Identifier == m.selectedModel {
					m.originalStreamingStatus = model.Status
					m.loadedModels[i].Status = "generating"
					break
				}
			}

			// Update the loaded list items to reflect the change
			loadedItems := make([]list.Item, len(m.loadedModels))
			for i, model := range m.loadedModels {
				loadedItems[i] = loadedModelItem{model: model}
			}
			m.loadedList.SetItems(loadedItems)
			// Send the message using streaming
			go func() {
				err := m.client.SendMessageStreamWithModel(m.ctx, message, m.selectedModel, func(chunk string, contentType string) {
					select {
					case m.streamChan <- streamChunkMsg{content: chunk, contentType: contentType}:
					default:
					}
				})
				// Only send error/completion if not cancelled
				if !m.client.IsCancelled() {
					if err != nil {
						select {
						case m.streamChan <- fmt.Sprintf("ERROR: %v", err):
						default:
						}
					}
					// Send completion signal
					select {
					case m.streamChan <- "STREAM_COMPLETE":
					default:
					}
				}
			}()
			return m, m.streamSubscription()
		}
		return m, nil
	case chatKeyMap.Cancel.Keys()[0]:
		// Cancel ongoing request
		if m.streaming {
			m.client.CancelRequest()
			m.streaming = false

			// Drain streamChan to prevent stale messages
			go func() {
				for {
					select {
					case <-m.streamChan:
						// Consume and discard
					default:
						return
					}
				}
			}()

			// Save partial response to chat history before cancelling
			if len(m.currentResponse.Segments) > 0 {
				// Build content from segments
				var responseContent strings.Builder
				for _, seg := range m.currentResponse.Segments {
					responseContent.WriteString(seg.Text)
				}
				// Append cancelled marker as output segment
				cancelledSegments := append([]rendering.ContentSegment{}, m.currentResponse.Segments...)
				cancelledSegments = append(cancelledSegments, rendering.ContentSegment{
					Text: " [cancelled]",
					Type: rendering.ContentTypeOutput,
				})
				aiMsg := rendering.ChatMessage{
					Type:     rendering.MessageTypeAI,
					Author:   m.selectedModel,
					Segments: cancelledSegments,
				}
				m.chatMessages = append(m.chatMessages, aiMsg)

				// Add partial response to client conversation for context
				m.client.AddAssistantMessage(responseContent.String())

				m.currentResponse.Reset()

				// Update viewport with cancelled message
				var renderedMessages []string
				for _, msg := range m.chatMessages {
					renderedMessages = append(renderedMessages, rendering.RenderChatMessage(msg, m.chatViewport.Width, m.logChan))
				}
				content := strings.Join(renderedMessages, "\n")
				m.chatViewport.SetContent(content)
				m.chatViewport.GotoBottom()
			}

			// Restore the original status of the selected model
			if m.originalStreamingStatus != "" {
				for i, model := range m.loadedModels {
					if model.Identifier == m.selectedModel {
						m.loadedModels[i].Status = m.originalStreamingStatus
						m.originalStreamingStatus = ""
						break
					}
				}

				// Update the loaded list items to reflect the change
				loadedItems := make([]list.Item, len(m.loadedModels))
				for i, model := range m.loadedModels {
					loadedItems[i] = loadedModelItem{model: model}
				}
				m.loadedList.SetItems(loadedItems)
			}

			m.chatInput.Placeholder = fmt.Sprintf(ChatWithModelPlaceholder, m.selectedModel)
			return m, tea.Cmd(func() tea.Msg { return logMsg("Request cancelled by user") })
		}
		return m, nil
	case chatKeyMap.ExitChat.Keys()[0]:
		m.currentView = "status"
		m.chatInput.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.chatInput, cmd = m.chatInput.Update(msg)
		return m, cmd
	}
}

// handleSystemInputKeys handles keys when system input is focused
func (m Model) handleSystemInputKeys(msg tea.KeyMsg, globalKeyMap keybindings.GlobalKeyMap) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case globalKeyMap.Quit.Keys()[0]:
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case globalKeyMap.ClearChat.Keys()[0]:
		// Clear chat messages
		m.chatMessages = []rendering.ChatMessage{}
		m.chatViewport.SetContent("")
		m.chatViewport.GotoTop()
		return m, tea.Cmd(func() tea.Msg { return logMsg("Chat cleared") })
	case "enter":
		// Set system prompt
		systemPrompt := m.systemInput.Value()
		if systemPrompt != "" {
			m.systemPrompt = systemPrompt
			m.client.SetSystemMessage(systemPrompt)
			m.systemInput.SetValue("")
			m.showSystemPopup = false
			m.systemInput.Blur()
			return m, nil
		}
		return m, nil
	case "esc":
		m.showSystemPopup = false
		m.systemInput.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.systemInput, cmd = m.systemInput.Update(msg)
		return m, cmd
	}
}

// handleLogsKeys handles keys for logs view
func (m Model) handleLogsKeys(msg tea.KeyMsg, chatKeyMap keybindings.ChatKeyMap) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case chatKeyMap.ScrollUp.Keys()[0], chatKeyMap.ScrollUp.Keys()[1]:
		var cmd tea.Cmd
		m.logsViewport, cmd = m.logsViewport.Update(msg)
		return m, cmd
	case chatKeyMap.ScrollDown.Keys()[0], chatKeyMap.ScrollDown.Keys()[1]:
		var cmd tea.Cmd
		m.logsViewport, cmd = m.logsViewport.Update(msg)
		return m, cmd
	case "pgup":
		var cmd tea.Cmd
		m.logsViewport, cmd = m.logsViewport.Update(msg)
		return m, cmd
	case "pgdown":
		var cmd tea.Cmd
		m.logsViewport, cmd = m.logsViewport.Update(msg)
		return m, cmd
	case "home":
		m.logsViewport.GotoTop()
		return m, nil
	case "end":
		m.logsViewport.GotoBottom()
		return m, nil
	default:
		return m, nil
	}
}

// handleChatViewportKeys handles keys for chat viewport
func (m Model) handleChatViewportKeys(msg tea.KeyMsg, chatKeyMap keybindings.ChatKeyMap) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case chatKeyMap.ScrollUp.Keys()[0]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.ScrollDown.Keys()[0]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.PageUp.Keys()[0], chatKeyMap.PageUp.Keys()[1]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.PageDown.Keys()[0], chatKeyMap.PageDown.Keys()[1]:
		var cmd tea.Cmd
		m.chatViewport, cmd = m.chatViewport.Update(msg)
		return m, cmd
	case chatKeyMap.Home.Keys()[0]:
		m.chatViewport.GotoTop()
		return m, nil
	case chatKeyMap.End.Keys()[0]:
		m.chatViewport.GotoBottom()
		return m, nil
	case "enter":
		// This case is now handled by the main input handling above
		return m, nil
	default:
		// Let the text input handle the key
		var cmd tea.Cmd
		m.chatInput, cmd = m.chatInput.Update(msg)
		return m, cmd
	}
}
