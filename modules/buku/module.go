package buku

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "buku",
		Icon:        "",
		Description: "Command-line bookmark manager",
		Category:    "Utility",
		Website:     "https://github.com/jarun/buku",
		Repo:        "https://github.com/jarun/buku",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "10MB",
		CheckCommand:  "buku --version",
	})
}
