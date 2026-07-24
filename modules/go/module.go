package golang

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "go",
		Icon:        "",
		Description: "Go programming language",
		Category:    "Language",
		Website:     "https://go.dev/",
		Repo:        "https://github.com/golang/go",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "500MB",
		CheckCommand:  "go version",
	})
}
