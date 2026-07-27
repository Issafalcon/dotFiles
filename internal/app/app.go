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
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/issafalcon/dotfiles-tui/internal/config"
	"github.com/issafalcon/dotfiles-tui/internal/detail"
	"github.com/issafalcon/dotfiles-tui/internal/docs"
	"github.com/issafalcon/dotfiles-tui/internal/installer"
	"github.com/issafalcon/dotfiles-tui/internal/module"
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
	// StateModulesSetup asks for the modules directory on first run.
	StateModulesSetup
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
	helpPopup     popup.HelpModel     // Help overlay (? key)
	confirmPopup  popup.ConfirmModel  // Install confirmation dialog
	scriptPopup   popup.ScriptModel   // Script review overlay
	categoryPopup popup.CategoryModel // Category filter picker
	inputPopup    popup.InputModel    // User input dialog

	// --- State ---
	showHelp     bool   // Whether the help overlay is visible
	showConfirm  bool   // Whether the confirm dialog is visible
	showScript   bool   // Whether the script review popup is visible
	showCategory bool   // Whether the category picker is visible
	showInput    bool   // Whether the input dialog is visible
	selectedMod  string // Currently selected module name

	// --- Streaming install state ---
	// The program reference is needed to call p.Send() from background
	// goroutines that stream install output to the Output pane.
	program *tea.Program

	// These fields track the install/uninstall sequence when sudo pre-auth is needed.
	// After RunSudoAuth completes, these are used to start the streaming operation.
	installingMod      string              // module currently being installed/uninstalled
	installScriptPath  string              // path to install.sh / uninstall.sh (may be empty)
	installStowEnabled bool                // whether to stow/unstow after the script finishes
	pendingAction      popup.ConfirmAction // tracks whether sudo pre-auth is for install or uninstall
	installQueue       []string            // remaining modules to install (deps then target)
	installPlan        []string            // full ordered plan for progress labels

	loadWarnings []string // module.yaml load issues shown once on dashboard
	showDocs     bool
	docsPopup    popup.ScriptModel
}

func buildSidebarItems() []sidebar.ModuleItem {
	allModules := module.DefaultRegistry.All()
	installedModules, _ := utils.GetInstalledModules()
	installedSet := make(map[string]bool)
	for _, name := range installedModules {
		installedSet[name] = true
	}
	items := make([]sidebar.ModuleItem, 0, len(allModules))
	for _, mod := range allModules {
		items = append(items, sidebar.ModuleItem{
			Name:        mod.Name,
			Icon:        mod.Icon,
			Description: mod.Description,
			Category:    mod.Category,
			Installed:   installedSet[mod.Name],
		})
	}
	return items
}

// NewModel creates and returns the initial application model.
func NewModel() Model {
	dir := utils.GetModulesDir()
	var warnings []string
	if dir != "" {
		res := module.LoadFromDir(dir)
		warnings = res.Errors
	}

	return Model{
		state:        StatePrereqCheck,
		focus:        FocusSidebar,
		prereqModel:  prereqs.New(),
		sidebarModel: sidebar.NewModel(buildSidebarItems(), 40, 30),
		detailModel:  detail.NewModel(60, 30),
		helpPopup:    popup.NewHelpPopup(nil),
		loadWarnings: warnings,
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

		if m.showDocs {
			var cmd tea.Cmd
			m.docsPopup, cmd = m.docsPopup.Update(msg)
			return m, cmd
		}

		if m.showScript {
			return m.updateScriptPopup(msg)
		}

		if m.showCategory {
			return m.updateCategoryPopup(msg)
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

	// PrereqsPassedMsg: All prerequisites met — configure modules path or open dashboard.
	case prereqs.PrereqsPassedMsg:
		if utils.GetModulesDir() == "" {
			m.state = StateModulesSetup
			m.inputPopup = popup.NewInputDialog(
				"Modules directory",
				"Path to your modules/ folder (created if missing):",
			)
			m.showInput = true
			return m, nil
		}
		m.reloadModules()
		m.state = StateDashboard
		if sel := m.sidebarModel.Selected(); sel != "" {
			m.selectedMod = sel
			m.updateDetailForModule(sel)
		}
		return m, nil

	case popup.InputSubmitMsg:
		m.showInput = false
		if m.state == StateModulesSetup {
			path := expandHome(msg.Value)
			if path == "" {
				m.showInput = true
				m.inputPopup = popup.NewInputDialog(
					"Modules directory",
					"Path cannot be empty. Enter a modules/ folder path:",
				)
				return m, nil
			}
			if err := os.MkdirAll(path, 0o755); err != nil {
				m.showInput = true
				m.inputPopup = popup.NewInputDialog(
					"Modules directory",
					fmt.Sprintf("Could not create dir (%v). Try another path:", err),
				)
				return m, nil
			}
			if err := config.SetModulesDir(path); err != nil {
				m.showInput = true
				m.inputPopup = popup.NewInputDialog(
					"Modules directory",
					fmt.Sprintf("Could not save config (%v). Try again:", err),
				)
				return m, nil
			}
			m.reloadModules()
			m.state = StateDashboard
			if sel := m.sidebarModel.Selected(); sel != "" {
				m.selectedMod = sel
				m.updateDetailForModule(sel)
			}
			return m, nil
		}
		return m, nil

	case popup.InputCancelMsg:
		m.showInput = false
		if m.state == StateModulesSetup {
			m.showInput = true
			m.inputPopup = popup.NewInputDialog(
				"Modules directory",
				"A modules path is required. Enter a folder path:",
			)
			return m, nil
		}
		return m, nil

	// ModuleSelectedMsg: User pressed enter on a module in the sidebar.
	case sidebar.ModuleSelectedMsg:
		m.selectedMod = msg.Name
		m.updateDetailForModule(msg.Name)
		m.focus = FocusDetail
		m.sidebarModel.SetFocused(false)
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
		case popup.ActionInstall, popup.ActionReinstall:
			m.state = StateInstalling
			m.detailModel.OutputModel().SetInstalling(msg.ModuleName, true)

			var queue []string
			var err error
			if msg.Action == popup.ActionReinstall {
				queue = []string{msg.ModuleName}
				m.detailModel.OutputModel().AppendLine(
					fmt.Sprintf("Force re-run: %s (deps skipped)", msg.ModuleName))
			} else {
				queue, err = planInstallQueue(msg.ModuleName)
				if err != nil {
					m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
					m.detailModel.OutputModel().AppendLine(
						fmt.Sprintf("✗ Could not plan install: %s", err))
					m.state = StateDashboard
					return m, nil
				}
				if len(queue) == 0 {
					m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
					m.detailModel.OutputModel().AppendLine(
						fmt.Sprintf("✓ %s and its dependencies are already installed", msg.ModuleName))
					m.detailModel.OutputModel().AppendLine(
						"  Tip: press r to force re-run this module's install script")
					m.sidebarModel.SetInstalled(msg.ModuleName, true)
					m.state = StateDashboard
					return m, nil
				}
			}

			m.installPlan = queue
			m.installQueue = queue[1:]
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("Install plan (%d): %s", len(queue), strings.Join(queue, " → ")))
			return m.beginModuleInstall(queue[0])

		// --- Uninstall flow ---
		case popup.ActionUninstall:
			m.state = StateInstalling
			m.detailModel.OutputModel().SetInstalling(msg.ModuleName, true)

			scriptPath := ""
			if utils.ModuleScriptExists(msg.ModuleName, "uninstall.sh") {
				scriptPath = utils.ModuleScriptPath(msg.ModuleName, "uninstall.sh")
			}
			modulesDir := utils.GetModulesDir()

			// Stow-only uninstall (no uninstall.sh): just remove symlinks.
			if scriptPath == "" {
				if mod.StowEnabled {
					if err := utils.Unstow(msg.ModuleName, modulesDir); err != nil {
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

			if installer.NeedsSudoScript(scriptPath) {
				m.installingMod = msg.ModuleName
				m.installScriptPath = scriptPath
				m.installStowEnabled = mod.StowEnabled
				m.pendingAction = popup.ActionUninstall
				return m, installer.RunSudoAuth(msg.ModuleName)
			}

			return m, installer.RunUninstallStreaming(
				m.program, msg.ModuleName, scriptPath,
				modulesDir, mod.StowEnabled)
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
			if len(m.installQueue) > 0 {
				m.detailModel.OutputModel().AppendLine(
					fmt.Sprintf("✗ Skipping remaining: %s", strings.Join(m.installQueue, ", ")))
			}
			m.installQueue = nil
			m.installPlan = nil
			m.installingMod = ""
			m.installScriptPath = ""
			m.pendingAction = ""
			return m, nil
		}
		modulesDir := utils.GetModulesDir()
		if m.pendingAction == popup.ActionUninstall {
			m.pendingAction = ""
			return m, installer.RunUninstallStreaming(
				m.program, m.installingMod, m.installScriptPath,
				modulesDir, m.installStowEnabled)
		}
		// Default: install flow
		m.pendingAction = ""
		m.detailModel.OutputModel().AppendLine(
			fmt.Sprintf("▸ Starting %s after sudo…", m.installingMod))
		return m, installer.RunInstallStreaming(
			m.program, m.installingMod, m.installScriptPath,
			modulesDir, m.installStowEnabled)

	// ConfirmNoMsg: User cancelled the action.
	case popup.ConfirmNoMsg:
		m.showConfirm = false
		return m, nil

	// ConfirmReviewMsg: open script viewer; keep confirm state for return.
	case popup.ConfirmReviewMsg:
		scriptName := "install.sh"
		if msg.Action == popup.ActionUninstall {
			scriptName = "uninstall.sh"
		}
		path := utils.ModuleScriptPath(msg.ModuleName, scriptName)
		data, err := os.ReadFile(path)
		content := ""
		if err != nil {
			content = fmt.Sprintf("(could not read %s: %v)", path, err)
		} else {
			content = string(data)
		}
		m.scriptPopup = popup.NewScriptViewer(scriptName+" — "+msg.ModuleName, content)
		m.showScript = true
		return m, nil

	// ScriptDismissMsg: return from script viewer to confirm dialog.
	case popup.ScriptDismissMsg:
		if m.showDocs {
			m.showDocs = false
			return m, nil
		}
		m.showScript = false
		return m, nil

	case popup.CategorySelectedMsg:
		m.showCategory = false
		m.sidebarModel.SetCategoryFilter(msg.Category)
		if name := m.sidebarModel.Selected(); name != "" {
			m.selectedMod = name
			m.updateDetailForModule(name)
		}
		return m, nil

	case popup.CategoryCancelMsg:
		m.showCategory = false
		return m, nil

	// Install/Uninstall output messages — forward to the detail panel's output tab.
	case installer.InstallOutputMsg:
		m.detailModel.OutputModel().AppendLine(msg.Line)
		return m, nil

	// InstallCompleteMsg: one module finished (may continue dep queue).
	case installer.InstallCompleteMsg:
		m.detailModel.OutputModel().SetInstalling(msg.ModuleName, false)
		if msg.Success {
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("\n✓ %s installed successfully!", msg.ModuleName))
			m.sidebarModel.SetInstalled(msg.ModuleName, true)
			if m.selectedMod != "" {
				m.updateDetailForModule(m.selectedMod)
			}

			if len(m.installQueue) > 0 {
				next := m.installQueue[0]
				m.installQueue = m.installQueue[1:]
				m.state = StateInstalling
				return m.beginModuleInstall(next)
			}

			m.state = StateDashboard
			m.installingMod = ""
			m.installScriptPath = ""
			m.installPlan = nil
			return m, nil
		}

		errMsg := "unknown error"
		if msg.Error != nil {
			errMsg = msg.Error.Error()
		}
		m.detailModel.OutputModel().AppendLine(
			fmt.Sprintf("\n✗ %s installation failed: %s", msg.ModuleName, errMsg))
		if len(m.installQueue) > 0 {
			m.detailModel.OutputModel().AppendLine(
				fmt.Sprintf("✗ Skipping remaining: %s", strings.Join(m.installQueue, ", ")))
		}
		m.installQueue = nil
		m.installPlan = nil
		m.state = StateDashboard
		m.installingMod = ""
		m.installScriptPath = ""
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
		m.installScriptPath = ""
		return m, nil
	}

	// --- State-specific message routing ---
	switch m.state {
	case StatePrereqCheck:
		return m.updatePrereqs(msg, cmds)
	case StateModulesSetup:
		// Waiting on input popup; keys already handled when showInput.
		return m, tea.Batch(cmds...)
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
			m.sidebarModel.SetFocused(true)
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
			m.sidebarModel.SetFocused(m.focus == FocusSidebar)
			return m, nil

		case key.Matches(keyMsg, DefaultKeyMap.FilterCategory):
			m.categoryPopup = popup.NewCategoryPicker(
				module.DefaultRegistry.Categories(),
				m.sidebarModel.CategoryFilter(),
			)
			m.showCategory = true
			return m, nil

		case key.Matches(keyMsg, DefaultKeyMap.Docs):
			m.docsPopup = popup.NewScriptViewer("Adding modules", docs.AddingModules)
			m.showDocs = true
			return m, nil

		// Install the selected module
		case key.Matches(keyMsg, DefaultKeyMap.Install):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil {
					items := []string{mod.Name + " — " + mod.Description}
					queue, err := planInstallQueue(m.selectedMod)
					if err != nil {
						items = append(items, "  ✗ "+err.Error())
					} else {
						for _, name := range queue {
							if name == m.selectedMod {
								continue
							}
							depMod := module.DefaultRegistry.Get(name)
							if depMod != nil {
								items = append(items, "  ▸ dep: "+depMod.Name+" — "+depMod.Description)
							} else {
								items = append(items, "  ▸ dep: "+name)
							}
						}
						if len(queue) == 0 {
							items = append(items, "  (already installed — press r to re-run)")
						}
					}
					hasScript := utils.ModuleScriptExists(m.selectedMod, "install.sh")
					if hasScript {
						items = append(items, "  ▸ Run install.sh")
					}
					if mod.StowEnabled {
						items = append(items, "  ▸ Create stow symlinks")
					}
					m.confirmPopup = popup.NewConfirmDialog(m.selectedMod, items, hasScript)
					m.showConfirm = true
					return m, nil
				}
			}

		// Force re-run install.sh for the selected module (even if already installed).
		case key.Matches(keyMsg, DefaultKeyMap.Reinstall):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil {
					hasScript := utils.ModuleScriptExists(m.selectedMod, "install.sh")
					if !hasScript && !mod.StowEnabled {
						return m, nil
					}
					m.confirmPopup = popup.NewReinstallDialog(m.selectedMod, hasScript)
					m.showConfirm = true
					return m, nil
				}
			}

		// Uninstall the selected module — show a confirmation dialog first.
		case key.Matches(keyMsg, DefaultKeyMap.Uninstall):
			if m.selectedMod != "" {
				mod := module.DefaultRegistry.Get(m.selectedMod)
				if mod != nil {
					var items []string
					items = append(items, mod.Name+" — "+mod.Description)
					hasScript := utils.ModuleScriptExists(m.selectedMod, "uninstall.sh")
					if hasScript {
						items = append(items, "  ▸ Run uninstall.sh")
					}
					if mod.StowEnabled {
						items = append(items, "  ▸ Remove stow symlinks")
					}

					m.confirmPopup = popup.NewUninstallDialog(m.selectedMod, items, hasScript)
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

// updateScriptPopup handles messages for the script review overlay.
func (m Model) updateScriptPopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.scriptPopup, cmd = m.scriptPopup.Update(msg)
	return m, cmd
}

// updateCategoryPopup handles messages for the category filter picker.
func (m Model) updateCategoryPopup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.categoryPopup, cmd = m.categoryPopup.Update(msg)
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
	m.sidebarModel.SetFocused(m.focus == FocusSidebar)

	contentWidth, contentHeight := m.contentDimensions()

	// Sidebar gets ~30% of width.
	sidebarWidth := int(float64(contentWidth) * 0.3)
	detailWidth := contentWidth - sidebarWidth - 1

	// Match viewDashboard: both panels always get a 1-cell border frame, so
	// inner content is (panelW-2) x (panelH-2) and total size never shifts.
	panelHeight := contentHeight - 5
	if panelHeight < 10 {
		panelHeight = 10
	}
	m.sidebarModel.SetSize(sidebarWidth-2, panelHeight-2)
	m.detailModel.SetSize(detailWidth-2, panelHeight-2)

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

	// Module dependencies (from module.yaml) plus external package deps.
	deps := make([]detail.DepStatus, 0, len(mod.Dependencies)+len(mod.ExternalDeps))
	for _, depName := range mod.Dependencies {
		depMod := module.DefaultRegistry.Get(depName)
		check := ""
		if depMod != nil {
			check = depMod.CheckCommand
		}
		deps = append(deps, detail.DepStatus{
			Name:      depName,
			Method:    "module",
			Installed: utils.ModuleSatisfied(depName, check),
			Checking:  false,
		})
	}
	for _, dep := range mod.ExternalDeps {
		installed := false
		if dep.CheckCommand != "" {
			installed = utils.ModuleSatisfied(dep.Name, dep.CheckCommand)
		} else {
			installed = utils.IsCommandAvailable(dep.Name)
		}
		method := dep.InstallMethod
		if method == "" {
			method = "external"
		}
		deps = append(deps, detail.DepStatus{
			Name:      dep.Name,
			Method:    method,
			Installed: installed,
			Checking:  false,
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

func (m *Model) reloadModules() {
	dir := utils.GetModulesDir()
	res := module.LoadFromDir(dir)
	m.loadWarnings = res.Errors
	m.sidebarModel.SetItems(buildSidebarItems())
}

func expandHome(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		if path == "~" {
			return home
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// planInstallQueue returns modules to install for target (deps first), skipping
// anything already satisfied via tracking file or check_command.
func planInstallQueue(target string) ([]string, error) {
	order, err := module.DefaultRegistry.GetInstallOrder([]string{target})
	if err != nil {
		return nil, err
	}
	var todo []string
	for _, name := range order {
		mod := module.DefaultRegistry.Get(name)
		if mod == nil {
			return nil, fmt.Errorf("unknown module %q", name)
		}
		if utils.ModuleSatisfied(name, mod.CheckCommand) {
			continue
		}
		todo = append(todo, name)
	}
	return todo, nil
}

// beginModuleInstall starts install.sh/stow for one module (possibly mid-queue).
func (m Model) beginModuleInstall(name string) (tea.Model, tea.Cmd) {
	mod := module.DefaultRegistry.Get(name)
	if mod == nil {
		return m, func() tea.Msg {
			return installer.InstallCompleteMsg{
				ModuleName: name,
				Success:    false,
				Error:      fmt.Errorf("unknown module"),
			}
		}
	}

	scriptPath := ""
	if utils.ModuleScriptExists(name, "install.sh") {
		scriptPath = utils.ModuleScriptPath(name, "install.sh")
	}
	modulesDir := utils.GetModulesDir()
	m.installingMod = name
	m.installScriptPath = scriptPath
	m.installStowEnabled = mod.StowEnabled

	step, total := installProgress(m.installPlan, m.installQueue, name)
	m.detailModel.OutputModel().SetInstalling(name, true)
	m.detailModel.OutputModel().AppendLine("")
	m.detailModel.OutputModel().AppendLine(
		fmt.Sprintf("━━━ [%d/%d] Installing %s ━━━", step, total, name))
	if len(m.installQueue) > 0 {
		m.detailModel.OutputModel().AppendLine(
			fmt.Sprintf("    queued next: %s", strings.Join(m.installQueue, ", ")))
	} else if total > 1 {
		m.detailModel.OutputModel().AppendLine("    (last module in plan)")
	}

	// Stow-only: finish via InstallCompleteMsg so the dep queue advances.
	if scriptPath == "" {
		if mod.StowEnabled {
			if err := utils.Stow(name, modulesDir); err != nil {
				return m, func() tea.Msg {
					return installer.InstallCompleteMsg{
						ModuleName: name,
						Success:    false,
						Error:      err,
					}
				}
			}
			m.detailModel.OutputModel().AppendLine("✓ Stow links created")
		}
		_ = utils.SetModuleInstalled(name)
		return m, func() tea.Msg {
			return installer.InstallCompleteMsg{ModuleName: name, Success: true}
		}
	}

	if installer.NeedsSudoScript(scriptPath) {
		m.pendingAction = popup.ActionInstall
		m.detailModel.OutputModel().AppendLine("▸ Authenticating sudo…")
		return m, installer.RunSudoAuth(name)
	}

	if m.program == nil {
		return m, func() tea.Msg {
			return installer.InstallCompleteMsg{
				ModuleName: name,
				Success:    false,
				Error:      fmt.Errorf("internal: tea program not ready"),
			}
		}
	}

	return m, installer.RunInstallStreaming(
		m.program, name, scriptPath, modulesDir, mod.StowEnabled)
}

// installProgress returns 1-based step and total for the current module.
func installProgress(plan, remaining []string, current string) (step, total int) {
	total = len(plan)
	if total == 0 {
		return 1, 1
	}
	// remaining is modules after current; completed = total - len(remaining) - 1
	step = total - len(remaining)
	if step < 1 {
		step = 1
	}
	if step > total {
		step = total
	}
	_ = current
	return step, total
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
	case StateModulesSetup:
		content = theme.Title.Render("Modules path setup") + "\n\n" +
			theme.NormalText.Render("Point the TUI at a folder that contains your modules") + "\n" +
			theme.DimText.Render("(each subfolder has module.yaml + optional install.sh).") + "\n\n" +
			theme.DimText.Render("Enter a path in the dialog (use ~ for your home directory).")
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
	if m.showScript {
		finalView = m.scriptPopup.Render(contentWidth+2, contentHeight+2)
	}
	if m.showDocs {
		finalView = m.docsPopup.Render(contentWidth+2, contentHeight+2)
	}
	if m.showCategory {
		finalView = m.categoryPopup.Render(contentWidth+2, contentHeight+2)
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
	focusLabel := "modules"
	if m.focus == FocusDetail {
		focusLabel = "detail"
	}
	title := theme.Title.Render("⚡ DotFiles Manager") +
		theme.DimText.Render("  · focus: ") +
		theme.Subtitle.Render(focusLabel)

	help := theme.HelpStyle.Render(
		"q: quit • ?: help • H: docs • j/k: navigate • shift+tab: switch panel • tab: switch tab • i: install • d: uninstall • o: open URL • s: search • c: category",
	)

	panelHeight := height - 5
	if panelHeight < 10 {
		panelHeight = 10
	}

	sidebarWidth := int(float64(width) * 0.3)
	detailWidth := width - sidebarWidth - 1

	sidebarContent := m.sidebarModel.View()
	detailContent := m.detailModel.View()
	if m.focus != FocusSidebar {
		sidebarContent = lipgloss.NewStyle().Faint(true).Render(sidebarContent)
	}
	if m.focus != FocusDetail {
		detailContent = lipgloss.NewStyle().Faint(true).Render(detailContent)
	}

	// Always frame both panels (same geometry). Focus = cyan border; other = muted.
	sidebarView := framePanel(sidebarContent, m.focus == FocusSidebar, sidebarWidth, panelHeight)
	detailView := framePanel(detailContent, m.focus == FocusDetail, detailWidth, panelHeight)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, " ", detailView)
	body = theme.Clip(body, width, panelHeight)

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, body, help)
}

// framePanel draws a fixed-size panel border. Width/Height are the TOTAL outer
// size. Content is clipped to the inner area before the border is applied so
// overflowing module cards cannot push the bottom border off-screen.
func framePanel(content string, focused bool, width, height int) string {
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	borderColor := theme.ColorSurface
	if focused {
		borderColor = theme.ColorCyan
	}

	clipped := theme.Clip(content, innerW, innerH)
	framed := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(innerW).
		Render(clipped)

	// Final clamp — keeps the bottom border visible even if lipgloss border
	// math differs slightly across terminals.
	return theme.Clip(framed, width, height)
}
