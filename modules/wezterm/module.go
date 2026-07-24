package wezterm

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "wezterm",
		Icon:        "",
		Description: "GPU-accelerated terminal emulator",
		Category:    "Utility",
		Website:     "https://wezfurlong.org/wezterm/",
		Repo:        "https://github.com/wez/wezterm",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "wezterm --version",
	})
}
