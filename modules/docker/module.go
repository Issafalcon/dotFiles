package docker

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "docker",
		Icon:        "",
		Description: "Docker container runtime and tools",
		Category:    "DevOps",
		Website:     "https://www.docker.com/",
		Repo:        "https://github.com/docker/cli",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "docker --version",
	})
}
