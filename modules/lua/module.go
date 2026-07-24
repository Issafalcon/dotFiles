package lua

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "lua",
		Icon:        "",
		Description: "Lua 5.4 scripting language",
		Category:    "Language",
		Website:     "https://www.lua.org/",
		Repo:        "https://github.com/lua/lua",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "10MB",
		CheckCommand:  "lua5.4 -v",
	})
}
