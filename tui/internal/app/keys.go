// Package app — keys.go defines all keyboard shortcuts for the application.
//
// This uses the Bubbles key package to define keybindings in a structured way.
// Each binding has:
//   - The actual keys that trigger it (e.g., "j", "down")
//   - Help text that's displayed in the help menu
//
// The key.Binding type works with the Bubbles help component to automatically
// generate a help view from your keybindings.
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/key
// See: https://pkg.go.dev/charm.land/bubbles/v2/help
package app

import (
	"charm.land/bubbles/v2/key"
)

// KeyMap defines all the keyboard shortcuts for the application.
// Each field is a key.Binding which associates one or more keys with an action.
//
// In Go, struct fields can have "tags" (the `json:"..."` or similar annotations).
// The Bubbles help component uses the help text from key.WithHelp() to build
// the help menu automatically.
//
// See: https://go.dev/ref/spec#Struct_types
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Select key.Binding

	// Actions
	Install   key.Binding
	Uninstall key.Binding
	Search    key.Binding
	OpenURL   key.Binding

	// UI
	SwitchTab key.Binding
	Help      key.Binding
	Cancel    key.Binding
	Quit      key.Binding
}

// DefaultKeyMap returns the default keybindings for the application.
//
// key.NewBinding() creates a new keybinding with:
//   - key.WithKeys(): The actual key(s) that trigger the binding
//   - key.WithHelp(): Short key name + description for the help menu
//
// See: https://pkg.go.dev/charm.land/bubbles/v2/key#NewBinding
var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install module"),
	),
	Uninstall: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "uninstall module"),
	),
	Search: key.NewBinding(
		key.WithKeys("s", "/"),
		key.WithHelp("s", "search modules"),
	),
	OpenURL: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open URL in browser"),
	),
	SwitchTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch tab"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel/close"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp returns the keybindings to show in the compact help view.
// This method satisfies the help.KeyMap interface from the Bubbles help package.
//
// In Go, interfaces are satisfied implicitly — you don't need to declare
// "implements". If a type has the right methods, it satisfies the interface.
// See: https://go.dev/doc/effective_go#interfaces
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Select, k.Install, k.Search, k.Help, k.Quit,
	}
}

// FullHelp returns the keybindings to show in the expanded help view.
// The outer slice creates groups (rendered as columns), inner slices are items.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},           // Navigation
		{k.Install, k.Uninstall, k.OpenURL}, // Actions
		{k.SwitchTab, k.Search},              // UI
		{k.Help, k.Cancel, k.Quit},           // App
	}
}
