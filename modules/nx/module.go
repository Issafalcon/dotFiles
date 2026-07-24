package nx

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "nx",
		Icon:        "󰝖",
		Description: "Nx monorepo build system CLI",
		Category:    "Application",
		Website:     "https://nx.dev/",
		Repo:        "https://github.com/nrwl/nx",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "nx --version",
	})
}
