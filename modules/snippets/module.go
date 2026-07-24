package snippets

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:          "snippets",
		Icon:          "",
		Description:   "Custom code snippet definitions",
		Category:      "Utility",
		Website:       "https://github.com/Issafalcon/dotFiles",
		Repo:          "https://github.com/Issafalcon/dotFiles",
		StowEnabled:   true,
		EstimatedTime: "5s",
		EstimatedSize: "10KB",
		CheckCommand:  "test -d $HOME/.config/snippets",
	})
}
