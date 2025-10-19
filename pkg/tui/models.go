package tui

import (
	"io"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/client"
	"github.com/Rugz007/lazylms/pkg/tui/animation"
	"github.com/Rugz007/lazylms/pkg/tui/layout"
	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// styleLoadedModelItem styles a loaded model item based on its state
func styleLoadedModelItem(title, desc string, highlighted, explicitlySelected bool) string {
	if highlighted {
		if explicitlySelected {
			// Highlighted and selected - green border with bold text
			styledTitle := lipgloss.NewStyle().Foreground(styles.ColorWhite).Bold(true).Render(title)
			styledDesc := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(desc)
			return lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(styles.ColorOrange).
				PaddingLeft(1).
				Bold(true).
				Render(styledTitle + "\n" + styledDesc)
		} else {
			// Just highlighted - gray border
			styledTitle := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(title)
			styledDesc := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(desc)
			return lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(styles.ColorGray).
				PaddingLeft(1).
				Render(styledTitle + "\n" + styledDesc)
		}
	} else {
		if explicitlySelected {
			// Selected but not highlighted - green left border with bold text
			styledTitle := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(title)
			styledDesc := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(desc)
			return lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(styles.ColorOrange).
				PaddingLeft(1).
				Render(styledTitle + "\n" + styledDesc)
		} else {
			// Not selected or highlighted - gray text
			styledTitle := lipgloss.NewStyle().Foreground(styles.ColorGray).Render(title)
			styledDesc := lipgloss.NewStyle().Foreground(styles.ColorGray).Render(desc)
			return lipgloss.NewStyle().PaddingLeft(2).Render(styledTitle + "\n" + styledDesc)
		}
	}
}

type loadedModelItem struct {
	model client.LMSLoadedListItem
}

type downloadedModelItem struct {
	model client.LMSDownloadedListItem
}

func (d downloadedModelItem) Title() string {
	return d.model.DisplayName
}

func (d downloadedModelItem) Description() string {
	format := d.model.Format
	if format == "safetensors" {
		format = "mlx"
	}
	return format + " - " + d.model.Quantization.Name
}

func (d downloadedModelItem) FilterValue() string { return d.model.DisplayName }

func (m loadedModelItem) Title() string {
	return m.model.Identifier
}
func (m loadedModelItem) Description() string {
	return m.model.Status
}
func (m loadedModelItem) FilterValue() string { return m.model.Identifier }

// Custom delegate for downloaded models
type customDownloadedDelegate struct {
	height  int
	spacing int
}

func (d customDownloadedDelegate) Height() int                             { return d.height }
func (d customDownloadedDelegate) Spacing() int                            { return d.spacing }
func (d customDownloadedDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d *customDownloadedDelegate) SetHeight(h int)                        { d.height = h }
func (d *customDownloadedDelegate) SetSpacing(s int)                       { d.spacing = s }

// Custom delegate for loaded models with animation support
type customLoadedDelegate struct {
	height                         int
	spacing                        int
	generatingModelIdentifierArray []string
	animationTime                  time.Time
	selectedModel                  string
}

func (d customLoadedDelegate) Height() int                             { return d.height }
func (d customLoadedDelegate) Spacing() int                            { return d.spacing }
func (d customLoadedDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d *customLoadedDelegate) SetHeight(h int)                        { d.height = h }
func (d *customLoadedDelegate) SetSpacing(s int)                       { d.spacing = s }
func (d *customLoadedDelegate) SetGeneratingArray(g []string)          { d.generatingModelIdentifierArray = g }
func (d *customLoadedDelegate) SetAnimationTime(t time.Time)           { d.animationTime = t }
func (d *customLoadedDelegate) SetSelectedModel(s string)              { d.selectedModel = s }

func (d customLoadedDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if i, ok := item.(loadedModelItem); ok {
		highlighted := index == m.Index()
		explicitlySelected := i.model.Identifier == d.selectedModel

		title := i.model.Identifier

		// Use animated status if generating
		var desc string
		isGenerating := slices.Contains(d.generatingModelIdentifierArray, i.model.Identifier)
		if isGenerating && highlighted {
			desc = animation.GeneratingStatus(d.animationTime)
		} else {
			desc = i.model.Status
		}

		styledContent := styleLoadedModelItem(title, desc, highlighted, explicitlySelected)
		io.WriteString(w, styledContent)
	}
}

func (d customDownloadedDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if i, ok := item.(downloadedModelItem); ok {
		selected := index == m.Index()
		var title string
		if !i.model.CanLoad {
			title = lipgloss.NewStyle().Foreground(styles.ColorYellow).Render("⚠") + " " + i.model.DisplayName
		} else {
			title = i.model.DisplayName
		}

		desc := i.Description()

		if selected {
			styledTitle := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(title)
			styledDesc := lipgloss.NewStyle().Foreground(styles.ColorWhite).Render(desc)
			styledContent := lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(styles.ColorGray).
				PaddingLeft(1).
				Render(styledTitle + "\n" + styledDesc)
			io.WriteString(w, styledContent)
		} else {
			if !i.model.CanLoad {
				title = lipgloss.NewStyle().Foreground(styles.ColorYellow).Render("⚠") + " " + lipgloss.NewStyle().Foreground(styles.ColorGray).Render(i.model.DisplayName)
			} else {
				title = lipgloss.NewStyle().Foreground(styles.ColorGray).Render(title)
			}
			desc = lipgloss.NewStyle().Foreground(styles.ColorGray).Render(desc)
			styledContent := lipgloss.NewStyle().PaddingLeft(2).Render(title + "\n" + desc)
			io.WriteString(w, styledContent)
		}
	}
}

// renderLoadedModelsView renders the loaded models view
func (m Model) renderLoadedModelsView(mainHeight, leftColumnWidth int) string {
	active := m.currentView == "loaded"

	var content string
	if !m.status {
		content = lipgloss.NewStyle().Padding(1).Render("Server is OFF")
	} else {
		m.loadedList.SetShowStatusBar(false)
		m.loadedList.SetFilteringEnabled(false)
		m.loadedList.SetShowTitle(false)
		m.loadedList.AdditionalShortHelpKeys = func() []key.Binding {
			selectBinding := key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			)
			unloadBinding := key.NewBinding(
				key.WithKeys("u"),
				key.WithHelp("u", "unload"),
			)
			unloadAllBinding := key.NewBinding(
				key.WithKeys("U"),
				key.WithHelp("U", "unload all"),
			)
			var keyBindings []key.Binding
			keyBindings = append(keyBindings, selectBinding, unloadBinding, unloadAllBinding)
			return keyBindings
		}

		// Create custom delegate with animation support
		loadedDelegate := &customLoadedDelegate{}
		loadedDelegate.SetHeight(2)
		loadedDelegate.SetSpacing(1)

		// Get all the models which are generating\
		var generatingModelIdentifiers []string
		for _, element := range m.loadedModels {
			if element.Status == "generating" {
				generatingModelIdentifiers = append(generatingModelIdentifiers, element.Identifier)
			}
		}

		loadedDelegate.SetGeneratingArray(generatingModelIdentifiers)
		loadedDelegate.SetAnimationTime(m.animationTime)
		loadedDelegate.SetSelectedModel(m.explicitlySelectedModel)
		m.loadedList.SetDelegate(loadedDelegate)

		content = lipgloss.NewStyle().Padding(1).Render(m.loadedList.View())
	}

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "[2] ▶ Loaded Models",
	}

	return layout.Borderize(content, active, leftColumnWidth-2, (mainHeight-6)/2, embeddedText)
}

// renderDownloadedModelsView renders the downloaded models view
func (m Model) renderDownloadedModelsView(mainHeight, leftColumnWidth int) string {
	active := m.currentView == "downloaded"

	var content string
	if !m.status {
		content = lipgloss.NewStyle().Padding(1).Render("Server is OFF")
	} else {
		downloadedDelegate := &customDownloadedDelegate{}
		downloadedDelegate.SetHeight(2)
		downloadedDelegate.SetSpacing(1)

		m.downloadedList.AdditionalShortHelpKeys = func() []key.Binding {
			keyBinding := key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "load model"),
			)
			var keyBindings []key.Binding
			keyBindings = append(keyBindings, keyBinding)
			return keyBindings
		}
		m.downloadedList.SetShowTitle(false)
		m.downloadedList.SetShowStatusBar(false)
		m.downloadedList.SetFilteringEnabled(false)
		m.downloadedList.SetDelegate(downloadedDelegate)

		content = lipgloss.NewStyle().Padding(1).Render(m.downloadedList.View())
	}

	embeddedText := map[layout.BorderPosition]string{
		layout.TopLeftBorder: "[3] ▼ Downloaded Models",
	}

	return layout.Borderize(content, active, leftColumnWidth-2, (mainHeight-6)/2, embeddedText)
}
