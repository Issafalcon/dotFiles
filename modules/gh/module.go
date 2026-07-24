package gh

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "gh",
		Icon:        "",
		Description: "GitHub CLI",
		Category:    "Utility",
		Website:     "https://cli.github.com/",
		Repo:        "https://github.com/cli/cli",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "gh --version",
	})
}
