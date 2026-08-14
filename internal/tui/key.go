// Package tui contains terminal user interface commands and related functionality.
package tui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines key bindings for navigating the TUI.
type keyMap struct {
	Start  key.Binding
	Rename key.Binding
	Kill   key.Binding
	Add    key.Binding
	Save   key.Binding
	Escape key.Binding
	Enter  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Filter key.Binding
}

// defaultKeyMap provides the default key bindings for moving up and down in the TUI.
var defaultKeyMap = keyMap{
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start session"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Kill: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "kill session"),
	),
	Save: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "save to config"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "switch"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit input"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Start}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Rename},
		{k.Save, k.Kill},
		{k.Enter, k.Filter},
	}
}
