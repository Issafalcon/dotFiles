package wslopen

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "wsl-open",
		Icon:        "󰖳",
		Description: "Open files in Windows apps from WSL",
		Category:    "Utility",
		Website:     "https://github.com/4U6U57/wsl-open",
		Repo:        "https://github.com/4U6U57/wsl-open",
		StowEnabled:   false,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "wsl-open --version",
	})
}
