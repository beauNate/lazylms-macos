package keybindings

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// GlobalKeyMap defines global key bindings for the application
type GlobalKeyMap struct {
	Quit         key.Binding
	NextView     key.Binding
	PrevView     key.Binding
	Help         key.Binding
	SystemPrompt key.Binding
	ClearChat    key.Binding
}

func DefaultGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		NextView: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next panel"),
		),
		PrevView: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous panel"),
		),
		Help: key.NewBinding(
			key.WithKeys("h", "?"),
			key.WithHelp("h", "help"),
		),
		SystemPrompt: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "system prompt"),
		),
		ClearChat: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "clear chat"),
		),
	}
}

// ViewKeyMap defines key bindings for view navigation
type ViewKeyMap struct {
	Status     key.Binding
	Loaded     key.Binding
	Downloaded key.Binding
	Chat       key.Binding
	Logs       key.Binding
}

func DefaultViewKeyMap() ViewKeyMap {
	return ViewKeyMap{
		Status: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "status"),
		),
		Loaded: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "loaded models"),
		),
		Downloaded: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "downloaded models"),
		),
		Chat: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "chat"),
		),
		Logs: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "logs"),
		),
	}
}

func HandleGlobalKey(msg tea.KeyMsg, keyMap GlobalKeyMap) (tea.Cmd, bool) {
	switch msg.String() {
	case keyMap.Quit.Keys()[0]:
		return tea.Quit, true
	case keyMap.NextView.Keys()[0]:
		return nil, true // Next view command
	case keyMap.Help.Keys()[0], keyMap.Help.Keys()[1]:
		return nil, true // Toggle help
	case keyMap.SystemPrompt.Keys()[0]:
		return nil, true // Toggle system prompt
	case keyMap.ClearChat.Keys()[0]:
		return nil, true
	}
	return nil, false
}
