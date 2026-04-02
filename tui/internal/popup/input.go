// This file implements a text input dialog popup.
//
// Some module installations require user input (e.g., a Git email address,
// a custom path, or an API key). This dialog presents a text field inside
// a popup overlay where the user can type a value and submit it.
//
// # Composing Bubble Tea Models
//
// InputModel embeds both popup.Model (for the overlay frame) and
// textinput.Model (from the Bubbles component library) for the actual
// text field. This is a common pattern in Bubble Tea apps: compose larger
// UIs from smaller, reusable sub-models.
//
// The key to composition is forwarding messages: InputModel.Update() first
// checks for its own keys (enter, esc), then passes any remaining messages
// to textinput.Model.Update() so the text field can handle character input,
// cursor movement, etc.
// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput
package popup

import (
	// tea is the Bubble Tea TUI framework.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2
	tea "charm.land/bubbletea/v2"

	// textinput provides a single-line text input Bubble (sub-component).
	// It handles character input, cursor display, placeholder text, and more.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput
	"charm.land/bubbles/v2/textinput"

	"github.com/issafalcon/dotfiles-tui/internal/theme"
)

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

// InputSubmitMsg is sent when the user presses Enter to submit their input.
// Value contains the text that was typed into the field.
type InputSubmitMsg struct {
	Value string
}

// InputCancelMsg is sent when the user presses Escape to cancel.
type InputCancelMsg struct{}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

// InputModel is a text input dialog that embeds both the base popup Model
// and a Bubbles textinput.Model.
//
// Go allows embedding multiple types in a single struct. Each embedded type's
// exported methods are promoted, but if two embedded types have methods with
// the same name, neither is promoted (the compiler will report an ambiguity
// if you try to call it without qualifying the type).
// See: https://go.dev/ref/spec#Struct_types
type InputModel struct {
	// Model is the base popup (provides overlay frame, title, visibility).
	Model

	// prompt is the instructional text shown above the input field
	// (e.g., "Enter your Git email address:").
	prompt string

	// textInput is the Bubbles text input sub-component. It handles all
	// the low-level details of text editing: cursor movement, character
	// insertion/deletion, clipboard, placeholder rendering, etc.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model
	textInput textinput.Model

	// submitted tracks whether the user has pressed Enter.
	submitted bool
}

// NewInputDialog creates a new text input dialog.
//
// Parameters:
//   - title:  the popup title (displayed in bold pink at the top)
//   - prompt: instructional text shown above the text field
//
// The dialog starts visible with the text input focused (ready to accept
// keystrokes). The text field has a placeholder to hint at what to type.
func NewInputDialog(title, prompt string) InputModel {
	// textinput.New() creates a new text input model with sensible defaults.
	// It returns a value type (not a pointer).
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#New
	ti := textinput.New()

	// Placeholder is shown in dim text when the field is empty, giving
	// the user a hint about what to type.
	ti.Placeholder = "Type here..."

	// CharLimit restricts how many characters the user can enter.
	// 0 means no limit; we set a reasonable maximum.
	ti.CharLimit = 256

	// Focus() tells the textinput to start accepting keyboard input.
	// In Bubble Tea, only focused components process key events.
	// Focus() returns a tea.Cmd that starts the cursor blink timer.
	// Since we're in a constructor (not Update), we can't run this command
	// directly — the cursor blink will start on the first Update cycle.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Focus
	ti.Focus()

	return InputModel{
		Model:     NewPopup(title, "", 50, 10).Show(),
		prompt:    prompt,
		textInput: ti,
		submitted: false,
	}
}

// Update handles messages for the input dialog.
//
// The update flow is:
//  1. Check for our own action keys (enter to submit, esc to cancel)
//  2. If no action key was pressed, forward the message to the embedded
//     textinput so it can handle character input, cursor movement, etc.
//
// This "intercept then forward" pattern is fundamental to composing
// Bubble Tea models. The outer model gets first dibs on messages, and
// whatever it doesn't handle flows down to sub-components.
func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	// First, check for our own key bindings.
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		// Enter submits the current input value.
		case "enter":
			m.submitted = true

			// Capture the current value before creating the closure.
			// textinput.Value() returns the text currently in the field.
			// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Value
			value := m.textInput.Value()

			return m, func() tea.Msg {
				return InputSubmitMsg{Value: value}
			}

		// Escape cancels without submitting.
		case "esc":
			return m, func() tea.Msg {
				return InputCancelMsg{}
			}
		}
	}

	// --- Forward to textinput ---
	// For any message we didn't handle above (regular characters, backspace,
	// cursor movement, etc.), pass it to the textinput sub-component.
	//
	// textinput.Update() returns a new textinput.Model (value semantics) and
	// an optional tea.Cmd. We reassign m.textInput to the updated value.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.Update
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View builds the inner content of the input dialog.
//
// Layout:
//   - Prompt text (e.g., "Enter your Git email address:")
//   - Text input field (rendered by textinput.View())
//   - Hint line showing available actions
func (m InputModel) View() string {
	// Build the content by joining styled strings with newlines.
	// Each section uses a theme style for consistent visual treatment.
	content := theme.NormalText.Render(m.prompt) + "\n\n"

	// textinput.View() renders the text field with cursor, placeholder, etc.
	// It returns a plain string that we compose into our layout.
	// See: https://pkg.go.dev/charm.land/bubbles/v2/textinput#Model.View
	content += m.textInput.View() + "\n\n"

	// Show keyboard hints so the user knows how to interact.
	// theme.DimText renders in a muted colour so hints don't dominate.
	content += theme.DimText.Render("enter submit • esc cancel")

	return content
}

// Render produces the final string: dialog content wrapped in the themed
// popup box, centered on the terminal screen.
func (m InputModel) Render(screenWidth, screenHeight int) string {
	if !m.Model.IsVisible() {
		return ""
	}
	return renderPopup(
		m.Model.title, m.View(),
		m.Model.width, m.Model.height,
		screenWidth, screenHeight,
	)
}
