package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/layout"
	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// renderChatView renders the chat view
func (m Model) renderChatView(mainHeight, rightColumnWidth int) string {
	active := m.currentView == "chat"

	var content string
	if !m.status {
		content = "Server is OFF"
	} else {
		content = m.chatViewport.View()
	}

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "[4] ◆ Chat",
	}

	content = lipgloss.NewStyle().Padding(1).Render(content)

	return layout.Borderize(content, active, rightColumnWidth-2, ((mainHeight-6)*3/4)-1, embeddedText)
}

// renderChatInputView renders the chat input field
func (m Model) renderChatInputView(rightColumnWidth int) string {
	active := m.currentView == "chat"

	var content string
	if !m.status {
		content = "Server is OFF"
	} else if m.currentView == "chat" {
		content = m.chatInput.View()
	} else {
		// Show placeholder text to indicate input area
		placeholder := lipgloss.NewStyle().
			Foreground(styles.ColorGray).
			Italic(true).
			Render("Chat input (press 4 to activate)")
		content = placeholder
	}

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "Input",
	}

	return layout.Borderize(content, active, rightColumnWidth-2, 3, embeddedText)
}
