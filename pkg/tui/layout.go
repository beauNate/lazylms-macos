package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/layout"
	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// renderHelpView renders the help view
func (m Model) renderHelpView() string {
	helpContent := `
LazyLMS - Model Management TUI

NAVIGATION:
   Tab          Cycle through panels
   1-5          Jump to specific panel
   h, ?         Toggle this help screen

MODEL MANAGEMENT:
   ↑/↓, j/k     Navigate in lists
   Enter        Load model (from downloaded list)
   Enter        Select model for chat (from loaded list)
   u            Unload single model (from loaded list)
   U            Unload all models (from loaded list)

CHAT:
    4            Enter chat mode (input field becomes active)
    Tab          Exit chat mode and switch to other views
    Esc          Exit chat mode to status view
    Enter        Send message (when in chat mode)
    ↑/↓          Scroll chat history (when in chat mode)
    PgUp/PgDown  Page up/down in chat history
    Home/End     Go to top/bottom of chat history

SYSTEM PROMPT:
   Ctrl+S       Open system prompt popup
   Enter        Set system prompt (in popup)
   Esc          Close system prompt popup

GENERAL:
    Ctrl+C    Quit application
    Ctrl+L    Clear chat history

PANELS:
   1 Status     Server connection status
   2 Loaded     Currently loaded models
   3 Downloaded All available models
   4 Chat       Chat interface (press 4 to activate input)
   5 Logs       Server logs

Press 'h' or '?' to close this help.
`
	content := helpContent

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "◐ Help",
	}

	return layout.Borderize(content, true, m.width-4, m.height-4, embeddedText)
}

// renderSystemPopup renders the system prompt popup
func (m Model) renderSystemPopup() string {
	popupWidth := 60

	instructions := "Enter system prompt and press Enter to set, or Esc to cancel"
	instructionsStyle := lipgloss.NewStyle().
		Foreground(styles.ColorGray).
		Align(lipgloss.Center).
		Width(popupWidth - 8)

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		m.systemInput.View(),
		"",
		instructionsStyle.Render(instructions),
		"",
	)

	// content is already set above

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "⚙ Set System Prompt",
	}

	popup := layout.Borderize(content, true, popupWidth, 8, embeddedText)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, popup)
}

// renderFooterView renders the footer
func (m Model) renderFooterView() string {
	foregroundColor := styles.ColorGray
	backgroundColor := styles.ColorBlack

	// If footer is somehow the current view (though it's not in the current implementation)
	// we could highlight it, but for now keep it consistent
	if m.currentView == "footer" { // This won't happen with current key handling
		foregroundColor = styles.ColorWhite
		backgroundColor = styles.ColorPurple
	}

	style := lipgloss.NewStyle().
		Foreground(foregroundColor).
		Background(backgroundColor).
		Width(m.width).
		Height(1).
		Align(lipgloss.Center)

	var content string
	if m.currentView == "chat" {
		if m.streaming {
			content = "tab: panels | ctrl+x: cancel | ctrl+l: clear chat | ↑↓/pgup/home: nav | esc: exit"
		} else {
			content = "tab: panels | enter: send | ctrl+l: clear chat | ↑↓/pgup/home: nav | esc: exit"
		}
	} else {
		content = "1-5: panels | ctrl+s: system prompt | ctrl+l: clear chat | enter: select | h: help | ctrl+c: exit | LazyLMS BETA"
	}
	return style.Render(content)
}
