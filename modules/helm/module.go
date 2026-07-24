package helm

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "helm",
		Icon:        "󱃾",
		Description: "Kubernetes package manager",
		Category:    "DevOps",
		Website:     "https://helm.sh/",
		Repo:        "https://github.com/helm/helm",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "50MB",
		CheckCommand:  "helm version",
	})
}
