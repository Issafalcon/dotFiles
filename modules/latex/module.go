package latex

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "latex",
		Icon:        "",
		Description:  "Full TeX Live distribution",
		Category:     "Utility",
		Website:      "https://www.latex-project.org/",
		Repo:         "https://github.com/latex3/latex3",
		Dependencies: []string{"python"},
		StowEnabled:   true,
		EstimatedTime: "10m",
		EstimatedSize: "5GB",
		CheckCommand:  "latex --version",
	})
}
