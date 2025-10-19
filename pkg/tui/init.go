package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		estimateTickCmd(),
		animationTickCmd(),
		m.startLogListening(),
		m.updateStatusCmd(),
		m.updateModelsCmd(),
		m.updateEstimatesCmd(),
	)
}
