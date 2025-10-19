package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// wrapText wraps text to the given width using utils.WrapText

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Show help view if requested
	if m.showHelp {
		return m.renderHelpView()
	}

	// Show system prompt popup if requested
	if m.showSystemPopup {
		return m.renderSystemPopup()
	}

	// Calculate dimensions
	footerHeight := 1
	mainHeight := m.height - footerHeight
	leftColumnWidth := m.width / 4      // 25% for left column
	rightColumnWidth := 3 * m.width / 4 // 75% for right column

	statusView := m.renderStatusView(mainHeight, leftColumnWidth)
	loadedView := m.renderLoadedModelsView(mainHeight, leftColumnWidth)
	downloadedView := m.renderDownloadedModelsView(mainHeight, leftColumnWidth)
	chatView := m.renderChatView(mainHeight, rightColumnWidth)
	logsView := m.renderLogsView(mainHeight, rightColumnWidth)
	footerView := m.renderFooterView()

	// Left column: status (6 lines) + loaded models + downloaded models
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		statusView,
		loadedView,
		downloadedView,
	)

	// Right column: chat + input + logs
	inputView := m.renderChatInputView(rightColumnWidth)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left,
		chatView,
		inputView,
		logsView,
	)

	// Main layout: left column + right column
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(leftColumnWidth).Render(leftColumn),
		lipgloss.NewStyle().Width(rightColumnWidth).Render(rightColumn),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		mainLayout,
		footerView,
	)
}
