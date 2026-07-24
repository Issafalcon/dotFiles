package terraform

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "terraform",
		Icon:         "󱁢",
		Description:  "Terraform and Terragrunt for IaC",
		Category:     "DevOps",
		Website:      "https://www.terraform.io/",
		Repo:         "https://github.com/hashicorp/terraform",
		Dependencies: []string{"homebrew"},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "200MB",
		CheckCommand:  "terraform --version",
	})
}
