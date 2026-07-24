package vagrant

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "vagrant",
		Icon:        "󰜐",
		Description: "HashiCorp Vagrant for dev environments",
		Category:    "DevOps",
		Website:     "https://www.vagrantup.com/",
		Repo:        "https://github.com/hashicorp/vagrant",
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "vagrant --version",
	})
}
