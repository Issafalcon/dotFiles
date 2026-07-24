package aitools

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
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
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "claude --version",
	})
}
