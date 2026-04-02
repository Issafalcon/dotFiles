package modules

// ai.go registers AI-related modules: ai-tools.
//
// This module installs AI coding assistants — Claude Code and GitHub Copilot
// CLI — along with the MCP Hub for managing Model Context Protocol servers.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- ai-tools ---
	// AI coding tools: Claude Code CLI, GitHub Copilot CLI, and MCP Hub.
	// Requires node (for npm packages) as a dependency.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "ai-tools",
		Icon:         "󰚩",
		Description:  "AI coding tools (Claude, Copilot, MCP Hub)",
		Category:     "AI",
		Website:      "https://claude.ai/",
		Repo:         "https://github.com/anthropics/claude-code",
		Dependencies: []string{"node"},
		ExternalDeps: []module.ExternalDep{
			{
				Name:           "npm",
				CheckCommand:   "npm --version",
				InstallCommand: "Install the 'node' module first",
				InstallMethod:  "npm",
			},
		},
		InstallCommands: []string{
			// Claude Code CLI
			`if ! command -v claude >/dev/null; then curl -fsSL https://claude.ai/install.sh | bash; fi`,
			// GitHub Copilot CLI
			`if ! command -v copilot >/dev/null; then npm install -g @github/copilot; fi`,
			// MCP Hub
			"npm install -g mcp-hub@latest",
		},
		UninstallCommands: []string{
			"npm uninstall -g @github/copilot",
			"npm uninstall -g mcp-hub",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "claude --version",
	})
}
