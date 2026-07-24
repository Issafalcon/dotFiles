package quarto

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "quarto",
		Icon:        "󰐗",
		Description: "Scientific and technical publishing",
		Category:    "Utility",
		Website:     "https://quarto.org/",
		Repo:        "https://github.com/quarto-dev/quarto-cli",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "quarto --version",
	})
}
