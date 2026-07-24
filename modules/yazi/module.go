package yazi

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "yazi",
		Icon:         "󰇥",
		Description:  "Blazing fast terminal file manager",
		Category:     "Utility",
		Website:      "https://yazi-rs.github.io/",
		Repo:         "https://github.com/sxyazi/yazi",
		Dependencies: []string{"homebrew"},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "100MB",
		CheckCommand:  "yazi --version",
	})
}
