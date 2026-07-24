package aws

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "aws",
		Icon:        "",
		Description: "AWS CLI v2",
		Category:    "Cloud",
		Website:     "https://aws.amazon.com/cli/",
		Repo:        "https://github.com/aws/aws-cli",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "aws --version",
	})
}
