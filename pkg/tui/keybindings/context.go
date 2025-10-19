package keybindings

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ChatKeyMap defines key bindings for the chat view
type ChatKeyMap struct {
	SendMessage key.Binding
	ScrollUp    key.Binding
	ScrollDown  key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Home        key.Binding
	End         key.Binding
	ExitChat    key.Binding
	Cancel      key.Binding
}

func DefaultChatKeyMap() ChatKeyMap {
	return ChatKeyMap{
		SendMessage: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "send message"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "scroll down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "pageup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "pagedown"),
			key.WithHelp("pgdown", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to bottom"),
		),
		ExitChat: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit chat"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "cancel request"),
		),
	}
}

// ListKeyMap defines key bindings for list views
type ListKeyMap struct {
	Select    key.Binding
	Unload    key.Binding
	UnloadAll key.Binding
}

func DefaultListKeyMap() ListKeyMap {
	return ListKeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select model"),
		),
		Unload: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "unload"),
		),
		UnloadAll: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "unload all"),
		),
	}
}

func HandleChatKey(msg tea.KeyMsg, keyMap ChatKeyMap) (tea.Cmd, bool) {
	switch msg.String() {
	case keyMap.SendMessage.Keys()[0]:
		return nil, true
	case keyMap.ScrollUp.Keys()[0]:
		return nil, true // Scroll up
	case keyMap.ScrollDown.Keys()[0]:
		return nil, true // Scroll down
	case keyMap.PageUp.Keys()[0], keyMap.PageUp.Keys()[1]:
		return nil, true // Page up
	case keyMap.PageDown.Keys()[0], keyMap.PageDown.Keys()[1]:
		return nil, true // Page down
	case keyMap.Home.Keys()[0]:
		return nil, true // Go to top
	case keyMap.End.Keys()[0]:
		return nil, true // Go to bottom
	case keyMap.ExitChat.Keys()[0]:
		return nil, true // Exit chat
	case keyMap.Cancel.Keys()[0]:
		return nil, true // Cancel request
	}
	return nil, false
}

func HandleListKey(msg tea.KeyMsg, keyMap ListKeyMap) (tea.Cmd, bool) {
	switch msg.String() {
	case keyMap.Select.Keys()[0]:
		return nil, true // Select model
	case keyMap.Unload.Keys()[0]:
		return nil, true // Unload model
	case keyMap.UnloadAll.Keys()[0]:
		return nil, true // Unload all
	}
	return nil, false
}
