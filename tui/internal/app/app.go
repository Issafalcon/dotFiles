// Package app contains the root application model for the DotFiles TUI.
//
// This is the top-level Bubble Tea model that owns all sub-models and
// orchestrates the overall application flow. It follows The Elm Architecture:
//
//   - Model: The App struct holds all application state
//   - Init(): Sets up initial state and kicks off prerequisite checks
//   - Update(): Routes messages to the appropriate sub-model
//   - View(): Composes the full UI from sub-model views
//
// # Struct Embedding and Composition
//
// Go doesn't have inheritance. Instead, it uses composition — you embed
// structs inside other structs. The root App model "owns" sub-models for
// the sidebar, detail panel, popup layer, etc. Each sub-model handles its
// own Update/View cycle, and the root model delegates messages to them.
//
// See: https://go.dev/doc/effective_go#embedding
// See: https://pkg.go.dev/charm.land/bubbletea/v2
package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/detail"
	"github.com/issafalcon/dotfiles-tui/internal/installer"
	"github.com/issafalcon/dotfiles-tui/internal/module"
	// Blank import to trigger all module init() registrations.
	// In Go, importing a package with _ means "run its init() functions but
	// don't use any of its exported names." This is how plugins/modules
	// self-register at startup.
	// See: https://go.dev/doc/effective_go#blank_import
	_ "github.com/issafalcon/dotfiles-tui/internal/module/modules"
	"github.com/issafalcon/dotfiles-tui/internal/popup"
	"github.com/issafalcon/dotfiles-tui/internal/prereqs"
	"github.com/issafalcon/dotfiles-tui/internal/sidebar"
	"github.com/issafalcon/dotfiles-tui/internal/theme"
	"github.com/issafalcon/dotfiles-tui/internal/utils"
)

// AppState represents which screen/phase the application is in.
// In Go, we use custom types based on int to create enumerations.
// The iota keyword auto-increments: PrereqCheck=0, Dashboard=1, Installing=2.
// See: https://go.dev/ref/spec#Iota
type AppState int

const (
	// StatePrereqCheck shows the prerequisites checking screen.
	StatePrereqCheck AppState = iota
	// StateDashboard shows the main module browsing interface.
	StateDashboard
	// StateInstalling indicates a module is being installed.
	StateInstalling
)

// FocusArea tracks which panel has keyboard focus in the dashboard.
type FocusArea int

const (
	FocusSidebar FocusArea = iota
	FocusDetail
)

// ProgramReadyMsg is sent from main.go after the *tea.Program is created.
// This provides the program reference needed for streaming install output
// via p.Send() from background goroutines.
type ProgramReadyMsg struct {
	Program *tea.Program
}

// Model is the root application model. It holds all state for the TUI app.
//
// In Bubble Tea, the model is any type that implements the tea.Model interface:
//
//	type Model interface {
//	    Init() Cmd
//	    Update(Msg) (Model, Cmd)
//	    View() View
//	}
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Model
type Model struct {
	// state tracks which screen we're currently showing.
	state AppState

	// focus tracks which panel currently has keyboard focus.
	focus FocusArea

	// Terminal dimensions, updated on resize events.
	// These are used to calculate the floating window size (~80% of terminal).
	width  int
	height int

	// ready indicates we've received the initial WindowSizeMsg.
	// Bubble Tea sends this automatically when the program starts.
	ready bool

	// --- Sub-models ---
	// Each sub-model handles its own Update/View cycle.
	// The root model delegates messages to the appropriate sub-model
	// based on the current state and focus area.

	prereqModel  prereqs.Model   // Prerequisites checking screen
	sidebarModel sidebar.Model   // Left panel: module list
	detailModel  detail.Model    // Right panel: tabs (overview/output/config)
	helpPopup    popup.HelpModel // Help overlay (? key)
	confirmPopup popup.ConfirmModel // Install confirmation dialog
	inputPopup   popup.InputModel   // User input dialog

	// --- State ---
	showHelp    bool   // Whether the help overlay is visible
	showConfirm bool   // Whether the confirm dialog is visible
	showInput   bool   // Whether the input dialog is visible
	selectedMod string // Currently selected module name

	// --- Streaming install state ---
	// The program reference is needed to call p.Send() from background
	// goroutines that stream install output to the Output pane.
	program *tea.Program

	// These fields track the install/uninstall sequence when sudo pre-auth is needed.
	// After RunSudoAuth completes, these are used to start the streaming operation.
	installingMod      string              // module currently being installed/uninstalled
	installCommands    []string            // full list of commands to run
	installStowEnabled bool                // whether to stow/unstow after all commands finish
	pendingAction      popup.ConfirmAction // tracks whether sudo pre-auth is for install or uninstall
}

// NewModel creates and returns the initial application model.
//
// In Go, constructor functions are conventionally named New<Type> or New.
// They return an initialized struct, since Go doesn't have constructors.
// See: https://go.dev/doc/effective_go#composite_literals
func NewModel() Model {
	// Build the sidebar items from the module registry.
	// The registry was populated by init() functions in the modules package.
	allModules := module.DefaultRegistry.All()
	installedModules, _ := utils.GetInstalledModules()
	installedSet := make(map[string]bool)
	for _, name := range installedModules {
		installedSet[name] = true
	}

	sidebarItems := make([]sidebar.ModuleItem, 0, len(allModules))
	for _, mod := range allModules {
		sidebarItems = append(sidebarItems, sidebar.ModuleItem{
			Name:        mod.Name,
			Icon:        mod.Icon,
			Description: mod.Description,
			Category:    mod.Category,
			Installed:   installedSet[mod.Name],
		})
	}

	return Model{
		state:        StatePrereqCheck,
		focus:        FocusSidebar,
		prereqModel:  prereqs.New(),
		sidebarModel: sidebar.NewModel(sidebarItems, 40, 30), // Sizes updated on first resize
		detailModel:  detail.NewModel(60, 30),                 // Sizes updated on first resize
		helpPopup:    popup.NewHelpPopup(nil),                 // Uses default bindings
	}
}

// Init is called once when the program starts. It returns an initial command.
//
// tea.Cmd is a function that performs I/O and returns a tea.Msg.
// tea.Batch() combines multiple commands to run concurrently — it takes
// any number of Cmds and returns a single Cmd that runs them all.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Cmd
// See: https://pkg.go.dev/charm.land/bubbletea/v2#Batch
func (m Model) Init() tea.Cmd {
	// Start the prerequisites check. The prereq model's Init() kicks off
	// async checks for each required tool and starts the spinner.
	return m.prereqModel.Init()
}

// Update is called whenever a message (event) arrives. It's the heart of TEA.
//
// Messages can be:
//   - tea.KeyPressMsg: A key was pressed
//   - tea.WindowSizeMsg: The terminal was resized
//   - Custom messages: Results from async operations (prereq checks, installs, etc.)
//
// The type switch (msg.(type)) is Go's way of checking which concrete type
// an interface value holds. This is called a "type assertion" or "type switch".
// See: https://go.dev/tour/methods/16
//
// Update returns the updated model and optionally a Cmd for more I/O.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// --- Global message handling (applies regardless of state) ---
	switch msg := msg.(type) {

	// ProgramReadyMsg provides the *tea.Program reference needed for
	// streaming install output via p.Send() from background goroutines.
	case ProgramReadyMsg:
		m.program = msg.Program
		return m, nil

	// tea.WindowSizeMsg is sent when the terminal is resized (and on startup).
	// We store the dimensions and propagate to sub-models so they resize too.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateSubModelSizes()
		return m, nil

	// tea.KeyPressMsg is sent when the user presses a key.
	case tea.KeyPressMsg:
		// If a popup is showing, handle its keys first.
		if m.showHelp {
			if msg.String() == "?" || msg.String() == "esc" || msg.String() == "q" {
				m.showHelp = false
				return m, nil
			}
			return m, nil // Consume all keys while help is open
		}

		if m.showConfirm {
			return m.updateConfirmPopup(msg)
		}

		if m.showInput {
			return m.updateInputPopup(msg)
		}

		// When sidebar search is active, only ctrl+c should work at app level.
		// All other keys must pass through to the sidebar's search input.
		if m.state == StateDashboard && m.sidebarModel.IsSearching() {
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.sidebarModel, cmd = m.sidebarModel.Update(msg)
			return m, cmd
		}

		// Global keys that work in any state.
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "?":
			m.showHelp = true
			return m, nil
		}

	// --- Cross-cutting messages from sub-models ---

	// PrereqsPassedMsg: All prerequisites met, transition to dashboard.
	case prereqs.PrereqsPassedMsg:
		m.state = StateDashboard
		// Select the first module to show its details.
		if sel := m.sidebarModel.Selected(); sel != "" {
			m.selectedMod = sel
			m.updateDetailForModule(sel)
		}
		return m, nil

	// ModuleSelectedMsg: User pressed enter on a module in the sidebar.
	case sidebar.ModuleSelectedMsg:
		m.selectedMod = msg.Name
		m.updateDetailForModule(msg.Name)
		m.focus = FocusDetail
		return m, nil

	// CursorChangedMsg: User navigated to a different module in the sidebar.
	case sidebar.CursorChangedMsg:
		m.selectedMod = msg.Name
		m.updateDetailForModule(msg.Name)
		return m, nil

	// ConfirmYesMsg: User confirmed an action (install or uninstall).
	// The Action field tells us which flow to execute.
	case popup.ConfirmYesMsg:
		m.showConfirm = false
		m.detailModel.SetActiveTab(detail.TabOutput)
		m.detailModel.OutputModel().Clear()

		mod := module.DefaultRegistry.Get(msg.ModuleName)
		if mod == nil {
			return m, nil
		}

		switch msg.Action {
		// --- Install flow ---
		case popup.ActionInstall:
			m.state = StateInstalling
			m.detailModel.OutputModel().SetInstalling(msg.ModuleName, true)

			// Handle stow-only modules (no install commands).
			// These only need symlinks created — no shell commands to run.
			if len(mod.InstallCommands) == 0 {
				if mod.StowEnabled {
					dotfilesDir := utils.GetDotfilesDir()
					if err := utils.Stow(msg.ModuleName, dotfilesDir); err != nil {
						m.detailModel.OutputModel().AppendLine(
							fmt.Sprintf("✗ Stow failed: %s", err))
					} else {
						m.detailModel.OutputModel().AppendLine("✓ Stow links created")
					}
				}
				_ = utils.SetModuleInstalled(msg.ModuleName)
				m.sidebarModel.SetInstalled(msg.ModuleName, true)
				m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
				m.detailModel.OutputModel().AppendLine(
					fmt.Sprintf("\n✓ %s installed successfully!", msg.ModuleName))
				m.state = StateDashboard
				return m, nil
			}

			// If commands need sudo, pre-authenticate first so the password
			// prompt can use the real terminal. After that, the streaming
			// install runs non-interactively with cached credentials.
			if installer.NeedsSudo(mod.InstallCommands) {
				m.installingMod = msg.ModuleName
				m.installCommands = mod.InstallCommands
				m.installStowEnabled = mod.StowEnabled
				return m, installer.RunSudoAuth(msg.ModuleName)
			}

			// No sudo needed — start streaming directly.
			return m, installer.RunInstallStreaming(
				m.program, msg.ModuleName, mod.InstallCommands,
				utils.GetDotfilesDir(), mod.StowEnabled)

		// --- Uninstall flow ---
		case popup.ActionUninstall:
			m.state = StateInstalling
			m.detailModel.OutputModel().SetInstalling(msg.ModuleName, true)

			// Handle stow-only modules (no uninstall commands).
			// Just remove symlinks and update tracking.
			if len(mod.UninstallCommands) == 0 {
				if mod.StowEnabled {
					dotfilesDir := utils.GetDotfilesDir()
					if err := utils.Unstow(msg.ModuleName, dotfilesDir); err != nil {
						m.detailModel.OutputModel().AppendLine(
							fmt.Sprintf("✗ Unstow failed: %s", err))
					} else {
						m.detailModel.OutputModel().AppendLine("✓ Stow links removed")
					}
				}
				_ = utils.SetModuleUninstalled(msg.ModuleName)
				m.sidebarModel.SetInstalled(msg.ModuleName, false)
				m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
				m.detailModel.OutputModel().AppendLine(
					fmt.Sprintf("\n✓ %s uninstalled successfully!", msg.ModuleName))
				m.state = StateDashboard
				return m, nil
			}

			// If uninstall commands need sudo, pre-authenticate first.
			if installer.NeedsSudo(mod.UninstallCommands) {
				m.installingMod = msg.ModuleName
				m.installCommands = mod.UninstallCommands
				m.installStowEnabled = mod.StowEnabled
				// Store the action so SudoAuthCompleteMsg knows whether to
				// install or uninstall after auth succeeds.
				m.pendingAction = popup.ActionUninstall
				return m, installer.RunSudoAuth(msg.ModuleName)
			}

			// No sudo needed — start streaming uninstall directly.
			return m, installer.RunUninstallStreaming(
				m.program, msg.ModuleName, mod.UninstallCommands,
				utils.GetDotfilesDir(), mod.StowEnabled)
		}
		return m, nil

	// SudoAuthCompleteMsg: sudo -v finished — now start the streaming operation.
	// This handles both install and uninstall flows, distinguished by m.pendingAction.
	case installer.SudoAuthCompleteMsg:
		if msg.Error != nil {
			m.state = StateDashboard
			m.detailModel.OutputModel().SetInstalling(m.installingMod, false)
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✗ sudo authentication failed: %s", msg.Error))
			m.installingMod = ""
			m.installCommands = nil
			m.pendingAction = ""
			return m, nil
		}
		if m.pendingAction == popup.ActionUninstall {
			m.pendingAction = ""
			return m, installer.RunUninstallStreaming(
				m.program, m.installingMod, m.installCommands,
				utils.GetDotfilesDir(), m.installStowEnabled)
		}
		// Default: install flow
		m.pendingAction = ""
		return m, installer.RunInstallStreaming(
			m.program, m.installingMod, m.installCommands,
			utils.GetDotfilesDir(), m.installStowEnabled)

	// ConfirmNoMsg: User cancelled the action.
	case popup.ConfirmNoMsg:
		m.showConfirm = false
		return m, nil

	// Install/Uninstall output messages — forward to the detail panel's output tab.
	case installer.InstallOutputMsg:
		m.detailModel.OutputModel().AppendLine(msg.Line)
		return m, nil

	// InstallCompleteMsg: All install commands finished (streaming or orchestrator path).
	case installer.InstallCompleteMsg:
		m.state = StateDashboard
		m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
		if msg.Success {
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✓ %s installed successfully!", msg.ModuleName))
			m.sidebarModel.SetInstalled(msg.ModuleName, true)
		} else {
			errMsg := "unknown error"
			if msg.Error != nil {
				errMsg = msg.Error.Error()
			}
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✗ %s installation failed: %s", msg.ModuleName, errMsg))
		}
		m.installingMod = ""
		m.installCommands = nil
		return m, nil

	// UninstallCompleteMsg: All uninstall commands finished.
	// Mirrors the install handler but updates the sidebar to mark the module
	// as NOT installed on success.
	case installer.UninstallCompleteMsg:
		m.state = StateDashboard
		m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
		if msg.Success {
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✓ %s uninstalled successfully!", msg.ModuleName))
			m.sidebarModel.SetInstalled(msg.ModuleName, false)
		} else {
			errMsg := "unknown error"
			if msg.Error != nil {
				errMsg = msg.Error.Error()
			}
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✗ %s uninstall failed: %s", msg.ModuleName, errMsg))
		}
		m.installingMod = ""
		m.installCommands = nil
		return m, nil
	}

	// --- State-specific message routing ---
	switch m.state {
	case StatePrereqCheck:
		return m.updatePrereqs(msg, cmds)
	case StateDashboard, StateInstalling:
		return m.updateDashboard(msg, cmds)
	}

	return m, nil
}

// updatePrereqs delegates messages to the prerequisites sub-model.
func (m Model) updatePrereqs(msg tea.Msg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	// Forward the message to the prereq model. It returns an updated model
	// and optionally a Cmd for more async work (e.g., next prereq check).
	var cmd tea.Cmd
	m.prereqModel, cmd = m.prereqModel.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// updateDashboard delegates messages to sidebar and detail sub-models
// based on which panel has focus.
func (m Model) updateDashboard(msg tea.Msg, cmds []tea.Cmd) (tea.Model, tea.Cmd) {
	// Handle dashboard-specific key presses.
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		// During search, forward all keys to sidebar — don't process dashboard shortcuts.
		if m.sidebarModel.IsSearching() {
			var cmd tea.Cmd
			m.sidebarModel, cmd = m.sidebarModel.Update(msg)
			return m, cmd
		}

		switch {
		// Search shortcut — works from any panel, not just the sidebar.
		case key.Matches(keyMsg, DefaultKeyMap.Search) && m.focus != FocusSidebar:
			m.focus = FocusSidebar
			cmd := m.sidebarModel.ActivateSearch()
			return m, cmd

		// Quit — only when not in search mode
		case key.Matches(keyMsg, DefaultKeyMap.Quit):
			return m, tea.Quit

		// Tab key with Shift swaps focus between sidebar and detail
		case keyMsg.String() == "shift+tab":
			if m.focus == FocusSidebar {
				m.focus = FocusDetail
			} else {
				m.focus = FocusSidebar
			}
			return m, nil

		// Install the selected module
		case key.Matches(keyMsg, DefaultKeyMap.Install):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil {
					// Build enriched list: module itself + dependencies with descriptions.
					items := []string{mod.Name + " — " + mod.Description}
					for _, depName := range mod.Dependencies {
						depMod := module.DefaultRegistry.Get(depName)
						if depMod != nil {
							items = append(items, depMod.Name+" — "+depMod.Description)
						} else {
							items = append(items, depName)
						}
					}
					m.confirmPopup = popup.NewConfirmDialog(m.selectedMod, items)
					m.showConfirm = true
					return m, nil
				}
			}

		// Uninstall the selected module — show a confirmation dialog first.
		// Builds an enriched list of what will be removed so the user knows
		// exactly which commands will run and whether stow links will be removed.
		case key.Matches(keyMsg, DefaultKeyMap.Uninstall):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil {
					// Build a summary of what the uninstall will do.
					var items []string
					items = append(items, mod.Name+" — "+mod.Description)
					if len(mod.UninstallCommands) > 0 {
						for _, cmd := range mod.UninstallCommands {
							items = append(items, "  ▸ "+cmd)
						}
					}
					if mod.StowEnabled {
						items = append(items, "  ▸ Remove stow symlinks")
					}

					m.confirmPopup = popup.NewUninstallDialog(m.selectedMod, items)
					m.showConfirm = true
					return m, nil
				}
			}

		// Open URL in browser
		case key.Matches(keyMsg, DefaultKeyMap.OpenURL):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil && mod.Website != "" {
					_ = utils.OpenURL(mod.Website)
				}
			}
			return m, nil
		}
	}

	// Forward messages to the focused sub-model.
	var cmd tea.Cmd
	switch m.focus {
	case FocusSidebar:
		m.sidebarModel, cmd = m.sidebarModel.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case FocusDetail:
		m.detailModel, cmd = m.detailModel.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// updateConfirmPopup handles messages for the confirmation dialog.
func (m Model) updateConfirmPopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.confirmPopup, cmd = m.confirmPopup.Update(msg)
	return m, cmd
}

// updateInputPopup handles messages for the input dialog.
func (m Model) updateInputPopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.inputPopup, cmd = m.inputPopup.Update(msg)
	return m, cmd
}

// updateSubModelSizes recalculates and propagates sizes to sub-models
// when the terminal is resized.
func (m *Model) updateSubModelSizes() {
	contentWidth, contentHeight := m.contentDimensions()

	// Sidebar gets ~30% of width.
	sidebarWidth := int(float64(contentWidth) * 0.3)
	// Detail gets the rest minus a gap.
	detailWidth := contentWidth - sidebarWidth - 1

	m.sidebarModel.SetSize(sidebarWidth, contentHeight-4)     // -4 for title + help bar
	m.detailModel.SetSize(detailWidth, contentHeight-4)

	// Calculate the sidebar's Y offset in the terminal for mouse support.
	// The floating window (contentHeight + 2 for border) is centered vertically.
	// Inside: 1 (top border) + ~3 lines (title with MarginBottom + "\n\n" gap).
	windowHeight := contentHeight + 2
	yPad := (m.height - windowHeight) / 2
	m.sidebarModel.SetYOffset(yPad + 1 + 3)
}

// updateDetailForModule updates the detail panel to show info for the given module.
func (m *Model) updateDetailForModule(name string) {
	mod := module.DefaultRegistry.Get(name)
	if mod == nil {
		return
	}

	// Build dependency status list for the overview tab.
	deps := make([]detail.DepStatus, 0, len(mod.ExternalDeps))
	for _, dep := range mod.ExternalDeps {
		deps = append(deps, detail.DepStatus{
			Name:     dep.Name,
			Method:   dep.InstallMethod,
			Installed: utils.IsCommandAvailable(dep.Name),
			Checking: false,
		})
	}

	m.detailModel.OverviewModel().SetModule(
		mod.Name,
		mod.Description,
		mod.Website,
		mod.Repo,
		deps,
	)

	// Build config options for the config tab.
	configOpts := make([]detail.ConfigOption, 0, len(mod.ConfigOptions))
	for _, opt := range mod.ConfigOptions {
		configOpts = append(configOpts, detail.ConfigOption{
			Name:        opt.Name,
			Description: opt.Description,
			Default:     opt.Default,
			Choices:     opt.Choices,
			Selected:    opt.Default,
		})
	}
	m.detailModel.ConfigModel().SetModule(mod.Name, configOpts)
}

// contentDimensions returns the inner content width and height
// for the floating window (~80% of terminal).
//
// The returned height is the usable interior space INSIDE the AppBorder.
// AppBorder adds a 1-cell border on each side (2 rows, 2 columns), so we
// subtract that overhead from the 80% allocation to prevent the bottom
// border from being cut off.
func (m Model) contentDimensions() (int, int) {
	contentWidth := int(float64(m.width) * 0.8)
	contentHeight := int(float64(m.height) * 0.8)

	// Subtract border overhead (top + bottom = 2 rows, left + right = 2 cols)
	// so that AppBorder.Width/Height refer to the interior and the total
	// rendered box still fits within the 80% allocation.
	contentWidth -= 2
	contentHeight -= 2

	if contentWidth < 60 {
		contentWidth = 60
	}
	if contentHeight < 20 {
		contentHeight = 20
	}
	return contentWidth, contentHeight
}

// View renders the entire UI as a string. Bubble Tea calls this after every Update.
//
// IMPORTANT: View() must be a pure function — it should only read from the model,
// never modify it or perform I/O. All side effects happen in Update() via Cmds.
//
// See: https://pkg.go.dev/charm.land/bubbletea/v2#View
func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	contentWidth, contentHeight := m.contentDimensions()

	// Build the content based on current state.
	var content string
	switch m.state {
	case StatePrereqCheck:
		content = m.prereqModel.View()
	case StateDashboard, StateInstalling:
		content = m.viewDashboard(contentWidth, contentHeight)
	}

	// Force content to exact fixed dimensions before applying the border.
	// Height sets minimum (pads short content), MaxHeight truncates overflow.
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		MaxWidth(contentWidth).
		MaxHeight(contentHeight)
	content = contentStyle.Render(content)

	// Apply the floating window border style.
	window := theme.AppBorder.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)

	// Overlay popups on top of the window if visible.
	finalView := window
	if m.showHelp {
		// Render help popup over the main window.
		finalView = m.helpPopup.Render(contentWidth+2, contentHeight+2)
	}
	if m.showConfirm {
		finalView = m.confirmPopup.Render(contentWidth+2, contentHeight+2)
	}
	if m.showInput {
		finalView = m.inputPopup.Render(contentWidth+2, contentHeight+2)
	}

	// Center the floating window in the terminal using lipgloss.Place().
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#Place
	view := tea.NewView(
		lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			finalView,
		),
	)

	// In Bubble Tea v2, AltScreen and MouseMode are set on the View struct.
	// See: https://pkg.go.dev/charm.land/bubbletea/v2#View
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

// viewDashboard renders the main module browsing dashboard with sidebar + detail.
func (m Model) viewDashboard(width, height int) string {
	title := theme.Title.Render("⚡ DotFiles Manager")

	// Help bar at the bottom showing available shortcuts.
	help := theme.HelpStyle.Render(
		"q: quit • ?: help • j/k: navigate • shift+tab: switch panel • tab: switch tab • i: install • d: uninstall • o: open URL • s: search",
	)

	// Reserve vertical space for title (1 line + MarginBottom 1 + "\n\n" = ~3 lines)
	// and help bar (~1 line + "\n\n" gap = ~3 lines). Total overhead ≈ 5 lines.
	panelHeight := height - 5
	if panelHeight < 10 {
		panelHeight = 10
	}

	// Force sidebar and detail views to exact heights using Height (min) + MaxHeight (max).
	panelStyle := lipgloss.NewStyle().
		Height(panelHeight).
		MaxHeight(panelHeight)

	sidebarView := panelStyle.Render(m.sidebarModel.View())
	detailView := panelStyle.Render(m.detailModel.View())

	// Render the sidebar and detail panel side by side.
	// lipgloss.JoinHorizontal places rendered strings horizontally.
	// See: https://pkg.go.dev/charm.land/lipgloss/v2#JoinHorizontal
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, " ", detailView)

	// Force body to fixed height so it doesn't vary with content.
	body = lipgloss.NewStyle().
		Height(panelHeight).
		MaxHeight(panelHeight).
		Render(body)

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, body, help)
}
