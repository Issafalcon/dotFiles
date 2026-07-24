package fzf

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "fzf",
		Icon:        "",
		Description: "Command-line fuzzy finder",
		Category:    "Utility",
		Website:     "https://junegunn.github.io/fzf/",
		Repo:        "https://github.com/junegunn/fzf",
		StowEnabled:   true,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "fzf --version",
	})
}
