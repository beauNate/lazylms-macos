package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/layout"
	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// renderStatusView renders the status view
func (m Model) renderStatusView(_, leftColumnWidth int) string {
	active := m.currentView == "status"

	// Convert boolean status to display text
	var content string
	if m.status {
		statusText := "ON"
		content = lipgloss.NewStyle().Background(styles.ColorGreen).Padding(0, 1).Foreground(styles.ColorForegroundInverted).Render(statusText)
	} else {
		statusText := "OFF"
		content = lipgloss.NewStyle().Background(styles.ColorYellow).Padding(0, 1).Foreground(styles.ColorForegroundInverted).Render(statusText)
	}
	title := lipgloss.NewStyle().Italic(true).Render("👾  lazylms")
	content = lipgloss.NewStyle().Padding(1).Render(title + "\nAPI Server: " + content)

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "[1] ◉ Status",
	}

	return layout.Borderize(content, active, leftColumnWidth-2, 6, embeddedText)
}
