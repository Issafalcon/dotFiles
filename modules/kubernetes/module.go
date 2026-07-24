package kubernetes

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "kubernetes",
		Icon:        "󱃾",
		Description: "kubectl CLI and k9s terminal UI",
		Category:    "DevOps",
		Website:     "https://kubernetes.io/",
		Repo:        "https://github.com/kubernetes/kubernetes",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "100MB",
		CheckCommand:  "kubectl version --client",
	})
}
