package git

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "git",
		Icon:         "",
		Description:  "Git with delta diff viewer",
		Category:     "Utility",
		Website:      "https://git-scm.com/",
		Repo:         "https://github.com/git/git",
		Dependencies: []string{"homebrew"},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "git --version",
	})
}
