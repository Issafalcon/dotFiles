package editorconfig

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:          "editorconfig",
		Icon:          "",
		Description:   "Editor configuration for consistent coding styles",
		Category:      "Utility",
		Website:       "https://editorconfig.org/",
		Repo:          "https://github.com/editorconfig/editorconfig",
		StowEnabled:   true,
		EstimatedTime: "5s",
		EstimatedSize: "1KB",
		CheckCommand:  "test -f $HOME/.editorconfig",
	})
}
