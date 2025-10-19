package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/layout"
)

// renderLogsView renders the logs view
func (m Model) renderLogsView(mainHeight, rightColumnWidth int) string {
	active := m.currentView == "logs"

	var content string
	if !m.status {
		content = "Server is OFF"
	} else {
		content = m.logsViewport.View()
	}

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "[5] ✎ Logs",
	}

	content = lipgloss.NewStyle().Padding(0, 1).Render(content)

	return layout.Borderize(content, active, rightColumnWidth-2, (mainHeight)/4, embeddedText)
}
