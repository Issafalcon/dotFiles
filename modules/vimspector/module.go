package vimspector

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:          "vimspector",
		Icon:          "",
		Description:   "Multi-language debugging for Vim/Neovim",
		Category:      "Editor",
		Website:       "https://puremourning.github.io/vimspector-web/",
		Repo:          "https://github.com/puremourning/vimspector",
		StowEnabled:   true,
		EstimatedTime: "10s",
		EstimatedSize: "1MB",
		CheckCommand:  "test -d $HOME/.config/vimspector",
	})
}
