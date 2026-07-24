package tmux

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "tmux",
		Icon:        "",
		Description: "Terminal multiplexer",
		Category:    "Utility",
		Website:     "https://github.com/tmux/tmux/wiki",
		Repo:        "https://github.com/tmux/tmux",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "tmux -V",
	})
}
