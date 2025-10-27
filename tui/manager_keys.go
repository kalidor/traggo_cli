package tui

import "github.com/charmbracelet/bubbles/key"

// mainKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type mainKeyMap struct {
	E     key.Binding // Edit
	R     key.Binding // Refresh
	C     key.Binding // Continue
	D     key.Binding // Delete
	Tab   key.Binding // Switch to Current/Complete tasks
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.C, k.D, k.E, k.R, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Up, k.Down}, // first column
		{k.C, k.D, k.E, k.R},  // second column
		{k.Help, k.Quit},      // third column
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
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "current/complete"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
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
