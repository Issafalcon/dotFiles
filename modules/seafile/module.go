package seafile

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "seafile",
		Icon:        "󰅟",
		Description: "Self-hosted file sync and share",
		Category:    "Application",
		Website:     "https://www.seafile.com/",
		Repo:        "https://github.com/haiwen/seafile",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "100MB",
		CheckCommand:  "test -f /usr/bin/seafile",
	})
}
