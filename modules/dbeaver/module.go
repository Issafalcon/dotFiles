package dbeaver

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "dbeaver",
		Icon:        "",
		Description: "DBeaver universal database manager",
		Category:    "Database",
		Website:     "https://dbeaver.io/",
		Repo:        "https://github.com/dbeaver/dbeaver",
		StowEnabled:   false,
		EstimatedTime: "2m",
		EstimatedSize: "300MB",
		CheckCommand:  "dbeaver --version",
	})
}
