package zathura

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "zathura",
		Icon:        "",
		Description: "Document viewer with vim-like keybindings",
		Category:    "Utility",
		Website:     "https://pwmt.org/projects/zathura/",
		Repo:        "https://git.pwmt.org/pwmt/zathura",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "zathura --version",
	})
}
