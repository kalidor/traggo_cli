package tui

import "github.com/charmbracelet/bubbles/key"

type searchKeyMap struct {
	Tab   key.Binding // Switch to Current/Complete tasks
	CtrlC key.Binding
	CtrlS key.Binding
	CtrlL key.Binding
	Esc   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k searchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CtrlC, k.CtrlS, k.CtrlL, k.Esc}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k searchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CtrlC, k.CtrlS, k.CtrlL, k.Esc},
	}
}

var searchKeys = searchKeyMap{
	CtrlC: key.NewBinding(
		key.WithKeys("Ctrl+c"),
		key.WithHelp("Ctrl+c", "Quit"),
	),
	CtrlS: key.NewBinding(
		key.WithKeys("Ctrl+s"),
		key.WithHelp("Ctrl+s", "Toogle search case"),
	),
	CtrlL: key.NewBinding(
		key.WithKeys("Ctrl+l"),
		key.WithHelp("Ctrl+l", "Clear search"),
	),
	Esc: key.NewBinding(
		key.WithKeys("Esc"),
		key.WithHelp("Esc", "Keep search result & go back"),
	),
}
