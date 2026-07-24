package pandoc

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "pandoc",
		Icon:        "",
		Description: "Universal document converter",
		Category:    "Utility",
		Website:     "https://pandoc.org/",
		Repo:        "https://github.com/jgm/pandoc",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "pandoc --version",
	})
}
