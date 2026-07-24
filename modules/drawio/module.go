package drawio

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "drawio",
		Icon:        "󰺷",
		Description: "Diagramming application (draw.io)",
		Category:    "Application",
		Website:     "https://www.drawio.com/",
		Repo:        "https://github.com/jgraph/drawio-desktop",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "150MB",
		CheckCommand:  "test -f /usr/bin/drawio",
	})
}
