package node

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "node",
		Icon:        "󰎙",
		Description: "Node.js via NVM (Node Version Manager)",
		Category:    "Language",
		Website:     "https://nodejs.org/",
		Repo:        "https://github.com/nvm-sh/nvm",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "200MB",
		CheckCommand:  "node --version",
	})
}
