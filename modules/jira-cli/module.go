package jiracli

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "jira-cli",
		Icon:        "󰌃",
		Description: "Interactive Jira CLI tool",
		Category:    "Application",
		Website:     "https://github.com/ankitpokhrel/jira-cli",
		Repo:        "https://github.com/ankitpokhrel/jira-cli",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "30MB",
		CheckCommand:  "test -d /usr/local/jira_1.1.0_linux_x86_64",
	})
}
