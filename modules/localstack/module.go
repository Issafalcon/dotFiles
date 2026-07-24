package localstack

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "localstack",
		Icon:         "󰸏",
		Description:  "Local AWS cloud stack for testing",
		Category:     "DevOps",
		Website:      "https://localstack.cloud/",
		Repo:         "https://github.com/localstack/localstack",
		Dependencies: []string{"python"},
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "localstack --version",
	})
}
