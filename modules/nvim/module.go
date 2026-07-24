package nvim

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "nvim",
		Icon:         "",
		Description:  "Neovim editor with full plugin ecosystem",
		Category:     "Editor",
		Website:      "https://neovim.io/",
		Repo:         "https://github.com/neovim/neovim",
		Dependencies: []string{"python", "node", "go", "homebrew", "yazi"},
		ExternalDeps: []module.ExternalDep{
			{
				Name:           "ripgrep",
				CheckCommand:   "rg --version",
				InstallCommand: "sudo apt-get install -y ripgrep",
				InstallMethod:  "apt",
			},
			{
				Name:           "cmake",
				CheckCommand:   "cmake --version",
				InstallCommand: "sudo apt-get install -y cmake",
				InstallMethod:  "apt",
			},
		},
		StowEnabled:   true,
		EstimatedTime: "5m",
		EstimatedSize: "500MB",
		CheckCommand:  "nvim --version",
	})
}
