package keybindings

import (
	"github.com/charmbracelet/bubbles/key"
)

func GetGlobalHelpKeys() []key.Binding {
	keyMap := DefaultGlobalKeyMap()
	return []key.Binding{
		keyMap.Quit,
		keyMap.NextView,
		keyMap.Help,
		keyMap.SystemPrompt,
		keyMap.ClearChat,
	}
}

func GetViewHelpKeys() []key.Binding {
	keyMap := DefaultViewKeyMap()
	return []key.Binding{
		keyMap.Status,
		keyMap.Loaded,
		keyMap.Downloaded,
		keyMap.Chat,
		keyMap.Logs,
	}
}

func GetChatHelpKeys() []key.Binding {
	keyMap := DefaultChatKeyMap()
	return []key.Binding{
		keyMap.SendMessage,
		keyMap.ScrollUp,
		keyMap.ScrollDown,
		keyMap.ExitChat,
		keyMap.Cancel,
	}
}

func GetListHelpKeys() []key.Binding {
	keyMap := DefaultListKeyMap()
	return []key.Binding{
		keyMap.Select,
		keyMap.Unload,
		keyMap.UnloadAll,
	}
}
