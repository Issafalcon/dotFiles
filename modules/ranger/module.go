package ranger

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "ranger",
		Icon:        "",
		Description: "Console file manager with VI bindings",
		Category:    "Utility",
		Website:     "https://ranger.github.io/",
		Repo:        "https://github.com/ranger/ranger",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "ranger --version",
	})
}
