package obsidian

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "obsidian",
		Icon:        "󰈙",
		Description: "Knowledge base on local Markdown files",
		Category:    "Application",
		Website:     "https://obsidian.md/",
		Repo:        "https://github.com/obsidianmd/obsidian-releases",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "test -f /usr/bin/obsidian",
	})
}
