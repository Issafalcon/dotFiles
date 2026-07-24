package gojira

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "gojira",
		Icon:        "󰌃",
		Description: "Command-line Jira client (go-jira)",
		Category:    "Application",
		Website:     "https://github.com/go-jira/jira",
		Repo:        "https://github.com/go-jira/jira",
		StowEnabled:   true,
		EstimatedTime: "15s",
		EstimatedSize: "20MB",
		CheckCommand:  "jira --version",
	})
}
