package packer

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "packer",
		Icon:        "󱁢",
		Description: "HashiCorp Packer for machine images",
		Category:    "DevOps",
		Website:     "https://www.packer.io/",
		Repo:        "https://github.com/hashicorp/packer",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "packer --version",
	})
}
