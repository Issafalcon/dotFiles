package plantuml

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "plantuml",
		Icon:        "󰈏",
		Description: "UML diagrams from text descriptions",
		Category:    "Utility",
		Website:     "https://plantuml.com/",
		Repo:        "https://github.com/plantuml/plantuml",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "plantuml -version",
	})
}
