package azure

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "azure",
		Icon:        "󰠅",
		Description: "Azure CLI with Functions Core Tools",
		Category:    "Cloud",
		Website:     "https://docs.microsoft.com/en-us/cli/azure/",
		Repo:        "https://github.com/Azure/azure-cli",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "az --version",
	})
}
