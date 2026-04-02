// Package theme defines the visual styling for the DotFiles TUI application.
//
// It uses Lip Gloss (https://pkg.go.dev/charm.land/lipgloss/v2) to create
// CSS-like style definitions for terminal output. Lip Gloss styles are
// immutable value types — each method call returns a new style rather than
// mutating the original. This is similar to how strings work in Go.
//
// The colour scheme is inspired by the Dracula theme (https://draculatheme.com/)
// with a cyberpunk/neon twist: deep purples, hot pinks, electric cyan, and
// neon green on dark backgrounds.
//
// # Lip Gloss Basics
//
// Styles are created with lipgloss.NewStyle() and chained with method calls:
//
//	style := lipgloss.NewStyle().
//	    Bold(true).
//	    Foreground(lipgloss.Color("#FF79C6"))
//
// Then rendered with style.Render("some text").
//
// See: https://pkg.go.dev/charm.land/lipgloss/v2
package theme

import (
	lipgloss "charm.land/lipgloss/v2"
)

// ---------------------------------------------------------------------------
// Colour Palette
// ---------------------------------------------------------------------------
// These are the core colours used throughout the app.
// lipgloss.Color() accepts hex strings (#RRGGBB), ANSI 256 colour codes ("201"),
// or ANSI 16 colour names ("5"). Lip Gloss automatically downgrades colours
// for terminals that don't support true colour.
// See: https://pkg.go.dev/charm.land/lipgloss/v2#Color

var (
	// Background colours
	ColorBackground     = lipgloss.Color("#282A36") // Deep dark purple (Dracula BG)
	ColorSurface        = lipgloss.Color("#44475A") // Slightly lighter surface
	ColorSurfaceHigh    = lipgloss.Color("#6272A4") // Highlighted surface / comments

	// Primary accent colours
	ColorPink           = lipgloss.Color("#FF79C6") // Hot pink — primary accent
	ColorPurple         = lipgloss.Color("#BD93F9") // Soft purple — secondary accent
	ColorCyan           = lipgloss.Color("#8BE9FD") // Electric cyan — highlights
	ColorGreen          = lipgloss.Color("#50FA7B") // Neon green — success states
	ColorRed            = lipgloss.Color("#FF5555") // Bright red — errors
	ColorYellow         = lipgloss.Color("#F1FA8C") // Yellow — warnings
	ColorOrange         = lipgloss.Color("#FFB86C") // Orange — attention

	// Text colours
	ColorForeground     = lipgloss.Color("#F8F8F2") // Primary text (light)
	ColorForegroundDim  = lipgloss.Color("#6272A4") // Dimmed text (comments)
	ColorForegroundMuted = lipgloss.Color("#44475A") // Very dim text
)

// ---------------------------------------------------------------------------
// Shared Styles
// ---------------------------------------------------------------------------
// These styles are used throughout the app to maintain visual consistency.
// Each style is a lipgloss.Style value type — assigning it to another variable
// creates a true copy, so you can safely extend styles without mutation.
// See: https://pkg.go.dev/charm.land/lipgloss/v2#Style

var (
	// AppBorder is the outermost border style for the floating window.
	// lipgloss.RoundedBorder() returns a Border struct with rounded corners (╭╮╰╯).
	// BorderForeground sets the colour of the border characters.
	AppBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPurple)

	// Title renders section titles with bold pink text.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPink).
		MarginBottom(1)

	// Subtitle renders secondary headings.
	Subtitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorCyan)

	// NormalText is the default text style.
	NormalText = lipgloss.NewStyle().
		Foreground(ColorForeground)

	// DimText is for less important information.
	DimText = lipgloss.NewStyle().
		Foreground(ColorForegroundDim)

	// SuccessText indicates something positive (e.g., "installed").
	SuccessText = lipgloss.NewStyle().
		Foreground(ColorGreen)

	// ErrorText indicates an error.
	ErrorText = lipgloss.NewStyle().
		Foreground(ColorRed)

	// WarningText indicates a warning or caution.
	WarningText = lipgloss.NewStyle().
		Foreground(ColorYellow)

	// ActiveTab style for the currently selected tab.
	ActiveTab = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPink).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorPink).
		Padding(0, 2)

	// InactiveTab style for tabs that aren't selected.
	InactiveTab = lipgloss.NewStyle().
		Foreground(ColorForegroundDim).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorSurface).
		Padding(0, 2)

	// SidebarItem renders a module entry in the sidebar.
	SidebarItem = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSurface).
		Padding(0, 1).
		MarginBottom(0)

	// SidebarItemActive is for the currently highlighted sidebar item.
	// Uses a bold border, bright background tint, and prominent colours
	// so the active item is unmistakably distinct from inactive items.
	SidebarItemActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPink).
		Bold(true).
		Foreground(ColorForeground).
		Padding(0, 1).
		MarginBottom(0)

	// SidebarItemInstalled is for modules that are already installed.
	SidebarItemInstalled = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorGreen).
		Padding(0, 1).
		MarginBottom(0)

	// StatusInstalled shows a green tick for installed items.
	StatusInstalled = lipgloss.NewStyle().
		Foreground(ColorGreen).
		SetString("✓")

	// StatusNotInstalled shows a red cross for missing items.
	StatusNotInstalled = lipgloss.NewStyle().
		Foreground(ColorRed).
		SetString("✗")

	// StatusChecking shows a spinner-like indicator for async checks.
	StatusChecking = lipgloss.NewStyle().
		Foreground(ColorYellow).
		SetString("⟳")

	// PopupStyle renders modal dialogs/popups.
	PopupStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPink).
		Padding(1, 2)

	// HelpStyle renders the help bar at the bottom.
	HelpStyle = lipgloss.NewStyle().
		Foreground(ColorForegroundDim).
		Padding(0, 1)

	// ProgressBarFilled is the style for the filled portion of progress bars.
	ProgressBarFilled = lipgloss.NewStyle().
		Foreground(ColorCyan)

	// ProgressBarEmpty is the style for the empty portion of progress bars.
	ProgressBarEmpty = lipgloss.NewStyle().
		Foreground(ColorSurface)

	// SearchBar renders the search input.
	SearchBar = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPurple).
		Padding(0, 1)

	// URLStyle renders clickable-looking URLs.
	URLStyle = lipgloss.NewStyle().
		Foreground(ColorCyan).
		Underline(true)

	// KeyStyle renders keyboard shortcuts in the help view.
	KeyStyle = lipgloss.NewStyle().
		Foreground(ColorPink).
		Bold(true)

	// DescStyle renders descriptions next to keyboard shortcuts.
	DescStyle = lipgloss.NewStyle().
		Foreground(ColorForegroundDim)
)

// Icon constants using Nerd Font glyphs.
// These require a Nerd Font to be installed in the terminal.
// See: https://www.nerdfonts.com/
const (
	IconInstalled   = "✓"
	IconNotInstalled = "✗"
	IconChecking    = "⟳"
	IconSearch      = ""
	IconFolder      = ""
	IconPackage     = ""
	IconTerminal    = ""
	IconGear        = ""
	IconWarning     = ""
	IconError       = ""
	IconInfo        = ""
	IconArrowRight  = ""
)

// ModuleIcons maps module names to their Nerd Font icons.
// If a module doesn't have a specific icon, a default package icon is used.
var ModuleIcons = map[string]string{
	"nvim":         "",
	"zsh":          "",
	"git":          "",
	"tmux":         "",
	"docker":       "",
	"node":         "",
	"python":       "",
	"go":           "",
	"rust":         "",
	"lua":          "",
	"kubernetes":   "󱃾",
	"terraform":    "󱁢",
	"aws":          "",
	"azure":        "󰠅",
	"google-cloud": "",
	"homebrew":     "",
	"fzf":          "",
	"ranger":       "",
	"yazi":         "",
	"lazygit":      "",
	"wezterm":      "",
	"postman":      "",
	"obsidian":     "󰮏",
	"latex":        "",
	"mysql":        "",
	"helm":         "󱃾",
	"vagrant":      "",
	"packer":       "",
	"dotnet":       "󰪮",
	"cpp":          "",
	"powershell":   "󰨊",
}

// GetModuleIcon returns the Nerd Font icon for a given module name.
// If no specific icon exists, it returns a default package icon.
//
// In Go, the "comma ok" idiom is used to check if a map key exists:
//
//	value, ok := myMap[key]
//
// If ok is true, the key was found. If false, value is the zero value for its type.
// See: https://go.dev/doc/effective_go#maps
func GetModuleIcon(name string) string {
	if icon, ok := ModuleIcons[name]; ok {
		return icon
	}
	return IconPackage
}
