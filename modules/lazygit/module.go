package lazygit

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "lazygit",
		Icon:         "",
		Description:  "Simple terminal UI for Git",
		Category:     "Utility",
		Website:      "https://github.com/jesseduffield/lazygit",
		Repo:         "https://github.com/jesseduffield/lazygit",
		Dependencies: []string{"homebrew"},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "30MB",
		CheckCommand:  "lazygit --version",
	})
}
