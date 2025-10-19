package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/beauNate/lazylms-macos/internal/api"
	"github.com/beauNate/lazylms-macos/internal/config"
)

// Mac OS 26 Liquid Glass color scheme
var (
	// Liquid Glass theme colors - translucent, depth-enhanced
	glassBackground = lipgloss.Color("#0A0A0AE6") // Semi-transparent dark
	glassAccent     = lipgloss.Color("#00D4FF")   // Bright cyan
	glassSecondary  = lipgloss.Color("#FF006E")   // Vibrant magenta
	glassText       = lipgloss.Color("#E8E8E8")   // Soft white
	glassDim        = lipgloss.Color("#666666")   // Dimmed text
	glassGlow       = lipgloss.Color("#00FFD4")   // Glow effect
)

// Liquid Glass style definitions
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(glassAccent).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(glassGlow)

	modelItemStyle = lipgloss.NewStyle().
		Foreground(glassText).
		PaddingLeft(2).
		MarginLeft(1)

	selectedModelStyle = lipgloss.NewStyle().
		Foreground(glassSecondary).
		Background(lipgloss.Color("#1A1A1A")).
		Bold(true).
		PaddingLeft(2).
		MarginLeft(1)

	inputPromptStyle = lipgloss.NewStyle().
		Foreground(glassAccent).
		Bold(true)

	inputTextStyle = lipgloss.NewStyle().
		Foreground(glassText)

	responseStyle = lipgloss.NewStyle().
		Foreground(glassText).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(glassDim)

	errorStyle = lipgloss.NewStyle().
		Foreground(glassSecondary).
		Bold(true)

	statusStyle = lipgloss.NewStyle().
		Foreground(glassDim).
		Italic(true)
)

type Model struct {
	client   *api.Client
	config   *config.Config
	models   []string
	selected int
	loading  bool
	error    string
	prompt   string
	response string
	mode     string // "list" or "chat"
}

func NewModel(client *api.Client, cfg *config.Config) Model {
	return Model{
		client:   client,
		config:   cfg,
		models:   []string{},
		selected: 0,
		loading:  true,
		mode:     "list",
	}
}

type modelsLoadedMsg []string
type responseMsg string
type errMsg error

func (m Model) Init() tea.Cmd {
	return m.loadModels()
}

func (m Model) loadModels() tea.Cmd {
	return func() tea.Msg {
		models, err := m.client.ListModels()
		if err != nil {
			return errMsg(err)
		}
		return modelsLoadedMsg(models)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.mode == "list" && m.selected > 0 {
				m.selected--
			}

		case "down", "j":
			if m.mode == "list" && m.selected < len(m.models)-1 {
				m.selected++
			}

		case "enter":
			if m.mode == "list" && len(m.models) > 0 {
				m.mode = "chat"
				m.prompt = ""
				m.response = ""
			}

		case "esc":
			if m.mode == "chat" {
				m.mode = "list"
				m.prompt = ""
				m.response = ""
			}
		}

	case modelsLoadedMsg:
		m.models = []string(msg)
		m.loading = false
		if len(m.models) > 0 {
			m.selected = 0
		}

	case responseMsg:
		m.response = string(msg)

	case errMsg:
		m.error = msg.Error()
		m.loading = false
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Title with Liquid Glass styling
	b.WriteString(titleStyle.Render("🔮 lazylms-macos"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(statusStyle.Render("⚡ Loading models..."))
		b.WriteString("\n")
		return b.String()
	}

	if m.error != "" {
		b.WriteString(errorStyle.Render(fmt.Sprintf("❌ Error: %s", m.error)))
		b.WriteString("\n\n")
		b.WriteString(statusStyle.Render("Press q to quit"))
		b.WriteString("\n")
		return b.String()
	}

	if m.mode == "list" {
		b.WriteString(statusStyle.Render(fmt.Sprintf("Connected to %s:%d", m.config.Host, m.config.Port)))
		b.WriteString("\n\n")
		b.WriteString(inputPromptStyle.Render("Available Models:"))
		b.WriteString("\n\n")

		for i, model := range m.models {
			if i == m.selected {
				b.WriteString(selectedModelStyle.Render("→ " + model))
			} else {
				b.WriteString(modelItemStyle.Render("  " + model))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(statusStyle.Render("↑/↓: navigate • enter: select • q: quit"))
	} else {
		b.WriteString(inputPromptStyle.Render(fmt.Sprintf("Model: %s", m.models[m.selected])))
		b.WriteString("\n\n")

		if m.response != "" {
			b.WriteString(responseStyle.Render(m.response))
			b.WriteString("\n\n")
		}

		b.WriteString(statusStyle.Render("esc: back • q: quit"))
	}

	b.WriteString("\n")
	return b.String()
}
