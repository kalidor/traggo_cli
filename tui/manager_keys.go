package tui

import "github.com/charmbracelet/bubbles/key"

// mainKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type mainKeyMap struct {
	E     key.Binding // Edit
	R     key.Binding // Refresh
	C     key.Binding // Continue
	D     key.Binding // Delete
	P     key.Binding // Period search
	S     key.Binding // Stop task
	N     key.Binding // New task
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Enter key.Binding
	Quit  key.Binding
	Slash key.Binding // search
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.N, k.S, k.C, k.D, k.E, k.Slash, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.N, k.S, k.Up, k.Down}, // first column
		{k.C, k.D, k.E, k.R},     // second column
		{k.Help, k.Quit},         // third column
	}
}

var mainKeys = mainKeyMap{
	C: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "continue"),
	),
	D: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	E: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "show/hide"),
	),
	R: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	N: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	S: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Slash: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	P: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "period search"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
