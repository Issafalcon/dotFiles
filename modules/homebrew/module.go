package homebrew

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "homebrew",
		Icon:        "🍺",
		Description: "Homebrew package manager for Linux",
		Category:    "Utility",
		Website:     "https://brew.sh/",
		Repo:        "https://github.com/Homebrew/brew",
		StowEnabled:   false,
		EstimatedTime: "3m",
		EstimatedSize: "500MB",
		CheckCommand:  "brew --version",
		RequiresInput: true,
	})
}
