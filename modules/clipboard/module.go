package clipboard

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "clipboard",
		Icon:        "󰅍",
		Description: "Clipboard integration (xclip)",
		Category:    "Utility",
		Website:     "https://github.com/astrand/xclip",
		Repo:        "https://github.com/astrand/xclip",
		StowEnabled:   false,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "xclip -version",
	})
}
