package zsh

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "zsh",
		Icon:        "",
		Description: "Z Shell with zinit plugin manager",
		Category:    "Shell",
		Website:     "https://www.zsh.org/",
		Repo:        "https://github.com/zsh-users/zsh",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "50MB",
		CheckCommand:  "zsh --version",
	})
}
