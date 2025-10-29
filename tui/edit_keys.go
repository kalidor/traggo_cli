package tui

import "github.com/charmbracelet/bubbles/key"

type editKeyMap struct {
	Tab   key.Binding // Switch to Current/Complete tasks
	CtrlC key.Binding
	CtrlL key.Binding
	Esc   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k editKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.CtrlC, k.CtrlL, k.Esc}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k editKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.CtrlC, k.CtrlL, k.Esc},
	}
}

var editKeys = editKeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab/MajTab/↑/↓", "change focus"),
	),
	CtrlC: key.NewBinding(
		key.WithKeys("Ctrl+c"),
		key.WithHelp("Ctrl+c", "Quit"),
	),
	CtrlL: key.NewBinding(
		key.WithKeys("Ctrl+l"),
		key.WithHelp("Ctrl+l", "Clear"),
	),
	Esc: key.NewBinding(
		key.WithKeys("Esc"),
		key.WithHelp("Esc", "Go back"),
	),
}
